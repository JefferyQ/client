// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"errors"
	"fmt"
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"github.com/keybase/go-framed-msgpack-rpc/rpc"
)

type PGPPushPrivate struct {
	arg keybase1.PGPPushPrivateArg
}

func (e *PGPPushPrivate) Name() string {
	return "PGPPushPrivate"
}

func (e *PGPPushPrivate) Prereqs() Prereqs {
	return Prereqs{}
}

func (e *PGPPushPrivate) RequiredUIs() []libkb.UIKind {
	return []libkb.UIKind{}
}

func (e *PGPPushPrivate) SubConsumers() []libkb.UIConsumer {
	return []libkb.UIConsumer{}
}

func NewPGPPushPrivate(arg keybase1.PGPPushPrivateArg) *PGPPushPrivate {
	return &PGPPushPrivate{arg}
}

func getCurrentUserPGPKeys(m libkb.MetaContext) ([]libkb.PGPFingerprint, error) {
	uid := m.CurrentUID()
	if uid.IsNil() {
		return nil, libkb.NewLoginRequiredError("for push/pull of PGP private keys to KBFS")
	}
	upk, _, err := m.G().GetUPAKLoader().LoadV2(libkb.NewLoadUserArgWithMetaContext(m).WithUID(uid))
	if err != nil {
		return nil, err
	}
	var res []libkb.PGPFingerprint
	for _, key := range upk.Current.PGPKeys {
		res = append(res, libkb.PGPFingerprint(key.Fingerprint))
	}
	return res, nil
}

func getPrivateFingerprints(m libkb.MetaContext, fphex string) ([]libkb.PGPFingerprint, error) {
	if len(fphex) == 0 {
		return getCurrentUserPGPKeys(m)
	}
	fp, err := libkb.PGPFingerprintFromHex(fphex)
	if err != nil {
		return nil, err
	}
	return []libkb.PGPFingerprint{*fp}, nil
}

func simpleFSClient(m libkb.MetaContext) (*keybase1.SimpleFSClient, error) {
	xp := m.G().ConnectionManager.LookupByClientType(keybase1.ClientType_KBFS)
	if xp == nil {
		return nil, libkb.KBFSNotRunningError{}
	}
	return &keybase1.SimpleFSClient{
		Cli: rpc.NewClient(xp, libkb.NewContextifiedErrorUnwrapper(m.G()), nil),
	}, nil
}

func (e *PGPPushPrivate) mkdir(m libkb.MetaContext, fs *keybase1.SimpleFSClient, path string) (err error) {
	opid, err := fs.SimpleFSMakeOpid(m.Ctx())
	if err != nil {
		return err
	}
	defer fs.SimpleFSClose(m.Ctx(), opid)
	err = fs.SimpleFSOpen(m.Ctx(), keybase1.SimpleFSOpenArg{
		OpID:  opid,
		Dest:  keybase1.NewPathWithKbfs(path),
		Flags: keybase1.OpenFlags_DIRECTORY,
	})
	return err
}

func (e *PGPPushPrivate) write(m libkb.MetaContext, fs *keybase1.SimpleFSClient, path string, data string) (err error) {
	opid, err := fs.SimpleFSMakeOpid(m.Ctx())
	if err != nil {
		return err
	}
	defer fs.SimpleFSClose(m.Ctx(), opid)
	err = fs.SimpleFSOpen(m.Ctx(), keybase1.SimpleFSOpenArg{
		OpID:  opid,
		Dest:  keybase1.NewPathWithKbfs(path),
		Flags: keybase1.OpenFlags_EXISTING,
	})
	if err != nil {
		return err
	}
	err = fs.SimpleFSWrite(m.Ctx(), keybase1.SimpleFSWriteArg{
		OpID:    opid,
		Offset:  0,
		Content: []byte(data),
	})
	if err != nil {
		return err
	}
	return nil
}

func (e *PGPPushPrivate) link(m libkb.MetaContext, fs *keybase1.SimpleFSClient, file string, link string) (err error) {
	err = fs.SimpleFSSymlink(m.Ctx(), keybase1.SimpleFSSymlinkArg{
		Target: file,
		Link:   keybase1.NewPathWithKbfs(link),
	})
	return err
}

func (e *PGPPushPrivate) push(m libkb.MetaContext, fp libkb.PGPFingerprint, tty string, fs *keybase1.SimpleFSClient) error {
	armored, err := m.G().GetGpgClient().ImportKeyArmored(true, fp, tty)
	if err != nil {
		return err
	}

	username := m.CurrentUsername()
	if username.IsNil() {
		return libkb.NewLoginRequiredError("no username found")
	}

	path := "/kebase/private/" + username.String() + ".keys"
	e.mkdir(m, fs, path)
	path = path + "/gpg"
	err = e.mkdir(m, fs, path)
	if err != nil {
		return err
	}

	path = path + "/" + fp.String()
	link := path + ".asc"
	file := path + "-" + fmt.Sprintf("%d", m.G().Clock().Now().Unix()) + ".asc"

	err = e.write(m, fs, file, armored)
	if err != nil {
		return err
	}

	err = e.link(m, fs, file, link)
	return err
}

func (e *PGPPushPrivate) Run(m libkb.MetaContext) (err error) {

	defer m.CTrace("PGPPushPrivate#Run", func() error { return err })()

	tty, err := m.UIs().GPGUI.GetTTY(m.Ctx())
	if err != nil {
		return err
	}

	fingerprints, err := getPrivateFingerprints(m, e.arg.Fingerprint)
	if err != nil {
		return err
	}

	fs, err := simpleFSClient(m)
	if err != nil {
		return err
	}

	if len(fingerprints) == 0 {
		return errors.New("no PGP keys provided")
	}

	for _, fp := range fingerprints {
		err = e.push(m, fp, tty, fs)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
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

func (e *PGPPushPrivate) Run(m libkb.MetaContext) error {
	return nil
}

func NewPGPPushPrivate(arg keybase1.PGPPushPrivateArg) *PGPPushPrivate {
	return &PGPPushPrivate{arg}
}

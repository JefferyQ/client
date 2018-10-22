// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol/keybase1"
)

type PGPPullPrivate struct {
	arg keybase1.PGPPullPrivateArg
}

func (e *PGPPullPrivate) Name() string {
	return "PGPPullPrivate"
}

func (e *PGPPullPrivate) Prereqs() Prereqs {
	return Prereqs{}
}

func (e *PGPPullPrivate) RequiredUIs() []libkb.UIKind {
	return []libkb.UIKind{}
}

func (e *PGPPullPrivate) SubConsumers() []libkb.UIConsumer {
	return []libkb.UIConsumer{}
}

func (e *PGPPullPrivate) Run(m libkb.MetaContext) error {
	return nil
}

func NewPGPPullPrivate(arg keybase1.PGPPullPrivateArg) *PGPPullPrivate {
	return &PGPPullPrivate{arg}
}

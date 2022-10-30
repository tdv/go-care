// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// The license can be found in the LICENSE file.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"sync/atomic"
)

type Switch interface {
	TurnOn()
	TurnOff()

	IsTurnedOn() bool
}

type switchImpl struct {
	isTurnedOn atomic.Bool
}

func (this *switchImpl) TurnOn() {
	this.isTurnedOn.Store(true)
}

func (this *switchImpl) TurnOff() {
	this.isTurnedOn.Store(false)
}

func (this *switchImpl) IsTurnedOn() bool {
	return this.isTurnedOn.Load()
}

func newSwitch() Switch {
	inst := switchImpl{}
	inst.TurnOn()
	return &inst
}

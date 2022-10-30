// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"sync/atomic"
)

// A feature switch.
type Switch interface {
	// Turns on the feature.
	TurnOn()
	// Turns off the feature.
	TurnOff()

	// Returns a state of the feature (turned on/off).
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

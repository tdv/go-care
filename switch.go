// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"sync/atomic"
)

// Switch is a switch feature.
type Switch interface {
	// TurnOn - turns on the feature.
	TurnOn()
	// TurnOff - turns off the feature.
	TurnOff()

	// IsTurnedOn - returns a state of the feature (turned on/off).
	IsTurnedOn() bool
}

type switchImpl struct {
	isTurnedOn atomic.Bool
}

func (s *switchImpl) TurnOn() {
	s.isTurnedOn.Store(true)
}

func (s *switchImpl) TurnOff() {
	s.isTurnedOn.Store(false)
}

func (s *switchImpl) IsTurnedOn() bool {
	return s.isTurnedOn.Load()
}

func newSwitch() Switch {
	inst := switchImpl{}
	inst.TurnOn()
	return &inst
}

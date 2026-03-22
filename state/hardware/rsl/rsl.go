package rsl

import (
	"sync/atomic"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
)

type (
	RSL struct {
		isEnabled atomic.Bool
		pin       *hardware.Pin
	}
)

func New(rslPin int) *RSL {
	rsl := &RSL{pin: hardware.NewPin(rslPin)}
	rsl.pin.OnClose(func() {
		rsl.pin.Set(false)
	})
	go rsl.loop()
	return rsl
}

func (rsl *RSL) loop() {
	t := time.NewTicker(time.Millisecond * 500)

	for range t.C {
		if rsl.isEnabled.Load() {
			rsl.pin.Set(!rsl.pin.Read())
		} else {
			rsl.pin.Set(hardware.PIN_HIGH)
		}
	}
}

func (rsl *RSL) SetEnabled(val bool) {
	rsl.isEnabled.Store(val)
}

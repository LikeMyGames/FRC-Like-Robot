package rsl

import (
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
)

type (
	RSL struct {
		isEnabled bool
		pin       *hardware.Pin
	}
)

func New(rslPin int) *RSL {
	rsl := &RSL{pin: hardware.NewPin(rslPin)}
	event.Listen("ROBOT_ENABLE_STATUS", "RSL", func(event any) {
		enabled, ok := event.(bool)
		if ok {
			rsl.isEnabled = enabled
		}
	})
	rsl.pin.OnClose(func() {
		rsl.pin.Set(false)
	})
	go rsl.loop()
	return rsl
}

func (rsl *RSL) loop() {
	t := time.NewTicker(time.Millisecond * 500)

	for range t.C {
		if rsl.isEnabled {
			rsl.pin.Set(!rsl.pin.Read())
		} else {
			rsl.pin.Set(hardware.PIN_HIGH)
		}
	}
}

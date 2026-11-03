package rsl

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
)

type (
	RSL struct {
		isEnabled atomic.Bool
		pin       *hardware.Pin
		cancel    context.CancelFunc
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
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	// ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	t := time.NewTicker(time.Millisecond * 500)

	for range t.C {
		select {
		case <-sigint:
			fmt.Println("stopping rsl")
			signal.Stop(sigint)
			t.Stop()
			return
		default:
			if rsl.isEnabled.Load() {
				rsl.pin.Set(!rsl.pin.Read())
			} else {
				rsl.pin.Set(hardware.PIN_HIGH)
			}
		}

	}
}

func (rsl *RSL) SetEnabled(val bool) {
	rsl.isEnabled.Store(val)
}

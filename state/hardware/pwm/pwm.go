package pwm

import (
	"fmt"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type (
	PWM struct {
		pin                *hardware.Pin
		Period             time.Duration
		dutyCycleTickCount int8
	}
)

// Creates a new 8-bit pwm object
func New(pinNum int, period time.Duration) *PWM {
	pwm := &PWM{pin: hardware.NewPin(pinNum), Period: period}
	return pwm
}

func (pwm *PWM) Start() {
	go func() {
		fmt.Println(pwm.Period)
		fmt.Println(pwm.Period / 256)
		t := time.NewTicker(pwm.Period / 256)
		tickCount := 0
		val := hardware.PIN_LOW

		for range t.C {
			if tickCount >= 256 {
				tickCount = 0
			}
			if val == hardware.PIN_LOW && tickCount < int(pwm.dutyCycleTickCount) {
				pwm.pin.Set(hardware.PIN_HIGH)
				val = hardware.PIN_HIGH
			} else if tickCount >= int(pwm.dutyCycleTickCount) {
				pwm.pin.Set(hardware.PIN_LOW)
				val = hardware.PIN_LOW
			}
			// fmt.Printf("%v: %d;%d\n", pwm.dutyCycleTickCount < int8(tickCount), tickCount, pwm.dutyCycleTickCount)

			tickCount++
		}
	}()
}

func (pwm *PWM) SetDutyCycleInPercent(percent float64) {
	pwm.dutyCycleTickCount = int8(mathutils.MapRange(percent, 0, 100, 0, 255))
	fmt.Println(pwm.dutyCycleTickCount)
}

// Num most be between 0 and 255
func (pwm *PWM) SetDutyCycle(duty int8) {
	pwm.dutyCycleTickCount = duty
}

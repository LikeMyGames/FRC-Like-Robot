//go:build !linux

package hardware

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
)

type (
	Pin struct {
	}
)

var (
	config constantTypes.Battery = constantTypes.Battery{}
)

func OpenSpi()

func SetBatteryConfig(conf constantTypes.Battery) {
	config = conf
}

func ReadBatteryPercentage() uint {
	return 0
}

func ReadBatteryVoltage() float64 {
	return -1
}

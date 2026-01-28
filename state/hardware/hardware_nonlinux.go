//go:build !linux

package hardware

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
)

type ()

var (
	config constantTypes.Battery = constantTypes.Battery{}
)

func SetConfig(cfg constantTypes.Battery) {
	config = cfg
}

func ReadBatteryPercentage() uint {
	return 0
}

func ReadBatteryVoltage() float64 {
	return -1
}

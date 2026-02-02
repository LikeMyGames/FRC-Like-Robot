//go:build !linux

package hardware

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils"
	"periph.io/x/conn/v3/spi"
)

type (
	Pin struct {
	}
)

var (
	config constantTypes.Battery = constantTypes.Battery{}
	Os     string                = "not linux"
)

const (
	PIN_HIGH bool = true
	PIN_LOW  bool = false
)

func OpenSpi() {}

func GetSpiConn() spi.Conn { return nil }

func CloseSpiPort() {}

func CheckStatus() bool { return false }

func NewPin(pinNum int) *Pin { return nil }

func (pin *Pin) Set(state bool) {}

func (pin *Pin) Read() bool { return false }

func (pin *Pin) OnClose(action func()) {}

func CloseAllPins() {}

func SetBatteryConfig(conf constantTypes.Battery) {
	config = conf
}

func ReadBatteryPercentage() float64 {
	return utils.TruncateFloat64(float64((ReadBatteryVoltage()/config.NominalVoltage)*100), 2)
}

func ReadBatteryVoltage() float64 {
	return 12
}

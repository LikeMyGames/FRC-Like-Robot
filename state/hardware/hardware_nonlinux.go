//go:build !linux

package hardware

import (
	"fmt"
	"math"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/pid"
)

type (
	CanBus struct {
		SPI_MISO_PIN int
		SPI_MOSI_PIN int
		SPI_SCLK_PIN int
	}

	CanNode struct {
		id          int
		registerMap map[string]int
	}

	Device struct {
		id     int64
		value  float64
		target float64
	}

	MotorController struct {
		device        Device
		PidController *pid.PIDController
	}
)

// var bus = NewCanBus(19, 21, 23)

var (
	config constantTypes.Battery = constantTypes.Battery{}
)

func NewCanBus(MOSI, MISO, SCLK int) *CanBus {
	return &CanBus{
		SPI_MISO_PIN: MISO,
		SPI_MOSI_PIN: MOSI,
		SPI_SCLK_PIN: SCLK,
	}
}

// func NewDevice(id int64) *Device {
// 	return &Device{id: id}
// }

// func (d *Device) SetTargetType(target string) {
// 	fmt.Println(d.id, "target type is", target)
// }

// func (d *Device) SetTargetValue(target float64) {
// 	fmt.Println(d.id, "target value is", target)
// }

func NewMotorController(config constantTypes.MotorController) *MotorController {
	return &MotorController{
		device:        Device{id: config.Id},
		PidController: pid.NewPIDController(config.PID),
	}
}

func (c *MotorController) SetTarget(val float64) {
	c.PidController.SetTarget(val)
}

func (c *MotorController) GetValue() float64 {
	return c.device.value
}

func (c *MotorController) GetTarget() float64 {
	return c.device.target
}

func (c *MotorController) GetId() int64 {
	return int64(c.device.target)
}

func (c *MotorController) Run() {
	c.PidController.Calculate(c.device.value)
}

func (c *MotorController) AtValue() bool {
	// replace less than val later
	return math.Abs(c.device.target-c.device.value) < 0.2
}

func (c *MotorController) Write(val float64) {
	fmt.Println("the write function for the MotorController struct only works in Linux, as such, it is only defined when the program is built in Linux")
}

func SetConfig(cfg constantTypes.Battery) {
	config = cfg
}

func ReadBatteryPercentage() uint {
	return 0
}

func ReadBatteryVoltage() float64 {
	return -1
}

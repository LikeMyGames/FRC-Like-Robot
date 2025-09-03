//go:build linux && cgo

package hardware

/*
#cgo CXXFLAGS: -std=c++11
#include "./interface.cpp"
#include "./interface.hpp"

// Declare the C function so cgo can link it
extern void Hello();
*/
import "C"

type (
	CANbus struct {
		CANname string
		BaseID  uint
		Devices []Device
	}

	Device interface {
		GetID() uint
		Read() []byte
		Write([]byte) bool
	}
	MotorController struct {
		Velocity     float64
		Displacement float64
		CandID       uint
	}

	Sensor struct {
		Value float64
	}
)

func NewBus() *CANbus {
	return &CANbus{
		CANname: "spidev0.0",
		BaseID:  0,
		Devices: []Device{},
	}
}

func (c *CANbus) AddDevice(d Device) {
}

func Hello() {
	C.Hello()
}

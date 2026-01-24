package single_motor_controller

import (
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
)

type (
	SingleMotorController struct {
		motor *motor.Motor
	}
)

var singleMotorRegisterMap map[string]int = map[string]int{
	"Status":       0x0,
	"Estop":        0x3,
	"Input_Pos":    0x0c,
	"Input_Vel":    0x0d,
	"Input_Torque": 0x0e,
	"Limits":       0x0f,
}

func NewSingleMotorController(CanId int) *SingleMotorController {
	controller := &SingleMotorController{}
	controller.motor.LoadRegisterMap(singleMotorRegisterMap)

	return controller
}

func (c *SingleMotorController) ReadAngle() float64 {
	return c.motor.ReadAngle()
}

func (c *SingleMotorController) SetAngle(angle float64) {
	c.motor.SetAngle(angle)
}

func (c *SingleMotorController) ReadVelocity() float64 {
	return c.motor.ReadAngle()
}

func (c *SingleMotorController) SetMotorVelocity(velocity float64) {
	c.motor.SetVelocity(velocity)
}

func (c *SingleMotorController) ReadTorque() float64 {
	return c.motor.ReadTorque()
}

func (c *SingleMotorController) SetTorque(torque float64) {
	c.motor.SetTorque(torque)
}

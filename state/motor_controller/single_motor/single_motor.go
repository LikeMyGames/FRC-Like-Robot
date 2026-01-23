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
	"": 0x1,
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

func (c *SingleMotorController) ReadAcceleration() float64 {
	return c.motor.ReadAcceleration()
}

func (c *SingleMotorController) SetAcceleration(acceleration float64) {
	c.motor.SetAcceleration(acceleration)
}

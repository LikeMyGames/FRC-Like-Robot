package double_motor_controller

import (
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
)

type (
	DoubleMotorController struct {
		primaryMotor   *motor.Motor
		secondaryMotor *motor.Motor
	}
)

func (c *DoubleMotorController) ReadPrimaryMotorAngle() float64 {
	return c.primaryMotor.ReadAngle()
}

func (c *DoubleMotorController) ReadSecondaryMotorAngle() float64 {
	return c.secondaryMotor.ReadAngle()
}

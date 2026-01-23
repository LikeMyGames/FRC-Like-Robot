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

var doubleMotorRegisterMap map[string]int = map[string]int{
	"": 0x1,
}

func NewDoubleMotorController() *DoubleMotorController {
	controller := &DoubleMotorController{}
	// controller.primaryMotor
	return controller
}

// func (c *DoubleMotorController) ReadPrimaryMotorAngle() float64 {
// 	return c.primaryMotor.ReadAngle()
// }

// func (c *DoubleMotorController) ReadSecondaryMotorAngle() float64 {
// 	return c.secondaryMotor.ReadAngle()
// }

// func (c *DoubleMotorController) SetPrimaryMotorAngle(angle float64) {
// 	c.primaryMotor.SetAngle(angle)
// }

// func (c *DoubleMotorController) SetSecondaryMotorAngle(angle float64) {
// 	c.secondaryMotor.SetAngle(angle)
// }

// func (c *DoubleMotorController) ReadPrimaryMotorVelocity() float64 {
// 	return c.primaryMotor.ReadAngle()
// }

// func (c *DoubleMotorController) ReadSecondaryMotorVelocity() float64 {
// 	return c.secondaryMotor.ReadAngle()
// }

// func (c *DoubleMotorController) SetPrimaryMotorVelocity(velocity float64) {
// 	c.primaryMotor.SetVelocity(velocity)
// }

// func (c *DoubleMotorController) SetSecondaryMotorVelocity(velocity float64) {
// 	c.secondaryMotor.SetVelocity(velocity)
// }

// func (c *DoubleMotorController) ReadPrimaryMotorAcceleration() float64 {
// 	return c.primaryMotor.ReadAcceleration()
// }

// func (c *DoubleMotorController) ReadSecondaryMotorAcceleration() float64 {
// 	return c.secondaryMotor.ReadAcceleration()
// }

// func (c *DoubleMotorController) SetPrimaryMotorAcceleration(acceleration float64) {
// 	c.primaryMotor.SetAcceleration(acceleration)
// }

// func (c *DoubleMotorController) SetSecondaryMotorAcceleration(acceleration float64) {
// 	c.secondaryMotor.SetAcceleration(acceleration)
// }

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
	"": 0,
}

func NewSingleMotorController(CanId int) *SingleMotorController {
	controller := &SingleMotorController{}
	controller.motor.LoadRegisterMap(singleMotorRegisterMap)

	return controller
}

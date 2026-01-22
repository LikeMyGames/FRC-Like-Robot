package single_motor_controller

import (
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
)

type (
	SingleMotorController struct {
		motor *motor.Motor
	}
)

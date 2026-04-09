package intake_types

import (
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type Constants struct {
	RollerMotorCanId    int
	ExtensionMotorCanId int

	RollerMotorP                      float64
	RollerMotorD                      float64
	RollerMotorFF                     float64
	RollerMotorPositionConversion     float64
	RollerMotorVelocityConversion     float64
	RollerMotorAccelerationConversion float64

	ExtensionMotorP                      float64
	ExtensionMotorI                      float64
	ExtensionMotorD                      float64
	ExtensionMotorCosFF                  float64
	ExtensionMotorPositionConversion     float64
	ExtensionMotorVelocityConversion     float64
	ExtensionMotorAccelerationConversion float64
}

type MotorConfigs struct {
	RollerMotorConfig    motor.Config
	ExtensionMotorConfig motor.Config
}

type IntakeSubsystem struct {
	RollerMotor    motor.Motor
	ExtensionMotor motor.Motor
	StateMachine   *state_machine.StateMachine
}

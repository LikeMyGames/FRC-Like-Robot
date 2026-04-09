package shooter_types

import (
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type (
	Constants struct {
		MaxFlyWheelVelocity     float64 // Max tip velocity of wheel, measured in meters/second
		MaxFlyWheelAcceleration float64 // Max tip acceleration of wheel, measured in meters/second/second
		MaxAzimuthVelocity      float64 // Max rotational velocity of azimuth action of shooter, measured in radians/second
		MaxAzimuthAcceleartion  float64 // Max rotational acceleration of azimuth action of shooter, measured in radians/second/second
		MinAzimuthOffset        float64
		MaxFeedVelocity         float64

		SpinnerMotorCanId                  int
		SpinnerMotorP                      float64
		SpinnerMotorD                      float64
		SpinnerMotorFF                     float64
		SpinnerMotorVelocityConversion     float64
		SpinnerMotorAccelerationConversion float64

		TiltMotorCanId                  int
		TiltMotorP                      float64
		TiltMotorI                      float64
		TiltMotorD                      float64
		TiltMotorPositionConversion     float64
		TiltMotorVelocityConversion     float64
		TiltMotorAccelerationConversion float64

		AzimuthMotorCanId                  int
		AzimuthMotorP                      float64
		AzimuthMotorI                      float64
		AzimuthMotorD                      float64
		AzimuthMotorPositionConversion     float64
		AzimuthMotorVelocityConversion     float64
		AzimuthMotorAccelerationConversion float64
	}

	MotorConfigs struct {
		SpinnerMotor motor.Config
		TiltMotor    motor.Config
		AzimuthMotor motor.Config
	}

	ShooterSubsystem struct {
		HasBall      bool
		ReadyToShoot bool
		SpinnerMotor *motor.Motor
		TiltMotor    *motor.Motor
		AzimuthMotor *motor.Motor
		StateMachine *state_machine.StateMachine
	}
)

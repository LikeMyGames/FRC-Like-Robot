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

		FlywheelMotorCanId                  int
		FlywheelMotorP                      float64
		FlywheelMotorD                      float64
		FlywheelMotorFF                     float64
		FlywheelMotorVelocityConversion     float64
		FlywheelMotorAccelerationConversion float64

		HoodMotorCanId                  int
		HoodMotorP                      float64
		HoodMotorI                      float64
		HoodMotorD                      float64
		HoodMotorPositionConversion     float64
		HoodMotorVelocityConversion     float64
		HoodMotorAccelerationConversion float64

		TurretMotorCanId                  int
		TurretMotorP                      float64
		TurretMotorI                      float64
		TurretMotorD                      float64
		TurretMotorPositionConversion     float64
		TurretMotorVelocityConversion     float64
		TurretMotorAccelerationConversion float64
	}

	MotorConfigs struct {
		FlywheelMotor *motor.Config
		HoodMotor     *motor.Config
		TurretMotor   *motor.Config
	}

	ShooterSubsystem struct {
		HasBall       bool
		ReadyToShoot  bool
		FlywheelMotor *motor.Motor
		HoodMotor     *motor.Motor
		TurretMotor   *motor.Motor
		StateMachine  *state_machine.StateMachine

		// RobotPose         mathutils.Pose2D // this is gets resolved by an event listener triggered from the drive subsystem
		// RobotPoseListener *event.Listener
	}
)

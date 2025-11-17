package shooter_types

import "github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"

type (
	ShooterConfig struct {
		MaxFlyWheelVelocity     float64 // Max tip velocity of wheel, measured in meters/second
		MaxFlyWheelAcceleration float64 // Max tip acceleration of wheel, measured in meters/second/second
		MaxAzimuthVelocity      float64 // Max rotational velocity of azimuth action of shooter, measured in radians/second
		MaxAzimuthAcceleartion  float64 // Max rotational acceleration of azimuth action of shooter, measured in radians/second/second
		MinAzimuthOffset        float64
		FlyWheelMotor           constantTypes.MotorController
		PitchMotor              constantTypes.MotorController
		FeedWheelMotor          constantTypes.MotorController
		AzimuthMotor            constantTypes.MotorController
	}
)

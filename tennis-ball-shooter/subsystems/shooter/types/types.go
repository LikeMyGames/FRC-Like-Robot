package shooter_types

type (
	ShooterConfig struct {
		MaxFlyWheelVelocity     float64 // Max tip velocity of wheel, measured in meters/second
		MaxFlyWheelAcceleration float64 // Max tip acceleration of wheel, measured in meters/second/second
		MaxAzimuthVelocity      float64 // Max rotational velocity of azimuth action of shooter, measured in radians/second
		MaxAzimuthAcceleartion  float64 // Max rotational acceleration of azimuth action of shooter, measured in radians/second/second
	}
)

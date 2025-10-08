package swerve

import (
	"time"

	constants "github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

type (
	SwerveDrive struct {
		Pose          mathutils.Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        constants.SwerveDriveConfig
		TimeInterval  time.Duration
	}

	SwerveDriveModules struct {
		FrontLeft  hardware.SwerveModule
		FrontRight hardware.SwerveModule
		BackLeft   hardware.SwerveModule
		BackRight  hardware.SwerveModule
	}

	SwerveDriveModulesVector struct {
		FrontLeft  mathutils.VectorTheta
		FrontRight mathutils.VectorTheta
		BackLeft   mathutils.VectorTheta
		BackRight  mathutils.VectorTheta
	}

	DriveProperties struct {
		TranslationalV mathutils.Vector2D
		RotationalV    float64
		TranslationalA mathutils.Vector2D
		RotationalA    float64
	}
)

func NewSwerveDrive(interval time.Duration, config constants.SwerveDriveConfig) *SwerveDrive {
	pose := mathutils.Pose2D{Location: mathutils.Vector2D{X: 0, Y: 0}, Angle: 0}
	config.MaxSpeed.RotationalV = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	swerve_modules := SwerveDriveModules{}

	return &SwerveDrive{
		Pose:          pose,
		Config:        config,
		SwerveModules: swerve_modules,
		TimeInterval:  interval,
	}
}

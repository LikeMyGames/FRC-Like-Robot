package constants

import (
	"math"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constants"
)

var Drive constants.DriveConstants = constants.SwerveDriveConfig{
	MaxSpeed: constants.DriveMaxes{
		TranslationalV: 1,   // Max Translational Velocity of the robot
		RotationalV:    180, // Max Rotational Velocity of the robot
		TranslationalA: 0.5, // Max Translational Acceleration of the robot
		RotationalA:    20,  // Max Rotational Acceleration of the robot
	},
	Modules: constants.Modules{
		FrontLeft: constants.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    math.Pi,
			DriveCanID:       10,
			AzimuthCanID:     11,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
		FrontRight: constants.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    math.Pi / 2,
			DriveCanID:       20,
			AzimuthCanID:     21,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
		BackLeft: constants.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    -math.Pi / 2,
			DriveCanID:       30,
			AzimuthCanID:     31,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
		BackRight: constants.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    0,
			DriveCanID:       40,
			AzimuthCanID:     41,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
	},
}

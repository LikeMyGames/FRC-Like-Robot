package constants

import (
	"math"

	constantTypes "github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
)

var Drive constantTypes.SwerveDriveConfig = constantTypes.SwerveDriveConfig{
	MaxSpeed: constantTypes.DriveMaxes{
		TranslationalV: 1,   // Max Translational Velocity of the robot
		RotationalV:    180, // Max Rotational Velocity of the robot
		TranslationalA: 0.5, // Max Translational Acceleration of the robot
		RotationalA:    20,  // Max Rotational Acceleration of the robot
	},
	Modules: constantTypes.Modules{
		FrontLeft: constantTypes.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    math.Pi,
			DriveCanID:       10,
			AzimuthCanID:     11,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
		FrontRight: constantTypes.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    math.Pi / 2,
			DriveCanID:       20,
			AzimuthCanID:     21,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
		BackLeft: constantTypes.ModuleConstants{
			OffsetX:          0.15,
			OffsetY:          0.15,
			AngularOffset:    -math.Pi / 2,
			DriveCanID:       30,
			AzimuthCanID:     31,
			DriveGearRatio:   9,
			AzimuthGearRatio: 16,
		},
		BackRight: constantTypes.ModuleConstants{
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

var Controller0 constantTypes.ControllerConfig = constantTypes.ControllerConfig{
	ControllerNum: 0,
	Precision:     2,
	Deadzones: constantTypes.ControllerDeadzones{
		ThumbL:   0.1,
		ThumbR:   0.1,
		TriggerL: 0.2,
		TriggerR: 0.2,
	},
}

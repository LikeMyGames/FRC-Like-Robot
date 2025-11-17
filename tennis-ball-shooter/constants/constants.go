// Do not edit in code, ONLY EDIT HERE
package constants

import (
	"math"

	// Importing the constant type from the FRC-Like-Robot State module
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
)

// The Drive contants defined for the robot
// Used in the Drive Subsystem
// Do not edit in code ONLY EDIT HERE
var Drive constantTypes.SwerveDriveConfig = constantTypes.SwerveDriveConfig{
	MaxSpeed: constantTypes.DriveMaxes{
		TranslationalV: 1,   // Max Translational Velocity of the robot
		RotationalV:    180, // Max Rotational Velocity of the robot in degrees per second
		TranslationalA: 0.5, // Max Translational Acceleration of the robot
		RotationalA:    20,  // Max Rotational Acceleration of the robot in degrees per second per second
	},
	Modules: constantTypes.Modules{
		FrontLeft: constantTypes.ModuleConstants{
			OffsetX:             0.15, // Offset from center (X direction)
			OffsetY:             0.15, // Offset from center (Y direction)
			AngularOffset:       math.Pi,
			DriveCanID:          10,
			AzimuthCanID:        11,
			DriveGearRatio:      23 / 44,
			DriveGearRatioInv:   44 / 23,
			AzimuthGearRatio:    4 / 57,
			AzimuthGearRatioInv: 57 / 4,
		},
		FrontRight: constantTypes.ModuleConstants{
			OffsetX:             0.15,
			OffsetY:             0.15,
			AngularOffset:       math.Pi / 2,
			DriveCanID:          20,
			AzimuthCanID:        21,
			DriveGearRatio:      23 / 44,
			DriveGearRatioInv:   44 / 23,
			AzimuthGearRatio:    4 / 57,
			AzimuthGearRatioInv: 57 / 4,
		},
		BackLeft: constantTypes.ModuleConstants{
			OffsetX:             0.15,
			OffsetY:             0.15,
			AngularOffset:       -math.Pi / 2,
			DriveCanID:          30,
			AzimuthCanID:        31,
			DriveGearRatio:      23 / 44,
			DriveGearRatioInv:   44 / 23,
			AzimuthGearRatio:    4 / 57,
			AzimuthGearRatioInv: 57 / 4,
		},
		BackRight: constantTypes.ModuleConstants{
			OffsetX:             0.15,
			OffsetY:             0.15,
			AngularOffset:       0,
			DriveCanID:          40,
			AzimuthCanID:        41,
			DriveGearRatio:      23 / 44,
			DriveGearRatioInv:   44 / 23,
			AzimuthGearRatio:    4 / 57,
			AzimuthGearRatioInv: 57 / 4,
		},
	},
}

// Contansts for the main controller used
// Used to instantiate a controller object in the main.go file
var Controller0 constantTypes.ControllerConfig = constantTypes.ControllerConfig{
	ControllerNum: 0,
	Precision:     2,
	Deadzones: constantTypes.ControllerDeadzones{
		ThumbL:   0.05,
		ThumbR:   0.05,
		TriggerL: 0.05,
		TriggerR: 0.05,
	},
	MinChange: 0.01,
}

var Shooter shooter_types.ShooterConfig = shooter_types.ShooterConfig{
	MaxFlyWheelVelocity:     10,
	MaxFlyWheelAcceleration: 1,
	MaxAzimuthVelocity:      math.Pi / 4,
	MaxAzimuthAcceleartion:  math.Pi / 16,
	MinAzimuthOffset:        math.Pi / 32,
	FlyWheelMotor: constantTypes.MotorController{
		Id:  52,
		PID: constantTypes.PidController{Kp: 0, Ki: 0, Kd: 0},
	},
	PitchMotor: constantTypes.MotorController{
		Id:  51,
		PID: constantTypes.PidController{Kp: 0, Ki: 0, Kd: 0},
	},
	FeedWheelMotor: constantTypes.MotorController{
		Id:  53,
		PID: constantTypes.PidController{Kp: 0, Ki: 0, Kd: 0},
	},
	AzimuthMotor: constantTypes.MotorController{
		Id:  50,
		PID: constantTypes.PidController{Kp: 0, Ki: 0, Kd: 0},
	},
}

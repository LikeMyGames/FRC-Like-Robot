// Do not edit in code, ONLY EDIT HERE
package constants

import (
	"math"
	"time"

	// Importing the constant type from the FRC-Like-Robot State module

	drive_types "tennis-ball-shooter/subsystems/drive/types"
	intake_types "tennis-ball-shooter/subsystems/intake/types"
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/drive/swerve"
)

type (
//	RobotType struct {
//		Frequency   time.Duration
//		Controllers []constantTypes.ControllerConfig
//		Drive       constantTypes.SwerveDriveConfig
//		Shooter     shooter_types.ShooterConfig
//		Battery     constantTypes.Battery
//	}
)

// var Robot RobotType = RobotType{
// 	Frequency: time.Millisecond * 100,
// 	Controllers: []constantTypes.ControllerConfig{
// 		Controller0,
// 	},
// 	Drive:   Drive,
// 	Shooter: Shooter,
// 	Battery: Battery,
// }

var Robot constantTypes.RobotConfig = constantTypes.RobotConfig{
	Period:        time.Millisecond * 100,
	StartingState: "power_on",
	RslPin:        22,
}

// The Drive contants defined for the robot
// Used in the Drive Subsystem
// Do not edit in code ONLY EDIT HERE
var Drive drive_types.Constants = drive_types.Constants{
	PositionalPidP: 0,
	PositionalPidI: 0,
	PositionalPidD: 0,
	Swerve: swerve.Config{
		MaxSpeed: swerve.DriveMaxes{
			TranslationalV: 2,   // Max Translational Velocity of the robot
			RotationalV:    180, // Max Rotational Velocity of the robot in degrees per second
			TranslationalA: 0.5, // Max Translational Acceleration of the robot
			RotationalA:    20,  // Max Rotational Acceleration of the robot in degrees per second per second
		},
		Modules: []swerve.ModuleConstants{
			{
				Name:                "FrontLeft",
				OffsetX:             0.15, // Offset from center (X direction)
				OffsetY:             0.15, // Offset from center (Y direction)
				AngularOffset:       0,
				DriveMotorCanId:     10,
				AzimuthMotorCanId:   11,
				DriveGearRatio:      23.0 / 44.0,
				DriveGearRatioInv:   44.0 / 23.0,
				AzimuthGearRatio:    4.0 / 57.0,
				AzimuthGearRatioInv: 57.0 / 4.0,
			},
			{
				Name:                "FrontRight",
				OffsetX:             0.15,
				OffsetY:             -0.15,
				AngularOffset:       -math.Pi / 2,
				DriveMotorCanId:     20,
				AzimuthMotorCanId:   21,
				DriveGearRatio:      23.0 / 44.0,
				DriveGearRatioInv:   44.0 / 23.0,
				AzimuthGearRatio:    4.0 / 57.0,
				AzimuthGearRatioInv: 57.0 / 4.0,
			},
			{
				Name:                "BackLeft",
				OffsetX:             -0.15,
				OffsetY:             0.15,
				AngularOffset:       math.Pi / 2,
				DriveMotorCanId:     30,
				AzimuthMotorCanId:   31,
				DriveGearRatio:      23.0 / 44.0,
				DriveGearRatioInv:   44.0 / 23.0,
				AzimuthGearRatio:    4.0 / 57.0,
				AzimuthGearRatioInv: 57.0 / 4.0,
			},
			{
				Name:                "BackRight",
				OffsetX:             -0.15,
				OffsetY:             -0.15,
				AngularOffset:       math.Pi,
				DriveMotorCanId:     40,
				AzimuthMotorCanId:   41,
				DriveGearRatio:      23.0 / 44.0,
				DriveGearRatioInv:   44.0 / 23.0,
				AzimuthGearRatio:    4.0 / 57.0,
				AzimuthGearRatioInv: 57.0 / 4.0,
			},
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
	MinChange: 0.05,
}

var Shooter shooter_types.Constants = shooter_types.Constants{
	MaxFlyWheelVelocity:     10,
	MaxFlyWheelAcceleration: 1,
	MaxAzimuthVelocity:      math.Pi / 4,
	MaxAzimuthAcceleartion:  math.Pi / 16,
	MinAzimuthOffset:        math.Pi / 32,
	MaxFeedVelocity:         2,

	FlywheelMotorCanId:                  52,
	FlywheelMotorP:                      0,
	FlywheelMotorD:                      0,
	FlywheelMotorFF:                     0,
	FlywheelMotorVelocityConversion:     1,
	FlywheelMotorAccelerationConversion: 1,

	HoodMotorCanId:                  51,
	HoodMotorP:                      0,
	HoodMotorI:                      0,
	HoodMotorD:                      0,
	HoodMotorPositionConversion:     1,
	HoodMotorVelocityConversion:     1,
	HoodMotorAccelerationConversion: 1,

	TurretMotorCanId:                  50,
	TurretMotorP:                      0,
	TurretMotorI:                      0,
	TurretMotorD:                      0,
	TurretMotorPositionConversion:     1,
	TurretMotorVelocityConversion:     1,
	TurretMotorAccelerationConversion: 1,
}

var Intake intake_types.Constants = intake_types.Constants{
	RollerMotorCanId:    50,
	ExtensionMotorCanId: 51,

	RollerMotorP:                      0,
	RollerMotorD:                      0,
	RollerMotorFF:                     0,
	RollerMotorVelocityConversion:     1,
	RollerMotorAccelerationConversion: 1,

	ExtensionMotorP:                      0,
	ExtensionMotorI:                      0,
	ExtensionMotorD:                      0,
	ExtensionMotorCosFF:                  0,
	ExtensionMotorVelocityConversion:     1,
	ExtensionMotorAccelerationConversion: 1,
}

var Battery constantTypes.Battery = constantTypes.Battery{
	NominalVoltage: 12,
}

const (
	DriverControllerNum int = 0
)

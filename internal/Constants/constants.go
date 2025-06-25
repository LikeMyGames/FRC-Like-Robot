package Constants

import "math"

type (
	ControllerConfig struct {
		ControllerNum int                 `json:"controllerNum"`
		Precision     int                 `json:"precision"`
		Deadzones     ControllerDeadzones `json:"deadzones"`
	}

	ControllerDeadzones struct {
		ThumbL   float32 `json:"thumbL"`
		ThumbR   float32 `json:"thumbR"`
		TriggerL float32 `json:"triggerL"`
		TriggerR float32 `json:"triggerR"`
	}

	SwerveDriveConfig struct {
		MaxSpeed DriveMaxes `json:"drive_maxes"`
		// ModuleConfig ModuleConfig `json:"module_config"`
		Modules Modules `json:"modules"`
	}

	DriveMaxes struct {
		TranslationalV float64 `json:"translationalV"`
		RotationalV    float64 `json:"rotationalV"`
		TranslationalA float64 `json:"translationalA"`
		RotationalA    float64 `json:"rotationalA"`
	}

	ModuleConfig struct {
		MaxTransSpeed   float64 `json:"maxTrans_speed"`
		MaxTransAccel   float64 `json:"maxTrans_accel"`
		MaxRotSpeed     float64 `json:"maxRot_speed"`
		MaxRotAccel     float64 `json:"maxRot_accel"`
		WheelRadius     float64 `json:"wheel_radius"`
		DriveBaseWidth  float64 `json:"driveBase_width"`
		DriveBaseLength float64 `json:"driveBase_length"`
	}

	Modules struct {
		FrontLeft  ModuleConstants `json:"frontLeft"`
		FrontRight ModuleConstants `json:"frontRight"`
		BackLeft   ModuleConstants `json:"backLeft"`
		BackRight  ModuleConstants `json:"backRight"`
	}

	ModuleConstants struct {
		OffsetX          float64 `json:"x"`
		OffsetY          float64 `json:"y"`
		AngularOffset    float64 `json:"angular_Offeset"`
		DriveCanID       uint    `json:"canID_drive"`
		AzimuthCanID     uint    `json:"canID_azimuth"`
		DriveGearRatio   float64 `json:"gearRatio_drive"`
		AzimuthGearRatio float64 `json:"gearRatio_azimuth"`
	}
)

func ControllerConstants() ControllerConfig {
	return ControllerConfig{
		ControllerNum: 0,
		Precision:     2,
		Deadzones: ControllerDeadzones{
			ThumbL:   0.1,
			ThumbR:   0.1,
			TriggerL: 0.2,
			TriggerR: 0.2,
		},
	}
}

func DriveConstants() SwerveDriveConfig {
	return SwerveDriveConfig{
		MaxSpeed: DriveMaxes{
			TranslationalV: 1,   // Max Translational Velocity of the robot
			RotationalV:    180, // Max Rotational Velocity of the robot
			TranslationalA: 0.5, // Max Translational Acceleration of the robot
			RotationalA:    20,  // Max Rotational Acceleration of the robot
		},
		// ModuleConfig: ModuleConfig{
		// 	MaxTransSpeed:   1, // Max
		// 	MaxTransAccel:   0.5,
		// 	MaxRotSpeed:     180,
		// 	MaxRotAccel:     20,
		// 	WheelRadius:     0.0375,
		// 	DriveBaseWidth:  0.3,
		// 	DriveBaseLength: 0.3,
		// },
		Modules: Modules{
			FrontLeft: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				AngularOffset:    math.Pi,
				DriveCanID:       10,
				AzimuthCanID:     11,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
			FrontRight: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				AngularOffset:    math.Pi / 2,
				DriveCanID:       20,
				AzimuthCanID:     21,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
			BackLeft: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				AngularOffset:    -math.Pi / 2,
				DriveCanID:       30,
				AzimuthCanID:     31,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
			BackRight: ModuleConstants{
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
}

package Constants

type (
	SwerveDriveConfig struct {
		MaxSpeed     DriveMaxes   `json:"drive_maxes"`
		ModuleConfig ModuleConfig `json:"module_config"`
		Modules      Modules      `json:"modules"`
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
		DriveCanID       uint    `json:"canID_drive"`
		AzimuthCanID     uint    `json:"canID_azimuth"`
		DriveGearRatio   float64 `json:"gearRatio_drive"`
		AzimuthGearRatio float64 `json:"gearRatio_azimuth"`
	}
)

func DriveConstants() SwerveDriveConfig {
	return SwerveDriveConfig{
		MaxSpeed: DriveMaxes{
			TranslationalV: 1,
			RotationalV:    180,
			TranslationalA: 0.5,
			RotationalA:    20,
		},
		ModuleConfig: ModuleConfig{
			MaxTransSpeed:   1,
			MaxTransAccel:   0.5,
			MaxRotSpeed:     180,
			MaxRotAccel:     20,
			WheelRadius:     0.0375,
			DriveBaseWidth:  0.3,
			DriveBaseLength: 0.3,
		},
		Modules: Modules{
			FrontLeft: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				DriveCanID:       10,
				AzimuthCanID:     11,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
			FrontRight: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				DriveCanID:       20,
				AzimuthCanID:     21,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
			BackLeft: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				DriveCanID:       30,
				AzimuthCanID:     31,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
			BackRight: ModuleConstants{
				OffsetX:          0.15,
				OffsetY:          0.15,
				DriveCanID:       40,
				AzimuthCanID:     41,
				DriveGearRatio:   9,
				AzimuthGearRatio: 16,
			},
		},
	}
}

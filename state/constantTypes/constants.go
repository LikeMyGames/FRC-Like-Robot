package constants

type (
	ControllerConfig struct {
		ControllerNum int                 `json:"controllerNum"`
		Precision     int                 `json:"precision"`
		Deadzones     ControllerDeadzones `json:"deadzones"`
		MinChange     float64
	}

	ControllerDeadzones struct {
		ThumbL   float64 `json:"thumbL"`
		ThumbR   float64 `json:"thumbR"`
		TriggerL float64 `json:"triggerL"`
		TriggerR float64 `json:"triggerR"`
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

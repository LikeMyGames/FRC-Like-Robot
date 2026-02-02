package constantTypes

import "time"

type (
	RobotConfig struct {
		StartingState string
		Period        time.Duration
		RslPin        int
	}
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
		// Modules Modules `json:"modules"`
		Modules []ModuleConstants
	}

	DriveMaxes struct {
		TranslationalV float64 `json:"translationalV"`
		RotationalV    float64 `json:"rotationalV"`
		TranslationalA float64 `json:"translationalA"`
		RotationalA    float64 `json:"rotationalA"`
	}

	ModuleConfig struct {
		MaxTransSpeed   float64 `json:"maxTrans_speed"`   // in meters per second
		MaxTransAccel   float64 `json:"maxTrans_accel"`   // in meters per second per second
		MaxRotSpeed     float64 `json:"maxRot_speed"`     // in radians per second
		MaxRotAccel     float64 `json:"maxRot_accel"`     // in radians per second per second
		WheelRadius     float64 `json:"wheel_radius"`     // in meters
		DriveBaseWidth  float64 `json:"driveBase_width"`  // in meters
		DriveBaseLength float64 `json:"driveBase_length"` // in meters
	}

	Modules struct {
		FrontLeft  ModuleConstants `json:"frontLeft"`
		FrontRight ModuleConstants `json:"frontRight"`
		BackLeft   ModuleConstants `json:"backLeft"`
		BackRight  ModuleConstants `json:"backRight"`
	}

	ModuleConstants struct {
		Name                string  `json:"name"`
		OffsetX             float64 `json:"x"`               // in meters
		OffsetY             float64 `json:"y"`               // in meters
		AngularOffset       float64 `json:"angular_Offeset"` // in radians
		DriveCanID          uint    `json:"canID_drive"`
		AzimuthCanID        uint    `json:"canID_azimuth"`
		DriveGearRatio      float64 `json:"gearRatio_drive"`
		DriveGearRatioInv   float64 `json:"gearRatio_driveInv"`
		AzimuthGearRatio    float64 `json:"gearRatio_azimuth"`
		AzimuthGearRatioInv float64 `json:"gearRatio_azimuthInv"`
	}

	MotorController struct {
		Id  int64         `json:"CanID"`
		PID PidController `json:"PID"`
	}

	PidController struct {
		Kp     float64
		Ki     float64
		Kd     float64
		MaxOut float64
		MinOut float64
	}

	Battery struct {
		NominalVoltage float64
	}
)

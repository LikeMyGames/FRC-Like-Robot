package swerve

import (
	"math"

	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"

	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type (
	SwerveModule struct {
		driveMotor    *motor.Motor
		turningMotor  *motor.Motor
		targetVector  mathutils.VectorTheta
		angularOffset float64
	}

	SwerveDrive struct {
		DriveProps    DriveProperties
		SwerveModules []SwerveModule
		Config        Config
	}

	SwerveDriveModules struct {
		FrontLeft  *SwerveModule
		FrontRight *SwerveModule
		BackLeft   *SwerveModule
		BackRight  *SwerveModule
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

	Config struct {
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

	ModuleConstants struct {
		Name                string  `json:"name"`
		OffsetX             float64 `json:"x"`               // in meters
		OffsetY             float64 `json:"y"`               // in meters
		AngularOffset       float64 `json:"angular_Offeset"` // in radians
		DriveMotorConfig    motor.Config
		AzimuthMotorConfig  motor.Config
		DriveGearRatio      float64 `json:"gearRatio_drive"`
		DriveGearRatioInv   float64 `json:"gearRatio_driveInv"`
		AzimuthGearRatio    float64 `json:"gearRatio_azimuth"`
		AzimuthGearRatioInv float64 `json:"gearRatio_azimuthInv"`
	}

	// ModuleConfig struct {
	// 	MaxTransSpeed   float64 `json:"maxTrans_speed"`   // in meters per second
	// 	MaxTransAccel   float64 `json:"maxTrans_accel"`   // in meters per second per second
	// 	MaxRotSpeed     float64 `json:"maxRot_speed"`     // in radians per second
	// 	MaxRotAccel     float64 `json:"maxRot_accel"`     // in radians per second per second
	// 	WheelRadius     float64 `json:"wheel_radius"`     // in meters
	// 	DriveBaseWidth  float64 `json:"driveBase_width"`  // in meters
	// 	DriveBaseLength float64 `json:"driveBase_length"` // in meters
	// }
)

var (
	// ctrl               *controller.Controller = nil
	ctrlTrans          mathutils.Vector2D = mathutils.Vector2D{}
	ctrlRot            mathutils.Vector2D = mathutils.Vector2D{}
	transEventTarget   string             = controller.LeftStick
	transEventListener *event.Listener    = nil
	rotEventTarget     string             = controller.RightStick
	rotEventListener   *event.Listener    = nil
)

func NewSwerveDrive(config Config) *SwerveDrive {
	config.MaxSpeed.RotationalV = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	// fmt.Println("Drive Config: ", config)
	swerve_modules := make([]SwerveModule, len(config.Modules))
	for i, v := range config.Modules {
		swerve_modules[i] = SwerveModule{
			driveMotor:    motor.New(v.DriveMotorConfig),
			turningMotor:  motor.New(v.AzimuthMotorConfig),
			targetVector:  mathutils.VectorTheta{},
			angularOffset: v.AngularOffset,
		}
		// fmt.Printf("Created new swerve drive with drive motor id: %d; and turning motor id: %d\n", v.DriveCanID, v.AzimuthCanID)
		swerve_modules[i].turningMotor.SetIsSecondaryMotorOnController(true)
	}

	return &SwerveDrive{
		Config:        config,
		SwerveModules: swerve_modules,
	}
}

func (drive *SwerveDrive) CalculateSwerveFromSavedControllerVals() {
	drive.CalculateSwerve(ctrlTrans.X, ctrlTrans.Y, ctrlRot.X)
}

func (drive *SwerveDrive) CalculateSwerve(x, y, rot float64) []mathutils.VectorTheta {
	x = mathutils.Clamp(x, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)
	y = mathutils.Clamp(y, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)
	rot = mathutils.Clamp(rot, drive.Config.MaxSpeed.RotationalV, -drive.Config.MaxSpeed.RotationalV)
	drive.DriveProps.TranslationalV = mathutils.Vector2D{X: x, Y: y}
	drive.DriveProps.RotationalV = rot

	states := make([]mathutils.VectorTheta, len(drive.SwerveModules))

	for i, v := range drive.Config.Modules {
		// offset := v.AngularOffset
		moduleOffsetAngle := mathutils.Vector2D{X: v.OffsetX, Y: v.OffsetY}.ToVectorTheta().Angle
		distance := math.Hypot(v.OffsetX, v.OffsetY)

		rotVector := mathutils.Vector2D{X: 0, Y: distance}
		rotVector.Rotate(moduleOffsetAngle).Multiply(rot)

		vector := mathutils.AddVector2D(mathutils.Vector2D{X: x, Y: y}, rotVector)
		newstate := vector.ToVectorTheta()
		oldstate := drive.SwerveModules[i].targetVector
		angleErr := newstate.Angle - oldstate.Angle
		if angleErr > math.Pi/2 || angleErr < -math.Pi/2 {
			newstate.Magnitude *= -1
		}
		states[i] = newstate
	}

	return states
}

// func (drive *SwerveDrive) Normalize(states *[]mathutils.VectorTheta, x, y, rot float64) {
// 	maxModuleSpeed := 0.0
// 	for _, v := range *states {
// 		maxModuleSpeed = math.Max(maxModuleSpeed, math.Abs(v.Magnitude))
// 	}

// 	if drive.Config.MaxSpeed.TranslationalV == 0 || drive.Config.MaxSpeed.RotationalV == 0 || maxModuleSpeed == 0 {
// 		return
// 	}

// 	// math.Hypot(x, y) / maxModuleSpeed
// 	for _, v := range *states {
// 		v.Magnitude *= maxModuleSpeed / drive.Config.MaxSpeed.TranslationalV
// 	}
// }

func (drive *SwerveDrive) Normalize(states *[]mathutils.VectorTheta, xSpeed, ySpeed, rotSpeed float64) {
	realMaxSpeed := 0.0
	for _, v := range *states {
		realMaxSpeed = math.Max(realMaxSpeed, math.Abs(v.Magnitude))
	}

	if drive.Config.MaxSpeed.TranslationalV == 0 || drive.Config.MaxSpeed.RotationalV == 0 || realMaxSpeed == 0 {
		return
	}

	translationalK := math.Hypot(xSpeed, ySpeed) / drive.Config.MaxSpeed.TranslationalV
	rotationalK := math.Abs(rotSpeed) / drive.Config.MaxSpeed.RotationalV
	k := math.Max(translationalK, rotationalK)
	scale := math.Min(k*drive.Config.MaxSpeed.TranslationalV/realMaxSpeed, 1)
	for _, v := range *states {
		v.Magnitude *= scale
	}
}

// func SetDriveController(controller *controller.Controller) {
// 	ctrl = controller
// }

func SetTransEventTarget(target string) {
	event.Remove(transEventListener)
	if target != "" {
		transEventListener = event.Listen(target, "DRIVE_SUBSYSTEM_TRANS", func(event any) {
			trans := event.(mathutils.Vector2D)
			ctrlTrans = trans
		})
	}
}

func GetTransEventTarget() string {
	return transEventTarget
}

func SetRotEventTarget(target string) {
	event.Remove(rotEventListener)
	if target != "" {
		rotEventListener = event.Listen(target, "DRIVE_SUBSYSTEM_ROT", func(event any) {
			rot := event.(mathutils.Vector2D)
			ctrlRot = rot
		})
	}
}

func GetRotEventTarget() string {
	return rotEventTarget
}

func (m *SwerveModule) ReadAzimuthAngle() float64 {

	return 0
}

func (m *SwerveModule) SetTarget(target mathutils.VectorTheta) {
	m.targetVector = target
	m.driveMotor.SetVelocity(m.targetVector.Magnitude)
	m.turningMotor.SetAngle(m.targetVector.Angle + m.angularOffset)
}

// func GetDriveVectorsFromController(ctrl *controller.Controller) (trans, rot mathutils.Vector2D) {
// 	return mathutils.Vector2D{}, mathutils.Vector2D{}
// }

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

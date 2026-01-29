package swerve

import (
	"fmt"
	"math"

	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

type (
	SwerveModule struct {
		driveMotor   *motor.Motor
		turningMotor *motor.Motor
		targetVector mathutils.VectorTheta
	}

	SwerveDrive struct {
		DriveProps    DriveProperties
		SwerveModules []SwerveModule
		Config        constantTypes.SwerveDriveConfig
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

func NewSwerveDrive(config constantTypes.SwerveDriveConfig) *SwerveDrive {
	config.MaxSpeed.RotationalV = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	// fmt.Println("Drive Config: ", config)
	swerve_modules := make([]SwerveModule, 0)
	for i, v := range config.Modules {
		swerve_modules[i] = SwerveModule{
			driveMotor:   motor.NewMotor(int(v.DriveCanID)),
			turningMotor: motor.NewMotor(int(v.AzimuthCanID)),
			targetVector: mathutils.VectorTheta{},
		}
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
	drive.DriveProps.TranslationalV = mathutils.Vector2D{X: mathutils.Clamp(x, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV), Y: mathutils.Clamp(y, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)}
	drive.DriveProps.RotationalV = mathutils.Clamp(rot, drive.Config.MaxSpeed.RotationalV, -drive.Config.MaxSpeed.RotationalV)

	states := make([]mathutils.VectorTheta, len(drive.SwerveModules))

	for i, v := range drive.Config.Modules {
		offset := v.AngularOffset
		distance := math.Hypot(v.OffsetX, v.OffsetY)

		rotVector := mathutils.Vector2D{X: 0, Y: distance}
		rotVector.Rotate(offset)

		vector := mathutils.VectorAdd(mathutils.Vector2D{X: x, Y: y}, rotVector)
		newstate := vector.ToVectorTheta()
		oldstate := drive.SwerveModules[i].targetVector
		angleErr := newstate.Angle - oldstate.Angle
		if angleErr > math.Pi/2 || angleErr < -math.Pi/2 {
			newstate.Magnitude *= -1
		}
		states[i] = newstate
	}

	// // Front Left Wheel Calculation
	// flOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.FrontLeft.OffsetX, Y: drive.Config.Modules.FrontLeft.OffsetY})
	// fl := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: flOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * flOffset.Magnitude})))
	// // drive.SwerveModules.FrontLeft

	// // Front Right Wheel Calculation
	// frOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.FrontRight.OffsetX, Y: drive.Config.Modules.FrontRight.OffsetY})
	// fr := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: frOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * frOffset.Magnitude})))

	// // Back Left Wheel Calculation
	// blOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.BackLeft.OffsetX, Y: drive.Config.Modules.BackLeft.OffsetY})
	// bl := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: blOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * blOffset.Magnitude})))

	// // Back Right Wheel Calculation
	// brOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.BackRight.OffsetX, Y: drive.Config.Modules.BackRight.OffsetY})
	// br := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: brOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * brOffset.Magnitude})))

	fmt.Println(states)

	// Setting the swerve calculations in the drive objects pointer
	return states
}

func (drive *SwerveDrive) Normalize(states []mathutils.VectorTheta, xSpeed, ySpeed, rotSpeed float64) {
	realMaxSpeed := 0.0
	for _, v := range states {
		realMaxSpeed = math.Max(realMaxSpeed, math.Abs(v.Magnitude))
	}

	if drive.Config.MaxSpeed.TranslationalV == 0 || drive.Config.MaxSpeed.RotationalV == 0 || realMaxSpeed == 0 {
		return
	}

	translationalK := math.Hypot(xSpeed, ySpeed) / drive.Config.MaxSpeed.TranslationalV
	rotationalK := math.Abs(rotSpeed) / drive.Config.MaxSpeed.RotationalV
	k := math.Max(translationalK, rotationalK)
	scale := math.Min(k*drive.Config.MaxSpeed.TranslationalV/realMaxSpeed, 1)
	for _, v := range states {
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
	m.turningMotor.SetAngle(m.targetVector.Angle)
}

// func GetDriveVectorsFromController(ctrl *controller.Controller) (trans, rot mathutils.Vector2D) {
// 	return mathutils.Vector2D{}, mathutils.Vector2D{}
// }

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

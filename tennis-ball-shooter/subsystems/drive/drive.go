package drive

import (
	"fmt"
	"math"

	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/drive/swerve"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

type (
	SwerveDrive struct {
		Pose          mathutils.Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        constantTypes.SwerveDriveConfig
	}

	SwerveDriveModules struct {
		FrontLeft  swerve.SwerveModule
		FrontRight swerve.SwerveModule
		BackLeft   swerve.SwerveModule
		BackRight  swerve.SwerveModule
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
	pose := mathutils.Pose2D{Location: mathutils.Vector2D{X: 0, Y: 0}, Angle: 0}
	config.MaxSpeed.RotationalV = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	// fmt.Println("Drive Config: ", config)
	swerve_modules := SwerveDriveModules{}

	return &SwerveDrive{
		Pose:          pose,
		Config:        config,
		SwerveModules: swerve_modules,
	}
}

func (drive *SwerveDrive) CalculateSwerveFromSavedControllerVals() {
	drive.CalculateSwerve(ctrlTrans, ctrlRot)
}

func (drive *SwerveDrive) CalculateSwerve(trans, rot mathutils.Vector2D) SwerveDriveModulesVector {
	drive.DriveProps.TranslationalV = mathutils.Vector2D{X: mathutils.Clamp(trans.X, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV), Y: mathutils.Clamp(trans.Y, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)}
	drive.DriveProps.RotationalV = mathutils.Clamp(rot.X, drive.Config.MaxSpeed.RotationalV, -drive.Config.MaxSpeed.RotationalV)

	// Front Left Wheel Calculation
	flOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.FrontLeft.OffsetX, Y: drive.Config.Modules.FrontLeft.OffsetY})
	fl := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: flOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * flOffset.Magnitude})))
	// drive.SwerveModules.FrontLeft

	// Front Right Wheel Calculation
	frOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.FrontRight.OffsetX, Y: drive.Config.Modules.FrontRight.OffsetY})
	fr := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: frOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * frOffset.Magnitude})))

	// Back Left Wheel Calculation
	blOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.BackLeft.OffsetX, Y: drive.Config.Modules.BackLeft.OffsetY})
	bl := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: blOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * blOffset.Magnitude})))

	// Back Right Wheel Calculation
	brOffset := mathutils.Vector2DtoVectorTheta(mathutils.Vector2D{X: drive.Config.Modules.BackRight.OffsetX, Y: drive.Config.Modules.BackRight.OffsetY})
	br := mathutils.Vector2DtoVectorTheta(mathutils.VectorAdd(drive.DriveProps.TranslationalV, mathutils.VectorThetatoVector2D(mathutils.VectorTheta{Angle: brOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * brOffset.Magnitude})))

	// fmt.Println(SwerveDriveModulesVector{
	// 	FrontLeft:  fl,
	// 	FrontRight: fr,
	// 	BackLeft:   bl,
	// 	BackRight:  br,
	// })

	// Setting the swerve calculations in the drive objects pointer
	return SwerveDriveModulesVector{
		FrontLeft:  fl,
		FrontRight: fr,
		BackLeft:   bl,
		BackRight:  br,
	}
}

func (drive *SwerveDrive) DriveWheels(modules SwerveDriveModules) {

}

func (drive *SwerveDrive) DriveToPose(pose mathutils.Pose2D) {
	diff := mathutils.Pose2D{}
	diff.Location = mathutils.Vector2D{X: pose.Location.X - drive.Pose.Location.X, Y: pose.Location.Y - drive.Pose.Location.Y}
	diff.Angle = pose.Angle - drive.Pose.Angle
	fmt.Println(diff)
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

// func GetDriveVectorsFromController(ctrl *controller.Controller) (trans, rot mathutils.Vector2D) {
// 	return mathutils.Vector2D{}, mathutils.Vector2D{}
// }

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

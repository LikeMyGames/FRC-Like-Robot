package drive

import (
	"fmt"
	"math"
	"test/constants"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"

	constantTypes "github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

type (
	SwerveDrive struct {
		Pose          mathutils.Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        constantTypes.SwerveDriveConfig
		TimeInterval  time.Duration
	}

	SwerveDriveModules struct {
		FrontLeft  hardware.SwerveModule
		FrontRight hardware.SwerveModule
		BackLeft   hardware.SwerveModule
		BackRight  hardware.SwerveModule
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

func NewSwerveDrive(interval time.Duration) *SwerveDrive {
	pose := mathutils.Pose2D{Location: mathutils.Vector2D{X: 0, Y: 0}, Angle: 0}
	config := constants.Drive
	config.MaxSpeed.RotationalV = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = mathutils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	fmt.Println("Drive Config: ", config)
	swerve_modules := SwerveDriveModules{}

	return &SwerveDrive{
		Pose:          pose,
		Config:        config,
		SwerveModules: swerve_modules,
		TimeInterval:  interval,
	}
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

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

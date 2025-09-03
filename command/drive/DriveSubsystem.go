package drive

import (
	"fmt"
	"frcrobot/constants"
	"frcrobot/hardware"
	"frcrobot/utils/MathUtils"
	"frcrobot/utils/Types"
	"frcrobot/utils/VectorMath"
	"math"
	"time"
)

type (
	SwerveDrive struct {
		Pose          Types.Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        constants.SwerveDriveConfig
		TimeInterval  time.Duration
	}

	SwerveDriveModules struct {
		FrontLeft  SwerveDriveModule
		FrontRight SwerveDriveModule
		BackLeft   SwerveDriveModule
		BackRight  SwerveDriveModule
	}

	SwerveDriveModulesVector struct {
		FrontLeft  Types.VectorTheta
		FrontRight Types.VectorTheta
		BackLeft   Types.VectorTheta
		BackRight  Types.VectorTheta
	}

	SwerveDriveModule struct {
		DriveVelocity float64
		AzimuthAngle  float64
		DriveMotor    hardware.MotorController
		AzimuthMotor  hardware.MotorController
	}

	DriveProperties struct {
		TranslationalV VectorMath.Vector2D
		RotationalV    float64
		TranslationalA VectorMath.Vector2D
		RotationalA    float64
	}
)

func NewSwerveDrive(interval time.Duration) *SwerveDrive {
	pose := Types.Pose2D{Location: VectorMath.Vector2D{X: 0, Y: 0}, Angle: 0}
	config := constants.DriveConstants()
	config.MaxSpeed.RotationalV = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	fmt.Println("Drive Config: ", config)
	swerve_modules := SwerveDriveModules{}

	return &SwerveDrive{
		Pose:          pose,
		Config:        config,
		SwerveModules: swerve_modules,
		TimeInterval:  interval,
	}
}

func (drive *SwerveDrive) CalculateSwerve(trans, rot Types.Vector2D) SwerveDriveModulesVector {
	drive.DriveProps.TranslationalV = Types.Vector2D{X: MathUtils.Clamp(trans.X, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV), Y: MathUtils.Clamp(trans.Y, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)}
	drive.DriveProps.RotationalV = MathUtils.Clamp(rot.X, drive.Config.MaxSpeed.RotationalV, -drive.Config.MaxSpeed.RotationalV)

	// Front Left Wheel Calculation
	flOffset := VectorMath.Vector2DtoVectorTheta(Types.Vector2D{X: drive.Config.Modules.FrontLeft.OffsetX, Y: drive.Config.Modules.FrontLeft.OffsetY})
	fl := VectorMath.Vector2DtoVectorTheta(VectorMath.VectorAdd(drive.DriveProps.TranslationalV, VectorMath.VectorThetatoVector2D(Types.VectorTheta{Angle: flOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * flOffset.Magnitude})))
	// drive.SwerveModules.FrontLeft

	// Front Right Wheel Calculation
	frOffset := VectorMath.Vector2DtoVectorTheta(Types.Vector2D{X: drive.Config.Modules.FrontRight.OffsetX, Y: drive.Config.Modules.FrontRight.OffsetY})
	fr := VectorMath.Vector2DtoVectorTheta(VectorMath.VectorAdd(drive.DriveProps.TranslationalV, VectorMath.VectorThetatoVector2D(Types.VectorTheta{Angle: frOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * frOffset.Magnitude})))

	// Back Left Wheel Calculation
	blOffset := VectorMath.Vector2DtoVectorTheta(Types.Vector2D{X: drive.Config.Modules.BackLeft.OffsetX, Y: drive.Config.Modules.BackLeft.OffsetY})
	bl := VectorMath.Vector2DtoVectorTheta(VectorMath.VectorAdd(drive.DriveProps.TranslationalV, VectorMath.VectorThetatoVector2D(Types.VectorTheta{Angle: blOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * blOffset.Magnitude})))

	// Back Right Wheel Calculation
	brOffset := VectorMath.Vector2DtoVectorTheta(Types.Vector2D{X: drive.Config.Modules.BackRight.OffsetX, Y: drive.Config.Modules.BackRight.OffsetY})
	br := VectorMath.Vector2DtoVectorTheta(VectorMath.VectorAdd(drive.DriveProps.TranslationalV, VectorMath.VectorThetatoVector2D(Types.VectorTheta{Angle: brOffset.Angle + (math.Pi / 2), Magnitude: drive.DriveProps.RotationalV * brOffset.Magnitude})))

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

func (drive *SwerveDrive) DriveToPose(pose Types.Pose2D) {
	diff := Types.Pose2D{}
	diff.Location = VectorMath.Vector2D{X: pose.Location.X - drive.Pose.Location.X, Y: pose.Location.Y - drive.Pose.Location.Y}
	diff.Angle = pose.Angle - drive.Pose.Angle
	fmt.Println(diff)
}

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

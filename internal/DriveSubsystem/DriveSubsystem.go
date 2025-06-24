package DriveSubsystem

import (
	"fmt"
	"frcrobot/internal/Constants"
	"frcrobot/internal/Utils/MathUtils"
	"frcrobot/internal/Utils/Types"
	"frcrobot/internal/Utils/VectorMath"
	"time"
)

type (
	SwerveDrive struct {
		Pose          Types.Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        Constants.SwerveDriveConfig
		TimeInterval  time.Duration
	}

	SwerveDriveModules struct {
		FrontLeft  Types.VectorTheta
		FrontRight Types.VectorTheta
		BackLeft   Types.VectorTheta
		BackRight  Types.VectorTheta
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
	config := Constants.DriveConstants()
	config.MaxSpeed.RotationalV = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	fmt.Println("Drive Config: ", config)
	swerve_modules := SwerveDriveModules{}

	return &SwerveDrive{Pose: pose, Config: config, SwerveModules: swerve_modules, TimeInterval: interval}
}

// func (drive *SwerveDrive) CalculateSwerveModules(trans VectorMath.Vector2D, rot float64) {
// 	drive.DriveProps.TranslationalA = trans
// 	drive.DriveProps.RotationalA = rot
// 	drive.SwerveModules.FrontLeft = VectorMath.VectorAddNormalized(drive.SwerveModules.FrontLeft, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: ((3 * math.Pi) / 4)}), 1), 1)
// 	drive.SwerveModules.FrontRight = VectorMath.VectorAddNormalized(drive.SwerveModules.FrontRight, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: ((5 * math.Pi) / 4)}), 1), 1)
// 	drive.SwerveModules.BackLeft = VectorMath.VectorAddNormalized(drive.SwerveModules.BackLeft, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: ((7 * math.Pi) / 4)}), 1), 1)
// 	drive.SwerveModules.BackRight = VectorMath.VectorAddNormalized(drive.SwerveModules.BackRight, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: (math.Pi / 4)}), 1), 1)
// 	fmt.Println(drive.SwerveModules)
// }

func (drive *SwerveDrive) CalculateSwerve(trans, rot Types.Vector2D) {
	drive.DriveProps.TranslationalV = Types.Vector2D{X: MathUtils.Clamp(trans.X, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV), Y: MathUtils.Clamp(trans.Y, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)}
	drive.DriveProps.RotationalV = MathUtils.Clamp(rot.X, drive.Config.MaxSpeed.TranslationalV, -drive.Config.MaxSpeed.TranslationalV)

	flOffset := VectorMath.Vector2DtoVectorTheta(Types.Vector2D{X: drive.Config.Modules.FrontLeft.OffsetX, Y: drive.Config.Modules.FrontLeft.OffsetY})
	fl := VectorMath.Vector2DtoVectorTheta(VectorMath.VectorAdd(drive.DriveProps.TranslationalV, VectorMath.VectorThetatoVector2D(Types.VectorTheta{Angle: flOffset.Angle + 90, Magnitude: drive.DriveProps.RotationalV})))
	fr := Types.VectorTheta{}
	bl := Types.VectorTheta{}
	br := Types.VectorTheta{}

	drive.SwerveModules = SwerveDriveModules{
		FrontLeft:  fl,
		FrontRight: fr,
		BackLeft:   bl,
		BackRight:  br,
	}
}

// func (drive *SwerveDrive) DriveWheels(left, right Types.Axis) {

// }

func (drive *SwerveDrive) DriveToPose(pose Types.Pose2D) {
	diff := Types.Pose2D{}
	diff.Location = VectorMath.Vector2D{X: pose.Location.X - drive.Pose.Location.X, Y: pose.Location.Y - drive.Pose.Location.Y}
	diff.Angle = pose.Angle - drive.Pose.Angle
	fmt.Println(diff)
}

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

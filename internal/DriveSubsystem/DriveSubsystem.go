package DriveSubsystem

import (
	"fmt"
	"frcrobot/internal/File"
	"frcrobot/internal/Utils/MathUtils"
	"frcrobot/internal/Utils/Types"
	"frcrobot/internal/Utils/VectorMath"
	"math"
	"time"
)

type (
	SwerveDrive struct {
		Pose          Types.Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        SwerveDriveConfig
		TimeInterval  time.Duration
	}

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
		MaxTransSpeed float64 `json:"maxTrans_speed"`
		MaxTransAccel float64 `json:"maxTrans_accel"`
		MaxRotSpeed   float64 `json:"maxRot_speed"`
		MaxRotAccel   float64 `json:"maxRot_accel"`
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

	SwerveDriveModules struct {
		FrontLeft  VectorMath.Vector2D
		FrontRight VectorMath.Vector2D
		BackLeft   VectorMath.Vector2D
		BackRight  VectorMath.Vector2D
	}

	DriveProperties struct {
		TranslationalV VectorMath.Vector2D
		RotationalV    float64
		TranslationalA VectorMath.Vector2D
		RotationalA    float64
	}
)

func NewSwerveDrive(constants string, interval time.Duration) *SwerveDrive {
	pose := Types.Pose2D{Location: VectorMath.Vector2D{X: 0, Y: 0}, Angle: 0}
	config := SwerveDriveConfig{}
	File.ReadJSON(constants, &config)
	config.MaxSpeed.RotationalV = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	fmt.Println("Drive Config: ", config)
	swerve_modules := SwerveDriveModules{VectorMath.Vector2D{X: 0, Y: 0}, VectorMath.Vector2D{X: 0, Y: 0}, VectorMath.Vector2D{X: 0, Y: 0}, VectorMath.Vector2D{X: 0, Y: 0}}

	return &SwerveDrive{Pose: pose, Config: config, SwerveModules: swerve_modules, TimeInterval: interval}
}

func (drive *SwerveDrive) CalculateSwerveModules(trans VectorMath.Vector2D, rot float64) {
	drive.DriveProps.TranslationalA = trans
	drive.DriveProps.RotationalA = rot
	drive.SwerveModules.FrontLeft = VectorMath.VectorAddNormalized(drive.SwerveModules.FrontLeft, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: ((3 * math.Pi) / 4)}), 1), 1)
	drive.SwerveModules.FrontRight = VectorMath.VectorAddNormalized(drive.SwerveModules.FrontRight, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: ((5 * math.Pi) / 4)}), 1), 1)
	drive.SwerveModules.BackLeft = VectorMath.VectorAddNormalized(drive.SwerveModules.BackLeft, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: ((7 * math.Pi) / 4)}), 1), 1)
	drive.SwerveModules.BackRight = VectorMath.VectorAddNormalized(drive.SwerveModules.BackRight, VectorMath.VectorAddNormalized(trans, VectorMath.VectorThetatoVector2D(VectorMath.VectorTheta{L: rot, T: (math.Pi / 4)}), 1), 1)
	fmt.Println(drive.SwerveModules)
}

func (drive *SwerveDrive) DriveWheels(left, right Types.Axis) {

}

func (drive *SwerveDrive) DriveToPose(pose Types.Pose2D) {
	diff := Types.Pose2D{}
	diff.Location = VectorMath.Vector2D{X: pose.Location.X - drive.Pose.Location.X, Y: pose.Location.Y - drive.Pose.Location.Y}
	diff.Angle = pose.Angle - drive.Pose.Angle
	fmt.Println(diff)
}

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

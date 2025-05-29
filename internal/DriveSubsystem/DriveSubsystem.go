package DriveSubsystem

import (
	"fmt"
	"frcrobot/internal/File"
	"frcrobot/internal/Utils/MathUtils"
	"frcrobot/internal/Utils/VectorMath"
	"math"
)

type (
	SwerveDrive struct {
		Pose          Pose2D
		DriveProps    DriveProperties
		SwerveModules SwerveDriveModules
		Config        SwerveDriveConfig
	}

	Pose2D struct {
		Location VectorMath.Vector2D
		Angle    float64
	}

	SwerveDriveConfig struct {
		MaxSpeed DriveMaxes `json:"drive_maxes"`
	}

	DriveMaxes struct {
		TranslationalV float64 `json:"translationalV"`
		RotationalV    float64 `json:"rotationalV"`
		TranslationalA float64 `json:"translationalA"`
		RotationalA    float64 `json:"rotationalA"`
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

func NewSwerveDrive(constants string) *SwerveDrive {
	pose := Pose2D{VectorMath.Vector2D{X: 0, Y: 0}, 0}
	config := SwerveDriveConfig{}
	File.ReadJSON(constants, &config)
	config.MaxSpeed.RotationalV = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalV))
	config.MaxSpeed.RotationalA = MathUtils.DegtoRad(float64(config.MaxSpeed.RotationalA))
	fmt.Println("Drive Config: ", config)
	swerve_modules := SwerveDriveModules{VectorMath.Vector2D{X: 0, Y: 0}, VectorMath.Vector2D{X: 0, Y: 0}, VectorMath.Vector2D{X: 0, Y: 0}, VectorMath.Vector2D{X: 0, Y: 0}}

	return &SwerveDrive{Pose: pose, Config: config, SwerveModules: swerve_modules}
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

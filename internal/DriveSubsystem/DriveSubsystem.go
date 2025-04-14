package DriveSubsystem

import (
	"fmt"
	"frcrobot/internal/File"
	"frcrobot/internal/Utils/VectorMath"
)

type (
	SwerveDrive struct {
		Pose          Pose2d
		SwerveModules SwerveDriveModules
		Config        SwerveDriveConfig
	}

	Pose2d struct {
		Location VectorMath.Vector2D
		Angle    float32
	}

	SwerveDriveConfig struct {
		MaxSpeed MaxSpeeds `json:"max_speeds"`
	}

	MaxSpeeds struct {
		Rranslational float32 `json:"translational"`
		Rotational    float32 `json:"rotational"`
	}

	SwerveDriveModules struct {
		FrontLeft  VectorMath.VectorTheta
		FrontRight VectorMath.VectorTheta
		BackLeft   VectorMath.VectorTheta
		BackRight  VectorMath.VectorTheta
	}
)

func NewSwerveDrive(constants string) SwerveDrive {
	pose := Pose2d{VectorMath.Vector2D{X: 0, Y: 0}, 0}
	config := SwerveDriveConfig{}
	File.ReadJSON(constants, &config)
	fmt.Println("Drive Config: ", config)

	return SwerveDrive{Pose: pose, Config: config}
}

func (drive *SwerveDrive) CalculateSwerveModules(trans VectorMath.Vector2D, rot float32) {
	fmt.Println("TransV: ", trans, "\tRotV: ", rot)
}

package drive_types

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/drive/swerve"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type Constants struct {
	PositionalPidP float64
	PositionalPidI float64
	PositionalPidD float64

	Swerve swerve.Config
}

type MotorConfigs struct {
}

type DriveSubsystem struct {
	SwerveDrive *swerve.SwerveDrive
	Pose        mathutils.Pose2D
}

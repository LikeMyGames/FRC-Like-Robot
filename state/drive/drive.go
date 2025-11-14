package drive

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

type (
	DriveBase interface {
		DriveWheels(*DriveBase, any)
		DriveToPose(*DriveBase, mathutils.Pose2D)
	}
)

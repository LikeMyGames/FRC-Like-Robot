package vision_types

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/photonvision"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type VisionSubsystem struct {
	FrontCamera *photonvision.Camera
	BackCamera  *photonvision.Camera
	LeftCamera  *photonvision.Camera
	RightCamera *photonvision.Camera
}

type Constants struct {
	FrontCameraName   string
	FrontCameraOffset mathutils.Transorm3D
}

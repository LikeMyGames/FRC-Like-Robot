package vision_types

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/photonvision"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type VisionSubsystem struct {
	FrontCamera   *photonvision.Camera
	BackCamera    *photonvision.Camera
	LeftCamera    *photonvision.Camera
	RightCamera   *photonvision.Camera
	PoseEstimator *photonvision.PoseEstimator
}

type Constants struct {
	FrontCameraName   string
	FrontCameraOffset mathutils.Transform3D
	BackCameraName    string
	BackCameraOffset  mathutils.Transform3D
	LeftCameraName    string
	LeftCameraOffset  mathutils.Transform3D
	RightCameraName   string
	RightCameraOffset mathutils.Transform3D
}

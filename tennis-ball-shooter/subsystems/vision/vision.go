package vision

import (
	vision_types "tennis-ball-shooter/subsystems/vision/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/photonvision"
)

type VisionSubsystem vision_types.VisionSubsystem

func Initialize() *VisionSubsystem {
	s := new(VisionSubsystem)
	s.FrontCamera = photonvision.NewCamera("front")

	return s
}

func Periodic() {

}

package field

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/conn/driver_station"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

func GetCurrentAllianceHubPose() mathutils.Pose2D {
	return GetHubPoseFromAlliance(driver_station.GetAlliance())
}

func GetHubPoseFromAlliance(alliance driver_station.Alliance) mathutils.Pose2D {
	if alliance.IsRed() {
		return mathutils.Pose2D{}
	}
	return mathutils.Pose2D{}
}

func GetRedHubPose() mathutils.Pose2D {
	return mathutils.Pose2D{}
}

func GetBlueHubPose() mathutils.Pose2D {
	return mathutils.Pose2D{}
}

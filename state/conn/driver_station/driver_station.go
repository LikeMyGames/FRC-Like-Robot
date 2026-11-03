package driver_station

import "sync/atomic"

type (
	Alliance struct {
		isRed bool
	}
)

var (
	RedAlliance  Alliance = Alliance{isRed: true}
	BlueAlliance Alliance = Alliance{isRed: false}
)

var (
	isRedAlliance        atomic.Bool
	driverStationChannel atomic.Int32
)

const (
	TestMatch = iota
	PraticeMatch
	QualificationMatch
	EliminationMath
)

func GetMatchNumber() int {
	return 0
}

func GetReplayNumber() int {
	return 0
}

func GetAlliance() Alliance {
	if isRedAlliance.Load() {
		return RedAlliance
	}
	return BlueAlliance

	// if isRed := utils.ReadChannelNonBlocking(allianceChannel); isRed == nil || *isRed {
	// 	return RedAlliance
	// }
	// return BlueAlliance
}

func SetAlliance(isRed bool) {
	// for len(allianceChannel) > 0 {
	// 	<-allianceChannel
	// }
	isRedAlliance.Store(isRed)
}

func GetTournamentType() int {
	return 0
}

func GetDriverStationPosition() int {
	return int(driverStationChannel.Load())
}

func SetDriverStationPosition(driverStation int) {
	driverStationChannel.Store(int32(driverStation))
}

func (a Alliance) IsRed() bool {
	return a.isRed
}

func (a Alliance) IsBlue() bool {
	return !a.isRed
}

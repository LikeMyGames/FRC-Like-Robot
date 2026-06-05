package driver_station

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
	allianceChannel      chan bool
	driverStationChannel chan int
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
	if <-allianceChannel {
		return RedAlliance
	}
	return BlueAlliance
}

func SetAlliance(isRed bool) {
	for len(allianceChannel) > 0 {
		<-allianceChannel
	}
	allianceChannel <- isRed
}

func GetTournamentType() int {
	return 0
}

func GetDriverStationPosition() int {
	return <-driverStationChannel
}

func SetDriverStationPosition(driverStation int) {
	for len(driverStationChannel) > 0 {
		<-driverStationChannel
	}
	driverStationChannel <- driverStation
}

func (a Alliance) IsRed() bool {
	return a.isRed
}

func (a Alliance) IsBlue() bool {
	return !a.isRed
}

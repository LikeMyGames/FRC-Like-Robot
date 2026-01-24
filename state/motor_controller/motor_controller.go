package motor

import (
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/can"
)

type (
	Motor struct {
		config      Config
		runningMode int

		angle     float64
		prevAngle float64
		angleErr  float64

		velocity     float64
		prevVelocity float64
		velocityErr  float64
		velocityFF   float64
		maxVelocity  float64

		torque     float64
		prevTorque float64
		torqueErr  float64
		torqueFF   float64
		maxTorque  float64
	}

	Config struct {
		canID  int
		regMap map[string]int
	}

	runningMode struct {
		Enabled  int
		Disabled int
		Estopped int
	}
)

var RunningMode runningMode = runningMode{
	Enabled:  0,
	Disabled: 1,
	Estopped: 3,
}

func (m *Motor) Configure(config Config) {
	m.config = config
}

func (m *Motor) LoadRegisterMap(regMap map[string]int) {
	m.config.regMap = regMap
}

func (m *Motor) SetRunningMode(newMode int) {
	m.runningMode = newMode
}

func (m *Motor) Estop() {
	m.runningMode = RunningMode.Estopped
	can.BuildAndSendFrame(m.config.canID, m.config.regMap["Estop"])
}

func (m *Motor) ReadAngle() float64 {
	return m.angle
}

func (m *Motor) SetAngle(angle float64) {
	m.angle = angle
	can.BuildAndSendFrame(m.config.canID, m.config.regMap["Input_Pos"], float64(m.angle))
}

func (m *Motor) ReadVelocity() float64 {
	return m.velocity
}

func (m *Motor) SetVelocity(velocity float64) {
	m.velocity = velocity
	can.BuildAndSendFrame(m.config.canID, m.config.regMap["Input_Vel"], float64(m.velocity))
}

func (m *Motor) ReadTorque() float64 {
	return m.torque
}

func (m *Motor) SetTorque(torque float64) {
	m.torque = torque
	can.BuildAndSendFrame(m.config.canID, m.config.regMap["Input_Torque"], float64(m.torque))
}

func (m *Motor) ReadLimits() (velocity, current float64) {
	can.BuildAndSendFrame()
}

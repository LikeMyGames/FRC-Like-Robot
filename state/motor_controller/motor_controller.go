package motor

type (
	Motor struct {
		config Config

		angle     float64
		prevAngle float64
		angleErr  float64

		velocity     float64
		prevVelocity float64
		velocityErr  float64

		acceleration     float64
		prevAcceleration float64
		accelerationErr  float64
	}

	Config struct {
		canID  int
		regMap map[string]int
	}
)

func (m *Motor) LoadRegisterMap(regMap map[string]int) {
	m.config.regMap = regMap
}

func (m *Motor) ReadAngle() float64 {
	return m.angle
}

func (m *Motor) SetAngle(angle float64) {

}

func (m *Motor) ReadVelocity() float64 {
	return m.angle
}

func (m *Motor) SetVelocity(angle float64) {

}

func (m *Motor) ReadAcceleration() float64 {
	return m.angle
}

func (m *Motor) SetAcceleration(angle float64) {

}

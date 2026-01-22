package motor

type (
	Motor struct {
		canId int
		angle float64
	}
)

func (m *Motor) ReadAngle() float64 {
	return m.angle
}

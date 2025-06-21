package Hardware

// will add cgo import

type (
	CANbus struct {
		CANname string
		BaseID  uint
		Devices []Device
	}

	Device interface {
		GetID() uint
		Read() []byte
		Write([]byte) bool
	}
	MotorController struct {
		Velocity     float64
		Displacement float64
		CandID       uint
	}

	Sensor struct {
		Value float64
	}
)

func NewBus() *CANbus {
	return &CANbus{
		CANname: "",
		BaseID:  0,
		Devices: []Device{},
	}
}

func (c *CANbus) AddDevice(d Device) {
	c.Devices = append(c.Devices, d)
}

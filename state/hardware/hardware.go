package hardware

import "fmt"

type (
	SwerveModule struct {
		DriveAngle        float64
		AzimuthAngle      float64
		DriveMotorCanID   int
		AzimuthMotorCanID int
	}

	CANbus struct {
		SPI_MISO_PIN int
		SPI_MOSI_PIN int
		SPI_SCLK_PIN int
	}

	Device struct {
		id    int64
		value float64
	}
)

// var bus = NewCanBus(19, 21, 23)

func NewCanBus(MOSI, MISO, SCLK int) *CANbus {
	return &CANbus{
		SPI_MISO_PIN: MISO,
		SPI_MOSI_PIN: MOSI,
		SPI_SCLK_PIN: SCLK,
	}
}

func (m *SwerveModule) ReadAzimuthAngle() float64 {

	return 0
}

func NewDevice(id int64) *Device {
	return &Device{id: id}
}

func (d *Device) SetTargetType(target string) {
	fmt.Println(d.id, "target type is", target)
}

func (d *Device) SetTargetValue(target float64) {
	fmt.Println(d.id, "target value is", target)
}

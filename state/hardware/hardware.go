package hardware

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
)

var bus = NewCanBus(19, 21, 23)

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

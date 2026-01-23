//go:build linux

package hardware

import (
	"log"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/pid"
	"github.com/warthog618/go-gpiocdev"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
)

type (
	SwerveModule struct {
		DriveAngle        float64
		AzimuthAngle      float64
		DriveMotorCanID   int
		AzimuthMotorCanID int
	}

	Device struct {
		id     int64
		value  float64
		target float64
	}

	MotorController struct {
		device        Device
		PidController *pid.PIDController
	}
)

var (
	bus    *CANbus               = nil
	config constantTypes.Battery = constantTypes.Battery{}
)

func CheckStatus() bool {
	err := gpiocdev.IsChip("gpiochip0")
	return err != nil
}

func NewCanBus() *CANbus {
	if bus != nil {
		return bus
	}
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	// Use spireg SPI port registry to find the first available SPI bus.
	p, err := spireg.Open("/dev/spidev0.0")
	if err != nil {
		log.Fatal(err)
	}

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		log.Fatal(err)
	}

	bus = &CANbus{
		spiPort:       c,
		spiPortCloser: p,
	}

	return bus
}

func (b *CANbus) Close() {
	b.spiPortCloser.Close()
}

func (m *SwerveModule) ReadAzimuthAngle() float64 {

	return 0
}

// func NewDevice(id int64) *Device {
// 	return &Device{id: id}
// }

// func (d *Device) SetTargetType(target string) {
// 	fmt.Println(d.id, "target type is", target)
// }

// func (d *Device) SetTargetValue(target float64) {
// 	fmt.Println(d.id, "target value is", target)
// }

// func NewMotorController(config constantTypes.MotorController) *MotorController {
// 	return &MotorController{
// 		device:        Device{id: config.Id},
// 		PidController: pid.NewPIDController(config.PID),
// 	}
// }

// func (c *MotorController) SetTarget(val float64) {
// 	c.Write(val)
// 	c.PidController.SetTarget(val)
// }

// func (c *MotorController) GetValue() float64 {
// 	return c.device.value
// }

// func (c *MotorController) GetTarget() float64 {
// 	return c.PidController.GetTarget()
// }

// func (c *MotorController) GetId() int64 {
// 	return int64(c.device.target)
// }

// func (c *MotorController) Run() {
// 	c.PidController.Calculate(c.device.value)
// }

// func (c *MotorController) AtValue() bool {
// 	// replace less than val later
// 	return math.Abs(c.device.target-c.device.value) < 0.2
// }

// idk if this method works
// func (c *MotorController) Write(val float64) {
// 	// // Write 0x10 to the device, and read a byte right after.
// 	// write := []byte{0x10, 0x00}
// 	// buf := new(bytes.Buffer)
// 	// err := binary.Write(buf, binary.LittleEndian, val)
// 	// if err != nil {
// 	// 	fmt.Println("could not convert float to bytes")
// 	// 	return
// 	// }
// 	// write := append(buf.Bytes(), 0x00)
// 	// read := make([]byte, len(write))
// 	// if err := bus.spiPort.Tx(write, read); err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Println(read[1:])
// 	// WriteToCan(uint(c.device.id), []byte{})
// }

func SetBatteryConfig(conf constantTypes.Battery) {
	config = conf
}

func ReadBatteryPercentage() uint {
	return uint((ReadBatteryVoltage() / config.NominalVoltage) * 100)
}

func ReadBatteryVoltage() float64 {
	return 12.0
}

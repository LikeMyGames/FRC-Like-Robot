package can

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/warthog618/go-gpiocdev"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
)

type (
	CanFrame struct {
		canId                 int
		cmd                   int
		id                    int
		data                  [8]byte
		callbackEventListener *event.Listener
	}

	CANbus struct {
		spiPort       spi.Conn
		spiPortCloser spi.PortCloser
	}
)

var (
	lowLine     *gpiocdev.Line
	highLine    *gpiocdev.Line
	writeBuffer *bytes.Buffer = new(bytes.Buffer)
	readBuffer  *bytes.Buffer = new(bytes.Buffer)
	bus         *CANbus       = nil
)

func (b *CANbus) Close() {
	b.spiPortCloser.Close()
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

func setUpGPIO() {
	low, err := gpiocdev.RequestLine("gpiochip0", 29)
	if err != nil {
		panic(err)
	}
	high, err := gpiocdev.RequestLine("gpiochip0", 29)
	if err != nil {
		panic(err)
	}
	lowLine = low
	highLine = high
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildFrame(canId, cmd int, data ...any) *CanFrame {
	frame := &CanFrame{}
	frame.canId = canId
	frame.cmd = cmd
	frame.id = (canId << 5) | cmd
	var dataBin []byte
	var err error
	for _, v := range data {
		dataBin, err = binary.Append(dataBin, binary.BigEndian, v)
		if err != nil {
			panic(err)
		}
	}
	fmt.Print("Can data (in binary)")
	for _, v := range dataBin {
		fmt.Printf("%b", v)
	}
	fmt.Println()
	return frame
}

func BuildFrameWithCallback(canId, cmd int, callbackEventTrigger string, callbackFunc func(event any), data ...any) *CanFrame {
	frame := BuildFrame(canId, cmd, data)
	frame.callbackEventListener = event.Listen(callbackEventTrigger, "", callbackFunc)
	return frame
}

func BuildAndSendFrame(canId, cmd int, data ...any) {
	frame := BuildFrame(canId, cmd, data)
	SendFrame(frame)
}

func BuildAndSendFrameWithCallback(canId, cmd int, callbackEventTrigger string, callbackFunc func(event any), data ...any) {
	frame := BuildFrameWithCallback(canId, cmd, callbackEventTrigger, callbackFunc, data)
	SendFrame(frame)
}

func SendFrame(frame *CanFrame) {

}

// need to figure out i want to do this
// func WriteToCan(id uint, data []byte) error {
// 	if bits.Len(id) > 11 {
// 		return errors.New("id must be less than 11 bits in length (<2047)")
// 	}
// 	write := 1
// 	for range 11 - bits.Len(id) {
// 		write = (write << 1) | 0
// 	}
// 	// adding id bits
// 	write = (write << bits.Len(id)) | int(id)
// 	// adding RTR, IDE, and r0 bits
// 	write <<= 3
// 	// identifies as CAN FD
// 	write = (write << 1) | 1
// 	// r1 bit
// 	write = (write << 1) | 0
// 	// BRS bit
// 	write = (write << 1) | 0
// 	// ESI bit
// 	write = (write << 1) | 0
// 	// DLC bits
// 	write = (write << 4) | 0b1111
// 	// payload bits
// 	// length of 64 bytes (512 bits) with above DLC bit layout
// 	fmt.Println("Setting device with id:", id)
// 	fmt.Printf("%b\n", write)
// 	return nil
// }

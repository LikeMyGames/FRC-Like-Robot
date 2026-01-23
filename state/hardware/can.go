package hardware

import (
	"bytes"
	"reflect"

	"github.com/warthog618/go-gpiocdev"
	"periph.io/x/conn/v3/spi"
)

type (
	CanFrame struct {
		canId int
		cmd   int
		id    int
		data  [8]uint8
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
)

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
func BuildFrame(canId, cmd int, data [8]any) *CanFrame {
	frame := &CanFrame{}
	frame.canId = canId
	frame.cmd = cmd
	frame.id = (canId << 5) | cmd
	for i, v := range data {
		// need to convert data to [8]uint8
		reflect.TypeOf(v).Kind().String()
	}
	return frame
}

func BuildSendFrame(canId, cmd int, data [8]uint8) {
	frame := BuildFrame(canId, cmd, data)
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

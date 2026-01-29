package can

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"periph.io/x/conn/v3/spi"
)

type (
	CanFrame struct {
		canId                 int
		cmd                   int
		id                    int
		data                  [8]byte
		callbackEventListener *event.Listener
	}

	CanBus struct {
		spiPort       spi.Conn
		spiPortCloser spi.PortCloser
		messageBuffer map[int]([]*CanFrame)
		devices       []CanDevice
	}

	CanDevice interface {
		Update()
		GetCanId() int
	}
)

var (
	writeBuffer *bytes.Buffer = new(bytes.Buffer)
	readBuffer  *bytes.Buffer = new(bytes.Buffer)
	bus         *CanBus       = nil
)

func (b *CanBus) Close() {
	b.spiPortCloser.Close()
}

func (b *CanBus) UpdateDevices() {
	for _, v := range b.devices {
		v.Update()
	}
}

func NewCanBus() *CanBus {
	if bus != nil {
		return bus
	}

	hardware.OpenSpi()

	bus = &CanBus{
		spiPort:       nil, // c
		spiPortCloser: nil, // p
	}

	return bus
}

func AddDeviceToBus(device CanDevice) {
	if bus == nil {
		return
	}
	bus.devices = append(bus.devices, device)
}

func GetCanMessageFromBuffer(canId, i int) [8]byte {
	bus := NewCanBus()
	data := bus.messageBuffer[canId][i].data
	bus.messageBuffer[canId] = append(bus.messageBuffer[canId][:i], bus.messageBuffer[canId][i+1:]...)
	return data
}

// func setUpGPIO() {
// 	low, err := gpiocdev.RequestLine("gpiochip0", 29)
// 	if err != nil {
// 		panic(err)
// 	}
// 	high, err := gpiocdev.RequestLine("gpiochip0", 29)
// 	if err != nil {
// 		panic(err)
// 	}
// 	lowLine = low
// 	highLine = high
// }

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildFrame(canId, cmd int, data ...any) *CanFrame {
	frame := &CanFrame{}
	frame.canId = canId
	frame.cmd = cmd
	frame.id = (canId << 5) | cmd
	buf := make([]byte, 8)
	var err error
	for _, v := range data {
		switch v := v.(type) {
		case int:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case int8:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case int16:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case int32:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case int64:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case uint:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case uint8:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case uint16:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case uint32:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case uint64:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case float32:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case float64:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		case bool:
			buf, err = binary.Append(buf, binary.BigEndian, v)
		}
		if err != nil {
			panic(err)
		}
	}
	frame.data = [8]byte(buf[:8])
	fmt.Printf("Can data (in binary): %s\n", frame.data)
	return frame
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
// the event parameter in the callback function is the index of the desired *CanFrame in the CanBus object's message buffer map (map[CanId]([]*CanFrame))
func BuildFrameWithCallback(canId, cmd int, callbackFunc func(event any), data ...any) *CanFrame {
	frame := BuildFrame(canId, cmd, data)
	callbackEventTrigger := fmt.Sprintf("CAN_CALLBACK_%v_%v", canId, cmd)
	frame.callbackEventListener = event.Listen(callbackEventTrigger, "", callbackFunc)
	return frame
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildAndSendFrame(canId, cmd int, data ...any) {
	frame := BuildFrame(canId, cmd, data)
	SendFrame(frame)
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildAndSendFrameWithCallback(canId, cmd int, callbackFunc func(event any), data ...any) {
	frame := BuildFrameWithCallback(canId, cmd, callbackFunc, data)
	SendFrame(frame)
}

func SendFrame(frame *CanFrame) {

}

func (f *CanFrame) GetCanId() int {
	return f.canId
}

func (f *CanFrame) GetCmd() int {
	return f.cmd
}

func (f *CanFrame) GetData() [8]byte {
	return f.data
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

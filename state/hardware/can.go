package hardware

import (
	"bytes"
	"errors"
	"fmt"
	"math/bits"

	"github.com/warthog618/go-gpiocdev"
)

type (
	CanFrame struct {
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

// need to figure out i want to do this
func WriteToCan(id uint, data []byte) error {
	if bits.Len(id) > 11 {
		return errors.New("id must be less than 11 bits in length (<2047)")
	}
	write := 1
	for range 11 - bits.Len(id) {
		write = (write << 1) | 0
	}
	// adding id bits
	write = (write << bits.Len(id)) | int(id)
	// adding RTR, IDE, and r0 bits
	write <<= 3
	// identifies as CAN FD
	write = (write << 1) | 1
	// r1 bit
	write = (write << 1) | 0
	// BRS bit
	write = (write << 1) | 0
	// ESI bit
	write = (write << 1) | 0
	// DLC bits
	write = (write << 4) | 0b1111
	// payload bits
	// length of 64 bytes (512 bits) with above DLC bit layout
	fmt.Println("Setting device with id:", id)
	fmt.Printf("%b\n", write)
	return nil
}

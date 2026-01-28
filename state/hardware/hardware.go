//go:build linux

package hardware

import (
	"fmt"
	"log"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/warthog618/go-gpiocdev"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type ()

var (
	config constantTypes.Battery = constantTypes.Battery{}
)

func OpenSpi() {
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use spireg SPI port registry to find the first available SPI bus.
	p, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	c, err := p.Connect(physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		log.Fatal(err)
	}

	// Write 0x10 to the device, and read a byte right after.
	write := []byte{0x10, 0x00}
	read := make([]byte, len(write))
	if err := c.Tx(write, read); err != nil {
		log.Fatal(err)
	}
	// Use read.
	fmt.Printf("%v\n", read[1:])
}

func Pin() {
	v := 0
	l, err := gpiocdev.RequestLine("gpiochip0", 4, gpiocdev.AsOutput(v))
	if err != nil {
		panic(err)
	}
	for {
		<-time.After(time.Second)
		v ^= 1
		l.SetValue(v)
	}
}

func CheckStatus() bool {
	err := gpiocdev.IsChip("gpiochip0")
	return err != nil
}

func SetBatteryConfig(conf constantTypes.Battery) {
	config = conf
}

func ReadBatteryPercentage() uint {
	return uint((ReadBatteryVoltage() / config.NominalVoltage) * 100)
}

func ReadBatteryVoltage() float64 {
	return 12.0
}

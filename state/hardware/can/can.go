package can

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"gitlab.com/Vitmark/can-go"
	"gitlab.com/Vitmark/can-go/slcan"
	"gitlab.com/Vitmark/canable-go/canable"
	"periph.io/x/conn/v3/spi"
)

type (
	CanFrame struct {
		// canId                 int
		// cmd                   int
		id                    int
		deviceType            int
		manufacturer          int
		apiClass              int
		apiIndex              int
		data                  [8]byte
		callbackEventListener *event.Listener
		timestamp             time.Time
	}

	CanBus struct {
		spiPort spi.Conn
		// messageBuffer map[int](chan *CanFrame)
		devices        []CanDevice
		buffer         map[int][]chan *CanFrame
		messageBuffers []*MessageBuffer
		connection     canable.Connection
		bitrate        slcan.DataBitrate
	}

	CanDevice interface {
		Status() bool
		Update()
		GetCanId() int
	}

	MessageBuffer struct {
		registeredIDs    []int
		channelLocations map[int]int
		buffers          map[int][]*CanFrame
		callbacks        map[int][]func(data []byte)
	}
)

var (
	// writeBuffer *bytes.Buffer = new(bytes.Buffer)
	// readBuffer  *bytes.Buffer = new(bytes.Buffer)
	bus *CanBus = nil
)

var (
	NOT_COMPATIBLE_FRAME_ERROR error = fmt.Errorf("Recieved CAN frame is not compatible with CAN system. CAN frame id is not extended and CAN frame is not a data frame.")
)

func NewCanBus() *CanBus {
	if bus != nil {
		return bus
	}

	bus = new(CanBus)

	bus.buffer = make(map[int][]chan *CanFrame)
	bus.messageBuffers = make([]*MessageBuffer, 0)

	var portName string
	var enableDebug bool
	flag.StringVar(&portName, "c", "", "Canable2 COM port (leave empty for automatic port search)")
	flag.BoolVar(&enableDebug, "d", false, "Set logging level to Debug")
	flag.Parse()

	if enableDebug {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	mainCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	slog.Info("Attempting to connect to Canable2 device...")
	connection, err := canable.ListAndConnect(mainCtx, portName, canable.DefaultConnectionConfig())
	if err != nil {
		go func() {
			intChan := make(chan os.Signal, 1)
			signal.Notify(intChan, os.Interrupt, syscall.SIGTERM)
			c := time.NewTicker(time.Second)
			for range c.C {

				select {
				case <-intChan:
					c.Stop()
					os.Exit(0)
					return
				default:
					slog.Info("Attempting to connect to Canable2 device...")
					connection, err := canable.ListAndConnect(mainCtx, portName, canable.DefaultConnectionConfig())
					if err == nil {
						slog.Info("Successfully connected to Canable2", "port", connection.PortName())
						connection.StartCAN(slcan.Bitrate(slcan.B5M))
						bus.connection = connection
						bus.bitrate = slcan.B5M
						c.Stop()

						return
					}
					slog.Info("Failed to connect to Canable2 device")
				}

			}
		}()
	} else {
		slog.Info("Successfully connected to Canable2", "port", connection.PortName())

		connection.StartCAN(slcan.Bitrate(slcan.B5M))
		bus.connection = connection
		bus.bitrate = slcan.B5M
	}

	return bus
}

// func framePrinter(ctx context.Context, portName string) error {

// 	// mainLoop:
// 	// 	for {
// 	// 		select {
// 	// 		case <-ctx.Done():
// 	// 			break mainLoop
// 	// 		case rxFrame, ok := <-connection.RxChannel():
// 	// 			if !ok {
// 	// 				break mainLoop
// 	// 			}
// 	// 			fmt.Printf("%s\n", rxFrame)
// 	// 		}
// 	// 	}

// 	// 	missedFrames := connection.RxMissedCount()
// 	// 	err = connection.Close()
// 	// 	slog.Info("exiting", "error", err, "missedFrames", missedFrames)
// 	return nil
// }

// func listenOnHardwareSpiBus() {
// 	canTicker := time.NewTicker(time.Millisecond)

// 	for range canTicker.C {

// 	}
// }

func (b *CanBus) Close() {
	b.connection.Close()
}

func (b *CanBus) UpdateDevices() {
	for _, v := range b.devices {
		v.Update()
	}
}

func (b *CanBus) DistributeMessagesToBuffers() {
	for rxFrame := range b.connection.RxChannel() {
		frame := new(CanFrame)

		canFrame := rxFrame.CanFrame
		canFrame.Validate()

		frame.timestamp = rxFrame.Timestamp
		frame.data = rxFrame.CanFrame.Data

		frame.id = int(canFrame.ID.GetId())
		var err error = nil
		frame.deviceType, frame.manufacturer, frame.apiClass, frame.apiIndex, err = processCanFrameId(canFrame.ID)
		if err != nil {
			panic(err)
		}

		slog.Info("Received Can Frame", "Device Type", frame.deviceType, "Manufacturer", frame.manufacturer, "API Class", frame.apiClass, "API Index", frame.apiIndex, "Data", frame.data)
		b.pushFrameToRegisterBuffer(frame)
	}

	// for {
	// 	select {
	// 	case rxFrame, ok := <-b.connection.RxChannel():

	// 	}
	// }
}

func (b *CanBus) CheckStatuses() bool {
	// fmt.Println("starting CAN device status check")
	bad := make([]int, 0)
	for _, v := range b.devices {
		// fmt.Printf("Checking status of CAN device %d\n", v.GetCanId())
		if !v.Status() {
			// fmt.Println("adding CAN device to bad status list")
			bad = append(bad, v.GetCanId())
		}
	}

	for _, i := range bad {
		fmt.Printf("Could not get good status from CAN device %d\n", i)
		// fmt.Println("Check connection to CAN chain and power status of device for troubleshooting")
	}

	// temporary
	return len(bad) == 0
}

func (bus *CanBus) pushFrameToRegisterBuffer(frame *CanFrame) {
	registeredIds := slices.Collect(maps.Keys(bus.buffer))
	if slices.Contains(registeredIds, frame.id) {
		for i := range bus.buffer[frame.id] {
			bus.buffer[frame.id][i] <- frame
		}
	}
	for _, buf := range bus.messageBuffers {
		buf.pullFromMainBufferToSubBuffer()
	}
}

// func ReceiveFrame(frame *CanFrame) {
// 	NewCanBus().messageBuffer[frame.id] <- frame
// }

func AddDeviceToBus(device CanDevice) {
	if bus == nil {
		return
	}
	bus.devices = append(bus.devices, device)
	fmt.Printf("Added new device with id %d to CanBus device list\n", device.GetCanId())
}

func NewMessageBuffer() *MessageBuffer {
	buf := new(MessageBuffer)
	buf.buffers = make(map[int][]*CanFrame)
	buf.callbacks = make(map[int][]func(data []byte))
	buf.channelLocations = make(map[int]int)
	buf.registeredIDs = make([]int, 0)

	bus.messageBuffers = append(bus.messageBuffers, buf)
	return buf
}

func (buf *MessageBuffer) RegisterId(id int) {
	buf.registeredIDs = append(buf.registeredIDs, id)
}

func (buf *MessageBuffer) RegisterCallbackOnID(id int, callback func(data []byte)) {
	if buf.callbacks[id] == nil {
		buf.callbacks[id] = make([]func([]byte), 0)
	}

	buf.callbacks[id] = append(buf.callbacks[id], callback)
}

func (buf *MessageBuffer) callIdCallbacks() {
	for id, frames := range buf.buffers {
		var latestFrame *CanFrame = nil
		for _, frame := range frames {
			if latestFrame.timestamp.Before(frame.timestamp) {
				latestFrame = frame
			}
		}
		for _, v := range buf.callbacks[id] {
			v(latestFrame.data[:])
		}
	}
}

func (buf *MessageBuffer) RegisterCallbackWithID(id int, callback func(data []byte)) {
	buf.RegisterId(id)
	buf.RegisterCallbackOnID(id, callback)
}

func (buf *MessageBuffer) pullFromMainBufferToSubBuffer() {
	for id, i := range buf.channelLocations {
		for len(bus.buffer[id][i]) > 0 {
			buf.buffers[id] = append(buf.buffers[id], <-bus.buffer[id][i])
		}
	}

	buf.callIdCallbacks()
}

func AddBufferToBus(buffer *MessageBuffer) {
	for id := range buffer.registeredIDs {
		if !slices.Contains(slices.Collect(maps.Keys(bus.buffer)), id) {
			bus.buffer[id] = make([]chan *CanFrame, 0)
		}
		bus.buffer[id] = append(bus.buffer[id], make(chan *CanFrame))
		buffer.channelLocations[id] = len(bus.buffer[id]) - 1
	}
}

// func GetCanMessageFromBuffer(id, cmd int) *[8]byte {
// 	bus := NewCanBus()
// 	frame := (utils.ReadChannelNonBlocking(bus.messageBuffer[id]))
// 	if frame == nil {
// 		return nil
// 	}
// 	return &(*frame).data
// }

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildFrame(id int, data ...any) *CanFrame {
	frame := &CanFrame{}
	frame.id = id
	// frame.id = (canId << 5) | cmd
	buf := make([]byte, 8)
	// buf, err = binary.Append(buf, binary.BigEndian, data)
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
	// fmt.Printf("Can data (in binary): %s\n", frame.data)
	return frame
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
// the event parameter in the callback function is the index of the desired *CanFrame in the CanBus object's message buffer map (map[CanId]([]*CanFrame))
func BuildFrameWithCallback(canId int, callbackFunc func(event any), data ...any) *CanFrame {
	frame := BuildFrame(canId, data)
	callbackEventTrigger := fmt.Sprintf("CAN_CALLBACK_%v", canId)
	frame.callbackEventListener = event.Listen(callbackEventTrigger, callbackFunc)
	return frame
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildAndSendFrame(canId int, data ...any) {
	frame := BuildFrame(canId, data)
	SendFrame(frame)
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildAndSendFrameWithCallback(canId, cmd int, callbackFunc func(event any), data ...any) {
	frame := BuildFrameWithCallback(canId, callbackFunc, data)
	SendFrame(frame)
}

func BuildId(deviceType, manufacturer, class, index, deviceId int) int {
	deviceType &= 0b11111
	manufacturer &= 0b11111111
	class &= 0b111111
	index &= 0b1111
	deviceId &= 0b11111
	return (((deviceType<<8|manufacturer)<<6|class)<<4|index)<<5 | deviceId
}

func SendFrame(frame *CanFrame) {

}

func SendSyncMessage(matchTime, matchNumber, replayNumber int64, redAlliance, enabled, autonomous, test, systemWatchdog bool, tournamentType int64) {
	now := time.Now()
	hour := now.Hour()
	hour &= (1 << 5) - 1

	minute := now.Minute()
	minute &= (1 << 6) - 1

	second := now.Second()
	second &= (1 << 6) - 1

	day := now.Day()
	day &= (1 << 5) - 1

	month := int(time.Month(now.Month()))
	month &= (1 << 4) - 1

	year := now.Year()
	year &= (1 << 6) - 1

	date := int64(((((year<<4|month)<<5|day)<<6|second)<<6|minute)<<6 | hour)

	matchTime &= (1 << 8) - 1
	matchNumber &= (1 << 10) - 1
	replayNumber &= (1 << 6) - 1
	var _redAlliance int64 = 0
	if redAlliance {
		_redAlliance = 1
	}
	var _enabled int64 = 0
	if enabled {
		_enabled = 1
	}
	var _autonomous int64 = 0
	if autonomous {
		_autonomous = 1
	}
	var _test int64 = 0
	if test {
		_test = 1
	}
	var _watchdog int64 = 0
	if systemWatchdog {
		_watchdog = 1
	}
	tournamentType &= (1 << 3) - 1

	data := ((((((((matchTime<<10|matchNumber)<<6|replayNumber)<<1|_redAlliance)<<1|_enabled)<<1|_autonomous)<<1|_test)<<1|_watchdog)<<3|tournamentType)<<32 | date

	buf := make([]byte, 8)
	binary.Encode(buf, binary.BigEndian, data)

	frame := &CanFrame{}
	frame.id = 0x01011840
	frame.data = [8]byte(buf)

	SendFrame(frame)
}

func (f *CanFrame) GetId() int {
	return f.id
}

func (f *CanFrame) GetData() [8]byte {
	return f.data
}

func (f *CanFrame) GetTimestamp() time.Time {
	return f.timestamp
}

func (f *CanFrame) GetDeviceType() int {
	return f.deviceType
}

func (f *CanFrame) GetManufacturer() int {
	return f.manufacturer
}

func (f *CanFrame) GetApiClass() int {
	return f.apiClass
}

func (f *CanFrame) GetApiIndex() int {
	return f.apiIndex
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

func processCanFrameId(id can.Id) (deviceType, manufacturer, apiClass, apiIndex int, err error) {
	err = nil
	if id.IsData() && id.IsExtended() {

	}
	err = NOT_COMPATIBLE_FRAME_ERROR
	return
}

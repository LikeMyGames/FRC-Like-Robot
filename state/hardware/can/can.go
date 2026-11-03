package can

import (
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"os/signal"
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
		buffer         map[int][]*CanFrame
		messageBuffers []*MessageBuffer
		connection     canable.Connection
		bitrate        slcan.DataBitrate
		// transmitBuffer []*CanFrame
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
	bus         *CanBus = nil
	portName    string  = ""
	enableDebug bool    = false
)

var (
	NOT_COMPATIBLE_FRAME_ERROR error = fmt.Errorf("Recieved CAN frame is not compatible with CAN system. CAN frame id is not extended and CAN frame is not a data frame.")
)

func init() {
	flag.StringVar(&portName, "c", "", "Canable2 COM port (leave empty for automatic port search)")
	flag.BoolVar(&enableDebug, "d", false, "Set logging level to Debug")
}

func NewCanBus() *CanBus {
	if bus != nil {
		return bus
	}

	bus = new(CanBus)
	bus.messageBuffers = make([]*MessageBuffer, 0)

	flag.Parse()

	if enableDebug {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	// slog.Info("Attempting to connect to Canable2 device...")
	// connection, err := canable.ListAndConnect(context.Background(), portName, canable.DefaultConnectionConfig())

	func() {
		ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
		c := time.NewTicker(time.Second)
		for range c.C {

			select {
			case <-ctx.Done():
				slog.Info("Stopping Canable2 device discovery loop")
				c.Stop()
				return
			default:
				slog.Info("Attempting to connect to Canable2 device...")
				connection, err := canable.ListAndConnect(ctx, portName, canable.DefaultConnectionConfig())
				if err == nil {
					slog.Info("Successfully connected to Canable2", "port", connection.PortName())
					connection.StartCAN(slcan.B1M)
					// connection.StartCAN(slcan.B250k)
					bus.connection = connection
					// bus.bitrate = slcan.DataBitrate(slcan.B250k)
					bus.bitrate = slcan.DataBitrate(slcan.B1M)
					c.Stop()

					return
				}
				slog.Info("Failed to connect to Canable2 device")
			}

		}
	}()

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
	if b.connection != nil {
		b.connection.Close()
	}
}

func (b *CanBus) UpdateDevices() {
	b.buffer = make(map[int][]*CanFrame)
	b.DistributeMessagesToBuffers()

	for _, buf := range b.messageBuffers {
		for _, id := range buf.registeredIDs {
			buf.buffers[id] = append(buf.buffers[id], b.buffer[id]...)
		}

		for id, callbacks := range buf.callbacks {
			for _, can_frame := range buf.buffers[id] {
				for _, v := range callbacks {
					v(can_frame.data[:])
				}
			}
		}
	}

	for _, v := range b.devices {
		v.Update()
	}

	slog.Info("Dropping any unused CAN frames")
}

func (b *CanBus) DistributeMessagesToBuffers() {

	for {
		select {
		case rxFrame := <-b.connection.RxChannel():
			frame := new(CanFrame)

			canFrame := rxFrame.CanFrame
			canFrame.Validate()

			frame.timestamp = rxFrame.Timestamp
			frame.data = rxFrame.CanFrame.Data

			frame.id = int(canFrame.ID.GetId())

			slog.Info("Recieved CAN Frame", "id", frame.id)

			b.buffer[frame.id] = append(b.buffer[frame.id], frame)
		default:
			return
		}
	}
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

	buf.RegisterId(id)
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildFrame(id int, data ...any) *CanFrame {
	frame := &CanFrame{}
	frame.id = id

	buf := make([]byte, 0)
	// buf, err = binary.Append(buf, binary.BigEndian, data)
	for _, v := range data {
		switch v := v.(type) {
		case int8:
			data := uint8(v)
			if v < 0 {
				data |= 0b10000000
			}
			buf = append(buf, data)
		case int16:
			data := uint16(v)
			if v < 0 {
				data |= 0b1000000000000000
			}
			buf = append(buf, 0, 0)
			binary.LittleEndian.PutUint16(buf[len(buf)-2:], data)
		case int32:
			data := uint32(v)
			if v < 0 {
				data |= 0b100000000000000000000000
			}
			buf = append(buf, 0, 0, 0, 0)
			binary.LittleEndian.PutUint32(buf[len(buf)-4:], data)
		case int64:
			data := uint64(v)
			if v < 0 {
				data |= 0b10000000000000000000000000000000
			}
			buf = append(buf, 0, 0, 0, 0, 0, 0, 0, 0)
			binary.LittleEndian.PutUint64(buf[len(buf)-8:], data)
		case uint8:
			buf = append(buf, v)
		case uint16:
			buf = append(buf, 0, 0)
			binary.LittleEndian.PutUint16(buf[len(buf)-2:], v)
		case uint32:
			buf = append(buf, 0, 0, 0, 0)
			binary.LittleEndian.PutUint32(buf[len(buf)-4:], v)
		case uint64:
			binary.LittleEndian.PutUint64(buf[len(buf)-8:], v)
		case float32:
			data := math.Float32bits(v)
			buf = append(buf, 0, 0, 0, 0)
			binary.LittleEndian.PutUint32(buf[len(buf)-4:], data)
		case float64:
			data := math.Float64bits(v)
			buf = append(buf, 0, 0, 0, 0, 0, 0, 0, 0)
			binary.LittleEndian.PutUint64(buf[len(buf)-8:], data)
		}
	}
	if len(buf) > 8 {
		frame.data = [8]byte(buf[:8])
	} else {
		len := len(buf)
		for i := 0; i < 8-len; i++ {
			buf = append(buf, 0)
		}
		// fmt.Println(buf)
		frame.data = [8]byte(buf)
	}
	// fmt.Printf("Can data (in binary): %s\n", frame.data)
	return frame
}

func BuildFrameFromBits(id int, data [8]byte) *CanFrame {
	frame := new(CanFrame)
	frame.id = id
	frame.timestamp = time.Now()
	frame.data = data

	return frame
}

func BuildAndSendFrameFromBits(id int, data [8]byte) {
	SendFrame(BuildFrameFromBits(id, data))
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
// the event parameter in the callback function is the index of the desired *CanFrame in the CanBus object's message buffer map (map[CanId]([]*CanFrame))
func BuildFrameWithCallback(canId int, callbackFunc func(event any), data ...any) *CanFrame {
	frame := BuildFrame(canId, data...)
	callbackEventTrigger := fmt.Sprintf("CAN_CALLBACK_%v", canId)
	frame.callbackEventListener = event.Listen(callbackEventTrigger, callbackFunc)
	return frame
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildAndSendFrame(canId int, data ...any) {
	frame := BuildFrame(canId, data...)
	SendFrame(frame)
}

// canId is a 6 bit max integer
// cmd is a 5 bit max integer
// The max sendable data size is 8 bytes (64 bits)
func BuildAndSendFrameWithCallback(canId, cmd int, callbackFunc func(event any), data ...any) {
	frame := BuildFrameWithCallback(canId, callbackFunc, data...)
	SendFrame(frame)
}

func BuildId(deviceType, manufacturer, class, index, deviceId int) int {
	deviceType &= 0b11111
	manufacturer &= 0b11111111
	class &= 0b111111
	index &= 0b1111
	deviceId &= 0b11111
	return (((deviceType<<8|manufacturer)<<6|class)<<4|index)<<6 | deviceId
}

func SendFrame(frame *CanFrame) error {
	if bus == nil {
		return errors.New("No CAN bus created yet. You must create a CAN bus singleton object using can.NewCanBus()")
	}
	if bus.connection == nil {
		// bus.transmitBuffer = append(bus.transmitBuffer, frame)
		return canable.ErrConnClosed
	}
	slog.Info("Sending CAN frame", "id", frame.id, "data", frame.data)
	return bus.connection.SendCanFrame(can.Frame{ID: can.NewId(uint32(frame.id), true, false), Data: frame.data, DLC: 8})
}

func SendSyncMessage(matchTime, matchNumber, replayNumber uint64, redAlliance, enabled, autonomous, test, systemWatchdog bool, tournamentType uint64) {
	now := time.Now()
	hour := now.Hour()
	hour &= (1 << 5) - 1
	// slog.Info("heartbeat section added", "section", "hour")

	minute := now.Minute()
	minute &= (1 << 6) - 1
	// slog.Info("heartbeat section added", "section", "minute")

	second := now.Second()
	second &= (1 << 6) - 1
	// slog.Info("heartbeat section added", "section", "second")

	day := now.Day()
	day &= (1 << 5) - 1
	// slog.Info("heartbeat section added", "section", "day")

	month := int(time.Month(now.Month()))
	month &= (1 << 4) - 1
	// slog.Info("heartbeat section added", "section", "month")

	year := now.Year()
	year &= (1 << 6) - 1
	// slog.Info("heartbeat section added", "section", "year")

	date := uint64(((((year<<4|month)<<5|day)<<6|second)<<6|minute)<<6 | hour)
	// slog.Info("heartbeat section added", "section", "date")

	matchTime &= (1 << 8) - 1
	// slog.Info("heartbeat section added", "section", "natchTime")
	matchNumber &= (1 << 10) - 1
	// slog.Info("heartbeat section added", "section", "matchNumber")
	replayNumber &= (1 << 6) - 1
	// slog.Info("heartbeat section added", "section", "replayNumber")
	var _redAlliance uint64 = 0
	if redAlliance {
		_redAlliance = 1
	}
	// slog.Info("heartbeat section added", "section", "redAlliance")
	var _enabled uint64 = 0
	if enabled {
		_enabled = 1
	}
	// slog.Info("heartbeat section added", "section", "enabled")
	var _autonomous uint64 = 0
	if autonomous {
		_autonomous = 1
	}
	// slog.Info("heartbeat section added", "section", "autonomous")
	var _test uint64 = 0
	if test {
		_test = 1
	}
	// slog.Info("heartbeat section added", "section", "test")
	var _watchdog uint64 = 0
	if systemWatchdog {
		_watchdog = 1
	}
	// slog.Info("heartbeat section added", "section", "watchdog")
	tournamentType &= (1 << 3) - 1
	// slog.Info("heartbeat section added", "section", "tournamentType")

	data := ((((((((matchTime<<10|matchNumber)<<6|replayNumber)<<1|_redAlliance)<<1|_enabled)<<1|_autonomous)<<1|_test)<<1|_watchdog)<<3|tournamentType)<<32 | date
	// slog.Info("heartbeat message combined into int64")

	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf[:], data)
	// binary.Encode(buf, binary., data)

	// fmt.Printf("heartbeat message data=%b\n", data)
	// for _, v := range buf {
	// 	fmt.Printf("%08b", v)
	// }
	// fmt.Println()

	frame := new(CanFrame)
	frame.timestamp = time.Now()
	frame.id = 0x01011840
	frame.data = [8]byte(buf)

	// slog.Info("heartbeat message frame built")

	// slog.Info("Sending heartbeat CAN message", "enabled", enabled)
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

func processCanFrameId(id can.Id) (deviceType, manufacturer, apiClass, apiIndex int, err error) {
	err = nil
	if id.IsData() && id.IsExtended() {

	}
	err = NOT_COMPATIBLE_FRAME_ERROR
	return
}

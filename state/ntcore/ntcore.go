package ntcore

const (
	PORT_NT4 = 1736
	PORT_NT3 = 1735

	PENDING   = "pending"
	LISTENING = "listening"
	READY     = "ready"
)

var (
	ProtocolVersion = [2]byte{0x04, 0x02}
)

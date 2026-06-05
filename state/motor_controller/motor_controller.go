package motor

import (
	"log/slog"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/can"
	"github.com/LikeMyGames/FRC-Like-Robot/state/pid"
)

type (
	Motor struct {
		canid       int
		config      Config
		runningMode int

		targetAngle        float64
		angle              float64
		prevAngle          float64
		angleErr           float64
		acceptableAngleErr float64
		atAngleTarget      bool

		targetVelocity        float64
		velocity              float64
		prevVelocity          float64
		velocityErr           float64
		velocityFF            float64
		maxVelocity           float64
		acceptableVelocityErr float64
		atVelocityTarget      bool

		targetTorque         float64
		torque               float64
		prevTorque           float64
		torqueErr            float64
		torqueFF             float64
		maxTorque            float64
		acceptableTorqueErr  float64
		atAccelerationTarget bool

		maxCurrent float64

		canMessageBuffer *can.MessageBuffer
	}

	Config struct {
		deviceType                 int
		manufacturer               int
		class                      int
		index                      int
		deviceId                   int
		pid_map                    map[int]pid.Constants
		Slot0, Slot1, Slot2, Slot3 Slot

		// Position_PID_0 pid.Constants
		// Velocity_PID_0 pid.Constants
		// Torque_PID_0   pid.Constants
		// Position_PID_1 pid.Constants
		// Velocity_PID_1 pid.Constants
		// Torque_PID_1   pid.Constants
		// Position_PID_2 pid.Constants
		// Velocity_PID_2 pid.Constants
		// Torque_PID_2   pid.Constants
		// Position_PID_3 pid.Constants
		// Velocity_PID_3 pid.Constants
		// Torque_PID_3   pid.Constants
		// PID_P_0        float64
		// PID_I_0        float64
		// PID_D_0        float64
		// PID_I_ZONE_0   float64
		// PID_FF_0       float64
		// PID_P_1        float64
		// PID_I_1        float64
		// PID_D_1        float64
		// PID_I_ZONE_1   float64
		// PID_FF_1       float64
		// PID_P_2        float64
		// PID_I_2        float64
		// PID_D_2        float64
		// PID_I_ZONE_2   float64
		// PID_FF_2       float64
		// PID_P_3        float64
		// PID_I_3        float64
		// PID_D_3        float64
		// PID_I_ZONE_3   float64
		// PID_FF_3       float64
		// PID_CosFF              float64

		PositionConversion     float64
		VelocityConversion     float64
		AccelerationConversion float64
		// regMap                 map[string]register
		IsSecondaryMotor  bool
		secondarySumValue int
	}

	runningMode struct {
		Enabled  int
		Disabled int
		Estopped int
	}

	// register struct {
	// 	class            int
	// 	index            int
	// 	sendStructure    []reflect.Kind
	// 	receiveStructure []reflect.Kind
	// }

	Slot struct {
		Position pid.Constants
		Velocity pid.Constants
		Torque   pid.Constants
	}

	ControlSlot int

	ControlType int
)

var RunningMode runningMode = runningMode{
	Enabled:  0,
	Disabled: 1,
	Estopped: 3,
}

const (
	NOT_RESPONDING int = iota
	NO_ERROR
	CAN_TIMEOUT
)

const (
	Slot_0 ControlSlot = iota
	Slot_1
	Slot_2
	Slot_3
)

const (
	PositionControl ControlType = iota
	VelocityControl
	TorqueControl
	VoltageControl
	PercentControl
)

const (
	slot0PositionPidClass int = 1
	slot0VelocityPidClass int = 2
	slot0TorquePidClass   int = 3

	slot1PositionPidClass int = 4
	slot1VelocityPidClass int = 5
	slot1TorquePidClass   int = 6

	slot2PositionPidClass int = 7
	slot2VelocityPidClass int = 8
	slot2TorquePidClass   int = 9

	slot3PositionPidClass int = 10
	slot3VelocityPidClass int = 11
	slot3TorquePidClass   int = 12

	setSetpointClass int = 14
)

const (
	deviceType   = 2
	manufacturer = 255
)

// var motorRegisterMap map[string]register = map[string]register{
// 	"Version": {
// 		class:            0x00,
// 		index:            0x00,
// 		sendStructure:    []reflect.Kind{},
// 		receiveStructure: []reflect.Kind{reflect.Int8, reflect.Int8, reflect.Int8},
// 	},
// 	"Set_P_Slot": {
// 		cmd:              0x01,
// 		sendStructure:    []reflect.Kind{},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Input_Pos": {
// 		cmd:              0x02,
// 		sendStructure:    []reflect.Kind{reflect.Float32, reflect.Int16, reflect.Int16},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Input_Vel": {
// 		cmd:              0x03,
// 		sendStructure:    []reflect.Kind{reflect.Float32, reflect.Int16},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Input_Torque": {
// 		cmd:              0x04,
// 		sendStructure:    []reflect.Kind{reflect.Float32},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Limits": {
// 		cmd:              0x05,
// 		sendStructure:    []reflect.Kind{reflect.Float32, reflect.Float32},
// 		receiveStructure: []reflect.Kind{reflect.Float32, reflect.Float32},
// 	},
// 	"Reboot": {
// 		cmd:              0x06,
// 		sendStructure:    []reflect.Kind{},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Get_Bus_Voltage_Current": {
// 		cmd:              0x07,
// 		sendStructure:    []reflect.Kind{},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Clear_Errors": {
// 		cmd:              0x08,
// 		sendStructure:    []reflect.Kind{},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Encoders": {
// 		cmd:              0x09,
// 		sendStructure:    []reflect.Kind{},
// 		receiveStructure: []reflect.Kind{},
// 	},
// 	"Running_Mode": {
// 		cmd:              0x0a,
// 		sendStructure:    []reflect.Kind{reflect.Int8},
// 		receiveStructure: []reflect.Kind{},
// 	},
// }

// func New(config Config) *Motor {
// 	return motor
// }

func New(CanId int, config Config) *Motor {
	motor := new(Motor)
	// motor.LoadRegisterMap(motorRegisterMap)
	motor.config.deviceId = CanId
	motor.canid = can.BuildId(deviceType, manufacturer, 0, 0, CanId)

	motor.canMessageBuffer = can.NewMessageBuffer()
	motor.registerIds()

	can.AddBufferToBus(motor.canMessageBuffer)
	slog.Info("Created new motor controller", "CanId", CanId)
	// fmt.Println("Created new motor with id: ", CanId)

	motor.Configure(config)

	return motor
}

func (m *Motor) registerIds() {
	// internal encoder
	m.canMessageBuffer.RegisterId(m.buildCanId(13, 0)) // position read
	m.canMessageBuffer.RegisterId(m.buildCanId(13, 1)) // velocity read
	m.canMessageBuffer.RegisterId(m.buildCanId(13, 2)) // acceleration read

	// external encoder
	m.canMessageBuffer.RegisterId(m.buildCanId(13, 3)) // position read
	m.canMessageBuffer.RegisterId(m.buildCanId(13, 4)) // velocity read
	m.canMessageBuffer.RegisterId(m.buildCanId(13, 5)) // acceleration read
}

func (m *Motor) GetCanId() int {
	return m.config.deviceId
}

func (m *Motor) Configure(config Config) {
	m.config = config
	// m.canid = can.BuildId(deviceType, manufacturer, 0, 0, config.deviceId)

	// slot 0
	for i, value := range config.Slot0.Position.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot0PositionPidClass, i), value)
	}

	for i, value := range config.Slot0.Velocity.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot0VelocityPidClass, i), value)
	}

	for i, value := range config.Slot0.Torque.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot0TorquePidClass, i), value)
	}

	// slot 1
	for i, value := range config.Slot1.Position.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot1PositionPidClass, i), value)
	}

	for i, value := range config.Slot1.Velocity.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot1VelocityPidClass, i), value)
	}

	for i, value := range config.Slot1.Torque.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot1TorquePidClass, i), value)
	}

	// slot 2
	for i, value := range config.Slot2.Position.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot2PositionPidClass, i), value)
	}

	for i, value := range config.Slot2.Velocity.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot2VelocityPidClass, i), value)
	}

	for i, value := range config.Slot2.Torque.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot2TorquePidClass, i), value)
	}

	// slot 3
	for i, value := range config.Slot3.Position.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot3PositionPidClass, i), value)
	}

	for i, value := range config.Slot3.Velocity.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot3VelocityPidClass, i), value)
	}

	for i, value := range config.Slot3.Torque.GetAsArray() {
		can.BuildAndSendFrame(m.buildCanId(slot3TorquePidClass, i), value)
	}
}

func (m *Motor) SetRunningMode(newMode int) {
	m.runningMode = newMode
}

// func (m *Motor) Estop() {
// 	m.runningMode = RunningMode.Estopped
// 	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Estop"].cmd|m.config.secondarySumValue)
// }

// func (m *Motor) Reboot() {
// 	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Reboot"].cmd|m.config.secondarySumValue)
// }

func (m *Motor) SetSetpoint() {}

func (m *Motor) ReadAngle() float64 {
	return m.angle
}

func (m *Motor) SetAngle(angle float64) {
	m.angle = angle
	can.BuildAndSendFrame(
		m.buildCanId(14, 0),
		float32(m.angle),
		PositionControl,
		Slot_0,
	)
}

func (m *Motor) SetAngleToSlot(angle float64, slot ControlSlot) {
	m.angle = angle
	can.BuildAndSendFrame(
		m.buildCanId(14, 0),
		PositionControl,
		slot,
	)
}

func (m *Motor) ReadVelocity() float64 {
	return m.velocity
}

func (m *Motor) SetVelocity(velocity float64) {
	m.velocity = velocity
	// can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Input_Vel"].cmd|m.config.secondarySumValue, float64(m.velocity))
}

func (m *Motor) ReadTorque() float64 {
	return m.torque
}

func (m *Motor) SetTorque(torque float64) {
	m.torque = torque
	// can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Input_Torque"].cmd|m.config.secondarySumValue, float64(m.torque))
}

func (m *Motor) ReadLimits() (velocity, current float64) {
	return m.maxVelocity, m.maxCurrent
}

func (m *Motor) ReceiveCanFrame() {

}

func (m *Motor) buildCanId(apiClass, apiIndex int) int {
	return (apiClass<<4|apiIndex)<<5 | m.canid
}

// func (m *Motor) Update() {
// 	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Input_Pos"].cmd|m.config.secondarySumValue, func(event any) {
// 		i, ok := event.(int)
// 		if !ok {
// 			return
// 		}
// 		m.prevAngle = m.angle
// 		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
// 		binary.Decode(buf[:8], binary.BigEndian, &m.angle)
// 		m.angleErr = m.targetAngle - m.angle
// 	})
// 	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Input_Vel"].cmd|m.config.secondarySumValue, func(event any) {
// 		i, ok := event.(int)
// 		if !ok {
// 			return
// 		}
// 		m.prevVelocity = m.velocity
// 		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
// 		binary.Decode(buf[:8], binary.BigEndian, &m.velocity)
// 		m.velocityErr = m.targetVelocity - m.velocity
// 	})
// 	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Input_Torque"].cmd|m.config.secondarySumValue, func(event any) {
// 		i, ok := event.(int)
// 		if !ok {
// 			return
// 		}
// 		m.prevTorque = m.torque
// 		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
// 		binary.Decode(buf[:8], binary.BigEndian, &m.torque)
// 		m.torqueErr = m.targetTorque - m.torque
// 	})
// 	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Limits"].cmd|m.config.secondarySumValue, func(event any) {
// 		i, ok := event.(int)
// 		if !ok {
// 			return
// 		}
// 		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
// 		binary.Decode(buf[:4], binary.BigEndian, &m.maxVelocity)
// 		binary.Decode(buf[4:], binary.BigEndian, &m.maxCurrent)
// 	})
// }

// func (m *Motor) Status() bool {
// 	// create structure for status return, then parse [8]byte into said structure
// 	// fmt.Println("Getting CAN message frame from buffer")
// 	buf := can.GetCanMessageFromBuffer(m.config.CanId, m.config.regMap["Status"].cmd|m.config.secondarySumValue)
// 	if buf == nil {
// 		return false
// 	}
// 	// fmt.Println("got message from buffer")
// 	var status int
// 	binary.Decode(buf[:1], binary.BigEndian, status)
// 	return status == 0
// }

// default
func (c *Config) Position() *pid.Constants {
	return &(c.Slot0.Position)
}

func (c *Config) Velocity() *pid.Constants {
	return &(c.Slot0.Velocity)
}

func (c *Config) Torque() *pid.Constants {
	return &(c.Slot0.Torque)
}

func (c *Config) SetSlot0(slot Slot) {
	c.Slot0 = slot
}

func (c *Config) SetSlot1(slot Slot) {
	c.Slot1 = slot
}

func (c *Config) SetSlot2(slot Slot) {
	c.Slot2 = slot
}

func (c *Config) SetSlot3(slot Slot) {
	c.Slot3 = slot
}

// per slot
func (s *Slot) SetPositionConstants(constants pid.Constants) {
	s.Position = constants
}

func (s *Slot) SetVelocityConstants(constants pid.Constants) {
	s.Velocity = constants
}

func (s *Slot) SetTorqueConstants(constants pid.Constants) {
	s.Torque = constants
}

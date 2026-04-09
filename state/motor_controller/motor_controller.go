package motor

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/can"
)

type (
	Motor struct {
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
	}

	Config struct {
		CanId                  int
		PID_P                  float64
		PID_I                  float64
		PID_D                  float64
		PID_FF                 float64
		PID_CosFF              float64
		PositionConversion     float64
		VelocityConversion     float64
		AccelerationConversion float64
		regMap                 map[string]register
		IsSecondaryMotor       bool
		secondarySumValue      int
	}

	runningMode struct {
		Enabled  int
		Disabled int
		Estopped int
	}

	register struct {
		cmd              int
		sendStructure    []reflect.Kind
		receiveStructure []reflect.Kind
	}
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

// a motor is marked as secondary, 0x10 is ored on to what every register is used
var motorRegisterMap map[string]register = map[string]register{
	"Status": {
		cmd:              0x00,
		sendStructure:    []reflect.Kind{},
		receiveStructure: []reflect.Kind{reflect.Int8},
	},
	"Estop": {
		cmd:              0x01,
		sendStructure:    []reflect.Kind{},
		receiveStructure: []reflect.Kind{},
	},
	"Input_Pos": {
		cmd:              0x02,
		sendStructure:    []reflect.Kind{reflect.Float32, reflect.Int16, reflect.Int16},
		receiveStructure: []reflect.Kind{},
	},
	"Input_Vel": {
		cmd:              0x03,
		sendStructure:    []reflect.Kind{reflect.Float32, reflect.Int16},
		receiveStructure: []reflect.Kind{},
	},
	"Input_Torque": {
		cmd:              0x04,
		sendStructure:    []reflect.Kind{reflect.Float32},
		receiveStructure: []reflect.Kind{},
	},
	"Limits": {
		cmd:              0x05,
		sendStructure:    []reflect.Kind{reflect.Float32, reflect.Float32},
		receiveStructure: []reflect.Kind{reflect.Float32, reflect.Float32},
	},
	"Reboot": {
		cmd:              0x06,
		sendStructure:    []reflect.Kind{},
		receiveStructure: []reflect.Kind{},
	},
	"Get_Bus_Voltage_Current": {
		cmd:              0x07,
		sendStructure:    []reflect.Kind{},
		receiveStructure: []reflect.Kind{},
	},
	"Clear_Errors": {
		cmd:              0x08,
		sendStructure:    []reflect.Kind{},
		receiveStructure: []reflect.Kind{},
	},
	"Encoders": {
		cmd:              0x09,
		sendStructure:    []reflect.Kind{},
		receiveStructure: []reflect.Kind{},
	},
	"Running_Mode": {
		cmd:              0x0a,
		sendStructure:    []reflect.Kind{reflect.Int8},
		receiveStructure: []reflect.Kind{},
	},
}

func New(config Config) *Motor {
	motor := new(Motor)
	motor.Configure(config)
	return motor
}

func NewFromCanId(CanId int) *Motor {
	motor := new(Motor)
	motor.LoadRegisterMap(motorRegisterMap)
	motor.config.CanId = CanId
	can.AddDeviceToBus(motor)
	fmt.Println("Created new motor with id: ", CanId)

	return motor
}

func (m *Motor) GetCanId() int {
	return m.config.CanId
}

func (m *Motor) Configure(config Config) {
	m.config = config
	m.LoadRegisterMap(motorRegisterMap)
}

func (m *Motor) LoadRegisterMap(regMap map[string]register) {
	m.config.regMap = regMap
}

func (m *Motor) SetIsSecondaryMotorOnController(val bool) {
	m.config.IsSecondaryMotor = val
	if val {
		m.config.secondarySumValue = 0b10000
	} else {
		m.config.secondarySumValue = 0b0
	}
}

func (m *Motor) SetRunningMode(newMode int) {
	m.runningMode = newMode
}

func (m *Motor) Estop() {
	m.runningMode = RunningMode.Estopped
	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Estop"].cmd|m.config.secondarySumValue)
}

func (m *Motor) Reboot() {
	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Reboot"].cmd|m.config.secondarySumValue)
}

func (m *Motor) ReadAngle() float64 {
	return m.angle
}

func (m *Motor) SetAngle(angle float64) {
	m.angle = angle
	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Input_Pos"].cmd|m.config.secondarySumValue, float64(m.angle))
}

func (m *Motor) ReadVelocity() float64 {
	return m.velocity
}

func (m *Motor) SetVelocity(velocity float64) {
	m.velocity = velocity
	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Input_Vel"].cmd|m.config.secondarySumValue, float64(m.velocity))
}

func (m *Motor) ReadTorque() float64 {
	return m.torque
}

func (m *Motor) SetTorque(torque float64) {
	m.torque = torque
	can.BuildAndSendFrame(m.config.CanId, m.config.regMap["Input_Torque"].cmd|m.config.secondarySumValue, float64(m.torque))
}

func (m *Motor) ReadLimits() (velocity, current float64) {
	return m.maxVelocity, m.maxCurrent
}

func (m *Motor) Update() {
	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Input_Pos"].cmd|m.config.secondarySumValue, func(event any) {
		i, ok := event.(int)
		if !ok {
			return
		}
		m.prevAngle = m.angle
		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
		binary.Decode(buf[:8], binary.BigEndian, &m.angle)
		m.angleErr = m.targetAngle - m.angle
	})
	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Input_Vel"].cmd|m.config.secondarySumValue, func(event any) {
		i, ok := event.(int)
		if !ok {
			return
		}
		m.prevVelocity = m.velocity
		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
		binary.Decode(buf[:8], binary.BigEndian, &m.velocity)
		m.velocityErr = m.targetVelocity - m.velocity
	})
	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Input_Torque"].cmd|m.config.secondarySumValue, func(event any) {
		i, ok := event.(int)
		if !ok {
			return
		}
		m.prevTorque = m.torque
		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
		binary.Decode(buf[:8], binary.BigEndian, &m.torque)
		m.torqueErr = m.targetTorque - m.torque
	})
	can.BuildAndSendFrameWithCallback(m.config.CanId, m.config.regMap["Limits"].cmd|m.config.secondarySumValue, func(event any) {
		i, ok := event.(int)
		if !ok {
			return
		}
		buf := can.GetCanMessageFromBuffer(m.config.CanId, i)
		binary.Decode(buf[:4], binary.BigEndian, &m.maxVelocity)
		binary.Decode(buf[4:], binary.BigEndian, &m.maxCurrent)
	})
}

func (m *Motor) Status() bool {
	// create structure for status return, then parse [8]byte into said structure
	// fmt.Println("Getting CAN message frame from buffer")
	buf := can.GetCanMessageFromBuffer(m.config.CanId, m.config.regMap["Status"].cmd|m.config.secondarySumValue)
	if buf == nil {
		return false
	}
	// fmt.Println("got message from buffer")
	var status int
	binary.Decode(buf[:1], binary.BigEndian, status)
	return status == 0
}

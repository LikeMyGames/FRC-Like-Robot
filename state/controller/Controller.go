package controller

import (
	"errors"
	"fmt"
	"math"

	// "github.com/LikeMyGames/FRC-Like-Robot/state/command"
	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	constants "github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

const (
	// Event target used get X, Y input from the left stick
	// [event] variable is of type JoystickInput
	LeftStick = "THUMB_L"
	// Event target used get X input from the left stick
	// [event] variable is of type float64
	LeftStickX = "THUMB_LX"
	// Event target used get Y input from the left stick
	// [event] variable is of type float64
	LeftStickY = "THUMB_LY"
	// Event target used get X, Y input from the right stick
	// [event] variable is of type JoystickInput
	RightStick = "THUMB_R"
	// Event target used get X input from the right stick
	// [event] variable is of type float64
	RightStickX = "THUMB_RX"
	// Event target used get Y input from the right stick
	// [event] variable is of type float64
	RightStickY     = "THUMB_RY"
	A               = "BUTTON_A"
	B               = "BUTTON_B"
	X               = "BUTTON_X"
	Y               = "BUTTON_Y"
	DpadUP          = "DPAD_UP"
	DpadDown        = "DPAD_DOWN"
	DpadLeft        = "DPAD_LEFT"
	DpadRight       = "DPAD_RIGHT"
	LeftStickPress  = "LEFT_THUMB"
	RightStickPress = "RIGHT_THUMB"
	LeftShoulder    = "LEFT_SHOULDER"
	RightShoulder   = "RIGHT_SHOULDER"
	Start           = "START"
	Back            = "BACK"
	LeftTrigger     = "LEFT_TRIGGER"
	RightTrigger    = "RIGHT_TRIGGER"
)

type (
	ControllerAction struct {
		// whileTrue   bool
		ListenValue string
		Action      *func(any)
	}

	Controller struct {
		Config       constants.ControllerConfig
		ControllerID uint
		State        *State
	}

	State struct {
		Buttons      uint16
		LeftTrigger  uint8
		RightTrigger uint8
		ThumbLX      int16
		ThumbLY      int16
		ThumbRX      int16
		ThumbRY      int16
	}

	JoystickInput struct {
		X float64
		Y float64
	}
)

var (
	Controllers []*Controller
)

func NewController(config constants.ControllerConfig) *Controller {
	ctrl := &Controller{
		Config:       config,
		ControllerID: uint(config.ControllerNum),
		State:        &State{},
	}
	Controllers = append(Controllers, ctrl)
	return ctrl
}

func (c *Controller) GetEventTarget(target string) string {
	if target == "" {
		return ""
	}
	return fmt.Sprintf("CONTROLLER_%v_%s", c.ControllerID, target)
}

func ReadController(ctrl *Controller) {
	state := conn.LastControllerState
	if state.ControllerID == ctrl.ControllerID {
		oldState := *ctrl.State
		ctrl.State = &State{
			Buttons:      state.Buttons,
			LeftTrigger:  state.TriggerL,
			RightTrigger: state.TriggerR,
			ThumbLX:      state.ThumbLX,
			ThumbLY:      state.ThumbLY,
			ThumbRX:      state.ThumbRX,
			ThumbRY:      state.ThumbRY,
		}
		buttons := getPressedButtons(ctrl.State.Buttons)
		for _, v := range buttons {
			// fmt.Println("triggering", v)
			event.Trigger(ctrl.GetEventTarget(v), nil)
		}

		TriggerL := mathutils.MapRange(float64(ctrl.State.LeftTrigger), 0.0, 255.0, 0.0, 1.0)
		TriggerL = math.Trunc(TriggerL*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		OldTriggerL := mathutils.MapRange(float64(oldState.LeftTrigger), 0.0, 255.0, 0.0, 1.0)
		OldTriggerL = math.Trunc(OldTriggerL*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		if math.Abs(TriggerL) < ctrl.Config.Deadzones.TriggerL {
			TriggerL = 0
		}

		TriggerR := mathutils.MapRange(float64(ctrl.State.RightTrigger), 0.0, 255.0, 0.0, 1.0)
		TriggerR = math.Trunc(TriggerR*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		OldTriggerR := mathutils.MapRange(float64(oldState.RightTrigger), 0.0, 255.0, 0.0, 1.0)
		OldTriggerR = math.Trunc(OldTriggerR*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		if math.Abs(TriggerR) < ctrl.Config.Deadzones.TriggerR {
			TriggerR = 0
		}

		ThumbLX := mathutils.MapRange(float64(ctrl.State.ThumbLX), -32768.0, 32768.0, -1.0, 1.0)
		ThumbLX = math.Trunc(ThumbLX*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		OldThumbLX := mathutils.MapRange(float64(oldState.ThumbLX), -32768.0, 32768.0, -1.0, 1.0)
		OldThumbLX = math.Trunc(OldThumbLX*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		if math.Abs(ThumbLX) < ctrl.Config.Deadzones.ThumbL {
			ThumbLX = 0
		}

		ThumbLY := mathutils.MapRange(float64(ctrl.State.ThumbLY), -32768.0, 32768.0, -1.0, 1.0)
		ThumbLY = math.Trunc(ThumbLY*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		OldThumbLY := mathutils.MapRange(float64(oldState.ThumbLY), -32768.0, 32768.0, -1.0, 1.0)
		OldThumbLY = math.Trunc(OldThumbLY*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		if math.Abs(ThumbLY) < ctrl.Config.Deadzones.ThumbL {
			ThumbLY = 0
		}

		ThumbRX := mathutils.MapRange(float64(ctrl.State.ThumbRX), -32768.0, 32768.0, -1.0, 1.0)
		ThumbRX = math.Trunc(ThumbRX*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		OldThumbRX := mathutils.MapRange(float64(oldState.ThumbRX), -32768.0, 32768.0, -1.0, 1.0)
		OldThumbRX = math.Trunc(OldThumbRX*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		if math.Abs(ThumbRX) < ctrl.Config.Deadzones.ThumbR {
			ThumbRX = 0
		}

		ThumbRY := mathutils.MapRange(float64(ctrl.State.ThumbRY), -32768.0, 32768.0, -1.0, 1.0)
		ThumbRY = math.Trunc(ThumbRY*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		OldThumbRY := mathutils.MapRange(float64(oldState.ThumbRY), -32768.0, 32768.0, -1.0, 1.0)
		OldThumbRY = math.Trunc(OldThumbRY*math.Pow(10, float64(ctrl.Config.Precision))) / math.Pow(10, float64(ctrl.Config.Precision))
		if math.Abs(ThumbRY) < ctrl.Config.Deadzones.ThumbR {
			ThumbRY = 0
		}

		if math.Abs(TriggerL-OldTriggerL) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(LeftTrigger), TriggerL)
		}
		if math.Abs(TriggerR-OldTriggerR) > ctrl.Config.Deadzones.TriggerR {
			event.Trigger(ctrl.GetEventTarget(RightTrigger), TriggerR)
		}
		if math.Abs(ThumbLX-OldThumbLX) > ctrl.Config.MinChange || math.Abs(ThumbLY)-math.Abs(OldThumbLY) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(LeftStick), mathutils.Vector2D{X: ThumbLX, Y: ThumbLY})
		}
		if math.Abs(ThumbLX-OldThumbLX) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(LeftStickX), ThumbLX)
		}
		if math.Abs(ThumbLY-OldThumbLY) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(LeftStickY), ThumbLY)
		}
		if math.Abs(ThumbRX-OldThumbRX) > ctrl.Config.MinChange || math.Abs(ThumbRY)-math.Abs(OldThumbRY) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(RightStick), mathutils.Vector2D{X: ThumbRX, Y: ThumbRY})
		}
		if math.Abs(ThumbRX-OldThumbRX) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(RightStickY), ThumbRY)
		}
		if math.Abs(ThumbRY-OldThumbRY) > ctrl.Config.MinChange {
			event.Trigger(ctrl.GetEventTarget(RightStickY), ThumbRY)
		}
	}

	// return &command.Command{
	// 	Required: struct {
	// 		ctrl      *Controller
	// 		scheduler *command.CommandScheduler
	// 	}{ctrl: ctrl, scheduler: scheduler},
	// 	FirstRun:   true,
	// 	Name:       fmt.Sprintf("Read Controller ID: %v", ctrl.ControllerID),
	// 	Initialize: func() {},
	// 	Execute: func(required any) bool {
	// 		req, ok := required.(struct {
	// 			ctrl      *Controller
	// 			scheduler *command.CommandScheduler
	// 		})
	// 		if ok {
	// 			state := conn.LastControllerState
	// 			if state.ControllerID == req.ctrl.ControllerID {
	// 				req.ctrl.State = &State{
	// 					Buttons:      state.Buttons,
	// 					LeftTrigger:  state.TriggerL,
	// 					RightTrigger: state.TriggerR,
	// 					ThumbLX:      state.ThumbLX,
	// 					ThumbLY:      state.ThumbLY,
	// 					ThumbRX:      state.ThumbRX,
	// 					ThumbRY:      state.ThumbRY,
	// 				}
	// 				buttons := GetPressedButtons(req.ctrl.State.Buttons)
	// 				// fmt.Println(buttons)
	// 				contains := false
	// 				for _, action := range req.ctrl.Actions {
	// 					contains = slices.Contains(buttons, action.ListenValue)
	// 					if (action.whileTrue && contains) || (!action.whileTrue && !contains) {
	// 						// Only schedule if not already scheduled
	// 						if !slices.ContainsFunc(req.scheduler.Commands, func(command *command.Command) bool {
	// 							has := command.Name == action.Command.Name && !command.End
	// 							// fmt.Println(has)
	// 							return has
	// 						}) {
	// 							command := new(command.Command)
	// 							*command = *action.Command
	// 							req.scheduler.ScheduleCommand(command)
	// 						}
	// 					}
	// 					if action.ListenValue == LeftStick || action.ListenValue == RightStick {
	// 						command := new(command.Command)
	// 						*command = *action.Command
	// 						req.scheduler.ScheduleCommand(command)
	// 					}
	// 				}
	// 			}
	// 		}
	// 		return false
	// 	},
	// }
}

func getPressedButtons(sum uint16) []string {
	nums := []uint16{0b1,
		0b10,
		0b100,
		0b1000,
		0b10000,
		0b100000,
		0b1000000,
		0b10000000,
		0b100000000,
		0b1000000000,
		0b1000000000000,
		0b10000000000000,
		0b100000000000000,
		0b1000000000000000,
	}
	str := []string{}
	for i := len(nums) - 1; i >= 0; i-- {
		if sum >= nums[i] {
			sum -= nums[i]
			str = append(str, buttonIntToString(nums[i]))
		}
	}
	return str
}

func buttonIntToString(num uint16) string {
	switch num {
	case 1:
		return DpadUP
	case 2:
		return DpadDown
	case 4:
		return DpadLeft
	case 8:
		return DpadRight
	case 16:
		return Start
	case 32:
		return Back
	case 64:
		return LeftStickPress
	case 128:
		return RightStickPress
	case 256:
		return LeftShoulder
	case 512:
		return RightShoulder
	case 4096:
		return A
	case 8192:
		return B
	case 16384:
		return X
	case 32768:
		return Y
	default:
		return ""
	}
}

func EventDataToJoystickInput(event any) (JoystickInput, error) {
	data, ok := event.(JoystickInput)
	if ok {
		return data, nil
	}
	return JoystickInput{}, errors.New("Event data could not be converted to JoystickInput type")
}

// func (ctrl *Controller) AddAction(listenVal string, command *command.Command) *ControllerAction {
// 	action := &ControllerAction{ListenValue: listenVal, Command: command}
// 	ctrl.Actions = append(ctrl.Actions, action)
// 	return action
// }

// func (input *ControllerAction) WhileTrue() *ControllerAction {
// 	input.whileTrue = true
// 	return input
// }

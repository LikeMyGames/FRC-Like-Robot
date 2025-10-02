package controller

import (
	"fmt"

	// "github.com/LikeMyGames/FRC-Like-Robot/state/command"
	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	constants "github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
)

const (
	LeftStick       = "THUMB_L"
	RightStick      = "THUMB_R"
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
		whileTrue   bool
		ListenValue string
		Action      *func(any)
	}

	Controller struct {
		Config       constants.ControllerConfig
		ControllerID uint
		Actions      []*ControllerAction
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
)

var (
	Controllers []*Controller
)

func NewController(config constants.ControllerConfig) *Controller {
	ctrl := &Controller{
		Config:       config,
		ControllerID: uint(config.ControllerNum),
		Actions:      nil,
		State:        &State{},
	}
	Controllers = append(Controllers, ctrl)
	return ctrl
}

func ReadController(ctrl *Controller) {
	state := conn.LastControllerState
	if state.ControllerID == ctrl.ControllerID {
		fmt.Println(state)
		ctrl.State = &State{
			Buttons:      state.Buttons,
			LeftTrigger:  state.TriggerL,
			RightTrigger: state.TriggerR,
			ThumbLX:      state.ThumbLX,
			ThumbLY:      state.ThumbLY,
			ThumbRX:      state.ThumbRX,
			ThumbRY:      state.ThumbRY,
		}
		buttons := GetPressedButtons(ctrl.State.Buttons)
		for _, v := range buttons {
			event.Trigger(v, nil)
		}

		ctrl.State.LeftTrigger = uint8(mathutils.MapRange(float64(ctrl.State.LeftTrigger), 0.0, 255.0, 0.0, 1.0))
		ctrl.State.RightTrigger = uint8(mathutils.MapRange(float64(ctrl.State.RightTrigger), 0.0, 255.0, 0.0, 1.0))
		ctrl.State.ThumbLX = int16(mathutils.MapRange(float64(ctrl.State.ThumbLX), -32768.0, 32768.0, -1.0, 1.0))
		ctrl.State.ThumbLY = int16(mathutils.MapRange(float64(ctrl.State.ThumbLY), -32768.0, 32768.0, -1.0, 1.0))
		ctrl.State.ThumbRX = int16(mathutils.MapRange(float64(ctrl.State.ThumbRX), -32768.0, 32768.0, -1.0, 1.0))
		ctrl.State.ThumbRY = int16(mathutils.MapRange(float64(ctrl.State.ThumbRY), -32768.0, 32768.0, -1.0, 1.0))
		event.Trigger(LeftTrigger, ctrl.State.LeftTrigger)
		event.Trigger(RightTrigger, ctrl.State.RightTrigger)
		event.Trigger(LeftStick, struct {
			x int16
			y int16
		}{x: ctrl.State.ThumbLX, y: ctrl.State.ThumbLY})
		event.Trigger(RightStick, struct {
			x int16
			y int16
		}{x: ctrl.State.ThumbRX, y: ctrl.State.ThumbRY})
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

func GetPressedButtons(sum uint16) []string {
	nums := []uint16{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 4096, 8192, 16384, 32768}
	str := []string{}
	for i := len(nums) - 1; i >= 0; i-- {
		if sum >= nums[i] {
			sum -= nums[i]
			str = append(str, ButtonIntToString(nums[i]))
		}
	}
	return str
}

func ButtonIntToString(num uint16) string {
	switch num {
	case 1:
		return "DPAD_UP"
	case 2:
		return "DPAD_DOWN"
	case 4:
		return "DPAD_LEFT"
	case 8:
		return "DPAD_RIGHT"
	case 16:
		return "START"
	case 32:
		return "BACK"
	case 64:
		return "LEFT_THUMB"
	case 128:
		return "RIGHT_THUMB"
	case 256:
		return "LEFT_SHOULDER"
	case 512:
		return "RIGHT_SHOULDER"
	case 4096:
		return "BUTTON_A"
	case 8192:
		return "BUTTON_B"
	case 16384:
		return "BUTTON_X"
	case 32768:
		return "BUTTON_Y"
	default:
		return ""
	}
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

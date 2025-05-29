package Controller

import (
	"fmt"
	"slices"

	"frcrobot/internal/Command"
	"frcrobot/internal/File"
	"frcrobot/internal/GUI"

	"github.com/tajtiattila/xinput"
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
)

type (
	ControllerConfig struct {
		ControllerNum int                  `json:"controllerNum"`
		Deadzones     ConstrollerDeadzones `json:"deadzones"`
	}

	ControllerAction struct {
		whileTrue   bool
		ListenValue string
		Command     *Command.Command
	}

	Controller struct {
		Config       ControllerConfig
		ControllerID uint
		Actions      []*ControllerAction
		State        *xinput.State
	}

	ConstrollerDeadzones struct {
		ThumbL   float32 `json:"thumbL"`
		ThumbR   float32 `json:"thumbR"`
		TriggerL float32 `json:"triggerL"`
		TriggerR float32 `json:"triggerR"`
	}
)

var (
	Controllers []*Controller
)

func StartController(controllerID uint, scheduler *Command.CommandScheduler) *Controller {
	config := ControllerConfig{}
	File.ReadJSON("controller.config", &config)
	fmt.Println("Controller Config: ", config)

	ctrl := &Controller{
		Config:       config,
		ControllerID: controllerID,
		Actions:      nil,
		State:        &xinput.State{PacketNumber: 0, Gamepad: xinput.Gamepad{}},
	}
	Controllers = append(Controllers, ctrl)
	scheduler.ScheduleCommand(NewReadControllerCommand(ctrl, scheduler))
	return ctrl
}

func NewReadControllerCommand(ctrl *Controller, scheduler *Command.CommandScheduler) *Command.Command {
	return &Command.Command{
		Required: struct {
			ctrl      *Controller
			scheduler *Command.CommandScheduler
		}{ctrl: ctrl, scheduler: scheduler},
		FirstRun:   true,
		Name:       fmt.Sprintf("Read Controller ID: %v", ctrl.ControllerID),
		Initialize: func() {},
		Execute: func(required any) bool {
			req, ok := required.(struct {
				ctrl      *Controller
				scheduler *Command.CommandScheduler
			})
			if ok {
				state := GUI.LastControllerState
				if state.ControllerID == req.ctrl.ControllerID {
					req.ctrl.State.Gamepad = xinput.Gamepad{
						Buttons:      state.Buttons,
						LeftTrigger:  state.TriggerL,
						RightTrigger: state.TriggerR,
						ThumbLX:      state.ThumbLX,
						ThumbLY:      state.ThumbLY,
						ThumbRX:      state.ThumbRX,
						ThumbRY:      state.ThumbRY,
					}
					buttons := GetPressedButtons(req.ctrl.State.Gamepad.Buttons)
					// fmt.Println(buttons)
					contains := false
					for _, action := range req.ctrl.Actions {
						contains = slices.Contains(buttons, action.ListenValue)
						if (action.whileTrue && contains) || (!action.whileTrue && !contains) {
							// Only schedule if not already scheduled
							if !slices.ContainsFunc(req.scheduler.Commands, func(command *Command.Command) bool {
								has := command.Name == action.Command.Name && !command.End
								// fmt.Println(has)
								return has
							}) {
								command := new(Command.Command)
								*command = *action.Command
								req.scheduler.ScheduleCommand(command)
							}
						}
						if action.ListenValue == LeftStick || action.ListenValue == RightStick {
							command := new(Command.Command)
							*command = *action.Command
							req.scheduler.ScheduleCommand(command)
						}
					}
				}
			}
			return false
		},
	}
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

func (ctrl *Controller) AddAction(listenVal string, command *Command.Command) *ControllerAction {
	action := &ControllerAction{ListenValue: listenVal, Command: command}
	ctrl.Actions = append(ctrl.Actions, action)
	return action
}

func (input *ControllerAction) WhileTrue() *ControllerAction {
	input.whileTrue = true
	return input
}

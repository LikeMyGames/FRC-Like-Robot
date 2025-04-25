package Controller

import (
	"fmt"

	"frcrobot/internal/Command"
	"frcrobot/internal/Command/Commands/Controller/ReadController"
	"frcrobot/internal/File"
	"frcrobot/internal/GUI"

	"github.com/orsinium-labs/gamepad"
)

type (
	ControllerConfig struct {
		ControllerNum int                  `json:"controllerNum"`
		Deadzones     ConstrollerDeadzones `json:"deadzones"`
	}

	ControllerInput struct {
		whileTrue   bool
		ListenValue string
		Command     *Command.Command
	}

	Controller struct {
		Config       ControllerConfig
		ControllerID int
		Inputs       []*ControllerInput
		GamePad      *gamepad.GamePad
	}
)

type ConstrollerDeadzones struct {
	ThumbL   float32 `json:"thumbL"`
	ThumbR   float32 `json:"thumbR"`
	TriggerL float32 `json:"triggerL"`
	TriggerR float32 `json:"triggerR"`
}

func StartController(controllerID int, scheduler *Command.CommandScheduler) *Controller {
	config := ControllerConfig{}
	File.ReadJSON("controller.config", &config)
	fmt.Println("Controller Config: ", config)
	controller := &Controller{
		Config:       config,
		ControllerID: controllerID,
	}
	scheduler.ScheduleCommand(NewConnectControllerCommand(controller, scheduler))
	return controller

	// var pressedButtons []string
	// var thumbL []float32
	// var thumbR []float32
	// var triggerL float32
	// var triggerR float32
	// var pressedButtonsOLD []string
	// var thumbLOld []float32
	// var thumbROld []float32
	// var triggerLOld float32
	// var triggerROld float32
	// first := true
	// for {
	// 	if !first {
	// 		controllerState, err = controller.State()
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		controllerState.
	// 			thumbL = []float32{float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbLX), -32768, 32768, -1, 1)), 4)), float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbLY), -32768, 32768, -1, 1)), 4))}
	// 		thumbR = []float32{float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbRX), -32768, 32768, -1, 1)), 4)), float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbRY), -32768, 32768, -1, 1)), 4))}
	// 		triggerL = float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.LeftTrigger), -32768, 32768, -1, 1)), 4))
	// 		triggerR = float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.RightTrigger), -32768, 32768, -1, 1)), 4))
	// 		// fmt.Println("ThumbL: ", thumbL, "\tThumbR: ", thumbR, "\tTriggerL:", triggerL, "\tTriggerR", triggerR, "\tButtons: ", pressedButtons)

	// 		//rate of change based analog axis EventListenerListener emitters
	// 		// if (math.Abs(float64(thumbLOld[0])-float64(thumbL[0]))) > 0.05 || (math.Abs(float64(thumbLOld[1])-float64(thumbL[1]))) > 0.05 {
	// 		// 	EventListener.Emit("THUMB_L", thumbL)
	// 		// }
	// 		// if (math.Abs(float64(thumbROld[0])-float64(thumbR[0]))) > 0.05 || (math.Abs(float64(thumbROld[1])-float64(thumbR[1]))) > 0.05 {
	// 		// 	EventListener.Emit("THUMB_R", thumbR)
	// 		// }
	// 		// if (math.Abs(float64(triggerLOld) - float64(triggerL))) > 0.05 {
	// 		// 	EventListener.Emit("TRIGGER_L", triggerL)
	// 		// }
	// 		// if (math.Abs(float64(triggerROld) - float64(triggerR))) > 0.05 {
	// 		// 	EventListener.Emit("TRIGGER_R", triggerR)
	// 		// }

	// 		//deadzone based analog axis EventListener emmiters
	// 		if math.Abs(float64(thumbL[0])) > float64(config.Deadzones.ThumbL) || math.Abs(float64(thumbL[1])) > float64(config.Deadzones.ThumbL) || thumbL[0] == 0 || thumbL[1] == 0 {
	// 			if (math.Abs(float64(thumbLOld[0])-float64(thumbL[0]))) > 0.005 || (math.Abs(float64(thumbLOld[1])-float64(thumbL[1]))) > 0.005 {
	// 				EventListener.Emit("THUMB_L", thumbL)
	// 			}
	// 		}
	// 		if math.Abs(float64(thumbR[0])) > float64(config.Deadzones.ThumbR) || math.Abs(float64(thumbR[1])) > float64(config.Deadzones.ThumbR) || thumbR[0] == 0 || thumbR[1] == 0 {
	// 			if (math.Abs(float64(thumbROld[0])-float64(thumbR[0]))) > 0.005 || (math.Abs(float64(thumbROld[1])-float64(thumbR[1]))) > 0.005 {
	// 				EventListener.Emit("THUMB_R", thumbR)
	// 			}
	// 		}
	// 		if math.Abs(float64(triggerL)) > float64(config.Deadzones.TriggerL) || triggerL == 0 {
	// 			if (math.Abs(float64(triggerLOld) - float64(triggerL))) > 0.005 {
	// 				EventListener.Emit("TRIGGER_L", triggerL)
	// 			}
	// 		}
	// 		if math.Abs(float64(triggerR)) > float64(config.Deadzones.TriggerR) || triggerR == 0 {
	// 			if (math.Abs(float64(triggerROld) - float64(triggerR))) > 0.005 {
	// 				EventListener.Emit("TRIGGER_R", triggerR)
	// 			}
	// 		}

	// 		for _, v := range pressedButtons {
	// 			if !slices.Contains(pressedButtonsOLD, v) {
	// 				EventListener.Emit(fmt.Sprintf("%v_PRESS", v))
	// 			}
	// 		}
	// 		for _, v := range pressedButtonsOLD {
	// 			if !slices.Contains(pressedButtons, v) {
	// 				EventListener.Emit(fmt.Sprintf("%v_RELEASE", v))
	// 			}
	// 		}
	// 		thumbLOld = thumbL
	// 		thumbROld = thumbR
	// 		triggerLOld = triggerL
	// 		triggerROld = triggerR
	// 		pressedButtonsOLD = pressedButtons
	// 	}
	// 	if first {
	// 		first = false
	// 		EventListener.Listen("BACK_PRESS", func(a ...any) any {
	// 			os.Exit(0)
	// 			return nil
	// 		})
	// 		thumbLOld = []float32{float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbLX), -32768, 32768, -1, 1)), 4)), float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbLY), -32768, 32768, -1, 1)), 4))}
	// 		thumbROld = []float32{float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbRX), -32768, 32768, -1, 1)), 4)), float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.ThumbRY), -32768, 32768, -1, 1)), 4))}
	// 		triggerLOld = float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.LeftTrigger), -32768, 32768, -1, 1)), 4))
	// 		triggerROld = float32(MathUtils.Trunc(float64(MathUtils.MapRange(float32(controller.RightTrigger), -32768, 32768, -1, 1)), 4))
	// 		pressedButtonsOLD = getPressedButtons(controller.Buttons)
	// 	}
	// }
}

func NewConnectControllerCommand(controller *Controller, scheduler *Command.CommandScheduler) *Command.Command {
	return &Command.Command{
		Required:   controller,
		Name:       "Connect Controller",
		FirstRun:   true,
		Initialize: func() {},
		Execute: func(controller any) {
			ctrl := controller.(*Controller)
			Gamepad, err := gamepad.NewGamepad(ctrl.ControllerID)
			if err != nil {
				GUI.SendData([]byte(`{"system_logger":{"type":"warn","message":"GamePad not connected"},"robot_status":{"type":"joysticks","value":false}}`))
				return
			}
			ctrl.GamePad = Gamepad
			scheduler.ScheduleCommand(ReadController.NewReadControllerCommand(ctrl.GamePad))
		},
		End: func() bool {
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

func (controller *Controller) AddControllerInput(value string, command *Command.Command) ControllerInput {
	input := ControllerInput{}
	controller.Inputs = append(controller.Inputs, &input)
	return input
}

func (input *ControllerInput) WhileTrue() *ControllerInput {
	input.whileTrue = true
	return input
}

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/tajtiattila/xinput"
)

type ControllerConfig struct {
	controllerNum int        `json:"controllerNum"`
	deadzones     []Deadzone `json:"deadzones"`
}

type Deadzone struct {
	name  string  `json:"name"`
	value float32 `json:"value"`
}

func StartController() {
	config := ControllerConfig{}
	readJSON("controllerConfig", config)
	fmt.Println(config)

	fmt.Println("starting controller")
	fmt.Println(runtime.GOOS)
	if runtime.GOOS != "windows" {
		fmt.Println("this program only works on windows")
		os.Exit(1)
	}
	err := xinput.Load()
	if err == nil {
		var controllerState xinput.State
		var controller xinput.Gamepad
		var pressedButtons []string
		var thumbL []float32
		var thumbR []float32
		var triggerL float32
		var triggerR float32
		first := true
		for {
			if !first {
				xinput.GetState(0, &controllerState)
				controller = controllerState.Gamepad
				pressedButtons = getPressedButtEventListener.Listen(controller.Buttons)
				thumbL = []float32{mapRange(float32(controller.ThumbLX), -32768, 32768, -1, 1, 4), mapRange(float32(controller.ThumbLY), -32768, 32768, -1, 1, 4)}
				thumbR = []float32{mapRange(float32(controller.ThumbRX), -32768, 32768, -1, 1, 4), mapRange(float32(controller.ThumbRY), -32768, 32768, -1, 1, 4)}
				triggerL = mapRange(float32(controller.LeftTrigger), 0, 255, 0, 1, 4)
				triggerR = mapRange(float32(controller.RightTrigger), 0, 255, 0, 1, 4)
				// fmt.Println("ThumbL: ", thumbL, "\tThumbR: ", thumbR, "\tTriggerL:", triggerL, "\tTriggerR", triggerR, "\tButtons: ", pressedButtons)
				EventListener.Listen("START", func(a ...any) any {
					os.Exit(0)
					return nil
				})
				EventListener.Emit("THUMB_L", thumbL)
				EventListener.Emit("THUMB_R", thumbR)
				EventListener.Emit("TRIGGER_L", triggerL)
				EventListener.Emit("TRIGGER_R", triggerR)
				for _, v := range pressedButtons {
					EventListener.Emit(v)
				}
			}
			if first {
				first = false
			}
			// fmt.Println("Left Trigger: ", float64(controller.LeftTrigger)*64, "\tRight Trigger: ", float64(controller.RightTrigger)*64)
			// xinput.SetState(0, &xinput.Vibration{LeftMotorSpeed: uint16(controller.LeftTrigger) * 64, RightMotorSpeed: uint16(controller.RightTrigger) * 64})
		}
	}

}

func mapRange(num, inLow, inHigh, outLow, outHigh float32, trunc int) float32 {
	numRet, err := strconv.ParseFloat(fmt.Sprintf("%f", (outLow + (((num - inLow) / (inHigh - inLow)) * (outHigh - outLow))))[:trunc+1], 32)
	if err != nil {
		os.Exit(2)
	}
	return float32(numRet)
}

func getPressedButtEventListener.Listen(sum uint16) []string {
	nums := []uint16{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 4096, 8192, 16384, 32768}
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

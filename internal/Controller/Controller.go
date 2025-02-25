package Controller

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"

	Event "internal/EventListener"
	File "internal/File"

	"github.com/tajtiattila/xinput"
)

type ControllerConfig struct {
	ControllerNum int                  `json:"controllerNum"`
	Deadzones     ConstrollerDeadzones `json:"deadzones"`
}

type ConstrollerDeadzones struct {
	ThumbL   float32 `json:"thumbL"`
	ThumbR   float32 `json:"thumbR"`
	TriggerL float32 `json:"triggerL"`
	TriggerR float32 `json:"triggerR"`
}

func StartController() {
	config := ControllerConfig{}
	File.ReadJSON("controller.config", &config)
	fmt.Println(config)

	fmt.Println("starting controller ...")
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
		var thumbLOld []float32
		var thumbROld []float32
		var triggerLOld float32
		var triggerROld float32
		first := true
		for {
			if !first {
				xinput.GetState(0, &controllerState)
				controller = controllerState.Gamepad
				pressedButtons = getPressedButtons(controller.Buttons)
				thumbL = []float32{mapRange(float32(controller.ThumbLX), -32768, 32768, -1, 1, 4), mapRange(float32(controller.ThumbLY), -32768, 32768, -1, 1, 4)}
				thumbR = []float32{mapRange(float32(controller.ThumbRX), -32768, 32768, -1, 1, 4), mapRange(float32(controller.ThumbRY), -32768, 32768, -1, 1, 4)}
				triggerL = mapRange(float32(controller.LeftTrigger), 0, 255, 0, 1, 4)
				triggerR = mapRange(float32(controller.RightTrigger), 0, 255, 0, 1, 4)
				// fmt.Println("ThumbL: ", thumbL, "\tThumbR: ", thumbR, "\tTriggerL:", triggerL, "\tTriggerR", triggerR, "\tButtons: ", pressedButtons)

				//rate of change based analog axis event emitters
				// if (math.Abs(float64(thumbLOld[0])-float64(thumbL[0]))) > 0.05 || (math.Abs(float64(thumbLOld[1])-float64(thumbL[1]))) > 0.05 {
				// 	Event.Emit("THUMB_L", thumbL)
				// }
				// if (math.Abs(float64(thumbROld[0])-float64(thumbR[0]))) > 0.05 || (math.Abs(float64(thumbROld[1])-float64(thumbR[1]))) > 0.05 {
				// 	Event.Emit("THUMB_R", thumbR)
				// }
				// if (math.Abs(float64(triggerLOld) - float64(triggerL))) > 0.05 {
				// 	Event.Emit("TRIGGER_L", triggerL)
				// }
				// if (math.Abs(float64(triggerROld) - float64(triggerR))) > 0.05 {
				// 	Event.Emit("TRIGGER_R", triggerR)
				// }

				//deadzone based analog axis event emmiters
				if math.Abs(float64(thumbL[0])) > float64(config.Deadzones.ThumbL) || math.Abs(float64(thumbL[1])) > float64(config.Deadzones.ThumbL) || thumbL[0] == 0 || thumbL[1] == 0 {
					if (math.Abs(float64(thumbLOld[0])-float64(thumbL[0]))) > 0.005 || (math.Abs(float64(thumbLOld[1])-float64(thumbL[1]))) > 0.005 {
						Event.Emit("THUMB_L", thumbL)
					}
				}
				if math.Abs(float64(thumbR[0])) > float64(config.Deadzones.ThumbR) || math.Abs(float64(thumbR[1])) > float64(config.Deadzones.ThumbR) || thumbR[0] == 0 || thumbR[1] == 0 {
					if (math.Abs(float64(thumbROld[0])-float64(thumbR[0]))) > 0.005 || (math.Abs(float64(thumbROld[1])-float64(thumbR[1]))) > 0.005 {
						Event.Emit("THUMB_R", thumbR)
					}
				}
				if math.Abs(float64(triggerL)) > float64(config.Deadzones.TriggerL) || triggerL == 0 {
					if (math.Abs(float64(triggerLOld) - float64(triggerL))) > 0.005 {
						Event.Emit("TRIGGER_L", triggerL)
					}
				}
				if math.Abs(float64(triggerR)) > float64(config.Deadzones.TriggerR) || triggerR == 0 {
					if (math.Abs(float64(triggerROld) - float64(triggerR))) > 0.005 {
						Event.Emit("TRIGGER_R", triggerR)
					}
				}

				for _, v := range pressedButtons {
					Event.Emit(v)
				}
				thumbLOld = thumbL
				thumbROld = thumbR
				triggerLOld = triggerL
				triggerROld = triggerR
			}
			if first {
				first = false
				Event.Listen("BACK", func(a ...any) any {
					os.Exit(0)
					return nil
				})
				thumbLOld = []float32{mapRange(float32(controller.ThumbLX), -32768, 32768, -1, 1, 4), mapRange(float32(controller.ThumbLY), -32768, 32768, -1, 1, 4)}
				thumbROld = []float32{mapRange(float32(controller.ThumbRX), -32768, 32768, -1, 1, 4), mapRange(float32(controller.ThumbRY), -32768, 32768, -1, 1, 4)}
				triggerLOld = mapRange(float32(controller.LeftTrigger), 0, 255, 0, 1, 4)
				triggerROld = mapRange(float32(controller.RightTrigger), 0, 255, 0, 1, 4)
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

func clamp(num, low, high float32) float32 {
	if num >= high {
		return high
	} else if num <= low {
		return low
	}
	return num
}

func getPressedButtons(sum uint16) []string {
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

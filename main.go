package main

import (
	"fmt"
	Controller "internal/Controller"
	Event "internal/EventListener"
	Webpage "internal/Webpage"
	"slices"
)

func main() {
	pressedButtons := make([]string, 0)
	thumbL := make([]float32, 0)
	thumbR := make([]float32, 0)
	triggerL := float32(0.0)
	triggerR := float32(0.0)

	// DPAD_UP Input
	Event.Listen("DPAD_UP_PRESS", func(a ...any) any {
		fmt.Println("DPAD_UP_PRESS")
		pressedButtons = append(pressedButtons, "DPAD_UP_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("DPAD_UP_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "DPAD_UP_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("DPAD_UP_RELEASE")
		return nil
	})

	// DPAD_DOWN Input
	Event.Listen("DPAD_DOWN_PRESS", func(a ...any) any {
		fmt.Println("DPAD_DOWN_PRESS")
		pressedButtons = append(pressedButtons, "DPAD_DOWN_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("DPAD_DOWN_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "DPAD_DOWN_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("DPAD_DOWN_RELEASE")
		return nil
	})

	// DPAD_LEFT Input
	Event.Listen("DPAD_LEFT_PRESS", func(a ...any) any {
		fmt.Println("DPAD_LEFT_PRESS")
		pressedButtons = append(pressedButtons, "DPAD_LEFT_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("DPAD_LEFT_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "DPAD_LEFT_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("DPAD_LEFT_RELEASE")
		return nil
	})

	// DPAD_RIGHT Input
	Event.Listen("DPAD_RIGHT_PRESS", func(a ...any) any {
		fmt.Println("DPAD_RIGHT_PRESS")
		pressedButtons = append(pressedButtons, "DPAD_RIGHT_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("DPAD_RIGHT_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "DPAD_RIGHT_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("DPAD_RIGHT_RELEASE")
		return nil
	})

	// START Input
	Event.Listen("START_PRESS", func(a ...any) any {
		fmt.Println("START_PRESS")
		Event.ListListeners()
		pressedButtons = append(pressedButtons, "START_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("START_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "START_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("START_RELEASE")
		return nil
	})

	// BACK Input
	Event.Listen("BACK_PRESS", func(a ...any) any {
		fmt.Println("BACK_PRESS")
		pressedButtons = append(pressedButtons, "BACK_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("BACK_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "BACK_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("BACK_RELEASE")
		return nil
	})

	// LEFT_THUMB Input
	Event.Listen("LEFT_THUMB_PRESS", func(a ...any) any {
		fmt.Println("LEFT_THUMB_PRESS")
		pressedButtons = append(pressedButtons, "LEFT_THUMB_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("LEFT_THUMB_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "LEFT_THUMB_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("LEFT_THUMB_RELEASE")
		return nil
	})

	// RIGHT_THUMB Input
	Event.Listen("RIGHT_THUMB_PRESS", func(a ...any) any {
		fmt.Println("RIGHT_THUMB_PRESS")
		pressedButtons = append(pressedButtons, "RIGHT_THUMB_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("RIGHT_THUMB_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "RIGHT_THUMB_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("RIGHT_THUMB_RELEASE")
		return nil
	})

	// LEFT_SHOULDER Input
	Event.Listen("LEFT_SHOULDER_PRESS", func(a ...any) any {
		fmt.Println("LEFT_SHOULDER_PRESS")
		pressedButtons = append(pressedButtons, "LEFT_SHOULDER_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("LEFT_SHOULDER_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "LEFT_SHOULDER_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("LEFT_SHOULDER_RELEASE")
		return nil
	})

	// RIGHT_SHOULDER Input
	Event.Listen("RIGHT_SHOULDER_PRESS", func(a ...any) any {
		fmt.Println("RIGHT_SHOULDER_PRESS")
		pressedButtons = append(pressedButtons, "RIGHT_SHOULDER_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("RIGHT_SHOULDER_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "RIGHT_SHOULDER_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("RIGHT_SHOULDER_RELEASE")
		return nil
	})

	// BUTTON_A Input
	Event.Listen("BUTTON_A_PRESS", func(a ...any) any {
		fmt.Println("BUTTON_A_PRESS")
		pressedButtons = append(pressedButtons, "BUTTON_A_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("BUTTON_A_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "BUTTON_A_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("BUTTON_A_RELEASE")
		return nil
	})

	// BUTTON_B Input
	Event.Listen("BUTTON_B_PRESS", func(a ...any) any {
		fmt.Println("BUTTON_B_PRESS")
		pressedButtons = append(pressedButtons, "BUTTON_B_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("BUTTON_B_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "BUTTON_B_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("BUTTON_B_RELEASE")
		return nil
	})

	// BUTTON_X Input
	Event.Listen("BUTTON_X_PRESS", func(a ...any) any {
		fmt.Println("BUTTON_X_PRESS")
		pressedButtons = append(pressedButtons, "BUTTON_X_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("BUTTON_X_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "BUTTON_X_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("BUTTON_X_RELEASE")
		return nil
	})

	// BUTTON_Y Input
	Event.Listen("BUTTON_Y_PRESS", func(a ...any) any {
		fmt.Println("BUTTON_Y_PRESS")
		pressedButtons = append(pressedButtons, "BUTTON_Y_PRESS")
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})
	Event.Listen("BUTTON_Y_RELEASE", func(a ...any) any {
		i, t := slices.BinarySearch(pressedButtons, "BUTTON_Y_PRESS")
		if t {
			pressedButtons = slices.Delete(pressedButtons, i, i+1)
		}
		fmt.Println("BUTTON_Y_RELEASE")
		return nil
	})

	// THUMB_L Input
	Event.Listen("THUMB_L", func(a ...any) any {
		fmt.Println("THUMB_L: ", a[0])
		thumbL = a[0].([]float32)
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})

	// THUMB_R Input
	Event.Listen("THUMB_R", func(a ...any) any {
		fmt.Println("THUMB_R: ", a[0])
		thumbR = a[0].([]float32)
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})

	// TRIGGER_L Input
	Event.Listen("TRIGGER_L", func(a ...any) any {
		fmt.Println("TRIGGER_L: ", a[0])
		triggerL = a[0].(float32)
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})

	// TRIGGER_R Input
	Event.Listen("TRIGGER_R", func(a ...any) any {
		fmt.Println("TRIGGER_R: ", a[0])
		triggerR = a[0].(float32)
		Webpage.Send(thumbL, thumbR, triggerL, triggerR, pressedButtons)
		return nil
	})

	go Controller.StartController()
	Webpage.Start()
	//  Webpage.SendVariables()
	//  Call at end of file
}

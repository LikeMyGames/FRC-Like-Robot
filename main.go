package main

import (
	"fmt"
	Controller "internal/Controller"
	Event "internal/EventListener"
	Webpage "internal/Webpage"
)

func main() {
	// var pressedButtons []string
	// var thumbL []float32
	// var thumbR []float32
	// var triggerL float32
	// var triggerR float32

	Event.Listen("DPAD_UP", func(a ...any) any {
		fmt.Println("DPAD_UP")
		return nil
	})
	Event.Listen("DPAD_DOWN", func(a ...any) any {
		fmt.Println("DPAD_DOWN")
		return nil
	})
	Event.Listen("DPAD_LEFT", func(a ...any) any {
		fmt.Println("DPAD_LEFT")
		return nil
	})
	Event.Listen("DPAD_RIGHT", func(a ...any) any {
		fmt.Println("DPAD_RIGHT")
		return nil
	})
	Event.Listen("START", func(a ...any) any {
		fmt.Println("START")
		Event.ListListeners()
		return nil
	})
	Event.Listen("BACK", func(a ...any) any {
		fmt.Println("BACK")
		return nil
	})
	Event.Listen("LEFT_THUMB", func(a ...any) any {
		fmt.Println("LEFT_THUMB")
		return nil
	})
	Event.Listen("RIGHT_THUMB", func(a ...any) any {
		fmt.Println("RIGHT_THUMB")
		return nil
	})
	Event.Listen("LEFT_SHOULDER", func(a ...any) any {
		fmt.Println("LEFT_SHOULDER")
		return nil
	})
	Event.Listen("RIGHT_SHOULDER", func(a ...any) any {
		fmt.Println("RIGHT_SHOULDER")
		return nil
	})
	Event.Listen("BUTTON_A", func(a ...any) any {
		fmt.Println("BUTTON_A")
		return nil
	})
	Event.Listen("BUTTON_B", func(a ...any) any {
		fmt.Println("BUTTON_B")
		return nil
	})
	Event.Listen("BUTTON_X", func(a ...any) any {
		fmt.Println("BUTTON_X")
		return nil
	})
	Event.Listen("BUTTON_Y", func(a ...any) any {
		fmt.Println("BUTTON_Y")
		return nil
	})
	Event.Listen("THUMB_L", func(a ...any) any {
		fmt.Println("THUMB_L: ", a[0])
		return nil
	})
	Event.Listen("THUMB_R", func(a ...any) any {
		fmt.Println("THUMB_R: ", a[0])
		return nil
	})
	Event.Listen("TRIGGER_L", func(a ...any) any {
		fmt.Println("TRIGGER_L: ", a[0])
		return nil
	})
	Event.Listen("TRIGGER_R", func(a ...any) any {
		fmt.Println("TRIGGER_R: ", a[0])
		return nil
	})
	go Controller.StartController()
	Webpage.Start(5000)
	// Webpage.SendVariables()
	// Call at end of file
}

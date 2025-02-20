package Robot

import (
	"github.com/LikeMyGames/FRC-Like-Robot/Controller"
	"github.com/LikeMyGames/FRC-Like-Robot/EventListener"
)

func main() {
	EventListener.Listen("DPAD_UP", func(a ...any) any {

		return nil
	})
	EventListener.Listen("DPAD_DOWN", func(a ...any) any {

		return nil
	})
	EventListener.Listen("DPAD_LEFT", func(a ...any) any {

		return nil
	})
	EventListener.Listen("DPAD_RIGHT", func(a ...any) any {

		return nil
	})
	EventListener.Listen("START", func(a ...any) any {

		return nil
	})
	EventListener.Listen("BACK", func(a ...any) any {

		return nil
	})
	EventListener.Listen("LEFT_THUMB", func(a ...any) any {

		return nil
	})
	EventListener.Listen("RIGHT_THUMB", func(a ...any) any {

		return nil
	})
	EventListener.Listen("LEFT_SHOULDER", func(a ...any) any {

		return nil
	})
	EventListener.Listen("RIGHT_SHOULDER", func(a ...any) any {

		return nil
	})
	EventListener.Listen("BUTTON_A", func(a ...any) any {

		return nil
	})
	EventListener.Listen("BUTTON_B", func(a ...any) any {

		return nil
	})
	EventListener.Listen("BUTTON_X", func(a ...any) any {

		return nil
	})
	EventListener.Listen("BUTTON_Y", func(a ...any) any {

		return nil
	})
	EventListener.Listen("THUMB_L", func(a ...any) any {

		return nil
	})
	EventListener.Listen("THUMB_R", func(a ...any) any {

		return nil
	})
	EventListener.Listen("TRIGGER_L", func(a ...any) any {

		return nil
	})
	EventListener.Listen("TRIGGER_R", func(a ...any) any {

		return nil
	})
	// Call at end of file
	go Controller.StartController()
}

package Robot

import (
	"github.com/LikeMyGames/FRC-Like-Robot/internal/Controller"
	Event "github.com/LikeMyGames/FRC-Like-Robot/internal/EventListener"
)

func main() {
	Event.Listen("DPAD_UP", func(a ...any) any {

		return nil
	})
	Event.Listen("DPAD_DOWN", func(a ...any) any {

		return nil
	})
	Event.Listen("DPAD_LEFT", func(a ...any) any {

		return nil
	})
	Event.Listen("DPAD_RIGHT", func(a ...any) any {

		return nil
	})
	Event.Listen("START", func(a ...any) any {

		return nil
	})
	Event.Listen("BACK", func(a ...any) any {

		return nil
	})
	Event.Listen("LEFT_THUMB", func(a ...any) any {

		return nil
	})
	Event.Listen("RIGHT_THUMB", func(a ...any) any {

		return nil
	})
	Event.Listen("LEFT_SHOULDER", func(a ...any) any {

		return nil
	})
	Event.Listen("RIGHT_SHOULDER", func(a ...any) any {

		return nil
	})
	Event.Listen("BUTTON_A", func(a ...any) any {

		return nil
	})
	Event.Listen("BUTTON_B", func(a ...any) any {

		return nil
	})
	Event.Listen("BUTTON_X", func(a ...any) any {

		return nil
	})
	Event.Listen("BUTTON_Y", func(a ...any) any {

		return nil
	})
	Event.Listen("THUMB_L", func(a ...any) any {

		return nil
	})
	Event.Listen("THUMB_R", func(a ...any) any {

		return nil
	})
	Event.Listen("TRIGGER_L", func(a ...any) any {

		return nil
	})
	Event.Listen("TRIGGER_R", func(a ...any) any {

		return nil
	})
	// Call at end of file
	go Controller.StartController()
}

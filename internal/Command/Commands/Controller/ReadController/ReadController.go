package ReadController

import (
	"frcrobot/internal/Command"
	"log"

	"github.com/orsinium-labs/gamepad"
)

func NewReadControllerCommand(controller *gamepad.GamePad) *Command.Command {
	return &Command.Command{
		Required:   controller,
		Name:       "Read Controller",
		FirstRun:   true,
		Initialize: func() {},
		Execute: func(required any) {
			_, err := required.(*gamepad.GamePad).State()
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println("Controller State: ", controllerState)
		},
		End: func() bool {
			return false
		},
	}
}

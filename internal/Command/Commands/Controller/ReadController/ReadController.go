package ReadController

import (
	"fmt"
	"frcrobot/internal/Command"
	"log"

	"github.com/orsinium-labs/gamepad"
)

func NewReadControllerCommand(controller *gamepad.GamePad) *Command.Command {
	return &Command.Command{
		Required:   controller,
		Name:       "GetControllerInputs",
		FirstRun:   true,
		Initialize: func() {},
		Execute: func(required any) {
			controllerState, err := controller.State()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Controller State: ", controllerState)
		},
		End: func() bool {
			return false
		},
	}
}

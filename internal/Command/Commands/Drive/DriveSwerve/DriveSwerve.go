package DriveSwerve

import (
	"frcrobot/internal/Command"
	"frcrobot/internal/DriveSubsystem"
	"log"

	"github.com/orsinium-labs/gamepad"
)

func NewDriveSwerveCommand(drive *DriveSubsystem.SwerveDrive) *Command.Command {
	return &Command.Command{
		Required:   drive,
		Name:       "Drive Swerve",
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

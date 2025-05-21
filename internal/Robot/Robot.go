package Robot

import (
	"fmt"
	"frcrobot/internal/Command"
	"frcrobot/internal/Controller"
	"frcrobot/internal/DriveSubsystem"
)

type (
	Robot struct {
		DriveSubsystem *DriveSubsystem.SwerveDrive
		Controllers    []*Controller.Controller
		Scheduler      *Command.CommandScheduler
	}
)

var (
	robot *Robot
)

func NewRobot(controllerID []uint) *Robot {
	// Create a new scheduler for the robot
	scheduler := Command.NewCommandScheduler()

	// Initialize the drive subsystem
	drive := DriveSubsystem.NewSwerveDrive("robot.constants")
	fmt.Println(drive)

	controllers := make([]*Controller.Controller, len(controllerID))

	for i, v := range controllerID {
		controllers[i] = Controller.StartController(v, scheduler)
	}

	controllers[0].AddAction("BUTTON_B", &Command.Command{
		Required:   nil,
		FirstRun:   true,
		Name:       "button b input",
		Initialize: func() {},
		Execute: func(required any) bool {
			fmt.Println("button b pressed")
			return true
		},
	}).WhileTrue()

	return &Robot{
		Scheduler:      scheduler,
		DriveSubsystem: drive,
		Controllers:    controllers,
	}
}

func (r *Robot) Start() {
	r.Scheduler.Start()
}

func GetRobot() *Robot {
	return robot
}

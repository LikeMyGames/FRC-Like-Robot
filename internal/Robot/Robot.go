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
		Enabled        bool
	}
)

var (
	robot *Robot
)

func AddControllerActions(ctrl *Controller.Controller) {
	ctrl.AddAction(Controller.B, &Command.Command{
		Required:   "button b pressed",
		FirstRun:   true,
		Name:       "button b input",
		Initialize: func() {},
		Execute: func(required any) bool {
			req, ok := required.(string)
			if ok {
				fmt.Println(req)
			}
			return true
		},
	}).WhileTrue()
}

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

	return &Robot{
		Scheduler:      scheduler,
		DriveSubsystem: drive,
		Controllers:    controllers,
		Enabled:        false,
	}
}

func (r *Robot) Start() {
	r.Scheduler.Start()
}

func GetRobot() *Robot {
	return robot
}

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
		Controllers    []Controller.Controller
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

	return &Robot{
		Scheduler:      scheduler,
		DriveSubsystem: drive,
		// Controllers:
	}
}

func (r *Robot) Start() {
	r.Scheduler.Start()
}

func GetRobot() *Robot {
	return robot
}

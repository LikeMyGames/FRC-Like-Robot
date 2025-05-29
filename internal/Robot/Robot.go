package Robot

import (
	"fmt"
	"frcrobot/internal/Command"
	"frcrobot/internal/Controller"
	"frcrobot/internal/DriveSubsystem"
	"frcrobot/internal/GUI"
	"frcrobot/internal/Utils/MathUtils"
	"math"
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
	// Axis commands
	ctrl.AddAction(Controller.LeftStick, &Command.Command{
		Required: struct {
			Ctrl *Controller.Controller
		}{
			Ctrl: ctrl,
		},
		Name:       "drive wheels",
		FirstRun:   true,
		End:        false,
		Initialize: func() {},
		Execute: func(required any) bool {
			req, ok := required.(struct {
				Ctrl *Controller.Controller
			})
			if ok {
				type Axis struct {
					X float64
					Y float64
				}
				axis := Axis{X: float64(req.State.Gamepad.ThumbLX), Y: float64(req.State.Gamepad.ThumbLY)}

				pres := 2

				axis.X = math.Round(MathUtils.MapRange(axis.X, -32768, 32768, -1, 1)*math.Pow10(pres)) / math.Pow10(pres)
				axis.Y = math.Round(MathUtils.MapRange(axis.Y, -32768, 32768, -1, 1)*math.Pow10(pres)) / math.Pow10(pres)

				fmt.Println(axis)
			}
			return true
		},
	})

	// Button commands
	ctrl.AddAction(Controller.B, &Command.Command{
		Required:   "button b pressed",
		FirstRun:   true,
		Name:       "button b input",
		End:        false,
		Initialize: func() {},
		Execute: func(required any) bool {
			req, ok := required.(string)
			if ok {
				fmt.Println(req)
				go GUI.Success(req)
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

	robot = &Robot{
		Scheduler:      scheduler,
		DriveSubsystem: drive,
		Controllers:    controllers,
		Enabled:        false,
	}

	for i := range robot.Controllers {
		AddControllerActions(robot.Controllers[i])
	}

	return robot
}

func (r *Robot) Start() {
	r.Scheduler.Start()
}

func GetRobot() *Robot {
	return robot
}

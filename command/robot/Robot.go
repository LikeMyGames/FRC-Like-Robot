package robot

import (
	"fmt"
	"frcrobot/command"
	"frcrobot/controller"
	"frcrobot/drive"
	"frcrobot/gui"
	"frcrobot/utils/MathUtils"
	"frcrobot/utils/Types"
	"math"
)

type (
	Robot struct {
		DriveSubsystem *drive.SwerveDrive
		Controllers    []*controller.Controller
		Scheduler      *command.CommandScheduler
		// HardwareConn   Hardware.Conn
		Enabled bool
	}
)

var (
	robot *Robot
)

func AddControllerActions(ctrl *controller.Controller) {
	// Axis commands
	ctrl.AddAction(controller.LeftStick, &command.Command{
		Required: struct {
			Ctrl  *controller.Controller
			Robot *Robot
		}{
			Ctrl:  ctrl,
			Robot: robot,
		},
		Name:       "drive wheels",
		FirstRun:   true,
		End:        false,
		Initialize: func() {},
		Execute: func(required any) bool {
			req, ok := required.(struct {
				Ctrl  *controller.Controller
				Robot *Robot
			})
			if ok {

				axis := Types.Axis{X: float64(req.Ctrl.State.ThumbLX), Y: float64(req.Ctrl.State.ThumbLY)}

				pres := 2

				axis.X = math.Round(MathUtils.MapRange(axis.X, -32768, 32768, -1, 1)*math.Pow10(pres)) / math.Pow10(pres)
				axis.Y = math.Round(MathUtils.MapRange(axis.Y, -32768, 32768, -1, 1)*math.Pow10(pres)) / math.Pow10(pres)

				// req.Robot.drive.DriveToRelativePose()
				fmt.Println(axis)
			}
			return true
		},
	})

	// Button commands
	ctrl.AddAction(controller.B, &command.Command{
		Required:   "button b pressed",
		FirstRun:   true,
		Name:       "button b input",
		End:        false,
		Initialize: func() {},
		Execute: func(required any) bool {
			req, ok := required.(string)
			if ok {
				fmt.Println(req)
				go gui.Success(req)
			}
			return true
		},
	}).WhileTrue()
}

func NewRobot(controllerID []uint) *Robot {
	// Create a new scheduler for the robot
	scheduler := command.NewCommandScheduler()

	// Initialize the drive subsystem
	drive := drive.NewSwerveDrive(scheduler.Interval)

	controllers := make([]*controller.Controller, len(controllerID))

	for i, v := range controllerID {
		controllers[i] = controller.StartController(v, scheduler)
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

	// robot.Scheduler.ScheduleCommand(&command.Command{
	// 	Required:   nil,
	// 	Name:       "connect to hardware interface",
	// 	FirstRun:   true,
	// 	Initialize: func() {},
	// 	Execute: func(a any) bool {
	// 		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8765", nil)
	// 		if err != nil {
	// 			log.Println(err)
	// 			return false
	// 		}
	// 		conn.WriteMessage(websocket.TextMessage, []byte("robot connected"))
	// 		_, data, err := conn.ReadMessage()
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		log.Println(string(data))
	// 		robot.HardwareConn = conn
	// 		return true
	// 	},
	// })

	return robot
}

func (r *Robot) Start() {
	r.Scheduler.Start()
}

func GetRobot() *Robot {
	return robot
}

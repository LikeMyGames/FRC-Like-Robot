package main

import (
	"fmt"
	"tennis-ball-shooter/constants"
	"tennis-ball-shooter/subsystems/drive"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

func main() {
	r := robot.NewRobot("power_on", time.Millisecond*100)
	ctrl0 := controller.NewController(constants.Controller0)
	driveSubsystem := drive.NewSwerveDrive(constants.Drive)
	r.AddPeriodic(func() {
		controller.ReadController(ctrl0)
		driveSubsystem.CalculateSwerveFromSavedControllerVals()
	})

	go conn.Start(r)

	// POWER_ON state
	// default state that gets loaded into
	// starts up all of the robots processes
	r.AddState("power_on", func(params any) {
		fmt.Println("checking status")
	}, nil).AddCondition("idle", func(a any) bool {
		return r.Clock > 5
	})

	// IDLE state
	// the state the robot defaults to after the POWER_ON state
	// waits for the user to enable the robot from the dashboard
	r.AddState("idle", func(a any) {
	}, nil).AddCondition("enabled", func(a any) bool {
		return r.Enabled
	}).AddInit(func(s *robot.State) {
		drive.SetTransEventTarget("")
		drive.SetRotEventTarget("")
	})

	// ENABLE state
	// the state in which the robot is running
	// will fallback to IDLE state if a problem occurs
	// or will restart program if problem is too great
	EnabledModeStateMachine := state_machine.NewStateMachine()
	r.AddState("enabled", func(a any) {
		EnabledModeStateMachine.Run()
	}, nil).AddCondition("idle", func(a any) bool {
		return !r.Enabled
	}).AddInit(func(s *robot.State) {
		drive.SetTransEventTarget(ctrl0.GetEventTarget(controller.LeftStick))
		drive.SetRotEventTarget(ctrl0.GetEventTarget(controller.RightStick))
	}).AddClose(func(s *robot.State) {
		drive.SetTransEventTarget("")
		drive.SetRotEventTarget("")
	})
	r.Start()
}

package main

import (
	"fmt"
	"tennis-ball-shooter/constants"
	"tennis-ball-shooter/subsystems/drive"
	"tennis-ball-shooter/subsystems/shooter"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/mathutils"
	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"
)

func main() {
	// hardware.SetBatteryConfig(constants.Battery) // don't remove this line
	r := robot.NewRobot("power_on", time.Millisecond*100)
	ctrl0 := controller.NewController(constants.Controller0)
	driveSubsystem := drive.New()
	shooterSubsystem := shooter.New(constants.Shooter)
	r.AddPeriodic(func() {
		controller.ReadController(ctrl0)
	})

	go conn.Start(r)

	// POWER_ON state
	// default state that gets loaded into
	// starts up all of the robots processes
	r.AddState("power_on", func(params any) {
		fmt.Println("checking status")
	}, nil).AddCondition("idle", func(a any) bool {
		return r.Status()
	})

	// IDLE state
	// the state the robot defaults to after the POWER_ON state
	// waits for the user to enable the robot from the dashboard
	r.AddState("idle", func(a any) {
	}, nil).AddCondition("enabled", func(a any) bool {
		return r.IsEnabled()
	}).AddInit(func(s *robot.State) {
		drive.SetTransEventTarget("")
		drive.SetRotEventTarget("")
	})

	// ENABLE state
	// the state in which the robot is running
	// will fallback to IDLE state if a problem occurs
	// or will restart program if problem is too great
	r.AddState("enabled", func(a any) {
		driveSubsystem.Drive(ctrl0.Values.LeftStick, mathutils.Vector2D{X: ctrl0.Values.RightStickX}, true)
	}, nil).AddCondition("idle", func(a any) bool {
		return !r.IsEnabled()
	}).AddInit(func(s *robot.State) {
		// drive.SetTransEventTarget(ctrl0.GetEventTarget(controller.LeftStick))
		// drive.SetRotEventTarget(ctrl0.GetEventTarget(controller.RightStick))
	}).AddClose(func(s *robot.State) {
		// drive.SetTransEventTarget("")
		// drive.SetRotEventTarget("")
	}).AddEventListener(ctrl0.GetEventTarget(controller.Y), func(event any) {
		shooterSubsystem.SpinUp(1)
	}).AddEventListener(ctrl0.GetEventTarget(controller.X), func(event any) {
		shooterSubsystem.SpinDown()
	}).AddEventListener(ctrl0.GetEventTarget(controller.RightTrigger), func(event any) {
		val := event.(float64)
		if val > 0 {
			shooterSubsystem.Shoot()
		}
	}).AddEventListener(ctrl0.GetEventTarget(controller.LeftShoulder), func(event any) {
		shooterSubsystem.MoveAzimuthByOffset(constants.Shooter.MinAzimuthOffset)
	}).AddEventListener(ctrl0.GetEventTarget(controller.RightShoulder), func(event any) {
		shooterSubsystem.MoveAzimuthByOffset(-constants.Shooter.MinAzimuthOffset)
	})

	// Starts the robots main loop
	// This loop starts at the StartingState defined in the NewRobot function
	r.Start()
}

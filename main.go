package main

import (
	"fmt"
	"frcrobot/internal/Command"
	"frcrobot/internal/Controller"
	"frcrobot/internal/DriveSubsystem"
	"frcrobot/internal/GUI"
)

func main() {
	go GUI.StartUI()
	drive := DriveSubsystem.NewSwerveDrive("robot.constants")
	fmt.Println(drive)
	scheduler := Command.NewCommandScheduler()
	Controller.StartController(0, scheduler)

	//make sure to keep this call last
	scheduler.Start()
}

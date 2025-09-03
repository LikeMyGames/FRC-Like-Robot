package main

import (
	"frcrobot/gui"
	"frcrobot/hardware"
	Robot "frcrobot/robot"
)

func main() {
	hardware.Hello()
	go gui.StartUI()
	// if runtime.GOOS == "linux" {
	// 	exec.Command("")
	// }
	robot := Robot.NewRobot([]uint{0})

	robot.Start()
}

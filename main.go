package main

import (
	"frcrobot/internal/GUI"
	"frcrobot/internal/Hardware"
	"frcrobot/internal/Robot"
)

func main() {
	Hardware.Hello()
	go GUI.StartUI()
	// if runtime.GOOS == "linux" {
	// 	exec.Command("")
	// }
	robot := Robot.NewRobot([]uint{0})

	robot.Start()
}

package main

import (
	"frcrobot/internal/GUI"
	"frcrobot/internal/Robot"
)

func main() {
	go GUI.StartUI()
	// if runtime.GOOS == "linux" {
	// 	exec.Command("")
	// }
	robot := Robot.NewRobot([]uint{0})

	robot.Start()
}

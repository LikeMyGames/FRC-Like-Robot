package main

import (
	"frcrobot/internal/GUI"
	"frcrobot/internal/Robot"
)

func main() {
	go GUI.StartUI()
	robot := Robot.NewRobot([]uint{0})

	robot.Start()
}

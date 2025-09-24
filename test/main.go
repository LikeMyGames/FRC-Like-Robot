package main

import (
	"fmt"
	"test/constants"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"
)

func main() {
	r := robot.NewRobot("power_on")

	conn.Start(r)
	fmt.Println(constants.Drive)

	r.AddState("power_on", func(a any) {
		fmt.Println("checking status")
	}, nil).AddCondition("idle", func(a any) bool {
		return true
	})

	r.AddState("idle", func(a any) {
		fmt.Println("robot idling")
	}, nil).AddCondition("enabled", func(a any) bool {
		return r.Enabled
	})

	r.AddState("enabled", func(a any) {
		fmt.Println("robot enabled")
	}, nil).AddCondition("idle", func(a any) bool {
		return r.Enabled
	})

	r.Start()
}

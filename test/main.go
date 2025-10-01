package main

import (
	"fmt"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"
)

func main() {
	r := robot.NewRobot("power_on")

	go conn.Start(r)

	r.AddState("power_on", func(params any) {
		fmt.Println("checking status")
	}, nil).AddCondition("idle", func(a any) bool {
		if r.Clock > 5 {
			return true
		}
		return false
	})

	r.AddState("idle", func(a any) {
		fmt.Println("robot idling")
	}, nil).AddCondition("enabled", func(a any) bool {
		return r.Enabled
	})

	r.AddState("enabled", func(a any) {
		fmt.Println("robot enabled")
	}, nil).AddCondition("idle", func(a any) bool {
		return !r.Enabled
	})

	r.Start()
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println(os.Args)
	cmd := exec.Command(fmt.Sprintf("./%s", os.Args[1]))

	fmt.Println("Starting robot code")
	for {
		cmd.Start()
		t := time.NewTicker(time.Second)
		for range t.C {
			if cmd.ProcessState.ExitCode() != -1 {
				break
			}
		}
		time.Sleep(time.Second * 1000)
		fmt.Println("Restarting robot code")
	}
}

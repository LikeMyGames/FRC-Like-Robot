package watchdog

import (
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	freq, err := strconv.Atoi(os.Args[1])
	if err != nil {
		freq = 10
	}
	t := time.NewTicker(time.Duration(freq) * time.Millisecond)

	running := false
	var cmd *exec.Cmd
	for range t.C {
		if checkRunning() {
			startRobot()
		}
	}
}

func startRobot() *exec.Cmd {
	cmd := exec.Command("./main.exe")
	cmd.Run()
	return cmd
}

func checkRunning(cmd *exec.Cmd) bool {

}

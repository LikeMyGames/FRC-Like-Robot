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

	for range t.C {

	}
}

func startRobot() {
	cmd := exec.Command("./main.exe")
	cmd.Run()
}

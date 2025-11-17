package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :5000")
	changingRobotExe := false
	exeStarted := false

	cmd := exec.Command("./robot.exe")

	runCommand := func(cmd *exec.Cmd) {
		exeStarted = true
		fmt.Println("Starting robot.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		out, err := cmd.CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			fmt.Println(err)
		}
	}

	go func(cmd *exec.Cmd) {
		t := time.NewTicker(time.Millisecond * 5000)

		runCommand(cmd)

		for range t.C {
			if exeStarted && !changingRobotExe {
				runCommand(cmd)
			}
		}
	}(cmd)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		changingRobotExe = true
		exeStarted = false

		// Kill robot synchronously, not in a goroutine.
		if cmd.Process != nil {
			err = cmd.Process.Kill()
			if err != nil {
				panic(err)
			}
			err = cmd.Wait()
			if err != nil {
				panic(err)
			}
		}

		// Receive file fully
		handleConnection(conn)

		// Recreate the command after overwriting robot.exe
		cmd = exec.Command("./robot.exe")

		changingRobotExe = false

		// Start robot.exe again
		go runCommand(cmd)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// In a real application, you'd receive metadata first (filename, size)
	// For simplicity, let's assume a fixed filename for now.
	fileName := "robot.exe"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	bytesCopied, err := io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error copying data:", err)
		return
	}

	exec.Command("chmod", "+x", "robot.exe")

	fmt.Printf("Received %d bytes and saved to %s\n", bytesCopied, fileName)
}

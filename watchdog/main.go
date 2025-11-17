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
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :8080")
	changingRobotExe := false
	var cmd *exec.Cmd = nil
	go WatchRobotExe(&changingRobotExe, cmd)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		changingRobotExe = true
		if cmd != nil {
			cmd.Process.Kill()
		}
		handleConnection(conn)
		changingRobotExe = false
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

	bytesCopied, err := io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error copying data:", err)
		return
	}
	file.Close()

	exec.Command("chmod", "+x", "robot.exe")

	fmt.Printf("Received %d bytes and saved to %s\n", bytesCopied, fileName)
}

func WatchRobotExe(changing *bool, cmd *exec.Cmd) {
	cmd = exec.Command("./robot.exe")

	for {
		if !*changing && cmd.ProcessState.Exited() {
			fmt.Println("Starting robot.exe")
			out, err := cmd.CombinedOutput()
			fmt.Println(out)
			if err != nil {

			}
		} else {
			time.Sleep(time.Millisecond * 1500)
		}
	}
}

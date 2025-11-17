package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

var cmd *exec.Cmd
var running bool

func startRobot() {
	cmd = exec.Command("./robot.exe")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting robot.exe:", err)
		return
	}

	running = true

	// optional watcher
	go func() {
		cmd.Wait()
		running = false
		fmt.Println("robot.exe exited")
	}()
}

func stopRobot() {
	if !running || cmd == nil || cmd.Process == nil {
		return
	}
	fmt.Println("Stopping robot.exe…")

	// Kill the process
	cmd.Process.Kill()

	// VERY IMPORTANT: Wait must be called to clean up
	cmd.Wait()

	running = false
}

func main() {
	listener, _ := net.Listen("tcp", ":5000")
	fmt.Println("Listening for updates...")

	startRobot() // start initially

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		fmt.Println("Incoming update!")

		// 1. STOP old process
		stopRobot()

		// 2. RECEIVE new file
		if err := receiveFile(conn, "robot.exe.tmp"); err != nil {
			fmt.Println("Error receiving file:", err)
			conn.Close()
			continue
		}
		conn.Close()

		// 3. ATOMIC SWAP: replace old exe
		os.Remove("robot.exe")
		os.Rename("robot.exe.tmp", "robot.exe")

		// 3.5. MAKE IT EXECUTABLE
		os.Chmod("robot.exe", 0755)

		// 4. RESTART new robot.exe
		startRobot()
	}

}

func receiveFile(conn net.Conn, tempName string) error {
	file, err := os.Create(tempName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy until sender closes connection
	_, err = io.Copy(file, conn)
	return err
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

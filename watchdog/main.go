package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
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
	fmt.Printf("Received %d bytes and saved to %s\n", bytesCopied, fileName)
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tajtiattila/xinput"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// return true // Accept connections from any origin
		return true
	},
}

func main() {
	http.HandleFunc("/", handleWebSocket)
	http.ListenAndServe(":3030", nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Read message from the client
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}
	fmt.Printf("Received: %s\n", msg)

	// Send a message back to the client
	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"system_logger":{"type":"success","message":"controller connected"}}`))
	if err != nil {
		fmt.Println("Error writing message:", err)
		return
	}

	xinput.Load()
	state := &xinput.State{}
	for {
		time.Sleep(time.Millisecond * 20)
		err := xinput.GetState(0, state)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(`{"system_logger":{"type":"warn","message":"controller not connected"}}`))
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"controller":{"ctrlID":0,"buttons":%v,"triggerL":%v,"triggerR":%v,"thumbLX":%v,"thumbLY":%v,"thumbRX":%v,"thumbRY":%v}}`, state.Gamepad.Buttons, state.Gamepad.LeftTrigger, state.Gamepad.RightTrigger, state.Gamepad.ThumbLX, state.Gamepad.ThumbLY, state.Gamepad.ThumbRX, state.Gamepad.ThumbRY)))
		}
	}
}

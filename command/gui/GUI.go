package gui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// return true // Accept connections from any origin
		return true
	},
}

var connection any

type (
	WebSocketData struct {
		SystemLog       *Logger          `json:"system_logger"`
		RobotStatus     *Status          `json:"robot_status"`
		RunningSettings *RunSettings     `json:"running_settings"`
		Controller      *ControllerState `json:"controller"`
	}

	Logger struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}

	Status struct {
		Comms bool    `json:"comms"`
		Code  bool    `json:"code"`
		Joy   bool    `json:"joy"`
		Msg   string  `json:"message"`
		BatP  uint    `json:"bat_p"`
		BatV  float64 `json:"bat_v"`
	}

	RunSettings struct {
		Enabled bool   `json:"enabled"`
		Mode    string `json:"running_mode"`
	}

	ControllerState struct {
		ControllerID uint   `json:"ctrlID"`
		Buttons      uint16 `json:"buttons"`
		TriggerL     uint8  `json:"triggerL"`
		TriggerR     uint8  `json:"triggerR"`
		ThumbLX      int16  `json:"thumbLX"`
		ThumbLY      int16  `json:"thumbLY"`
		ThumbRX      int16  `json:"thumbRX"`
		ThumbRY      int16  `json:"thumbRY"`
	}
)

var LastMessage *WebSocketData

var LastControllerState *ControllerState = &ControllerState{
	ControllerID: 0,
	Buttons:      0,
	TriggerL:     0,
	TriggerR:     0,
	ThumbLX:      0,
	ThumbLY:      0,
	ThumbRX:      0,
	ThumbRY:      0,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Host)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	connection = conn
	// defer conn.Close()

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("received: %s\n", p)

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"system_logger":{"type":"success","message":"Backend Connected"},"robot_status":{"comms":true,"code":true}}`))
	if err != nil {
		log.Println(err)
		return
	}

	data := &WebSocketData{}
	for {
		err := conn.ReadJSON(data)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if data.Controller != nil {
			LastControllerState = data.Controller
		} else {

		}
	}
}

func StartUI() {
	http.HandleFunc("/", handleWebSocket)

	log.Println("Server started on localhost:8080")
	go log.Fatal(http.ListenAndServe(":8080", nil))
}

func SendData(data []byte) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		log.Println("connection not established with GUI")
		return
	}
	err := ws.WriteMessage(1, data)
	if err != nil {
		log.Println(err)
	}
}

func SendJSONData(v any) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		log.Println("connection not established with GUI")
		return
	}
	err := ws.WriteJSON(v)
	if err != nil {
		log.Println(err)
	}
}

func Log(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, fmt.Append(nil, `{"type":"log","message":"`+data+`"}`))
}

func Comment(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, fmt.Append(nil, `{"type":"comment","message":"`+data+`"}`))
}

func Success(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		log.Println("connection not established with GUI")
		return
	}
	fmt.Println(ws.LocalAddr())
	ws.WriteMessage(1, []byte(fmt.Sprint(`{"type":"success","message":"`+data+`"}`)))
}

func Warn(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, fmt.Append(nil, `{"type":"warn","message":"`+data+`"}`))
}

func Error(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, fmt.Append(nil, `{"type":"error","message":"`+data+`"}`))
}

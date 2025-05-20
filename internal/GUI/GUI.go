package GUI

import (
	"encoding/json"
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
		SystemLog   *Logger          `json:"system_logger"`
		RobotStatus *Status          `json:"robot_status"`
		Controller  *ControllerState `json:"controller"`
	}

	Logger struct {
		Type    string `json:"type,omitempty"`
		Message string `json:"message,omitempty"`
	}

	Status struct {
		Comms bool    `json:"comms,omitempty"`
		Code  bool    `json:"code,omitempty"`
		Joy   bool    `json:"joy,omitempty"`
		Msg   string  `json:"message,omitempty"`
		BatP  uint    `json:"bat_p,omitempty"`
		BatV  float64 `json:"bat_v,omitempty"`
	}

	ControllerState struct {
		Buttons  uint16 `json:"buttons,omitempty"`
		TriggerL uint8  `json:"triggerL,omitempty"`
		TriggerR uint8  `json:"triggerR,omitempty"`
		ThumbLX  uint8  `json:"thumbLX,omitempty"`
		ThumbLY  uint8  `json:"thumbLY,omitempty"`
		ThumbRX  uint8  `json:"thumbRX,omitempty"`
		ThumbRY  uint8  `json:"thumbRY,omitempty"`
	}
)

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

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"system_logger":{"type":"success","message":"Backend Connected"}}`))
	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log.Printf("recv: %s", msg)
		data := &WebSocketData{}
		err = json.Unmarshal(msg, data)
		if err != nil {
			log.Fatal(err)
		}
		if data.Controller != nil {
			Log(fmt.Sprintf("%v", data.Controller))
		}

		// err = conn.WriteMessage(mt, msg)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
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
		// log.Println("connection not established with GUI")
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
		// log.Println("connection not established with GUI")
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
	ws.WriteMessage(1, []byte(`{"type":"log","`+data+`"}`))
}

func Comment(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, []byte(`{"type":"comment","`+data+`"}`))
}

func Success(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, []byte(`{"type":"success","`+data+`"}`))
}

func Warn(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, []byte(`{"type":"warn","`+data+`"}`))
}

func Error(data string) {
	ws, ok := connection.(*websocket.Conn)
	if !ok {
		return
	}
	ws.WriteMessage(1, []byte(`{"type":"error","`+data+`"}`))
}

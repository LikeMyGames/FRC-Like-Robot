package conn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"

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

var connection *websocket.Conn

var bot *robot.Robot

type (
	WebSocketData struct {
		SystemLog   *Logger          `json:"system_logger"`
		RobotStatus *Status          `json:"robot_status"`
		RunSettings *RunSettings     `json:"run_settings"`
		Controller  *ControllerState `json:"controller"`
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
		Mode    string `json:"mode"`
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
	defer func() {
		conn.Close()
		connection.Close()
	}()

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

	for {
		_, data, err := connection.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		socketData := WebSocketData{}
		err = json.Unmarshal(data, &socketData)
		if err != nil {
			log.Fatal(err)
		}
		if socketData.Controller != nil {
			LastControllerState = socketData.Controller
		}
		if socketData.RunSettings != nil {
			bot.Enabled = socketData.RunSettings.Enabled
			bot.RunningMode = socketData.RunSettings.Mode
		}
	}
}

func Start(r *robot.Robot) {
	bot = r
	http.HandleFunc("/", handleWebSocket)

	log.Println("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SendData(data []byte) {
	err := connection.WriteMessage(1, data)
	if err != nil {
		log.Println(err)
	}
}

func SendJSONData(v any) {
	err := connection.WriteJSON(v)
	if err != nil {
		log.Println(err)
	}
}

func Log(data string) {
	connection.WriteMessage(1, fmt.Append(nil, `{"type":"log","message":"`+data+`"}`))
}

func Comment(data string) {
	connection.WriteMessage(1, fmt.Append(nil, `{"type":"comment","message":"`+data+`"}`))
}

func Success(data string) {
	fmt.Println(connection.LocalAddr())
	connection.WriteMessage(1, []byte(fmt.Sprint(`{"type":"success","message":"`+data+`"}`)))
}

func Warn(data string) {
	connection.WriteMessage(1, fmt.Append(nil, `{"type":"warn","message":"`+data+`"}`))
}

func Error(data string) {
	connection.WriteMessage(1, fmt.Append(nil, `{"type":"error","message":"`+data+`"}`))
}

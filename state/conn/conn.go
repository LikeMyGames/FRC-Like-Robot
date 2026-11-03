package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn/driver_station"
	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"

	"github.com/gorilla/websocket"
	nt4 "github.com/levifitzpatrick1/go-nt4"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// return true // Accept connections from any origin
			return true
		},
	}
	connection *websocket.Conn
	bot        *robot.Robot
	NT4        *nt4.Client
)

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
		BatP  float64 `json:"bat_p"`
		BatV  float64 `json:"bat_v"`
	}

	RunSettings struct {
		Enabled               bool   `json:"enabled"`
		Mode                  string `json:"mode"`
		RedAlliance           bool   `json:"alliance"`
		DriverStationPosition int    `json:"driverStationPosition"`
	}

	ControllerState struct {
		Id           uint   `json:"ctrlID"`
		Buttons      uint16 `json:"buttons"`
		LeftTrigger  uint8  `json:"triggerL"`
		RightTrigger uint8  `json:"triggerR"`
		ThumbLX      int16  `json:"thumbLX"`
		ThumbLY      int16  `json:"thumbLY"`
		ThumbRX      int16  `json:"thumbRX"`
		ThumbRY      int16  `json:"thumbRY"`
	}
)

var LastMessage *WebSocketData = new(WebSocketData)

// var LastControllerState *ControllerState = &ControllerState{
// 	ControllerID: 0,
// 	Buttons:      0,
// 	TriggerL:     0,
// 	TriggerR:     0,
// 	ThumbLX:      0,
// 	ThumbLY:      0,
// 	ThumbRX:      0,
// 	ThumbRY:      0,
// }

func OpenNT4Connection() {
	opts := nt4.DefaultClientOptions(nt4.TeamNumberToAddress(0))
	NT4 := nt4.NewClient(opts)
	if err := NT4.Connect(); err != nil {
		panic(err)
	}
}

func GetLatestNT4Update(updates chan nt4.TopicUpdate) (update nt4.TopicUpdate) {
	update = <-updates
	for v := range updates {
		if v.Timestamp > update.Timestamp {
			update = v
		}
	}

	return update
}

// func GetNT4Client() *nt4.Client {
// 	if nt4Client == nil {
// 		opts := nt4.DefaultClientOptions("")
// 		nt4Client = nt4.NewClient(opts)
// 		if err := nt4Client.Connect(); err != nil {
// 			panic(err)
// 		}
// 	}
// 	return nt4Client
// }

func Close() {
	if NT4 != nil {
		NT4.Disconnect()
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	defer func() {
		bot.Disable()
	}()
	fmt.Println(r.URL.Host)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	connection = conn
	defer func() {
		fmt.Println("closing websocket handler")
		conn.Close()
		// connection.Close()
	}()

	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("received: %s\n", p)
	// err = conn.WriteMessage(websocket.TextMessage, []byte(`{"system_logger":{"type":"success","message":"Backend Connected"},"robot_status":{"comms":true,"code":true,"message":"Robot Connected"}}`))

	SendJSONData(WebSocketData{
		SystemLog: &Logger{
			Type:    "success",
			Message: "Robot Connected",
		},
		RobotStatus: &Status{
			Comms: true,
			Code:  true,
			Msg:   "Robot Connected",
		},
	})

	logFile, err := os.Create("log_file.txt")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	os.Stdout = logFile

	LastMessage.RunSettings = new(RunSettings)
	LastMessage.RunSettings.Enabled = false
	LastMessage.RunSettings.Mode = ""
	LastMessage.RunSettings.RedAlliance = true
	LastMessage.RunSettings.DriverStationPosition = 0

	for {
		SendJSONData(WebSocketData{
			RobotStatus: &Status{
				Comms: true,
				Code:  true,
				BatP:  hardware.ReadBatteryPercentage(),
				BatV:  hardware.ReadBatteryVoltage(),
			},
			RunSettings: &RunSettings{
				Enabled:               bot.IsEnabled(),
				Mode:                  LastMessage.RunSettings.Mode,
				RedAlliance:           LastMessage.RunSettings.RedAlliance,
				DriverStationPosition: LastMessage.RunSettings.DriverStationPosition,
			},
		})
		_, data, err := connection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		socketData := WebSocketData{}
		err = json.Unmarshal(data, &socketData)
		if err != nil {
			log.Fatal(err)
			return
		}
		if socketData.Controller != nil {
			ctrl := socketData.Controller
			controller.AddControllerState(int(ctrl.Id), &controller.State{
				Buttons:      ctrl.Buttons,
				ThumbLX:      ctrl.ThumbLX,
				ThumbLY:      ctrl.ThumbLY,
				ThumbRX:      ctrl.ThumbRX,
				ThumbRY:      ctrl.ThumbRY,
				LeftTrigger:  ctrl.LeftTrigger,
				RightTrigger: ctrl.RightTrigger,
			})
		}
		if socketData.RunSettings != nil {
			if socketData.RunSettings.Enabled {
				slog.Info("robot enabled")
				bot.Enable()
			} else {
				slog.Info("robot disabled")
				bot.Disable()
			}
			bot.RunningMode = socketData.RunSettings.Mode
			driver_station.SetAlliance(socketData.RunSettings.RedAlliance)
			driver_station.SetDriverStationPosition(socketData.RunSettings.DriverStationPosition)
		}
	}
}

func Start(r *robot.Robot) {
	bot = r
	http.HandleFunc("/", handleWebSocket)

	robot.AddDeferFunction(func() { Close() })

	go func() {
		server := http.Server{
			Addr:    ":8080",
			Handler: nil,
		}

		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			signal.Stop(sigint)

			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			if err := server.Shutdown(shutdownCtx); err != nil {
				fmt.Println(err)
			}
			fmt.Println("server shutdown")
			cancel()
		}()

		for {
			fmt.Println("Server started on localhost:8080")
			err := server.ListenAndServe()
			if err != nil {
				fmt.Println(err)
				if err == http.ErrServerClosed {
					return
				}
			}
			fmt.Println("Driver Station disconnected")
		}
	}()
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

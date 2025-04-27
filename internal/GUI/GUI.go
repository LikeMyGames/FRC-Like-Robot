package GUI

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

	err = conn.WriteMessage(1, []byte(`{"system_logger":{"type":"success","message":"Backend Connected"}}`))
	if err != nil {
		log.Println(err)
		return
	}

}

// func handleAPI(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Print(err)
// 			return
// 		}
// 		fmt.Println("Received from frontend:", string(body))
// 		w.Write([]byte("Data received"))
// 	}
// }

func StartUI() {
	http.Handle("/", http.FileServer(http.Dir("./StaticGUI")))
	http.HandleFunc("/ws", handleWebSocket)

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

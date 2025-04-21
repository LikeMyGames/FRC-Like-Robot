package WebSocket

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Accept connections from any origin
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	messageType, p, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("received: %s\n", p)

	err = ws.WriteMessage(messageType, p)
	if err != nil {
		log.Println(err)
		return
	}

}

func StartUI() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Working Directory: ", dir)
	fileServer := http.FileServer(http.Dir("./internal/WebSocket/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	http.HandleFunc("/ws", handleConnections)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

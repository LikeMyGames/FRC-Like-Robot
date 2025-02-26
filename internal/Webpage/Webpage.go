package webpage

import (
	"encoding/json"
	"fmt"
	"internal/File"
	"log"
	"net/http"
)

var pressedButtons []string
var thumbL []float32
var thumbR []float32
var triggerL float32
var triggerR float32

type API_DATA struct {
	PressedButtons []string  `json:"pressedButtons"`
	ThumbL         []float32 `json:"thumbL"`
	ThumbR         []float32 `json:"thumbR"`
	TriggerL       float32   `json:"triggerL"`
	TriggerR       float32   `json:"triggerR"`
}

func Start() {
	router := http.NewServeMux()

	router.HandleFunc("GET /webpage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		fmt.Fprint(w, string(File.ReadBytes("./frontend/index.html")))
	})

	router.HandleFunc("GET /api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		data, err := json.Marshal(API_DATA{PressedButtons: pressedButtons, ThumbL: thumbL, ThumbR: thumbR, TriggerL: triggerL, TriggerR: triggerR})
		if err != nil {
			return
		}
		w.Write(data)
	})

	log.Printf("Serving running on localhost:%v\n", 5000)
	log.Fatal(http.ListenAndServe(":5000", router))
}

func Send(INthumbL, INthumbR []float32, INtriggerL, INtriggerR float32, INpressedButtons []string) {
	thumbL = INthumbL
	thumbR = INthumbR
	triggerL = INtriggerL
	triggerR = INtriggerR
	pressedButtons = INpressedButtons
}

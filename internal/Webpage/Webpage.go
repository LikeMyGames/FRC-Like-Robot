package webpage

import (
	"fmt"
	"internal/File"
	"log"
	"net/http"
)

func Start(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		fmt.Fprint(w, string(File.ReadBytes("./frontend/index.html")))
	})

	log.Printf("Serving running on localhost:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

// func SendVariables()

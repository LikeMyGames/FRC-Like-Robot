package JSON

import (
	"encoding/json"
	"log"
	"os"
)

func ReadJSON(filename string, jsonData interface{}) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

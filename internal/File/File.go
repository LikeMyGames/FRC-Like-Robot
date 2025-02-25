package File

import (
	"encoding/json"
	"log"
	"os"
)

func ReadJSON(filename string, jsonData any) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(data, jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func ReadBytes(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}
	return data
}

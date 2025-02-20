package main

import (
	"encoding/json"
	"log"
	"os"
)

func readJSON(filename string, jsonData interface{}) {
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

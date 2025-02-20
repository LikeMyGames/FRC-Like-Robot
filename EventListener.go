package main

import (
	"fmt"
)

var (
	eventListeners map[string]func(...any) any = map[string]func(...any) any{}
)

// type (
// 	eventCallback func(...any) (any)
// )

func on(name string, callback func(...any) any) {
	eventListeners[name] = callback
}

func trigger(name string, params ...any) any {
	if eventListeners[name] != nil {
		return eventListeners[name](params...)
	}
	return nil
}

func listListeners() {
	fmt.Println(eventListeners)
}

func removeListener(name string) {
	if eventListeners[name] != nil {
		delete(eventListeners, name)
	}
}

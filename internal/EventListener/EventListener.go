package Event

import (
	"fmt"
)

var (
	eventListeners map[string][]func(...any) any = map[string][]func(...any) any{}
)

type EventCallback struct {
	CallbackFunc func(...any) any
	Settings     EventCallbackSettings
}

type EventCallbackSettings struct {
}

// type (
// 	eventCallback func(...any) (any)
// )

func Listen(name string, callback func(...any) any) {
	eventListeners[name] = append(eventListeners[name], callback)
}

func Emit(name string, params ...any) {
	if eventListeners[name] != nil {
		for _, v := range eventListeners[name] {
			v(params...)
		}
	}
}

func ListListeners() {
	fmt.Println("\n-----------------------------------------------")
	fmt.Println("Event Listeners")
	fmt.Println("-----------------------------------------------")
	for i, v := range eventListeners {
		fmt.Println("Event: ", i, "\tFunctions: ", v)
		fmt.Println("-----------------------------------------------")
	}
	fmt.Println("")
}

func RemoveListener(name string) {
	if eventListeners[name] != nil {
		eventListeners[name] = eventListeners[name][:len(eventListeners[name])-1]
	}
}

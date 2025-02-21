package Event

import (
	"fmt"
)

var (
	eventListeners map[string]func(...any) any = map[string]func(...any) any{}
)

// type (
// 	eventCallback func(...any) (any)
// )

func Listen(name string, callback func(...any) any) {
	eventListeners[name] = callback
}

func Emit(name string, params ...any) any {
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

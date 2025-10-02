package event

import "fmt"

var listeners map[string]([]func(event any))

func Listen(target string, a func(event any)) {
	listeners[target] = append(listeners[target], a)
}

func Trigger(target string, event any) {
	if listeners[target] != nil {
		fmt.Println("No listeners with that target exist")
		return
	}
	for _, a := range listeners[target] {
		a(event)
	}
}

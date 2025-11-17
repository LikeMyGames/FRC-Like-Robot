package event

import (
	"fmt"
	"time"
)

type (
	Listener struct {
		from     string
		target   string
		id       uint64
		callback func(event any)
	}
)

var (
	listeners                = make(map[string][]*Listener)
	removingListeners        = false
	nextId            uint64 = 0
)

func Listen(target, from string, callback func(event any)) *Listener {
	if removingListeners {
		listener := &Listener{}
		go func() {
			time.Sleep(time.Millisecond * 10)
			listener = Listen(target, from, callback)
		}()
		return listener
	}
	listener := &Listener{from: from, target: target, callback: callback}
	listeners[target] = append(listeners[target], listener)
	return listener
}

func Trigger(target string, event any) {
	if removingListeners {
		go func() {
			time.Sleep(time.Millisecond * 10)
			Trigger(target, event)
		}()
		return
	}
	if listeners[target] == nil {
		return
	}
	fmt.Println("Triggering:", target)
	for _, a := range listeners[target] {
		a.callback(event)
	}
}

func Remove(listener *Listener) {
	removingListeners = true
	if listener == nil {
		removingListeners = false
		return
	}
	for i, v := range listeners[listener.target] {
		if v.from == listener.from {
			listeners[listener.target] = append(listeners[listener.target][:i], listeners[listener.target][i+1:]...)
		}
	}
	removingListeners = false
}

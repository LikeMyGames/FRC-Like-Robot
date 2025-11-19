package event

import (
	"time"
)

type (
	Listener struct {
		From     string
		Target   string
		id       uint64
		Callback func(event any)
	}
)

var (
	listeners                = make(map[string][]*Listener)
	removingListeners        = false
	nextId            uint64 = 1
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
	listener := &Listener{From: from, Target: target, Callback: callback, id: nextId}
	nextId++
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
	if listeners[target] == nil || len(listeners[target]) == 0 {
		return
	}

	for _, a := range listeners[target] {
		a.Callback(event)
	}
}

func Remove(listener *Listener) {
	removingListeners = true
	if listener == nil {
		removingListeners = false
		return
	}
	for i, v := range listeners[listener.Target] {
		if v.id == listener.id {
			listeners[listener.Target] = append(listeners[listener.Target][:i], listeners[listener.Target][i+1:]...)
		}
	}
	removingListeners = false
	listener.id = 0
	return
}

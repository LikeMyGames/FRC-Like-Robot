package event

type (
	Listener struct {
		from     string
		target   string
		callback func(event any)
	}
)

var listeners = make(map[string][]*Listener)

func Listen(target, from string, callback func(event any)) *Listener {
	listener := &Listener{from: from, target: target, callback: callback}
	listeners[target] = append(listeners[target], listener)
	return listener
}

func Trigger(target string, event any) {
	if listeners[target] == nil {
		return
	}
	for _, a := range listeners[target] {
		a.callback(event)
	}
}

func Remove(listener *Listener) {
	for i, v := range listeners[listener.target] {
		if v.from == listener.from {
			listeners[listener.target] = append(listeners[listener.target][:i], listeners[listener.target][i+1:]...)
		}
	}
}

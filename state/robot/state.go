package robot

import (
	"fmt"

	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
)

type (
	State struct {
		name       string
		action     func(any)
		switches   map[string]func(any) bool
		parameters any
		init       func(*State)
		close      func(*State)
		Listeners  []*event.Listener
	}
)

func (s *State) CheckCondition() *string {
	for i, v := range s.switches {
		if v(nil) {
			return &i
		}
	}
	return nil
}

func (s *State) AddCondition(target string, condition func(any) bool) *State {
	s.switches[target] = condition
	return s
}

func (s *State) AddEventListener(target string, callback func(event any)) *State {
	s.Listeners = append(s.Listeners, &event.Listener{Target: target, From: fmt.Sprintf("STATE_%s", s.name), Callback: callback})
	return s
}

func (s *State) AddInit(action func(*State)) *State {
	s.init = action
	return s
}

func (s *State) AddClose(action func(*State)) *State {
	s.close = action
	return s
}

func (s *State) loadEventListeners() {
	for i, v := range s.Listeners {
		s.Listeners[i] = event.Listen(v.Target, fmt.Sprintf("STATE_%s", s.name), v.Callback)
	}
}

func (s *State) unLoadEventListeners() {
	for _, v := range s.Listeners {
		event.Remove(v)
	}
}

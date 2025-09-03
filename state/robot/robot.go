package robot

import (
	"log"
	"time"
)

type (
	Robot struct {
		TeamNum   uint8
		Addr      string
		States    map[string]*State
		State     string
		Frequency time.Duration
		Enabled   bool
	}

	State struct {
		action     func(any)
		switches   map[string]func(any) bool
		parameters any
	}
)

func NewRobot(StartState string) *Robot {
	return &Robot{
		States:    LoadStates(),
		State:     StartState,
		Frequency: time.Millisecond * 1000,
	}
}

func (r *Robot) AddState(name string, action func(any), params any) *State {
	s := &State{
		action:     action,
		parameters: params,
	}
	return s
}

func (r *Robot) SetState(newState string) *State {
	r.State = newState
	return r.States[r.State]
}

func LoadStates() map[string]*State {
	return map[string]*State{
		"power_on": {
			action: func(any) {

			},
			switches: map[string]func(any) bool{
				"enabled": func(any) bool {
					return false
				},
			},
		},
	}
}

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

func (r *Robot) Start() {
	t := time.NewTicker(r.Frequency)

	for range t.C {
		s := r.States[r.State]
		if ns := s.CheckCondition(); ns != nil {
			r.SetState(*ns)
		}
		s.action(s.parameters)
		log.Println(r.State)
	}
}

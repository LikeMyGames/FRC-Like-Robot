package robot

import (
	"fmt"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
)

type (
	Robot struct {
		TeamNum     uint8
		Addr        string
		States      map[string]*State
		State       string
		Frequency   time.Duration
		Enabled     bool
		Clock       int64
		RunningMode string
		PeriodFuncs []func()
	}

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

var RobotRef *Robot = nil

func NewRobot(StartState string, freq time.Duration) *Robot {
	if RobotRef == nil {
		RobotRef = &Robot{
			States:    map[string]*State{},
			State:     StartState,
			Frequency: freq,
		}
	}
	return RobotRef
}

func (r *Robot) AddState(name string, action func(any), params any) *State {
	s := &State{
		name:       name,
		action:     action,
		parameters: params,
		switches:   map[string]func(any) bool{},
	}
	r.States[name] = s
	return s
}

func (r *Robot) SetState(newState string) *State {
	r.State = newState
	return r.States[r.State]
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

func (s *State) AddEventListener(target string, callback func(event any)) *State {
	s.Listeners = append(s.Listeners, event.Listen(target, fmt.Sprintf("STATE_%s", s.name), callback))
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

func (r *Robot) GetState() *State {
	return r.States[r.State]
}

func (r *Robot) AddPeriodic(a func()) *Robot {
	r.PeriodFuncs = append(r.PeriodFuncs, a)
	return r
}

func (r *Robot) Start() {
	t := time.NewTicker(r.Frequency)

	for range t.C {
		r.Clock++
		s := r.States[r.State]
		if ns := s.CheckCondition(); ns != nil {
			if r.GetState().close != nil {
				r.GetState().close(r.GetState())
			}
			fmt.Println("Switching to", *ns)
			r.SetState(*ns)
			if r.GetState().init != nil {
				r.GetState().init(r.GetState())
			}
			continue
		}
		for _, v := range r.PeriodFuncs {
			v()
		}
		s.action(s.parameters)
	}
}

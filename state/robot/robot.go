package robot

import (
	"fmt"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/can"
)

type (
	Robot struct {
		TeamNum     uint8
		Addr        string
		States      map[string]*State
		State       string
		Frequency   time.Duration
		enabled     bool
		Clock       int64
		RunningMode string
		PeriodFuncs []func()
		CanBus      *can.CanBus
		rsl         *hardware.Pin
	}
)

var RobotRef *Robot = nil

func NewRobot(StartState string, freq time.Duration) *Robot {
	if RobotRef == nil {
		RobotRef = &Robot{
			States:    map[string]*State{},
			State:     StartState,
			Frequency: freq,
			CanBus:    can.NewCanBus(),
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

func (r *Robot) Enable() {
	r.enabled = true
}

func (r *Robot) Disable() {
	r.enabled = false
}

func (r *Robot) IsEnabled() bool {
	return r.enabled
}

func (r *Robot) Status() bool {
	return r.CanBus.CheckStatuses()
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
		r.CanBus.UpdateDevices()
		if ns := s.CheckCondition(); ns != nil {
			if r.GetState().close != nil {
				r.GetState().unLoadEventListeners()
				r.GetState().close(r.GetState())
			}
			fmt.Println("Switching to", *ns)
			r.SetState(*ns)
			if r.GetState().init != nil {
				r.GetState().loadEventListeners()
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

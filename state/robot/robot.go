package robot

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/can"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/rsl"
)

type (
	Robot struct {
		TeamNum     uint8
		Addr        string
		States      map[string]*State
		State       string
		Frequency   time.Duration
		enabled     atomic.Bool
		Clock       int64
		RunningMode string
		PeriodFuncs []func()
		CanBus      *can.CanBus
		rsl         *rsl.RSL
	}
)

var RobotRef *Robot = nil
var DeferList []func() = make([]func(), 0)

func NewRobot(config constantTypes.RobotConfig) *Robot {
	hardware.OpenSpi()
	if RobotRef == nil {
		RobotRef = &Robot{
			States:    map[string]*State{},
			State:     config.StartingState,
			Frequency: config.Period,
			CanBus:    can.NewCanBus(),
			rsl:       rsl.New(config.RslPin),
		}
	}
	return RobotRef
}

func AddDeferFunction(function func()) {
	DeferList = append(DeferList, function)
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
	r.enabled.Store(true)
}

func (r *Robot) Disable() {
	r.enabled.Store(false)
}

func (r *Robot) IsEnabled() bool {
	return r.enabled.Load()
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

	event.Trigger("ROBOT_ENABLE_STATUS", false)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	// figure out why the pins aren't being closed properly
	defer func() {
		hardware.CloseAllPins()
		hardware.CloseSpiPort()
		for _, v := range DeferList {
			v()
		}
	}()

	for range t.C {
		r.Clock++
		r.rsl.SetEnabled(r.IsEnabled())
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

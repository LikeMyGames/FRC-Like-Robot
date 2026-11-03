package robot

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn/driver_station"
	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/can"
	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware/rsl"
	"github.com/LikeMyGames/FRC-Like-Robot/state/subsystem"
)

type (
	Robot struct {
		TeamNum     uint8
		Addr        string
		States      map[string]*State
		State       string
		Period      time.Duration
		enabled     atomic.Bool
		Clock       uint64
		MatchTime   uint64
		RunningMode string
		PeriodFuncs []func()
		CanBus      *can.CanBus
		rsl         *rsl.RSL
		subsystems  []*subsystem.Subsystem
	}
)

var RobotRef *Robot = nil
var DeferList []func() = make([]func(), 0)

func NewRobot(config constantTypes.RobotConfig) *Robot {
	// hardware.OpenSpi()
	if RobotRef == nil {
		RobotRef = &Robot{
			States: map[string]*State{},
			State:  config.StartingState,
			Period: config.Period,
			CanBus: can.NewCanBus(),
			rsl:    rsl.New(config.RslPin),
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
	can.BuildAndSendFrameFromBits(0, [8]byte{})
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

func (r *Robot) AddSubsystem(subsystem *subsystem.Subsystem) {
	r.subsystems = append(r.subsystems, subsystem)
}

func (r *Robot) Start() {
	t := time.NewTicker(r.Period)

	event.Trigger("ROBOT_ENABLE_STATUS", false)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// figure out why the pins aren't being closed properly
	deferFunc := func() {
		hardware.CloseAllPins()
		// hardware.CloseSpiPort()
		for _, v := range DeferList {
			v()
		}
		r.CanBus.Close()
		signal.Stop(quit)
	}

	defer deferFunc()

	for _, v := range r.subsystems {
		(*v).Initialize()
	}

	for range t.C {
		select {
		case <-quit:
			slog.Info("Stopping main robot loop")
			t.Stop()
			deferFunc()
			return
		default:
			r.MatchTime = r.Clock * uint64(r.Period.Milliseconds()) / 1000
			controller.ReadAllControllers()
			r.rsl.SetEnabled(r.IsEnabled())
			r.CanBus.UpdateDevices()
			s := r.States[r.State]
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
			}
			for _, v := range r.PeriodFuncs {
				v()
			}
			for _, v := range r.subsystems {
				(*v).Periodic()
			}
			s.action(s.parameters)
			can.SendSyncMessage(
				r.MatchTime,
				uint64(driver_station.GetMatchNumber()),
				uint64(driver_station.GetReplayNumber()),
				driver_station.GetAlliance().IsRed(),
				r.IsEnabled(),
				r.RunningMode == "Autonomous",
				r.RunningMode == "Test",
				true,
				uint64(driver_station.GetTournamentType()))
			r.Clock++
		}
	}
}

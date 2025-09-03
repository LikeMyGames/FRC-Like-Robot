package command

import (
	"frcrobot/gui"
	"log"
	"slices"
	"time"
)

type (
	CommandScheduler struct {
		Interval time.Duration
		Commands []*Command
	}

	Command struct {
		Initialize func()
		Execute    func(any) bool
		End        bool
		Required   any
		FirstRun   bool
		Name       string
	}
)

func NewCommandScheduler() *CommandScheduler {
	scheduler := &CommandScheduler{
		Interval: time.Millisecond * 100,
		Commands: make([]*Command, 0),
	}
	log.Println("Created Scheduler: ", &scheduler)
	return scheduler
}

func (scheduler *CommandScheduler) Start() {
	ticker := time.NewTicker(scheduler.Interval)
	gui.Success("Scheduler started")
	for range ticker.C {
		for i := len(scheduler.Commands) - 1; i >= 0; i-- {
			v := scheduler.Commands[i]
			if v == nil {
				continue
			}
			if !v.End {
				if v.FirstRun {
					v.Initialize()
					v.FirstRun = false
				}
				v.End = v.Execute(v.Required)
				if v.End {
					scheduler.Commands = slices.Delete(scheduler.Commands, i, i+1)
				}
			} else {
				scheduler.Commands = slices.Delete(scheduler.Commands, i, i+1)
			}
		}
	}
}

func (scheduler *CommandScheduler) ScheduleCommand(commands *Command) {
	scheduler.Commands = append(scheduler.Commands, commands)
}

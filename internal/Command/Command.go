package Command

import (
	"frcrobot/internal/GUI"
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
	GUI.Success("Scheduler started")
	for range ticker.C {
		for i, v := range scheduler.Commands {
			if v == nil {
				continue
			}
			if !v.End {
				// log.Println("running command: ", v.Name, &v)
				if v.FirstRun {
					v.Initialize()
					v.FirstRun = false
				}
				v.End = v.Execute(v.Required)
			} else {
				scheduler.Commands = slices.Delete(scheduler.Commands, i, i+1)
			}
		}
	}
}

func (scheduler *CommandScheduler) ScheduleCommand(commands *Command) {
	log.Println("Scheduled commands: ", commands.Name)
	scheduler.Commands = append(scheduler.Commands, commands)
}

// func (command *DefaultCommand) Initialize() {
// 	if command.Initialize != nil {
// 		command.InitFunc()
// 	}
// }

// func (command *DefaultCommand) Execute() {
// 	if command.Initialize != nil {
// 		command.ExecFunc()
// 	}
// }

// func (command *DefaultCommand) End() {
// 	if command.Initialize != nil {
// 		command.EndFunc()
// 	}
// }

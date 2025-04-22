package Command

import (
	"time"
)

type (
	CommandScheduler struct {
		Interval time.Duration
		Commands map[string]*Command
	}

	// Command struct {
	// 	CommandInterface
	// 	FirstRun bool
	// 	Name     string
	// 	Required any
	// }

	// CommandInterface interface {
	// 	Initialize()
	// 	Execute(any)
	// 	End() bool
	// 	getRequired() any
	// }

	Command struct {
		Initialize func()
		Execute    func(any)
		End        func() bool
		Required   any
		FirstRun   bool
		Name       string
	}
)

func NewCommandScheduler() *CommandScheduler {
	return &CommandScheduler{
		Interval: time.Second / 20,
		Commands: make(map[string]*Command),
	}
}

func (scheduler *CommandScheduler) Start() {
	ticker := time.NewTicker(scheduler.Interval)

	for range ticker.C {
		for _, v := range scheduler.Commands {
			// log.Println("Running Command: ", v.Name)
			if !v.End() {
				if v.FirstRun {
					v.Initialize()
					v.FirstRun = false
				}
				v.Execute(v.Required)
			} else {
				delete(scheduler.Commands, v.Name)
			}
		}
	}
}

func (scheduler *CommandScheduler) ScheduleCommand(commands ...*Command) {
	for _, v := range commands {
		scheduler.Commands[v.Name] = v
	}
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

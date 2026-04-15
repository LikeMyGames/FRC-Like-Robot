package state_machine

import "fmt"

type (
	StateMachine struct {
		States     map[string]*State
		State      string
		parameters any
	}

	// State struct {
	// 	name       string
	// 	action     func(any)
	// 	switches   map[string]func(any) bool
	// 	parameters any
	// 	init       func(*State)
	// 	close      func(*State)
	// }

	State interface {
		GetName() string
		Initialize()
		Execute()
		End()
		GetSwitches() map[string]func(any) bool
	}

	// State interface {
	// 	GetName() string
	// 	Action(any)
	// 	GetSwitches() map[string]func(any) bool
	// 	GetParameters() any
	// 	Init() *State
	// 	Close() *State
	// }
)

func NewStateMachine(states ...State) *StateMachine {
	machine := &StateMachine{
		States: map[string]*State{},
	}
	for _, s := range states {
		machine.States[s.GetName()] = &s
	}
	return machine
}

// func NewState(name string, action func(any), switches map[string]func(any) bool, parameters any, init func(*State), close func(*State)) *State {
// 	return &State{
// 		name:       name,
// 		action:     action,
// 		switches:   switches,
// 		parameters: parameters,
// 		init:       init,
// 		close:      close,
// 	}
// }

// func NewStateFromInterface(state StateInterface) {

// }

func (m *StateMachine) AddState(state State) *StateMachine {
	m.States[state.GetName()] = &state
	return m
}

func (m *StateMachine) Run() {
	if m.State == "" {
		return
	}
	s := m.States[m.State]
	if ns := checkStateCondititon(*s); ns != nil {
		(*s).End()
		fmt.Println("Switching to", *ns)
		m.State = *ns
		(*s).Initialize()
	}
	(*s).Execute()
}

func checkStateCondititon(s State) *string {
	for i, v := range s.GetSwitches() {
		if v(nil) {
			return &i
		}
	}
	return nil
}

package state_machine

import "fmt"

type (
	StateMachine struct {
		States map[string]*State
		State  string
	}

	State struct {
		name       string
		action     func(any)
		switches   map[string]func(any) bool
		parameters any
		init       func(*State)
		close      func(*State)
	}
)

func NewStateMachine(states ...*State) *StateMachine {
	machine := &StateMachine{
		States: map[string]*State{},
	}
	for _, s := range states {
		machine.States[s.name] = s
	}
	return machine
}

func NewState(name string, action func(any), switches map[string]func(any) bool, parameters any, init func(*State), close func(*State)) *State {
	return &State{
		name:       name,
		action:     action,
		switches:   switches,
		parameters: parameters,
		init:       init,
		close:      close,
	}
}

func (m *StateMachine) AddState(state *State) *StateMachine {
	m.States[state.name] = state
	return m
}

func (m *StateMachine) Run() {
	if m.State == "" {
		return
	}
	s := m.States[m.State]
	if ns := s.CheckCondition(); ns != nil {
		if s.close != nil {
			s.close(s)
		}
		fmt.Println("Switching to", *ns)
		m.State = *ns
		if s.init != nil {
			s.init(s)
		}
	}
	if s.action != nil {
		s.action(s.parameters)
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

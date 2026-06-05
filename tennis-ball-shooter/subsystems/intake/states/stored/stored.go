package stored

import intake_types "tennis-ball-shooter/subsystems/intake/types"

type Stored struct {
	name            string
	intakeSubsystem *intake_types.IntakeSubsystem
}

func Get(intakeSubsystem *intake_types.IntakeSubsystem) *Stored {
	s := new(Stored)
	s.name = "STORED"
	s.intakeSubsystem = intakeSubsystem
	return s
}

func (s *Stored) GetName() string {
	return s.name
}

func (s *Stored) Initialize() {

}

func (s *Stored) Execute() {

}

func (s *Stored) End() {

}

func (s *Stored) GetSwitches() map[string]func() bool {
	return map[string]func() bool{
		"EXTENDED": func() bool {
			return false
		},
		"AGITATING": func() bool {
			return false
		},
	}
}

package stored

import intake_types "tennis-ball-shooter/subsystems/intake/types"

type Stored struct {
	name            string
	intakeSubsystem *intake_types.IntakeSubsystem
}

func Get(intakeSubsystem *intake_types.IntakeSubsystem) *Stored {
	s := new(Stored)
	s.name = "INTAKE_STORED"
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

func (s *Stored) GetSwitches() map[string]func(any) bool {
	return map[string]func(any) bool{
		"EXTENDED": func(a any) bool {
			return false
		},
		"AGITATING": func(a any) bool {
			return false
		},
	}
}

package agitating

import intake_types "tennis-ball-shooter/subsystems/intake/types"

type Agitating struct {
	name            string
	intakeSubsystem *intake_types.IntakeSubsystem
}

func Get(intakeSubsystem *intake_types.IntakeSubsystem) *Agitating {
	s := new(Agitating)
	s.name = "INTAKE_AGITATING"
	s.intakeSubsystem = intakeSubsystem

	return s
}

func (s *Agitating) GetName() string {
	return s.name
}

func (s *Agitating) Initialize() {

}

func (s *Agitating) Execute() {

}

func (s *Agitating) End() {

}

func (s *Agitating) GetSwitches() map[string]func(any) bool {
	return map[string]func(any) bool{}
}

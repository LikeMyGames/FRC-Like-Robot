package extended

import intake_types "tennis-ball-shooter/subsystems/intake/types"

type Extended struct {
	name            string
	intakeSubsystem *intake_types.IntakeSubsystem
}

func Get(intakeSubsystem *intake_types.IntakeSubsystem) *Extended {
	s := new(Extended)
	s.name = "INTAKE_EXTENDED"
	s.intakeSubsystem = intakeSubsystem

	return s
}

func (s *Extended) GetName() string {
	return s.name
}

func (s *Extended) Initialize() {

}

func (s *Extended) Execute() {

}

func (s *Extended) End() {

}

func (s *Extended) GetSwitches() map[string]func(any) bool {

	return map[string]func(any) bool{}
}

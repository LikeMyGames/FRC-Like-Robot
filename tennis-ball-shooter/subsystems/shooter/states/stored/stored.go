package stored

import shooter_types "tennis-ball-shooter/subsystems/shooter/types"

type Stored struct {
	name             string
	shooterSubsystem *shooter_types.ShooterSubsystem
}

func Get(shooterSubsystem *shooter_types.ShooterSubsystem) *Stored {
	s := new(Stored)
	s.name = "STORED"
	s.shooterSubsystem = shooterSubsystem

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
		"SHOOTING": toShooting,
	}
}

func toShooting() bool {
	return false
}

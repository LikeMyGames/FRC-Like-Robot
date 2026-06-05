package shooting

import shooter_types "tennis-ball-shooter/subsystems/shooter/types"

type Shooting struct {
	name             string
	shooterSubsystem *shooter_types.ShooterSubsystem
}

func Get(shooterSubsystem *shooter_types.ShooterSubsystem) *Shooting {
	s := new(Shooting)
	s.name = "SHOOTING"
	s.shooterSubsystem = shooterSubsystem

	return s
}

func (s *Shooting) GetName() string {
	return s.name
}

func (s *Shooting) Initialize() {

}

func (s *Shooting) Execute() {

}

func (s *Shooting) End() {

}

func (s *Shooting) GetSwitches() map[string]func() bool {
	return map[string]func() bool{
		"STORED": toStored,
	}
}

func toStored() bool {
	return false
}

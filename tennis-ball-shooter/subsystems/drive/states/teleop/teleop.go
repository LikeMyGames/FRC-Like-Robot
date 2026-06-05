package teleop

import drive_types "tennis-ball-shooter/subsystems/drive/types"

type Teleop struct {
	name           string
	driveSubsystem *drive_types.DriveSubsystem
}

func Get(driveSubsystem *drive_types.DriveSubsystem) *Teleop {
	s := new(Teleop)
	s.name = "TELEOP"
	s.driveSubsystem = driveSubsystem

	return s
}

func (s *Teleop) GetName() string {
	return s.name
}

func (s *Teleop) Initialize() {

}

func (s *Teleop) Execute() {

}

func (s *Teleop) End() {

}

func (s *Teleop) GetSwitches() map[string]func() bool {
	return map[string]func() bool{
		"AUTO": toAuto,
	}
}

func toAuto() bool {
	return true
}

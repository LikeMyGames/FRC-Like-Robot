package auto

import drive_types "tennis-ball-shooter/subsystems/drive/types"

type Auto struct {
	name           string
	driveSubsystem *drive_types.DriveSubsystem
}

func Get(driveSubsystem *drive_types.DriveSubsystem) *Auto {
	s := new(Auto)
	s.name = "AUTO"
	s.driveSubsystem = driveSubsystem

	return s
}

func (s *Auto) GetName() string {
	return s.name
}

func (s *Auto) Initialize() {

}

func (s *Auto) Execute() {

}

func (s *Auto) End() {

}

func (s *Auto) GetSwitches() map[string]func() bool {
	return map[string]func() bool{
		"TELEOP": toTeleop,
	}
}

func toTeleop() bool {
	return true
}

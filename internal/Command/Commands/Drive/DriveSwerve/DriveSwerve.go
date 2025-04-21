package DriveSwerve

import (
	"frcrobot/internal/Command"
	"frcrobot/internal/DriveSubsystem"
)

type (
	DriveSwerve struct {
		Command Command.Command
		Drive   *DriveSubsystem.SwerveDrive
	}
)

func NewDriverSwerveCommand(drive *DriveSubsystem.SwerveDrive) *DriveSwerve {
	return &DriveSwerve{
		Command: Command.Command{Name: "DriveSwerve", FirstRun: true},
		Drive:   drive,
	}
}

func (d *DriveSwerve) Initialize() {

}

func (d *DriveSwerve) Execute() {
	// d.drive.CalculateSwerveModules()
}

func (d *DriveSwerve) End(interrupted bool) bool {
	return false
}

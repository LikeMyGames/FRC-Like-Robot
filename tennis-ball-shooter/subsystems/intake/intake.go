package intake

import (
	"tennis-ball-shooter/configs"
	"tennis-ball-shooter/subsystems/intake/states/agitating"
	"tennis-ball-shooter/subsystems/intake/states/extended"
	"tennis-ball-shooter/subsystems/intake/states/stored"
	intake_types "tennis-ball-shooter/subsystems/intake/types"

	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type IntakeSubsystem intake_types.IntakeSubsystem

func New() *intake_types.IntakeSubsystem {
	s := new(intake_types.IntakeSubsystem)

	s.RollerMotor = *motor.New(configs.IntakeMotors.RollerMotorConfig)
	s.ExtensionMotor = *motor.New(configs.IntakeMotors.ExtensionMotorConfig)

	s.StateMachine = state_machine.NewStateMachine()
	s.StateMachine.AddState(stored.Get(s))
	s.StateMachine.AddState(extended.Get(s))
	s.StateMachine.AddState(agitating.Get(s))

	return s
}

// func (s *IntakeSubsystem) Initialize() {
// 	// s.machine.AddState()
// }

func (s *IntakeSubsystem) Periodic() {
	s.StateMachine.Run()
}

package intake

import (
	"tennis-ball-shooter/configs"
	"tennis-ball-shooter/constants"
	"tennis-ball-shooter/subsystems/intake/states/agitating"
	"tennis-ball-shooter/subsystems/intake/states/extended"
	"tennis-ball-shooter/subsystems/intake/states/stored"
	intake_types "tennis-ball-shooter/subsystems/intake/types"

	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type IntakeSubsystem intake_types.IntakeSubsystem

var instance *IntakeSubsystem

func New() *IntakeSubsystem {
	s := new(IntakeSubsystem)

	s.RollerMotor = *motor.New(constants.Intake.RollerMotorCanId, configs.IntakeMotors.RollerMotorConfig)
	s.ExtensionMotor = *motor.New(constants.Intake.ExtensionMotorCanId, configs.IntakeMotors.ExtensionMotorConfig)

	s.StateMachine = state_machine.NewStateMachine()
	s.StateMachine.AddState(stored.Get(s.purify()))
	s.StateMachine.AddState(extended.Get(s.purify()))
	s.StateMachine.AddState(agitating.Get(s.purify()))

	instance = s

	return s
}

func GetInstance() *IntakeSubsystem {
	if instance != nil {
		return instance
	}
	return nil
}

func (s *IntakeSubsystem) purify() *intake_types.IntakeSubsystem {
	return (*intake_types.IntakeSubsystem)(s)
}

func (s *IntakeSubsystem) Initialize() {

}

func (s *IntakeSubsystem) Periodic() {
	s.StateMachine.Run()
}

func (s *IntakeSubsystem) SetState(stateName string) {

}

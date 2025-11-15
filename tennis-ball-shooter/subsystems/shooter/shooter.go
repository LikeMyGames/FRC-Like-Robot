package shooter

import (
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type (
	Shooter struct {
		config         shooter_types.ShooterConfig
		HasBall        bool
		ReadyToShoot   bool
		TopFlyWheel    *hardware.MotorController
		BottomFlyWheel *hardware.MotorController
		FeedWheel      *hardware.MotorController
	}
)

func New(config shooter_types.ShooterConfig) *Shooter {
	return &Shooter{
		config: config,
	}
}

// func (s *Shooter) NewStateMachine() {

// }

func (s *Shooter) SpinUp(speedPercent float64) {
	s.SpinUpTop(speedPercent)
	s.SpinUpBottom(speedPercent)
}

func (s *Shooter) SpinUpTop(speedPercent float64) {
	s.TopFlyWheel.SetTarget(s.config.MaxFlyWheelVelocity)
}

func (s *Shooter) SpinUpBottom(speedPercent float64) {
	s.BottomFlyWheel.SetTarget(s.config.MaxFlyWheelVelocity)
}

func (s *Shooter) Shoot() {
	s.SpinUp(1)
	s.FeedBall()
}

func (s *Shooter) FeedBall() {
	if s.HasBall {
		s.FeedWheel.SetTarget(0)
		return
	}
	s.FeedWheel.SetTarget(1)
}

func (s *Shooter) GetStates() []*state_machine.State {
	return []*state_machine.State{
		state_machine.NewState(
			"PREP_SHOOTER",
			func(a any) {

			},
			map[string]func(any) bool{},
			nil,
			func(st *state_machine.State) {

			},
			nil,
		),
		state_machine.NewState(
			"SHOOTING",
			nil,
			map[string]func(any) bool{},
			nil,
			nil,
			nil,
		),
	}
}

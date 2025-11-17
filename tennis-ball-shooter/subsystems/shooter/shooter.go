package shooter

import (
	"fmt"
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type (
	Shooter struct {
		config         shooter_types.ShooterConfig
		HasBall        bool
		ReadyToShoot   bool
		FlyWheelMotor  *hardware.MotorController
		PitchMotor     *hardware.MotorController
		FeedWheelMotor *hardware.MotorController
		AzimuthMotor   *hardware.MotorController
	}
)

func New(config shooter_types.ShooterConfig) *Shooter {
	return &Shooter{
		config:         config,
		HasBall:        false,
		ReadyToShoot:   false,
		FlyWheelMotor:  hardware.NewMotorController(config.FlyWheelMotor),
		PitchMotor:     hardware.NewMotorController(config.PitchMotor),
		FeedWheelMotor: hardware.NewMotorController(config.FeedWheelMotor),
		AzimuthMotor:   hardware.NewMotorController(config.AzimuthMotor),
	}
}

// func (s *Shooter) NewStateMachine() {

// }

func (s *Shooter) SpinUp(speedPercent float64) {
	s.FlyWheelMotor.SetTarget(s.config.MaxFlyWheelVelocity * speedPercent)
}

func (s *Shooter) Shoot() {
	fmt.Println("Shooting Ball")
	s.SpinUp(1)
	s.FeedBall()
}

func (s *Shooter) FeedBall() {
	if s.HasBall {
		s.FeedWheelMotor.SetTarget(0)
		return
	}
	s.FeedWheelMotor.SetTarget(1)
}

func (s *Shooter) MoveAzimuthByOffset(offset float64) {
	fmt.Println("Moving azimuth to", (s.AzimuthMotor.GetTarget() + offset))
	s.AzimuthMotor.SetTarget(s.AzimuthMotor.GetTarget() + offset)
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

package shooter

import (
	"fmt"
	"tennis-ball-shooter/constants"
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type (
	Shooter struct {
		config         shooter_types.ShooterConfig
		HasBall        bool
		ReadyToShoot   bool
		FlyWheelMotor  *motor.Motor
		PitchMotor     *motor.Motor
		FeedWheelMotor *motor.Motor
		AzimuthMotor   *motor.Motor
	}
)

func New(config shooter_types.ShooterConfig) *Shooter {
	return &Shooter{
		config:         config,
		HasBall:        false,
		ReadyToShoot:   false,
		FlyWheelMotor:  motor.New(int(config.FlyWheelMotor.Id)),
		PitchMotor:     motor.New(int(config.PitchMotor.Id)),
		FeedWheelMotor: motor.New(int(config.FeedWheelMotor.Id)),
		AzimuthMotor:   motor.New(int(config.AzimuthMotor.Id)),
	}
}

// func (s *Shooter) NewStateMachine() {

// }

func (s *Shooter) SpinUp(speed float64) {
	fmt.Println("Spinning up Shooter")
	s.FlyWheelMotor.SetVelocity(s.config.MaxFlyWheelVelocity * speed)
}

func (s *Shooter) SpinDown() {
	s.FlyWheelMotor.SetTorque(0)
}

func (s *Shooter) BrakeFlyWheel() {
	s.FlyWheelMotor.SetVelocity(0)
}

func (s *Shooter) Shoot() {
	fmt.Println("Feeding ball for Shooting")
	// s.SpinUp(1)
	s.feedBall()
}

func (s *Shooter) feedBall() {
	if s.HasBall {
		s.FeedWheelMotor.SetVelocity(0)
		return
	}
	s.FeedWheelMotor.SetVelocity(s.config.MaxFeedVelocity)
}

func (s *Shooter) MoveAzimuthByOffset(offset float64) {
	fmt.Println("Moving azimuth to", (s.AzimuthMotor.ReadAngle() + constants.Shooter.MinAzimuthOffset))
	s.AzimuthMotor.SetAngle(s.AzimuthMotor.ReadAngle() + offset)
}

func (s *Shooter) GetStates() []*state_machine.State {
	return []*state_machine.State{
		state_machine.NewState(
			"PREP_SHOOTER",
			func(a any) {

			},
			map[string]func(any) bool{
				"SHOOTING": func(a any) bool {
					return true
				},
			},
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

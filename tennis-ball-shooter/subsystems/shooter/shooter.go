package shooter

import (
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/hardware"
)

type (
	Shooter struct {
		config         shooter_types.ShooterConfig
		TopFlyWheel    *hardware.Device
		BottomFlyWheel *hardware.Device
	}
)

func New(config shooter_types.ShooterConfig) *Shooter {
	return &Shooter{
		config: config,
	}
}

func (s *Shooter) SpinUp() {
	s.SpinUpTop()
	s.SpinUpBottom()
}

func (s *Shooter) SpinUpTop() {
	s.TopFlyWheel.SetTargetValue(s.config.MaxFlyWheelVelocity / 2)
}

func (s *Shooter) SpinUpBottom() {
	s.BottomFlyWheel.SetTargetValue(s.config.MaxFlyWheelVelocity / 2)
}

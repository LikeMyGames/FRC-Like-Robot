package shooter

import (
	"tennis-ball-shooter/configs"
	"tennis-ball-shooter/constants"
	"tennis-ball-shooter/subsystems/drive"
	"tennis-ball-shooter/subsystems/shooter/states/shooting"
	"tennis-ball-shooter/subsystems/shooter/states/stored"
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/game/field"
	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
)

type ShooterSubsystem shooter_types.ShooterSubsystem

var instance *ShooterSubsystem

func New() *ShooterSubsystem {
	s := new(ShooterSubsystem)

	s.FlywheelMotor = motor.New(constants.Shooter.FlywheelMotorCanId, configs.FlywheelMotor())
	s.HoodMotor = motor.New(constants.Shooter.HoodMotorCanId, configs.HoodMotor())
	s.TurretMotor = motor.New(constants.Shooter.TurretMotorCanId, configs.TurretMotor())

	s.StateMachine = state_machine.NewStateMachine(
		shooting.Get(s.purify()),
		stored.Get(s.purify()),
	)

	instance = s

	return s
}

func GetInstance() *ShooterSubsystem {
	if instance != nil {
		return instance
	}
	return nil
}

func (s *ShooterSubsystem) purify() *shooter_types.ShooterSubsystem {
	return (*shooter_types.ShooterSubsystem)(s)
}

func (s *ShooterSubsystem) Initialize() {

}

func (s *ShooterSubsystem) Periodic() {
	s.StateMachine.Run()
}

func (s *ShooterSubsystem) AutoAim() {
	pose := drive.GetInstance().Pose
	target := field.GetCurrentAllianceHubPose()

	pose.Location.Subtract(&target.Location)

}

func (s *ShooterSubsystem) AutoRange() {

}

// func (s *ShooterSubsystem) SetState(stateName string) {

// }

// func (s *ShooterSubsystem) SpinUp(speed float64) {
// 	fmt.Println("Spinning up Shooter")
// 	s.FlyWheelMotor.SetVelocity(constants.Shooter.MaxFlyWheelVelocity * speed)
// }

// func (s *ShooterSubsystem) SpinDown() {
// 	s.FlyWheelMotor.SetTorque(0)
// }

// func (s *ShooterSubsystem) BrakeFlyWheel() {
// 	s.FlyWheelMotor.SetVelocity(0)
// }

// func (s *ShooterSubsystem) Shoot() {
// 	fmt.Println("Feeding ball for Shooting")
// 	// s.SpinUp(1)
// 	s.feedBall()
// }

// func (s *ShooterSubsystem) feedBall() {
// 	if s.HasBall {
// 		s.FeedWheelMotor.SetVelocity(0)
// 		return
// 	}
// 	s.FeedWheelMotor.SetVelocity(s.Config.MaxFeedVelocity)
// }

// func (s *ShooterSubsystem) MoveAzimuthByOffset(offset float64) {
// 	fmt.Println("Moving azimuth to", (s.AzimuthMotor.ReadAngle() + constants.Shooter.MinAzimuthOffset))
// 	s.AzimuthMotor.SetAngle(s.AzimuthMotor.ReadAngle() + offset)
// }

// func (s *ShooterSubsystem) GetStates() []*state_machine.State {
// 	return []*state_machine.State{
// 		state_machine.NewState(
// 			"PREP_SHOOTER",
// 			func(a any) {

// 			},
// 			map[string]func(any) bool{
// 				"SHOOTING": func(a any) bool {
// 					return true
// 				},
// 			},
// 			nil,
// 			func(st *state_machine.State) {

// 			},
// 			nil,
// 		),
// 		state_machine.NewState(
// 			"SHOOTING",
// 			nil,
// 			map[string]func(any) bool{},
// 			nil,
// 			nil,
// 			nil,
// 		),
// 	}
// }

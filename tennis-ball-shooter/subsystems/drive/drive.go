package drive

import (
	"fmt"
	"tennis-ball-shooter/configs"
	"tennis-ball-shooter/constants"
	"tennis-ball-shooter/subsystems/drive/states/auto"
	"tennis-ball-shooter/subsystems/drive/states/teleop"
	drive_types "tennis-ball-shooter/subsystems/drive/types"

	"github.com/LikeMyGames/FRC-Like-Robot/state/controller"
	"github.com/LikeMyGames/FRC-Like-Robot/state/drive/swerve"
	"github.com/LikeMyGames/FRC-Like-Robot/state/event"
	"github.com/LikeMyGames/FRC-Like-Robot/state/state_machine"
	"github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"
)

type (
	DriveSubsystem drive_types.DriveSubsystem
)

var (
	// ctrl               *controller.Controller = nil
	ctrlTrans          mathutils.Vector2D = mathutils.Vector2D{}
	ctrlRot            mathutils.Vector2D = mathutils.Vector2D{}
	transEventTarget   string             = controller.LeftStick
	transEventListener *event.Listener    = nil
	rotEventTarget     string             = controller.RightStick
	rotEventListener   *event.Listener    = nil
)

const (
	PoseListenerTarget string = "DRIVE_SUBSYSTEM_POSE"
)

var instance *DriveSubsystem

func New() *DriveSubsystem {
	s := new(DriveSubsystem)
	s.SwerveDrive = swerve.NewSwerveDriveWithConfigs(constants.Drive.Swerve, configs.DriveMotor(), configs.AzimuthMotor())
	s.StateMachine = state_machine.NewStateMachine(
		teleop.Get(s.purify()),
		auto.Get(s.purify()),
	)

	s.DriverController = controller.GetController(constants.DriverControllerNum)

	instance = s

	return s
}

func GetInstance() *DriveSubsystem {
	if instance != nil {
		return instance
	}
	return nil
}

func (s *DriveSubsystem) purify() *drive_types.DriveSubsystem {
	return (*drive_types.DriveSubsystem)(s)
}

func (s *DriveSubsystem) Initialize() {

}

func (s *DriveSubsystem) Periodic() {

}

func (drive *DriveSubsystem) Drive(trans, rot mathutils.Vector2D, worldRelative bool) {
	xSpeed := trans.X * drive.SwerveDrive.Config.MaxSpeed.TranslationalV
	ySpeed := trans.Y * drive.SwerveDrive.Config.MaxSpeed.TranslationalV
	rotSpeed := rot.X * drive.SwerveDrive.Config.MaxSpeed.RotationalV

	// xSpeed := trans.X
	// ySpeed := trans.Y
	// rotSpeed := rot.X

	if !worldRelative {
		trans.Rotate(-drive.Pose.Angle)
		xSpeed = trans.X
		ySpeed = trans.Y
	}

	states := drive.SwerveDrive.CalculateSwerve(xSpeed, ySpeed, rotSpeed)

	fmt.Println()
	for i, v := range states {
		fmt.Println(i, ": ", v)
	}
	drive.SwerveDrive.Normalize(&states, xSpeed, ySpeed, rotSpeed)

	fmt.Println()
	for i, v := range states {
		fmt.Println(i, ": ", v)
	}

	for i, v := range drive.SwerveDrive.SwerveModules {
		v.SetTarget(states[i])
	}
}

func (drive *DriveSubsystem) DriveToPose(pose mathutils.Pose2D) {
	diff := mathutils.Pose2D{}
	diff.Location = mathutils.Vector2D{X: pose.Location.X - drive.Pose.Location.X, Y: pose.Location.Y - drive.Pose.Location.Y}
	diff.Angle = pose.Angle - drive.Pose.Angle
	fmt.Println(diff)
}

// func SetDriveController(controller *controller.Controller) {
// 	ctrl = controller
// }

func (s *DriveSubsystem) GetPose2D() mathutils.Pose2D {
	return s.Pose
}

func SetTransEventTarget(target string) {
	event.Remove(transEventListener)
	if target != "" {
		// transEventListener = event.Listen(target, "DRIVE_SUBSYSTEM_TRANS", func(event any) {
		// 	trans := event.(mathutils.Vector2D)
		// 	ctrlTrans = trans
		// })
		transEventListener = event.Listen(target, func(event any) {
			trans := event.(mathutils.Vector2D)
			ctrlTrans = trans
		})
	}
}

func GetTransEventTarget() string {
	return transEventTarget
}

func SetRotEventTarget(target string) {
	event.Remove(rotEventListener)
	if target != "" {
		// rotEventListener = event.Listen(target, "DRIVE_SUBSYSTEM_ROT", func(event any) {
		// 	rot := event.(mathutils.Vector2D)
		// 	ctrlRot = rot
		// })
		rotEventListener = event.Listen(target, func(event any) {
			rot := event.(mathutils.Vector2D)
			ctrlRot = rot
		})
	}
}

func GetRotEventTarget() string {
	return rotEventTarget
}

// func GetDriveVectorsFromController(ctrl *controller.Controller) (trans, rot mathutils.Vector2D) {
// 	return mathutils.Vector2D{}, mathutils.Vector2D{}
// }

// func (drive *SwerveDrive) DriveToRelativePose(diff Types.Pose2D) {

// }

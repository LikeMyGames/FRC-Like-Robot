package main

import (
	"frcrobot/internal/Controller"
	"frcrobot/internal/DriveSubsystem"
	"frcrobot/internal/EventListener"
	"frcrobot/internal/Utils/VectorMath"
	// Webpage "robot/internal/Webpage"
)

func main() {
	drive := DriveSubsystem.NewSwerveDrive("robot.constants")

	EventListener.Listen("THUMB_L", func(a ...any) any {
		thumbL := a[0].([]float32)
		v := VectorMath.Vector2D{X: thumbL[0], Y: thumbL[1]}
		drive.CalculateSwerveModules(v, drive.Pose.Angle)

		return nil
	})

	EventListener.Listen("THUMB_R", func(a ...any) any {
		thumbR := a[0].([]float32)
		drive.CalculateSwerveModules(drive.Pose.Location, thumbR[0])

		return nil
	})

	Controller.StartController()
	// Webpage.Start()
	//  Webpage.SendVariables()
	//  Call at end of file
}

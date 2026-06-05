package configs

import (
	"tennis-ball-shooter/constants"
	intake_types "tennis-ball-shooter/subsystems/intake/types"
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
)

func DriveMotor() motor.Config {
	config := new(motor.Config)

	return *config
}

func AzimuthMotor() motor.Config {
	config := new(motor.Config)

	return *config
}

func FlywheelMotor() motor.Config {
	config := new(motor.Config)

	config.Velocity().SetP(constants.Shooter.FlywheelMotorP).SetD(constants.Shooter.FlywheelMotorD).SetFF(constants.Shooter.FlywheelMotorFF)
	config.VelocityConversion = constants.Shooter.FlywheelMotorVelocityConversion
	config.AccelerationConversion = constants.Shooter.FlywheelMotorAccelerationConversion

	return *config
}

func HoodMotor() motor.Config {
	config := new(motor.Config)

	config.Position().SetP(constants.Shooter.HoodMotorP).SetI(constants.Shooter.HoodMotorI).SetD(constants.Shooter.HoodMotorD)
	config.PositionConversion = constants.Shooter.HoodMotorPositionConversion
	config.VelocityConversion = constants.Shooter.HoodMotorVelocityConversion
	config.AccelerationConversion = constants.Shooter.HoodMotorAccelerationConversion

	return *config
}

func TurretMotor() motor.Config {
	config := new(motor.Config)

	config.Position().SetP(constants.Shooter.TurretMotorP).SetI(constants.Shooter.TurretMotorI).SetD(constants.Shooter.TurretMotorD)
	config.PositionConversion = constants.Shooter.TurretMotorPositionConversion
	config.VelocityConversion = constants.Shooter.TurretMotorVelocityConversion
	config.AccelerationConversion = constants.Shooter.TurretMotorAccelerationConversion

	return *config
}

var IntakeMotors intake_types.MotorConfigs = intake_types.MotorConfigs{
	// RollerMotorConfig: motor.Config{
	// 	CanId:                  constants.Intake.RollerMotorCanId,
	// 	PID_P:                  constants.Intake.RollerMotorP,
	// 	PID_D:                  constants.Intake.RollerMotorD,
	// 	PID_FF:                 constants.Intake.RollerMotorFF,
	// 	VelocityConversion:     constants.Intake.RollerMotorVelocityConversion,
	// 	AccelerationConversion: constants.Intake.RollerMotorAccelerationConversion,
	// },
	// ExtensionMotorConfig: motor.Config{
	// 	CanId:                  constants.Intake.ExtensionMotorCanId,
	// 	PID_P:                  constants.Intake.ExtensionMotorP,
	// 	PID_I:                  constants.Intake.ExtensionMotorI,
	// 	PID_D:                  constants.Intake.ExtensionMotorD,
	// 	VelocityConversion:     constants.Intake.ExtensionMotorVelocityConversion,
	// 	AccelerationConversion: constants.Intake.ExtensionMotorAccelerationConversion,
	// },
}

var ShooterMotors shooter_types.MotorConfigs = shooter_types.MotorConfigs{
	// SpinnerMotor: motor.Config{
	// 	CanId:                  constants.Shooter.SpinnerMotorCanId,
	// 	PID_P:                  constants.Shooter.SpinnerMotorP,
	// 	PID_D:                  constants.Shooter.SpinnerMotorD,
	// 	PID_FF:                 constants.Shooter.SpinnerMotorFF,
	// 	VelocityConversion:     constants.Shooter.SpinnerMotorVelocityConversion,
	// 	AccelerationConversion: constants.Shooter.SpinnerMotorAccelerationConversion,
	// },
	// FlywheelMotor: ,
}

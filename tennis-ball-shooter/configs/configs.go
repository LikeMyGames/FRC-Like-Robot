package configs

import (
	"tennis-ball-shooter/constants"
	intake_types "tennis-ball-shooter/subsystems/intake/types"
	shooter_types "tennis-ball-shooter/subsystems/shooter/types"

	motor "github.com/LikeMyGames/FRC-Like-Robot/state/motor_controller"
)

func DriveMotor() motor.Config {
	config := motor.NewConfig()

	return config
}

func AzimuthMotor() motor.Config {
	config := motor.NewConfig()

	return config
}

func FlywheelMotor() motor.Config {
	config := motor.NewConfig()

	config.Velocity().SetP(constants.Shooter.FlywheelMotorP).SetD(constants.Shooter.FlywheelMotorD).SetFF(constants.Shooter.FlywheelMotorFF)
	config.VelocityConversion = constants.Shooter.FlywheelMotorVelocityConversion
	config.AccelerationConversion = constants.Shooter.FlywheelMotorAccelerationConversion

	return config
}

func HoodMotor() motor.Config {
	config := motor.NewConfig()

	config.Position().SetP(constants.Shooter.HoodMotorP).SetI(constants.Shooter.HoodMotorI).SetD(constants.Shooter.HoodMotorD)
	config.PositionConversion = constants.Shooter.HoodMotorPositionConversion
	config.VelocityConversion = constants.Shooter.HoodMotorVelocityConversion
	config.AccelerationConversion = constants.Shooter.HoodMotorAccelerationConversion

	return config
}

func TurretMotor() motor.Config {
	config := motor.NewConfig()

	config.Position().SetP(constants.Shooter.TurretMotorP).SetI(constants.Shooter.TurretMotorI).SetD(constants.Shooter.TurretMotorD)
	config.PositionConversion = constants.Shooter.TurretMotorPositionConversion
	config.VelocityConversion = constants.Shooter.TurretMotorVelocityConversion
	config.AccelerationConversion = constants.Shooter.TurretMotorAccelerationConversion

	return config
}

var IntakeMotors intake_types.MotorConfigs = intake_types.MotorConfigs{
	RollerMotorConfig:    Intake_RollerMotorConfig,
	ExtensionMotorConfig: Intake_ExtensionMotorConfig,
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

func Intake_RollerMotorConfig() motor.Config {
	config := motor.NewConfig()

	config.Velocity().SetP(constants.Intake.RollerMotorP).SetD(constants.Intake.RollerMotorD).SetFF(constants.Intake.RollerMotorFF)
	config.VelocityConversion = constants.Intake.RollerMotorVelocityConversion
	config.AccelerationConversion = constants.Intake.RollerMotorAccelerationConversion

	return config
}

func Intake_ExtensionMotorConfig() motor.Config {
	config := motor.NewConfig()

	config.Velocity().SetP(constants.Intake.ExtensionMotorP).SetI(constants.Intake.ExtensionMotorI).SetD(constants.Intake.ExtensionMotorD).SetFF(constants.Intake.ExtensionMotorCosFF)
	config.PositionConversion = constants.Intake.ExtensionMotorPositionConversion
	config.VelocityConversion = constants.Intake.ExtensionMotorVelocityConversion
	config.AccelerationConversion = constants.Intake.ExtensionMotorAccelerationConversion

	return config
}

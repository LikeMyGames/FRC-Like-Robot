package pid

import (
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/constantTypes"
	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"
)

// PIDController represents a PID controller.
type PIDController struct {
	Kp, Ki, Kd float64       // PID gains
	prevError  float64       // Previous error for derivative term
	integral   float64       // Integral sum for integral term
	setpoint   float64       // Desired target value
	minOutput  float64       // Minimum allowed output
	maxOutput  float64       // Maximum allowed output
	dt         time.Duration // Time step for calculations
}

// NewPIDController creates and returns a new PIDController instance.
// func NewPIDController(Kp, Ki, Kd, setpoint, minOutput, maxOutput float64) *PIDController {
func NewPIDController(config constantTypes.PidController) *PIDController {
	return &PIDController{
		Kp:        config.Kp,
		Ki:        config.Ki,
		Kd:        config.Kd,
		minOutput: config.MinOut,
		maxOutput: config.MaxOut,
		dt:        robot.RobotRef.Frequency,
	}
}

func (pid *PIDController) SetTarget(target float64) {
	pid.setpoint = target
}

// Calculate computes the control output based on the current process value.
func (pid *PIDController) Calculate(processValue float64) float64 {
	err := pid.setpoint - processValue

	// Proportional term
	proportionalTerm := pid.Kp * err

	// Integral term
	pid.integral += err * pid.dt.Seconds()
	// Anti-windup (optional, but recommended)
	if pid.integral > pid.maxOutput {
		pid.integral = pid.maxOutput
	} else if pid.integral < pid.minOutput {
		pid.integral = pid.minOutput
	}
	integralTerm := pid.Ki * pid.integral

	// Derivative term
	derivativeTerm := pid.Kd * (err - pid.prevError) / pid.dt.Seconds()
	pid.prevError = err

	// Total output
	output := proportionalTerm + integralTerm + derivativeTerm

	// Clamp output to defined limits
	if output > pid.maxOutput {
		output = pid.maxOutput
	} else if output < pid.minOutput {
		output = pid.minOutput
	}

	return output
}

// SetSetpoint updates the desired target value.
// func (pid *PIDController) SetSetpoint(newSetpoint float64) {
// 	pid.setpoint = newSetpoint
// }

// func main() {
// 	// Example usage
// 	Kp, Ki, Kd := 0.5, 0.1, 0.2         // PID gains
// 	setpoint := 100.0                   // Desired value
// 	minOutput, maxOutput := -50.0, 50.0 // Output limits
// 	dt := 100 * time.Millisecond        // Time step

// 	pid := NewPIDController(Kp, Ki, Kd, setpoint, minOutput, maxOutput, dt)

// 	processValue := 0.0 // Initial process value

// 	fmt.Println("Simulating PID control:")
// 	for i := 0; i < 50; i++ {
// 		controlOutput := pid.Calculate(processValue)
// 		// Simulate system response (replace with your actual system dynamics)
// 		// For this example, we'll just add a fraction of the output to the process value
// 		processValue += controlOutput * 0.1

// 		fmt.Printf("Iteration %d: Process Value = %.2f, Control Output = %.2f\n", i, processValue, controlOutput)
// 		time.Sleep(dt)
// 	}
// }

package pid

import (
	"time"

	"github.com/LikeMyGames/FRC-Like-Robot/state/robot"
)

// PIDController represents a PID controller.
type PIDController struct {
	constants Constants
	prevError float64       // Previous error for derivative term
	integral  float64       // Integral sum for integral term
	setpoint  float64       // Desired target value
	minOutput float64       // Minimum allowed output
	maxOutput float64       // Maximum allowed output
	dt        time.Duration // Time step for calculations
}

type Constants struct {
	kP, kI, kD, kIzone, kFF float64
}

type Config struct {
	constants    Constants
	min, max, dt float64
}

// NewPIDController creates and returns a new PIDController instance.
// func NewPIDController(Kp, Ki, Kd, setpoint, minOutput, maxOutput float64) *PIDController {
func NewFromConfig(config Config) *PIDController {
	controller := new(PIDController)
	controller.constants = config.constants
	controller.maxOutput = config.max
	controller.minOutput = config.min
	controller.dt = robot.RobotRef.Frequency

	return controller
}

func New() *PIDController {
	return &PIDController{
		constants: Constants{},
		minOutput: 0,
		maxOutput: 0,
		dt:        0,
	}
}

func NewFromConstants(constants Constants) *PIDController {
	controller := new(PIDController)
	controller.constants = constants

	return controller
}

func (p *PIDController) SetP(kP float64) {
	// p.Kp = kP
	p.constants.SetP(kP)
}

func (p *PIDController) SetI(kI float64) {
	// p.Kp = kI
	p.constants.SetI(kI)
}

func (p *PIDController) SetD(kD float64) {
	// p.Kd = kD
	p.constants.SetD(kD)
}

func (p *PIDController) GetP() float64 {
	// return p.Kp
	return p.constants.kP
}

func (p *PIDController) GetI() float64 {
	// return p.Ki
	return p.constants.kI
}

func (p *PIDController) GetD() float64 {
	// return p.Kd
	return p.constants.kD
}

func (p *PIDController) SetMaxOuput(max float64) {
	p.maxOutput = max
}

func (p *PIDController) SetMinOuput(min float64) {
	p.minOutput = min
}

func (p *PIDController) GetMaxOuput() float64 {
	return p.maxOutput
}

func (p *PIDController) GetMinOuput() float64 {
	return p.minOutput
}

func (pid *PIDController) SetTarget(target float64) {
	pid.setpoint = target
}

func (pid *PIDController) GetTarget() float64 {
	return pid.setpoint
}

// Calculate computes the control output based on the current process value.
func (pid *PIDController) Calculate(processValue float64) float64 {
	err := pid.setpoint - processValue

	// Proportional term
	proportionalTerm := pid.constants.kP * err

	// Integral term
	pid.integral += err * pid.dt.Seconds()
	// Anti-windup (optional, but recommended)
	if pid.integral > pid.maxOutput {
		pid.integral = pid.maxOutput
	} else if pid.integral < pid.minOutput {
		pid.integral = pid.minOutput
	}
	integralTerm := pid.constants.kI * pid.integral

	// Derivative term
	derivativeTerm := pid.constants.kD * (err - pid.prevError) / pid.dt.Seconds()
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

func NewConstants(p, i, d float64) Constants {
	return Constants{
		kP: p,
		kI: i,
		kD: d,
	}
}

func (c *Constants) GetAsArray() []float64 {
	return []float64{c.kP, c.kI, c.kD, c.kIzone, c.kFF}
}

func (c *Constants) SetP(p float64) *Constants {
	c.kP = p
	return c
}

func (c *Constants) SetI(i float64) *Constants {
	c.kI = i
	return c
}

func (c *Constants) SetD(d float64) *Constants {
	c.kD = d
	return c
}

func (c *Constants) SetIZone(izone float64) *Constants {
	c.kIzone = izone
	return c
}

func (c *Constants) SetFF(ff float64) *Constants {
	c.kFF = ff
	return c
}

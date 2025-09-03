package MathUtils

import "math"

// from x = (v_0) * t + 0.5 * a * t^2 kinematic equation
func DistanceFromAccelerationTimeVelocity(a, t, v0 float64) float64 {
	return (v0 * t) + (0.5 * a * (t * t))
}

func AccelerationFromDistanceTimeVelocity(x, t, v0 float64) float64 {
	return (x - (v0 * t)) / (0.5 * (t * t))
}

func VelocityFromDistanceAccelerationTime(x, a, t float64) float64 {
	return (x - (0.5 * a * (t * t))) / t
}

func TimeFromDistanceAccelerationVelocity(x, a, v0 float64) float64 {
	discriminant := ((v0 * v0) - (2 * a * (-x)))
	if discriminant < 0 {
		return math.NaN()
	} else if discriminant == 0 {
		return (-v0) / (a)
	}
	return (-v0 + (math.Sqrt(discriminant))) / (a)
}

// from (v_f)^2 = (v_0)^2 + 2ax kinematic equation
func InitialVelocityFromAccelerationFinalVelocityDistance(a, vf, x float64) float64 {
	return math.Sqrt(math.Abs((vf * vf) - (2 * a * x)))
}

func FinalVelocityFromAccelerationInitialVelocityDistance(a, v0, x float64) float64 {
	return math.Sqrt(math.Abs((v0 * v0) + (2 * a * x)))
}

func AccelerationFromInitialVelocityFinalVelocityDistance(v0, vf, x float64) float64 {
	return ((vf * vf) - (v0 * v0)) / (2 * x)
}

func DistanceFromInitialVelocityFinalVelocityAcceleration(v0, vf, a float64) float64 {
	return ((vf * vf) - (v0 * v0)) / (2 * a)
}

// from v_f = v_t + a * t
func FinalVelocityFromInitalVelocityAccelerationTime(v0, a, t float64) float64 {
	return (v0 + (a * t))
}

func InitalVelocityFromFinalVelocityAccelerationTime(vf, a, t float64) float64 {
	return (vf - (a * t))
}

func AccelerationFromInitialVelocityFinalVelocityTime(v0, vf, t float64) float64 {
	return ((vf - v0) / t)
}

func TimeFromInitialVelocityFinalVelocityAccerlation(v0, vf, a float64) float64 {
	return ((vf - v0) / a)
}

package mathutils

import (
	"math"
)

type (
	Vector2D struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	VectorTheta struct {
		Magnitude float64 `json:"a"`
		Angle     float64 `json:"theta"`
	}

	Pose3D struct {
		XY Pose2D
		Z  float64
	}

	Pose2D struct {
		Location Vector2D
		Angle    float64
	}

	Location3D struct {
		XY Location2D
		Z  float64
	}

	Location2D struct {
		X float64
		Y float64
	}

	Axis struct {
		X float64
		Y float64
	}
)

func Vector2DtoVectorTheta(v Vector2D) VectorTheta {
	theta := 0.0
	if v.X == 0 {
		if v.Y > 0 {
			theta = (3 * (math.Pi / 2))
		} else if v.Y < 0 {
			theta = (3 * (math.Pi / 2))
		}
	} else {
		theta = math.Atan(math.Abs(v.Y / v.X))
		if v.X < 0 && v.Y > 0 {
			theta = math.Atan(math.Abs(v.Y/v.X)) + (math.Pi / 2)
		} else if v.X < 0 && v.Y < 0 {
			theta = math.Atan(math.Abs(v.Y/v.X)) + math.Pi
		} else if v.X > 0 && v.Y < 0 {
			theta = math.Atan(math.Abs(v.Y/v.X)) + (3 * (math.Pi / 2))
		}
	}
	l := Clamp(math.Sqrt((v.Y*v.Y)+(v.X*v.X)), 1, 0)
	if math.IsNaN(l) {
		l = 0
	}

	return VectorTheta{Magnitude: l, Angle: theta}
}

func VectorThetatoVector2D(v VectorTheta) Vector2D {
	x := math.Cos(v.Angle) * v.Magnitude
	y := math.Sin(v.Angle) * v.Magnitude

	return Vector2D{X: x, Y: y}
}

func VectorAdd(v1, v2 Vector2D) Vector2D {
	x := v1.X + v2.X
	y := v1.Y + v2.Y
	return Vector2D{X: x, Y: y}
}

func VectorSubtract(v1, v2 Vector2D) Vector2D {
	x := v1.X - v2.X
	y := v1.Y - v2.Y
	return Vector2D{X: x, Y: y}
}
func VectorThetaAdd(v1, v2 VectorTheta) VectorTheta {
	return Vector2DtoVectorTheta(VectorAdd(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
}

func VectorThetaSubtract(v1, v2 VectorTheta) VectorTheta {
	return Vector2DtoVectorTheta(VectorSubtract(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
}

func VectorAddNormalized(v1, v2 Vector2D, maxLen float64) Vector2D {
	v1theta := Vector2DtoVectorTheta(v1)
	v2theta := Vector2DtoVectorTheta(v2)
	v3 := VectorAdd(v1, v2)
	v3theta := Vector2DtoVectorTheta(v3)
	v3theta.Magnitude = MapRange(v3theta.Magnitude, 0, v1theta.Magnitude+v2theta.Magnitude, 0, maxLen)
	return VectorThetatoVector2D(v3theta)
}

func VectorSubtractNormalized(v1, v2 Vector2D, maxLen float64) Vector2D {
	v1theta := Vector2DtoVectorTheta(v1)
	v2theta := Vector2DtoVectorTheta(v2)
	v3 := VectorSubtract(v1, v2)
	v3theta := Vector2DtoVectorTheta(v3)
	v3theta.Magnitude = MapRange(v3theta.Magnitude, 0, v1theta.Magnitude-v2theta.Magnitude, 0, maxLen)
	return VectorThetatoVector2D(v3theta)
}

func VectorThetaAddNormalized(v1, v2 VectorTheta, maxLen float64) VectorTheta {
	v3 := Vector2DtoVectorTheta(VectorAdd(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
	v3.Magnitude = MapRange(v3.Magnitude, 0, v1.Magnitude+v2.Magnitude, 0, maxLen)
	return v3
}

func VectorThetaSubtractNormalized(v1, v2 VectorTheta, maxLen float64) VectorTheta {
	v3 := Vector2DtoVectorTheta(VectorSubtract(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
	v3.Magnitude = MapRange(v3.Magnitude, 0, v1.Magnitude-v2.Magnitude, 0, maxLen)
	return v3
}

func ClampVector(val, max, min Vector2D) Vector2D {
	if val.X > max.X {
		val.X = max.X
	} else if val.X < min.X {
		val.X = min.X
	}
	if val.Y > max.Y {
		val.Y = max.Y
	} else if val.Y < min.Y {
		val.Y = min.Y
	}
	return val
}

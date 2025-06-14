package VectorMath

import (
	"frcrobot/internal/Utils/MathUtils"
	"frcrobot/internal/Utils/Types"
	"math"
)

type (
	Vector2D    = Types.Vector2D
	VectorTheta = Types.VectorTheta
)

func Vector2DtoVectorTheta(v Vector2D) VectorTheta {
	theta := float64(0.0)
	if v.X == 0 {
		if v.Y > 0 {
			theta = float64(math.Pi / 2)
		} else if v.Y < 0 {
			theta = float64(-math.Pi / 2)
		}
	} else {
		theta = float64(math.Atan(math.Abs(float64(v.Y / v.X))))
		if v.X < 0 && v.Y > 0 {
			theta = float64(math.Atan(math.Abs(float64(v.X/v.Y)))) + (math.Pi / 2)
		} else if v.X < 0 && v.Y < 0 {
			theta = float64(math.Atan(math.Abs(float64(v.Y/v.X)))) + math.Pi
		} else if v.X > 0 && v.Y < 0 {
			theta = float64(math.Atan(math.Abs(float64(v.X/v.Y)))) + (3 * (math.Pi / 2))
		}
	}
	l := float64(MathUtils.Clamp(float64(v.Y/math.Sin(float64(theta))), 1, 0))
	if math.IsNaN(l) {
		l = 0
	}

	return VectorTheta{L: l, T: theta}
}

func VectorThetatoVector2D(v VectorTheta) Vector2D {
	x := math.Cos(float64(v.T)) * v.L
	y := math.Sin(float64(v.T)) * v.L

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
	v3theta.L = MathUtils.MapRange(v3theta.L, 0, v1theta.L+v2theta.L, 0, maxLen)
	return VectorThetatoVector2D(v3theta)
}

func VectorSubtractNormalized(v1, v2 Vector2D, maxLen float64) Vector2D {
	v1theta := Vector2DtoVectorTheta(v1)
	v2theta := Vector2DtoVectorTheta(v2)
	v3 := VectorSubtract(v1, v2)
	v3theta := Vector2DtoVectorTheta(v3)
	v3theta.L = MathUtils.MapRange(v3theta.L, 0, v1theta.L-v2theta.L, 0, maxLen)
	return VectorThetatoVector2D(v3theta)
}

func VectorThetaAddNormalized(v1, v2 VectorTheta, maxLen float64) VectorTheta {
	v3 := Vector2DtoVectorTheta(VectorAdd(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
	v3.L = MathUtils.MapRange(v3.L, 0, v1.L+v2.L, 0, maxLen)
	return v3
}

func VectorThetaSubtractNormalized(v1, v2 VectorTheta, maxLen float64) VectorTheta {
	v3 := Vector2DtoVectorTheta(VectorSubtract(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
	v3.L = MathUtils.MapRange(v3.L, 0, v1.L-v2.L, 0, maxLen)
	return v3
}

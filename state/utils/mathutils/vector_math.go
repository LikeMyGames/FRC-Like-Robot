package mathutils

import (
	"fmt"
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
		Location Location3D
		Angle    float64
	}

	Pose2D struct {
		Location Vector2D
		Angle    float64
	}

	Location3D struct {
		X float64
		Y float64
		Z float64
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
	theta := math.Atan2(v.Y, v.X)
	l := math.Hypot(v.Y, v.X)
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

// VectorAdd adds the two Vector2D's provided in a head to tail fashion
// v2 is essentially translated so that its tail is on the same exact point as the head of v1.
// The vector output from this function is then acquaried by creating a new vector that goes
// from the tail of v1 to the head of v2.
// The way the function actually works is by adding the X's individually and the Y's individually.
// This has the same effect as the explaination first provided.
func AddVector2D(v1, v2 Vector2D) *Vector2D {
	x := v1.X + v2.X
	y := v1.Y + v2.Y
	return &Vector2D{X: x, Y: y}
}

func SubtractVector2D(v1, v2 Vector2D) *Vector2D {
	x := v1.X - v2.X
	y := v1.Y - v2.Y
	return &Vector2D{X: x, Y: y}
}

func MultiplyVector2D(v Vector2D, num float64) *Vector2D {
	return v.Multiply(num)
}

func DivideVector2D(v Vector2D, num float64) *Vector2D {
	return v.Divide(num)
}

func VectorThetaAdd(v1, v2 VectorTheta) VectorTheta {
	return Vector2DtoVectorTheta(*AddVector2D(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
}

func VectorThetaSubtract(v1, v2 VectorTheta) VectorTheta {
	return Vector2DtoVectorTheta(*SubtractVector2D(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
}

func VectorAddNormalized(v1, v2 Vector2D, maxLen float64) Vector2D {
	v1theta := Vector2DtoVectorTheta(v1)
	v2theta := Vector2DtoVectorTheta(v2)
	v3 := AddVector2D(v1, v2)
	v3theta := Vector2DtoVectorTheta(*v3)
	v3theta.Magnitude = MapRange(v3theta.Magnitude, 0, v1theta.Magnitude+v2theta.Magnitude, 0, maxLen)
	return VectorThetatoVector2D(v3theta)
}

func VectorSubtractNormalized(v1, v2 Vector2D, maxLen float64) Vector2D {
	v1theta := Vector2DtoVectorTheta(v1)
	v2theta := Vector2DtoVectorTheta(v2)
	v3 := SubtractVector2D(v1, v2)
	v3theta := Vector2DtoVectorTheta(*v3)
	v3theta.Magnitude = MapRange(v3theta.Magnitude, 0, v1theta.Magnitude-v2theta.Magnitude, 0, maxLen)
	return VectorThetatoVector2D(v3theta)
}

func VectorThetaAddNormalized(v1, v2 VectorTheta, maxLen float64) VectorTheta {
	v3 := Vector2DtoVectorTheta(*AddVector2D(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
	v3.Magnitude = MapRange(v3.Magnitude, 0, v1.Magnitude+v2.Magnitude, 0, maxLen)
	return v3
}

func VectorThetaSubtractNormalized(v1, v2 VectorTheta, maxLen float64) VectorTheta {
	v3 := Vector2DtoVectorTheta(*SubtractVector2D(VectorThetatoVector2D(v1), VectorThetatoVector2D(v2)))
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

func DotProductVector2D(v1, v2 Vector2D) float64 {
	return (v1.X * v2.X) + (v1.Y * v2.Y)
}

func CrossProductVector2D(v1, v2 Vector2D) float64 {
	return (v1.X * v2.Y) - (v1.Y * v2.X)
}

// angle should be in radians
func (v *Vector2D) Rotate(angle float64) *Vector2D {
	sin, cos := math.Sincos(angle)
	v.X = (v.X * cos) - (v.Y * sin)
	v.Y = (v.X * sin) + (v.Y * cos)
	return v
}

func (v Vector2D) ToVectorTheta() VectorTheta {
	return Vector2DtoVectorTheta(v)
}

func (v Vector2D) DotProduct(v2 Vector2D) float64 {
	return DotProductVector2D(v, v2)
}

func (v Vector2D) CrossProduct(v2 Vector2D) float64 {
	return CrossProductVector2D(v, v2)
}

func (o *Vector2D) Add(v *Vector2D) *Vector2D {
	o.X += v.X
	o.Y += v.X
	return o
}

func (o *Vector2D) Subtract(v *Vector2D) *Vector2D {
	o.X -= v.X
	o.Y -= v.X
	return o
}

func (v *Vector2D) Multiply(num float64) *Vector2D {
	v.X *= num
	v.Y *= num
	return v
}

func (v *Vector2D) Divide(num float64) *Vector2D {
	v.X /= num
	v.Y /= num
	return v
}

func (v *Vector2D) Length() float64 {
	return math.Sqrt(v.DotProduct(*v))
}

func (v Vector2D) String() string {
	return fmt.Sprintf("{X: %f Y: %f}", v.X, v.Y)
}

func (v VectorTheta) String() string {
	return fmt.Sprintf("{Magnitude: %f Angle: %f}", v.Magnitude, v.Angle)
}

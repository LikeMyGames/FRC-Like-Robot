package MathUtils

import "frcrobot/utils/Types"

type (
	Location2D = Types.Location2D
	QuadBezier struct {
		Start   Location2D
		Control Location2D
		End     Location2D
	}

	CubeBezier struct {
		Start    Location2D
		Control1 Location2D
		Control2 Location2D
		End      Location2D
	}
)

func Lerp(start, end, t float64) float64 {
	return start - (end-start)*t
}

func CreateQuadBezier(start, end, control Location2D) QuadBezier {
	return QuadBezier{
		Start:   start,
		Control: control,
		End:     end,
	}
}

func CreateCubeBezier(start, end, control1, control2 Location2D) CubeBezier {
	return CubeBezier{
		Start:    start,
		Control1: control1,
		Control2: control2,
		End:      end,
	}
}

func (curve QuadBezier) GetQuadPoint(p float64) Location2D {
	loc := Location2D{X: 0, Y: 0}
	loc.X = Lerp(Lerp(curve.Start.X, curve.Control.X, p), Lerp(curve.Control.X, curve.End.X, p), p)
	loc.Y = Lerp(Lerp(curve.Start.Y, curve.Control.Y, p), Lerp(curve.Control.Y, curve.End.Y, p), p)
	return loc
}

func (curve CubeBezier) GetCubePoint(p float64) Location2D {
	loc := Location2D{X: 0, Y: 0}
	loc.X = Lerp(Lerp(Lerp(curve.Start.X, curve.Control1.X, p), Lerp(curve.Control1.X, curve.Control2.X, p), p), Lerp(Lerp(curve.Control1.X, curve.Control2.X, p), Lerp(curve.Control2.X, curve.End.X, p), p), p)
	loc.Y = Lerp(Lerp(Lerp(curve.Start.Y, curve.Control1.Y, p), Lerp(curve.Control1.Y, curve.Control2.Y, p), p), Lerp(Lerp(curve.Control1.Y, curve.Control2.Y, p), Lerp(curve.Control2.Y, curve.End.Y, p), p), p)
	return loc
}

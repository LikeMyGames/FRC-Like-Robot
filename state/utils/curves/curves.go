package curves

import "github.com/LikeMyGames/FRC-Like-Robot/state/utils/mathutils"

type (
	CubicBezier struct {
		start           *mathutils.Vector2D
		control_point_1 *mathutils.Vector2D
		control_point_2 *mathutils.Vector2D
		end             *mathutils.Vector2D
		length          float64
	}

	QuadraticBezier struct {
		start         *mathutils.Vector2D
		control_point *mathutils.Vector2D
		end           *mathutils.Vector2D
		length        float64
	}
)

func NewCubicBezier(start, cp1, cp2, end *mathutils.Vector2D) *CubicBezier {
	return &CubicBezier{
		start:           start,
		control_point_1: cp1,
		control_point_2: cp2,
		end:             end,
	}
}

func NewQuadraticBezier(start, cp, end *mathutils.Vector2D) *QuadraticBezier {
	return &QuadraticBezier{
		start:         start,
		control_point: cp,
		end:           end,
	}
}

func (c *CubicBezier) GetPointByPercent(percent float64) *mathutils.Vector2D {
	X := mathutils.Lerp(
		mathutils.Lerp(
			mathutils.Lerp(
				c.start.X,
				c.control_point_1.X,
				percent,
			),
			mathutils.Lerp(
				c.control_point_1.X,
				c.control_point_2.X,
				percent,
			),
			percent,
		),
		mathutils.Lerp(
			mathutils.Lerp(
				c.control_point_1.X,
				c.control_point_2.X,
				percent,
			),
			mathutils.Lerp(
				c.control_point_2.X,
				c.end.X,
				percent,
			),
			percent,
		),
		percent,
	)
	Y := mathutils.Lerp(
		mathutils.Lerp(
			mathutils.Lerp(
				c.start.Y,
				c.control_point_1.Y,
				percent,
			),
			mathutils.Lerp(
				c.control_point_1.Y,
				c.control_point_2.Y,
				percent,
			),
			percent,
		),
		mathutils.Lerp(
			mathutils.Lerp(
				c.control_point_1.Y,
				c.control_point_2.Y,
				percent,
			),
			mathutils.Lerp(
				c.control_point_2.Y,
				c.end.Y,
				percent,
			),
			percent,
		),
		percent,
	)
	return &mathutils.Vector2D{X: X, Y: Y}
}

func (c *CubicBezier) GetDistanceBetweenPoints(startPercent, endPercent float64) float64 {
	return c.GetPointByPercent(endPercent).Subtract(c.GetPointByPercent(startPercent)).Length()
}

func (c *CubicBezier) Length() float64 {
	chord_length := mathutils.SubtractVector2D(*c.end, *c.start).Length()
	control_length := mathutils.SubtractVector2D(*c.end, *c.control_point_2).Length() + mathutils.SubtractVector2D(*c.control_point_2, *c.control_point_1).Length() + mathutils.SubtractVector2D(*c.control_point_1, *c.start).Length()
	c.length = ((2 * chord_length) + (2 * control_length)) / 2
	return c.length
}

func (c *QuadraticBezier) GetPointByPercent(percent float64) *mathutils.Vector2D {
	X := mathutils.Lerp(
		mathutils.Lerp(
			c.start.X,
			c.control_point.X,
			percent,
		),
		mathutils.Lerp(
			c.control_point.X,
			c.end.X,
			percent,
		),
		percent,
	)
	Y := mathutils.Lerp(
		mathutils.Lerp(
			c.start.Y,
			c.control_point.Y,
			percent,
		),
		mathutils.Lerp(
			c.control_point.Y,
			c.end.Y,
			percent,
		),
		percent,
	)
	return &mathutils.Vector2D{X: X, Y: Y}
}

func (c *QuadraticBezier) Length() float64 {
	chord_length := mathutils.SubtractVector2D(*c.end, *c.start).Length()
	control_length := mathutils.SubtractVector2D(*c.end, *c.control_point).Length() + mathutils.SubtractVector2D(*c.control_point, *c.start).Length()
	c.length = (2*chord_length + control_length) / 2
	return c.length
}

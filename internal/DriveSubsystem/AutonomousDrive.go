package DriveSubsystem

import (
	"frcrobot/internal/Utils/MathUtils"
	"frcrobot/internal/Utils/Types"
)

type (
	Auto struct {
		Name  string
		Paths []Path
	}

	Path struct {
		Name       string
		Resolution int
		Curve      MathUtils.CubeBezier
		Points     []Point
	}

	Point struct {
		Location MathUtils.Location2D
		// InitialVelocity float64
		// FinalVelocity   float64
	}
)

func NewAuto(name string) *Auto {
	return &Auto{Name: name}
}

// func (auto *Auto) RunAuto() {

// }

func NewPath(name string, start, control1, control2, end Types.Location2D, resolution int) Path {
	path := Path{Name: name, Resolution: resolution}
	path.Curve = MathUtils.CreateCubeBezier(start, control1, control2, end)
	for i := 0; i < resolution; i++ {
		path.Points[i] = Point{Location: path.Curve.GetCubePoint(float64(i) / float64(resolution))}
	}
	return path
}

// func NewPathFromFile(fileName string) Path {

// }

func (auto *Auto) AddPath(paths ...Path) {
	auto.Paths = append(auto.Paths, paths...)
}

package Types

type (
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

	Vector2D struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	VectorTheta struct {
		L float64 `json:"a"`
		T float64 `json:"theta"`
	}

	Axis struct {
		X float64
		Y float64
	}
)

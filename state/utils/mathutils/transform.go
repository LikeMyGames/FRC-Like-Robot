package mathutils

type (
	Translation2D struct {
		x, y float64
	}
	Translation3D struct {
		x, y, z float64
	}

	Rotation3D struct {
		x, y, z        float64
		qx, qy, qz, qw float64
	}

	Transform3D struct {
		trans Translation3D
		rot   Rotation3D
	}
)

func (t *Transform3D) GetTranslation3D() Translation3D {
	return t.trans
}

func (t *Transform3D) GetRotation3D() Rotation3D {
	return t.rot
}

func (t *Translation3D) GetX() float64 {
	return t.x
}

func (t *Translation3D) GetY() float64 {
	return t.y
}

func (t *Translation3D) GetZ() float64 {
	return t.z
}

func (t *Translation2D) ToTransform3D() *Translation3D {
	return &Translation3D{x: t.x, y: t.y, z: 0}
}

func (r *Rotation3D) GetX() float64 {
	return r.x
}

func (r *Rotation3D) GetY() float64 {
	return r.y
}

func (r *Rotation3D) GetZ() float64 {
	return r.z
}

func (r *Rotation3D) GetQuaternionX() float64 {
	return r.qx
}

func (r *Rotation3D) GetQuaternionY() float64 {
	return r.qy
}

func (r *Rotation3D) GetQuaternionZ() float64 {
	return r.qz
}

func (r *Rotation3D) GetQuaternionW() float64 {
	return r.qw
}

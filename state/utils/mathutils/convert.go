package mathutils

import (
	"encoding/binary"
	"math"
)

func BytesToFloat32(data []byte) float32 {
	int_data := binary.LittleEndian.Uint32(data[:4])
	return math.Float32frombits(int_data)
}

func WriteBytesToFloat32(data []byte, float *float32) {
	int_data := binary.LittleEndian.Uint32(data[:4])
	float_data := math.Float32frombits(int_data)
	float = &float_data
}

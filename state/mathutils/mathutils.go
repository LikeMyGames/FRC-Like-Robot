package mathutils

func MapRange(num, inLow, inHigh, outLow, outHigh float64) float64 {
	newNum := (outLow + (((num - inLow) / (inHigh - inLow)) * (outHigh - outLow)))
	return newNum
}

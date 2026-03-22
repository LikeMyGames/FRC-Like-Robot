package utils

import "math"

/*
Gets value from channel ch and return a pointer to that value.
Will return nil if channel ch is empty.
This creates a non-blocking affect (hence the naming of the function)
compared to the normal blocking behavior of reading a channel
*/
func ReadChannelNonBlocking[T any](ch chan T) *T {
	var val T
	select {
	case val = <-ch:
	default:
		// fmt.Println("Channel is empty")
		return nil
	}
	return &val
}

func TruncateFloat64(num float64, precision int) float64 {
	num *= math.Pow(10, float64(precision))
	num = math.Trunc(num)
	num /= math.Pow(10, float64(precision))
	return num
}

func ConvertEdgeAlignedAngleToCenterAligned(angle float64) float64 {
	if angle >= 180 {
		angle = 180 - angle
	} else if angle < -180 {
		angle = 360 + angle
	}
	return angle
}

func ConvertCenterAlignedAngleToEdgeAligned(angle float64) float64 {
	if angle < 0 {
		angle = 180 - angle
	}
	return angle
}

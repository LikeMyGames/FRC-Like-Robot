package mathutils

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func RadtoDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func DegtoRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func Clamp(v, max, min float64) float64 {
	if v > max {
		return max
	} else if v < min {
		return min
	}
	return v
}

func MapRange(num, inLow, inHigh, outLow, outHigh float64) float64 {
	newNum := (outLow + (((num - inLow) / (inHigh - inLow)) * (outHigh - outLow)))
	return newNum
}

func Trunc(v float64, len int) float64 {
	str := fmt.Sprintf("%f", v)
	val := float64(0.0)
	var err error
	if len == 0 {
		val, err = strconv.ParseFloat(str[:strings.Index(str, ".")], 64)
	} else {
		val, err = strconv.ParseFloat(str[:strings.Index(str, ".")+1+len], 64)
	}
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func Avg(n ...float64) float64 {
	sum := 0.0
	for _, v := range n {
		sum += v
	}
	return sum / float64(len(n))
}

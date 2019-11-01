package utils

import (
	"math"
	"math/rand"
)

// Round f to nearest number of decimal points
func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}

// Sample norm distrib with mean and sd given
func Norm(mean float64, sd float64) float64 {
	return rand.NormFloat64()*sd + mean
}

// Sample chi-squarre distrib
func ChiSq(mean float64, sd float64) float64 {
	tmp := rand.NormFloat64()*sd + mean
	tmp = tmp * tmp
	return tmp
}

// Round a float to the nearest int
func Round(f float64) float64 {
	return math.Floor(f + .5)
}

// generate a random number between two floats
func RandRangeFloat(min, max float64) float64 {
	return MapNumFloat(rand.Float64(), 0, 1, min, max)
}

// map a number from one range to another range
func MapNumFloat(value, inMin, inMax, outMin, outMax float64) float64 {
	if value >= inMax {
		return outMax
	}
	if value <= inMin {
		return outMin
	}
	return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func RandRangInt(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
func MapNumInt(value, inMin, inMax, outMin, outMax int64) int64 {
	if value >= inMax {
		return outMax
	}
	if value <= inMin {
		return outMin
	}
	return (value-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

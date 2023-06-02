package math64

import "math"

const (

	// PI and its commonly needed variants.
	PI     float64 = math.Pi
	PIx2   float64 = PI * 2
	HalfPi float64 = PIx2 * 0.25
	DegRad float64 = PIx2 / 360.0 // X degrees * DEG_RAD = Y radians
	RadDeg float64 = 360.0 / PIx2 // Y radians * RAD_DEG = X degrees

	// Convenience numbers.
	Large float64 = math.MaxFloat32
	Sqrt2 float64 = math.Sqrt2
	Sqrt3 float64 = 1.73205

	// Epsilon is used to distinguish when a float is close enough to a number.
	// Wikipedia: "In set theory epsilon is the limit ordinal of the sequence..."
	Epsilon float64 = 0.000001
)

func MinMax(values ...float64) (float64, float64) {
	max := values[0]
	min := values[0]
	for _, value := range values {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func Min(values ...float64) float64 {
	//max := values[0]
	min := values[0]
	for _, value := range values {
		//if max < value {
		//	max = value
		//}
		if min > value {
			min = value
		}
	}
	return min
}

func Max(values ...float64) float64 {
	max := values[0]
	for _, value := range values {
		if max < value {
			max = value
		}
	}
	return max
}

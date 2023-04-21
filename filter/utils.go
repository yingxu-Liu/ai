package filter

import "math"

// Max3 returns the maximum of three uint8 values.
func Max3(a, b, c uint8) float64 {
	if a > b {
		return math.Max(float64(a), float64(c))
	}
	return math.Max(float64(b), float64(c))
}

// Min3 returns the minimum of three uint8 values.
func Min3(a, b, c uint8) float64 {
	if a < b {
		return math.Min(float64(a), float64(c))
	}
	return math.Min(float64(b), float64(c))
}

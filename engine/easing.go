package engine

import "math"

// EaseOutQuad decelerates smoothly. t should be 0..1
func EaseOutQuad(t float64) float64 {
	return 1 - (1-t)*(1-t)
}

// EaseInOutQuad accelerates then decelerates. t should be 0..1
func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - (-2*t+2)*(-2*t+2)/2
}

// EaseOutElastic bouncy settle effect. t should be 0..1
func EaseOutElastic(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	return math.Pow(2, -10*t)*math.Sin((t*10-0.75)*(2*math.Pi/3)) + 1
}

// Lerp linearly interpolates between a and b by t (0..1)
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}


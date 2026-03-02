package engine

const (
	ScreenWidth  = 1400
	ScreenHeight = 800
	TargetFPS    = 60
	FrameDelay   = 1000 / TargetFPS
)

func Clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

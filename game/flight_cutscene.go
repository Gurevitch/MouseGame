package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// flightCutscene owns the 4-second airplane-flight scene state — the
// destination the player is flying to, the elapsed timer, and the biplane
// animation frames. Previously these were 5 flat fields on Game; extracted
// during Phase 6 so the cutscene is self-contained.
//
// Not wired through SequencePlayer because the flight carries a parameter
// (destination) that the current JSON sequence schema doesn't support, and
// the logic is tiny (~15 lines) — keeping it as a typed struct is simpler
// than adding variable substitution to sequences.
type flightCutscene struct {
	destination string
	timer       float64
	frames      []npcFrame
	frameIdx    int
	frameTimer  float64
}

const flightDurationSeconds = 4.0
const flightFrameIntervalSeconds = 0.12

// loadAirplaneFrames loads the PP-in-biplane animation frames from the
// pp_airplane sheet (4x3 grid per the existing author layout). Variable
// row lengths are tolerated so sheets with trailing short rows don't
// produce nil textures.
func loadAirplaneFrames(renderer *sdl.Renderer) []npcFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, "assets/images/player/pp_airplane.png", 4, 3, spriteInset)
	var frames []npcFrame
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			gf := grid[r][c]
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}
	return frames
}

// Start kicks off a new flight toward `dest`. Caller is responsible for
// transitioning the scene to airplane_flight after calling Start.
func (f *flightCutscene) Start(dest string) {
	f.destination = dest
	f.timer = 0
	f.frameIdx = 0
	f.frameTimer = 0
}

// Active reports whether a flight is in progress.
func (f *flightCutscene) Active() bool { return f.destination != "" }

// Update advances the timer and animation. Returns (true, dest) when the
// cutscene finishes — the caller should transitionTo(dest). Otherwise
// returns (false, "").
func (f *flightCutscene) Update(dt float64) (bool, string) {
	if !f.Active() {
		return false, ""
	}
	f.timer += dt
	f.frameTimer += dt
	if f.frameTimer >= flightFrameIntervalSeconds && len(f.frames) > 0 {
		f.frameTimer -= flightFrameIntervalSeconds
		f.frameIdx = (f.frameIdx + 1) % len(f.frames)
	}
	if f.timer < flightDurationSeconds {
		return false, ""
	}
	dest := f.destination
	f.destination = ""
	f.timer = 0
	return true, dest
}

// Draw renders the current airplane frame centered on screen, scaled 3x
// with a subtle sinusoidal bob so the biplane feels airborne.
func (f *flightCutscene) Draw(renderer *sdl.Renderer) {
	if !f.Active() || len(f.frames) == 0 {
		return
	}
	frame := f.frames[f.frameIdx%len(f.frames)]
	if frame.tex == nil {
		return
	}
	const scale = 3.0
	bob := math.Sin(f.timer*2.0) * 8
	dstW := int32(float64(frame.w) * scale)
	dstH := int32(float64(frame.h) * scale)
	dst := sdl.Rect{
		X: engine.ScreenWidth/2 - dstW/2,
		Y: int32(float64(engine.ScreenHeight)/2 - float64(dstH)/2 + bob),
		W: dstW,
		H: dstH,
	}
	renderer.Copy(frame.tex, nil, &dst)
}

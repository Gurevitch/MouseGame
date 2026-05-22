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
// pp_airplane sheet. User 2026-05-12: the real layout is **6 cols × 2
// rows** (12 frames, 295×443 each on the 1774×887 sheet). Previously
// loaded as 4×3 which cropped half-of-frame + half-of-next per cell so
// the plane rendered mis-aligned. Do not apply spriteInset here: this sheet
// has no cell grid lines, and trimming each edge clips the propeller/tail.
func loadAirplaneFrames(renderer *sdl.Renderer) []npcFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, "assets/images/player/pp_airplane.png", 6, 2, 0)
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

// Draw renders the current airplane frame centered on screen with a
// subtle sinusoidal bob so the biplane feels airborne. User 2026-05-12:
// scale dropped from 3.0 to 1.5 — the cells are already 295×443 each,
// so 1.5× gives 443×665 on screen (was 885×1329 at 3.0 which dwarfed
// the 1400×800 background).
func (f *flightCutscene) Draw(renderer *sdl.Renderer) {
	if !f.Active() || len(f.frames) == 0 {
		return
	}
	idx := f.frameIdx % len(f.frames)
	frame := f.frames[idx]
	if frame.tex == nil {
		return
	}
	const scale = 1.5
	bob := math.Sin(f.timer*2.0) * 8
	dstW := int32(float64(frame.w) * scale)
	dstH := int32(float64(frame.h) * scale)
	// User 2026-05-20: pp_airplane is a 6×2 sheet where the row 1 cells
	// have the plane drawn ~85px lower than row 0 within the same cell
	// box. Without compensation, the plane bounces top-to-bottom as the
	// animation cycles 0→11. Compensate by lifting row-1 frames (idx 6-11)
	// up by the per-cell offset, scaled to render size. Final art regen
	// (locking fuselage centerline per EXTRA_PROMPTS) will let us drop
	// this table to all-zeros.
	var rowYOffsetPx int32
	if idx >= 6 {
		v := 85.0 * float64(scale)
		rowYOffsetPx = int32(v)
	}
	dst := sdl.Rect{
		X: engine.ScreenWidth/2 - dstW/2,
		Y: int32(float64(engine.ScreenHeight)/2-float64(dstH)/2+bob) - rowYOffsetPx,
		W: dstW,
		H: dstH,
	}
	renderer.Copy(frame.tex, nil, &dst)
}

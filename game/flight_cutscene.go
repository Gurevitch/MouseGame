package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// flightCutscene owns the 4-second airplane-flight scene state - the
// destination the player is flying to, the elapsed timer, and the biplane
// animation frames. Previously these were 5 flat fields on Game; extracted
// during Phase 6 so the cutscene is self-contained.
//
// Not wired through SequencePlayer because the flight carries a parameter
// (destination) that the current JSON sequence schema doesn't support, and
// the logic is tiny (~15 lines) - keeping it as a typed struct is simpler
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
	// User playtest #18: the sheet is a 5×2 grid (counted from the art), not
	// 6×2. Loading 6 cols sliced the 307px frames at 256px, so the plane slid
	// and "jumped between rows". Load 5×2.
	grid := engine.SpriteGridFromPNGClean(renderer, "assets/images/player/pp_airplane.png", 5, 2, 0)
	var frames []npcFrame
	for r := 0; r < len(grid); r++ {
		for c := 0; c < len(grid[r]); c++ {
			gf := grid[r][c]
			// Keep the opaque box so Draw can center the PLANE itself per
			// frame - the two grid rows place the plane at different cell
			// heights, which read as "jumping between two lines" (#14).
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH})
		}
	}
	return trimBlankTail(frames)
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
// cutscene finishes - the caller should transitionTo(dest). Otherwise
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
// scale dropped from 3.0 to 1.5 - the cells are already 295×443 each,
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
	// 2026-06-11 #14: draw by the frame's OPAQUE BOX, not the full cell.
	// The plane sits at different heights inside row-0 vs row-1 cells, so
	// full-cell centering bounced it between two lines. Centering the
	// content box pins the fuselage to one spot regardless of sheet layout.
	var src *sdl.Rect
	cw, ch := frame.w, frame.h
	if frame.ow > 0 && frame.oh > 0 {
		s := sdl.Rect{X: frame.ox, Y: frame.oy, W: frame.ow, H: frame.oh}
		src = &s
		cw, ch = frame.ow, frame.oh
	}
	dstW := int32(float64(cw) * scale)
	dstH := int32(float64(ch) * scale)
	dst := sdl.Rect{
		X: engine.ScreenWidth/2 - dstW/2,
		Y: int32(float64(engine.ScreenHeight)/2 - float64(dstH)/2 + bob),
		W: dstW,
		H: dstH,
	}
	renderer.Copy(frame.tex, src, &dst)
}

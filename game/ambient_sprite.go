package game

import (
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// ambientSprite is a sheet-backed background mover — the Paris biker, the
// Jerusalem worshippers, the camp crow. Unlike `particle` (which draws
// procedural shapes: dots, triangles, rects), an ambientSprite plays an
// animated sprite strip behind the actors. It loads from a TRANSPARENT PNG
// via SpriteGridFromPNGRaw (NO color-key) so cream robes and white-striped
// shirts survive — color-keying white would punch holes in them.
//
// Each frame is drawn cropped to its opaque box (ox/oy/ow/oh) and anchored
// by a foot point (a.x, a.y = where the figure stands on the ground), so the
// big transparent margins in a 1536-wide cell don't shove the figure around.
type ambientKind int

const (
	ambientTravel ambientKind = iota // crosses the screen horizontally, then wraps
	ambientSway                      // stays put, just cycles frames in place (worshippers)
	ambientPerch                     // flies in -> perches -> holds -> flies off -> repeats (crow)
)

type ambientSprite struct {
	frames []npcFrame
	kind   ambientKind
	scale  float64

	// x, y is the foot anchor: the ground point the figure stands on.
	x, y float64

	// travel / fly speed in px/sec. Sign sets facing (negative = face left).
	vx      float64
	wrapPad float64 // travel: how far off-screen before wrapping

	// animation
	frameIdx   int
	frameTimer float64
	frameSec   float64 // seconds per frame

	// perch state machine (crow)
	perchX, perchY float64
	startX, startY float64
	state          int // 0 fly-in, 1 perched, 2 fly-off
	stateTimer     float64
	perchHold      float64 // seconds to stay perched
	flyFrames      int     // leading frames that are the wing-flap loop (rest = perched pose)
}

// loadAmbientStrip cuts a single-row transparent sprite strip into frames
// WITHOUT color-keying. Returns nil if the PNG isn't on disk yet — callers
// then no-op, so the wiring can ship ahead of the art (same pattern as the
// bird/cloud ambient sheets in scene.go).
func loadAmbientStrip(renderer *sdl.Renderer, path string, frames int) []npcFrame {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	grid := engine.SpriteGridFromPNGRaw(renderer, path, frames, 1)
	return framesFromGrid(grid, frames, 1, path)
}

// --- constructors ---

// newAmbientBiker is the Parisian cyclist that drifts past near the Eiffel
// (item 9 / §AMB2). 8-frame ride loop, faces right, rides left->right.
func newAmbientBiker(renderer *sdl.Renderer, startX, groundY, vx, scale float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStrip(renderer, "assets/images/locations/paris/npc/outside/biker.png", 8),
		kind:     ambientTravel,
		x:        startX,
		y:        groundY,
		vx:       vx,
		scale:    scale,
		wrapPad:  240,
		frameSec: 0.09,
	}
}

// newAmbientWorshippers is the cluster of tiny figures swaying at the Western
// Wall (item 9 / §AMB1). 6-frame in-place sway loop.
func newAmbientWorshippers(renderer *sdl.Renderer, x, y, scale float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStrip(renderer, "assets/images/locations/jerusalem/npc/wall/people_pray.png", 6),
		kind:     ambientSway,
		x:        x,
		y:        y,
		scale:    scale,
		frameSec: 0.4,
	}
}

// newAmbientCrow is the camp crow that flaps in, lands on the camp sign, sits
// a beat, then flaps away — and repeats. Art is pending (assets/images/ambient/
// crow.png, 8-frame: 0-5 flap, 6-7 perched). No-ops until the PNG lands.
func newAmbientCrow(renderer *sdl.Renderer, perchX, perchY float64) *ambientSprite {
	startX := -120.0
	startY := perchY - 160
	return &ambientSprite{
		frames:    loadAmbientStrip(renderer, "assets/images/ambient/crow.png", 8),
		kind:      ambientPerch,
		scale:     0.5,
		x:         startX,
		y:         startY,
		startX:    startX,
		startY:    startY,
		perchX:    perchX,
		perchY:    perchY,
		vx:        150,
		frameSec:  0.08,
		flyFrames: 6,
		perchHold: 2.5,
	}
}

// --- update ---

func (a *ambientSprite) update(dt float64) {
	if len(a.frames) == 0 {
		return
	}
	switch a.kind {
	case ambientTravel:
		a.x += a.vx * dt
		if a.vx > 0 && a.x > float64(engine.ScreenWidth)+a.wrapPad {
			a.x = -a.wrapPad
		} else if a.vx < 0 && a.x < -a.wrapPad {
			a.x = float64(engine.ScreenWidth) + a.wrapPad
		}
		a.advanceFrame(dt, 0, len(a.frames))
	case ambientSway:
		a.advanceFrame(dt, 0, len(a.frames))
	case ambientPerch:
		a.updatePerch(dt)
	}
}

// advanceFrame ticks the frame index, wrapping inside [lo, hi).
func (a *ambientSprite) advanceFrame(dt float64, lo, hi int) {
	if hi <= lo {
		return
	}
	a.frameTimer += dt
	if a.frameTimer >= a.frameSec {
		a.frameTimer = 0
		a.frameIdx++
		if a.frameIdx < lo || a.frameIdx >= hi {
			a.frameIdx = lo
		}
	}
}

func (a *ambientSprite) updatePerch(dt float64) {
	flyHi := a.flyFrames
	if flyHi <= 0 || flyHi > len(a.frames) {
		flyHi = len(a.frames)
	}
	perchFrame := len(a.frames) - 1 // last cell = the standing/perched pose

	switch a.state {
	case 0: // gliding in toward the perch
		a.advanceFrame(dt, 0, flyHi)
		dx := a.perchX - a.x
		step := a.vx * dt
		if afAbs(dx) <= step {
			a.x = a.perchX
			a.y = a.perchY
			a.state = 1
			a.stateTimer = 0
			a.frameIdx = perchFrame
		} else {
			a.x += afSign(dx) * step
			a.y += (a.perchY - a.y) * afMin(1, dt*2.5) // ease down onto the perch
		}
	case 1: // perched, holding
		a.frameIdx = perchFrame
		a.stateTimer += dt
		if a.stateTimer >= a.perchHold {
			a.state = 2
		}
	case 2: // flapping away (continues right, rising)
		a.advanceFrame(dt, 0, flyHi)
		a.x += a.vx * dt
		a.y -= 45 * dt
		if a.x > float64(engine.ScreenWidth)+120 {
			a.x = a.startX
			a.y = a.startY
			a.state = 0
		}
	}
}

// --- draw ---

func (a *ambientSprite) draw(renderer *sdl.Renderer) {
	if len(a.frames) == 0 {
		return
	}
	idx := a.frameIdx % len(a.frames)
	if idx < 0 {
		idx += len(a.frames)
	}
	f := a.frames[idx]
	if f.tex == nil || f.ow <= 0 || f.oh <= 0 {
		return
	}
	src := sdl.Rect{X: f.ox, Y: f.oy, W: f.ow, H: f.oh}
	dstW := int32(float64(f.ow) * a.scale)
	dstH := int32(float64(f.oh) * a.scale)
	dst := sdl.Rect{X: int32(a.x) - dstW/2, Y: int32(a.y) - dstH, W: dstW, H: dstH}
	flip := sdl.FLIP_NONE
	if a.vx < 0 {
		flip = sdl.FLIP_HORIZONTAL
	}
	renderer.CopyEx(f.tex, &src, &dst, 0, nil, flip)
}

// small float helpers (kept local so they don't collide with the Go 1.21
// builtin min/max used elsewhere in the package).
func afAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

func afSign(v float64) float64 {
	if v < 0 {
		return -1
	}
	return 1
}

func afMin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

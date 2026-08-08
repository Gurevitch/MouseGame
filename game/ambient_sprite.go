package game

import (
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// ambientSprite is a sheet-backed background mover - the Paris biker, the
// Jerusalem worshippers, the camp crow. Unlike `particle` (which draws
// procedural shapes: dots, triangles, rects), an ambientSprite plays an
// animated sprite strip behind the actors. It loads from a TRANSPARENT PNG
// via SpriteGridFromPNGRaw (NO color-key) so cream robes and white-striped
// shirts survive - color-keying white would punch holes in them.
//
// Each frame is drawn cropped to its opaque box (ox/oy/ow/oh) and anchored
// by a foot point (a.x, a.y = where the figure stands on the ground), so the
// big transparent margins in a 1536-wide cell don't shove the figure around.
type ambientKind int

const (
	ambientTravel ambientKind = iota // crosses the screen horizontally, then wraps
	ambientSway                      // stays put, just cycles frames in place (worshippers)
	ambientPerch                     // flies in -> perches -> holds -> flies off -> repeats (crow)
	ambientFlyOff                    // starts at (x,y), flaps, drifts up+away, self-clears (PR#29 pot pigeon)
	ambientFall                      // drifts DOWN (+ gentle sideways), wraps to the top - a live falling loop (sakura/leaves)
)

type ambientSprite struct {
	frames []npcFrame
	kind   ambientKind
	scale  float64

	// Interactive crossing (2026-06-11 #16, the Paris biker): when onClick is
	// set, the sprite is clickable. paused freezes movement (he brakes
	// mid-street); the click handler resumes it when its dialog ends.
	onClick  func()
	paused   bool
	lastRect sdl.Rect
	// pauseFrame (>0): frame index shown while paused - the §BK1 biker sheet
	// has dedicated BRAKED poses in frames 7-8 (index 6). 0 = freeze current.
	pauseFrame int
	// loopFrames (>0): the ride/sway loop cycles only frames [0,loopFrames) -
	// keeps the biker's braked poses out of the riding animation. 0 = all.
	loopFrames int

	// x, y is the foot anchor: the ground point the figure stands on.
	x, y float64

	// foreground (2026-08-07 #8): draw AFTER the footY-sorted actors instead
	// of in the background ambient pass — the Paris biker crosses the street
	// in FRONT of everything (he was riding behind the flower pot).
	foreground bool

	// travel / fly speed in px/sec. Sign sets facing (negative = face left).
	vx      float64
	vyUp    float64 // ambientFlyOff: upward drift px/sec
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

// loadAmbientStripKeyed is for ambient sheets that shipped with a BAKED
// white background instead of transparency (the current biker.png rendered
// as a white box riding the street, 2026-06-11 #16). The edge-connected key
// strips the background while keeping interior whites. Tol 24 (2026-06-12
// #12): tol 8 left a white halo around the biker's anti-aliased outline.
func loadAmbientStripKeyed(renderer *sdl.Renderer, path string, frames int) []npcFrame {
	return loadAmbientStripKeyedTol(renderer, path, frames, 24)
}

// loadAmbientStripKeyedTol is loadAmbientStripKeyed with a caller-chosen
// connected-key tolerance. PR#7: the biker uses tol 40 to shave the last
// anti-aliased edge halo. (Its leftover INTERIOR white - bike-frame
// triangles, wheel gaps - is enclosed by the outline so the edge flood can't
// reach it; a full fix needs a transparent-bg re-roll, queued at §BK2.)
func loadAmbientStripKeyedTol(renderer *sdl.Renderer, path string, frames int, tol uint8) []npcFrame {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	grid := engine.SpriteGridFromPNGCleanConnectedTol(renderer, path, frames, 1, 0, tol)
	return framesFromGrid(grid, frames, 1, path)
}

// loadAmbientStripGlobalKeyed removes the sampled background colour everywhere,
// including enclosed gaps. Use only when the keyed colour is absent from the
// foreground art.
func loadAmbientStripGlobalKeyed(renderer *sdl.Renderer, path string, frames int) []npcFrame {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	grid := engine.SpriteGridFromPNGClean(renderer, path, frames, 1, 0)
	return framesFromGrid(grid, frames, 1, path)
}

// loadAmbientStripGlobalKeyedTol is loadAmbientStripGlobalKeyed with a
// caller-chosen tolerance (see newAmbientSwayGlobalTol).
func loadAmbientStripGlobalKeyedTol(renderer *sdl.Renderer, path string, frames int, tol uint8) []npcFrame {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	grid := engine.SpriteGridFromPNGCleanTol(renderer, path, frames, 1, 0, tol)
	return framesFromGrid(grid, frames, 1, path)
}

// containsPoint hit-tests the sprite's last drawn rect (clickable ambients).
func (a *ambientSprite) containsPoint(x, y int32) bool {
	if a.onClick == nil || a.lastRect.W <= 0 {
		return false
	}
	pt := sdl.Point{X: x, Y: y}
	return pt.InRect(&a.lastRect)
}

// --- constructors ---

// newAmbientBiker is the Parisian cyclist that drifts past near the Eiffel
// (item 9 / §AMB2). 8-frame ride loop, faces right, rides left->right.
func newAmbientBiker(renderer *sdl.Renderer, startX, groundY, vx, scale float64) *ambientSprite {
	return &ambientSprite{
		// 2026-06-21: biker.png is NOT actually transparent - it ships with a
		// near-white OPAQUE background (every pixel alpha=255, corners ~233-253),
		// so the old RAW load drew a white box around him. Use the EDGE-CONNECTED
		// white key (tol 40): it floods the background away from the sheet edges
		// while leaving his enclosed white striped shirt intact.
		// 2026-07-18 (user: "biker is missing colors"): the sheet is now on the
		// flat BLUE bg, and the old tol-40 connected key bled into his light-blue
		// clothing. The bg is uniform, so the tight GLOBAL sampled key (tol 8)
		// strips it — including the enclosed wheel gaps — without touching him.
		frames:   loadAmbientStripGlobalKeyed(renderer, "assets/images/locations/paris/npc/outside/biker.png", 8),
		kind:     ambientTravel,
		x:        startX,
		y:        groundY,
		vx:       vx,
		scale:    scale,
		wrapPad:  240,
		frameSec: 0.09,
		// §BK1 sheet layout: frames 1-6 ride loop, 7-8 braked poses.
		loopFrames: 6,
		pauseFrame: 6,
	}
}

// newAmbientPigeonFlyUp (PR#29): the flower-pot pigeon lifts off when Pierre
// shoos it, flapping up and to the right toward the rooftops, then clears
// itself. Reuses the existing `npc_pierre_pigeon_lands.png` - an 8-frame takeoff
// strip (perched → flapping → climbing up-right), which is exactly a fly-up.
// 2026-06-21: this sheet is opaque near-white (NOT transparent), so it loads
// through the edge-connected white key (tol 40) - otherwise it drew a white box
// around the pigeon. Returns nil if the art is ever absent.
func newAmbientPigeonFlyUp(renderer *sdl.Renderer, x, y float64) *ambientSprite {
	frames := loadAmbientStripKeyedTol(renderer, "assets/images/locations/paris/npc/outside/npc_pierre_pigeon_lands.png", 8, 40)
	if len(frames) == 0 {
		return nil
	}
	return &ambientSprite{
		frames:   frames,
		kind:     ambientFlyOff,
		scale:    0.35,
		x:        x,
		y:        y,
		vx:       90,  // drift right
		vyUp:     130, // rise toward the rooftops
		frameSec: 0.08,
	}
}

// newAmbientWorshippers is the cluster of tiny figures swaying at the Western
// Wall (item 9 / §AMB1). 6-frame in-place sway loop.
// newAmbientWorshippersSheet is newAmbientWorshippers with a caller-chosen
// sheet (2026-07-18 user #26: the entrance crowd alternates two groups ABAB).
// 2026-07-23 (user #24): GLOBAL sampled key — the flat-blue people_pray2
// re-roll has no background-colored detail on the figures (tallit stripes are
// navy), and the global key also clears any enclosed pockets between them.
func newAmbientWorshippersSheet(renderer *sdl.Renderer, sheet string, x, y, scale float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStripGlobalKeyed(renderer, sheet, 6),
		kind:     ambientSway,
		x:        x,
		y:        y,
		scale:    scale,
		frameSec: 0.5,
	}
}

// newAmbientSway is the generic in-place looping flavor figure (retro plan
// #5, 2026-06-12): an N-frame strip cycling at frameSec, anchored at (x, y),
// no movement. Used for the Paris street-density ambients (§AMB5/§AMB6);
// any future "someone doing something along the walk line" reuses this
// instead of a bespoke constructor. Keyed load - these sheets ship with a
// baked white background like the biker's.
func newAmbientSway(renderer *sdl.Renderer, sheet string, cols int, x, y, scale, frameSec float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStripKeyed(renderer, sheet, cols),
		kind:     ambientSway,
		x:        x,
		y:        y,
		scale:    scale,
		frameSec: frameSec,
	}
}

// newAmbientSwayGlobal (2026-07-25 user): newAmbientSway with the GLOBAL
// sampled key — for the flat-blue §-prompt landscape props (Louvre/Japan)
// whose enclosed background pockets the connected key preserves ("blue
// spots around the couple"). Their canvas rule bans bg-coloured detail on
// the figures, so the global key is always safe here.
func newAmbientSwayGlobal(renderer *sdl.Renderer, sheet string, cols int, x, y, scale, frameSec float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStripGlobalKeyed(renderer, sheet, cols),
		kind:     ambientSway,
		x:        x,
		y:        y,
		scale:    scale,
		frameSec: frameSec,
	}
}

// newAmbientSwayGlobalTol is newAmbientSwayGlobal with a caller-chosen key
// tolerance. Tol 24 shaves the anti-aliased seams the tol-8 global key leaves
// where the flat blue meets the prop's outline (the teahouse kettle's "blue
// spot", 2026-07-26).
func newAmbientSwayGlobalTol(renderer *sdl.Renderer, sheet string, cols int, x, y, scale, frameSec float64, tol uint8) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStripGlobalKeyedTol(renderer, sheet, cols, tol),
		kind:     ambientSway,
		x:        x,
		y:        y,
		scale:    scale,
		frameSec: frameSec,
	}
}

// newAmbientPropKeyed is newAmbientProp with a global sampled colour key.
// Use for props whose key colour also appears in enclosed background gaps.
func newAmbientPropKeyed(renderer *sdl.Renderer, sheet string, x, y, scale float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStripGlobalKeyed(renderer, sheet, 1),
		kind:     ambientSway,
		x:        x,
		y:        y,
		scale:    scale,
		frameSec: 999,
	}
}

// newAmbientCrow is the camp crow that flaps in, lands on the camp sign, sits
// a beat, then flaps away - and repeats. Art landed 2026-07-14 (8-frame:
// 0-5 flap, 6-7 perched) with a BAKED-IN checkerboard background, so it must
// go through the KEYED loader (the checker squares are near-white and
// edge-connected — same treatment as the ramen stall props).
func newAmbientCrow(renderer *sdl.Renderer, perchX, perchY float64) *ambientSprite {
	startX := -120.0
	startY := perchY - 160
	return &ambientSprite{
		// 2026-07-24 (user #18): GLOBAL key — the connected key left the
		// enclosed checker pocket between the perched crow's legs (the
		// corner samples catch both checker colours, so global clears it).
		frames:    loadAmbientStripGlobalKeyed(renderer, "assets/images/ambient/crow.png", 8),
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
	if a.paused {
		// Braked mid-street (clicked biker) - hold position until the click
		// handler's dialog resumes us, showing the dedicated braked pose
		// when the sheet has one (§BK1 frames 7-8).
		if a.pauseFrame > 0 && a.pauseFrame < len(a.frames) {
			a.frameIdx = a.pauseFrame
		}
		return
	}
	loopHi := len(a.frames)
	if a.loopFrames > 0 && a.loopFrames < loopHi {
		loopHi = a.loopFrames
	}
	switch a.kind {
	case ambientTravel:
		a.x += a.vx * dt
		if a.vx > 0 && a.x > float64(engine.ScreenWidth)+a.wrapPad {
			a.x = -a.wrapPad
		} else if a.vx < 0 && a.x < -a.wrapPad {
			a.x = float64(engine.ScreenWidth) + a.wrapPad
		}
		a.advanceFrame(dt, 0, loopHi)
	case ambientSway:
		a.advanceFrame(dt, 0, loopHi)
	case ambientPerch:
		a.updatePerch(dt)
	case ambientFlyOff:
		// PR#29: pigeon lifts off the flower pot and drifts up+away toward the
		// rooftops, then clears itself (nil frames → update/draw no-op).
		a.x += a.vx * dt
		a.y -= a.vyUp * dt
		a.advanceFrame(dt, 0, loopHi)
		if a.y < -140 || a.x > float64(engine.ScreenWidth)+140 {
			a.frames = nil
		}
	case ambientFall:
		// vyUp carries the DOWNWARD speed here; vx the sideways drift. Wraps back
		// to the top (and across) so the leaf keeps falling forever.
		a.y += a.vyUp * dt
		a.x += a.vx * dt
		if a.y > float64(engine.ScreenHeight)+a.wrapPad {
			a.y = -a.wrapPad
		}
		if a.x > float64(engine.ScreenWidth)+a.wrapPad {
			a.x = -a.wrapPad
		} else if a.x < -a.wrapPad {
			a.x = float64(engine.ScreenWidth) + a.wrapPad
		}
		a.advanceFrame(dt, 0, loopHi)
	}
}

// newAmbientLeafFall is a sprite leaf/petal that drifts down the screen and
// wraps back to the top - a "live" falling loop (Japan ramen-store tree). The
// sheet is a short N-frame flutter cycle; missing art no-ops (nil frames).
func newAmbientLeafFall(renderer *sdl.Renderer, sheet string, cols int, x, y, scale, fallSpeed, driftX, frameSec float64) *ambientSprite {
	return &ambientSprite{
		frames:   loadAmbientStripKeyed(renderer, sheet, cols),
		kind:     ambientFall,
		x:        x,
		y:        y,
		scale:    scale,
		vyUp:     fallSpeed,
		vx:       driftX,
		wrapPad:  90,
		frameSec: frameSec,
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
	flip := sdl.FLIP_NONE
	if a.vx < 0 {
		flip = sdl.FLIP_HORIZONTAL
	}
	// 2026-08-08 D-2: anchor on the engine's per-frame FEET anchors
	// (fcx/fry, deadband-stabilized like the NPCs) instead of the opaque
	// bounding-box centre — a drifting steam plume / raised tail shifts the
	// bbox per frame, which slid the kettle (±9px), the cat, and the other
	// props side to side every frame. Falls back to the old bbox anchor
	// when the loader left the anchors unset.
	ax := float64(f.ow) / 2 // anchor offset inside the opaque box
	ay := float64(f.oh)
	if f.fcx > 0 {
		ax = float64(f.fcx - f.ox)
	}
	if f.fry > f.oy {
		ay = float64(f.fry - f.oy)
	}
	if flip == sdl.FLIP_HORIZONTAL {
		ax = float64(f.ow) - ax
	}
	dst := sdl.Rect{
		X: int32(a.x - ax*a.scale),
		Y: int32(a.y - ay*a.scale),
		W: dstW, H: dstH,
	}
	renderer.CopyEx(f.tex, &src, &dst, 0, nil, flip)
	a.lastRect = dst
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

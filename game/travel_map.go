package game

import (
	"fmt"
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type travelLocation struct {
	id           string
	name         string
	scene        string
	pinX         int32
	pinY         int32
	unlocked     bool
	info         string   // legacy one-liner; facts is now the canonical list
	facts        []string // multi-line facts shown as bullet paragraphs in the info panel
	landmarkPath string
	landmarkTex  *sdl.Texture
	landmarkW    int32
	landmarkH    int32
	// audio is an optional path to a voice clip describing the place.
	// If set and the file exists, it plays on click alongside the info popup.
	// Empty string = popup only (current default).
	audio string
	// relevantWhen is an expression (see game/npc_rules.go:evalCondition)
	// that evaluates to true when THIS location is the current story target.
	// Drives the "bright glow" tier on the map. Empty string = never relevant
	// (stays at the unlocked-but-not-current tier even when unlocked).
	relevantWhen string
}

type travelMap struct {
	locations []travelLocation
	renderer  *sdl.Renderer
	bgTex     *sdl.Texture
	// visible + returnScene were previously flags on Game (showTravelMap,
	// travelMapFrom). Moved here during the Phase 6 god-object collapse —
	// they're display state that belongs with the widget, not the game loop.
	visible     bool
	returnScene string
	// game is a back-reference set by Game.New so drawOverlay can evaluate
	// each location's relevantWhen expression against VarStore each frame.
	// nil during construction (before Game.New finishes); draw handles nil
	// by treating all locations as not-currently-relevant.
	game *Game
	// panelLoc is the location whose info panel is currently overlaid on
	// the map. nil = no panel. Set by openInfoPanel, cleared by click or Esc.
	panelLoc *travelLocation
}

// Visible reports whether the travel-map modal is currently open.
func (tm *travelMap) Visible() bool { return tm.visible }

// ReturnScene is the scene name the player came from (flight cutscenes
// deposit the player back here if they close the map without travelling).
func (tm *travelMap) ReturnScene() string { return tm.returnScene }

// Show opens the travel-map modal and records which scene to return to.
func (tm *travelMap) Show(fromScene string) {
	tm.visible = true
	tm.returnScene = fromScene
}

// Hide closes the travel-map modal.
func (tm *travelMap) Hide() { tm.visible = false }

// Toggle flips visibility. When opening, records the fromScene as the
// return-to scene. When closing, leaves returnScene intact so callers can
// still read it during the close animation frame.
func (tm *travelMap) Toggle(fromScene string) {
	tm.visible = !tm.visible
	if tm.visible {
		tm.returnScene = fromScene
	}
}

// travelMapDataPath is the canonical location of the authoritative travel-
// map data. Moved out of the constructor so save/load + tests can override
// it if ever needed. Camp Chilly Wa Wa is intentionally NOT in this file —
// see docs/FIXME.md "Deferred to follow-up" for why and when to re-add.
const travelMapDataPath = "assets/data/travel_map.json"

func newTravelMap(renderer *sdl.Renderer) *travelMap {
	tm := &travelMap{renderer: renderer}

	locs, err := loadTravelLocations(travelMapDataPath)
	if err != nil {
		fmt.Printf("travel_map: falling back to empty locations list: %v\n", err)
	} else {
		tm.locations = locs
	}

	// Load landmark textures. SafeTextureFromPNGRaw preserves alpha from the
	// source PNG and does NOT run the corner-sample color-key, which was
	// eating pale colors inside landmarks (Eiffel steel, colosseum stone).
	// Landmark PNGs are authored with transparent backgrounds already.
	for i := range tm.locations {
		if tm.locations[i].landmarkPath != "" {
			tex, w, h := engine.SafeTextureFromPNGRaw(renderer, tm.locations[i].landmarkPath)
			if tex != nil {
				tex.SetBlendMode(sdl.BLENDMODE_BLEND)
				tm.locations[i].landmarkTex = tex
				tm.locations[i].landmarkW = w
				tm.locations[i].landmarkH = h
			}
		}
	}

	// Try to load globe image, fall back to procedural map
	globeTex, _, _ := engine.SafeTextureFromPNGRaw(renderer, "assets/images/ui/travel_globe.png")
	if globeTex != nil {
		tm.bgTex = globeTex
	} else {
		tm.bgTex = tm.generateMapTexture()
	}
	return tm
}

// attachGame stores a Game back-reference so drawOverlay can evaluate each
// location's relevantWhen expression against VarStore. Called from Game.New
// after the Game struct is wired up.
func (tm *travelMap) attachGame(g *Game) { tm.game = g }

// isRelevant returns true when the location's relevantWhen expression
// evaluates true in the current VarStore state. nil game (pre-wire) returns
// false so nothing flashes before state is ready.
func (tm *travelMap) isRelevant(loc *travelLocation) bool {
	if tm.game == nil || loc.relevantWhen == "" {
		return false
	}
	return evalCondition(loc.relevantWhen, ruleContext{game: tm.game})
}

func (tm *travelMap) generateMapTexture() *sdl.Texture {
	w := int32(engine.ScreenWidth)
	h := int32(engine.ScreenHeight)
	surface, err := sdl.CreateRGBSurface(0, w, h, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		panic(err)
	}
	defer surface.Free()
	f := surface.Format

	// Ocean gradient
	for y := int32(0); y < h; y++ {
		t := float64(y) / float64(h)
		r := uint8(30 + t*20)
		g := uint8(60 + t*30)
		b := uint8(120 + t*40)
		surface.FillRect(&sdl.Rect{X: 0, Y: y, W: w, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}

	// Decorative wave lines
	for i := 0; i < 20; i++ {
		baseY := 100 + int32(i*35)
		for x := int32(0); x < w; x += 3 {
			wy := baseY + int32(math.Sin(float64(x)*0.015+float64(i)*1.5)*8)
			surface.FillRect(&sdl.Rect{X: x, Y: wy, W: 2, H: 1},
				sdl.MapRGBA(f, 50, 80, 150, 40))
		}
	}

	// --- Landmasses ---

	// North America (Camp Chilly Wa Wa area) — left side
	landColor := sdl.MapRGBA(f, 85, 130, 70, 255)
	landDark := sdl.MapRGBA(f, 70, 110, 55, 255)

	// Main continent blob
	fillEllipseTM(surface, 350, 320, 200, 140, landColor)
	fillEllipseTM(surface, 280, 260, 120, 80, landColor)
	fillEllipseTM(surface, 420, 380, 100, 70, landDark)
	fillEllipseTM(surface, 300, 360, 80, 50, landDark)

	// Europe (London & Paris area) — right side
	fillEllipseTM(surface, 750, 270, 180, 110, landColor)
	fillEllipseTM(surface, 850, 230, 100, 70, landColor)
	fillEllipseTM(surface, 680, 310, 80, 60, landDark)
	fillEllipseTM(surface, 900, 280, 120, 80, landDark)
	fillEllipseTM(surface, 1050, 300, 130, 90, landColor)
	fillEllipseTM(surface, 1000, 260, 80, 50, landDark)

	// Small islands
	fillEllipseTM(surface, 560, 400, 25, 15, landColor)
	fillEllipseTM(surface, 580, 450, 15, 10, landDark)

	// Mountain details on continents
	mountColor := sdl.MapRGBA(f, 100, 90, 70, 255)
	snowColor := sdl.MapRGBA(f, 220, 220, 230, 200)
	drawMountainTM(surface, f, 310, 290, 20, mountColor, snowColor)
	drawMountainTM(surface, f, 370, 310, 15, mountColor, snowColor)
	drawMountainTM(surface, f, 800, 250, 18, mountColor, snowColor)
	drawMountainTM(surface, f, 1020, 280, 16, mountColor, snowColor)

	// Forest dots near Camp
	forestColor := sdl.MapRGBA(f, 40, 80, 35, 200)
	for i := 0; i < 30; i++ {
		fx := int32(270 + math.Sin(float64(i)*2.1)*70)
		fy := int32(300 + math.Cos(float64(i)*1.7)*50)
		fillEllipseTM(surface, fx, fy, 6, 4, forestColor)
	}

	// City dots near London
	cityColor := sdl.MapRGBA(f, 160, 150, 140, 180)
	for i := 0; i < 8; i++ {
		cx := int32(680 + math.Sin(float64(i)*3.1)*30)
		cy := int32(265 + math.Cos(float64(i)*2.3)*20)
		surface.FillRect(&sdl.Rect{X: cx, Y: cy, W: 4, H: 6}, cityColor)
	}

	// City dots near Paris
	for i := 0; i < 6; i++ {
		cx := int32(1030 + math.Sin(float64(i)*2.7)*25)
		cy := int32(305 + math.Cos(float64(i)*2.0)*15)
		surface.FillRect(&sdl.Rect{X: cx, Y: cy, W: 4, H: 5}, cityColor)
	}

	// Title banner area at top
	bannerColor := sdl.MapRGBA(f, 40, 30, 20, 180)
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: w, H: 70}, bannerColor)
	surface.FillRect(&sdl.Rect{X: 0, Y: 70, W: w, H: 3}, sdl.MapRGBA(f, 180, 150, 80, 200))

	// Bottom decorative border
	surface.FillRect(&sdl.Rect{X: 0, Y: h - 60, W: w, H: 60}, bannerColor)
	surface.FillRect(&sdl.Rect{X: 0, Y: h - 63, W: w, H: 3}, sdl.MapRGBA(f, 180, 150, 80, 200))

	// Compass rose (simple)
	compassX, compassY := int32(1250), int32(600)
	compassColor := sdl.MapRGBA(f, 200, 180, 120, 200)
	surface.FillRect(&sdl.Rect{X: compassX, Y: compassY - 30, W: 2, H: 60}, compassColor)
	surface.FillRect(&sdl.Rect{X: compassX - 30, Y: compassY, W: 60, H: 2}, compassColor)
	// N marker
	surface.FillRect(&sdl.Rect{X: compassX - 3, Y: compassY - 35, W: 8, H: 8},
		sdl.MapRGBA(f, 220, 50, 40, 220))

	tex, err := tm.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	return tex
}

func (tm *travelMap) drawOverlay(renderer *sdl.Renderer, font *engine.BitmapFont, mx, my int32) {
	ticks := float64(sdl.GetTicks())

	for i := range tm.locations {
		loc := &tm.locations[i]
		px, py := loc.pinX, loc.pinY

		switch {
		case loc.unlocked && tm.isRelevant(loc):
			tm.drawRelevantPin(renderer, font, loc, px, py, ticks, mx, my)
		case loc.unlocked:
			tm.drawUnlockedIdlePin(renderer, font, loc, px, py, ticks, mx, my)
		default:
			tm.drawLockedPin(renderer, font, loc, px, py)
		}
	}

	// Instructions at bottom
	font.DrawText(renderer, "Click a glowing destination to travel  |  Click a locked pin for info  |  ESC to close",
		engine.ScreenWidth/2-420, engine.ScreenHeight-35, 2,
		sdl.Color{R: 170, G: 160, B: 140, A: 200})
}

// drawRelevantPin renders the bright-gold "go here next" state: full-size
// landmark, big pulsing glow, always-on label. Only one pin should be in
// this state at a time in a well-authored story flow.
func (tm *travelMap) drawRelevantPin(renderer *sdl.Renderer, font *engine.BitmapFont, loc *travelLocation, px, py int32, ticks float64, mx, my int32) {
	glowPulse := 0.55 + 0.45*math.Sin(ticks*0.005)
	for r := int32(58); r > 8; r -= 2 {
		a := uint8(float64(68-r) * glowPulse)
		renderer.SetDrawColor(255, 220, 60, a)
		for dy := -r; dy <= r; dy++ {
			hw := int32(math.Sqrt(float64(r*r - dy*dy)))
			renderer.FillRect(&sdl.Rect{X: px - hw, Y: py + dy, W: hw * 2, H: 1})
		}
	}
	tm.drawLandmark(renderer, loc, px, py, 70, 255, 255)

	nameW := font.TextWidth(loc.name, 2)
	labelX := px - nameW/2
	labelY := py - 50
	renderer.SetDrawColor(30, 25, 18, 220)
	renderer.FillRect(&sdl.Rect{X: labelX - 8, Y: labelY - 4, W: nameW + 16, H: 24})
	renderer.SetDrawColor(220, 190, 60, 220)
	renderer.FillRect(&sdl.Rect{X: labelX - 8, Y: labelY - 4, W: nameW + 16, H: 2})
	renderer.FillRect(&sdl.Rect{X: labelX - 8, Y: labelY + 18, W: nameW + 16, H: 2})

	hoverRect := sdl.Rect{X: px - 50, Y: py - 55, W: 100, H: 110}
	pt := sdl.Point{X: mx, Y: my}
	if pt.InRect(&hoverRect) {
		renderer.SetDrawColor(255, 255, 180, 30)
		renderer.FillRect(&sdl.Rect{X: px - 55, Y: py - 55, W: 110, H: 110})
		font.DrawText(renderer, loc.name, labelX, labelY, 2, sdl.Color{R: 255, G: 255, B: 200, A: 255})
	} else {
		font.DrawText(renderer, loc.name, labelX, labelY, 2, sdl.Color{R: 255, G: 240, B: 160, A: 255})
	}
}

// drawUnlockedIdlePin is the "visited / not currently relevant" state:
// landmark at full color but smaller, a soft low-alpha halo, and the label
// only shows on hover. Keeps previously-unlocked cities discoverable without
// competing for attention with the current story target.
func (tm *travelMap) drawUnlockedIdlePin(renderer *sdl.Renderer, font *engine.BitmapFont, loc *travelLocation, px, py int32, ticks float64, mx, my int32) {
	haloPulse := 0.35 + 0.25*math.Sin(ticks*0.003)
	for r := int32(34); r > 12; r -= 3 {
		a := uint8(float64(40-r) * haloPulse)
		renderer.SetDrawColor(200, 200, 220, a)
		for dy := -r; dy <= r; dy++ {
			hw := int32(math.Sqrt(float64(r*r - dy*dy)))
			renderer.FillRect(&sdl.Rect{X: px - hw, Y: py + dy, W: hw * 2, H: 1})
		}
	}
	tm.drawLandmark(renderer, loc, px, py, 55, 255, 255)

	hoverRect := sdl.Rect{X: px - 45, Y: py - 40, W: 90, H: 90}
	pt := sdl.Point{X: mx, Y: my}
	if pt.InRect(&hoverRect) {
		nameW := font.TextWidth(loc.name, 2)
		labelX := px - nameW/2
		labelY := py - 45
		renderer.SetDrawColor(30, 25, 18, 200)
		renderer.FillRect(&sdl.Rect{X: labelX - 6, Y: labelY - 3, W: nameW + 12, H: 22})
		font.DrawText(renderer, loc.name, labelX, labelY, 2, sdl.Color{R: 220, G: 220, B: 240, A: 255})
	}
}

// drawLockedPin renders the gray, dimmed state used for cities that aren't
// unlocked yet. Clicking still opens an info popup.
func (tm *travelMap) drawLockedPin(renderer *sdl.Renderer, font *engine.BitmapFont, loc *travelLocation, px, py int32) {
	if loc.landmarkTex != nil {
		tm.drawLandmark(renderer, loc, px, py, 45, 180, 180)
	} else {
		renderer.SetDrawColor(100, 95, 80, 150)
		for dy := int32(-4); dy <= 4; dy++ {
			hw := int32(math.Sqrt(float64(16 - dy*dy)))
			renderer.FillRect(&sdl.Rect{X: px - hw, Y: py + dy, W: hw * 2, H: 1})
		}
	}

	nameW := font.TextWidth(loc.name, 1)
	labelX := px - nameW/2
	labelY := py + 30
	font.DrawText(renderer, loc.name, labelX, labelY, 1, sdl.Color{R: 140, G: 130, B: 110, A: 160})
}

// drawLandmark is the shared landmark-blit helper. `targetH` sets the on-
// screen height (scaling preserves aspect); `colorMod` tints and `alphaMod`
// sets translucency. 255/255 = untouched colors (relevant/idle); 180/180 =
// desaturated grayed (locked).
func (tm *travelMap) drawLandmark(renderer *sdl.Renderer, loc *travelLocation, px, py int32, targetH int32, colorMod, alphaMod uint8) {
	if loc.landmarkTex == nil {
		return
	}
	scale := float64(targetH) / float64(loc.landmarkH)
	dstW := int32(float64(loc.landmarkW) * scale)
	dstH := targetH
	dstX := px - dstW/2
	dstY := py - dstH/2
	loc.landmarkTex.SetColorMod(colorMod, colorMod, colorMod)
	loc.landmarkTex.SetAlphaMod(alphaMod)
	renderer.Copy(loc.landmarkTex, nil, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
	// Reset so subsequent draws aren't tinted.
	loc.landmarkTex.SetColorMod(255, 255, 255)
	loc.landmarkTex.SetAlphaMod(255)
}

// Hit rectangles center on the landmark sprite with enough padding that
// adjacent pins in Europe don't overlap. 90x110 covers the ~70px landmark
// plus a bit of its glow, and deliberately excludes the label sitting 50px
// above the pin so a click on "Rome" doesn't bleed up into "Paris".
func (tm *travelMap) pinHitRect(loc *travelLocation) sdl.Rect {
	return sdl.Rect{X: loc.pinX - 45, Y: loc.pinY - 55, W: 90, H: 110}
}

// distanceSqFromPin returns the squared distance from (mx, my) to the
// location's pin center. Used to tie-break when two hit rects overlap —
// the click snaps to the closest pin, not the first one in slice order.
func (tm *travelMap) distanceSqFromPin(loc *travelLocation, mx, my int32) int64 {
	dx := int64(mx - loc.pinX)
	dy := int64(my - loc.pinY)
	return dx*dx + dy*dy
}

// hitTest returns the TRAVEL-TARGET at (mx, my). Only story-relevant pins
// are valid travel targets — unlocked-but-not-currently-relevant pins fall
// through and behave like locked pins (they open the info popup instead).
// This keeps the player from accidentally skipping ahead to a previously
// visited city when a new story target is glowing.
//
// When two hit rects overlap (Europe is crowded) the closest pin center
// wins so clicks on Rome don't bleed into Paris.
func (tm *travelMap) hitTest(mx, my int32) *travelLocation {
	pt := sdl.Point{X: mx, Y: my}
	var best *travelLocation
	var bestDist int64
	for i := range tm.locations {
		loc := &tm.locations[i]
		if !loc.unlocked || !tm.isRelevant(loc) {
			continue
		}
		hit := tm.pinHitRect(loc)
		if !pt.InRect(&hit) {
			continue
		}
		d := tm.distanceSqFromPin(loc, mx, my)
		if best == nil || d < bestDist {
			best = loc
			bestDist = d
		}
	}
	return best
}

func (tm *travelMap) hitTestAny(mx, my int32) *travelLocation {
	pt := sdl.Point{X: mx, Y: my}
	var best *travelLocation
	var bestDist int64
	for i := range tm.locations {
		loc := &tm.locations[i]
		hit := tm.pinHitRect(loc)
		if !pt.InRect(&hit) {
			continue
		}
		d := tm.distanceSqFromPin(loc, mx, my)
		if best == nil || d < bestDist {
			best = loc
			bestDist = d
		}
	}
	return best
}

func (tm *travelMap) setUnlocked(scene string, unlocked bool) {
	for i := range tm.locations {
		if tm.locations[i].scene == scene {
			tm.locations[i].unlocked = unlocked
		}
	}
}

// Drawing helpers for travel map surface generation
func fillEllipseTM(surface *sdl.Surface, cx, cy, rx, ry int32, color uint32) {
	for dy := -ry; dy <= ry; dy++ {
		halfW := int32(float64(rx) * math.Sqrt(1.0-float64(dy*dy)/float64(ry*ry)))
		if halfW > 0 {
			surface.FillRect(&sdl.Rect{X: cx - halfW, Y: cy + dy, W: halfW * 2, H: 1}, color)
		}
	}
}

func drawMountainTM(surface *sdl.Surface, f *sdl.PixelFormat, x, y, size int32, mColor, sColor uint32) {
	for dy := int32(0); dy <= size; dy++ {
		t := float64(dy) / float64(size)
		w := int32(float64(size) * t)
		if w < 1 {
			w = 1
		}
		surface.FillRect(&sdl.Rect{X: x - w/2, Y: y + dy - size, W: w, H: 1}, mColor)
	}
	// Snow cap
	for dy := int32(0); dy <= size/3; dy++ {
		t := float64(dy) / float64(size)
		w := int32(float64(size) * t)
		if w < 1 {
			w = 1
		}
		surface.FillRect(&sdl.Rect{X: x - w/2, Y: y + dy - size, W: w, H: 1}, sColor)
	}
}

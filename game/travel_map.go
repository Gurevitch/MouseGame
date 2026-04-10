package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type travelLocation struct {
	name     string
	scene    string
	pinX     int32
	pinY     int32
	unlocked bool
	info     string // city facts shown when clicking locked cities
}

type travelMap struct {
	locations []travelLocation
	renderer  *sdl.Renderer
	bgTex     *sdl.Texture
}

func newTravelMap(renderer *sdl.Renderer) *travelMap {
	tm := &travelMap{
		renderer: renderer,
		locations: []travelLocation{
			{name: "Camp Chilly Wa Wa", scene: "camp_entrance", pinX: 310, pinY: 280, unlocked: true, info: "Camp Chilly Wa Wa - A summer camp in the mountains. Home base for PP and the kids."},
			{name: "Paris", scene: "paris_street", pinX: 646, pinY: 296, unlocked: false, info: "Paris, France. City of lights! Home of the Eiffel Tower (1889) and the Louvre museum with over 380,000 artworks."},
			{name: "Jerusalem", scene: "jerusalem_street", pinX: 782, pinY: 349, unlocked: false, info: "Jerusalem, Israel. One of the oldest cities in the world. Home of the Western Wall and ancient underground tunnels."},
			{name: "Tokyo", scene: "tokyo_street", pinX: 1164, pinY: 328, unlocked: false, info: "Tokyo, Japan. A city of ancient temples and modern towers. Famous for cherry blossoms, torii gates, and Senso-ji temple."},
			{name: "Rome", scene: "rome_street", pinX: 730, pinY: 330, unlocked: false, info: "Rome, Italy. The Eternal City! Home of the Colosseum (72 AD), where gladiators once fought before 50,000 spectators."},
			{name: "Rio de Janeiro", scene: "rio_street", pinX: 431, pinY: 504, unlocked: false, info: "Rio de Janeiro, Brazil. Famous for Christ the Redeemer statue, Copacabana beach, and the world's biggest Carnival."},
			{name: "Egypt", scene: "egypt_street", pinX: 755, pinY: 369, unlocked: false, info: "Egypt. Home of the Great Pyramids of Giza, the Sphinx, and the ancient pharaohs. The Nile River runs through it all."},
			{name: "India", scene: "india_street", pinX: 932, pinY: 399, unlocked: false, info: "India. Home of the Taj Mahal, a monument of love built in 1632. A land of ancient temples, spices, and vibrant culture."},
			{name: "Thailand", scene: "thailand_street", pinX: 1000, pinY: 397, unlocked: false, info: "Thailand. Land of golden temples, floating markets, and ancient Buddhist monasteries. Known as the Land of Smiles."},
			{name: "China", scene: "china_street", pinX: 1049, pinY: 344, unlocked: false, info: "China. Home of the Great Wall (over 13,000 miles long!), the Forbidden City, and thousands of years of civilization."},
			{name: "Australia", scene: "australia_street", pinX: 1139, pinY: 569, unlocked: false, info: "Australia. Home of the Sydney Opera House, the Great Barrier Reef, and unique wildlife like kangaroos and koalas."},
		},
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

	for _, loc := range tm.locations {
		px, py := loc.pinX, loc.pinY

		if loc.unlocked {
			// === UNLOCKED: Yellow pulsing glow + bigger pin ===

			// Large yellow pulsing glow
			glowPulse := 0.5 + 0.5*math.Sin(ticks*0.004)
			for r := int32(35); r > 8; r -= 2 {
				a := uint8(float64(55-r) * glowPulse)
				renderer.SetDrawColor(255, 220, 50, a)
				for dy := -r; dy <= r; dy++ {
					hw := int32(math.Sqrt(float64(r*r - dy*dy)))
					renderer.FillRect(&sdl.Rect{X: px - hw, Y: py + dy, W: hw * 2, H: 1})
				}
			}

			// Pin base (bigger red circle)
			renderer.SetDrawColor(200, 50, 40, 255)
			for dy := int32(-10); dy <= 10; dy++ {
				hw := int32(math.Sqrt(float64(100 - dy*dy)))
				renderer.FillRect(&sdl.Rect{X: px - hw, Y: py + dy, W: hw * 2, H: 1})
			}
			// Pin highlight
			renderer.SetDrawColor(240, 100, 80, 200)
			for dy := int32(-7); dy <= 3; dy++ {
				hw := int32(math.Sqrt(float64(49 - dy*dy)))
				renderer.FillRect(&sdl.Rect{X: px - hw + 1, Y: py + dy - 1, W: hw, H: 1})
			}
			// Pin point
			renderer.SetDrawColor(200, 50, 40, 255)
			renderer.FillRect(&sdl.Rect{X: px - 1, Y: py + 10, W: 3, H: 12})
			renderer.FillRect(&sdl.Rect{X: px, Y: py + 22, W: 1, H: 4})

			// Location name label
			nameW := font.TextWidth(loc.name, 2)
			labelX := px - nameW/2
			labelY := py - 35

			// Label background
			renderer.SetDrawColor(30, 25, 18, 220)
			renderer.FillRect(&sdl.Rect{X: labelX - 8, Y: labelY - 4, W: nameW + 16, H: 24})
			renderer.SetDrawColor(220, 190, 60, 220)
			renderer.FillRect(&sdl.Rect{X: labelX - 8, Y: labelY - 4, W: nameW + 16, H: 2})
			renderer.FillRect(&sdl.Rect{X: labelX - 8, Y: labelY + 18, W: nameW + 16, H: 2})

			// Hover effect
			hoverRect := sdl.Rect{X: px - 45, Y: py - 40, W: 90, H: 80}
			pt := sdl.Point{X: mx, Y: my}
			if pt.InRect(&hoverRect) {
				renderer.SetDrawColor(255, 255, 180, 40)
				renderer.FillRect(&sdl.Rect{X: px - 50, Y: py - 45, W: 100, H: 90})
				font.DrawText(renderer, loc.name, labelX, labelY, 2, sdl.Color{R: 255, G: 255, B: 200, A: 255})
			} else {
				font.DrawText(renderer, loc.name, labelX, labelY, 2, sdl.Color{R: 255, G: 240, B: 160, A: 255})
			}
		} else {
			// === LOCKED: Small gray dot + city name ===

			// Small gray dot
			renderer.SetDrawColor(100, 95, 80, 150)
			for dy := int32(-4); dy <= 4; dy++ {
				hw := int32(math.Sqrt(float64(16 - dy*dy)))
				renderer.FillRect(&sdl.Rect{X: px - hw, Y: py + dy, W: hw * 2, H: 1})
			}

			// City name in muted text
			nameW := font.TextWidth(loc.name, 1)
			labelX := px - nameW/2
			labelY := py + 10
			font.DrawText(renderer, loc.name, labelX, labelY, 1, sdl.Color{R: 140, G: 130, B: 110, A: 160})
		}
	}

	// Dotted travel routes between unlocked locations
	renderer.SetDrawColor(220, 190, 60, 50)
	for i := 0; i < len(tm.locations); i++ {
		if !tm.locations[i].unlocked {
			continue
		}
		// Draw route from Camp to each unlocked city
		if tm.locations[i].scene == "camp_entrance" {
			continue
		}
		x1, y1 := float64(tm.locations[0].pinX), float64(tm.locations[0].pinY)
		x2, y2 := float64(tm.locations[i].pinX), float64(tm.locations[i].pinY)
		dx := x2 - x1
		dy := y2 - y1
		dist := math.Sqrt(dx*dx + dy*dy)
		steps := int(dist / 10)
		for s := 0; s < steps; s++ {
			if s%2 != 0 {
				continue
			}
			t := float64(s) / float64(steps)
			px := int32(x1 + dx*t)
			py := int32(y1 + dy*t)
			renderer.FillRect(&sdl.Rect{X: px, Y: py, W: 3, H: 3})
		}
	}

	// Instructions at bottom
	font.DrawText(renderer, "Click a destination to travel  |  ESC to go back",
		engine.ScreenWidth/2-310, engine.ScreenHeight-35, 2,
		sdl.Color{R: 170, G: 160, B: 140, A: 200})
}

func (tm *travelMap) hitTest(mx, my int32) *travelLocation {
	pt := sdl.Point{X: mx, Y: my}
	for i := range tm.locations {
		if !tm.locations[i].unlocked {
			continue
		}
		loc := &tm.locations[i]
		hitRect := sdl.Rect{X: loc.pinX - 40, Y: loc.pinY - 35, W: 80, H: 70}
		if pt.InRect(&hitRect) {
			return loc
		}
	}
	return nil
}

func (tm *travelMap) hitTestAny(mx, my int32) *travelLocation {
	pt := sdl.Point{X: mx, Y: my}
	for i := range tm.locations {
		loc := &tm.locations[i]
		hitRect := sdl.Rect{X: loc.pinX - 40, Y: loc.pinY - 35, W: 80, H: 70}
		if pt.InRect(&hitRect) {
			return loc
		}
	}
	return nil
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

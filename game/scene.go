package game

import (
	"math"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type arrowDir int

const (
	arrowNone arrowDir = iota
	arrowRight
	arrowLeft
	arrowUp
)

type hotspot struct {
	bounds      sdl.Rect
	targetScene string
	name        string
	arrow       arrowDir
}

type particle struct {
	x, y    float64
	vx, vy  float64
	alpha   uint8
	size    int32
	baseY   float64
	twinkle bool
}

type glowEffect struct {
	x, y    int32
	w, h    int32
	r, g, b uint8
	alpha   uint8
	pulse   float64
	timer   float64
}

type scene struct {
	name      string
	bg        *background
	npcs      []*npc
	hotspots  []hotspot
	particles []particle
	glows     []glowEffect
	blockers  []sdl.Rect
	spawnX    float64
	spawnY    float64
	musicPath string
}

type sceneManager struct {
	scenes        map[string]*scene
	currentName   string
	transitioning bool
	fadeAlpha     float64
	fadeIn        bool
	nextScene     string
	transPlayer   *player
}

func newSceneManager(renderer *sdl.Renderer) *sceneManager {
	sm := &sceneManager{
		scenes:      make(map[string]*scene),
		currentName: "street",
	}

	// ===== Street (London Day) =====
	street := &scene{
		name:   "street",
		bg:     newPNGBackground(renderer, "assets/images/locations/london/background/street_V2.png"),
		npcs:   []*npc{newPaparMan(renderer), newGrumpyKid(renderer), newStreetTalkers(renderer)},
		spawnX: 460,
		spawnY: 460,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "pub",
				name:        "Go Right",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 450, H: engine.ScreenHeight},
		},
	}

	// Dust motes drifting in sunlight
	for i := 0; i < 8; i++ {
		street.particles = append(street.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.3) * 8,
			vy:    -rand.Float64()*1.5 - 0.3,
			alpha: uint8(rand.Intn(15) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}

	street.glows = []glowEffect{
		{x: 400, y: 0, w: 500, h: 400, r: 255, g: 245, b: 210, alpha: 8, pulse: 0.2},
		{x: 0, y: 460, w: engine.ScreenWidth, h: 20, r: 200, g: 190, b: 160, alpha: 10, pulse: 0.4},
	}
	sm.scenes["street"] = street

	// ===== Interior (Wooden Cabin) =====
	interior := &scene{
		name:      "interior",
		bg:        newPNGBackground(renderer, "assets/images/backgrounds/bg_interior.png"),
		npcs:      []*npc{newCryingKid(renderer), newProfessor(renderer)},
		spawnX:    600,
		spawnY:    230,
		musicPath: "assets/sounds/The Pink Panther's Passport to Peril OST #08 - Camp Chilly Wa-Wa (Day 2 & 3) [HQ].mp3",
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 520, Y: 50, W: 260, H: 300},
				targetScene: "pub",
				name:        "Go Back",
				arrow:       arrowUp,
			},
		},
	}

	// Dust motes drifting through sunbeams
	for i := 0; i < 20; i++ {
		interior.particles = append(interior.particles, particle{
			x:     450 + rand.Float64()*500,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 6,
			vy:    -rand.Float64()*2 - 0.5,
			alpha: uint8(rand.Intn(30) + 10),
			size:  int32(rand.Intn(2) + 1),
		})
	}

	interior.glows = []glowEffect{
		{x: 540, y: 100, w: 310, h: 500, r: 255, g: 240, b: 200, alpha: 12, pulse: 0.3},
		{x: 60, y: 150, w: 160, h: 200, r: 200, g: 220, b: 240, alpha: 10, pulse: 0.4},
		{x: 1180, y: 150, w: 160, h: 200, r: 200, g: 220, b: 240, alpha: 10, pulse: 0.4},
		{x: 470, y: 550, w: 460, h: 100, r: 255, g: 220, b: 170, alpha: 8, pulse: 0.5},
	}
	sm.scenes["interior"] = interior

	// ===== Pub =====
	pub := &scene{
		name:   "pub",
		bg:     newPNGBackground(renderer, "assets/images/locations/london/background/pub 8K.jpg"),
		npcs:   []*npc{},
		spawnX: 600,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "street",
				name:        "Go Left",
				arrow:       arrowLeft,
			},
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "interior",
				name:        "Go Right",
				arrow:       arrowRight,
			},
		},
	}

	for i := 0; i < 10; i++ {
		pub.particles = append(pub.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * float64(engine.ScreenHeight),
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.0 - 0.2,
			alpha: uint8(rand.Intn(15) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}

	pub.glows = []glowEffect{
		{x: 300, y: 100, w: 200, h: 300, r: 255, g: 200, b: 120, alpha: 10, pulse: 0.3},
		{x: 800, y: 100, w: 200, h: 300, r: 255, g: 200, b: 120, alpha: 10, pulse: 0.4},
	}
	sm.scenes["pub"] = pub

	return sm
}

func (sm *sceneManager) current() *scene {
	return sm.scenes[sm.currentName]
}

func (sm *sceneManager) transitionTo(sceneName string, plr *player) {
	if _, ok := sm.scenes[sceneName]; !ok {
		return
	}
	sm.transitioning = true
	sm.fadeAlpha = 0
	sm.fadeIn = false
	sm.nextScene = sceneName
	sm.transPlayer = plr
}

func (sm *sceneManager) update(dt float64) {
	if !sm.transitioning {
		return
	}

	if !sm.fadeIn {
		sm.fadeAlpha += dt * 400
		if sm.fadeAlpha >= 255 {
			sm.fadeAlpha = 255
			sm.currentName = sm.nextScene
			sm.fadeIn = true
			if sm.transPlayer != nil {
				s := sm.scenes[sm.currentName]
				sm.transPlayer.x = s.spawnX
				sm.transPlayer.y = s.spawnY
				sm.transPlayer.moving = false
				sm.transPlayer.state = stateIdle
			}
		}
	} else {
		sm.fadeAlpha -= dt * 300
		if sm.fadeAlpha <= 0 {
			sm.fadeAlpha = 0
			sm.transitioning = false
		}
	}
}

func (sm *sceneManager) drawTransition(renderer *sdl.Renderer) {
	if !sm.transitioning {
		return
	}
	renderer.SetDrawColor(0, 0, 0, uint8(sm.fadeAlpha))
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})
}

func (s *scene) checkNPCClick(x, y int32) *npc {
	for _, n := range s.npcs {
		if n.silent {
			continue
		}
		if n.containsPoint(x, y) {
			if n.groupID != "" {
				return s.rightmostInGroup(n.groupID)
			}
			return n
		}
	}
	return nil
}

func (s *scene) rightmostInGroup(groupID string) *npc {
	var best *npc
	for _, n := range s.npcs {
		if n.groupID == groupID {
			if best == nil || n.bounds.X+n.bounds.W > best.bounds.X+best.bounds.W {
				best = n
			}
		}
	}
	return best
}

func (s *scene) checkHotspotClick(x, y int32) *hotspot {
	pt := sdl.Point{X: x, Y: y}
	for i := range s.hotspots {
		if pt.InRect(&s.hotspots[i].bounds) {
			return &s.hotspots[i]
		}
	}
	return nil
}

func (s *scene) drawBackground(renderer *sdl.Renderer, playerX float64) {
	s.bg.draw(renderer, playerX)
}

func (s *scene) drawHotspots(renderer *sdl.Renderer, hoverName string, mx, my int32) {
	pulse := 0.6 + 0.4*math.Sin(float64(sdl.GetTicks())*0.003)
	for _, hs := range s.hotspots {
		if hs.arrow == arrowNone {
			continue
		}
		hovered := hs.name == hoverName && hoverName != ""
		if hovered {
			alpha := uint8(float64(220) * pulse)
			drawArrow(renderer, mx, my, 40, 40, hs.arrow, alpha)
		} else {
			alpha := uint8(float64(60) * pulse)
			cx := hs.bounds.X + hs.bounds.W/2
			cy := hs.bounds.Y + hs.bounds.H/2
			drawArrow(renderer, cx, cy, 30, 30, hs.arrow, alpha)
		}
	}
}

func drawArrow(renderer *sdl.Renderer, cx, cy, w, h int32, dir arrowDir, alpha uint8) {
	renderer.SetDrawColor(255, 220, 100, alpha)
	switch dir {
	case arrowRight:
		// Right-pointing triangle: tip at right, base on left
		tipX := cx + w/2
		tipY := cy
		baseTop := cy - h/2
		baseBot := cy + h/2
		baseX := cx - w/2
		for y := baseTop; y <= baseBot; y++ {
			t := float64(y-baseTop) / float64(baseBot-baseTop)
			var x0, x1 int32
			if t <= 0.5 {
				x0 = baseX
				x1 = baseX + int32(float64(tipX-baseX)*t*2)
			} else {
				x0 = baseX
				x1 = baseX + int32(float64(tipX-baseX)*(1.0-t)*2)
			}
			if x1 > x0 {
				renderer.DrawLine(x0, y, x1, y)
			}
		}
		// Outline
		renderer.SetDrawColor(255, 240, 180, alpha)
		renderer.DrawLine(baseX, baseTop, tipX, tipY)
		renderer.DrawLine(tipX, tipY, baseX, baseBot)
		renderer.DrawLine(baseX, baseBot, baseX, baseTop)
	case arrowLeft:
		tipX := cx - w/2
		tipY := cy
		baseTop := cy - h/2
		baseBot := cy + h/2
		baseX := cx + w/2
		for y := baseTop; y <= baseBot; y++ {
			t := float64(y-baseTop) / float64(baseBot-baseTop)
			var x0, x1 int32
			if t <= 0.5 {
				x1 = baseX
				x0 = baseX - int32(float64(baseX-tipX)*t*2)
			} else {
				x1 = baseX
				x0 = baseX - int32(float64(baseX-tipX)*(1.0-t)*2)
			}
			if x1 > x0 {
				renderer.DrawLine(x0, y, x1, y)
			}
		}
		renderer.SetDrawColor(255, 240, 180, alpha)
		renderer.DrawLine(baseX, baseTop, tipX, tipY)
		renderer.DrawLine(tipX, tipY, baseX, baseBot)
		renderer.DrawLine(baseX, baseBot, baseX, baseTop)
	case arrowUp:
		// Up-pointing triangle: tip at top, base on bottom
		tipX := cx
		tipY := cy - h/2
		baseLeft := cx - w/2
		baseRight := cx + w/2
		baseY := cy + h/2
		for x := baseLeft; x <= baseRight; x++ {
			t := float64(x-baseLeft) / float64(baseRight-baseLeft)
			var y0, y1 int32
			if t <= 0.5 {
				y1 = baseY
				y0 = baseY - int32(float64(baseY-tipY)*t*2)
			} else {
				y1 = baseY
				y0 = baseY - int32(float64(baseY-tipY)*(1.0-t)*2)
			}
			if y1 > y0 {
				renderer.DrawLine(x, y0, x, y1)
			}
		}
		// Outline
		renderer.SetDrawColor(255, 240, 180, alpha)
		renderer.DrawLine(baseLeft, baseY, tipX, tipY)
		renderer.DrawLine(tipX, tipY, baseRight, baseY)
		renderer.DrawLine(baseRight, baseY, baseLeft, baseY)
	}
}

func (s *scene) drawNPCs(renderer *sdl.Renderer) {
	for _, n := range s.npcs {
		n.draw(renderer)
	}
}

func (s *scene) updateAmbient(dt float64) {
	for i := range s.particles {
		p := &s.particles[i]

		if p.twinkle {
			continue
		}

		p.x += p.vx * dt
		p.y += p.vy * dt

		if p.vy < 0 && p.y < -10 {
			p.y = float64(engine.ScreenHeight) + 10
			p.x = rand.Float64() * float64(engine.ScreenWidth)
		}
		if p.vy > 0 && p.y > float64(engine.ScreenHeight)+10 {
			p.y = p.baseY
			p.x += (rand.Float64() - 0.5) * 60
		}
		if p.vx > 0 && p.x > float64(engine.ScreenWidth)+float64(p.size) {
			p.x = -float64(p.size)
		}
		if p.vx < 0 && p.x < -float64(p.size) {
			p.x = float64(engine.ScreenWidth) + float64(p.size)
		}
	}

	for i := range s.glows {
		s.glows[i].timer += dt
	}

	for _, n := range s.npcs {
		n.update(dt)
	}
}

func (s *scene) drawAmbient(renderer *sdl.Renderer) {
	// Glow effects
	for _, g := range s.glows {
		base := 0.7 + 0.3*math.Sin(g.timer*g.pulse)
		if g.pulse > 3.0 {
			// High-pulse glows get random flicker jitter
			jitter := (rand.Float64() - 0.5) * 0.15
			base += jitter
			if base < 0.4 {
				base = 0.4
			}
		}
		a := float64(g.alpha) * base
		if a > 255 {
			a = 255
		}
		renderer.SetDrawColor(g.r, g.g, g.b, uint8(a))
		renderer.FillRect(&sdl.Rect{X: g.x, Y: g.y, W: g.w, H: g.h})
	}

	// Particles
	for i := range s.particles {
		p := &s.particles[i]

		if p.twinkle {
			// Twinkling star: alpha fades in and out using a unique phase per star
			phase := p.x*0.1 + p.y*0.07
			a := float64(p.alpha) * (0.3 + 0.7*math.Abs(math.Sin(phase+float64(sdl.GetTicks())*0.002)))
			renderer.SetDrawColor(255, 255, 240, uint8(a))
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
			continue
		}

		renderer.SetDrawColor(255, 255, 255, p.alpha)
		if p.size > 5 {
			// Wide fog particle
			renderer.SetDrawColor(200, 205, 215, p.alpha)
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size / 3})
		} else {
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
		}
	}
}

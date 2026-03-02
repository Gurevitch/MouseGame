package game

import (
	"math"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type hotspot struct {
	bounds      sdl.Rect
	targetScene string
	name        string
	r, g, b     uint8
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
	spawnX    float64
	spawnY    float64
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

	// ===== Street (London Night) =====
	street := &scene{
		name:   "street",
		bg:     newLondonBackground(renderer),
		npcs:   []*npc{newPaparMan(renderer)},
		spawnX: 200,
		spawnY: float64(engine.ScreenHeight) - playerDstH - 160,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1065, Y: 395, W: 55, H: 45},
				targetScene: "interior",
				name:        "Enter Building",
				r:           180, g: 140, b: 60,
			},
		},
	}

	// Twinkling star overlay particles (subtle flicker on top of static stars)
	for i := 0; i < 20; i++ {
		street.particles = append(street.particles, particle{
			x:       rand.Float64() * float64(engine.ScreenWidth),
			y:       rand.Float64() * 350,
			alpha:   uint8(rand.Intn(100) + 50),
			size:    int32(rand.Intn(2) + 1),
			twinkle: true,
		})
	}

	// Fog near ground level
	for i := 0; i < 25; i++ {
		street.particles = append(street.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     420 + rand.Float64()*60,
			baseY: 420 + rand.Float64()*60,
			vx:    rand.Float64()*20 + 5,
			vy:    0,
			alpha: uint8(rand.Intn(20) + 8),
			size:  int32(rand.Intn(60) + 30),
		})
	}

	// Lamp glow (warm, at the lamp head position)
	street.glows = []glowEffect{
		{x: 90, y: 200, w: 100, h: 120, r: 255, g: 210, b: 130, alpha: 35, pulse: 2.5},
		// Moon glow
		{x: 940, y: 50, w: 120, h: 120, r: 180, g: 190, b: 220, alpha: 15, pulse: 0.5},
		// Horizon warm glow
		{x: 0, y: 420, w: engine.ScreenWidth, h: 40, r: 50, g: 45, b: 35, alpha: 20, pulse: 1.0},
	}
	sm.scenes["street"] = street

	// ===== Interior (Wooden Cabin) =====
	interior := &scene{
		name:   "interior",
		bg:     newPNGBackground(renderer, "assets/images/bg_interior.png"),
		npcs:   []*npc{newCryingKid(renderer), newProfessor(renderer)},
		spawnX: 500,
		spawnY: float64(engine.ScreenHeight) - playerDstH - 120,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 510, Y: 120, W: 180, H: 350},
				targetScene: "street",
				name:        "Exit to Street",
				r:           180, g: 160, b: 100,
			},
		},
	}

	// Dust motes drifting through sunbeams
	for i := 0; i < 20; i++ {
		interior.particles = append(interior.particles, particle{
			x:     400 + rand.Float64()*400,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 6,
			vy:    -rand.Float64()*2 - 0.5,
			alpha: uint8(rand.Intn(30) + 10),
			size:  int32(rand.Intn(2) + 1),
		})
	}

	interior.glows = []glowEffect{
		// Sunlight from the central doorway
		{x: 460, y: 100, w: 280, h: 500, r: 255, g: 240, b: 200, alpha: 12, pulse: 0.3},
		// Left window light
		{x: 50, y: 150, w: 140, h: 200, r: 200, g: 220, b: 240, alpha: 10, pulse: 0.4},
		// Right window light
		{x: 1010, y: 150, w: 140, h: 200, r: 200, g: 220, b: 240, alpha: 10, pulse: 0.4},
		// Warm floor glow from doorway light
		{x: 400, y: 550, w: 400, h: 100, r: 255, g: 220, b: 170, alpha: 8, pulse: 0.5},
	}
	sm.scenes["interior"] = interior

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
		if n.containsPoint(x, y) {
			return n
		}
	}
	return nil
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

func (s *scene) drawBackground(renderer *sdl.Renderer) {
	s.bg.draw(renderer)
}

func (s *scene) drawHotspots(renderer *sdl.Renderer) {
	for _, hs := range s.hotspots {
		renderer.SetDrawColor(hs.r, hs.g, hs.b, 40)
		renderer.FillRect(&hs.bounds)
		renderer.SetDrawColor(hs.r, hs.g, hs.b, 80)
		renderer.DrawRect(&hs.bounds)
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
			// Stationary twinkling star: alpha oscillates
			continue
		}

		p.x += p.vx * dt
		p.y += p.vy * dt

		if p.vy < 0 && p.y < -10 {
			p.y = float64(engine.ScreenHeight) + 10
			p.x = rand.Float64() * float64(engine.ScreenWidth)
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
		a := float64(g.alpha) * (0.7 + 0.3*math.Sin(g.timer*g.pulse))
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
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size/3})
		} else {
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
		}
	}
}

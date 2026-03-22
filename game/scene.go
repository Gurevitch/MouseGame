package game

import (
	"math"
	"math/rand"
	"sort"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type arrowDir int

const (
	arrowNone arrowDir = iota
	arrowRight
	arrowLeft
	arrowUp
	arrowDown
)

type hotspot struct {
	bounds      sdl.Rect
	targetScene string
	name        string
	arrow       arrowDir
	onInteract  func() bool
}

type particle struct {
	x, y    float64
	vx, vy  float64
	alpha   uint8
	size    int32
	baseY   float64
	homeX   float64
	twinkle bool
	r, g, b uint8
	fire    bool
	bird    bool
}

type glowEffect struct {
	x, y    int32
	w, h    int32
	r, g, b uint8
	alpha   uint8
	pulse   float64
	timer   float64
}

type floorItem struct {
	tex      *sdl.Texture
	srcW     int32
	srcH     int32
	bounds   sdl.Rect
	name     string
	visible  bool
	onPickup func()
}

type scene struct {
	name       string
	bg         *background
	npcs       []*npc
	hotspots   []hotspot
	floorItems []*floorItem
	particles  []particle
	glows      []glowEffect
	blockers   []sdl.Rect
	spawnX     float64
	spawnY     float64
	musicPath  string
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
		currentName: "camp_entrance",
	}

	// ===== Street (London Day) =====
	street := &scene{
		name:   "street",
		bg:     newPNGBackground(renderer, "assets/images/locations/london/background/street_V2.png"),
		npcs:   []*npc{newPaperMan(renderer), newGrumpyKid(renderer), newStreetTalkers(renderer)},
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
		bg:     newPNGBackground(renderer, "assets/images/locations/london/background/pub_no_tables.png"),
		npcs:   []*npc{newBarmaid(renderer), newButler(renderer), newPoliceman(renderer)},
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

	// ===== Camp Chilly Wa Wa: Entrance =====
	campEntrance := &scene{
		name:   "camp_entrance",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_entrance.png"),
		npcs:   []*npc{newDirectorHiggins(renderer)},
		spawnX: 300,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "camp_grounds",
				name:        "Enter Camp",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 120, H: engine.ScreenHeight},
		},
	}
	for i := 0; i < 6; i++ {
		campEntrance.particles = append(campEntrance.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 350,
			vx:    (rand.Float64() - 0.3) * 5,
			vy:    -rand.Float64()*1.0 - 0.2,
			alpha: uint8(rand.Intn(12) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	for i := 0; i < 3; i++ {
		dir := 15 + rand.Float64()*20
		if rand.Float64() < 0.5 {
			dir = -dir
		}
		campEntrance.particles = append(campEntrance.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     30 + rand.Float64()*60,
			vx:    dir,
			baseY: 30 + rand.Float64()*60,
			alpha: uint8(rand.Intn(30) + 50),
			size:  3,
			bird:  true,
		})
	}
	campEntrance.glows = []glowEffect{
		{x: 300, y: 0, w: 600, h: 350, r: 255, g: 245, b: 210, alpha: 10, pulse: 0.25},
	}
	sm.scenes["camp_entrance"] = campEntrance

	// ===== Camp Chilly Wa Wa: Grounds =====
	campGrounds := &scene{
		name:   "camp_grounds",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_grounds.png"),
		npcs:   []*npc{newTommy(renderer), newJake(renderer), newLily(renderer), newMarcus(renderer), newDanny(renderer)},
		spawnX: 100,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "camp_entrance",
				name:        "Camp Entrance",
				arrow:       arrowLeft,
			},
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "camp_messhall",
				name:        "Mess Hall",
				arrow:       arrowRight,
			},
			{
				bounds:      sdl.Rect{X: 560, Y: 50, W: 280, H: 200},
				targetScene: "camp_lake",
				name:        "To the Lake",
				arrow:       arrowUp,
			},
		},
	}
	for i := 0; i < 10; i++ {
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*0.8 - 0.1,
			alpha: uint8(rand.Intn(10) + 3),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	// Campfire particles
	fireColors := [][3]uint8{{255, 140, 20}, {255, 180, 40}, {255, 100, 10}, {255, 200, 60}, {240, 80, 10},
		{255, 160, 30}, {255, 120, 15}, {255, 190, 50}}
	for i := 0; i < 8; i++ {
		c := fireColors[i%len(fireColors)]
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     680 + (rand.Float64()-0.5)*20,
			y:     515 - rand.Float64()*30,
			vx:    (rand.Float64() - 0.5) * 12,
			vy:    -rand.Float64()*35 - 15,
			alpha: uint8(rand.Intn(50) + 30),
			size:  int32(rand.Intn(2) + 1),
			baseY: 518,
			homeX: 680,
			fire:  true,
			r:     c[0], g: c[1], b: c[2],
		})
	}
	// Birds in the sky
	for i := 0; i < 3; i++ {
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     40 + rand.Float64()*80,
			vx:    15 + rand.Float64()*20,
			baseY: 40 + rand.Float64()*80,
			alpha: uint8(rand.Intn(30) + 50),
			size:  3,
			bird:  true,
		})
	}
	campGrounds.glows = []glowEffect{
		{x: 200, y: 50, w: 400, h: 300, r: 255, g: 245, b: 200, alpha: 8, pulse: 0.2},
		{x: 500, y: 400, w: 300, h: 100, r: 255, g: 200, b: 120, alpha: 6, pulse: 0.35},
		{x: 650, y: 490, w: 60, h: 40, r: 255, g: 160, b: 40, alpha: 18, pulse: 4.0},
	}
	sm.scenes["camp_grounds"] = campGrounds

	// ===== Camp Chilly Wa Wa: Mess Hall =====
	campMessHall := &scene{
		name:   "camp_messhall",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_messhall.png"),
		npcs:   []*npc{newCookMarge(renderer)},
		spawnX: 640,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 300, Y: 650, W: 700, H: 100},
				targetScene: "camp_grounds",
				name:        "Back to Camp",
				arrow:       arrowDown,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: engine.ScreenWidth, H: 300},
		},
	}
	campMessHall.glows = []glowEffect{
		{x: 300, y: 100, w: 600, h: 400, r: 255, g: 230, b: 180, alpha: 8, pulse: 0.3},
	}
	sm.scenes["camp_messhall"] = campMessHall

	// ===== Camp Chilly Wa Wa: Lake =====
	campLake := &scene{
		name:   "camp_lake",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_lake.png"),
		npcs:   []*npc{},
		spawnX: 200,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "camp_grounds",
				name:        "Back to Camp",
				arrow:       arrowLeft,
			},
		},
	}
	for i := 0; i < 5; i++ {
		campLake.particles = append(campLake.particles, particle{
			x:     400 + rand.Float64()*500,
			y:     300 + rand.Float64()*200,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(8) + 3),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	for i := 0; i < 3; i++ {
		dir := 12 + rand.Float64()*18
		if rand.Float64() < 0.5 {
			dir = -dir
		}
		campLake.particles = append(campLake.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     20 + rand.Float64()*60,
			vx:    dir,
			baseY: 20 + rand.Float64()*60,
			alpha: uint8(rand.Intn(25) + 40),
			size:  3,
			bird:  true,
		})
	}
	campLake.glows = []glowEffect{
		{x: 400, y: 250, w: 500, h: 200, r: 255, g: 200, b: 120, alpha: 6, pulse: 0.2},
		{x: 0, y: 0, w: engine.ScreenWidth, h: 200, r: 180, g: 150, b: 200, alpha: 5, pulse: 0.15},
	}
	sm.scenes["camp_lake"] = campLake

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
				sm.transPlayer.allowOffscreen = false
				sm.transPlayer.facingLeft = false
				sm.transPlayer.dir = dirDown
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

func (s *scene) checkFloorItemClick(x, y int32) *floorItem {
	pt := sdl.Point{X: x, Y: y}
	for _, fi := range s.floorItems {
		if fi.visible && pt.InRect(&fi.bounds) {
			return fi
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

func (s *scene) drawBackground(renderer *sdl.Renderer, playerX float64) {
	s.bg.draw(renderer, playerX)
}

func (s *scene) drawHotspots(renderer *sdl.Renderer, hoverName string, _, _ int32) {
	pulse := 0.5 + 0.5*math.Sin(float64(sdl.GetTicks())*0.004)
	for _, hs := range s.hotspots {
		if hs.arrow == arrowNone {
			continue
		}
		hovered := hs.name == hoverName && hoverName != ""
		if hovered {
			continue
		}

		cx := hs.bounds.X + hs.bounds.W/2
		cy := hs.bounds.Y + hs.bounds.H/2

		switch hs.arrow {
		case arrowLeft:
			cx = hs.bounds.X + 20
		case arrowRight:
			cx = hs.bounds.X + hs.bounds.W - 20
		case arrowUp:
			cy = hs.bounds.Y + 20
		case arrowDown:
			cy = hs.bounds.Y + hs.bounds.H - 20
		}

		a := uint8(40 + float64(40)*pulse)
		for r := int32(4); r >= 1; r-- {
			renderer.SetDrawColor(255, 220, 100, a)
			renderer.FillRect(&sdl.Rect{X: cx - r, Y: cy - r, W: r * 2, H: r * 2})
			a = uint8(float64(a) * 0.6)
		}
	}
}

func (s *scene) drawActors(renderer *sdl.Renderer, plr *player) {
	type actorDraw struct {
		footY int32
		order int
		draw  func()
	}

	actors := make([]actorDraw, 0, len(s.npcs)+len(s.floorItems)+1)

	for i := range s.floorItems {
		fi := s.floorItems[i]
		if !fi.visible {
			continue
		}
		actors = append(actors, actorDraw{
			footY: fi.bounds.Y + fi.bounds.H,
			order: i,
			draw: func() {
				renderer.Copy(fi.tex, nil, &fi.bounds)
			},
		})
	}

	base := len(s.floorItems)
	for i := range s.npcs {
		n := s.npcs[i]
		actors = append(actors, actorDraw{
			footY: n.footY(),
			order: base + i,
			draw: func() {
				n.draw(renderer)
			},
		})
	}

	if plr != nil {
		actors = append(actors, actorDraw{
			footY: plr.footY(),
			order: base + len(s.npcs),
			draw: func() {
				plr.draw(renderer)
			},
		})
	}

	sort.SliceStable(actors, func(i, j int) bool {
		if actors[i].footY == actors[j].footY {
			return actors[i].order < actors[j].order
		}
		return actors[i].footY < actors[j].footY
	})

	for _, actor := range actors {
		actor.draw()
	}
}

func (s *scene) updateAmbient(dt float64) {
	for i := range s.particles {
		p := &s.particles[i]

		if p.twinkle {
			continue
		}

		if p.fire {
			p.x += p.vx * dt
			p.y += p.vy * dt
			fadeRate := 80 * dt
			if float64(p.alpha) > fadeRate {
				p.alpha -= uint8(fadeRate)
			} else {
				p.alpha = 0
			}
			if p.alpha < 5 || p.y < p.baseY-70 {
				p.x = p.homeX + (rand.Float64()-0.5)*24
				p.y = p.baseY + rand.Float64()*4
				p.alpha = uint8(rand.Intn(50) + 30)
				p.vx = (rand.Float64() - 0.5) * 12
				p.vy = -rand.Float64()*35 - 15
			}
			continue
		}

		if p.bird {
			p.x += p.vx * dt
			p.y = p.baseY + math.Sin(p.x*0.02)*8
			if p.vx > 0 && p.x > float64(engine.ScreenWidth)+20 {
				p.x = -20
			}
			if p.vx < 0 && p.x < -20 {
				p.x = float64(engine.ScreenWidth) + 20
			}
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
			phase := p.x*0.1 + p.y*0.07
			a := float64(p.alpha) * (0.3 + 0.7*math.Abs(math.Sin(phase+float64(sdl.GetTicks())*0.002)))
			renderer.SetDrawColor(255, 255, 240, uint8(a))
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
			continue
		}

		if p.fire {
			renderer.SetDrawColor(p.r, p.g, p.b, p.alpha)
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
			continue
		}

		if p.bird {
			renderer.SetDrawColor(30, 25, 20, p.alpha)
			px := int32(p.x)
			py := int32(p.y)
			renderer.FillRect(&sdl.Rect{X: px, Y: py, W: 3, H: 1})
			renderer.FillRect(&sdl.Rect{X: px - 1, Y: py - 1, W: 1, H: 1})
			renderer.FillRect(&sdl.Rect{X: px + 3, Y: py - 1, W: 1, H: 1})
			continue
		}

		renderer.SetDrawColor(255, 255, 255, p.alpha)
		if p.size > 5 {
			renderer.SetDrawColor(200, 205, 215, p.alpha)
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size / 3})
		} else {
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
		}
	}
}

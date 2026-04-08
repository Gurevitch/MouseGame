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
	smoke   bool
	insect  bool
	cloud   bool
	water   bool
	timer   float64
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

type walkSegment struct {
	x1, y1, x2, y2 float64
}

type scene struct {
	name         string
	bg           *background
	npcs         []*npc
	hotspots     []hotspot
	floorItems   []*floorItem
	particles    []particle
	glows        []glowEffect
	blockers     []sdl.Rect
	walkSegments []walkSegment
	spawnX       float64
	spawnY       float64
	musicPath    string
	minY         float64
	maxY         float64
}

// snapToPath finds the nearest point on any walk segment. If no segments, returns input.
func (s *scene) snapToPath(x, y float64) (float64, float64) {
	if len(s.walkSegments) == 0 {
		return x, y
	}
	bestX, bestY := x, y
	bestDist := math.MaxFloat64
	for _, seg := range s.walkSegments {
		px, py := nearestPointOnSegment(x, y, seg.x1, seg.y1, seg.x2, seg.y2)
		dx, dy := px-x, py-y
		d := dx*dx + dy*dy
		if d < bestDist {
			bestDist = d
			bestX, bestY = px, py
		}
	}
	return bestX, bestY
}

func nearestPointOnSegment(px, py, x1, y1, x2, y2 float64) (float64, float64) {
	dx, dy := x2-x1, y2-y1
	lenSq := dx*dx + dy*dy
	if lenSq == 0 {
		return x1, y1
	}
	t := ((px-x1)*dx + (py-y1)*dy) / lenSq
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return x1 + t*dx, y1 + t*dy
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

	// ===== Camp Chilly Wa Wa: Entrance =====
	campEntrance := &scene{
		name:   "camp_entrance",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_entrance.png"),
		npcs:   []*npc{newDirectorHiggins(renderer)},
		spawnX: 500,
		spawnY: 570,
		minY:   400.0,
		maxY:   620.0,
		hotspots: []hotspot{
			{
				bounds: sdl.Rect{X: 500, Y: 100, W: 300, H: 250},
				name:   "Enter Camp",
				arrow:  arrowUp,
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
	// Butterflies near flowers/bushes
	butterflyColors := [][3]uint8{{240, 200, 80}, {180, 120, 200}, {100, 180, 220}}
	for i := 0; i < 3; i++ {
		c := butterflyColors[i%len(butterflyColors)]
		campEntrance.particles = append(campEntrance.particles, particle{
			x:      300 + rand.Float64()*600,
			baseY:  200 + rand.Float64()*200,
			vx:     (rand.Float64() - 0.5) * 12,
			alpha:  uint8(rand.Intn(40) + 50),
			insect: true,
			r:      c[0], g: c[1], b: c[2],
			timer: rand.Float64() * 10,
		})
	}
	// Clouds
	for i := 0; i < 2; i++ {
		campEntrance.particles = append(campEntrance.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     30 + rand.Float64()*60,
			vx:    2 + rand.Float64()*3,
			alpha: uint8(rand.Intn(6) + 3),
			size:  int32(50 + rand.Intn(40)),
			cloud: true,
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
		spawnX: 650,
		spawnY: 650,
		walkSegments: []walkSegment{
			{650, 650, 650, 550},
			{120, 480, 380, 480},
			{380, 480, 620, 520},
			{620, 520, 900, 480},
			{900, 480, 1200, 500},
			{620, 520, 650, 650},
			{120, 480, 50, 500},
			{1200, 500, 1350, 520},
		},
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 80, H: 500},
				targetScene: "camp_lake",
				name:        "To the Lake",
				arrow:       arrowLeft,
			},
			{
				bounds:      sdl.Rect{X: 1320, Y: 200, W: 80, H: 500},
				targetScene: "camp_office",
				name:        "Higgins' Office",
				arrow:       arrowRight,
			},
			{
				bounds:      sdl.Rect{X: 70, Y: 340, W: 120, H: 100},
				targetScene: "tommy_room",
				name:        "Tommy's Cabin",
				arrow:       arrowUp,
			},
			{
				bounds:      sdl.Rect{X: 320, Y: 310, W: 120, H: 100},
				targetScene: "jake_room",
				name:        "Jake's Cabin",
				arrow:       arrowUp,
			},
			{
				bounds:      sdl.Rect{X: 560, Y: 290, W: 120, H: 100},
				targetScene: "lily_room",
				name:        "Lily's Cabin",
				arrow:       arrowUp,
			},
			{
				bounds:      sdl.Rect{X: 860, Y: 310, W: 120, H: 100},
				targetScene: "marcus_room",
				name:        "Marcus's Cabin",
				arrow:       arrowUp,
			},
			{
				bounds:      sdl.Rect{X: 1160, Y: 340, W: 120, H: 100},
				targetScene: "danny_room",
				name:        "Danny's Cabin",
				arrow:       arrowUp,
			},
		},
	}
	for i := 0; i < 10; i++ {
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 350,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*0.8 - 0.1,
			alpha: uint8(rand.Intn(10) + 3),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	fireColors := [][3]uint8{{255, 140, 20}, {255, 180, 40}, {255, 100, 10}, {255, 200, 60}, {240, 80, 10},
		{255, 160, 30}, {255, 120, 15}, {255, 190, 50}}
	for i := 0; i < 8; i++ {
		c := fireColors[i%len(fireColors)]
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     620 + (rand.Float64()-0.5)*20,
			y:     520 - rand.Float64()*30,
			vx:    (rand.Float64() - 0.5) * 12,
			vy:    -rand.Float64()*35 - 15,
			alpha: uint8(rand.Intn(50) + 30),
			size:  int32(rand.Intn(2) + 1),
			baseY: 525,
			homeX: 620,
			fire:  true,
			r:     c[0], g: c[1], b: c[2],
		})
	}
	for i := 0; i < 3; i++ {
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     30 + rand.Float64()*60,
			vx:    15 + rand.Float64()*20,
			baseY: 30 + rand.Float64()*60,
			alpha: uint8(rand.Intn(30) + 50),
			size:  3,
			bird:  true,
		})
	}
	for i := 0; i < 5; i++ {
		campGrounds.particles = append(campGrounds.particles, particle{
			x:     620 + (rand.Float64()-0.5)*12,
			y:     495 - rand.Float64()*20,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*12 - 6,
			alpha: uint8(rand.Intn(15) + 8),
			size:  int32(rand.Intn(3) + 2),
			baseY: 495,
			homeX: 620,
			smoke: true,
			r:     140, g: 130, b: 120,
			timer: rand.Float64() * 10,
		})
	}
	for i := 0; i < 6; i++ {
		campGrounds.particles = append(campGrounds.particles, particle{
			x:       100 + rand.Float64()*1100,
			y:       350 + rand.Float64()*150,
			twinkle: true,
			alpha:   uint8(rand.Intn(30) + 20),
			size:    1,
			r:       255, g: 255, b: 150,
		})
	}
	campGrounds.glows = []glowEffect{
		{x: 200, y: 0, w: 800, h: 300, r: 255, g: 245, b: 200, alpha: 8, pulse: 0.2},
		{x: 450, y: 400, w: 300, h: 100, r: 255, g: 200, b: 120, alpha: 6, pulse: 0.35},
		{x: 590, y: 500, w: 60, h: 40, r: 255, g: 160, b: 40, alpha: 18, pulse: 4.0},
	}
	sm.scenes["camp_grounds"] = campGrounds

	// ===== Camp Chilly Wa Wa: Higgins' Office =====
	campOffice := &scene{
		name:   "camp_office",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_office.png"),
		npcs:   []*npc{newOfficeHiggins(renderer)},
		spawnX: 150,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "camp_grounds",
				name:        "Back to Camp",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{
			{X: 500, Y: 0, W: 900, H: 360},
		},
	}
	campOffice.glows = []glowEffect{
		{x: 600, y: 200, w: 300, h: 300, r: 255, g: 230, b: 170, alpha: 10, pulse: 0.3},
		{x: 0, y: 100, w: 300, h: 400, r: 255, g: 245, b: 210, alpha: 8, pulse: 0.2},
	}
	sm.scenes["camp_office"] = campOffice

	// ===== Camp Chilly Wa Wa: Night (campfire scene) =====
	campNight := &scene{
		name:   "camp_night",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/camp_night.png"),
		npcs:   []*npc{},
		spawnX: 650,
		spawnY: 550,
	}
	for i := 0; i < 16; i++ {
		campNight.particles = append(campNight.particles, particle{
			x:       80 + rand.Float64()*1200,
			y:       200 + rand.Float64()*300,
			twinkle: true,
			alpha:   uint8(rand.Intn(40) + 30),
			size:    1,
			r:       255, g: 255, b: 150,
		})
	}
	for i := 0; i < 6; i++ {
		campNight.particles = append(campNight.particles, particle{
			x:     620 + (rand.Float64()-0.5)*30,
			y:     500 - rand.Float64()*20,
			vx:    (rand.Float64() - 0.5) * 8,
			vy:    -rand.Float64()*25 - 10,
			alpha: uint8(rand.Intn(40) + 20),
			size:  int32(rand.Intn(2) + 1),
			baseY: 505,
			homeX: 620,
			fire:  true,
			r:     255, g: 160, b: 40,
		})
	}
	campNight.glows = []glowEffect{
		{x: 480, y: 400, w: 280, h: 180, r: 255, g: 160, b: 40, alpha: 15, pulse: 3.5},
		{x: 0, y: 0, w: engine.ScreenWidth, h: engine.ScreenHeight, r: 20, g: 15, b: 40, alpha: 12, pulse: 0.1},
		{x: 820, y: 280, w: 180, h: 120, r: 255, g: 200, b: 100, alpha: 8, pulse: 1.5},
	}
	sm.scenes["camp_night"] = campNight

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
		blockers: []sdl.Rect{
			{X: 500, Y: 0, W: 900, H: engine.ScreenHeight},
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
	// Water shimmer ripples
	for i := 0; i < 12; i++ {
		baseX := 350 + float64(i)*60
		baseY := 380 + rand.Float64()*120
		campLake.particles = append(campLake.particles, particle{
			homeX: baseX,
			baseY: baseY,
			x:     baseX,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(20) + 15),
			water: true,
			timer: rand.Float64() * 10,
		})
	}
	// Dragonflies over water
	for i := 0; i < 2; i++ {
		campLake.particles = append(campLake.particles, particle{
			x:      500 + rand.Float64()*400,
			baseY:  320 + rand.Float64()*60,
			vx:     6 + rand.Float64()*8,
			alpha:  uint8(rand.Intn(40) + 40),
			insect: true,
			r:      80, g: 160, b: 200,
			timer: rand.Float64() * 10,
		})
	}
	campLake.glows = []glowEffect{
		{x: 400, y: 250, w: 500, h: 200, r: 255, g: 200, b: 120, alpha: 6, pulse: 0.2},
		{x: 0, y: 0, w: engine.ScreenWidth, h: 200, r: 180, g: 150, b: 200, alpha: 5, pulse: 0.15},
		// Water surface reflection shimmer
		{x: 350, y: 370, w: 600, h: 20, r: 200, g: 230, b: 255, alpha: 5, pulse: 1.5},
	}
	sm.scenes["camp_lake"] = campLake

	// ===== Paris: Street (Eiffel Tower) =====
	parisStreet := &scene{
		name:   "paris_street",
		bg:     newPNGBackground(renderer, "assets/images/locations/paris/background/paris_street.png"),
		npcs:   []*npc{newFrenchGuide(renderer)},
		spawnX: 200,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "paris_louvre",
				name:        "To the Louvre",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 180, H: 500},
		},
	}
	// Pigeons
	for i := 0; i < 4; i++ {
		dir := 10 + rand.Float64()*15
		if rand.Float64() < 0.5 {
			dir = -dir
		}
		parisStreet.particles = append(parisStreet.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     50 + rand.Float64()*80,
			vx:    dir,
			baseY: 50 + rand.Float64()*80,
			alpha: uint8(rand.Intn(25) + 40),
			size:  3,
			bird:  true,
		})
	}
	// Cafe steam
	for i := 0; i < 3; i++ {
		baseX := 150 + float64(i)*80
		parisStreet.particles = append(parisStreet.particles, particle{
			x:     baseX + (rand.Float64()-0.5)*8,
			y:     350 - rand.Float64()*15,
			vx:    (rand.Float64() - 0.5) * 2,
			vy:    -rand.Float64()*10 - 5,
			alpha: uint8(rand.Intn(12) + 6),
			size:  int32(rand.Intn(2) + 2),
			baseY: 350,
			homeX: baseX,
			smoke: true,
			r:     230, g: 225, b: 220,
			timer: rand.Float64() * 10,
		})
	}
	// Dust motes
	for i := 0; i < 6; i++ {
		parisStreet.particles = append(parisStreet.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.3) * 5,
			vy:    -rand.Float64()*1.0 - 0.2,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	parisStreet.glows = []glowEffect{
		{x: 300, y: 0, w: 600, h: 400, r: 255, g: 245, b: 210, alpha: 10, pulse: 0.25},
		{x: 50, y: 300, w: 200, h: 150, r: 255, g: 220, b: 160, alpha: 8, pulse: 0.3},
	}
	sm.scenes["paris_street"] = parisStreet

	// ===== Paris: Louvre Interior =====
	parisLouvre := &scene{
		name:   "paris_louvre",
		bg:     newPNGBackground(renderer, "assets/images/locations/paris/background/paris_museum.png"),
		npcs:   []*npc{newMuseumCurator(renderer)},
		spawnX: 200,
		spawnY: 400,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "paris_street",
				name:        "Back to Street",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{
			{X: 400, Y: 0, W: 500, H: 340},
		},
	}
	// Dust motes in sunbeams
	for i := 0; i < 15; i++ {
		parisLouvre.particles = append(parisLouvre.particles, particle{
			x:     400 + rand.Float64()*500,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.5 - 0.3,
			alpha: uint8(rand.Intn(20) + 8),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	parisLouvre.glows = []glowEffect{
		{x: 400, y: 50, w: 400, h: 500, r: 255, g: 240, b: 200, alpha: 10, pulse: 0.2},
		{x: 200, y: 100, w: 150, h: 300, r: 255, g: 230, b: 180, alpha: 8, pulse: 0.4},
		{x: 900, y: 100, w: 150, h: 300, r: 255, g: 230, b: 180, alpha: 8, pulse: 0.4},
	}
	sm.scenes["paris_louvre"] = parisLouvre

	// ===== Cabin Interiors =====

	// --- Tommy's Room ---
	tommyRoom := &scene{
		name:   "tommy_room",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/tommy_room.png"),
		npcs:   []*npc{},
		spawnX: 1100,
		spawnY: 500,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1250, Y: 300, W: 130, H: 400},
				targetScene: "camp_grounds",
				name:        "Exit Cabin",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 500, H: 450},
			{X: 500, Y: 0, W: 500, H: 350},
		},
	}
	for i := 0; i < 8; i++ {
		tommyRoom.particles = append(tommyRoom.particles, particle{
			x:     300 + rand.Float64()*700,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	tommyRoom.glows = []glowEffect{
		{x: 500, y: 150, w: 400, h: 350, r: 255, g: 240, b: 190, alpha: 8, pulse: 0.2},
	}
	sm.scenes["tommy_room"] = tommyRoom

	// --- Jake's Room ---
	jakeRoom := &scene{
		name:   "jake_room",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/jake_room.png"),
		npcs:   []*npc{},
		spawnX: 1100,
		spawnY: 500,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1250, Y: 300, W: 130, H: 400},
				targetScene: "camp_grounds",
				name:        "Exit Cabin",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 450, H: 420},
			{X: 500, Y: 0, W: 500, H: 350},
		},
	}
	for i := 0; i < 6; i++ {
		jakeRoom.particles = append(jakeRoom.particles, particle{
			x:     300 + rand.Float64()*700,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(8) + 3),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	jakeRoom.glows = []glowEffect{
		{x: 550, y: 150, w: 350, h: 300, r: 255, g: 240, b: 200, alpha: 7, pulse: 0.2},
	}
	sm.scenes["jake_room"] = jakeRoom

	// --- Lily's Room ---
	lilyRoom := &scene{
		name:   "lily_room",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/lily_room.png"),
		npcs:   []*npc{},
		spawnX: 700,
		spawnY: 550,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 600, Y: 650, W: 250, H: 150},
				targetScene: "camp_grounds",
				name:        "Exit Cabin",
				arrow:       arrowDown,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 300, H: 500},
			{X: 400, Y: 0, W: 600, H: 350},
			{X: 1000, Y: 0, W: 400, H: 500},
		},
	}
	for i := 0; i < 6; i++ {
		lilyRoom.particles = append(lilyRoom.particles, particle{
			x:     400 + rand.Float64()*500,
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(10) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	butterflyColorsLily := [][3]uint8{{240, 180, 200}, {200, 160, 220}, {180, 220, 200}}
	for i := 0; i < 3; i++ {
		c := butterflyColorsLily[i%len(butterflyColorsLily)]
		lilyRoom.particles = append(lilyRoom.particles, particle{
			x:      500 + rand.Float64()*400,
			baseY:  200 + rand.Float64()*150,
			vx:     (rand.Float64() - 0.5) * 10,
			alpha:  uint8(rand.Intn(35) + 40),
			insect: true,
			r:      c[0], g: c[1], b: c[2],
			timer: rand.Float64() * 10,
		})
	}
	lilyRoom.glows = []glowEffect{
		{x: 400, y: 100, w: 500, h: 400, r: 255, g: 245, b: 220, alpha: 8, pulse: 0.2},
	}
	sm.scenes["lily_room"] = lilyRoom

	// --- Marcus's Room ---
	marcusRoom := &scene{
		name:   "marcus_room",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/marcus_room_day.png"),
		npcs:   []*npc{newMarcus(renderer)},
		spawnX: 700,
		spawnY: 550,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 600, Y: 650, W: 250, H: 150},
				targetScene: "camp_grounds",
				name:        "Exit Cabin",
				arrow:       arrowDown,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 400, H: 450},
			{X: 400, Y: 0, W: 500, H: 300},
			{X: 1000, Y: 0, W: 400, H: 450},
		},
	}
	for i := 0; i < 10; i++ {
		marcusRoom.particles = append(marcusRoom.particles, particle{
			x:     300 + rand.Float64()*700,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	marcusRoom.glows = []glowEffect{
		{x: 500, y: 100, w: 400, h: 400, r: 255, g: 240, b: 190, alpha: 9, pulse: 0.25},
	}
	sm.scenes["marcus_room"] = marcusRoom

	// --- Danny's Room ---
	dannyRoom := &scene{
		name:   "danny_room",
		bg:     newPNGBackground(renderer, "assets/images/locations/camp/background/danny_room.png"),
		npcs:   []*npc{},
		spawnX: 1100,
		spawnY: 500,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1250, Y: 300, W: 130, H: 400},
				targetScene: "camp_grounds",
				name:        "Exit Cabin",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 400, H: 450},
			{X: 400, Y: 0, W: 600, H: 300},
		},
	}
	for i := 0; i < 8; i++ {
		dannyRoom.particles = append(dannyRoom.particles, particle{
			x:     300 + rand.Float64()*700,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(8) + 3),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	dannyRoom.glows = []glowEffect{
		{x: 400, y: 50, w: 300, h: 300, r: 255, g: 245, b: 210, alpha: 7, pulse: 0.2},
	}
	sm.scenes["danny_room"] = dannyRoom

	return sm
}

func (sm *sceneManager) current() *scene {
	return sm.scenes[sm.currentName]
}

func (s *scene) effectiveMinY() float64 {
	if s.minY > 0 {
		return s.minY
	}
	return playerMinY
}

func (s *scene) effectiveMaxY() float64 {
	if s.maxY > 0 {
		return s.maxY
	}
	return playerMaxY
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
				sm.transPlayer.sceneMinY = s.minY
				sm.transPlayer.sceneMaxY = s.maxY
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

func (s *scene) drawHotspots(renderer *sdl.Renderer, hoverName string, mx, my int32) {
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
			cx = hs.bounds.X + 14
		case arrowRight:
			cx = hs.bounds.X + hs.bounds.W - 14
		case arrowUp:
			cy = hs.bounds.Y + 14
		case arrowDown:
			cy = hs.bounds.Y + hs.bounds.H - 14
		}

		dx := float64(mx) - float64(cx)
		dy := float64(my) - float64(cy)
		dist := math.Sqrt(dx*dx + dy*dy)
		proximityFade := 1.0
		if dist > 200 {
			proximityFade = math.Max(0.15, 1.0-(dist-200)/400)
		}

		baseA := (30 + float64(35)*pulse) * proximityFade
		a := uint8(baseA)

		sz := int32(5)
		switch hs.arrow {
		case arrowLeft:
			for i := int32(0); i < sz; i++ {
				renderer.SetDrawColor(255, 220, 100, a)
				renderer.FillRect(&sdl.Rect{X: cx + i, Y: cy - (sz - i), W: 1, H: (sz-i)*2 + 1})
			}
		case arrowRight:
			for i := int32(0); i < sz; i++ {
				renderer.SetDrawColor(255, 220, 100, a)
				renderer.FillRect(&sdl.Rect{X: cx - i, Y: cy - (sz - i), W: 1, H: (sz-i)*2 + 1})
			}
		case arrowUp:
			for i := int32(0); i < sz; i++ {
				renderer.SetDrawColor(255, 220, 100, a)
				renderer.FillRect(&sdl.Rect{X: cx - (sz - i), Y: cy + i, W: (sz-i)*2 + 1, H: 1})
			}
		case arrowDown:
			for i := int32(0); i < sz; i++ {
				renderer.SetDrawColor(255, 220, 100, a)
				renderer.FillRect(&sdl.Rect{X: cx - (sz - i), Y: cy - i, W: (sz-i)*2 + 1, H: 1})
			}
		}
	}
}

func (s *scene) drawActorsNoPlayer(renderer *sdl.Renderer) {
	s.drawActors(renderer, nil)
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
		if n.silent {
			continue
		}
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

		if p.smoke {
			p.timer += dt
			p.x += p.vx*dt + math.Sin(p.timer*2)*3*dt
			p.y += p.vy * dt
			fadeRate := 25 * dt
			if float64(p.alpha) > fadeRate {
				p.alpha -= uint8(fadeRate)
			} else {
				p.alpha = 0
			}
			if p.alpha < 3 || p.y < p.baseY-120 {
				p.x = p.homeX + (rand.Float64()-0.5)*16
				p.y = p.baseY
				p.alpha = uint8(rand.Intn(25) + 10)
				p.timer = rand.Float64() * 10
			}
			continue
		}

		if p.insect {
			p.timer += dt
			p.x += math.Sin(p.timer*3.5)*40*dt + p.vx*dt
			p.y = p.baseY + math.Sin(p.timer*2.7)*15
			if p.x > float64(engine.ScreenWidth)+20 {
				p.x = -20
			}
			if p.x < -20 {
				p.x = float64(engine.ScreenWidth) + 20
			}
			continue
		}

		if p.cloud {
			p.x += p.vx * dt
			if p.x > float64(engine.ScreenWidth)+float64(p.size) {
				p.x = -float64(p.size) * 2
			}
			if p.x < -float64(p.size)*2 {
				p.x = float64(engine.ScreenWidth) + float64(p.size)
			}
			continue
		}

		if p.water {
			p.timer += dt
			p.x = p.homeX + math.Sin(p.timer*1.5+p.baseY*0.1)*8
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
			tr, tg, tb := uint8(255), uint8(255), uint8(240)
			if p.r != 0 || p.g != 0 || p.b != 0 {
				tr, tg, tb = p.r, p.g, p.b
			}
			renderer.SetDrawColor(tr, tg, tb, uint8(a))
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
			continue
		}

		if p.fire {
			renderer.SetDrawColor(p.r, p.g, p.b, p.alpha)
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
			continue
		}

		if p.smoke {
			renderer.SetDrawColor(p.r, p.g, p.b, p.alpha)
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.y), W: p.size, H: p.size})
			if p.size > 2 {
				renderer.FillRect(&sdl.Rect{X: int32(p.x) + 1, Y: int32(p.y) - 1, W: p.size - 1, H: p.size + 1})
			}
			continue
		}

		if p.insect {
			px := int32(p.x)
			py := int32(p.y)
			wingSpread := int32(2 + math.Abs(math.Sin(p.timer*8))*2)
			renderer.SetDrawColor(p.r, p.g, p.b, p.alpha)
			renderer.FillRect(&sdl.Rect{X: px, Y: py, W: 2, H: 2})
			renderer.SetDrawColor(p.r, p.g, p.b, p.alpha / 2)
			renderer.FillRect(&sdl.Rect{X: px - wingSpread, Y: py - 1, W: wingSpread, H: 1})
			renderer.FillRect(&sdl.Rect{X: px + 2, Y: py - 1, W: wingSpread, H: 1})
			continue
		}

		if p.cloud {
			renderer.SetDrawColor(255, 255, 255, p.alpha)
			cx := int32(p.x)
			cy := int32(p.y)
			s := p.size
			renderer.FillRect(&sdl.Rect{X: cx, Y: cy, W: s, H: s / 3})
			renderer.FillRect(&sdl.Rect{X: cx + s/4, Y: cy - s/6, W: s / 2, H: s / 3})
			renderer.FillRect(&sdl.Rect{X: cx + s/6, Y: cy - s/4, W: s / 3, H: s / 4})
			continue
		}

		if p.water {
			a := uint8(float64(p.alpha) * (0.5 + 0.5*math.Sin(p.timer*2+p.baseY*0.1)))
			renderer.SetDrawColor(200, 220, 255, a)
			renderer.FillRect(&sdl.Rect{X: int32(p.x), Y: int32(p.baseY), W: p.size, H: 1})
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

package game

import (
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
)

// Ambient decorators. Once scenes load from JSON, the static parts
// (bg, npcs, hotspots, blockers, walk, spawn) live in assets/data/scenes/*.json
// but the procedural ambient — birds, butterflies, floating motes, clouds,
// warm glow — still reads best as Go code with rand + loops. Each scene's
// ambient lives in its own function here.
//
// JSON-authored spawners (kind + count + ranges) are a future refactor: the
// ambient is stable across passes and its data shape is bespoke enough per
// scene that authoring it declaratively would blow up the schema without
// saving meaningful lines.

// --- Paris ---

// decorateParisStreet adds cafe steam drifting up from the table cluster at
// left, plus dust motes in the afternoon air and a warm upper-sky glow. Paris
// skips birds/clouds on purpose — the plan wants it to feel urban, not rural.
func decorateParisStreet(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 3; i++ {
		baseX := 150 + float64(i)*80
		s.particles = append(s.particles, particle{
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
	for i := 0; i < 6; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.3) * 5,
			vy:    -rand.Float64()*1.0 - 0.2,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	s.glows = append(s.glows,
		glowEffect{x: 300, y: 0, w: 600, h: 400, r: 255, g: 245, b: 210, alpha: 10, pulse: 0.25},
		glowEffect{x: 50, y: 300, w: 200, h: 150, r: 255, g: 220, b: 160, alpha: 8, pulse: 0.3},
	)
}

// decorateParisLouvre adds the museum mood: dust motes swirling in the
// sunbeams from the glass pyramid, the main central light column, and two
// narrower side-window beams.
func decorateParisLouvre(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 15; i++ {
		s.particles = append(s.particles, particle{
			x:     400 + rand.Float64()*500,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.5 - 0.3,
			alpha: uint8(rand.Intn(20) + 8),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	s.glows = append(s.glows,
		glowEffect{x: 400, y: 50, w: 400, h: 500, r: 255, g: 240, b: 200, alpha: 10, pulse: 0.2},
		glowEffect{x: 200, y: 100, w: 150, h: 300, r: 255, g: 230, b: 180, alpha: 8, pulse: 0.4},
		glowEffect{x: 900, y: 100, w: 150, h: 300, r: 255, g: 230, b: 180, alpha: 8, pulse: 0.4},
	)
}

// decorateAirplaneFlight adds streaming cloud puffs whizzing past at speed —
// the cutscene sells travel by making the clouds move, not the plane.
func decorateAirplaneFlight(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 10; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     100 + rand.Float64()*400,
			vx:    -50 - rand.Float64()*40,
			alpha: uint8(rand.Intn(12) + 5),
			size:  int32(50 + rand.Intn(60)),
			cloud: true,
		})
	}
}

// --- Cabins ---
//
// All five cabins share the same basic shape: light dust motes in the
// cabin air + a warm window glow. Each kid's room differs in mote count
// and glow position; Lily's also has soft pastel butterflies. The shared
// helper cuts ~30 lines per room.

func cabinMotes(s *scene, n int) {
	for i := 0; i < n; i++ {
		s.particles = append(s.particles, particle{
			x:     300 + rand.Float64()*700,
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.5 - 0.1,
			alpha: uint8(rand.Intn(10) + 3),
			size:  int32(rand.Intn(2) + 1),
		})
	}
}

func decorateTommyRoom(s *scene) {
	if s == nil {
		return
	}
	cabinMotes(s, 8)
	s.glows = append(s.glows, glowEffect{
		x: 500, y: 150, w: 400, h: 350, r: 255, g: 240, b: 190, alpha: 8, pulse: 0.2,
	})
}

func decorateJakeRoom(s *scene) {
	if s == nil {
		return
	}
	cabinMotes(s, 6)
	s.glows = append(s.glows, glowEffect{
		x: 550, y: 150, w: 350, h: 300, r: 255, g: 240, b: 200, alpha: 7, pulse: 0.2,
	})
}

func decorateLilyRoom(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 6; i++ {
		s.particles = append(s.particles, particle{
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
		s.particles = append(s.particles, particle{
			x:      500 + rand.Float64()*400,
			baseY:  200 + rand.Float64()*150,
			vx:     (rand.Float64() - 0.5) * 10,
			alpha:  uint8(rand.Intn(35) + 40),
			insect: true,
			r:      c[0], g: c[1], b: c[2],
			timer: rand.Float64() * 10,
		})
	}
	s.glows = append(s.glows, glowEffect{
		x: 400, y: 100, w: 500, h: 400, r: 255, g: 245, b: 220, alpha: 8, pulse: 0.2,
	})
}

func decorateMarcusRoom(s *scene) {
	if s == nil {
		return
	}
	cabinMotes(s, 10)
	s.glows = append(s.glows, glowEffect{
		x: 500, y: 100, w: 400, h: 400, r: 255, g: 240, b: 190, alpha: 9, pulse: 0.25,
	})
}

func decorateDannyRoom(s *scene) {
	if s == nil {
		return
	}
	cabinMotes(s, 8)
	s.glows = append(s.glows, glowEffect{
		x: 400, y: 50, w: 300, h: 300, r: 255, g: 245, b: 210, alpha: 7, pulse: 0.2,
	})
}

// decorateCampOffice adds the lamp-lit indoor mood: a warm central glow from
// Higgins's desk lamp and a softer ambient glow from the window at left.
func decorateCampOffice(s *scene) {
	if s == nil {
		return
	}
	s.glows = append(s.glows,
		glowEffect{x: 600, y: 200, w: 300, h: 300, r: 255, g: 230, b: 170, alpha: 10, pulse: 0.3},
		glowEffect{x: 0, y: 100, w: 300, h: 400, r: 255, g: 245, b: 210, alpha: 8, pulse: 0.2},
	)
}

// decorateCampNight adds the nighttime campfire mood: fireflies twinkling in
// the darkness, an orange flame puff above the fire pit, the main fire glow,
// an all-over cool-blue darkness tint, and a soft warm patch near the cabins.
func decorateCampNight(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 16; i++ {
		s.particles = append(s.particles, particle{
			x:       80 + rand.Float64()*1200,
			y:       200 + rand.Float64()*300,
			twinkle: true,
			alpha:   uint8(rand.Intn(40) + 30),
			size:    1,
			r:       255, g: 255, b: 150,
		})
	}
	for i := 0; i < 6; i++ {
		s.particles = append(s.particles, particle{
			x:     622 + (rand.Float64()-0.5)*30,
			y:     568 - rand.Float64()*20,
			vx:    (rand.Float64() - 0.5) * 8,
			vy:    -rand.Float64()*25 - 10,
			alpha: uint8(rand.Intn(40) + 20),
			size:  int32(rand.Intn(2) + 1),
			baseY: 573,
			homeX: 622,
			fire:  true,
			r:     255, g: 160, b: 40,
		})
	}
	s.glows = append(s.glows,
		glowEffect{x: 510, y: 470, w: 260, h: 180, r: 255, g: 160, b: 40, alpha: 15, pulse: 3.5},
		glowEffect{x: 0, y: 0, w: engine.ScreenWidth, h: engine.ScreenHeight, r: 20, g: 15, b: 40, alpha: 12, pulse: 0.1},
		glowEffect{x: 820, y: 280, w: 180, h: 120, r: 255, g: 200, b: 100, alpha: 8, pulse: 1.5},
	)
}

// decorateCampLake adds the lake mood — drifting motes, birds, water shimmer
// ripples across the surface, dragonflies patrolling the water, and shoreline
// butterflies hovering near the flower patch.
func decorateCampLake(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 5; i++ {
		s.particles = append(s.particles, particle{
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
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     20 + rand.Float64()*60,
			vx:    dir,
			baseY: 20 + rand.Float64()*60,
			alpha: uint8(rand.Intn(25) + 40),
			size:  3,
			bird:  true,
		})
	}
	for i := 0; i < 12; i++ {
		baseX := 350 + float64(i)*60
		baseY := 380 + rand.Float64()*120
		s.particles = append(s.particles, particle{
			homeX: baseX,
			baseY: baseY,
			x:     baseX,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(20) + 15),
			water: true,
			timer: rand.Float64() * 10,
		})
	}
	for i := 0; i < 2; i++ {
		s.particles = append(s.particles, particle{
			x:      500 + rand.Float64()*400,
			baseY:  320 + rand.Float64()*60,
			vx:     6 + rand.Float64()*8,
			alpha:  uint8(rand.Intn(40) + 40),
			insect: true,
			r:      80, g: 160, b: 200,
			timer: rand.Float64() * 10,
		})
	}
	butterflyLakeColors := [][3]uint8{{232, 136, 43}, {240, 180, 80}, {210, 120, 60}}
	for i := 0; i < 3; i++ {
		c := butterflyLakeColors[i%len(butterflyLakeColors)]
		s.particles = append(s.particles, particle{
			x:      float64(150 + i*30),
			baseY:  float64(440 + i*6),
			vx:     4 + rand.Float64()*4,
			alpha:  uint8(rand.Intn(40) + 55),
			insect: true,
			r:      c[0], g: c[1], b: c[2],
			timer: rand.Float64() * 10,
		})
	}
	s.glows = append(s.glows,
		glowEffect{x: 400, y: 250, w: 500, h: 200, r: 255, g: 200, b: 120, alpha: 6, pulse: 0.2},
		glowEffect{x: 0, y: 0, w: engine.ScreenWidth, h: 200, r: 180, g: 150, b: 200, alpha: 5, pulse: 0.15},
		glowEffect{x: 350, y: 370, w: 600, h: 20, r: 200, g: 230, b: 255, alpha: 5, pulse: 1.5},
	)
}

// decorateCampGrounds adds the campsite mood — drifting motes, the campfire
// flame + smoke column, birds, clouds, fireflies at dusk, plus warm-light
// glows around the fire and the upper sky.
func decorateCampGrounds(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 10; i++ {
		s.particles = append(s.particles, particle{
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
		s.particles = append(s.particles, particle{
			x:     622 + (rand.Float64()-0.5)*20,
			y:     568 - rand.Float64()*30,
			vx:    (rand.Float64() - 0.5) * 12,
			vy:    -rand.Float64()*35 - 15,
			alpha: uint8(rand.Intn(50) + 30),
			size:  int32(rand.Intn(2) + 1),
			baseY: 573,
			homeX: 622,
			fire:  true,
			r:     c[0], g: c[1], b: c[2],
		})
	}
	for i := 0; i < 3; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     30 + rand.Float64()*60,
			vx:    15 + rand.Float64()*20,
			baseY: 30 + rand.Float64()*60,
			alpha: uint8(rand.Intn(30) + 50),
			size:  3,
			bird:  true,
		})
	}
	for i := 0; i < 2; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     40 + rand.Float64()*70,
			vx:    4 + rand.Float64()*3,
			alpha: uint8(rand.Intn(6) + 4),
			size:  int32(50 + rand.Intn(40)),
			cloud: true,
		})
	}
	for i := 0; i < 5; i++ {
		s.particles = append(s.particles, particle{
			x:     622 + (rand.Float64()-0.5)*12,
			y:     543 - rand.Float64()*20,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*12 - 6,
			alpha: uint8(rand.Intn(15) + 8),
			size:  int32(rand.Intn(3) + 2),
			baseY: 543,
			homeX: 622,
			smoke: true,
			r:     140, g: 130, b: 120,
			timer: rand.Float64() * 10,
		})
	}
	for i := 0; i < 6; i++ {
		s.particles = append(s.particles, particle{
			x:       100 + rand.Float64()*1100,
			y:       350 + rand.Float64()*150,
			twinkle: true,
			alpha:   uint8(rand.Intn(30) + 20),
			size:    1,
			r:       255, g: 255, b: 150,
		})
	}
	s.glows = append(s.glows,
		glowEffect{x: 200, y: 0, w: 800, h: 300, r: 255, g: 245, b: 200, alpha: 8, pulse: 0.2},
		glowEffect{x: 450, y: 400, w: 300, h: 100, r: 255, g: 200, b: 120, alpha: 6, pulse: 0.35},
		glowEffect{x: 560, y: 555, w: 130, h: 45, r: 255, g: 160, b: 40, alpha: 18, pulse: 4.0},
	)
}

// decorateCampEntrance adds the forest-entrance mood — drifting motes, a few
// birds, butterflies near the flower beds, slow clouds, and a soft warm glow
// at the top of the scene.
func decorateCampEntrance(s *scene) {
	if s == nil {
		return
	}
	for i := 0; i < 6; i++ {
		s.particles = append(s.particles, particle{
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
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     30 + rand.Float64()*60,
			vx:    dir,
			baseY: 30 + rand.Float64()*60,
			alpha: uint8(rand.Intn(30) + 50),
			size:  3,
			bird:  true,
		})
	}
	butterflyColors := [][3]uint8{{240, 200, 80}, {180, 120, 200}, {100, 180, 220}}
	for i := 0; i < 3; i++ {
		c := butterflyColors[i%len(butterflyColors)]
		s.particles = append(s.particles, particle{
			x:      300 + rand.Float64()*600,
			baseY:  200 + rand.Float64()*200,
			vx:     (rand.Float64() - 0.5) * 12,
			alpha:  uint8(rand.Intn(40) + 50),
			insect: true,
			r:      c[0], g: c[1], b: c[2],
			timer: rand.Float64() * 10,
		})
	}
	for i := 0; i < 2; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     30 + rand.Float64()*60,
			vx:    2 + rand.Float64()*3,
			alpha: uint8(rand.Intn(6) + 3),
			size:  int32(50 + rand.Intn(40)),
			cloud: true,
		})
	}
	s.glows = append(s.glows, glowEffect{
		x: 300, y: 0, w: 600, h: 350, r: 255, g: 245, b: 210, alpha: 10, pulse: 0.25,
	})
}

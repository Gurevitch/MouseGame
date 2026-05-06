package game

import (
	"math"
	"math/rand"
	"os"
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
	arrowDownRight
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
	// frameIdx / frameTimer drive the sprite-based ambient upgrade
	// (birds, butterflies). They stay zero for particle kinds that
	// still render as filled rectangles so nothing regresses.
	frameIdx   int
	frameTimer float64
}

// Ambient sprite sheets, loaded once per sceneManager construction.
// They're optional — if the PNGs aren't on disk yet the particle
// renderer falls back to the original filled-rect drawing. That lets
// us ship the code change ahead of the art landing.
var (
	ambientBirdFrames      []npcFrame
	ambientButterflyFrames []npcFrame
	ambientCloudTex        *sdl.Texture
	ambientCloudW          int32
	ambientCloudH          int32
	ambientLoaded          bool
)

// fileExists returns true when a sprite path is present and readable.
// We use os.Stat rather than trying to load the PNG first because a
// missing optional asset should be silent, not panic-inducing.
func ambientSpriteExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// loadAmbientSprites pre-loads the three shared sprites for camp Day 1
// flyovers. Safe to call multiple times; only the first call does real
// work. Missing files are ignored — the engine falls back to dot-style
// drawing.
func loadAmbientSprites(renderer *sdl.Renderer) {
	if ambientLoaded {
		return
	}
	ambientLoaded = true

	birdPath := "assets/images/ambient/bird_silhouette.png"
	if ambientSpriteExists(birdPath) {
		ambientBirdFrames = loadNPCGrid(renderer, birdPath, 8, 1)
	}

	butterflyPath := "assets/images/ambient/butterfly_flutter.png"
	if ambientSpriteExists(butterflyPath) {
		ambientButterflyFrames = loadNPCGrid(renderer, butterflyPath, 6, 1)
	}

	cloudPath := "assets/images/ambient/cloud_puff.png"
	if ambientSpriteExists(cloudPath) {
		tex, w, h := engine.TextureFromPNG(renderer, cloudPath)
		ambientCloudTex = tex
		ambientCloudW = w
		ambientCloudH = h
	}
}

// isCampOutdoorScene returns true for the three Day 1 camp outdoor
// scenes that host ambient life. Keeping the list in one place makes
// it easy to exclude interior/indoor/night scenes without touching
// every call site.
func isCampOutdoorScene(name string) bool {
	switch name {
	case "camp_entrance", "camp_grounds", "camp_lake":
		return true
	}
	return false
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
	// characterScale multiplies PP and NPC draw sizes (not hitboxes) so
	// tight/indoor scenes can render characters smaller without having
	// to redraw sheets. 1.0 = outdoor default; 0.85-0.9 for rooms. See
	// CHARACTERS.md.
	characterScale float64
}

// charScale returns the effective character-scale multiplier, defaulting
// to 1.0 when the scene wasn't explicitly configured. Keeps older scene
// definitions from rendering at 0.
func (s *scene) charScale() float64 {
	if s == nil || s.characterScale <= 0 {
		return 1.0
	}
	return s.characterScale
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

	loadAmbientSprites(renderer)

	// Phase 3 migration: scenes that have flipped to JSON are loaded here and
	// decorated with any procedural ambient. Hardcoded scene blocks below are
	// being emptied one at a time.
	sceneDefs := newSceneConfigStore("assets/data/scenes")

	// ===== Camp Chilly Wa Wa: Entrance (JSON-driven) =====
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "camp_entrance"); s != nil {
		decorateCampEntrance(s)
		sm.scenes["camp_entrance"] = s
	}

	// ===== Camp Chilly Wa Wa: Grounds (JSON-driven) =====
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "camp_grounds"); s != nil {
		decorateCampGrounds(s)
		sm.scenes["camp_grounds"] = s
	}

	// ===== Camp Chilly Wa Wa: Higgins' Office (JSON-driven) =====
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "camp_office"); s != nil {
		decorateCampOffice(s)
		sm.scenes["camp_office"] = s
	}

	// ===== Camp Chilly Wa Wa: Night (JSON-driven) =====
	// Higgins is silent here — the night cutscene drives his dialog via the
	// g.dialog helper so he appears to speak "in place" at the campfire.
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "camp_night"); s != nil {
		decorateCampNight(s)
		sm.scenes["camp_night"] = s
	}

	// ===== Camp Chilly Wa Wa: Lake (JSON-driven) =====
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "camp_lake"); s != nil {
		decorateCampLake(s)
		sm.scenes["camp_lake"] = s
	}

	// ===== Paris: Street + Louvre + Bakery (JSON-driven) =====
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "paris_street"); s != nil {
		decorateParisStreet(s)
		sm.scenes["paris_street"] = s
	}
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "paris_louvre"); s != nil {
		decorateParisLouvre(s)
		sm.scenes["paris_louvre"] = s
	}
	// User 2026-04-26: Paris bakery interior. NPC = bakery_woman, floor item
	// = rolling pin (registered in setupParisCallbacks). No decorate hook
	// yet — the JSON + factory are enough; add one later if ambient is
	// needed (e.g. flour particles).
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "paris_bakery"); s != nil {
		sm.scenes["paris_bakery"] = s
	}

	// ===== Cabin Interiors (JSON-driven) =====
	type roomSpec struct {
		name     string
		decorate func(*scene)
	}
	for _, r := range []roomSpec{
		{"tommy_room", decorateTommyRoom},
		{"jake_room", decorateJakeRoom},
		{"lily_room", decorateLilyRoom},
		{"marcus_room", decorateMarcusRoom},
		{"danny_room", decorateDannyRoom},
	} {
		if s := sm.loadSceneFromJSON(renderer, sceneDefs, r.name); s != nil {
			r.decorate(s)
			sm.scenes[r.name] = s
		}
	}

	// ===== Airplane Flight Cutscene (JSON-driven) =====
	if s := sm.loadSceneFromJSON(renderer, sceneDefs, "airplane_flight"); s != nil {
		decorateAirplaneFlight(s)
		sm.scenes["airplane_flight"] = s
	}

	// --- Extra city chapters (defined in their own files) ---
	addJerusalemScenes(sm, renderer)
	addTokyoScenes(sm, renderer)
	addRioScenes(sm, renderer)
	addRomeScenes(sm, renderer)
	addMexicoScenes(sm, renderer)

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
		if n.silent || n.hidden {
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
		case arrowDownRight:
			cx = hs.bounds.X + hs.bounds.W - 14
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
		case arrowDownRight:
			for i := int32(0); i < sz; i++ {
				renderer.SetDrawColor(255, 220, 100, a)
				renderer.FillRect(&sdl.Rect{X: cx - sz + i, Y: cy - sz + i, W: 1, H: 1})
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

	scale := s.charScale()

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
				n.drawScaled(renderer, scale)
			},
		})
	}

	if plr != nil {
		actors = append(actors, actorDraw{
			footY: plr.footY(),
			order: base + len(s.npcs),
			draw: func() {
				plr.drawScaled(renderer, scale)
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

// updateAmbient ticks all ambient particles. When showAmbientLife is
// false (Day 2+, night scenes, indoors, cities), bird / insect / cloud
// particles are skipped so the camp-specific wildlife freezes off
// without having to rebuild the particle list. Shimmer / smoke / fire
// / water particles keep running because they're scene-scoped effects
// (lake shimmer, campfire flames, cabin light).
func (s *scene) updateAmbient(dt float64, showAmbientLife bool) {
	for i := range s.particles {
		p := &s.particles[i]

		if p.twinkle {
			continue
		}

		if (p.bird || p.insect || p.cloud) && !showAmbientLife {
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
			if len(ambientButterflyFrames) > 0 {
				p.frameTimer += dt
				if p.frameTimer >= 0.12 {
					p.frameTimer = 0
					p.frameIdx = (p.frameIdx + 1) % len(ambientButterflyFrames)
				}
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
			if len(ambientBirdFrames) > 0 {
				p.frameTimer += dt
				if p.frameTimer >= 0.10 {
					p.frameTimer = 0
					p.frameIdx = (p.frameIdx + 1) % len(ambientBirdFrames)
				}
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

// drawAmbient draws glows and particles. Like updateAmbient, when
// showAmbientLife is false the flyover wildlife (birds, butterflies,
// clouds) is skipped so interior/Day-2/night scenes stay still.
func (s *scene) drawAmbient(renderer *sdl.Renderer, showAmbientLife bool) {
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
			if !showAmbientLife {
				continue
			}
			if len(ambientButterflyFrames) > 0 {
				f := ambientButterflyFrames[p.frameIdx%len(ambientButterflyFrames)]
				if f.tex != nil {
					dst := sdl.Rect{
						X: int32(p.x) - f.w/2,
						Y: int32(p.y) - f.h/2,
						W: f.w,
						H: f.h,
					}
					renderer.Copy(f.tex, nil, &dst)
					continue
				}
			}
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
			if !showAmbientLife {
				continue
			}
			if ambientCloudTex != nil {
				ratio := float64(p.size) / 80.0
				if ratio < 0.6 {
					ratio = 0.6
				}
				dstW := int32(float64(ambientCloudW) * ratio)
				dstH := int32(float64(ambientCloudH) * ratio)
				dst := sdl.Rect{
					X: int32(p.x),
					Y: int32(p.y),
					W: dstW,
					H: dstH,
				}
				modAlpha := 160 + int(p.alpha)
				if modAlpha > 255 {
					modAlpha = 255
				}
				ambientCloudTex.SetAlphaMod(uint8(modAlpha))
				renderer.Copy(ambientCloudTex, nil, &dst)
				continue
			}
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
			if !showAmbientLife {
				continue
			}
			if len(ambientBirdFrames) > 0 {
				f := ambientBirdFrames[p.frameIdx%len(ambientBirdFrames)]
				if f.tex != nil {
					flip := sdl.FLIP_NONE
					if p.vx < 0 {
						flip = sdl.FLIP_HORIZONTAL
					}
					dst := sdl.Rect{
						X: int32(p.x) - f.w/2,
						Y: int32(p.y) - f.h/2,
						W: f.w,
						H: f.h,
					}
					renderer.CopyEx(f.tex, nil, &dst, 0, nil, flip)
					continue
				}
			}
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

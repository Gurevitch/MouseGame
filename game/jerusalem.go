package game

import (
	"image/color"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Jerusalem chapter: Jake's courage arc.
//
// Three scenes, with the entrance plaza as a hub (user redesign 2026-06-08):
//   - jerusalem_entrance: the Western Wall plaza where PP lands. Open stone
//     courtyard with the Dome of the Rock behind. Two ways on: LEFT through
//     the arch to the market, RIGHT along the wall to get up close.
//   - jerusalem_wall: up at the Western Wall. Miriam the archeologist works
//     the tunnels behind the stones, Dov (her kid brother) hunts "glow-bugs"
//     with a flashlight, and worshippers sway at the base of the wall.
//   - jerusalem_market: the covered Old City souk — a stone tunnel of stalls.
//     Eli the spice seller works the SPICES stall and Gary the tourist
//     wanders here.
//
// Chapter anchor item: "Coin Rubbing" — Miriam hands it over at the wall once
// PP asks for something the boy can hold. Bringing it back to Jake's cabin
// triggers his healing dialogue. The "tunnels under a wall" thread from Jake's
// nightmare pays off here as the real Western Wall tunnels + the souk tunnel.
//
// The NPC roster follows the Hokus Pokus Pink / Passport to Peril template:
// *locals with their own problems*, not tour guides.

// ---------- Sprite paths ----------
//
// Jerusalem-specific NPC sheets aren't authored yet, so the NPCs borrow
// existing Paris/camp art (same cols×rows) via loadNPCGridPath: it picks the
// city sheet when it lands on disk, otherwise the fallback. Backgrounds and
// the worshipper overlay ARE authored (2026-06-08).

const (
	// Eli the spice seller — runs the SPICES stall on the left of the souk.
	// Art vendor sheet (8x2 split idle/talk) as placeholder until a Jerusalem
	// spice-merchant sheet is authored (see §JN1 in EXTRA_PROMPTS.md).
	jerArtEli     = "assets/images/locations/jerusalem/npc/npc_eli_idle.png"
	jerArtEliBack = "assets/images/locations/paris/npc/outside/npc_art_vendor.png"

	// Gary the tourist — reuse the security guard sheet (6x2)
	jerArtGary     = "assets/images/locations/jerusalem/npc/npc_gary_tourist.png"
	jerArtGaryBack = "assets/images/locations/paris/npc/outside/npc_security_guard.png"

	// Miriam the archeologist — French guide sheet layout
	jerArtMiriamIdle     = "assets/images/locations/jerusalem/npc/npc_miriam_idle.png"
	jerArtMiriamIdleBack = "assets/images/locations/paris/npc/outside/npc_french_guide_idle.png"
	jerArtMiriamTalk     = "assets/images/locations/jerusalem/npc/npc_miriam_talk.png"
	jerArtMiriamTalkBack = "assets/images/locations/paris/npc/outside/npc_french_guide_talk.png"

	// Dov, a curious local kid — borrows the camp kid idle sheet (8x2)
	jerArtDov         = "assets/images/locations/jerusalem/npc/npc_dov_idle.png"
	jerArtDovBack     = "assets/images/locations/camp/npc/kids/jake/npc_jake_idle.png"
	jerArtDovTalk     = "assets/images/locations/jerusalem/npc/npc_dov_talk.png"
	jerArtDovTalkBack = "assets/images/locations/camp/npc/kids/jake/npc_jake_talk.png"

	// Backgrounds (authored 2026-06-08)
	jerBgEntrance = "assets/images/locations/jerusalem/background/wall_enterence.png"
	jerBgWall     = "assets/images/locations/jerusalem/background/wall_close.png"
	jerBgMarket   = "assets/images/locations/jerusalem/background/market.png"
)

// Placeholder palette — warm limestone for the plaza/wall, dim amber for the souk
var (
	jerPlazaBase  = color.NRGBA{R: 214, G: 182, B: 140, A: 255}
	jerWallBase   = color.NRGBA{R: 224, G: 190, B: 120, A: 255}
	jerMarketBase = color.NRGBA{R: 120, G: 96, B: 68, A: 255}
)

// ---------- NPC dialogs ----------

var eliSpiceDialog = []dialogEntry{
	{speaker: "Eli", text: "Shalom! Za'atar, sumac, cardamom — ze finest spices in ze whole souk! Here, smell zis one."},
	{speaker: "Pink Panther", text: "*sniff* ...My eyes are watering and I have never been happier."},
	{speaker: "Eli", text: "Zat is ze real Jerusalem gold. Forget ze tourists and their little postcards."},
	{speaker: "Pink Panther", text: "Speaking of which — I'm chasing a boy's nightmare. Tunnels. A face in an old wall."},
	{speaker: "Eli", text: "Ah — ze Western Wall, and ze tunnels behind it. Miriam digs there."},
	{speaker: "Eli", text: "Out ze arch to ze plaza, then along to ze Wall. Tell her Eli sent you. She owes me for a kilo of cumin."},
}

var eliSpicePostDialog = []dialogEntry{
	{speaker: "Eli", text: "Miriam. Ze Wall. And my cumin money. Do not forget ze cumin money."},
}

var garyTouristDialog = []dialogEntry{
	{speaker: "Gary", text: "You again! Wait, have we met? I met a purple cat in Paris last Tuesday."},
	{speaker: "Pink Panther", text: "I am pink. I have never once been purple."},
	{speaker: "Gary", text: "Right! Gary. Retired dentist. Tourist 'til Thursday."},
	{speaker: "Gary", text: "Did you know these market stones are older than my fillings? Astonishing."},
	{speaker: "Pink Panther", text: "A truly remarkable comparison."},
	{speaker: "Gary", text: "If you find the old tunnels, mind your head. Jerusalem's tunnels are low and very opinionated."},
}

var garyTouristPostDialog = []dialogEntry{
	{speaker: "Gary", text: "Thursday! I fly home Thursday! Don't let me forget!"},
}

var miriamArchDialog = []dialogEntry{
	{speaker: "Miriam", text: "Careful by the scaffolding. The Wall's waited two thousand years; it can wait while you watch your step."},
	{speaker: "Pink Panther", text: "I'm chasing a kid's nightmare. Tunnels. A face staring out of the stones."},
	{speaker: "Miriam", text: "A face? That's not a nightmare. That's memory."},
	{speaker: "Miriam", text: "Last week I lifted a Roman coin from a crack in the tunnel wall. Emperor Hadrian, sharp as the day it was struck."},
	{speaker: "Miriam", text: "Kids come back from the tunnels with it in their dreams. The Wall remembers you right back."},
	{speaker: "Pink Panther", text: "Could I take something to the boy? Something he can hold?"},
	{speaker: "Miriam", text: "Here — a pencil rubbing of the coin. Paper carries what stone keeps quiet."},
	{speaker: "Miriam", text: "Tell him the face in his dream isn't chasing him. It's only remembering."},
}

var miriamArchPostDialog = []dialogEntry{
	{speaker: "Miriam", text: "Keep the rubbing dry. Wet paper forgets."},
}

var dovKidDialog = []dialogEntry{
	{speaker: "Dov", text: "Shhh! I'm hunting GLOW-BUGS. They only come where my flashlight hits the old stones."},
	{speaker: "Pink Panther", text: "How many have you caught?"},
	{speaker: "Dov", text: "Zero. But YESTERDAY I almost had one. My sister says I have to be brave like her."},
	{speaker: "Pink Panther", text: "Is Miriam your sister?"},
	{speaker: "Dov", text: "Yes! And she says fear is just a bug. You catch it and put it in your pocket."},
	{speaker: "Pink Panther", text: "I'll remember that."},
}

var dovKidPostDialog = []dialogEntry{
	{speaker: "Dov", text: "Fear is a bug. POCKET IT."},
}

// ---------- NPC constructors ----------

// newEli builds the souk spice seller — stands behind the SPICES stall on the
// left of the market and points PP toward Miriam at the Wall.
func newEli(renderer *sdl.Renderer, x int32) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtEli, jerArtEliBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtEli, jerArtEliBack, 8, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 400, W: 130, H: 240},
		name:           "Eli",
		dialog:         eliSpiceDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newGaryTourist(renderer *sdl.Renderer, x int32) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtGary, jerArtGaryBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtGary, jerArtGaryBack, 6, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 390, W: 120, H: 240},
		name:           "Gary",
		dialog:         garyTouristDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newMiriam(renderer *sdl.Renderer, x int32) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, jerArtMiriamIdle, jerArtMiriamIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, jerArtMiriamTalk, jerArtMiriamTalkBack, 8, 1),
		bounds:         sdl.Rect{X: x, Y: 470, W: 130, H: 230},
		name:           "Miriam",
		dialog:         miriamArchDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newDov(renderer *sdl.Renderer, x int32) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtDov, jerArtDovBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtDovTalk, jerArtDovTalkBack, 8, 2, 0),
		bounds:         sdl.Rect{X: x, Y: 500, W: 100, H: 200},
		name:           "Dov",
		dialog:         dovKidDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

// ---------- Scene builders ----------

// addJerusalemScenes wires the three Jerusalem scenes into the given manager.
// Called from newSceneManager so load is deterministic.
func addJerusalemScenes(sm *sceneManager, renderer *sdl.Renderer) {
	// ===== Entrance plaza (hub — PP lands here) =====
	entrance := &scene{
		name:   "jerusalem_entrance",
		bg:     newPNGBackgroundOr(renderer, jerBgEntrance, jerPlazaBase),
		spawnX: 250,
		spawnY: 600,
		hotspots: []hotspot{
			{
				// LEFT through the arch to the souk
				bounds:      sdl.Rect{X: 20, Y: 330, W: 190, H: 250},
				targetScene: "jerusalem_market",
				name:        "To the Market",
				arrow:       arrowLeft,
			},
			{
				// RIGHT along the wall to get up close
				bounds:      sdl.Rect{X: 1180, Y: 150, W: 196, H: 520},
				targetScene: "jerusalem_wall",
				name:        "To the Wall",
				arrow:       arrowRight,
			},
		},
		minY: 470,
		maxY: 660,
	}
	// Warm dust hanging in the plaza air + a soft top-of-sky glow
	for i := 0; i < 10; i++ {
		entrance.particles = append(entrance.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.2 - 0.2,
			alpha: uint8(rand.Intn(12) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	entrance.glows = []glowEffect{
		{x: 0, y: 0, w: 1400, h: 280, r: 255, g: 240, b: 200, alpha: 10, pulse: 0.25},
	}
	// Tiny worshippers at the wall in the mid-distance (right side, small)
	entrance.ambientSprites = append(entrance.ambientSprites,
		newAmbientWorshippers(renderer, 1000, 470, 0.45),
	)
	sm.scenes["jerusalem_entrance"] = entrance

	// ===== Up at the Western Wall =====
	wall := &scene{
		name:           "jerusalem_wall",
		bg:             newPNGBackgroundOr(renderer, jerBgWall, jerWallBase),
		npcs:           []*npc{newMiriam(renderer, 470), newDov(renderer, 980)},
		spawnX:         220,
		spawnY:         680,
		characterScale: 0.85,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 250, W: 110, H: 420},
				targetScene: "jerusalem_entrance",
				name:        "Back to the Plaza",
				arrow:       arrowLeft,
			},
		},
		minY: 640,
		maxY: 710,
	}
	// Worshippers swaying at the foot of the wall + warm light wash
	wall.ambientSprites = append(wall.ambientSprites,
		newAmbientWorshippers(renderer, 1120, 720, 0.7),
	)
	wall.glows = []glowEffect{
		{x: 0, y: 0, w: engine.ScreenWidth, h: 300, r: 255, g: 235, b: 180, alpha: 10, pulse: 0.2},
	}
	for i := 0; i < 6; i++ {
		wall.particles = append(wall.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.9 - 0.1,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	sm.scenes["jerusalem_wall"] = wall

	// ===== Old City market (souk tunnel) =====
	market := &scene{
		name:           "jerusalem_market",
		bg:             newPNGBackgroundOr(renderer, jerBgMarket, jerMarketBase),
		npcs:           []*npc{newEli(renderer, 200), newGaryTourist(renderer, 980)},
		spawnX:         260,
		spawnY:         640,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 250, W: 120, H: 420},
				targetScene: "jerusalem_entrance",
				name:        "Back to the Plaza",
				arrow:       arrowLeft,
			},
		},
		minY: 580,
		maxY: 700,
	}
	// Hanging-lamp glow down the souk + faint incense/dust haze
	market.glows = []glowEffect{
		{x: 100, y: 50, w: 300, h: 250, r: 255, g: 190, b: 110, alpha: 12, pulse: 1.4},
		{x: 980, y: 50, w: 300, h: 250, r: 255, g: 180, b: 100, alpha: 10, pulse: 1.0},
		{x: 560, y: 120, w: 256, h: 300, r: 255, g: 240, b: 200, alpha: 8, pulse: 0.3},
	}
	for i := 0; i < 8; i++ {
		market.particles = append(market.particles, particle{
			x:     200 + rand.Float64()*1000,
			y:     rand.Float64() * 450,
			vx:    (rand.Float64() - 0.5) * 2,
			vy:    -rand.Float64()*0.7 - 0.1,
			alpha: uint8(rand.Intn(12) + 6),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	sm.scenes["jerusalem_market"] = market
}

// setupJerusalemCallbacks wires NPC post-dialog swaps, Miriam's coin rubbing
// handoff, and the travel-map return on the entrance plaza. Call after scenes
// are built.
func (g *Game) setupJerusalemCallbacks() {
	game := g

	// Entrance plaza: travel-map return (this is where PP arrives, so it's
	// where he leaves Jerusalem from — open the map at the top of the plaza).
	if entrance, ok := g.sceneMgr.scenes["jerusalem_entrance"]; ok {
		entrance.hotspots = append(entrance.hotspots, hotspot{
			bounds: sdl.Rect{X: 540, Y: 0, W: 300, H: 90},
			name:   "Travel Map",
			arrow:  arrowUp,
			onInteract: func() bool {
				game.travelMap.Show("jerusalem_entrance")
				return true
			},
		})
	}

	// Market: Eli + Gary post-dialog swaps.
	if market, ok := g.sceneMgr.scenes["jerusalem_market"]; ok {
		for _, n := range market.npcs {
			switch n.name {
			case "Eli":
				eli := n
				eli.onDialogEnd = func() {
					eli.dialog = eliSpicePostDialog
				}
			case "Gary":
				gary := n
				gary.onDialogEnd = func() {
					gary.dialog = garyTouristPostDialog
				}
			}
		}
	}

	// Wall: Miriam hands over the coin rubbing; Dov post-dialog swap.
	if wall, ok := g.sceneMgr.scenes["jerusalem_wall"]; ok {
		for _, n := range wall.npcs {
			switch n.name {
			case "Miriam":
				miriam := n
				miriam.onDialogEnd = func() {
					// Hand PP the coin rubbing on the first conversation
					miriam.dialog = miriamArchPostDialog
					if !game.inv.hasItem("Coin Rubbing") {
						if item := game.items.createItem("coin_rubbing"); item != nil {
							game.inv.addItem(item)
						}
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "A rubbing of an ancient coin. Jake needs to see this."},
						})
					}
				}
			case "Dov":
				dov := n
				dov.onDialogEnd = func() {
					dov.dialog = dovKidPostDialog
				}
			}
		}
	}
}

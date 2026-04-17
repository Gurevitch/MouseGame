package game

import (
	"image/color"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Tokyo chapter: Lily's voice arc.
//
// Two scenes:
//   - tokyo_street: torii gate street with ramen stalls and petals.
//     Hiro (ramen cook) and Kenji (calligraphy student) live here.
//   - tokyo_temple: quiet sakura garden behind the temple. Oba-chan, an
//     elderly flower-arranger, works beneath a pressed-petal awning, and
//     Gary the tourist is back, lost in a guidebook.
//
// Chapter anchor item: "Pressed Sakura" — Oba-chan presses a petal into
// parchment and hands it over when PP asks about quiet things.
//
// Bring the pressed sakura to Lily's cabin to finish her arc.

const (
	// Background art
	tokBgStreet     = "assets/images/locations/tokyo/background/tokyo_street.png"
	tokBgTemple     = "assets/images/locations/tokyo/background/tokyo_temple.png"

	// NPC art: each Tokyo local has a preferred path and a fallback that
	// reuses an existing sheet with matching cols/rows.
	tokArtHiroIdle     = "assets/images/locations/tokyo/npc/npc_hiro_idle.png"
	tokArtHiroIdleBack = "assets/images/locations/paris/npc/npc_art_vendor.png"

	tokArtKenjiIdle     = "assets/images/locations/tokyo/npc/npc_kenji_idle.png"
	tokArtKenjiIdleBack = "assets/images/locations/paris/npc/npc_security_guard.png"

	tokArtObaIdle     = "assets/images/locations/tokyo/npc/npc_obachan_idle.png"
	tokArtObaIdleBack = "assets/images/locations/paris/npc/npc_french_guide_idle.png"
	tokArtObaTalk     = "assets/images/locations/tokyo/npc/npc_obachan_talk.png"
	tokArtObaTalkBack = "assets/images/locations/paris/npc/npc_french_guide_talk.png"

	tokArtGaryIdle     = "assets/images/locations/tokyo/npc/npc_gary_idle.png"
	tokArtGaryIdleBack = "assets/images/locations/paris/npc/npc_security_guard.png"
)

var (
	tokStreetBase = color.NRGBA{R: 255, G: 196, B: 208, A: 255} // soft sakura pink
	tokTempleBase = color.NRGBA{R: 140, G: 102, B: 84, A: 255}  // temple wood
)

// ---------- Dialogs ----------

var hiroRamenDialog = []dialogEntry{
	{speaker: "Hiro", text: "Irasshaimase! One bowl tonkotsu ramen for ze pink-haired gentleman?"},
	{speaker: "Pink Panther", text: "Not tonight. I'm looking for someone who helps quiet people find their words."},
	{speaker: "Hiro", text: "Ah. Ze garden behind ze temple. Oba-chan arranges flowers zere."},
	{speaker: "Hiro", text: "She does not say much. She lets ze petals speak."},
	{speaker: "Hiro", text: "If she likes you, she will press a sakura into paper for you."},
}

var hiroRamenPostDialog = []dialogEntry{
	{speaker: "Hiro", text: "Come back for noodles when ze heart is heavy, panther-san."},
}

var kenjiStudentDialog = []dialogEntry{
	{speaker: "Kenji", text: "Please, do not nudge — my brush is mid-stroke."},
	{speaker: "Pink Panther", text: "What are you writing?"},
	{speaker: "Kenji", text: "Ze kanji for 'voice'. My sister said hers has gone quiet."},
	{speaker: "Kenji", text: "I write it every morning. Ink is how my family speaks across ze ocean."},
}

var kenjiStudentPostDialog = []dialogEntry{
	{speaker: "Kenji", text: "Ink is just a voice zat takes its time."},
}

var obachanDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "..."},
	{speaker: "Pink Panther", text: "Hello, madame. A friend is losing her voice — not literally, but inside. Do you know a gift for that?"},
	{speaker: "Oba-chan", text: "Mm. Ze quiet ones cannot hear their own words. We must give zem something to hold."},
	{speaker: "Oba-chan", text: "Here — a petal I pressed zis morning. Sakura keeps its shape but loses its weight."},
	{speaker: "Oba-chan", text: "Give her ze petal. Let her practice being light."},
}

var obachanPostDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "..."},
	{speaker: "Pink Panther", text: "Thank you."},
}

var garyTokyoDialog = []dialogEntry{
	{speaker: "Gary", text: "Pink panther! Jerusalem last week, Tokyo zis week! Are you a spy?"},
	{speaker: "Pink Panther", text: "Retired dentist, remember?"},
	{speaker: "Gary", text: "Oh right, zat's me. Guidebook says ze temple is zat way. Or zat way. Guidebook is upside down."},
}

var garyTokyoPostDialog = []dialogEntry{
	{speaker: "Gary", text: "Flipped ze book! Now it says Tokyo is in Peru! Remarkable!"},
}

// ---------- NPC constructors ----------

func newHiroRamen(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, tokArtHiroIdle, tokArtHiroIdleBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, tokArtHiroIdle, tokArtHiroIdleBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 320, Y: 360, W: 130, H: 240},
		name:           "Hiro",
		dialog:         hiroRamenDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newKenjiStudent(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, tokArtKenjiIdle, tokArtKenjiIdleBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, tokArtKenjiIdle, tokArtKenjiIdleBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 940, Y: 380, W: 120, H: 240},
		name:           "Kenji",
		dialog:         kenjiStudentDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newObachan(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, tokArtObaIdle, tokArtObaIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, tokArtObaTalk, tokArtObaTalkBack, 8, 1),
		bounds:         sdl.Rect{X: 500, Y: 370, W: 130, H: 240},
		name:           "Oba-chan",
		dialog:         obachanDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newGaryTokyo(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, tokArtGaryIdle, tokArtGaryIdleBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, tokArtGaryIdle, tokArtGaryIdleBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 1080, Y: 370, W: 120, H: 240},
		name:           "Gary",
		dialog:         garyTokyoDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

// ---------- Scene builders ----------

func addTokyoScenes(sm *sceneManager, renderer *sdl.Renderer) {
	street := &scene{
		name:   "tokyo_street",
		bg:     newPNGBackgroundOr(renderer, tokBgStreet, tokStreetBase),
		npcs:   []*npc{newHiroRamen(renderer), newKenjiStudent(renderer)},
		spawnX: 200,
		spawnY: 450,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "tokyo_temple",
				name:        "To the Temple Garden",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 150, H: 500},
		},
		minY: 380,
		maxY: 640,
	}
	// Falling sakura petals
	for i := 0; i < 18; i++ {
		street.particles = append(street.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * float64(engine.ScreenHeight),
			vx:    (rand.Float64() - 0.4) * 20,
			vy:    15 + rand.Float64()*10,
			alpha: uint8(rand.Intn(30) + 50),
			size:  int32(rand.Intn(3) + 2),
			r:     255, g: 180, b: 200,
		})
	}
	sm.scenes["tokyo_street"] = street

	temple := &scene{
		name:           "tokyo_temple",
		bg:             newPNGBackgroundOr(renderer, tokBgTemple, tokTempleBase),
		npcs:           []*npc{newObachan(renderer), newGaryTokyo(renderer)},
		spawnX:         200,
		spawnY:         450,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "tokyo_street",
				name:        "Back to the Street",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{
			{X: 1300, Y: 0, W: 100, H: 500},
		},
		minY: 380,
		maxY: 640,
	}
	temple.glows = []glowEffect{
		{x: 200, y: 100, w: 500, h: 300, r: 255, g: 220, b: 230, alpha: 12, pulse: 0.3},
	}
	for i := 0; i < 12; i++ {
		temple.particles = append(temple.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.4) * 10,
			vy:    10 + rand.Float64()*8,
			alpha: uint8(rand.Intn(20) + 40),
			size:  int32(rand.Intn(2) + 2),
			r:     255, g: 180, b: 200,
		})
	}
	sm.scenes["tokyo_temple"] = temple
}

func (g *Game) setupTokyoCallbacks() {
	game := g

	if street, ok := g.sceneMgr.scenes["tokyo_street"]; ok {
		for _, n := range street.npcs {
			switch n.name {
			case "Hiro":
				hiro := n
				hiro.onDialogEnd = func() {
					hiro.dialog = hiroRamenPostDialog
				}
			case "Kenji":
				kenji := n
				kenji.onDialogEnd = func() {
					kenji.dialog = kenjiStudentPostDialog
				}
			}
		}
		street.hotspots = append(street.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.showTravelMap = true
				game.travelMapFrom = "tokyo_street"
				return true
			},
		})
	}

	if temple, ok := g.sceneMgr.scenes["tokyo_temple"]; ok {
		for _, n := range temple.npcs {
			switch n.name {
			case "Oba-chan":
				oba := n
				oba.onDialogEnd = func() {
					oba.dialog = obachanPostDialog
					if !game.inv.hasItem("Pressed Sakura") {
						if item := game.items.createItem("pressed_sakura"); item != nil {
							game.inv.addItem(item)
						}
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "A pressed sakura. Lily will hold it and understand."},
						})
					}
				}
			case "Gary":
				gary := n
				gary.onDialogEnd = func() {
					gary.dialog = garyTokyoPostDialog
				}
			}
		}
	}
}

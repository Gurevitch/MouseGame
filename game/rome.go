package game

import (
	"image/color"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Rome chapter: Danny's identity arc.
//
// Two scenes:
//   - rome_street: Nonna Rosa's pasta stall out in front of the Colosseum.
//     Luca (street musician) and Garibaldi the cat hang around here.
//   - rome_colosseum: Inside the arena. Dottor Bianchi the classicist
//     runs a pencil-rubbing tour, and the player can collect the anchor
//     item: "Inscription Rubbing".
//
// Bring the rubbing home, show it to Danny, watch him realise the letters
// spell his own name in Latin. Healing unlocks Mexico City (finale).

const (
	romeBgStreet    = "assets/images/locations/rome/background/rome_street.png"
	romeBgColosseum = "assets/images/locations/rome/background/rome_colosseum.png"

	romeArtNonnaIdle     = "assets/images/locations/rome/npc/npc_nonna_idle.png"
	romeArtNonnaIdleBack = "assets/images/locations/paris/npc/npc_french_guide_idle.png"
	romeArtNonnaTalk     = "assets/images/locations/rome/npc/npc_nonna_talk.png"
	romeArtNonnaTalkBack = "assets/images/locations/paris/npc/npc_french_guide_talk.png"

	romeArtLucaIdle     = "assets/images/locations/rome/npc/npc_luca_idle.png"
	romeArtLucaIdleBack = "assets/images/locations/paris/npc/npc_art_vendor.png"

	romeArtDottorIdle     = "assets/images/locations/rome/npc/npc_dottor_idle.png"
	romeArtDottorIdleBack = "assets/images/locations/paris/npc/npc_security_guard.png"

	romeArtGariIdle     = "assets/images/locations/rome/npc/npc_garibaldi_cat_idle.png"
	romeArtGariIdleBack = "assets/images/locations/paris/npc/npc_art_vendor.png"
)

var (
	romeStreetBase = color.NRGBA{R: 235, G: 196, B: 140, A: 255}
	romeArenaBase  = color.NRGBA{R: 120, G: 100, B: 80, A: 255}
)

// ---------- Dialogs ----------

var nonnaRosaDialog = []dialogEntry{
	{speaker: "Nonna Rosa", text: "Rosa! Carbonara! Five euro! You eat, you smile, you go!"},
	{speaker: "Pink Panther", text: "I'm looking for someone who can read Latin inscriptions. A boy back home keeps drawing one."},
	{speaker: "Nonna Rosa", text: "Ah. Dottor Bianchi, ze Colosseum. He is old but his eyes still read ze stones."},
	{speaker: "Nonna Rosa", text: "Tell him Nonna sent you. He owes me two bowls of pasta."},
}

var nonnaRosaPostDialog = []dialogEntry{
	{speaker: "Nonna Rosa", text: "TWO bowls, Bianchi! Tell him!"},
}

var lucaMusicianDialog = []dialogEntry{
	{speaker: "Luca", text: "Oi! Want a song? One euro. 'O Sole Mio' or Metallica — your choice!"},
	{speaker: "Pink Panther", text: "Metallica. On an accordion. I'd pay to see that."},
	{speaker: "Luca", text: "Zat is ze SPIRIT, panther! Maybe later. First zou must find Dottor Bianchi. He is in ze arena."},
}

var lucaMusicianPostDialog = []dialogEntry{
	{speaker: "Luca", text: "Metallica. Accordion. Later!"},
}

var dottorBianchiDialog = []dialogEntry{
	{speaker: "Dottor Bianchi", text: "Shh. Chalk work. Ze marble cannot tolerate ze breath of ze impatient."},
	{speaker: "Pink Panther", text: "Sorry. I'm looking for a specific inscription. A kid back at camp keeps drawing Roman arches around his own name."},
	{speaker: "Dottor Bianchi", text: "Around his name? Hm. Most Roman arches were dedicated to emperors. A few to families."},
	{speaker: "Dottor Bianchi", text: "Tell me ze child's name."},
	{speaker: "Pink Panther", text: "Danny. But he signs his drawings 'D.M.' in weird curly letters."},
	{speaker: "Dottor Bianchi", text: "D.M. — ze letters appear on every tomb. Dis Manibus. 'To ze spirits of ze departed.'"},
	{speaker: "Dottor Bianchi", text: "But also — see? Danillus Marcus. Ze arch on ze east side bears zat name."},
	{speaker: "Dottor Bianchi", text: "Here. A fresh rubbing, still warm from ze chalk. Take it to your boy."},
}

var dottorBianchiPostDialog = []dialogEntry{
	{speaker: "Dottor Bianchi", text: "Protect ze paper. Ze chalk crumbles."},
}

var garibaldiCatDialog = []dialogEntry{
	{speaker: "Garibaldi", text: "Mrrrow."},
	{speaker: "Pink Panther", text: "Excellent point."},
	{speaker: "Garibaldi", text: "Prrrt."},
	{speaker: "Pink Panther", text: "I entirely agree."},
}

var garibaldiCatPostDialog = []dialogEntry{
	{speaker: "Garibaldi", text: "..."},
}

// ---------- NPC constructors ----------

func newNonnaRosa(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, romeArtNonnaIdle, romeArtNonnaIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, romeArtNonnaTalk, romeArtNonnaTalkBack, 8, 1),
		bounds:         sdl.Rect{X: 320, Y: 360, W: 130, H: 240},
		name:           "Nonna Rosa",
		dialog:         nonnaRosaDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newLucaMusician(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, romeArtLucaIdle, romeArtLucaIdleBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, romeArtLucaIdle, romeArtLucaIdleBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 900, Y: 380, W: 130, H: 240},
		name:           "Luca",
		dialog:         lucaMusicianDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newDottorBianchi(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, romeArtDottorIdle, romeArtDottorIdleBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, romeArtDottorIdle, romeArtDottorIdleBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 520, Y: 370, W: 130, H: 240},
		name:           "Dottor Bianchi",
		dialog:         dottorBianchiDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newGaribaldiCat(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, romeArtGariIdle, romeArtGariIdleBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, romeArtGariIdle, romeArtGariIdleBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 1080, Y: 510, W: 90, H: 100},
		name:           "Garibaldi",
		dialog:         garibaldiCatDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.16,
	}
}

// ---------- Scene builders ----------

func addRomeScenes(sm *sceneManager, renderer *sdl.Renderer) {
	street := &scene{
		name:   "rome_street",
		bg:     newPNGBackgroundOr(renderer, romeBgStreet, romeStreetBase),
		npcs:   []*npc{newNonnaRosa(renderer), newLucaMusician(renderer), newGaribaldiCat(renderer)},
		spawnX: 200,
		spawnY: 450,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "rome_colosseum",
				name:        "Into the Colosseum",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{{X: 0, Y: 0, W: 150, H: 500}},
		minY:     380, maxY: 640,
	}
	for i := 0; i < 8; i++ {
		street.particles = append(street.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.0 - 0.2,
			alpha: uint8(rand.Intn(12) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	sm.scenes["rome_street"] = street

	colosseum := &scene{
		name:           "rome_colosseum",
		bg:             newPNGBackgroundOr(renderer, romeBgColosseum, romeArenaBase),
		npcs:           []*npc{newDottorBianchi(renderer)},
		spawnX:         200,
		spawnY:         450,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "rome_street",
				name:        "Back to Street",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{{X: 1300, Y: 0, W: 100, H: 500}},
		minY:     380, maxY: 640,
	}
	colosseum.glows = []glowEffect{
		{x: 300, y: 50, w: 800, h: 400, r: 255, g: 220, b: 180, alpha: 12, pulse: 0.25},
	}
	sm.scenes["rome_colosseum"] = colosseum
}

func (g *Game) setupRomeCallbacks() {
	game := g
	if s, ok := g.sceneMgr.scenes["rome_street"]; ok {
		for _, n := range s.npcs {
			switch n.name {
			case "Nonna Rosa":
				nr := n
				nr.onDialogEnd = func() { nr.dialog = nonnaRosaPostDialog }
			case "Luca":
				luca := n
				luca.onDialogEnd = func() { luca.dialog = lucaMusicianPostDialog }
			case "Garibaldi":
				gc := n
				gc.onDialogEnd = func() { gc.dialog = garibaldiCatPostDialog }
			}
		}
		s.hotspots = append(s.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.showTravelMap = true
				game.travelMapFrom = "rome_street"
				return true
			},
		})
	}
	if s, ok := g.sceneMgr.scenes["rome_colosseum"]; ok {
		for _, n := range s.npcs {
			if n.name == "Dottor Bianchi" {
				doc := n
				doc.onDialogEnd = func() {
					doc.dialog = dottorBianchiPostDialog
					if !game.inv.hasItem("Inscription Rubbing") {
						if item := game.items.createItem("inscription_rubbing"); item != nil {
							game.inv.addItem(item)
						}
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "'Danillus Marcus.' That's Danny's name in Latin. He's been writing himself onto walls."},
						})
					}
				}
			}
		}
	}
}

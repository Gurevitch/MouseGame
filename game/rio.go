package game

import (
	"image/color"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Rio + Buenos Aires chapter: Tommy's "missing sister" arc.
// Two cities, one chapter. The anchor item is a two-half dance card:
// PP finds the first half in Rio (Tio Jorge's bar) and the second half
// in Buenos Aires (Don Rafa's tango school). Joined, the card names the
// sister Tommy has been screaming about in his sleep.

const (
	rioBgStreet = "assets/images/locations/rio/background/rio_street.png"
	rioBgBar    = "assets/images/locations/rio/background/rio_bar.png"
	baBgStreet  = "assets/images/locations/ba/background/buenos_aires_street.png"
	baBgTango   = "assets/images/locations/ba/background/ba_tango_school.png"

	// --- Rio locals ---
	rioArtJorgeIdle     = "assets/images/locations/rio/npc/npc_tio_jorge_idle.png"
	rioArtJorgeIdleBack = "assets/images/locations/paris/npc/npc_art_vendor.png"
	rioArtMarisaIdle     = "assets/images/locations/rio/npc/npc_marisa_idle.png"
	rioArtMarisaIdleBack = "assets/images/locations/paris/npc/npc_french_guide_idle.png"
	rioArtMarisaTalk     = "assets/images/locations/rio/npc/npc_marisa_talk.png"
	rioArtMarisaTalkBack = "assets/images/locations/paris/npc/npc_french_guide_talk.png"
	rioArtPadreIdle     = "assets/images/locations/rio/npc/npc_padre_idle.png"
	rioArtPadreIdleBack = "assets/images/locations/paris/npc/npc_security_guard.png"
	rioArtBrunoIdle     = "assets/images/locations/rio/npc/npc_bruno_kid_idle.png"
	rioArtBrunoIdleBack = "assets/images/locations/camp/npc/kids/tommy/npc_tommy_idle.png"
	rioArtBrunoTalk     = "assets/images/locations/rio/npc/npc_bruno_kid_talk.png"
	rioArtBrunoTalkBack = "assets/images/locations/camp/npc/kids/tommy/npc_tommy_talk.png"

	// --- Buenos Aires locals ---
	baArtRafaIdle     = "assets/images/locations/ba/npc/npc_don_rafa_idle.png"
	baArtRafaIdleBack = "assets/images/locations/paris/npc/npc_art_vendor.png"
	baArtLuciaIdle     = "assets/images/locations/ba/npc/npc_lucia_idle.png"
	baArtLuciaIdleBack = "assets/images/locations/paris/npc/npc_french_guide_idle.png"
	baArtLuciaTalk     = "assets/images/locations/ba/npc/npc_lucia_talk.png"
	baArtLuciaTalkBack = "assets/images/locations/paris/npc/npc_french_guide_talk.png"
	baArtPacoIdle     = "assets/images/locations/ba/npc/npc_paco_idle.png"
	baArtPacoIdleBack = "assets/images/locations/paris/npc/npc_security_guard.png"
	baArtGaryIdle     = "assets/images/locations/ba/npc/npc_gary_idle.png"
	baArtGaryIdleBack = "assets/images/locations/paris/npc/npc_security_guard.png"
)

var (
	rioStreetBase = color.NRGBA{R: 255, G: 170, B: 100, A: 255} // sunset palette
	rioBarBase    = color.NRGBA{R: 85, G: 55, B: 80, A: 255}
	baStreetBase  = color.NRGBA{R: 180, G: 140, B: 210, A: 255} // violet dusk
	baTangoBase   = color.NRGBA{R: 70, G: 40, B: 55, A: 255}
)

// ---------- Dialogs ----------

var tioJorgeDialog = []dialogEntry{
	{speaker: "Tio Jorge", text: "Eh, rosado! Sit, sit. One caipirinha? No? Guarana then."},
	{speaker: "Pink Panther", text: "I'm looking for anything a boy back at camp has been screaming in his sleep. 'Marisa.'"},
	{speaker: "Tio Jorge", text: "Marisa! Nu, every bar in Copacabana has a Marisa. Which one?"},
	{speaker: "Tio Jorge", text: "Ask my daughter in ze back room. She keeps all ze old dance cards from Carnival."},
}

var tioJorgePostDialog = []dialogEntry{
	{speaker: "Tio Jorge", text: "Marisa in ze back, rosado."},
}

var marisaBartenderDialog = []dialogEntry{
	{speaker: "Marisa", text: "Every Carnival since 1987 I keep ze cards. My father says I am sentimental."},
	{speaker: "Pink Panther", text: "A boy at camp keeps shouting 'Marisa' like he knows her. His handwriting looks like yours."},
	{speaker: "Marisa", text: "Dios meu. Take zis — half of a dance card I tore in 1991. ZE other half... I never found."},
	{speaker: "Marisa", text: "Ze other half has ze partner's name. Maybe in Buenos Aires. My cousin ran away zere."},
}

var marisaBartenderPostDialog = []dialogEntry{
	{speaker: "Marisa", text: "If you find ze other half... tell me."},
}

var padreAntonioDialog = []dialogEntry{
	{speaker: "Padre Antonio", text: "Peace, senhor. Ze chapel bell is heavy today."},
	{speaker: "Pink Panther", text: "Do families often get split between Rio and Buenos Aires, padre?"},
	{speaker: "Padre Antonio", text: "Too often. Ze river is kinder zan ze border. Ze children carry both sides in ze chest."},
}

var padreAntonioPostDialog = []dialogEntry{
	{speaker: "Padre Antonio", text: "Go with peace, senhor."},
}

var brunoKidDialog = []dialogEntry{
	{speaker: "Bruno", text: "Can you juggle? My uncle can juggle FIVE ORANGES."},
	{speaker: "Pink Panther", text: "Alas, only three paws. How about a pink panther who can dance?"},
	{speaker: "Bruno", text: "Dance like Carnival! NOT like opera! Opera is BORING."},
}

var brunoKidPostDialog = []dialogEntry{
	{speaker: "Bruno", text: "Not opera. Carnival!"},
}

// Buenos Aires

var donRafaDialog = []dialogEntry{
	{speaker: "Don Rafa", text: "Ah, a panther in my school? Ze tango does not discriminate!"},
	{speaker: "Pink Panther", text: "I'm looking for a torn dance card. The other half is in Rio."},
	{speaker: "Don Rafa", text: "Rio and Buenos Aires — always zey share ze halves. Check ze lost-and-found box by ze piano."},
}

var donRafaPostDialog = []dialogEntry{
	{speaker: "Don Rafa", text: "Ze piano box, panther. Go on."},
}

var luciaTangoDialog = []dialogEntry{
	{speaker: "Lucia", text: "Ay, who taught you to stand like zat? Knees soft, senhor."},
	{speaker: "Pink Panther", text: "I'm more of a leans-against-a-wall dancer."},
	{speaker: "Lucia", text: "Pity. Here — if you want ze piano box, it is under ze red shawl. Do not touch ze shawl."},
	{speaker: "Lucia", text: "Ze shawl is my grandmother's. She is ninety-six and she will know."},
}

var luciaTangoPostDialog = []dialogEntry{
	{speaker: "Lucia", text: "Do NOT touch ze shawl."},
}

var pacoTangoDialog = []dialogEntry{
	{speaker: "Paco", text: "You look lost, amigo. Piano box is red shawl, knees soft, heart open."},
	{speaker: "Pink Panther", text: "That's a lot of instructions for one box."},
}

var pacoTangoPostDialog = []dialogEntry{
	{speaker: "Paco", text: "Knees soft!"},
}

var garyBADialog = []dialogEntry{
	{speaker: "Gary", text: "PINK PANTHER! Buenos Aires! I knew ze guidebook was wrong about Tokyo!"},
	{speaker: "Pink Panther", text: "Gary, how are you everywhere at once?"},
	{speaker: "Gary", text: "One of ze great mysteries of retirement, friend. Also I have a Eurail pass."},
}

var garyBAPostDialog = []dialogEntry{
	{speaker: "Gary", text: "Eurail pass! Also works in Argentina apparently!"},
}

// ---------- NPC constructors ----------

func newTioJorge(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, rioArtJorgeIdle, rioArtJorgeIdleBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, rioArtJorgeIdle, rioArtJorgeIdleBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 320, Y: 360, W: 130, H: 240},
		name:           "Tio Jorge",
		dialog:         tioJorgeDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newMarisaBartender(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, rioArtMarisaIdle, rioArtMarisaIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, rioArtMarisaTalk, rioArtMarisaTalkBack, 8, 1),
		bounds:         sdl.Rect{X: 480, Y: 370, W: 130, H: 240},
		name:           "Marisa",
		dialog:         marisaBartenderDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newPadreAntonio(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, rioArtPadreIdle, rioArtPadreIdleBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, rioArtPadreIdle, rioArtPadreIdleBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 940, Y: 380, W: 120, H: 240},
		name:           "Padre Antonio",
		dialog:         padreAntonioDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newBrunoKid(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, rioArtBrunoIdle, rioArtBrunoIdleBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, rioArtBrunoTalk, rioArtBrunoTalkBack, 8, 2, 0),
		bounds:         sdl.Rect{X: 1080, Y: 400, W: 100, H: 200},
		name:           "Bruno",
		dialog:         brunoKidDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newDonRafa(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, baArtRafaIdle, baArtRafaIdleBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, baArtRafaIdle, baArtRafaIdleBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 320, Y: 360, W: 130, H: 240},
		name:           "Don Rafa",
		dialog:         donRafaDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newLuciaTango(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, baArtLuciaIdle, baArtLuciaIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, baArtLuciaTalk, baArtLuciaTalkBack, 8, 1),
		bounds:         sdl.Rect{X: 520, Y: 370, W: 130, H: 240},
		name:           "Lucia",
		dialog:         luciaTangoDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newPacoTango(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, baArtPacoIdle, baArtPacoIdleBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, baArtPacoIdle, baArtPacoIdleBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 900, Y: 380, W: 120, H: 240},
		name:           "Paco",
		dialog:         pacoTangoDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newGaryBA(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, baArtGaryIdle, baArtGaryIdleBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, baArtGaryIdle, baArtGaryIdleBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 1080, Y: 380, W: 120, H: 240},
		name:           "Gary",
		dialog:         garyBADialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

// ---------- Scene builders ----------

func addRioScenes(sm *sceneManager, renderer *sdl.Renderer) {
	rioStreet := &scene{
		name:   "rio_street",
		bg:     newPNGBackgroundOr(renderer, rioBgStreet, rioStreetBase),
		npcs:   []*npc{newTioJorge(renderer), newPadreAntonio(renderer)},
		spawnX: 200,
		spawnY: 450,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "rio_bar",
				name:        "Into the Bar",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{{X: 0, Y: 0, W: 150, H: 500}},
		minY:     380, maxY: 640,
	}
	for i := 0; i < 8; i++ {
		rioStreet.particles = append(rioStreet.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.4) * 10,
			vy:    -rand.Float64()*0.6 - 0.2,
			alpha: uint8(rand.Intn(12) + 6),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	rioStreet.glows = []glowEffect{
		{x: 0, y: 0, w: 1400, h: 300, r: 255, g: 190, b: 130, alpha: 14, pulse: 0.25},
	}
	sm.scenes["rio_street"] = rioStreet

	rioBar := &scene{
		name:           "rio_bar",
		bg:             newPNGBackgroundOr(renderer, rioBgBar, rioBarBase),
		npcs:           []*npc{newMarisaBartender(renderer), newBrunoKid(renderer)},
		spawnX:         200,
		spawnY:         450,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "rio_street",
				name:        "Back to Street",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{{X: 1300, Y: 0, W: 100, H: 500}},
		minY:     380, maxY: 640,
	}
	rioBar.glows = []glowEffect{
		{x: 400, y: 100, w: 400, h: 400, r: 255, g: 180, b: 120, alpha: 14, pulse: 0.35},
	}
	sm.scenes["rio_bar"] = rioBar

	baStreet := &scene{
		name:   "buenos_aires_street",
		bg:     newPNGBackgroundOr(renderer, baBgStreet, baStreetBase),
		npcs:   []*npc{newDonRafa(renderer), newGaryBA(renderer)},
		spawnX: 200,
		spawnY: 450,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "ba_tango_school",
				name:        "Into the Tango School",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{{X: 0, Y: 0, W: 150, H: 500}},
		minY:     380, maxY: 640,
	}
	sm.scenes["buenos_aires_street"] = baStreet

	baTango := &scene{
		name:           "ba_tango_school",
		bg:             newPNGBackgroundOr(renderer, baBgTango, baTangoBase),
		npcs:           []*npc{newLuciaTango(renderer), newPacoTango(renderer)},
		spawnX:         200,
		spawnY:         450,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "buenos_aires_street",
				name:        "Back to Street",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{{X: 1300, Y: 0, W: 100, H: 500}},
		minY:     380, maxY: 640,
	}
	baTango.glows = []glowEffect{
		{x: 200, y: 50, w: 800, h: 300, r: 255, g: 200, b: 150, alpha: 12, pulse: 0.4},
	}
	sm.scenes["ba_tango_school"] = baTango
}

func (g *Game) setupRioCallbacks() {
	game := g
	// Rio street: post-dialog swaps + travel map pin
	if s, ok := g.sceneMgr.scenes["rio_street"]; ok {
		for _, n := range s.npcs {
			switch n.name {
			case "Tio Jorge":
				jorge := n
				jorge.onDialogEnd = func() { jorge.dialog = tioJorgePostDialog }
			case "Padre Antonio":
				padre := n
				padre.onDialogEnd = func() { padre.dialog = padreAntonioPostDialog }
			}
		}
		s.hotspots = append(s.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.travelMap.Show("rio_street")
				return true
			},
		})
	}

	// Rio bar: Marisa hands out the first half of the dance card
	if s, ok := g.sceneMgr.scenes["rio_bar"]; ok {
		for _, n := range s.npcs {
			switch n.name {
			case "Marisa":
				marisa := n
				marisa.onDialogEnd = func() {
					marisa.dialog = marisaBartenderPostDialog
					// Only hand out the card once (we need the BA half before
					// the player can claim a complete "Dance Card").
					if !game.vars.GetBool(ScopeGame, "dance_card_rio_done") {
						game.vars.SetBool(ScopeGame, "dance_card_rio_done", true)
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "Half a dance card with 'Marisa' on it. Need the other half now."},
						})
					}
				}
			case "Bruno":
				bruno := n
				bruno.onDialogEnd = func() { bruno.dialog = brunoKidPostDialog }
			}
		}
	}

	// BA street
	if s, ok := g.sceneMgr.scenes["buenos_aires_street"]; ok {
		for _, n := range s.npcs {
			switch n.name {
			case "Don Rafa":
				rafa := n
				rafa.onDialogEnd = func() { rafa.dialog = donRafaPostDialog }
			case "Gary":
				gary := n
				gary.onDialogEnd = func() { gary.dialog = garyBAPostDialog }
			}
		}
		s.hotspots = append(s.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.travelMap.Show("buenos_aires_street")
				return true
			},
		})
	}

	// BA tango school: Lucia holds the second half. Only now does PP actually
	// get the "Dance Card" inventory item, representing both halves joined.
	if s, ok := g.sceneMgr.scenes["ba_tango_school"]; ok {
		for _, n := range s.npcs {
			switch n.name {
			case "Lucia":
				lucia := n
				lucia.onDialogEnd = func() {
					lucia.dialog = luciaTangoPostDialog
					if game.vars.GetBool(ScopeGame, "dance_card_rio_done") && !game.inv.hasItem("Dance Card") {
						if item := game.items.createItem("dance_card"); item != nil {
							game.inv.addItem(item)
						}
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "Both halves. The card says 'Tommy y Marisa — Carnival 1991'."},
							{speaker: "Pink Panther", text: "Time to go home."},
						})
					}
				}
			case "Paco":
				paco := n
				paco.onDialogEnd = func() { paco.dialog = pacoTangoPostDialog }
			}
		}
	}
}

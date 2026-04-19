package game

import (
	"image/color"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Jerusalem chapter: Jake's courage arc.
//
// Two scenes:
//   - jerusalem_street: sunny stone plaza near the Western Wall. Locals
//     Eli (kiosk owner) and Gary the Tourist hang out here.
//   - jerusalem_tunnel: cool underground dig tunnels. Miriam the amateur
//     archeologist works the wall, and Dov (kid with a flashlight) pokes
//     around behind her.
//
// Chapter anchor item: "Coin Rubbing" — Miriam hands it over once PP says
// the magic word "tunnels". Bringing it back to Jake's cabin triggers his
// healing dialogue.
//
// The NPC roster intentionally follows the Hokus Pokus Pink / Passport to
// Peril template: *locals with their own problems*, not tour guides.

// ---------- Fallback sprite paths ----------
//
// Until Jerusalem-specific PNGs are authored the NPCs borrow existing art
// so the city is playable. loadNPCGridPath picks the city-specific sheet
// when it shows up on disk; until then we reuse Paris sheets that share
// the same (cols, rows) layout.

const (
	// Eli kiosk owner — art vendor sheet (8x2 split idle/talk)
	jerArtEli     = "assets/images/locations/jerusalem/npc/npc_eli_idle.png"
	jerArtEliBack = "assets/images/locations/paris/npc/npc_art_vendor.png"

	// Gary the tourist — we reuse the security guard sheet (6x2)
	jerArtGary     = "assets/images/locations/jerusalem/npc/npc_gary_tourist.png"
	jerArtGaryBack = "assets/images/locations/paris/npc/npc_security_guard.png"

	// Miriam the archeologist — French guide sheet layout
	jerArtMiriamIdle     = "assets/images/locations/jerusalem/npc/npc_miriam_idle.png"
	jerArtMiriamIdleBack = "assets/images/locations/paris/npc/npc_french_guide_idle.png"
	jerArtMiriamTalk     = "assets/images/locations/jerusalem/npc/npc_miriam_talk.png"
	jerArtMiriamTalkBack = "assets/images/locations/paris/npc/npc_french_guide_talk.png"

	// Dov, a curious local kid — borrows the camp kid idle sheet (8x2)
	jerArtDov     = "assets/images/locations/jerusalem/npc/npc_dov_idle.png"
	jerArtDovBack = "assets/images/locations/camp/npc/kids/jake/npc_jake_idle.png"
	jerArtDovTalk     = "assets/images/locations/jerusalem/npc/npc_dov_talk.png"
	jerArtDovTalkBack = "assets/images/locations/camp/npc/kids/jake/npc_jake_talk.png"

	// Backgrounds
	jerBgStreet     = "assets/images/locations/jerusalem/background/jerusalem_street.png"
	jerBgStreetBack = ""                                                                  // use placeholder
	jerBgTunnel     = "assets/images/locations/jerusalem/background/jerusalem_tunnel.png"
	jerBgTunnelBack = ""
)

// Placeholder palette — warm limestone for the plaza, cool stone for the tunnels
var (
	jerStreetBase = color.NRGBA{R: 214, G: 182, B: 140, A: 255}
	jerTunnelBase = color.NRGBA{R: 68, G: 62, B: 78, A: 255}
)

// ---------- NPC dialogs ----------

var eliKioskDialog = []dialogEntry{
	{speaker: "Eli", text: "Shalom! Five shekel for ze best mango juice in ze Old City."},
	{speaker: "Pink Panther", text: "Mango juice on top of a postcard stand? What won't you sell?"},
	{speaker: "Eli", text: "Sunglasses. I do not sell sunglasses. Too much drama."},
	{speaker: "Pink Panther", text: "I'm looking for a boy who dreams about tunnels. Under a wall."},
	{speaker: "Eli", text: "Tunnels! Nu, of course. Miriam is digging beneath ze plaza right now."},
	{speaker: "Eli", text: "Take ze staircase past ze gate. Tell her Eli sent you. She still owes me juice money."},
}

var eliKioskPostDialog = []dialogEntry{
	{speaker: "Eli", text: "Miriam. Tunnels. Juice. Don't forget ze juice."},
}

var garyTouristDialog = []dialogEntry{
	{speaker: "Gary", text: "You again! Wait, have we met? I met a purple cat in Paris last week."},
	{speaker: "Pink Panther", text: "I am pink. And I have never been purple."},
	{speaker: "Gary", text: "Oh! Right. Gary. Tourist. Retired dentist. Here 'til Thursday."},
	{speaker: "Gary", text: "Did you know every stone in this wall is older than my fillings? Amazing."},
	{speaker: "Pink Panther", text: "That is... a remarkable comparison."},
	{speaker: "Gary", text: "If you go underground, watch your head. Tunnels in Jerusalem are short and opinionated."},
}

var garyTouristPostDialog = []dialogEntry{
	{speaker: "Gary", text: "Thursday! Remember! Thursday I fly home!"},
}

var miriamArchDialog = []dialogEntry{
	{speaker: "Miriam", text: "Mind the low beam. Archeologists don't respect tall mammals."},
	{speaker: "Pink Panther", text: "I'm looking for anything connected to a kid's nightmares. Tunnels, a face in the wall..."},
	{speaker: "Miriam", text: "Hm. A face? Not nightmares. Memory."},
	{speaker: "Miriam", text: "Last month I pulled a Roman coin out of this wall. Emperor Hadrian's face, still crisp."},
	{speaker: "Miriam", text: "Kids come home from here with it in their dreams. That wall remembers you back."},
	{speaker: "Pink Panther", text: "Can I take something back for the boy? Something he can touch?"},
	{speaker: "Miriam", text: "Here — a pencil rubbing of the coin. Paper holds what stone can't say."},
	{speaker: "Miriam", text: "Tell him the face in his dream isn't hunting him. It's remembering."},
}

var miriamArchPostDialog = []dialogEntry{
	{speaker: "Miriam", text: "Don't let the paper get wet. The rubbing fades."},
}

var dovKidDialog = []dialogEntry{
	{speaker: "Dov", text: "Shhhh! I'm hunting GLOW-BUGS. They only come out in the dig lamp."},
	{speaker: "Pink Panther", text: "How many have you caught?"},
	{speaker: "Dov", text: "Zero. But YESTERDAY I almost had one. My sister says I have to be brave like Miriam."},
	{speaker: "Pink Panther", text: "Is Miriam your sister?"},
	{speaker: "Dov", text: "Yes! And she says fear is just a bug. You just catch it and put it in your pocket."},
	{speaker: "Pink Panther", text: "I'll remember that."},
}

var dovKidPostDialog = []dialogEntry{
	{speaker: "Dov", text: "Fear is a bug. POCKET IT."},
}

// ---------- NPC constructors ----------

func newEli(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtEli, jerArtEliBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtEli, jerArtEliBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 320, Y: 360, W: 130, H: 240},
		name:           "Eli",
		dialog:         eliKioskDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newGaryTourist(renderer *sdl.Renderer, x int32) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtGary, jerArtGaryBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtGary, jerArtGaryBack, 6, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 370, W: 120, H: 240},
		name:           "Gary",
		dialog:         garyTouristDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newMiriam(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, jerArtMiriamIdle, jerArtMiriamIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, jerArtMiriamTalk, jerArtMiriamTalkBack, 8, 1),
		bounds:         sdl.Rect{X: 520, Y: 360, W: 130, H: 240},
		name:           "Miriam",
		dialog:         miriamArchDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newDov(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtDov, jerArtDovBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtDovTalk, jerArtDovTalkBack, 8, 2, 0),
		bounds:         sdl.Rect{X: 900, Y: 390, W: 100, H: 200},
		name:           "Dov",
		dialog:         dovKidDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

// ---------- Scene builders ----------

// addJerusalemScenes wires the two Jerusalem scenes into the given manager.
// Called from newSceneManager so load is deterministic.
func addJerusalemScenes(sm *sceneManager, renderer *sdl.Renderer) {
	street := &scene{
		name:   "jerusalem_street",
		bg:     newPNGBackgroundOr(renderer, jerBgStreet, jerStreetBase),
		npcs:   []*npc{newEli(renderer), newGaryTourist(renderer, 1080)},
		spawnX: 200,
		spawnY: 450,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 1300, Y: 200, W: 100, H: 400},
				targetScene: "jerusalem_tunnel",
				name:        "To the Tunnels",
				arrow:       arrowRight,
			},
		},
		blockers: []sdl.Rect{
			{X: 0, Y: 0, W: 150, H: 500},
		},
		minY: 380,
		maxY: 640,
	}
	// Dust motes in the stone plaza
	for i := 0; i < 10; i++ {
		street.particles = append(street.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.2 - 0.2,
			alpha: uint8(rand.Intn(12) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	street.glows = []glowEffect{
		{x: 0, y: 0, w: 1400, h: 300, r: 255, g: 240, b: 200, alpha: 10, pulse: 0.25},
	}
	sm.scenes["jerusalem_street"] = street

	tunnel := &scene{
		name:           "jerusalem_tunnel",
		bg:             newPNGBackgroundOr(renderer, jerBgTunnel, jerTunnelBase),
		npcs:           []*npc{newMiriam(renderer), newDov(renderer)},
		spawnX:         200,
		spawnY:         450,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
				targetScene: "jerusalem_street",
				name:        "Back to the Plaza",
				arrow:       arrowLeft,
			},
		},
		blockers: []sdl.Rect{
			{X: 1300, Y: 0, W: 100, H: 500},
		},
		minY: 380,
		maxY: 640,
	}
	// Flickering lamp glow
	tunnel.glows = []glowEffect{
		{x: 400, y: 100, w: 250, h: 400, r: 255, g: 200, b: 140, alpha: 14, pulse: 1.2},
		{x: 800, y: 50, w: 200, h: 400, r: 255, g: 180, b: 120, alpha: 10, pulse: 0.8},
	}
	for i := 0; i < 6; i++ {
		tunnel.particles = append(tunnel.particles, particle{
			x:     400 + rand.Float64()*600,
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.5) * 2,
			vy:    -rand.Float64()*0.8 - 0.2,
			alpha: uint8(rand.Intn(15) + 8),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	sm.scenes["jerusalem_tunnel"] = tunnel
}

// setupJerusalemCallbacks wires NPC post-dialog swaps and Miriam's coin
// rubbing handoff. Call after scenes are built.
func (g *Game) setupJerusalemCallbacks() {
	game := g

	if street, ok := g.sceneMgr.scenes["jerusalem_street"]; ok {
		for _, n := range street.npcs {
			switch n.name {
			case "Eli":
				eli := n
				eli.onDialogEnd = func() {
					eli.dialog = eliKioskPostDialog
				}
			case "Gary":
				gary := n
				gary.onDialogEnd = func() {
					gary.dialog = garyTouristPostDialog
				}
			}
		}
		// Travel map hotspot (return to map on the left edge, same as Paris)
		street.hotspots = append(street.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.travelMap.Show("jerusalem_street")
				return true
			},
		})
	}

	if tunnel, ok := g.sceneMgr.scenes["jerusalem_tunnel"]; ok {
		for _, n := range tunnel.npcs {
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

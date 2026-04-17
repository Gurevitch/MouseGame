package game

import (
	"image/color"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Mexico City chapter: the finale.
//
// Only one scene — the plaza in front of the Palacio de Bellas Artes,
// where all five healed kids (Marcus, Jake, Lily, Tommy, Danny) meet up
// in a single frame for the ending monologue. Director Higgins stands
// to the side wearing a straw hat.
//
// This chapter's "anchor" is simply the player arriving with every heal
// flag set. No item puzzle — the story has already been paid for.

const (
	mxBgPlaza       = "assets/images/locations/mexico/background/mexico_plaza.png"
	mxArtMariachi      = "assets/images/locations/mexico/npc/npc_mariachi_idle.png"
	mxArtMariachiBack  = "assets/images/locations/paris/npc/npc_art_vendor.png"
	mxArtAbuelaIdle     = "assets/images/locations/mexico/npc/npc_abuela_idle.png"
	mxArtAbuelaIdleBack = "assets/images/locations/paris/npc/npc_french_guide_idle.png"
	mxArtAbuelaTalk     = "assets/images/locations/mexico/npc/npc_abuela_talk.png"
	mxArtAbuelaTalkBack = "assets/images/locations/paris/npc/npc_french_guide_talk.png"
	mxArtVendor     = "assets/images/locations/mexico/npc/npc_vendor_idle.png"
	mxArtVendorBack = "assets/images/locations/paris/npc/npc_security_guard.png"
)

var mxPlazaBase = color.NRGBA{R: 240, G: 170, B: 120, A: 255}

// ---------- Dialogs ----------

var mariachiDialog = []dialogEntry{
	{speaker: "Mariachi", text: "Hola, amigo! One song for ze lady? Oh — no lady. One song for ze panther!"},
	{speaker: "Pink Panther", text: "I'm not much of a singer, but I appreciate the offer."},
	{speaker: "Mariachi", text: "Zen you must DANCE. Ze plaza is waiting."},
}

var mariachiPostDialog = []dialogEntry{
	{speaker: "Mariachi", text: "Dance, amigo, dance!"},
}

var abuelaDialog = []dialogEntry{
	{speaker: "Abuela", text: "Niño, are you lost? Ze plaza is always full of lost people."},
	{speaker: "Pink Panther", text: "Not lost. I'm bringing five kids home."},
	{speaker: "Abuela", text: "Home! Home is a word zat costs nothing and weighs everything. Go."},
}

var abuelaPostDialog = []dialogEntry{
	{speaker: "Abuela", text: "Home, niño."},
}

var mxVendorDialog = []dialogEntry{
	{speaker: "Vendor", text: "Tamales, elote, agua fresca — everything for ze reunion!"},
	{speaker: "Pink Panther", text: "You already know about the reunion?"},
	{speaker: "Vendor", text: "Five kids, one panther — ze whole city has been waiting for you."},
}

var mxVendorPostDialog = []dialogEntry{
	{speaker: "Vendor", text: "Elote! Five pesos!"},
}

var finaleMonologue = []dialogEntry{
	{speaker: "Pink Panther", text: "Five cities. Five names. Five kids."},
	{speaker: "Pink Panther", text: "Turns out they weren't broken. They were connected — to places they had never seen, to histories they had never lived."},
	{speaker: "Pink Panther", text: "Marcus to Paris. Jake to Jerusalem. Lily to Tokyo. Tommy to Rio and Buenos Aires. Danny to Rome."},
	{speaker: "Pink Panther", text: "And all of them, somehow, to this plaza."},
	{speaker: "Director Higgins", text: "PP. You did it. Every cabin is smiling again."},
	{speaker: "Pink Panther", text: "I didn't do it, Higgins. They did. I just carried the postcards."},
	{speaker: "Director Higgins", text: "Stay for dinner?"},
	{speaker: "Pink Panther", text: "Only if there's pie."},
}

// ---------- NPC constructors ----------

func newMariachi(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, mxArtMariachi, mxArtMariachiBack, 8, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, mxArtMariachi, mxArtMariachiBack, 8, 2, 1),
		bounds:         sdl.Rect{X: 320, Y: 360, W: 130, H: 240},
		name:           "Mariachi",
		dialog:         mariachiDialog,
		bobAmount:      2,
		talkFrameSpeed: 0.12,
	}
}

func newAbuela(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridPath(renderer, mxArtAbuelaIdle, mxArtAbuelaIdleBack, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, mxArtAbuelaTalk, mxArtAbuelaTalkBack, 8, 1),
		bounds:         sdl.Rect{X: 500, Y: 370, W: 130, H: 240},
		name:           "Abuela",
		dialog:         abuelaDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

func newMexicanVendor(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGridRowPath(renderer, mxArtVendor, mxArtVendorBack, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, mxArtVendor, mxArtVendorBack, 6, 2, 1),
		bounds:         sdl.Rect{X: 1080, Y: 370, W: 120, H: 240},
		name:           "Vendor",
		dialog:         mxVendorDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
}

// ---------- Scene builder ----------

func addMexicoScenes(sm *sceneManager, renderer *sdl.Renderer) {
	plaza := &scene{
		name:   "mexico_street",
		bg:     newPNGBackgroundOr(renderer, mxBgPlaza, mxPlazaBase),
		npcs:   []*npc{newMariachi(renderer), newAbuela(renderer), newMexicanVendor(renderer)},
		spawnX: 200,
		spawnY: 450,
		blockers: []sdl.Rect{{X: 0, Y: 0, W: 150, H: 500}},
		minY:     380, maxY: 640,
	}
	for i := 0; i < 14; i++ {
		plaza.particles = append(plaza.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.4) * 10,
			vy:    -rand.Float64()*0.8 - 0.2,
			alpha: uint8(rand.Intn(18) + 8),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	plaza.glows = []glowEffect{
		{x: 0, y: 0, w: 1400, h: 300, r: 255, g: 200, b: 150, alpha: 14, pulse: 0.3},
	}
	sm.scenes["mexico_street"] = plaza
}

func (g *Game) setupMexicoCallbacks() {
	game := g
	if s, ok := g.sceneMgr.scenes["mexico_street"]; ok {
		for _, n := range s.npcs {
			switch n.name {
			case "Mariachi":
				m := n
				m.onDialogEnd = func() { m.dialog = mariachiPostDialog }
			case "Abuela":
				a := n
				a.onDialogEnd = func() { a.dialog = abuelaPostDialog }
			case "Vendor":
				v := n
				v.onDialogEnd = func() { v.dialog = mxVendorPostDialog }
			}
		}
		s.hotspots = append(s.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.showTravelMap = true
				game.travelMapFrom = "mexico_street"
				return true
			},
		})
	}
}

// triggerFinaleMonologue runs the ending speech once the player enters
// mexico_street with every heal flag set. Hooked from Update() so it
// plays once, the moment the plaza fades in.
func (g *Game) triggerFinaleMonologue() {
	if g.vars.GetBool(ScopeGame, "finale_monologue_played") {
		return
	}
	if !g.vars.GetBool(ScopeGame, VarMarcusHealed) ||
		!g.vars.GetBool(ScopeGame, VarJakeHealed) ||
		!g.vars.GetBool(ScopeGame, VarLilyHealed) ||
		!g.vars.GetBool(ScopeGame, VarTommyHealed) ||
		!g.vars.GetBool(ScopeGame, VarDannyHealed) {
		return
	}
	g.vars.SetBool(ScopeGame, "finale_monologue_played", true)
	g.SetChapter(ChapterFinale)
	g.dialog.startDialog(finaleMonologue)
}

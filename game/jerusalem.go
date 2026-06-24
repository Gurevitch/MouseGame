package game

import (
	"image/color"
	"math/rand"
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Jerusalem chapter: Jake's courage arc (rebuilt 2026-06-21, #26).
//
// New retro daisy-chain (replaces the old trivial "Miriam hands the coin
// rubbing" stub). Three scenes, the entrance plaza is the hub:
//   - jerusalem_entrance (plaza): PP lands here. Shimon stands by the
//     separation fence and directs the player: UP-arrow (right) goes to the
//     Wall, LEFT-arrow to the market. A bagel (ka'ak) seller works the plaza.
//   - jerusalem_market (souk): a coffee seller (centre) and a spice seller.
//   - jerusalem_wall: a praying man (his idle IS praying; he turns to talk,
//     then turns back) and two kids, one rehearsing his bar-mitzvah portion.
//     Worshippers sway along the base.
//
// The chain:
//   spices seller -> Cardamom -> coffee seller (sits + teaches) -> Coffee ->
//   bagel seller -> Bagel -> praying man (note custom) -> Note Paper ->
//   Shimon -> Pen -> write + place the note in the Wall (jer_note_placed) ->
//   Shimon -> COIN.
// The Coin is Jake's anchor item (it replaced the old "Coin Rubbing"); placing
// the note gates the flight home.
//
// All Jerusalem NPC art is still to be authored, so the NPCs borrow existing
// Paris/camp sheets via loadNPCGridRowPath (it prefers the city sheet when it
// lands on disk). New PP/NPC one-shots and item icons no-op gracefully until
// their art lands (prompts queued at EXTRA_PROMPTS §JERUSALEM).

// ---------- Sprite paths (city sheet preferred, fallback to existing art) ----------

const (
	jerBgEntrance = "assets/images/locations/jerusalem/background/wall_enterence.png"
	jerBgWall     = "assets/images/locations/jerusalem/background/wall_close.png"
	jerBgMarket   = "assets/images/locations/jerusalem/background/market.png"

	// Art is organised into wall/ (plaza + Wall NPCs) and market/ (souk NPCs).
	// Shimon's full-body 6x2 sheet has landed; the rest borrow Paris/camp art
	// via the fallbacks below until their sheets are authored (§JERUSALEM).

	jerNPCWall   = "assets/images/locations/jerusalem/npc/wall/"
	jerNPCMarket = "assets/images/locations/jerusalem/npc/market/"

	// --- wall / plaza ---
	jerArtShimon     = jerNPCWall + "npc_shimon.png"      // 6x2 (idle row0, talk row1) - LANDED
	jerArtShimonGive = jerNPCWall + "npc_shimon_give.png" // 8x1 give one-shot (pen / coin)
	jerArtBagel      = jerNPCWall + "npc_bagel_seller.png"
	jerArtBagelGive  = jerNPCWall + "npc_bagel_seller_give.png"
	jerArtPrayIdle   = jerNPCWall + "npc_praying_man_idle.png"
	jerArtPrayTalk   = jerNPCWall + "npc_praying_man_talk.png"
	jerArtPrayGive   = jerNPCWall + "npc_praying_man_give.png"
	jerArtKidIdle    = jerNPCWall + "npc_wall_kid_idle.png" // #separate idle + talk
	jerArtKidTalk    = jerNPCWall + "npc_wall_kid_talk.png"

	// --- market (souk): full body, SEPARATE idle/talk per the user ---
	jerArtSpiceIdle  = jerNPCMarket + "npc_spice_seller_idle.png"
	jerArtSpiceTalk  = jerNPCMarket + "npc_spice_seller_talk.png"
	jerArtSpiceGive  = jerNPCMarket + "npc_spice_seller_give.png"
	jerArtCoffeeIdle = jerNPCMarket + "npc_coffee_seller_idle.png"
	jerArtCoffeeTalk = jerNPCMarket + "npc_coffee_seller_talk.png"
	jerArtCoffeeGive = jerNPCMarket + "npc_coffee_seller_give.png"

	// Separation-fence prop in the plaza (static overlay; keyed load).
	jerArtFence = "assets/images/locations/jerusalem/props/fence.png"

	// Placeholder fallbacks (existing Paris/camp sheets).
	jerFbkGuard6x2     = "assets/images/locations/paris/npc/outside/npc_security_guard.png"
	jerFbkVendor8x2    = "assets/images/locations/paris/npc/outside/npc_art_vendor.png"
	jerFbkGuideIdle8x2 = "assets/images/locations/paris/npc/outside/npc_french_guide_idle.png"
	jerFbkGuideTalk8x1 = "assets/images/locations/paris/npc/outside/npc_french_guide_talk.png"
	jerFbkKid8x2       = "assets/images/locations/camp/npc/kids/jake/npc_jake_idle.png"
	// A DIFFERENT placeholder for the coffee seller so he doesn't look identical
	// to the spice seller (both used the art-vendor sheet) until his art lands.
	jerFbkPhotog8x2 = "assets/images/locations/paris/npc/outside/npc_press_photographer.png"
)

// loadJerNPCSheet prefers a city sheet at `pref` (cut prefCols×prefRows) and,
// until it lands, falls back to one ROW of an existing Paris/camp placeholder
// sheet. Lets the Jerusalem NPCs use proper SEPARATE full-body idle/talk sheets
// when authored while still showing a placeholder today.
func loadJerNPCSheet(renderer *sdl.Renderer, pref string, prefCols, prefRows int, fbk string, fbkCols, fbkRows, fbkRow int) []npcFrame {
	if _, err := os.Stat(pref); err == nil {
		return loadNPCGrid(renderer, pref, prefCols, prefRows)
	}
	return loadNPCGridRow(renderer, fbk, fbkCols, fbkRows, fbkRow)
}

// registerJerGive loads an optional NPC give one-shot (no-ops if the art is
// absent), so the trade callbacks can play it without a missing-file load.
func registerJerGive(renderer *sdl.Renderer, n *npc, path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	if f := loadNPCGrid(renderer, path, 8, 1); len(f) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["give"] = f
	}
}

// Placeholder palette — warm limestone for the plaza/wall, dim amber for the souk.
var (
	jerPlazaBase  = color.NRGBA{R: 214, G: 182, B: 140, A: 255}
	jerWallBase   = color.NRGBA{R: 224, G: 190, B: 120, A: 255}
	jerMarketBase = color.NRGBA{R: 120, G: 96, B: 68, A: 255}
)

// ---------- NPC dialogs ----------

var shimonIntroDialog = []dialogEntry{
	{speaker: "Shimon", text: "Shalom, friend. First time at the Kotel - the Western Wall?"},
	{speaker: "Pink Panther", text: "I'm chasing a boy's nightmare. A face in old stones. Tunnels."},
	{speaker: "Shimon", text: "Then you've come to the right place. The Wall is just to my right - take the path up."},
	{speaker: "Shimon", text: "But mind your manners. You don't take from the Wall without leaving something behind."},
	{speaker: "Shimon", text: "The market's through the arch to the left, if you need... fortification. The souk coffee is famous."},
}

var shimonWaitDialog = []dialogEntry{
	{speaker: "Shimon", text: "The Wall is up to my right. The market, through the arch on the left. Go on."},
}

var shimonPenDialog = []dialogEntry{
	{speaker: "Shimon", text: "A note for the Wall? Here - take my pen. Everyone deserves a pen for a wish."},
	{speaker: "Pink Panther", text: "Thank you, Shimon."},
	{speaker: "Shimon", text: "Write it true. Then tuck it deep in the stones."},
}

var shimonCoinDialog = []dialogEntry{
	{speaker: "Shimon", text: "You left your note. Good. The Wall always answers - sometimes slowly."},
	{speaker: "Shimon", text: "Here. I found this in the dust by the gate years ago. An old, old coin. Take it - a memory of this place."},
	{speaker: "Pink Panther", text: "That face... that's HIM. That's exactly the face from Jake's dream!"},
	{speaker: "Shimon", text: "Then carry it to him. Tell him the face was never chasing him. It was only remembering."},
}

var shimonDoneDialog = []dialogEntry{
	{speaker: "Shimon", text: "Safe travels, friend. The Wall will be here when you return."},
}

var spiceIntroDialog = []dialogEntry{
	{speaker: "Spice Seller", text: "Za'atar, sumac, CARDAMOM - the finest in the souk! Here, for you, a pinch of cardamom."},
	{speaker: "Pink Panther", text: "For me? What's the catch?"},
	{speaker: "Spice Seller", text: "No catch! But if you want a real Jerusalem coffee, take it to my cousin's stall in the middle. Cardamom makes the coffee."},
}

var spicePostDialog = []dialogEntry{
	{speaker: "Spice Seller", text: "Cardamom in the coffee, my friend. Tell my cousin Shimon sent... no, tell him the spice man sent you."},
}

var coffeeNeedCardamomDialog = []dialogEntry{
	{speaker: "Coffee Seller", text: "Coffee? Of course - but a proper one needs cardamom. Get a pinch from the spice stall and come sit with me."},
}

var coffeeTradeDialog = []dialogEntry{
	{speaker: "Coffee Seller", text: "Cardamom! Perfect. Sit, sit. Let it brew."},
	{speaker: "Coffee Seller", text: "You feel that quiet? Three thousand years of people sitting exactly here. Romans, pilgrims, traders."},
	{speaker: "Coffee Seller", text: "The boy in your story - the face he draws is on an old coin from the tunnels. The Wall keeps such things."},
	{speaker: "Pink Panther", text: "So the nightmare is really a memory."},
	{speaker: "Coffee Seller", text: "Just so. Here - take a cup for the road. And a tip: the ka'ak seller in the plaza trades bread for good coffee."},
}

var coffeePostDialog = []dialogEntry{
	{speaker: "Coffee Seller", text: "Enjoy the coffee. And take some to the bread man in the plaza - he has a sweet tooth for it."},
}

var bagelNeedCoffeeDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Fresh ka'ak! Sesame ka'ak! ...but ahh, what I'd give for a real souk coffee with it."},
}

var bagelTradeDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Is that cardamom coffee I smell? Trade you - a warm ka'ak for that cup!"},
	{speaker: "Pink Panther", text: "Deal."},
	{speaker: "Bagel Seller", text: "Bless you. Take the ka'ak to the old man praying at the Wall - he hasn't eaten all morning, stubborn soul."},
}

var bagelPostDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Go on, take the ka'ak to the old man at the Wall."},
}

var prayingIntroDialog = []dialogEntry{
	{speaker: "Praying Man", text: "(he turns from the stones)  ...Shalom. You stand at the oldest mailbox in the world, you know."},
	{speaker: "Pink Panther", text: "Mailbox?"},
	{speaker: "Praying Man", text: "People write what's in their hearts and tuck it in the cracks. Wishes. Fears. The Wall holds them all."},
	{speaker: "Praying Man", text: "If you carry someone's fear, leave it here. But first - has the bread man been by? I am faint with hunger."},
}

var prayingBagelDialog = []dialogEntry{
	{speaker: "Praying Man", text: "Ah, a warm ka'ak! Bless you, bless you. Now sit a moment."},
	{speaker: "Praying Man", text: "Here is a slip of paper. Write the boy's fear on it - name it - and place it in the Wall."},
	{speaker: "Praying Man", text: "Naming a fear is the first courage. The rest follows."},
	{speaker: "Pink Panther", text: "I'll need something to write with."},
	{speaker: "Praying Man", text: "Shimon by the gate always has a pen. Ask him."},
}

var prayingPostDialog = []dialogEntry{
	{speaker: "Praying Man", text: "Write the boy's fear, then leave it in the stones. Shimon has a pen."},
}

var kidPrepDialog = []dialogEntry{
	{speaker: "Kid", text: "I'm practicing for my bar mitzvah. I have to read in front of EVERYONE next week."},
	{speaker: "Pink Panther", text: "Nervous?"},
	{speaker: "Kid", text: "Terrified! But Saba says - you write the scary thing down, leave it in the Wall, and walk away lighter."},
	{speaker: "Kid", text: "I left mine yesterday. I feel a little braver already. You should try it."},
}

var kidPostDialog = []dialogEntry{
	{speaker: "Kid", text: "Write it down. Leave it in the Wall. Walk away lighter."},
}

var wallCrackBlockedDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Cracks full of folded notes... I should leave one too. I'll need paper and something to write with first."},
}

// ---------- NPC constructors ----------

func newShimon(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtShimon, jerFbkGuard6x2, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtShimon, jerFbkGuard6x2, 6, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 430, W: 120, H: 230},
		name:           "Shimon",
		dialog:         shimonIntroDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.2,
	}
	registerJerGive(renderer, n, jerArtShimonGive)
	return n
}

func newSpiceSeller(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		// SEPARATE idle/talk sheets, FULL BODY (#user 2026-06-21).
		idleGrid:       loadJerNPCSheet(renderer, jerArtSpiceIdle, 8, 1, jerFbkVendor8x2, 8, 2, 0),
		talkGrid:       loadJerNPCSheet(renderer, jerArtSpiceTalk, 8, 1, jerFbkVendor8x2, 8, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 250, W: 140, H: 230},
		name:           "Spice Seller",
		dialog:         spiceIntroDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
	registerJerGive(renderer, n, jerArtSpiceGive)
	return n
}

func newCoffeeSeller(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		// SEPARATE idle/talk sheets, FULL BODY (#user 2026-06-21). Distinct
		// placeholder from the spice seller (was identical art-vendor) until art lands.
		idleGrid:       loadJerNPCSheet(renderer, jerArtCoffeeIdle, 8, 1, jerFbkPhotog8x2, 8, 2, 0),
		talkGrid:       loadJerNPCSheet(renderer, jerArtCoffeeTalk, 8, 1, jerFbkPhotog8x2, 8, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 250, W: 140, H: 230},
		name:           "Coffee Seller",
		dialog:         coffeeNeedCardamomDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
		ppFacePlayer:   true,
	}
	registerJerGive(renderer, n, jerArtCoffeeGive)
	return n
}

func newBagelSeller(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		idleGrid:       loadNPCGridRowPath(renderer, jerArtBagel, jerFbkGuard6x2, 6, 2, 0),
		talkGrid:       loadNPCGridRowPath(renderer, jerArtBagel, jerFbkGuard6x2, 6, 2, 1),
		bounds:         sdl.Rect{X: x, Y: 430, W: 120, H: 230},
		name:           "Bagel Seller",
		dialog:         bagelNeedCoffeeDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
	registerJerGive(renderer, n, jerArtBagelGive)
	return n
}

func newPrayingMan(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		idleGrid:       loadNPCGridPath(renderer, jerArtPrayIdle, jerFbkGuideIdle8x2, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, jerArtPrayTalk, jerFbkGuideTalk8x1, 8, 1),
		bounds:         sdl.Rect{X: x, Y: 470, W: 130, H: 230},
		name:           "Praying Man",
		dialog:         prayingIntroDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.2,
	}
	registerJerGive(renderer, n, jerArtPrayGive)
	return n
}

func newWallKid(renderer *sdl.Renderer, x int32) *npc {
	// SEPARATE idle/talk sheets (#user 2026-06-21).
	return &npc{
		idleGrid:       loadJerNPCSheet(renderer, jerArtKidIdle, 8, 1, jerFbkKid8x2, 8, 2, 0),
		talkGrid:       loadJerNPCSheet(renderer, jerArtKidTalk, 8, 1, jerFbkKid8x2, 8, 2, 0),
		bounds:         sdl.Rect{X: x, Y: 500, W: 100, H: 200},
		name:           "Kid",
		dialog:         kidPrepDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
}

// ---------- Scene builders ----------

func addJerusalemScenes(sm *sceneManager, renderer *sdl.Renderer) {
	// ===== Entrance plaza (hub - PP lands here) =====
	// Shimon at the separation fence (centre-right), bagel seller on the left.
	entrance := &scene{
		name:   "jerusalem_entrance",
		bg:     newPNGBackgroundOr(renderer, jerBgEntrance, jerPlazaBase),
		npcs:   []*npc{newShimon(renderer, 760), newBagelSeller(renderer, 250)},
		spawnX: 640,
		spawnY: 560,
		hotspots: []hotspot{
			{
				// LEFT through the arch to the souk.
				bounds:      sdl.Rect{X: 20, Y: 330, W: 180, H: 260},
				targetScene: "jerusalem_market",
				name:        "To the Market",
				arrow:       arrowLeft,
			},
			{
				// To Shimon's RIGHT, UP the path to the Wall (#24).
				bounds:      sdl.Rect{X: 980, Y: 120, W: 320, H: 300},
				targetScene: "jerusalem_wall",
				name:        "To the Wall",
				arrow:       arrowUp,
			},
		},
		minY: 470,
		maxY: 600,
		// Paris-style flat walk line across the plaza (#23).
		walkSegments: []walkSegment{
			{x1: 150, y1: 560, x2: 1150, y2: 560},
		},
	}
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
	// #22: a CROWD of worshippers at the wall in the mid-distance (multiplied).
	for i := 0; i < 4; i++ {
		entrance.ambientSprites = append(entrance.ambientSprites,
			newAmbientWorshippers(renderer, float64(940+i*70), 470, 0.40+rand.Float64()*0.12))
	}
	// Separation-fence prop Shimon stands by (user 2026-06-21). Static 1-frame
	// overlay; no-ops until props/fence.png lands (§JERUSALEM). Placed at the
	// plaza ground beside Shimon (x≈760) - F3-tune its x/y/scale to the art.
	entrance.ambientSprites = append(entrance.ambientSprites,
		newAmbientSway(renderer, jerArtFence, 1, 700, 600, 1.0, 1.0))
	sm.scenes["jerusalem_entrance"] = entrance

	// ===== Up at the Western Wall =====
	wall := &scene{
		name:           "jerusalem_wall",
		bg:             newPNGBackgroundOr(renderer, jerBgWall, jerWallBase),
		npcs:           []*npc{newPrayingMan(renderer, 470), newWallKid(renderer, 760), newWallKid(renderer, 980)},
		spawnX:         220,
		spawnY:         640,
		characterScale: 0.85,
		hotspots: []hotspot{
			{
				bounds:      sdl.Rect{X: 0, Y: 250, W: 110, H: 420},
				targetScene: "jerusalem_entrance",
				name:        "Back to the Plaza",
				arrow:       arrowLeft,
			},
			{
				// The crack in the Wall where notes are placed (wired in
				// setupJerusalemCallbacks onInteract).
				bounds: sdl.Rect{X: 560, Y: 180, W: 220, H: 320},
				name:   "A crack in the Wall",
			},
		},
		minY: 600,
		maxY: 690,
		walkSegments: []walkSegment{
			{x1: 150, y1: 640, x2: 1150, y2: 640},
		},
	}
	// #22: multiplied worshippers swaying at the foot of the Wall.
	for i := 0; i < 5; i++ {
		wall.ambientSprites = append(wall.ambientSprites,
			newAmbientWorshippers(renderer, float64(150+i*230), 700, 0.6+rand.Float64()*0.15))
	}
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
	// User #25: enter from the far centre, walk down to the centre; the exit
	// back to the plaza is an UP square. Coffee seller at the centre (578,455),
	// spice seller on the left. Walk lines authored as foot-135 (CENTER) of the
	// user's foot coords: bottom foot 733 -> 598, top foot 441 -> 306.
	market := &scene{
		name:           "jerusalem_market",
		bg:             newPNGBackgroundOr(renderer, jerBgMarket, jerMarketBase),
		npcs:           []*npc{newSpiceSeller(renderer, 200), newCoffeeSeller(renderer, 510)},
		spawnX:         760,
		spawnY:         306,
		characterScale: 0.9,
		hotspots: []hotspot{
			{
				// #25: UP square back to the plaza (780,270)-(780,417).
				bounds:      sdl.Rect{X: 700, Y: 270, W: 160, H: 150},
				targetScene: "jerusalem_entrance",
				name:        "Back to the Plaza",
				arrow:       arrowUp,
			},
		},
		minY: 171,
		maxY: 463,
		walkSegments: []walkSegment{
			{x1: 319, y1: 598, x2: 970, y2: 598},
			{x1: 659, y1: 306, x2: 859, y2: 306},
		},
	}
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

// setupJerusalemCallbacks wires the full Jerusalem daisy-chain. Called after
// scenes are built. Linear state lives in captured bools (like the Paris chain).
func (g *Game) setupJerusalemCallbacks() {
	game := g

	held := func(name string) bool {
		return game.inv.heldItem != nil && game.inv.heldItem.name == name
	}
	give := func(id, name string) {
		if item := game.items.createItem(id); item != nil {
			game.inv.addItem(item)
		}
		_ = name
	}

	// ===== Plaza: Shimon (pen + coin) + bagel seller + travel-map return =====
	if entrance, ok := g.sceneMgr.scenes["jerusalem_entrance"]; ok {
		// Travel-map return at the top of the plaza.
		entrance.hotspots = append(entrance.hotspots, hotspot{
			bounds: sdl.Rect{X: 540, Y: 0, W: 300, H: 80},
			name:   "Travel Map",
			arrow:  arrowUp,
			onInteract: func() bool {
				game.openTravelMap("jerusalem_entrance")
				return true
			},
		})
		for _, n := range entrance.npcs {
			switch n.name {
			case "Shimon":
				shimon := n
				gavePen := false
				shimon.onDialogEnd = func() {
					if shimon.dialog == nil {
						return
					}
					// First chat → switch to the short "go on" reminder.
					if !gavePen && !game.vars.GetBool(ScopeGame, VarJerNotePlaced) {
						shimon.dialog = shimonWaitDialog
					}
				}
				shimon.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					// Stage 3: note placed → give the Coin (Jake's anchor).
					if game.vars.GetBool(ScopeGame, VarJerNotePlaced) && !game.inv.hasItem("Coin") {
						return shimonCoinDialog, func() {
							shimon.playOneShotAnim("give", 1.2)
							game.player.playOneShot("receive_item", 1.2, nil)
							give("coin", "Coin")
							shimon.dialog = shimonDoneDialog
							shimon.altDialogFunc = nil
						}, nil
					}
					// Stage 2: PP has the note paper but no pen → give the Pen.
					if !gavePen && game.inv.hasItem("Note Paper") && !game.inv.hasItem("Pen") {
						return shimonPenDialog, func() {
							shimon.playOneShotAnim("give", 1.2)
							game.player.playOneShot("receive_item", 1.2, nil)
							give("pen", "Pen")
							gavePen = true
							shimon.dialog = shimonWaitDialog
						}, nil
					}
					return nil, nil, nil
				}
			case "Bagel Seller":
				bagel := n
				bagel.altDialogRequiresHeld = true
				bagel.altDialogRequiresItem = "Coffee"
				bagel.altDialogStrictMissingHint = bagelNeedCoffeeDialog
				bagel.onDialogEnd = func() { bagel.dialog = bagelNeedCoffeeDialog }
				bagel.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !held("Coffee") || game.inv.hasItem("Bagel") {
						return nil, nil, nil
					}
					return bagelTradeDialog, func() {
						game.inv.giveItemTo("Coffee", "bagel_seller")
						bagel.playOneShotAnim("give", 1.2)
						game.player.playOneShot("receive_item", 1.2, nil)
						give("bagel", "Bagel")
						bagel.dialog = bagelPostDialog
						bagel.altDialogFunc = nil
						bagel.altDialogRequiresHeld = false
						bagel.altDialogRequiresItem = ""
					}, &handOff{item: "Coffee", npcAnim: "receive_item"}
				}
			}
		}
	}

	// ===== Market: spice seller (cardamom) + coffee seller (coffee) =====
	if market, ok := g.sceneMgr.scenes["jerusalem_market"]; ok {
		for _, n := range market.npcs {
			switch n.name {
			case "Spice Seller":
				spice := n
				spice.onDialogEnd = func() {
					if !game.inv.hasItem("Cardamom") && !game.inv.hasItem("Coffee") {
						give("cardamom", "Cardamom")
						spice.playOneShotAnim("give", 1.2)
						game.player.playOneShot("receive_item", 1.0, nil)
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "A pinch of cardamom. Off to the coffee stall, then."},
						})
					}
					spice.dialog = spicePostDialog
				}
			case "Coffee Seller":
				coffee := n
				coffee.altDialogRequiresHeld = true
				coffee.altDialogRequiresItem = "Cardamom"
				coffee.altDialogStrictMissingHint = coffeeNeedCardamomDialog
				coffee.onDialogEnd = func() { coffee.dialog = coffeeNeedCardamomDialog }
				coffee.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !held("Cardamom") || game.inv.hasItem("Coffee") {
						return nil, nil, nil
					}
					return coffeeTradeDialog, func() {
						game.inv.giveItemTo("Cardamom", "coffee_seller")
						coffee.playOneShotAnim("give", 1.2)
						game.player.playOneShot("receive_item", 1.2, nil)
						give("jerusalem_coffee", "Coffee")
						coffee.dialog = coffeePostDialog
						coffee.altDialogFunc = nil
						coffee.altDialogRequiresHeld = false
						coffee.altDialogRequiresItem = ""
					}, &handOff{item: "Cardamom", npcAnim: "receive_item"}
				}
			}
		}
	}

	// ===== Wall: praying man (bagel → note paper), kids, the note-crack =====
	if wall, ok := g.sceneMgr.scenes["jerusalem_wall"]; ok {
		for _, n := range wall.npcs {
			switch n.name {
			case "Praying Man":
				pray := n
				pray.altDialogRequiresHeld = true
				pray.altDialogRequiresItem = "Bagel"
				pray.altDialogStrictMissingHint = prayingIntroDialog
				pray.onDialogEnd = func() {
					if !game.inv.hasItem("Note Paper") {
						pray.dialog = prayingIntroDialog
					} else {
						pray.dialog = prayingPostDialog
					}
				}
				pray.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !held("Bagel") || game.inv.hasItem("Note Paper") {
						return nil, nil, nil
					}
					return prayingBagelDialog, func() {
						game.inv.giveItemTo("Bagel", "praying_man")
						pray.playOneShotAnim("give", 1.2)
						game.player.playOneShot("receive_item", 1.2, nil)
						give("note_paper", "Note Paper")
						pray.dialog = prayingPostDialog
						pray.altDialogFunc = nil
						pray.altDialogRequiresHeld = false
						pray.altDialogRequiresItem = ""
					}, &handOff{item: "Bagel", npcAnim: "receive_item"}
				}
			case "Kid":
				kid := n
				kid.onDialogEnd = func() { kid.dialog = kidPostDialog }
			}
		}

		// The crack-in-the-Wall hotspot: write + place the note once PP has both
		// the Note Paper and the Pen. Sets jer_note_placed (gates the flight home
		// and Shimon's coin). Plays PP write + put one-shots (no-op until art).
		for i := range wall.hotspots {
			if wall.hotspots[i].name != "A crack in the Wall" {
				continue
			}
			wall.hotspots[i].onInteract = func() bool {
				if game.vars.GetBool(ScopeGame, VarJerNotePlaced) {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "My note's tucked safely in the stones. Jake's fear can stay here now."},
					})
					return true
				}
				if !game.inv.hasItem("Note Paper") || !game.inv.hasItem("Pen") {
					game.dialog.startDialog(wallCrackBlockedDialog)
					return true
				}
				// Write, then place.
				game.player.playOneShot("write_note", 1.4, func() {
					game.player.playOneShot("put_note", 1.4, func() {
						game.inv.removeItem("Note Paper")
						game.vars.SetBool(ScopeGame, VarJerNotePlaced, true)
						game.dialog.startDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "There. Jake's fear, named and left in the Wall."},
							{speaker: "Pink Panther", text: "Shimon said the Wall always answers - I should go see him."},
						})
					})
				})
				return true
			}
			break
		}
	}
}

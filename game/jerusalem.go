package game

import (
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
	// Art is organised into wall/ (plaza + Wall NPCs) and market/ (souk NPCs).
	// Shimon's full-body 6x2 sheet has landed; the rest borrow Paris/camp art
	// via the fallbacks below until their sheets are authored (§JERUSALEM).

	jerNPCWall   = "assets/images/locations/jerusalem/npc/wall/"
	jerNPCMarket = "assets/images/locations/jerusalem/npc/market/"

	// --- wall / plaza ---
	jerArtShimon     = jerNPCWall + "npc_shimon.png"      // 6x2 (idle row0, talk row1) - LANDED
	jerArtShimonGive = jerNPCWall + "npc_shimon_give.png" // 8x1 give one-shot (pen / coin)
	// Per-item Shimon gives (#35/#37) - upgrade the two-stage hand-off when they land.
	jerArtShimonGivePen  = jerNPCWall + "npc_shimon_give_pen.png"
	jerArtShimonGiveCoin = jerNPCWall + "npc_shimon_give_coin.png"
	jerArtBagel          = jerNPCWall + "npc_bagel_seller.png" // legacy combined fallback
	jerArtBagelIdle      = jerNPCWall + "npc_bagel_seller_idle.png"
	jerArtBagelTalk      = jerNPCWall + "npc_bagel_seller_talk.png"
	jerArtBagelGive      = jerNPCWall + "npc_bagel_seller_give.png"
	jerArtPrayIdle       = jerNPCWall + "npc_praying_man_idle.png"
	jerArtPrayTalk       = jerNPCWall + "npc_praying_man_talk.png"
	jerArtPrayGive       = jerNPCWall + "npc_praying_man_give.png"
	jerArtPrayGivePaper  = jerNPCWall + "npc_praying_man_give_paper.png" // #34
	// #33 distinct wall-worshipper sway sheets (4-frame each, seen from behind
	// at the Wall). Used by the ambient sway system, NOT as NPC idle/talk.
	// 2026-07-24 (user #32): the user is moving praying_man{,2,3}.png into
	// props/ — resolve via jerWallPrayerSheet() so both homes work.
	jerArtWallPrayer1 = jerNPCWall + "praying_man.png"
	jerArtWallPrayer2 = jerNPCWall + "praying_man2.png"
	jerProps          = "assets/images/locations/jerusalem/props/"
	jerArtKidIdle     = jerNPCWall + "npc_wall_kid_idle.png" // #separate idle + talk
	jerArtKidTalk     = jerNPCWall + "npc_wall_kid_talk.png"

	// --- market (souk): full body, SEPARATE idle/talk per the user ---
	jerArtSpiceIdle  = jerNPCMarket + "npc_spice_seller_idle.png"
	jerArtSpiceTalk  = jerNPCMarket + "npc_spice_seller_talk.png"
	jerArtSpiceGive  = jerNPCMarket + "npc_spice_seller_give.png"
	jerArtCoffeeIdle = jerNPCMarket + "npc_coffee_seller_idle.png"
	jerArtCoffeeTalk = jerNPCMarket + "npc_coffee_seller_talk.png"
	jerArtCoffeeGive = jerNPCMarket + "npc_coffee_seller_give.png"

	// --- market antiques stall (#28): a kid minding the stall for her grandpa,
	// who dozes on a chair beside it. SEPARATE idle/talk; art queued in
	// EXTRA_PROMPTS.md, placeholders until then.
	jerArtAntiqueKidIdle = jerNPCMarket + "kid_antique_idle.png"
	jerArtAntiqueKidTalk = jerNPCMarket + "kid_antique_speak.png"
	jerArtAntiqueKidAlt  = jerNPCMarket + "kid_antique_idle_alter.png"
	// Grandpa's talk sheet preserves the same seated chair orientation.
	jerArtOldManIdle = jerNPCMarket + "grandpa_idle.png"
	jerArtOldManTalk = jerNPCMarket + "grandpa_idle_talk.png"

	// Separation-fence prop in the plaza (static overlay; keyed load).
	jerArtFence = "assets/images/locations/jerusalem/props/fence.png"

	// Placeholder fallbacks (existing Paris/camp sheets).
	jerFbkGuard6x2     = "assets/images/locations/paris/npc/outside/npc_security_guard.png"
	jerFbkVendor8x2    = "assets/images/locations/paris/npc/outside/npc_art_vendor.png"
	jerFbkGuideIdle8x2 = "assets/images/locations/paris/npc/outside/npc_madame_colette_idle.png"
	jerFbkGuideTalk8x1 = "assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png"
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
	registerJerGiveGrid(renderer, n, path, 8, 1)
}

// registerJerGiveGrid is the geometry-aware form used by the 6×1
// blue-background re-rolls.
func registerJerGiveGrid(renderer *sdl.Renderer, n *npc, path string, cols, rows int) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	if f := loadNPCGrid(renderer, path, cols, rows); len(f) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["give"] = f
	}
}

// registerJerGiveNamedGrid is registerJerGiveNamed with a caller grid (the
// blue-bg 6×1 sheets).
func registerJerGiveNamedGrid(renderer *sdl.Renderer, n *npc, key, path string, cols, rows int) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	if f := loadNPCGrid(renderer, path, cols, rows); len(f) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims[key] = f
	}
}

// registerJerGiveNamed registers an optional per-item NPC give one-shot under
// `key` (e.g. "give_paper", "give_pen", "give_coin") so the two-stage hand-off
// upgrades from the generic "give" reach the moment that art lands. No-ops when
// the sheet is absent (2026-06-24 #34/#35/#37).
func registerJerGiveNamed(renderer *sdl.Renderer, n *npc, key, path string) {
	registerJerGiveNamedGrid(renderer, n, key, path, 8, 1)
}

// ---------- NPC dialogs ----------

var shimonIntroDialog = []dialogEntry{
	{speaker: "Shimon", text: "Shalom, friend. First time at the Kotel - the Western Wall?"},
	{speaker: "Pink Panther", text: "I'm chasing a boy's nightmare. A face staring out of the old stones."},
	{speaker: "Shimon", text: "Then you've come to the right place. The Wall is just to my right - take the path up."},
	{speaker: "Shimon", text: "But mind your manners. You don't take from the Wall without leaving something behind."},
	{speaker: "Shimon", text: "The market's through the arch to the left, if you need... fortification. The souk coffee is famous."},
}

var shimonWaitDialog = []dialogEntry{
	{speaker: "Shimon", text: "The Wall is up to my right. The market, through the arch on the left. Go on."},
}

var shimonPenDialog = []dialogEntry{
	// 2026-07-15 (user #26): PP ASKS to borrow the pen first.
	{speaker: "Pink Panther", text: "Shimon - I have paper for a note, but nothing to write with. Could I borrow a pen?"},
	{speaker: "Shimon", text: "A note for the Wall? Then you will need this - my pen. Write what is in your heart."},
	{speaker: "Pink Panther", text: "Thank you, Shimon."},
	{speaker: "Shimon", text: "No thanks needed. Tuck it deep in the stones - the Wall keeps everything."},
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

// 2026-07-18 (user #27): the spice seller no longer gifts on first talk —
// he pitches his wares; the cardamom only comes once the coffee seller has
// sent PP (jer_coffee_asked).
var spiceIntroDialog = []dialogEntry{
	{speaker: "Spice Seller", text: "Za'atar, sumac, CARDAMOM - the finest spices in the souk! Smell that? That is the whole market in one breath."},
	{speaker: "Pink Panther", text: "Impressive stock. I'll keep it in mind."},
	{speaker: "Spice Seller", text: "Do that, friend. Nothing in Jerusalem tastes right without something from my table."},
}

var spicePostDialog = []dialogEntry{
	{speaker: "Spice Seller", text: "Cardamom in the coffee, my friend. Tell my cousin Shimon sent... no, tell him the spice man sent you."},
}

var coffeeNeedCardamomDialog = []dialogEntry{
	// 2026-07-15 (user #20): PP asks for the cup FIRST — the request is his.
	{speaker: "Pink Panther", text: "Shalom! One coffee, please. Your smallest, strongest cup."},
	{speaker: "Coffee Seller", text: "Coffee? Of course - but a proper one needs cardamom. Get a pinch from the spice stall and come sit with me."},
}

var coffeeTradeDialog = []dialogEntry{
	// 2026-07-18 (user #28): PP opens the exchange.
	{speaker: "Pink Panther", text: "Here it is - cardamom, for the coffee."},
	{speaker: "Coffee Seller", text: "Cardamom! Perfect. Sit, sit. Let it brew, and stay a while with us."},
	{speaker: "Pink Panther", text: "Ahh, I'd love to, truly - but I'm in a bit of a hurry. There's a frightened boy back at camp counting on me."},
	{speaker: "Pink Panther", text: "Any chance I could take it to go? One coffee, to travel?"},
	{speaker: "Coffee Seller", text: "A panther in a hurry! Of course, of course. A little cup for ze road, zen - but you must still hear one thing while I pour."},
	{speaker: "Coffee Seller", text: "You feel that quiet? Three thousand years of people sitting exactly here. Romans, pilgrims, traders."},
	{speaker: "Coffee Seller", text: "The boy in your story - the face he draws is on an old coin from the Wall's own stones. The Wall keeps such things."},
	{speaker: "Pink Panther", text: "So the nightmare is really a memory."},
	{speaker: "Coffee Seller", text: "Just so. Here - your cup for the road. And a tip: the ka'ak seller in the plaza trades bread for good coffee."},
}

var coffeePostDialog = []dialogEntry{
	{speaker: "Coffee Seller", text: "Enjoy the coffee. And take some to the bread man in the plaza - he has a sweet tooth for it."},
}

var bagelNeedCoffeeDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Fresh ka'ak! Sesame ka'ak! ...but ahh, what I'd give for a real souk coffee with it."},
}

// 2026-08-01 (user #17): the coffee no longer buys the ka'ak on the spot —
// the seller drinks it and OWES PP A FAVOR. Only after Avi asks for a ka'ak
// does PP come back and call the favor in.
var bagelCoffeeThanksDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Is that cardamom coffee I smell? For me? Friend, you just saved my morning."},
	{speaker: "Pink Panther", text: "Enjoy it while it's hot."},
	{speaker: "Bagel Seller", text: "Ahh... perfect. I owe you one now. If you ever need a favor - anything at all - you come and talk to me."},
}

var bagelFavorWaitDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Best coffee I've had all week. Don't forget - you have a favor waiting with me."},
}

var bagelFavorCallDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Remember when you said I could ask you for a favor?"},
	{speaker: "Bagel Seller", text: "Of course I remember! Name it, friend."},
	{speaker: "Pink Panther", text: "One warm ka'ak. It's for the old man praying at the Wall - he hasn't eaten all morning."},
	{speaker: "Bagel Seller", text: "For Avi? Take the warmest one on the cart. Go, go - before it cools!"},
}

var bagelPostDialog = []dialogEntry{
	{speaker: "Bagel Seller", text: "Go on, take the ka'ak to the old man at the Wall."},
}

var prayingIntroDialog = []dialogEntry{
	// 2026-07-24 (user #25): hungry version — he's been praying all day.
	{speaker: "Avi", text: "Shalom, friend... forgive me if I sway. I have been praying here since sunrise, and my stomach is praying louder than I am."},
	{speaker: "Pink Panther", text: "Sounds like you could use a break. Anything I can do?"},
	{speaker: "Avi", text: "You can help me with exactly one thing: a ka'ak from the bagel man in the plaza. Warm, if he has it."},
	{speaker: "Avi", text: "And while I eat, I will teach you what people leave in the Wall - wishes, fears, all of it. A fair trade, no?"},
}

var prayingBagelDialog = []dialogEntry{
	{speaker: "Avi", text: "Ah, a warm ka'ak! Bless you, bless you. Now sit a moment."},
	{speaker: "Avi", text: "Here is a slip of paper. Write the boy's fear on it - name it - and place it in the Wall."},
	{speaker: "Avi", text: "Naming a fear is the first courage. The rest follows."},
	{speaker: "Pink Panther", text: "I'll need something to write with."},
	{speaker: "Avi", text: "Shimon by the gate always has a pen. Ask him."},
}

var prayingPostDialog = []dialogEntry{
	{speaker: "Avi", text: "Write the boy's fear, then leave it in the stones. Shimon has a pen."},
}

var kidPrepDialog = []dialogEntry{
	{speaker: "Noam", text: "I'm practicing for my bar mitzvah. I have to read in front of EVERYONE next week."},
	{speaker: "Pink Panther", text: "Nervous?"},
	{speaker: "Noam", text: "Terrified! But Saba says - you write the scary thing down, leave it in the Wall, and walk away lighter."},
	{speaker: "Noam", text: "I left mine yesterday. I feel a little braver already. You should try it."},
}

var kidPostDialog = []dialogEntry{
	{speaker: "Noam", text: "Write it down. Leave it in the Wall. Walk away lighter."},
}

var wallCrackBlockedDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Cracks full of folded notes... I should leave one too. I'll need paper and something to write with first."},
}

// #28: the antiques stall - a girl keeping shop for her dozing grandfather.
var antiqueKidDialog = []dialogEntry{
	{speaker: "Antiques Girl", text: "Careful with the lamps, mister! Everything here belongs to the family."},
	{speaker: "Pink Panther", text: "You run this whole stall yourself?"},
	{speaker: "Antiques Girl", text: "I mind it while Saba rests. I'm saving it all for him - it's his life's work."},
	{speaker: "Antiques Girl", text: "That's him on the chair. Don't wake him - he guards the oldest pieces even in his sleep."},
}

var oldManDialog = []dialogEntry{
	{speaker: "Old Antiques Man", text: "*snoozes, a brass coin turning slowly between his fingers*"},
	{speaker: "Old Antiques Man", text: "Hm? ...Ah. Every old thing here has a face, and every face has a story. My granddaughter knows them all now."},
	{speaker: "Pink Panther", text: "She's keeping it safe for you."},
	{speaker: "Old Antiques Man", text: "She is. That is how a thing outlives its maker. *drifts back to sleep*"},
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
	// 2026-07-24 (user #29): the generic multi-item npc_shimon_give.png is
	// RETIRED — both his beats name an explicit per-item give (give_pen /
	// give_coin), and the generic reach was leaking in as the pen-take
	// fallback ("he gives some items then the coin").
	registerJerGiveNamed(renderer, n, "give_pen", jerArtShimonGivePen)
	registerJerGiveNamed(renderer, n, "give_coin", jerArtShimonGiveCoin)
	// 2026-08-01 (user #18): Shimon visibly TAKES the pen back once
	// §SHIMON-RECEIVE-PEN lands (optional; the coin trade skips the take
	// until then — see the conditional skipNPCTake at the trade).
	registerJerGiveNamedGrid(renderer, n, "receive_pen", jerNPCWall+"npc_shimon_receive_pen.png", 6, 1)
	return n
}

func newSpiceSeller(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		// SEPARATE idle/talk sheets, FULL BODY (#user 2026-06-21).
		// §SPICE-SIDE (landed 2026-07-14): when talking he turns RIGHT toward
		// the market centre (stall table baked into the sheet, same as idle).
		idleGrid: loadJerNPCSheet(renderer, jerArtSpiceIdle, 8, 1, jerFbkVendor8x2, 8, 2, 0),
		// 2026-07-15 (user): the REDESIGNED talk sheet wins; the old side-talk
		// (previous design) stays as fallback only.
		// 2026-07-23 (user #25): the sheet is 8 frames (idle + both checker
		// manifests agree) but was read 6×1 — each cell spanned 1.33 real
		// frames, the "frames swiping" the user saw.
		talkGrid: loadJerNPCSheet(renderer, firstExisting(jerArtSpiceTalk, jerNPCMarket+"npc_spice_seller_talk_side.png"), 8, 1, jerFbkVendor8x2, 8, 2, 1),
		// 2026-06-24 (#26): feet on the market floor (~598) at x≈319, not floating
		// mid-scene. Y = 598 - H. (Leg "jitter" between idle/talk is a sheet
		// baseline issue queued for sprite-check.)
		bounds:         sdl.Rect{X: x, Y: 368, W: 140, H: 230},
		name:           "Spice Seller",
		dialog:         spiceIntroDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
		// D5: spice seller is on the left wall — keep him facing right (toward
		// the market centre) rather than flipping to face PP every time.
		fixedFacing: true,
		flipped:     false,
	}
	registerJerGiveGrid(renderer, n, jerArtSpiceGive, 6, 1)
	// #15: pin talk + the give one-shot to the idle's scale so the spice seller
	// doesn't shrink when he talks or hands over the cardamom.
	n.anchorRefH = maxOpaqueH(n.idleGrid)
	n.anchorRefHOneShots = true
	// #14: 3D depth — the spice seller is at the FRONT of the souk, so render him
	// a touch bigger than the farther coffee seller (kept modest, not huge).
	n.extraScale = 1.08
	return n
}

func newCoffeeSeller(renderer *sdl.Renderer, x int32) *npc {
	// The talk sheet (npc_coffee_seller_talk.png) is half-body (bust only).
	// Load the full-body idle for both until a full-body talk sheet replaces it.
	idle := loadJerNPCSheet(renderer, jerArtCoffeeIdle, 8, 1, jerFbkPhotog8x2, 8, 2, 0)
	n := &npc{
		idleGrid:       idle,
		talkGrid:       idle,
		bounds:         sdl.Rect{X: x, Y: 250, W: 140, H: 230},
		name:           "Coffee Seller",
		dialog:         coffeeNeedCardamomDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
		// 2026-07-15 (user #20): PP stands to his RIGHT and talks facing
		// LEFT (was ppFacePlayer/front). Seller centre ≈ x+70 → mark at +180.
		approachXOverride: x + 180,
	}
	// §COFFEE-GIVE (landed 2026-07-14): the dedicated pour-and-offer sheet;
	// the older generic give stays as fallback.
	// 2026-08-07 #23: tol 24 — the near-white bg's AA halos bridged figure 8's
	// extended cup-arm into figure 7's zone under the tol-8 key, so the
	// waist-splitter cut through the offered cup. Registered under BOTH the
	// generic "give" and the derived give_jerusalem_coffee key that
	// stageReturn tries first, so the hand-back resolves without fallback.
	if p := firstExisting(jerNPCMarket+"npc_coffee_seller_give_coffee.png", jerArtCoffeeGive); p != "" {
		if f := loadNPCGridTol(renderer, p, 8, 1, 24); len(f) > 0 {
			n.oneShotAnims = map[string][]npcFrame{
				"give":                  f,
				"give_jerusalem_coffee": f,
			}
		}
	}
	// #15/#16: pin the give one-shot to the idle's scale so the coffee seller
	// doesn't shrink when he pours/hands the coffee.
	n.anchorRefH = maxOpaqueH(n.idleGrid)
	n.anchorRefHOneShots = true
	// #14: 3D depth — the coffee seller sits FARTHER back than the spice seller,
	// so render him smaller for the perspective.
	n.extraScale = 0.78
	return n
}

func newBagelSeller(renderer *sdl.Renderer, x int32) *npc {
	idle := loadNPCGridRowPath(renderer, jerArtBagel, jerFbkGuard6x2, 6, 2, 0)
	if _, err := os.Stat(jerArtBagelIdle); err == nil {
		idle = loadNPCGrid(renderer, jerArtBagelIdle, 6, 1)
	}
	talk := loadNPCGridRowPath(renderer, jerArtBagel, jerFbkGuard6x2, 6, 2, 1)
	if _, err := os.Stat(jerArtBagelTalk); err == nil {
		talk = loadNPCGrid(renderer, jerArtBagelTalk, 6, 1)
	}
	n := &npc{
		idleGrid:       idle,
		talkGrid:       talk,
		bounds:         sdl.Rect{X: x, Y: 430, W: 120, H: 230},
		name:           "Bagel Seller",
		dialog:         bagelNeedCoffeeDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
	registerJerGiveGrid(renderer, n, jerArtBagelGive, 6, 1)
	// 2026-07-29 (user): the redesigned give sheet reaches AWAY from PP —
	// mirror it (same fix as Avi's receive_bagel).
	n.oneShotFlip = map[string]bool{"give": true}
	// 2026-07-24 (user #24): he visibly TAKES the finjan on the trade once
	// §BAGEL-RECEIVE-COFFEE lands (optional; the handOff skips the take
	// until then).
	registerJerGiveNamedGrid(renderer, n, "receive_coffee", jerNPCWall+"npc_bagel_seller_receive_coffee.png", 6, 1)
	return n
}

func newPrayingMan(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		idleGrid:       loadNPCGridPath(renderer, jerArtPrayIdle, jerFbkGuideIdle8x2, 8, 2),
		talkGrid:       loadNPCGridPath(renderer, jerArtPrayTalk, jerFbkGuideTalk8x1, 8, 1),
		bounds:         sdl.Rect{X: x, Y: 470, W: 130, H: 230},
		name:           "Avi",
		dialog:         prayingIntroDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.2,
	}
	registerJerGive(renderer, n, jerArtPrayGive)
	// 2026-08-07 #1: the give-paper sheet holds SEVEN figures (gaps 60-75px),
	// not 8 — the 8×1 load waist-split a real figure into a 34×103 runt that
	// dropMalformedFrames warned about and deleted every boot.
	registerJerGiveNamedGrid(renderer, n, "give_paper", jerArtPrayGivePaper, 7, 1)
	registerJerGiveNamedGrid(renderer, n, "receive_bagel", jerNPCWall+"npc_praying_man_receive_bagel.png", 6, 1) // §AVI-RECEIVE-BAGEL
	// 2026-07-24 (user #27): the receive sheet reaches away from PP — mirror it.
	n.oneShotFlip = map[string]bool{"receive_bagel": true}
	return n
}

func newWallKid(renderer *sdl.Renderer, x int32) *npc {
	// SEPARATE idle/talk sheets (#user 2026-06-21).
	return &npc{
		idleGrid:       loadJerNPCSheet(renderer, jerArtKidIdle, 8, 1, jerFbkKid8x2, 8, 2, 0),
		talkGrid:       loadJerNPCSheet(renderer, jerArtKidTalk, 8, 1, jerFbkKid8x2, 8, 2, 0),
		bounds:         sdl.Rect{X: x, Y: 525, W: 100, H: 175}, // 2026-07-23 (user #27): H150 read tiny next to Avi — up to 175 (~76% of Avi), feet kept at 700
		name:           "Noam",
		dialog:         kidPrepDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
}

// newAntiquesKid (#28): a girl minding the family antiques stall in the souk.
func newAntiquesKid(renderer *sdl.Renderer, x int32) *npc {
	n := &npc{
		idleGrid:       loadJerNPCSheet(renderer, jerArtAntiqueKidIdle, 6, 1, jerFbkKid8x2, 8, 2, 0),
		talkGrid:       loadJerNPCSheet(renderer, jerArtAntiqueKidTalk, 6, 1, jerFbkKid8x2, 8, 2, 0),
		bounds:         sdl.Rect{X: x, Y: 408, W: 100, H: 190},
		name:           "Antiques Girl",
		dialog:         antiqueKidDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
	// Optional second idle (dusting variation) — plays as the periodic alt-idle.
	// The v2 sheet matches the redesigned 6×1 idle and omits the old table.
	if f := loadNPCGrid(renderer, jerNPCMarket+"kid_antique_idle_alter_v2.png", 6, 1); len(f) > 0 {
		n.altIdleGrid = f
		n.altIdleAfterSec = 5.0
	}
	return n
}

// newAntiquesOldMan (#28): the grandpa dozing on a chair by the antiques stall.
func newAntiquesOldMan(renderer *sdl.Renderer, x int32) *npc {
	return &npc{
		idleGrid: loadJerNPCSheet(renderer, jerArtOldManIdle, 8, 1, jerFbkGuard6x2, 6, 2, 0),
		talkGrid: loadJerNPCSheet(renderer, jerArtOldManTalk, 8, 1, jerFbkGuard6x2, 6, 2, 1),
		// Seated on a chair, so he reads shorter than a standing vendor.
		bounds:         sdl.Rect{X: x, Y: 426, W: 120, H: 160}, // 2026-07-18 (user #29): raised 12px
		name:           "Old Antiques Man",
		dialog:         oldManDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.16,
		// 2026-07-15 (user #22): mirroring him mirrored the CHAIR too — hold
		// the authored facing; a dedicated talk sheet is queued (§GRANDPA-TALK).
		fixedFacing: true,
		flipped:     false,
	}
}

// ---------- Scene decorators ----------
// Static scene data (background, spawn, walk-segments, hotspots, NPCs) now
// lives in assets/data/scenes/jerusalem_*.json. These functions add the
// procedural layer (particles, glows, ambient sprites) that JSON can't express.

func decorateJerusalemEntrance(s *scene, renderer *sdl.Renderer) {
	// Warm dusty-light haze rising off the stones.
	for i := 0; i < 10; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 500,
			vx:    (rand.Float64() - 0.5) * 4,
			vy:    -rand.Float64()*1.2 - 0.2,
			alpha: uint8(rand.Intn(12) + 5),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	s.glows = []glowEffect{
		{x: 0, y: 0, w: 1400, h: 280, r: 255, g: 240, b: 200, alpha: 10, pulse: 0.25},
	}
	// Distant worshippers at the wall in the mid-ground.
	// #13 (2026-06-30): moved DOWN (470 -> 525) so the praying group sits lower
	// against the wall in the BG, as asked.
	// 2026-07-23 (user #24): both sheets are now figures-only flat-blue strips.
	// Alternate the distinct four-person groups for ABAB crowd variety.
	crowds := []string{
		firstExisting(jerProps+"people_pray.png", jerArtWallPrayer1),
		firstExisting(jerProps+"people_pray2.png", jerNPCWall+"people_pray2.png"),
	}
	// 2026-07-24 (user #22): SIX smaller groups, ABABAB, along the same line
	// (was 4 at 0.40-0.52). Spread starts at 880 so the crowd stays right of
	// Shimon (x760-880) and inside the frame.
	for i := 0; i < 6; i++ {
		s.ambientSprites = append(s.ambientSprites,
			newAmbientWorshippersSheet(renderer, crowds[i%len(crowds)], float64(880+i*62), 585, 0.30+rand.Float64()*0.09))
	}
	// Separation fence: three globally-keyed sections end-to-end across the
	// foreground. The flat blue key is absent from the grey fence, so enclosed
	// background between its bars is removed safely.
	// #12 (2026-06-30): smaller (0.42 -> 0.32). The white BETWEEN the iron bars is
	// enclosed background the old edge key could not reach.
	// 2026-07-15 (user #19): smaller still (0.20) and packed into a tight
	// line on the left/centre — the RIGHT side stays open for the wall path.
	for _, fx := range []float64{160, 420, 680} {
		s.ambientSprites = append(s.ambientSprites,
			newAmbientPropKeyed(renderer, jerArtFence, fx, 620, 0.20))
	}
}

func decorateJerusalemWall(s *scene, renderer *sdl.Renderer) {
	// Warm golden glow across the top of the Wall.
	s.glows = []glowEffect{
		{x: 0, y: 0, w: engine.ScreenWidth, h: 300, r: 255, g: 235, b: 180, alpha: 10, pulse: 0.2},
	}
	// Dust motes drifting upward in the sacred light.
	for i := 0; i < 6; i++ {
		s.particles = append(s.particles, particle{
			x:     rand.Float64() * float64(engine.ScreenWidth),
			y:     rand.Float64() * 400,
			vx:    (rand.Float64() - 0.5) * 3,
			vy:    -rand.Float64()*0.9 - 0.1,
			alpha: uint8(rand.Intn(10) + 4),
			size:  int32(rand.Intn(2) + 1),
		})
	}
	// Two distinct worshipper figures at the foot of the Wall (D11: was 4).
	// 2026-07-15 (user #24): the two sheets read as the SAME man — the right
	// worshipper now uses the distinct praying_man3 sway (4×1, from-behind).
	// 2026-07-18 (user #31): the from-behind praying_man3 sheet carries its
	// content higher in-frame, so it needs a LOWER anchor than the left man.
	// 2026-07-23 (user #27): the two sway figures rendered ~250-275px tall —
	// bigger than Avi (~196px) and even PP. Scaled down to Avi's band so the
	// wall group reads as one crowd; feet stay planted (sway anchors at y).
	// 2026-07-24 (user #32): sheets moving to props/ — props path preferred,
	// old wall path as fallback.
	wallPrayerSheets := []string{
		firstExisting(jerProps+"praying_man.png", jerArtWallPrayer1),
		firstExisting(jerProps+"praying_man3.png", jerNPCWall+"praying_man3.png"),
	}
	// 2026-07-24 (user pre-PR check): measured content ~525px/frame → both
	// at 0.37 ≈ Avi's ~195px; FEET ON THE STANDING LINE (Avi/Noam foot 700,
	// the left one floated 40px above it); the right one moved off Noam
	// (x880) to the open Wall face at x1100 — everyone prays INTO the Wall
	// on the same line: worshipper 280 · Avi 470 · Noam 880 · worshipper 1100.
	for i, sp := range []struct{ x, y, scale float64 }{{280, 700, 0.37}, {1100, 705, 0.37}} {
		s.ambientSprites = append(s.ambientSprites,
			newAmbientSway(renderer, wallPrayerSheets[i%2], 4, sp.x, sp.y, sp.scale, 0.5))
	}
}

func decorateJerusalemMarket(s *scene) {
	// Warm lantern glows: left stall, right stall, centre shaft of light.
	s.glows = []glowEffect{
		{x: 100, y: 50, w: 300, h: 250, r: 255, g: 190, b: 110, alpha: 12, pulse: 1.4},
		{x: 980, y: 50, w: 300, h: 250, r: 255, g: 180, b: 100, alpha: 10, pulse: 1.0},
		{x: 560, y: 120, w: 256, h: 300, r: 255, g: 240, b: 200, alpha: 8, pulse: 0.3},
	}
	// Spice dust drifting through the souk.
	for i := 0; i < 8; i++ {
		s.particles = append(s.particles, particle{
			x:     200 + rand.Float64()*1000,
			y:     rand.Float64() * 450,
			vx:    (rand.Float64() - 0.5) * 2,
			vy:    -rand.Float64()*0.7 - 0.1,
			alpha: uint8(rand.Intn(12) + 6),
			size:  int32(rand.Intn(2) + 1),
		})
	}
}

// addJerusalemScenes is called from newSceneManager after the JSON scenes have
// been loaded and registered. It applies the procedural decorators; the static
// data (spawn, walk-segs, hotspots, NPCs) now lives in the JSON files.
func addJerusalemScenes(sm *sceneManager, renderer *sdl.Renderer) {
	if s, ok := sm.scenes["jerusalem_entrance"]; ok {
		decorateJerusalemEntrance(s, renderer)
	}
	if s, ok := sm.scenes["jerusalem_wall"]; ok {
		decorateJerusalemWall(s, renderer)
	}
	if s, ok := sm.scenes["jerusalem_market"]; ok {
		decorateJerusalemMarket(s)
	}
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
		// D9: the path UP to the Wall reads as walking into the distance, not
		// marching off the top of the screen. Recede (shrink + drift up) then
		// transition, same as the market→plaza step-out.
		for i := range entrance.hotspots {
			if entrance.hotspots[i].targetScene != "jerusalem_wall" {
				continue
			}
			entrance.hotspots[i].onInteract = func() bool {
				game.player.playRecede(1.0, 0.5, 80, func() {
					game.sceneMgr.transitionTo("jerusalem_wall", game.player)
				})
				return true
			}
			break
		}
		for _, n := range entrance.npcs {
			switch n.name {
			case "Shimon":
				shimon := n
				penGiven := func() bool { return game.vars.GetBool(ScopeGame, VarJerPenGiven) }
				penReturnArmed := false
				shimon.onDialogEnd = func() {
					if shimon.dialog == nil {
						return
					}
					// First chat → switch to the short "go on" reminder.
					if !penGiven() && !game.vars.GetBool(ScopeGame, VarJerNotePlaced) {
						shimon.dialog = shimonWaitDialog
					}
				}
				// 2026-07-15 (user #30 REGRESSION FIX): the old coinStarted
				// latch fired on ui.go's HOVER probe (it calls altDialogFunc
				// every frame) and permanently blocked the pen→coin trade.
				// Selectors must be SIDE-EFFECT-FREE; re-click replay (#29) is
				// prevented by the in-flight guard below instead.
				shimon.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					// 2026-07-18 (user #35): the coin stage needs the pen
					// deliberately CARRIED back, not just sitting in the bag.
					if game.vars.GetBool(ScopeGame, VarJerNotePlaced) && !penReturnArmed {
						penReturnArmed = true
						shimon.altDialogRequiresHeld = true
						shimon.altDialogRequiresItem = "Pen"
					}
					// In-flight guard: while a dialog/sequence/one-shot is
					// already playing, don't start another trade beat.
					if game.dialog.active || game.seqPlayer.IsPlaying() || game.player.activeOneShot != "" {
						return nil, nil, nil
					}
					// Guard: coin already given → nothing more to do.
					if game.vars.GetBool(ScopeGame, VarJerNotePlaced) && penGiven() && !game.inv.hasItem("Pen") && !game.inv.hasItem("Coin") {
						return nil, nil, nil
					}
					// Stage 3 (user #30): PP GIVES THE PEN BACK — holding it is
					// the deliberate act that earns the Coin.
					if game.vars.GetBool(ScopeGame, VarJerNotePlaced) && game.inv.heldItem != nil && game.inv.heldItem.name == "Pen" {
						return shimonCoinDialog, func() {
							game.inv.removeItem("Pen")
							give("coin", "Coin")
							shimon.dialog = shimonDoneDialog
							shimon.altDialogFunc = nil
							// 2026-07-24 (user #29): skipNPCTake — the pen-take
							// stage fell back to the generic multi-item give
							// sheet before give_coin played ("he gives some
							// items then the coin"). PP's own give_pen beat
							// carries the hand-over; Shimon only gives the coin.
							// 2026-08-01 (user #18): once §SHIMON-RECEIVE-PEN
							// lands he takes the pen for real; PP's give_pen now
							// aliases the pencil give so the hand-back animates.
						}, &handOff{item: "Pen", returnItem: "Coin", npcGiveAnim: "give_coin",
							skipNPCTake: !shimon.hasOneShotAnim("receive_pen"), dialogFirst: true}
					}
					// Stage 2: PP has the note paper but no pen → give the Pen (once).
					if !penGiven() && game.inv.hasItem("Note Paper") && !game.inv.hasItem("Pen") {
						return shimonPenDialog, func() {
							give("pen", "Pen")
							game.vars.SetBool(ScopeGame, VarJerPenGiven, true)
							shimon.dialog = shimonWaitDialog
						}, &handOff{returnItem: "Pen", npcGiveAnim: "give_pen", dialogFirst: true}
					}
					return nil, nil, nil
				}
			case "Bagel Seller":
				// 2026-08-01 (user #17): TWO-stage favor rework. The coffee no
				// longer buys the ka'ak on the spot — the seller drinks it and
				// owes PP a favor; only after Avi ASKS for a ka'ak does PP come
				// back and call the favor in. No requires-fields: the branching
				// below runs on every click (the Camille multi-branch pattern).
				bagel := n
				bagel.onDialogEnd = func() {
					if game.vars.GetBool(ScopeGame, "jer_bagel_favor_owed") &&
						!game.vars.GetBool(ScopeGame, "jer_bagel_given") {
						bagel.dialog = bagelFavorWaitDialog
					}
				}
				bagel.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					favorOwed := game.vars.GetBool(ScopeGame, "jer_bagel_favor_owed")
					// Stage 1: PP hands the coffee over → the favor is banked.
					if held("Coffee") && !favorOwed {
						return bagelCoffeeThanksDialog, func() {
								game.inv.giveItemTo("Coffee", "bagel_seller")
								game.vars.SetBool(ScopeGame, "jer_bagel_favor_owed", true)
								bagel.dialog = bagelFavorWaitDialog
								// 2026-07-24 (user #24): once §BAGEL-RECEIVE-COFFEE
								// lands he visibly TAKES the finjan; until then the
								// take stays skipped so his give sheet can't double
								// as it (#25).
							}, &handOff{item: "Coffee", npcAnim: "receive_coffee",
								skipNPCTake: !bagel.hasOneShotAnim("receive_coffee")}
					}
					// Stage 2: Avi has asked for a ka'ak → PP calls the favor in
					// on a plain click (no item needed).
					if favorOwed && game.vars.GetBool(ScopeGame, "jer_avi_asked") &&
						!game.vars.GetBool(ScopeGame, "jer_bagel_given") &&
						!game.inv.hasItem("Bagel") {
						return bagelFavorCallDialog, func() {
							give("bagel", "Bagel")
							game.vars.SetBool(ScopeGame, "jer_bagel_given", true)
							bagel.dialog = bagelPostDialog
							bagel.altDialogFunc = nil
						}, &handOff{returnItem: "Bagel", dialogFirst: true}
					}
					return nil, nil, nil
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
					// user #27: the cardamom is only given once the coffee seller
					// has ASKED for it (jer_coffee_asked) — before that he just
					// pitches; after the gift his dialog moves on.
					if game.vars.GetBool(ScopeGame, "jer_coffee_asked") &&
						!game.inv.hasItem("Cardamom") && !game.inv.hasItem("Coffee") {
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Spice Seller", text: "Ah - so my cousin sent you! Coffee without cardamom is just sad water. Here, a pinch, on the house."},
							{speaker: "Pink Panther", text: "A pinch of cardamom. Off to the coffee stall, then."},
						})
						give("cardamom", "Cardamom")
						spice.playOneShotAnimThen(spice.giveAnimOr("give_cardamom"), 1.2, func() {
							game.player.playReceive("cardamom", false, 1.0, nil)
						})
						spice.dialog = spicePostDialog
					}
				}
			case "Coffee Seller":
				coffee := n
				coffee.altDialogRequiresHeld = true
				coffee.altDialogRequiresItem = "Cardamom"
				// 2026-07-23 (user #26, STORY SOFT-LOCK): no strict-missing-hint
				// here. The hint branch in startNPCDialog returns WITHOUT running
				// onDialogEnd, so `jer_coffee_asked` could never flip and the
				// spice seller never handed the cardamom — the chain dead-ended
				// on every save. Without the hint, a no-cardamom click falls
				// through to the base dialog below (same need-cardamom text)
				// and onDialogEnd runs.
				coffee.onDialogEnd = func() {
					// user #27: his ask is what unlocks the spice seller's pinch.
					game.vars.SetBool(ScopeGame, "jer_coffee_asked", true)
					coffee.dialog = coffeeNeedCardamomDialog
				}
				coffee.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !held("Cardamom") || game.inv.hasItem("Coffee") {
						return nil, nil, nil
					}
					return coffeeTradeDialog, func() {
						game.inv.giveItemTo("Cardamom", "coffee_seller")
						give("jerusalem_coffee", "Coffee")
						coffee.dialog = coffeePostDialog
						coffee.altDialogFunc = nil
						coffee.altDialogRequiresHeld = false
						coffee.altDialogRequiresItem = ""
					}, &handOff{item: "Cardamom", returnItem: "Coffee", skipNPCTake: true} // user #23: same single-give-sheet dedup
				}
			}
		}
		// 2026-06-24 (#31): leaving the souk back to the plaza should feel like
		// stepping OUT of an enclosed space - PP recedes/shrinks up through the
		// exit instead of just walking straight up. Mirrors the cabin-door recede.
		for i := range market.hotspots {
			if market.hotspots[i].targetScene != "jerusalem_entrance" {
				continue
			}
			market.hotspots[i].onInteract = func() bool {
				game.player.playRecede(0.9, 0.45, 70, func() {
					game.sceneMgr.transitionTo("jerusalem_entrance", game.player)
				})
				return true
			}
			break
		}
	}

	// ===== Wall: praying man (bagel → note paper), kids, the note-crack =====
	if wall, ok := g.sceneMgr.scenes["jerusalem_wall"]; ok {
		for _, n := range wall.npcs {
			switch n.name {
			case "Avi":
				pray := n
				pray.altDialogRequiresHeld = true
				pray.altDialogRequiresItem = "Bagel"
				// 2026-08-01 (user #17): NO strict-missing-hint here — that
				// branch returns without running onDialogEnd (the #26 coffee
				// soft-lock), so `jer_avi_asked` could never flip. A no-bagel
				// click falls through to his base dialog (the same ask text)
				// and onDialogEnd runs.
				pray.onDialogEnd = func() {
					// 2026-08-01 (user #17): his ka'ak ask is what lets PP call
					// in the bagel seller's favor (the bagel is no longer handed
					// out with the coffee trade).
					game.vars.SetBool(ScopeGame, "jer_avi_asked", true)
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
						// 2026-07-18 (user #33): the paper comes AFTER the talk —
						// he receives the bagel, speaks, THEN hands the slip.
						pray.playOneShotAnimThen(pray.giveAnimOr("give_paper"), 1.3, func() {
							game.player.playReceive("paper", false, 1.2, nil)
							give("note_paper", "Note Paper")
						})
						pray.dialog = prayingPostDialog
						pray.altDialogFunc = nil
						pray.altDialogRequiresHeld = false
						pray.altDialogRequiresItem = ""
					}, &handOff{item: "Bagel", npcAnim: "receive_bagel"}
				}
			case "Noam":
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
				// 2026-06-24 (#36): walk PP up to the foot of the Wall directly
				// under the crack FIRST, so the write/put-note one-shots play on
				// the same line as the stones. Crack centre ≈ x670; wall line y660.
				game.player.walkToAndDo(670, 660, func() {
					// Write, then place.
					game.player.playOneShot("write_note", 1.4, func() {
						game.player.playOneShot("put_note", 1.4, func() {
							game.inv.removeItem("Note Paper")
							game.vars.SetBool(ScopeGame, VarJerNotePlaced, true)
							// 2026-07-15 (user #31): once the note is in the
							// stones, Avi's chat closes the beat instead of
							// re-teaching the custom.
							if wall, ok := game.sceneMgr.scenes["jerusalem_wall"]; ok {
								for _, wn := range wall.npcs {
									if wn.name == "Avi" {
										wn.dialog = []dialogEntry{
											{speaker: "Avi", text: "I saw. The Wall holds it now - the fear stays HERE, not with the boy."},
											{speaker: "Avi", text: "Go home to him, friend. And walk lightly - you leave lighter than you came."},
										}
										break
									}
								}
							}
							game.dialog.startDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "There. Jake's fear, named and left in the Wall."},
								{speaker: "Pink Panther", text: "Shimon said the Wall always answers - I should go see him."},
							})
						})
					})
				})
				return true
			}
			break
		}
	}
}

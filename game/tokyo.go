package game

import (
	"math/rand"
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Japan / Kyoto chapter: Lily's arc. Art lives under
// assets/images/locations/japan/. Three scenes, left -> right:
//
//   - tokyo_torii  (torii gate corridor): PP lands here. Right -> ramen street.
//   - tokyo_street (ramen store; the tree drops leaves): Hiro the cook + Gary
//     the tourist. Left -> torii, right -> the flower grove.
//   - tokyo_temple (flower store by the forest): Oba-chan presses the sakura
//     (the anchor) + Kiku the dresser spins PP into a kimono for a gag. Left ->
//     the ramen street.
//
// Anchor item "Pressed Sakura" -> PP carries it home and heals Lily at the lake.
//
// The Japan art has been re-saved under several different filenames during
// authoring, so every asset is resolved from a CANDIDATE LIST (firstExisting):
// whichever name is on disk wins, and a missing asset degrades gracefully
// (flat-colour BG / placeholder NPC). Scene keys stay "tokyo_*".

const (
	jpNPCDir  = "assets/images/locations/japan/npc/"
	jpBGDir   = "assets/images/locations/japan/background/"
	jpPropDir = "assets/images/locations/japan/props/"

	jpLeafFall = jpPropDir + "leaf_fall.png" // §JP-LEAVES (pending)

	// Generic fallback row if an NPC's idle sheet is missing entirely.
	jpFbkVendor8x2 = "assets/images/locations/paris/npc/outside/npc_art_vendor.png"
)

// firstExisting returns the first candidate path that exists on disk, or "" if
// none do (callers then fall back to a placeholder).
func firstExisting(paths ...string) string {
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// loadJapanNPC loads idle + talk as SEPARATE 8x1 sheets (user rule 2026-06-24:
// every NPC gets its own idle and talk sheet). Each is resolved from a
// candidate list; talk falls back to idle until its sheet lands; idle falls
// back to a generic vendor row if even idle is missing.
func loadJapanNPC(renderer *sdl.Renderer, idleCands, talkCands []string) (idle, talk []npcFrame) {
	if p := firstExisting(idleCands...); p != "" {
		idle = loadNPCGrid(renderer, p, 8, 1)
	}
	if len(idle) == 0 {
		idle = loadNPCGridRow(renderer, jpFbkVendor8x2, 8, 2, 0)
	}
	if p := firstExisting(talkCands...); p != "" {
		talk = loadNPCGrid(renderer, p, 8, 1)
	}
	if len(talk) == 0 {
		talk = idle
	}
	return
}

// ---------- Dialogs ----------

// --- Gary (torii): overjoyed to finally be in Kyoto; his ramen tip OPENS the
// stall down the street (the book-upside-down gag plays mid-chat). ---
var garyTokyoDialog = []dialogEntry{
	{speaker: "Gary", text: "PINK PANTHER! Can you BELIEVE it - KYOTO! I have dreamed of zis since I was a boy. Ze temples, ze gardens, ze cherry blossoms..."},
	{speaker: "Gary", text: "I read EVERYTHING about zis place - every shrine, every festival, every bowl of noodles. It is ALL here in my guidebook!"},
	{speaker: "Pink Panther", text: "...Gary. You're holding the book upside down."},
	{speaker: "Gary", text: "WHAT? ...Oh. Oh my. ...Ahh. MUCH better. Now it makes far more sense!"},
	{speaker: "Gary", text: "Right - rule one of any guidebook: you MUST taste ze ramen at ze little stall down ze street. Go, go - tell zem Gary sent you, zey'll open right up!"},
	{speaker: "Pink Panther", text: "Ramen it is. Thanks, Gary."},
}

var garyTokyoPostDialog = []dialogEntry{
	{speaker: "Gary", text: "(flips ze book again) Now it says Kyoto is in PERU! Remarkable little book."},
}

// --- Hiro (street): OPEN for business (Gary's tip), but the SACRED hearth for
// the blessed offering bowl still needs his crow-stolen fire-striker ---
var hiroRamenDialog = []dialogEntry{
	{speaker: "Hiro", text: "Irasshaimase, panther-san! You have ze look of a hungry traveller - but my stall is dark today, I am afraid."},
	{speaker: "Pink Panther", text: "Actually, I need an OFFERING bowl - one blessed at your hearth, for the old cherry tree."},
	{speaker: "Hiro", text: "Ahh, ze Whispering Cherry. For zat I must light ze SACRED hearth - but my fire-striker, ze flint, a CROW stole it zis morning!"},
	{speaker: "Hiro", text: "Bring me my striker and I will bless your offering bowl in ze first flame."},
}

var hiroRamenPostDialog = []dialogEntry{
	{speaker: "Hiro", text: "No striker, no sacred flame, panther-san. Zat thieving crow! Bring it back and ze blessed bowl is yours."},
}

var hiroOpenDialog = []dialogEntry{
	{speaker: "Hiro", text: "MY STRIKER! You found it! Stand back - "},
	{speaker: "Hiro", text: "(steel on flint - a spark - ze sacred hearth flares blue-gold)"},
	{speaker: "Hiro", text: "Zere. A bowl blessed in ze first flame. Carry it to ze old tree, with respect."},
}

// --- Kenji (street): saw where the crow dropped it; needs well-water first ---
var kenjiStudentDialog = []dialogEntry{
	{speaker: "Kenji", text: "Please - do not nudge ze table... oh, ze panther. You have ze look of a man hunting a crow's hiding place."},
	{speaker: "Pink Panther", text: "Hiro's fire-striker. You saw where the crow dropped it?"},
	{speaker: "Kenji", text: "I did. But my ink has dried to dust and I cannot think with a dry brush. Bring me water from ze temple well - just zere - and I will tell you."},
}

var kenjiWaterDialog = []dialogEntry{
	{speaker: "Kenji", text: "Ahh, cool well-water. Now ze ink flows... and so does my memory."},
	{speaker: "Kenji", text: "Ze crow dropped ze striker in ze flower store's eaves. Oba-chan keeps every shiny lost thing - ask her for it."},
	{speaker: "Kenji", text: "And here - I brushed you ze kanji for 'voice'. For ze quiet girl. A heart should carry one."},
}

var kenjiStudentPostDialog = []dialogEntry{
	{speaker: "Kenji", text: "Ink is just a voice zat takes its time. Ask Oba-chan for ze striker, panther-san."},
}

// --- Oba-chan (flower store): initial gate, gives the striker, then leads ---
var obachanDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "..."},
	{speaker: "Pink Panther", text: "Hello, madame. I'm looking for a blossom for a girl who's lost her voice inside."},
	{speaker: "Oba-chan", text: "Ze Whispering Cherry can mend such a heart. But it blooms only for a true offering, and ze path is not for strangers. Earn it first."},
}

var obachanStrikerDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "Hiro's fire-striker? Hah - a crow dropped it in my eaves zis morning, ze little thief."},
	{speaker: "Oba-chan", text: "Here. Take it back to him, and ze street will eat again."},
}

var obachanLeadDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "A bowl blessed at ze first flame, and a voice charm besides. Now you carry something to GIVE."},
	{speaker: "Oba-chan", text: "Come - follow me. I open ze path to ze old tree. Pick ze blossom yourself; it means more zat way."},
}

var obachanPostDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "Ze path is open, panther-san. Past ze shop, into ze pink grove - ze oldest tree. Pick gently."},
}

// --- The old tree in the grove: place the offering, then pick the blossom ---
var groveTreeNeedOfferingDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "The old tree's branches are bare. Oba-chan said it blooms for an offering blessed at the hearth - I should bring one first."},
}

var groveTreeDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Here, old tree. A bowl from the first flame, and a charm for a quiet voice. Please... one blossom, for Lily."},
	{speaker: "Pink Panther", text: "(he sets the offering at the roots - and the whole tree shivers awake into bloom)"},
}

var groveTreeDoneDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "I've got my blossom. Light as a breath. Time to take it home."},
}

// Kiku the dresser-geisha: she dresses PP (the kimono-spin gag) AND teaches him
// the way of tea. PP can't begin the matcha quest until he's heard her - she's
// what unlocks the matcha + bowl shelves (jp_tea_learned).
var dresserDialog = []dialogEntry{
	{speaker: "Kiku", text: "Ara! A pink panther in MY shop and not one stitch of silk on him? Unforgivable. Hold still - SPIN!"},
	{speaker: "Pink Panther", text: "Whoa - okay, okay, I'm dressed. ...Actually, this is rather nice."},
	{speaker: "Kiku", text: "Of course it is. And now, properly dressed, you must learn ze way of TEA. Ze old grove does not open its heart to a restless guest."},
	{speaker: "Kiku", text: "Take matcha and a bowl from my shelves, draw fresh water at ze street well, whisk it - zen kneel with ze tea master up in ze temple house. THAT is how you still a racing heart."},
}

var dresserPostDialog = []dialogEntry{
	{speaker: "Kiku", text: "Matcha, a bowl, well-water - zen ze tea master in ze temple house. Go, go, panther-san!"},
}

// --- Tea master (flower store): the matcha ceremony that gates the grove ---
var teaMasterDialog = []dialogEntry{
	{speaker: "Tea Master", text: "You wish to enter ze old grove? Hm. Ze Whispering Cherry does not open for a racing heart."},
	{speaker: "Tea Master", text: "Bring me matcha from ze shelf, a bowl of your choosing, and water from ze street well. We will share a bowl, and your heart will be still enough."},
}

var teaMasterNeedDialog = []dialogEntry{
	{speaker: "Tea Master", text: "Not yet. Matcha from ze shelf, a bowl, water from ze well - whisk zem together, zen return to me."},
}

var teaMasterReadyDialog = []dialogEntry{
	{speaker: "Tea Master", text: "Ah - whisked just so. Into a kimono with you, and kneel. We drink in silence."},
}

// Plays while PP is SEATED (after the spin-and-sit one-shot).
var teaMasterSippingDialog = []dialogEntry{
	{speaker: "Tea Master", text: "(the whisk hums, the froth settles, you each take a slow sip... and the noise inside you quiets)"},
	{speaker: "Tea Master", text: "Zere - your heart is still now. Ze grove will welcome you. Go gently, panther-san."},
}

var teaMasterPostDialog = []dialogEntry{
	{speaker: "Tea Master", text: "Carry ze stillness with you into ze grove."},
}

// Flavor names for the "random cup" pickup (cosmetic only).
var teaBowlNames = []string{"ze crane", "ze pine branch", "ze grey wave", "ze autumn moon", "ze persimmon", "ze plum blossom"}

// dannyPhoneCallDialog (user beat): Danny rings PP in Kyoto right after the
// pressed sakura. He isn't calling from a city - he's dug up an old phone (and
// some very Roman-looking treasure) back at camp, foreshadowing his arc. PP's
// running gag: "I didn't know you guys had phones in the camp."
var dannyPhoneCallDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "...is my paw buzzing? (digs out a phone) ...Hello?"},
	{speaker: "Danny", text: "PP!! It's Danny! Whatcha doin'? You sound super far away - is that a TEMPLE BELL?"},
	{speaker: "Pink Panther", text: "Danny? ...Hang on. I didn't even know you kids HAD phones at this camp."},
	{speaker: "Danny", text: "We don't! I dug this one up by the flagpole! Found a buncha old gold coins too, and a map, and half a stone sword. There's like a WHOLE buried city under the soccer field, PP!"},
	{speaker: "Pink Panther", text: "A buried city. Under the camp. Of course there is."},
	{speaker: "Danny", text: "Gotta go, the hole's getting bigger! BYEEE-"},
	{speaker: "Pink Panther", text: "(to you) ...Gold coins. Ruins. That kid's trouble is going to have a Roman accent. But first - Lily's flower."},
}

// ---------- NPC constructors ----------

func newRamenSeller(renderer *sdl.Renderer) *npc {
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_hiro_idle.png", jpNPCDir + "ramen_seller.png"},
		[]string{jpNPCDir + "npc_hiro_talk.png", jpNPCDir + "ramen_seller_talk.png"})
	return &npc{
		idleGrid: idle, talkGrid: talk,
		// 2026-07-14: stood at the LEFT edge of his stall (the prop now sits at
		// screen 471-818); once the stall opens he moves behind the counter
		// (openRamenStall + §JP-HIRO-COUNTER upper-body sheets).
		bounds:         sdl.Rect{X: 470, Y: 360, W: 150, H: 250},
		name:           "Hiro",
		dialog:         hiroRamenDialog,
		talkFrameSpeed: 0.12,
	}
}

func newTouristTokyo(renderer *sdl.Renderer) *npc {
	// START state: Gary holds the guidebook UPSIDE-DOWN (the "opposite book"
	// sheets). After he flips it (onDialogEnd) the callback swaps him to the
	// plain npc_gary_idle (book correct).
	// Gary has light-coloured eyes — use connected-edge key so the engine
	// doesn't erase them when colour-keying white.
	var idle, talk []npcFrame
	if p := firstExisting(jpNPCDir+"npc_gary_idle_oposite_book.png", jpNPCDir+"npc_tourist.png"); p != "" {
		idle = loadNPCGridConnected(renderer, p, 8, 1)
	}
	if len(idle) == 0 {
		idle = loadNPCGridRow(renderer, jpFbkVendor8x2, 8, 2, 0)
	}
	if p := firstExisting(jpNPCDir+"npc_gary_talk_oposite_book.png", jpNPCDir+"npc_tourist_talk.png"); p != "" {
		talk = loadNPCGridConnected(renderer, p, 8, 1)
	}
	if len(talk) == 0 {
		talk = idle
	}
	n := &npc{
		idleGrid: idle, talkGrid: talk,
		// Placed at the torii arrival now (was the ramen street).
		bounds:         sdl.Rect{X: 560, Y: 370, W: 130, H: 240},
		name:           "Gary",
		dialog:         garyTokyoDialog,
		talkFrameSpeed: 0.12,
	}
	// "Flip the book" gag (§JP-TOURIST): one-shot of him turning the book over.
	if p := firstExisting(jpNPCDir+"npc_gary_flip_his_book.png", jpNPCDir+"npc_gary_flip.png", jpNPCDir+"npc_tourist_flip.png"); p != "" {
		if f := loadNPCGrid(renderer, p, 8, 1); len(f) > 0 {
			n.oneShotAnims = map[string][]npcFrame{"flip": f}
		}
	}
	return n
}

func newKenjiStudent(renderer *sdl.Renderer) *npc {
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_kenji_idle.png"},
		[]string{jpNPCDir + "npc_kenji_talk.png"})
	return &npc{
		idleGrid: idle, talkGrid: talk,
		bounds:         sdl.Rect{X: 940, Y: 380, W: 130, H: 240},
		name:           "Kenji",
		dialog:         kenjiStudentDialog,
		talkFrameSpeed: 0.12,
	}
}

func newTeaMaster(renderer *sdl.Renderer) *npc {
	// The tea master is a woman. npc_tea_master_talk.png is the canonical
	// (female, flat-cartoon) sheet. Prefer a matching female idle when present;
	// fall back to the talk sheet for the idle until the regen lands (§TEA-MASTER-IDLE).
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_tea_master_idle.png", jpNPCDir + "npc_tea_master_idle_f.png", jpNPCDir + "npc_tea_master_talk.png"},
		[]string{jpNPCDir + "npc_tea_master_talk.png"})
	return &npc{
		idleGrid: idle, talkGrid: talk,
		bounds:         sdl.Rect{X: 640, Y: 372, W: 140, H: 248},
		name:           "Tea Master",
		dialog:         teaMasterDialog,
		talkFrameSpeed: 0.12,
	}
}

func newObachan(renderer *sdl.Renderer) *npc {
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_obachan_idle.png", jpNPCDir + "old_lady_idle.png"},
		[]string{jpNPCDir + "npc_obachan_talk.png", jpNPCDir + "old_lady.png"})
	return &npc{
		idleGrid: idle, talkGrid: talk,
		bounds:         sdl.Rect{X: 380, Y: 370, W: 140, H: 250},
		name:           "Oba-chan",
		dialog:         obachanDialog,
		talkFrameSpeed: 0.12,
	}
}

func newDresser(renderer *sdl.Renderer) *npc {
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_geisha_idle.png", jpNPCDir + "npc_dresser_idle.png", jpNPCDir + "drawer.png"},
		[]string{jpNPCDir + "npc_geisha_talk.png", jpNPCDir + "npc_dresser_talk.png", jpNPCDir + "drawer_talk.png"})
	return &npc{
		idleGrid: idle, talkGrid: talk,
		bounds:         sdl.Rect{X: 980, Y: 360, W: 150, H: 250},
		name:           "Kiku",
		dialog:         dresserDialog,
		talkFrameSpeed: 0.20,
	}
}

// ---------- Scene builders ----------

// addTokyoScenes decorates the 5 Japan/Kyoto scenes with particles, glows, and
// ambient leaf sprites. The static scene data (bg, spawn, hotspots, NPCs,
// walkSegments) is loaded from JSON before this function runs.
func addTokyoScenes(sm *sceneManager, renderer *sdl.Renderer) {
	if torii, ok := sm.scenes["tokyo_torii"]; ok {
		for i := 0; i < 16; i++ {
			torii.particles = append(torii.particles, particle{
				x: rand.Float64() * float64(engine.ScreenWidth), y: rand.Float64() * float64(engine.ScreenHeight),
				vx: (rand.Float64() - 0.4) * 18, vy: 14 + rand.Float64()*10,
				alpha: uint8(rand.Intn(30) + 50), size: int32(rand.Intn(3) + 2), r: 255, g: 180, b: 200,
			})
		}
	}

	// Live falling leaves over the tree (user: "leaves fall like a live
	// animation"). No-ops until the leaf sheet lands (§JP-LEAVES).
	if street, ok := sm.scenes["tokyo_street"]; ok {
		leafSpots := []struct{ x, y, scale, speed, drift, sec float64 }{
			{420, -40, 0.35, 55, 14, 0.18}, {560, -160, 0.28, 42, -10, 0.22},
			{700, -90, 0.40, 64, 8, 0.16}, {300, -200, 0.25, 48, 18, 0.2},
		}
		for _, l := range leafSpots {
			street.ambientSprites = append(street.ambientSprites,
				newAmbientLeafFall(renderer, jpLeafFall, 3, l.x, l.y, l.scale, l.speed, l.drift, l.sec))
		}
	}

	if grove, ok := sm.scenes["tokyo_temple"]; ok {
		grove.glows = []glowEffect{{x: 200, y: 100, w: 600, h: 320, r: 255, g: 220, b: 230, alpha: 12, pulse: 0.3}}
		for i := 0; i < 14; i++ {
			grove.particles = append(grove.particles, particle{
				x: rand.Float64() * float64(engine.ScreenWidth), y: rand.Float64() * 420,
				vx: (rand.Float64() - 0.4) * 10, vy: 10 + rand.Float64()*8,
				alpha: uint8(rand.Intn(20) + 40), size: int32(rand.Intn(2) + 2), r: 255, g: 180, b: 200,
			})
		}
	}

	if sakura, ok := sm.scenes["tokyo_sakura"]; ok {
		sakura.glows = []glowEffect{{x: 300, y: 80, w: 700, h: 360, r: 255, g: 200, b: 220, alpha: 14, pulse: 0.25}}
		for i := 0; i < 26; i++ {
			sakura.particles = append(sakura.particles, particle{
				x: rand.Float64() * float64(engine.ScreenWidth), y: rand.Float64() * float64(engine.ScreenHeight),
				vx: (rand.Float64() - 0.4) * 16, vy: 16 + rand.Float64()*12,
				alpha: uint8(rand.Intn(30) + 55), size: int32(rand.Intn(3) + 2), r: 255, g: 175, b: 200,
			})
		}
		for _, l := range []struct{ x, y, scale, speed, drift, sec float64 }{{500, -60, 0.7, 50, 12, 0.18}, {820, -180, 0.6, 60, -8, 0.2}} {
			sakura.ambientSprites = append(sakura.ambientSprites,
				newAmbientLeafFall(renderer, jpLeafFall, 3, l.x, l.y, l.scale, l.speed, l.drift, l.sec))
		}
	}

	if teahouse, ok := sm.scenes["tokyo_teahouse"]; ok {
		teahouse.glows = []glowEffect{{x: 250, y: 120, w: 600, h: 300, r: 245, g: 225, b: 190, alpha: 10, pulse: 0.2}}
	}
}

// openRamenStall swaps the stall prop to its "open" frame and seats the waiting
// line at the counter (the dynamic open→sit beat). Graceful: any missing art
// just leaves that sprite as-is. Sets jp_ramen_open.
func (g *Game) openRamenStall() {
	if g.ramenStoreProp != nil && len(g.ramenOpenFrames) > 0 {
		g.ramenStoreProp.frames = g.ramenOpenFrames
	}
	// The open prop's stools sit at screen x≈570/621/682/733, ground y≈530
	// (prop at 645/530 scale 0.30). Seat the two queue members on stools 1 and
	// 3, shrunk to the stall's mid-ground depth (their sit sheet includes the
	// stool, so they anchor on the same ground line as the prop).
	stoolX := []float64{570, 682}
	for i, c := range g.ramenQueue {
		if c == nil {
			continue
		}
		if len(g.ramenSitFrames) > 0 {
			c.frames = g.ramenSitFrames
		}
		if i < len(stoolX) {
			c.x = stoolX[i]
		} else {
			c.x = 733
		}
		c.y = 530
		c.scale = 0.23
	}
	// Hiro takes his place BEHIND the counter, upper body only (user
	// 2026-07-14). Graceful: until the §JP-HIRO-COUNTER sheets land he keeps
	// standing beside the stall with his full-body art.
	if street, ok := g.sceneMgr.scenes["tokyo_street"]; ok {
		for _, n := range street.npcs {
			if n.name != "Hiro" {
				continue
			}
			counterIdle := firstExisting(jpNPCDir + "npc_hiro_counter_idle.png")
			if counterIdle == "" {
				break
			}
			idle := loadNPCGridConnected(g.renderer, counterIdle, 8, 1)
			if len(idle) == 0 {
				break
			}
			n.idleGrid = idle
			n.talkGrid = idle
			if p := firstExisting(jpNPCDir + "npc_hiro_counter_talk.png"); p != "" {
				if talk := loadNPCGridConnected(g.renderer, p, 8, 1); len(talk) > 0 {
					n.talkGrid = talk
				}
			}
			// Waist-cut sprite in the counter window: bottom edge lands on the
			// counter top (screen y≈428) so the counter occludes him naturally.
			n.bounds = sdl.Rect{X: 600, Y: 353, W: 95, H: 75}
			n.approachXOverride = 645
			n.approachYOverride = 385 // PP stays on the street walk line
			break
		}
	}
	g.vars.SetBool(ScopeGame, VarJpRamenOpen, true)
}

func (g *Game) setupTokyoCallbacks() {
	game := g
	give := func(id string) {
		if item := game.items.createItem(id); item != nil {
			game.inv.addItem(item)
		}
	}

	if torii, ok := g.sceneMgr.scenes["tokyo_torii"]; ok {
		// No travel-map hotspot here — torii is a walk-through gate, not a hub.
		// Gary greets PP at the gates. He starts holding the guidebook
		// UPSIDE-DOWN; when PP talks to him he flips it and KEEPS it right-way-up
		// (swap idle+talk to the "book correct" sheets after the flip one-shot).
		// All graceful - no swap until the flipped sheets land (§JP-TOURIST).
		for _, n := range torii.npcs {
			if n.name != "Gary" {
				continue
			}
			gary := n
			var garyFlipIdle, garyFlipTalk []npcFrame
			// G2: after the flip Gary holds the book RIGHT-WAY-UP. The
			// npc_gary_idle_normal_book_idle/_talk sheets are the correct-book
			// poses (8×1). The old npc_gary_idle.png actually showed the book
			// UPSIDE-DOWN (6×2), which is why he "reverted to flipped" after a chat.
			if p := firstExisting(jpNPCDir+"npc_gary_idle_normal_book_idle.png", jpNPCDir+"npc_gary_idle_flipped.png"); p != "" {
				garyFlipIdle = loadNPCGridConnected(game.renderer, p, 8, 1)
			}
			if p := firstExisting(jpNPCDir+"npc_gary_idle_normal_book_talk.png", jpNPCDir+"npc_gary_talk_flipped.png", jpNPCDir+"npc_gary_talk.png"); p != "" {
				garyFlipTalk = loadNPCGridConnected(game.renderer, p, 8, 1)
			}
			gary.onDialogEnd = func() {
				gary.playOneShotAnim("flip", 1.0)
				if len(garyFlipIdle) > 0 {
					gary.idleGrid = garyFlipIdle
				}
				if len(garyFlipTalk) > 0 {
					gary.talkGrid = garyFlipTalk
				}
				gary.dialog = garyTokyoPostDialog
				// His ramen tip opens the stall: by the time PP reaches the street
				// it's lit and the waiting line has sat down at the counter.
				if !game.vars.GetBool(ScopeGame, VarJpRamenOpen) {
					g.openRamenStall()
				}
			}
		}
	}

	if street, ok := g.sceneMgr.scenes["tokyo_street"]; ok {
		// Dynamic ramen stall + waiting line. A closed/open prop over the stall
		// and a static line of 4 customers that SIT at the counter when Hiro
		// opens (openRamenStall). Art pending → invisible until it lands; state
		// restored on load.
		jpProp := "assets/images/locations/japan/props/"
		// Art landed 2026-07-14 (1536×1024, content ~1158×885 after the baked
		// checkerboard is keyed): bottom-center anchored over the painted-in BG
		// stall (BG px 555-875 → screen ~506-798, ground ≈527), scale 0.30 so the
		// prop covers it edge to edge.
		g.ramenStoreProp = newAmbientSway(g.renderer, firstExisting(jpProp+"ramen_closed.png"), 1, 645, 530, 0.30, 0.4)
		street.ambientSprites = append(street.ambientSprites, g.ramenStoreProp)
		g.ramenOpenFrames = loadAmbientStripKeyed(g.renderer, firstExisting(jpProp+"ramen_open.png"), 1)
		// customer_sit.png is a 4-frame eating loop (stool included in-frame) —
		// loading it as 8 cols split every figure in half (2026-07-14 fix).
		g.ramenSitFrames = loadAmbientStripKeyed(g.renderer, firstExisting(jpNPCDir+"customer_sit.png"), 4)
		// Two distinct queue members (different art) so the line doesn't look
		// like clones. #25: spaced ~180px apart (was 80 -> they overlapped into a
		// clone-pile) and stood in front of the shuttered stall (Hiro waits at
		// its left edge, x470), so they read as two people queuing on the
		// street, not a folder dump.
		queueEntries := []struct {
			path  string
			x, y  float64
			scale float64
		}{
			{firstExisting(jpNPCDir+"npc_ramen_tourist_wait.png", jpNPCDir+"customer_wait.png"), 520, 590, 0.62},
			{firstExisting(jpNPCDir+"npc_ramen_local_wait.png", jpNPCDir+"customer_wait.png"), 700, 590, 0.62},
		}
		for _, q := range queueEntries {
			// #25: the wait sheets are an 8-FRAME sway of ONE person. Loading them
			// as 1 column drew the entire 8-figure strip as a single sprite — that
			// was the "crowd of identical men". Load 8×1 so each entry is ONE
			// animated person (the engine cuts + cycles the 8 poses).
			c := newAmbientSway(g.renderer, q.path, 8, q.x, q.y, q.scale, 0.5)
			street.ambientSprites = append(street.ambientSprites, c)
			g.ramenQueue = append(g.ramenQueue, c)
		}
		if game.vars.GetBool(ScopeGame, VarJpRamenOpen) {
			g.openRamenStall()
		}

		// The temple well (Kenji's water errand) - a hotspot in the street.
		street.hotspots = append(street.hotspots, hotspot{
			bounds: sdl.Rect{X: 120, Y: 430, W: 130, H: 180}, name: "The temple well",
			onInteract: func() bool {
				// Matcha ceremony: with the powder + a bowl in hand, whisk a proper
				// Matcha Bowl at the well's cool water.
				if game.inv.hasItem("Matcha") && game.inv.hasItem("Tea Bowl") && !game.inv.hasItem("Matcha Bowl") {
					game.dialog.startDialogWithCallback([]dialogEntry{
						{speaker: "Pink Panther", text: "Cool well-water, a scoop of matcha, a brisk whisk... a proper bowl. Now to find the tea master."},
					}, func() {
						game.inv.removeItem("Matcha")
						game.inv.removeItem("Tea Bowl")
						give("matcha_bowl")
					})
					return true
				}
				if game.inv.hasItem("Well-Water") {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "I've already got a cup of well-water for Kenji."},
					})
					return true
				}
				game.dialog.startDialogWithCallback([]dialogEntry{
					{speaker: "Pink Panther", text: "A cold stone well. I'll draw a cup - Kenji needs it for his ink."},
				}, func() {
					give("well_water")
				})
				return true
			},
		})

		for _, n := range street.npcs {
			switch n.name {
			case "Hiro":
				hiro := n
				hiro.onDialogEnd = func() {
					// Until PP earns the blessed bowl, the reminder is "bring my striker".
					if !game.inv.hasItem("Offering Bowl") {
						hiro.dialog = hiroRamenPostDialog
					}
				}
				// The stall is already OPEN (Gary's tip). Bring Hiro his fire-striker
				// → he lights the SACRED hearth and blesses the Offering Bowl.
				hiro.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !game.inv.hasItem("Fire-Striker") || game.inv.hasItem("Offering Bowl") {
						return nil, nil, nil
					}
					return hiroOpenDialog, func() {
						game.inv.removeItem("Fire-Striker")
						give("offering_bowl")
						hiro.dialog = []dialogEntry{
							{speaker: "Hiro", text: "Take ze blessed bowl to ze old tree, panther-san - and come back for noodles when your heart is light."},
						}
						hiro.altDialogFunc = nil
					}, &handOff{item: "Fire-Striker", returnItem: "Offering Bowl"}
				}
			case "Kenji":
				kenji := n
				kenji.onDialogEnd = func() {} // stays on the water ask until traded
				// Bring Kenji well-water → he points to Oba-chan's eaves + gives the Voice Charm.
				kenji.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !game.inv.hasItem("Well-Water") || game.inv.hasItem("Voice Charm") {
						return nil, nil, nil
					}
					return kenjiWaterDialog, func() {
						game.inv.removeItem("Well-Water")
						give("voice_charm")
						kenji.dialog = kenjiStudentPostDialog
						kenji.altDialogFunc = nil
					}, &handOff{item: "Well-Water", returnItem: "Voice Charm"}
				}
			}
		}
	}

	if grove, ok := g.sceneMgr.scenes["tokyo_temple"]; ok {
		for _, n := range grove.npcs {
			switch n.name {
			case "Oba-chan":
				oba := n
				// Multi-stage: (1) once PP has Kenji's clue (Voice Charm) she hands
				// over the crow-dropped Fire-Striker; (2) once PP carries the
				// blessed Offering Bowl she "follow me"s him - opening the grove.
				oba.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if game.inv.hasItem("Offering Bowl") && !game.vars.GetBool(ScopeGame, VarJpGroveRevealed) {
						return obachanLeadDialog, func() {
							game.vars.SetBool(ScopeGame, VarJpGroveRevealed, true)
							oba.dialog = obachanPostDialog
						}, nil
					}
					if game.inv.hasItem("Voice Charm") && !game.inv.hasItem("Fire-Striker") &&
						!game.vars.GetBool(ScopeGame, VarJpRamenOpen) {
						return obachanStrikerDialog, func() {
							give("fire_striker")
						}, &handOff{returnItem: "Fire-Striker"}
					}
					return nil, nil, nil
				}
			case "Kiku":
				kiku := n
				kiku.onDialogEnd = func() {
					// She spins PP into a kimono (the gag) AND teaches the tea
					// ceremony - talking to her unlocks the matcha + bowl shelves
					// (jp_tea_learned). One-shot no-ops until PP_kimono_spin.png lands.
					game.player.playOneShot("kimono_spin", 1.8, nil)
					game.vars.SetBool(ScopeGame, VarJpTeaLearned, true)
					kiku.dialog = dresserPostDialog
				}
			}
		}
		// Tea-ceremony shelves in the flower store: a matcha tin + a shelf of
		// chawan (you're handed a random one).
		grove.hotspots = append(grove.hotspots, hotspot{
			bounds: sdl.Rect{X: 150, Y: 250, W: 120, H: 150}, name: "A tin of matcha",
			onInteract: func() bool {
				if !game.vars.GetBool(ScopeGame, VarJpTeaLearned) {
					game.dialog.startDialog([]dialogEntry{{speaker: "Pink Panther", text: "Pretty green powder... but I wouldn't know what to do with it. Maybe the kimono lady can teach me."}})
					return true
				}
				if game.inv.hasItem("Matcha") || game.vars.GetBool(ScopeGame, VarJpTeaDone) {
					game.dialog.startDialog([]dialogEntry{{speaker: "Pink Panther", text: "I've already got the matcha."}})
					return true
				}
				game.dialog.startDialogWithCallback([]dialogEntry{
					{speaker: "Pink Panther", text: "Bright green matcha powder. Just like Kiku said - the tea master will want this."},
				}, func() { give("matcha") })
				return true
			},
		})
		grove.hotspots = append(grove.hotspots, hotspot{
			bounds: sdl.Rect{X: 290, Y: 250, W: 120, H: 150}, name: "A shelf of tea bowls",
			onInteract: func() bool {
				if !game.vars.GetBool(ScopeGame, VarJpTeaLearned) {
					game.dialog.startDialog([]dialogEntry{{speaker: "Pink Panther", text: "Lovely bowls. I shouldn't just grab one - I should learn the proper way first. The kimono lady, maybe."}})
					return true
				}
				if game.inv.hasItem("Tea Bowl") || game.inv.hasItem("Matcha Bowl") || game.vars.GetBool(ScopeGame, VarJpTeaDone) {
					game.dialog.startDialog([]dialogEntry{{speaker: "Pink Panther", text: "I've got a bowl already."}})
					return true
				}
				bowl := teaBowlNames[rand.Intn(len(teaBowlNames))]
				game.dialog.startDialogWithCallback([]dialogEntry{
					{speaker: "Pink Panther", text: "So many chawan to choose from... I'll take " + bowl + " one today."},
				}, func() { give("tea_bowl") })
				return true
			},
		})
		// Exit INTO the hidden grove - needs BOTH Oba-chan's opened path AND a
		// still heart (the tea ceremony). Right edge, opposite the street exit.
		grove.hotspots = append(grove.hotspots, hotspot{
			bounds: sdl.Rect{X: 1300, Y: 180, W: 100, H: 460}, name: "Into the sakura grove", arrow: arrowRight,
			onInteract: func() bool {
				if !game.vars.GetBool(ScopeGame, VarJpGroveRevealed) {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "Just trees and a fence this way. Oba-chan said she'd show me the path - I should talk to her first."},
					})
					return true
				}
				if !game.vars.GetBool(ScopeGame, VarJpTeaDone) {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "The path's open... but my head's still racing. The tea master said to share a bowl and still my heart before the grove. Not yet."},
					})
					return true
				}
				game.sceneMgr.transitionTo("tokyo_sakura", game.player)
				return true
			},
		})
		// G9: the path UP to the temple tea-house should recede (shrink + drift
		// up) instead of marching PP off the top of the screen.
		for i := range grove.hotspots {
			if grove.hotspots[i].targetScene != "tokyo_teahouse" {
				continue
			}
			grove.hotspots[i].onInteract = func() bool {
				game.player.playRecede(1.0, 0.5, 80, func() {
					game.sceneMgr.transitionTo("tokyo_teahouse", game.player)
				})
				return true
			}
			break
		}
	}

	// Temple tea-house: share the whisked Matcha Bowl with the tea master →
	// jp_tea_done (the grove gate). No reward item; just the moment.
	if teahouse, ok := g.sceneMgr.scenes["tokyo_teahouse"]; ok {
		for _, n := range teahouse.npcs {
			if n.name != "Tea Master" {
				continue
			}
			tea := n
			tea.onDialogEnd = func() {
				if !game.vars.GetBool(ScopeGame, VarJpTeaDone) {
					tea.dialog = teaMasterNeedDialog
				}
			}
			tea.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
				if !game.inv.hasItem("Matcha Bowl") || game.vars.GetBool(ScopeGame, VarJpTeaDone) {
					return nil, nil, nil
				}
				return teaMasterReadyDialog, func() {
					game.inv.removeItem("Matcha Bowl")
					// PP spins fast into a kimono and kneels (tea_sit one-shot); then
					// the seated ceremony dialog plays with him in his SITTING poses;
					// then he stands and the grove gate opens.
					game.player.playOneShot("tea_sit", 2.2, func() {
						game.player.seated = true
						game.dialog.startDialogWithCallback(teaMasterSippingDialog, func() {
							game.player.seated = false
							game.vars.SetBool(ScopeGame, VarJpTeaDone, true)
							tea.dialog = teaMasterPostDialog
							tea.altDialogFunc = nil
						})
					})
				}, &handOff{item: "Matcha Bowl"}
			}
		}
	}

	// Hidden sakura grove: the old tree is the pick-the-blossom payoff. Picking
	// gives the Pressed Sakura (the anchor) + fires Danny's foreshadow call.
	if sakura, ok := g.sceneMgr.scenes["tokyo_sakura"]; ok {
		for i := range sakura.hotspots {
			if sakura.hotspots[i].name != "The oldest cherry tree" {
				continue
			}
			sakura.hotspots[i].onInteract = func() bool {
				if game.inv.hasItem("Pressed Sakura") {
					game.dialog.startDialog(groveTreeDoneDialog)
					return true
				}
				// The tree only blooms once PP sets the blessed offering at its roots.
				if !game.inv.hasItem("Offering Bowl") {
					game.dialog.startDialog(groveTreeNeedOfferingDialog)
					return true
				}
				game.dialog.startDialogWithCallback(groveTreeDialog, func() {
					// Place the offering + the voice charm (both consumed here, so no
					// Kyoto item lingers), then PP picks a blossom (reuses the
					// flower-grab one-shot); the petal lands in the bag + Danny calls.
					game.inv.removeItem("Offering Bowl")
					game.inv.removeItem("Voice Charm")
					game.player.playOneShot("grab_flower", 0.9, func() {
						give("pressed_sakura")
						game.dialog.startDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "A real sakura blossom. Light as a breath. Lily will hold this and come back to herself."},
						})
						game.dialog.queueDialog(dannyPhoneCallDialog)
					})
				})
				return true
			}
			break
		}
	}
}

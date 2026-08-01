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
//   - kyoto_torii  (torii gate corridor): PP lands here. Right -> ramen street.
//   - kyoto_street (ramen store; the tree drops leaves): Hiro the cook + Gary
//     the tourist. Left -> torii, right -> the flower grove.
//   - kyoto_temple (flower store by the forest): Oba-chan presses the sakura
//     (the anchor) + Kiku the dresser spins PP into a kimono for a gag. Left ->
//     the ramen street.
//
// Anchor item "Pressed Sakura" -> PP carries it home and heals Lily at the lake.
//
// The Japan art has been re-saved under several different filenames during
// authoring, so every asset is resolved from a CANDIDATE LIST (firstExisting):
// whichever name is on disk wins, and a missing asset degrades gracefully
// (flat-colour BG / placeholder NPC). Scene keys stay "kyoto_*".

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
	// 2026-07-24 (user #37/#40): CONNECTED key for every Japan NPC — the
	// global tol-8 key punched their white eye sclera into holes (Hiro lost
	// his eyes in both sheets; same class for Kenji/Kiku/Tea Master).
	if p := firstExisting(idleCands...); p != "" {
		idle = loadNPCGridConnected(renderer, p, 8, 1)
	}
	if len(idle) == 0 {
		idle = loadNPCGridRow(renderer, jpFbkVendor8x2, 8, 2, 0)
	}
	if p := firstExisting(talkCands...); p != "" {
		talk = loadNPCGridConnected(renderer, p, 8, 1)
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
	{speaker: "Gary", text: "PINK PANTHER! Can you BELIEVE it - KYOTO! I've dreamed about this since I was a kid. The temples, the gardens, the cherry blossoms..."},
	{speaker: "Gary", text: "I read EVERYTHING about this place - every shrine, every festival, every bowl of noodles. It's ALL here in my guidebook!"},
	{speaker: "Pink Panther", text: "...Gary. You're holding the book upside down."},
	{speaker: "Gary", text: "WHAT? ...Oh. Oh my. ...Ahh. MUCH better. Now it makes far more sense!"},
	{speaker: "Gary", text: "Listen - the book says the old cherry in the hidden grove only blooms for an offering blessed at a sacred hearth. The ramen cook down the street keeps one!"},
	{speaker: "Pink Panther", text: "A blessed offering. Ramen first, then. Thanks, Gary."},
}

// 2026-07-18 (user #45/#46): Gary is the chapter's CUSTOMS ORACLE — whenever
// PP doesn't know the local etiquette, Gary + guidebook has the answer. His
// repeat dialog advances with the chapter flags (see setupKyotoCallbacks).
var garyTokyoPostDialog = []dialogEntry{
	{speaker: "Gary", text: "Need a custom checked? My book knows everything. Now it says Kyoto is in PERU, but I'm reading between the lines."},
}

var garyDressCodeDialog = []dialogEntry{
	{speaker: "Gary", text: "The book is very clear: you do NOT visit the Whispering Cherry in street clothes. It's a dress-code thing."},
	{speaker: "Gary", text: "The lady at the flower store dresses visitors properly. Tell her the tree sent you. Well - don't, actually. That sounds weird."},
}

var garyTeaCustomDialog = []dialogEntry{
	{speaker: "Gary", text: "Before the grove, the book says you need 'a still heart' - that's the tea ceremony. Matcha, a bowl, cool well water."},
	{speaker: "Gary", text: "The tea master at the temple tea-house hosts it. Very serene. I cried a little."},
}

// --- Hiro (street): OPEN for business (Gary's tip), but the SACRED hearth for
// the blessed offering bowl still needs his kite-stolen fire-striker ---
var hiroRamenDialog = []dialogEntry{
	{speaker: "Hiro", text: "Irasshaimase, panther-san! You have the look of a hungry traveller - but my stall is dark today, I am afraid."},
	{speaker: "Pink Panther", text: "Actually, I need an OFFERING bowl - one blessed at your hearth, for the old cherry tree."},
	{speaker: "Hiro", text: "Ahh, the Whispering Cherry. For that I must light the SACRED hearth - but my fire-striker, the flint, a KITE swooped down and stole it this morning!"},
	{speaker: "Hiro", text: "Bring me my striker and I will bless your offering bowl in the first flame."},
}

var hiroRamenPostDialog = []dialogEntry{
	{speaker: "Hiro", text: "No striker, no sacred flame, panther-san. That thieving kite! Bring it back and the blessed bowl is yours."},
}

var hiroOpenDialog = []dialogEntry{
	{speaker: "Hiro", text: "MY STRIKER! You found it! Stand back - "},
	{speaker: "Hiro", text: "(steel on flint - a spark - the sacred hearth flares blue-gold)"},
	{speaker: "Hiro", text: "There. A bowl blessed in the first flame. Carry it to the old tree, with respect."},
}

// --- Kenji (street): saw where the kite dropped it; needs well-water first ---
var kenjiStudentDialog = []dialogEntry{
	{speaker: "Kenji", text: "Please - do not nudge the table... oh, the panther. You have the look of a man hunting a kite's hiding place."},
	{speaker: "Pink Panther", text: "Hiro's fire-striker. You saw where the kite dropped it?"},
	{speaker: "Kenji", text: "I did. But my ink has dried to dust and I cannot think with a dry brush. Bring me water from the temple well - down in the ramen street, at the foot of the stairs - and I will tell you."},
}

var kenjiWaterDialog = []dialogEntry{
	{speaker: "Kenji", text: "Ahh, cool well-water. Now the ink flows... and so does my memory."},
	{speaker: "Kenji", text: "The kite dropped the striker in the flower store's eaves. Oba-chan keeps every shiny lost thing - ask her for it."},
	{speaker: "Kenji", text: "And here - I brushed you the kanji for 'voice'. For the quiet girl. A heart should carry one."},
}

var kenjiStudentPostDialog = []dialogEntry{
	{speaker: "Kenji", text: "Ink is just a voice that takes its time. Ask Oba-chan for the striker, panther-san."},
}

// --- Takeshi (street): the KAMISHIBAI storyteller — culture flavor, no
// quest (user 2026-07-27, v2: "just tell us to eat loud" was flat). The
// classic street paper-theater man with a little wooden card stage. Each
// chat is one picture-tale; the tales teach real Japan AND touch things
// the player can see or will do (the cherry legend primes the grove quest,
// the cat tale points at the roof-cat prop, the gates tale covers the
// torii bow). onDialogEnd rotates to the next tale.
var takeshiCherryTaleDialog = []dialogEntry{
	{speaker: "Takeshi", text: "Ah, a customer for my paper theater! Sit, sit. Today's tale... (flips a card) ...The Whispering Cherry."},
	{speaker: "Takeshi", text: "Deep in a hidden grove stands a tree older than the temple bells. Most years it will not bloom at all. Gold cannot buy its flowers. Princes have tried."},
	{speaker: "Takeshi", text: "(flips the card) But bring it a true gift - something made with care, given with a whole heart - and the old tree wakes... and WHISPERS back."},
	{speaker: "Pink Panther", text: "A tree that answers. I know a little girl who'd love that story."},
	{speaker: "Takeshi", text: "Heh. It is no story, pink one. Ask any old man in this street."},
}

var takeshiCatTaleDialog = []dialogEntry{
	{speaker: "Takeshi", text: "(flips a card) The Beckoning Cat! Long ago, a poor shopkeeper shared his last rice ball with a stray cat..."},
	{speaker: "Takeshi", text: "The cat sat at his door and raised one paw - and travelers followed it in, curious. The shop never stood empty again. We call it maneki-neko, the beckoning cat."},
	{speaker: "Takeshi", text: "Look up - the stray on Hiro's roof. Left paw tucked, waiting. When that stall opens again, watch her raise it. Luck knows where soup is coming."},
}

var takeshiGatesTaleDialog = []dialogEntry{
	{speaker: "Takeshi", text: "(flips a card) The Red Gates. A torii is a doorway with no wall - because the wall it passes through cannot be seen."},
	{speaker: "Takeshi", text: "On this side, the everyday world. Through the gate - the sacred. So we bow once as we pass. Small manners for big neighbours."},
	{speaker: "Takeshi", text: "And since you ask about manners: at the noodle stand, SLURP LOUDLY. Quiet noodles insult the cook. That card I painted from experience."},
}

// --- Oba-chan (flower store): initial gate, gives the striker, then leads ---
var obachanDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "..."},
	{speaker: "Pink Panther", text: "Hello, madame. I'm looking for a blossom for a girl who's lost her voice inside."},
	{speaker: "Oba-chan", text: "The Whispering Cherry can mend such a heart. But it blooms only for a true offering, and the path is not for strangers. Earn it first."},
}

var obachanStrikerDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "Hiro's fire-striker? Hah - a kite dropped it in my eaves this morning, the little thief."},
	{speaker: "Oba-chan", text: "Here. Take it back to him, and the street will eat again."},
}

var obachanLeadDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "A bowl blessed at the first flame, and a voice charm besides. Now you carry something to GIVE."},
	{speaker: "Oba-chan", text: "Come - follow me. I open the path to the old tree. Pick the blossom yourself; it means more that way."},
}

var obachanPostDialog = []dialogEntry{
	{speaker: "Oba-chan", text: "The path is open, panther-san. Past the shop, into the pink grove - the oldest tree. Pick gently."},
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
	{speaker: "Kiku", text: "Of course it is. And now, properly dressed, you must learn the way of TEA. The old grove does not open its heart to a restless guest."},
	{speaker: "Kiku", text: "Take matcha and a bowl from my shelves, draw fresh water at the street well, whisk it - then kneel with the tea master up in the temple house. THAT is how you still a racing heart."},
}

var dresserPostDialog = []dialogEntry{
	{speaker: "Kiku", text: "Matcha, a bowl, well-water - then the tea master in the temple house. Go, go, panther-san!"},
}

// --- Tea master (flower store): the matcha ceremony that gates the grove ---
var teaMasterDialog = []dialogEntry{
	{speaker: "Tea Master", text: "You wish to enter the old grove? Hm. The Whispering Cherry does not open for a racing heart."},
	{speaker: "Tea Master", text: "Bring me matcha from the shelf, a bowl of your choosing, and water from the street well. We will share a bowl, and your heart will be still enough."},
}

var teaMasterNeedDialog = []dialogEntry{
	{speaker: "Tea Master", text: "Not yet. Matcha from the shelf, a bowl, water from the well - whisk them together, then return to me."},
}

var teaMasterReadyDialog = []dialogEntry{
	{speaker: "Tea Master", text: "Ah - whisked just so. Into a kimono with you, and kneel. We drink in silence."},
}

// Plays while PP is SEATED (after the spin-and-sit one-shot).
var teaMasterSippingDialog = []dialogEntry{
	{speaker: "Tea Master", text: "(the whisk hums, the froth settles, you each take a slow sip... and the noise inside you quiets)"},
	{speaker: "Tea Master", text: "There - your heart is still now. The grove will welcome you. Go gently, panther-san."},
}

var teaMasterPostDialog = []dialogEntry{
	{speaker: "Tea Master", text: "Carry the stillness with you into the grove."},
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

func newTouristKyoto(renderer *sdl.Renderer) *npc {
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
	// 2026-07-24 (user #40): CONNECTED key — the global key ate his eye
	// whites mid-flip (his idle/talk were already connected).
	if p := firstExisting(jpNPCDir+"npc_gary_flip_his_book.png", jpNPCDir+"npc_gary_flip.png", jpNPCDir+"npc_tourist_flip.png"); p != "" {
		if f := loadNPCGridConnected(renderer, p, 8, 1); len(f) > 0 {
			n.oneShotAnims = map[string][]npcFrame{"flip": f}
		}
	}
	return n
}

func newKenjiStudent(renderer *sdl.Renderer) *npc {
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_kenji_idle.png"},
		[]string{jpNPCDir + "npc_kenji_talk.png"})
	n := &npc{
		idleGrid: idle, talkGrid: talk,
		// 2026-07-27 (user): moved into the TORII scene — a calligraphy
		// student's natural pitch is at the shrine gates, selling brushwork
		// to arriving visitors. Gives Gary's scene a second beat, and his
		// well errand now ping-pongs between two scenes (torii <-> street),
		// matching the Paris bakery<->street rhythm.
		// 2026-08-01 (user #24): he rendered huge — a young student should sit
		// well under PP. H 240 → 190 (W scaled with it), feet kept at 620.
		bounds:         sdl.Rect{X: 260, Y: 430, W: 105, H: 190},
		name:           "Kenji",
		dialog:         kenjiStudentDialog,
		talkFrameSpeed: 0.12,
	}
	// §KENJI-GIVE-CHARM (landed 2026-07-18): his brush-the-charm one-shot.
	if p := firstExisting(jpNPCDir + "npc_kenji_give_charm.png"); p != "" {
		if f := loadNPCGridConnected(renderer, p, 8, 1); len(f) > 0 {
			n.oneShotAnims = map[string][]npcFrame{"give": f}
		}
	}
	return n
}

func newTakeshi(renderer *sdl.Renderer) *npc {
	// §JP-TAKESHI (queued): the elderly KAMISHIBAI storyteller behind his
	// little wooden card stage. Placeholder vendor row until the sheets land.
	idle, talk := loadJapanNPC(renderer,
		[]string{jpNPCDir + "npc_takeshi_idle.png"},
		[]string{jpNPCDir + "npc_takeshi_talk.png"})
	n := &npc{
		idleGrid: idle, talkGrid: talk,
		// Mid-right of the street (Kenji's old area, now free): queue at
		// 340-470, stall window 510-760.
		// 2026-08-01 (user #23): shifted right 885 → 920 — as far as the
		// street allows without covering the well hotspot (1015-1100; NPC
		// clicks win over hotspots, so his box must leave it reachable).
		bounds:         sdl.Rect{X: 920, Y: 385, W: 125, H: 235},
		name:           "Takeshi",
		dialog:         takeshiCherryTaleDialog,
		talkFrameSpeed: 0.14,
	}
	// §JP-TAKESHI-DOORS (user 2026-07-29): the kamishibai stage gets working
	// DOORS — idle re-rolls with the doors closed, and these one-shots open
	// them when a tale starts / shut them when it ends. Graceful no-ops
	// until the sheets land.
	for key, file := range map[string]string{
		"open_doors":  "npc_takeshi_open_doors.png",
		"close_doors": "npc_takeshi_close_doors.png",
	} {
		if p := firstExisting(jpNPCDir + file); p != "" {
			if f := loadNPCGridConnected(renderer, p, 8, 1); len(f) > 0 {
				if n.oneShotAnims == nil {
					n.oneShotAnims = map[string][]npcFrame{}
				}
				n.oneShotAnims[key] = f
			}
		}
	}
	return n
}

func newTeaMaster(renderer *sdl.Renderer) *npc {
	// The regenerated blue tea-master sheets are dedicated, cleanly separated
	// 6x1 strips, so use the gap-aware cutter instead of clipping them to equal
	// mathematical cells.
	loadTeaMaster := func(cands ...string) []npcFrame {
		if p := firstExisting(cands...); p != "" {
			return loadNPCGridConnected(renderer, p, 6, 1)
		}
		return nil
	}
	idle := loadTeaMaster(jpNPCDir+"npc_tea_master_idle.png", jpNPCDir+"npc_tea_master_idle_f.png", jpNPCDir+"npc_tea_master_talk.png")
	if len(idle) == 0 {
		idle = loadNPCGridRow(renderer, jpFbkVendor8x2, 8, 2, 0)
	}
	talk := loadTeaMaster(jpNPCDir + "npc_tea_master_talk.png")
	if len(talk) == 0 {
		talk = idle
	}
	return &npc{
		idleGrid: idle, talkGrid: talk,
		// 2026-07-24 (user #43): interim size bump (H 248 → 280, foot kept
		// at 620) — her placeholder sheets draw the figure tiny in-cell.
		// The real fix stays the queued §TEA-MASTER-BLUE-v2 re-roll.
		// 2026-07-26 PR #6: nudged a little left (X 640 → 615) per user.
		bounds:         sdl.Rect{X: 615, Y: 340, W: 140, H: 280},
		name:           "Tea Master",
		dialog:         teaMasterDialog,
		talkFrameSpeed: 0.12,
		// 2026-08-01 (user #28): he kneels on the raised tatami platform —
		// foot-aligning PP to his bounds walked PP way up the platform
		// ("climbing to the roof"). elevated keeps PP on his own floor row
		// and he just walks across to the platform edge.
		elevated: true,
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
	// 2026-07-24 (user D7): landscape props — Kyoto had NO living ambience
	// (only petals/leaves/glows). Each is a silent no-op until its art lands
	// under japan/props/ (§JP-* prompts in EXTRA_PROMPTS).
	addJapanProp := func(sceneName, file string, x, y, scale, frameSec float64, tol uint8) {
		s, ok := sm.scenes[sceneName]
		if !ok || s == nil {
			return
		}
		if p := firstExisting(jpPropDir + file); p != "" {
			s.ambientSprites = append(s.ambientSprites,
				newAmbientSwayGlobalTol(renderer, p, 6, x, y, scale, frameSec, tol))
		}
	}
	// 2026-07-26 PR #2: the torii pilgrims are REMOVED — Gary's arrival scene
	// reads cleaner without figures behind him (the petal particles stay).
	// 2026-07-26 PR #8: cat repositioned onto the street awning per user.
	addJapanProp("kyoto_street", "street_cat.png", 866, 309, 0.3, 0.7, 8)
	// the thieving KITE (tobi) perched up on the temple roof — the bird
	// Hiro and Kenji blame for the fire-striker (§JP-KITE-SHINY).
	// 2026-07-26 PR #5: foot moved to (657,217) per user.
	addJapanProp("kyoto_temple", "kite_shiny.png", 657, 217, 0.32, 0.6, 8)
	// stone lantern with incense wisps by the flower store
	// 2026-07-26 PR #4: foot moved to (334,473) per user.
	addJapanProp("kyoto_temple", "incense_lantern.png", 334, 473, 0.42, 0.5, 8)
	// steaming kettle beside the tea master. 2026-07-26 PR #7: tol 24 — the
	// tol-8 global key left a blue AA seam pocket on the kettle.
	addJapanProp("kyoto_teahouse", "teahouse_kettle.png", 850, 610, 0.34, 0.45, 24)

	if torii, ok := sm.scenes["kyoto_torii"]; ok {
		for i := 0; i < 16; i++ {
			torii.particles = append(torii.particles, particle{
				x: rand.Float64() * float64(engine.ScreenWidth), y: rand.Float64() * float64(engine.ScreenHeight),
				vx: (rand.Float64() - 0.4) * 18, vy: 14 + rand.Float64()*10,
				alpha: uint8(rand.Intn(30) + 50), size: int32(rand.Intn(3) + 2), r: 255, g: 180, b: 200,
			})
		}
	}

	// 2026-07-27 (user): a waiting LINE at the LEFT side of the ramen stall —
	// three line PROPS (they live under japan/props/, like every other
	// landscape prop): two standing + one in the classic deep Asian squat.
	// Until the §JP-LINE-* props land, the two standing slots fall back to
	// the old npc/ checkerboard wait sheets (edge-connected keyed loader
	// strips either background). Staggered frame speeds so the three don't
	// sway in sync.
	if street, ok := sm.scenes["kyoto_street"]; ok {
		// 2026-08-01 (user #22): per-figure foot lines (were all 520) so the
		// three stand ON the stall's ground line: businessman 561, grandma
		// 558, squat man 557.
		for _, q := range []struct {
			cands    []string
			x, y     float64
			frameSec float64
		}{
			{[]string{jpPropDir + "npc_ramen_wait_salaryman.png", jpNPCDir + "npc_ramen_tourist_wait.png"}, 340, 561, 0.60},
			{[]string{jpPropDir + "npc_ramen_wait_grandma.png", jpNPCDir + "npc_ramen_local_wait.png"}, 405, 558, 0.50},
			{[]string{jpPropDir + "npc_ramen_local_squat.png"}, 470, 557, 0.55},
		} {
			if p := firstExisting(q.cands...); p != "" {
				street.ambientSprites = append(street.ambientSprites,
					newAmbientSway(renderer, p, 8, q.x, q.y, 0.34, q.frameSec))
			}
		}
	}

	// Live falling leaves over the tree (user: "leaves fall like a live
	// animation"). No-ops until the leaf sheet lands (§JP-LEAVES).
	if street, ok := sm.scenes["kyoto_street"]; ok {
		leafSpots := []struct{ x, y, scale, speed, drift, sec float64 }{
			{420, -40, 0.35, 55, 14, 0.18}, {560, -160, 0.28, 42, -10, 0.22},
			{700, -90, 0.40, 64, 8, 0.16}, {300, -200, 0.25, 48, 18, 0.2},
		}
		for _, l := range leafSpots {
			street.ambientSprites = append(street.ambientSprites,
				newAmbientLeafFall(renderer, jpLeafFall, 3, l.x, l.y, l.scale, l.speed, l.drift, l.sec))
		}
	}

	if grove, ok := sm.scenes["kyoto_temple"]; ok {
		grove.glows = []glowEffect{{x: 200, y: 100, w: 600, h: 320, r: 255, g: 220, b: 230, alpha: 12, pulse: 0.3}}
		for i := 0; i < 14; i++ {
			grove.particles = append(grove.particles, particle{
				x: rand.Float64() * float64(engine.ScreenWidth), y: rand.Float64() * 420,
				vx: (rand.Float64() - 0.4) * 10, vy: 10 + rand.Float64()*8,
				alpha: uint8(rand.Intn(20) + 40), size: int32(rand.Intn(2) + 2), r: 255, g: 180, b: 200,
			})
		}
	}

	if sakura, ok := sm.scenes["kyoto_sakura"]; ok {
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

	if teahouse, ok := sm.scenes["kyoto_teahouse"]; ok {
		teahouse.glows = []glowEffect{{x: 250, y: 120, w: 600, h: 300, r: 245, g: 225, b: 190, alpha: 10, pulse: 0.2}}
	}
}

func (g *Game) setupKyotoCallbacks() {
	game := g
	give := func(id string) {
		if item := game.items.createItem(id); item != nil {
			game.inv.addItem(item)
		}
	}

	if torii, ok := g.sceneMgr.scenes["kyoto_torii"]; ok {
		// No travel-map hotspot here — torii is a walk-through gate, not a hub.
		// Gary greets PP at the gates. He starts holding the guidebook
		// UPSIDE-DOWN; when PP talks to him he flips it and KEEPS it right-way-up
		// (swap idle+talk to the "book correct" sheets after the flip one-shot).
		// All graceful - no swap until the flipped sheets land (§JP-TOURIST).
		// Kenji lives at the gates now (2026-07-27) — his well-water trade
		// wiring moved here from the street block.
		for _, n := range torii.npcs {
			if n.name != "Kenji" {
				continue
			}
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
			break
		}

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
			// 2026-07-26 PR #3: the flip gag plays ONCE, on the first chat.
			// Every later dialog finds the book already right-way-up (the
			// swapped normal-book idle/talk sheets stay in place).
			garyFlipped := false
			gary.onDialogEnd = func() {
				if !garyFlipped {
					garyFlipped = true
					gary.playOneShotAnim("flip", 1.0)
					if len(garyFlipIdle) > 0 {
						gary.idleGrid = garyFlipIdle
					}
					if len(garyFlipTalk) > 0 {
						gary.talkGrid = garyFlipTalk
					}
				}
				// 2026-07-18 (user #45/#46): Gary = the customs oracle. His
				// repeat dialog advances with the chapter: dress code once the
				// offering path opens, tea etiquette once the grove is revealed
				// but the ceremony isn't done, generic book gag otherwise.
				switch {
				case game.vars.GetBool(ScopeGame, VarJpGroveRevealed) && !game.vars.GetBool(ScopeGame, VarJpTeaDone):
					gary.dialog = garyTeaCustomDialog
				case game.inv.hasItem("Offering Bowl") && !game.vars.GetBool(ScopeGame, VarJpTeaLearned):
					gary.dialog = garyDressCodeDialog
				default:
					gary.dialog = garyTokyoPostDialog
				}
			}

		}
	}

	if street, ok := g.sceneMgr.scenes["kyoto_street"]; ok {
		// 2026-07-18 (user #43/#44): the stall overlay prop + waiting-line
		// ambients are REMOVED — they never fit the painted stall and stacked
		// on the NPCs. The painted stall in the BG is the stall; Hiro stands
		// at it (and moves behind the counter permanently once the
		// §JP-HIRO-COUNTER waist-cut sheets land).
		for _, n := range street.npcs {
			if n.name != "Hiro" {
				continue
			}
			if counterIdle := firstExisting(jpNPCDir + "npc_hiro_counter_idle.png"); counterIdle != "" {
				if idle := loadNPCGridConnected(g.renderer, counterIdle, 6, 1); len(idle) > 0 {
					n.idleGrid = idle
					n.talkGrid = idle
					if p := firstExisting(jpNPCDir + "npc_hiro_counter_talk.png"); p != "" {
						if talk := loadNPCGridConnected(g.renderer, p, 6, 1); len(talk) > 0 {
							n.talkGrid = talk
						}
					}
					// Waist-cut bust in the stall window; PP stays on the walk line.
					n.bounds = sdl.Rect{X: 600, Y: 353, W: 95, H: 75}
					n.approachXOverride = 645
					n.approachYOverride = 385
					// 2026-07-29 (user): "random of him working" — after a few
					// idle seconds the work strip (chopping, ladling, wiping)
					// plays as the alt-idle punctuation (Marcus-freakout
					// mechanism). Graceful until §JP-HIRO-COUNTER's work
					// sheet lands.
					if p := firstExisting(jpNPCDir + "npc_hiro_counter_work.png"); p != "" {
						if work := loadNPCGridConnected(g.renderer, p, 6, 1); len(work) > 0 {
							n.altIdleGrid = work
							n.altIdleAfterSec = 6.0
						}
					}
					// 2026-07-31: the landed counter GIVE sheet (ladles broth,
					// tops the bowl, slides it forward) was on disk but wired
					// to NOTHING, so the Fire-Striker → Offering Bowl trade
					// animated on PP's side only. Registered under the
					// auto-derived key AND the generic "give" so playHandOff
					// finds it either way.
					if p := firstExisting(jpNPCDir + "npc_hiro_counter_give.png"); p != "" {
						if give := loadNPCGridConnected(g.renderer, p, 6, 1); len(give) > 0 {
							if n.oneShotAnims == nil {
								n.oneShotAnims = map[string][]npcFrame{}
							}
							n.oneShotAnims["give_offering_bowl"] = give
							n.oneShotAnims["give"] = give
						}
					}
				}
			}
			break
		}

		// The temple well (Kenji's water errand) - a hotspot in the street.
		// 2026-07-26 PR #9: moved from the empty left edge onto the little
		// ROOFED WELL painted at the base of the temple stairs (right side),
		// so the pickup has a visible anchor the player can find.
		street.hotspots = append(street.hotspots, hotspot{
			bounds: sdl.Rect{X: 1015, Y: 420, W: 85, H: 75}, name: "The temple well",
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

		// 2026-07-26 PR #10: the flower-store up-exit climbs the painted
		// STAIRS — PP first walks to the stair base (foot ≈ 1184,459; above
		// the walk band, so unclamped), then recedes (shrink + drift up) into
		// the transition, like the Jerusalem market→plaza exit.
		for i := range street.hotspots {
			if street.hotspots[i].targetScene != "kyoto_temple" {
				continue
			}
			street.hotspots[i].onInteract = func() bool {
				game.player.walkToAndDoUnclamped(1184, 324, func() {
					game.player.playRecede(0.9, 0.45, 70, func() {
						game.sceneMgr.transitionTo("kyoto_temple", game.player)
					})
				})
				return true
			}
			break
		}

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
							{speaker: "Hiro", text: "Take the blessed bowl to the old tree, panther-san - and come back for noodles when your heart is light."},
						}
						hiro.altDialogFunc = nil
					}, &handOff{item: "Fire-Striker", returnItem: "Offering Bowl", back: true}
				}
				// 2026-08-01 (user #25): Hiro is BEHIND his counter — chats with
				// him use the same walk-to-mark + recede-shrink choreography as
				// Pierre/Poulain, and PP shows his BACK to the camera. The
				// recedeHeld guard talks in place on chained clicks.
				hiro.ppFaceBack = true
				hiro.onClickOverride = func() {
					if game.player == nil || game.dialog == nil {
						return
					}
					talk := func() {
						game.player.state = stateTalking
						game.player.dir = dirUp
						game.player.facingLeft = false
						releaseRecede := func() {
							game.player.state = stateIdle
							game.player.holdRecede()
						}
						if hiro.altDialogFunc != nil {
							entries, cb, ho := hiro.altDialogFunc()
							if entries != nil {
								game.inv.heldItem = nil
								start := func() {
									game.dialog.startDialogWithCallback(entries, func() {
										if cb != nil {
											cb()
										}
										releaseRecede()
									})
								}
								if ho != nil {
									game.player.playHandOff(hiro, ho, start)
								} else {
									start()
								}
								return
							}
						}
						d := hiro.dialog
						game.dialog.startDialogWithCallback(d, func() {
							if hiro.onDialogEnd != nil {
								hiro.onDialogEnd()
							}
							releaseRecede()
						})
					}
					if game.player.recedeHeld {
						talk()
						return
					}
					// Stand mark: centre of the counter window, on the street lane.
					game.player.walkToAndDo(640, 520, func() {
						game.player.playRecede(1.0, 0.65, 50, talk)
					})
				}
			case "Takeshi":
				// Culture flavor only (user 2026-07-27): the kamishibai man
				// rotates his picture-tales — cherry legend → beckoning cat →
				// red gates → repeat.
				// 2026-07-29 (user): the stage DOORS open when a tale begins
				// and close when it ends, and each tale shows ITS OWN card —
				// per-tale talk sheets swap in when they land (§JP-TAKESHI-
				// DOORS); until then the landed single talk sheet plays.
				takeshi := n
				loadTaleTalk := func(file string) []npcFrame {
					if p := firstExisting(jpNPCDir + file); p != "" {
						return loadNPCGridConnected(game.renderer, p, 8, 1)
					}
					return nil
				}
				baseTalk := takeshi.talkGrid
				type kamishibaiTale struct {
					dialog []dialogEntry
					talk   []npcFrame
				}
				tales := []kamishibaiTale{
					{takeshiCherryTaleDialog, loadTaleTalk("npc_takeshi_talk_cherry.png")},
					{takeshiCatTaleDialog, loadTaleTalk("npc_takeshi_talk_cat.png")},
					{takeshiGatesTaleDialog, loadTaleTalk("npc_takeshi_talk_gates.png")},
				}
				if len(tales[0].talk) > 0 {
					takeshi.talkGrid = tales[0].talk
				}
				taleIdx := 0
				// 2026-08-01 (user #23): the doors STAY OPEN for the whole
				// dialog — after the open one-shot plays, hold its final
				// (doors-open) frame as the idle, so PP's lines between
				// Takeshi's talk beats don't flash the closed-door idle.
				baseIdle := takeshi.idleGrid
				takeshi.onDialogStart = func() {
					if f, ok := takeshi.oneShotAnims["open_doors"]; ok && len(f) > 0 {
						takeshi.idleGrid = f[len(f)-1:]
					}
					takeshi.playOneShotAnim("open_doors", 1.0)
				}
				takeshi.onDialogEnd = func() {
					takeshi.idleGrid = baseIdle
					takeshi.playOneShotAnim("close_doors", 1.0)
					taleIdx++
					next := tales[taleIdx%len(tales)]
					takeshi.dialog = next.dialog
					if len(next.talk) > 0 {
						takeshi.talkGrid = next.talk
					} else {
						takeshi.talkGrid = baseTalk
					}
				}
			}
		}
	}

	if grove, ok := g.sceneMgr.scenes["kyoto_temple"]; ok {
		for _, n := range grove.npcs {
			switch n.name {
			case "Oba-chan":
				oba := n
				// Multi-stage: (1) once PP has Kenji's clue (Voice Charm) she hands
				// over the kite-dropped Fire-Striker; (2) once PP carries the
				// blessed Offering Bowl she "follow me"s him - opening the grove.
				oba.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if game.inv.hasItem("Offering Bowl") && !game.vars.GetBool(ScopeGame, VarJpGroveRevealed) {
						return obachanLeadDialog, func() {
							game.vars.SetBool(ScopeGame, VarJpGroveRevealed, true)
							oba.dialog = obachanPostDialog
						}, nil
					}
					// 2026-07-26 PR: the old !jp_ramen_open guard was DEAD (the
					// flag is never set since the stall-open rework), so a chat
					// after the Hiro trade handed a SECOND striker that could
					// never leave the bag. Gate on the trade's outcome instead:
					// she only has the striker while PP hasn't earned the bowl.
					if game.inv.hasItem("Voice Charm") && !game.inv.hasItem("Fire-Striker") &&
						!game.inv.hasItem("Offering Bowl") &&
						!game.vars.GetBool(ScopeGame, VarJpGroveRevealed) {
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
					// (jp_tea_learned). One-shots no-op until their sheets land.
					// 2026-07-24 (user #42): the single 4.5s spin read as slow
					// motion and never HELD the kimono. New chain: quick spin
					// in → hold the kimono pose a beat → quick spin back.
					// 2026-07-24 (user, take 2): THREE dedicated beats/sheets —
					// spin INTO the costume, MODEL it, spin BACK. Each key is
					// its own sheet once §PP-KIMONO-SPLIT lands; until then a
					// slice of the combined spin sheet (see player.go).
					game.player.playOneShot("kimono_spin_in", 1.3, func() {
						game.player.playOneShot("kimono_model", 1.5, func() {
							game.player.playOneShot("kimono_spin_out", 0.8, nil)
						})
					})
					game.vars.SetBool(ScopeGame, VarJpTeaLearned, true)
					kiku.dialog = dresserPostDialog
				}
			}
		}
		// Tea-ceremony shelves in the flower store: a matcha tin + a shelf of
		// chawan (you're handed a random one).
		// 2026-07-26 PR #9: both moved from the bare upper-left wall onto the
		// STALL TABLE with the painted pots and bowls (right of Kiku, left of
		// the grove exit), so the pickups sit on something the player can see.
		grove.hotspots = append(grove.hotspots, hotspot{
			bounds: sdl.Rect{X: 1135, Y: 430, W: 80, H: 85}, name: "A tin of matcha",
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
			bounds: sdl.Rect{X: 1220, Y: 430, W: 75, H: 85}, name: "A shelf of tea bowls",
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
		// still heart (the tea ceremony).
		// 2026-08-01 (user #27): the exit is INVISIBLE until Oba-chan reveals
		// the path (visibleWhen — no hover name, no arrow cursor), and it
		// moved to the path opening around (373,350) as an UP arrow.
		grove.hotspots = append(grove.hotspots, hotspot{
			bounds: sdl.Rect{X: 293, Y: 270, W: 160, H: 160}, name: "Into the sakura grove", arrow: arrowUp,
			visibleWhen: func() bool {
				return game.vars.GetBool(ScopeGame, VarJpGroveRevealed)
			},
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
				game.sceneMgr.transitionTo("kyoto_sakura", game.player)
				return true
			},
		})
		// G9: the path UP to the temple tea-house should recede (shrink + drift
		// up) instead of marching PP off the top of the screen.
		// 2026-07-26 PR #11: PP first walks to the temple STAIRS' base
		// (foot ≈ 838,384; above the walk band, so unclamped), THEN recedes —
		// same beat as the street→flower-store stairs.
		for i := range grove.hotspots {
			if grove.hotspots[i].targetScene != "kyoto_teahouse" {
				continue
			}
			grove.hotspots[i].onInteract = func() bool {
				game.player.walkToAndDoUnclamped(838, 249, func() {
					game.player.playRecede(1.0, 0.5, 80, func() {
						game.sceneMgr.transitionTo("kyoto_teahouse", game.player)
					})
				})
				return true
			}
			break
		}
	}

	// Temple tea-house: share the whisked Matcha Bowl with the tea master →
	// jp_tea_done (the grove gate). No reward item; just the moment.
	if teahouse, ok := g.sceneMgr.scenes["kyoto_teahouse"]; ok {
		// 2026-08-01 (user #29): leaving the tea-house goes THROUGH the shoji
		// door — walk to the door base first, then recede through it. The
		// generic up-arrow path (walkToExit) marched PP straight up from
		// wherever he stood.
		for i := range teahouse.hotspots {
			if teahouse.hotspots[i].targetScene != "kyoto_temple" {
				continue
			}
			teahouse.hotspots[i].onInteract = func() bool {
				game.player.walkToAndDoUnclamped(193, 430, func() {
					game.player.playRecede(0.8, 0.5, 60, func() {
						game.sceneMgr.transitionTo("kyoto_temple", game.player)
					})
				})
				return true
			}
			break
		}
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
							// 2026-07-27 (user): the shared sip is VISIBLE now —
							// the sit_drink one-shot (sliced from PP_sit_talk's
							// drink arc) plays before PP rises. Graceful: a
							// missing slice falls straight through to standing.
							game.player.playOneShot("sit_drink", 2.0, func() {
								game.player.playOneShot("tea_stand", 1.8, func() {
									game.player.seated = false
									game.vars.SetBool(ScopeGame, VarJpTeaDone, true)
									tea.dialog = teaMasterPostDialog
									tea.altDialogFunc = nil
								})
							})
						})
					})
				}, &handOff{item: "Matcha Bowl"}
			}
		}
	}

	// Hidden sakura grove: the old tree is the pick-the-blossom payoff. Picking
	// gives the Pressed Sakura (the anchor) + fires Danny's foreshadow call.
	if sakura, ok := g.sceneMgr.scenes["kyoto_sakura"]; ok {
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
					// 2026-07-29 (user): a DEDICATED pick — PP reaches up to
					// the branch (§PP-PICK-SAKURA) — replaces the reused camp
					// flower grab once its sheet lands.
					pick := "grab_flower"
					if _, ok := game.player.oneShotAnims["pick_sakura"]; ok {
						pick = "pick_sakura"
					}
					game.player.playOneShot(pick, 0.9, func() {
						give("pressed_sakura")
						game.player.playOneShot("sakura_selfie", 2.4, func() {
							game.dialog.startDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "A real sakura blossom. Light as a breath. Lily will hold this and come back to herself."},
							})
							game.dialog.queueDialog(dannyPhoneCallDialog)
						})
					})
				})
				return true
			}
			break
		}
	}
}

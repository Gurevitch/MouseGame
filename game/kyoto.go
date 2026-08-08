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

// --- Gary (torii). 2026-08-01 (user): the chapter opening is now a
// MEET-BOTH gate — Gary's arrival chat is the pure book gag, Hiro's stall
// is simply dark; once BOTH are met, Gary's book produces the ramen-place
// tip, and only then Hiro names what he's missing (fire-striker + two
// groceries), which starts the item chain. Gary's book then knows where
// each grocery is sold. ---
var garyArrivalDialog = []dialogEntry{
	{speaker: "Gary", text: "PINK PANTHER! Can you BELIEVE it - KYOTO! I've dreamed about this since I was a kid. The temples, the gardens, the cherry blossoms..."},
	{speaker: "Gary", text: "I read EVERYTHING about this place - every shrine, every festival, every bowl of noodles. It's ALL here in my guidebook!"},
	{speaker: "Pink Panther", text: "...Gary. You're holding the book upside down."},
	{speaker: "Gary", text: "WHAT? ...Oh. Oh my. ...Ahh. MUCH better. Now it makes far more sense!"},
}

var garyRamenTipDialog = []dialogEntry{
	{speaker: "Gary", text: "You met the ramen cook down the street? THE temple-street ramen cook? Hold on - my book has a whole page on him!"},
	{speaker: "Gary", text: "Here: 'the old cherry in the hidden grove blooms only for an offering blessed over a sacred FIRE - and the ramen cook of the temple street keeps the last one.'"},
	{speaker: "Pink Panther", text: "So Hiro's dark little stall is the key to the whole grove. Back down the street, then. Thanks, Gary."},
}

// 2026-08-01 (user, round 2): the shopping list is a PING-PONG — Hiro asks
// for ONE grocery at a time, and Gary's book locates each one separately.
var garyKatsuobushiDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Gary - where does a panther buy KATSUOBUSHI in this town?"},
	{speaker: "Gary", text: "Katsuo-what? Oh - the smoked bonito flakes! The book has a market index! Let me check, let me check..."},
	{speaker: "Gary", text: "Here it is - KATSUOBUSHI: 'the flower-store lady by the temple stocks a little of everything.' The shop up the street stairs!"},
	{speaker: "Pink Panther", text: "The flower store it is. Thanks, Gary."},
}

var garyKombuDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "One more, Gary - KOMBU. Where do I find dried kelp?"},
	{speaker: "Gary", text: "Back to the index! Let me check, let me check... kimonos... koi... KOMBU - here it is!"},
	{speaker: "Gary", text: "'Hung to dry by the water.' The pond at the bottom of the ramen street, by the little bridge!"},
	{speaker: "Pink Panther", text: "Drying by the pond. This book is starting to earn its keep."},
}

// 2026-07-18 (user #45/#46): Gary is the chapter's CUSTOMS ORACLE — whenever
// PP doesn't know the local etiquette, Gary + guidebook has the answer. His
// repeat dialog advances with the chapter flags (see setupKyotoCallbacks).
// 2026-08-08 #18 (user): the "Kyoto is in PERU" gag read as nonsense —
// replaced with a coherent guidebook joke.
var garyTokyoPostDialog = []dialogEntry{
	{speaker: "Gary", text: "Need a custom checked? My book knows everything - I just found a whole chapter I'd been holding upside down. TWICE."},
}

var garyDressCodeDialog = []dialogEntry{
	{speaker: "Gary", text: "The book is very clear: you do NOT visit the Whispering Cherry in street clothes. It's a dress-code thing."},
	{speaker: "Gary", text: "The lady at the flower store dresses visitors properly. Tell her the tree sent you. Well - don't, actually. That sounds weird."},
}

var garyTeaCustomDialog = []dialogEntry{
	{speaker: "Gary", text: "Before the grove, the book says you need 'a still heart' - that's the tea ceremony. Matcha, a bowl, cool well water."},
	{speaker: "Gary", text: "The tea master at the temple tea-house hosts it. Very serene. I cried a little."},
}

// --- Hiro (street): stall dark on arrival; once Gary's tip lands he names
// everything he's missing — the kite-stolen fire-striker AND two groceries
// (katsuobushi + kombu). All three open the sacred hearth. ---
var hiroRamenDialog = []dialogEntry{
	{speaker: "Hiro", text: "Irasshaimase, panther-san... forgive the dark stall. No noodles today, I am afraid."},
	{speaker: "Pink Panther", text: "Rough morning?"},
	{speaker: "Hiro", text: "The roughest in years. But you did not come all this way to hear an old cook's troubles. Enjoy the street, my friend."},
}

var hiroMissingDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Still closed, Hiro? I was counting on your fire today."},
	{speaker: "Hiro", text: "Closed, closed, CLOSED - and I am racing the morning as it is, panther-san!"},
	{speaker: "Hiro", text: "This morning my KATSUOBUSHI delivery did not come - the first time in twenty years! No flakes, no broth, no ramen!"},
	{speaker: "Hiro", text: "And on top of it all my FIRE-STRIKER has vanished - it was by the stove at dawn, and now, poof! Gone!"},
	{speaker: "Pink Panther", text: "A flint and fish flakes. Keep your pots ready - I'll do the running."},
}

// Delivering the katsuobushi IS the second ask (the ping-pong middle beat).
var hiroKatsuTradeDialog = []dialogEntry{
	{speaker: "Hiro", text: "KATSUOBUSHI! Fresh flakes, the good kind - into the pantry they go."},
	{speaker: "Hiro", text: "Now all that is left for the ramen is KOMBU - dried kelp, for the deep of the broth."},
	// 2026-08-08 #25 (user): PP's "One more trip to Gary's book, then." line removed.
}

var hiroRamenPostDialog = []dialogEntry{
	{speaker: "Hiro", text: "Still closed, panther-san - my striker is still missing, the pantry still waits for katsuobushi, and the morning will not wait at all!"},
}

var hiroPostKatsuDialog = []dialogEntry{
	{speaker: "Hiro", text: "The flakes rest in my pantry, panther-san. Now KOMBU - and my lost striker - and the stall opens again."},
}

// 2026-08-08 #29: Hiro accepts the kelp ON ITS OWN (PP from behind at the
// counter, the bakery framing) instead of demanding kelp + striker together.
var hiroKombuTradeDialog = []dialogEntry{
	{speaker: "Hiro", text: "KOMBU! Thick and dark - the pond dries the best kelp in Kyoto. Into the pantry, next to the flakes."},
	{speaker: "Hiro", text: "The broth is READY to be born, panther-san. Only my fire-striker is still out there - a kite took it, shiny thief."},
}

var hiroPostKombuDialog = []dialogEntry{
	{speaker: "Hiro", text: "The pantry is full, the hearth is cold. Find my STRIKER, panther-san - ask around, someone always sees the kites."},
}

// The cherry agenda is PP's (from Gary's tip), not Hiro's — after the hearth
// lights, PP is the one who asks for the blessed bowl.
var hiroOpenDialog = []dialogEntry{
	{speaker: "Hiro", text: "My STRIKER! With the flakes and the kelp already in the pantry, the broth is complete. Stand back - "},
	{speaker: "Pink Panther", text: "Beautiful. Now, about that favor - one bowl, blessed in the first flame. It's for the old cherry in the grove."},
	{speaker: "Hiro", text: "For the one who saved my morning? Gladly. Here - carry it with respect."},
}

// --- Kenji (torii): saw where the kite dropped it; needs well-water first ---
// Kenji is the WITNESS: Hiro only knows his striker vanished; Kenji is the
// one who saw the kite take it (2026-08-01 — the theft claim moved here,
// where it makes sense). Before Hiro's ask exists, Kenji is pure flavor —
// the witness beat can't fire until PP actually knows a striker is missing.
var kenjiIdleDialog = []dialogEntry{
	{speaker: "Kenji", text: "Please - do not nudge the table. Ink and stillness, panther-san. A calligrapher needs both, and the street offers neither."},
}

var kenjiNeedWaterDialog = []dialogEntry{
	{speaker: "Kenji", text: "Water from the temple well, panther-san - down the ramen street, at the foot of the stairs. Then the ink flows, and so does my memory."},
}

var kenjiStudentDialog = []dialogEntry{
	{speaker: "Kenji", text: "Please - do not nudge the table... oh, the panther. You have the look of a man who lost something small and shiny."},
	{speaker: "Pink Panther", text: "Hiro's fire-striker. It vanished from his stall at dawn."},
	{speaker: "Kenji", text: "Not vanished - TAKEN. A kite swooped the street this morning and flew off with a glint in its grip. I saw exactly where it came down..."},
	{speaker: "Kenji", text: "...but my ink has dried to dust and I cannot think with a dry brush. Bring me water from the temple well - down in the ramen street, at the foot of the stairs - and I will tell you."},
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
	{speaker: "Takeshi", text: "Ah, a customer for my paper theater! Sit, sit. Today's tale... The Whispering Cherry."},
	{speaker: "Takeshi", text: "Deep in a hidden grove stands a tree older than the temple bells. Most years it will not bloom at all. Gold cannot buy its flowers. Princes have tried."},
	{speaker: "Takeshi", text: "But bring it a true gift - something made with care, given with a whole heart - and the old tree wakes... and WHISPERS back."},
	{speaker: "Pink Panther", text: "A tree that answers. I know a little girl who'd love that story."},
	{speaker: "Takeshi", text: "Heh. It is no story, pink one. Ask any old man in this street."},
}

var takeshiCatTaleDialog = []dialogEntry{
	{speaker: "Takeshi", text: "The Beckoning Cat! Long ago, a poor shopkeeper shared his last rice ball with a stray cat..."},
	{speaker: "Takeshi", text: "The cat sat at his door and raised one paw - and travelers followed it in, curious. The shop never stood empty again. We call it maneki-neko, the beckoning cat."},
	{speaker: "Takeshi", text: "Look up - the stray on Hiro's roof. Left paw tucked, waiting. When that stall opens again, watch her raise it. Luck knows where soup is coming."},
}

var takeshiGatesTaleDialog = []dialogEntry{
	{speaker: "Takeshi", text: "The Red Gates. A torii is a doorway with no wall - because the wall it passes through cannot be seen."},
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

// 2026-08-07 #31: the katsuobushi is now a DIALOG hand-over from Oba-chan
// (Gary's tip says "the flower-store lady stocks it"), replacing the silent
// stall-table hotspot whose click rect sat entirely INSIDE Kiku's box — NPC
// clicks win over hotspots, so the packet was unreachable and the whole
// hearth chain dead-ended.
var obachanKatsuobushiDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Oba-chan - Hiro's pantry needs KATSUOBUSHI, and a certain guidebook says you stock a little of everything."},
	{speaker: "Oba-chan", text: "The book is right for once. Smoked bonito, shaved this morning - the last packet on the shelf."},
	{speaker: "Oba-chan", text: "For Hiro's broth? Then take it as a neighbour's favor. Tell him Oba-chan keeps the street fed TOO."},
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
	{speaker: "Pink Panther", text: "The old tree's branches are bare. Oba-chan said it blooms for an offering blessed over a sacred fire - I should bring one first."},
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
// 2026-08-08 #21 (user: "does the clothes change make sense?"): rewritten so
// the kimono has a REASON before it happens — she names the grove's dress
// code first (echoing Gary's tip), teaches the tea path, and only then
// "SPIN!" as the LAST line, so the spin one-shot (played at dialog end)
// lands right on cue instead of PP saying "I'm dressed" before any spin.
var dresserDialog = []dialogEntry{
	{speaker: "Kiku", text: "Ara... a pink panther, in MY shop. And WHERE do you think you are going, dressed like ZAT?"},
	{speaker: "Pink Panther", text: "To the old cherry grove - once it opens its heart to me, anyway."},
	{speaker: "Kiku", text: "The WHISPERING CHERRY? It does not receive guests in street fur. There is a dress code, panther-san - and I am its keeper."},
	{speaker: "Kiku", text: "And silk alone will not do: the grove refuses a RESTLESS heart. Take matcha from my stall, draw cool water at the street well, whisk it - then kneel with the tea master up in the temple house."},
	{speaker: "Kiku", text: "A proper chawan? You will have one of MINE - ask me when the matcha is in your paw. Now hold still - SPIN!"},
}

var dresserPostDialog = []dialogEntry{
	{speaker: "Kiku", text: "Matcha from my stall, a chawan from ME, well-water from the street - then the tea master in the temple house. Go, go, panther-san!"},
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
	// 2026-08-07 #1: when the counter sheets exist (setupKyotoCallbacks swaps
	// them in before anything draws), skip the legacy full-body sheets
	// entirely — the old npc_hiro_talk has merged figure pairs that spammed
	// dropMalformedFrames boot warnings for frames the game never showed.
	var idle, talk []npcFrame
	if firstExisting(jpNPCDir+"npc_hiro_counter_idle.png") != "" {
		// Equal-cell load: the counter idle's poses touch (5 runs for 6, one
		// ~500px merged blob), so the gap cutter moved its boundaries per
		// frame — the office-Higgins jitter class.
		idle = loadNPCGridEqualConnectedTol(renderer, jpNPCDir+"npc_hiro_counter_idle.png", 6, 1, 4, 24)
		talk = idle
		if p := firstExisting(jpNPCDir + "npc_hiro_counter_talk.png"); p != "" {
			if t := loadNPCGridConnected(renderer, p, 6, 1); len(t) > 0 {
				talk = t
			}
		}
	} else {
		idle, talk = loadJapanNPC(renderer,
			[]string{jpNPCDir + "npc_hiro_idle.png", jpNPCDir + "ramen_seller.png"},
			[]string{jpNPCDir + "npc_hiro_talk.png", jpNPCDir + "ramen_seller_talk.png"})
	}
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
		dialog:         garyArrivalDialog,
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
	// 2026-08-07 #26 (user): TWO idles, both on the same seated-calligrapher
	// design — a SLOW nature-gazing base idle (§KENJI-IDLE-NATURE) punctuated
	// by a drawing loop (§KENJI-IDLE-DRAWING, the alt-idle mechanism).
	// 2026-08-08 #17 (user: "kenji is jittering"): ALL his sheets (incl. the
	// freshly landed nature + drawing) carry one merged figure pair, so the
	// gap cutter waist-split a different boundary per frame. Their figures
	// sit ≤13px off the 192 grid → EQUAL-cell loads (the counter-Hiro /
	// office-Higgins fix) hold the boundaries still.
	kenjiSheet := func(cands ...string) []npcFrame {
		if p := firstExisting(cands...); p != "" {
			return loadNPCGridEqualConnectedTol(renderer, p, 8, 1, 4, 24)
		}
		return nil
	}
	idle := kenjiSheet(jpNPCDir+"npc_kenji_idle_nature.png", jpNPCDir+"npc_kenji_idle.png")
	if len(idle) == 0 {
		idle = loadNPCGridRow(renderer, jpFbkVendor8x2, 8, 2, 0)
	}
	talk := kenjiSheet(jpNPCDir + "npc_kenji_talk.png")
	if len(talk) == 0 {
		talk = idle
	}
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
		// #26: the gaze reads contemplative — slow the idle cadence well
		// under the derived talk*2 (Pierre precedent 0.40).
		idleFrameSpeed: 0.5,
	}
	// §KENJI-GIVE-CHARM (landed 2026-07-18): his brush-the-charm one-shot.
	if p := firstExisting(jpNPCDir + "npc_kenji_give_charm.png"); p != "" {
		if f := loadNPCGridConnected(renderer, p, 8, 1); len(f) > 0 {
			n.oneShotAnims = map[string][]npcFrame{"give": f}
		}
	}
	// #26: the drawing punctuation — after ~6 idle seconds he brushes a few
	// strokes, then returns to gazing (one full alt cycle per fire).
	// #17: equal-cell like the base idle (its pair 7+8 also merges).
	if f := kenjiSheet(jpNPCDir + "npc_kenji_idle_drawing.png"); len(f) > 0 {
		n.altIdleGrid = f
		n.altIdleAfterSec = 6.0
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
	// 2026-08-08 #24 (user: idle still shows two frames at once): the idle
	// sheet is measured PHASE-SHIFTED +19..40px off the 256 grid with dirty
	// valleys — equal cells slice a sliver of the previous pose into every
	// frame, and the gap cutter only finds 3 runs. NO cutter can win, so the
	// clean TALK sheet (6/6 gap-detected runs) doubles as the idle until the
	// §TEA-MASTER-IDLE-v3 regen lands at its own preferred path.
	idle := []npcFrame{}
	if p := firstExisting(jpNPCDir + "npc_tea_master_idle_v3.png"); p != "" {
		idle = loadNPCGridConnected(renderer, p, 6, 1)
	}
	talk := []npcFrame{}
	if p := firstExisting(jpNPCDir + "npc_tea_master_talk.png"); p != "" {
		talk = loadNPCGridConnected(renderer, p, 6, 1)
	}
	if len(idle) == 0 {
		idle = talk
	}
	if len(idle) == 0 {
		idle = loadNPCGridRow(renderer, jpFbkVendor8x2, 8, 2, 0)
	}
	if len(talk) == 0 {
		talk = idle
	}
	return &npc{
		idleGrid: idle, talkGrid: talk,
		// 2026-07-24 (user #43): interim size bump (H 248 → 280, foot kept
		// at 620) — her placeholder sheets draw the figure tiny in-cell.
		// The real fix stays the queued §TEA-MASTER-BLUE-v2 re-roll.
		// 2026-07-26 PR #6: nudged a little left (X 640 → 615) per user.
		bounds: sdl.Rect{X: 615, Y: 340, W: 140, H: 280},
		name:   "Tea Master",
		dialog: teaMasterDialog,
		// 2026-08-07 #29 (user: "she need to act slowly"): talk 0.12 → 0.20
		// and an explicit slow idle (unset derived talk*2 = 0.24 → 0.5).
		talkFrameSpeed: 0.20,
		idleFrameSpeed: 0.5,
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
	n := &npc{
		idleGrid: idle, talkGrid: talk,
		bounds:         sdl.Rect{X: 380, Y: 370, W: 140, H: 250},
		name:           "Oba-chan",
		dialog:         obachanDialog,
		talkFrameSpeed: 0.12,
	}
	// 2026-08-08 #20 (user): she visibly HANDS things over (the striker, the
	// katsuobushi packet) — graceful until §OBACHAN-GIVE lands.
	if p := firstExisting(jpNPCDir + "npc_obachan_give.png"); p != "" {
		if f := loadNPCGridConnected(renderer, p, 6, 1); len(f) > 0 {
			n.oneShotAnims = map[string][]npcFrame{"give": f}
		}
	}
	return n
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
	// 2026-08-08 #15 (user): a little smaller — 0.3 → 0.26.
	addJapanProp("kyoto_street", "street_cat.png", 866, 309, 0.26, 0.7, 8)
	// the thieving KITE (tobi) perched up on the temple roof — the bird
	// Hiro and Kenji blame for the fire-striker (§JP-KITE-SHINY).
	// 2026-07-26 PR #5: foot moved to (657,217) per user.
	addJapanProp("kyoto_temple", "kite_shiny.png", 657, 217, 0.32, 0.6, 8)
	// stone lantern with incense wisps by the flower store
	// 2026-07-26 PR #4: foot moved to (334,473) per user.
	addJapanProp("kyoto_temple", "incense_lantern.png", 334, 473, 0.42, 0.5, 8)
	// steaming kettle beside the tea master. 2026-07-26 PR #7: tol 24 — the
	// tol-8 global key left a blue AA seam pocket on the kettle.
	// 2026-08-07 #29: kettle moved (850,610) → (780,530) per user.
	addJapanProp("kyoto_teahouse", "teahouse_kettle.png", 780, 530, 0.34, 0.45, 24)

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
			// 2026-08-01 (coherence pass): Kenji is FLAVOR until PP actually
			// knows a striker is missing (Hiro's ask) — his default witness
			// dialog used to spill the kite story unprompted.
			kenji.dialog = kenjiIdleDialog
			kenji.onDialogEnd = func() {
				if len(kenji.dialog) > 0 && &kenji.dialog[0] == &kenjiStudentDialog[0] {
					// The witness beat just played — he's now waiting on water.
					game.vars.SetBool(ScopeGame, "jp_kenji_water_asked", true)
					kenji.dialog = kenjiNeedWaterDialog
				}
			}
			kenji.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
				// Witness beat: only once Hiro's ask exists and it hasn't played.
				if game.vars.GetBool(ScopeGame, "jp_hiro_ask_done") &&
					!game.vars.GetBool(ScopeGame, "jp_kenji_water_asked") &&
					!game.inv.hasItem("Voice Charm") {
					return kenjiStudentDialog, func() {
						game.vars.SetBool(ScopeGame, "jp_kenji_water_asked", true)
						kenji.dialog = kenjiNeedWaterDialog
					}, nil
				}
				// Bring Kenji well-water → he points to Oba-chan's eaves + gives the Voice Charm.
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
			// 2026-08-08 #26 (user: "I click him multiple times to reach the
			// relevant line"): the customs-oracle ladder used to be selected
			// in onDialogEnd — one step BEHIND the state, so every chapter
			// advance cost a wasted chat replaying the stale line. The switch
			// now runs at CLICK time (onDialogStart fires in startNPCDialog
			// before the first line shows — the Takeshi doors precedent), so
			// the first click after any state change already speaks the right
			// line. Skipped for his very first chat (the arrival dialog).
			garySelectDialog := func() {
				if !game.vars.GetBool(ScopeGame, "jp_met_gary") {
					return // arrival dialog plays as authored
				}
				switch {
				case game.vars.GetBool(ScopeGame, VarJpGroveRevealed) && !game.vars.GetBool(ScopeGame, VarJpTeaDone):
					gary.dialog = garyTeaCustomDialog
				case game.inv.hasItem("Offering Bowl") && !game.vars.GetBool(ScopeGame, VarJpTeaLearned):
					gary.dialog = garyDressCodeDialog
				case game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") &&
					!game.vars.GetBool(ScopeGame, "jp_kombu_given") &&
					!game.inv.hasItem("Kombu") && !game.inv.hasItem("Offering Bowl"):
					gary.dialog = garyKombuDialog
				case game.vars.GetBool(ScopeGame, "jp_hiro_ask_done") &&
					!game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") &&
					!game.inv.hasItem("Katsuobushi"):
					gary.dialog = garyKatsuobushiDialog
				case game.vars.GetBool(ScopeGame, "jp_met_hiro") && !game.vars.GetBool(ScopeGame, "jp_gary_tip_done"):
					gary.dialog = garyRamenTipDialog
				default:
					gary.dialog = garyTokyoPostDialog
				}
			}
			gary.onDialogStart = garySelectDialog
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
				// 2026-08-01 (user): the opening is a meet-both gate — record
				// that Gary was met, and note when the tip was actually HEARD
				// (his dialog was the tip when this chat ended).
				game.vars.SetBool(ScopeGame, "jp_met_gary", true)
				if len(gary.dialog) > 0 && &gary.dialog[0] == &garyRamenTipDialog[0] {
					game.vars.SetBool(ScopeGame, "jp_gary_tip_done", true)
				}
				// #26: keep the resting dialog fresh for save/load too.
				garySelectDialog()
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
				// 2026-08-07 #1: EQUAL-cell load — the counter idle's poses
				// touch, so the gap cutter jittered its boundaries per frame.
				if idle := loadNPCGridEqualConnectedTol(g.renderer, counterIdle, 6, 1, 4, 24); len(idle) > 0 {
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
					// 2026-08-08 #20: his counter TAKE (reaches over, receives,
					// stows under the counter) — the grocery hand-overs skip
					// the take entirely until this lands (§HIRO-COUNTER-RECEIVE),
					// because the only fallback was the broth-ladle give.
					if p := firstExisting(jpNPCDir + "npc_hiro_counter_receive.png"); p != "" {
						if take := loadNPCGridConnected(g.renderer, p, 6, 1); len(take) > 0 {
							if n.oneShotAnims == nil {
								n.oneShotAnims = map[string][]npcFrame{}
							}
							n.oneShotAnims["receive_katsuobushi"] = take
							n.oneShotAnims["receive_kombu"] = take
						}
					}
				}
			}
			break
		}

		// 2026-08-01 (user): Hiro's grocery #2 — kombu drying by the water at
		// the pond bridge (bottom-right of the street). Gated on Hiro's ask
		// like the katsuobushi beat.
		// 2026-08-08 #28 (user): reworked from a silent hotspot to a hidden
		// FLOOR ITEM — grab cursor for free — whose stand mark is ON THE
		// BRIDGE DECK (the old hotspot walk aimed at the rect centre, foot
		// 825: PP marched into the pond, feet off-screen; the water is now
		// foot-blocked in kyoto_street.json). The pick plays a kneel-and-pull
		// beat: §PP-KNEEL-SCOOP once it lands, the crouch-grab until then.
		street.floorItems = append(street.floorItems, &floorItem{
			bounds:  sdl.Rect{X: 1150, Y: 650, W: 130, H: 80},
			name:    "Kombu",
			visible: true,
			hidden:  true, // the drying rack is painted in the BG
			// PP CENTER mark: (1195, 600) → foot (1195, 735), on the planks.
			standXOverride: 1195,
			standYOverride: 600,
			onPickup: func() {
				// Gated on the SECOND ask — Hiro only names kombu after the
				// katsuobushi is in his pantry (user: ping-pong, one at a time).
				if !game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") {
					game.player.dir = dirDown
					game.player.facingLeft = false
					game.player.state = stateTalking
					game.dialog.startDialogWithCallback([]dialogEntry{
						{speaker: "Pink Panther", text: "Long dark ribbons drying by the water... seaweed? Someone's dinner, not mine."},
					}, func() { game.player.state = stateIdle })
					return
				}
				if game.inv.hasItem("Kombu") || game.vars.GetBool(ScopeGame, "jp_kombu_given") ||
					game.inv.hasItem("Offering Bowl") {
					game.dialog.startDialog([]dialogEntry{{speaker: "Pink Panther", text: "I've got the kelp already."}})
					return
				}
				// Kneel at the plank edge, dip, pull the kelp out.
				scoop := "grab_flower" // crouch-and-reach, the closest landed beat
				if game.player.hasOneShot("kneel_scoop") {
					scoop = "kneel_scoop" // §PP-KNEEL-SCOOP
				}
				game.player.facingLeft = true
				game.player.dir = dirLeft
				game.player.playOneShot(scoop, 1.2, func() {
					give("kombu")
					game.player.dir = dirDown
					game.player.facingLeft = false
					game.player.state = stateTalking
					game.dialog.startDialogWithCallback([]dialogEntry{
						{speaker: "Pink Panther", text: "Kombu - dried kelp, hung by the water just like Gary's book said. That's the other half of Hiro's broth."},
					}, func() { game.player.state = stateIdle })
				})
			},
		})

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
				// 2026-08-01 (coherence pass): no reason to draw water before
				// Kenji has actually asked for it.
				if !game.vars.GetBool(ScopeGame, "jp_kenji_water_asked") {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "A cold stone well. Good clear water... but I've got no cup and no reason to draw one."},
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
				// 2026-08-01 (user): opening rework. First chat = the stall is
				// just dark (no ask). Once Gary's tip lands, the next chat has
				// Hiro name EVERYTHING he's missing — the kite-stolen striker
				// plus katsuobushi + kombu — which starts the item chain.
				hiro.onDialogEnd = func() {
					game.vars.SetBool(ScopeGame, "jp_met_hiro", true)
					if len(hiro.dialog) > 0 && &hiro.dialog[0] == &hiroMissingDialog[0] {
						game.vars.SetBool(ScopeGame, "jp_hiro_ask_done", true)
					}
					if game.inv.hasItem("Offering Bowl") {
						return // post-trade dialog was set by the trade callback
					}
					switch {
					case game.vars.GetBool(ScopeGame, "jp_katsuobushi_given"):
						hiro.dialog = hiroPostKatsuDialog
					case game.vars.GetBool(ScopeGame, "jp_hiro_ask_done"):
						hiro.dialog = hiroRamenPostDialog
					case game.vars.GetBool(ScopeGame, "jp_gary_tip_done"):
						hiro.dialog = hiroMissingDialog
					default:
						hiro.dialog = hiroRamenDialog
					}
				}
				// 2026-08-01 (user, round 2): PING-PONG trades. Stage 1 — PP
				// delivers the katsuobushi and Hiro immediately names the
				// SECOND grocery (kombu). Stage 2 — kombu + the fire-striker
				// light the sacred hearth → the Offering Bowl.
				hiro.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if game.inv.hasItem("Katsuobushi") && !game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") {
						return hiroKatsuTradeDialog, func() {
							game.inv.removeItem("Katsuobushi")
							game.vars.SetBool(ScopeGame, "jp_katsuobushi_given", true)
							hiro.dialog = hiroPostKatsuDialog
							// 2026-08-08 #20: skip his take-reach only while
							// the dedicated counter receive hasn't landed —
							// the fallback played his broth-LADLE give as
							// the take, which read wrong.
						}, &handOff{item: "Katsuobushi", back: true,
							skipNPCTake: !hiro.hasOneShotAnim("receive_katsuobushi")}
					}
					// 2026-08-08 #29 (user: "I can't bring the kombu to
					// Hiro"): there was NO kombu-alone branch — carrying the
					// kelp without the striker hit nothing and replayed a
					// generic reminder. Hiro now TAKES the kelp on its own
					// (PP from behind, the bakery-counter framing) and his
					// post dialog names what's still missing.
					if game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") &&
						game.inv.hasItem("Kombu") &&
						!game.vars.GetBool(ScopeGame, "jp_kombu_given") {
						return hiroKombuTradeDialog, func() {
							game.inv.removeItem("Kombu")
							game.vars.SetBool(ScopeGame, "jp_kombu_given", true)
							hiro.dialog = hiroPostKombuDialog
						}, &handOff{item: "Kombu", back: true,
							skipNPCTake: !hiro.hasOneShotAnim("receive_kombu")}
					}
					// The blessed bowl: pantry stocked (both groceries GIVEN,
					// no longer carried) + the striker in paw.
					if game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") &&
						game.vars.GetBool(ScopeGame, "jp_kombu_given") &&
						game.inv.hasItem("Fire-Striker") &&
						!game.inv.hasItem("Offering Bowl") {
						return hiroOpenDialog, func() {
							game.inv.removeItem("Fire-Striker")
							give("offering_bowl")
							hiro.dialog = []dialogEntry{
								{speaker: "Hiro", text: "Take the blessed bowl to the old tree, panther-san - and come back for noodles when your heart is light."},
							}
							hiro.altDialogFunc = nil
						}, &handOff{item: "Fire-Striker", returnItem: "Offering Bowl", back: true}
					}
					return nil, nil, nil
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
					// 2026-08-07 #31: the katsuobushi hand-over (was an
					// unreachable stall-table hotspot buried under Kiku's
					// click box). Gated on Hiro having NAMED the list, like
					// the old hotspot; one packet only.
					if game.vars.GetBool(ScopeGame, "jp_hiro_ask_done") &&
						!game.inv.hasItem("Katsuobushi") &&
						!game.vars.GetBool(ScopeGame, "jp_katsuobushi_given") {
						return obachanKatsuobushiDialog, func() {
							give("katsuobushi")
						}, &handOff{returnItem: "Katsuobushi"}
					}
					return nil, nil, nil
				}
			case "Kiku":
				kiku := n
				// 2026-08-08 #27 (user: two pickups in one spot "looks like
				// stealing"): the TEA BOWL is now KIKU's hand-over — she
				// lends one of her own chawan once the ceremony is taught
				// and the matcha is in the bag. The shelf hotspot is gone.
				kiku.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if game.vars.GetBool(ScopeGame, VarJpTeaLearned) &&
						game.inv.hasItem("Matcha") &&
						!game.inv.hasItem("Tea Bowl") && !game.inv.hasItem("Matcha Bowl") &&
						!game.vars.GetBool(ScopeGame, VarJpTeaDone) {
						bowl := teaBowlNames[rand.Intn(len(teaBowlNames))]
						return []dialogEntry{
							{speaker: "Kiku", text: "Matcha in paw? Good. And no, you will NOT paw through shelves for a bowl like a market cat."},
							{speaker: "Kiku", text: "Here - the " + bowl + " one, from my own set. It has seen a hundred ceremonies; do not let it see the ground."},
						}, func() { give("tea_bowl") }, &handOff{returnItem: "Tea Bowl"}
					}
					return nil, nil, nil
				}
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
					// 2026-08-07 #32 (user): the spin chain + teaching fire
					// ONLY on the FIRST chat — jp_tea_learned itself is the
					// guard, so it's save-safe (a reloaded game skips it too).
					if !game.vars.GetBool(ScopeGame, VarJpTeaLearned) {
						game.player.playOneShot("kimono_spin_in", 1.3, func() {
							game.player.playOneShot("kimono_model", 1.5, func() {
								game.player.playOneShot("kimono_spin_out", 0.8, nil)
							})
						})
						game.vars.SetBool(ScopeGame, VarJpTeaLearned, true)
					}
					kiku.dialog = dresserPostDialog
				}
			}
		}
		// 2026-08-08 #27 (user): the two shelf HOTSPOTS (matcha + chawan, 5px
		// apart, no cursor, no animation — "looks like stealing") are GONE.
		// The TEA BOWL is Kiku's dialog hand-over now (see her altDialog),
		// and the MATCHA is a hidden floorItem on the painted tin: the grab
		// cursor, the beside-the-item walk, and PP's grab one-shot all come
		// free with the floor-item path (the rolling-pin-basket pattern).
		grove.floorItems = append(grove.floorItems, &floorItem{
			bounds:  sdl.Rect{X: 1135, Y: 430, W: 80, H: 85},
			name:    "Matcha",
			visible: true,
			hidden:  true, // painted tin IS the art; cursor-only presence
			standRight: false,
			standGapX:  20,
			onPickup: func() {
				if !game.vars.GetBool(ScopeGame, VarJpTeaLearned) {
					game.player.dir = dirDown
					game.player.facingLeft = false
					game.player.state = stateTalking
					game.dialog.startDialogWithCallback([]dialogEntry{
						{speaker: "Pink Panther", text: "Pretty green powder... but I wouldn't know what to do with it. Maybe the kimono lady can teach me."},
					}, func() { game.player.state = stateIdle })
					return
				}
				if game.inv.hasItem("Matcha") || game.vars.GetBool(ScopeGame, VarJpTeaDone) {
					game.dialog.startDialog([]dialogEntry{{speaker: "Pink Panther", text: "I've already got the matcha."}})
					return
				}
				// Waist-high stall reach: the basket-grab sheet is the closest
				// existing beat ("grab" is a registered-nowhere no-op).
				matchaGrab := "grab"
				if game.player.hasOneShot("grab_rolling_pin") {
					matchaGrab = "grab_rolling_pin"
				}
				game.player.playOneShot(matchaGrab, 0.9, func() {
					give("matcha")
					game.player.dir = dirDown
					game.player.facingLeft = false
					game.player.state = stateTalking
					game.dialog.startDialogWithCallback([]dialogEntry{
						{speaker: "Pink Panther", text: "Bright green matcha powder. Just like Kiku said - the tea master will want this."},
					}, func() { game.player.state = stateIdle })
				})
			},
		})
		// The katsuobushi-packet hotspot (2026-08-01) was REMOVED 2026-08-07
		// #31: its click rect {1050,435,70,75} sat entirely inside Kiku's
		// bounds {980,360,150,250}, and NPC clicks beat hotspots — the packet
		// could never be clicked, dead-ending the whole hearth chain. The
		// hand-over is now Oba-chan's dialog beat (obachanKatsuobushiDialog),
		// which is also what Gary's tip promised ("the flower-store LADY
		// stocks it").
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
		// (above the walk band, so unclamped), THEN recedes — same beat as
		// the street→flower-store stairs.
		// 2026-08-07 #28: recede drift 80 → 40 (the "roof climb").
		// 2026-08-08 #22 (user): the shrink starts LOWER — anchor foot
		// (844,353) → (852,422), the painted stair BASE, matching the
		// re-centred arrow above it.
		for i := range grove.hotspots {
			if grove.hotspots[i].targetScene != "kyoto_teahouse" {
				continue
			}
			grove.hotspots[i].onInteract = func() bool {
				game.player.walkToAndDoUnclamped(852, 287, func() {
					game.player.playRecede(1.0, 0.5, 40, func() {
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

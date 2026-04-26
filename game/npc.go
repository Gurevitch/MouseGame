package game

import (
	"math"
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type npcFrame struct {
	tex *sdl.Texture
	w   int32
	h   int32
	// src is the source rectangle inside tex. nil means "draw the whole
	// texture" (legacy per-frame loaders produce one texture per frame so
	// they leave this nil). Atlas-backed frames share one texture and set
	// src to the frame's rect within the atlas.
	src *sdl.Rect
}

type npc struct {
	bounds    sdl.Rect
	dialog    []dialogEntry
	name      string
	bobTimer  float64
	bobAmount float64
	flipped   bool
	// preTalkFlipped snapshots n.flipped before a dialog starts so
	// startNPCDialog can flip the NPC to face PP and then the wrapCb
	// can restore the authored pose when the conversation ends. Without
	// this, NPCs like Danny (authored flipped=true so he faces the camp
	// center) would stay stuck in whatever direction they were last
	// turned during talk.
	preTalkFlipped bool
	hovered        bool
	itemMatch      bool
	elevated       bool
	silent         bool
	// hidden skips the draw pass for this NPC. Used for story-timed
	// arrivals (e.g. Higgins appearing next to Lily only after her shy
	// dialog) so the NPC can sit in the scene list from load without
	// being visible or clickable until his cue.
	hidden  bool
	groupID string

	dialogDone    bool
	onDialogEnd   func()
	altDialogFunc func() ([]dialogEntry, func())
	// altDialogRequiresHeld gates altDialogFunc behind the player
	// actively carrying a specific item (altDialogRequiresItem). Without
	// this, the alt dialog would fire on any click once its condition
	// passed — breaking "give-item" flows where the player needs to
	// explicitly offer the item (e.g. Lily's flower). The default is
	// off (false) so existing altDialogFunc attachments keep working.
	altDialogRequiresHeld bool
	altDialogRequiresItem string
	// hintState is a small per-NPC dialog progression counter. Lily uses
	// 0 = has not been spoken to, 1 = shy dialog played (waiting for
	// flower), 2 = flower given. Storing this on the NPC instead of a
	// closure variable keeps the state deterministic across scene
	// re-entry and save/load (closures would reset back to zero when
	// setupCampCallbacks ran again).
	hintState int
	sm    *npcStateMachine   // optional state machine (named states: default/post/strange/post_strange)
	rules []InteractionRule  // optional rule list for data-driven interactions (see npc_rules.go)
	// game is a back-reference set by spawnNPCs so rule-driven NPCs can
	// call g.fireTrigger without threading *Game through every handler.
	// Not set for NPCs built via legacy callbacks — the rules slice stays
	// empty for those and fireTrigger is a no-op.
	game *Game

	idleGrid       []npcFrame
	talkGrid       []npcFrame
	talkFrameSpeed float64
	curFrame       int
	frameTimer     float64
	idleCurFrame   int
	idleFrameTimer float64
	animState      int

	strangeIdle []npcFrame
	strangeTalk []npcFrame
	normalIdle  []npcFrame
	normalTalk  []npcFrame
	isStrange   bool
	// strangeTalkFrameSpeed slows the talk animation while the NPC is in
	// strange state (Marcus's freakout looked too flickery at the default
	// 0.10 s/frame). 0 = inherit talkFrameSpeed unchanged.
	strangeTalkFrameSpeed float64
}

func (n *npc) setStrange(strange bool) {
	if strange == n.isStrange {
		return
	}
	n.isStrange = strange
	if strange && len(n.strangeIdle) > 0 {
		n.normalIdle = n.idleGrid
		n.normalTalk = n.talkGrid
		n.idleGrid = n.strangeIdle
		n.talkGrid = n.strangeTalk
	} else if !strange && len(n.normalIdle) > 0 {
		n.idleGrid = n.normalIdle
		n.talkGrid = n.normalTalk
	}
	n.curFrame = 0
	n.frameTimer = 0
	n.idleCurFrame = 0
	n.idleFrameTimer = 0
	n.animState = npcAnimIdle
}


// ===== Camp Chilly Wa Wa NPCs =====

// npcSpriteInset matches the trim used for player sheets. Keeps cell seams
// from leaking into the NPC idle/talk animations.
const npcSpriteInset = 3

// framesFromGrid flattens a rows x cols GridFrame matrix into an
// npcFrame list and trims trailing frames whose texture is nil (loader
// bailed on a missing cell). We do not attempt to trim "empty" frames
// whose texture is valid but fully transparent — measuring that per
// frame would require a GPU readback, and authored sheets that have
// 5-7 real cells in a row of 8 usually keep the last slot either fully
// transparent or a duplicate of the last pose, neither of which hurts
// the idle loop as much as getting the grid geometry wrong.
func framesFromGrid(grid [][]engine.GridFrame, cols, rows int) []npcFrame {
	var frames []npcFrame
	for r := 0; r < rows && r < len(grid); r++ {
		for c := 0; c < cols && c < len(grid[r]); c++ {
			gf := grid[r][c]
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}
	for len(frames) > 1 && frames[len(frames)-1].tex == nil {
		frames = frames[:len(frames)-1]
	}
	return frames
}

func loadNPCGrid(renderer *sdl.Renderer, path string, cols, rows int) []npcFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, path, cols, rows, npcSpriteInset)
	return framesFromGrid(grid, cols, rows)
}

func loadNPCGridRow(renderer *sdl.Renderer, path string, cols, rows, row int) []npcFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, path, cols, rows, npcSpriteInset)
	var frames []npcFrame
	if row < len(grid) {
		for c := 0; c < cols && c < len(grid[row]); c++ {
			gf := grid[row][c]
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}
	for len(frames) > 1 && frames[len(frames)-1].tex == nil {
		frames = frames[:len(frames)-1]
	}
	return frames
}

// loadNPCGridPath picks the right sprite sheet: the preferred city-specific
// one if its PNG exists, otherwise the given fallback path. Both sheets
// must have the same (cols, rows) geometry so the animation frame counts
// line up.
func loadNPCGridPath(renderer *sdl.Renderer, preferred, fallback string, cols, rows int) []npcFrame {
	if _, err := os.Stat(preferred); err == nil {
		return loadNPCGrid(renderer, preferred, cols, rows)
	}
	return loadNPCGrid(renderer, fallback, cols, rows)
}

// loadNPCGridRowPath is the row-indexed twin of loadNPCGridPath.
func loadNPCGridRowPath(renderer *sdl.Renderer, preferred, fallback string, cols, rows, row int) []npcFrame {
	if _, err := os.Stat(preferred); err == nil {
		return loadNPCGridRow(renderer, preferred, cols, rows, row)
	}
	return loadNPCGridRow(renderer, fallback, cols, rows, row)
}

// --- Director Higgins ---

var higginsDefaultDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Ah, you must be the new counselor! Finally!"},
	{speaker: "Pink Panther", text: "Good afternoon. Pink Panther, at your service."},
	{speaker: "Director Higgins", text: "Yes, yes. Welcome to Camp Chilly Wa Wa."},
	{speaker: "Director Higgins", text: "The kids are through the gate. Go introduce yourself."},
	{speaker: "Director Higgins", text: "They're a good bunch. A little... eccentric, but good."},
}

var higginsPostDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Go on, the kids are waiting in the camp grounds!"},
	{speaker: "Pink Panther", text: "On my way."},
}

var higginsWorriedDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Something is wrong with the kids."},
	{speaker: "Director Higgins", text: "Marcus has been up all night drawing things he's never seen."},
	{speaker: "Director Higgins", text: "Buildings, paintings, rooftops... from places he's never been!"},
	{speaker: "Pink Panther", text: "I saw him last night by the campfire. He was... not himself."},
	{speaker: "Director Higgins", text: "I've seen this kind of thing before... well, no I haven't. But it's NOT normal!"},
	{speaker: "Director Higgins", text: "A glass pyramid, a woman's face... it sounds like Paris. The Louvre."},
	{speaker: "Director Higgins", text: "Here, take this travel map. Camp Chilly Wa Wa Air can get you there."},
	{speaker: "Pink Panther", text: "A camp... airline?"},
	{speaker: "Director Higgins", text: "Don't ask questions. Just go find out what Marcus is connected to."},
}

// higginsLilyHintDialog runs when the camp-grounds Higgins appears next
// to Lily after her shy dialog. Gives the player the flower clue without
// them needing to guess.
var higginsLilyHintDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Ah, counselor. Lily's a quiet one, isn't she."},
	{speaker: "Pink Panther", text: "She barely said two words."},
	{speaker: "Director Higgins", text: "She loves flowers. Try the lake — daisies grow wild by the water."},
	{speaker: "Director Higgins", text: "Bring her one and you'll see a different girl."},
}

var higginsPostWorriedDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "I already gave you the map, Panther."},
	{speaker: "Director Higgins", text: "Come on — we need to fix this up. The kids are counting on us."},
	{speaker: "Director Higgins", text: "Marcus is in the camp grounds. Start there."},
}

func newDirectorHiggins(renderer *sdl.Renderer) *npc {
	// Bounds sized to 200x265 so the aspect-preserve draw produces
	// ~225-235 px of actual sprite on camp_entrance — matches the
	// "adult NPC" row in CHARACTERS.md (PP is 170x235 for reference).
	// Do not shrink below 200x260 or Higgins reads as a kid.
	//
	// Both sheets are clean single-row grids per PROMPTS.md:
	//   idle: 7x1 at 172x384 per cell
	//   talk: 6x1 (clipboard lowered, mouth open)
	return &npc{
		idleGrid:       loadNPCGrid(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 7, 1),
		talkGrid:       loadNPCGridRow(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png", 8, 2, 0),
		bounds:         sdl.Rect{X: 660, Y: 345, W: 200, H: 265},
		name:           "Director Higgins",
		dialog:         higginsDefaultDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.25,
	}
}

func newOfficeHiggins(renderer *sdl.Renderer) *npc {
	// Office Higgins bounds were 180x280 which rendered him at ~35% of
	// screen height — too tall vs the PTP reference. Dropped to 160x225
	// to put him in the 210-225 band from CHARACTERS.md; camp_office's
	// characterScale 0.9 shaves the final render to ~200 which sits
	// comfortably in the tight indoor shot.
	return &npc{
		idleGrid:       loadNPCGridRow(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png", 6, 2, 0),
		talkGrid:       loadNPCGrid(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png", 6, 2),
		// User spec 2026-04-17: office Higgins top-left at (1062, 357),
		// sitting behind the desk. Sized so head lands at ~y=357 and feet
		// rest on the desk chair around y=640.
		bounds:         sdl.Rect{X: 1062, Y: 357, W: 220, H: 280},
		name:           "Director Higgins",
		dialog:         higginsWorriedDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.25,
		silent:         true,
	}
}

// newGroundsHiggins is the hidden Higgins that appears next to the cabin path
// after Lily's shy dialog ends (see setupCampCallbacks). He delivers the
// "she needs a flower" hint. Starts hidden + silent; the callback flips both
// flags when Lily's first dialog completes.
func newGroundsHiggins(renderer *sdl.Renderer) *npc {
	// Positioned by the cabin path near the office entrance, not stacked on
	// top of Marcus (whose bounds start at x=890). 1060x and 570y puts him
	// visible below/right of the kid row without any overlap.
	h := newDirectorHiggins(renderer)
	h.bounds = sdl.Rect{X: 1060, Y: 570, W: 180, H: 200}
	h.hidden = true
	h.silent = true
	h.dialog = higginsLilyHintDialog
	return h
}

// newRoomTommy / newRoomJake / newRoomLily / newRoomDanny return the kid's
// cabin-scene instance: positioned at the room's "bed" spot and silent by
// default. Callbacks flip .silent off when Day 2 story beats start — that's
// how the kid "shows up" in their room after Higgins points PP at them.
//
// Marcus's room NPC is slightly different: he is not silent (Day 1 flow lets
// PP peek in on him immediately) and is drawn larger to fill the room. Kept
// in its own factory to make that intent explicit.
func newRoomTommy(renderer *sdl.Renderer) *npc {
	n := newTommy(renderer)
	n.bounds = sdl.Rect{X: 760, Y: 430, W: 170, H: 260}
	n.silent = true
	return n
}

func newRoomJake(renderer *sdl.Renderer) *npc {
	n := newJake(renderer)
	n.bounds = sdl.Rect{X: 760, Y: 420, W: 170, H: 260}
	n.silent = true
	return n
}

func newRoomLily(renderer *sdl.Renderer) *npc {
	n := newLily(renderer)
	n.bounds = sdl.Rect{X: 666, Y: 461, W: 170, H: 260}
	n.silent = true
	return n
}

func newRoomMarcus(renderer *sdl.Renderer) *npc {
	n := newMarcus(renderer)
	n.bounds = sdl.Rect{X: 526, Y: 181, W: 280, H: 380}
	return n
}

func newRoomDanny(renderer *sdl.Renderer) *npc {
	n := newDanny(renderer)
	n.bounds = sdl.Rect{X: 760, Y: 430, W: 170, H: 260}
	n.silent = true
	return n
}

// newNightHiggins is the campfire Higgins — silent by default so he doesn't
// block exploration, but driven directly by the night cutscene so he appears
// to deliver the "lights out" speech in-place, not at camp grounds.
func newNightHiggins(renderer *sdl.Renderer) *npc {
	return &npc{
		idleGrid:       loadNPCGrid(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 7, 1),
		talkGrid:       loadNPCGridRow(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png", 8, 2, 0),
		bounds:         sdl.Rect{X: 1120, Y: 430, W: 200, H: 260},
		name:           "Director Higgins",
		bobAmount:      0,
		talkFrameSpeed: 0.25,
		silent:         true,
	}
}

// --- Tommy (Homesick Kid) ---

var tommyDialog = []dialogEntry{
	{speaker: "Tommy", text: "Hi! I'm Tommy! Are you the new counselor?"},
	{speaker: "Pink Panther", text: "That's right. Nice to meet you, Tommy."},
	{speaker: "Tommy", text: "I love telling stories! Did you know there's a legend about a treasure at this camp?"},
	{speaker: "Tommy", text: "Probably not true though... I like making things sound more exciting than they are!"},
	{speaker: "Pink Panther", text: "A natural storyteller. I like that."},
}

var tommyPostDialog = []dialogEntry{
	{speaker: "Tommy", text: "Want to hear another story? I've got HUNDREDS!"},
	{speaker: "Pink Panther", text: "Maybe later, Tommy."},
}

var tommyStrangeDialog = []dialogEntry{
	{speaker: "Tommy", text: "Do you hear that? The music?"},
	{speaker: "Pink Panther", text: "Music? I don't hear anything."},
	{speaker: "Tommy", text: "It's drums and singing... and there's a GIANT STATUE watching over everyone!"},
	{speaker: "Tommy", text: "People are dancing in the streets! It's like the biggest party in the world!"},
	{speaker: "Tommy", text: "And then... tango? Somewhere else, a different city, a wide road with a tall white tower..."},
	{speaker: "Pink Panther", text: "Tommy, are you alright? You've never been to any of these places."},
	{speaker: "Tommy", text: "I KNOW! That's what's so weird! But I can SEE it!"},
}

var tommyPostStrangeDialog = []dialogEntry{
	{speaker: "Tommy", text: "The music won't stop... a giant statue, parades, dancing..."},
	{speaker: "Tommy", text: "It feels like two places at once. I can't explain it."},
}

func newTommy(renderer *sdl.Renderer) *npc {
	n := &npc{
		bounds:         sdl.Rect{X: 130, Y: 405, W: 150, H: 180},
		name:           "Tommy",
		dialog:         tommyDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	applyKidAtlas(renderer, n, "tommy")
	return n
}

// --- Jake (Bully Kid) ---

var jakeDialog = []dialogEntry{
	{speaker: "Jake", text: "Hey! You the new guy?"},
	{speaker: "Pink Panther", text: "That's right. And you are?"},
	{speaker: "Jake", text: "Jake. I'm the toughest kid at camp. Don't forget it."},
	{speaker: "Jake", text: "I collect stuff. Rocks, coins, anything shiny. Check out this coin my dad brought from Israel."},
	{speaker: "Pink Panther", text: "That's a beautiful coin. Where exactly is it from?"},
	{speaker: "Jake", text: "Some old city with tunnels underneath. Jerusalem, I think. Dad said the tunnels are ANCIENT."},
	{speaker: "Pink Panther", text: "Fascinating collection you've got there."},
}

var jakePostDialog = []dialogEntry{
	{speaker: "Jake", text: "Don't touch my collection. I'm watching you."},
	{speaker: "Pink Panther", text: "Wouldn't dream of it."},
}

var jakeStrangeDialog = []dialogEntry{
	{speaker: "Jake", text: "Something's happening to my coins..."},
	{speaker: "Pink Panther", text: "What do you mean?"},
	{speaker: "Jake", text: "I keep hearing echoes. Like tunnels underground. Voices bouncing off stone walls."},
	{speaker: "Jake", text: "And I can't stop rubbing every metal surface for symbols. Look at this bench - I KNOW there's something underneath!"},
	{speaker: "Pink Panther", text: "Jake, that's just a wooden bench."},
	{speaker: "Jake", text: "NO! There are tunnels! Old ones! Under an ancient city! I can FEEL them!"},
}

var jakePostStrangeDialog = []dialogEntry{
	{speaker: "Jake", text: "The echoes won't stop... tunnels under old stone walls..."},
	{speaker: "Jake", text: "It's like I can see a huge wall... and something hidden behind it."},
}

func newJake(renderer *sdl.Renderer) *npc {
	n := &npc{
		bounds:         sdl.Rect{X: 370, Y: 400, W: 150, H: 180},
		name:           "Jake",
		dialog:         jakeDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	applyKidAtlas(renderer, n, "jake")
	return n
}

// --- Lily (Shy Girl) ---

var lilyShyDialog = []dialogEntry{
	{speaker: "Lily", text: "..."},
	{speaker: "Pink Panther", text: "Hello there. I'm the new counselor."},
	{speaker: "Lily", text: "..."},
	{speaker: "Pink Panther", text: "Not much of a talker, huh?"},
}

var lilyFlowerDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "I found this flower by the lake. Would you like it?"},
	{speaker: "Lily", text: "...! A daisy! It's beautiful!"},
	{speaker: "Lily", text: "...thank you... nobody ever brings me flowers..."},
	{speaker: "Pink Panther", text: "I'm the new counselor. What's your name?"},
	{speaker: "Lily", text: "...Lily... I like flowers... and quiet places..."},
	{speaker: "Pink Panther", text: "Nice to meet you, Lily. Those are beautiful flowers you're arranging."},
	{speaker: "Lily", text: "...thank you... you're nice..."},
}

var lilyDialog = []dialogEntry{
	{speaker: "Lily", text: "...hi again..."},
	{speaker: "Pink Panther", text: "Hello, Lily. Beautiful day, isn't it?"},
	{speaker: "Lily", text: "*small nod*"},
}

var lilyPostDialog = lilyDialog

var lilyStrangeDialog = []dialogEntry{
	{speaker: "Lily", text: "...the flowers are glowing..."},
	{speaker: "Pink Panther", text: "Glowing? They look normal to me."},
	{speaker: "Lily", text: "Not these flowers... OTHER flowers. In a garden far away..."},
	{speaker: "Lily", text: "I keep arranging petals into shapes... symbols I don't understand..."},
	{speaker: "Lily", text: "And I hear bells. Temple bells. Very old ones."},
	{speaker: "Lily", text: "There's a red gate... and cherry blossoms falling everywhere..."},
	{speaker: "Pink Panther", text: "That sounds like Japan, Lily. Have you ever been there?"},
	{speaker: "Lily", text: "...never... but I can see it when I close my eyes..."},
}

var lilyPostStrangeDialog = []dialogEntry{
	{speaker: "Lily", text: "...the bells again... and glowing petals..."},
	{speaker: "Lily", text: "...a temple in the mountains... I can almost touch it..."},
}

func newLily(renderer *sdl.Renderer) *npc {
	n := &npc{
		bounds:         sdl.Rect{X: 600, Y: 395, W: 150, H: 180},
		name:           "Lily",
		dialog:         lilyShyDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	applyKidAtlas(renderer, n, "lily")
	return n
}

// --- Marcus (Know-It-All) ---

var marcusDialog = []dialogEntry{
	{speaker: "Marcus", text: "Ah, a new counselor! Did you know this camp was founded in 1968?"},
	{speaker: "Pink Panther", text: "I did not. And you are?"},
	{speaker: "Marcus", text: "Marcus. I know everything about everything. Ask me anything!"},
	{speaker: "Pink Panther", text: "Alright. What's the most interesting thing about this camp?"},
	{speaker: "Marcus", text: "Statistically, the mess hall food has a 73 percent chance of being inedible."},
	{speaker: "Marcus", text: "But I also love drawing! Want to see my sketches? I drew the whole campfire!"},
	{speaker: "Pink Panther", text: "Very impressive. You've got talent, Marcus."},
}

var marcusPostDialog = []dialogEntry{
	{speaker: "Marcus", text: "Did you know butterflies taste with their feet? It's TRUE!"},
	{speaker: "Pink Panther", text: "You never stop, do you?"},
}

var marcusStrangeDialog = []dialogEntry{
	{speaker: "Marcus", text: "It's WRONG! The picture is WRONG!"},
	{speaker: "Pink Panther", text: "Marcus? What's going on?"},
	{speaker: "Marcus", text: "I keep drawing this woman's face... but I've NEVER seen her before!"},
	{speaker: "Marcus", text: "And these frames... ornate golden frames... and rooftops I've never visited!"},
	{speaker: "Marcus", text: "It's a museum. A HUGE museum. The biggest in the world!"},
	{speaker: "Marcus", text: "There's a glass pyramid in front of it... and inside, a painting that everyone stares at..."},
	{speaker: "Marcus", text: "But something is MISSING from the picture! I can feel it!"},
	{speaker: "Pink Panther", text: "A glass pyramid... the biggest museum... That sounds like the Louvre in Paris."},
	{speaker: "Marcus", text: "I've never been to Paris! But I can't stop drawing it!"},
}

var marcusPostStrangeDialog = []dialogEntry{
	{speaker: "Marcus", text: "The woman's face again... the golden frames... something is missing..."},
	{speaker: "Marcus", text: "I filled twelve pages last night. I can't stop."},
}

func newMarcus(renderer *sdl.Renderer) *npc {
	n := &npc{
		bounds:         sdl.Rect{X: 890, Y: 395, W: 150, H: 180},
		name:           "Marcus",
		dialog:         marcusDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		// Freakout feels frantic if it runs at normal talk speed — slow it
		// down so the strange dialogue has room to breathe.
		strangeTalkFrameSpeed: 0.16,
	}
	applyKidAtlas(renderer, n, "marcus")
	return n
}

// --- Danny (Prankster) ---

var dannyDialog = []dialogEntry{
	{speaker: "Danny", text: "Psst! Hey! Over here!"},
	{speaker: "Pink Panther", text: "Hmm? And who might you be?"},
	{speaker: "Danny", text: "I'm Danny, master of mischief! I'm setting up the ULTIMATE prank!"},
	{speaker: "Danny", text: "I love treasure stories. My cousin went to Italy once and saw REAL ancient ruins!"},
	{speaker: "Danny", text: "The Colosseum! Gladiators fought there! How cool is that?!"},
	{speaker: "Pink Panther", text: "Very cool, Danny. Try not to prank anyone too badly."},
}

var dannyPostDialog = []dialogEntry{
	{speaker: "Danny", text: "Psst! Want to help me put a frog in Higgins' coffee?"},
	{speaker: "Pink Panther", text: "I'll pretend I didn't hear that."},
}

var dannyStrangeDialog = []dialogEntry{
	{speaker: "Danny", text: "Dude! DUDE! There's treasure EVERYWHERE!"},
	{speaker: "Pink Panther", text: "Danny, calm down. What are you talking about?"},
	{speaker: "Danny", text: "I've been mapping the whole camp! It's just like ancient ruins!"},
	{speaker: "Danny", text: "There are gold paths under the ground... I can FEEL them!"},
	{speaker: "Danny", text: "A huge round arena... with arches... thousands of people cheering..."},
	{speaker: "Danny", text: "And tunnels underneath with hidden rooms full of treasure!"},
	{speaker: "Pink Panther", text: "An arena with arches... that sounds like the Colosseum in Rome."},
	{speaker: "Danny", text: "I've never been to Rome! But I drew a map of it! Look!"},
}

var dannyPostStrangeDialog = []dialogEntry{
	{speaker: "Danny", text: "The gold paths are getting clearer... arches and tunnels everywhere..."},
	{speaker: "Danny", text: "I've dug three holes behind the cabin already. Higgins is NOT happy."},
}

func newDanny(renderer *sdl.Renderer) *npc {
	n := &npc{
		bounds:         sdl.Rect{X: 1110, Y: 400, W: 150, H: 180},
		name:           "Danny",
		dialog:         dannyDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		flipped:        true,
	}
	applyKidAtlas(renderer, n, "danny")
	return n
}

const (
	npcAnimIdle  = 0
	npcAnimTalk  = 1
	npcAnimDrink = 2
)

func (n *npc) setAnimState(state int) {
	if n.animState == state {
		return
	}
	n.animState = state
	n.curFrame = 0
	n.frameTimer = 0
}

func (n *npc) update(dt float64) {
	n.bobTimer += dt

	speed := n.talkFrameSpeed
	if speed <= 0 {
		speed = 0.12
	}
	// Strange state gets its own talk speed so the freakout doesn't strobe
	// (Marcus). NPCs that don't override stay on the default speed.
	if n.isStrange && n.strangeTalkFrameSpeed > 0 {
		speed = n.strangeTalkFrameSpeed
	}

	if len(n.idleGrid) > 1 {
		n.idleFrameTimer += dt
		idleSpeed := speed * 2.5 // idle cycles slower than talk
		if n.idleFrameTimer >= idleSpeed {
			n.idleFrameTimer -= idleSpeed
			n.idleCurFrame = (n.idleCurFrame + 1) % len(n.idleGrid)
		}
	}

	if n.animState == npcAnimTalk && len(n.talkGrid) > 0 {
		n.frameTimer += dt
		if n.frameTimer >= speed {
			n.frameTimer -= speed
			n.curFrame = (n.curFrame + 1) % len(n.talkGrid)
		}
	}
}

func (n *npc) draw(renderer *sdl.Renderer) {
	n.drawScaled(renderer, 1.0)
}

// drawScaled renders the NPC with an additional character-scale factor
// applied to the on-screen size. The hitbox (n.bounds) stays at its
// authored dimensions so click targets don't shrink with the scene
// scale. The visible sprite is anchored at foot-center so shrinking
// only trims from the head and shoulders.
func (n *npc) drawScaled(renderer *sdl.Renderer, charScale float64) {
	if n.hidden {
		return
	}
	if charScale <= 0 {
		charScale = 1.0
	}
	bobOffset := int32(math.Sin(n.bobTimer*1.5) * n.bobAmount)
	breathScale := 1.0

	shadowCX := n.bounds.X + n.bounds.W/2
	shadowFY := n.bounds.Y + n.bounds.H
	drawShadow(renderer, shadowCX, shadowFY, int32(float64(n.bounds.W-10)*charScale))

	flip := sdl.FLIP_NONE
	if n.flipped {
		flip = sdl.FLIP_HORIZONTAL
	}

	var frame npcFrame
	if n.animState == npcAnimTalk && len(n.talkGrid) > 0 {
		frame = n.talkGrid[n.curFrame%len(n.talkGrid)]
	} else if len(n.idleGrid) > 0 {
		frame = n.idleGrid[n.idleCurFrame%len(n.idleGrid)]
	}

	if frame.tex == nil {
		return
	}

	targetW := float64(n.bounds.W) * charScale
	targetH := float64(n.bounds.H) * charScale
	scaleW := targetW * breathScale / float64(frame.w)
	scaleH := targetH * breathScale / float64(frame.h)
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}
	dstW := int32(float64(frame.w) * scale)
	dstH := int32(float64(frame.h) * scale)
	dstX := n.bounds.X + (n.bounds.W-dstW)/2
	dstY := n.bounds.Y + bobOffset + (n.bounds.H - dstH)

	dst := sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH}
	renderer.CopyEx(frame.tex, frame.src, &dst, 0, nil, flip)
}

// containsPoint is used for both cursor hover (showing the "talk" icon) and
// actual click detection. Keeping them unified means: wherever the cursor
// shows "talk", a click always lands. We pad generously so small sprites or
// slightly-missed clicks still register as an interaction.
func (n *npc) containsPoint(x, y int32) bool {
	// Strict-bounds hit test (user request 2026-04-17): the cursor must
	// land inside the NPC's authored rect. No radius padding.
	// The old padX=70, padY=50 expansion made clicks snap to the wrong
	// NPC when two kids stood close (Danny vs Marcus at camp_grounds)
	// and let clicks on empty ground behind an NPC trigger dialog.
	pt := sdl.Point{X: x, Y: y}
	return pt.InRect(&n.bounds)
}

func (n *npc) footY() int32 {
	return n.bounds.Y + n.bounds.H
}

// ===== Paris NPCs =====

var frenchGuideDialog = []dialogEntry{
	{speaker: "Madame Colette", text: "Bonjour, monsieur! Welcome to Paris!"},
	{speaker: "Pink Panther", text: "Bonjour, madame. I'm looking for information about the Louvre."},
	{speaker: "Madame Colette", text: "Ah, ze Louvre! Ze largest art museum in ze world!"},
	{speaker: "Madame Colette", text: "It was originally a royal palace, built in ze 12th century."},
	{speaker: "Madame Colette", text: "Today it holds over 380,000 objects and 35,000 works of art!"},
	{speaker: "Pink Panther", text: "Impressive. And what about that glass pyramid?"},
	{speaker: "Madame Colette", text: "Ah, ze Pyramid! Designed by I.M. Pei in 1989. Very controversial at first!"},
	{speaker: "Madame Colette", text: "People said it did not belong. Now it is ze most famous entrance in ze world."},
	{speaker: "Madame Colette", text: "And of course, ze Eiffel Tower behind you — built in 1889 for ze World Fair."},
	{speaker: "Madame Colette", text: "Gustave Eiffel designed it. It was meant to be temporary — just 20 years!"},
	{speaker: "Madame Colette", text: "But zey kept it because it was perfect for radio transmissions."},
	{speaker: "Pink Panther", text: "A temporary tower that became permanent. How fitting."},
	{speaker: "Madame Colette", text: "Ze museum is just down ze street, to ze right. Enjoy, monsieur!"},
}

var frenchGuidePostDialog = []dialogEntry{
	{speaker: "Madame Colette", text: "Ze Louvre is to ze right, monsieur. You cannot miss ze pyramid!"},
	{speaker: "Pink Panther", text: "Merci, madame."},
}

// --- Bakery Woman (pre-Louvre quest, step 1) ---
// Sells PP a baguette, which he trades to Pierre for a press pass, which
// he shows Claude to get the museum ticket that unlocks the Louvre. Retro-
// style "collect props before the main door opens" chain.
var bakeryWomanDialog = []dialogEntry{
	{speaker: "Madame Poulain", text: "Bonjour, monsieur! Fresh baguettes, straight from ze oven!"},
	{speaker: "Pink Panther", text: "They smell wonderful. I'd love one."},
	{speaker: "Madame Poulain", text: "For you, a compliment, and ze bread is yours. Non?"},
	{speaker: "Pink Panther", text: "Madame, your boulangerie smells like Paris itself."},
	{speaker: "Madame Poulain", text: "*laughs* Charmant! Here, take a baguette. Tell your friends!"},
}

var bakeryWomanPostDialog = []dialogEntry{
	{speaker: "Madame Poulain", text: "Enjoy ze baguette, monsieur! Zhere's a photographer near ze museum — he loves fresh bread."},
}

func newBakeryWoman(renderer *sdl.Renderer) *npc {
	// Dedicated Bakery Woman sheet (see docs/EXTRA_PROMPTS.md §8). 8×2
	// canvas: row 0 = idle (mouth closed), row 1 = talk (mouth open).
	// Packed atlas at assets/sprites/paris/bakery_woman.(png|json) is the
	// preferred path; legacy per-row PNG slicing stays as a fallback so
	// the NPC still spawns if pack_atlas.py hasn't been run.
	n := &npc{
		bounds:         sdl.Rect{X: 540, Y: 440, W: 140, H: 240},
		name:           "Madame Poulain",
		dialog:         bakeryWomanDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
		flipped:        false, // sheet draws her facing right already
	}
	if !applyNPCAtlas(renderer, n, "paris/bakery_woman") {
		const sheet = "assets/images/locations/paris/npc/npc_bakery_woman.png"
		n.idleGrid = loadNPCGridRow(renderer, sheet, 8, 2, 0)
		n.talkGrid = loadNPCGridRow(renderer, sheet, 8, 2, 1)
	}
	return n
}

// --- Press Photographer (flavor NPC near the Louvre steps) ---
// Madame Poulain's post-baguette dialog name-drops a photographer near the
// museum. Nicolas is that flavor NPC — chatty Parisian with a camera slung
// over his shoulder. He is not on the critical quest chain; Pierre still
// hands over the press pass in exchange for the baguette.
var pressPhotographerDialog = []dialogEntry{
	{speaker: "Nicolas", text: "Ah, a visitor! Hold still — ze light is perfect."},
	{speaker: "Pink Panther", text: "Are you... photographing me?"},
	{speaker: "Nicolas", text: "Non, non, I photograph Paris. You happen to be in ze frame."},
	{speaker: "Nicolas", text: "I have been here twenty years. I have seen ze Louvre in every weather."},
	{speaker: "Pink Panther", text: "Any advice for a curious traveler?"},
	{speaker: "Nicolas", text: "Talk to Pierre ze painter and Claude ze gendarme. Zey know ze street better zhan ze guidebooks."},
}

var pressPhotographerPostDialog = []dialogEntry{
	{speaker: "Nicolas", text: "Bonne chance, monsieur! Smile for ze camera."},
}

func newPressPhotographer(renderer *sdl.Renderer) *npc {
	// Dedicated Press Photographer sheet (see docs/EXTRA_PROMPTS.md §9). 8×2
	// canvas: row 0 = idle (mouth closed), row 1 = talk (mouth open).
	// Positioned between Pierre (X=880) and Claude (X=1120) — fits the
	// Bakery Woman's "photographer near ze museum" breadcrumb. Tight cluster
	// of Paris street characters by the Louvre entrance hotspot (x=1300).
	// Packed atlas at assets/sprites/paris/press_photographer.(png|json)
	// is preferred; legacy PNG slicing stays as a fallback.
	n := &npc{
		bounds:         sdl.Rect{X: 1010, Y: 440, W: 110, H: 240},
		name:           "Nicolas",
		dialog:         pressPhotographerDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
		flipped:        false, // sheet draws him facing right already
	}
	if !applyNPCAtlas(renderer, n, "paris/press_photographer") {
		const sheet = "assets/images/locations/paris/npc/npc_press_photographer.png"
		n.idleGrid = loadNPCGridRow(renderer, sheet, 8, 2, 0)
		n.talkGrid = loadNPCGridRow(renderer, sheet, 8, 2, 1)
	}
	return n
}

func newFrenchGuide(renderer *sdl.Renderer) *npc {
	// Packed atlas at assets/sprites/paris/french_guide.(png|json) is the
	// preferred path; legacy per-sheet PNG loading stays as a fallback.
	// Feet land at y≈680 on the paris_street floor line; user reported
	// the previous Y=350 (feet ~590) had NPCs floating above the ground.
	n := &npc{
		bounds:         sdl.Rect{X: 300, Y: 440, W: 140, H: 240},
		name:           "Madame Colette",
		dialog:         frenchGuideDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
	if !applyNPCAtlas(renderer, n, "paris/french_guide") {
		n.idleGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/npc_french_guide_idle.png", 8, 2)
		n.talkGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/npc_french_guide_talk.png", 8, 1)
	}
	return n
}

var museumCuratorDialog = []dialogEntry{
	{speaker: "Curator Beaumont", text: "Ah, a visitor! Welcome to ze Musee du Louvre."},
	{speaker: "Pink Panther", text: "Thank you. I'm investigating something... unusual."},
	{speaker: "Curator Beaumont", text: "Unusual? In ze Louvre, everything is extraordinary!"},
	{speaker: "Curator Beaumont", text: "Zis hall contains some of ze finest works in history."},
	{speaker: "Curator Beaumont", text: "Ze Mona Lisa, of course — painted by Leonardo da Vinci around 1503."},
	{speaker: "Curator Beaumont", text: "Her smile has puzzled visitors for over 500 years!"},
	{speaker: "Curator Beaumont", text: "And ze Venus de Milo — a Greek sculpture from around 100 BC."},
	{speaker: "Pink Panther", text: "Actually, I'm looking for a specific painting. A boy back at camp keeps drawing it."},
	{speaker: "Curator Beaumont", text: "A boy... drawing paintings he has never seen? How peculiar."},
	{speaker: "Curator Beaumont", text: "Describe what he draws, and perhaps I can identify it."},
	{speaker: "Pink Panther", text: "A woman's face. Ornate golden frames. He says something is 'missing' from it."},
	{speaker: "Curator Beaumont", text: "Mon Dieu... zat sounds like ze portrait in Room 7."},
	{speaker: "Curator Beaumont", text: "A painting zat was recently restored. Ze restorer found a hidden symbol underneath."},
	{speaker: "Curator Beaumont", text: "Perhaps your boy senses what was hidden. Take zis postcard of ze painting."},
	{speaker: "Curator Beaumont", text: "If he sees ze complete image, perhaps his mind will settle."},
}

var museumCuratorPostDialog = []dialogEntry{
	{speaker: "Curator Beaumont", text: "Ze postcard should help your young friend."},
	{speaker: "Curator Beaumont", text: "Ze mysteries of art connect us in ways we do not understand."},
}

func newMuseumCurator(renderer *sdl.Renderer) *npc {
	// Packed atlas at assets/sprites/paris/museum_curator.(png|json) is the
	// preferred path; legacy per-sheet PNG loading stays as a fallback.
	n := &npc{
		bounds:         sdl.Rect{X: 500, Y: 320, W: 130, H: 250},
		name:           "Curator Beaumont",
		dialog:         museumCuratorDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
	if !applyNPCAtlas(renderer, n, "paris/museum_curator") {
		n.idleGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/npc_museum_curator_idle.png", 8, 1)
		n.talkGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/npc_museum_curator_talk.png", 4, 2)
	}
	return n
}

// --- Pierre the Street Artist ---
// A friendly beret-wearing painter who sells portraits on the sidewalk.
// Typical retro-adventure "local" NPC — adds flavour and drops a casual
// clue, but isn't a guide. Uses npc_art_vendor.png (8x2 grid).
var pierreArtistDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Bonjour! You're painting... Pink cats?"},
	{speaker: "Pierre", text: "Oui! Pink, blue, ze panther-colors. Monet himself loved ze violet shadows."},
	{speaker: "Pierre", text: "I am Pierre. Zis sidewalk, zis easel — zat is my whole world since 1982."},
	{speaker: "Pink Panther", text: "Quite a view. The tower, the cafe, the pigeons."},
	{speaker: "Pierre", text: "Ze pigeons are ze real critics. If zey do not land on ze canvas, ze painting is no good."},
	{speaker: "Pink Panther", text: "I'm looking for a boy who keeps drawing a woman's face. Something missing from it."},
	{speaker: "Pierre", text: "Hm. Ze Curator inside ze Louvre, she knows every face in Paris. Ask her."},
	{speaker: "Pierre", text: "Tell her Pierre sent you. She still owes me a coffee from ze '89 restoration."},
}

var pierreArtistPostDialog = []dialogEntry{
	{speaker: "Pierre", text: "Don't forget — ze pigeons approve of your pink, monsieur!"},
}

func newPierreArtist(renderer *sdl.Renderer) *npc {
	// Packed atlas at assets/sprites/paris/pierre_artist.(png|json) is the
	// preferred path; legacy per-row PNG slicing stays as a fallback.
	n := &npc{
		bounds:         sdl.Rect{X: 880, Y: 440, W: 130, H: 240},
		name:           "Pierre",
		dialog:         pierreArtistDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
	if !applyNPCAtlas(renderer, n, "paris/pierre_artist") {
		const sheet = "assets/images/locations/paris/npc/npc_art_vendor.png"
		n.idleGrid = loadNPCGridRow(renderer, sheet, 8, 2, 0)
		n.talkGrid = loadNPCGridRow(renderer, sheet, 8, 2, 1)
	}
	return n
}

// --- Gendarme Claude ---
// Friendly Parisian police officer stationed near the Louvre entrance.
// Adds a second local on the street and can warn about pickpockets so the
// player gets a reason to clutch the postcard on the way back.
var gendarmeDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Officer. Beautiful evening."},
	{speaker: "Claude", text: "Bonsoir, monsieur. Gendarme Claude, at your service."},
	{speaker: "Claude", text: "Watch out for ze pickpockets near ze tower. Zey move like cats."},
	{speaker: "Claude", text: "And ze mimes! Ze mimes are ze worst — zey steal your attention, zen your wallet."},
	{speaker: "Pink Panther", text: "I'll keep both eyes on my pocket. Is the Louvre still open?"},
	{speaker: "Claude", text: "Oui, ze curator stays late on Fridays. Tell her Claude said bonjour."},
	{speaker: "Claude", text: "Bon courage, monsieur panther."},
}

var gendarmePostDialog = []dialogEntry{
	{speaker: "Claude", text: "Pickpockets — eyes open, monsieur!"},
}

func newGendarmeClaude(renderer *sdl.Renderer) *npc {
	// Packed atlas at assets/sprites/paris/gendarme_claude.(png|json) is
	// the preferred path; legacy per-row PNG slicing stays as a fallback.
	n := &npc{
		bounds:         sdl.Rect{X: 1120, Y: 430, W: 120, H: 250},
		name:           "Claude",
		dialog:         gendarmeDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.12,
	}
	if !applyNPCAtlas(renderer, n, "paris/gendarme_claude") {
		const sheet = "assets/images/locations/paris/npc/npc_security_guard.png"
		n.idleGrid = loadNPCGridRow(renderer, sheet, 6, 2, 0)
		n.talkGrid = loadNPCGridRow(renderer, sheet, 6, 2, 1)
	}
	return n
}

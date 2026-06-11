package game

import (
	"fmt"
	"math"
	"os"
	"sort"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type npcFrame struct {
	tex *sdl.Texture
	w   int32
	h   int32
	// ox/oy/ow/oh is the opaque content box within the frame (frame-local).
	// drawScaled scales by this so idle/talk/walk render at one size and the
	// sprite anchors by feet + horizontal centre. Zero ow/oh = use full w/h.
	ox int32
	oy int32
	ow int32
	oh int32
	// fcx/fry: per-frame FEET anchors from the engine (foot-band centre X +
	// feet-line Y, both cell-local; thin tail strands excluded). drawScaled
	// anchors by these — deadband-snapped to the sheet median by
	// stabilizeNPCAnchors — so each frame's feet land on the same screen
	// spot even when the art drifts inside the cells. 0 = unset → fall back
	// to the opaque-box centre/bottom.
	fcx int32
	fry int32
	// src is the source rectangle inside tex. nil means "draw the whole
	// texture" (legacy per-frame loaders produce one texture per frame so
	// they leave this nil). Atlas-backed frames share one texture and set
	// src to the frame's rect within the atlas.
	src *sdl.Rect
	// srcPath is the on-disk PNG this frame came from. Used by the click
	// probe diagnostic (F2) to re-decode the file and sample the alpha
	// channel at the clicked pixel — that's how we tell a "real" hit on
	// the cartoon outline apart from a click in a transparent halo of a
	// sloppy BG-cut. Empty when the loader didn't track the path.
	srcPath string
}

type npc struct {
	bounds sdl.Rect
	// drawFootY, when > 0, overrides the foot-Y used for back-to-front draw
	// sorting (scene.drawActors). Lets a seated NPC drawn high on screen (small
	// bounds.Y) still sort to the FRONT — e.g. café patrons at the front tables
	// should render in front of PP even though their bounds are near the top (#27).
	drawFootY int32
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
	// approachRight forces PP to walk to this NPC's RIGHT side. Used for the
	// far-right kid (Danny) whose left side overlaps Marcus, so the default
	// "approach right-half NPCs from the left" rule would stand PP on Marcus (#7).
	approachRight bool
	// fixedFacing keeps the NPC's authored `flipped` during dialog instead of
	// auto-turning to face PP. For seated NPCs (office Higgins) who must hold a
	// fixed orientation behind a desk (#16).
	fixedFacing bool
	silent      bool
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
	// altDialogStrictMissingHint, when non-empty, replaces the regular
	// kid.dialog playback when the player clicks the NPC but doesn't yet
	// hold/own the altDialogRequiresItem. Without this, an NPC like Lily
	// who has progressed past the hint stage would happily replay her
	// "thanks for the flower" dialog every click even if PP doesn't have
	// the flower in hand — the user 2026-05-21 reported this as "the
	// flower dialog plays without giving her the flower". With this set,
	// clicking on a gated NPC without the item plays a short hint instead
	// (e.g. "She won't look up — maybe she needs something special.")
	// until the trade actually happens.
	altDialogStrictMissingHint []dialogEntry
	// onClickOverride, when set, COMPLETELY replaces the normal click→walk-to-
	// →dialog flow for this NPC. Pierre uses it for the depth-walk choreography
	// (PP walks to middle of road → playRecede shrink → talk). The handler is
	// responsible for everything (movement, dialog, scale restore). Set to nil
	// (default) for standard NPC click flow. User 2026-05-22.
	onClickOverride func()

	// altIdleGrid is an optional alt-idle frame strip that the engine swaps
	// idleGrid for periodically while the player isn't interacting with this
	// NPC. Marcus's `_strange_alt` sheet uses this for ambient "freakout
	// punctuation" — every altIdleAfterSec the alt frames play one cycle,
	// then idleGrid restores. User 2026-05-22.
	altIdleGrid     []npcFrame
	idleAccumSec    float64 // accumulates while in npcAnimIdle; reset by dialog start
	altIdleAfterSec float64 // 0 disables; otherwise seconds before swap fires
	altIdleActive   bool    // engine-private: tracks "currently in alt cycle"
	altIdleBackup   []npcFrame

	// srcCropBottomFrac, if in (0, 1.0), tells drawScaled to render only the
	// TOP portion of each frame (e.g. 0.55 = top 55%). Used for café patrons
	// whose source sheet is full-body but the BG already has chair art under
	// them — clipping the bottom hides duplicate legs. 0 means draw the full
	// frame (default). User 2026-05-22.
	srcCropBottomFrac float64
	// extraScale is an extra render-size multiplier lerped by a SeqNPCMove
	// endScale, so an NPC genuinely shrinks as he walks "into" the scene
	// (Jake stepping back into his cabin). 0 or 1 = no extra scaling.
	extraScale float64
	// hintState is a small per-NPC dialog progression counter. Lily uses
	// 0 = has not been spoken to, 1 = shy dialog played (waiting for
	// flower), 2 = flower given. Storing this on the NPC instead of a
	// closure variable keeps the state deterministic across scene
	// re-entry and save/load (closures would reset back to zero when
	// setupCampCallbacks ran again).
	hintState int
	sm        *npcStateMachine  // optional state machine (named states: default/post/strange/post_strange)
	rules     []InteractionRule // optional rule list for data-driven interactions (see npc_rules.go)
	// game is a back-reference set by spawnNPCs so rule-driven NPCs can
	// call g.fireTrigger without threading *Game through every handler.
	// Not set for NPCs built via legacy callbacks — the rules slice stays
	// empty for those and fireTrigger is a no-op.
	game *Game

	idleGrid []npcFrame
	talkGrid []npcFrame
	// postGiveTalkGrid, if set, replaces talkGrid after this NPC receives the
	// quest item it was waiting for — e.g. Lily holding the daisy while she
	// talks once PP has handed it over (#4). Swapped in by the give callback.
	postGiveTalkGrid []npcFrame
	talkFrameSpeed   float64
	curFrame         int
	frameTimer       float64
	idleCurFrame     int
	idleFrameTimer   float64
	animState        int

	strangeIdle []npcFrame
	strangeTalk []npcFrame
	normalIdle  []npcFrame
	normalTalk  []npcFrame
	isStrange   bool
	// strangeTalkFrameSpeed slows the talk animation while the NPC is in
	// strange state (Marcus's freakout looked too flickery at the default
	// 0.10 s/frame). 0 = inherit talkFrameSpeed unchanged.
	strangeTalkFrameSpeed float64

	// oneShotAnims holds named, non-loop animations the sequence player
	// can trigger (e.g. Higgins's "give_map"). When activeOneShot != "" the
	// draw loop renders from oneShotAnims[activeOneShot] using oneShotIdx
	// instead of idleGrid/talkGrid. Duration is enforced by the sequence
	// player which calls endOneShotAnim() when the timeline ends.
	oneShotAnims    map[string][]npcFrame
	activeOneShot   string
	oneShotIdx      int
	oneShotTimer    float64
	oneShotDuration float64

	// swappedIdleBackup holds the original idleGrid when a sequence step
	// temporarily swaps it for a looping named animation (e.g. Higgins's
	// "walk_back" cycle during an npc_move). The next idle/talk anim step
	// restores it via restoreSwappedIdle(). Unlike one-shots this path
	// uses the existing idle frame cycler so the animation loops at the
	// natural talkFrameSpeed × 2.5 pace.
	swappedIdleBackup []npcFrame

	// lastDrawRect caches the on-screen rect from the most recent
	// drawScaled call (after characterScale + aspect-preserve). containsPoint
	// uses this so hover + click only register on the visible sprite, not
	// the wider bounds rect (which is sized for design-time anchoring).
	// Zero until the first frame; containsPoint falls back to bounds in
	// that case so initial-frame clicks aren't lost.
	lastDrawRect sdl.Rect
	// lastDrawnFrame mirrors lastDrawRect for the source side: the exact
	// npcFrame that was rendered most recently (idle vs talk vs one-shot,
	// frame index baked in). The click probe samples this frame's PNG
	// alpha to validate the BG cut.
	lastDrawnFrame npcFrame
	lastDrawnFlip  bool
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
func framesFromGrid(grid [][]engine.GridFrame, cols, rows int, srcPath string) []npcFrame {
	var frames []npcFrame
	for r := 0; r < rows && r < len(grid); r++ {
		for c := 0; c < cols && c < len(grid[r]); c++ {
			gf := grid[r][c]
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH,
				fcx: gf.FCX, fry: gf.FRY, srcPath: srcPath})
		}
	}
	stabilizeNPCAnchors(frames)
	for len(frames) > 1 && frames[len(frames)-1].tex == nil {
		frames = frames[:len(frames)-1]
	}
	// A missing sheet now yields a shaped-but-empty grid (engine.emptyGrid)
	// instead of a panic — collapse it to nil so `len(frames) > 0` guards
	// don't register invisible animations.
	if len(frames) == 1 && frames[0].tex == nil {
		return nil
	}
	return frames
}

func loadNPCGrid(renderer *sdl.Renderer, path string, cols, rows int) []npcFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, path, cols, rows, npcSpriteInset)
	return framesFromGrid(grid, cols, rows, path)
}

func loadNPCGridRow(renderer *sdl.Renderer, path string, cols, rows, row int) []npcFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, path, cols, rows, npcSpriteInset)
	var frames []npcFrame
	if row < len(grid) {
		for c := 0; c < cols && c < len(grid[row]); c++ {
			gf := grid[row][c]
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH,
				fcx: gf.FCX, fry: gf.FRY, srcPath: path})
		}
	}
	for len(frames) > 1 && frames[len(frames)-1].tex == nil {
		frames = frames[:len(frames)-1]
	}
	stabilizeNPCAnchors(frames)
	return frames
}

// loadNPCGridClean is loadNPCGrid with a wider color-key tolerance (16 vs the
// default 8). Use for sheets whose background leaves a visible halo at the
// default tolerance — Higgins idle/talk are the canonical adult case.
func loadNPCGridClean(renderer *sdl.Renderer, path string, cols, rows int) []npcFrame {
	grid := engine.SpriteGridFromPNGCleanKids(renderer, path, cols, rows, npcSpriteInset)
	return framesFromGrid(grid, cols, rows, path)
}

func loadNPCGridConnected(renderer *sdl.Renderer, path string, cols, rows int) []npcFrame {
	grid := engine.SpriteGridFromPNGCleanConnected(renderer, path, cols, rows, npcSpriteInset)
	return framesFromGrid(grid, cols, rows, path)
}

// loadNPCGridRowClean is the row-indexed twin of loadNPCGridClean.
func loadNPCGridRowClean(renderer *sdl.Renderer, path string, cols, rows, row int) []npcFrame {
	grid := engine.SpriteGridFromPNGCleanKids(renderer, path, cols, rows, npcSpriteInset)
	var frames []npcFrame
	if row < len(grid) {
		for c := 0; c < cols && c < len(grid[row]); c++ {
			gf := grid[row][c]
			frames = append(frames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH,
				fcx: gf.FCX, fry: gf.FRY, srcPath: path})
		}
	}
	for len(frames) > 1 && frames[len(frames)-1].tex == nil {
		frames = frames[:len(frames)-1]
	}
	stabilizeNPCAnchors(frames)
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

// higginsPostMarcusHealedDialog plays after Marcus has been healed by the
// Louvre postcard. It is the narrative BRIDGE into Jake's chapter (Jerusalem):
// Marcus's heal wakes Jake into the strange state and lights up Jerusalem on
// the map, so Higgins must point PP at Jake next — NOT Lily (that's a later
// chapter). User playtest 2026-06-05 (#39): "we need to start talking about
// Jake... fill it with more lines."
var higginsPostMarcusHealedDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Marcus is finally sleeping soundly. Whatever you brought back from Paris, it worked. Good work, Panther."},
	{speaker: "Director Higgins", text: "But I'm afraid it's not over. The moment Marcus settled... Jake started up."},
	{speaker: "Pink Panther", text: "Jake? The tough kid who never says much?"},
	{speaker: "Director Higgins", text: "That's the one. Now he won't stop. Muttering about ancient tunnels, a great stone wall, coins buried under the city."},
	{speaker: "Director Higgins", text: "He keeps scratching the same symbol into the dirt. I've never seen anything like it. I don't understand any of this, Panther."},
	{speaker: "Pink Panther", text: "A wall, old coins, tunnels under a city... that sounds like Jerusalem."},
	{speaker: "Director Higgins", text: "Then that's where you're headed. The travel map lit up Jerusalem on its own — same as it did Paris for Marcus."},
	{speaker: "Director Higgins", text: "Go talk to Jake first, in his cabin — see what he's fixated on. Then take the map and find whatever he's missing. The kids are counting on us."},
}

func newDirectorHiggins(renderer *sdl.Renderer) *npc {
	// Bounds sized to 200x265 so the aspect-preserve draw produces
	// ~225-235 px of actual sprite on camp_entrance — matches the
	// "adult NPC" row in CHARACTERS.md (PP is 170x235 for reference).
	// Do not shrink below 200x260 or Higgins reads as a kid.
	//
	// Both sheets are clean grids per PROMPTS.md:
	//   idle: 6x1
	//   talk: 6x1 (clipboard lowered, mouth open)
	return &npc{
		idleGrid: loadNPCGridClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 6, 1),
		talkGrid: loadNPCGridRowClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png", 8, 2, 0),
		// User 2026-05-18: shifted X 660 → 760 so PP's walk-up-to-talk
		// position lands clear of the left gate post / fence rail. PP
		// resting spot (post walk-in) also shifted to keep the same gap.
		bounds:         sdl.Rect{X: 760, Y: 390, W: 168, H: 220},
		name:           "Director Higgins",
		dialog:         higginsDefaultDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.18,
	}
}

func newOfficeHiggins(renderer *sdl.Renderer) *npc {
	// Office Higgins bounds were 180x280 which rendered him at ~35% of
	// screen height — too tall vs the PTP reference. Dropped to 160x225
	// to put him in the 210-225 band from CHARACTERS.md; camp_office's
	// characterScale 0.9 shaves the final render to ~200 which sits
	// comfortably in the tight indoor shot.
	n := &npc{
		idleGrid: loadNPCGridRowClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png", 6, 2, 0),
		talkGrid: loadNPCGridClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png", 6, 2),
		// User spec 2026-04-17: office Higgins top-left at (1062, 357),
		// sitting behind the desk. Sized so head lands at ~y=357 and feet
		// rest on the desk chair around y=640. 2026-05-12 (revised after
		// screenshot showed NPCs dwarfing the bg): moderate scale instead
		// of the full retro-proportion bump.
		// User 2026-05-19: anchor at (1106, 480) per playtest — Higgins
		// "sits" with his torso/head visible at this top-left. Foot now
		// at y=715 (below PP foot max 665); PP walks up to him and
		// stands in front of the desk.
		// User 2026-05-20: move to (1091, 365) so Higgins is framed lower
		// in the desk window; PP also needs grounding (camp_office.json
		// spawnY adjusted) and the back-arrow flipped to "left".
		// User 2026-05-21: refined to "sitting behind the desk" pose.
		// Desk surface is at y≈490 in the office BG, chair centered ~1015-1180.
		// New bounds (990, 290, 220, 200) → top at 290 (head clearly above
		// desk), foot at 490 (sprite bottom rests on desk surface, lower body
		// naturally clipped by desk art). Aspect-preserve renders him as
		// 129×200 centered horizontally in the 220 bounds.
		// User 2026-05-23: Y nudged 310 → 300 (a few pixels up so the
		// head sits cleaner above the desk per the playtest report).
		// User 2026-06-02 (#16): nudged up 300 → 280 so the head sits higher in
		// the desk window.
		bounds:         sdl.Rect{X: 990, Y: 280, W: 220, H: 200},
		name:           "Director Higgins",
		dialog:         higginsWorriedDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.08,
		// User playtest #12: flip office Higgins 180° (he now faces LEFT toward
		// PP / the desk). fixedFacing keeps him from turning back toward PP each
		// time a dialog starts, so he holds the flipped orientation.
		flipped:     true,
		fixedFacing: true,
		silent:      true,
	}
	// Register the give-map one-shot animation. User 2026-05-31 (#14): the
	// sheet is a 6×2 grid (detect_grid), not 8×1 — cutting it 8×1 made cellH
	// span BOTH rows so each frame drew Higgins twice, one above the other.
	giveFrames := loadNPCGridClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_give_map.png", 6, 2)
	if len(giveFrames) > 0 {
		n.oneShotAnims = map[string][]npcFrame{"give_map": giveFrames}
	}
	return n
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
	h.bounds = sdl.Rect{X: 1060, Y: 560, W: 180, H: 210}
	h.hidden = true
	h.silent = true
	h.dialog = higginsLilyHintDialog
	// Register the back-walk one-shot used by the higgins_walk_in sequence
	// when he enters from the right edge after Lily's shy dialog. PNG is
	// 1376×768; load as 8×2 take_row=0 to mirror the talk sheet's geometry.
	walkBackFrames := loadNPCGridRowClean(renderer,
		"assets/images/locations/camp/npc/higgins/npc_director_higgins_walk_back.png",
		8, 2, 0)
	if len(walkBackFrames) > 0 {
		h.oneShotAnims = map[string][]npcFrame{"walk_back": walkBackFrames}
	}
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
	n.bounds = sdl.Rect{X: 760, Y: 445, W: 162, H: 245}
	n.silent = true
	return n
}

func newRoomJake(renderer *sdl.Renderer) *npc {
	n := newJake(renderer)
	n.bounds = sdl.Rect{X: 760, Y: 435, W: 162, H: 245}
	n.silent = true
	return n
}

func newRoomLily(renderer *sdl.Renderer) *npc {
	n := newLily(renderer)
	n.bounds = sdl.Rect{X: 666, Y: 476, W: 162, H: 245}
	n.silent = true
	return n
}

func newRoomMarcus(renderer *sdl.Renderer) *npc {
	// User feedback 2026-04-26: room Marcus was reading huge — bounds 280x380
	// + characterScale 0.85 still rendered him oversize next to PP. Shrunk to
	// 200x300. 2026-05-12 (revised): aligned with the moderate global scale
	// so Marcus matches PP's 270 height (he's the freakout-giant silhouette
	// but the room-internal composition can't take him much bigger).
	n := newMarcus(renderer)
	// User 2026-05-19: Y 290 → 350 so Marcus's foot drops to y=620
	// (cabin floor line) instead of mid-room.
	// User 2026-05-20: nudge down another 35px so feet rest on the cabin
	// floor instead of hovering above it.
	// User playtest #10: room Marcus read as "way bigger than PP". The old
	// 187×270 matched PP's height (the "looming freakout giant" intent), but the
	// user wants him clearly shorter than PP. Shrunk to 150×205 with the foot
	// kept on the cabin floor line (Y+H ≈ 655), so he now reads as a kid.
	n.bounds = sdl.Rect{X: 615, Y: 450, W: 150, H: 205}
	// Hidden until the night freakout cutscene unhides him. Without this,
	// peeking into Marcus's cabin on Day 1 already shows him there even
	// though Day-1 Marcus belongs on the camp grounds.
	n.hidden = true
	// User 2026-05-22: load the alt-idle "strange_alt" frames so the
	// engine swaps Marcus's idle for one cycle every ~5 seconds of
	// inactivity in the cabin scene. Ambient "freakout punctuation".
	// User 2026-06-02: this sheet has 8 poses per row, not 7 — loading 7×2
	// drifted each cell off its pose and sliced two half-Marcuses into one
	// frame ("two frames at once" during the strange idle). Load 8×2 to match
	// the regenerated uniform sheet.
	altFrames := loadNPCGrid(renderer, "assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_alt.png", 8, 2)
	if len(altFrames) > 0 {
		n.altIdleGrid = altFrames
		n.altIdleAfterSec = 5.0
	}
	return n
}

func newRoomDanny(renderer *sdl.Renderer) *npc {
	n := newDanny(renderer)
	n.bounds = sdl.Rect{X: 760, Y: 445, W: 162, H: 245}
	n.silent = true
	return n
}

// newNightHiggins is the campfire Higgins — silent by default so he doesn't
// block exploration, but driven directly by the night cutscene so he appears
// to deliver the "lights out" speech in-place, not at camp grounds.
func newNightHiggins(renderer *sdl.Renderer) *npc {
	n := &npc{
		// Idle sheet is 2304px wide = 6 frames of 384 (2304/7 is not whole);
		// loading 7×1 sliced mid-character and slid Higgins sideways.
		idleGrid:       loadNPCGridClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 6, 1),
		talkGrid:       loadNPCGridRowClean(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png", 8, 2, 0),
		bounds:         sdl.Rect{X: 1120, Y: 430, W: 172, H: 220},
		name:           "Director Higgins",
		bobAmount:      0,
		talkFrameSpeed: 0.18,
		silent:         true,
	}
	// Register "shout" one-shot so night_bedtime can play the angry-lights-out
	// frames instead of the default talk cycle. User 2026-05-22: previous
	// attempts used loadNPCGridClean (tighter inset + tighter chroma-key)
	// which produced 0 frames — log + fallback to loadNPCGrid (lenient).
	// User 2026-05-31 (#9): shout sheet is 8×2, not 8×1 — loading 8×1 stacked
	// both rows into each cell so the swap played a garbled double-Higgins and
	// read as "not shouting". Load 8×2.
	shoutFrames := loadNPCGrid(renderer, "assets/images/locations/camp/npc/higgins/npc_director_higgins_shout.png", 8, 2)
	fmt.Printf("[newNightHiggins] shout frames loaded: %d\n", len(shoutFrames))
	if len(shoutFrames) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["shout"] = shoutFrames
		// User playtest #7: the shout still "wasn't activating" because the
		// night sequence's dialog step put Higgins into his normal TALK anim,
		// overriding the one-shot idle-swap. The night Higgins exists ONLY to
		// bellow "lights out", so just make his idle AND talk the shout frames —
		// now he's shouting whether the sequence has him idle or talking.
		n.idleGrid = shoutFrames
		n.talkGrid = shoutFrames
	}
	return n
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
	// User 2026-05-23: reverted to 145-wide click rect (X=130). The earlier
	// W-shrink-to-100 left the rect too narrow — depending on which animation
	// frame is showing, the kid's body extends past the trimmed bounds and
	// clicks miss. 145 wide gives a forgiving target while still hugging
	// the visible character.
	n := &npc{
		bounds:         sdl.Rect{X: 130, Y: 465, W: 145, H: 120}, // #5: kids ~55% of PP (feet kept at y=585)
		name:           "Tommy",
		dialog:         tommyDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	applyKidAtlasOrFallback(renderer, n, "tommy")
	// User 2026-05-21: register "walk_left" one-shot so the tommy_exit
	// sequence can swap idle for the walking sheet during the move.
	// The PNG ships as 1536×1024 with kid content in the MIDDLE band
	// (rows 324-678 — about 35% of canvas height); the rest is empty
	// white padding. Loading as 8×1 would give 192×1024 cells where the
	// kid takes only ~35% — engine's aspect-preserve renders him at
	// ~37px wide, "very very small" per user. Loading as 8×3 take_row=1
	// gives 192×341 cells centered on the middle band, so the kid fills
	// the cell and renders at ~112px wide instead — much closer to his
	// idle visual size. Full art regen tracked in EXTRA_PROMPTS §E.
	// User 2026-05-31 (#7): walk-left is a clean 8×1 full-body strip. The old
	// 8×3 take_row=1 grabbed only the middle band and CROPPED Tommy's head.
	// With opaque-box normalization the full 8×1 cell is trimmed to his body,
	// so he renders whole and at a good size.
	walkLeftFrames := loadNPCGridRow(renderer,
		"assets/images/locations/camp/npc/kids/tommy/npc_tommy_walk_left.png",
		8, 1, 0)
	if len(walkLeftFrames) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["walk_left"] = walkLeftFrames
	}
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
	// User 2026-05-23: reverted to 145-wide (see Tommy comment).
	n := &npc{
		bounds:         sdl.Rect{X: 395, Y: 460, W: 145, H: 120}, // #5: kids ~55% of PP (feet kept at y=580)
		name:           "Jake",
		dialog:         jakeDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	applyKidAtlasOrFallback(renderer, n, "jake")
	// User 2026-05-24: kid content is at PNG y=231-660 (1672×941 sheet).
	// Previous 8×3 take_row=1 (cell y=313-627) chopped the top of his
	// head (82 px above cell top) and bottom of feet — user reported
	// "head is cutted". Going back to 8×1 (full-cell 209×941) so the
	// whole kid fits, even though the render is narrower; the head and
	// feet are now both visible. Final regen tracked in EXTRA_PROMPTS §F.
	walkBackFrames := loadNPCGrid(renderer,
		"assets/images/locations/camp/npc/kids/jake/npc_jake_walk_back.png",
		8, 1)
	if len(walkBackFrames) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["walk_back"] = walkBackFrames
	}
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
	// User 2026-05-23: reverted to 145-wide (see Tommy comment).
	n := &npc{
		bounds:         sdl.Rect{X: 600, Y: 440, W: 145, H: 120}, // #5: Lily ~55% of PP (feet kept at y=560)
		name:           "Lily",
		dialog:         lilyShyDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	applyKidAtlasOrFallback(renderer, n, "lily")
	// Receive-flower one-shot played when PP hands her the flower. Sheet
	// is 1024×128 = 8×1 single-row strip per the file dims (cells 128×128).
	receiveFlower := loadNPCGrid(renderer,
		"assets/images/locations/camp/npc/kids/lily/npc_lily_receive_flower.png",
		8, 1)
	if len(receiveFlower) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["receive_flower"] = receiveFlower
	}
	// #4: "talking after getting the flower" — Lily holds the daisy while she
	// talks once PP has handed it over. Loaded here, swapped into talkGrid by
	// the flower handoff callback in setupCampCallbacks. 8×2 like her other sheets.
	if withFlower := loadNPCGridConnected(renderer,
		"assets/images/locations/camp/npc/kids/lily/npc_lily_talk_with_flower.png", 8, 2); len(withFlower) > 0 {
		n.postGiveTalkGrid = withFlower
	}
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
	// User 2026-05-23: reverted to 145-wide (see Tommy comment).
	n := &npc{
		bounds:         sdl.Rect{X: 890, Y: 455, W: 145, H: 120}, // #5: kids ~55% of PP (feet kept at y=575)
		name:           "Marcus",
		dialog:         marcusDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		// Freakout feels frantic if it runs at normal talk speed — slow it
		// down so the strange dialogue has room to breathe.
		strangeTalkFrameSpeed: 0.16,
	}
	applyKidAtlasOrFallback(renderer, n, "marcus")
	// JIT regen (2026-06-10) restored these to clean 8x2 sheets; keep the
	// explicit reload so packed-atlas fallbacks cannot reuse the old 7-column cut.
	marcusDir := "assets/images/locations/camp/npc/kids/marcus/"
	if frames := loadNPCGridConnected(renderer, marcusDir+"npc_marcus_idle.png", 8, 2); len(frames) > 0 {
		n.idleGrid = frames
	}
	if frames := loadNPCGridConnected(renderer, marcusDir+"npc_marcus_talk.png", 8, 2); len(frames) > 0 {
		n.talkGrid = frames
	}
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
	// User 2026-05-23: third iteration on Danny click rect — user still
	// reports "clicked on him, nothing happen; clicked right, worked".
	// Going extra-generous: W=180 (vs the kid baseline 145) and X=1090 so
	// the rect spans 1090-1270 — covers the visible kid no matter which
	// animation frame is showing AND a forgiveness margin on both sides.
	// The NPC > hotspot click priority (set in HandleClick) means Danny's
	// dialog wins over both Lily-cabin (1017-1137) and Danny-cabin
	// (1183-1303) when click lands in the overlap zone.
	n := &npc{
		bounds:         sdl.Rect{X: 1110, Y: 460, W: 160, H: 120}, // #5/#7: smaller; shifted right (1110-1270) to clear Marcus so PP doesn't stand on him
		name:           "Danny",
		dialog:         dannyDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		flipped:        true,
		approachRight:  true, // #7: stand to Danny's right, not on Marcus
	}
	applyKidAtlasOrFallback(renderer, n, "danny")
	// User playtest #8: the wider "kids"/clean key was stripping Danny's WHITE
	// EYES along with the background. Use the connected-edge key (only removes
	// background pixels touching the cell edge) so his interior whites — eyes,
	// teeth — survive. (If a soft halo returns, fix it in the art with off-white.)
	dannyIdle := "assets/images/locations/camp/npc/kids/danny/npc_danny_idle.png"
	dannyTalk := "assets/images/locations/camp/npc/kids/danny/npc_danny_talk.png"
	if _, err := os.Stat(dannyIdle); err == nil {
		if f := loadNPCGridConnected(renderer, dannyIdle, 8, 2); len(f) > 0 {
			n.idleGrid = f
		}
	}
	if _, err := os.Stat(dannyTalk); err == nil {
		if f := loadNPCGridConnected(renderer, dannyTalk, 8, 2); len(f) > 0 {
			n.talkGrid = f
		}
	}
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

// playOneShotAnim starts a named non-looping animation registered under
// n.oneShotAnims[name]. duration is wall-clock seconds; the sequence player
// is expected to call endOneShotAnim when its own timer expires. If the
// requested anim isn't registered, this is a silent no-op so the sequence
// keeps moving.
func (n *npc) playOneShotAnim(name string, duration float64) {
	if n.oneShotAnims == nil {
		return
	}
	if _, ok := n.oneShotAnims[name]; !ok {
		return
	}
	n.activeOneShot = name
	n.oneShotIdx = 0
	n.oneShotTimer = 0
	n.oneShotDuration = duration
}

func (n *npc) endOneShotAnim() {
	n.activeOneShot = ""
	n.oneShotIdx = 0
	n.oneShotTimer = 0
}

// swapIdleForOneShot temporarily replaces idleGrid with the frames of a
// registered one-shot animation, so the existing idle frame cycler loops it
// at natural pace. Use this for looping named animations like Higgins's
// "walk_back" during an npc_move — the one-shot pathway alone freezes at
// the last frame (it's authored as fire-and-forget). Restored by
// restoreSwappedIdle() on the next idle/talk anim step.
func (n *npc) swapIdleForOneShot(name string) {
	if n.oneShotAnims == nil {
		return
	}
	frames, ok := n.oneShotAnims[name]
	if !ok || len(frames) == 0 {
		return
	}
	if n.swappedIdleBackup == nil {
		n.swappedIdleBackup = n.idleGrid
	}
	n.idleGrid = frames
	n.idleCurFrame = 0
	n.idleFrameTimer = 0
	n.animState = npcAnimIdle
}

// restoreSwappedIdle undoes swapIdleForOneShot if a swap is active.
func (n *npc) restoreSwappedIdle() {
	if n.swappedIdleBackup == nil {
		return
	}
	n.idleGrid = n.swappedIdleBackup
	n.swappedIdleBackup = nil
	n.idleCurFrame = 0
	n.idleFrameTimer = 0
}

func (n *npc) update(dt float64) {
	n.bobTimer += dt

	// One-shot anim (e.g. Higgins's give_map) overrides idle/talk while
	// active. Frames advance at the standard talkFrameSpeed; the sequence
	// player owns end-of-anim cleanup via endOneShotAnim.
	if n.activeOneShot != "" {
		frames := n.oneShotAnims[n.activeOneShot]
		n.oneShotTimer += dt
		stepLen := 0.12
		if n.talkFrameSpeed > 0 {
			stepLen = n.talkFrameSpeed
		}
		// Spread frames evenly across the anim's wall-clock duration so a
		// short asset (6 frames) over a long timeline (1.4 s) doesn't loop
		// twice. Falls back to talkFrameSpeed if duration unset.
		if n.oneShotDuration > 0 && len(frames) > 0 {
			stepLen = n.oneShotDuration / float64(len(frames))
		}
		if n.oneShotTimer >= stepLen && len(frames) > 0 {
			n.oneShotTimer -= stepLen
			if n.oneShotIdx < len(frames)-1 {
				n.oneShotIdx++
			}
		}
	}

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
		// User 2026-05-31 (#4/#13): ×2.5 (≈0.375s/frame) read as a slow,
		// visible "swish" between idle frames. ×2.0 (≈0.30s) is smoother.
		idleSpeed := speed * 2.0 // idle cycles a little slower than talk
		// User 2026-06-02 (#15): the strange/freakout idle read as "way too
		// fast" — the poses jump hard so even the normal cadence strobes. Slow
		// the strange idle (and its alt-idle beat) right down so it reads as an
		// uneasy fidget, not a flicker.
		if n.isStrange {
			idleSpeed = speed * 3.5
		}
		// User 2026-05-20: when a walk/named-anim is swapped into idleGrid
		// (via swapIdleForOneShot), cycle at walk-cadence (~0.10s) instead
		// of the slow idle cadence. Higgins's walk_back was previously
		// playing at 0.45s/frame which made the 8-frame cycle drag and
		// look unsmooth during the 1.8s walk_in move.
		if n.swappedIdleBackup != nil {
			idleSpeed = 0.10
		}
		if n.idleFrameTimer >= idleSpeed {
			n.idleFrameTimer -= idleSpeed
			n.idleCurFrame = (n.idleCurFrame + 1) % len(n.idleGrid)
		}
	}

	// User 2026-05-22: inactivity alt-idle swap. While the NPC is idling
	// (not talking, not in a sequence-driven swap, not currently in an alt
	// cycle), accumulate dt. When threshold passes, swap idleGrid for the
	// altIdleGrid for ONE full cycle, then restore. Reset accumulator on
	// talk-start (setAnimState calls reset elsewhere; safety-reset here on
	// state change too). Marcus uses this for the ambient strange_alt beat.
	if n.altIdleAfterSec > 0 && len(n.altIdleGrid) > 0 &&
		n.animState == npcAnimIdle && n.activeOneShot == "" && n.swappedIdleBackup == nil {
		if n.altIdleActive {
			// In an alt cycle — check if we've played one full loop already.
			// idleCurFrame returns to 0 once per cycle; flip back to normal.
			if n.idleCurFrame == 0 && n.idleFrameTimer < dt*1.5 {
				// freshly wrapped to frame 0 — restore normal idle
				if n.altIdleBackup != nil {
					n.idleGrid = n.altIdleBackup
					n.altIdleBackup = nil
				}
				n.altIdleActive = false
				n.idleAccumSec = 0
			}
		} else {
			n.idleAccumSec += dt
			if n.idleAccumSec >= n.altIdleAfterSec {
				// Trigger the alt cycle now.
				n.altIdleBackup = n.idleGrid
				n.idleGrid = n.altIdleGrid
				n.idleCurFrame = 0
				n.idleFrameTimer = 0
				n.altIdleActive = true
			}
		}
	}

	if n.animState == npcAnimTalk && len(n.talkGrid) > 0 {
		// #2: while a dialog is active, this NPC's mouth animates only on ITS
		// lines (speaker != Pink Panther) and only while the line is still
		// revealing — so the mouth tracks the words and holds closed (frame 0)
		// during PP's lines or once the text is fully shown. If no dialog is
		// active (e.g. a sequence-driven talk pose), keep the old free-run.
		speaking := true
		if n.game != nil && n.game.dialog != nil && n.game.dialog.active {
			ds := n.game.dialog
			speaking = ds.isRevealing() && ds.currentSpeaker() != "Pink Panther"
		}
		if speaking {
			n.frameTimer += dt
			if n.frameTimer >= speed {
				n.frameTimer -= speed
				n.curFrame = (n.curFrame + 1) % len(n.talkGrid)
			}
		} else {
			n.frameTimer = 0
			n.curFrame = 0
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
// activeFrames returns the frame slice currently playing (one-shot / talk /
// idle), used to compute the animation's reference height.
func (n *npc) activeFrames() []npcFrame {
	if n.activeOneShot != "" {
		if frames, ok := n.oneShotAnims[n.activeOneShot]; ok && len(frames) > 0 {
			return frames
		}
	}
	if n.animState == npcAnimTalk && len(n.talkGrid) > 0 {
		return n.talkGrid
	}
	return n.idleGrid
}

// maxOpaqueH is the tallest opaque content height across a frame slice. Scaling
// every frame so this maps to bounds.H keeps the character one consistent size
// across the animation (tallest pose fills the box; shorter poses keep planted
// feet) instead of pulsing.
func maxOpaqueH(frames []npcFrame) int32 {
	var m int32
	for _, f := range frames {
		if f.oh > m {
			m = f.oh
		}
	}
	return m
}

// stabilizeNPCAnchors applies a DEADBAND to every frame's feet anchors (fcx +
// fry): values within ±6px of the animation median snap to the median —
// killing foot-detection noise so well-aligned frames are rock-stable — while
// larger deviations keep the frame's OWN feet position, compensating art that
// genuinely drifts inside the cells (user 2026-06-10: "the frames place in
// the same spot"; a constant median anchor made drifting sheets jump).
func stabilizeNPCAnchors(frames []npcFrame) {
	const deadband = 6
	cxs := make([]int, 0, len(frames))
	frys := make([]int, 0, len(frames))
	for i := range frames {
		if frames[i].ow <= 0 || frames[i].oh <= 0 {
			continue
		}
		if frames[i].fcx <= 0 {
			frames[i].fcx = frames[i].ox + frames[i].ow/2
		}
		if frames[i].fry <= frames[i].oy {
			frames[i].fry = frames[i].oy + frames[i].oh
		}
		cxs = append(cxs, int(frames[i].fcx))
		frys = append(frys, int(frames[i].fry))
	}
	if len(cxs) == 0 {
		return
	}
	sort.Ints(cxs)
	sort.Ints(frys)
	medCX := int32(cxs[len(cxs)/2])
	medFRY := int32(frys[len(frys)/2])
	for i := range frames {
		if frames[i].ow <= 0 || frames[i].oh <= 0 {
			continue
		}
		if d := frames[i].fcx - medCX; d >= -deadband && d <= deadband {
			frames[i].fcx = medCX
		}
		if d := frames[i].fry - medFRY; d >= -deadband && d <= deadband {
			frames[i].fry = medFRY
		}
	}
}

func (n *npc) drawScaled(renderer *sdl.Renderer, charScale float64) {
	if n.hidden {
		return
	}
	if charScale <= 0 {
		charScale = 1.0
	}
	if n.extraScale > 0 {
		charScale *= n.extraScale
	}
	bobOffset := int32(math.Sin(n.bobTimer*1.5) * n.bobAmount)

	shadowCX := n.bounds.X + n.bounds.W/2
	shadowFY := n.bounds.Y + n.bounds.H
	drawShadow(renderer, shadowCX, shadowFY, int32(float64(n.bounds.W-10)*charScale))

	flip := sdl.FLIP_NONE
	if n.flipped {
		flip = sdl.FLIP_HORIZONTAL
	}

	frames := n.activeFrames()
	var frame npcFrame
	if n.activeOneShot != "" {
		if fr, ok := n.oneShotAnims[n.activeOneShot]; ok && len(fr) > 0 {
			frame = fr[n.oneShotIdx%len(fr)]
		}
	} else if n.animState == npcAnimTalk && len(n.talkGrid) > 0 {
		frame = n.talkGrid[n.curFrame%len(n.talkGrid)]
	} else if len(n.idleGrid) > 0 {
		frame = n.idleGrid[n.idleCurFrame%len(n.idleGrid)]
	}

	if frame.tex == nil {
		return
	}

	var dstW, dstH, dstX, dstY int32
	var src *sdl.Rect

	// Opaque-box normalization (the size/cut fix): scale so the animation's
	// tallest opaque pose fills bounds.H, then anchor by feet (bottom) +
	// horizontal centre. Makes idle/talk/walk render at one consistent size
	// and stops empty-padding cells rendering tiny / head-cropped. Café patrons
	// (srcCropBottomFrac>0) keep the legacy top-anchored crop path below.
	if frame.ow > 0 && frame.oh > 0 && n.srcCropBottomFrac == 0 {
		refH := maxOpaqueH(frames)
		if refH <= 0 {
			refH = frame.oh
		}
		scale := float64(n.bounds.H) * charScale / float64(refH)
		base := sdl.Rect{X: 0, Y: 0, W: frame.w, H: frame.h}
		if frame.src != nil {
			base = *frame.src
		}
		s := sdl.Rect{X: base.X + frame.ox, Y: base.Y + frame.oy, W: frame.ow, H: frame.oh}
		src = &s
		dstW = int32(float64(frame.ow) * scale)
		dstH = int32(float64(frame.oh) * scale)
		// Anchor by the frame's FEET (per-frame fcx/fry, deadband-snapped to
		// the animation median by stabilizeNPCAnchors). Every frame plants
		// its feet on the same screen spot: art drift inside the cells is
		// compensated per frame, while a gesturing arm or a tail dipping
		// below the feet extends naturally past the anchor (user 2026-06-10:
		// "the frames place in the same spot").
		fcx := frame.fcx
		if fcx <= 0 {
			fcx = frame.ox + frame.ow/2
		}
		fry := frame.fry
		if fry <= frame.oy {
			fry = frame.oy + frame.oh
		}
		anchorX := float64(n.bounds.X) + float64(n.bounds.W)/2
		colFromLeft := (float64(fcx) - float64(frame.ox)) * scale
		if n.flipped {
			// CopyEx mirrors within the dst rect — mirror the anchor too.
			dstX = int32(anchorX - (float64(dstW) - colFromLeft))
		} else {
			dstX = int32(anchorX - colFromLeft)
		}
		footLine := n.bounds.Y + n.bounds.H + bobOffset
		dstY = footLine - int32((float64(fry)-float64(frame.oy))*scale)
	} else {
		breathScale := 1.0
		targetW := float64(n.bounds.W) * charScale
		targetH := float64(n.bounds.H) * charScale
		scaleW := targetW * breathScale / float64(frame.w)
		scaleH := targetH * breathScale / float64(frame.h)
		scale := scaleW
		if scaleH < scale {
			scale = scaleH
		}
		dstW = int32(float64(frame.w) * scale)
		dstH = int32(float64(frame.h) * scale)
		dstX = n.bounds.X + (n.bounds.W-dstW)/2
		dstY = n.bounds.Y + bobOffset + (n.bounds.H - dstH)
		src = frame.src
		if n.srcCropBottomFrac > 0 && n.srcCropBottomFrac < 1.0 {
			baseSrc := frame.src
			if baseSrc == nil {
				baseSrc = &sdl.Rect{X: 0, Y: 0, W: int32(frame.w), H: int32(frame.h)}
			}
			newH := int32(float64(baseSrc.H) * n.srcCropBottomFrac)
			src = &sdl.Rect{X: baseSrc.X, Y: baseSrc.Y, W: baseSrc.W, H: newH}
			dstH = int32(float64(dstH) * n.srcCropBottomFrac)
		}
	}

	dst := sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH}
	renderer.CopyEx(frame.tex, src, &dst, 0, nil, flip)
	// lastDrawRect now hugs the rendered body (opaque-anchored), so the click
	// hit-test (containsPoint) uses it + a small margin.
	n.lastDrawRect = dst
	n.lastDrawnFrame = frame
	n.lastDrawnFlip = n.flipped
}

// containsPoint is used for both cursor hover (showing the "talk" icon) and
// actual click detection. Keeping them unified means: wherever the cursor
// shows "talk", a click always lands.
//
// User 2026-05-24: hybrid hit-test. Past iterations toggled between
// "use bounds rect" (too wide — click lands in empty rect space, user
// reported "I had to click to the right of every NPC to talk") and
// "use lastDrawRect" (too narrow — edge-only clickable). Hybrid:
// expand lastDrawRect by ±25 px horizontally and ±15 px vertically
// as forgiveness, intersected with the authored bounds rect so the
// click region never extends past the design-time max. Falls back
// to bounds when lastDrawRect isn't set yet (first frame).
func (n *npc) containsPoint(x, y int32) bool {
	pt := sdl.Point{X: x, Y: y}
	if n.lastDrawRect.W <= 0 || n.lastDrawRect.H <= 0 {
		return pt.InRect(&n.bounds)
	}
	// Expand the actual draw rect by forgiveness padding.
	const padX = int32(25)
	const padY = int32(15)
	hit := sdl.Rect{
		X: n.lastDrawRect.X - padX,
		Y: n.lastDrawRect.Y - padY,
		W: n.lastDrawRect.W + 2*padX,
		H: n.lastDrawRect.H + 2*padY,
	}
	// Don't let the expanded rect escape the authored bounds (keeps
	// click region predictable for NPCs whose bounds are intentionally
	// shrunk like Pierre back-of-line).
	if hit.X < n.bounds.X {
		hit.X = n.bounds.X
	}
	if hit.Y < n.bounds.Y {
		hit.Y = n.bounds.Y
	}
	if hit.X+hit.W > n.bounds.X+n.bounds.W {
		hit.W = n.bounds.X + n.bounds.W - hit.X
	}
	if hit.Y+hit.H > n.bounds.Y+n.bounds.H {
		hit.H = n.bounds.Y + n.bounds.H - hit.Y
	}
	return pt.InRect(&hit)
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
// bakeryWomanLostPinDialog is the new initial beat — Madame Poulain has lost
// her rolling pin somewhere on the floor and won't bake until it's recovered.
// Replaces the old "free baguette on first click" beat so the Paris arc has
// a real intro puzzle (user 2026-04-26 retro-style rework).
var bakeryWomanLostPinDialog = []dialogEntry{
	{speaker: "Madame Poulain", text: "Mon dieu! Bonjour, monsieur — but I cannot serve you!"},
	{speaker: "Pink Panther", text: "What's wrong, madame?"},
	{speaker: "Madame Poulain", text: "My rolling pin! It rolled off ze counter and I cannot find it!"},
	{speaker: "Madame Poulain", text: "Without it, no dough, no bread, no baguette."},
	{speaker: "Pink Panther", text: "I'll take a look around."},
	{speaker: "Madame Poulain", text: "Merci, monsieur! Find it and ze first baguette is yours."},
}

// bakeryWomanPinTradeDialog fires once PP returns the rolling pin (altDialog).
var bakeryWomanPinTradeDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "I think I found what you were looking for, madame."},
	{speaker: "Madame Poulain", text: "My rolling pin! Bless you, monsieur!"},
	{speaker: "Madame Poulain", text: "Here — your baguette, fresh and warm. Tell Pierre I send my regards."},
}

var bakeryWomanPostDialog = []dialogEntry{
	{speaker: "Madame Poulain", text: "Enjoy ze baguette, monsieur! Zhere's a photographer near ze museum — he loves fresh bread."},
}

// bakeryWomanLouvreSouvenirDialog is the next-anchor beat: after Marcus
// is healed and PP returns to Paris, Madame Poulain asks him to bring
// her a Louvre postcard for her grandson back in Lyon. The postcard is
// the same `postcard` item that heals Marcus, so PP brings TWO. User
// 2026-05-20 story step forward.
var bakeryWomanLouvreSouvenirDialog = []dialogEntry{
	{speaker: "Madame Poulain", text: "Ah, ze panther returns! Tell me, did you visit ze Louvre?"},
	{speaker: "Pink Panther", text: "I did. Quite an experience."},
	{speaker: "Madame Poulain", text: "Mon petit-fils in Lyon, he collects postcards of ze museum."},
	{speaker: "Madame Poulain", text: "If you bring me one, I will send it to him. A small kindness, no?"},
	{speaker: "Pink Panther", text: "I'll see what I can do, madame."},
}

// bakeryWomanCoffeeRefillDialog — Poulain pours another café au lait when a
// quest needs one (Henri's trade still pending, or Lucien asleep on
// Camille's pencil) and PP isn't already carrying a cup.
var bakeryWomanCoffeeRefillDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Madame, could I trouble you for another café au lait?"},
	{speaker: "Madame Poulain", text: "For ze panther who found my rolling pin? Bien sûr!"},
	{speaker: "Madame Poulain", text: "Zere — hot and fresh. Don't let zis one get cold, hm?"},
}

// bakeryWomanHeelDialog — Pierre needs crumbs for the pigeon critics; Poulain
// donates the day-old baguette heel.
var bakeryWomanHeelDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Madame, Pierre needs crumbs. Pigeon business. It's a long story."},
	{speaker: "Madame Poulain", text: "Pierre and his pigeon critics! (laughs) Here — yesterday's baguette heel."},
	{speaker: "Madame Poulain", text: "Ze ends are for ze birds anyway. Tell him ze bakery expects a good review."},
}

// bakeryWomanSouvenirThanksDialog fires when PP hands over the signed
// postcard for her grandson (closes the post-Marcus souvenir loop).
var bakeryWomanSouvenirThanksDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "One Louvre postcard, madame — signed by ze curator himself."},
	{speaker: "Madame Poulain", text: "Signed?! Oh, mon petit-fils will not sleep for a week!"},
	{speaker: "Madame Poulain", text: "You have ze biggest heart in Paris, monsieur. From today, ze pink éclair in my window is called 'Le Panthère Rose'."},
	{speaker: "Pink Panther", text: "Fame at last. And it smells better than a press pass."},
}

var bakeryWomanSouvenirDoneDialog = []dialogEntry{
	{speaker: "Madame Poulain", text: "Ze postcard is already in ze mail to Lyon, monsieur. Merci, from both of us."},
}

func newBakeryWoman(renderer *sdl.Renderer) *npc {
	// Dedicated Bakery Woman sheet (see docs/EXTRA_PROMPTS.md §8). 8×2
	// canvas: row 0 = idle (mouth closed), row 1 = talk (mouth open).
	// Packed atlas at assets/sprites/paris/bakery_woman.(png|json) is the
	// preferred path; legacy per-row PNG slicing stays as a fallback so
	// the NPC still spawns if pack_atlas.py hasn't been run.
	n := &npc{
		// User 2026-05-31: Y=182 placed her up in the shelves/menu area where
		// she blended into the busy BG ("not in the game"). Moved her DOWN to
		// stand behind the counter glass: foot/waist-cutoff at y≈430 (counter
		// top), head at y≈250, centred behind the display case. Fine-tune with
		// in-game coords if needed.
		// User 2026-06-02 (#20): raise her (Y 250 → 215) so more of her sits
		// above the counter glass instead of sinking behind it.
		// User playtest #20/#21: reposition so she sits on the counter line at
		// y≈308 (behind the desk).
		bounds:         sdl.Rect{X: 605, Y: 308, W: 170, H: 180},
		name:           "Madame Poulain",
		dialog:         bakeryWomanLostPinDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		flipped:        false, // sheet draws her facing right already
	}
	if !applyNPCAtlas(renderer, n, "paris/bakery_woman") {
		// User 2026-05-31 (#19): use Madame Poulain's dedicated new sheets
		// (npc_madame_poulain_idle/_talk.png, full 8×2 each, behind-counter
		// upper-body pose). Fall back to the old combined npc_bakery_woman.png
		// (row0 idle / row1 talk) if the new sheets aren't present.
		// User playtest 2026-06-05: Poulain's sheets now live in the coffee/
		// folder (with the café patrons). Load idle/talk/work from there.
		poulainIdle := "assets/images/locations/paris/npc/coffee/npc_madame_poulain_idle.png"
		poulainTalk := "assets/images/locations/paris/npc/coffee/npc_madame_poulain_talk.png"
		if _, err := os.Stat(poulainIdle); err == nil {
			n.idleGrid = loadNPCGridConnected(renderer, poulainIdle, 8, 2)
			n.talkGrid = loadNPCGridConnected(renderer, poulainTalk, 8, 2)
			// User playtest #21: the work alt-idle wasn't showing. Trigger it
			// sooner (every ~3s idle) so she visibly kneads/works between chats.
			poulainWork := "assets/images/locations/paris/npc/coffee/npc_madame_poulain_work.png"
			if _, werr := os.Stat(poulainWork); werr == nil {
				if frames := loadNPCGridConnected(renderer, poulainWork, 8, 2); len(frames) > 0 {
					n.altIdleGrid = frames
					n.altIdleAfterSec = 3.0
				}
			}
		} else if _, err := os.Stat("assets/images/locations/paris/npc/npc_bakery_woman.png"); err == nil {
			// Legacy fallback only if the old combined sheet still exists —
			// guarded so a moved/deleted file can't panic the grid loader.
			const sheet = "assets/images/locations/paris/npc/npc_bakery_woman.png"
			n.idleGrid = loadNPCGridRow(renderer, sheet, 8, 2, 0)
			n.talkGrid = loadNPCGridRow(renderer, sheet, 8, 2, 1)
		}
	}
	// #25: Poulain handing the baguette over the counter — 8-frame give one-shot.
	// User playtest 2026-06-05: the give sheet moved to coffee/ with the rest of
	// Poulain's art; the old outside/ path loaded 0 frames so the hand-over
	// animation silently stopped playing.
	if f := loadNPCGridConnected(renderer, "assets/images/locations/paris/npc/coffee/npc_madame_poulain_give.png", 8, 1); len(f) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["give"] = f
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

// nicolasPencilHintDialog — Nicolas saw where Camille's pencil rolled
// ("Camille and the Sold-Out Postcard" street hop).
var nicolasPencilHintDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Nicolas — Camille lost her charcoal pencil here at sunrise. Did your lens catch it?"},
	{speaker: "Nicolas", text: "Ze lens catches everything, monsieur. It rolled off ze curb — straight into ze flower pot by ze Louvre steps."},
	{speaker: "Nicolas", text: "Ze pigeons have been guarding it like ze crown jewels. Good luck."},
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
		// User 2026-05-22: width 86 was an outlier vs the other Paris
		// front-line NPCs (Colette 135, Claude 115). Unified at 120×235.
		bounds:         sdl.Rect{X: 950, Y: 490, W: 120, H: 235},
		name:           "Nicolas",
		dialog:         pressPhotographerDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		flipped:        false, // sheet draws him facing right already
	}
	if !applyNPCAtlas(renderer, n, "paris/press_photographer") {
		const sheet = "assets/images/locations/paris/npc/outside/npc_press_photographer.png"
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
		// User 2026-05-22: unified Paris front-line NPCs at 120×235.
		// User 2026-05-22: talkFrameSpeed 0.10 → 0.08 for smoother
		// animation cadence (was choppy on Colette specifically).
		bounds:    sdl.Rect{X: 300, Y: 490, W: 120, H: 235},
		name:      "Madame Colette",
		dialog:    frenchGuideDialog,
		bobAmount: 0,
		// User 2026-05-31 (#21): her talk was too fast at 0.08; 0.13 slows the
		// cadence so it reads smoothly (#20).
		talkFrameSpeed: 0.13,
	}
	if !applyNPCAtlas(renderer, n, "paris/french_guide") {
		// User 2026-05-31 (#18): use Colette's dedicated new sheets
		// (npc_madame_colette_*) instead of the old generic french_guide art,
		// and load them CONNECTED (only edge-connected background is keyed) so
		// her light shirt keeps its colour instead of being chroma-keyed out.
		// User playtest 2026-06-05: Colette's sheets moved to the outside/ folder.
		coletteIdle := "assets/images/locations/paris/npc/outside/npc_madame_colette_idle.png"
		coletteTalk := "assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png"
		if _, err := os.Stat(coletteIdle); err == nil {
			n.idleGrid = loadNPCGridConnected(renderer, coletteIdle, 8, 2)
			// JIT regen restored Colette's talk sheet to a clean 8x2 grid.
			n.talkGrid = loadNPCGridConnected(renderer, coletteTalk, 8, 2)
		} else {
			n.idleGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/outside/npc_french_guide_idle.png", 8, 2)
			n.talkGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/outside/npc_french_guide_talk.png", 8, 2)
		}
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
	{speaker: "Curator Beaumont", text: "Perhaps your boy senses what was hidden. Ze gift shop sells a postcard of ze restored painting..."},
	{speaker: "Pink Panther", text: "Perfect. I'll take one."},
	{speaker: "Curator Beaumont", text: "...sold out, monsieur. Ze whole city wants one since ze news. New prints arrive in two weeks."},
	{speaker: "Curator Beaumont", text: "But! I keep ze last one in ze archive. Bring me a replica sketch of ze portrait for ze archive wall, and it is yours."},
	{speaker: "Curator Beaumont", text: "Mademoiselle Camille at ze café — ze fastest charcoal in Paris. Tell her Beaumont asks."},
}

// curatorWaitingDialog loops while PP owes Beaumont the replica sketch.
var curatorWaitingDialog = []dialogEntry{
	{speaker: "Curator Beaumont", text: "Ze archive wall waits for Camille's sketch, monsieur — and your postcard waits with it."},
}

// curatorSketchTradeDialog fires when PP hands over Camille's Sketch (altDialog).
var curatorSketchTradeDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "One Room 7 portrait — by ze fastest charcoal in Paris."},
	{speaker: "Curator Beaumont", text: "Magnifique... look at ze linework! Zis goes straight to ze archive wall."},
	{speaker: "Curator Beaumont", text: "Camille drew zis? Tell her ze Louvre may have a commission for her."},
	{speaker: "Curator Beaumont", text: "And as promised — ze last postcard of ze restored painting. For your young friend."},
	{speaker: "Curator Beaumont", text: "If he sees ze complete image, perhaps his mind will settle."},
}

var museumCuratorPostDialog = []dialogEntry{
	{speaker: "Curator Beaumont", text: "Ze postcard should help your young friend."},
	{speaker: "Curator Beaumont", text: "Ze mysteries of art connect us in ways we do not understand."},
}

// curatorSouvenirDialog fires after Madame Poulain asks for a postcard for
// her grandson (post-Marcus-heal souvenir loop). Beaumont signs a second one.
var curatorSouvenirDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Madame Poulain at ze bakery — her grandson in Lyon collects postcards of ze museum."},
	{speaker: "Curator Beaumont", text: "Madame Poulain! Her croissants kept zis museum running through ze '89 restoration."},
	{speaker: "Curator Beaumont", text: "And good timing — ze new prints arrived zis morning. For her grandson... zere. Signed by ze curator."},
	{speaker: "Curator Beaumont", text: "Every collection needs a rare piece."},
	{speaker: "Pink Panther", text: "You just made a small boy in Lyon very famous at school."},
}

func newMuseumCurator(renderer *sdl.Renderer) *npc {
	// Packed atlas at assets/sprites/paris/museum_curator.(png|json) is the
	// preferred path; legacy per-sheet PNG loading stays as a fallback.
	n := &npc{
		// User playtest #29: flip Beaumont 180° (face left toward PP entering
		// from the tunnel) and stand him at ~(546, 599) — bounds bottom (foot)
		// lands at Y+H=599. Scene characterScale 0.7 shrinks him so he reads as
		// "far down the gallery".
		bounds:         sdl.Rect{X: 546, Y: 359, W: 125, H: 240},
		name:           "Curator Beaumont",
		dialog:         museumCuratorDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
		flipped:        true,
	}
	if !applyNPCAtlas(renderer, n, "paris/museum_curator") {
		// User playtest 2026-06-05: curator sheets moved to museum/; the new
		// talk sheet is a clean 8×1 strip (was an uneven 4×2).
		n.idleGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/museum/npc_museum_curator_idle.png", 8, 1)
		n.talkGrid = loadNPCGrid(renderer, "assets/images/locations/paris/npc/museum/npc_museum_curator_talk.png", 8, 1)
	}
	if f := loadNPCGrid(renderer, "assets/images/locations/paris/npc/museum/npc_beaumont_give.png", 8, 1); len(f) > 0 {
		n.oneShotAnims = map[string][]npcFrame{"give": f}
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

// pierrePencilFavorDialog — after the baguette + confiture trades Pierre
// owes PP one (user 2026-06-10). The favor: the pigeons guarding Camille's
// pencil in the flower pot obey crumbs — and Pierre speaks fluent pigeon.
// Also seeds the Pigeon Critic gag ("anything for crumbs, except ze canvas").
var pierrePencilFavorDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Pierre — Camille's lucky pencil is in ze flower pot by ze steps, and ze pigeons won't let me near it. You owe me one, remember?"},
	{speaker: "Pierre", text: "A baguette AND confiture — oui, Pierre pays his debts. Watch zis, mon ami."},
	{speaker: "Pierre", text: "(he whistles and scatters a pinch of crumbs by his easel — ze pigeons abandon ze flower pot at once)"},
	{speaker: "Pierre", text: "Zey will do ANYTHING for crumbs... except land on my canvas. Critics! Ze pot is yours, monsieur."},
}

// --- "The Pigeon Critic" side quest (2026-06-10) ---
// After the press-pass trade, Pierre's running gag about pigeon critics
// becomes playable: his masterpiece is done but no pigeon will land to
// approve it. PP fetches the day-old Baguette Heel from Poulain, crumbs
// bring the critic, and Pierre rewards PP with a mini portrait.
var pierrePigeonAskDialog = []dialogEntry{
	{speaker: "Pierre", text: "Ze masterpiece is finished, mon ami... but look. Ze canvas, it is EMPTY of pigeons."},
	{speaker: "Pierre", text: "No critic will land! Twenty years on zis sidewalk and ze birds choose TODAY to boycott."},
	{speaker: "Pierre", text: "Crumbs, monsieur. Find me crumbs — Madame Poulain always has a stale heel for ze birds."},
}

var pierrePigeonWaitDialog = []dialogEntry{
	{speaker: "Pierre", text: "Ze critics circle, monsieur, but zey do not land. Crumbs! Ask Madame Poulain."},
}

// pierrePigeonLandDialog fires when PP hands over the Baguette Heel (altDialog).
var pierrePigeonLandDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "One day-old baguette heel. The bakery expects a good review."},
	{speaker: "Pierre", text: "Parfait! Now, crumble it... gently... by ze easel..."},
	{speaker: "Pierre", text: "(a plump pigeon flutters down, struts across ze cobblestones, and settles on ze corner of ze canvas)"},
	{speaker: "Pierre", text: "IT LANDED. Ze painting is approved! Monet himself never had a finer critic."},
	{speaker: "Pierre", text: "Monet painted outside like zis too — 'plein air', we call it. Ze light changes every minute, so you must paint fast."},
	{speaker: "Pierre", text: "Here — a little portrait of you, mon ami. Ze pigeon posed for ze background."},
}

var pierrePigeonDoneDialog = []dialogEntry{
	{speaker: "Pierre", text: "Ze critic returns every morning now. I pay him in crumbs — cheaper zhan ze newspapers."},
}

func newPierreArtist(renderer *sdl.Renderer) *npc {
	// Packed atlas at assets/sprites/paris/pierre_artist.(png|json) is the
	// preferred path; legacy per-row PNG slicing stays as a fallback.
	n := &npc{
		// Moved back in perspective (Y up, W/H shrunk ~25%) so Pierre stands
		// further down the Paris street line. PP's existing depthScale (driven
		// by player.y) shrinks PP automatically as he walks up to talk and
		// restores when he walks back to the front. User 2026-05-20.
		// User 2026-05-21: move left 40 + down 80 so Pierre is visually
		// grounded in the mid-distance cobblestones instead of floating.
		// User playtest #29: shrunk a little (95×175 → 84×156, foot kept at y=645).
		bounds:         sdl.Rect{X: 780, Y: 489, W: 84, H: 156},
		name:           "Pierre",
		dialog:         pierreArtistDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	if !applyNPCAtlas(renderer, n, "paris/pierre_artist") {
		const sheet = "assets/images/locations/paris/npc/outside/npc_art_vendor.png"
		n.idleGrid = loadNPCGridRow(renderer, sheet, 8, 2, 0)
		n.talkGrid = loadNPCGridRow(renderer, sheet, 8, 2, 1)
	}
	if f := loadNPCGrid(renderer, "assets/images/locations/paris/npc/outside/npc_pierre_pigeon_lands.png", 8, 1); len(f) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["pigeon"] = f
	}
	if f := loadNPCGrid(renderer, "assets/images/locations/paris/npc/outside/npc_pierre_give.png", 8, 1); len(f) > 0 {
		if n.oneShotAnims == nil {
			n.oneShotAnims = map[string][]npcFrame{}
		}
		n.oneShotAnims["give"] = f
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
		// User 2026-05-22: unified Paris front-line NPCs at 120×235.
		// User playtest #31: nudged down a little (Y 480 → 510).
		bounds:    sdl.Rect{X: 1180, Y: 510, W: 120, H: 235},
		name:      "Claude",
		dialog:    gendarmeDialog,
		bobAmount: 0,
		// User playtest #19: Claude's talk cycle flickered too fast. Slowed
		// from 0.10 to 0.16 s/frame so the mouth animation reads smoothly.
		talkFrameSpeed: 0.16,
	}
	if !applyNPCAtlas(renderer, n, "paris/gendarme_claude") {
		const sheet = "assets/images/locations/paris/npc/outside/npc_security_guard.png"
		n.idleGrid = loadNPCGridRow(renderer, sheet, 6, 2, 0)
		n.talkGrid = loadNPCGridRow(renderer, sheet, 6, 2, 1)
	}
	return n
}

// =====================================================================
// Café patrons (paris_bakery interior) — 6 seated NPCs at 3 tables.
//
// One of them (Henri) is on the main quest chain: he asks PP to fetch a
// coffee on first visit, then trades the coffee for homemade Confiture
// out of his bag — which PP needs to give to Pierre so Pierre will eat
// the baguette (otherwise the baguette is too dry and Pierre refuses
// the press-pass trade). The other 5 are pure flavor.
//
// Sheets: assets/images/locations/paris/npc/coffee/cafe_patron_<name>_<idle|talk>.png
// 8×1 strip per file, 100×170 per cell. Falls back gracefully if the
// PNG isn't on disk yet — the engine no-ops on missing textures.
// =====================================================================

// loadCafePatronGrids is a small helper to keep the 6 factories DRY.
//
// User 2026-05-22 update: prefer the SPLIT format (two files per patron) —
// `cafe_patron_<name>_idle.png` (800×170, 8×1) and `cafe_patron_<name>_talk.png`
// (800×170, 8×1) — because it lets ChatGPT regen idle or talk independently
// without re-rolling the whole 16-frame sheet. Falls back to the legacy
// combined `cafe_patron_<name>.png` (1376×768, 8×2 with row 0 idle + row 1
// talk) if the split files aren't on disk yet.
//
// Both paths use the CLEAN loader (tolerance 16) so off-white fringe pixels
// at clothing/cup edges chroma-key away. Logs help catch silent load failures.
func loadCafePatronGrids(renderer *sdl.Renderer, name string) ([]npcFrame, []npcFrame) {
	base := "assets/images/locations/paris/npc/coffee/cafe_patron_" + name
	// User 2026-05-31 (#20): the 2026-05 art ships one 8×1 sheet per state at
	// 1536×1024 — idle is "<name>_idle.png" if present else the bare
	// "<name>.png"; talk is "<name>_talking.png" (new) or "<name>_talk.png"
	// (legacy). The old loader looked for _idle/_talk only, so it fell into the
	// "combined 8×2" branch and loaded half-figures (patrons "not imported").
	// Connected color-key keeps interior whites (cups, collars).
	firstExisting := func(paths ...string) string {
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
		return ""
	}
	var idle, talk []npcFrame
	if p := firstExisting(base+"_idle.png", base+".png"); p != "" {
		idle = loadNPCGridConnected(renderer, p, 8, 1)
	}
	if p := firstExisting(base+"_talking.png", base+"_talk.png"); p != "" {
		talk = loadNPCGridConnected(renderer, p, 8, 1)
	}
	if len(talk) == 0 {
		talk = idle // fall back to idle so the patron still renders while talking
	}
	if len(idle) == 0 {
		fmt.Printf("[cafePatron %s] WARNING: no idle frames found under %s*\n", name, base)
	}
	return idle, talk
}

// --- Madame Yvette (beret + pearls, sipping tea) — flavor + foreshadow ---
var yvetteDialog = []dialogEntry{
	{speaker: "Madame Yvette", text: "Ze museum restoration is all anyone talks about, monsieur."},
	{speaker: "Madame Yvette", text: "A hidden symbol under ze portrait! Imagine — five hundred years and we still find new things."},
	{speaker: "Pink Panther", text: "Five hundred years is a long time to keep a secret."},
	{speaker: "Madame Yvette", text: "Ze gift shop sold out of ze restoration postcards in ONE day. Ze whole city wants one."},
}

func newCafePatronYvette(renderer *sdl.Renderer) *npc {
	idle, talk := loadCafePatronGrids(renderer, "yvette")
	return &npc{
		idleGrid: idle,
		talkGrid: talk,
		// User 2026-05-22: anchored to LEFT chair of left café table.
		// srcCropBottomFrac clips the lower body so only head+shoulders
		// render — chair art from the BG fills the lower half.
		bounds:         sdl.Rect{X: 80, Y: 339, W: 90, H: 160},
		name:           "Madame Yvette",
		dialog:         yvetteDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
}

// --- Monsieur Bernard (bearded, Le Figaro reader) — flavor ---
var bernardDialog = []dialogEntry{
	{speaker: "Monsieur Bernard", text: "(rustles paper) Le Figaro headline today, monsieur — restorer found a symbol under ze famous portrait."},
	{speaker: "Monsieur Bernard", text: "Marvelous. Five hundred years of cleaning, and ze paint still has secrets."},
}

func newCafePatronBernard(renderer *sdl.Renderer) *npc {
	idle, talk := loadCafePatronGrids(renderer, "bernard")
	return &npc{
		idleGrid: idle,
		talkGrid: talk,
		// User playtest #23: moved left + a little down.
		bounds:         sdl.Rect{X: 195, Y: 355, W: 90, H: 160},
		name:           "Monsieur Bernard",
		dialog:         bernardDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
}

// --- Mademoiselle Camille (red beret, art student) — QUEST NPC ---
//
// "Camille and the Sold-Out Postcard" (2026-06-10, reworked same day —
// user wanted a lighter tone and back-and-forth between street, bakery
// and museum, woven into the MAIN postcard chain):
//  1. Beaumont's restoration postcards are SOLD OUT (ties into Yvette /
//     Bernard's gossip). He asks for a replica sketch by Camille for the
//     archive wall — in trade for his own archive postcard.
//  2. Camille is thrilled, but she lost her lucky charcoal pencil
//     sketching the Louvre at sunrise. Nicolas the photographer saw
//     where it rolled (street hop).
//  3. PP fishes the pencil out of the flower pot by the Louvre steps,
//     returns it → Camille sketches the Room 7 portrait (one-shot from
//     npc_camille_sketching.png) → "Camille's Sketch".
//  4. Beaumont trades the sketch for the Postcard (main chain resumes).
var camilleFlavorDialog = []dialogEntry{
	{speaker: "Mademoiselle Camille", text: "Ze light in zis café is perfect for sketching, monsieur. I drew ze Louvre at sunrise today — magnifique!"},
	{speaker: "Mademoiselle Camille", text: "One day my sketches will hang INSIDE ze museum, not just outside of it."},
}

// camilleSketchAskDialog plays once Beaumont has asked for her sketch.
var camilleSketchAskDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Mademoiselle — Curator Beaumont needs a sketch of ze Room 7 portrait. He says you have ze fastest charcoal in Paris."},
	{speaker: "Mademoiselle Camille", text: "Beaumont said ZAT? About ME?! Monsieur, I would sketch ze whole Louvre for him!"},
	{speaker: "Mademoiselle Camille", text: "But — oh non. My lucky charcoal pencil. I lost it zis morning, sketching ze museum at sunrise."},
	{speaker: "Mademoiselle Camille", text: "Ask Nicolas, ze photographer by ze steps. Nothing happens on zat street without his lens seeing it."},
}

var camillePencilReminderDialog = []dialogEntry{
	{speaker: "Mademoiselle Camille", text: "No pencil, no masterpiece, monsieur. Nicolas sees everything — ask him where it rolled!"},
}

// camilleSketchTradeDialog fires when PP hands her the pencil (altDialog).
var camilleSketchTradeDialog = []dialogEntry{
	{speaker: "Mademoiselle Camille", text: "My lucky pencil! You are a hero of ze arts, monsieur."},
	{speaker: "Mademoiselle Camille", text: "Now watch — ze Room 7 portrait. I know every brushstroke by heart..."},
	{speaker: "Mademoiselle Camille", text: "(charcoal flies across ze sketchpad — quick, confident strokes)"},
	{speaker: "Mademoiselle Camille", text: "Voilà! Tell Beaumont zis one comes with interest — and zat Camille is ready for a commission."},
}

var camillePostSketchDialog = []dialogEntry{
	{speaker: "Mademoiselle Camille", text: "Imagine — MY sketch, on ze Louvre archive wall. Keep posing for ze world, monsieur!"},
}

func newCafePatronCamille(renderer *sdl.Renderer) *npc {
	idle, talk := loadCafePatronGrids(renderer, "camille")
	n := &npc{
		idleGrid: idle,
		talkGrid: talk,
		// User playtest #23: nudged right (420 → 500) so her legs tuck behind
		// the café table in the BG instead of showing below it.
		bounds:         sdl.Rect{X: 500, Y: 339, W: 90, H: 160},
		name:           "Mademoiselle Camille",
		dialog:         camilleFlavorDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	// Sketching one-shot (EXTRA_PROMPTS §T): ends with Camille turning the
	// sketchpad toward the camera, revealing the drawing of PP.
	if f := loadNPCGrid(renderer, "assets/images/locations/paris/npc/coffee/npc_camille_sketching.png", 8, 1); len(f) > 0 {
		n.oneShotAnims = map[string][]npcFrame{"sketch": f}
	}
	return n
}

// --- Monsieur Henri (silver mustache, croissant + bag) — QUEST NPC ---
//
// Henri's flow:
//  1. First visit: asks PP to fetch a café au lait. Promises something
//     from his bag in return.
//  2. PP brings the Café au Lait → altDialog fires: Henri remembers his
//     promise, hands PP homemade Confiture from his bag.
//  3. PP can now trade the Confiture to Pierre.
var henriInitialDialog = []dialogEntry{
	{speaker: "Monsieur Henri", text: "Ah, mon ami! A pink panther — zere's a sight."},
	{speaker: "Monsieur Henri", text: "Could you do an old man a favor? Madame Poulain is overworked. Fetch me a coffee?"},
	{speaker: "Monsieur Henri", text: "If you bring me one, I'll dig something nice from my bag for you. I keep treasures in here."},
}

var henriCoffeeTradeDialog = []dialogEntry{
	{speaker: "Monsieur Henri", text: "Ah, mon ami! Café au lait, parfait!"},
	{speaker: "Monsieur Henri", text: "Remember I said I had something from my bag for you?"},
	{speaker: "Monsieur Henri", text: "Here — homemade strawberry confiture, made it zis morning. Goes well on a fresh baguette."},
	{speaker: "Pink Panther", text: "Merci, Henri. Smells incredible."},
}

var henriPostTradeDialog = []dialogEntry{
	{speaker: "Monsieur Henri", text: "Enjoy ze confiture, mon ami! And if you see Pierre, tell him to eat properly."},
}

func newCafePatronHenri(renderer *sdl.Renderer) *npc {
	idle, talk := loadCafePatronGrids(renderer, "henri")
	n := &npc{
		idleGrid:       idle,
		talkGrid:       talk,
		bounds:         sdl.Rect{X: 580, Y: 339, W: 90, H: 160},
		name:           "Monsieur Henri",
		dialog:         henriInitialDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
	// #25: Henri hands PP the jam — 6-frame give-jam one-shot.
	if f := loadNPCGrid(renderer, "assets/images/locations/paris/npc/coffee/npc_henri_give_jam.png", 6, 1); len(f) > 0 {
		n.oneShotAnims = map[string][]npcFrame{"give_jam": f}
	}
	return n
}

// --- Lucien (turtleneck, espresso) — flavor with Tokyo foreshadow ---
var lucienDialog = []dialogEntry{
	{speaker: "Lucien", text: "Ze world... it does not feel right zis week, monsieur."},
	{speaker: "Lucien", text: "I had ze same dream three nights now — a tower covered in flowers, and bells ringing far away."},
	{speaker: "Lucien", text: "Probably ze coffee. Or someone is calling, somewhere."},
}

func newCafePatronLucien(renderer *sdl.Renderer) *npc {
	idle, talk := loadCafePatronGrids(renderer, "lucien")
	return &npc{
		idleGrid:       idle,
		talkGrid:       talk,
		bounds:         sdl.Rect{X: 920, Y: 339, W: 90, H: 160},
		name:           "Lucien",
		dialog:         lucienDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
}

// --- Madame Élise (auburn hair, autumn scarf) — flavor warmth ---
var eliseDialog = []dialogEntry{
	{speaker: "Madame Élise", text: "Autumn is coming early zis year, monsieur."},
	{speaker: "Madame Élise", text: "Ze geese are already heading south, and ze wind has teeth."},
	{speaker: "Madame Élise", text: "Wear your scarf. Even a panther catches cold."},
}

func newCafePatronElise(renderer *sdl.Renderer) *npc {
	idle, talk := loadCafePatronGrids(renderer, "elise")
	return &npc{
		idleGrid:       idle,
		talkGrid:       talk,
		bounds:         sdl.Rect{X: 660, Y: 540, W: 90, H: 160},
		name:           "Madame Élise",
		dialog:         eliseDialog,
		bobAmount:      0,
		talkFrameSpeed: 0.10,
	}
}

package game

import (
	"math"
	"os"
	"sort"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type spriteFrame struct {
	tex *sdl.Texture
	w   int32
	h   int32
	// ox/oy/ow/oh: opaque content box within the cell (frame-local). drawScaled
	// samples it and anchors by feet + centre so PP doesn't "jump"/drift and
	// idle/talk/walk stay the same on-screen size.
	ox int32
	oy int32
	ow int32
	oh int32
	// footCX: cell-local X of the feet (bottom-band centre of mass). drawScaled
	// plants this at the same screen X every frame so PP stands/walks/talks in
	// one place even when the art drifts the body or an arm/leg extends out.
	footCX int32
	// footRow: cell-local Y of the foot LINE, stabilized to the sheet median
	// (stabilizeFrames). drawScaled plants this row on the screen foot Y, so a
	// dipped tail or a frame whose art sits higher in the cell can no longer
	// lift/drop the whole body (user 2026-06-10: "the object stays in the
	// same spot"). Zero = unset → fall back to the frame's own content bottom.
	footRow int32
}

const (
	playerBaseSpeed = 250.0
	// User 2026-05-12: PP made the tallest character on screen by design
	// (see docs/CHARACTERS.md + retro_frames reference). Initial rebalance
	// targeted H=340 to match raw retro 54% screen-fill, but that
	// overshot - our background art is composed at a smaller character
	// scale than retro PtP, so PP at 340 dwarfed the cabins and the kids
	// rose above the cabin doorways. Pulled back to H=270, which keeps PP
	// the tallest on screen (still > all NPCs) without overwhelming the
	// background composition. Width 200 keeps the same cell aspect.
	playerDstW    = 200
	playerDstH    = 270
	playerMinX    = 10.0
	playerMaxX    = engine.ScreenWidth - playerDstW - 10.0
	playerMinY    = 265.0
	playerMaxY    = 395.0
	walkFrameTime = 0.08
	// 2026-07-15 (user): the front/back walk sheets read too fast at the side
	// cadence — tick them slower.
	walkFrameTimeFrontBack = 0.10
	talkFrameTime          = 0.07
	actionFrameTime        = 0.10
)

// walkFrameSecs is the per-frame cadence for the ACTIVE walk sheet: front/back
// walks cycle slower than the side walk (user 2026-07-15).
func (p *player) walkFrameSecs() float64 {
	if p.dir == dirUp || p.dir == dirDown {
		return walkFrameTimeFrontBack
	}
	return walkFrameTime
}

type playerState int

const (
	stateIdle playerState = iota
	stateWalking
	stateTalking
	stateGrabbing
	stateUsing
	stateExamining
	stateReacting
	stateShowInventory
)

type direction int

const (
	dirRight direction = iota
	dirLeft
	dirUp
	dirDown
)

type player struct {
	walkSideFrames  []spriteFrame
	walkUpFrames    []spriteFrame
	walkDownFrames  []spriteFrame
	idleFrontFrames []spriteFrame
	idleSideFrames  []spriteFrame
	idleBackFrames  []spriteFrame
	talkFrames      []spriteFrame
	talkSideFrames  []spriteFrame
	talkBackFrames  []spriteFrame
	// talkSideScaredFrames + scaredTalk (2026-07-24 #34): frightened side-talk
	// for the rude-Higgins beat (§PP-TALK-SIDE-SCARED, optional).
	talkSideScaredFrames []spriteFrame
	scaredTalk           bool
	grabFrames           []spriteFrame
	useItemFrames        []spriteFrame
	// seated mode (Japan tea ceremony): while `seated`, idle/talk render the
	// kneeling poses instead of the standing ones. Set after the spin-and-sit
	// entry one-shot, cleared when PP stands again.
	seatedIdle    []spriteFrame
	seatedTalk    []spriteFrame
	seated        bool
	examineFrames []spriteFrame
	reactFrames   []spriteFrame
	showInvFrames []spriteFrame

	x, y    float64
	targetX float64
	targetY float64
	// lWalkApproach (2026-07-15 user #13, paris_bakery): NPC talk approaches
	// walk the horizontal leg along the aisle FIRST, then the vertical leg at
	// the target x — instead of an oblique beeline across the tables.
	lWalkApproach bool
	// lastFootBlockers (2026-07-23 #9): the current scene's footBlockers,
	// stashed each update so startLWalk can test the horizontal corridor
	// for table crossings before committing to a straight leg.
	lastFootBlockers []sdl.Rect
	// oneShotFlip (2026-07-24 #23/#28): one-shot anims whose sheet faces the
	// wrong way for their beat — mirrored at draw time (NPC pattern).
	oneShotFlip map[string]bool
	// movementLocked (2026-07-25 user): freeze PP through an NPC's
	// multi-beat give chain (Poulain's counter hand-backs) — world clicks
	// and walks are ignored until the chain's final beat unlocks, so PP
	// can't wander off or turn out of the back-facing pose mid-hand-over.
	movementLocked bool
	moving         bool
	allowOffscreen bool
	facingLeft     bool
	dir            direction
	state          playerState

	breathTimer    float64
	walkCycleIdx   int
	walkTimer      float64
	talkCycleIdx   int
	talkTimer      float64
	actionIdx      int
	actionTimer    float64
	actionCallback func()

	interactTarget *npc
	dialogSys      *dialogSystem
	// inv is the shared inventory pointer used to gate item-specific
	// alt dialogs (e.g. Lily's flower handoff). Set once in Game.New
	// and read on every startNPCDialog call. Nil-safe.
	inv       *inventory
	onArrival func()

	sceneMinY float64
	sceneMaxY float64

	// Recede-into-distance tween used by the "walk to camp" transition.
	// While recedeActive: state stays stateWalking with dir=dirUp so the
	// back-walk frames cycle, position drifts up by (recedeDyUp * t), and
	// recedeScale lerps from 1.0 to recedeEndScale. drawScaled multiplies
	// recedeScale into the final draw rect. On completion, recedeOnDone
	// fires (typically a sceneMgr.transitionTo).
	recedeActive   bool
	recedeStartX   float64
	recedeStartY   float64
	recedeEndScale float64
	recedeDyUp     float64
	recedeDuration float64
	recedeElapsed  float64
	recedeScale    float64
	recedeOnDone   func()

	// Smooth recede release (#16). When a dialog that used playRecede ends,
	// snapping recedeScale back to 1.0 made PP visibly "pop" to full size.
	// releaseRecedeSmooth instead un-freezes movement and lerps recedeScale
	// from its current (shrunk) value back to 1.0 over recedeReleaseDur, so
	// PP grows back smoothly as the conversation ends.
	recedeReleasing      bool
	recedeReleaseFrom    float64
	recedeReleaseElapsed float64
	recedeReleaseDur     float64

	// recedeHeld (#28): PP keeps rendering at the receded (shrunk) scale after a
	// recede-framed dialog ends, with movement re-enabled. He grows back to full
	// size (releaseRecedeSmooth) on the next setTarget - i.e. the next click that
	// moves him elsewhere. Used for Pierre: PP stays small at Pierre's depth
	// until the player walks him away.
	recedeHeld bool

	// Walk-in tween (opening-cutscene PP entering from off-screen-left).
	// Drives p.x directly each frame via lerp, bypassing the moving /
	// allowOffscreen / clamp machinery so the engine can't shove PP back
	// on-screen mid-walk. Side-walk frames cycle exactly like the
	// recede tween. On completion, walkInOnDone fires (typically starts
	// the opening monologue).
	walkInActive   bool
	walkInStartX   float64
	walkInEndX     float64
	walkInY        float64
	walkInDuration float64
	walkInElapsed  float64
	walkInOnDone   func()

	// One-shot named animations triggered by the sequence player. While
	// activeOneShot != "" the player draws frames from oneShotAnims[name]
	// instead of the state-based idle/walk/talk cycle. Used for the
	// give-map sequence (PP plays "receive_map") and similar short clips.
	oneShotAnims    map[string][]spriteFrame
	activeOneShot   string
	oneShotIdx      int
	oneShotTimer    float64
	oneShotDuration float64
	oneShotElapsed  float64 // wall-clock elapsed since playOneShot start; drives completion
	oneShotOnDone   func()
}

// spriteInset is the pixel margin trimmed off each cell at load time to strip
// the black grid-line borders that AI-generated sheets bake in between cells.
const spriteInset = 3

func gridFrames(renderer *sdl.Renderer, path string, cols, rows int) []spriteFrame {
	grid := engine.SpriteGridFromPNGClean(renderer, path, cols, rows, spriteInset)
	var frames []spriteFrame
	for r := 0; r < rows && r < len(grid); r++ {
		for c := 0; c < cols && c < len(grid[r]); c++ {
			gf := grid[r][c]
			if gf.Tex == nil {
				// Missing sheet → engine.emptyGrid; skip so `len(frames) > 0`
				// guards don't register invisible animations.
				continue
			}
			frames = append(frames, spriteFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH, footCX: gf.FCX, footRow: gf.FRY})
		}
	}
	stabilizeFootCX(frames)
	return frames
}

// gridFramesTol (2026-07-23 #5/#13/#15): GLOBAL key with a caller tolerance.
// The tol-8 global key missed the off-white/off-blue ENCLOSED pockets (232-246
// whites) generators paint between an arm and the body, and the connected key
// preserves them by design - both left visible squares in the arm gaps. Tol 24
// clears the pockets while the art-rule colors survive (ivory belly #F2EFE5 is
// 26 off pure white, cream props 49+ off the flat blue). ONLY for sheets whose
// props are cream/tan - a white prop (the Paris porcelain cup) would be eaten.
func gridFramesTol(renderer *sdl.Renderer, path string, cols, rows int, tol uint8) []spriteFrame {
	grid := engine.SpriteGridFromPNGCleanTol(renderer, path, cols, rows, spriteInset, tol)
	var frames []spriteFrame
	for r := 0; r < rows && r < len(grid); r++ {
		for c := 0; c < cols && c < len(grid[r]); c++ {
			gf := grid[r][c]
			if gf.Tex == nil {
				continue
			}
			frames = append(frames, spriteFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH, footCX: gf.FCX, footRow: gf.FRY})
		}
	}
	stabilizeFootCX(frames)
	return frames
}

// gridFramesConnected is gridFrames with the EDGE-CONNECTED color key (the
// same loader the NPCs use): only background white reachable from the sheet
// edges is cleared, so enclosed whites survive. Used for the item hand-off
// one-shots (give_*/receive/grab flower) - their hand-held props are WHITE
// (daisy petals, coffee cup, postcard, sketch page, jam-jar label) and the
// global key punched see-through holes in the very item being handed over
// (user 2026-06-12 PR#2: "give flower not smooth" - the flower flickered).
func gridFramesConnected(renderer *sdl.Renderer, path string, cols, rows int) []spriteFrame {
	// Default connected-key tolerance (8) leaves an anti-aliased white fringe on
	// softly-drawn sheets. gridFramesConnectedTol lets a caller widen it.
	return gridFramesConnectedTol(renderer, path, cols, rows, 8)
}

func gridFramesConnectedTol(renderer *sdl.Renderer, path string, cols, rows int, tol uint8) []spriteFrame {
	grid := engine.SpriteGridFromPNGCleanConnectedTol(renderer, path, cols, rows, spriteInset, tol)
	var frames []spriteFrame
	for r := 0; r < rows && r < len(grid); r++ {
		for c := 0; c < cols && c < len(grid[r]); c++ {
			gf := grid[r][c]
			if gf.Tex == nil {
				continue
			}
			frames = append(frames, spriteFrame{tex: gf.Tex, w: gf.W, h: gf.H,
				ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH, footCX: gf.FCX, footRow: gf.FRY})
		}
	}
	stabilizeFootCX(frames)
	return frames
}

// stabilizeFootCX applies a DEADBAND to every frame's foot anchors (foot
// centre X + foot row Y): values within ±6px of the sheet median snap to the
// median - killing detection noise so well-aligned frames are rock-stable -
// while larger deviations keep the frame's OWN feet position, so art that
// genuinely drifts inside the cell (e.g. row 0 drawn higher than row 1) is
// compensated per frame instead of rendering as a jump (user 2026-06-10:
// a constant median row made PP "jump between two rows" on the idle sheet).
func stabilizeFootCX(frames []spriteFrame) {
	const deadband = 6
	vals := make([]int, 0, len(frames))
	rowVals := make([]int, 0, len(frames))
	for _, f := range frames {
		if f.ow > 0 {
			vals = append(vals, int(f.footCX))
			rowVals = append(rowVals, int(f.footRow))
		}
	}
	if len(vals) == 0 {
		return
	}
	sort.Ints(vals)
	sort.Ints(rowVals)
	med := int32(vals[len(vals)/2])
	medRow := int32(rowVals[len(rowVals)/2])
	for i := range frames {
		if frames[i].ow <= 0 {
			continue
		}
		if d := frames[i].footCX - med; d >= -deadband && d <= deadband {
			frames[i].footCX = med
		}
		if d := frames[i].footRow - medRow; d >= -deadband && d <= deadband {
			frames[i].footRow = medRow
		}
	}
}

func newPlayer(renderer *sdl.Renderer) *player {
	p := &player{
		x: 630,
		y: float64(engine.ScreenHeight) - playerDstH - 100,
	}

	// User playtest 2026-06-05: walk-left is now a SINGLE ROW of 8 frames (the
	// whole walk cycle). The old 8×2 read only row 0, so a cycle split across
	// two rows showed as half a stride ("not a full circle"). Load 8×1.
	// §AB (2026-06-10): the side-walk strip exists in two generations - the
	// current 10×1 (2400px wide) and the regen spec 8×1 (1536px wide). Pick
	// the column count from the sheet on disk so dropping the new art in
	// just works ("read the sprite properly" - no code change needed).
	walkSideCols := 10
	if w, _, err := engine.PNGSize("assets/images/player/PP walk left.png"); err == nil && w <= 1600 {
		walkSideCols = 8
	}
	p.walkSideFrames = gridFrames(renderer, "assets/images/player/PP walk left.png", walkSideCols, 1)
	// 2026-07-15: front/back walks are now 12-frame 6x2 cycles, played row-major.
	p.walkDownFrames = gridFrames(renderer, "assets/images/player/PP walk front.png", 6, 2)
	p.walkUpFrames = gridFrames(renderer, "assets/images/player/PP walk back.png", 6, 2)

	// Idle images - use all frames for animated idle
	p.idleFrontFrames = gridFrames(renderer, "assets/images/player/PP idle front.png", 8, 2)
	p.idleSideFrames = gridFrames(renderer, "assets/images/player/PP idle side.png", 8, 2)
	p.idleBackFrames = gridFrames(renderer, "assets/images/player/PP idle back.png", 8, 2)

	p.talkFrames = gridFrames(renderer, "assets/images/player/PP talk front.png", 8, 2)
	p.talkSideFrames = gridFrames(renderer, "assets/images/player/PP talk side.png", 8, 2)
	// Back-facing talk sheet (PP seen from behind, talking up toward a
	// counter-height NPC like Madame Poulain). Optional - guarded with
	// firstExisting so a missing file isn't re-opened every startup; until the
	// art lands, currentTalkFrames falls back to the back IDLE so PP still
	// shows his back (not his face) during a ppFaceBack dialog. See
	// EXTRA_PROMPTS §PP-TALK-BACK.
	if path := firstExisting("assets/images/player/PP talk back.png", "assets/images/player/PP_talk_back.png"); path != "" {
		p.talkBackFrames = gridFrames(renderer, path, 8, 2)
	}
	// 2026-07-24 (user #34): frightened side-talk for the rude-Higgins beat
	// (§PP-TALK-SIDE-SCARED). Optional — setScaredTalk no-ops until it lands.
	if path := firstExisting("assets/images/player/PP talk side scared.png", "assets/images/player/PP_talk_side_scared.png"); path != "" {
		p.talkSideScaredFrames = gridFrames(renderer, path, 8, 2)
	}

	// Seated poses for the Japan tea ceremony (kneeling idle + talk). Optional -
	// guarded so a missing sheet isn't re-opened every startup; while seated with
	// no art, currentIdle/TalkFrames fall back to the standing poses (§JP-TEA-SIT).
	for _, ns := range []struct {
		dst   *[]spriteFrame
		cands []string
	}{
		{&p.seatedIdle, []string{"assets/images/player/PP_sit_idle.png", "assets/images/player/PP_seated_idle.png"}},
		{&p.seatedTalk, []string{"assets/images/player/PP_sit_talk.png", "assets/images/player/PP_seated_talk.png"}},
	} {
		if path := firstExisting(ns.cands...); path != "" {
			*ns.dst = gridFrames(renderer, path, 8, 1)
		}
	}
	// User 2026-05-12: swapped from "PP grab flower.png" (square 128×128
	// cells) to the canonical "PP grab.png" (portrait 172×384 cells). The
	// square cells made the grab anim render shorter than idle inside the
	// 245×340 bounds - full size now matches idle. PP grab.png shows the
	// same crouch-and-rise cycle with a flower in the last frames.
	p.grabFrames = gridFrames(renderer, "assets/images/player/PP grab.png", 8, 2)

	// 2026-07-18 (user #1): PP celebrate.png + PP sneak examine.png were
	// deleted from disk — react/showInv/examine fall back to the grab frames
	// (their states barely appear in play) instead of warning every startup.
	p.reactFrames = p.grabFrames
	if len(p.grabFrames) >= 2 {
		p.showInvFrames = p.grabFrames[0:2]
	}
	p.examineFrames = p.grabFrames
	// 2026-07-18 (user): PP sneak use.png removed — same fallback as examine.
	p.useItemFrames = p.grabFrames

	// One-shot named animations for sequence playback. receive_map drives
	// the give-map handoff so PP visibly takes the map from Higgins instead
	// of having it appear in the inventory bar. grab_flower plays when the
	// flower floor item at the lake is picked up (and re-used when handing
	// it over to Lily).
	p.oneShotAnims = map[string][]spriteFrame{}
	// 2026-07-27 (user: "speak and drink tea when needed"): the landed
	// PP_sit_talk.png is really TWO beats on one sheet — frames 1-2 kneeling
	// talk, frames 3-8 a full lift→sip→lower drink arc. Slice it: the talk
	// loop keeps only the mouth beats (drinking mid-sentence read as random),
	// and the drink arc becomes the `sit_drink` one-shot the tea ceremony
	// plays at the shared-sip moment.
	if len(p.seatedTalk) == 8 {
		p.oneShotAnims["sit_drink"] = p.seatedTalk[2:8]
		p.seatedTalk = p.seatedTalk[0:2]
	}
	// The spin→sit TRANSFORM: PP standing → fast spin into the tea clothes →
	// drops into a kneel (ends seated). No-op until the art lands.
	if path := firstExisting(
		"assets/images/player/PP_spin_to_sit.png",
		"assets/images/player/PP_tea_ceremony.png",
		"assets/images/player/PP_sit_down.png",
	); path != "" {
		if f := gridFramesConnected(renderer, path, 8, 1); len(f) > 0 {
			p.oneShotAnims["tea_sit"] = f
		}
	}
	// 2026-07-28: the reverse ceremony beat is a separate 6-frame strip so PP
	// visibly rises from seiza before the seated mode is cleared.
	if path := firstExisting("assets/images/player/PP_sit_to_stand.png"); path != "" {
		if f := gridFrames(renderer, path, 6, 1); len(f) > 0 {
			p.oneShotAnims["tea_stand"] = f
		}
	}
	// Sakura payoff: camera out of the invisible pocket, selfie, camera back in.
	if path := firstExisting("assets/images/player/pp_sakura_selfie.png"); path != "" {
		if f := gridFrames(renderer, path, 8, 1); len(f) > 0 {
			p.oneShotAnims["sakura_selfie"] = f
		}
	}
	// User 2026-05-31 (#15): receive-map sheet is 4×2, not 8×2 - cutting 8×2
	// split every pose down the middle so PP "blinked" the whole catch.
	receiveMap := gridFrames(renderer, "assets/images/player/PP receive map.png", 4, 2)
	if len(receiveMap) > 0 {
		p.oneShotAnims["receive_map"] = receiveMap
	}
	// Dedicated six-frame flower pickup strip: neutral, crouch, reach, grab,
	// stand holding low, stand holding chest-high. The sheet now uses off-white
	// petals/body details, so the normal global key can remove both the outer
	// background and the white gap between PP's hand and body.
	grabFlower := gridFrames(renderer, "assets/images/player/PP grab flower.png", 6, 1)
	if len(grabFlower) > 0 {
		p.oneShotAnims["grab_flower"] = grabFlower
	}
	// User playtest #4/#14/#19: prefer the dedicated waist-high basket pickup;
	// retain the older ground-reach strip as a graceful fallback.
	if path := firstExisting(
		"assets/images/player/PP grab basket.png",
		"assets/images/player/PP grab rolling pin.png",
	); path != "" {
		if f := gridFrames(renderer, path, 6, 1); len(f) > 0 {
			p.oneShotAnims["grab_rolling_pin"] = f
		}
	}
	// #25: PP receiving the baguette / the jam from the bakery NPCs. The
	// baguette receive was re-rolled 6×1 on a sampled blue background; use the
	// global key so enclosed blue gaps disappear without erasing the cream prop.
	if f := gridFrames(renderer, "assets/images/player/PP get bagguette.png", 6, 1); len(f) > 0 {
		p.oneShotAnims["get_baguette"] = f
	}
	// 2026-07-24 (user #10): tol 8 → 24 — white AA fringe/spots (flower fix).
	if f := gridFramesConnectedTol(renderer, "assets/images/player/PP get jam.png", 8, 1, 24); len(f) > 0 {
		p.oneShotAnims["get_jam"] = f
	}
	// #7/#9/#15/#16 (user 2026-06-30): do NOT register a generic "receive_item"
	// fallback. It used the crouch-and-lift grab frames, which (a) read as the
	// wrong "get item" pose and (b) normalised TINY. With no generic fallback a
	// receive with no item-specific sheet simply no-ops (PP stays idle, the item
	// still lands in the bag). Missing per-item receive sheets are queued as
	// prompts in docs/EXTRA_PROMPTS.md rather than faked with the grab.
	// §CAM-RECEIVE: PP takes the Camille sketch and slides it into his pocket.
	// The dropped art is named pp_get_skatch.png; accept that and the canonical
	// name. Falls back to grabFrames (generic reach) if neither is present.
	if rp := firstExisting("assets/images/player/PP receive sketch.png", "assets/images/player/pp_get_skatch.png"); rp != "" {
		// 2026-07-25 (user): tol 8 → 24 — white AA fringe around PP ("white
		// spots"); the connected flood still protects the bone sketch page.
		if f := gridFramesConnectedTol(renderer, rp, 8, 1, 24); len(f) > 0 {
			p.oneShotAnims["receive_sketch"] = f
		}
	}
	// 2026-07-15 (user #9/#14): pp_get_card.png is the PRESS-PASS ("Card")
	// receive — register it under "receive_card" so it plays when Pierre hands
	// the pass over. It must NOT double as the Beaumont postcard receive.
	if rp := firstExisting("assets/images/player/pp_get_card.png"); rp != "" {
		// 2026-07-25 (user): GLOBAL tol-24 — the blue-bg sheet kept enclosed
		// blue arm-gap pockets under the connected key ("blue spots around
		// him"); the card is cream, nothing blue on PP.
		if f := gridFramesTol(renderer, rp, 8, 1, 24); len(f) > 0 {
			p.oneShotAnims["receive_card"] = f
		}
	}
	// §PP-GET-POSTCARD: dedicated Louvre-postcard receive beat.
	if rp := firstExisting("assets/images/player/PP get postcard.png"); rp != "" {
		if f := gridFrames(renderer, rp, 6, 1); len(f) > 0 {
			p.oneShotAnims["receive_postcard"] = f
		}
	}
	// No grab fallback for receive_sketch (#9): if the dedicated sheet is absent
	// PP just takes it in idle rather than the tiny crouch-grab.
	// Generic "grab" one-shot (the crouch-and-lift cycle). Floor-item pickups
	// must run through a ONE-SHOT, never playAction(stateGrabbing): player.update
	// force-resets any non-talking state to idle every frame while !moving, so a
	// playAction grab is killed before it completes and its add-item callback
	// never fires (2026-06-15 #19/#20 - the pot pencil softlock). One-shots are
	// tracked separately (activeOneShot) and always fire onDone.
	if len(p.grabFrames) > 0 {
		p.oneShotAnims["grab"] = p.grabFrames
	}
	// §PG1 (2026-06-11 #33, refined same day): PER-ITEM give one-shots - the
	// user wants each hand-over to SHOW the actual item leaving PP's paw.
	// Each sheet is optional; playGive falls back item -> generic -> grab.
	// Connected key (PR#2): the handed items are largely WHITE (daisy, cup,
	// postcard, sketch page, jar label) - the global key erased them mid-give.
	for key, path := range map[string]string{
		"flower":        "assets/images/player/PP give flower.png",
		"confiture":     "assets/images/player/PP give confiture.png",
		"baguette_heel": "assets/images/player/PP give heel.png",
		// "pencil" moved to the global tol-24 load below (2026-07-24 #16:
		// blue-bg sheet; the connected key kept its enclosed blue pockets).
		// "rolling_pin" RETIRED (2026-07-24, user): the pin's only hand-over
		// is Poulain's back-facing trade, which always resolves to
		// give_rolling_pin_back — the front sheet never played.
	} {
		// #5: tol 24 (was 8) — the give poses touch with a faint near-white
		// bridge that survived the low-tol key, so the gap-detector mis-cut them
		// (the office-Higgins bug). tol 24 breaks the bridge for a clean 1×8 cut;
		// PP's pink/black and the enclosed held item are far from white and
		// survive.
		if f := gridFramesConnectedTol(renderer, path, 8, 1, 24); len(f) > 0 {
			p.oneShotAnims["give_"+key] = f
		}
	}
	// 2026-07-24 (user #15): dedicated pot-pickup for the charcoal pencil
	// (§PP-GET-PENCIL-POT, the rolling-pin-basket pattern). Optional — the
	// pot pickup resolves to the generic grab until it lands.
	if f := gridFrames(renderer, "assets/images/player/pp_get_pencil_pot.png", 6, 1); len(f) > 0 {
		p.oneShotAnims["grab_pencil_pot"] = f
	}
	// 2026-07-24 (user #16): give-pencil is a BLUE-bg sheet — the connected
	// key kept its enclosed blue arm-gap pockets. Global tol-24 clears them
	// (nothing on PP is near the pale blue; the dark pencil is far off it).
	if f := gridFramesTol(renderer, "assets/images/player/PP give pencil.png", 8, 1, 24); len(f) > 0 {
		p.oneShotAnims["give_pencil"] = f
	}
	// 2026-07-17: wide-cell give re-rolls use 6×1 and a sampled blue
	// background. Global keying removes the blue even from enclosed arm gaps.
	for key, path := range map[string]string{
		"baguette":     "assets/images/player/PP give baguette.png",
		"sketch":       "assets/images/player/PP give sketch.png",
		"postcard":     "assets/images/player/PP give postcard.png",
		"cafe_au_lait": "assets/images/player/PP give coffee.png",
	} {
		if f := gridFrames(renderer, path, 6, 1); len(f) > 0 {
			p.oneShotAnims["give_"+key] = f
		}
	}
	// 2026-07-15 (user #25/#27): Jerusalem per-item sheets — all optional,
	// no-op until the art lands (§PP-GIVE-BAGEL / §PP-GET-PEN etc.).
	for key, path := range map[string]string{
		"receive_coin": "assets/images/player/pp_get_coin.png",
		"give_coin":    "assets/images/player/PP give coin.png",
	} {
		if _, err := os.Stat(path); err == nil {
			if f := gridFramesConnectedTol(renderer, path, 8, 1, 24); len(f) > 0 {
				p.oneShotAnims[key] = f
			}
		}
	}
	// 2026-07-17: blue-background Jerusalem sheets are all wide-cell 6×1.
	for key, path := range map[string]string{
		"give_bagel":               "assets/images/player/PP give bagel.png",
		"receive_bagel":            "assets/images/player/pp_get_bagel.png",
		"give_card":                "assets/images/player/PP give card.png",
		"give_fire_striker":        "assets/images/player/PP give fire striker.png",
		"receive_well_water":       "assets/images/player/pp_get_well_water.png",
		"receive_voice_charm":      "assets/images/player/pp_get_voice_charm.png",
		"receive_offering_bowl":    "assets/images/player/pp_get_offering_bowl.png",
		"give_matcha_bowl":         "assets/images/player/PP give matcha bowl.png",
		"receive_pen":              "assets/images/player/pp_get_pen.png",
		"receive_cardamom":         "assets/images/player/pp_get_cardamom.png",
		"give_cardamom":            "assets/images/player/PP give cardamom.png",
		"receive_jerusalem_coffee": "assets/images/player/pp_get_coffee.png",
		"give_jerusalem_coffee":    "assets/images/player/PP give jerusalem coffee.png",
		"receive_paper":            "assets/images/player/pp_get_paper.png",
	} {
		if _, err := os.Stat(path); err == nil {
			if f := gridFrames(renderer, path, 6, 1); len(f) > 0 {
				p.oneShotAnims[key] = f
			}
		}
	}
	// 2026-07-24 (user #23/#28): these receive sheets face away from their
	// giver in-scene — mirror them at draw time.
	p.oneShotFlip = map[string]bool{
		"receive_jerusalem_coffee": true,
		"receive_paper":            true,
		// 2026-07-25 (user): the postcard take reaches away from Beaumont.
		"receive_postcard": true,
		// 2026-07-25 (user): the bagel take reaches away from the seller.
		"receive_bagel": true,
	}
	// 2026-07-25 (user): the pot pickup — until §PP-GET-PENCIL-POT lands,
	// alias the generic grab MIRRORED (the grab leans LEFT; the pot sits to
	// PP's RIGHT at the shared Pierre-depth mark). The dedicated sheet is
	// drawn reaching right, so it loads un-flipped and replaces the alias.
	if dp := firstExisting("assets/images/player/pp_get_pencil_pot.png"); dp != "" {
		if f := gridFrames(renderer, dp, 6, 1); len(f) > 0 {
			p.oneShotAnims["grab_pencil_pot"] = f
		}
	} else if len(p.grabFrames) > 0 {
		p.oneShotAnims["grab_pencil_pot"] = p.grabFrames
		p.oneShotFlip["grab_pencil_pot"] = true
	}
	// §PM1 (2026-06-12): PP pulls the travel map out of his invisible hip
	// pocket BEFORE the globe screen opens (user: "a new sprite of us pull
	// it out from pocket and then the map screen"). Optional - until the
	// sheet lands, openTravelMap's playOneShot no-ops straight to the map.
	// 2026-07-23 (#5): GLOBAL tol-24 — the blue re-roll left off-blue pockets
	// in the arm gaps under the connected key; the map prop is cream, safe.
	if f := gridFramesTol(renderer, "assets/images/player/PP pull map.png", 8, 1, 24); len(f) > 0 {
		p.oneShotAnims["pull_map"] = f
	}
	// PR#9 (2026-06-12): PP recoils/hops backward when the Paris biker bumps
	// him. Optional - the bump falls back to the generic flinch until this
	// sheet lands (§JUMPBACK).
	if f := gridFramesConnected(renderer, "assets/images/player/PP jump back.png", 8, 1); len(f) > 0 {
		p.oneShotAnims["jump_back"] = f
	}
	// #7/#16 (user 2026-06-30): no generic "give_item" fallback either — it used
	// the tiny grab frames and read as "get item". A give with no per-item sheet
	// no-ops (PP idle); the item still changes hands. Add the per-item give sheet
	// (above) or queue a prompt instead of faking it with the grab.
	// Jerusalem (#26): PP writes a note and tucks it in the Wall. Art is
	// pending (§JW1/§JW2) - guard the load with os.Stat so a missing file isn't
	// re-opened every startup (slow under on-access AV), and fall back to the
	// grab frames so the beat still animates.
	for _, ns := range []struct{ key, path string }{
		{"write_note", "assets/images/player/PP write note.png"},
	} {
		if _, err := os.Stat(ns.path); err == nil {
			// #5: tol 24 to break touching-pose bridges (put-note fell back to
			// no clean 6×1 split at the low tol).
			if f := gridFramesConnectedTol(renderer, ns.path, 6, 1, 24); len(f) > 0 {
				p.oneShotAnims[ns.key] = f
				continue
			}
		}
		if len(p.grabFrames) > 0 {
			p.oneShotAnims[ns.key] = p.grabFrames
		}
	}
	// The put-note re-roll uses the 6×1 blue-background convention.
	if f := gridFrames(renderer, "assets/images/player/PP put note in wall.png", 6, 1); len(f) > 0 {
		p.oneShotAnims["put_note"] = f
	} else if len(p.grabFrames) > 0 {
		p.oneShotAnims["put_note"] = p.grabFrames
	}
	// 2026-06-24: PP's back-facing hand-off sheets at Poulain's counter (#11/#7).
	// These are ITEM-SPECIFIC (user: PP takes the actual baguette/coffee, no
	// generic back sprite): the rolling pin he hands over and the baguette +
	// coffee he takes back, all drawn from behind. Guarded with os.Stat; until
	// each lands the trade falls back to the front sheet for that item.
	// 2026-07-23 (#13): the Poulain-counter back sheets carried off-white
	// enclosed pockets in the arm gaps (the connected key preserves them by
	// design). GLOBAL tol-24 key clears them; their props are wood/tan/cream,
	// no white detail to protect on a back view.
	for _, ns := range []struct{ key, path string }{
		{"give_rolling_pin_back", "assets/images/player/PP_give_rolling_pin_back.png"},
		{"get_baguette_back", "assets/images/player/PP_get_baguette_back.png"},
		{"receive_cafe_au_lait_back", "assets/images/player/PP_get_coffee_back.png"},
	} {
		if _, err := os.Stat(ns.path); err == nil {
			if f := gridFramesTol(renderer, ns.path, 8, 1, 24); len(f) > 0 {
				p.oneShotAnims[ns.key] = f
			}
		}
	}
	for _, ns := range []struct{ key, path string }{
		// Japan: the dresser-geisha gag - PP spins, ends up in a kimono, spins
		// again, back to normal. One continuous one-shot (no persistent state).
		{"kimono_spin", "assets/images/player/PP_kimono_spin.png"},
		// §PP-RUN-JAPAN: four sheets for the kimono gag.
		// 1) normal sprint off-right, 2) kimono sprint in from right (facing left),
		// 3) stop centre-stage and show off the kimono (pose beat), 4) kimono
		// sprint off-right again. Falls back gracefully until art lands.
		{"run_right", "assets/images/player/PP_run_right.png"},
		{"run_kimono_left", "assets/images/player/PP_run_kimono_left.png"},
		{"kimono_pose", "assets/images/player/PP_kimono_pose.png"},
		{"run_kimono_right", "assets/images/player/PP_run_kimono_right.png"},
	} {
		if _, err := os.Stat(ns.path); err == nil {
			if f := gridFramesConnected(renderer, ns.path, 8, 1); len(f) > 0 {
				p.oneShotAnims[ns.key] = f
			}
		}
	}
	// 2026-07-24 (user): the kimono gag is THREE separate beats — spin INTO
	// the costume, MODEL it, spin BACK to normal. Dedicated 6×1 blue sheets
	// are queued (§PP-KIMONO-SPLIT); until each lands, its phase is sliced
	// out of the combined PP_kimono_spin.png (frames: 1-2 spin blur →
	// 3 kimono forming → 4-6 dressed poses → 7 spin blur → 8 normal).
	for _, ks := range []struct {
		key, path string
		lo, hi    int
	}{
		{"kimono_spin_in", "assets/images/player/PP_kimono_spin_in.png", 0, 4},
		{"kimono_model", "assets/images/player/PP_kimono_model.png", 3, 6},
		{"kimono_spin_out", "assets/images/player/PP_kimono_spin_out.png", 6, 8},
	} {
		if _, err := os.Stat(ks.path); err == nil {
			if f := gridFrames(renderer, ks.path, 6, 1); len(f) > 0 {
				p.oneShotAnims[ks.key] = f
				continue
			}
		}
		if full := p.oneShotAnims["kimono_spin"]; len(full) >= 8 {
			p.oneShotAnims[ks.key] = full[ks.lo:ks.hi]
		}
	}

	p.dir = dirDown

	return p
}

// giveAnimKeyForItem maps inventory item NAMES to the per-item give one-shot
// keys registered in newPlayer ("give_<key>"). Unknown names return "item",
// which playGive resolves through the generic give/grab fallback chain.
func giveAnimKeyForItem(name string) string {
	switch name {
	case "Flower":
		return "flower"
	case "Rolling Pin":
		return "rolling_pin"
	case "Baguette":
		return "baguette"
	case "Confiture":
		return "confiture"
	case "Cafe au Lait":
		return "cafe_au_lait"
	case "Baguette Heel":
		return "baguette_heel"
	case "Charcoal Pencil":
		return "pencil"
	case "Camille's Sketch":
		return "sketch"
	case "Postcard", "Signed Postcard":
		return "postcard"
	case "Card":
		return "card"
	// 2026-07-15 (user #23/#25/#27): Jerusalem chain items were falling to the
	// generic "item" key, so their hand-overs played NO PP animation at all.
	case "Coffee":
		// 2026-07-18 (user #30): the Jerusalem finjan is NOT the Paris paper
		// cup — its own key; falls back to the generic beat until
		// §PP-GIVE-JER-COFFEE lands.
		return "jerusalem_coffee"
	case "Bagel":
		return "bagel"
	case "Pen":
		return "pen"
	case "Coin":
		return "coin"
	case "Cardamom":
		return "cardamom"
	case "Note Paper":
		return "paper"
	}
	return "item"
}

func receiveAnimKeyForItem(name string) string {
	switch name {
	case "Card":
		return "card"
	}
	return giveAnimKeyForItem(name)
}

// playHandOff runs the physical hand-over beat for an alt dialog (user
// 2026-06-12 PR#1): PP's give one-shot plays, then the NPC's receive
// one-shot, and only then `then` fires (the dispatcher starts the dialog
// text there). Every step degrades to an immediate `then` when art is
// missing, so the flow can't deadlock.
func (p *player) playHandOff(target *npc, ho *handOff, then func()) {
	if ho == nil || target == nil {
		if then != nil {
			then()
		}
		return
	}
	// Stage 2: after PP hands his item over (or immediately, for an NPC-only
	// give) the NPC reaches into frame and hands `returnItem` back, then PP
	// plays his receive beat. Folded into a closure so the one-way path (no
	// returnItem) just calls `then` directly.
	stageReturn := func() {
		if ho.returnItem == "" {
			then()
			return
		}
		rkey := receiveAnimKeyForItem(ho.returnItem)
		// The NPC gives the return item: prefer an explicit give one-shot, then
		// give_<key>, then the NPC's generic "give" reach.
		nanim := ho.npcGiveAnim
		if nanim == "" {
			nanim = "give_" + rkey
		}
		// Fall back to the NPC's generic "give" reach until the per-item give
		// sheet (give_paper, give_pen, give_coin...) lands, so the hand-back
		// still animates instead of snapping.
		if !target.hasOneShotAnim(nanim) {
			nanim = "give"
		}
		recvDur := ho.ppReceiveDur
		if recvDur <= 0 {
			recvDur = 1.3
		}
		target.playOneShotAnimThen(nanim, recvDur, func() {
			p.playReceive(rkey, ho.back, recvDur, then)
		})
	}

	giveDur := ho.giveDur
	if giveDur <= 0 {
		giveDur = 1.3
	}
	// Stage 1: PP hands his item over and the NPC takes it. When ho.item is
	// empty this is a pure NPC->PP give, so skip straight to stage 2.
	if ho.item == "" {
		stageReturn()
		return
	}
	key := giveAnimKeyForItem(ho.item)
	p.playGiveFacing(key, ho.back, giveDur, func() {
		// #10: some trades don't want the NPC's take-reach (it falls back to the
		// NPC's "give" sheet and double-plays the hand-back). Skip straight to the
		// return/give stage.
		if ho.skipNPCTake {
			stageReturn()
			return
		}
		anim := ho.npcAnim
		if anim == "" {
			anim = "receive_" + key
			if !target.hasOneShotAnim(anim) {
				anim = "receive_item"
			}
			// 2026-06-15 #11/#17: some NPCs register only a "give" reaching
			// one-shot (no receive_*/receive_item), so the NPC half of the
			// hand-over was an instant no-op. Fall back to the NPC's own "give"
			// reach so they visibly take the item.
			if !target.hasOneShotAnim(anim) && target.hasOneShotAnim("give") {
				anim = "give"
			}
			// 2026-07-15 (user #25, bagel seller): when there IS a return item
			// whose give-back will ALSO fall back to the same generic "give"
			// sheet, playing it here too reads as the NPC handing the item
			// TWICE. Skip the take-reach and go straight to the hand-back.
			if anim == "give" && ho.returnItem != "" && ho.npcGiveAnim == "" &&
				!target.hasOneShotAnim("give_"+receiveAnimKeyForItem(ho.returnItem)) {
				stageReturn()
				return
			}
		}
		dur := ho.npcAnimDur
		if dur <= 0 {
			dur = 1.3
		}
		target.playOneShotAnimThen(anim, dur, stageReturn)
	})
}

// playReceive plays PP taking `itemKey` into his paw: tries receive_<key>
// (and its _back variant when PP faces away), then the generic receive_item.
func (p *player) playReceive(itemKey string, back bool, dur float64, onDone func()) {
	name := p.resolveOneShot("receive_"+itemKey, back, "receive_item")
	if dur <= 0 {
		dur = 1.3
	}
	p.playOneShot(name, dur, onDone)
}

// playGiveFacing is playGive with a back-facing variant preference (counter /
// Wall trades where PP is shown from behind).
func (p *player) playGiveFacing(itemKey string, back bool, dur float64, onDone func()) {
	if !back {
		p.playGive(itemKey, dur, onDone)
		return
	}
	name := p.resolveOneShot("give_"+itemKey, true, "give_item")
	p.playOneShot(name, dur, onDone)
}

// resolveOneShot picks the best registered one-shot for a base key: the
// "_back" variant when `back` and it exists, then the bare key, then a
// fallback key. Used so back-facing trades degrade to the front sheets until
// the dedicated back art lands.
func (p *player) resolveOneShot(base string, back bool, fallback string) string {
	if back && p.hasOneShot(base+"_back") {
		return base + "_back"
	}
	if p.hasOneShot(base) {
		return base
	}
	return fallback
}

// playGive plays the ITEM-SPECIFIC give one-shot ("give_<key>") when its
// sheet exists, falling back to the generic "give_item" (which itself falls
// back to the grab frames at load). 2026-06-11: each hand-over should show
// the actual item leaving PP's paw - drop a "PP give <item>.png" sheet in
// and the matching trade upgrades automatically, no code change.
func (p *player) playGive(itemKey string, dur float64, onDone func()) {
	name := "give_" + itemKey
	if p.oneShotAnims == nil {
		if onDone != nil {
			onDone()
		}
		return
	}
	if _, ok := p.oneShotAnims[name]; !ok {
		name = "give_item"
	}
	p.playOneShot(name, dur, onDone)
}

// playOneShot runs a registered one-shot animation for `dur` seconds, then
// fires onDone. If the named anim isn't registered, onDone is invoked
// immediately so sequences don't deadlock on missing assets.
// hasOneShot reports whether a one-shot animation is registered (its art
// landed). Callers fall back to a generic action when it isn't.
func (p *player) hasOneShot(name string) bool {
	_, ok := p.oneShotAnims[name]
	return ok
}

func (p *player) playOneShot(name string, dur float64, onDone func()) {
	if p.oneShotAnims == nil {
		if onDone != nil {
			onDone()
		}
		return
	}
	if _, ok := p.oneShotAnims[name]; !ok {
		if onDone != nil {
			onDone()
		}
		return
	}
	if dur <= 0 {
		dur = 1.0
	}
	p.activeOneShot = name
	p.oneShotIdx = 0
	p.oneShotTimer = 0
	p.oneShotElapsed = 0
	p.oneShotDuration = dur
	p.oneShotOnDone = onDone
	p.state = stateIdle
	p.moving = false
}

func (p *player) currentWalkFrames() []spriteFrame {
	switch p.dir {
	case dirUp:
		return p.walkUpFrames
	case dirDown:
		return p.walkDownFrames
	default:
		return p.walkSideFrames
	}
}

func (p *player) currentIdleFrames() []spriteFrame {
	if p.seated && len(p.seatedIdle) > 0 {
		return p.seatedIdle
	}
	switch p.dir {
	case dirUp:
		return p.idleBackFrames
	case dirLeft, dirRight:
		return p.idleSideFrames
	default:
		return p.idleFrontFrames
	}
}

// setScaredTalk (2026-07-24 user #34): swaps PP's side-talk to the frightened
// variant for the rude-Higgins exchange. No-ops until the §PP-TALK-SIDE-SCARED
// sheet lands (the normal side talk plays meanwhile).
func (p *player) setScaredTalk(on bool) {
	p.scaredTalk = on && len(p.talkSideScaredFrames) > 0
}

func (p *player) currentTalkFrames() []spriteFrame {
	if p.seated && len(p.seatedTalk) > 0 {
		return p.seatedTalk
	}
	switch p.dir {
	case dirLeft, dirRight:
		if p.scaredTalk {
			return p.talkSideScaredFrames
		}
		return p.talkSideFrames
	case dirUp:
		// Talking with his back to the camera (ppFaceBack NPCs). Use the
		// dedicated back-talk sheet once it lands; until then fall back to the
		// back IDLE so PP keeps facing away rather than snapping to front talk.
		if len(p.talkBackFrames) > 0 {
			return p.talkBackFrames
		}
		if len(p.idleBackFrames) > 0 {
			return p.idleBackFrames
		}
		return p.talkFrames
	default:
		return p.talkFrames
	}
}

func firstAvailableFrame(groups ...[]spriteFrame) spriteFrame {
	for _, group := range groups {
		if len(group) > 0 {
			return group[0]
		}
	}
	return spriteFrame{}
}

func (p *player) actionFrames() []spriteFrame {
	switch p.state {
	case stateGrabbing:
		return p.grabFrames
	case stateUsing:
		return p.useItemFrames
	case stateExamining:
		return p.examineFrames
	case stateReacting:
		return p.reactFrames
	case stateShowInventory:
		return p.showInvFrames
	}
	return nil
}

func (p *player) currentSprite() spriteFrame {
	// One-shot anim (sequence player) overrides state-based selection.
	if p.activeOneShot != "" {
		if frames, ok := p.oneShotAnims[p.activeOneShot]; ok && len(frames) > 0 {
			return frames[p.oneShotIdx%len(frames)]
		}
	}
	switch p.state {
	case stateWalking:
		frames := p.currentWalkFrames()
		if len(frames) == 0 {
			return firstAvailableFrame(p.currentIdleFrames(), p.walkDownFrames, p.walkSideFrames, p.walkUpFrames)
		}
		return frames[p.walkCycleIdx%len(frames)]
	case stateTalking:
		frames := p.currentTalkFrames()
		if len(frames) > 0 {
			return frames[p.talkCycleIdx%len(frames)]
		}
		return firstAvailableFrame(p.currentIdleFrames(), p.walkDownFrames, p.walkSideFrames, p.walkUpFrames)
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		frames := p.actionFrames()
		if len(frames) > 0 {
			return frames[p.actionIdx%len(frames)]
		}
		return firstAvailableFrame(p.currentIdleFrames(), p.walkDownFrames)
	default:
		frames := p.currentIdleFrames()
		if len(frames) > 0 {
			idx := int(p.breathTimer*4) % len(frames)
			return frames[idx]
		}
		return firstAvailableFrame(p.walkDownFrames, p.walkSideFrames, p.walkUpFrames)
	}
}

func (p *player) minY() float64 {
	if p.sceneMinY > 0 {
		return p.sceneMinY
	}
	return playerMinY
}

func (p *player) maxY() float64 {
	if p.sceneMaxY > 0 {
		return p.sceneMaxY
	}
	return playerMaxY
}

func (p *player) setTarget(x, y float64) {
	// 2026-07-25 (user): frozen through a counter give chain — ignore walks.
	if p.movementLocked {
		return
	}
	// #28: a move-elsewhere click grows PP back to full size if he was holding
	// the shrunk Pierre-depth pose after a dialog.
	// 2026-06-11 #19: TWO-STAGE exit from the receded depth - PP first walks
	// straight TOWARD the camera (down to the front row, growing to full
	// size), and only then heads for the clicked point. The old diagonal
	// beeline read wrong ("walk front until full size, then to the dot").
	if p.recedeHeld {
		p.releaseRecedeSmooth(0.5)
		fx, fy := x, y
		p.targetX = engine.Clamp(p.x, playerMinX, playerMaxX)
		p.targetY = p.maxY()
		p.moving = true
		p.allowOffscreen = false
		p.state = stateWalking
		p.dir = dirDown
		p.interactTarget = nil
		p.onArrival = func() {
			p.setTarget(fx, fy)
		}
		return
	}
	tx := engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	ty := engine.Clamp(y-playerDstH/2, p.minY(), p.maxY())
	p.targetX = tx
	p.targetY = ty
	p.interactTarget = nil
	p.onArrival = nil
	// 2026-07-15 (user #7): plain clicks (incl. the bakery exit) also walk the
	// aisle L instead of cutting across the tables.
	if p.startLWalk(tx, ty, tx-p.x, ty-p.y, func() {}) {
		return
	}
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
}

// talkApproachPos returns the world top-left (tx,ty) PP should stand at to talk
// to target - picking the interior side, avoiding overlap, honoring approachRight,
// and foot-aligning to the NPC. Shared by walkToAndInteract and the held-item
// hand-off so both place PP identically.
func (p *player) talkApproachPos(target *npc) (float64, float64) {
	npcLeft := float64(target.bounds.X)
	npcRight := float64(target.bounds.X + target.bounds.W)
	npcCenter := float64(target.bounds.X + target.bounds.W/2)
	gap := float64(target.approachGapX)
	if gap <= 0 {
		gap = 10
	}
	pickSide := func(preferRight bool) float64 {
		if preferRight {
			return npcRight + gap
		}
		return npcLeft - playerDstW - gap
	}
	preferred := npcCenter >= engine.ScreenWidth/2
	tx := engine.Clamp(pickSide(!preferred), playerMinX, playerMaxX)
	if tx < npcRight && tx+playerDstW > npcLeft {
		tx = engine.Clamp(pickSide(preferred), playerMinX, playerMaxX)
	}
	// #7: far-right kid (Danny) - force the right side so PP doesn't approach
	// from the left and end up standing on Marcus.
	if target.approachRight {
		tx = engine.Clamp(pickSide(true), playerMinX, playerMaxX)
	}
	// 2026-06-11 #2: mirror flag - Jake is approached from his LEFT so PP
	// doesn't stand on Lily (and faces Jake's back at the campfire).
	if target.approachLeft {
		tx = engine.Clamp(pickSide(false), playerMinX, playerMaxX)
	}
	ty := p.y
	if !target.elevated {
		ty = float64(target.bounds.Y+target.bounds.H) - playerDstH + 4
	}
	// 2026-06-11 #25: per-NPC stand row (bakery patrons - PP stands in the
	// aisle between Poulain's counter and the tables instead of below them).
	if target.approachYOverride > 0 {
		ty = float64(target.approachYOverride)
	}
	// 2026-06-24 (#11/#14): pin PP's foot-center x to an exact mark when set.
	if target.approachXOverride > 0 {
		tx = engine.Clamp(float64(target.approachXOverride)-playerDstW/2, playerMinX, playerMaxX)
	}
	return tx, engine.Clamp(ty, p.minY(), p.maxY())
}

// walkToTalkPos walks PP to the approach position and fires onArrive there -
// even when PP is already adjacent (snaps + fires immediately). Used by the
// held-item hand-off (interactTarget is NOT set), so the give-item beat works at
// any distance (#10), not only when PP has to walk in. Does not touch
// interactTarget, so it won't also trigger the normal NPC dialog.
// faceNPC turns PP toward the NPC he's talking/handing to — one source of
// truth for the dialog AND the held-item hand-off (2026-07-15 user #5: the
// hand-off path never set facing, so PP kept his walk direction and gave
// Margaux the heel facing away).
func (p *player) faceNPC(n *npc) {
	npcCenter := float64(n.bounds.X + n.bounds.W/2)
	p.facingLeft = p.x+playerDstW/2 > npcCenter
	if n.ppTalkFlip {
		p.facingLeft = !p.facingLeft
	}
	if p.facingLeft {
		p.dir = dirLeft
	} else {
		p.dir = dirRight
	}
}

func (p *player) walkToTalkPos(target *npc, onArrive func()) {
	// 2026-06-15 #13: same as walkToAndInteract - a held-item hand-over to
	// another NPC (e.g. handing Claude the Press Pass) must release any held
	// Pierre-depth recede so PP isn't stuck shrunk.
	if p.recedeHeld {
		p.releaseRecedeSmooth(0.5)
	}
	tx, ty := p.talkApproachPos(target)
	p.targetX = tx
	p.targetY = ty
	dx := tx - p.x
	dy := ty - p.y
	if dx*dx+dy*dy < 30*30 {
		p.x = tx
		p.y = ty
		p.moving = false
		p.state = stateIdle
		p.faceNPC(target)
		if onArrive != nil {
			onArrive()
		}
		return
	}
	arrive := func() {
		p.faceNPC(target)
		if onArrive != nil {
			onArrive()
		}
	}
	// user #13: L-shaped approach in aisle scenes.
	if p.startLWalk(tx, ty, dx, dy, func() { p.onArrival = arrive }) {
		return
	}
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
	p.onArrival = arrive
}

// startLWalk (2026-07-15 user #13, paris_bakery): when approaching a talk
// mark on a different row, walk the two legs of an L instead of an oblique
// beeline across the tables — moving UP: horizontal leg along the current
// (aisle) row first, then straight up at the target x; moving DOWN: straight
// down into the aisle first, then horizontal. `arm` re-arms the caller's
// arrival behavior (onArrival / interactTarget) for the second leg. Returns
// false when the scene doesn't use L-walks or the move is single-axis.
func (p *player) startLWalk(tx, ty, dx, dy float64, arm func()) bool {
	if !p.lWalkApproach {
		return false
	}
	// 2026-07-23 (#9): near-horizontal moves used to skip the L entirely —
	// but Poulain's counter row and the patron row differ by only ~7px, so a
	// counter→patron click beelined straight ACROSS the table band and got
	// shoved into a table. If the horizontal corridor at the current foot row
	// crosses a table footBlocker, detour through the front lane instead:
	// down, across, back up.
	if dy*dy <= 40*40 || dx*dx <= 40*40 {
		if dx*dx <= 40*40 {
			return false // vertical-ish moves can't cross the table band sideways
		}
		footY := p.y + playerDstH
		x0, x1 := p.x+playerDstW/2, tx+playerDstW/2
		if x1 < x0 {
			x0, x1 = x1, x0
		}
		laneFoot := 0.0
		for _, b := range p.lastFootBlockers {
			if footY >= float64(b.Y) && footY <= float64(b.Y+b.H) &&
				x1 >= float64(b.X) && x0 <= float64(b.X+b.W) {
				if bottom := float64(b.Y+b.H) + 15; bottom > laneFoot {
					laneFoot = bottom
				}
			}
		}
		if laneFoot == 0 {
			return false // corridor clear — the straight walk is fine
		}
		laneY := math.Min(laneFoot-playerDstH, p.maxY())
		p.targetX = p.x
		p.targetY = laneY
		p.onArrival = func() {
			p.targetX = tx
			p.targetY = laneY
			p.moving = true
			p.state = stateWalking
			p.onArrival = func() {
				p.targetX = tx
				p.targetY = ty
				p.moving = true
				p.state = stateWalking
				arm()
			}
		}
		p.moving = true
		p.allowOffscreen = false
		p.state = stateWalking
		return true
	}
	if dy < 0 {
		p.targetX = tx
		p.targetY = p.y
	} else {
		p.targetX = p.x
		p.targetY = ty
	}
	p.onArrival = func() {
		p.targetX = tx
		p.targetY = ty
		p.moving = true
		p.state = stateWalking
		arm()
	}
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
	return true
}

func (p *player) walkToAndInteract(target *npc, ds *dialogSystem) {
	// 2026-06-15 #13: PP could carry the held Pierre-depth recede (recedeHeld)
	// into a click on ANOTHER NPC (Claude). Only setTarget released it, so PP
	// stayed shrunk through and after that NPC's dialog. Release it here too so
	// walking to any NPC grows him back to full size.
	if p.recedeHeld {
		p.releaseRecedeSmooth(0.5)
	}
	tx, ty := p.talkApproachPos(target)
	p.targetX = tx
	p.targetY = ty
	p.interactTarget = target
	p.dialogSys = ds
	p.onArrival = nil
	// User 2026-05-10: clicking an NPC must open the dialog as the direct result
	// of THE CLICK. If PP is already adjacent, snap + fire now instead of a 1-2s
	// walk that reads as "my click did nothing"; otherwise the movement-arrival
	// handler fires startNPCDialog via interactTarget.
	dx := tx - p.x
	dy := ty - p.y
	if dx*dx+dy*dy < 30*30 {
		p.x = tx
		p.y = ty
		p.moving = false
		p.state = stateIdle
		p.startNPCDialog()
		return
	}
	// user #13: L-shaped approach in aisle scenes — interactTarget must only
	// arm on the SECOND leg, or the dialog would fire at the corner point.
	if p.lWalkApproach {
		p.interactTarget = nil
		if p.startLWalk(tx, ty, dx, dy, func() {
			p.interactTarget = target
			p.dialogSys = ds
		}) {
			return
		}
		p.interactTarget = target
	}
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
}

func (p *player) walkToAndDo(x, y float64, action func()) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(y-playerDstH/2, p.minY(), p.maxY())
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = action
}

// walkToAndDoViaLanes (2026-07-24 #11): walkToAndDo that honors the bakery
// lane detour. Exit-hotspot walks used the straight walkToAndDo and beelined
// across the table footBlockers ("he walked through the table"); this runs
// the same corridor/L check startLWalk performs and falls back to the plain
// walk in scenes without lWalkApproach.
func (p *player) walkToAndDoViaLanes(x, y float64, action func()) {
	tx := engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	ty := engine.Clamp(y-playerDstH/2, p.minY(), p.maxY())
	p.interactTarget = nil
	if p.startLWalk(tx, ty, tx-p.x, ty-p.y, func() { p.onArrival = action }) {
		return
	}
	p.walkToAndDo(x, y, action)
}

// walkToAndDoUnclamped (2026-07-23 #3): like walkToAndDo but WITHOUT the
// walk-band Y clamp. Scripted exit walks only (room door at the bottom of
// the frame sits far below the room's maxY, so the clamped walk barely moved
// PP - "just the frame changed"). X stays clamped to the screen.
// 2026-07-24 (#2/#31): allowOffscreen=true — the unclamped TARGET was
// useless while update() still clamped the MOVEMENT to the walk band each
// frame, so PP jammed at the band edge and the arrival (scene transition)
// never fired (the jake_room "stuck" oscillation).
func (p *player) walkToAndDoUnclamped(x, y float64, action func()) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = y - playerDstH/2
	p.moving = true
	p.allowOffscreen = true
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = action
}

// playRecede starts the "walk into the distance" tween used by the camp
// entrance transition. PP stays anchored near his current X (no left-drift),
// faces upward so the back-walk frames cycle, drifts up by dyUp pixels over
// dur seconds, and shrinks from 1.0 to endScale. onDone fires once the tween
// completes - typically a sceneMgr.transitionTo call.
//
// User 2026-04-26: replaces the old walkToAndDo(599, 200, ...) that read as
// "walking left then up". Recede is intended to read as "walking away into
// the camp" without a strafe.
// clearRecede resets the recede tween state so PP renders at full size
// again. Call after a scene transition repositions the player so the
// next scene starts fresh - without this PP would keep rendering at
// recedeEndScale after the cabin/camp-entrance transition completes.
func (p *player) clearRecede() {
	p.recedeActive = false
	p.recedeReleasing = false
	p.recedeHeld = false
	p.recedeScale = 1.0
	p.recedeElapsed = 0
	p.recedeOnDone = nil
	if p.state == stateWalking && !p.moving {
		p.state = stateIdle
	}
}

// holdRecede (#28) freezes PP at the current receded scale but RE-ENABLES
// movement. He renders shrunk until the next setTarget grows him back. Use it
// when a recede-framed dialog ends and PP should stay small (Pierre) rather
// than ease back immediately.
func (p *player) holdRecede() {
	p.recedeActive = false
	p.recedeReleasing = false
	p.recedeHeld = true
	if p.state == stateWalking && !p.moving {
		p.state = stateIdle
	}
}

// releaseRecedeSmooth ends a recede freeze the gentle way (#16): movement is
// re-enabled immediately, but PP's rendered scale eases from its current
// shrunk value back to full size over dur seconds instead of snapping. Use
// this when a recede-framed dialog ends and PP should walk away without a
// jarring size pop.
func (p *player) releaseRecedeSmooth(dur float64) {
	if !p.recedeActive && !p.recedeReleasing && !p.recedeHeld {
		return
	}
	if dur <= 0 {
		dur = 0.6
	}
	p.recedeReleaseFrom = p.recedeScale
	p.recedeReleaseElapsed = 0
	p.recedeReleaseDur = dur
	p.recedeReleasing = true
	p.recedeActive = false // un-freeze movement; the release tween owns scale
	p.recedeHeld = false
	if p.state == stateWalking && !p.moving {
		p.state = stateIdle
	}
}

func (p *player) playRecede(dur, endScale, dyUp float64, onDone func()) {
	if dur <= 0 {
		dur = 1.0
	}
	if endScale <= 0 || endScale > 1 {
		endScale = 0.4
	}
	p.recedeActive = true
	p.recedeStartX = p.x
	p.recedeStartY = p.y
	p.recedeEndScale = endScale
	p.recedeDyUp = dyUp
	p.recedeDuration = dur
	p.recedeElapsed = 0
	p.recedeScale = 1.0
	p.recedeReleasing = false
	p.recedeHeld = false
	p.recedeOnDone = onDone
	p.dir = dirUp
	p.facingLeft = false
	p.state = stateWalking
	p.moving = false // movement is driven by the tween, not setTarget
	p.interactTarget = nil
	p.onArrival = nil
}

// playWalkIn lerps PP horizontally from startX to endX at fixed y over
// `dur` seconds, with the side-walk frames cycling. Used for the opening
// "PP enters camp from off-screen-left" cutscene. Drives position
// directly so it's immune to the clamp / blocker / setTarget machinery.
// On completion, fires `onDone` (typically the monologue start).
func (p *player) playWalkIn(startX, endX, y, dur float64, onDone func()) {
	if dur <= 0 {
		dur = 2.0
	}
	p.walkInActive = true
	p.walkInStartX = startX
	p.walkInEndX = endX
	p.walkInY = y
	p.walkInDuration = dur
	p.walkInElapsed = 0
	p.walkInOnDone = onDone
	p.x = startX
	p.y = y
	p.dir = dirRight
	p.facingLeft = false
	p.state = stateWalking
	p.moving = false
	p.allowOffscreen = false
	p.interactTarget = nil
	p.onArrival = nil
}

func (p *player) walkToExit(dir arrowDir, action func()) {
	p.targetY = engine.Clamp(p.y, p.minY(), p.maxY())
	switch dir {
	case arrowLeft:
		p.targetX = -playerDstW
		p.dir = dirLeft
		p.facingLeft = true
	case arrowRight:
		p.targetX = engine.ScreenWidth + playerDstW
		p.dir = dirRight
		p.facingLeft = false
	case arrowDown:
		p.targetX = p.x
		p.targetY = engine.ScreenHeight + playerDstH
		p.dir = dirDown
		p.facingLeft = false
	case arrowUp:
		p.targetX = p.x
		p.targetY = -playerDstH
		p.dir = dirUp
		p.facingLeft = false
	case arrowDownRight:
		p.targetX = engine.ScreenWidth + playerDstW
		p.targetY = engine.ScreenHeight + playerDstH
		p.dir = dirDown
		p.facingLeft = false
	default:
		p.allowOffscreen = false
		p.walkToAndDo(p.x+playerDstW/2, p.y+playerDstH/2, action)
		return
	}
	p.moving = true
	p.allowOffscreen = true
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = action
}

func (p *player) update(dt float64, blockers []sdl.Rect, footBlockers []sdl.Rect) {
	p.breathTimer += dt
	p.lastFootBlockers = footBlockers

	// One-shot named anim (sequence player). Frames advance evenly across
	// the requested duration; on completion the registered callback fires.
	// Short-circuits the rest of update so the anim can't be interrupted by
	// state machine churn while it's running.
	if p.activeOneShot != "" {
		frames := p.oneShotAnims[p.activeOneShot]
		// User 2026-05-12: track wall-clock elapsed directly so the one-shot
		// completes after `oneShotDuration` seconds regardless of frame
		// pacing. The previous totalElapsed-from-idx formula stopped
		// advancing once idx hit the last frame (idx capped at len-1 while
		// the timer kept decrementing), so the one-shot never crossed the
		// completion threshold and PP froze on the final grab frame.
		p.oneShotTimer += dt
		stepLen := actionFrameTime
		if p.oneShotDuration > 0 && len(frames) > 0 {
			stepLen = p.oneShotDuration / float64(len(frames))
		}
		p.oneShotElapsed += dt
		// Frame index follows elapsed time, capped at the last frame so
		// the final pose holds until completion.
		if p.oneShotDuration > 0 && len(frames) > 0 {
			idx := int(p.oneShotElapsed / stepLen)
			if idx >= len(frames) {
				idx = len(frames) - 1
			}
			p.oneShotIdx = idx
		} else if p.oneShotTimer >= stepLen && len(frames) > 0 {
			p.oneShotTimer -= stepLen
			if p.oneShotIdx < len(frames)-1 {
				p.oneShotIdx++
			}
		}
		if p.oneShotDuration > 0 && p.oneShotElapsed >= p.oneShotDuration {
			cb := p.oneShotOnDone
			p.activeOneShot = ""
			p.oneShotIdx = 0
			p.oneShotTimer = 0
			p.oneShotElapsed = 0
			p.oneShotOnDone = nil
			// User 2026-05-12: reset state to idle BEFORE the callback so PP
			// doesn't appear "stuck" mid-walk after a grab/receive one-shot
			// completes. Without this, if PP was stateWalking when the
			// one-shot started (he walked to the flower), state stayed
			// walking after - the user sees PP frozen in a walk pose.
			p.state = stateIdle
			p.moving = false
			if cb != nil {
				cb()
			}
		}
		return
	}

	// Recede tween (camp-entrance transition) - drives position+scale itself
	// and short-circuits the rest of update so movement / blocker code can't
	// fight the tween. Walk-frame ticking still runs so PP cycles his back
	// walk during the recede.
	if p.recedeActive {
		p.recedeElapsed += dt
		t := p.recedeElapsed / p.recedeDuration
		if t > 1 {
			t = 1
		}
		p.y = p.recedeStartY - p.recedeDyUp*t
		p.recedeScale = 1.0 - (1.0-p.recedeEndScale)*t

		p.walkTimer += dt
		if p.walkTimer >= walkFrameTimeFrontBack {
			p.walkTimer -= walkFrameTimeFrontBack
			frames := p.currentWalkFrames()
			if len(frames) > 0 {
				p.walkCycleIdx = (p.walkCycleIdx + 1) % len(frames)
			}
		}

		if t >= 1 {
			// User 2026-05-10: don't reset scale to 1.0 here - that lets
			// one frame render at full size between this update and the
			// scene transition's first fade tick, which the user sees as
			// "shrink, then a beat at full size, then transition". Clamp
			// to endScale and keep recedeActive=true so subsequent draws
			// stay scaled-down. The scene transition fires from cb(), and
			// scene.transitionTo's new-spawn handler clears the recede
			// state (player.clearRecede) once the player is repositioned.
			p.recedeScale = p.recedeEndScale
			if cb := p.recedeOnDone; cb != nil {
				p.recedeOnDone = nil
				cb()
			}
		}
		return
	}

	// Smooth recede release (#16): eases recedeScale back to 1.0 without
	// freezing movement, so PP grows back to full size gradually after a
	// recede-framed dialog (e.g. Pierre) instead of popping instantly.
	if p.recedeReleasing {
		p.recedeReleaseElapsed += dt
		t := p.recedeReleaseElapsed / p.recedeReleaseDur
		if t >= 1 {
			t = 1
			p.recedeReleasing = false
		}
		p.recedeScale = p.recedeReleaseFrom + (1.0-p.recedeReleaseFrom)*t
	}

	// Walk-in tween (opening cutscene). Drives p.x directly via lerp +
	// cycles the side-walk frames. Short-circuits the rest of update so
	// the normal moving/clamp/blocker pipeline can't fight it (the same
	// pattern used by playRecede above).
	if p.walkInActive {
		p.walkInElapsed += dt
		t := p.walkInElapsed / p.walkInDuration
		if t > 1 {
			t = 1
		}
		p.x = p.walkInStartX + (p.walkInEndX-p.walkInStartX)*t
		p.y = p.walkInY

		p.walkTimer += dt
		if p.walkTimer >= walkFrameTime {
			p.walkTimer -= walkFrameTime
			frames := p.currentWalkFrames()
			if len(frames) > 0 {
				p.walkCycleIdx = (p.walkCycleIdx + 1) % len(frames)
			}
		}

		if t >= 1 {
			p.walkInActive = false
			p.state = stateIdle
			if cb := p.walkInOnDone; cb != nil {
				p.walkInOnDone = nil
				cb()
			}
		}
		return
	}

	if p.moving {
		p.walkTimer += dt
		if ft := p.walkFrameSecs(); p.walkTimer >= ft {
			p.walkTimer -= ft
			frames := p.currentWalkFrames()
			if len(frames) > 0 {
				p.walkCycleIdx = (p.walkCycleIdx + 1) % len(frames)
			}
		}
	} else {
		p.walkCycleIdx = 0
		p.walkTimer = 0
		p.allowOffscreen = false
	}

	if p.state == stateTalking {
		// #2: PP's mouth animates only while HE is the speaker and the line is
		// still revealing; otherwise it holds frame 0 (mouth closed). So the
		// mouth tracks the text and stops when it's the NPC's turn or the line
		// has finished appearing.
		ppSpeaking := p.dialogSys != nil && p.dialogSys.isRevealing() &&
			p.dialogSys.currentSpeaker() == "Pink Panther"
		if ppSpeaking {
			p.talkTimer += dt
			if p.talkTimer >= talkFrameTime {
				p.talkTimer -= talkFrameTime
				frames := p.currentTalkFrames()
				if len(frames) > 0 {
					p.talkCycleIdx = (p.talkCycleIdx + 1) % len(frames)
				}
			}
		} else {
			p.talkTimer = 0
			p.talkCycleIdx = 0
		}
	} else {
		p.talkTimer = 0
		p.talkCycleIdx = 0
	}

	switch p.state {
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		p.actionTimer += dt
		if p.actionTimer >= actionFrameTime {
			p.actionTimer -= actionFrameTime
			frames := p.actionFrames()
			if len(frames) > 0 {
				p.actionIdx++
				if p.actionIdx >= len(frames) {
					p.actionIdx = 0
					p.state = stateIdle
					if p.actionCallback != nil {
						fn := p.actionCallback
						p.actionCallback = nil
						fn()
					}
				}
			}
		}
	default:
		p.actionIdx = 0
		p.actionTimer = 0
	}

	if !p.moving {
		if p.state != stateTalking {
			p.state = stateIdle
			p.dir = dirDown
			p.facingLeft = false
		}
		return
	}

	dx := p.targetX - p.x
	dy := p.targetY - p.y
	dist := math.Sqrt(dx*dx + dy*dy)

	// Determine direction from movement delta. User 2026-05-22: threshold
	// dropped 1.2 → 0.8 so vertical motion wins more often. Previously
	// a click at down-right (dx, dy both positive and roughly equal)
	// fell into the horizontal branch and PP played his side-walk sprite
	// even when the user expected a forward (down) approach.
	if math.Abs(dy) > math.Abs(dx)*0.8 {
		if dy < 0 {
			p.dir = dirUp
		} else {
			p.dir = dirDown
		}
	} else {
		if dx < 0 {
			p.dir = dirLeft
		} else {
			p.dir = dirRight
		}
	}
	p.facingLeft = dx < 0

	if dist < 3.0 {
		p.x = p.targetX
		p.y = p.targetY
		p.moving = false
		p.state = stateIdle
		if p.interactTarget != nil && p.dialogSys != nil {
			p.startNPCDialog()
		}
		if p.onArrival != nil {
			fn := p.onArrival
			p.onArrival = nil
			fn()
		}
		return
	}

	speed := playerBaseSpeed
	if dist < 100 {
		speed = playerBaseSpeed * (0.3 + 0.7*dist/100.0)
	}
	nextX := p.x + (dx/dist)*speed*dt
	nextY := p.y + (dy/dist)*speed*dt
	if p.allowOffscreen {
		p.x = nextX
		p.y = nextY
	} else {
		p.x = engine.Clamp(nextX, playerMinX, playerMaxX)
		p.y = engine.Clamp(nextY, p.minY(), p.maxY())
	}

	// footBlockers (2026-07-15 user #10, bakery tables): tested against PP's
	// FOOT POINT only, so the tall body box can still overlap tables visually
	// while his feet are kept off the cloth. Same horizontal-shove resolve as
	// body blockers.
	for _, b := range footBlockers {
		fx := int32(p.x) + playerDstW/2
		fy := int32(p.y) + playerDstH
		if fx >= b.X && fx <= b.X+b.W && fy >= b.Y && fy <= b.Y+b.H {
			blockerCX := float64(b.X) + float64(b.W)/2
			if p.x+playerDstW/2 < blockerCX {
				p.x = float64(b.X) - playerDstW/2
			} else {
				p.x = float64(b.X+b.W) - playerDstW/2
			}
			if p.targetX+playerDstW/2 > float64(b.X) && p.targetX+playerDstW/2 < float64(b.X+b.W) &&
				p.targetY+playerDstH > float64(b.Y) && p.targetY+playerDstH < float64(b.Y+b.H) {
				p.moving = false
			}
		}
	}

	for _, b := range blockers {
		pr := sdl.Rect{X: int32(p.x), Y: int32(p.y), W: playerDstW, H: playerDstH}
		if pr.HasIntersection(&b) {
			playerCX := p.x + playerDstW/2
			blockerCX := float64(b.X) + float64(b.W)/2
			if playerCX < blockerCX {
				p.x = float64(b.X) - playerDstW
			} else {
				p.x = float64(b.X + b.W)
			}
			if p.targetX < float64(b.X+b.W) && p.targetX+playerDstW > float64(b.X) {
				p.moving = false
				if p.interactTarget != nil && p.dialogSys != nil {
					p.startNPCDialog()
				} else if p.onArrival != nil {
					fn := p.onArrival
					p.onArrival = nil
					p.state = stateIdle
					fn()
				} else {
					p.state = stateIdle
				}
			}
		}
	}
}

// canTriggerAltDialog gates give-item alt dialogs so they only fire when
// the required item is reachable. Two independent flags shape the rule:
//
//   - altDialogRequiresHeld: the player must actively carry the item
//     (inv.heldItem). Use this for flows that visually hand the item over
//     from the cursor, like drag-and-drop puzzles.
//   - altDialogRequiresItem: the item must exist in inventory; holding is
//     optional. Most NPCs (Lily, Jake, Tommy, Danny) use this so the
//     player just needs to click on the NPC after pickup.
//
// Both flags unset -> altDialogFunc is called unconditionally (legacy).
func (p *player) canTriggerAltDialog(n *npc) bool {
	if n == nil {
		return true
	}
	if p.inv == nil {
		return !n.altDialogRequiresHeld && n.altDialogRequiresItem == ""
	}
	if n.altDialogRequiresHeld {
		if p.inv.heldItem == nil {
			return false
		}
		if n.altDialogRequiresItem == "" {
			return true
		}
		return p.inv.heldItem.name == n.altDialogRequiresItem
	}
	if n.altDialogRequiresItem != "" {
		if p.inv.heldItem != nil && p.inv.heldItem.name == n.altDialogRequiresItem {
			return true
		}
		return p.inv.hasItem(n.altDialogRequiresItem)
	}
	return true
}

func (p *player) startNPCDialog() {
	n := p.interactTarget
	ds := p.dialogSys
	if n == nil || ds == nil {
		return
	}
	// User 2026-05-22: NPC may want to fully script the click flow (Pierre's
	// walk-to-middle → recede → talk choreography). If onClickOverride is
	// set, hand off control and return - the override owns dialog start +
	// any post-dialog state.
	if n.onClickOverride != nil {
		cb := n.onClickOverride
		p.interactTarget = nil
		cb()
		return
	}
	p.state = stateTalking
	p.talkCycleIdx = 0
	p.talkTimer = 0
	p.faceNPC(n)
	// 2026-06-20 #10: NPCs behind a back-of-scene counter (Madame Poulain) sit
	// ABOVE PP, so the default left/right facing reads as PP showing his back.
	// ppFacePlayer makes PP face the camera (front) for this NPC's dialog so we
	// see his face instead.
	if n.ppFacePlayer {
		p.dir = dirDown
		p.facingLeft = false
	}
	// 2026-06-30: ppFaceBack is the opposite of ppFacePlayer - PP stands below
	// the NPC and talks UP toward a back-of-scene counter (Madame Poulain), so
	// we want his BACK to the camera, not his face. Takes precedence.
	if n.ppFaceBack {
		p.dir = dirUp
		p.facingLeft = false
	}

	// Snapshot the NPC's authored facing so we can restore it when the
	// dialog ends, then flip the NPC to face PP. Sprite sheets are
	// drawn facing right, so flipped=true means "face left". If PP is
	// to the left of the NPC's center, the NPC needs to face left.
	n.preTalkFlipped = n.flipped
	// #16: seated NPCs (office Higgins) hold their authored facing instead of
	// turning to face PP.
	if !n.fixedFacing {
		n.flipped = p.x+playerDstW/2 < float64(n.bounds.X+n.bounds.W/2)
	}

	// 2026-07-29: one-time per-dialog hook (Takeshi opens his kamishibai
	// stage doors over the first line). Fires for rule, alt and plain paths.
	if n.onDialogStart != nil {
		n.onDialogStart()
	}

	if len(n.talkGrid) > 0 {
		n.setAnimState(npcAnimTalk)
	}

	wrapCb := func(inner func()) func() {
		target := n
		return func() {
			if len(target.talkGrid) > 0 {
				target.setAnimState(npcAnimIdle)
			}
			target.flipped = target.preTalkFlipped
			if inner != nil {
				inner()
			}
		}
	}

	// Data-driven rules fire first. If any click-trigger rule matches, the
	// NPC was already handled (dialog, give, state change, etc.) and we
	// skip the legacy closure path. Rule-less NPCs fall through unchanged.
	if n.game != nil && len(n.rules) > 0 {
		if n.game.fireTrigger(n, "click", n.rules) {
			if len(n.talkGrid) > 0 {
				n.setAnimState(npcAnimIdle)
			}
			n.flipped = n.preTalkFlipped
			p.interactTarget = nil
			return
		}
	}

	if n.altDialogFunc != nil && p.canTriggerAltDialog(n) {
		entries, cb, ho := n.altDialogFunc()
		if entries != nil {
			// PR#1 (2026-06-12): hand-over dialogs play the give beat FIRST -
			// PP hands the item, the NPC takes it, THEN the text starts. The
			// NPC drops back to idle for the beat (talk anim was armed above)
			// and re-enters talk when the dialog actually begins.
			start := func() {
				p.state = stateTalking
				if len(n.talkGrid) > 0 {
					n.setAnimState(npcAnimTalk)
				}
				ds.startDialogWithCallback(entries, wrapCb(cb))
			}
			switch {
			case ho != nil && ho.dialogFirst:
				// #18/#19: "PP asks → NPC gives" — play the dialog first, then the
				// hand-off (and the inventory callback) after the text. The give
				// runs inside the dialog's end callback so PP doesn't receive the
				// item before he's even asked.
				p.state = stateTalking
				if len(n.talkGrid) > 0 {
					n.setAnimState(npcAnimTalk)
				}
				ds.startDialogWithCallback(entries, wrapCb(func() {
					p.playHandOff(n, ho, func() {
						if cb != nil {
							cb()
						}
					})
				}))
			case ho != nil:
				if len(n.talkGrid) > 0 {
					n.setAnimState(npcAnimIdle)
				}
				p.playHandOff(n, ho, start)
			default:
				start()
			}
			p.interactTarget = nil
			return
		}
	}

	// Strict gate (user 2026-05-21): if the NPC has a give-item gate but PP
	// doesn't satisfy it (no held item / not in inventory), AND the NPC has
	// a strict-missing hint dialog set, play the hint instead of falling
	// through to the regular n.dialog. Lily is the canonical case - without
	// this, clicking her after the hint phase replays her flower-thanks
	// dialog every time, even when PP doesn't have a flower.
	if n.altDialogFunc != nil &&
		(n.altDialogRequiresHeld || n.altDialogRequiresItem != "") &&
		!p.canTriggerAltDialog(n) &&
		len(n.altDialogStrictMissingHint) > 0 {
		ds.startDialogWithCallback(n.altDialogStrictMissingHint, wrapCb(nil))
		p.interactTarget = nil
		return
	}

	cb := n.onDialogEnd
	ds.startDialogWithCallback(n.dialog, wrapCb(func() {
		if cb != nil {
			cb()
		}
		n.dialogDone = true
	}))
	p.interactTarget = nil
}

func (p *player) playAction(s playerState, cb func()) {
	p.state = s
	p.actionIdx = 0
	p.actionTimer = 0
	p.actionCallback = cb
}

// containsPoint is the "click is on PP" test used to open the inventory bag.
// User playtest #22: the full 200×270 dst rect was too wide - the bag opened
// when clicking well to the side of PP's slim body, so it felt like a big
// invisible radius. Tightened to a 100-px-wide band centered on his body and
// dropped 24px from the very top (his head is narrow), so the bag only opens
// when the click is actually on PP. Clicks just outside now fall through to
// hotspots / NPCs behind him.
func (p *player) containsPoint(x, y int32) bool {
	pt := sdl.Point{X: x, Y: y}
	const bodyW = 100
	r := sdl.Rect{
		X: int32(p.x) + (playerDstW-bodyW)/2,
		Y: int32(p.y) + 24,
		W: bodyW,
		H: playerDstH - 24,
	}
	return pt.InRect(&r)
}

func (p *player) footCenter() (int32, int32) {
	cx := int32(p.x) + playerDstW/2
	fy := int32(p.y) + playerDstH
	return cx, fy
}

func (p *player) footY() int32 {
	_, fy := p.footCenter()
	return fy
}

func (p *player) depthScale() float64 {
	progress := engine.Clamp((p.y-playerMinY)/(playerMaxY-playerMinY), 0, 1)
	// User 2026-06-02 (#7): the old 0.88..1.06 range shrank PP noticeably when
	// he walked up-screen to talk to a kid (Lily). Narrow to 0.95..1.05 so
	// depth still reads but PP no longer visibly "shrinks" mid-conversation.
	return 0.95 + progress*0.10
}

// playerRenderFillFrac: the tallest opaque pose in the active animation is
// scaled to this fraction of playerDstH. Normalising by opaque height (not the
// raw cell height) keeps PP one consistent on-screen size across idle/talk/walk
// in every direction, even though the front sheets use 384px cells and the side
// sheets 512px cells with different amounts of blank padding (#1/#2/#3).
const playerRenderFillFrac = 0.78

// maxOpaqueHeightP is the tallest opaque content height across a frame slice -
// the player-side mirror of npc.maxOpaqueH.
func maxOpaqueHeightP(frames []spriteFrame) int32 {
	var m int32
	for _, f := range frames {
		if f.oh > m {
			m = f.oh
		}
	}
	return m
}

// currentGroup returns the frame slice the active animation is cycling, so the
// draw path can normalise size by the tallest pose in the set.
func (p *player) currentGroup() []spriteFrame {
	if p.activeOneShot != "" {
		if frames, ok := p.oneShotAnims[p.activeOneShot]; ok && len(frames) > 0 {
			return frames
		}
	}
	switch p.state {
	case stateWalking:
		if f := p.currentWalkFrames(); len(f) > 0 {
			return f
		}
	case stateTalking:
		if f := p.currentTalkFrames(); len(f) > 0 {
			return f
		}
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		if f := p.actionFrames(); len(f) > 0 {
			return f
		}
	}
	return p.currentIdleFrames()
}

// drawScaled renders PP with a character-scale multiplier so tight
// indoor scenes can shrink PP to match the PTP "pub shot" ratios
// without altering the underlying 170x235 hitbox.
func (p *player) drawScaled(renderer *sdl.Renderer, charScale float64) {
	if charScale <= 0 {
		charScale = 1.0
	}
	// Recede tween multiplies into the same render scale path. The smooth
	// release (#16) keeps applying the easing-back scale after the freeze ends,
	// and the held state (#28) keeps PP shrunk after a Pierre dialog until he
	// moves, so he doesn't pop to full size in either case.
	if (p.recedeActive || p.recedeReleasing || p.recedeHeld) && p.recedeScale > 0 {
		charScale *= p.recedeScale
	}
	frame := p.currentSprite()
	if frame.tex == nil || frame.h == 0 {
		return
	}

	// Keep the historical scale (full cell → playerDstH) so PP's overall size
	// is unchanged, but sample only the opaque content box and anchor it by
	// feet (bottom) + horizontal centre. This removes the per-frame "jump"
	// (each pose plants its feet on the same line) and the left-of-path drift,
	// and makes idle/talk/walk read as the same size. Falls back to the full
	// cell when no opaque data is present.
	scaledHeight := float64(playerDstH) * p.depthScale() * charScale
	var src *sdl.Rect
	var dstW, dstH, dstX int32
	// Normalise by the tallest opaque pose in the active animation so every
	// sheet (front 384px cells, side 512px cells, old/new art) renders PP at the
	// same standing height - fixes "not the same size" across walk/talk/idle and
	// the per-frame size jump that read as "two frames at once" (#1/#2/#3).
	refH := maxOpaqueHeightP(p.currentGroup())
	anchorX := float64(p.x) + float64(playerDstW)/2
	var footOffset float64

	// Decide the horizontal flip up-front so the foot-anchor below can account
	// for it (CopyEx mirrors within the dst rect).
	// 2026-07-18 (user #2): shift the flower grab LEFT so the crouch reach
	// meets the flower at PP's stand spot (he stands to its right).
	grabFlowerXShift := int32(0)
	if p.activeOneShot == "grab_flower" {
		grabFlowerXShift = -25
	}
	flip := sdl.FLIP_NONE
	if p.state == stateWalking && (p.dir == dirLeft || p.dir == dirRight) {
		// The side-walk sheet (PP walk left.png) is drawn FACING LEFT: show it
		// as-is when walking left, mirror it when walking right.
		if p.dir == dirRight {
			flip = sdl.FLIP_HORIZONTAL
		}
	} else if p.dir == dirLeft {
		// 2026-07-15 (user: "flip him by default"): the CURRENT regenerated
		// idle/talk SIDE sheets (PP idle side.png / PP talk side.png) are
		// drawn FACING RIGHT — verified visually (muzzle, whiskers and the
		// pointing paw all aim right). So mirror when PP faces LEFT and show
		// as-is when he faces right. (This inverts the 2026-06-30 rule, which
		// was written for the older left-facing art and made PP turn his back
		// on every NPC he talked to — Higgins included.)
		flip = sdl.FLIP_HORIZONTAL
	}
	// 2026-07-24 (#23/#28): per-anim one-shot mirror, same mechanism the
	// NPCs have (npc.oneShotFlip) — some receive sheets face the wrong way
	// relative to the giver (Jerusalem coffee, the note paper). Inverts
	// whatever the rules above decided.
	if p.activeOneShot != "" && p.oneShotFlip[p.activeOneShot] {
		if flip == sdl.FLIP_HORIZONTAL {
			flip = sdl.FLIP_NONE
		} else {
			flip = sdl.FLIP_HORIZONTAL
		}
	}

	if frame.ow > 0 && frame.oh > 0 && refH > 0 {
		targetTall := scaledHeight * playerRenderFillFrac
		frameScale := targetTall / float64(refH)
		s := sdl.Rect{X: frame.ox, Y: frame.oy, W: frame.ow, H: frame.oh}
		src = &s
		dstW = int32(float64(frame.ow) * frameScale)
		dstH = int32(float64(frame.oh) * frameScale)
		// Anchor by the animation's CONSTANT foot-centre. footCX is the SAME for
		// every frame of a sheet (stabilizeFootCX sets it to the median), and the
		// math cancels the per-frame box offset, so a fixed footCX pins the body at
		// one screen X while only the limbs/tail move. This kills the jitter from a
		// per-frame foot-centre swinging with the swishing tail. Idle/talk AND walk.
		footFromLeft := (float64(frame.footCX) - float64(frame.ox)) * frameScale
		if flip == sdl.FLIP_HORIZONTAL {
			dstX = int32(anchorX - (float64(dstW) - footFromLeft))
		} else {
			dstX = int32(anchorX - footFromLeft)
		}
		dstX += grabFlowerXShift
		// Anchor Y by the animation's CONSTANT foot ROW (sheet median, set by
		// stabilizeFootCX) instead of each frame's own content bottom - a
		// dipped tail or a higher-sitting pose extends past the line instead
		// of lifting the whole body (user 2026-06-10: stay in the same spot).
		footRow := frame.footRow
		if footRow <= frame.oy {
			footRow = frame.oy + frame.oh
		}
		footOffset = (float64(footRow) - float64(frame.oy)) * frameScale
	} else {
		// No opaque data -- fall back to the legacy full-cell scale.
		frameScale := scaledHeight / float64(frame.h)
		dstW = int32(float64(frame.w) * frameScale)
		dstH = int32(scaledHeight)
		dstX = int32(anchorX - float64(dstW)/2)
		footOffset = float64(dstH)
	}

	dstY := p.footY() - int32(footOffset)
	// User 2026-06-02 (#10): lift the flower-grab pose so PP's reach lines up
	// with the flower on the ground instead of bending below it.
	if p.activeOneShot == "grab_flower" {
		dstY -= 38
	}

	// PR#16: nudge the rolling-pin grab pose down so PP's reach meets the bike
	// basket. 2026-06-30 (#2): the old +60 (tuned for the lower spawn band)
	// shoved PP's legs under the bottom of the screen. 2026-07-15 (user #10):
	// +12 read as "picking from the air" over the basket (PP stands taller
	// than the bike) — dropped to +45 so the crouch reach lands IN the basket
	// (top y≈614) while the feet stay on-screen. F3-tune.
	if p.activeOneShot == "grab_rolling_pin" {
		dstY += 45
	}

	switch p.state {
	case stateIdle:
		breathVal := math.Sin(p.breathTimer * 2.0)
		dstY += int32(breathVal * 2.0)
	case stateTalking:
		bob := math.Sin(p.breathTimer*3.0) * 1.5
		dstY += int32(bob)
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		bob := math.Sin(p.breathTimer*2.5) * 1.0
		dstY += int32(bob)
	}

	cx, fy := p.footCenter()
	drawShadow(renderer, cx, fy, int32(float64(playerDstW-20)*p.depthScale()*charScale))

	renderer.CopyEx(frame.tex, src,
		&sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH},
		0, nil, flip)
}

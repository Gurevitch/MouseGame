package game

import (
	"fmt"
	"strings"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// Fallback monologues (used if JSON not loaded)
var openingMonologue = []dialogEntry{
	{speaker: "Pink Panther", text: "Camp Chilly Wa Wa... it's been a while."},
	{speaker: "Pink Panther", text: "The old sign is barely standing, the cabins have seen better days..."},
	{speaker: "Pink Panther", text: "But a job is a job. Time to meet the kids."},
}

var day2Monologue = []dialogEntry{
	{speaker: "Pink Panther", text: "What a night... I could hear Marcus from across the camp."},
	{speaker: "Pink Panther", text: "I need to find him. He must be in his cabin."},
	{speaker: "Pink Panther", text: "And then I should tell Higgins about this."},
}

var parisStreetMonologue = []dialogEntry{
	{speaker: "Pink Panther", text: "Ah, Paris! The city of lights, love, and... mysteries, apparently."},
	{speaker: "Pink Panther", text: "Marcus kept drawing a museum with a glass pyramid. That must be the Louvre."},
	{speaker: "Pink Panther", text: "Time to find out what he's been seeing."},
}

// louvreArrivalMonologue is PP's first-arrival beat inside the museum (#28).
// Plays once, after he walks in from the tunnel on the left.
var louvreArrivalMonologue = []dialogEntry{
	{speaker: "Pink Panther", text: "So this is the Louvre... my first time inside."},
	{speaker: "Pink Panther", text: "Marble floors, quiet halls, and somewhere in here - the painting Marcus keeps drawing."},
	{speaker: "Pink Panther", text: "That must be the curator over there. Let's have a word."},
}

// getMonologue returns dialog from JSON store, falling back to hardcoded
func (g *Game) getMonologue(name string) []dialogEntry {
	if d := g.dialogs.Get("monologues", name); d != nil {
		return d
	}
	switch name {
	case "opening":
		return openingMonologue
	case "day2":
		return day2Monologue
	case "paris_street":
		return parisStreetMonologue
	}
	return nil
}

// getNPCDialog returns dialog from JSON store, falling back to hardcoded vars
func (g *Game) getNPCDialog(file, name string) []dialogEntry {
	if d := g.dialogs.Get(file, name); d != nil {
		return d
	}
	return nil
}

type Game struct {
	renderer   *sdl.Renderer
	sceneMgr   *sceneManager
	player     *player
	dialog     *dialogSystem
	ui         *uiManager
	audio      *audioManager
	inv        *inventory
	travelMap  *travelMap
	items      *itemRegistry
	dialogs    *dialogStore
	npcDefs    *npcConfigStore
	sceneDefs  *sceneConfigStore
	vars       *VarStore
	seqPlayer  *SequencePlayer
	seqStore   *sequenceStore
	eventBus   *EventBus
	menu       *gameMenu
	devMenu    *devMenu
	clickProbe *clickProbe
	walkDbg    *walkDebug
	// bikerBumpCheck is the paris_street biker-encounter proximity trigger
	// (2026-06-12 #12); set by setupParisCallbacks, run from Update.
	bikerBumpCheck func()
	// higginsRudeStarted guards the one-time Japan-opening rude-Higgins beat so
	// the camp_grounds Update trigger fires it exactly once (transient, not saved
	// - VarHigginsRudeDone is the persistent gate).
	higginsRudeStarted bool
	// higginsWalk drives Higgins striding halfway down the path toward PP during
	// the rude intercept (a simple x-lerp run from Update); nil when idle.
	higginsWalk *higginsWalkState
	// Japan ramen stall: the closed/open prop + the waiting line that SITS at the
	// counter when Hiro opens (built in setupTokyoCallbacks, swapped by
	// openRamenStall). Art is pending so these are invisible until it lands.
	ramenStoreProp  *ambientSprite
	ramenQueue      []*ambientSprite
	ramenOpenFrames []npcFrame
	ramenSitFrames  []npcFrame
	font           *engine.BitmapFont
	lastScene      string
	mouseX         int32
	mouseY         int32

	// Story progression
	monologuePlayed bool
	day             int  // 1 = arrival/normal, 2 = weirdness begins
	day2Started     bool // Day 2 transition played
	metKids         int  // How many kids PP has talked to on Day 1
	talkedToMarcus  bool // Talked to Marcus on Day 2 (strange)
	parisUnlocked   bool // Paris available on travel map
	nightSceneDone  bool // Night campfire scene completed
	// day1BedtimeStarted latches the "Higgins bedtime speech on grounds"
	// beat so the Lily-flower handoff can't retrigger the fade-to-night
	// sequence if the callback runs twice.
	day1BedtimeStarted bool
	marcusHealed       bool // Postcard given, strange flip off, chapter 2 closed

	// Night scene
	playerSleeping   bool
	sleepingFrames   []npcFrame
	sleepingFrameIdx int
	sleepingTimer    float64
	wakingFrames     []npcFrame
	wakingPhase      int // 0=sleeping, 1=waking, 2=done
	campfireFrames   []npcFrame
	campfireFrameIdx int
	campfireTimer    float64
	// nightHidePlayer suppresses PP rendering during phase 3 (inside
	// Marcus's cabin) so the cutscene shows only Marcus freaking out,
	// even though PP is technically "present" in the marcus_room scene.
	// playerSleeping alone doesn't cover this: we flipped it to false
	// when transitioning so the sleep sprite wouldn't follow PP into the
	// cabin, but that left the walking PP visible there.
	nightHidePlayer bool

	// Flight cutscene - 4-second biplane transition between cities.
	flight *flightCutscene

	// City monologues
	parisMonologuePlayed bool

	// sceneAltBGs holds pre-loaded alternate backgrounds keyed by
	// "scene_name/variant" (e.g. "marcus_room/day", "marcus_room/night").
	// SeqSetSceneBG looks up through setSceneAltBG; reverting day-mode on
	// save load uses the same map. Replaces what used to be two named
	// *background fields - new alts drop in by adding a load call in
	// Game.New and a JSON sequence step.
	sceneAltBGs map[string]*background
}

func New(renderer *sdl.Renderer, font *engine.BitmapFont) *Game {
	g := &Game{
		renderer: renderer,
		sceneMgr: newSceneManager(renderer),
		player:   newPlayer(renderer),
		dialog:   newDialogSystem(font),
		ui:       newUIManager(font),
		audio:    newAudioManager(),
		inv:      newInventory(font, renderer),
		day:      1,
	}
	g.player.inv = g.inv
	g.lastScene = g.sceneMgr.currentName
	// Wire per-line voice playback (no-op until dialog entries have an
	// audio field set + the file lands at the path).
	g.dialog.audio = g.audio
	g.audio.playMusic(g.sceneMgr.current().musicPath)

	// Travel Map item: clicking it in the inventory opens the travel map
	// from the current scene. Replaces the camp_entrance / paris_street
	// scene-edge map hotspots (user 2026-04-26 retro-style cleanup).
	g.inv.onSelectItem = func(it *inventoryItem) bool {
		if it == nil || it.name != "Travel Map" {
			return false
		}
		g.openTravelMap(g.sceneMgr.currentName)
		return true
	}

	g.travelMap = newTravelMap(renderer)
	g.travelMap.attachGame(g)
	g.items = newItemRegistry(renderer, "assets/data/items.json")
	g.dialogs = newDialogStore("assets/data/dialog")
	g.npcDefs = newNPCConfigStore("assets/data/npc")
	g.sceneDefs = newSceneConfigStore("assets/data/scenes")
	g.vars = newVarStore()
	g.vars.Set(ScopeGame, VarChapter, ChapterCampDay1)
	g.vars.Set(ScopeGame, VarDay, 1)
	g.seqPlayer = newSequencePlayer(g)
	g.seqStore = newSequenceStore("assets/data/sequences", g)
	g.eventBus = newEventBus()
	g.menu = newGameMenu()
	g.devMenu = newDevMenu()
	g.clickProbe = newClickProbe()
	g.walkDbg = newWalkDebug()
	g.font = font
	g.attachGameToNPCs()
	g.setupCampCallbacks()
	g.setupParisCallbacks()
	g.setupJerusalemCallbacks()
	g.setupTokyoCallbacks()
	g.setupRioCallbacks()
	g.setupRomeCallbacks()
	g.setupMexicoCallbacks()
	g.setupTravelHotspots()
	g.ui.initCursors(renderer)

	// User playtest 2026-06-05: sleep/wake sheets are now authored as a SINGLE
	// ROW of 8 frames each (the old 8×2 layout had a duplicate second row, so
	// the wake-up replayed the same 8 poses twice and never read as one full
	// wake-up). Load 8×1. Aggressive color-key (tol 24) + inset 4 strips the
	// cream-white rim these sheets ship with.
	// NOTE: requires the regenerated single-row PNGs (EXTRA_PROMPTS §AD); a
	// legacy 2-row sheet loaded as 8×1 will stack both rows into each cell.
	// Store the opaque box (gf.OX/OY/OW/OH) so the draw can size PP by his actual
	// pixels, not the 192×1024 cell - the new sheets fill only part of each tall
	// cell, so the old full-cell scale drew him huge/out of place.
	// 2026-06-11 #6/#11: the AGGRESSIVE key (tol 32) ate PP's cream face and
	// chest on these sheets - the CONNECTED key only removes background
	// reachable from the cell edges, so interior light colors survive.
	sleepGrid := engine.SpriteGridFromPNGCleanConnected(renderer, "assets/images/player/pp_sleeping.png", 8, 1, 4)
	for c := 0; c < 8 && len(sleepGrid) > 0 && c < len(sleepGrid[0]); c++ {
		gf := sleepGrid[0][c]
		g.sleepingFrames = append(g.sleepingFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H, ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH})
	}
	wakeGrid := engine.SpriteGridFromPNGCleanConnected(renderer, "assets/images/player/pp_waking.png", 8, 1, 4)
	for c := 0; c < 8 && len(wakeGrid) > 0 && c < len(wakeGrid[0]); c++ {
		gf := wakeGrid[0][c]
		g.wakingFrames = append(g.wakingFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H, ox: gf.OX, oy: gf.OY, ow: gf.OW, oh: gf.OH})
	}

	// User feedback 2026-04-26: switched from the bulky campfire_idle.png
	// (8x4, rows 1-3 had bg drift) to the dedicated campfire_small.png
	// generated for the 2026-04-19 campaign - clean 6x1 grid sized to land
	// the visible flame inside the (581,592)-(702,594) target band at 1×
	// draw. Aggressive color-key + inset 4 still strip any white halo.
	fireGrid := engine.SpriteGridFromPNGCleanAggressive(renderer, "assets/images/locations/camp/campfire_small.png", 6, 1, 4)
	if len(fireGrid) > 0 {
		for c := 0; c < len(fireGrid[0]); c++ {
			gf := fireGrid[0][c]
			g.campfireFrames = append(g.campfireFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}

	g.flight = &flightCutscene{frames: loadAirplaneFrames(renderer)}

	g.sceneAltBGs = map[string]*background{
		"marcus_room/day":   newPNGBackground(renderer, "assets/images/locations/camp/background/day1/marcus_room.png"),
		"marcus_room/night": newPNGBackground(renderer, "assets/images/locations/camp/background/day1/marcus_room_night.png"),
	}

	startScene := g.sceneMgr.current()
	g.player.sceneMinY = startScene.minY
	g.player.sceneMaxY = startScene.maxY

	return g
}

// startDay2 transitions all NPCs to their "strange" dialogs
func (g *Game) startDay2() {
	if g.day >= 2 {
		return
	}
	g.day = 2

	// Make Marcus strange in his room and (re-)reveal him there. Day-1 he
	// starts hidden so peeking into the cabin doesn't pre-spoil the
	// freakout; night cutscene + Day 2 are when he's actually in his room.
	if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
		for _, n := range marcusRoom.npcs {
			if n.name == "Marcus" {
				n.hidden = false
				n.dialog = marcusStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
				break
			}
		}
	}

	// User 2026-05-23: hide all camp_grounds NPCs on Day 2 - they're each
	// in their cabin (Marcus freaking out, Lily / Danny still inside, etc.)
	// or at the office (Higgins). The Day-2 story flow runs through their
	// individual room scenes, not from the grounds. Without this, the user
	// reported "lily danny and higgins are outside on day 2".
	//
	// Tommy + Jake are already `hidden=true` after their Day-1 exit
	// sequences (tommy_exit / jake_exit) fired, but we set hidden again
	// here so Day-2 is consistent regardless of whether those sequences
	// actually ran.
	//
	// Restoring `hidden` (not `silent`) - silent leaves the sprite drawn
	// but un-clickable, which we don't want. hidden=true removes them
	// from the scene entirely.
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			switch n.name {
			case "Marcus", "Tommy", "Jake", "Lily", "Danny", "Director Higgins":
				n.hidden = true
			}
		}
	}

	if office, ok := g.sceneMgr.scenes["camp_office"]; ok {
		for _, n := range office.npcs {
			if n.name == "Director Higgins" {
				n.silent = false
				n.dialog = higginsWorriedDialog
				n.dialogDone = false
				break
			}
		}
	}

	// Restore Marcus room to day background (night cutscene is over). Routed
	// through setSceneAltBG so Marcus's strange-idle also flips to the DAY
	// lighting variant (2026-06-12).
	g.setSceneAltBG("marcus_room", "day")
}

func (g *Game) setupCampCallbacks() {
	game := g

	// --- Day 1: Higgins intro, swaps to post-dialog ---
	for _, n := range g.sceneMgr.scenes["camp_entrance"].npcs {
		if n.name == "Director Higgins" {
			higgins := n
			higgins.onDialogEnd = func() {
				if game.day == 1 {
					higgins.dialog = higginsPostDialog
				}
			}
			break
		}
	}

	// Jake in his cabin: silent by default, turns on once Marcus is healed
	// (i.e. when the Jerusalem chapter opens up). Uses his strange idle/talk
	// sheets. Healed by the Coin Rubbing from Jerusalem.
	if jakeRoom, ok := g.sceneMgr.scenes["jake_room"]; ok {
		for _, n := range jakeRoom.npcs {
			if n.name == "Jake" {
				jake := n
				jake.dialog = jakeStrangeDialog
				jake.onDialogEnd = func() {
					jake.dialog = jakePostStrangeDialog
				}
				jake.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					// #26: Jake's anchor is now the COIN from Jerusalem (was the
					// old "Coin Rubbing" stub).
					if !game.inv.hasItem("Coin") {
						return nil, nil, nil
					}
					return []dialogEntry{
							{speaker: "Pink Panther", text: "Jake. I went to Jerusalem. I found your wall, your tunnels."},
							{speaker: "Jake", text: "The... the wall? You were THERE?"},
							{speaker: "Pink Panther", text: "Here - an old coin from the gate. Look at the face on it."},
							{speaker: "Jake", text: "That's HIM. That's the face in my head. It was REAL the whole time!"},
							{speaker: "Jake", text: "The echoes... they're... quieter. Like someone closed a door in my skull."},
							{speaker: "Pink Panther", text: "Rest easy, tough guy. Two down."},
						}, func() {
							game.inv.giveItemTo("Coin", "jake")
							jake.setStrange(false)
							game.vars.SetBool(ScopeGame, VarJakeHealed, true)
							// 2026-06-24 (#39): Jake settles and drifts off, same as
							// Marcus - play the go-to-sleep one-shot, kill the strange
							// alt-idle so it can't flash back, swap to the sleeping
							// idle, and lock the pose so later chats don't wake the
							// talk sheet. Keeps the current (dark) BG. All no-op
							// gracefully until Jake's sleep art lands.
							jake.playOneShotAnim("sleep", 1.4)
							jake.altIdleAfterSec = 0
							jake.altIdleGrid = nil
							jake.altIdleActive = false
							jake.altIdleBackup = nil
							if len(jake.sleepIdleGrid) > 0 {
								jake.idleGrid = jake.sleepIdleGrid
							}
							jake.lockIdleInDialog = true
							jake.dialog = []dialogEntry{
								{speaker: "Jake", text: "Thanks for bringing me the coin... *yawn* My collection just got LEGENDARY."},
								{speaker: "Jake", text: "Everything's quiet now... I think I'll just... rest my eyes... zzz."},
							}
							// Lily's arc opens at the LAKE (Japan chapter): reveal the
							// sad Lily sitting at the end of the dock. Tokyo stays
							// LOCKED until the lake beat + the rude-Higgins exchange
							// (the camera aside is what lights up the map).
							game.vars.SetBool(ScopeGame, VarLilyArcStarted, true)
							if lk, ok := game.sceneMgr.scenes["camp_lake"]; ok {
								for _, n := range lk.npcs {
									if n.name == "Lily" {
										n.hidden = false
										n.silent = false
										n.setStrange(true)
										break
									}
								}
							}
							// Hide the cabin Lily so she isn't in two places (her whole
							// arc now happens at the lake).
							if lr, ok := game.sceneMgr.scenes["lily_room"]; ok {
								for _, n := range lr.npcs {
									if n.name == "Lily" {
										n.hidden = true
										n.silent = true
										break
									}
								}
							}
							game.dialog.queueDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "Lily's not in her cabin. Tommy says she's been sitting down by the lake for hours, just staring at the water."},
								{speaker: "Pink Panther", text: "I should go find her."},
							})
						}, &handOff{item: "Coin"}
				}
				break
			}
		}
	}

	// Lily in her cabin: silent until Jake heals. Healed by Pressed Sakura.
	if lilyRoom, ok := g.sceneMgr.scenes["lily_room"]; ok {
		for _, n := range lilyRoom.npcs {
			if n.name == "Lily" {
				lily := n
				lily.dialog = lilyStrangeDialog
				lily.onDialogEnd = func() {
					lily.dialog = lilyPostStrangeDialog
				}
				lily.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !game.inv.hasItem("Pressed Sakura") {
						return nil, nil, nil
					}
					return []dialogEntry{
							{speaker: "Pink Panther", text: "Lily. Oba-chan sent you a petal."},
							{speaker: "Lily", text: "..."},
							{speaker: "Pink Panther", text: "She said you could practice being light. Like a petal."},
							{speaker: "Lily", text: "...thank you."},
							{speaker: "Lily", text: "I... I can hear myself think again. The wind stopped whispering for me."},
							{speaker: "Pink Panther", text: "Three down."},
						}, func() {
							game.inv.giveItemTo("Pressed Sakura", "lily")
							lily.setStrange(false)
							game.vars.SetBool(ScopeGame, VarLilyHealed, true)
							lily.dialog = []dialogEntry{
								{speaker: "Lily", text: "Thank you for the petal. I'll keep it forever."},
							}
							game.travelMap.setUnlocked("rio_street", true)
							game.travelMap.setUnlocked("buenos_aires_street", true)
							if tRoom, ok := game.sceneMgr.scenes["tommy_room"]; ok {
								for _, n := range tRoom.npcs {
									if n.name == "Tommy" {
										n.silent = false
										n.setStrange(true)
										break
									}
								}
							}
							game.dialog.queueDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "Tommy next. He keeps yelling about a sister he doesn't have."},
								{speaker: "Pink Panther", text: "Rio AND Buenos Aires both lit up. This one might take two stops."},
							})
						}, &handOff{item: "Pressed Sakura"}
				}
				break
			}
		}
	}

	// Tommy in his cabin: silent until Lily heals. Healed by Dance Card.
	if tommyRoom, ok := g.sceneMgr.scenes["tommy_room"]; ok {
		for _, n := range tommyRoom.npcs {
			if n.name == "Tommy" {
				tommy := n
				tommy.dialog = tommyStrangeDialog
				tommy.onDialogEnd = func() {
					tommy.dialog = tommyPostStrangeDialog
				}
				tommy.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !game.inv.hasItem("Dance Card") {
						return nil, nil, nil
					}
					return []dialogEntry{
							{speaker: "Pink Panther", text: "Tommy. I found both halves of the dance card. Rio and Buenos Aires."},
							{speaker: "Tommy", text: "BOTH halves? That means... my sister's name is written on ONE of them?"},
							{speaker: "Pink Panther", text: "Marisa. Just the one name. But the handwriting matches yours."},
							{speaker: "Tommy", text: "Marisa! I've been screaming that name for weeks, I didn't know why!"},
							{speaker: "Tommy", text: "She's... she's real. She's across the ocean. I'm not crazy."},
							{speaker: "Pink Panther", text: "Four down, rockstar."},
						}, func() {
							game.inv.giveItemTo("Dance Card", "tommy")
							tommy.setStrange(false)
							game.vars.SetBool(ScopeGame, VarTommyHealed, true)
							tommy.dialog = []dialogEntry{
								{speaker: "Tommy", text: "You brought me my sister. I'll never forget that."},
							}
							game.travelMap.setUnlocked("rome_street", true)
							if dRoom, ok := game.sceneMgr.scenes["danny_room"]; ok {
								for _, n := range dRoom.npcs {
									if n.name == "Danny" {
										n.silent = false
										n.setStrange(true)
										break
									}
								}
							}
							game.dialog.queueDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "Danny's last. He keeps drawing Roman arches in the dirt."},
								{speaker: "Pink Panther", text: "Rome just unlocked. Let's finish this."},
							})
						}, &handOff{item: "Dance Card"}
				}
				break
			}
		}
	}

	// Danny in his cabin: silent until Tommy heals. Healed by Inscription Rubbing.
	if dannyRoom, ok := g.sceneMgr.scenes["danny_room"]; ok {
		for _, n := range dannyRoom.npcs {
			if n.name == "Danny" {
				danny := n
				danny.dialog = dannyStrangeDialog
				danny.onDialogEnd = func() {
					danny.dialog = dannyPostStrangeDialog
				}
				danny.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if !game.inv.hasItem("Inscription Rubbing") {
						return nil, nil, nil
					}
					return []dialogEntry{
							{speaker: "Pink Panther", text: "Danny. Look at the letters I copied from the Roman monument."},
							{speaker: "Danny", text: "That's... that's my NAME. My real name. In Latin."},
							{speaker: "Danny", text: "I've been drawing that arch for weeks because I was trying to draw my own name."},
							{speaker: "Pink Panther", text: "You found yourself in the inscription."},
							{speaker: "Danny", text: "Yeah. I did. All of us - we're spread across the world in pieces."},
							{speaker: "Pink Panther", text: "Five down. All home. Let's tell Higgins."},
						}, func() {
							game.inv.giveItemTo("Inscription Rubbing", "danny")
							danny.setStrange(false)
							game.vars.SetBool(ScopeGame, VarDannyHealed, true)
							danny.dialog = []dialogEntry{
								{speaker: "Danny", text: "I'm drawing something new now. All of us at camp, together."},
							}
							// Danny's heal unlocks Mexico City for the finale.
							game.travelMap.setUnlocked("mexico_street", true)
							game.dialog.queueDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "One pin left - Mexico City. That's where this ends."},
							})
						}, &handOff{item: "Inscription Rubbing"}
				}
				break
			}
		}
	}

	if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
		for _, n := range marcusRoom.npcs {
			if n.name == "Marcus" {
				marcus := n
				// 2026-06-24 (#20a): the postcard heal must be an intentional
				// hand-over - PP has to SELECT/hold the Postcard and click Marcus,
				// not have it fire just from owning it on scene entry. Gate the
				// alt dialog on the HELD item like the other trades.
				marcus.altDialogRequiresHeld = true
				marcus.altDialogRequiresItem = "Postcard"
				marcus.onDialogEnd = func() {
					// 2026-06-20 #19: never revert to the strange dialog once Marcus
					// is healed (this overwrite was leaving him on the pre-Paris
					// dialog after the postcard heal).
					if game.marcusHealed {
						return
					}
					if game.day == 2 && !game.talkedToMarcus {
						game.talkedToMarcus = true
						marcus.dialog = marcusPostStrangeDialog
					}
				}
				marcus.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if game.inv.heldItem == nil || game.inv.heldItem.name != "Postcard" {
						return nil, nil, nil
					}
					return []dialogEntry{
							{speaker: "Pink Panther", text: "Marcus, look at this postcard from the Louvre."},
							{speaker: "Marcus", text: "That's... that's the painting! The woman's face!"},
							{speaker: "Marcus", text: "And the glass pyramid... it's all real!"},
							{speaker: "Marcus", text: "I feel... calmer. The lines are fading..."},
							{speaker: "Marcus", text: "The whispers are stopping. It's like... waking up."},
							{speaker: "Pink Panther", text: "Rest now, Marcus. One kid down."},
						}, func() {
							game.inv.giveItemTo("Postcard", "marcus")
							marcus.setStrange(false)
							game.marcusHealed = true
							// #19: Marcus is calm now and finally drowsy - play the
							// go-to-sleep one-shot, then loop his sleeping idle so he's
							// actually asleep (matches Higgins's "sleeping soundly"
							// line). Both no-op gracefully until the art lands.
							marcus.playOneShotAnim("sleep", 1.4)
							// 2026-06-21 #23: the sleep MUST stick. The strange-alt
							// "freakout" punctuation (altIdleGrid/altIdleAfterSec) was
							// still firing after the heal and flashing him back to the
							// strange pose. Disable it and clear any in-flight alt cycle
							// so the sleeping idle persists.
							marcus.altIdleAfterSec = 0
							marcus.altIdleGrid = nil
							marcus.altIdleActive = false
							marcus.altIdleBackup = nil
							if len(marcus.sleepIdleGrid) > 0 {
								marcus.idleGrid = marcus.sleepIdleGrid
							}
							// 2026-06-24 (#21): keep the sleeping pose even when PP
							// talks to him again - don't let dialog swap him to the
							// talk sheet.
							marcus.lockIdleInDialog = true
							// Post-heal: a sleepy line, then he's out (the onDialogEnd
							// guard above keeps this from reverting to the strange dialog).
							marcus.dialog = []dialogEntry{
								{speaker: "Marcus", text: "*yawn* ...thanks, counselor. The lines are gone..."},
								{speaker: "Marcus", text: "I'm so tired... maybe tomorrow I'll feel better... zzz."},
							}
							// 2026-06-24 (#20b): do NOT brighten the room back to the
							// "day" BG when Marcus sleeps - the camp is in its darkened
							// mood by now, so the room must stay in that same BG.
							game.travelMap.setUnlocked("jerusalem_entrance", true)
							// User playtest #39: update office Higgins's dialog the
							// moment Marcus heals, so a return office visit shows the
							// Jake bridge (not the stale pre-Marcus "worried" line,
							// which only swapped after talking to him once).
							if office, ok := game.sceneMgr.scenes["camp_office"]; ok {
								for _, hn := range office.npcs {
									if hn.name == "Director Higgins" {
										hn.dialog = higginsPostMarcusHealedDialog
										break
									}
								}
							}
							// Wake up Jake so the player can talk to him in his cabin
							if jRoom, ok := game.sceneMgr.scenes["jake_room"]; ok {
								for _, n := range jRoom.npcs {
									if n.name == "Jake" {
										n.silent = false
										n.setStrange(true)
										break
									}
								}
							}
							game.dialog.queueDialog([]dialogEntry{
								{speaker: "Pink Panther", text: "One down. Jake's next - he keeps muttering about tunnels and a wall."},
								{speaker: "Pink Panther", text: "The travel map just lit up Jerusalem. That can't be a coincidence."},
							})
							// #22: Marcus visibly takes/looks at the postcard during the
							// hand-over (npcAnim); no-ops until the art lands.
						}, &handOff{item: "Postcard", npcAnim: "receive_postcard"}
				}
				break
			}
		}
	}

	// --- Day 2: Office Higgins gives map and unlocks Paris ---
	if office, ok := g.sceneMgr.scenes["camp_office"]; ok {
		for _, n := range office.npcs {
			if n.name == "Director Higgins" {
				officeHiggins := n
				officeHiggins.onDialogEnd = func() {
					if !game.parisUnlocked {
						game.parisUnlocked = true
						game.travelMap.setUnlocked("paris_street", true)
						// User 2026-04-26: replace the silent giveMapItem
						// with the higgins_give_map sequence - Higgins
						// plays his give-map anim, PP plays receive_map,
						// then the item drops into inventory. No more
						// inventory-bar pop on map handover.
						if seq := game.seqStore.Get("higgins_give_map"); seq != nil {
							game.seqPlayer.Play(seq)
						} else {
							game.giveMapItem()
						}
					}
					if game.marcusHealed {
						officeHiggins.dialog = higginsPostMarcusHealedDialog
					} else {
						officeHiggins.dialog = higginsPostWorriedDialog
					}
				}
				break
			}
		}
	}

	// --- Day 1: Kids normal intros, count meetings ---
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			kid := n
			switch kid.name {
			case "Marcus":
				// metInDay1 latches once so re-clicking Marcus after the post
				// dialog swap doesn't bump metKids twice (user 2026-04-26:
				// "i spoke with the kids and then the night scene just
				// started" - the unguarded Day 1 branch let metKids hit 5
				// from one or two kids alone).
				marcusMet := false
				kid.onDialogEnd = func() {
					if game.day == 1 && !marcusMet {
						marcusMet = true
						game.metKids++
						kid.dialog = marcusPostDialog
						game.checkDay1Complete()
					} else if game.day >= 2 {
						if !game.talkedToMarcus {
							game.talkedToMarcus = true
						}
						kid.dialog = marcusPostStrangeDialog
					}
				}
			case "Tommy":
				tommyMet := false
				tommyExited := false
				kid.onDialogEnd = func() {
					if game.day == 1 && !tommyMet {
						tommyMet = true
						game.metKids++
						kid.dialog = tommyPostDialog
						game.checkDay1Complete()
						// User 2026-05-21: after the Day-1 intro, Tommy walks
						// off-left and exits the camp grounds. Plays the
						// tommy_exit sequence if available; otherwise just
						// hides him (so re-entering the scene doesn't re-spawn
						// him in the same spot).
						if !tommyExited {
							tommyExited = true
							if seq := game.seqStore.Get("tommy_exit"); seq != nil {
								game.seqPlayer.Play(seq)
							} else {
								kid.hidden = true
							}
						}
					} else if game.day >= 2 && !kid.dialogDone {
						kid.dialogDone = true
						kid.dialog = tommyPostStrangeDialog
					}
				}
			case "Jake":
				jakeMet := false
				jakeExited := false
				kid.onDialogEnd = func() {
					if game.day == 1 && !jakeMet {
						jakeMet = true
						game.metKids++
						kid.dialog = jakePostDialog
						game.checkDay1Complete()
						// User 2026-05-21: after the Day-1 intro, Jake walks
						// back into his cabin and disappears.
						if !jakeExited {
							jakeExited = true
							if seq := game.seqStore.Get("jake_exit"); seq != nil {
								game.seqPlayer.Play(seq)
							} else {
								kid.hidden = true
							}
						}
					} else if game.day >= 2 && !kid.dialogDone {
						kid.dialogDone = true
						kid.dialog = jakePostStrangeDialog
					}
				}
			case "Lily":
				// hintState progression (see npc.hintState):
				//   0 -> 1: shy dialog done, arm the flower alt-dialog
				//   1 -> 2: flower handed over, swap to post-dialog
				// Day 2+ uses the flat dialogDone flag and lilyPostStrangeDialog.
				kid.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if kid.hintState != 1 || !game.inv.hasItem("Flower") {
						return nil, nil, nil
					}
					// NOTE: altDialogFunc must stay SIDE-EFFECT FREE before the
					// returned callback - the cursor hover probe (ui.go
					// itemMatch) calls it every frame to test for a match.
					// Her receive-flower one-shot used to fire HERE, which made
					// Lily react the moment PP merely HOVERED her with the
					// flower (2026-06-12 #5).
					// PR#1: the hand-off beat (PP's give + her receive_flower)
					// moved to the handOff return - it plays BEFORE the text.
					return lilyFlowerDialog, func() {
						game.inv.giveItemTo("Flower", "lily")
						game.metKids++
						kid.hintState = 2
						kid.dialog = lilyDialog
						// #4: from now on she talks holding the daisy.
						if len(kid.postGiveTalkGrid) > 0 {
							kid.talkGrid = kid.postGiveTalkGrid
						}
						kid.altDialogFunc = nil
						kid.altDialogRequiresHeld = false
						kid.altDialogRequiresItem = ""
						game.checkDay1Complete()
					}, &handOff{item: "Flower", npcAnim: "receive_flower", npcAnimDur: 1.4}
				}
				kid.onDialogEnd = func() {
					if game.day == 1 && kid.hintState == 0 {
						kid.hintState = 1
						// User 2026-05-23: strict HELD gate. The flower beat
						// only fires when PP has actively pulled the flower
						// out of inventory (it's heldItem on the cursor) and
						// clicked Lily. Just having the flower in the bag is
						// NOT enough - user wants the give-item motion
						// required. Same pattern will apply to Marcus
						// postcard, Jake coin, etc.
						kid.altDialogRequiresHeld = true
						kid.altDialogRequiresItem = "Flower"
						kid.altDialogStrictMissingHint = []dialogEntry{
							{speaker: "Lily", text: "..."},
							{speaker: "Pink Panther", text: "She won't look up. If I had something for her, I'd need to bring it out and offer it."},
						}
						// Higgins walks in from offscreen-right to deliver the
						// flower hint. The sequence teleports him out of view,
						// un-hides, lerps to his hint spot, then plays dialog.
						// Replaces the old "teleport in + queueDialog" code.
						if seq := game.seqStore.Get("higgins_walk_in"); seq != nil {
							game.seqPlayer.Play(seq)
						} else {
							// Fallback if the sequence file is missing so the
							// story isn't blocked: reveal Higgins at his hint
							// position and queue the dialog directly.
							if grounds, ok := game.sceneMgr.scenes["camp_grounds"]; ok {
								for _, n := range grounds.npcs {
									if n.name == "Director Higgins" {
										n.hidden = false
										n.silent = false
										break
									}
								}
							}
							game.dialog.queueDialog(higginsLilyHintDialog)
						}
					} else if game.day >= 2 && !kid.dialogDone {
						kid.dialogDone = true
						kid.dialog = lilyPostStrangeDialog
					}
				}
			case "Danny":
				dannyMet := false
				kid.onDialogEnd = func() {
					if game.day == 1 && !dannyMet {
						dannyMet = true
						game.metKids++
						kid.dialog = dannyPostDialog
						game.checkDay1Complete()
					} else if game.day >= 2 && !kid.dialogDone {
						kid.dialogDone = true
						kid.dialog = dannyPostStrangeDialog
					}
				}
			}
		}
	}

	// --- Lake: Flower pickup for Lily ---
	if lake, ok := g.sceneMgr.scenes["camp_lake"]; ok {
		flowerDef, _ := g.items.getDef("flower")
		flowerTex, flowerW, flowerH := engine.SafeTextureFromPNGKeyed(g.renderer, flowerDef.Texture)
		flower := &floorItem{
			tex:  flowerTex,
			srcW: flowerW,
			srcH: flowerH,
			// User 2026-05-19: bounds widened from 50×50 → 100×100 so the
			// pickup feels findable. Cursor change (in updateHover) +
			// click hit-test both use these bounds, so the icon and
			// action align across the full 100×100 patch.
			bounds:  sdl.Rect{X: 150, Y: 440, W: 100, H: 100},
			name:    "Flower",
			visible: true,
			// 2026-06-15 (#1 follow-up): the grab_flower sheet has PP leaning
			// LEFT toward a daisy on his left. walkToFloorItem stands PP to the
			// item's LEFT by default, which puts the flower on his RIGHT - so he
			// bent AWAY from it ("looks at the other side"). Stand to the flower's
			// RIGHT so the daisy is on his left, matching the animation.
			standRight: true,
			onPickup: func() {
				// Hide flower in scene first so the grab anim doesn't play
				// over a still-visible daisy on the ground.
				if lake, ok := game.sceneMgr.scenes["camp_lake"]; ok {
					for _, fi := range lake.floorItems {
						if fi.name == "Flower" {
							fi.visible = false
							break
						}
					}
				}
				// Play the grab one-shot before the inventory pulse so the
				// player visibly bends, picks, and rises. Falls through
				// instantly if the asset isn't registered.
				game.player.playOneShot("grab_flower", 0.9, func() {
					item := game.items.createItem("flower")
					if item != nil {
						game.inv.addItem(item)
					}
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "A pretty daisy. I bet Lily would like this."},
					})
				})
			},
		}
		lake.floorItems = append(lake.floorItems, flower)

		// #10: clicking the water gives a flavour line instead of nothing. PP
		// walks to the deck edge (the hotspot centre, snapped to the deck path)
		// then declines to swim. Bounds cover the lake surface above the deck;
		// tune in-game if the water sits elsewhere in the art.
		lake.hotspots = append(lake.hotspots, hotspot{
			bounds: sdl.Rect{X: 360, Y: 250, W: 700, H: 210},
			name:   "Lake",
			onInteract: func() bool {
				game.dialog.startDialog([]dialogEntry{
					{speaker: "Pink Panther", text: "Cool, clear water... but I'm not really in the mood for a swim right now."},
				})
				return true
			},
		})

		// Japan chapter: sad Lily at the dock (revealed after Jake's heal). The
		// discovery beat plays first; later, holding the Pressed Sakura from
		// Kyoto heals her right here at the lake.
		for _, n := range lake.npcs {
			if n.name != "Lily" {
				continue
			}
			lily := n
			lily.dialog = []dialogEntry{
				{speaker: "Pink Panther", text: "Lily? It's me. ...You've been out here a long time."},
				{speaker: "Lily", text: "The flowers are all wrong. The colours won't sit still. There's a tower of orange gates, and bells I can't quite hear."},
				{speaker: "Lily", text: "I just want it to be quiet and pretty again."},
				{speaker: "Pink Panther", text: "Orange gates... bells... I think I know where that is. Hang on, Lily."},
			}
			lily.onDialogEnd = func() {
				game.vars.SetBool(ScopeGame, VarLilyLakeMet, true)
				lily.dialog = []dialogEntry{
					{speaker: "Lily", text: "..."},
					{speaker: "Pink Panther", text: "Stay here, Lily. I'll be back."},
				}
			}
			lily.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
				if !game.inv.hasItem("Pressed Sakura") {
					return nil, nil, nil
				}
				return []dialogEntry{
						{speaker: "Pink Panther", text: "Lily. I went to Kyoto - the place with the orange gates."},
						{speaker: "Lily", text: "You... you saw it? It was real?"},
						{speaker: "Pink Panther", text: "Here. A real cherry blossom, pressed flat. Hold it."},
						{speaker: "Lily", text: "...it's so light. And it's the right pink. The exact right pink."},
						{speaker: "Lily", text: "The wind stopped whispering. I can hear myself again."},
						{speaker: "Pink Panther", text: "Three down."},
					}, func() {
						game.inv.giveItemTo("Pressed Sakura", "lily")
						lily.setStrange(false)
						lily.lockIdleInDialog = false
						game.vars.SetBool(ScopeGame, VarLilyHealed, true)
						lily.dialog = []dialogEntry{
							{speaker: "Lily", text: "Thank you for the flower. I'll keep it forever."},
						}
						game.travelMap.setUnlocked("rio_street", true)
						game.travelMap.setUnlocked("buenos_aires_street", true)
						if tRoom, ok := game.sceneMgr.scenes["tommy_room"]; ok {
							for _, tn := range tRoom.npcs {
								if tn.name == "Tommy" {
									tn.silent = false
									tn.setStrange(true)
									break
								}
							}
						}
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "Tommy next. He keeps yelling about a sister he doesn't have."},
							{speaker: "Pink Panther", text: "Rio AND Buenos Aires both lit up. This one might take two stops."},
						})
					}, &handOff{item: "Pressed Sakura"}
			}
			break
		}
	}

}

// playHigginsRudeBeat (Japan opening #8): once PP has found Lily at the lake,
// Higgins "comes half the way" down the grounds path, brushes off the worry,
// and leaves. PP walks up to meet him; the curt exchange plays; then the camera
// aside connects flowers + Japan and lights up Tokyo. One-shot (higginsRudeStarted).
func (g *Game) playHigginsRudeBeat() {
	grounds, ok := g.sceneMgr.scenes["camp_grounds"]
	if !ok {
		g.finishHigginsRude()
		return
	}
	var higgins *npc
	for _, n := range grounds.npcs {
		if n.name == "Director Higgins" {
			higgins = n
			break
		}
	}
	if higgins == nil {
		g.finishHigginsRude()
		return
	}
	// He comes OUT of his office (the bottom-right corner, where the office
	// hotspot is) and STRIDES across the grounds toward PP (front-facing walk
	// loop), then stops and delivers the curt line.
	higgins.hidden = false
	higgins.silent = false
	higgins.bounds.X = 1180
	higgins.bounds.Y = 540
	higgins.dialog = higginsRudeDialog
	higgins.swapIdleForOneShot("walk_front") // loop the walk while he approaches
	// PP comes up a little to meet him.
	g.player.walkToAndDo(620, float64(higgins.bounds.Y+higgins.bounds.H)-float64(playerDstH)/2, nil)
	g.higginsWalk = &higginsWalkState{
		n: higgins, fromX: 1180, toX: 800, dur: 1.6,
		onArrive: func() {
			higgins.restoreSwappedIdle()
			g.dialog.startDialogWithCallback(higginsRudeDialog, func() {
				g.finishHigginsRude()
			})
		},
	}
}

// higginsWalkState is the simple x-lerp that strides Higgins halfway down the
// camp-grounds path for the rude intercept (driven from Update).
type higginsWalkState struct {
	n        *npc
	fromX    float64
	toX      float64
	elapsed  float64
	dur      float64
	onArrive func()
}

// finishHigginsRude plays PP's camera aside, hides Higgins, and unlocks Tokyo.
func (g *Game) finishHigginsRude() {
	g.player.dir = dirDown
	g.player.facingLeft = false
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			if n.name == "Director Higgins" {
				n.hidden = true
				n.silent = true
				break
			}
		}
	}
	g.vars.SetBool(ScopeGame, VarHigginsRudeDone, true)
	g.travelMap.setUnlocked("tokyo_torii", true)
	g.vars.SetBool(ScopeGame, VarTokyoUnlocked, true)
	g.dialog.startDialog(higginsRudeAsideDialog)
}

// checkDay1Complete triggers the Day 1 -> Night transition once PP has met
// all 5 kids. Higgins delivers a short "it's getting late" beat on
// camp_grounds first so the fade to night isn't abrupt (user complaint
// 2026-04-17: "higgins didnt say it time to sleep" - the previous version
// cut straight to the campfire and the Lily-flower handoff felt swallowed
// by the transition).
func (g *Game) checkDay1Complete() {
	if g.metKids < 5 || g.day != 1 || g.sceneMgr.transitioning || g.day1BedtimeStarted {
		return
	}
	// Belt-and-braces (user 2026-04-26): require Lily's flower handoff to
	// have completed before night triggers, even if metKids somehow hits 5
	// without it. hintState reaches 2 only inside Lily's altDialogFunc
	// callback (see setupCampCallbacks "Lily" case).
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			if n.name == "Lily" && n.hintState < 2 {
				return
			}
		}
	}
	// Latch so the handoff can't re-fire the bedtime sequence.
	g.day1BedtimeStarted = true
	game := g
	// User 2026-06-12 (#6): this is the bellow the player actually watches -
	// swap BOTH idle and talk to the shout frames for the dialog (the dialog
	// system puts the speaker into TALK state, which would override a swapped
	// idle alone; that's why the camp_night ctor sets idle=talk=shout too).
	// Restored in the callback right before the night transition.
	var hg *npc
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			if n.name == "Director Higgins" {
				hg = n
				break
			}
		}
	}
	var savedTalk []npcFrame
	if hg != nil && len(hg.oneShotAnims["shout"]) > 0 {
		hg.swapIdleForOneShot("shout")
		savedTalk = hg.talkGrid
		hg.talkGrid = hg.oneShotAnims["shout"]
	}
	g.dialog.startDialogWithCallback([]dialogEntry{
		{speaker: "Director Higgins", text: "Ahem! It's getting very late, counselor."},
		{speaker: "Director Higgins", text: "All campers to their cabins. NOW."},
	}, func() {
		if hg != nil {
			if savedTalk != nil {
				hg.talkGrid = savedTalk
			}
			hg.restoreSwappedIdle()
		}
		game.sceneMgr.transitionTo("camp_night", game.player)
	})
}

// findNightHiggins returns the silent night-campfire Higgins NPC, if present.
// Kept because other setup code may still reference him for hidden-NPC tweaks;
// the night cutscene itself now runs through assets/data/sequences/night_bedtime.json.
func (g *Game) findNightHiggins() *npc {
	scene, ok := g.sceneMgr.scenes["camp_night"]
	if !ok {
		return nil
	}
	for _, n := range scene.npcs {
		if n.name == "Director Higgins" {
			return n
		}
	}
	return nil
}

func (g *Game) giveMapItem() {
	if g.inv.hasItem("Travel Map") {
		return
	}
	item := g.items.createItem("travel_map")
	if item == nil {
		return
	}

	// User request 2026-04-17: no more big "map grows on screen" reveal.
	// The PP-takes-map gesture already plays in the give-map animation;
	// just drop the item into inventory silently.
	g.inv.addItem(item)
}

// genericPickupLines rotate through PP's pocket-quips. SKILL.md §8c (user
// 2026-06-11 #18): pickup lines stay GENERIC - "who needs this" hints belong
// in NPC dialogs, not in the pocket beat, so items stay reusable across
// quests and the player keeps discovering uses by talking to people.
var genericPickupLines = []string{
	"This might come in handy.",
	"Straight into the pocket.",
	"A panther never leaves a clue behind.",
	"You never know when you'll need one of these.",
}
var genericPickupIdx int

// genericPickupDialog builds the standard one-line pickup beat: the caller's
// flavor description of WHAT was found + a rotating generic quip.
func genericPickupDialog(flavor string) []dialogEntry {
	line := genericPickupLines[genericPickupIdx%len(genericPickupLines)]
	genericPickupIdx++
	return []dialogEntry{
		{speaker: "Pink Panther", text: flavor + " " + line},
	}
}

func (g *Game) setupParisCallbacks() {
	// --- Paris side-quest state (2026-06-10) ---
	// Shared across the street / louvre / bakery blocks below. Closure state
	// follows the existing pattern (souvenirArmed, pierre.hintState); none of
	// it needs to survive a save/load any more than those do.
	// PR (2026-06-12) Paris chain reorder — strictly linear, softlock-proof:
	//   rolling pin → Poulain (baguette+coffee) → Henri (confiture) →
	//   Pierre (baguette+confiture → PRESS PASS) → Claude (press pass →
	//   louvreUnlocked) → Beaumont (asks for Camille's sketch) → Camille
	//   (needs her pencil; the flower-pot pigeon guards it) → Poulain gives
	//   the day-old Baguette Heel → Pierre shoos the pot pigeon with it
	//   (pigeonsCleared) → pot reveals the pencil → Camille sketches →
	//   Beaumont trades the Postcard. The old easel "pigeon critic" +
	//   mini_portrait beat is removed; the heel's job is shooing the pot
	//   pigeon. Closure state (none needs save/load):
	var (
		sketchAsked    bool // Beaumont asked for Camille's replica sketch (postcards sold out)
		camilleAsked   bool // Camille sent PP after her lost lucky pencil
		louvreUnlocked bool // Claude took the press pass and waved PP in
		pigeonsCleared bool // Pierre shooed the flower-pot pigeon with the heel
		pencilTaken    bool // pencil fished out of the flower pot by the Louvre steps
		sketchDone     bool // Camille drew the Room 7 replica, sketch in PP's bag
		souvenirAsked  bool // Poulain asked for the grandson postcard (post-heal)
		souvenirDone   bool // signed postcard delivered to Poulain
	)
	// Bakery NPC handles needed by Poulain's counter-service branch below.
	var bakeryHenri *npc
	if bakery, ok := g.sceneMgr.scenes["paris_bakery"]; ok {
		for _, n := range bakery.npcs {
			if n.name == "Monsieur Henri" {
				bakeryHenri = n
				break
			}
		}
	}

	// French Guide + the two locals: all swap to post-dialog after first chat.
	if parisStreet, ok := g.sceneMgr.scenes["paris_street"]; ok {
		game := g
		// clearPotPigeon (2026-06-12): the pigeon lady lures the flower-pot
		// guard pigeon off with the heel. Swaps the pot prop to the exposed
		// pencil and flaps the pigeon up-and-away. Shared so the pencil
		// pickup gate (pigeonsCleared) and the lady's hand-off agree.
		clearPotPigeon := func() {
			pigeonsCleared = true
			for _, fi := range parisStreet.floorItems {
				if fi.name == "Charcoal Pencil" {
					tex, w, h := engine.SafeTextureFromPNGKeyed(game.renderer, "assets/images/locations/paris/props/flower_pot_pencil.png")
					if tex != nil {
						fi.tex = tex
						fi.srcW = w
						fi.srcH = h
					}
					break
				}
			}
			// 2026-06-20 #8: fly-up follows the pot to its new spot beside Pierre
			// (pot top-centre ~909,565).
			if amb := newAmbientPigeonFlyUp(game.renderer, 909, 565); amb != nil {
				parisStreet.ambientSprites = append(parisStreet.ambientSprites, amb)
			}
		}
		for _, n := range parisStreet.npcs {
			switch n.name {
			case "Madame Colette":
				guide := n
				guide.onDialogEnd = func() {
					guide.dialog = frenchGuidePostDialog
				}
			case "Pierre":
				// Quest step 2 (user 2026-05-21): rewritten as TWO STAGES so PP
				// brings the baguette and the spread on two separate trips.
				// Stage 1 (hintState 0): holds Baguette → Pierre takes it,
				//   says "it's dry, bring me a spread", swaps dialog. No
				//   press-pass yet. Inventory loses the baguette.
				// Stage 2 (hintState 1): holds Confiture → Pierre takes it,
				//   hands over the Press Pass.
				// PP cannot present both items at once (one of them is no
				// longer in the inventory at stage 2).
				pierre := n
				// 2026-06-20 #8: PP must look SMALLER (mid-distance) at Pierre, the
				// back-of-line painter. depthScale can't do it (PP's y is clamped to
				// the same street line everywhere), so Pierre keeps a recede
				// choreography - but the CLEAN version: walk to the line + smooth
				// shrink, and if PP is ALREADY receded here (stage-2 click) talk in
				// place instead of re-walking (the recedeHeld guard, which fixed the
				// old "jump back to the road then trudge in" #11). Margaux no longer
				// recedes, so there's no inter-NPC size pop. The held-item drop path
				// in HandleClick routes through this override too.
				pierre.onClickOverride = func() {
					if game.player == nil || game.dialog == nil {
						return
					}
					talk := func() {
						game.player.state = stateTalking
						// Talk to the SIDE (Pierre's easel is to PP's right).
						game.player.facingLeft = false
						game.player.dir = dirRight
						releaseRecede := func() {
							game.player.state = stateIdle
							// Stay shrunk at Pierre's depth until PP next moves
							// (setTarget grows him back smoothly).
							game.player.holdRecede()
						}
						if pierre.altDialogFunc != nil {
							entries, cb, ho := pierre.altDialogFunc()
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
									game.player.playHandOff(pierre, ho, start)
								} else {
									start()
								}
								return
							}
						}
						game.dialog.startDialogWithCallback(pierre.dialog, func() {
							if pierre.hintState == 0 {
								pierre.dialog = pierreArtistPostDialog
							}
							releaseRecede()
						})
					}
					// Already lined up + receded from a prior click -> talk in place
					// (no jump back), the #11 guard.
					if game.player.recedeHeld {
						talk()
						return
					}
					game.player.walkToAndDo(690, 510, func() {
						game.player.playRecede(1.0, 0.65, 50, talk)
					})
				}
				pierre.onDialogEnd = func() {
					// Don't overwrite if we already pivoted to the stage-1
					// "I need a spread" dialog. Pierre's `dialog` is updated
					// in-stage below.
					if pierre.hintState == 0 {
						pierre.dialog = pierreArtistPostDialog
					}
				}
				pierre.altDialogRequiresHeld = false
				pierre.altDialogRequiresItem = "Baguette"
				pierre.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					// User playtest #25: the trade only fires when PP actively
					// HANDS OVER the item (it's on the cursor), never just from
					// having it in the bag. heldMatches checks the carried item.
					heldMatches := func(name string) bool {
						return game.inv.heldItem != nil && game.inv.heldItem.name == name
					}
					// Stage 1: PP brings the baguette for the first time.
					if pierre.hintState == 0 && heldMatches("Baguette") {
						return []dialogEntry{
								{speaker: "Pierre", text: "Mon Dieu! Is zat a fresh baguette from Madame Poulain?"},
								{speaker: "Pink Panther", text: "It can be yours, Pierre. I need a favor."},
								{speaker: "Pierre", text: "Bread is good, mon ami, but it is so... dry."},
								{speaker: "Pierre", text: "Bring me a spread - beurre or confiture - and we can talk press passes."},
							}, func() {
								game.inv.giveItemTo("Baguette", "pierre")
								pierre.hintState = 1
								pierre.dialog = []dialogEntry{
									{speaker: "Pierre", text: "Still waiting on zat spread, mon ami. Beurre or confiture, anything."},
								}
								pierre.altDialogRequiresItem = "Confiture"
							}, &handOff{item: "Baguette", npcAnim: "receive_baguette"}
					}
					// Stage 2: PP brings the spread.
					if pierre.hintState == 1 && heldMatches("Confiture") {
						return []dialogEntry{
								{speaker: "Pierre", text: "Magnifique! Strawberries from ze south - perfect with my baguette."},
								{speaker: "Pierre", text: "Here, ze press pass - Claude owes me a favor anyway."},
								{speaker: "Pink Panther", text: "Bon appetit, Pierre."},
							}, func() {
								// §8b: the confiture hand-over plays pre-dialog (PR#1).
								// Pierre's press pass is now handed back via the
								// two-stage handOff below (#16), not a parallel receive.
								game.inv.giveItemTo("Confiture", "pierre")
								if item := game.items.createItem("press_pass"); item != nil {
									game.inv.addItem(item)
								}
								pierre.hintState = 2
								// Reorder (2026-06-12): Pierre is done questing after
								// the press pass - he just points PP at the Louvre. The
								// flower-pot pigeon is now handled by Madame Margaux,
								// the pigeon lady across the street (PP brings HER the
								// heel), so Pierre no longer takes it.
								pierre.dialog = pierreArtistPostDialog
								pierre.altDialogFunc = nil
								pierre.altDialogRequiresItem = ""
							}, &handOff{item: "Confiture", npcAnim: "receive_confiture", returnItem: "Press Pass", npcGiveAnim: "give_ticket"}
					}
					return nil, nil, nil
				}
			case "Madame Margaux":
				// PR (2026-06-12): the pigeon lady lures the flower-pot guard
				// pigeon off when PP brings her the day-old Baguette Heel.
				// Gated on camilleAsked (so it only matters once Camille has
				// sent PP after the pencil) and !pigeonsCleared.
				margaux := n
				margaux.onDialogEnd = func() { margaux.dialog = pigeonLadyPostDialog }
				margaux.altDialogRequiresHeld = true
				margaux.altDialogRequiresItem = "Baguette Heel"
				margaux.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					held := game.inv.heldItem
					if !camilleAsked || pigeonsCleared ||
						held == nil || held.name != "Baguette Heel" {
						return nil, nil, nil
					}
					return pigeonLadyHeelDialog, func() {
						// §8b: the heel hand-over plays pre-dialog (PR#1); she
						// scatters it and the pot pigeon flutters off.
						game.inv.giveItemTo("Baguette Heel", "pigeon_lady")
						clearPotPigeon()
						margaux.dialog = pigeonLadyPostDialog
						margaux.altDialogFunc = nil
						margaux.altDialogRequiresHeld = false
						margaux.altDialogRequiresItem = ""
					}, &handOff{item: "Baguette Heel"}
				}
				// 2026-06-24 (#9): restore Margaux's recede (mirror of Pierre's) so
				// PP renders at the SAME mid-distance SIZE at her as at Pierre - the
				// narrowed depthScale (0.95-1.05) can't shrink him enough on its own,
				// so she gets the walk-to-line + smooth-shrink choreography too.
				margaux.onClickOverride = func() {
					if game.player == nil || game.dialog == nil {
						return
					}
					talk := func() {
						game.player.state = stateTalking
						game.player.facingLeft = false
						game.player.dir = dirRight
						releaseRecede := func() {
							game.player.state = stateIdle
							game.player.holdRecede()
						}
						if margaux.altDialogFunc != nil {
							entries, cb, ho := margaux.altDialogFunc()
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
									game.player.playHandOff(margaux, ho, start)
								} else {
									start()
								}
								return
							}
						}
						game.dialog.startDialogWithCallback(margaux.dialog, func() {
							releaseRecede()
						})
					}
					if game.player.recedeHeld {
						talk()
						return
					}
					game.player.walkToAndDo(560, 510, func() {
						game.player.playRecede(1.0, 0.65, 50, talk)
					})
				}
			case "Nicolas":
				// "Camille and the Sold-Out Postcard" street hop: once
				// Camille asks about her lost pencil, Nicolas's lens knows
				// exactly where it rolled.
				nicolas := n
				nicolas.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if camilleAsked && !pencilTaken {
						return nicolasPencilHintDialog, nil, nil
					}
					return nil, nil, nil
				}
			case "Claude":
				// PR#24 (2026-06-12): the press pass was a silent key (the gate
				// just checked ownership), so it felt stuck in the bag. Now PP
				// HANDS it to Claude, who waves him in - that sets louvreUnlocked
				// and consumes the pass. The Louvre hotspot checks the flag.
				claude := n
				claude.onDialogEnd = func() {
					claude.dialog = gendarmePostDialog
				}
				claude.altDialogRequiresHeld = true
				claude.altDialogRequiresItem = "Press Pass"
				claude.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					if louvreUnlocked || !(game.inv.heldItem != nil && game.inv.heldItem.name == "Press Pass") {
						return nil, nil, nil
					}
					return claudePressPassDialog, func() {
						game.inv.giveItemTo("Press Pass", "claude") // consumed at the door
						louvreUnlocked = true
						claude.dialog = gendarmePostDialog
						claude.altDialogFunc = nil
						claude.altDialogRequiresHeld = false
						claude.altDialogRequiresItem = ""
					}, &handOff{item: "Press Pass"}
				}
			}
		}

		// Louvre entrance gate: needs the Press Pass (#37 - single credential).
		for i := range parisStreet.hotspots {
			if parisStreet.hotspots[i].name != "To the Louvre" {
				continue
			}
			h := &parisStreet.hotspots[i]
			h.onInteract = func() bool {
				// PR#24: opens only after PP HANDS Claude the pass (louvreUnlocked).
				if !louvreUnlocked {
					if game.inv.hasItem("Press Pass") {
						game.dialog.startDialog([]dialogEntry{
							{speaker: "Gendarme", text: "A press pass, oui? Bring it HERE, monsieur - hand it to me and I wave you straight in."},
						})
						return true
					}
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Gendarme", text: "Monsieur, ze museum is press and pass-holders only today."},
						{speaker: "Gendarme", text: "Find a press pass - Pierre ze painter knows everyone."},
					})
					return true
				}
				// 2026-06-11 #34: single walk-in path - PP enters from the
				// left tunnel via entryWalkPending (no double-spawn flicker).
				game.sceneMgr.entryWalkPending = true
				game.sceneMgr.transitionTo("paris_louvre", game.player)
				return true
			}
			break
		}
	}

	// Museum Curator: gives postcard on first dialog, then post-dialog
	if parisLouvre, ok := g.sceneMgr.scenes["paris_louvre"]; ok {
		game := g
		for _, n := range parisLouvre.npcs {
			if n.name == "Curator Beaumont" {
				curator := n
				postcardGiven := false
				curator.onDialogEnd = func() {
					// 2026-06-10 rework: the postcards are SOLD OUT, so the
					// first chat no longer hands the postcard over - Beaumont
					// asks for Camille's replica sketch and the postcard moves
					// to the sketch trade below. The flag flip also guards
					// against the old duplicate-postcard repeat-chat bug.
					if !sketchAsked {
						sketchAsked = true
						curator.dialog = curatorWaitingDialog
					}
				}
				curator.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
					// Sketch → postcard trade (main chain resumes here).
					if held := game.inv.heldItem; held != nil && held.name == "Camille's Sketch" && !postcardGiven {
						return curatorSketchTradeDialog, func() {
							postcardGiven = true
							game.inv.giveItemTo("Camille's Sketch", "curator_beaumont")
							item := game.items.createItem("postcard")
							if item != nil {
								game.inv.addItem(item)
							}
							curator.dialog = museumCuratorPostDialog
							// User playtest #32: getting the postcard is the Paris
							// goal. Mark Paris done and unlock the "fly back to
							// Camp" travel pin so the player can return to heal
							// Marcus.
							game.vars.SetBool(ScopeGame, VarParisDone, true)
							game.travelMap.setUnlocked("camp_entrance", true)
							// #34: the camp turns "wrong" the moment the France
							// trip is behind us - darken the grounds bg so the
							// return landing reads as ominous.
							game.applyCampMood()
							// §PR3: the sketch hand-over plays pre-dialog (PR#1).
							// After Beaumont's lines he visibly hands the postcard
							// back, PP takes it, then the homeward monologue plays.
							curator.playOneShotAnimThen(curator.giveAnimOr("give_postcard"), 1.0, func() {
								game.player.playReceive("postcard", false, 1.0, func() {
									game.dialog.startDialog([]dialogEntry{
										{speaker: "Pink Panther", text: "A postcard of the painting... this is what Marcus has been drawing."},
										{speaker: "Pink Panther", text: "I should head back outside, then take the travel map home to camp. Marcus needs this."},
									})
								})
							})
						}, &handOff{item: "Camille's Sketch"}
					}
					// Grandson souvenir loop: once Poulain has asked, Beaumont
					// signs a second postcard (the new prints have arrived).
					if souvenirAsked && !souvenirDone && !game.inv.hasItem("Signed Postcard") {
						return curatorSouvenirDialog, func() {
							// §PR3: reuse Beaumont's postcard hand-over for the signed card too.
							curator.playOneShotAnimThen(curator.giveAnimOr("give_postcard"), 1.0, func() {
								game.player.playReceive("postcard", false, 1.0, func() {
									if item := game.items.createItem("postcard_grandson"); item != nil {
										game.inv.addItem(item)
									}
								})
							})
						}, nil
					}
					return nil, nil, nil
				}
				break
			}
		}

		// User playtest #31: removed the bottom-right "Travel Map" hotspot in
		// the museum. The travel map is opened from the Travel Map inventory
		// item (g.inv.onSelectItem), so the on-screen button was redundant.
	}

	// User 2026-04-26: removed the paris_street left-arrow that opened the
	// travel map. The map now opens by clicking the Travel Map item in the
	// inventory (see inventory.handleClick). The left-arrow on paris_street
	// is repurposed in paris_street.json to open paris_bakery.

	// --- Paris Bakery: Madame Poulain's rolling-pin quest ---
	// Retro-style intro beat (user 2026-04-26): the baguette is no longer a
	// freebie. PP has to find Madame Poulain's lost rolling pin on the bakery
	// floor, hand it over, and only then does she trade the baguette.
	if bakery, ok := g.sceneMgr.scenes["paris_bakery"]; ok {
		// Playtest 2026-06-15 (#7/#10): PP "disappeared" after talking to ANY
		// bakery NPC. Root cause: PR#12 moved PP's stand row up into the aisle
		// (foot ~470-480), which lands him squarely in the seated patrons' bust
		// band (~y352-487). The old #27 hack forced every patron to draw IN
		// FRONT of PP (drawFootY=900), so a patron bust overlapping PP's spot
		// swallowed him whole. PP's foot in this scene is always >=470 (minY
		// 200) and always ABOVE the tablecloth line (~y536+), so drawing PP in
		// FRONT of the seated patrons keeps him visible without putting him "on
		// the cloths." Pin the patrons' sort-foot BELOW PP's minimum foot so the
		// roaming protagonist always renders on top of the seated regulars.
		// (Poulain keeps her natural low foot-Y → she renders BEHIND PP, behind
		// the counter, which already read correctly.)
		for _, bn := range bakery.npcs {
			if bn.name != "Madame Poulain" {
				bn.drawFootY = 400
			}
		}
		game := g
		// #25: exit via the door. Arrow is "up"; the hotspot's onInteract path
		// first walks PP UP to the door (walkToAndDo to the hotspot centre runs
		// before onInteract), then onInteract walks him OUT to the right into the
		// street before transitioning.
		for i := range bakery.hotspots {
			if bakery.hotspots[i].targetScene == "paris_street" {
				bakery.hotspots[i].onInteract = func() bool {
					// User playtest #26: do NOT shrink on exit. PP has already
					// walked to the door (the hotspot walks him to its centre
					// before onInteract fires), so just transition out at full
					// size - no recede.
					game.sceneMgr.transitionTo("paris_street", game.player)
					return true
				}
				break
			}
		}
		for _, n := range bakery.npcs {
			if n.name != "Madame Poulain" {
				continue
			}
			poulain := n
			souvenirArmed := false
			poulain.onDialogEnd = func() {
				// Subsequent clicks while the rolling pin is still missing
				// just replay the lost-pin beat (no flag flip yet).
				// User 2026-05-20: once Marcus is healed, Poulain pivots
				// to the next anchor beat - asking for a Louvre postcard
				// for her grandson. Wires the next chapter so the bakery
				// stops looping the trade-complete line forever.
				if game.marcusHealed && !souvenirArmed {
					poulain.dialog = bakeryWomanLouvreSouvenirDialog
					souvenirArmed = true
					return
				}
				// The souvenir ask itself just played → Beaumont will sign a
				// second postcard from now on (see the curator's altDialog).
				if souvenirArmed && !souvenirAsked {
					souvenirAsked = true
				}
			}
			// Counter service (2026-06-10): after the rolling-pin trade Poulain
			// becomes the renewable source the side quests draw on. Branches
			// are checked in priority order; nil falls through to her regular
			// dialog. Armed by the trade callback below.
			poulainCounterService := func() ([]dialogEntry, func(), *handOff) {
				// 1) Signed postcard hand-in (grandson souvenir loop).
				if held := game.inv.heldItem; held != nil && held.name == "Signed Postcard" && !souvenirDone {
					return bakeryWomanSouvenirThanksDialog, func() {
						game.inv.giveItemTo("Signed Postcard", "madame_poulain")
						souvenirDone = true
						poulain.dialog = bakeryWomanSouvenirDoneDialog
					}, &handOff{item: "Signed Postcard"}
				}
				// 2) Day-old heel to shoo the flower-pot pigeon (PR#29). Offered
				// once Camille has sent PP after the pencil, until Pierre has
				// shooed the bird.
				if camilleAsked && !pigeonsCleared && !game.inv.hasItem("Baguette Heel") {
					return bakeryWomanHeelDialog, func() {
						// SKILL.md §8b: hand-overs are animated on both sides.
						// Her give sheet + PP's baguette receive are the closest
						// existing sheets (it IS a baguette end). PR#18: face PP.
						poulain.flipped = (game.player.x + playerDstW/2) < float64(poulain.bounds.X+poulain.bounds.W/2)
						poulain.playOneShotAnimThen(poulain.giveAnimOr("give_heel"), 1.5, func() {
							game.player.playOneShot(game.player.resolveOneShot("get_baguette", true, "get_baguette"), 1.6, func() {
								if item := game.items.createItem("baguette_heel"); item != nil {
									game.inv.addItem(item)
								}
							})
						})
					}, nil
				}
				// 3) Coffee refill while Henri's confiture trade is still
				//    pending - keeps the chain unstuckable if the first cup
				//    goes cold (or a later quest borrows it).
				henriWaiting := bakeryHenri != nil && bakeryHenri.altDialogFunc != nil
				if !game.inv.hasItem("Cafe au Lait") && henriWaiting {
					return bakeryWomanCoffeeRefillDialog, func() {
						// §8b / #13: sequence the hand-over - she lifts the cup over
						// the counter, THEN PP takes it (was give+receive in parallel).
						poulain.playOneShotAnimThen(poulain.giveAnimOr("give_coffee"), 1.0, func() {
							game.player.playReceive("cafe_au_lait", true, 1.0, func() {
								if item := game.items.createItem("cafe_au_lait"); item != nil {
									game.inv.addItem(item)
								}
							})
						})
					}, nil
				}
				return nil, nil, nil
			}
			poulain.altDialogRequiresItem = "Rolling Pin"
			// User playtest #25: PP must actively hand the rolling pin over (pull
			// it from the bag and drop it on Poulain) - having it in the bag is no
			// longer enough. Clicking her without holding it just replays her
			// lost-pin line, which is the nudge to bring it.
			poulain.altDialogRequiresHeld = true
			poulain.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
				if !game.inv.hasItem("Rolling Pin") || game.inv.hasItem("Baguette") {
					return nil, nil, nil
				}
				// User 2026-05-21: Poulain now hands out BOTH a Baguette and
				// a Cafe au Lait when PP returns her rolling pin. The coffee
				// is for Henri (who's been waiting in the cafe corner), and
				// the baguette goes to Pierre via the new 2-stage trade.
				return []dialogEntry{
						{speaker: "Pink Panther", text: "I think I found what you were looking for, madame."},
						{speaker: "Madame Poulain", text: "My rolling pin! Bless you, monsieur!"},
						{speaker: "Madame Poulain", text: "Here - your baguette, fresh and warm."},
						{speaker: "Madame Poulain", text: "And take zis cafe au lait too. Henri has been waiting for one all morning - give it to him for me, oui?"},
					}, func() {
						// #25: the pin hand-over plays pre-dialog (PR#1). After her
						// lines Poulain hands the baguette over the counter and PP
						// takes it (cosmetic hand-off; items still land in the bag).
						// PR#18: face her toward PP so she hands it in his direction
						// (her sheet draws facing right; flip if PP is to her left).
						poulain.flipped = (game.player.x + playerDstW/2) < float64(poulain.bounds.X+poulain.bounds.W/2)
						game.inv.giveItemTo("Rolling Pin", "madame_poulain")
						// 2026-06-24 (#12/#13): sequence the two hand-backs like
						// Henri's working jam trade instead of firing give + receive
						// in parallel (which read as "broken"). Poulain hands the
						// baguette → PP takes it → she hands the coffee → PP takes it.
						// Items still land in the bag at the end of the chain.
						poulain.playOneShotAnimThen(poulain.giveAnimOr("give_baguette"), 1.5, func() {
							game.player.playOneShot(game.player.resolveOneShot("get_baguette", true, "get_baguette"), 1.6, func() {
								if b := game.items.createItem("baguette"); b != nil {
									game.inv.addItem(b)
								}
								poulain.playOneShotAnimThen(poulain.giveAnimOr("give_coffee"), 1.0, func() {
									game.player.playReceive("cafe_au_lait", true, 1.0, func() {
										if c := game.items.createItem("cafe_au_lait"); c != nil {
											game.inv.addItem(c)
										}
									})
								})
							})
						})
						poulain.dialog = bakeryWomanPostDialog
						// Pivot from the rolling-pin gate to open counter
						// service (coffee refills, the pigeon heel, the
						// signed-postcard hand-in).
						poulain.altDialogRequiresItem = ""
						poulain.altDialogRequiresHeld = false
						poulain.altDialogFunc = poulainCounterService
					}, &handOff{item: "Rolling Pin", back: true}
			}
			break
		}

		// --- Henri (cafe patron): coffee → confiture trade ---
		// User 2026-05-21: PP brings the Cafe au Lait → Henri trades it for
		// homemade Confiture from his bag. The confiture is then traded to
		// Pierre (stage 2 of his 2-stage quest) for the press pass.
		for _, n := range bakery.npcs {
			if n.name != "Monsieur Henri" {
				continue
			}
			henri := n
			henri.onDialogEnd = func() {
				// While PP doesn't have the coffee yet, just replay the
				// "fetch me a coffee" beat on subsequent clicks.
			}
			// User playtest #25: PP must hand the cafe au lait over (held),
			// not just carry it, before Henri trades the confiture.
			henri.altDialogRequiresHeld = true
			henri.altDialogRequiresItem = "Cafe au Lait"
			henri.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
				if !game.inv.hasItem("Cafe au Lait") || game.inv.hasItem("Confiture") {
					return nil, nil, nil
				}
				return henriCoffeeTradeDialog, func() {
					// 2026-06-20 #13: clean BRING-then-PICK-UP order. The coffee was
					// already handed over in the pre-dialog handoff, so finalize that
					// first; THEN Henri digs the jam out of his bag (give_jam), and
					// only when that finishes does PP take it (get_jam) and the
					// Confiture lands in the bag. Previously give_jam + get_jam fired
					// simultaneously and the item was added instantly, so the receive
					// read as broken/instant.
					game.inv.giveItemTo("Cafe au Lait", "monsieur_henri")
					henri.playOneShotAnimThen("give_jam", 1.3, func() {
						game.player.playOneShot("get_jam", 1.5, func() {
							if c := game.items.createItem("confiture"); c != nil {
								game.inv.addItem(c)
							}
						})
					})
					henri.dialog = henriPostTradeDialog
					henri.altDialogFunc = nil
					henri.altDialogRequiresItem = ""
				}, &handOff{item: "Cafe au Lait"}
			}
			break
		}

		// --- "Camille and the Sold-Out Postcard" (2026-06-10 rework) ---
		// Main-chain quest gate: Beaumont's postcards are sold out, so the
		// postcard now flows museum → bakery → street → bakery → museum.
		// Camille's altDialog is a multi-branch func (no requires fields, so
		// it's consulted on every click and falls through to her dialog):
		//   1. PP holds the Charcoal Pencil → sketch one-shot + Camille's Sketch
		//   2. Beaumont asked + Camille hasn't yet → the lost-pencil ask
		//   3. otherwise → nil (flavor dialog / reminder)
		for _, n := range bakery.npcs {
			if n.name != "Mademoiselle Camille" {
				continue
			}
			camille := n
			// User 2026-06-10: the sketching one-shot
			// (npc_camille_sketching.png) already exists - show it off on her
			// first regular chat, so she's seen mid-sketch from the start.
			sketchShown := false
			camille.onDialogEnd = func() {
				if !sketchShown {
					sketchShown = true
					camille.playOneShotAnimHold("sketch", 2.0, 1.0) // PR#14: slower + hold the reveal
				}
			}
			camille.altDialogFunc = func() ([]dialogEntry, func(), *handOff) {
				held := game.inv.heldItem
				holdingPencil := held != nil && held.name == "Charcoal Pencil"
				// (debug print removed 2026-06-12 - the hover probe calls this
				// func every frame, so it spammed the console on mere hover)
				// Branch 1: pencil hand-over → she sketches the Room 7 replica.
				// Gate on EITHER quest flag (2026-06-11 #38: the story dead-ended
				// here - accept the pencil as soon as anyone asked for it).
				if holdingPencil && (sketchAsked || camilleAsked) && !sketchDone {
					return camilleSketchTradeDialog, func() {
						// §8b: the pencil hand-over plays pre-dialog (PR#1). Her
						// "watch zis" lines play, then the sketching one-shot is
						// the give-back and PP takes the page.
						// PR#14: slower (2.6s) + hold the reveal 1.2s so the
						// finished sketch is readable before she reverts to idle.
						camille.playOneShotAnimHold("sketch", 2.6, 1.2)
						game.player.playOneShot("receive_item", 1.0, nil)
						game.inv.giveItemTo("Charcoal Pencil", "camille")
						if item := game.items.createItem("camille_sketch"); item != nil {
							game.inv.addItem(item)
						}
						sketchDone = true
						camille.dialog = camillePostSketchDialog
					}, &handOff{item: "Charcoal Pencil"}
				}
				// Holding the pencil before the quest is active - never silent.
				if holdingPencil && !sketchDone {
					return []dialogEntry{
						{speaker: "Mademoiselle Camille", text: "Zat is a fine charcoal pencil, monsieur. Mine is still lost out on ze street somewhere..."},
					}, nil, nil
				}
				// Pencil in the BAG but not on the cursor: nudge the hand-over
				// motion instead of replaying the stale "ask Nicolas" reminder
				// (2026-06-11 #38 - this read as the quest being stuck).
				if camilleAsked && !sketchDone && game.inv.hasItem("Charcoal Pencil") {
					return []dialogEntry{
						{speaker: "Mademoiselle Camille", text: "You FOUND it?! Don't tease an artist, monsieur - take it from your bag and hand it here!"},
					}, nil, nil
				}
				// Branch 2: Beaumont has asked → Camille's lost-pencil ask.
				if sketchAsked && !camilleAsked {
					return camilleSketchAskDialog, func() {
						camilleAsked = true
						camille.dialog = camillePencilReminderDialog
						// PR#23: her dismay one-shot (§CAM2) plays on the ask.
						camille.playOneShotAnimHold("lost_pencil", 2.0, 0.8)
					}, nil
				}
				return nil, nil, nil
			}
			break
		}

	}

	// User playtest #14: the rolling pin is HIDDEN inside the bicycle basket on
	// the cobblestone street (~539,644). Its sprite is NOT drawn - the player
	// only discovers it because the cursor changes to the grab hand over the
	// basket. Picking it up plays PP's dedicated "grab rolling pin" animation
	// (reach into basket, lift overhead) before the item lands in the bag.
	if parisStreet, ok := g.sceneMgr.scenes["paris_street"]; ok {
		game := g

		// 2026-06-12 #12 (encounter v2): clicking the biker no longer brakes
		// him on the spot. PP walks into the lane AHEAD of him while he keeps
		// riding; the moment the bike reaches PP it brakes (dedicated braked
		// pose), PP flinches (stateReacting holds through the dialog because
		// player.update is gated on !dialog.active), the apology plays, and
		// he rides on when it closes. If PP can't reach the lane in time the
		// wrap-around brings the biker back and the bump fires on his next
		// pass. bikerBumpCheck runs from Game.Update while on paris_street.
		for _, amb := range parisStreet.ambientSprites {
			a := amb
			bumpArmed := false
			a.onClick = func() {
				if bumpArmed || a.paused {
					return
				}
				bumpArmed = true
				// Meet point: ahead of the biker in his lane, clamped
				// on-screen. walkToAndDo's y is PP's CENTER - a.y minus half
				// the player box plants PP's feet on the biker's ground line.
				meetX := a.x + 320
				if meetX < 300 {
					meetX = 300
				} else if meetX > 1150 {
					meetX = 1150
				}
				game.player.walkToAndDo(meetX, a.y-float64(playerDstH)/2, nil)
			}
			g.bikerBumpCheck = func() {
				if !bumpArmed || a.paused || game.dialog.active || game.seqPlayer.IsPlaying() {
					return
				}
				px, py := game.player.footCenter()
				if afAbs(float64(py)-a.y) > 70 {
					return // PP hasn't reached the lane yet
				}
				if afAbs(a.x-float64(px)) > 50 {
					return
				}
				bumpArmed = false
				a.paused = true // brake: ambient update holds the braked pose
				game.player.moving = false
				// PR#9: PP recoils/hops backward from the bump. Plays the
				// dedicated jump-back one-shot if its art has landed, else the
				// generic flinch. Dialog starts after the recoil.
				bikerLines := []dialogEntry{
					{speaker: "Biker", text: "Pardon, pardon! Sorry monsieur, but you are blocking ze way!"},
					{speaker: "Pink Panther", text: "My apologies. Nice bell."},
					{speaker: "Biker", text: "Merci! Bonne journee!"},
				}
				if game.player.hasOneShot("jump_back") {
					game.player.playOneShot("jump_back", 0.7, func() {
						game.dialog.startDialogWithCallback(bikerLines, func() { a.paused = false })
					})
				} else {
					game.player.playAction(stateReacting, nil)
					game.dialog.startDialogWithCallback(bikerLines, func() { a.paused = false })
				}
			}
		}

		pin := &floorItem{
			// Hidden: no sprite drawn. Bounds cover the bike basket so the
			// grab cursor lights up there.
			bounds:  sdl.Rect{X: 499, Y: 614, W: 80, H: 60},
			name:    "Rolling Pin",
			visible: false,
			hidden:  true,
			// 2026-06-12 #15: stand to the basket's RIGHT so the grab anim's
			// reach hand lands inside the basket instead of past it.
			standRight: true,
			onPickup: func() {
				// Mark it taken so it can't be re-grabbed, then play the grab
				// one-shot and only add the item + dialog when the anim ends.
				if ps, ok := game.sceneMgr.scenes["paris_street"]; ok {
					for _, fi := range ps.floorItems {
						if fi.name == "Rolling Pin" {
							fi.visible = false
							fi.hidden = false
							break
						}
					}
				}
				// PR#16: PP stands to the basket's RIGHT (standRight) so he must
				// face LEFT to reach into it - flip 180 from the default. The
				// grab_rolling_pin draw offset (drawScaled) also drops him down so
				// the reach lands in the basket.
				game.player.facingLeft = true
				game.player.dir = dirLeft
				game.player.playOneShot("grab_rolling_pin", 1.0, func() {
					if item := game.items.createItem("rolling_pin"); item != nil {
						game.inv.addItem(item)
					}
					// 2026-06-11 #18 / SKILL.md §8c: pickup lines are GENERIC -
					// the "who needs this" hint lives in NPC dialogs instead.
					game.dialog.startDialog(genericPickupDialog(
						"A wooden rolling pin, tucked away in someone's bike basket."))
				})
			},
		}
		parisStreet.floorItems = append(parisStreet.floorItems, pin)

		// "Camille and the Sold-Out Postcard": her lucky charcoal pencil sits
		// in the flower pot by the Louvre steps (right side of the street,
		// below the museum hotspot). The pot is visible; pigeons guard it until
		// Pierre repays his favor, then the art swaps to the exposed-pencil state.
		potTex, potW, potH := engine.SafeTextureFromPNGKeyed(g.renderer, "assets/images/locations/paris/props/flower_pot_pigeon.png")
		pencil := &floorItem{
			// 2026-06-20 #8: the pot sat big in the middle of the street. Shrunk
			// and moved beside Pierre's easel (Pierre x780 foot 645), where a
			// street-painter's flower pot reads logically and at his mid-distance
			// scale. bounds foot 650 (≈ Pierre's foot), just to his right.
			tex:     potTex,
			srcW:    potW,
			srcH:    potH,
			bounds:  sdl.Rect{X: 868, Y: 562, W: 82, H: 88},
			name:    "Charcoal Pencil",
			visible: true,
			hidden:  true,
			onPickup: func() {
				// The pigeons guard the pot until Pierre repays his favor
				// (user 2026-06-10): trying early plays a blocked beat and
				// leaves the pencil in place.
				if !pigeonsCleared {
					if !camilleAsked {
						game.dialog.startDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "A pencil, deep in a flower pot... and a very protective pigeon sitting on it. I'll leave it be for now."},
						})
						return
					}
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "Ow! This pigeon guards Camille's pencil like the crown jewels."},
						{speaker: "Pink Panther", text: "Madame Margaux, the pigeon lady across the street, can coax it off - I'll bring her a day-old baguette heel from Madame Poulain."},
					})
					return
				}
				if ps, ok := game.sceneMgr.scenes["paris_street"]; ok {
					for _, fi := range ps.floorItems {
						if fi.name == "Charcoal Pencil" {
							fi.visible = false
							fi.hidden = false
							break
						}
					}
				}
				// 2026-06-15 #19/#20: was playAction(stateGrabbing, ...) which
				// player.update kills before the callback runs, so the pencil was
				// never added and PP appeared stuck. Use the guaranteed "grab"
				// one-shot (fires onDone even if art is missing) like the rolling pin.
				game.player.playOneShot("grab", 1.0, func() {
					pencilTaken = true
					if item := game.items.createItem("charcoal_pencil"); item != nil {
						game.inv.addItem(item)
					}
					// §8c: generic pickup line - Camille's own dialog already
					// nudges where the pencil goes.
					game.dialog.startDialog(genericPickupDialog(
						"A charcoal pencil, rescued from the pigeons."))
				})
			},
		}
		parisStreet.floorItems = append(parisStreet.floorItems, pencil)
	}
}

func (g *Game) setupTravelHotspots() {
	game := g

	if campEntrance, ok := g.sceneMgr.scenes["camp_entrance"]; ok {
		for i := range campEntrance.hotspots {
			if campEntrance.hotspots[i].name == "Enter Camp" {
				campEntrance.hotspots[i].onInteract = func() bool {
					// User 2026-04-26: PP no longer strafes left to (599,200);
					// instead he walks back-facing in place and shrinks into
					// the distance over ~1.6s, then the scene transitions.
					game.player.playRecede(1.6, 0.35, 80, func() {
						if grounds, ok := game.sceneMgr.scenes["camp_grounds"]; ok {
							// Destination PP walks IN to (on the path), not the
							// far-left edge - he enters from off-screen left and
							// strolls onto the path (#2).
							// User playtest #35: on arrival PP walks in from the
							// off-screen left to ~(245, 652) - foot at y≈652 means a
							// top-left spawnY of 652-270=382. (Was overshooting to
							// ~861,781.) Starting values; tune against the path in-game.
							grounds.spawnX = 245
							grounds.spawnY = 382
						}
						// Flag THIS transition as the camp arrival so PP walks in
						// from the left; room exits leave it unset (#2/#14).
						game.sceneMgr.entryWalkPending = true
						game.sceneMgr.transitionTo("camp_grounds", game.player)
					})
					return true
				}
				break
			}
		}
	}

	// User 2026-05-21: removed the "Camp Chilly Wa Wa Air" left-arrow
	// hotspot from camp_entrance. Travel is opened by clicking the
	// Travel Map item in the inventory (see g.inv.onSelectItem above) -
	// the scene-edge hotspot was a stale leftover from the pre-2026-04-26
	// design and showed up as a confusing left arrow on first entry to
	// the game.
}

// openTravelMap is THE way to open the travel globe (user 2026-06-12): PP
// first plays the §PM1 "pull_map" one-shot - pulling the map out of his
// invisible hip pocket - and the map screen opens when it finishes. Every
// open path (inventory item click, held-item drop, the city street
// hotspots) routes through here so the beat is consistent. While the §PM1
// sheet isn't on disk, playOneShot fires the callback immediately, so the
// map still opens (just without the flourish).
func (g *Game) openTravelMap(fromScene string) {
	g.player.playOneShot("pull_map", 0.9, func() {
		g.travelMap.Show(fromScene)
	})
}

func (g *Game) Close() {
	g.audio.close()
}

// walkToFloorItem (2026-06-12 #14/#15, SKILL.md §8c): PP walks to a stand
// point BESIDE a floor item instead of on top of it - feet aligned with the
// item's base, to its LEFT by default (standRight flips sides for grabs whose
// reach hand points the other way). On arrival he squares up to the camera
// (dir front) so blocked/observation lines play facing the player, then the
// pickup action runs.
func (g *Game) walkToFloorItem(fi *floorItem, action func()) {
	itemCX := float64(fi.bounds.X) + float64(fi.bounds.W)/2
	itemBase := float64(fi.bounds.Y) + float64(fi.bounds.H)
	offset := float64(fi.bounds.W)/2 + 80
	standX := itemCX - offset
	if fi.standRight {
		standX = itemCX + offset
	}
	plr := g.player
	// walkToAndDo's y is PP's CENTER: itemBase minus half the player box
	// plants his feet on the item's base line (scene minY/maxY still clamp).
	plr.walkToAndDo(standX, itemBase-float64(playerDstH)/2, func() {
		plr.dir = dirDown
		plr.facingLeft = false
		if action != nil {
			action()
		}
	})
}

func (g *Game) HandleClick(x, y int32) {
	g.ui.triggerClick()

	// Dev menu (F1 toggleable) sits above everything else so a click on a
	// scenario row routes here before the world / inventory.
	if g.devMenu != nil && g.devMenu.Visible() {
		if g.devMenu.handleClick(x, y, g) {
			return
		}
	}

	// Pause menu sits above everything - click routes here first.
	if g.menuHandleClick(x, y) {
		return
	}

	// Click probe (F2) intercepts clicks before any game logic runs:
	// validate the bbox-hit against the source PNG alpha and place a
	// green/red marker. Always swallows the click so probing doesn't also
	// trigger talk dialog or movement.
	if g.clickProbe != nil && g.clickProbe.active {
		if scene := g.sceneMgr.current(); scene != nil {
			g.clickProbe.recordClick(scene, x, y)
		} else {
			g.clickProbe.pushMarker(x, y, sdl.Color{R: 160, G: 160, B: 160, A: 255}, "no scene")
		}
		return
	}

	if g.travelMap.Visible() {
		// Info panel eats clicks: any click while open dismisses it,
		// leaving the map visible underneath.
		if g.travelMap.panelHandleClick() {
			return
		}
		if loc := g.travelMap.hitTest(x, y); loc != nil {
			if loc.scene == g.travelMap.ReturnScene() {
				g.travelMap.Hide()
				return
			}
			g.travelMap.Hide()
			if loc.scene == "camp_entrance" {
				// User 2026-06-08 (#34): returning home now arrives BY PLANE at
				// the airstrip (camp_landing) rather than popping straight to the
				// forest gate. PP then walks right through the gate into the
				// (now darkened) camp grounds. The pin/unlock still key on
				// "camp_entrance"; only the arrival scene is redirected.
				g.flight.Start("camp_landing")
				g.sceneMgr.transitionTo("airplane_flight", g.player)
				return
			}
			g.flight.Start(loc.scene)
			g.sceneMgr.transitionTo("airplane_flight", g.player)
			return
		}

		// Non-travel click: open info panel for locked pins AND for
		// unlocked pins that aren't the current story target. Only the
		// relevant pin gets through to hitTest above and actually travels.
		if anyLoc := g.travelMap.hitTestAny(x, y); anyLoc != nil {
			if len(anyLoc.facts) > 0 || anyLoc.info != "" {
				if anyLoc.audio != "" {
					g.audio.playSFX(anyLoc.audio)
				}
				g.travelMap.openInfoPanel(anyLoc)
			}
			return
		}
		return
	}

	if g.inv.open {
		g.inv.handleClick(x, y)
		return
	}
	if g.dialog.active {
		g.dialog.advance()
		return
	}
	if g.sceneMgr.transitioning {
		return
	}
	// User playtest #4: while a cutscene sequence is playing (e.g. Jake walking
	// off after his Day-1 dialog), swallow world clicks so clicking the
	// still-moving NPC can't restart his dialog. (Dialog advance is handled
	// above, so scripted dialog steps still work.)
	if g.seqPlayer.IsPlaying() {
		return
	}
	scene := g.sceneMgr.current()

	// PP-click handling. Two cases:
	//   (a) Player is NOT holding an item → click on PP opens inventory.
	//       Empty-inventory click is silently eaten (no fall-through to
	//       hotspot behind PP, per user 2026-05-23).
	//   (b) Player IS holding an item → fall through to the held-item drop
	//       logic below. User 2026-05-24: previously this block toggled
	//       inventory for ANY click on PP, which meant trying to deliver a
	//       flower to Lily standing next to PP would hit PP's rect first
	//       and re-pocket the item - making the give-item beat
	//       undeliverable ("can't move item to Lily, game stack").
	if g.player.containsPoint(x, y) && g.inv.heldItem == nil {
		if len(g.inv.items) > 0 {
			g.inv.toggle()
		}
		return
	}

	if g.inv.heldItem != nil {
		if clickedNPC := scene.checkNPCClick(x, y); clickedNPC != nil {
			// #33 (2026-06-11): NPCs with a scripted click flow (Pierre's
			// walk-up + recede choreography) must get it for hand-overs too -
			// the inline path below skipped the choreography, so giving Pierre
			// the baguette traded at full size with no shrink. The override
			// reads g.inv.heldItem itself; giveItemTo clears it on a match.
			if clickedNPC.onClickOverride != nil {
				ov := clickedNPC.onClickOverride
				g.player.interactTarget = nil
				ov()
				return
			}
			if clickedNPC.altDialogFunc != nil &&
				(clickedNPC.altDialogRequiresItem == "" ||
					g.inv.heldItem.name == clickedNPC.altDialogRequiresItem) {
				entries, cb, ho := clickedNPC.altDialogFunc()
				if entries != nil {
					g.inv.heldItem = nil
					ds := g.dialog
					target := clickedNPC
					// #10: walk to the talk spot and fire the hand-off there -
					// walkToTalkPos runs the callback whether PP walks in or is
					// already adjacent, so giving Lily the flower works standing
					// next to her (the old walkToAndInteract snapped and opened
					// her NORMAL dialog when close, skipping this hand-off).
					g.player.interactTarget = nil
					g.player.walkToTalkPos(target, func() {
						targetCenter := float64(target.bounds.X + target.bounds.W/2)
						playerCenter := g.player.x + playerDstW/2
						g.player.facingLeft = playerCenter > targetCenter
						if g.player.facingLeft {
							g.player.dir = dirLeft
						} else {
							g.player.dir = dirRight
						}
						// Same face-toward-PP flip used by startNPCDialog - the
						// drag-onto-NPC path runs its own dialog start so it
						// needs its own snapshot/restore. Otherwise dropping
						// Lily's flower would leave her still facing into the
						// flower patch even though she's mid-conversation.
						target.preTalkFlipped = target.flipped
						target.flipped = playerCenter < targetCenter
						wrappedCb := func() {
							if len(target.talkGrid) > 0 {
								target.setAnimState(npcAnimIdle)
							}
							target.flipped = target.preTalkFlipped
							if cb != nil {
								cb()
							}
						}
						// PR#1 (2026-06-12): the item changes hands BEFORE the
						// text - PP's give one-shot, the NPC's receive one-shot,
						// then talk state + dialog.
						start := func() {
							// PR#1 (#1): re-assert side-facing at the moment talk
							// begins. The give one-shot runs for ~2s after the
							// arrival callback set dir, and can leave it stale so
							// the front-talk sheet played (user: "PP talks front
							// not side after giving Lily the flower").
							pc := g.player.x + playerDstW/2
							g.player.facingLeft = pc > targetCenter
							if g.player.facingLeft {
								g.player.dir = dirLeft
							} else {
								g.player.dir = dirRight
							}
							g.player.state = stateTalking
							if len(target.talkGrid) > 0 {
								target.setAnimState(npcAnimTalk)
							}
							ds.startDialogWithCallback(entries, wrappedCb)
						}
						if ho != nil {
							g.player.playHandOff(target, ho, start)
						} else {
							start()
						}
					})
					return
				}
			}
		}
		// Travel Map: drop anywhere to open globe
		if g.inv.heldItem.name == "Travel Map" {
			g.inv.heldItem = nil
			g.openTravelMap(g.sceneMgr.currentName)
			return
		}
		// PR#27: clicking a FLOOR ITEM (e.g. the flower pot) while holding
		// something used to fall through and silently pocket the held item
		// ("clicked the pot with the pencil and it just disappeared"). Keep the
		// item on the cursor instead - a floor-item click isn't a hand-over.
		if scene.checkFloorItemClick(x, y) != nil {
			return
		}
		g.inv.heldItem = nil
		return
	}

	// Clickable ambient sprites (2026-06-11 #16: the crossing biker). Checked
	// before floor items / NPCs - he's a moving target, the click should win.
	// 2026-06-12 #12: no auto-pause here anymore - the encounter handler
	// decides when the biker actually brakes (he keeps riding until he
	// reaches PP in the lane).
	for _, amb := range scene.ambientSprites {
		if amb.containsPoint(x, y) && !amb.paused {
			cb := amb.onClick
			if cb != nil {
				cb()
			}
			return
		}
	}

	// User 2026-05-22: floor items BEFORE npcs so a pickable item that sits
	// under an NPC's bounds rect (e.g. the rolling pin near Nicolas) routes
	// the click to the pickup, not to the NPC dialog. Matches the cursor
	// priority in ui.go.
	if fi := scene.checkFloorItemClick(x, y); fi != nil {
		fiLocal := fi
		g.walkToFloorItem(fiLocal, func() {
			if fiLocal.onPickup != nil {
				fiLocal.onPickup()
			}
		})
		return
	}
	if npc := scene.checkNPCClick(x, y); npc != nil {
		// User 2026-05-31 (#3): log which NPC a click started a dialog with,
		// and where the click landed, to make talk issues debuggable.
		fmt.Printf("[dialog] click (%d,%d) in %q started dialog with %q (npc bounds %+v)\n",
			x, y, g.sceneMgr.currentName, npc.name, npc.bounds)
		g.player.walkToAndInteract(npc, g.dialog)
		return
	}
	if hs := scene.checkHotspotClick(x, y); hs != nil {
		// PR#25: if PP is frozen at Pierre's shrunk depth (recedeHeld), a
		// hotspot click used walkToAndDo and marched him off still tiny.
		// Release the recede first so he grows back to full size as he heads
		// to the exit ("move to the main road, then walk left").
		if g.player.recedeHeld {
			g.player.releaseRecedeSmooth(0.5)
		}
		if hs.onInteract != nil {
			hsLocal := hs
			g.player.walkToAndDo(
				float64(hsLocal.bounds.X+hsLocal.bounds.W/2),
				float64(hsLocal.bounds.Y+hsLocal.bounds.H/2),
				func() {
					hsLocal.onInteract()
				},
			)
			return
		}
		tgt := hs.targetScene
		// 2026-06-21: the FIRST landing (day-1 arrival, before any city trip) leads
		// to the camp ENTRANCE to meet Higgins; returns from a city (paris_done)
		// land straight at the (dark) grounds. Same up-arrow, different target.
		if g.sceneMgr.currentName == "camp_landing" && tgt == "camp_grounds" &&
			!g.vars.GetBool(ScopeGame, VarParisDone) {
			tgt = "camp_entrance"
		}
		plr := g.player
		sm := g.sceneMgr
		onArrival := func() { sm.transitionTo(tgt, plr) }
		// Cabin doors: walk to the door anchor, then shrink-and-rise into the
		// frame instead of marching off the top of the screen. The hotspot's
		// arrow is "up" but walkToExit("up") drives Y to -playerDstH which
		// reads as PP "flying to the sky". playRecede holds X, drifts up by
		// dyUp and shrinks 1.0 -> endScale, which reads as "stepping through
		// the door". See FIXME.md trailing "new PR" block.
		if hs.arrow == arrowUp && strings.HasSuffix(tgt, "_room") {
			doorX := float64(hs.bounds.X + hs.bounds.W/2)
			doorY := float64(hs.bounds.Y + hs.bounds.H/2)
			plr.walkToAndDo(doorX, doorY, func() {
				plr.playRecede(0.7, 0.45, 60, onArrival)
			})
			return
		}
		// User 2026-06-12 (#10): the office exit used walkToExit(downRight)
		// which marched PP diagonally off-screen from wherever he stood.
		// Route him along the painted path instead: mid-path on the main
		// trail first, then down the right-hand trail, then transition.
		if g.sceneMgr.currentName == "camp_grounds" && tgt == "camp_office" {
			plr.walkToAndDo(900, 483, func() {
				plr.walkToAndDo(1050, 640, onArrival)
			})
			return
		}
		// 2026-06-20 #18: camp_landing exit walks PP UP the dirt road to the gate
		// (waypoints along the painted road) instead of straight up from wherever
		// he stands, which read as "walking to the side."
		if g.sceneMgr.currentName == "camp_landing" && (tgt == "camp_grounds" || tgt == "camp_entrance") {
			// 2026-06-24 (#1): at the camp ENTRANCE, PP strolls in from off-screen
			// left to his spawn mark (same arrival beat as the grounds) instead of
			// popping in at the gate.
			if tgt == "camp_entrance" {
				sm.entryWalkPending = true
			}
			plr.walkToAndDo(820, 540, func() {
				plr.walkToAndDo(1160, 440, onArrival)
			})
			return
		}
		if hs.arrow == arrowLeft || hs.arrow == arrowRight || hs.arrow == arrowDown || hs.arrow == arrowUp || hs.arrow == arrowDownRight {
			plr.walkToExit(hs.arrow, onArrival)
		} else {
			plr.walkToAndDo(
				float64(hs.bounds.X+hs.bounds.W/2),
				float64(hs.bounds.Y+hs.bounds.H/2),
				onArrival,
			)
		}
		return
	}
	// If a pending NPC interaction is still in flight (PP is walking to talk
	// to someone), don't let a stray floor click clear it. The user often
	// double-clicks NPCs while waiting; the second click would normally
	// snap-to-path and call setTarget which nukes interactTarget.
	if g.player.interactTarget != nil && g.player.moving {
		return
	}
	tx, ty := float64(x), float64(y)
	tx, ty = scene.snapToPath(tx, ty)
	if g.walkDbg != nil {
		g.walkDbg.recordSnap(tx, ty)
	}
	g.player.setTarget(tx, ty)
}

func (g *Game) HandleKey(scancode sdl.Scancode) {
	// F1 toggles the dev/chapter-jump menu from anywhere. Keep this above
	// the menu modal-eat so the dev can dismiss the dev menu even when the
	// pause menu is open.
	if scancode == sdl.SCANCODE_F1 {
		if g.devMenu != nil {
			g.devMenu.toggle()
		}
		return
	}
	// F2 toggles the click probe - a dev diagnostic that, while active,
	// turns clicks into alpha-channel hit-tests on NPC sprites and drops
	// a green/red marker at the click point. See click_probe.go.
	if scancode == sdl.SCANCODE_F2 {
		if g.clickProbe != nil {
			g.clickProbe.toggle()
		}
		return
	}
	// F3 toggles the walk-debug overlay - draws walkSegments, PP's foot
	// point and the last snapped click target with live coordinates, so the
	// painted-path tuning is done with exact numbers (2026-06-12 #2).
	if scancode == sdl.SCANCODE_F3 {
		if g.walkDbg != nil {
			g.walkDbg.toggle()
		}
		return
	}
	if g.devMenu != nil && g.devMenu.Visible() {
		// Eat other input while dev menu is up.
		return
	}
	// Esc: first press closes the travel map if open; otherwise toggles the
	// pause menu. This keeps the existing "Esc to leave globe" behavior
	// while adding a universal pause option.
	if scancode == sdl.SCANCODE_ESCAPE {
		switch {
		case g.travelMap.Visible() && g.travelMap.panelVisible():
			// First Esc closes the map info panel; a second Esc closes the map.
			g.travelMap.closeInfoPanel()
		case g.travelMap.Visible():
			g.travelMap.Hide()
		case g.menu.Visible():
			g.menu.Hide()
		default:
			g.menu.Show()
		}
		return
	}
	if g.menu.Visible() {
		// Menu is modal - eat other input until it closes.
		return
	}
	// 2026-06-20 #6: the "m" shortcut to open the travel map was removed (the
	// map opens by clicking the Travel Map inventory item / city hotspots).
	if scancode == sdl.SCANCODE_SPACE && g.dialog.active {
		g.dialog.advance()
	}
	if scancode == sdl.SCANCODE_F5 {
		g.SaveGame("savegame.json")
	}
	if scancode == sdl.SCANCODE_F9 {
		if err := g.LoadGame("savegame.json"); err != nil {
			fmt.Printf("Load failed: %v\n", err)
		}
	}
}

func (g *Game) Update(dt float64, mx, my int32) {
	g.mouseX = mx
	g.mouseY = my

	// Mirror the flat progression flags into the VarStore once a frame so
	// sequences and save files can treat VarStore as the source of truth.
	g.syncFlagsToVars()

	// Pause menu freezes game state - only update its own hover and return.
	if g.menu.Visible() {
		g.menu.UpdateHover(mx, my)
		g.ui.cursor = cursorNormal
		g.ui.hoverName = ""
		return
	}

	// User 2026-05-22: inventory is modal - freeze the world while it's
	// open. Only the inventory's own ticks (pulse animation in inv.draw)
	// keep running. Clicks behind are already blocked at HandleClick:1235.
	// This stops NPC frames, floor-item visibility, ambient effects, and
	// scene-trigger checks from advancing while the player browses items.
	if g.inv.open {
		// 2026-06-11 #4: inside the open bag every click is an inventory
		// action (page/pick) - show the pink POINT hand, not the arrow.
		g.ui.cursor = cursorPoint
		g.ui.hoverName = ""
		return
	}

	if g.travelMap.Visible() {
		g.ui.hoverName = ""
		g.ui.cursor = cursorNormal
		if loc := g.travelMap.hitTestAny(mx, my); loc != nil {
			g.ui.hoverName = loc.name
			if loc.unlocked {
				// User playtest #6: only the pin that is the current relevant
				// travel destination gets the active "pointing" cursor; other
				// unlocked (info-only / already-visited) pins keep the talk
				// icon so the player can tell where the story wants them next.
				if g.travelMap.isRelevant(loc) {
					g.ui.cursor = cursorPoint
				} else {
					g.ui.cursor = cursorTalk
				}
			}
		}
		return
	}

	scene := g.sceneMgr.current()

	// 2026-06-21: the game now STARTS at the airstrip (camp_landing) - PP gets
	// off the plane, delivers the opening monologue, then walks up to the camp
	// ENTRANCE to meet Higgins, then into the grounds. (The monologue used to
	// play at camp_entrance.) PP strolls in from off-screen-left to the landing
	// spawn; walkY = the landing's spawn row so he's on the airstrip ground.
	if !g.monologuePlayed && g.sceneMgr.currentName == "camp_landing" && !g.sceneMgr.transitioning {
		g.monologuePlayed = true
		startX := -float64(playerDstW)
		endX := scene.spawnX - float64(playerDstW)/2
		walkY := scene.spawnY
		g.player.playWalkIn(startX, endX, walkY, 2.5, func() {
			g.player.state = stateTalking
			g.player.dir = dirDown
			g.dialog.startDialogWithCallback(openingMonologue, func() {
				g.player.state = stateIdle
			})
		})
	}

	// Night cutscene: when the player first reaches camp_night on Day 1 after
	// meeting all 5 kids, kick off the JSON sequence. It handles Higgins
	// bedtime, PP sleeping, Marcus freakout, waking, and Day-2 transition.
	if !g.nightSceneDone && g.day == 1 && g.metKids >= 5 &&
		g.sceneMgr.currentName == "camp_night" && !g.sceneMgr.transitioning &&
		!g.dialog.active && !g.seqPlayer.IsPlaying() {
		g.nightSceneDone = true
		if seq := g.seqStore.Get("night_bedtime"); seq != nil {
			g.seqPlayer.Play(seq)
		}
	}

	// User 2026-05-22: Day-2 monologue removed from here - it was duplicating
	// the wake-up dialog the night_bedtime.json sequence already plays at
	// lines 38-42 (the "*yawn* What a night..." beat). Now plays exactly
	// once, in-context, at the end of the sequence.
	_ = g.day2Started // keep field for save-state compat; no longer drives dialog

	if !g.parisMonologuePlayed && g.sceneMgr.currentName == "paris_street" && !g.sceneMgr.transitioning {
		g.parisMonologuePlayed = true
		g.dialog.startDialog(parisStreetMonologue)
	}

	// Biker bump encounter (2026-06-12 #12): armed by clicking the biker,
	// fires when he reaches PP standing in his lane.
	if g.bikerBumpCheck != nil && g.sceneMgr.currentName == "paris_street" && !g.sceneMgr.transitioning {
		g.bikerBumpCheck()
	}

	// Japan opening (#8): after PP finds Lily at the lake, the next time he's in
	// the grounds Higgins intercepts him (rude), then PP's camera aside unlocks
	// Tokyo. Fires once; VarHigginsRudeDone is the persistent gate.
	if !g.higginsRudeStarted && g.sceneMgr.currentName == "camp_grounds" &&
		g.vars.GetBool(ScopeGame, VarLilyLakeMet) &&
		!g.vars.GetBool(ScopeGame, VarHigginsRudeDone) &&
		!g.sceneMgr.transitioning && !g.dialog.active && !g.seqPlayer.IsPlaying() {
		g.higginsRudeStarted = true
		g.playHigginsRudeBeat()
	}
	// Drive Higgins's stride-in for the rude intercept (simple x-lerp).
	if g.higginsWalk != nil {
		w := g.higginsWalk
		w.elapsed += dt
		t := w.elapsed / w.dur
		if t >= 1 {
			t = 1
		}
		w.n.bounds.X = int32(w.fromX + (w.toX-w.fromX)*t)
		if t >= 1 {
			g.higginsWalk = nil
			if w.onArrive != nil {
				w.onArrive()
			}
		}
	}

	// Museum first arrival (#28/#30): the first time PP reaches the Louvre, walk
	// him in from the tunnel on the left, then play the arrival monologue. Gated
	// on a VarStore flag so it fires exactly once (survives save/load).
	// 2026-06-11 #34: the old playWalkIn here raced the transition spawn -
	// PP popped on at the spawn point, vanished off-screen, then walked in
	// again. The Louvre gate hotspot now sets entryWalkPending (the single
	// walk-in path, scene.go), so this block only waits for that walk to
	// finish and plays the monologue once.
	if !g.vars.GetBool(ScopeGame, VarMonologueLouvre) && g.sceneMgr.currentName == "paris_louvre" &&
		!g.sceneMgr.transitioning && !g.player.moving && !g.dialog.active {
		g.vars.SetBool(ScopeGame, VarMonologueLouvre, true)
		g.player.state = stateTalking
		// 2026-06-20 #14: PP faces FRONT for the arrival monologue (he's taking in
		// the gallery and talking to the player), not off to the side.
		g.player.dir = dirDown
		g.dialog.startDialogWithCallback(louvreArrivalMonologue, func() {
			g.player.state = stateIdle
		})
	}

	// Finale: plays once the player lands in Mexico City with every heal flag set.
	if g.sceneMgr.currentName == "mexico_street" && !g.sceneMgr.transitioning && !g.dialog.active {
		g.triggerFinaleMonologue()
	}

	// Animate the active scene's background (no-op for static scenes).
	// Runs every frame so animated sheets like airplane_flight's cloud sky
	// keep looping even while the flight cutscene drives the foreground.
	if cur := g.sceneMgr.current(); cur != nil && cur.bg != nil {
		cur.bg.update(dt)
	}

	// Click-probe markers fade with TTL.
	if g.clickProbe != nil {
		g.clickProbe.update(dt)
	}

	// Airplane flight cutscene
	if g.sceneMgr.currentName == "airplane_flight" && !g.sceneMgr.transitioning {
		if done, dest := g.flight.Update(dt); done {
			g.sceneMgr.transitionTo(dest, g.player)
		}
	}

	// Animate campfire
	if g.sceneMgr.currentName == "camp_night" && len(g.campfireFrames) > 0 {
		g.campfireTimer += dt
		if g.campfireTimer >= 0.12 {
			g.campfireTimer -= 0.12
			g.campfireFrameIdx = (g.campfireFrameIdx + 1) % len(g.campfireFrames)
		}
	}

	// Animate sleeping/waking frames (same speed as walking)
	if g.playerSleeping {
		g.sleepingTimer += dt
		if g.sleepingTimer >= 0.15 {
			g.sleepingTimer -= 0.15
			if g.wakingPhase == 0 {
				// Sleeping: loop the breathing cycle.
				if n := len(g.sleepingFrames); n > 0 {
					g.sleepingFrameIdx = (g.sleepingFrameIdx + 1) % n
				}
			} else {
				// User playtest #14: waking plays ONCE and holds the last frame
				// (standing, eyes open) instead of looping ("runs twice and
				// again"). wakingPhase=2 marks it finished.
				if n := len(g.wakingFrames); n > 0 {
					if g.sleepingFrameIdx < n-1 {
						g.sleepingFrameIdx++
					} else {
						g.wakingPhase = 2
					}
				}
			}
		}
	}

	if !g.dialog.active && g.player.state == stateTalking {
		g.player.state = stateIdle
	}
	if g.dialog.active && g.player.state == stateTalking {
		g.player.breathTimer += dt
		speaker := g.dialog.currentSpeaker()
		isPPSpeaking := speaker == "Pink Panther"

		// PP only animates talk when PP is speaking
		if isPPSpeaking {
			g.player.talkTimer += dt
			ppTalkSpeed := talkFrameTime * 2.0 // slower talk for PP
			if g.player.talkTimer >= ppTalkSpeed {
				g.player.talkTimer -= ppTalkSpeed
				frames := g.player.currentTalkFrames()
				if len(frames) > 0 {
					g.player.talkCycleIdx = (g.player.talkCycleIdx + 1) % len(frames)
				}
			}
		} else {
			g.player.talkCycleIdx = 0
			g.player.talkTimer = 0
		}

		// NPC animates talk only when THEY are speaking
		for _, n := range scene.npcs {
			hasTalk := len(n.talkGrid) > 0
			if !hasTalk {
				continue
			}
			npcIsSpeaker := speaker == n.name
			if npcIsSpeaker && n.animState != npcAnimTalk {
				n.setAnimState(npcAnimTalk)
			} else if !npcIsSpeaker && n.animState == npcAnimTalk {
				n.setAnimState(npcAnimIdle)
			}
		}
	}
	if !g.dialog.active && !g.sceneMgr.transitioning {
		g.player.update(dt, scene.blockers)
	}
	g.dialog.update(dt)
	g.seqPlayer.Update(dt)
	g.sceneMgr.update(dt)
	showAmbient := g.day == 1 && isCampOutdoorScene(g.sceneMgr.currentName)
	scene.updateAmbient(dt, showAmbient)
	g.ui.updateHover(scene, mx, my, g.inv, g.player, dt)
	g.inv.update(dt)

	if g.sceneMgr.currentName != g.lastScene {
		g.lastScene = g.sceneMgr.currentName
		g.audio.playMusic(g.sceneMgr.current().musicPath)
	}
}

func (g *Game) Draw(renderer *sdl.Renderer) {
	if g.travelMap.Visible() {
		renderer.Copy(g.travelMap.bgTex, nil,
			&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})
		g.travelMap.drawOverlay(renderer, g.ui.font, g.mouseX, g.mouseY)
		// Info panel sits on top of the map so the player can keep spatial
		// context while reading. drawInfoPanel is a no-op when no pin is
		// being inspected.
		g.travelMap.drawInfoPanel(renderer, g.ui.font)
		drawVignette(renderer)
		g.ui.drawCursor(renderer, g.mouseX, g.mouseY)
		return
	}

	scene := g.sceneMgr.current()

	scene.drawBackground(renderer, g.player.x)
	showAmbient := g.day == 1 && isCampOutdoorScene(g.sceneMgr.currentName)
	scene.drawAmbient(renderer, showAmbient)
	scene.drawHotspots(renderer, g.ui.hoverName, g.mouseX, g.mouseY)

	// Draw campfire animation in night scene. User 2026-05-19: scale
	// 2.5 → 1.5 (fire was reading way too large vs the camp bg), anchor
	// (622, 573) → (646, 598) - slight shift right + down so the fire
	// sits on the actual fire-pit position in the camp_night.png art.
	if g.sceneMgr.currentName == "camp_night" && len(g.campfireFrames) > 0 {
		f := g.campfireFrames[g.campfireFrameIdx%len(g.campfireFrames)]
		if f.tex != nil {
			fireScale := 1.5
			dstW := int32(float64(f.w) * fireScale)
			dstH := int32(float64(f.h) * fireScale)
			fireX := int32(646) - dstW/2
			fireY := int32(615) - dstH
			renderer.Copy(f.tex, nil, &sdl.Rect{X: fireX, Y: fireY, W: dstW, H: dstH})
		}
	}

	// Draw airplane in flight scene
	if g.sceneMgr.currentName == "airplane_flight" {
		g.flight.Draw(renderer)
	}

	if g.playerSleeping && g.sceneMgr.currentName == "camp_night" {
		// User 2026-05-20: gate the sleep overlay to camp_night only.
		// Without this, the transient frame between the night sequence's
		// `player_sleep true` step and the camp_night transition rendered
		// PP sleeping on top of the marcus_room scene.
		scene.drawActorsNoPlayer(renderer)
		var frames []npcFrame
		if g.wakingPhase == 0 {
			frames = g.sleepingFrames
		} else {
			frames = g.wakingFrames
		}
		if len(frames) > 0 {
			idx := g.sleepingFrameIdx % len(frames)
			f := frames[idx]
			if f.tex != nil {
				// User request 2026-04-17: PP sleeping was huge. 1.1 lands
				// him at a realistic size next to the fire.
				// User 2026-05-17: anchor (335, 591) → (335, 615) so the
				// sleep+wake foot aligns with the new fire Y=615 anchor.
				// User 2026-06-02 (#14): nudge sleeping/waking PP down to 650 so
				// he rests lower by the fire instead of floating above it.
				// User playtest #12: place the sleeping/waking PP at (337,565)
				// (centre-X 337, foot/bottom at Y=565) - "just like before". The
				// wake-up draws at the exact same spot so there's no jump.
				// Size by PP's pixels (opaque box), not the tall 192x1024 cell.
				// One shared scale (tallest wake-up pose ~220px) for both sleeping
				// and waking, anchored bottom-centre at (337,565). (#12)
				// 2026-06-12 #7: match PP's RENDERED idle height - drawScaled
				// fills only playerRenderFillFrac (0.78) of playerDstH, so
				// scaling the sleep pose to the raw 270 made it ~28% bigger
				// than the idle next to it.
				const sleepStandH = float64(playerDstH) * playerRenderFillFrac
				refH := 0.0
				for _, wf := range g.wakingFrames {
					if float64(wf.oh) > refH {
						refH = float64(wf.oh)
					}
				}
				if refH <= 0 {
					refH = float64(f.h)
				}
				scale := sleepStandH / refH
				var src *sdl.Rect
				dstW := int32(float64(f.w) * scale)
				dstH := int32(float64(f.h) * scale)
				if f.ow > 0 && f.oh > 0 {
					s := sdl.Rect{X: f.ox, Y: f.oy, W: f.ow, H: f.oh}
					src = &s
					dstW = int32(float64(f.ow) * scale)
					dstH = int32(float64(f.oh) * scale)
				}
				dstX := int32(337) - dstW/2
				dstY := int32(565) - dstH
				renderer.Copy(f.tex, src, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
			}
		}
	} else if g.nightHidePlayer || g.sceneMgr.currentName == "airplane_flight" {
		// airplane_flight: PP is "inside" the plane sprite, so his
		// standing idle must not render over the cutscene.
		scene.drawActorsNoPlayer(renderer)
	} else {
		scene.drawActors(renderer, g.player)
	}

	// Sequence-owned projectile sprites (thrown map, etc.) render on top
	// of the scene actors but below the dialog / HUD. No-op when no
	// SeqTweenItem step is active. See game/sequence.go:SequencePlayer.Draw.
	g.seqPlayer.Draw(renderer)

	if !(g.sceneMgr.currentName == "camp_night" && g.playerSleeping) {
		drawWarmTint(renderer)
	}
	drawVignette(renderer)

	g.dialog.draw(renderer)
	g.ui.draw(renderer, g.mouseX, g.mouseY)
	g.inv.draw(renderer)
	g.inv.drawHeld(renderer, g.mouseX, g.mouseY)
	g.sceneMgr.drawTransition(renderer)
	// Pause menu goes on top of everything except the cursor.
	g.menu.Draw(renderer, g.ui.font, g.mouseX, g.mouseY)
	// Dev menu sits above the pause menu so it can always be dismissed.
	if g.devMenu != nil {
		g.devMenu.draw(renderer, g.ui.font)
	}
	// Click-probe markers + banner sit on top of the dev menu so the
	// markers stay visible even if the menu was just used to jump scenes.
	if g.clickProbe != nil {
		g.clickProbe.draw(renderer, g.ui.font)
	}
	if g.walkDbg != nil {
		g.walkDbg.draw(renderer, g.ui.font, g.sceneMgr.current(), g.player)
	}
	g.ui.drawCursor(renderer, g.mouseX, g.mouseY)
}

func drawWarmTint(renderer *sdl.Renderer) {
	renderer.SetDrawColor(255, 230, 180, 8)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})
}

func drawVignette(renderer *sdl.Renderer) {
	w := int32(engine.ScreenWidth)
	h := int32(engine.ScreenHeight)

	layers := []struct {
		inset int32
		alpha uint8
	}{
		{0, 22},
		{30, 15},
		{70, 10},
		{120, 5},
	}

	for _, l := range layers {
		renderer.SetDrawColor(0, 0, 0, l.alpha)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: w, H: l.inset + 30})
		renderer.FillRect(&sdl.Rect{X: 0, Y: h - l.inset - 30, W: w, H: l.inset + 30})
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: l.inset + 40, H: h})
		renderer.FillRect(&sdl.Rect{X: w - l.inset - 40, Y: 0, W: l.inset + 40, H: h})
	}
}

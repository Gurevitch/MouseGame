package game

import (
	"fmt"

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
	renderer  *sdl.Renderer
	sceneMgr  *sceneManager
	player    *player
	dialog    *dialogSystem
	ui        *uiManager
	audio     *audioManager
	inv       *inventory
	travelMap *travelMap
	items     *itemRegistry
	dialogs   *dialogStore
	npcDefs   *npcConfigStore
	sceneDefs *sceneConfigStore
	vars      *VarStore
	seqPlayer *SequencePlayer
	seqStore  *sequenceStore
	eventBus  *EventBus
	menu      *gameMenu
	devMenu   *devMenu
	font      *engine.BitmapFont
	lastScene string
	mouseX    int32
	mouseY    int32

	// Story progression
	monologuePlayed bool
	day             int // 1 = arrival/normal, 2 = weirdness begins
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
	playerSleeping    bool
	sleepingFrames    []npcFrame
	sleepingFrameIdx  int
	sleepingTimer     float64
	wakingFrames      []npcFrame
	wakingPhase       int // 0=sleeping, 1=waking, 2=done
	campfireFrames    []npcFrame
	campfireFrameIdx  int
	campfireTimer     float64
	// nightHidePlayer suppresses PP rendering during phase 3 (inside
	// Marcus's cabin) so the cutscene shows only Marcus freaking out,
	// even though PP is technically "present" in the marcus_room scene.
	// playerSleeping alone doesn't cover this: we flipped it to false
	// when transitioning so the sleep sprite wouldn't follow PP into the
	// cabin, but that left the walking PP visible there.
	nightHidePlayer bool

	// Flight cutscene — 4-second biplane transition between cities.
	flight *flightCutscene

	// City monologues
	parisMonologuePlayed bool

	// sceneAltBGs holds pre-loaded alternate backgrounds keyed by
	// "scene_name/variant" (e.g. "marcus_room/day", "marcus_room/night").
	// SeqSetSceneBG looks up through setSceneAltBG; reverting day-mode on
	// save load uses the same map. Replaces what used to be two named
	// *background fields — new alts drop in by adding a load call in
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
	g.audio.playMusic(g.sceneMgr.current().musicPath)

	// Travel Map item: clicking it in the inventory opens the travel map
	// from the current scene. Replaces the camp_entrance / paris_street
	// scene-edge map hotspots (user 2026-04-26 retro-style cleanup).
	g.inv.onSelectItem = func(it *inventoryItem) bool {
		if it == nil || it.name != "Travel Map" {
			return false
		}
		g.travelMap.Show(g.sceneMgr.currentName)
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

	// User reported a white rim on the sleep/wake sprites. The default
	// color-key in SpriteGridFromPNGClean uses tolerance 8 which leaves
	// cream-white halos on these two sheets; aggressive mode (tol 24)
	// + larger inset (4) strips them cleanly.
	sleepGrid := engine.SpriteGridFromPNGCleanAggressive(renderer, "assets/images/player/pp_sleeping.png", 8, 2, 4)
	for c := 0; c < 8; c++ {
		gf := sleepGrid[0][c]
		g.sleepingFrames = append(g.sleepingFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
	}
	wakeGrid := engine.SpriteGridFromPNGCleanAggressive(renderer, "assets/images/player/pp_waking.png", 8, 2, 4)
	for r := 0; r < 2; r++ {
		for c := 0; c < 8; c++ {
			gf := wakeGrid[r][c]
			g.wakingFrames = append(g.wakingFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}

	// User feedback 2026-04-26: switched from the bulky campfire_idle.png
	// (8x4, rows 1-3 had bg drift) to the dedicated campfire_small.png
	// generated for the 2026-04-19 campaign — clean 6x1 grid sized to land
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
		"marcus_room/day":   newPNGBackground(renderer, "assets/images/locations/camp/background/marcus_room_day.png"),
		"marcus_room/night": newPNGBackground(renderer, "assets/images/locations/camp/background/marcus_room_night.png"),
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

	// Make Marcus strange in his room
	if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
		for _, n := range marcusRoom.npcs {
			if n.name == "Marcus" {
				n.dialog = marcusStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
				break
			}
		}
	}

	// Make other kids silent on camp_grounds (they're in their rooms now)
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			n.silent = true
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

	// Restore Marcus room to day background (night cutscene is over)
	if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
		marcusRoom.bg = newPNGBackground(g.renderer, "assets/images/locations/camp/background/marcus_room_day.png")
	}
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
				jake.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Coin Rubbing") {
						return nil, nil
					}
					return []dialogEntry{
						{speaker: "Pink Panther", text: "Jake. I went to Jerusalem. I found your tunnels."},
						{speaker: "Jake", text: "The... the wall? You were THERE?"},
						{speaker: "Pink Panther", text: "Here — a rubbing from one of Miriam's coins. Emperor Hadrian, look."},
						{speaker: "Jake", text: "That's HIM. That's the face in my head. You put him on PAPER!"},
						{speaker: "Jake", text: "The echoes... they're... quieter. Like someone closed a door in my skull."},
						{speaker: "Pink Panther", text: "Rest easy, tough guy. Two down."},
					}, func() {
						game.inv.giveItemTo("Coin Rubbing", "jake")
						jake.setStrange(false)
						game.vars.SetBool(ScopeGame, VarJakeHealed, true)
						jake.dialog = []dialogEntry{
							{speaker: "Jake", text: "Thanks for bringing me the coin. My collection just got LEGENDARY."},
							{speaker: "Jake", text: "But how did I know about that wall before I ever saw it?"},
						}
						game.travelMap.setUnlocked("tokyo_street", true)
						// Wake Lily up in her cabin, ready for chapter 4
						if lRoom, ok := game.sceneMgr.scenes["lily_room"]; ok {
							for _, n := range lRoom.npcs {
								if n.name == "Lily" {
									n.silent = false
									n.setStrange(true)
									break
								}
							}
						}
						game.dialog.queueDialog([]dialogEntry{
							{speaker: "Pink Panther", text: "Lily next. She's been drawing pink petals for a week."},
							{speaker: "Pink Panther", text: "Tokyo just lit up on the map."},
						})
					}
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
				lily.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Pressed Sakura") {
						return nil, nil
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
					}
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
				tommy.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Dance Card") {
						return nil, nil
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
					}
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
				danny.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Inscription Rubbing") {
						return nil, nil
					}
					return []dialogEntry{
						{speaker: "Pink Panther", text: "Danny. Look at the letters I copied from the Roman monument."},
						{speaker: "Danny", text: "That's... that's my NAME. My real name. In Latin."},
						{speaker: "Danny", text: "I've been drawing that arch for weeks because I was trying to draw my own name."},
						{speaker: "Pink Panther", text: "You found yourself in the inscription."},
						{speaker: "Danny", text: "Yeah. I did. All of us — we're spread across the world in pieces."},
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
							{speaker: "Pink Panther", text: "One pin left — Mexico City. That's where this ends."},
						})
					}
				}
				break
			}
		}
	}

	if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
		for _, n := range marcusRoom.npcs {
			if n.name == "Marcus" {
				marcus := n
				marcus.onDialogEnd = func() {
					if game.day == 2 && !game.talkedToMarcus {
						game.talkedToMarcus = true
						marcus.dialog = marcusPostStrangeDialog
					}
				}
				marcus.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Postcard") {
						return nil, nil
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
						marcus.dialog = []dialogEntry{
							{speaker: "Marcus", text: "Thanks for the postcard, counselor. I feel like me again."},
							{speaker: "Marcus", text: "But I still wonder... how did I know about that painting?"},
							{speaker: "Marcus", text: "Go check on the other kids. Something's up with all of us."},
						}
						if mRoom, ok := game.sceneMgr.scenes["marcus_room"]; ok {
							if day, ok := game.sceneAltBGs["marcus_room/day"]; ok {
								mRoom.bg = day
							}
						}
						game.travelMap.setUnlocked("jerusalem_street", true)
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
							{speaker: "Pink Panther", text: "One down. Jake's next — he keeps muttering about tunnels and a wall."},
							{speaker: "Pink Panther", text: "The travel map just lit up Jerusalem. That can't be a coincidence."},
						})
					}
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
						// with the higgins_give_map sequence — Higgins
						// plays his give-map anim, PP plays receive_map,
						// then the item drops into inventory. No more
						// inventory-bar pop on map handover.
						if seq := game.seqStore.Get("higgins_give_map"); seq != nil {
							game.seqPlayer.Play(seq)
						} else {
							game.giveMapItem()
						}
					}
					officeHiggins.dialog = higginsPostWorriedDialog
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
				// started" — the unguarded Day 1 branch let metKids hit 5
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
			kid.onDialogEnd = func() {
				if game.day == 1 && !tommyMet {
					tommyMet = true
					game.metKids++
					kid.dialog = tommyPostDialog
					game.checkDay1Complete()
				} else if game.day >= 2 && !kid.dialogDone {
					kid.dialogDone = true
					kid.dialog = tommyPostStrangeDialog
				}
			}
		case "Jake":
			jakeMet := false
			kid.onDialogEnd = func() {
				if game.day == 1 && !jakeMet {
					jakeMet = true
					game.metKids++
					kid.dialog = jakePostDialog
					game.checkDay1Complete()
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
			kid.altDialogFunc = func() ([]dialogEntry, func()) {
				if kid.hintState != 1 || !game.inv.hasItem("Flower") {
					return nil, nil
				}
				return lilyFlowerDialog, func() {
					game.inv.giveItemTo("Flower", "lily")
					game.metKids++
					kid.hintState = 2
					kid.dialog = lilyDialog
					kid.altDialogFunc = nil
					kid.altDialogRequiresHeld = false
					kid.altDialogRequiresItem = ""
					game.checkDay1Complete()
				}
			}
			kid.onDialogEnd = func() {
				if game.day == 1 && kid.hintState == 0 {
					kid.hintState = 1
					kid.altDialogRequiresHeld = false
					kid.altDialogRequiresItem = "Flower"
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
			tex:     flowerTex,
			srcW:    flowerW,
			srcH:    flowerH,
			bounds:  sdl.Rect{X: 180, Y: 456, W: 50, H: 50},
			name:    "Flower",
			visible: true,
			onPickup: func() {
				item := game.items.createItem("flower")
				if item != nil {
					game.inv.addItem(item)
				}
				// Hide flower in scene
				if lake, ok := game.sceneMgr.scenes["camp_lake"]; ok {
					for _, fi := range lake.floorItems {
						if fi.name == "Flower" {
							fi.visible = false
							break
						}
					}
				}
				game.dialog.startDialog([]dialogEntry{
					{speaker: "Pink Panther", text: "A pretty daisy. I bet Lily would like this."},
				})
			},
		}
		lake.floorItems = append(lake.floorItems, flower)
	}

}

// checkDay1Complete triggers the Day 1 -> Night transition once PP has met
// all 5 kids. Higgins delivers a short "it's getting late" beat on
// camp_grounds first so the fade to night isn't abrupt (user complaint
// 2026-04-17: "higgins didnt say it time to sleep" — the previous version
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
	g.dialog.startDialogWithCallback([]dialogEntry{
		{speaker: "Director Higgins", text: "Ahem! It's getting very late, counselor."},
		{speaker: "Director Higgins", text: "All campers to their cabins. NOW."},
		{speaker: "Pink Panther", text: "Goodnight, Director. I'll turn in by the fire."},
	}, func() {
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

func (g *Game) setupParisCallbacks() {
	// French Guide + the two locals: all swap to post-dialog after first chat.
	if parisStreet, ok := g.sceneMgr.scenes["paris_street"]; ok {
		game := g
		for _, n := range parisStreet.npcs {
			switch n.name {
			case "Madame Colette":
				guide := n
				guide.onDialogEnd = func() {
					guide.dialog = frenchGuidePostDialog
				}
			case "Pierre":
				// Quest step 2: once PP is carrying the baguette, Pierre
				// trades his press pass for it. altDialogFunc fires only
				// when the Baguette is in the bag.
				pierre := n
				pierre.onDialogEnd = func() {
					pierre.dialog = pierreArtistPostDialog
				}
				pierre.altDialogRequiresItem = "Baguette"
				pierre.altDialogRequiresHeld = false
				pierre.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Baguette") || game.inv.hasItem("Press Pass") {
						return nil, nil
					}
					return []dialogEntry{
						{speaker: "Pierre", text: "Mon Dieu! Is that a fresh baguette from Madame Poulain?"},
						{speaker: "Pink Panther", text: "It can be yours, Pierre. I need a favor."},
						{speaker: "Pierre", text: "Anything for bread, mon ami! Take my press pass — it gets you past ze gendarme."},
					}, func() {
						game.inv.giveItemTo("Baguette", "pierre")
						if item := game.items.createItem("press_pass"); item != nil {
							game.inv.addItem(item)
						}
						pierre.altDialogFunc = nil
						pierre.altDialogRequiresItem = ""
					}
				}
			case "Claude":
				// Quest step 3: press pass → museum ticket. Claude waves PP
				// past the queue and hands over the ticket that opens the
				// Louvre entrance hotspot.
				claude := n
				claude.onDialogEnd = func() {
					claude.dialog = gendarmePostDialog
				}
				claude.altDialogRequiresItem = "Press Pass"
				claude.altDialogRequiresHeld = false
				claude.altDialogFunc = func() ([]dialogEntry, func()) {
					if !game.inv.hasItem("Press Pass") || game.inv.hasItem("Museum Ticket") {
						return nil, nil
					}
					return []dialogEntry{
						{speaker: "Claude", text: "A press pass? Ah, press, very well. I zink I have a ticket for ze museum here..."},
						{speaker: "Claude", text: "Pardon ze queue. Follow ze line around ze pyramid."},
						{speaker: "Pink Panther", text: "Merci, Claude. Very kind."},
					}, func() {
						if item := game.items.createItem("museum_ticket"); item != nil {
							game.inv.addItem(item)
						}
						claude.altDialogFunc = nil
						claude.altDialogRequiresItem = ""
					}
				}
			}
		}

		// Louvre entrance gate: needs Museum Ticket.
		for i := range parisStreet.hotspots {
			if parisStreet.hotspots[i].name != "To the Louvre" {
				continue
			}
			h := &parisStreet.hotspots[i]
			h.onInteract = func() bool {
				if !game.inv.hasItem("Museum Ticket") {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Gendarme", text: "Monsieur, you need a ticket to enter the museum."},
						{speaker: "Gendarme", text: "Ask around the street — someone always has a spare."},
					})
					return true
				}
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
				curator.onDialogEnd = func() {
					curator.dialog = museumCuratorPostDialog
					item := game.items.createItem("postcard")
					if item != nil {
						game.inv.addItem(item)
					}
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "A postcard of the painting... this is what Marcus has been drawing."},
						{speaker: "Pink Panther", text: "I need to bring this back to camp. It might help him."},
					})
				}
				break
			}
		}

		// Travel map hotspot to return to camp
		parisLouvre.hotspots = append(parisLouvre.hotspots, hotspot{
			bounds: sdl.Rect{X: 1300, Y: 600, W: 100, H: 150},
			name:   "Travel Map",
			arrow:  arrowDown,
			onInteract: func() bool {
				game.travelMap.Show("paris_louvre")
				return true
			},
		})
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
		game := g
		for _, n := range bakery.npcs {
			if n.name != "Madame Poulain" {
				continue
			}
			poulain := n
			poulain.onDialogEnd = func() {
				// Subsequent clicks while the rolling pin is still missing
				// just replay the lost-pin beat (no flag flip yet).
			}
			poulain.altDialogRequiresItem = "Rolling Pin"
			poulain.altDialogRequiresHeld = false
			poulain.altDialogFunc = func() ([]dialogEntry, func()) {
				if !game.inv.hasItem("Rolling Pin") || game.inv.hasItem("Baguette") {
					return nil, nil
				}
				return bakeryWomanPinTradeDialog, func() {
					game.inv.giveItemTo("Rolling Pin", "madame_poulain")
					if item := game.items.createItem("baguette"); item != nil {
						game.inv.addItem(item)
					}
					poulain.dialog = bakeryWomanPostDialog
					poulain.altDialogFunc = nil
					poulain.altDialogRequiresItem = ""
				}
			}
			break
		}

		// Floor item: rolling pin lying on the bakery floor. Same mechanism
		// as the lake's flower (camp_lake floorItems above). Clickable hit
		// area is intentionally generous so the puzzle reads as findable.
		pinDef, ok := game.items.getDef("rolling_pin")
		if ok {
			pinTex, pinW, pinH := engine.SafeTextureFromPNGKeyed(g.renderer, pinDef.Texture)
			pin := &floorItem{
				tex:     pinTex,
				srcW:    pinW,
				srcH:    pinH,
				bounds:  sdl.Rect{X: 420, Y: 660, W: 90, H: 60},
				name:    "Rolling Pin",
				visible: true,
				onPickup: func() {
					if item := game.items.createItem("rolling_pin"); item != nil {
						game.inv.addItem(item)
					}
					if b, ok := game.sceneMgr.scenes["paris_bakery"]; ok {
						for _, fi := range b.floorItems {
							if fi.name == "Rolling Pin" {
								fi.visible = false
								break
							}
						}
					}
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "A wooden rolling pin. Madame Poulain will be happy."},
					})
				},
			}
			bakery.floorItems = append(bakery.floorItems, pin)
		}
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
							grounds.spawnX = 30
							grounds.spawnY = 490
						}
						game.sceneMgr.transitionTo("camp_grounds", game.player)
					})
					return true
				}
				break
			}
		}
	}

	// Camp Entrance: bus stop to open travel map
	if campEntrance, ok := g.sceneMgr.scenes["camp_entrance"]; ok {
		campEntrance.hotspots = append(campEntrance.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 250, W: 130, H: 350},
			name:   "Camp Chilly Wa Wa Air",
			arrow:  arrowLeft,
			onInteract: func() bool {
				if !game.inv.hasItem("Travel Map") {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "An old airstrip. 'Camp Chilly Wa Wa Air' — how quaint."},
						{speaker: "Pink Panther", text: "I don't have a travel map yet. Maybe Higgins can help."},
					})
					return true
				}
				game.travelMap.Show("camp_entrance")
				return true
			},
		})
	}
}

func (g *Game) Close() {
	g.audio.close()
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

	// Pause menu sits above everything — click routes here first.
	if g.menuHandleClick(x, y) {
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
				g.sceneMgr.transitionTo(loc.scene, g.player)
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
	scene := g.sceneMgr.current()

	if g.player.containsPoint(x, y) {
		if g.inv.heldItem != nil {
			g.inv.toggle()
			return
		}
		if len(g.inv.items) > 0 {
			g.inv.toggle()
			return
		}
	}

	if g.inv.heldItem != nil {
		if clickedNPC := scene.checkNPCClick(x, y); clickedNPC != nil {
			if clickedNPC.altDialogFunc != nil &&
				(clickedNPC.altDialogRequiresItem == "" ||
					g.inv.heldItem.name == clickedNPC.altDialogRequiresItem) {
				entries, cb := clickedNPC.altDialogFunc()
				if entries != nil {
					g.inv.heldItem = nil
					ds := g.dialog
					target := clickedNPC
					g.player.walkToAndInteract(target, ds)
					g.player.interactTarget = nil
				g.player.onArrival = func() {
					g.player.state = stateTalking
					targetCenter := float64(target.bounds.X + target.bounds.W/2)
					playerCenter := g.player.x + playerDstW/2
					g.player.facingLeft = playerCenter > targetCenter
					if g.player.facingLeft {
						g.player.dir = dirLeft
					} else {
						g.player.dir = dirRight
					}
					// Same face-toward-PP flip used by startNPCDialog — the
					// drag-onto-NPC path runs its own dialog start so it
					// needs its own snapshot/restore. Otherwise dropping
					// Lily's flower would leave her still facing into the
					// flower patch even though she's mid-conversation.
					target.preTalkFlipped = target.flipped
					target.flipped = playerCenter < targetCenter
					if len(target.talkGrid) > 0 {
						target.setAnimState(npcAnimTalk)
					}
					wrappedCb := func() {
						if len(target.talkGrid) > 0 {
							target.setAnimState(npcAnimIdle)
						}
						target.flipped = target.preTalkFlipped
						if cb != nil {
							cb()
						}
					}
					ds.startDialogWithCallback(entries, wrappedCb)
				}
					return
				}
			}
		}
		// Travel Map: drop anywhere to open globe
		if g.inv.heldItem.name == "Travel Map" {
			g.inv.heldItem = nil
			g.travelMap.Show(g.sceneMgr.currentName)
			return
		}
		g.inv.heldItem = nil
		return
	}

	if npc := scene.checkNPCClick(x, y); npc != nil {
		g.player.walkToAndInteract(npc, g.dialog)
		return
	}
	if fi := scene.checkFloorItemClick(x, y); fi != nil {
		fiLocal := fi
		g.player.walkToAndDo(
			float64(fiLocal.bounds.X+fiLocal.bounds.W/2),
			float64(fiLocal.bounds.Y+fiLocal.bounds.H/2),
			func() {
				if fiLocal.onPickup != nil {
					fiLocal.onPickup()
				}
			},
		)
		return
	}
	if hs := scene.checkHotspotClick(x, y); hs != nil {
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
		plr := g.player
		sm := g.sceneMgr
		onArrival := func() { sm.transitionTo(tgt, plr) }
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
	tx, ty := float64(x), float64(y)
	tx, ty = scene.snapToPath(tx, ty)
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
		// Menu is modal — eat other input until it closes.
		return
	}
	if scancode == sdl.SCANCODE_M && !g.dialog.active && !g.sceneMgr.transitioning {
		g.travelMap.Toggle(g.sceneMgr.currentName)
		return
	}
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

	// Pause menu freezes game state — only update its own hover and return.
	if g.menu.Visible() {
		g.menu.UpdateHover(mx, my)
		g.ui.cursor = cursorNormal
		g.ui.hoverName = ""
		return
	}

	if g.travelMap.Visible() {
		g.ui.hoverName = ""
		g.ui.cursor = cursorNormal
		if loc := g.travelMap.hitTestAny(mx, my); loc != nil {
			g.ui.hoverName = loc.name
			if loc.unlocked {
				g.ui.cursor = cursorTalk
			}
		}
		return
	}

	scene := g.sceneMgr.current()

	if !g.monologuePlayed && g.sceneMgr.currentName == "camp_entrance" && !g.sceneMgr.transitioning {
		g.monologuePlayed = true
		g.player.state = stateTalking
		g.player.dir = dirDown
		g.dialog.startDialogWithCallback(openingMonologue, func() {
			g.player.state = stateIdle
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

	if g.day2Started && g.sceneMgr.currentName == "camp_grounds" && !g.sceneMgr.transitioning && !g.dialog.active {
		g.day2Started = false
		g.dialog.startDialog(day2Monologue)
	}

	if !g.parisMonologuePlayed && g.sceneMgr.currentName == "paris_street" && !g.sceneMgr.transitioning {
		g.parisMonologuePlayed = true
		g.dialog.startDialog(parisStreetMonologue)
	}

	// Finale: plays once the player lands in Mexico City with every heal flag set.
	if g.sceneMgr.currentName == "mexico_street" && !g.sceneMgr.transitioning && !g.dialog.active {
		g.triggerFinaleMonologue()
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
			var frameCount int
			if g.wakingPhase == 0 {
				frameCount = len(g.sleepingFrames)
			} else {
				frameCount = len(g.wakingFrames)
			}
			if frameCount > 0 {
				g.sleepingFrameIdx = (g.sleepingFrameIdx + 1) % frameCount
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
	g.ui.updateHover(scene, mx, my, g.inv, dt)
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

	// Draw campfire animation in night scene
	if g.sceneMgr.currentName == "camp_night" && len(g.campfireFrames) > 0 {
		f := g.campfireFrames[g.campfireFrameIdx%len(g.campfireFrames)]
		if f.tex != nil {
			fireScale := 2.5
			dstW := int32(float64(f.w) * fireScale)
			dstH := int32(float64(f.h) * fireScale)
			fireX := int32(622) - dstW/2
			fireY := int32(573) - dstH
			renderer.Copy(f.tex, nil, &sdl.Rect{X: fireX, Y: fireY, W: dstW, H: dstH})
		}
	}

	// Draw airplane in flight scene
	if g.sceneMgr.currentName == "airplane_flight" {
		g.flight.Draw(renderer)
	}

	if g.playerSleeping {
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
				// User request 2026-04-17: PP sleeping was huge. 1.8 put
				// him at ~2.5x the campfire; 1.1 lands him at a realistic
				// size next to the fire at (622,573).
				scale := 1.1
				dstW := int32(float64(f.w) * scale)
				dstH := int32(float64(f.h) * scale)
				dstX := int32(335) - dstW/2
				dstY := int32(591) - dstH
				renderer.Copy(f.tex, nil, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
			}
		}
	} else if g.nightHidePlayer || g.sceneMgr.currentName == "airplane_flight" {
		// airplane_flight: PP is "inside" the plane sprite, so his
		// standing idle must not render over the cutscene.
		scene.drawActorsNoPlayer(renderer)
	} else {
		scene.drawActors(renderer, g.player)
	}

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

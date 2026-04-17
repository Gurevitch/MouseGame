package game

import (
	"fmt"
	"math"

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
	lastScene string
	mouseX    int32
	mouseY    int32

	// Story progression
	monologuePlayed bool
	showTravelMap   bool
	travelMapFrom   string
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
	playerSleeping    bool
	sleepingFrames    []npcFrame
	sleepingFrameIdx  int
	sleepingTimer     float64
	wakingFrames      []npcFrame
	wakingPhase       int // 0=sleeping, 1=waking, 2=done
	campfireFrames    []npcFrame
	campfireFrameIdx  int
	campfireTimer     float64
	nightPhase            int     // 0=not started, 1=higgins speech, 2=sleeping+marcus audio, 3=marcus room, 4=wake, 5=day2
	nightTimer            float64
	marcusFreakoutStarted bool
	// nightHidePlayer suppresses PP rendering during phase 3 (inside
	// Marcus's cabin) so the cutscene shows only Marcus freaking out,
	// even though PP is technically "present" in the marcus_room scene.
	// playerSleeping alone doesn't cover this: we flipped it to false
	// when transitioning so the sleep sprite wouldn't follow PP into the
	// cabin, but that left the walking PP visible there.
	nightHidePlayer bool

	// Map reveal animation
	mapRevealing  bool
	mapRevealScale float64
	mapRevealTex   *sdl.Texture
	mapRevealW     int32
	mapRevealH     int32

	// Flight cutscene
	flightDestination  string
	flightTimer        float64
	airplaneFrames     []npcFrame
	airplaneFrameIdx   int
	airplaneFrameTimer float64

	// City monologues
	parisMonologuePlayed bool

	marcusRoomBg      *background
	marcusRoomNightBg *background
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

	g.travelMap = newTravelMap(renderer)
	g.items = newItemRegistry(renderer, "assets/data/items.json")
	g.dialogs = newDialogStore("assets/data/dialog")
	g.npcDefs = newNPCConfigStore("assets/data/npc")
	g.sceneDefs = newSceneConfigStore("assets/data/scenes")
	g.vars = newVarStore()
	g.vars.Set(ScopeGame, VarChapter, ChapterCampDay1)
	g.vars.Set(ScopeGame, VarDay, 1)
	g.seqPlayer = newSequencePlayer(g)
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

	// Campfire sheet is authored as 8x4 but in practice rows 1-3 vary
	// wildly in background tint — flipbooking all 32 frames created
	// strobing fringes and a visible halo at draw scale 2.5. We use the
	// aggressive color-key variant (wider tolerance) and only row 0,
	// giving us a clean 8-frame flame loop without background bleed.
	// Inset is bumped to 4 so the outer fringe pixels never survive the
	// trim.
	fireGrid := engine.SpriteGridFromPNGCleanAggressive(renderer, "assets/images/locations/camp/campfire_idle.png", 8, 4, 4)
	if len(fireGrid) > 0 {
		for c := 0; c < len(fireGrid[0]); c++ {
			gf := fireGrid[0][c]
			g.campfireFrames = append(g.campfireFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}

	airGrid := engine.SpriteGridFromPNGClean(renderer, "assets/images/player/pp_airplane.png", 4, 3, spriteInset)
	for r := 0; r < len(airGrid); r++ {
		for c := 0; c < len(airGrid[r]); c++ {
			gf := airGrid[r][c]
			g.airplaneFrames = append(g.airplaneFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}

	g.marcusRoomBg = newPNGBackground(renderer, "assets/images/locations/camp/background/marcus_room_day.png")
	g.marcusRoomNightBg = newPNGBackground(renderer, "assets/images/locations/camp/background/marcus_room_night.png")

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
						if mRoom, ok := game.sceneMgr.scenes["marcus_room"]; ok && game.marcusRoomBg != nil {
							mRoom.bg = game.marcusRoomBg
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
						game.giveMapItem()
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
				kid.onDialogEnd = func() {
					if game.day == 1 {
						game.metKids++
						kid.dialog = marcusPostDialog
						game.checkDay1Complete()
					} else {
						if !game.talkedToMarcus {
							game.talkedToMarcus = true
						}
						kid.dialog = marcusPostStrangeDialog
					}
				}
		case "Tommy":
			kid.onDialogEnd = func() {
				if game.day == 1 {
					game.metKids++
					kid.dialog = tommyPostDialog
					game.checkDay1Complete()
				} else if !kid.dialogDone {
					kid.dialogDone = true
					kid.dialog = tommyPostStrangeDialog
				}
			}
		case "Jake":
			kid.onDialogEnd = func() {
				if game.day == 1 {
					game.metKids++
					kid.dialog = jakePostDialog
					game.checkDay1Complete()
				} else if !kid.dialogDone {
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
					// Reveal the camp-grounds Higgins so he can deliver
					// the flower clue. He stays visible for the rest of
					// Day 1 so the player can re-ask for the hint.
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
				} else if game.day >= 2 && !kid.dialogDone {
					kid.dialogDone = true
					kid.dialog = lilyPostStrangeDialog
				}
			}
		case "Danny":
			kid.onDialogEnd = func() {
				if game.day == 1 {
					game.metKids++
					kid.dialog = dannyPostDialog
					game.checkDay1Complete()
				} else if !kid.dialogDone {
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

// nightSceneArrival kicks off the campfire cutscene. Both actors are already
// placed by the scene definition (PP at spawn 335,591 — same coords used by
// the sleeping sprite — and night Higgins at 820,420). We just start phase 1,
// which runs Higgins' "get some rest" speech before PP lies down.
func (g *Game) nightSceneArrival() {
	if g.nightSceneDone {
		return
	}
	g.nightSceneDone = true
	g.nightPhase = 1
	g.nightTimer = 0
	g.playerSleeping = false
	g.wakingPhase = 0
}

// nightSceneUpdate drives the 6-phase campfire cutscene. Each phase is a small
// state-machine step; when a phase ends (dialog finishes or timer elapses) it
// flips to the next and zeroes the timer.
//
//   1. Higgins speaks at the fire, PP stands.
//   2. PP lies down; short beat; we only HEAR Marcus (no scene change yet).
//   3. Transition to Marcus's room (night bg, strange state) and see him.
//   4. Back to the campfire, still sleeping; PP wakes up.
//   5. PP Day-2 monologue; switch day; transition to camp_grounds.
func (g *Game) nightSceneUpdate(dt float64) {
	if g.nightPhase == 0 || g.nightPhase >= 6 {
		return
	}
	g.nightTimer += dt

	switch g.nightPhase {
	case 1:
		if g.nightTimer < 0.6 || g.dialog.active {
			return
		}
		higgins := g.findNightHiggins()
		if higgins != nil {
			higgins.setAnimState(npcAnimTalk)
		}
		g.dialog.startDialogWithCallback([]dialogEntry{
			{speaker: "Director Higgins", text: "Ahem! It's getting very late, counselor."},
			{speaker: "Director Higgins", text: "All campers to their cabins. NOW."},
			{speaker: "Director Higgins", text: "And you — get some rest by the fire. Big day tomorrow."},
			{speaker: "Pink Panther", text: "Goodnight, Director."},
		}, func() {
			if h := g.findNightHiggins(); h != nil {
				h.setAnimState(npcAnimIdle)
			}
			g.nightPhase = 2
			g.nightTimer = 0
			g.playerSleeping = true
			g.wakingPhase = 0
			g.sleepingFrameIdx = 0
			g.sleepingTimer = 0
		})

	case 2:
		if g.nightTimer < 3.0 || g.dialog.active {
			return
		}
		g.dialog.startDialogWithCallback([]dialogEntry{
			{speaker: "Marcus", text: "*from Marcus's cabin, muffled* No no no... the lines won't stop..."},
			{speaker: "Marcus", text: "A GLASS PYRAMID! The building is ENORMOUS!"},
			{speaker: "Marcus", text: "The painting is WRONG! Something is MISSING!"},
			{speaker: "Pink Panther", text: "*sleepily* Marcus...? That doesn't sound right..."},
		}, func() {
			g.playerSleeping = false
			g.nightHidePlayer = true
			g.nightPhase = 3
			g.nightTimer = 0
			g.marcusFreakoutStarted = false
			if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok && g.marcusRoomNightBg != nil {
				marcusRoom.bg = g.marcusRoomNightBg
			}
			if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
				for _, n := range marcusRoom.npcs {
					if n.name == "Marcus" {
						n.setStrange(true)
						break
					}
				}
			}
			g.sceneMgr.transitionTo("marcus_room", g.player)
		})

	case 3:
		if g.sceneMgr.transitioning || g.marcusFreakoutStarted {
			return
		}
		g.marcusFreakoutStarted = true
		g.dialog.startDialogWithCallback([]dialogEntry{
			{speaker: "Marcus", text: "I can't stop... the lines keep coming..."},
			{speaker: "Marcus", text: "A woman's face... golden frames everywhere..."},
			{speaker: "Marcus", text: "I have to draw it... I HAVE to draw it ALL!"},
			{speaker: "Pink Panther", text: "*whispers* He's not even awake... Higgins needs to know."},
		}, func() {
			g.nightHidePlayer = false
			g.nightPhase = 4
			g.nightTimer = 0
			g.playerSleeping = true
			g.wakingPhase = 0
			g.sleepingFrameIdx = 0
			g.sceneMgr.transitionTo("camp_night", g.player)
		})

	case 4:
		if g.sceneMgr.transitioning {
			return
		}
		if g.wakingPhase == 0 && g.nightTimer >= 1.0 {
			g.wakingPhase = 1
			g.sleepingFrameIdx = 0
			g.sleepingTimer = 0
			g.nightTimer = 0
			return
		}
		if g.wakingPhase == 1 && g.nightTimer >= 2.0 {
			g.wakingPhase = 2
			g.playerSleeping = false
			g.nightPhase = 5
			g.nightTimer = 0
			g.startDay2()
			g.dialog.startDialogWithCallback([]dialogEntry{
				{speaker: "Pink Panther", text: "*yawn* What a night..."},
				{speaker: "Pink Panther", text: "I heard Marcus freaking out in his cabin. Something about paintings and pyramids."},
				{speaker: "Pink Panther", text: "I need to check on him. His cabin is in the camp grounds."},
			}, func() {
				g.day2Started = true
				g.sceneMgr.transitionTo("camp_grounds", g.player)
			})
		}
	}
}

// findNightHiggins returns the silent night-campfire Higgins NPC, if present.
// Used by the cutscene so Higgins' talk animation syncs with his dialog.
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
		for _, n := range parisStreet.npcs {
			switch n.name {
			case "Madame Colette":
				guide := n
				guide.onDialogEnd = func() {
					guide.dialog = frenchGuidePostDialog
				}
			case "Pierre":
				pierre := n
				pierre.onDialogEnd = func() {
					pierre.dialog = pierreArtistPostDialog
				}
			case "Claude":
				claude := n
				claude.onDialogEnd = func() {
					claude.dialog = gendarmePostDialog
				}
			}
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
				game.showTravelMap = true
				game.travelMapFrom = "paris_louvre"
				return true
			},
		})
	}

	// Paris street: travel map access
	if parisStreet, ok := g.sceneMgr.scenes["paris_street"]; ok {
		game := g
		parisStreet.hotspots = append(parisStreet.hotspots, hotspot{
			bounds: sdl.Rect{X: 0, Y: 200, W: 100, H: 400},
			name:   "Travel Map",
			arrow:  arrowLeft,
			onInteract: func() bool {
				game.showTravelMap = true
				game.travelMapFrom = "paris_street"
				return true
			},
		})
	}
}

func (g *Game) setupTravelHotspots() {
	game := g

	if campEntrance, ok := g.sceneMgr.scenes["camp_entrance"]; ok {
		for i := range campEntrance.hotspots {
			if campEntrance.hotspots[i].name == "Enter Camp" {
				campEntrance.hotspots[i].onInteract = func() bool {
					game.player.dir = dirUp
					game.player.allowOffscreen = true
					game.player.walkToAndDo(599, 200, func() {
						game.player.allowOffscreen = false
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
				game.showTravelMap = true
				game.travelMapFrom = "camp_entrance"
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

	if g.showTravelMap {
		if loc := g.travelMap.hitTest(x, y); loc != nil {
			if loc.scene == g.travelMapFrom {
				g.showTravelMap = false
				return
			}
			g.showTravelMap = false
			if loc.scene == "camp_entrance" {
				g.sceneMgr.transitionTo(loc.scene, g.player)
				return
			}
			g.flightDestination = loc.scene
			g.flightTimer = 0
			g.airplaneFrameIdx = 0
			g.sceneMgr.transitionTo("airplane_flight", g.player)
			return
		}

		if anyLoc := g.travelMap.hitTestAny(x, y); anyLoc != nil && !anyLoc.unlocked && anyLoc.info != "" {
			g.showTravelMap = false
			g.dialog.startDialog([]dialogEntry{
				{speaker: anyLoc.name, text: anyLoc.info},
			})
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
			g.showTravelMap = true
			g.travelMapFrom = g.sceneMgr.currentName
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
	if g.showTravelMap && scancode == sdl.SCANCODE_ESCAPE {
		g.showTravelMap = false
		return
	}
	if scancode == sdl.SCANCODE_M && !g.dialog.active && !g.sceneMgr.transitioning {
		g.showTravelMap = !g.showTravelMap
		if g.showTravelMap {
			g.travelMapFrom = g.sceneMgr.currentName
		}
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

	if g.showTravelMap {
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

	if !g.nightSceneDone && g.day == 1 && g.metKids >= 5 && g.sceneMgr.currentName == "camp_night" && !g.sceneMgr.transitioning && !g.dialog.active {
		g.nightSceneArrival()
	}

	// Multi-phase night cutscene
	g.nightSceneUpdate(dt)

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
	if g.sceneMgr.currentName == "airplane_flight" && !g.sceneMgr.transitioning && g.flightDestination != "" {
		g.flightTimer += dt
		g.airplaneFrameTimer += dt
		if g.airplaneFrameTimer >= 0.12 && len(g.airplaneFrames) > 0 {
			g.airplaneFrameTimer -= 0.12
			g.airplaneFrameIdx = (g.airplaneFrameIdx + 1) % len(g.airplaneFrames)
		}
		if g.flightTimer >= 4.0 {
			dest := g.flightDestination
			g.flightDestination = ""
			g.flightTimer = 0
			g.sceneMgr.transitionTo(dest, g.player)
		}
	}

	// Map reveal animation
	if g.mapRevealing {
		g.mapRevealScale += dt * 1.5 // ~0.67 seconds to full
		if g.mapRevealScale >= 1.0 {
			g.mapRevealScale = 1.0
			g.mapRevealing = false
			g.dialog.startDialog([]dialogEntry{
				{speaker: "Pink Panther", text: "A travel map! Now I can use Camp Chilly Wa Wa Air."},
				{speaker: "Pink Panther", text: "I just need to take it from my inventory and use it anywhere."},
			})
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
	if g.showTravelMap {
		renderer.Copy(g.travelMap.bgTex, nil,
			&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})
		g.travelMap.drawOverlay(renderer, g.ui.font, g.mouseX, g.mouseY)
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
	if g.sceneMgr.currentName == "airplane_flight" && len(g.airplaneFrames) > 0 {
		f := g.airplaneFrames[g.airplaneFrameIdx%len(g.airplaneFrames)]
		if f.tex != nil {
			bob := math.Sin(g.flightTimer*2.0) * 8
			scale := 3.0
			dstW := int32(float64(f.w) * scale)
			dstH := int32(float64(f.h) * scale)
			dstX := engine.ScreenWidth/2 - dstW/2
			dstY := int32(float64(engine.ScreenHeight)/2 - float64(dstH)/2 + bob)
			renderer.Copy(f.tex, nil, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
		}
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

	// Map reveal animation overlay
	if g.mapRevealing && g.mapRevealTex != nil {
		// Dim background
		renderer.SetDrawColor(0, 0, 0, uint8(140*g.mapRevealScale))
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})

		scale := g.mapRevealScale
		dstW := int32(float64(g.mapRevealW) * scale * 3)
		dstH := int32(float64(g.mapRevealH) * scale * 3)
		dstX := engine.ScreenWidth/2 - dstW/2
		dstY := engine.ScreenHeight/2 - dstH/2
		renderer.Copy(g.mapRevealTex, nil, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
	}

	g.dialog.draw(renderer)
	g.ui.draw(renderer, g.mouseX, g.mouseY)
	g.inv.draw(renderer)
	g.inv.drawHeld(renderer, g.mouseX, g.mouseY)
	g.sceneMgr.drawTransition(renderer)
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

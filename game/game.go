package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

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

type Game struct {
	renderer  *sdl.Renderer
	sceneMgr  *sceneManager
	player    *player
	dialog    *dialogSystem
	ui        *uiManager
	audio     *audioManager
	inv       *inventory
	travelMap *travelMap
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

	// Map reveal animation
	mapRevealing  bool
	mapRevealScale float64
	mapRevealTex   *sdl.Texture
	mapRevealW     int32
	mapRevealH     int32

	// City monologues
	parisMonologuePlayed bool

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
	g.lastScene = g.sceneMgr.currentName
	g.audio.playMusic(g.sceneMgr.current().musicPath)

	g.travelMap = newTravelMap(renderer)
	g.setupCampCallbacks()
	g.setupParisCallbacks()
	g.setupTravelHotspots()
	g.ui.initCursors(renderer)

	// Load sleeping/waking sprites (use first row only)
	sleepGrid := engine.SpriteGridFromPNG(renderer, "assets/images/player/pp_sleeping.png", 8, 2)
	for c := 0; c < 8; c++ {
		gf := sleepGrid[0][c]
		g.sleepingFrames = append(g.sleepingFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
	}
	wakeGrid := engine.SpriteGridFromPNG(renderer, "assets/images/player/pp_waking.png", 8, 2)
	for r := 0; r < 2; r++ {
		for c := 0; c < 8; c++ {
			gf := wakeGrid[r][c]
			g.wakingFrames = append(g.wakingFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}

	fireGrid := engine.SpriteGridFromPNGRaw(renderer, "assets/images/locations/camp/campfire_idle.png", 8, 4)
	for r := 0; r < len(fireGrid); r++ {
		for c := 0; c < len(fireGrid[r]); c++ {
			gf := fireGrid[r][c]
			g.campfireFrames = append(g.campfireFrames, npcFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}

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

	// --- Marcus in his room (Day 2) ---
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
			kid.onDialogEnd = func() {
				if game.day == 1 && !kid.dialogDone {
					// First talk: shy dialog — Higgins hint
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Director Higgins", text: "She's very shy, that one."},
						{speaker: "Director Higgins", text: "Maybe try bringing her something she likes? She loves flowers."},
						{speaker: "Pink Panther", text: "Flowers... I think I saw some by the lake."},
					})
					kid.dialogDone = true // mark so hint only plays once
				} else if game.day >= 2 && !kid.dialogDone {
					kid.dialogDone = true
					kid.dialog = lilyPostStrangeDialog
				}
			}
			// Flower item interaction
			kid.altDialogFunc = func() ([]dialogEntry, func()) {
				return lilyFlowerDialog, func() {
					game.metKids++
					kid.dialog = lilyDialog
					kid.altDialogFunc = nil
					game.checkDay1Complete()
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
		flowerTex, flowerW, flowerH := engine.SafeTextureFromPNGRaw(g.renderer, "assets/images/items/flower.png")
		flower := &floorItem{
			tex:     flowerTex,
			srcW:    flowerW,
			srcH:    flowerH,
			bounds:  sdl.Rect{X: 180, Y: 456, W: 50, H: 50},
			name:    "Flower",
			visible: true,
			onPickup: func() {
				if flowerTex != nil {
					game.inv.addItem(&inventoryItem{
						name: "Flower",
						tex:  flowerTex,
						srcW: flowerW,
						srcH: flowerH,
						desc: "A pretty daisy from the lake. Lily might like this.",
					})
				}
				// Remove flower from scene
				if lake, ok := game.sceneMgr.scenes["camp_lake"]; ok {
					for i, fi := range lake.floorItems {
						if fi.name == "Flower" {
							fi.visible = false
							_ = i
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

// checkDay1Complete triggers Day 2 once PP has met all 5 kids
func (g *Game) checkDay1Complete() {
	if g.metKids >= 5 && g.day == 1 {
		g.dialog.startDialogWithCallback([]dialogEntry{
			{speaker: "Director Higgins", text: "Ahem! It's getting late, counselor."},
			{speaker: "Director Higgins", text: "The kids need to be in their cabins. And so do you!"},
			{speaker: "Pink Panther", text: "Already? Time flies when you're meeting interesting people."},
			{speaker: "Director Higgins", text: "Off to bed! Big day tomorrow."},
			{speaker: "Pink Panther", text: "Goodnight, Higgins."},
		}, func() {
			// Swap Marcus room to night background before entering
			if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok && g.marcusRoomNightBg != nil {
				marcusRoom.bg = g.marcusRoomNightBg
			}
			g.sceneMgr.transitionTo("marcus_room", g.player)
		})
	}
}

// nightSceneArrival — automatic cutscene: shows Marcus freaking out in his room
func (g *Game) nightSceneArrival() {
	if g.nightSceneDone {
		return
	}
	g.nightSceneDone = true
	g.playerSleeping = true // hide normal player during cutscene

	// Show Marcus's room at night with his frantic dialog
	g.dialog.startDialogWithCallback([]dialogEntry{
		{speaker: "Marcus", text: "I can't stop... the lines keep coming..."},
		{speaker: "Marcus", text: "A face... golden frames... a glass pyramid..."},
		{speaker: "Marcus", text: "I don't KNOW! I've never seen these things!"},
		{speaker: "Marcus", text: "A woman... she's smiling... and the building is ENORMOUS!"},
		{speaker: "Marcus", text: "The picture is WRONG! Something is MISSING!"},
	}, func() {
		// Morning: transition to camp_grounds, player must find Marcus
		g.playerSleeping = false
		g.startDay2()
		g.sceneMgr.transitionTo("camp_grounds", g.player)
		g.day2Started = true
	})
}

func (g *Game) giveMapItem() {
	if g.inv.hasItem("Travel Map") {
		return
	}
	tex, w, h := engine.TextureFromPNG(g.renderer, "assets/images/items/travel_map_icon.png")
	if tex != nil {
		tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	// Start map reveal animation
	g.mapRevealing = true
	g.mapRevealScale = 0.0
	g.mapRevealTex = tex
	g.mapRevealW = w
	g.mapRevealH = h

	g.inv.addItem(&inventoryItem{
		name: "Travel Map",
		tex:  tex,
		srcW: w,
		srcH: h,
		desc: "Use the map to travel by dropping it anywhere on the screen.",
	})
}

func (g *Game) setupParisCallbacks() {
	// French Guide: post-dialog after first conversation
	if parisStreet, ok := g.sceneMgr.scenes["paris_street"]; ok {
		for _, n := range parisStreet.npcs {
			if n.name == "Madame Colette" {
				guide := n
				guide.onDialogEnd = func() {
					guide.dialog = frenchGuidePostDialog
				}
				break
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
					// Give the postcard anchor item
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

	// Camp Entrance: "Enter Camp" — walk up the road, shrink, transition
	if campEntrance, ok := g.sceneMgr.scenes["camp_entrance"]; ok {
		for i := range campEntrance.hotspots {
			if campEntrance.hotspots[i].name == "Enter Camp" {
				campEntrance.hotspots[i].onInteract = func() bool {
					game.player.dir = dirUp
					game.player.allowOffscreen = true
					game.player.walkToAndDo(599, 200, func() {
						game.player.allowOffscreen = false
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
		loc := g.travelMap.hitTest(x, y)
		if loc != nil && loc.scene != g.travelMapFrom {
			g.showTravelMap = false
			g.sceneMgr.transitionTo(loc.scene, g.player)
		} else if loc == nil {
			// Check if they clicked a locked city — show info popup
			anyLoc := g.travelMap.hitTestAny(x, y)
			if anyLoc != nil && !anyLoc.unlocked && anyLoc.info != "" {
				g.showTravelMap = false
				g.dialog.startDialog([]dialogEntry{
					{speaker: anyLoc.name, text: anyLoc.info},
				})
			} else {
				g.showTravelMap = false
			}
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
			if clickedNPC.altDialogFunc != nil {
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
						if len(target.talkGrid) > 0 {
							target.setAnimState(npcAnimTalk)
						}
						wrappedCb := func() {
							if len(target.talkGrid) > 0 {
								target.setAnimState(npcAnimIdle)
							}
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
}

func (g *Game) Update(dt float64, mx, my int32) {
	g.mouseX = mx
	g.mouseY = my

	if g.showTravelMap {
		return
	}

	scene := g.sceneMgr.current()

	if !g.monologuePlayed && g.sceneMgr.currentName == "camp_entrance" && !g.sceneMgr.transitioning {
		g.monologuePlayed = true
		// PP faces camera and talks during opening monologue
		g.player.state = stateTalking
		g.player.dir = dirDown
		g.dialog.startDialogWithCallback(openingMonologue, func() {
			g.player.state = stateIdle
			// Auto-walk to Higgins after monologue
			for _, n := range scene.npcs {
				if n.name == "Director Higgins" {
					g.player.walkToAndInteract(n, g.dialog)
					break
				}
			}
		})
	}

	if !g.nightSceneDone && g.day == 1 && g.metKids >= 5 && g.sceneMgr.currentName == "marcus_room" && !g.sceneMgr.transitioning && !g.dialog.active {
		g.nightSceneArrival()
	}

	if g.day2Started && g.sceneMgr.currentName == "camp_grounds" && !g.sceneMgr.transitioning && !g.dialog.active {
		g.day2Started = false
		g.dialog.startDialog(day2Monologue)
	}

	if !g.parisMonologuePlayed && g.sceneMgr.currentName == "paris_street" && !g.sceneMgr.transitioning {
		g.parisMonologuePlayed = true
		g.dialog.startDialog(parisStreetMonologue)
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
	g.sceneMgr.update(dt)
	scene.updateAmbient(dt)
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
	scene.drawAmbient(renderer)
	scene.drawHotspots(renderer, g.ui.hoverName, g.mouseX, g.mouseY)

	// Draw campfire animation in night scene
	if g.sceneMgr.currentName == "camp_night" && len(g.campfireFrames) > 0 {
		f := g.campfireFrames[g.campfireFrameIdx%len(g.campfireFrames)]
		if f.tex != nil {
			fireScale := 2.5
			dstW := int32(float64(f.w) * fireScale)
			dstH := int32(float64(f.h) * fireScale)
			fireX := int32(600) - dstW/2
			fireY := int32(450) - dstH
			renderer.Copy(f.tex, nil, &sdl.Rect{X: fireX, Y: fireY, W: dstW, H: dstH})
		}
	}

	if g.playerSleeping {
		// Draw sleeping/waking sprite instead of player
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
				scale := 1.5
				dstW := int32(float64(f.w) * scale)
				dstH := int32(float64(f.h) * scale)
				dstX := int32(500) - dstW/2
				dstY := int32(420) - dstH
				renderer.Copy(f.tex, nil, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
			}
		}
	} else {
		scene.drawActors(renderer, g.player)
	}

	drawWarmTint(renderer)
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

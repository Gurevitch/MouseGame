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
	{speaker: "Pink Panther", text: "Something feels different this morning..."},
	{speaker: "Pink Panther", text: "The kids... they don't seem like themselves."},
	{speaker: "Pink Panther", text: "I should talk to everyone and find out what's going on."},
}

var parisStreetMonologue = []dialogEntry{
	{speaker: "Pink Panther", text: "Ah, Paris! The city of lights, love, and... mysteries, apparently."},
	{speaker: "Pink Panther", text: "Marcus kept drawing a museum with a glass pyramid. That must be the Louvre."},
	{speaker: "Pink Panther", text: "Time to find out what he's been seeing."},
}

type Game struct {
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

	// City monologues
	parisMonologuePlayed bool
}

func New(renderer *sdl.Renderer, font *engine.BitmapFont) *Game {
	g := &Game{
		sceneMgr: newSceneManager(renderer),
		player:   newPlayer(renderer),
		dialog:   newDialogSystem(font),
		ui:       newUIManager(font),
		audio:    newAudioManager(),
		inv:      newInventory(font),
		day:      1,
	}
	g.lastScene = g.sceneMgr.currentName
	g.audio.playMusic(g.sceneMgr.current().musicPath)

	g.travelMap = newTravelMap(renderer)
	g.setupCampCallbacks()
	g.setupParisCallbacks()
	g.setupTravelHotspots()
	g.ui.initCursors(renderer)

	return g
}

// startDay2 transitions all NPCs to their "strange" dialogs
func (g *Game) startDay2() {
	if g.day >= 2 {
		return
	}
	g.day = 2

	// Swap all kid dialogs and sprites to strange versions
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			switch n.name {
			case "Marcus":
				n.dialog = marcusStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
			case "Tommy":
				n.dialog = tommyStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
			case "Jake":
				n.dialog = jakeStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
			case "Lily":
				n.dialog = lilyStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
			case "Danny":
				n.dialog = dannyStrangeDialog
				n.dialogDone = false
				n.setStrange(true)
			}
		}
	}

	// Swap Higgins dialog
	if entrance, ok := g.sceneMgr.scenes["camp_entrance"]; ok {
		for _, n := range entrance.npcs {
			if n.name == "Director Higgins" {
				n.dialog = higginsWorriedDialog
				n.dialogDone = false
				break
			}
		}
	}

	// Swap Cook Marge dialog
	if messHall, ok := g.sceneMgr.scenes["camp_messhall"]; ok {
		for _, n := range messHall.npcs {
			if n.name == "Cook Marge" {
				n.dialog = cookMargeWorriedDialog
				n.dialogDone = false
				break
			}
		}
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
				} else {
					higgins.dialog = higginsPostWorriedDialog
				}
			}
			break
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
						// Day 2: talked to strange Marcus -> unlock Paris
						if !game.talkedToMarcus {
							game.talkedToMarcus = true
							kid.dialog = marcusPostStrangeDialog
							game.tryUnlockParis()
						}
					}
				}
			case "Tommy":
				kid.onDialogEnd = func() {
					if game.day == 1 {
						game.metKids++
						kid.dialog = tommyPostDialog
						game.checkDay1Complete()
					} else {
						kid.dialog = tommyPostStrangeDialog
					}
				}
			case "Jake":
				kid.onDialogEnd = func() {
					if game.day == 1 {
						game.metKids++
						kid.dialog = jakePostDialog
						game.checkDay1Complete()
					} else {
						kid.dialog = jakePostStrangeDialog
					}
				}
			case "Lily":
				kid.onDialogEnd = func() {
					if game.day == 1 {
						game.metKids++
						kid.dialog = lilyPostDialog
						game.checkDay1Complete()
					} else {
						kid.dialog = lilyPostStrangeDialog
					}
				}
			case "Danny":
				kid.onDialogEnd = func() {
					if game.day == 1 {
						game.metKids++
						kid.dialog = dannyPostDialog
						game.checkDay1Complete()
					} else {
						kid.dialog = dannyPostStrangeDialog
					}
				}
			}
		}
	}

	// --- Cook Marge ---
	if messHall, ok := g.sceneMgr.scenes["camp_messhall"]; ok {
		for _, n := range messHall.npcs {
			if n.name == "Cook Marge" {
				marge := n
				marge.onDialogEnd = func() {
					if game.day == 1 {
						marge.dialog = cookMargePostDialog
					} else {
						marge.dialog = cookMargePostWorriedDialog
					}
				}
				break
			}
		}
	}
}

// checkDay1Complete triggers Day 2 once PP has met all 5 kids
func (g *Game) checkDay1Complete() {
	if g.metKids >= 5 && g.day == 1 {
		// All kids met! Trigger Day 2 transition
		g.dialog.startDialogWithCallback([]dialogEntry{
			{speaker: "Pink Panther", text: "Well, I've met everyone. Seems like a nice group."},
			{speaker: "Pink Panther", text: "Time to get some rest. Tomorrow is a new day."},
		}, func() {
			g.startDay2()
			// Show Day 2 morning monologue
			g.dialog.startDialog(day2Monologue)
		})
	}
}

func (g *Game) tryUnlockParis() {
	if g.parisUnlocked {
		return
	}
	g.parisUnlocked = true
	g.travelMap.setUnlocked("paris_street", true)
	g.dialog.startDialog([]dialogEntry{
		{speaker: "Pink Panther", text: "A glass pyramid... the biggest museum in the world... a woman's face..."},
		{speaker: "Pink Panther", text: "Marcus is seeing the Louvre. In Paris."},
		{speaker: "Pink Panther", text: "I need to go there and find out what he's connected to."},
		{speaker: "Pink Panther", text: "The travel map should have Paris available now."},
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

	// Camp Entrance: bus stop to open travel map
	if campEntrance, ok := g.sceneMgr.scenes["camp_entrance"]; ok {
		campEntrance.hotspots = append(campEntrance.hotspots, hotspot{
			bounds: sdl.Rect{X: 120, Y: 250, W: 130, H: 200},
			name:   "Camp Chilly Wa Wa Air",
			arrow:  arrowLeft,
			onInteract: func() bool {
				if game.day < 2 {
					game.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "An old airstrip. 'Camp Chilly Wa Wa Air' — how quaint."},
						{speaker: "Pink Panther", text: "I don't think I need to go anywhere just yet."},
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
			g.showTravelMap = false
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
						if len(target.talkFrames) > 0 {
							target.setAnimState(npcAnimTalk)
						}
						wrappedCb := func() {
							if len(target.talkFrames) > 0 {
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
		if hs.arrow == arrowLeft || hs.arrow == arrowRight || hs.arrow == arrowDown {
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
	g.player.setTarget(float64(x), float64(y))
}

func (g *Game) HandleKey(scancode sdl.Scancode) {
	if g.showTravelMap && scancode == sdl.SCANCODE_ESCAPE {
		g.showTravelMap = false
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
		g.dialog.startDialog(openingMonologue)
	}

	if !g.parisMonologuePlayed && g.sceneMgr.currentName == "paris_street" && !g.sceneMgr.transitioning {
		g.parisMonologuePlayed = true
		g.dialog.startDialog(parisStreetMonologue)
	}

	if !g.dialog.active && g.player.state == stateTalking {
		g.player.state = stateIdle
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
	scene.drawActors(renderer, g.player)

	drawWarmTint(renderer)
	drawVignette(renderer)

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

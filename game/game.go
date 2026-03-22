package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

var openingMonologue = []dialogEntry{
	{speaker: "Pink Panther", text: "So many things have changed since the last time I was here..."},
	{speaker: "Pink Panther", text: "The old sign is barely standing, the cabins look like they've seen better days..."},
	{speaker: "Pink Panther", text: "...and is that a llama?"},
	{speaker: "Pink Panther", text: "Hmm. Well, a job is a job."},
	{speaker: "Pink Panther", text: "Time to see what Camp Chilly Wa Wa has in store for me this time."},
}

type Game struct {
	sceneMgr        *sceneManager
	player          *player
	dialog          *dialogSystem
	ui              *uiManager
	audio           *audioManager
	inv             *inventory
	lastScene       string
	mouseX          int32
	mouseY          int32
	monologuePlayed bool
	bullyDefeated   bool
	cookHelped      bool
}

func New(renderer *sdl.Renderer, font *engine.BitmapFont) *Game {
	g := &Game{
		sceneMgr: newSceneManager(renderer),
		player:   newPlayer(renderer),
		dialog:   newDialogSystem(font),
		ui:       newUIManager(font),
		audio:    newAudioManager(),
		inv:      newInventory(font),
	}
	g.lastScene = g.sceneMgr.currentName
	g.audio.playMusic(g.sceneMgr.current().musicPath)

	letterItem := createLetterTexture(renderer)
	comicItem := createComicBookTexture(renderer)
	beerItem := createBeerTexture(renderer)
	fishingRodItem := createFishingRodTexture(renderer)
	messNoteItem := createMessHallNoteTexture(renderer)
	marshmallowItem := createBurntMarshmallowTexture(renderer)

	if letterItem != nil {
		g.inv.addItem(letterItem)
	}

	g.setupLondonCallbacks(comicItem, beerItem)
	g.setupCampCallbacks(letterItem, comicItem, fishingRodItem, messNoteItem, marshmallowItem)
	g.ui.initCursors(renderer)

	return g
}

func (g *Game) setupLondonCallbacks(comicItem, beerItem *inventoryItem) {
	// Paper Man gives comic book on first dialog completion
	for _, n := range g.sceneMgr.scenes["street"].npcs {
		if n.name == "Paper Man" {
			pm := n
			inv := g.inv
			pm.onDialogEnd = func() {
				if comicItem != nil && !inv.hasItem("Comic Book") {
					inv.addItem(comicItem)
					pm.dialog = paperManPostComicDialog
				}
			}
			break
		}
	}

	// Crying Kid accepts the comic book when player is holding it
	for _, n := range g.sceneMgr.scenes["interior"].npcs {
		if n.name == "Crying Kid" {
			kid := n
			inv := g.inv
			kid.altDialogFunc = func() ([]dialogEntry, func()) {
				if inv.heldItem != nil && inv.heldItem.name == "Comic Book" {
					return cryingKidComicDialog, func() {
						inv.removeItem("Comic Book")
						inv.heldItem = nil
						kid.dialog = cryingKidHappyDialog
						kid.name = "Happy Kid"
					}
				}
				return nil, nil
			}
			break
		}
	}

	// Barmaid gives beer on first dialog completion
	for _, n := range g.sceneMgr.scenes["pub"].npcs {
		if n.name == "Barmaid" {
			bm := n
			inv := g.inv
			bm.onDialogEnd = func() {
				if beerItem != nil && !inv.hasItem("Pint of Beer") {
					inv.addItem(beerItem)
					bm.dialog = barmaidPostBeerDialog
				}
			}
			break
		}
	}

	// Bobby accepts beer and reveals clue
	for _, n := range g.sceneMgr.scenes["pub"].npcs {
		if n.name == "Bobby" {
			bobby := n
			inv := g.inv
			bobby.altDialogFunc = func() ([]dialogEntry, func()) {
				if inv.heldItem != nil && inv.heldItem.name == "Pint of Beer" {
					return bobbyBeerDialog, func() {
						inv.removeItem("Pint of Beer")
						inv.heldItem = nil
						bobby.animState = npcAnimDrink
						bobby.animOnce = true
						bobby.frameIdx = 0
						bobby.frameTimer = 0
						af := bobby.activeFrames()
						if len(af) > 0 {
							bobby.srcRect = af[0]
						}
						bobby.dialog = bobbyPostBeerDialog
					}
				}
				return nil, nil
			}
			break
		}
	}
}

func (g *Game) setupCampCallbacks(letterItem, comicItem, fishingRodItem, messNoteItem, marshmallowItem *inventoryItem) {
	inv := g.inv

	// Director Higgins: show letter to gain entry
	for _, n := range g.sceneMgr.scenes["camp_entrance"].npcs {
		if n.name == "Director Higgins" {
			higgins := n
			higgins.altDialogFunc = func() ([]dialogEntry, func()) {
				if inv.heldItem != nil && inv.heldItem.name == "Appointment Letter" {
					return higginsLetterDialog, func() {
						inv.removeItem("Appointment Letter")
						inv.heldItem = nil
						higgins.dialog = higginsPostLetterDialog
					}
				}
				return nil, nil
			}
			break
		}
	}

	// Tommy: give comic book -> reveals shortcut
	for _, n := range g.sceneMgr.scenes["camp_grounds"].npcs {
		if n.name == "Tommy" {
			kid := n
			kid.altDialogFunc = func() ([]dialogEntry, func()) {
				if inv.heldItem != nil && inv.heldItem.name == "Comic Book" {
					return tommyComicDialog, func() {
						inv.removeItem("Comic Book")
						inv.heldItem = nil
						kid.dialog = tommyHappyDialog
						kid.name = "Tommy"
					}
				}
				return nil, nil
			}
			break
		}
	}

	// Lily: opens up after first dialog
	for _, n := range g.sceneMgr.scenes["camp_grounds"].npcs {
		if n.name == "Lily" {
			girl := n
			girl.onDialogEnd = func() {
				if !girl.dialogDone {
					girl.dialog = lilySecondDialog
				}
			}
			break
		}
	}

	// Cook Marge: gives food item (marshmallow or stew) on first completion
	for _, n := range g.sceneMgr.scenes["camp_messhall"].npcs {
		if n.name == "Cook Marge" {
			marge := n
			game := g
			marge.onDialogEnd = func() {
				if !game.cookHelped && marshmallowItem != nil {
					game.cookHelped = true
					inv.addItem(marshmallowItem)
					marge.dialog = cookMargePostHelpDialog
				}
			}
			break
		}
	}

	// Jake: give food to pass
	for _, n := range g.sceneMgr.scenes["camp_grounds"].npcs {
		if n.name == "Jake" {
			bully := n
			game := g
			bully.altDialogFunc = func() ([]dialogEntry, func()) {
				if inv.heldItem != nil && inv.heldItem.name == "Burnt Marshmallow" {
					return jakeFedDialog, func() {
						inv.removeItem("Burnt Marshmallow")
						inv.heldItem = nil
						game.bullyDefeated = true
						bully.dialog = jakePostFedDialog
					}
				}
				return nil, nil
			}
			break
		}
	}

	// Floor item: Comic Book on the ground at camp_grounds
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok && comicItem != nil {
		fi := &floorItem{
			tex:     comicItem.tex,
			srcW:    comicItem.srcW,
			srcH:    comicItem.srcH,
			bounds:  sdl.Rect{X: 350, Y: 460, W: 36, H: 44},
			name:    "Comic Book",
			visible: true,
		}
		fi.onPickup = func() {
			if !fi.visible {
				return
			}
			fi.visible = false
			inv.addItem(comicItem)
			g.player.playAction(stateGrabbing, nil)
			g.dialog.startDialog([]dialogEntry{
				{speaker: "Pink Panther", text: "A comic book! Someone must have dropped it."},
			})
		}
		grounds.floorItems = append(grounds.floorItems, fi)
	}

	// Floor item: Fishing Rod at the lake dock
	if lake, ok := g.sceneMgr.scenes["camp_lake"]; ok && fishingRodItem != nil {
		fi := &floorItem{
			tex:     fishingRodItem.tex,
			srcW:    fishingRodItem.srcW,
			srcH:    fishingRodItem.srcH,
			bounds:  sdl.Rect{X: 650, Y: 370, W: 44, H: 50},
			name:    "Fishing Rod",
			visible: true,
		}
		fi.onPickup = func() {
			if !fi.visible {
				return
			}
			fi.visible = false
			inv.addItem(fishingRodItem)
			g.player.playAction(stateGrabbing, nil)
			g.dialog.startDialog([]dialogEntry{
				{speaker: "Pink Panther", text: "An old fishing rod. Might come in handy."},
			})
		}
		lake.floorItems = append(lake.floorItems, fi)

		footprintsExamined := false
		lake.hotspots = append(lake.hotspots, hotspot{
			bounds: sdl.Rect{X: 200, Y: 500, W: 120, H: 60},
			name:   "Footprints",
			onInteract: func() bool {
				if footprintsExamined {
					g.dialog.startDialog([]dialogEntry{
						{speaker: "Pink Panther", text: "Strange footprints... someone was here recently."},
					})
					return true
				}
				footprintsExamined = true
				g.dialog.startDialog([]dialogEntry{
					{speaker: "Pink Panther", text: "Hmm... footprints in the sand. Fresh ones."},
					{speaker: "Pink Panther", text: "They lead toward the water and then... vanish."},
					{speaker: "Pink Panther", text: "Someone was here recently. Very recently."},
				})
				return true
			},
		})
	}

	// Floor item: Mess Hall Note under a table
	if messHall, ok := g.sceneMgr.scenes["camp_messhall"]; ok && messNoteItem != nil {
		fi := &floorItem{
			tex:     messNoteItem.tex,
			srcW:    messNoteItem.srcW,
			srcH:    messNoteItem.srcH,
			bounds:  sdl.Rect{X: 400, Y: 475, W: 40, H: 32},
			name:    "Crumpled Note",
			visible: true,
		}
		fi.onPickup = func() {
			if !fi.visible {
				return
			}
			fi.visible = false
			inv.addItem(messNoteItem)
			g.player.playAction(stateGrabbing, nil)
			g.dialog.startDialog([]dialogEntry{
				{speaker: "Pink Panther", text: "What's this? A note hidden under the table..."},
				{speaker: "Pink Panther", text: "'Meet me at the lake. Midnight. Come alone.'"},
				{speaker: "Pink Panther", text: "Now THAT is suspicious."},
			})
		}
		messHall.floorItems = append(messHall.floorItems, fi)
	}
}

func (g *Game) Close() {
	g.audio.close()
}

func (g *Game) HandleClick(x, y int32) {
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
	if scancode == sdl.SCANCODE_SPACE && g.dialog.active {
		g.dialog.advance()
	}
}

func (g *Game) Update(dt float64, mx, my int32) {
	g.mouseX = mx
	g.mouseY = my
	scene := g.sceneMgr.current()

	if !g.monologuePlayed && g.sceneMgr.currentName == "camp_entrance" && !g.sceneMgr.transitioning {
		g.monologuePlayed = true
		g.dialog.startDialog(openingMonologue)
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
		// top
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: w, H: l.inset + 30})
		// bottom
		renderer.FillRect(&sdl.Rect{X: 0, Y: h - l.inset - 30, W: w, H: l.inset + 30})
		// left
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: l.inset + 40, H: h})
		// right
		renderer.FillRect(&sdl.Rect{X: w - l.inset - 40, Y: 0, W: l.inset + 40, H: h})
	}
}

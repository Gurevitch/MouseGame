package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	sceneMgr  *sceneManager
	player    *player
	dialog    *dialogSystem
	ui        *uiManager
	audio     *audioManager
	inv       *inventory
	lastScene string
	mouseX    int32
	mouseY    int32
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

	comicItem := createComicBookTexture(renderer)
	beerItem := createBeerTexture(renderer)
	g.setupNPCCallbacks(comicItem, beerItem)

	return g
}

func (g *Game) setupNPCCallbacks(comicItem, beerItem *inventoryItem) {
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
	if hs := scene.checkHotspotClick(x, y); hs != nil {
		tgt := hs.targetScene
		plr := g.player
		sm := g.sceneMgr
		onArrival := func() { sm.transitionTo(tgt, plr) }
		if hs.arrow == arrowLeft || hs.arrow == arrowRight {
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
	if !g.dialog.active && g.player.state == stateTalking {
		g.player.state = stateIdle
	}
	if !g.dialog.active && !g.sceneMgr.transitioning {
		g.player.update(dt, scene.blockers)
	}
	g.dialog.update(dt)
	g.sceneMgr.update(dt)
	scene.updateAmbient(dt)
	g.ui.updateHover(scene, mx, my, g.inv)
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

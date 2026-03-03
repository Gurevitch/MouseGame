package game

import (
	"fmt"
	"os"
	"time"

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
	g.setupNPCCallbacks(comicItem)

	// #region agent log -- log NPC positions at startup
	for sceneName, s := range g.sceneMgr.scenes {
		for _, n := range s.npcs {
			debugLog("game.go:New", "npc-pos", fmt.Sprintf(`{"scene":"%s","name":"%s","x":%d,"y":%d,"w":%d,"h":%d,"feetY":%d}`, sceneName, n.name, n.bounds.X, n.bounds.Y, n.bounds.W, n.bounds.H, n.bounds.Y+n.bounds.H))
		}
	}
	debugLog("game.go:New", "player-bounds", fmt.Sprintf(`{"minY":%.0f,"maxY":%.0f,"dstH":%d,"feetMinY":%.0f,"feetMaxY":%.0f}`, playerMinY, playerMaxY, playerDstH, playerMinY+playerDstH, playerMaxY+playerDstH))
	// #endregion
	return g
}

func (g *Game) setupNPCCallbacks(comicItem *inventoryItem) {
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
}

func (g *Game) Close() {
	g.audio.close()
}

func debugLog(loc, msg string, data string) {
	// #region agent log
	f, err := os.OpenFile("debug-e6d985.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	ts := time.Now().UnixMilli()
	fmt.Fprintf(f, `{"sessionId":"e6d985","location":"%s","message":"%s","data":%s,"timestamp":%d}`+"\n", loc, msg, data, ts)
	// #endregion
}

func (g *Game) HandleClick(x, y int32) {
	// #region agent log
	debugLog("game.go:HandleClick", "click", fmt.Sprintf(`{"x":%d,"y":%d,"scene":"%s"}`, x, y, g.sceneMgr.currentName))
	// #endregion

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
						g.player.facingLeft = g.player.x > float64(target.bounds.X)
						ds.startDialogWithCallback(entries, cb)
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
		plr.walkToAndDo(
			float64(hs.bounds.X+hs.bounds.W/2),
			float64(hs.bounds.Y+hs.bounds.H/2),
			func() { sm.transitionTo(tgt, plr) },
		)
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
	scene.drawNPCs(renderer)
	g.player.draw(renderer)

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

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
	}
	g.lastScene = g.sceneMgr.currentName
	g.audio.playMusic(g.sceneMgr.current().musicPath)
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
	if g.dialog.active {
		g.dialog.advance()
		return
	}
	if g.sceneMgr.transitioning {
		return
	}
	scene := g.sceneMgr.current()
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
	g.ui.updateHover(scene, mx, my)

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

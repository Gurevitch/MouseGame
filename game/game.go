package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	sceneMgr *sceneManager
	player   *player
	dialog   *dialogSystem
	ui       *uiManager
}

func New(renderer *sdl.Renderer, font *engine.BitmapFont) *Game {
	return &Game{
		sceneMgr: newSceneManager(renderer),
		player:   newPlayer(renderer),
		dialog:   newDialogSystem(font),
		ui:       newUIManager(font),
	}
}

func (g *Game) HandleClick(x, y int32) {
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
	scene := g.sceneMgr.current()
	if !g.dialog.active && !g.sceneMgr.transitioning {
		g.player.update(dt)
	}
	g.dialog.update(dt)
	g.sceneMgr.update(dt)
	scene.updateAmbient(dt)
	g.ui.updateHover(scene, mx, my)
}

func (g *Game) Draw(renderer *sdl.Renderer) {
	scene := g.sceneMgr.current()
	scene.drawBackground(renderer)
	scene.drawAmbient(renderer)
	scene.drawHotspots(renderer)
	scene.drawNPCs(renderer)
	g.player.draw(renderer)
	g.dialog.draw(renderer)
	g.ui.draw(renderer)
	g.sceneMgr.drawTransition(renderer)
}

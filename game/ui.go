package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type uiManager struct {
	font      *engine.BitmapFont
	hoverName string
}

func newUIManager(font *engine.BitmapFont) *uiManager {
	return &uiManager{font: font}
}

func (ui *uiManager) updateHover(s *scene, mx, my int32) {
	ui.hoverName = ""
	for _, n := range s.npcs {
		if n.containsPoint(mx, my) {
			ui.hoverName = n.name
			return
		}
	}
	pt := sdl.Point{X: mx, Y: my}
	for _, hs := range s.hotspots {
		if pt.InRect(&hs.bounds) {
			ui.hoverName = hs.name
			return
		}
	}
}

func (ui *uiManager) draw(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 120)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: 30})
	ui.font.DrawText(renderer,
		"CLICK TO WALK | CLICK CHARACTERS TO TALK | CLICK DOORS TO CHANGE ROOMS",
		8, 6, 2, sdl.Color{R: 255, G: 255, B: 255, A: 200})
	if ui.hoverName != "" {
		w := ui.font.TextWidth(ui.hoverName, 3)
		x := (engine.ScreenWidth - w) / 2
		renderer.SetDrawColor(0, 0, 0, 150)
		renderer.FillRect(&sdl.Rect{X: x - 6, Y: 34, W: w + 12, H: 28})
		ui.font.DrawText(renderer, ui.hoverName, x, 38, 3,
			sdl.Color{R: 255, G: 220, B: 100, A: 255})
	}
}

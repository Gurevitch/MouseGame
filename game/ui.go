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
		n.hovered = false
	}
	for _, n := range s.npcs {
		if n.containsPoint(mx, my) {
			ui.hoverName = n.name
			n.hovered = true
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
	renderer.SetDrawColor(0, 0, 0, 140)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: 36})

	txt := "CLICK TO WALK | CLICK CHARACTERS TO TALK | CLICK ARROWS TO CHANGE ROOMS"
	ui.font.DrawText(renderer, txt, 10, 11, 2,
		sdl.Color{R: 0, G: 0, B: 0, A: 120})
	ui.font.DrawText(renderer, txt, 9, 10, 2,
		sdl.Color{R: 255, G: 255, B: 255, A: 210})

	if ui.hoverName != "" {
		w := ui.font.TextWidth(ui.hoverName, 3)
		x := (engine.ScreenWidth - w) / 2
		renderer.SetDrawColor(0, 0, 0, 160)
		renderer.FillRect(&sdl.Rect{X: x - 10, Y: 38, W: w + 20, H: 32})
		renderer.SetDrawColor(255, 220, 100, 100)
		renderer.DrawRect(&sdl.Rect{X: x - 10, Y: 38, W: w + 20, H: 32})

		ui.font.DrawText(renderer, ui.hoverName, x+1, 44, 3,
			sdl.Color{R: 0, G: 0, B: 0, A: 120})
		ui.font.DrawText(renderer, ui.hoverName, x, 43, 3,
			sdl.Color{R: 255, G: 220, B: 100, A: 255})
	}
}

package game

import (
	"fmt"

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

func (ui *uiManager) updateHover(s *scene, mx, my int32, inv *inventory) {
	ui.hoverName = ""
	for _, n := range s.npcs {
		n.hovered = false
		n.itemMatch = false
	}
	for _, n := range s.npcs {
		if n.containsPoint(mx, my) {
			ui.hoverName = n.name
			n.hovered = true
			if inv.heldItem != nil && n.altDialogFunc != nil {
				entries, _ := n.altDialogFunc()
				if entries != nil {
					n.itemMatch = true
				}
			}
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

func (ui *uiManager) draw(renderer *sdl.Renderer, mx, my int32) {
	renderer.SetDrawColor(0, 0, 0, 140)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: 36})

	txt := "CLICK TO WALK | CLICK CHARACTERS TO TALK | CLICK ARROWS TO CHANGE ROOMS | mouse: (" + fmt.Sprintf("%d", mx) + ", " + fmt.Sprintf("%d", my) + ")"
	ui.font.DrawText(renderer, txt, 10, 11, 2,
		sdl.Color{R: 0, G: 0, B: 0, A: 120})
	ui.font.DrawText(renderer, txt, 9, 10, 2,
		sdl.Color{R: 255, G: 255, B: 255, A: 210})

	drawTaskIcon(renderer, engine.ScreenWidth-80, 6)
	ui.font.DrawText(renderer, "TASKS", engine.ScreenWidth-56, 13, 2,
		sdl.Color{R: 255, G: 180, B: 200, A: 200})

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

func drawTaskIcon(renderer *sdl.Renderer, x, y int32) {
	// Clipboard body
	renderer.SetDrawColor(180, 160, 140, 200)
	renderer.FillRect(&sdl.Rect{X: x, Y: y + 4, W: 18, H: 22})
	// Clipboard clip at top
	renderer.SetDrawColor(200, 180, 160, 220)
	renderer.FillRect(&sdl.Rect{X: x + 5, Y: y, W: 8, H: 6})
	// Task lines
	renderer.SetDrawColor(60, 50, 40, 220)
	renderer.FillRect(&sdl.Rect{X: x + 3, Y: y + 9, W: 12, H: 2})
	renderer.FillRect(&sdl.Rect{X: x + 3, Y: y + 14, W: 12, H: 2})
	renderer.FillRect(&sdl.Rect{X: x + 3, Y: y + 19, W: 8, H: 2})
	// Outline
	renderer.SetDrawColor(100, 80, 60, 180)
	renderer.DrawRect(&sdl.Rect{X: x, Y: y + 4, W: 18, H: 22})
}

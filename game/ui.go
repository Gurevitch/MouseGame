package game

import (
	"fmt"
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type cursorState int

const (
	cursorNormal cursorState = iota
	cursorTalk
	cursorGrab
	cursorArrowLeft
	cursorArrowRight
	cursorArrowUp
	cursorArrowDown
	cursorCount
)

type uiManager struct {
	font        *engine.BitmapFont
	hoverName   string
	cursor      cursorState
	cursorTex   [cursorCount]*sdl.Texture
	cursorW     [cursorCount]int32
	cursorH     [cursorCount]int32
	cursorTimer float64
}

func newUIManager(font *engine.BitmapFont) *uiManager {
	return &uiManager{font: font}
}

func (ui *uiManager) initCursors(renderer *sdl.Renderer) {
	mkSurf := func(w, h int32) *sdl.Surface {
		s, _ := sdl.CreateRGBSurface(0, w, h, 32,
			0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
		if s != nil {
			s.FillRect(nil, 0)
		}
		return s
	}
	px := func(s *sdl.Surface, x, y int32, r, g, b, a uint8) {
		if x < 0 || y < 0 || x >= s.W || y >= s.H {
			return
		}
		s.FillRect(&sdl.Rect{X: x, Y: y, W: 1, H: 1},
			sdl.MapRGBA(s.Format, r, g, b, a))
	}
	block := func(s *sdl.Surface, x, y, w, h int32, r, g, b, a uint8) {
		s.FillRect(&sdl.Rect{X: x, Y: y, W: w, H: h},
			sdl.MapRGBA(s.Format, r, g, b, a))
	}

	// Normal pointer -- classic angled arrow in pink with dark outline
	{
		const w, h int32 = 16, 20
		s := mkSurf(w, h)
		if s != nil {
			outline := [20][]int32{
				{0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5},
				{0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10},
				{0, 11}, {0, 7, 8}, {0, 1, 8, 9}, {1, 2, 9, 10},
				{2, 3, 10, 11}, {3, 4, 11, 12}, {4, 5, 12, 13},
				{5, 6, 13, 14}, {6, 14},
			}
			fill := [20][]int32{
				{}, {}, {1}, {1, 2}, {1, 2, 3}, {1, 2, 3, 4},
				{1, 2, 3, 4, 5}, {1, 2, 3, 4, 5, 6}, {1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7, 8}, {1, 2, 3, 4, 5, 6, 7, 8, 9},
				{1, 2, 3, 4, 5, 6}, {1, 2, 3, 4, 5, 6}, {2, 3, 4, 5, 6, 7},
				{3, 4, 5, 6, 7, 8}, {4, 5, 6, 7, 8, 9}, {5, 6, 7, 8, 9, 10},
				{6, 7, 8, 9, 10, 11}, {7, 8, 9, 10, 11, 12}, {7, 8, 9, 10, 11, 12, 13},
			}
			for row := int32(0); row < h; row++ {
				for _, col := range outline[row] {
					px(s, col, row, 30, 15, 25, 255)
				}
				for _, col := range fill[row] {
					if row < 3 {
						px(s, col, row, 255, 210, 225, 255)
					} else {
						px(s, col, row, 255, 150, 180, 255)
					}
				}
			}
			px(s, 1, 1, 255, 230, 240, 255)

			tex, _ := renderer.CreateTextureFromSurface(s)
			if tex != nil {
				tex.SetBlendMode(sdl.BLENDMODE_BLEND)
			}
			ui.cursorTex[cursorNormal] = tex
			ui.cursorW[cursorNormal] = w
			ui.cursorH[cursorNormal] = h
			s.Free()
		}
	}

	// Talk cursor (speech bubble)
	{
		const w, h = 22, 20
		s := mkSurf(w, h)
		if s != nil {
			block(s, 2, 0, 18, 2, 40, 20, 30, 255)
			block(s, 0, 2, 2, 10, 40, 20, 30, 255)
			block(s, 20, 2, 2, 10, 40, 20, 30, 255)
			block(s, 2, 12, 18, 2, 40, 20, 30, 255)
			block(s, 2, 2, 18, 10, 255, 250, 220, 240)
			block(s, 6, 14, 2, 2, 40, 20, 30, 255)
			block(s, 4, 16, 2, 2, 40, 20, 30, 255)
			block(s, 2, 18, 2, 2, 40, 20, 30, 255)
			block(s, 5, 5, 2, 4, 40, 20, 30, 200)
			block(s, 10, 5, 2, 4, 40, 20, 30, 200)
			block(s, 15, 5, 2, 4, 40, 20, 30, 200)

			tex, _ := renderer.CreateTextureFromSurface(s)
			if tex != nil {
				tex.SetBlendMode(sdl.BLENDMODE_BLEND)
			}
			ui.cursorTex[cursorTalk] = tex
			ui.cursorW[cursorTalk] = w
			ui.cursorH[cursorTalk] = h
			s.Free()
		}
	}

	// Grab hand
	{
		const w, h = 20, 22
		s := mkSurf(w, h)
		if s != nil {
			block(s, 4, 0, 4, 6, 40, 20, 30, 255)
			block(s, 5, 1, 2, 4, 255, 200, 170, 255)
			block(s, 9, 1, 4, 6, 40, 20, 30, 255)
			block(s, 10, 2, 2, 4, 255, 200, 170, 255)
			block(s, 14, 2, 4, 6, 40, 20, 30, 255)
			block(s, 15, 3, 2, 4, 255, 200, 170, 255)
			block(s, 2, 7, 16, 8, 40, 20, 30, 255)
			block(s, 3, 8, 14, 6, 255, 200, 170, 255)
			block(s, 2, 15, 16, 7, 40, 20, 30, 255)
			block(s, 3, 16, 14, 5, 255, 200, 170, 255)

			tex, _ := renderer.CreateTextureFromSurface(s)
			if tex != nil {
				tex.SetBlendMode(sdl.BLENDMODE_BLEND)
			}
			ui.cursorTex[cursorGrab] = tex
			ui.cursorW[cursorGrab] = w
			ui.cursorH[cursorGrab] = h
			s.Free()
		}
	}

	// Arrow cursors -- clean filled triangles with dark outline
	arrowSurf := func(dir int) (*sdl.Texture, int32, int32) {
		var sw, sh int32
		if dir == 2 {
			sw, sh = 16, 12
		} else {
			sw, sh = 12, 16
		}
		s := mkSurf(sw, sh)
		if s == nil {
			return nil, 0, 0
		}

		switch dir {
		case 0: // left -- tip at left edge, base at right
			mid := sh / 2
			for row := int32(0); row < sh; row++ {
				dist := row - mid
				if dist < 0 {
					dist = -dist
				}
				w := sw - int32(float64(dist)*float64(sw)/float64(mid))
				if w < 1 {
					w = 1
				}
				block(s, 0, row, w, 1, 255, 220, 100, 255)
				px(s, 0, row, 30, 15, 25, 255)
				if w > 1 {
					px(s, w-1, row, 30, 15, 25, 255)
				}
			}
			for col := int32(0); col < sw; col++ {
				px(s, col, 0, 30, 15, 25, 255)
				px(s, col, sh-1, 30, 15, 25, 255)
			}
		case 1: // right -- tip at right edge, base at left
			mid := sh / 2
			for row := int32(0); row < sh; row++ {
				dist := row - mid
				if dist < 0 {
					dist = -dist
				}
				w := sw - int32(float64(dist)*float64(sw)/float64(mid))
				if w < 1 {
					w = 1
				}
				block(s, sw-w, row, w, 1, 255, 220, 100, 255)
				px(s, sw-1, row, 30, 15, 25, 255)
				if w > 1 {
					px(s, sw-w, row, 30, 15, 25, 255)
				}
			}
			for col := int32(0); col < sw; col++ {
				px(s, col, 0, 30, 15, 25, 255)
				px(s, col, sh-1, 30, 15, 25, 255)
			}
		case 2: // up -- tip at top, base at bottom
			mid := sw / 2
			for col := int32(0); col < sw; col++ {
				dist := col - mid
				if dist < 0 {
					dist = -dist
				}
				h := sh - int32(float64(dist)*float64(sh)/float64(mid))
				if h < 1 {
					h = 1
				}
				block(s, col, 0, 1, h, 255, 220, 100, 255)
				px(s, col, 0, 30, 15, 25, 255)
				if h > 1 {
					px(s, col, h-1, 30, 15, 25, 255)
				}
			}
			for row := int32(0); row < sh; row++ {
				px(s, 0, row, 30, 15, 25, 255)
				px(s, sw-1, row, 30, 15, 25, 255)
			}
		case 3: // down -- tip at bottom, base at top
			mid := sw / 2
			for col := int32(0); col < sw; col++ {
				dist := col - mid
				if dist < 0 {
					dist = -dist
				}
				h := sh - int32(float64(dist)*float64(sh)/float64(mid))
				if h < 1 {
					h = 1
				}
				block(s, col, sh-h, 1, h, 255, 220, 100, 255)
				px(s, col, sh-1, 30, 15, 25, 255)
				if h > 1 {
					px(s, col, sh-h, 30, 15, 25, 255)
				}
			}
			for row := int32(0); row < sh; row++ {
				px(s, 0, row, 30, 15, 25, 255)
				px(s, sw-1, row, 30, 15, 25, 255)
			}
		}

		tex, _ := renderer.CreateTextureFromSurface(s)
		if tex != nil {
			tex.SetBlendMode(sdl.BLENDMODE_BLEND)
		}
		s.Free()
		return tex, sw, sh
	}

	ui.cursorTex[cursorArrowLeft], ui.cursorW[cursorArrowLeft], ui.cursorH[cursorArrowLeft] = arrowSurf(0)
	ui.cursorTex[cursorArrowRight], ui.cursorW[cursorArrowRight], ui.cursorH[cursorArrowRight] = arrowSurf(1)
	ui.cursorTex[cursorArrowUp], ui.cursorW[cursorArrowUp], ui.cursorH[cursorArrowUp] = arrowSurf(2)
	ui.cursorTex[cursorArrowDown], ui.cursorW[cursorArrowDown], ui.cursorH[cursorArrowDown] = arrowSurf(3)
}

func (ui *uiManager) updateHover(s *scene, mx, my int32, inv *inventory, dt float64) {
	ui.hoverName = ""
	ui.cursor = cursorNormal
	ui.cursorTimer += dt
	for _, n := range s.npcs {
		n.hovered = false
		n.itemMatch = false
	}
	for _, n := range s.npcs {
		if n.containsPoint(mx, my) {
			ui.hoverName = n.name
			n.hovered = true
			ui.cursor = cursorTalk
			if inv.heldItem != nil && n.altDialogFunc != nil {
				entries, _ := n.altDialogFunc()
				if entries != nil {
					n.itemMatch = true
				}
			}
			return
		}
	}
	for _, fi := range s.floorItems {
		if fi.visible {
			pt := sdl.Point{X: mx, Y: my}
			if pt.InRect(&fi.bounds) {
				ui.hoverName = fi.name
				ui.cursor = cursorGrab
				return
			}
		}
	}
	pt := sdl.Point{X: mx, Y: my}
	for _, hs := range s.hotspots {
		if pt.InRect(&hs.bounds) {
			ui.hoverName = hs.name
		switch hs.arrow {
		case arrowLeft:
			ui.cursor = cursorArrowLeft
		case arrowRight:
			ui.cursor = cursorArrowRight
		case arrowUp:
			ui.cursor = cursorArrowUp
		case arrowDown:
			ui.cursor = cursorArrowDown
		}
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

func (ui *uiManager) drawCursor(renderer *sdl.Renderer, mx, my int32) {
	c := ui.cursor
	tex := ui.cursorTex[c]
	if tex == nil {
		tex = ui.cursorTex[cursorNormal]
		c = cursorNormal
	}
	if tex == nil {
		return
	}
	w := ui.cursorW[c] * 2
	h := ui.cursorH[c] * 2

	bob := int32(math.Sin(ui.cursorTimer*3.0) * 2.0)

	var dx, dy int32
	switch c {
	case cursorArrowLeft:
		dx = -w
		dy = -h / 2
	case cursorArrowRight:
		dy = -h / 2
	case cursorArrowUp:
		dx = -w / 2
		dy = -h
	case cursorArrowDown:
		dx = -w / 2
		dy = 0
	case cursorTalk:
		dx = -w / 2
		dy = -h - 4
	case cursorGrab:
		dx = -w / 2
		dy = -h / 2
	default:
		dx = 0
		dy = 0
	}

	dst := sdl.Rect{X: mx + dx, Y: my + dy + bob, W: w, H: h}
	renderer.Copy(tex, nil, &dst)
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

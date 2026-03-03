package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type inventoryItem struct {
	name string
	tex  *sdl.Texture
	srcW int32
	srcH int32
	desc string
}

type inventory struct {
	font        *engine.BitmapFont
	items       []*inventoryItem
	open        bool
	selectedIdx int
	pulse       float64
	heldItem    *inventoryItem
}

func newInventory(font *engine.BitmapFont) *inventory {
	return &inventory{font: font}
}

func (inv *inventory) addItem(item *inventoryItem) {
	inv.items = append(inv.items, item)
}

func (inv *inventory) removeItem(name string) {
	for i, it := range inv.items {
		if it.name == name {
			inv.items = append(inv.items[:i], inv.items[i+1:]...)
			if inv.selectedIdx >= len(inv.items) && inv.selectedIdx > 0 {
				inv.selectedIdx--
			}
			return
		}
	}
}

func (inv *inventory) hasItem(name string) bool {
	for _, it := range inv.items {
		if it.name == name {
			return true
		}
	}
	return false
}

func (inv *inventory) toggle() {
	if len(inv.items) == 0 {
		return
	}
	inv.open = !inv.open
}

func (inv *inventory) update(dt float64) {
	inv.pulse += dt
}

const (
	invOvalW      = 360
	invOvalH      = 300
	invArrowSize  = 30
	invArrowPad   = 50
)

func invOvalCenter() (int32, int32) {
	return engine.ScreenWidth / 2, engine.ScreenHeight / 2
}

func (inv *inventory) handleClick(x, y int32) bool {
	if !inv.open {
		return false
	}
	cx, cy := invOvalCenter()

	leftArrowX := cx - invOvalW/2 - invArrowPad
	rightArrowX := cx + invOvalW/2 + invArrowPad
	arrowY := cy

	if len(inv.items) > 1 {
		if x >= leftArrowX-invArrowSize && x <= leftArrowX+invArrowSize &&
			y >= arrowY-invArrowSize && y <= arrowY+invArrowSize {
			inv.selectedIdx--
			if inv.selectedIdx < 0 {
				inv.selectedIdx = len(inv.items) - 1
			}
			return true
		}
		if x >= rightArrowX-invArrowSize && x <= rightArrowX+invArrowSize &&
			y >= arrowY-invArrowSize && y <= arrowY+invArrowSize {
			inv.selectedIdx = (inv.selectedIdx + 1) % len(inv.items)
			return true
		}
	}

	// Click outside oval closes it without selecting
	dx := float64(x-cx) / float64(invOvalW/2)
	dy := float64(y-cy) / float64(invOvalH/2)
	if dx*dx+dy*dy > 1.0 {
		inv.open = false
		return true
	}

	// Click inside oval selects the current item as held
	if len(inv.items) > 0 {
		inv.heldItem = inv.items[inv.selectedIdx]
		inv.open = false
	}
	return true
}

func (inv *inventory) draw(renderer *sdl.Renderer) {
	if !inv.open || len(inv.items) == 0 {
		return
	}

	cx, cy := invOvalCenter()

	// Dim background
	renderer.SetDrawColor(0, 0, 0, 140)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})

	// Outer oval (pink border)
	drawFilledOval(renderer, cx, cy, invOvalW/2+8, invOvalH/2+8, 255, 180, 200, 220)
	// Inner oval (dark background)
	drawFilledOval(renderer, cx, cy, invOvalW/2, invOvalH/2, 30, 20, 40, 230)
	// Inner highlight rim
	drawOvalOutline(renderer, cx, cy, invOvalW/2, invOvalH/2, 255, 200, 220, 120)

	// Draw current item
	item := inv.items[inv.selectedIdx]
	if item.tex != nil {
		maxW := int32(200)
		maxH := int32(180)
		scale := float64(maxW) / float64(item.srcW)
		if sh := float64(maxH) / float64(item.srcH); sh < scale {
			scale = sh
		}
		dstW := int32(float64(item.srcW) * scale)
		dstH := int32(float64(item.srcH) * scale)
		dst := sdl.Rect{X: cx - dstW/2, Y: cy - dstH/2 - 10, W: dstW, H: dstH}
		renderer.Copy(item.tex, nil, &dst)
	}

	// Item name below
	nameW := inv.font.TextWidth(item.name, 3)
	inv.font.DrawText(renderer, item.name, cx-nameW/2+1, cy+invOvalH/2-50+1, 3,
		sdl.Color{R: 0, G: 0, B: 0, A: 180})
	inv.font.DrawText(renderer, item.name, cx-nameW/2, cy+invOvalH/2-50, 3,
		sdl.Color{R: 255, G: 220, B: 100, A: 255})

	// Item count
	countTxt := ""
	if len(inv.items) > 1 {
		countTxt = string(rune('1'+inv.selectedIdx)) + "/" + string(rune('0'+len(inv.items)))
	}
	if countTxt != "" {
		cw := inv.font.TextWidth(countTxt, 2)
		inv.font.DrawText(renderer, countTxt, cx-cw/2, cy+invOvalH/2-25, 2,
			sdl.Color{R: 200, G: 200, B: 200, A: 200})
	}

	// Navigation arrows
	if len(inv.items) > 1 {
		pulse := uint8(180 + int(40*math.Sin(inv.pulse*3.0)))
		leftX := cx - invOvalW/2 - invArrowPad
		rightX := cx + invOvalW/2 + invArrowPad
		drawInvArrow(renderer, leftX, cy, invArrowSize, true, pulse)
		drawInvArrow(renderer, rightX, cy, invArrowSize, false, pulse)
	}
}

func (inv *inventory) drawHeld(renderer *sdl.Renderer, mx, my int32) {
	if inv.heldItem == nil || inv.heldItem.tex == nil {
		return
	}
	item := inv.heldItem
	const sz = 48
	scale := float64(sz) / float64(item.srcW)
	if sh := float64(sz) / float64(item.srcH); sh < scale {
		scale = sh
	}
	dstW := int32(float64(item.srcW) * scale)
	dstH := int32(float64(item.srcH) * scale)
	dst := sdl.Rect{X: mx + 12, Y: my + 12, W: dstW, H: dstH}

	renderer.SetDrawColor(255, 220, 100, 160)
	pad := int32(3)
	renderer.DrawRect(&sdl.Rect{X: dst.X - pad, Y: dst.Y - pad, W: dst.W + pad*2, H: dst.H + pad*2})
	renderer.Copy(item.tex, nil, &dst)
}

func drawInvArrow(renderer *sdl.Renderer, cx, cy, size int32, leftFacing bool, alpha uint8) {
	renderer.SetDrawColor(255, 220, 100, alpha)
	if leftFacing {
		tipX := cx - size/2
		baseX := cx + size/2
		for y := cy - size; y <= cy+size; y++ {
			t := float64(y-(cy-size)) / float64(2*size)
			var x0, x1 int32
			if t <= 0.5 {
				x1 = baseX
				x0 = baseX - int32(float64(baseX-tipX)*t*2)
			} else {
				x1 = baseX
				x0 = baseX - int32(float64(baseX-tipX)*(1.0-t)*2)
			}
			renderer.DrawLine(x0, y, x1, y)
		}
	} else {
		tipX := cx + size/2
		baseX := cx - size/2
		for y := cy - size; y <= cy+size; y++ {
			t := float64(y-(cy-size)) / float64(2*size)
			var x0, x1 int32
			if t <= 0.5 {
				x0 = baseX
				x1 = baseX + int32(float64(tipX-baseX)*t*2)
			} else {
				x0 = baseX
				x1 = baseX + int32(float64(tipX-baseX)*(1.0-t)*2)
			}
			renderer.DrawLine(x0, y, x1, y)
		}
	}
}

func createComicBookTexture(renderer *sdl.Renderer) *inventoryItem {
	w := int32(80)
	h := int32(100)
	surface, err := sdl.CreateRGBSurface(0, w, h, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		return nil
	}
	defer surface.Free()

	surface.FillRect(nil, sdl.MapRGBA(surface.Format, 0, 0, 0, 0))

	// Book cover (yellow-ish)
	surface.FillRect(&sdl.Rect{X: 4, Y: 4, W: w - 8, H: h - 8}, sdl.MapRGBA(surface.Format, 240, 210, 80, 255))
	// Spine
	surface.FillRect(&sdl.Rect{X: 4, Y: 4, W: 8, H: h - 8}, sdl.MapRGBA(surface.Format, 200, 160, 40, 255))
	// Title area
	surface.FillRect(&sdl.Rect{X: 16, Y: 10, W: w - 24, H: 20}, sdl.MapRGBA(surface.Format, 220, 50, 50, 255))
	// Illustration area
	surface.FillRect(&sdl.Rect{X: 16, Y: 36, W: w - 24, H: 40}, sdl.MapRGBA(surface.Format, 255, 255, 255, 255))
	// Pink panther silhouette (simple)
	surface.FillRect(&sdl.Rect{X: 30, Y: 42, W: 12, H: 28}, sdl.MapRGBA(surface.Format, 255, 150, 180, 255))
	surface.FillRect(&sdl.Rect{X: 32, Y: 38, W: 8, H: 8}, sdl.MapRGBA(surface.Format, 255, 150, 180, 255))

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	return &inventoryItem{
		name: "Comic Book",
		tex:  tex,
		srcW: w,
		srcH: h,
		desc: "A colorful comic book from the Paper Man.",
	}
}

func drawFilledOval(renderer *sdl.Renderer, cx, cy, rx, ry int32, r, g, b, a uint8) {
	renderer.SetDrawColor(r, g, b, a)
	for y := -ry; y <= ry; y++ {
		halfW := int32(float64(rx) * math.Sqrt(1.0-float64(y*y)/float64(ry*ry)))
		renderer.DrawLine(cx-halfW, cy+y, cx+halfW, cy+y)
	}
}

func drawOvalOutline(renderer *sdl.Renderer, cx, cy, rx, ry int32, r, g, b, a uint8) {
	renderer.SetDrawColor(r, g, b, a)
	steps := 80
	for i := 0; i < steps; i++ {
		angle := float64(i) * 2.0 * math.Pi / float64(steps)
		nextAngle := float64(i+1) * 2.0 * math.Pi / float64(steps)
		x1 := cx + int32(float64(rx)*math.Cos(angle))
		y1 := cy + int32(float64(ry)*math.Sin(angle))
		x2 := cx + int32(float64(rx)*math.Cos(nextAngle))
		y2 := cy + int32(float64(ry)*math.Sin(nextAngle))
		renderer.DrawLine(x1, y1, x2, y2)
	}
}

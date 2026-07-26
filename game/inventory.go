package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type inventoryItem struct {
	name  string
	tex   *sdl.Texture
	srcW  int32
	srcH  int32
	desc  string
	owner string // "player", "lily", "curator", "none" (on ground)
	// content box within the texture (post color-key artwork bounds).
	// Zero cw/ch = unknown -> draw the full texture. The bag and the held
	// cursor draw by this box so icons center on the ART, not the canvas
	// (generated icons often float off-center in big margins; 2026-06-11 #37).
	cbX, cbY, cbW, cbH int32
	// iconScale (2026-06-20 #17) multiplies the bag fit-scale so an oversized
	// icon (the charcoal pencil) can be drawn smaller. 0/1 = full fit.
	iconScale float64
}

// contentSrc returns the source rect + dimensions to draw an item by: the
// content box when known, else the full texture.
func (it *inventoryItem) contentSrc() (*sdl.Rect, int32, int32) {
	if it.cbW > 0 && it.cbH > 0 {
		return &sdl.Rect{X: it.cbX, Y: it.cbY, W: it.cbW, H: it.cbH}, it.cbW, it.cbH
	}
	return nil, it.srcW, it.srcH
}

type inventory struct {
	font        *engine.BitmapFont
	items       []*inventoryItem
	open        bool
	selectedIdx int
	pulse       float64
	heldItem    *inventoryItem
	circleTex   *sdl.Texture
	circleW     int32
	circleH     int32
	// onSelectItem lets Game intercept item-clicks for items that should
	// trigger an action instead of being held (e.g. Travel Map opens the
	// globe). Return true to consume the click and skip the held-item path.
	onSelectItem func(*inventoryItem) bool
}

func newInventory(font *engine.BitmapFont, renderer *sdl.Renderer) *inventory {
	inv := &inventory{font: font}
	tex, w, h := engine.SafeTextureFromPNGKeyed(renderer, "assets/images/ui/inv_circle.png")
	if tex != nil {
		tex.SetBlendMode(sdl.BLENDMODE_BLEND)
		inv.circleTex = tex
		inv.circleW = w
		inv.circleH = h
	}
	return inv
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

// giveItemTo transfers an item's ownership to a new owner and removes it from inventory
func (inv *inventory) giveItemTo(name, newOwner string) {
	for _, it := range inv.items {
		if it.name == name {
			it.owner = newOwner
			break
		}
	}
	// If the given item is riding the cursor, drop it too - trade paths that
	// don't go through the inline hand-off (Pierre's onClickOverride) used to
	// leave a ghost item stuck on the cursor (2026-06-11).
	if inv.heldItem != nil && inv.heldItem.name == name {
		inv.heldItem = nil
	}
	inv.removeItem(name)
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
	// User playtest #18: bag felt too small on a 1400×800 screen. Bumped from
	// 720×600 to 816×680 (kept the inv_circle.png 1.2:1 aspect so the art does
	// not stretch). Centered, it spans x≈276..1124 / y≈4..716 - fills more of
	// the screen without clipping the top edge.
	invOvalW = 816
	invOvalH = 680
)

func invOvalCenter() (int32, int32) {
	return engine.ScreenWidth / 2, engine.ScreenHeight / 2
}

func (inv *inventory) handleClick(x, y int32) bool {
	if !inv.open {
		return false
	}
	cx, cy := invOvalCenter()

	dx := float64(x-cx) / float64(invOvalW/2)
	dy := float64(y-cy) / float64(invOvalH/2)
	if dx*dx+dy*dy > 1.0 {
		inv.open = false
		return true
	}

	// Left/right click bands page through items (the inv_circle hand grips sit
	// on these bands). ~10% of the oval width on each side; the wide center band
	// picks the current item. (#18 oval grew to 816 wide.)
	leftCut := int32(invOvalW / 10)
	rightCut := int32(invOvalW / 10)
	inArrowBand := x < cx-leftCut || x > cx+rightCut
	if len(inv.items) > 1 {
		if x < cx-leftCut {
			inv.selectedIdx--
			if inv.selectedIdx < 0 {
				inv.selectedIdx = len(inv.items) - 1
			}
			return true
		}
		if x > cx+rightCut {
			inv.selectedIdx = (inv.selectedIdx + 1) % len(inv.items)
			return true
		}
	} else if inArrowBand {
		// 2026-06-15 #19: with a single item, a click on the L/R paging bands
		// used to fall through and SELECT the only item (e.g. opening the Travel
		// Map even though the player was clicking the arrows). Consume the click
		// without selecting - only the center band picks the item up.
		return true
	}

	if len(inv.items) > 0 {
		picked := inv.items[inv.selectedIdx]
		// Some items (Travel Map) should fire an action immediately instead
		// of being held for the next world click. Game registers the hook
		// via onSelectItem in Game.New; if the hook consumes the click, we
		// close the inventory without setting heldItem.
		if inv.onSelectItem != nil && inv.onSelectItem(picked) {
			inv.open = false
			return true
		}
		inv.heldItem = picked
		inv.open = false
	}
	return true
}

func (inv *inventory) draw(renderer *sdl.Renderer) {
	if !inv.open || len(inv.items) == 0 {
		return
	}

	cx, cy := invOvalCenter()

	// Dim background - user 2026-05-22: bumped from 140 to 190 so the modal
	// freeze reads more clearly visually (matches Update-side gate above).
	renderer.SetDrawColor(0, 0, 0, 190)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})

	if inv.circleTex != nil {
		ovalW := int32(invOvalW + 32)
		ovalH := int32(invOvalH + 32)
		dst := sdl.Rect{X: cx - ovalW/2, Y: cy - ovalH/2, W: ovalW, H: ovalH}
		renderer.Copy(inv.circleTex, nil, &dst)
	} else {
		drawFilledOval(renderer, cx, cy, invOvalW/2+8, invOvalH/2+8, 255, 180, 200, 220)
		drawFilledOval(renderer, cx, cy, invOvalW/2, invOvalH/2, 30, 20, 40, 230)
		drawOvalOutline(renderer, cx, cy, invOvalW/2, invOvalH/2, 255, 200, 220, 120)
	}

	item := inv.items[inv.selectedIdx]
	if item.tex != nil {
		// Scaled up with the larger oval (#18) so the item art keeps filling
		// the bag rather than looking lost in the middle. Drawn by the
		// content box so the artwork itself is centered (#37).
		src, sw, sh0 := item.contentSrc()
		maxW := int32(440)
		maxH := int32(400)
		scale := float64(maxW) / float64(sw)
		if sh := float64(maxH) / float64(sh0); sh < scale {
			scale = sh
		}
		if item.iconScale > 0 {
			scale *= item.iconScale
		}
		dstW := int32(float64(sw) * scale)
		dstH := int32(float64(sh0) * scale)
		dst := sdl.Rect{X: cx - dstW/2, Y: cy - dstH/2 - 20, W: dstW, H: dstH}
		renderer.Copy(item.tex, src, &dst)
	}

	nameW := inv.font.TextWidth(item.name, 4)
	inv.font.DrawText(renderer, item.name, cx-nameW/2+2, cy+invOvalH/2-90+2, 4,
		sdl.Color{R: 0, G: 0, B: 0, A: 200})
	inv.font.DrawText(renderer, item.name, cx-nameW/2, cy+invOvalH/2-90, 4,
		sdl.Color{R: 255, G: 220, B: 120, A: 255})

	if len(inv.items) > 1 {
		countTxt := string(rune('1'+inv.selectedIdx)) + "/" + string(rune('0'+len(inv.items)))
		cw := inv.font.TextWidth(countTxt, 3)
		inv.font.DrawText(renderer, countTxt, cx-cw/2, cy+invOvalH/2-40, 3,
			sdl.Color{R: 220, G: 220, B: 220, A: 220})

		// User playtest #17: the drawn chevron "logos" are removed. The
		// inv_circle.png art already shows hand grips on the left/right, and
		// the left/right click bands (handleClick) still page through items -
		// the extra drawn arrows were redundant clutter.
	}
}

func (inv *inventory) drawHeld(renderer *sdl.Renderer, mx, my int32) {
	if inv.heldItem == nil || inv.heldItem.tex == nil {
		return
	}
	item := inv.heldItem
	src, sw, sh0 := item.contentSrc()
	const sz = 48
	scale := float64(sz) / float64(sw)
	if sh := float64(sz) / float64(sh0); sh < scale {
		scale = sh
	}
	dstW := int32(float64(sw) * scale)
	dstH := int32(float64(sh0) * scale)
	dst := sdl.Rect{X: mx + 12, Y: my + 12, W: dstW, H: dstH}

	cx := dst.X + dst.W/2
	cy := dst.Y + dst.H/2
	pulse := 0.7 + 0.3*math.Sin(inv.pulse*3.0)
	for i := int32(5); i >= 1; i-- {
		rx := dst.W/2 + i + 2
		ry := dst.H/2 + i + 2
		a := uint8(float64(20-i*3) * pulse)
		drawFilledOval(renderer, cx, cy, rx, ry, 255, 220, 100, a)
	}
	renderer.Copy(item.tex, src, &dst)
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

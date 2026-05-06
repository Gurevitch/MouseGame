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
	inv.removeItem(name)
}

// itemOwner returns who owns a specific item, or "" if not found
func (inv *inventory) itemOwner(name string) string {
	for _, it := range inv.items {
		if it.name == name {
			return it.owner
		}
	}
	return ""
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
	invOvalW = 720
	invOvalH = 600
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

	if len(inv.items) > 1 {
		leftCut := int32(float64(invOvalW) * 0.20)
		rightCut := int32(float64(invOvalW) * 0.20)
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

	// Dim background
	renderer.SetDrawColor(0, 0, 0, 140)
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
		maxW := int32(360)
		maxH := int32(320)
		scale := float64(maxW) / float64(item.srcW)
		if sh := float64(maxH) / float64(item.srcH); sh < scale {
			scale = sh
		}
		dstW := int32(float64(item.srcW) * scale)
		dstH := int32(float64(item.srcH) * scale)
		dst := sdl.Rect{X: cx - dstW/2, Y: cy - dstH/2 - 20, W: dstW, H: dstH}
		renderer.Copy(item.tex, nil, &dst)
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

		pulse := 0.5 + 0.5*math.Sin(inv.pulse*2.4)
		alpha := uint8(120 + pulse*90)
		chevSize := int32(28)
		leftX := cx - invOvalW/2 + 70
		rightX := cx + invOvalW/2 - 70
		drawChevron(renderer, leftX, cy, chevSize, true, alpha)
		drawChevron(renderer, rightX, cy, chevSize, false, alpha)
	}
}

// drawChevron renders a simple two-stroke ">" or "<" pointer inside the
// inventory oval to hint at click-to-cycle zones. No fill, just thick lines
// so it feels clean instead of arrow-button heavy.
func drawChevron(renderer *sdl.Renderer, cx, cy, size int32, leftFacing bool, alpha uint8) {
	renderer.SetDrawColor(255, 230, 160, alpha)
	thickness := int32(4)
	for t := -thickness; t <= thickness; t++ {
		if leftFacing {
			for i := int32(0); i < size; i++ {
				renderer.DrawPoint(cx+i+t, cy-size+i)
				renderer.DrawPoint(cx+i+t, cy+size-i)
			}
		} else {
			for i := int32(0); i < size; i++ {
				renderer.DrawPoint(cx-i+t, cy-size+i)
				renderer.DrawPoint(cx-i+t, cy+size-i)
			}
		}
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

	cx := dst.X + dst.W/2
	cy := dst.Y + dst.H/2
	pulse := 0.7 + 0.3*math.Sin(inv.pulse*3.0)
	for i := int32(5); i >= 1; i-- {
		rx := dst.W/2 + i + 2
		ry := dst.H/2 + i + 2
		a := uint8(float64(20-i*3) * pulse)
		drawFilledOval(renderer, cx, cy, rx, ry, 255, 220, 100, a)
	}
	renderer.Copy(item.tex, nil, &dst)
}

func createComicBookTexture(renderer *sdl.Renderer) *inventoryItem {
	if item := createItemFromPNG(renderer, "assets/images/items/comic_book.png",
		"Comic Book", "A colorful comic book from the Paper Man."); item != nil {
		return item
	}
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

func createBeerTexture(renderer *sdl.Renderer) *inventoryItem {
	w := int32(60)
	h := int32(90)
	surface, err := sdl.CreateRGBSurface(0, w, h, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		return nil
	}
	defer surface.Free()

	surface.FillRect(nil, sdl.MapRGBA(surface.Format, 0, 0, 0, 0))

	glass := sdl.MapRGBA(surface.Format, 200, 200, 210, 120)
	surface.FillRect(&sdl.Rect{X: 12, Y: 15, W: 36, H: 60}, glass)

	beer := sdl.MapRGBA(surface.Format, 210, 170, 50, 230)
	surface.FillRect(&sdl.Rect{X: 14, Y: 25, W: 32, H: 48}, beer)

	foam := sdl.MapRGBA(surface.Format, 255, 250, 230, 250)
	surface.FillRect(&sdl.Rect{X: 12, Y: 15, W: 36, H: 14}, foam)

	handle := sdl.MapRGBA(surface.Format, 180, 180, 190, 200)
	surface.FillRect(&sdl.Rect{X: 48, Y: 30, W: 8, H: 30}, handle)
	surface.FillRect(&sdl.Rect{X: 44, Y: 28, W: 12, H: 6}, handle)
	surface.FillRect(&sdl.Rect{X: 44, Y: 56, W: 12, H: 6}, handle)

	base := sdl.MapRGBA(surface.Format, 180, 180, 190, 220)
	surface.FillRect(&sdl.Rect{X: 8, Y: 75, W: 44, H: 8}, base)

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	return &inventoryItem{
		name: "Pint of Beer",
		tex:  tex,
		srcW: w,
		srcH: h,
		desc: "A frothy pint from the barmaid.",
	}
}

func createItemFromPNG(renderer *sdl.Renderer, path, name, desc string) *inventoryItem {
	tex, w, h := engine.SafeTextureFromPNGKeyed(renderer, path)
	if tex == nil {
		return nil
	}
	return &inventoryItem{name: name, tex: tex, srcW: w, srcH: h, desc: desc}
}

func createLetterTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/letter.png",
		"Appointment Letter", "Official letter appointing PP as substitute counselor.")
}

func createFishingRodTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/fishing_rod.png",
		"Fishing Rod", "An old fishing rod found at the dock.")
}

func createMessHallNoteTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/mess_note.png",
		"Mess Hall Note", "A crumpled note found in the mess hall.")
}

func createBurntMarshmallowTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/marshmallow.png",
		"Burnt Marshmallow", "A charred marshmallow. Might be useful as a distraction.")
}

func createMuseumTicketTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/museum_ticket.png",
		"Museum Ticket", "An admission ticket to the Musee d'Art in Paris.")
}

func createMagnifyingGlassTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/magnifying_glass.png",
		"Magnifying Glass", "A brass magnifying glass. Perfect for examining details.")
}

func createFakePaintingTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/fake_painting.png",
		"Fake Painting", "A convincing forgery of a famous masterpiece.")
}

func createCatacombKeyTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/catacomb_key.png",
		"Catacomb Key", "An ancient iron key with a skull emblem. Opens something underground.")
}

func createGoldenArtifactTexture(renderer *sdl.Renderer) *inventoryItem {
	return createItemFromPNG(renderer, "assets/images/items/golden_artifact.png",
		"Golden Artifact", "A priceless golden idol. The prize everyone has been searching for.")
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

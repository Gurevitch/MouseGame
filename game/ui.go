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
	cursorArrowDownRight
	// cursorPoint is the pink "action" pointing-hand. Its PNG is a 2-frame
	// strip (idle finger | pressed finger); drawCursor shows the left frame
	// normally and the right frame during a click, so it never renders "two
	// hands at once."
	cursorPoint
	cursorCount
)

type uiManager struct {
	font             *engine.BitmapFont
	hoverName        string
	cursor           cursorState
	cursorTex        [cursorCount]*sdl.Texture
	cursorW          [cursorCount]int32
	cursorH          [cursorCount]int32
	cursorTimer      float64
	cursorClickTimer float64
	cursorClicking   bool
}

func newUIManager(font *engine.BitmapFont) *uiManager {
	return &uiManager{font: font}
}

func (ui *uiManager) initCursors(renderer *sdl.Renderer) {
	// Cursor PNGs come out of the image generator with an off-white halo
	// the default tol=8 color-key can't fully lift; use the wider tol=16
	// pass so cursor_talk in particular doesn't render with a visible
	// white square around it.
	load := func(path string) (*sdl.Texture, int32, int32) {
		tex, w, h := engine.TextureFromPNGAggressive(renderer, path)
		if tex != nil {
			tex.SetBlendMode(sdl.BLENDMODE_BLEND)
		}
		return tex, w, h
	}

	ui.cursorTex[cursorNormal], ui.cursorW[cursorNormal], ui.cursorH[cursorNormal] = load("assets/images/ui/cursors/cursor_normal.png")
	ui.cursorTex[cursorTalk], ui.cursorW[cursorTalk], ui.cursorH[cursorTalk] = load("assets/images/ui/cursors/cursor_talk.png")
	ui.cursorTex[cursorGrab], ui.cursorW[cursorGrab], ui.cursorH[cursorGrab] = load("assets/images/ui/cursors/cursor_grab.png")
	ui.cursorTex[cursorArrowLeft], ui.cursorW[cursorArrowLeft], ui.cursorH[cursorArrowLeft] = load("assets/images/ui/cursors/cursor_arrow_left.png")
	ui.cursorTex[cursorArrowRight], ui.cursorW[cursorArrowRight], ui.cursorH[cursorArrowRight] = load("assets/images/ui/cursors/cursor_arrow_right.png")
	ui.cursorTex[cursorArrowUp], ui.cursorW[cursorArrowUp], ui.cursorH[cursorArrowUp] = load("assets/images/ui/cursors/cursor_arrow_up.png")
	ui.cursorTex[cursorArrowDown], ui.cursorW[cursorArrowDown], ui.cursorH[cursorArrowDown] = load("assets/images/ui/cursors/cursor_arrow_down.png")
	ui.cursorTex[cursorArrowDownRight], ui.cursorW[cursorArrowDownRight], ui.cursorH[cursorArrowDownRight] = load("assets/images/ui/cursors/cursor_arrow_down_right.png")
	// cursor_point.png is a 2-frame strip (idle | pressed). Store the FRAME
	// width (half the sheet) so the scaler sizes one finger, not both.
	ui.cursorTex[cursorPoint], ui.cursorW[cursorPoint], ui.cursorH[cursorPoint] = load("assets/images/ui/cursors/cursor_point.png")
	ui.cursorW[cursorPoint] /= 2
}

func (ui *uiManager) triggerClick() {
	ui.cursorClicking = true
	ui.cursorClickTimer = 0.15
}

func (ui *uiManager) updateHover(s *scene, mx, my int32, inv *inventory, plr *player, dt float64) {
	ui.hoverName = ""
	// Default to the grab cursor whenever the player is carrying
	// something so the pointer itself reflects "you are holding an
	// item" even over empty space. Specific hover branches below
	// override this (talk/use/arrow). Without this, the only held-item
	// feedback was the ghost icon drawn beside the cursor, which was
	// easy to miss in busy scenes like camp_grounds.
	if inv != nil && inv.heldItem != nil {
		ui.cursor = cursorGrab
	} else {
		ui.cursor = cursorNormal
	}
	ui.cursorTimer += dt
	if ui.cursorClicking {
		ui.cursorClickTimer -= dt
		if ui.cursorClickTimer <= 0 {
			ui.cursorClicking = false
		}
	}
	for _, n := range s.npcs {
		n.hovered = false
		n.itemMatch = false
	}

	// #9: hovering PP (when not carrying an item) means "open the bag" — the
	// same priority HandleClick uses (PP is checked first). Without this the
	// cursor over PP showed the entrance arrow of an overlapping cabin hotspot,
	// so the player thought a click would enter the cabin instead of opening
	// the inventory. Show the bag/grab cursor so the action reads correctly.
	if plr != nil && inv != nil && inv.heldItem == nil && plr.containsPoint(mx, my) {
		if len(inv.items) > 0 {
			ui.hoverName = "Open Bag"
			ui.cursor = cursorPoint // the pink "action" pointing hand
		}
		return
	}

	// User 2026-05-22: floor items BEFORE npcs so the grab cursor wins
	// when a pickable item sits under an NPC's bounds rect (e.g. the
	// rolling pin on a café table near Nicolas). Previously NPCs were
	// checked first, so hovering a floor item next to an NPC showed
	// the talk cursor — wrong action signal for the player.
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
	for _, n := range s.npcs {
		if n.silent || n.hidden {
			continue
		}
		if n.containsPoint(mx, my) {
			ui.hoverName = n.name
			n.hovered = true
			ui.cursor = cursorTalk
			// itemMatch feedback: pulse on the NPC + swap the cursor tint
			// whenever clicking will actually run the alt dialog. The held
			// path still wins (so drag-onto-NPC gives a stronger cue), but
			// we also light up when the required item is just in the bag —
			// this is what tells the player "you've got what Lily needs,
			// click her to give it" without forcing them to manually draw
			// the flower out first.
			if n.altDialogFunc != nil && inv != nil {
				heldMatches := inv.heldItem != nil &&
					(n.altDialogRequiresItem == "" ||
						inv.heldItem.name == n.altDialogRequiresItem)
				bagMatches := !n.altDialogRequiresHeld &&
					n.altDialogRequiresItem != "" &&
					inv.hasItem(n.altDialogRequiresItem)
				if heldMatches || bagMatches {
					entries, _ := n.altDialogFunc()
					if entries != nil {
						n.itemMatch = true
					}
				}
			}
			return
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
			case arrowDownRight:
				ui.cursor = cursorArrowDownRight
			}
			return
		}
	}
}

func (ui *uiManager) draw(renderer *sdl.Renderer, mx, my int32) {
	renderer.SetDrawColor(0, 0, 0, 140)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: 36})

	txt := "(" + fmt.Sprintf("%d", mx) + ", " + fmt.Sprintf("%d", my) + ")"
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

	srcW := ui.cursorW[c]
	srcH := ui.cursorH[c]

	// Cursor PNGs were updated 2026-05-09 to single-frame portrait images
	// (cursor_normal is a single arrow, not the old 2-frame idle|click
	// strip). Click feedback now comes from the bob/pulse animation only,
	// so we render the full texture for every cursor state.
	var src *sdl.Rect
	// cursorPoint is a 2-frame strip: show the left (idle) frame, or the right
	// (pressed) frame while clicking, so only ONE finger ever renders (#9).
	if c == cursorPoint {
		frame := int32(0)
		if ui.cursorClicking {
			frame = 1
		}
		src = &sdl.Rect{X: frame * srcW, Y: 0, W: srcW, H: srcH}
	}

	targetW := int32(40)
	scale := float64(targetW) / float64(srcW)
	if scale > 1.5 {
		scale = 1.5
	}
	// Brief scale-up while clicking so the click still has a visual cue
	// now that the cursor sheet is single-frame. cursorClickTimer counts
	// down from 0.15s; pulse peaks at +20% in the middle of that window.
	if ui.cursorClicking {
		t := ui.cursorClickTimer / 0.15
		if t < 0 {
			t = 0
		} else if t > 1 {
			t = 1
		}
		// Triangle pulse: 0 -> 1 -> 0 across the window
		pulse := 1.0 - 2.0*math.Abs(t-0.5)
		scale *= 1.0 + 0.20*pulse
	}
	w := int32(float64(srcW) * scale)
	h := int32(float64(srcH) * scale)

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
	case cursorArrowDownRight:
		dy = 0
	case cursorTalk:
		dx = -w / 2
		dy = -h - 4
	case cursorGrab:
		dx = -w / 2
		dy = -h / 2
	case cursorPoint:
		// Fingertip points up; anchor it near the click point.
		dx = -w / 3
		dy = 0
	default:
		dx = 0
		dy = 0
	}

	dst := sdl.Rect{X: mx + dx, Y: my + dy + bob, W: w, H: h}
	renderer.Copy(tex, src, &dst)
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

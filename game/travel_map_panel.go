package game

import (
	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// Info-panel state lives on travelMap so the click path can check
// "is a panel open?" without a second overlay subsystem.
//
// Layout: a 720x380 card centered on screen. Left half shows the location's
// landmark on a cream card; right half shows the name + bulleted facts. The
// globe stays rendered behind the panel so the player keeps spatial context.

const (
	panelW           int32 = 720
	panelH           int32 = 400
	panelImageW      int32 = 300
	panelPad         int32 = 22
	panelFactLineGap int32 = 6
	panelFactScale   int32 = 2
	panelTitleScale  int32 = 3
)

// openInfoPanel switches the map into "reading about this place" mode.
// Click a relevant pin → travels (no panel). Click any other pin → this.
func (tm *travelMap) openInfoPanel(loc *travelLocation) {
	tm.panelLoc = loc
}

// closeInfoPanel dismisses the panel. The map stays visible underneath.
func (tm *travelMap) closeInfoPanel() {
	tm.panelLoc = nil
}

// panelVisible reports whether the info panel is currently overlaid.
func (tm *travelMap) panelVisible() bool {
	return tm.panelLoc != nil
}

// panelHandleClick consumes any click while the panel is open: click
// anywhere dismisses it and falls through to the map (no further action).
// Returns true so the caller stops processing the click.
func (tm *travelMap) panelHandleClick() bool {
	if tm.panelVisible() {
		tm.closeInfoPanel()
		return true
	}
	return false
}

// drawInfoPanel renders the overlay. Called from Game.Draw after the map
// overlay so the panel sits on top of the globe.
func (tm *travelMap) drawInfoPanel(renderer *sdl.Renderer, font *engine.BitmapFont) {
	if !tm.panelVisible() {
		return
	}
	loc := tm.panelLoc

	// Dim the map behind the panel so the text reads cleanly.
	renderer.SetDrawColor(0, 0, 0, 140)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})

	x := (engine.ScreenWidth - panelW) / 2
	y := (engine.ScreenHeight - panelH) / 2

	// Card background: cream paper with a dark charcoal border.
	renderer.SetDrawColor(240, 228, 200, 245)
	renderer.FillRect(&sdl.Rect{X: x, Y: y, W: panelW, H: panelH})
	renderer.SetDrawColor(40, 30, 22, 255)
	renderer.DrawRect(&sdl.Rect{X: x, Y: y, W: panelW, H: panelH})
	renderer.DrawRect(&sdl.Rect{X: x + 2, Y: y + 2, W: panelW - 4, H: panelH - 4})

	imgX := x + panelPad
	imgY := y + panelPad
	imgW := panelImageW
	imgH := panelH - 2*panelPad

	// Left: landmark image on a slightly deeper cream background, with a
	// hairline frame. Falls back to a compass-rose placeholder if the
	// location has no landmark texture (camp, BA, Mexico).
	renderer.SetDrawColor(220, 204, 170, 255)
	renderer.FillRect(&sdl.Rect{X: imgX, Y: imgY, W: imgW, H: imgH})
	renderer.SetDrawColor(80, 60, 40, 200)
	renderer.DrawRect(&sdl.Rect{X: imgX, Y: imgY, W: imgW, H: imgH})

	if loc.landmarkTex != nil {
		// Scale landmark to fill the image area while preserving aspect.
		srcW := float64(loc.landmarkW)
		srcH := float64(loc.landmarkH)
		scale := float64(imgW-20) / srcW
		if h := float64(imgH-20) / srcH; h < scale {
			scale = h
		}
		dstW := int32(srcW * scale)
		dstH := int32(srcH * scale)
		dstX := imgX + (imgW-dstW)/2
		dstY := imgY + (imgH-dstH)/2
		loc.landmarkTex.SetColorMod(255, 255, 255)
		loc.landmarkTex.SetAlphaMod(255)
		renderer.Copy(loc.landmarkTex, nil, &sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH})
	} else {
		// Compass-rose stand-in for landmark-less cities.
		cx := imgX + imgW/2
		cy := imgY + imgH/2
		renderer.SetDrawColor(120, 80, 40, 220)
		renderer.FillRect(&sdl.Rect{X: cx - 3, Y: cy - 60, W: 6, H: 120})
		renderer.FillRect(&sdl.Rect{X: cx - 60, Y: cy - 3, W: 120, H: 6})
		renderer.SetDrawColor(220, 60, 50, 220)
		renderer.FillRect(&sdl.Rect{X: cx - 6, Y: cy - 65, W: 12, H: 12})
	}

	// Right: title + facts.
	textX := x + panelPad*2 + panelImageW
	textMaxW := panelW - panelPad*3 - panelImageW
	titleY := y + panelPad

	// Title underline.
	font.DrawText(renderer, loc.name, textX+1, titleY+1, panelTitleScale,
		sdl.Color{R: 0, G: 0, B: 0, A: 160})
	font.DrawText(renderer, loc.name, textX, titleY, panelTitleScale,
		sdl.Color{R: 40, G: 30, B: 22, A: 255})
	titleH := font.LineHeight(panelTitleScale)
	renderer.SetDrawColor(120, 80, 40, 220)
	renderer.FillRect(&sdl.Rect{X: textX, Y: titleY + titleH + 6, W: textMaxW, H: 2})

	// Facts: word-wrapped per line. Characters per wrap computed from the
	// text area width and the font's char width at the chosen scale.
	factsY := titleY + titleH + 20
	charW := font.TextWidth("M", panelFactScale)
	if charW <= 0 {
		charW = 1
	}
	maxChars := int(textMaxW / charW)
	if maxChars < 10 {
		maxChars = 10
	}

	facts := loc.facts
	if len(facts) == 0 && loc.info != "" {
		facts = []string{loc.info}
	}

	cy := factsY
	for _, fact := range facts {
		for _, line := range wrapLine(fact, maxChars) {
			font.DrawText(renderer, "• "+line, textX+1, cy+1, panelFactScale,
				sdl.Color{R: 0, G: 0, B: 0, A: 140})
			font.DrawText(renderer, "• "+line, textX, cy, panelFactScale,
				sdl.Color{R: 50, G: 38, B: 28, A: 255})
			cy += font.LineHeight(panelFactScale) + panelFactLineGap
		}
		cy += panelFactLineGap
	}

	// Footer hint.
	hint := "Click anywhere or press Esc to close"
	hw := font.TextWidth(hint, 2)
	hx := x + (panelW-hw)/2
	hy := y + panelH - panelPad - font.LineHeight(2)
	font.DrawText(renderer, hint, hx, hy, 2, sdl.Color{R: 100, G: 80, B: 50, A: 200})
}

// wrapLine breaks `s` into lines no longer than `maxChars` runes, breaking
// at spaces when possible. Keeps it simple — no hyphenation or fancy
// word-break heuristics.
func wrapLine(s string, maxChars int) []string {
	if maxChars <= 0 {
		return []string{s}
	}
	var lines []string
	for len(s) > maxChars {
		// Look backwards from maxChars for a space to break at.
		cut := maxChars
		for cut > 0 && s[cut-1] != ' ' {
			cut--
		}
		if cut == 0 {
			cut = maxChars
		}
		lines = append(lines, s[:cut])
		s = s[cut:]
		// Trim leading space on the next line.
		for len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}
	if len(s) > 0 {
		lines = append(lines, s)
	}
	return lines
}

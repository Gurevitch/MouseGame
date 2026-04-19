package game

import (
	"fmt"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// gameMenu is the pause/system menu overlay — Save, Load, Exit. Opened with
// ESC, dismissed with ESC again or by clicking a button. Renders as a
// centered column of large buttons over a dimmed background so the current
// scene stays visible underneath.
//
// Kept local to the game package rather than in ui.go because the buttons
// call back into Game for their actions; a separate subsystem keeps those
// hooks obvious without growing uiManager's surface.
type gameMenu struct {
	visible      bool
	hoveredIndex int // 0..len(items)-1, -1 = none hovered
	status       string
	items        []menuItem
	savePath     string
}

type menuItem struct {
	label  string
	action menuAction
}

type menuAction int

const (
	menuSave menuAction = iota
	menuLoad
	menuResume
	menuExit
)

const (
	menuBtnW = 360
	menuBtnH = 64
	menuBtnGap = 18
)

func newGameMenu() *gameMenu {
	return &gameMenu{
		hoveredIndex: -1,
		savePath:     "savegame.json",
		items: []menuItem{
			{label: "Save Game", action: menuSave},
			{label: "Load Game", action: menuLoad},
			{label: "Resume", action: menuResume},
			{label: "Exit", action: menuExit},
		},
	}
}

func (m *gameMenu) Visible() bool { return m.visible }

func (m *gameMenu) Show() {
	m.visible = true
	m.status = ""
	m.hoveredIndex = -1
}

func (m *gameMenu) Hide() {
	m.visible = false
	m.status = ""
}

func (m *gameMenu) Toggle() {
	if m.visible {
		m.Hide()
	} else {
		m.Show()
	}
}

// buttonRect returns the destination rect for item i, centered horizontally
// with the column vertically centered on the screen.
func (m *gameMenu) buttonRect(i int) sdl.Rect {
	n := int32(len(m.items))
	totalH := n*menuBtnH + (n-1)*menuBtnGap
	top := engine.ScreenHeight/2 - totalH/2
	return sdl.Rect{
		X: engine.ScreenWidth/2 - menuBtnW/2,
		Y: top + int32(i)*(menuBtnH+menuBtnGap),
		W: menuBtnW,
		H: menuBtnH,
	}
}

// UpdateHover refreshes which button the cursor is over. Called from the
// main update loop while the menu is visible.
func (m *gameMenu) UpdateHover(mx, my int32) {
	if !m.visible {
		return
	}
	m.hoveredIndex = -1
	for i := range m.items {
		r := m.buttonRect(int(i))
		if mx >= r.X && mx < r.X+r.W && my >= r.Y && my < r.Y+r.H {
			m.hoveredIndex = i
			return
		}
	}
}

// HandleClick dispatches a click while the menu is visible. Returns true if
// the click was inside a menu button (caller should stop processing),
// false if the click hit dead space (menu stays open, click ignored).
func (g *Game) menuHandleClick(x, y int32) bool {
	if !g.menu.Visible() {
		return false
	}
	for i := range g.menu.items {
		r := g.menu.buttonRect(i)
		if x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H {
			g.applyMenuAction(g.menu.items[i].action)
			return true
		}
	}
	// Click missed all buttons — swallow it so a stray click doesn't
	// also hit the scene behind.
	return true
}

func (g *Game) applyMenuAction(a menuAction) {
	switch a {
	case menuSave:
		if err := g.SaveGame(g.menu.savePath); err != nil {
			g.menu.status = "Save failed: " + err.Error()
		} else {
			g.menu.status = "Saved."
		}
	case menuLoad:
		if err := g.LoadGame(g.menu.savePath); err != nil {
			g.menu.status = "Load failed: " + err.Error()
		} else {
			g.menu.status = "Loaded."
			g.menu.Hide()
		}
	case menuResume:
		g.menu.Hide()
	case menuExit:
		// Post a quit event so main.go's event loop returns cleanly,
		// triggering the deferred sdl.Quit and window.Destroy calls.
		if _, err := sdl.PushEvent(&sdl.QuitEvent{Type: sdl.QUIT, Timestamp: sdl.GetTicks()}); err != nil {
			fmt.Printf("menu: could not post quit event: %v\n", err)
		}
	}
}

// Draw renders the menu overlay. Call last in Game.Draw so the menu sits on
// top of everything else.
func (m *gameMenu) Draw(renderer *sdl.Renderer, font *engine.BitmapFont, mx, my int32) {
	if !m.visible {
		return
	}
	// Dim backdrop
	renderer.SetDrawColor(0, 0, 0, 180)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})

	// Title
	const titleScale = 4
	title := "PAUSED"
	tw := font.TextWidth(title, titleScale)
	titleX := engine.ScreenWidth/2 - tw/2
	titleY := m.buttonRect(0).Y - 90
	font.DrawText(renderer, title, titleX+2, titleY+2, titleScale, sdl.Color{R: 0, G: 0, B: 0, A: 200})
	font.DrawText(renderer, title, titleX, titleY, titleScale, sdl.Color{R: 255, G: 230, B: 210, A: 255})

	// Buttons
	for i, item := range m.items {
		r := m.buttonRect(i)
		fill := sdl.Color{R: 40, G: 30, B: 55, A: 230}
		border := sdl.Color{R: 200, G: 180, B: 220, A: 255}
		text := sdl.Color{R: 240, G: 230, B: 255, A: 255}
		if m.hoveredIndex == i {
			fill = sdl.Color{R: 120, G: 60, B: 140, A: 240}
			border = sdl.Color{R: 255, G: 230, B: 210, A: 255}
			text = sdl.Color{R: 255, G: 255, B: 255, A: 255}
		}
		renderer.SetDrawColor(fill.R, fill.G, fill.B, fill.A)
		renderer.FillRect(&r)
		renderer.SetDrawColor(border.R, border.G, border.B, border.A)
		renderer.DrawRect(&r)
		// Shift inward so the border isn't 1px wide
		inner := sdl.Rect{X: r.X + 2, Y: r.Y + 2, W: r.W - 4, H: r.H - 4}
		renderer.DrawRect(&inner)

		const btnScale = 3
		bw := font.TextWidth(item.label, btnScale)
		bx := r.X + (r.W-bw)/2
		by := r.Y + (r.H-int32(font.LineHeight(btnScale)))/2
		font.DrawText(renderer, item.label, bx+1, by+1, btnScale, sdl.Color{R: 0, G: 0, B: 0, A: 200})
		font.DrawText(renderer, item.label, bx, by, btnScale, text)
	}

	// Status line
	if m.status != "" {
		const s = 2
		sw := font.TextWidth(m.status, s)
		sx := engine.ScreenWidth/2 - sw/2
		sy := m.buttonRect(len(m.items)-1).Y + menuBtnH + 20
		font.DrawText(renderer, m.status, sx+1, sy+1, s, sdl.Color{R: 0, G: 0, B: 0, A: 200})
		font.DrawText(renderer, m.status, sx, sy, s, sdl.Color{R: 230, G: 230, B: 255, A: 255})
	}

	// Hint at bottom
	const hint = "Esc to resume"
	const hs = 2
	hw := font.TextWidth(hint, hs)
	hx := engine.ScreenWidth/2 - hw/2
	hy := m.buttonRect(len(m.items)-1).Y + menuBtnH + 60
	font.DrawText(renderer, hint, hx+1, hy+1, hs, sdl.Color{R: 0, G: 0, B: 0, A: 200})
	font.DrawText(renderer, hint, hx, hy, hs, sdl.Color{R: 200, G: 200, B: 220, A: 220})
}

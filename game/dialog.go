package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	dialogPanelHeight = 150
	dialogPadding     = 20
	typingSpeed       = 30.0
)

type dialogEntry struct {
	speaker string
	text    string
}

type dialogSystem struct {
	font          *engine.BitmapFont
	queue         []dialogEntry
	currentIndex  int
	displayedLen  int
	typingTimer   float64
	active        bool
	showContinue  bool
	continueTimer float64
}

func newDialogSystem(font *engine.BitmapFont) *dialogSystem {
	return &dialogSystem{font: font}
}

func (ds *dialogSystem) startDialog(entries []dialogEntry) {
	if len(entries) == 0 {
		return
	}
	ds.queue = entries
	ds.currentIndex = 0
	ds.displayedLen = 0
	ds.typingTimer = 0
	ds.active = true
	ds.showContinue = false
	ds.continueTimer = 0
}

func (ds *dialogSystem) advance() {
	if !ds.active || len(ds.queue) == 0 {
		return
	}
	current := ds.queue[ds.currentIndex]
	if ds.displayedLen < len(current.text) {
		ds.displayedLen = len(current.text)
		ds.showContinue = true
		return
	}
	ds.currentIndex++
	if ds.currentIndex >= len(ds.queue) {
		ds.active = false
		ds.queue = nil
		return
	}
	ds.displayedLen = 0
	ds.typingTimer = 0
	ds.showContinue = false
}

func (ds *dialogSystem) update(dt float64) {
	if !ds.active || len(ds.queue) == 0 {
		return
	}
	current := ds.queue[ds.currentIndex]
	if ds.displayedLen < len(current.text) {
		ds.typingTimer += dt
		ds.displayedLen = int(ds.typingTimer * typingSpeed)
		if ds.displayedLen >= len(current.text) {
			ds.displayedLen = len(current.text)
			ds.showContinue = true
		}
	}
	ds.continueTimer += dt
}

func (ds *dialogSystem) draw(renderer *sdl.Renderer) {
	if !ds.active || len(ds.queue) == 0 {
		return
	}
	panelY := int32(engine.ScreenHeight - dialogPanelHeight)
	renderer.SetDrawColor(0, 0, 0, 200)
	renderer.FillRect(&sdl.Rect{X: 0, Y: panelY, W: engine.ScreenWidth, H: dialogPanelHeight})
	renderer.SetDrawColor(255, 180, 200, 255)
	renderer.DrawRect(&sdl.Rect{X: 2, Y: panelY + 2, W: engine.ScreenWidth - 4, H: dialogPanelHeight - 4})

	current := ds.queue[ds.currentIndex]
	ds.font.DrawText(renderer, current.speaker, dialogPadding, panelY+14, 3,
		sdl.Color{R: 255, G: 180, B: 200, A: 255})

	visibleText := current.text
	if ds.displayedLen < len(current.text) {
		visibleText = current.text[:ds.displayedLen]
	}
	if len(visibleText) > 0 {
		ds.font.DrawText(renderer, visibleText, dialogPadding, panelY+50, 2,
			sdl.Color{R: 255, G: 255, B: 255, A: 255})
	}
	if ds.showContinue && int(ds.continueTimer*2)%2 == 0 {
		ds.font.DrawText(renderer, ">>>", engine.ScreenWidth-80, panelY+dialogPanelHeight-30, 3,
			sdl.Color{R: 255, G: 180, B: 200, A: 255})
	}
}

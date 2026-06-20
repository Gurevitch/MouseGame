package game

import (
	"fmt"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// walkDebug (F3) overlays the scene's walk geometry so path issues get tuned
// with exact on-screen coordinates instead of guesswork (2026-06-12 #2 -
// third round of "PP is not walking over the painted line"). It draws:
//   - every walkSegment as a yellow line with endpoint coordinates
//   - PP's FOOT point as a green crosshair with live coordinates
//   - the last click's snapped walk target as a red dot with coordinates
//
// Workflow: hit F3, walk the path, screenshot - the labels show exactly
// where the segments are vs where the painted path is, and the segment
// coords can then be set once, correctly.
type walkDebug struct {
	active    bool
	lastSnapX float64
	lastSnapY float64
	haveSnap  bool
}

func newWalkDebug() *walkDebug { return &walkDebug{} }

func (wd *walkDebug) toggle() {
	wd.active = !wd.active
	if wd.active {
		fmt.Println("[walk-debug] ON - F3 again to disable")
	} else {
		fmt.Println("[walk-debug] off")
	}
}

// recordSnap stores the latest click's snapped target while the overlay is on.
func (wd *walkDebug) recordSnap(x, y float64) {
	if wd.active {
		wd.lastSnapX, wd.lastSnapY = x, y
		wd.haveSnap = true
	}
}

func (wd *walkDebug) draw(renderer *sdl.Renderer, font *engine.BitmapFont, s *scene, p *player) {
	if !wd.active || s == nil {
		return
	}
	yellow := sdl.Color{R: 255, G: 220, B: 60, A: 255}
	// Walk segments (drawn twice, 1px apart, so the line reads thick).
	renderer.SetDrawColor(255, 220, 60, 230)
	for _, seg := range s.walkSegments {
		renderer.DrawLine(int32(seg.x1), int32(seg.y1), int32(seg.x2), int32(seg.y2))
		renderer.DrawLine(int32(seg.x1), int32(seg.y1)+1, int32(seg.x2), int32(seg.y2)+1)
	}
	for _, seg := range s.walkSegments {
		font.DrawText(renderer, fmt.Sprintf("%d,%d", int(seg.x1), int(seg.y1)),
			int32(seg.x1)+4, int32(seg.y1)-14, 1, yellow)
	}
	// PP's foot point.
	if p != nil {
		fx, fy := p.footCenter()
		renderer.SetDrawColor(60, 255, 90, 255)
		renderer.DrawLine(fx-12, fy, fx+12, fy)
		renderer.DrawLine(fx, fy-12, fx, fy+12)
		font.DrawText(renderer, fmt.Sprintf("foot %d,%d", fx, fy),
			fx+14, fy-10, 2, sdl.Color{R: 60, G: 255, B: 90, A: 255})
	}
	// Last snapped click target.
	if wd.haveSnap {
		renderer.SetDrawColor(255, 80, 80, 255)
		renderer.FillRect(&sdl.Rect{X: int32(wd.lastSnapX) - 4, Y: int32(wd.lastSnapY) - 4, W: 8, H: 8})
		font.DrawText(renderer, fmt.Sprintf("snap %d,%d", int(wd.lastSnapX), int(wd.lastSnapY)),
			int32(wd.lastSnapX)+10, int32(wd.lastSnapY)-10, 2, sdl.Color{R: 255, G: 100, B: 100, A: 255})
	}
	font.DrawText(renderer, "WALK DEBUG (F3)", engine.ScreenWidth-300, 12, 2, yellow)
}

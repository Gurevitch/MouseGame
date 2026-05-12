package game

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// clickProbe is a dev-mode diagnostic that validates whether NPC sprites have
// a clean transparent background-cut. The current click hit-test
// (npc.containsPoint -> lastDrawRect) is rectangular, so a sprite with a
// sloppy halo of semi-opaque or accidentally-opaque pixels around the
// character will register clicks in dead air. The probe samples the source
// PNG's alpha channel at the clicked screen pixel — if the hit-test said yes
// but the alpha is transparent, the cut is bad and the NPC name is logged.
//
// Toggle with F2. While active, clicks no longer talk to NPCs; instead they
// drop a green or red marker at the click point and print a one-line report.
//
// Marker colors:
//   GREEN — clicked an opaque pixel of an NPC's drawn frame (good cut)
//   RED   — bbox-hit an NPC but the source pixel is transparent (bad cut)
//   GREY  — no NPC under the click (background only)
type clickProbe struct {
	active  bool
	markers []probeMarker
}

type probeMarker struct {
	x, y  int32
	color sdl.Color
	label string
	ttl   float64
}

const (
	probeMarkerLifetimeSeconds = 4.0
	// alpha threshold below which a pixel is considered transparent for
	// hit-test purposes. 16 forgives anti-aliased edges (alpha ~30-200)
	// while still flagging real transparent-halo bugs (alpha 0).
	probeAlphaThreshold uint8 = 16
)

// alphaCache memoizes decoded NRGBA bitmaps keyed by the source PNG path.
// Decoding happens once per file; subsequent samples are O(1) reads.
var (
	alphaCacheMu sync.Mutex
	alphaCache   = map[string]image.Image{}
)

func newClickProbe() *clickProbe { return &clickProbe{} }

func (cp *clickProbe) toggle() {
	cp.active = !cp.active
	if cp.active {
		fmt.Println("[click-probe] ON — F2 again to disable. Clicks won't talk to NPCs.")
		cp.markers = cp.markers[:0]
	} else {
		fmt.Println("[click-probe] off")
	}
}

func (cp *clickProbe) update(dt float64) {
	if !cp.active && len(cp.markers) == 0 {
		return
	}
	out := cp.markers[:0]
	for _, m := range cp.markers {
		m.ttl -= dt
		if m.ttl > 0 {
			out = append(out, m)
		}
	}
	cp.markers = out
}

// recordClick is invoked by the game loop whenever the probe is active and
// the player clicks. It runs the same hit-test the real click handler uses,
// then samples the source PNG to decide whether the bbox-hit is a true hit
// on the cartoon outline.
func (cp *clickProbe) recordClick(s *scene, x, y int32) {
	if !cp.active {
		return
	}
	n := topmostNPCAt(s, x, y)
	if n == nil {
		cp.pushMarker(x, y, sdl.Color{R: 160, G: 160, B: 160, A: 255}, "no NPC")
		fmt.Printf("[click-probe] (%d,%d) no NPC under click\n", x, y)
		return
	}
	frame := n.lastDrawnFrame
	dst := n.lastDrawRect
	alpha, ok := sampleFrameAlpha(frame, dst, n.lastDrawnFlip, x, y)
	if !ok {
		cp.pushMarker(x, y, sdl.Color{R: 200, G: 200, B: 60, A: 255},
			fmt.Sprintf("%s: no sprite path", n.name))
		fmt.Printf("[click-probe] %q hit but no srcPath — can't validate cut\n", n.name)
		return
	}
	if alpha >= probeAlphaThreshold {
		cp.pushMarker(x, y, sdl.Color{R: 60, G: 220, B: 100, A: 255},
			fmt.Sprintf("%s OK a=%d", n.name, alpha))
		fmt.Printf("[click-probe] %q OK at (%d,%d) alpha=%d\n", n.name, x, y, alpha)
		return
	}
	cp.pushMarker(x, y, sdl.Color{R: 230, G: 60, B: 60, A: 255},
		fmt.Sprintf("%s BAD CUT a=%d", n.name, alpha))
	fmt.Printf("[click-probe] BAD CUT on %q: clicked (%d,%d), source pixel alpha=%d, sheet=%s\n",
		n.name, x, y, alpha, frame.srcPath)
}

func (cp *clickProbe) pushMarker(x, y int32, c sdl.Color, label string) {
	cp.markers = append(cp.markers, probeMarker{
		x: x, y: y, color: c, label: label, ttl: probeMarkerLifetimeSeconds,
	})
}

// draw renders the persistent click markers and a small status banner so
// it's obvious the probe is active. Call after the scene + UI so markers
// sit on top.
func (cp *clickProbe) draw(renderer *sdl.Renderer, font *engine.BitmapFont) {
	if !cp.active && len(cp.markers) == 0 {
		return
	}
	if cp.active {
		// Banner top-right
		banner := "CLICK PROBE  (F2 to exit)"
		bw := int32(len(banner))*8*2 + 24
		renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		renderer.SetDrawColor(20, 20, 30, 200)
		renderer.FillRect(&sdl.Rect{X: engine.ScreenWidth - bw - 12, Y: 12, W: bw, H: 36})
		renderer.SetDrawColor(255, 200, 120, 255)
		renderer.DrawRect(&sdl.Rect{X: engine.ScreenWidth - bw - 12, Y: 12, W: bw, H: 36})
		if font != nil {
			font.DrawText(renderer, banner, engine.ScreenWidth-bw, 22, 2,
				sdl.Color{R: 255, G: 220, B: 140, A: 255})
		}
	}
	for _, m := range cp.markers {
		// alpha fades over the marker's lifetime so old clicks dim out.
		alpha := uint8(255.0 * (m.ttl / probeMarkerLifetimeSeconds))
		c := m.color
		c.A = alpha
		renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
		renderer.SetDrawColor(c.R, c.G, c.B, c.A)
		renderer.FillRect(&sdl.Rect{X: m.x - 8, Y: m.y - 8, W: 16, H: 16})
		// White outline for visibility against any background.
		renderer.SetDrawColor(255, 255, 255, alpha)
		renderer.DrawRect(&sdl.Rect{X: m.x - 8, Y: m.y - 8, W: 16, H: 16})
		if font != nil && m.label != "" {
			font.DrawText(renderer, m.label, m.x+12, m.y-6, 2, sdl.Color{R: 255, G: 255, B: 255, A: alpha})
		}
	}
}

// topmostNPCAt mirrors checkNPCClick's iteration but ignores the silent /
// hidden filters. The probe wants to flag bad cuts on every clickable NPC
// regardless of story state, so we iterate everything and skip only NPCs
// that haven't drawn yet (lastDrawRect is zero) — those have no sprite to
// validate.
func topmostNPCAt(s *scene, x, y int32) *npc {
	pt := sdl.Point{X: x, Y: y}
	// Iterate in reverse so a later-drawn NPC (rendered on top) wins.
	for i := len(s.npcs) - 1; i >= 0; i-- {
		n := s.npcs[i]
		if n == nil || n.hidden {
			continue
		}
		if n.lastDrawRect.W == 0 || n.lastDrawRect.H == 0 {
			continue
		}
		if pt.InRect(&n.lastDrawRect) {
			return n
		}
	}
	return nil
}

// sampleFrameAlpha returns the alpha channel value of the source PNG pixel
// that lives under the screen point (sx, sy). dst is the on-screen rect the
// frame was drawn into; flipped indicates horizontal flip so we can mirror
// the lookup. ok=false when the frame has no srcPath or the file is missing.
//
// Mapping math: position the click within dst as a fraction, multiply by the
// frame size, add the frame's source-rect offset (atlas-backed sheets) to
// land on the PNG pixel.
func sampleFrameAlpha(frame npcFrame, dst sdl.Rect, flipped bool, sx, sy int32) (uint8, bool) {
	if frame.srcPath == "" || dst.W <= 0 || dst.H <= 0 {
		return 0, false
	}
	fx, fy, ok := mapScreenToFramePixel(frame, dst, flipped, sx, sy)
	if !ok {
		return 0, false
	}
	img, err := loadAlphaImage(frame.srcPath)
	if err != nil {
		return 0, false
	}
	b := img.Bounds()
	if fx < b.Min.X || fy < b.Min.Y || fx >= b.Max.X || fy >= b.Max.Y {
		return 0, false
	}
	_, _, _, a := img.At(fx, fy).RGBA()
	return uint8(a >> 8), true
}

// mapScreenToFramePixel inverts the dst-rect math in npc.drawScaled to turn
// a screen point into a pixel inside the source PNG. Pulled out so it can
// be unit-tested without an SDL renderer.
func mapScreenToFramePixel(frame npcFrame, dst sdl.Rect, flipped bool, sx, sy int32) (int, int, bool) {
	if dst.W <= 0 || dst.H <= 0 || frame.w <= 0 || frame.h <= 0 {
		return 0, 0, false
	}
	relX := float64(sx-dst.X) / float64(dst.W)
	relY := float64(sy-dst.Y) / float64(dst.H)
	if relX < 0 || relX >= 1 || relY < 0 || relY >= 1 {
		return 0, 0, false
	}
	fx := int(relX * float64(frame.w))
	fy := int(relY * float64(frame.h))
	if flipped {
		fx = int(frame.w) - 1 - fx
	}
	// If this frame is a sub-rect of an atlas, offset into the sheet.
	if frame.src != nil {
		fx += int(frame.src.X)
		fy += int(frame.src.Y)
	}
	return fx, fy, true
}

// loadAlphaImage decodes a PNG once and caches it. Subsequent lookups for
// the same path are O(1). Concurrency-safe because clicks come from the
// game loop's main thread, but the mutex is cheap insurance.
func loadAlphaImage(path string) (image.Image, error) {
	alphaCacheMu.Lock()
	defer alphaCacheMu.Unlock()
	if img, ok := alphaCache[path]; ok {
		return img, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	alphaCache[path] = img
	return img, nil
}

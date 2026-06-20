package game

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

// mapScreenToFramePixel is the load-bearing math for the probe - if it's
// wrong, every alpha sample reads from the wrong pixel and bad cuts go
// undetected (or good cuts are flagged as bad). These tests pin the
// inverse-mapping behaviour the real probe relies on.

func TestMapScreenToFramePixel_StandaloneFrameCenter(t *testing.T) {
	// Frame is 100x200, drawn into dst 200x400 at (50, 60). A click at
	// dst center should map to (50, 100) in the frame.
	frame := npcFrame{w: 100, h: 200}
	dst := sdl.Rect{X: 50, Y: 60, W: 200, H: 400}
	fx, fy, ok := mapScreenToFramePixel(frame, dst, false, 50+100, 60+200)
	if !ok {
		t.Fatalf("expected ok, got false")
	}
	if fx != 50 || fy != 100 {
		t.Errorf("center: got (%d,%d), want (50,100)", fx, fy)
	}
}

func TestMapScreenToFramePixel_OutOfDstReturnsFalse(t *testing.T) {
	frame := npcFrame{w: 100, h: 200}
	dst := sdl.Rect{X: 50, Y: 60, W: 200, H: 400}
	cases := []struct {
		name   string
		sx, sy int32
	}{
		{"left of dst", 49, 100},
		{"right of dst", 250, 100},
		{"above dst", 100, 59},
		{"below dst", 100, 460},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, ok := mapScreenToFramePixel(frame, dst, false, tc.sx, tc.sy); ok {
				t.Errorf("(%d,%d) outside dst should be ok=false", tc.sx, tc.sy)
			}
		})
	}
}

func TestMapScreenToFramePixel_FlippedMirrors(t *testing.T) {
	// At dst 25% from the left, fx should be 25% of frame.w when not
	// flipped, and (frame.w - 1 - 25% of frame.w) when flipped.
	frame := npcFrame{w: 100, h: 200}
	dst := sdl.Rect{X: 0, Y: 0, W: 400, H: 800}
	fx, _, ok := mapScreenToFramePixel(frame, dst, false, 100, 100) // 25% across
	if !ok || fx != 25 {
		t.Fatalf("not flipped: got fx=%d ok=%v, want 25", fx, ok)
	}
	fxFlip, _, ok := mapScreenToFramePixel(frame, dst, true, 100, 100)
	if !ok || fxFlip != 100-1-25 {
		t.Errorf("flipped: got fx=%d ok=%v, want %d", fxFlip, ok, 100-1-25)
	}
}

func TestMapScreenToFramePixel_AtlasFrameOffsetsIntoSheet(t *testing.T) {
	// Atlas-backed frame at sheet rect (300,400)-100x200. Center click
	// should land at sheet pixel (300+50, 400+100) = (350, 500).
	srcRect := sdl.Rect{X: 300, Y: 400, W: 100, H: 200}
	frame := npcFrame{w: 100, h: 200, src: &srcRect}
	dst := sdl.Rect{X: 0, Y: 0, W: 200, H: 400}
	fx, fy, ok := mapScreenToFramePixel(frame, dst, false, 100, 200)
	if !ok {
		t.Fatalf("expected ok")
	}
	if fx != 350 || fy != 500 {
		t.Errorf("atlas offset: got (%d,%d), want (350,500)", fx, fy)
	}
}

func TestMapScreenToFramePixel_ZeroSizeReturnsFalse(t *testing.T) {
	cases := []struct {
		name  string
		frame npcFrame
		dst   sdl.Rect
	}{
		{"zero frame width", npcFrame{w: 0, h: 200}, sdl.Rect{W: 100, H: 200}},
		{"zero frame height", npcFrame{w: 100, h: 0}, sdl.Rect{W: 100, H: 200}},
		{"zero dst width", npcFrame{w: 100, h: 200}, sdl.Rect{W: 0, H: 200}},
		{"zero dst height", npcFrame{w: 100, h: 200}, sdl.Rect{W: 100, H: 0}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, ok := mapScreenToFramePixel(tc.frame, tc.dst, false, 10, 10); ok {
				t.Errorf("zero-sized input should return ok=false")
			}
		})
	}
}

// Marker lifecycle tests - make sure stale markers actually get evicted.

func TestClickProbe_MarkerExpires(t *testing.T) {
	cp := newClickProbe()
	cp.active = true
	cp.pushMarker(100, 100, sdl.Color{R: 60, G: 220, B: 100, A: 255}, "test")
	if len(cp.markers) != 1 {
		t.Fatalf("expected 1 marker, got %d", len(cp.markers))
	}
	// Tick forward past the lifetime - marker should be gone.
	cp.update(probeMarkerLifetimeSeconds + 0.1)
	if len(cp.markers) != 0 {
		t.Errorf("marker should have expired, still %d alive", len(cp.markers))
	}
}

func TestClickProbe_MarkerStaysWithinLifetime(t *testing.T) {
	cp := newClickProbe()
	cp.active = true
	cp.pushMarker(100, 100, sdl.Color{R: 60, G: 220, B: 100, A: 255}, "test")
	cp.update(probeMarkerLifetimeSeconds * 0.5) // half life
	if len(cp.markers) != 1 {
		t.Errorf("marker should survive half its lifetime, got %d", len(cp.markers))
	}
}

func TestClickProbe_ToggleClearsMarkers(t *testing.T) {
	cp := newClickProbe()
	cp.active = true
	cp.pushMarker(100, 100, sdl.Color{}, "a")
	cp.pushMarker(200, 200, sdl.Color{}, "b")
	cp.toggle() // off
	cp.toggle() // on - should also clear leftovers
	if len(cp.markers) != 0 {
		t.Errorf("re-enabling probe should clear stale markers, got %d", len(cp.markers))
	}
}

func TestSampleFrameAlpha_NoSrcPathReturnsNotOk(t *testing.T) {
	// A frame that came from a loader without path tracking shouldn't
	// crash the probe; it should just decline to sample.
	frame := npcFrame{w: 100, h: 100, srcPath: ""}
	dst := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	if _, ok := sampleFrameAlpha(frame, dst, false, 50, 50); ok {
		t.Errorf("expected ok=false when srcPath is empty")
	}
}

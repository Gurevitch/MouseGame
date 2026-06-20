package game

import "testing"

// background.update is pure logic - no SDL calls - so we can construct a
// background literal directly and exercise the frame-advance loop without a
// renderer. These tests pin the contract the airplane_flight cloud loop
// depends on.

func TestBackgroundUpdate_StaticIsNoOp(t *testing.T) {
	cases := []struct {
		name         string
		frames       int
		frameSeconds float64
	}{
		{"zero frames", 0, 0.15},
		{"single frame", 1, 0.15},
		{"frames set but no interval", 6, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &background{frames: tc.frames, frameSeconds: tc.frameSeconds}
			b.update(1.0)
			if b.frameIdx != 0 || b.frameTimer != 0 {
				t.Errorf("static bg mutated: idx=%d timer=%v", b.frameIdx, b.frameTimer)
			}
		})
	}
}

func TestBackgroundUpdate_AdvancesOneFrame(t *testing.T) {
	b := &background{frames: 6, frameSeconds: 0.15}
	b.update(0.15)
	if b.frameIdx != 1 {
		t.Fatalf("frameIdx after one tick: got %d, want 1", b.frameIdx)
	}
}

func TestBackgroundUpdate_AccumulatesSubFrameDt(t *testing.T) {
	b := &background{frames: 6, frameSeconds: 0.15}
	b.update(0.10)
	if b.frameIdx != 0 {
		t.Errorf("should not advance before threshold: idx=%d", b.frameIdx)
	}
	b.update(0.06)
	if b.frameIdx != 1 {
		t.Errorf("should advance after accumulated 0.16 >= 0.15: idx=%d", b.frameIdx)
	}
}

func TestBackgroundUpdate_WrapsAtFrameCount(t *testing.T) {
	b := &background{frames: 6, frameSeconds: 0.15}
	// Advance 6 frames + a hair → should land back at index 0.
	for i := 0; i < 6; i++ {
		b.update(0.15)
	}
	if b.frameIdx != 0 {
		t.Errorf("expected wrap to 0 after 6 ticks, got %d", b.frameIdx)
	}
}

func TestBackgroundUpdate_LargeDtAdvancesMultipleFrames(t *testing.T) {
	// One huge dt (e.g. window was minimized then restored) should not
	// drop frames - the for-loop in update() catches up.
	b := &background{frames: 6, frameSeconds: 0.15}
	b.update(0.46) // 3 full frames + 0.01 leftover
	if b.frameIdx != 3 {
		t.Errorf("expected idx=3 after 0.46s, got %d", b.frameIdx)
	}
	// Floating point: leftover should be ~0.01.
	if b.frameTimer < 0.005 || b.frameTimer > 0.02 {
		t.Errorf("leftover timer drifted: %v", b.frameTimer)
	}
}

func TestBackgroundUpdate_LargeDtWrapsCleanly(t *testing.T) {
	// dt big enough to wrap past the frame count in a single tick - the
	// loop must keep going, not stop at the 6-frame boundary. Using a
	// non-exact multiple of frameSeconds (1.25 s vs 0.15) avoids the
	// float-subtraction edge where the residual lands a hair under
	// threshold; in real frames dt is never an exact multiple anyway.
	b := &background{frames: 6, frameSeconds: 0.15}
	b.update(1.25) // 1.25 / 0.15 = 8.33 advances → 8 mod 6 → idx=2
	if b.frameIdx != 2 {
		t.Errorf("expected idx=2 after ~8 advances mod 6, got %d", b.frameIdx)
	}
}

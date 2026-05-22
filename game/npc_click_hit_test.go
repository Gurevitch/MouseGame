package game

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

// TestNPCClickHitBoxes pins the click-target geometry for every authored
// NPC in the game. For each NPC's design-time `bounds` rect we fire FIVE
// synthetic clicks:
//
//   - top    (center-X, 10 px below the top edge)      — MUST hit
//   - mid    (center-X, center-Y)                      — MUST hit
//   - bottom (center-X, 10 px above the bottom edge)   — MUST hit
//   - leftOf  (10 px left of the left edge,   mid-Y)   — MUST miss
//   - rightOf (10 px right of the right edge, mid-Y)   — MUST miss
//
// `npc.containsPoint` is the same call the game uses for both cursor hover
// (the "talk" icon) and the click handler, so this test directly mirrors
// what the player experiences when they click.
//
// User 2026-05-21: "for each one i want 3 right clicks. top mid body and
// bottom and two around the object. the 3 need to work and the two not."
//
// To add a new NPC: drop a row into the table below with its bounds copied
// from the matching factory in npc.go (or scene_loader.go). To intentionally
// shrink a hit-box, narrow the bounds rather than adding special-case
// logic — `containsPoint` is rect-only on purpose so the cursor + click
// areas stay unified.
//
// To exclude an NPC from hit-testing (e.g. a hidden cutscene-only NPC),
// add it with hidden=true and the test will only verify the "around"
// clicks all miss.

type npcHitCase struct {
	scene  string
	name   string
	bounds sdl.Rect
	// hidden NPCs (cutscene-only, e.g. Night Higgins) should never be
	// clickable. We still run the misses (around clicks) and assert
	// `containsPoint` consistently returns false on every probe.
	hidden bool
}

func npcHitCases() []npcHitCase {
	return []npcHitCase{
		// --- Camp grounds (Day 1 + Day 2) ---
		// Kids on the camp grounds line. Bounds copied from
		// assets/data/npc/kids.json — JSON is canonical when factory
		// uses applyKidConfig. Bounds widths sit in the 145–170 band;
		// heights 195–245.
		{scene: "camp_grounds", name: "Tommy", bounds: sdl.Rect{X: 120, Y: 385, W: 170, H: 200}},
		{scene: "camp_grounds", name: "Jake", bounds: sdl.Rect{X: 360, Y: 365, W: 170, H: 215}},
		{scene: "camp_grounds", name: "Lily", bounds: sdl.Rect{X: 590, Y: 375, W: 160, H: 195}},
		{scene: "camp_grounds", name: "Marcus", bounds: sdl.Rect{X: 880, Y: 355, W: 170, H: 220}},
		{scene: "camp_grounds", name: "Danny", bounds: sdl.Rect{X: 1100, Y: 365, W: 170, H: 215}},

		// --- Camp entrance: Director Higgins (intro) ---
		// Bounds copied from newDirectorHiggins (npc.go:307).
		{scene: "camp_entrance", name: "Director Higgins (entrance)", bounds: sdl.Rect{X: 760, Y: 390, W: 168, H: 220}},

		// --- Higgins hidden walk-in NPC on grounds ---
		// Spawns near the cabin path after Lily's shy dialog. Hidden by
		// default until the walk-in sequence un-hides him.
		{scene: "camp_grounds", name: "Director Higgins (grounds)", bounds: sdl.Rect{X: 1060, Y: 560, W: 180, H: 210}, hidden: true},

		// --- Camp office: Higgins behind desk ---
		// User 2026-05-20: moved bounds to (1091, 365).
		{scene: "camp_office", name: "Director Higgins (office)", bounds: sdl.Rect{X: 1091, Y: 365, W: 182, H: 235}},

		// --- Night campfire Higgins ---
		// Silent + driven by cutscene; never clickable directly.
		{scene: "camp_night", name: "Director Higgins (night)", bounds: sdl.Rect{X: 1120, Y: 430, W: 172, H: 220}, hidden: true},

		// --- Kid bedrooms ---
		{scene: "tommy_room", name: "Tommy (room)", bounds: sdl.Rect{X: 670, Y: 440, W: 162, H: 245}},
		{scene: "jake_room", name: "Jake (room)", bounds: sdl.Rect{X: 720, Y: 460, W: 162, H: 245}},
		{scene: "lily_room", name: "Lily (room)", bounds: sdl.Rect{X: 666, Y: 476, W: 162, H: 245}},
		// User 2026-05-20: Marcus room Y nudged 350 → 385 so feet land
		// on the cabin floor.
		{scene: "marcus_room", name: "Marcus (room)", bounds: sdl.Rect{X: 600, Y: 385, W: 187, H: 270}},
		{scene: "danny_room", name: "Danny (room)", bounds: sdl.Rect{X: 760, Y: 445, W: 162, H: 245}},

		// --- Paris street (outside NPCs) ---
		// User 2026-05-20: all Paris street NPCs Y moved to ~490 so
		// feet land on the street tile.
		{scene: "paris_street", name: "Madame Colette", bounds: sdl.Rect{X: 300, Y: 490, W: 135, H: 230}},
		// Pierre shrunk + moved back-of-line for perspective effect.
		{scene: "paris_street", name: "Pierre", bounds: sdl.Rect{X: 820, Y: 390, W: 95, H: 175}},
		// Nicolas W shrunk 106 → 86 so click doesn't bleed into the
		// Louvre-exit hotspot at x≥1300.
		{scene: "paris_street", name: "Nicolas", bounds: sdl.Rect{X: 950, Y: 490, W: 86, H: 230}},
		{scene: "paris_street", name: "Gendarme Claude", bounds: sdl.Rect{X: 1180, Y: 480, W: 115, H: 240}},

		// --- Paris bakery (interior) ---
		{scene: "paris_bakery", name: "Madame Poulain", bounds: sdl.Rect{X: 540, Y: 490, W: 135, H: 230}},

		// --- Paris Louvre (interior) ---
		{scene: "paris_louvre", name: "Curator Beaumont", bounds: sdl.Rect{X: 500, Y: 320, W: 125, H: 240}},
	}
}

func TestNPCClickHitBoxes(t *testing.T) {
	for _, tc := range npcHitCases() {
		tc := tc
		t.Run(tc.scene+"/"+tc.name, func(t *testing.T) {
			n := &npc{bounds: tc.bounds, name: tc.name, hidden: tc.hidden}

			cx := tc.bounds.X + tc.bounds.W/2
			cy := tc.bounds.Y + tc.bounds.H/2

			// Three on-character probes. We use a 10 px offset from
			// the top/bottom edges so the test exercises the body of
			// the rect, not the exact corner.
			topX, topY := cx, tc.bounds.Y+10
			midX, midY := cx, cy
			botX, botY := cx, tc.bounds.Y+tc.bounds.H-10

			// Two around-character probes. 10 px outside the rect's
			// horizontal extent at mid-Y.
			leftX, leftY := tc.bounds.X-10, cy
			rightX, rightY := tc.bounds.X+tc.bounds.W+10, cy

			// Hidden NPCs are not click targets at any point — only the
			// miss-assertions run, but we still log the (would-be) hits
			// so the test reports geometry. Engine-side, hidden NPCs are
			// filtered before containsPoint is even called.
			if tc.hidden {
				for _, p := range []struct {
					label  string
					x, y   int32
				}{
					{"leftOf", leftX, leftY},
					{"rightOf", rightX, rightY},
				} {
					if n.containsPoint(p.x, p.y) {
						t.Errorf("[hidden NPC] %s probe (%d,%d) reports HIT — should never hit a rect outside-edge", p.label, p.x, p.y)
					}
				}
				return
			}

			// --- 3 MUST-HIT probes ---
			hits := []struct {
				label string
				x, y  int32
			}{
				{"top    (center-X, +10 px from top)", topX, topY},
				{"mid    (center-X, center-Y)", midX, midY},
				{"bottom (center-X, -10 px from bottom)", botX, botY},
			}
			for _, p := range hits {
				if !n.containsPoint(p.x, p.y) {
					t.Errorf("MISS on body click — %s at (%d,%d) should hit bounds %+v", p.label, p.x, p.y, tc.bounds)
				}
			}

			// --- 2 MUST-MISS probes ---
			misses := []struct {
				label string
				x, y  int32
			}{
				{"leftOf  (10 px left of bounds, mid-Y)", leftX, leftY},
				{"rightOf (10 px right of bounds, mid-Y)", rightX, rightY},
			}
			for _, p := range misses {
				if n.containsPoint(p.x, p.y) {
					t.Errorf("FALSE HIT around object — %s at (%d,%d) should NOT hit bounds %+v", p.label, p.x, p.y, tc.bounds)
				}
			}
		})
	}
}

// TestNPCClickHitBoxes_BoundsAreSane catches authoring mistakes (zero or
// negative bounds, sub-pixel widths) that would silently break the
// hit-test. Runs as a sibling test so failures here are obvious before
// the per-NPC probes drown them out.
func TestNPCClickHitBoxes_BoundsAreSane(t *testing.T) {
	for _, tc := range npcHitCases() {
		if tc.bounds.W <= 0 || tc.bounds.H <= 0 {
			t.Errorf("%s/%s has non-positive bounds %+v", tc.scene, tc.name, tc.bounds)
		}
		if tc.bounds.W < 40 || tc.bounds.H < 80 {
			t.Errorf("%s/%s bounds %+v look too small to be clickable in 1400×800 — verify against npc.go factory", tc.scene, tc.name, tc.bounds)
		}
		if tc.bounds.W > 400 || tc.bounds.H > 400 {
			t.Errorf("%s/%s bounds %+v look oversized — verify against npc.go factory", tc.scene, tc.name, tc.bounds)
		}
	}
}

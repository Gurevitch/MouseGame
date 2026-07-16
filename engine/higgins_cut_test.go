package engine

import "testing"

// TestOfficeHigginsCutClean locks in the office-Higgins fix: his idle/talk
// sheets are 12 touching bust poses that mis-gap-detect, so the game cuts them
// as EQUAL cells with a high colour-key tol (24) to break the faint bridge
// between neighbouring poses. This test asserts that cut stays clean — every
// one of the 12 cells holds a single, properly-sized pose:
//   - no SLIVER (a cell with almost no content — the old "blink"), and
//   - no full-width BLEED (content spanning edge-to-edge = a neighbour pose
//     leaking in, the "tiny parts from the frame to his left").
// Keep the inset/tol in sync with higginsOfficeInset / higginsOfficeTol in
// game/npc.go.
func TestOfficeHigginsCutClean(t *testing.T) {
	const (
		inset = 4
		tol   = 24
		cols  = 6
		rows  = 2
	)
	sheets := []string{
		"../assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png",
		"../assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png",
	}
	for _, p := range sheets {
		img, err := loadPNG(p)
		if err != nil {
			t.Skipf("skip %s: %v", p, err)
			continue
		}
		applyColorKeyConnectedTol(img, tol)
		eraseGridLines(img, cols, rows)
		bounds := img.Bounds()
		bleed := 0
		for r := 0; r < rows; r++ {
			for c := 0; c < cols; c++ {
				cell := gridCellRect(bounds, cols, rows, c, r, inset)
				ox, _, ow, _ := opaqueLocal(img, cell)
				cw := int32(cell.Dx())
				// Empty / sliver are the unambiguous "blink" defects (the
				// gap-detect regression cut some cells to ~26px). A real bust
				// fills well over a third of the cell.
				if ow <= 0 {
					t.Errorf("%s f%d,%d: empty cell (blank/sliver)", p, r, c)
					continue
				}
				if ow < cw/3 {
					t.Errorf("%s f%d,%d: sliver — ow=%d of cellW=%d", p, r, c, ow, cw)
				}
				// Full-width content = a neighbour pose leaking across the
				// boundary OR a legitimately wide gesture (one open-palm pose,
				// f0,3, spans the cell at any tol). Count them: at the game's
				// tol 24 only that one pose qualifies; if the colour-key tol is
				// ever lowered the bridge returns and 4+ frames span — that's
				// the regression we want to catch.
				if ox <= 2 && int32(ox+ow) >= cw-2 {
					bleed++
				}
			}
		}
		if bleed > 1 {
			t.Errorf("%s: %d full-width frames (expected <=1 wide gesture) — the touching-pose bridge is back; raise higginsOfficeTol", p, bleed)
		}
	}
}

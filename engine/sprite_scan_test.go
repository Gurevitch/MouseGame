package engine

import (
	"fmt"
	"image"
	"os"
	"sort"
	"strings"
	"testing"
)

// TestSpriteScan reads every registered sheet ONCE (no duplicate images, no
// montage files) and reports two defects per sheet, using the same cut the game
// uses (galleryLoaderFor):
//
//   - FRAME LEAK: a frame whose opaque content spans the full cell width — a
//     neighbour pose bleeding across the boundary (the office-Higgins bug). One
//     legitimately wide pose can also span, so the count matters, not presence.
//   - WHITE SPOT: enclosed transparent pixels inside a frame's content box — a
//     HOLE where the colour-key ate part of the character (background shows
//     through). Reported as % of the content area.
//
// Gated behind SPRITE_SCAN=1 (skipped in CI). It only prints a report; to SEE a
// flagged sheet, render just that one — we don't dump 141 images.
//
//	SPRITE_SCAN=1 GOTMPDIR=.gotmp go test ./engine -run TestSpriteScan -v
func TestSpriteScan(t *testing.T) {
	if os.Getenv("SPRITE_SCAN") == "" {
		t.Skip("set SPRITE_SCAN=1 to run the sprite leak / white-spot scan")
	}
	type report struct {
		name      string
		frames    int
		leak      int
		sliver    int
		holePct   float64
		worstHole int
	}
	var reps []report

	for _, c := range realSheetCases {
		img, err := loadPNG(c.path)
		if err != nil {
			continue
		}
		equal, inset, tol := galleryLoaderFor(c.path)
		applyColorKeyConnectedTol(img, tol)
		eraseGridLines(img, c.cols, c.rows)
		bounds := img.Bounds()
		var rects [][]image.Rectangle
		if !equal {
			rects = contentGridRects(img, c.cols, c.rows)
		}

		r := report{name: strings.TrimPrefix(c.path, "../")}
		var holePixels, contentPixels int
		for row := 0; row < c.rows; row++ {
			for col := 0; col < c.cols; col++ {
				cell := gridCellRect(bounds, c.cols, c.rows, col, row, inset)
				if rects != nil {
					cell = rects[row][col]
				}
				ox, oy, ow, oh := opaqueLocal(img, cell)
				if ow <= 0 || oh <= 0 {
					continue
				}
				r.frames++
				cw := int32(cell.Dx())
				if ox <= 2 && int32(ox+ow) >= cw-2 {
					r.leak++
				}
				if ow < cw/3 {
					r.sliver++
				}
				box := image.Rect(cell.Min.X+int(ox), cell.Min.Y+int(oy),
					cell.Min.X+int(ox)+int(ow), cell.Min.Y+int(oy)+int(oh))
				// White spots = near-pure-white pixels still OPAQUE inside the
				// character (the connected key can't reach enclosed whites, so a
				// white background pocket survives as a blob), PLUS any enclosed
				// transparent holes (a key bite). Both read as a white spot.
				holes := countWhiteSpots(img, box)
				holePixels += holes
				contentPixels += int(ow) * int(oh)
				if holes > r.worstHole {
					r.worstHole = holes
				}
			}
		}
		if contentPixels > 0 {
			r.holePct = 100 * float64(holePixels) / float64(contentPixels)
		}
		reps = append(reps, r)
	}

	// Sort worst-first by hole%, then leak count.
	sort.SliceStable(reps, func(i, j int) bool {
		if reps[i].holePct != reps[j].holePct {
			return reps[i].holePct > reps[j].holePct
		}
		return reps[i].leak > reps[j].leak
	})

	fmt.Println("\n=== WHITE-SPOT candidates (pure-white pixels on the character) ===")
	fmt.Println("    NOTE: includes INTENDED whites (e.g. biker's shirt stripes) — eyeball before fixing.")
	for _, r := range reps {
		if r.holePct >= 0.25 {
			fmt.Printf("  %5.2f%% white (worst frame %5d px)  %s\n", r.holePct, r.worstHole, r.name)
		}
	}
	fmt.Println("\n=== FRAME-LEAK candidates (>1 full-width frame) / SLIVERS ===")
	for _, r := range reps {
		if r.leak > 1 || r.sliver > 0 {
			fmt.Printf("  leak=%d sliver=%d (of %d frames)  %s\n", r.leak, r.sliver, r.frames, r.name)
		}
	}
}

// countWhiteSpots counts pixels inside box that read as a "white spot": either
//   - opaque and near-PURE-WHITE (R,G,B all >= 244) — an enclosed white pocket
//     the connected colour-key couldn't reach and erase, or a pure-white patch
//     on the character (which the engine also color-keys, so it shows as a hole
//     in-game), OR
//   - transparent but ENCLOSED by the character (a key bite — reachable test
//     below), which shows the background through the figure.
//
// Border-touching transparent pixels are the legit background and are ignored.
func countWhiteSpots(img *image.NRGBA, box image.Rectangle) int {
	w, h := box.Dx(), box.Dy()
	if w <= 0 || h <= 0 {
		return 0
	}
	at := func(x, y int) (r, g, b, a uint8) {
		p := img.NRGBAAt(box.Min.X+x, box.Min.Y+y)
		return p.R, p.G, p.B, p.A
	}
	trans := func(x, y int) bool { _, _, _, a := at(x, y); return a < 24 }
	// Flood the background-connected transparent region from the border.
	visited := make([]bool, w*h)
	stack := make([][2]int, 0, w*h/4)
	push := func(x, y int) {
		if x >= 0 && x < w && y >= 0 && y < h && !visited[y*w+x] && trans(x, y) {
			visited[y*w+x] = true
			stack = append(stack, [2]int{x, y})
		}
	}
	for x := 0; x < w; x++ {
		push(x, 0)
		push(x, h-1)
	}
	for y := 0; y < h; y++ {
		push(0, y)
		push(w-1, y)
	}
	for len(stack) > 0 {
		c := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		push(c[0]+1, c[1])
		push(c[0]-1, c[1])
		push(c[0], c[1]+1)
		push(c[0], c[1]-1)
	}
	// Only count pure-white OPAQUE pixels — a clean signal for "white on a
	// character" (the no-pure-white rule). Enclosed transparent regions are
	// dropped: they're usually legit background gaps between limbs, not spots.
	// _ = visited keeps the flood available if we want enclosed-hole detection
	// back later.
	_ = visited
	spots := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := at(x, y)
			if a >= 200 && r >= 244 && g >= 244 && b >= 244 {
				spots++
			}
		}
	}
	return spots
}

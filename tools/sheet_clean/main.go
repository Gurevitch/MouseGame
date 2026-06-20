// sheet_clean erases "ghost" pieces from sprite sheets: detached limbs the
// generator painted inside a cell (a hand from the previous frame, a duplicate
// arm) and spill-over from shapes that cross cell borders. For each cell it
// keeps ONLY the largest connected foreground component and paints everything
// else back to the background color.
//
// ONLY run on sheets where the character is a single connected body with NO
// legitimate separate props (no thrown maps, handed baguettes, pigeons...).
// The manifest below is that allowlist.
//
// Run: go run ./tools/sheet_clean
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type sheet struct {
	path string
	cols int
	rows int
}

// Allowlist: prop-free sheets flagged with GHOST PIECES / border spill by
// tools/jitter_audit on 2026-06-10.
var sheets = []sheet{
	{"assets/images/player/PP talk front.png", 8, 2},
	{"assets/images/player/PP idle side.png", 8, 2},
	{"assets/images/player/PP celebrate.png", 8, 2},
	// Marcus sheets REMOVED from the allowlist (2026-06-12 #4): these are the
	// restored OLD sheets whose figures stray outside the fixed cells the
	// cleaner assumes - cleaning erased real body parts (talk frames went
	// blank in-game). Only re-add once the §JIT-MARCUS regens land.
	{"assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png", 8, 2},
	{"assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png", 6, 2},
	{"assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png", 6, 2},
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func near(a, b uint8, tol int) bool {
	return absInt(int(a)-int(b)) <= tol
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	_ = os.Chdir(root)

	for _, s := range sheets {
		f, err := os.Open(s.path)
		if err != nil {
			fmt.Printf("SKIP (missing): %s\n", s.path)
			continue
		}
		src, err := png.Decode(f)
		f.Close()
		if err != nil {
			fmt.Printf("SKIP (bad png): %s\n", s.path)
			continue
		}
		b := src.Bounds()
		img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
		W, H := img.Bounds().Dx(), img.Bounds().Dy()
		cw, ch := W/s.cols, H/s.rows

		bg := img.RGBAAt(0, 0)
		tol := 24
		isFg := func(x, y int) bool {
			p := img.RGBAAt(x, y)
			if p.A < 0x40 {
				return false
			}
			return !(near(p.R, bg.R, tol) && near(p.G, bg.G, tol) && near(p.B, bg.B, tol))
		}

		erasedTotal := 0
		cellsTouched := 0
		for r := 0; r < s.rows; r++ {
			for c := 0; c < s.cols; c++ {
				x0, y0 := c*cw, r*ch
				// label components, tracking each component's bbox
				label := make([]int, cw*ch)
				type comp struct {
					area                   int
					minX, minY, maxX, maxY int
				}
				var comps []comp
				var stack []int
				next := 0
				for i := 0; i < cw*ch; i++ {
					if label[i] != 0 {
						continue
					}
					if !isFg(x0+i%cw, y0+i/cw) {
						label[i] = -1
						continue
					}
					next++
					cp := comp{minX: cw, minY: ch, maxX: -1, maxY: -1}
					stack = append(stack[:0], i)
					label[i] = next
					for len(stack) > 0 {
						p := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						cp.area++
						x, y := p%cw, p/cw
						if x < cp.minX {
							cp.minX = x
						}
						if x > cp.maxX {
							cp.maxX = x
						}
						if y < cp.minY {
							cp.minY = y
						}
						if y > cp.maxY {
							cp.maxY = y
						}
						try := func(q, qx, qy int) {
							if label[q] == 0 {
								if isFg(x0+qx, y0+qy) {
									label[q] = next
									stack = append(stack, q)
								} else {
									label[q] = -1
								}
							}
						}
						if x > 0 {
							try(p-1, x-1, y)
						}
						if x < cw-1 {
							try(p+1, x+1, y)
						}
						if y > 0 {
							try(p-cw, x, y-1)
						}
						if y < ch-1 {
							try(p+cw, x, y+1)
						}
					}
					comps = append(comps, cp)
				}
				if len(comps) < 2 {
					continue
				}
				// find largest (= the character's body)
				largest := 0
				for i := range comps {
					if comps[i].area > comps[largest].area {
						largest = i
					}
				}
				body := comps[largest]
				// Erase ONLY components that stick OUTSIDE the body's bbox
				// (a margin inside it is fine). Interior details — belly
				// shading, eye glints, shirt patches — live fully INSIDE the
				// body bbox and must never be touched (user 2026-06-10: v1
				// erased PP's belly details, leaving see-through holes).
				erase := make([]bool, len(comps))
				const margin = 4
				for i := range comps {
					if i == largest {
						continue
					}
					cp := comps[i]
					inside := cp.minX >= body.minX-margin && cp.maxX <= body.maxX+margin &&
						cp.minY >= body.minY-margin && cp.maxY <= body.maxY+margin
					if !inside {
						erase[i] = true
					}
				}
				erased := 0
				for i := 0; i < cw*ch; i++ {
					if label[i] > 0 && label[i]-1 != largest && erase[label[i]-1] {
						img.SetRGBA(x0+i%cw, y0+i/cw, color.RGBA{bg.R, bg.G, bg.B, 255})
						erased++
					}
				}
				if erased > 0 {
					cellsTouched++
					erasedTotal += erased
				}
			}
		}

		if erasedTotal == 0 {
			fmt.Printf("CLEAN already: %s\n", s.path)
			continue
		}
		out, err := os.Create(s.path)
		if err != nil {
			fmt.Printf("ERROR writing %s: %v\n", s.path, err)
			continue
		}
		if err := png.Encode(out, img); err != nil {
			fmt.Printf("ERROR encoding %s: %v\n", s.path, err)
		}
		out.Close()
		fmt.Printf("CLEANED %-70s erased %6d px of ghost pieces in %d cells\n", s.path, erasedTotal, cellsTouched)
	}
}

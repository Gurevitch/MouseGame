// bg_bake bakes background removal INTO sprite-sheet PNGs (writes alpha in
// place), for sheets whose backgrounds the runtime color-key can't fully
// lift: streaky/gradient backgrounds, enclosed bg-colored pockets between a
// character's arm and body, and near-white AA halo specks floating around
// the figure (2026-08-01 PR: pp_sleeping "white spots", pp_waking "blue
// spots", PP talk side bg, npc_pierre_get_baguette "blue spots").
//
// Per sheet it runs three passes:
//  1. edge flood — background-colored pixels reachable from the sheet edge
//     go transparent (same rule as the runtime connected key);
//  2. enclosed pockets — CLOSED regions still matching the bg color go
//     transparent only when SMALL (area <= enclosedMax), so big legit art
//     like PP's ivory belly or the blue blanket survives;
//  3. specks — tiny leftover floating components (area <= speckMax) that are
//     near-white or bg-colored go transparent (AA halo junk).
//
// Usage:  go run ./tools/bg_bake .          # bake the allowlist below
//         go run ./tools/bg_bake -dry .     # report only, write nothing
//
// The allowlist is intentionally explicit — add a sheet only after eyeballing
// which passes it needs (enclosedMax 0 disables pass 2 for sheets whose
// enclosed whites must survive).
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

type sheetCfg struct {
	path        string
	tol         uint8 // per-channel bg match tolerance
	enclosedMax int   // max area of an enclosed bg-colored pocket to remove; 0 = skip pass
	speckMax    int   // max area of a floating near-white/bg speck to remove; 0 = skip pass
	fringeIter  int   // defringe iterations: opaque px next to transparency that
	// still match the bg within tol+20 go transparent. 0 = skip; -1 = AUTO
	// (2 iterations when the bg is bluish, 0 on near-white backgrounds where
	// defringe would eat thin pale strokes like whiskers).
	global bool // true = key EVERY bg-matching pixel (no connectivity), the
	// runtime GLOBAL-key equivalent. Use for sheets whose loader is
	// gridFrames/gridFramesTol — after baking, the runtime key no-ops on the
	// transparent corners, so a flood-only bake would REGRESS the enclosed
	// pockets the global key used to remove (user 2026-08-01: the tiny white
	// block between arm and body on PP talk side frame r?c4).
	pocketTol uint8 // 0 = pass 2 uses `tol`. Set a STRICTER value when the
	// character's own light fills sit close to the bg (PP's ivory belly is
	// ~12/channel off the near-white bg): pockets are painted in the EXACT
	// bg color, so a tight tolerance kills them without touching the belly.
}

// Policy (2026-08-01, revised same day): bake ONLY sheets with a REPORTED,
// verified background problem — never as a proactive sweep. The first sweep
// over the older connected-key action sheets was rolled back by the user:
// even the 400px pocket cap punched pixels out of PP's ivory (give-flower
// "missing a lot of pixels"), and those sheets were already rendering fine
// through their runtime loaders. When a sheet DOES need baking: choose flood
// vs global by its runtime loader, use pocketTol 6 on near-white bgs, and
// ALWAYS eyeball the result (Read the PNG) before calling it done.
var sheets = []sheetCfg{
	{"assets/images/player/pp_sleeping.png", 24, 4000, 120, 2, false, 0},
	{"assets/images/player/pp_waking.png", 24, 4000, 120, 2, false, 0},
	// PP talk side: FLOOD tol 20 clears the streaky bg; the enclosed
	// arm/body pockets are painted in the EXACT bg color, so the strict
	// pocketTol 6 kills them while PP's ivory belly (~12/channel off the
	// bg) survives. A global tol-20 bake speckled the belly — rolled back
	// (user 2026-08-01).
	{"assets/images/player/PP talk side.png", 20, 400, 120, 0, false, 6},
	// Pierre: enclosedMax 30000 — frame 6 traps a BIG bg pocket between his
	// apron and the easel (well over the default 4000). Nothing legit on him
	// is near the pale sky blue, so the wide cap is safe.
	{"assets/images/locations/paris/npc/outside/npc_pierre_get_baguette.png", 24, 30000, 120, 2, false, 0},
	{"assets/images/locations/jerusalem/npc/wall/npc_bagel_seller_give.png", 24, 4000, 120, 2, false, 0},
}

func main() {
	dry := flag.Bool("dry", false, "report what would change, write nothing")
	flag.Parse()
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	for _, cfg := range sheets {
		full := filepath.Join(root, cfg.path)
		if err := bake(full, cfg, *dry); err != nil {
			fmt.Printf("FAIL %-70s %v\n", cfg.path, err)
		}
	}
}

func bake(path string, cfg sheetCfg, dry bool) error {
	img, err := loadNRGBA(path)
	if err != nil {
		return err
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	bgColors := sampleBG(img)
	if len(bgColors) == 0 {
		fmt.Printf("SKIP %-70s already transparent (or no opaque bg sample)\n", cfg.path)
		return nil
	}
	// AUTO fringe: bluish backgrounds get the defringe pass; near-white
	// backgrounds keep it off (it would eat thin pale strokes).
	if cfg.fringeIter < 0 {
		bg := bgColors[0]
		if int(bg.B) > int(bg.R)+15 {
			cfg.fringeIter = 2
		} else {
			cfg.fringeIter = 0
		}
	}

	matchesBG := func(c color.NRGBA) bool {
		if c.A < 200 {
			return false
		}
		for _, bg := range bgColors {
			if absd(c.R, bg.R) <= cfg.tol && absd(c.G, bg.G) <= cfg.tol && absd(c.B, bg.B) <= cfg.tol {
				return true
			}
		}
		return false
	}

	// --- pass 1: edge flood over bg-colored pixels (or GLOBAL key) ---
	mask := make([]bool, w*h) // true = make transparent
	seen := make([]bool, w*h)
	idx := func(x, y int) int { return (y-b.Min.Y)*w + (x - b.Min.X) }
	floodN := 0
	if cfg.global {
		// runtime-global-key equivalent: every matching pixel goes,
		// connectivity ignored. Passes 2-4 still run on what's left.
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				if matchesBG(img.NRGBAAt(x, y)) {
					mask[idx(x, y)] = true
					floodN++
				}
			}
		}
	} else {
		var stack []image.Point
		push := func(x, y int) {
			if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
				return
			}
			i := idx(x, y)
			if seen[i] || !matchesBG(img.NRGBAAt(x, y)) {
				return
			}
			seen[i] = true
			stack = append(stack, image.Point{X: x, Y: y})
		}
		for x := b.Min.X; x < b.Max.X; x++ {
			push(x, b.Min.Y)
			push(x, b.Max.Y-1)
		}
		for y := b.Min.Y; y < b.Max.Y; y++ {
			push(b.Min.X, y)
			push(b.Max.X-1, y)
		}
		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			mask[idx(p.X, p.Y)] = true
			floodN++
			push(p.X+1, p.Y)
			push(p.X-1, p.Y)
			push(p.X, p.Y+1)
			push(p.X, p.Y-1)
		}
	}

	// --- pass 2: small ENCLOSED bg-colored pockets ---
	pocketN, pocketPx := 0, 0
	if cfg.enclosedMax > 0 {
		ptol := cfg.pocketTol
		if ptol == 0 {
			ptol = cfg.tol
		}
		matchesPocket := func(c color.NRGBA) bool {
			if c.A < 200 {
				return false
			}
			for _, bg := range bgColors {
				if absd(c.R, bg.R) <= ptol && absd(c.G, bg.G) <= ptol && absd(c.B, bg.B) <= ptol {
					return true
				}
			}
			return false
		}
		comp := make([]bool, w*h)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				i := idx(x, y)
				if comp[i] || mask[i] || !matchesPocket(img.NRGBAAt(x, y)) {
					continue
				}
				pts := collect(img, b, comp, mask, x, y, matchesPocket)
				if len(pts) <= cfg.enclosedMax {
					for _, p := range pts {
						mask[idx(p.X, p.Y)] = true
					}
					pocketN++
					pocketPx += len(pts)
				}
			}
		}
	}

	// --- pass 3: tiny floating near-white / bg-ish specks ---
	speckN, speckPx := 0, 0
	if cfg.speckMax > 0 {
		comp := make([]bool, w*h)
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				i := idx(x, y)
				c := img.NRGBAAt(x, y)
				if comp[i] || mask[i] || c.A < 200 {
					continue
				}
				pts := collectOpaque(img, b, comp, mask, x, y)
				if len(pts) > cfg.speckMax {
					continue
				}
				// only kill it when it reads as halo junk: bright or bg-ish
				var rs, gs, bs int
				for _, p := range pts {
					pc := img.NRGBAAt(p.X, p.Y)
					rs += int(pc.R)
					gs += int(pc.G)
					bs += int(pc.B)
				}
				n := len(pts)
				mean := color.NRGBA{R: uint8(rs / n), G: uint8(gs / n), B: uint8(bs / n), A: 255}
				bright := (int(mean.R) + int(mean.G) + int(mean.B)) / 3
				if bright >= 228 || matchesBG(mean) {
					for _, p := range pts {
						mask[idx(p.X, p.Y)] = true
					}
					speckN++
					speckPx += n
				}
			}
		}
	}

	// --- pass 4: defringe — bg-tinted AA halo attached to the outline ---
	fringePx := 0
	if cfg.fringeIter > 0 {
		fringeTol := cfg.tol + 20
		matchesFringe := func(c color.NRGBA) bool {
			if c.A < 200 {
				return false
			}
			for _, bg := range bgColors {
				if absd(c.R, bg.R) <= fringeTol && absd(c.G, bg.G) <= fringeTol && absd(c.B, bg.B) <= fringeTol {
					return true
				}
			}
			return false
		}
		for it := 0; it < cfg.fringeIter; it++ {
			var add []int
			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					i := idx(x, y)
					c := img.NRGBAAt(x, y)
					if mask[i] || (c.A >= 200 && !matchesFringe(c)) {
						continue
					}
					// transparent-ish or fringe-colored: fringe-colored px must
					// touch existing transparency (mask or low alpha) to go.
					if c.A >= 200 {
						touches := false
						for _, d := range [4]image.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
							nx, ny := x+d.X, y+d.Y
							if nx < b.Min.X || nx >= b.Max.X || ny < b.Min.Y || ny >= b.Max.Y {
								touches = true
								break
							}
							if mask[idx(nx, ny)] || img.NRGBAAt(nx, ny).A < 200 {
								touches = true
								break
							}
						}
						if touches {
							add = append(add, i)
						}
					}
				}
			}
			for _, i := range add {
				mask[i] = true
			}
			fringePx += len(add)
		}
	}

	fmt.Printf("%-70s bg=%d flood=%dpx pockets=%d(%dpx) specks=%d(%dpx) fringe=%dpx\n",
		cfg.path, len(bgColors), floodN, pocketN, pocketPx, speckN, speckPx, fringePx)
	if dry {
		return nil
	}

	transparent := color.NRGBA{}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if mask[idx(x, y)] {
				img.SetNRGBA(x, y, transparent)
			}
		}
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, img)
}

// collect gathers a 4-connected component of bg-matching pixels.
func collect(img *image.NRGBA, b image.Rectangle, comp, mask []bool, sx, sy int, match func(color.NRGBA) bool) []image.Point {
	w := b.Dx()
	idx := func(x, y int) int { return (y-b.Min.Y)*w + (x - b.Min.X) }
	var pts []image.Point
	stack := []image.Point{{X: sx, Y: sy}}
	comp[idx(sx, sy)] = true
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		pts = append(pts, p)
		for _, d := range [4]image.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			x, y := p.X+d.X, p.Y+d.Y
			if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
				continue
			}
			i := idx(x, y)
			if comp[i] || mask[i] || !match(img.NRGBAAt(x, y)) {
				continue
			}
			comp[i] = true
			stack = append(stack, image.Point{X: x, Y: y})
		}
	}
	return pts
}

// collectOpaque gathers a 4-connected component of any opaque pixels.
func collectOpaque(img *image.NRGBA, b image.Rectangle, comp, mask []bool, sx, sy int) []image.Point {
	w := b.Dx()
	idx := func(x, y int) int { return (y-b.Min.Y)*w + (x - b.Min.X) }
	var pts []image.Point
	stack := []image.Point{{X: sx, Y: sy}}
	comp[idx(sx, sy)] = true
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		pts = append(pts, p)
		for _, d := range [4]image.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			x, y := p.X+d.X, p.Y+d.Y
			if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
				continue
			}
			i := idx(x, y)
			if comp[i] || mask[i] || img.NRGBAAt(x, y).A < 200 {
				continue
			}
			comp[i] = true
			stack = append(stack, image.Point{X: x, Y: y})
		}
	}
	return pts
}

func sampleBG(img *image.NRGBA) []color.NRGBA {
	b := img.Bounds()
	midX := (b.Min.X + b.Max.X) / 2
	midY := (b.Min.Y + b.Max.Y) / 2
	samples := []color.NRGBA{
		img.NRGBAAt(b.Min.X, b.Min.Y),
		img.NRGBAAt(b.Max.X-1, b.Min.Y),
		img.NRGBAAt(b.Min.X, b.Max.Y-1),
		img.NRGBAAt(b.Max.X-1, b.Max.Y-1),
		img.NRGBAAt(midX, b.Min.Y),
		img.NRGBAAt(b.Min.X, midY),
		img.NRGBAAt(midX, b.Max.Y-1),
		img.NRGBAAt(b.Max.X-1, midY),
	}
	var out []color.NRGBA
	for _, s := range samples {
		if s.A < 200 {
			continue
		}
		dup := false
		for _, o := range out {
			if absd(s.R, o.R) <= 5 && absd(s.G, o.G) <= 5 && absd(s.B, o.B) <= 5 {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, s)
		}
	}
	return out
}

func loadNRGBA(path string) (*image.NRGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	if n, ok := img.(*image.NRGBA); ok {
		return n, nil
	}
	n := image.NewNRGBA(img.Bounds())
	draw.Draw(n, img.Bounds(), img, img.Bounds().Min, draw.Src)
	return n, nil
}

func absd(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// sheet_repair fixes see-through holes inside character bodies on sheets the
// engine loads with a GLOBAL color key (player sheets): any pure-white region
// fully ENCLOSED by the body becomes transparent in-game. Such regions come
// from the v1 sheet_clean erasures (and from generators painting highlights as
// exact white). The repair: flood-fill the real background from the canvas
// borders, then repaint every remaining near-white enclosed region with the
// dominant color of its boundary — the hole melts into whatever surrounds it.
//
// Run: go run ./tools/sheet_repair
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

// Player sheets (global color-key loader) that had interior holes.
var paths = []string{
	"assets/images/player/PP talk front.png",
	"assets/images/player/PP idle side.png",
	"assets/images/player/PP celebrate.png",
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// keyWhite matches the engine's global color key (tolerance 8 per channel).
func keyWhite(p color.RGBA) bool {
	return p.A >= 0x40 &&
		absInt(int(p.R)-255) <= 8 && absInt(int(p.G)-255) <= 8 && absInt(int(p.B)-255) <= 8
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	_ = os.Chdir(root)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			fmt.Printf("SKIP (missing): %s\n", path)
			continue
		}
		src, err := png.Decode(f)
		f.Close()
		if err != nil {
			fmt.Printf("SKIP (bad png): %s\n", path)
			continue
		}
		b := src.Bounds()
		img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
		W, H := img.Bounds().Dx(), img.Bounds().Dy()

		// 1) flood-fill the OUTSIDE background from all border pixels.
		outside := make([]bool, W*H)
		var stack []int
		push := func(x, y int) {
			i := y*W + x
			if !outside[i] && keyWhite(img.RGBAAt(x, y)) {
				outside[i] = true
				stack = append(stack, i)
			}
		}
		for x := 0; x < W; x++ {
			push(x, 0)
			push(x, H-1)
		}
		for y := 0; y < H; y++ {
			push(0, y)
			push(W-1, y)
		}
		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x, y := p%W, p/W
			if x > 0 {
				push(x-1, y)
			}
			if x < W-1 {
				push(x+1, y)
			}
			if y > 0 {
				push(x, y-1)
			}
			if y < H-1 {
				push(x, y+1)
			}
		}

		// 2) every key-white pixel NOT reached is an enclosed hole. Group into
		// components, repaint each with the dominant boundary color.
		visited := make([]bool, W*H)
		repaired, regions := 0, 0
		for i := 0; i < W*H; i++ {
			x0, y0 := i%W, i/W
			if visited[i] || outside[i] || !keyWhite(img.RGBAAt(x0, y0)) {
				continue
			}
			// collect the component + boundary color histogram
			var comp []int
			colorCount := map[color.RGBA]int{}
			stack = append(stack[:0], i)
			visited[i] = true
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				comp = append(comp, p)
				x, y := p%W, p/W
				check := func(qx, qy int) {
					q := qy*W + qx
					pc := img.RGBAAt(qx, qy)
					if keyWhite(pc) {
						if !visited[q] && !outside[q] {
							visited[q] = true
							stack = append(stack, q)
						}
						return
					}
					if pc.A >= 0x40 {
						colorCount[pc]++
					}
				}
				if x > 0 {
					check(x-1, y)
				}
				if x < W-1 {
					check(x+1, y)
				}
				if y > 0 {
					check(x, y-1)
				}
				if y < H-1 {
					check(x, y+1)
				}
			}
			// dominant boundary color (skip near-black ink outlines so the
			// fill matches the fur/belly, not the linework)
			var best color.RGBA
			bestN := -1
			for cc, n := range colorCount {
				if int(cc.R)+int(cc.G)+int(cc.B) < 150 {
					continue // ink outline
				}
				if n > bestN {
					best, bestN = cc, n
				}
			}
			if bestN < 0 { // only ink around — use it anyway
				for cc, n := range colorCount {
					if n > bestN {
						best, bestN = cc, n
					}
				}
			}
			if bestN < 0 {
				continue
			}
			for _, p := range comp {
				img.SetRGBA(p%W, p/W, best)
			}
			repaired += len(comp)
			regions++
		}

		if repaired == 0 {
			fmt.Printf("NO HOLES: %s\n", path)
			continue
		}
		out, err := os.Create(path)
		if err != nil {
			fmt.Printf("ERROR writing %s: %v\n", path, err)
			continue
		}
		if err := png.Encode(out, img); err != nil {
			fmt.Printf("ERROR encoding %s: %v\n", path, err)
		}
		out.Close()
		fmt.Printf("REPAIRED %-55s %d enclosed white regions, %d px refilled\n", path, regions, repaired)
	}
}

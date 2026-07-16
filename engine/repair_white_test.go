package engine

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestRepairEnclosedWhite normalizes generated white pixels:
//   - edge-connected white/near-white is made pure white so the chroma key removes
//     the background fringe cleanly;
//   - enclosed white/near-white on the sprite (eyes, teeth, props) is recolored to
//     a soft vanilla so it does not key out as a hole.
//
// Interim fix until the proper no-white re-roll lands. Gated + reversible (git).
//
//	REPAIR_WHITE=1 GOTMPDIR=.gotmp go test ./engine -run TestRepairEnclosedWhite -v
func TestRepairEnclosedWhite(t *testing.T) {
	if os.Getenv("REPAIR_WHITE") == "" {
		t.Skip("set REPAIR_WHITE=1 to recolor enclosed pure-white to off-white")
	}
	vanilla := color.NRGBA{R: 242, G: 239, B: 229, A: 255}
	sheets := []struct {
		path        string
		replacement color.NRGBA
		clearGaps   bool
	}{
		{"../assets/images/player/PP walk front.png", vanilla, false},
		{"../assets/images/player/PP walk back.png", vanilla, false},
		{"../assets/images/player/PP grab flower.png", vanilla, true},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_idle.png", vanilla, false},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_talk.png", vanilla, false},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_get_baguette.png", vanilla, false},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_get_jam.png", vanilla, false},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_give_pass.png", vanilla, false},
	}
	for _, sh := range sheets {
		img, err := loadPNG(sh.path)
		if err != nil {
			t.Errorf("load %s: %v", sh.path, err)
			continue
		}
		b := img.Bounds()
		w, h := b.Dx(), b.Dy()
		isWhiteish := func(x, y int) bool {
			c := img.NRGBAAt(b.Min.X+x, b.Min.Y+y)
			return c.A >= 200 && c.R >= 235 && c.G >= 235 && c.B >= 235
		}
		// Flood white/near-white from the sheet borders = the background fringe
		// (+ the gaps between frames, all one connected region).
		bg := make([]bool, w*h)
		stack := make([][2]int, 0, w*h/4)
		push := func(x, y int) {
			if x >= 0 && x < w && y >= 0 && y < h && !bg[y*w+x] && isWhiteish(x, y) {
				bg[y*w+x] = true
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
		normalizedBG, recolored := 0, 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if !isWhiteish(x, y) {
					continue
				}
				if bg[y*w+x] {
					img.SetNRGBA(b.Min.X+x, b.Min.Y+y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
					normalizedBG++
				} else {
					img.SetNRGBA(b.Min.X+x, b.Min.Y+y, sh.replacement)
					recolored++
				}
			}
		}
		clearedGaps := 0
		if sh.clearGaps {
			clearedGaps = clearPPGrabFlowerBodyGaps(img)
		}
		f, err := os.Create(sh.path)
		if err != nil {
			t.Errorf("write %s: %v", sh.path, err)
			continue
		}
		_ = png.Encode(f, img)
		f.Close()
		fmt.Printf("%s: normalized %d bg/fringe px, recolored %d enclosed sprite px, cleared %d body-gap px\n", sh.path, normalizedBG, recolored, clearedGaps)
	}
}

func clearPPGrabFlowerBodyGaps(img interface {
	Bounds() image.Rectangle
	NRGBAAt(x, y int) color.NRGBA
	SetNRGBA(x, y int, c color.NRGBA)
}) int {
	b := img.Bounds()
	isVanillaGap := func(c color.NRGBA) bool {
		return c.A < 200 || (c.R >= 235 && c.G >= 230 && c.B >= 220)
	}
	clearRect := image.Rect(1390, 430, 1438, 585).Intersect(b)
	cleared := 0
	for y := clearRect.Min.Y; y < clearRect.Max.Y; y++ {
		for x := clearRect.Min.X; x < clearRect.Max.X; x++ {
			if isVanillaGap(img.NRGBAAt(x, y)) {
				img.SetNRGBA(x, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
				cleared++
			}
		}
	}
	return cleared
}

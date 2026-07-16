package engine

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestSpriteCleanExample is the PROTOTYPE white-spot cleaner, run on a few
// sheets before we generalise it. Each sheet becomes a green-square montage of
// FOUR rows so the result is eyeballable:
//
//	row 0 = original
//	row 1 = pure-white pixels painted MAGENTA (locate the spots)
//	row 2 = white REMOVED (→ transparent)
//	row 3 = white RECOLORED to cream
//
// Gated behind SPRITE_CLEAN=1; writes PNGs to the temp gallery dir (no repo
// clutter).
//
//	SPRITE_CLEAN=1 GOTMPDIR=.gotmp go test ./engine -run TestSpriteCleanExample -v
func TestSpriteCleanExample(t *testing.T) {
	if os.Getenv("SPRITE_CLEAN") == "" {
		t.Skip("set SPRITE_CLEAN=1 to render the white-spot cleaner prototype")
	}
	sheets := []struct {
		name       string
		path       string
		cols, rows int
		tol        uint8
	}{
		{"flower", "../assets/images/player/PP grab flower.png", 6, 1, 36},
		{"pierre", "../assets/images/locations/paris/npc/outside/npc_pierre_idle.png", 8, 1, 8},
	}
	outDir := filepath.Join(os.TempDir(), "clonedpp_sprite_gallery")
	_ = os.MkdirAll(outDir, 0o755)

	const thumbW, thumbH, pad = 300, 560, 12 // large squares → frames render near native res
	cream := color.NRGBA{R: 236, G: 229, B: 206, A: 255}

	for _, sh := range sheets {
		img, err := loadPNG(sh.path)
		if err != nil {
			t.Logf("skip %s: %v", sh.name, err)
			continue
		}
		applyColorKeyConnectedTol(img, sh.tol)
		eraseGridLines(img, sh.cols, sh.rows)
		bounds := img.Bounds()
		rects := contentGridRects(img, sh.cols, sh.rows)

		highlighted := image.NewNRGBA(bounds)
		transparent := image.NewNRGBA(bounds)
		recolored := image.NewNRGBA(bounds)
		draw.Draw(highlighted, bounds, img, bounds.Min, draw.Src)
		draw.Draw(transparent, bounds, img, bounds.Min, draw.Src)
		draw.Draw(recolored, bounds, img, bounds.Min, draw.Src)
		whitePixels := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				p := img.NRGBAAt(x, y)
				if p.A >= 200 && p.R >= 244 && p.G >= 244 && p.B >= 244 {
					whitePixels++
					highlighted.SetNRGBA(x, y, color.NRGBA{R: 255, G: 0, B: 220, A: 255})
					transparent.SetNRGBA(x, y, color.NRGBA{})
					recolored.SetNRGBA(x, y, cream)
				}
			}
		}

		type cellBox struct {
			box  image.Rectangle
			leak int
		}
		var cells []cellBox
		for r := 0; r < sh.rows; r++ {
			for c := 0; c < sh.cols; c++ {
				cell := gridCellRect(bounds, sh.cols, sh.rows, c, r, 3)
				if rects != nil {
					cell = rects[r][c]
				}
				ox, oy, ow, oh := opaqueLocal(img, cell)
				cb := cellBox{}
				if ow > 0 && oh > 0 {
					cb.box = image.Rect(cell.Min.X+int(ox), cell.Min.Y+int(oy),
						cell.Min.X+int(ox)+int(ow), cell.Min.Y+int(oy)+int(oh))
					if ox <= 2 && int32(ox+ow) >= int32(cell.Dx())-2 {
						cb.leak = 1
					}
				}
				cells = append(cells, cb)
			}
		}

		rowSrc := []*image.NRGBA{img, highlighted, transparent, recolored}
		n := len(cells)
		mw := n*(thumbW+pad) + pad
		mh := len(rowSrc)*(thumbH+pad) + pad
		mont := image.NewRGBA(image.Rect(0, 0, mw, mh))
		draw.Draw(mont, mont.Bounds(), &image.Uniform{color.RGBA{32, 32, 38, 255}}, image.Point{}, draw.Src)
		green := color.RGBA{60, 200, 90, 255}
		red := color.RGBA{230, 60, 60, 255}
		for i, cb := range cells {
			x0 := pad + i*(thumbW+pad)
			for row, src := range rowSrc {
				y0 := pad + row*(thumbH+pad)
				sq := image.Rect(x0, y0, x0+thumbW, y0+thumbH)
				bord := green
				if row == 0 && cb.leak > 0 {
					bord = red
				}
				drawBorder(mont, sq, bord, 3)
				if cb.box.Dx() > 0 {
					blitFit(mont, sq.Inset(4), src.SubImage(cb.box))
				}
			}
		}
		out := filepath.Join(outDir, sh.name+"_clean.png")
		f, err := os.Create(out)
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		_ = png.Encode(f, mont)
		f.Close()
		fmt.Printf("%-8s white=%6d px  rows: original / magenta / removed / cream  -> %s\n", sh.name, whitePixels, out)
	}
}

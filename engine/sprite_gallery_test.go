package engine

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestSpriteGallery renders every registered sheet (realSheetCases) the way the
// ENGINE cuts it — gap-detect into content cells, falling back to equal cells —
// into a labeled montage PNG per sheet, plus an index.html. Open the index with
// a display to eyeball every sprite's frames at once; frames that cut badly
// (sliver / double-wide / fallback) are bordered RED so problems jump out (this
// is exactly what made office-Higgins's talk "blink"). It's a generator, not a
// pass/fail test, so it's gated behind SPRITE_GALLERY=1 and skipped in CI:
//
//	SPRITE_GALLERY=1 GOTMPDIR=.gotmp go test ./engine -run TestSpriteGallery
//
// Output dir: $SPRITE_GALLERY_DIR (default ./sprite_gallery).
func TestSpriteGallery(t *testing.T) {
	if os.Getenv("SPRITE_GALLERY") == "" {
		t.Skip("set SPRITE_GALLERY=1 to generate the visual sprite gallery")
	}
	// Default OUT OF the repo (OS temp) so we never litter the project with
	// generated montages. Override with SPRITE_GALLERY_DIR.
	outDir := os.Getenv("SPRITE_GALLERY_DIR")
	if outDir == "" {
		outDir = filepath.Join(os.TempDir(), "clonedpp_sprite_gallery")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir gallery: %v", err)
	}

	const thumbW, thumbH, pad = 130, 170, 6

	type entry struct {
		title   string
		png     string
		mode    string
		frames  int
		anomaly string
	}
	var entries []entry

	for _, c := range realSheetCases {
		img, err := loadPNG(c.path)
		if err != nil {
			continue
		}
		// Match the game's per-sheet loader for the special cases so the gallery
		// shows what actually renders. Office Higgins uses EQUAL cells (his poses
		// touch and mis-gap-detect) with tol 4 + a wide inset.
		equal, inset, tol := galleryLoaderFor(c.path)
		applyColorKeyConnectedTol(img, tol)
		eraseGridLines(img, c.cols, c.rows)
		bounds := img.Bounds()
		var rects [][]image.Rectangle
		if !equal {
			rects = contentGridRects(img, c.cols, c.rows)
		}
		mode := "gap-detect"
		if equal {
			mode = "equal (forced)"
		} else if rects == nil {
			mode = "equal (fallback)"
		}

		// Gather the opaque-box crop of every cell, in row-major order.
		type frame struct {
			crop   image.Image
			ow, oh int
		}
		var frames []frame
		var widths []int
		for r := 0; r < c.rows; r++ {
			for cc := 0; cc < c.cols; cc++ {
				cell := gridCellRect(bounds, c.cols, c.rows, cc, r, inset)
				if rects != nil {
					cell = rects[r][cc]
				}
				ox, oy, ow, oh := opaqueLocal(img, cell)
				if ow <= 0 || oh <= 0 {
					frames = append(frames, frame{nil, 0, 0})
					continue
				}
				crop := image.Rect(cell.Min.X+int(ox), cell.Min.Y+int(oy),
					cell.Min.X+int(ox)+int(ow), cell.Min.Y+int(oy)+int(oh))
				frames = append(frames, frame{img.SubImage(crop), int(ow), int(oh)})
				widths = append(widths, int(ow))
			}
		}

		// Median opaque width → flag slivers (<55%) and doubles (>160%).
		med := medianInt(widths)
		anomalies := 0
		bad := make([]bool, len(frames))
		for i, f := range frames {
			if f.crop == nil {
				bad[i] = true
				anomalies++
				continue
			}
			if med > 0 && (f.ow < med*55/100 || f.ow > med*160/100) {
				bad[i] = true
				anomalies++
			}
		}

		// Montage: one row, each frame in a thumbW×thumbH box, content scaled to
		// fit, border green (ok) or red (bad).
		n := len(frames)
		mw := n*(thumbW+pad) + pad
		mh := thumbH + 2*pad
		mont := image.NewRGBA(image.Rect(0, 0, mw, mh))
		draw.Draw(mont, mont.Bounds(), &image.Uniform{color.RGBA{32, 32, 38, 255}}, image.Point{}, draw.Src)
		for i, f := range frames {
			x0 := pad + i*(thumbW+pad)
			box := image.Rect(x0, pad, x0+thumbW, pad+thumbH)
			border := color.RGBA{60, 200, 90, 255}
			if bad[i] {
				border = color.RGBA{230, 60, 60, 255}
			}
			drawBorder(mont, box, border, 3)
			if f.crop != nil {
				blitFit(mont, box.Inset(4), f.crop)
			}
		}

		base := sanitizeName(c.path)
		pngName := base + ".png"
		fout, err := os.Create(filepath.Join(outDir, pngName))
		if err != nil {
			t.Fatalf("create montage: %v", err)
		}
		_ = png.Encode(fout, mont)
		fout.Close()

		an := ""
		if anomalies > 0 {
			an = fmt.Sprintf("%d frame(s) sliver/double/blank — median w=%d", anomalies, med)
		}
		entries = append(entries, entry{strings.TrimPrefix(c.path, "../"), pngName, mode, n, an})
	}

	// Sort: anomalies first, so the broken sheets are at the top of the page.
	sort.SliceStable(entries, func(i, j int) bool {
		ai, aj := entries[i].anomaly != "", entries[j].anomaly != ""
		if ai != aj {
			return ai
		}
		return entries[i].title < entries[j].title
	})

	var b strings.Builder
	b.WriteString("<!doctype html><meta charset=utf-8><title>Sprite gallery</title>")
	b.WriteString("<style>body{background:#16161c;color:#ddd;font:13px system-ui;margin:20px}")
	b.WriteString("h2{font-size:14px;margin:18px 0 4px} .warn{color:#ff6b6b;font-weight:bold} .ok{color:#7fd07f}")
	b.WriteString("img{display:block;background:#202028;border-radius:4px;max-width:100%}</style>")
	b.WriteString(fmt.Sprintf("<h1>Sprite gallery — %d sheets</h1>", len(entries)))
	b.WriteString("<p>Red border = sliver/double/blank cut (engine gap-detect). Green = clean.</p>")
	for _, e := range entries {
		cls, note := "ok", "clean"
		if e.anomaly != "" {
			cls, note = "warn", "⚠ "+e.anomaly
		}
		b.WriteString(fmt.Sprintf("<h2>%s &nbsp;<small>[%s, %d frames] <span class=%s>%s</span></small></h2>",
			e.title, e.mode, e.frames, cls, note))
		b.WriteString(fmt.Sprintf("<img src=\"%s\">", e.png))
	}
	if err := os.WriteFile(filepath.Join(outDir, "index.html"), []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	abs, _ := filepath.Abs(filepath.Join(outDir, "index.html"))
	t.Logf("sprite gallery: %d sheets -> %s", len(entries), abs)
}

// galleryLoaderFor returns how the game cuts a given sheet, so the gallery is
// faithful for sheets that use a non-default loader. Keep in sync with the
// factories in game/npc.go. Default: gap-detect, inset 3, tol 8.
func galleryLoaderFor(path string) (equal bool, inset int, tol uint8) {
	switch {
	case strings.Contains(path, "npc_director_higgins_office_idle"),
		strings.Contains(path, "npc_director_higgins_office_talk"):
		return true, 4, 24 // game/npc.go newOfficeHiggins: equal cells, higginsOfficeInset 4, higginsOfficeTol 24 (high tol breaks the touching-pose bridge)
	case strings.Contains(path, "/PP give "), strings.Contains(path, "/PP put note"):
		return false, 3, 24 // player.go: give/note one-shots load gap-detect at tol 24 (#5)
	}
	return false, 3, 8
}

func medianInt(v []int) int {
	if len(v) == 0 {
		return 0
	}
	s := append([]int(nil), v...)
	sort.Ints(s)
	return s[len(s)/2]
}

func drawBorder(dst *image.RGBA, r image.Rectangle, c color.Color, w int) {
	for i := 0; i < w; i++ {
		rr := r.Inset(i)
		for x := rr.Min.X; x < rr.Max.X; x++ {
			dst.Set(x, rr.Min.Y, c)
			dst.Set(x, rr.Max.Y-1, c)
		}
		for y := rr.Min.Y; y < rr.Max.Y; y++ {
			dst.Set(rr.Min.X, y, c)
			dst.Set(rr.Max.X-1, y, c)
		}
	}
}

// blitFit nearest-neighbour scales src to fit inside dst (preserving aspect),
// centred. Good enough for a contact-sheet thumbnail.
func blitFit(dst *image.RGBA, dst2 image.Rectangle, src image.Image) {
	sb := src.Bounds()
	sw, sh := sb.Dx(), sb.Dy()
	if sw == 0 || sh == 0 {
		return
	}
	scale := float64(dst2.Dx()) / float64(sw)
	if s := float64(dst2.Dy()) / float64(sh); s < scale {
		scale = s
	}
	// Never upscale — keep crisp native pixels when the square is large enough
	// (the high-res clean montage). Small gallery squares still downscale.
	if scale > 1 {
		scale = 1
	}
	dw, dh := int(float64(sw)*scale), int(float64(sh)*scale)
	offX := dst2.Min.X + (dst2.Dx()-dw)/2
	offY := dst2.Min.Y + (dst2.Dy()-dh)/2
	for y := 0; y < dh; y++ {
		sy := sb.Min.Y + int(float64(y)/scale)
		for x := 0; x < dw; x++ {
			sx := sb.Min.X + int(float64(x)/scale)
			dst.Set(offX+x, offY+y, src.At(sx, sy))
		}
	}
}

func sanitizeName(p string) string {
	p = strings.TrimPrefix(p, "../assets/images/")
	p = strings.NewReplacer("/", "_", " ", "_", ".png", "").Replace(p)
	return p
}

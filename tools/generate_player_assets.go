package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
)

type frameSpec struct {
	path string
	cols int
	rows int
	col  int
	row  int
}

type frameImage struct {
	img image.Image
	w   int
	h   int
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	const playerStripSheetPath = "assets/images/player/strip_4x5.png"

	files := map[string]*image.NRGBA{}
	load := func(rel string) *image.NRGBA {
		if img, ok := files[rel]; ok {
			return img
		}
		full := filepath.Join(root, rel)
		img := loadPNG(full)
		files[rel] = img
		return img
	}

	walkAwaySpecs := rowSpecs(playerStripSheetPath, 8, 4, 3)
	buildStrip(root, "assets/images/player/pink_panther_walk_away.png", walkAwaySpecs, load)

	walkFrontSpecs := rowSpecs(playerStripSheetPath, 8, 4, 2)
	buildStrip(root, "assets/images/player/pink_panther_walk_front.png", walkFrontSpecs, load)

	walkSideSpecs := rowSpecs(playerStripSheetPath, 8, 4, 0)
	buildStrip(root, "assets/images/player/pink_panther_walk_side.png", walkSideSpecs, load)

	talkFrontSpecs := pingPongRowSpecs(playerStripSheetPath, 8, 4, 1)
	buildStrip(root, "assets/images/player/pink_panther_talk_front.png", talkFrontSpecs, load)

	talkSideSpecs := pingPongRowSpecs(playerStripSheetPath, 8, 4, 1)
	buildStrip(root, "assets/images/player/pink_panther_talk_side.png", talkSideSpecs, load)

	idleSpecs := []frameSpec{
		{path: playerStripSheetPath, cols: 8, rows: 4, col: 3, row: 2},
		{path: playerStripSheetPath, cols: 8, rows: 4, col: 6, row: 0},
		{path: playerStripSheetPath, cols: 8, rows: 4, col: 3, row: 3},
	}
	buildStrip(root, "assets/images/player/pink_panther_idle.png", idleSpecs, load)

	fullStripSpecs := make([]frameSpec, 0, len(idleSpecs)+len(walkSideSpecs)+len(walkAwaySpecs)+len(walkFrontSpecs)+len(talkFrontSpecs)+len(talkSideSpecs))
	fullStripSpecs = append(fullStripSpecs, idleSpecs...)
	fullStripSpecs = append(fullStripSpecs, walkSideSpecs...)
	fullStripSpecs = append(fullStripSpecs, walkAwaySpecs...)
	fullStripSpecs = append(fullStripSpecs, walkFrontSpecs...)
	fullStripSpecs = append(fullStripSpecs, talkFrontSpecs...)
	fullStripSpecs = append(fullStripSpecs, talkSideSpecs...)
	buildStrip(root, "assets/images/player/pink_panther_full_strip.png", fullStripSpecs, load)
}

func rowSpecs(path string, cols, rows, row int) []frameSpec {
	specs := make([]frameSpec, 0, cols)
	for col := 0; col < cols; col++ {
		specs = append(specs, frameSpec{path: path, cols: cols, rows: rows, col: col, row: row})
	}
	return specs
}

func pingPongRowSpecs(path string, cols, rows, row int) []frameSpec {
	specs := make([]frameSpec, 0, cols*2-2)
	for col := 0; col < cols; col++ {
		specs = append(specs, frameSpec{path: path, cols: cols, rows: rows, col: col, row: row})
	}
	for col := cols - 2; col >= 1; col-- {
		specs = append(specs, frameSpec{path: path, cols: cols, rows: rows, col: col, row: row})
	}
	return specs
}

func loadPNG(path string) *image.NRGBA {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	return toNRGBA(img)
}

func toNRGBA(img image.Image) *image.NRGBA {
	if nrgba, ok := img.(*image.NRGBA); ok {
		return nrgba
	}
	b := img.Bounds()
	dst := image.NewNRGBA(b)
	draw.Draw(dst, b, img, b.Min, draw.Src)
	return dst
}

func buildStrip(root string, outRel string, specs []frameSpec, load func(string) *image.NRGBA) {
	frames := make([]frameImage, 0, len(specs))
	maxW := 0
	maxH := 0

	for _, spec := range specs {
		cell := extractCell(load(spec.path), spec.cols, spec.rows, spec.col, spec.row)
		probe := cloneNRGBA(cell)
		applyColorKey(probe)
		trimmed := cropToRect(probe, opaqueBounds(probe))
		frames = append(frames, frameImage{img: trimmed, w: trimmed.Bounds().Dx(), h: trimmed.Bounds().Dy()})
		if trimmed.Bounds().Dx() > maxW {
			maxW = trimmed.Bounds().Dx()
		}
		if trimmed.Bounds().Dy() > maxH {
			maxH = trimmed.Bounds().Dy()
		}
	}

	const padX = 14
	const padTop = 14
	const padBottom = 10

	cellW := maxW + padX*2
	cellH := maxH + padTop + padBottom
	dst := image.NewNRGBA(image.Rect(0, 0, cellW*len(frames), cellH))

	for i, frame := range frames {
		x := i*cellW + (cellW-frame.w)/2
		y := cellH - padBottom - frame.h
		draw.Draw(dst, image.Rect(x, y, x+frame.w, y+frame.h), frame.img, frame.img.Bounds().Min, draw.Over)
	}

	outPath := filepath.Join(root, outRel)
	outFile, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	if err := png.Encode(outFile, dst); err != nil {
		panic(err)
	}
	fmt.Println("wrote", outRel)
}

func extractCell(img *image.NRGBA, cols, rows, col, row int) *image.NRGBA {
	b := img.Bounds()
	cellW := b.Dx() / cols
	cellH := b.Dy() / rows

	x0 := b.Min.X + col*cellW
	y0 := b.Min.Y + row*cellH
	x1 := x0 + cellW
	y1 := y0 + cellH
	if col == cols-1 {
		x1 = b.Max.X
	}
	if row == rows-1 {
		y1 = b.Max.Y
	}
	if x1-x0 > 6 {
		x0 += 3
		x1 -= 3
	}
	if y1-y0 > 6 {
		y0 += 3
		y1 -= 3
	}

	rect := image.Rect(x0, y0, x1, y1)
	dst := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
	return dst
}

func applyColorKey(img *image.NRGBA) {
	b := img.Bounds()
	samples := []color.NRGBA{
		img.NRGBAAt(b.Min.X, b.Min.Y),
		img.NRGBAAt(b.Max.X-1, b.Min.Y),
		img.NRGBAAt(b.Min.X, b.Max.Y-1),
		img.NRGBAAt(b.Max.X-1, b.Max.Y-1),
	}

	backgrounds := make([]color.NRGBA, 0, len(samples))
	for _, sample := range samples {
		if sample.A < 200 {
			continue
		}
		duplicate := false
		for _, bg := range backgrounds {
			if similar(sample, bg, 8) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			backgrounds = append(backgrounds, sample)
		}
	}

	if len(backgrounds) == 0 {
		return
	}

	transparent := color.NRGBA{0, 0, 0, 0}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			if c.A < 200 {
				continue
			}
			for _, bg := range backgrounds {
				if similar(c, bg, 30) {
					img.SetNRGBA(x, y, transparent)
					break
				}
			}
		}
	}
}

func similar(a, b color.NRGBA, tol uint8) bool {
	return absDiff(a.R, b.R) <= tol && absDiff(a.G, b.G) <= tol && absDiff(a.B, b.B) <= tol
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

func cloneNRGBA(img *image.NRGBA) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)
	return dst
}

func opaqueBounds(img *image.NRGBA) image.Rectangle {
	b := img.Bounds()
	minX, minY := b.Max.X, b.Max.Y
	maxX, maxY := b.Min.X, b.Min.Y

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if img.NRGBAAt(x, y).A > 10 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if minX > maxX || minY > maxY {
		return image.Rect(0, 0, 1, 1)
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func cropToRect(img *image.NRGBA, rect image.Rectangle) *image.NRGBA {
	dst := image.NewNRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), img, rect.Min, draw.Src)
	return dst
}

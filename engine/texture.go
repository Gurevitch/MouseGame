package engine

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// --- BMP loading ---

func TextureFromBMP(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading BMP %s: %v", filename, err))
	}
	defer img.Free()

	key := GetPixelColor(img, 0, 0)
	img.SetColorKey(true, key)

	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture from %s: %v", filename, err))
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	return tex
}

func TextureFromBMPRaw(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading BMP %s: %v", filename, err))
	}
	defer img.Free()

	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture from %s: %v", filename, err))
	}
	return tex
}

func GetPixelColor(s *sdl.Surface, x, y int32) uint32 {
	bpp := int(s.Format.BytesPerPixel)
	px := s.Pixels()
	off := int(y)*int(s.Pitch) + int(x)*bpp
	if off+bpp > len(px) {
		return 0
	}
	switch bpp {
	case 1:
		return uint32(px[off])
	case 2:
		return uint32(px[off]) | uint32(px[off+1])<<8
	case 3:
		return uint32(px[off]) | uint32(px[off+1])<<8 | uint32(px[off+2])<<16
	case 4:
		return uint32(px[off]) | uint32(px[off+1])<<8 | uint32(px[off+2])<<16 | uint32(px[off+3])<<24
	}
	return 0
}

// --- PNG loading with auto-crop ---

func loadPNG(filename string) (*image.NRGBA, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if nrgba, ok := img.(*image.NRGBA); ok {
		return nrgba, nil
	}

	bounds := img.Bounds()
	nrgba := image.NewNRGBA(bounds)
	draw.Draw(nrgba, bounds, img, bounds.Min, draw.Src)
	return nrgba, nil
}

// applyColorKey detects opaque background colors by sampling corners and edges,
// then makes all matching pixels fully transparent. Handles solid backgrounds,
// white backgrounds, and checkerboard patterns baked into the image.
func applyColorKey(img *image.NRGBA) {
	applyColorKeyTol(img, 8)
}

// applyColorKeyTol is the same as applyColorKey but with an adjustable
// per-channel tolerance. Use higher values (20-32) for sheets where the
// background bleeds into anti-aliased edges — e.g. the campfire loop
// where flame colors vary dramatically across frames and the corner
// samples sit near fringe pixels.
func applyColorKeyTol(img *image.NRGBA, matchTol uint8) {
	b := img.Bounds()
	midX := (b.Min.X + b.Max.X) / 2
	midY := (b.Min.Y + b.Max.Y) / 2

	samples := []color.NRGBA{
		img.NRGBAAt(b.Min.X, b.Min.Y),
		img.NRGBAAt(b.Max.X-1, b.Min.Y),
		img.NRGBAAt(b.Min.X, b.Max.Y-1),
		img.NRGBAAt(b.Max.X-1, b.Max.Y-1),
		img.NRGBAAt(b.Min.X+1, b.Min.Y),
		img.NRGBAAt(b.Min.X, b.Min.Y+1),
		img.NRGBAAt(midX, b.Min.Y),
		img.NRGBAAt(b.Min.X, midY),
	}

	var bgColors []color.NRGBA
	const dedupTol = 5
	for _, s := range samples {
		if s.A < 200 {
			continue
		}
		dup := false
		for _, bg := range bgColors {
			if absDiffU8(s.R, bg.R) <= dedupTol && absDiffU8(s.G, bg.G) <= dedupTol && absDiffU8(s.B, bg.B) <= dedupTol {
				dup = true
				break
			}
		}
		if !dup {
			bgColors = append(bgColors, s)
		}
	}

	if len(bgColors) == 0 {
		return
	}

	transparent := color.NRGBA{0, 0, 0, 0}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			if c.A < 200 {
				continue
			}
			for _, bg := range bgColors {
				if absDiffU8(c.R, bg.R) <= matchTol && absDiffU8(c.G, bg.G) <= matchTol && absDiffU8(c.B, bg.B) <= matchTol {
					img.SetNRGBA(x, y, transparent)
					break
				}
			}
		}
	}
}

func absDiffU8(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// findOpaqueBounds returns the tightest rectangle containing all pixels
// with alpha above a small threshold within the given region.
func findOpaqueBounds(img *image.NRGBA, region image.Rectangle) image.Rectangle {
	minX, minY := region.Max.X, region.Max.Y
	maxX, maxY := region.Min.X, region.Min.Y

	for y := region.Min.Y; y < region.Max.Y; y++ {
		for x := region.Min.X; x < region.Max.X; x++ {
			if img.NRGBAAt(x, y).A > 10 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if minX > maxX || minY > maxY {
		return region
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

// TextureFromNRGBA uploads the full NRGBA image to a new SDL texture.
// Handy for synthetic/placeholder imagery built in code.
func TextureFromNRGBA(renderer *sdl.Renderer, img *image.NRGBA) (*sdl.Texture, int32, int32) {
	return nrgbaToTexture(renderer, img, img.Bounds())
}

// nrgbaToTexture creates an SDL texture from a cropped region of an NRGBA image.
func nrgbaToTexture(renderer *sdl.Renderer, img *image.NRGBA, crop image.Rectangle) (*sdl.Texture, int32, int32) {
	w := int32(crop.Dx())
	h := int32(crop.Dy())

	surface, err := sdl.CreateRGBSurface(0, w, h, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		panic(fmt.Errorf("creating surface: %v", err))
	}
	defer surface.Free()

	pixels := surface.Pixels()
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			c := img.NRGBAAt(int(x)+crop.Min.X, int(y)+crop.Min.Y)
			off := int(y)*int(surface.Pitch) + int(x)*4
			pixels[off] = c.R
			pixels[off+1] = c.G
			pixels[off+2] = c.B
			pixels[off+3] = c.A
		}
	}

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(fmt.Errorf("creating texture: %v", err))
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	return tex, w, h
}

// TextureFromPNG loads a PNG file, auto-crops to the non-transparent bounding box,
// and returns the texture along with the cropped width and height.
func TextureFromPNG(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG %s: %v", filename, err))
	}
	applyColorKey(img)
	crop := findOpaqueBounds(img, img.Bounds())
	return nrgbaToTexture(renderer, img, crop)
}

// TextureFromPNGKeyed loads a PNG, applies color-key background removal, and
// returns the full image as a texture without auto-cropping.
func TextureFromPNGKeyed(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG %s: %v", filename, err))
	}
	applyColorKey(img)
	return nrgbaToTexture(renderer, img, img.Bounds())
}

// TextureFromPNGRaw loads a PNG without auto-cropping, returning the full image
// as a texture. Useful for backgrounds and full-scene art.
func TextureFromPNGRaw(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG %s: %v", filename, err))
	}
	return nrgbaToTexture(renderer, img, img.Bounds())
}

// SafeTextureFromPNGRaw loads a PNG like TextureFromPNGRaw but returns nil
// instead of panicking if the file is missing or unreadable.
func SafeTextureFromPNGRaw(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		fmt.Printf("Warning: could not load %s: %v\n", filename, err)
		return nil, 0, 0
	}
	return nrgbaToTexture(renderer, img, img.Bounds())
}

// SafeTextureFromPNGKeyed loads a PNG with color-key background removal,
// returning nil instead of panicking if the file is missing.
func SafeTextureFromPNGKeyed(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		fmt.Printf("Warning: could not load %s: %v\n", filename, err)
		return nil, 0, 0
	}
	applyColorKey(img)
	return nrgbaToTexture(renderer, img, img.Bounds())
}

// TextureFromPNGRawClean loads a PNG, removes the bottom-right watermark, and
// returns the full image as a texture without auto-cropping.
func TextureFromPNGRawClean(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG %s: %v", filename, err))
	}
	return nrgbaToTexture(renderer, img, img.Bounds())
}

// SurfaceFromPNG loads a PNG and returns it as an SDL surface (caller must Free).
// Useful for window icons and other non-renderer uses.
func SurfaceFromPNG(filename string) (*sdl.Surface, error) {
	img, err := loadPNG(filename)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	w := int32(b.Dx())
	h := int32(b.Dy())
	surface, err := sdl.CreateRGBSurface(0, w, h, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		return nil, err
	}
	pixels := surface.Pixels()
	for y := int32(0); y < h; y++ {
		for x := int32(0); x < w; x++ {
			c := img.NRGBAAt(int(x)+b.Min.X, int(y)+b.Min.Y)
			off := int(y)*int(surface.Pitch) + int(x)*4
			pixels[off] = c.R
			pixels[off+1] = c.G
			pixels[off+2] = c.B
			pixels[off+3] = c.A
		}
	}
	return surface, nil
}

// SpriteFramesFromPNG loads a PNG sprite sheet, splits it into numCols equal
// columns, auto-crops each column, and returns per-frame textures + dimensions.
func SpriteFramesFromPNG(renderer *sdl.Renderer, filename string, numCols int) ([]*sdl.Texture, []int32, []int32) {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG sheet %s: %v", filename, err))
	}
	applyColorKey(img)

	bounds := img.Bounds()
	colW := bounds.Dx() / numCols

	texs := make([]*sdl.Texture, numCols)
	ws := make([]int32, numCols)
	hs := make([]int32, numCols)

	for i := 0; i < numCols; i++ {
		colRect := image.Rect(
			bounds.Min.X+i*colW, bounds.Min.Y,
			bounds.Min.X+(i+1)*colW, bounds.Max.Y,
		)
		if i == numCols-1 {
			colRect.Max.X = bounds.Max.X
		}
		crop := findOpaqueBounds(img, colRect)
		texs[i], ws[i], hs[i] = nrgbaToTexture(renderer, img, crop)
	}

	return texs, ws, hs
}

// GridFrame holds a single frame extracted from a grid sprite sheet.
type GridFrame struct {
	Tex *sdl.Texture
	W   int32
	H   int32
}

// blankCornerLogo makes pixels in the bottom-right corner region transparent,
// removing watermarks/logos that image generators embed.
func blankCornerLogo(img *image.NRGBA, w, h int) {
	b := img.Bounds()
	transparent := color.NRGBA{0, 0, 0, 0}
	startX := b.Max.X - w
	startY := b.Max.Y - h
	if startX < b.Min.X {
		startX = b.Min.X
	}
	if startY < b.Min.Y {
		startY = b.Min.Y
	}
	for y := startY; y < b.Max.Y; y++ {
		for x := startX; x < b.Max.X; x++ {
			img.SetNRGBA(x, y, transparent)
		}
	}
}

// SpriteGridFromPNG loads a PNG sprite sheet arranged in a grid of cols x rows,
// removes the background via color-keying and any bottom-right watermark, and
// returns frames indexed [row][col]. Each cell uses its full grid dimensions
// (no auto-crop) so all frames share the same size.
// SpriteGridFromPNGRaw loads a PNG grid without color-key removal.
// Uses the PNG's native alpha channel. Each cell is its own texture.
func SpriteGridFromPNGRaw(renderer *sdl.Renderer, filename string, cols, rows int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG grid %s: %v", filename, err))
	}

	bounds := img.Bounds()
	cellW := bounds.Dx() / cols
	cellH := bounds.Dy() / rows

	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := image.Rect(
				bounds.Min.X+c*cellW, bounds.Min.Y+r*cellH,
				bounds.Min.X+(c+1)*cellW, bounds.Min.Y+(r+1)*cellH,
			)
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h}
		}
	}
	return grid
}

func SpriteGridFromPNG(renderer *sdl.Renderer, filename string, cols, rows int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG grid %s: %v", filename, err))
	}
	applyColorKey(img)

	bounds := img.Bounds()
	cellW := bounds.Dx() / cols
	cellH := bounds.Dy() / rows

	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := image.Rect(
				bounds.Min.X+c*cellW, bounds.Min.Y+r*cellH,
				bounds.Min.X+(c+1)*cellW, bounds.Min.Y+(r+1)*cellH,
			)
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h}
		}
	}

	return grid
}

// eraseGridLines scans horizontal and vertical strips of the image looking
// for grid-lines: thin rows/columns of near-uniform dark pixels that span
// most of the width/height. These dividers get baked into AI-generated
// sprite sheets. Any line detected is erased to fully transparent.
//
// History: earlier versions swept a ±6-pixel window and called anything
// RGB<80 "dark". That ate dark outlines on the sprite itself — boots,
// lashes, hairlines, the black rim of a coin. The tighter values below
// stay focused on the actual divider pixel without chewing into artwork:
//
//   - scanThickness ±2 instead of ±6: just enough slack for a 1-2 px
//     divider that's slightly off the mathematical cell boundary.
//   - isDark threshold RGB<50 instead of <80: divider black is near-zero;
//     sprite outlines are often 40-70. The old threshold spanned both.
//   - outer-edge wipe only touches pixels whose A >= 90. The old version
//     wiped every dark pixel along the outermost 3 rows/columns, which
//     meant any sprite that stepped right up to the cell edge lost its
//     foot or hat silhouette.
func eraseGridLines(img *image.NRGBA, cols, rows int) {
	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()
	cellW := w / cols
	cellH := h / rows

	isDark := func(c color.NRGBA) bool {
		if c.A < 40 {
			return false
		}
		return c.R < 50 && c.G < 50 && c.B < 50
	}

	scanThickness := 2
	transparent := color.NRGBA{0, 0, 0, 0}

	for c := 1; c < cols; c++ {
		centerX := b.Min.X + c*cellW
		for dx := -scanThickness; dx <= scanThickness; dx++ {
			x := centerX + dx
			if x < b.Min.X || x >= b.Max.X {
				continue
			}
			darkCount := 0
			for y := b.Min.Y; y < b.Max.Y; y++ {
				if isDark(img.NRGBAAt(x, y)) {
					darkCount++
				}
			}
			if float64(darkCount) >= float64(h)*0.70 {
				for y := b.Min.Y; y < b.Max.Y; y++ {
					if isDark(img.NRGBAAt(x, y)) {
						img.SetNRGBA(x, y, transparent)
					}
				}
			}
		}
	}

	for r := 1; r < rows; r++ {
		centerY := b.Min.Y + r*cellH
		for dy := -scanThickness; dy <= scanThickness; dy++ {
			y := centerY + dy
			if y < b.Min.Y || y >= b.Max.Y {
				continue
			}
			darkCount := 0
			for x := b.Min.X; x < b.Max.X; x++ {
				if isDark(img.NRGBAAt(x, y)) {
					darkCount++
				}
			}
			if float64(darkCount) >= float64(w)*0.70 {
				for x := b.Min.X; x < b.Max.X; x++ {
					if isDark(img.NRGBAAt(x, y)) {
						img.SetNRGBA(x, y, transparent)
					}
				}
			}
		}
	}

	// Outer-edge wipe: only touch pixels that are simultaneously dark
	// AND mostly opaque. A pixel with A < 90 is already fading into the
	// background and is almost certainly a fringe artifact, not a real
	// sprite pixel we'd want to erase.
	edgeDarkOpaque := func(c color.NRGBA) bool {
		if c.A < 90 {
			return false
		}
		return isDark(c)
	}
	edgeScan := 3
	for y := b.Min.Y; y < b.Min.Y+edgeScan && y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if edgeDarkOpaque(img.NRGBAAt(x, y)) {
				img.SetNRGBA(x, y, transparent)
			}
		}
	}
	for y := b.Max.Y - edgeScan; y < b.Max.Y; y++ {
		if y < b.Min.Y {
			continue
		}
		for x := b.Min.X; x < b.Max.X; x++ {
			if edgeDarkOpaque(img.NRGBAAt(x, y)) {
				img.SetNRGBA(x, y, transparent)
			}
		}
	}
	for x := b.Min.X; x < b.Min.X+edgeScan && x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if edgeDarkOpaque(img.NRGBAAt(x, y)) {
				img.SetNRGBA(x, y, transparent)
			}
		}
	}
	for x := b.Max.X - edgeScan; x < b.Max.X; x++ {
		if x < b.Min.X {
			continue
		}
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if edgeDarkOpaque(img.NRGBAAt(x, y)) {
				img.SetNRGBA(x, y, transparent)
			}
		}
	}
}

// SpriteGridFromPNGClean loads a PNG grid and cleans it thoroughly:
// 1. Removes the white/solid background via color-key sampling.
// 2. Detects and erases horizontal/vertical grid-lines between cells.
// 3. Trims each cell by `inset` pixels to drop any leftover seam.
// 4. Slices the image into [rows][cols] GridFrames with fixed cell sizes so
//    frame-to-frame Y positions stay stable (no apparent floating).
// Use inset=2 for typical AI-generated sheets with visible gridlines; inset=0
// for already-clean sheets.
func SpriteGridFromPNGClean(renderer *sdl.Renderer, filename string, cols, rows, inset int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG grid %s: %v", filename, err))
	}
	applyColorKey(img)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()
	cellW := bounds.Dx() / cols
	cellH := bounds.Dy() / rows

	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := image.Rect(
				bounds.Min.X+c*cellW+inset,
				bounds.Min.Y+r*cellH+inset,
				bounds.Min.X+(c+1)*cellW-inset,
				bounds.Min.Y+(r+1)*cellH-inset,
			)
			if cellRect.Max.X <= cellRect.Min.X || cellRect.Max.Y <= cellRect.Min.Y {
				cellRect = image.Rect(
					bounds.Min.X+c*cellW, bounds.Min.Y+r*cellH,
					bounds.Min.X+(c+1)*cellW, bounds.Min.Y+(r+1)*cellH,
				)
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h}
		}
	}

	return grid
}

// SpriteGridFromPNGCleanKids sits between the default (tol=8) and the
// aggressive (tol=32) variants. Tol=16 is wide enough to clear the
// soft-gradient backgrounds that kid sheets were authored with (cream,
// beige, pale-pink cells) without eating the saturated shirt colors or
// the skin-tone anti-aliasing that the aggressive path chews up.
// Audit the sheet: if default leaves a halo around the character and
// aggressive turns the shirt into swiss cheese, this is the one to use.
func SpriteGridFromPNGCleanKids(renderer *sdl.Renderer, filename string, cols, rows, inset int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG grid %s: %v", filename, err))
	}
	applyColorKeyTol(img, 16)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()
	cellW := bounds.Dx() / cols
	cellH := bounds.Dy() / rows

	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := image.Rect(
				bounds.Min.X+c*cellW+inset,
				bounds.Min.Y+r*cellH+inset,
				bounds.Min.X+(c+1)*cellW-inset,
				bounds.Min.Y+(r+1)*cellH-inset,
			)
			if cellRect.Max.X <= cellRect.Min.X || cellRect.Max.Y <= cellRect.Min.Y {
				cellRect = image.Rect(
					bounds.Min.X+c*cellW, bounds.Min.Y+r*cellH,
					bounds.Min.X+(c+1)*cellW, bounds.Min.Y+(r+1)*cellH,
				)
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h}
		}
	}

	return grid
}

// SpriteGridFromPNGCleanAggressive is the same pipeline as
// SpriteGridFromPNGClean but with a wider color-key tolerance (32 per
// channel vs the default 8). Use this for sheets where the background
// bleeds into anti-aliased edges and the default color-key leaves a
// visible halo — the campfire loop is the canonical case: dramatic
// flame color swings between frames move the corner-sample averages,
// so a gentle tolerance misses the fringe pixels against a saturated
// background. Do not use as a default — it can eat near-white pixels
// inside the character itself.
func SpriteGridFromPNGCleanAggressive(renderer *sdl.Renderer, filename string, cols, rows, inset int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG grid %s: %v", filename, err))
	}
	applyColorKeyTol(img, 32)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()
	cellW := bounds.Dx() / cols
	cellH := bounds.Dy() / rows

	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := image.Rect(
				bounds.Min.X+c*cellW+inset,
				bounds.Min.Y+r*cellH+inset,
				bounds.Min.X+(c+1)*cellW-inset,
				bounds.Min.Y+(r+1)*cellH-inset,
			)
			if cellRect.Max.X <= cellRect.Min.X || cellRect.Max.Y <= cellRect.Min.Y {
				cellRect = image.Rect(
					bounds.Min.X+c*cellW, bounds.Min.Y+r*cellH,
					bounds.Min.X+(c+1)*cellW, bounds.Min.Y+(r+1)*cellH,
				)
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h}
		}
	}

	return grid
}

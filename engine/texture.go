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

// applyColorKeyConnectedTol removes only background-colored pixels connected
// to the image edge. This keeps enclosed whites, such as cartoon eye whites,
// from being treated as background.
func applyColorKeyConnectedTol(img *image.NRGBA, matchTol uint8) {
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

	matchesBG := func(c color.NRGBA) bool {
		if c.A < 200 {
			return false
		}
		for _, bg := range bgColors {
			if absDiffU8(c.R, bg.R) <= matchTol && absDiffU8(c.G, bg.G) <= matchTol && absDiffU8(c.B, bg.B) <= matchTol {
				return true
			}
		}
		return false
	}

	w := b.Dx()
	h := b.Dy()
	seen := make([]bool, w*h)
	stack := make([]image.Point, 0, 2*w+2*h)
	push := func(x, y int) {
		if x < b.Min.X || x >= b.Max.X || y < b.Min.Y || y >= b.Max.Y {
			return
		}
		idx := (y-b.Min.Y)*w + (x - b.Min.X)
		if seen[idx] || !matchesBG(img.NRGBAAt(x, y)) {
			return
		}
		seen[idx] = true
		stack = append(stack, image.Point{X: x, Y: y})
	}

	for x := b.Min.X; x < b.Max.X; x++ {
		push(x, b.Min.Y)
		push(x, b.Max.Y-1)
	}
	for y := b.Min.Y + 1; y < b.Max.Y-1; y++ {
		push(b.Min.X, y)
		push(b.Max.X-1, y)
	}

	transparent := color.NRGBA{0, 0, 0, 0}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		img.SetNRGBA(p.X, p.Y, transparent)
		push(p.X+1, p.Y)
		push(p.X-1, p.Y)
		push(p.X, p.Y+1)
		push(p.X, p.Y-1)
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

// TextureFromPNGAggressive is TextureFromPNG with a wider color-key tolerance
// (16 per channel vs default 8). For UI assets like cursor PNGs whose
// generator leaves an off-white halo the default tolerance can't lift.
func TextureFromPNGAggressive(renderer *sdl.Renderer, filename string) (*sdl.Texture, int32, int32) {
	img, err := loadPNG(filename)
	if err != nil {
		panic(fmt.Errorf("loading PNG %s: %v", filename, err))
	}
	applyColorKeyTol(img, 16)
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
//
// W/H are the full cell size. OX/OY/OW/OH are the tight non-transparent
// ("opaque") content box within the cell, in cell-local coordinates. The
// renderer scales every animation frame by its opaque box so idle/talk/walk
// render at one consistent on-screen size, anchored by feet + horizontal
// centre — and so cells with empty padding (e.g. a kid drawn in the middle
// band of a tall sheet) don't render tiny or head-cropped.
type GridFrame struct {
	Tex *sdl.Texture
	W   int32
	H   int32
	OX  int32
	OY  int32
	OW  int32
	OH  int32
	// FCX is the cell-local X of the horizontal centre of the opaque pixels in
	// the bottom band of the frame — i.e. where the character's FEET are.
	// Anchoring an animation by FCX keeps a standing/walking character planted
	// even when the art drifts the body within the cell or an arm/leg extends
	// to one side (which would skew the full-box centre). 0 when no opaque data.
	FCX int32
	// FRY is the cell-local Y of the FEET LINE (exclusive bottom): the lowest
	// row of the opaque box wide enough to read as the solid feet/legs block.
	// A thin tail strand dipping below the feet is skipped, so anchoring by
	// FRY plants the feet on the ground line per frame — cancelling vertical
	// art drift WITHOUT letting the tail lift the body (user 2026-06-10:
	// "the frames place in the same spot"). 0 when no opaque data.
	FRY int32
}

// opaqueImgCache memoizes decoded NRGBA images for OpaqueBox so packed atlases
// (one PNG, many frames) are read off disk once. nil means a prior decode failed.
var opaqueImgCache = map[string]*image.NRGBA{}

// OpaqueBox returns the tight non-transparent bounding box of the (x,y,w,h)
// region inside the PNG at filename, relative to that region's top-left, using
// the PNG's native alpha (for packed atlases, which ship transparent). Falls
// back to the full region on read error. Decoded images are cached.
func OpaqueBox(filename string, x, y, w, h int32) (ox, oy, ow, oh int32) {
	im, ok := opaqueImgCache[filename]
	if !ok {
		var err error
		if im, err = loadPNG(filename); err != nil {
			im = nil
		}
		opaqueImgCache[filename] = im
	}
	if im == nil {
		return 0, 0, w, h
	}
	region := image.Rect(int(x), int(y), int(x+w), int(y+h)).Intersect(im.Bounds())
	if region.Empty() {
		return 0, 0, w, h
	}
	ob := findOpaqueBounds(im, region)
	return int32(ob.Min.X) - x, int32(ob.Min.Y) - y, int32(ob.Dx()), int32(ob.Dy())
}

// opaqueLocal returns the non-transparent bounding box of cellRect within img,
// in cell-local coordinates (relative to cellRect's top-left).
func opaqueLocal(img *image.NRGBA, cellRect image.Rectangle) (ox, oy, ow, oh int32) {
	ob := findOpaqueBounds(img, cellRect)
	return int32(ob.Min.X - cellRect.Min.X), int32(ob.Min.Y - cellRect.Min.Y),
		int32(ob.Dx()), int32(ob.Dy())
}

// footCenterLocal returns the cell-local X of the character's FEET in the bottom
// band (~bottom 12%) of the opaque box. It averages only the DENSE columns (the
// solid feet/legs block), ignoring thin strands like PP's trailing tail that dip
// into the band — those are sparse and would otherwise drag the anchor sideways
// (the "PP drifts left while talking" bug). ox/oy/ow/oh are cell-local opaque
// bounds; falls back to the box centre when there's no opaque data.
func footCenterLocal(img *image.NRGBA, cellRect image.Rectangle, ox, oy, ow, oh int32) int32 {
	if ow <= 0 || oh <= 0 {
		return ox + ow/2
	}
	bandH := oh * 12 / 100
	if bandH < 2 {
		bandH = 2
	}
	bottom := cellRect.Min.Y + int(oy+oh) // image-space bottom of opaque box (exclusive)
	top := bottom - int(bandH)
	x0 := cellRect.Min.X + int(ox)
	x1 := x0 + int(ow)
	// Per-column opaque counts in the band, and the peak (densest column).
	cols := make([]int, x1-x0)
	maxc := 0
	for y := top; y < bottom; y++ {
		for x := x0; x < x1; x++ {
			if img.NRGBAAt(x, y).A > 10 {
				cols[x-x0]++
				if cols[x-x0] > maxc {
					maxc = cols[x-x0]
				}
			}
		}
	}
	if maxc == 0 {
		return ox + ow/2
	}
	// Average only columns that are at least 45% as tall as the peak — the feet
	// and legs. The tail (a thin diagonal strand, 1-3 px per column) falls below
	// this and is excluded, so the anchor stays planted on the feet.
	thresh := maxc * 45 / 100
	if thresh < 1 {
		thresh = 1
	}
	sum, n := 0, 0
	for i, c := range cols {
		if c >= thresh {
			sum += x0 + i
			n++
		}
	}
	if n == 0 {
		return ox + ow/2
	}
	return int32(sum/n) - int32(cellRect.Min.X)
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

// ContentBoxKeyed returns the opaque content box of a PNG after the standard
// white color-key - the rectangle around the actual artwork. Inventory icons
// use it to center items by their CONTENT: generated icons often sit
// off-center inside large keyed-out margins (2026-06-11 #37).
func ContentBoxKeyed(filename string) (ox, oy, ow, oh int32) {
	img, err := loadPNG(filename)
	if err != nil {
		return 0, 0, 0, 0
	}
	applyColorKey(img)
	b := findOpaqueBounds(img, img.Bounds())
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return 0, 0, 0, 0
	}
	return int32(b.Min.X - img.Bounds().Min.X), int32(b.Min.Y - img.Bounds().Min.Y),
		int32(b.Dx()), int32(b.Dy())
}

// PNGSize returns a PNG's pixel dimensions without decoding the pixel data.
// Used by loaders that pick a grid layout from the sheet on disk (e.g. PP's
// side-walk strip shipped as 10×1 and its regen spec is 8×1 — the caller
// chooses the column count by width so either sheet just works).
func PNGSize(filename string) (int, int, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// footRowLocal returns the cell-local Y (exclusive bottom) of the character's
// FEET line: scanning the bottom ~35% of the opaque box from the bottom up,
// the first row whose opaque pixel count reads as a solid feet/legs block
// (≥ max(6, 12% of box width)) wins. Thin strands — PP's tail dipping below
// the feet — are 1-4 px per row and get skipped. Falls back to the box bottom
// when nothing qualifies (seated/cropped poses).
func footRowLocal(img *image.NRGBA, cellRect image.Rectangle, ox, oy, ow, oh int32) int32 {
	if ow <= 0 || oh <= 0 {
		return oy + oh
	}
	minWide := int(ow) * 12 / 100
	if minWide < 6 {
		minWide = 6
	}
	x0 := cellRect.Min.X + int(ox)
	x1 := x0 + int(ow)
	bottom := cellRect.Min.Y + int(oy+oh) // exclusive
	limit := bottom - int(oh)*35/100
	for y := bottom - 1; y >= limit; y-- {
		cnt := 0
		for x := x0; x < x1; x++ {
			if img.NRGBAAt(x, y).A > 10 {
				cnt++
			}
		}
		if cnt >= minWide {
			return int32(y+1) - int32(cellRect.Min.Y)
		}
	}
	return oy + oh
}

// emptyGrid returns a rows×cols grid of empty frames (Tex nil). The grid
// loaders return it when a sheet is missing or unreadable, so the game
// degrades to an invisible animation + console warning instead of panicking
// on boot (user 2026-06-10: a one-shot wired before its sheet was generated
// — "PP receive.png" — crashed the game at startup).
func emptyGrid(cols, rows int) [][]GridFrame {
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
	}
	return grid
}

// splitRuns finds runs of non-empty histogram entries and — when exactly
// `want` substantial runs can be resolved — returns `want` spans whose
// boundaries sit at the MIDPOINTS of the empty gaps between runs (the first
// span starts at 0, the last ends at len(hist)).
//
// Two repairs are attempted before giving up:
//   - too MANY runs: drop ghost specks (runs with tiny opaque mass).
//   - too FEW runs: two figures are BRIDGED by a thin connection (a tail tip
//     or fingertip touching the neighbour — the cause of cut limbs, user
//     2026-06-10). The widest run is split at its thinnest interior WAIST,
//     cutting only the bridge pixel column instead of a body.
//
// Returns nil when the count still doesn't match (genuinely missing frames),
// so the caller can fall back to proportional slicing.
func splitRuns(hist []int, want int) [][2]int {
	const noise = 2 // ≤2 opaque px in a line ≈ chroma residue, counts as empty
	type run struct{ s, e, mass, peak int }
	mkRun := func(s, e int) run {
		r := run{s: s, e: e}
		for i := s; i < e; i++ {
			r.mass += hist[i]
			if hist[i] > r.peak {
				r.peak = hist[i]
			}
		}
		return r
	}
	var runs []run
	in, s := false, 0
	for i, v := range hist {
		if v > noise {
			if !in {
				in, s = true, i
			}
		} else if in {
			runs = append(runs, mkRun(s, i))
			in = false
		}
	}
	if in {
		runs = append(runs, mkRun(s, len(hist)))
	}
	// Drop speck runs (≤2 px wide) — dust between real figures.
	kept := runs[:0]
	for _, r := range runs {
		if r.e-r.s > 2 {
			kept = append(kept, r)
		}
	}
	runs = kept
	if len(runs) == 0 {
		return nil
	}
	// Too many runs: drop ghost specks — runs whose opaque mass is tiny
	// compared to a real figure's share.
	for len(runs) > want {
		avg := 0
		for _, r := range runs {
			avg += r.mass
		}
		avg /= len(runs)
		min := 0
		for i := 1; i < len(runs); i++ {
			if runs[i].mass < runs[min].mass {
				min = i
			}
		}
		if runs[min].mass >= avg/8 {
			return nil // smallest run is substantial — genuinely extra content
		}
		runs = append(runs[:min], runs[min+1:]...)
	}
	// Too few runs: split bridged figures at the thinnest interior waist.
	for len(runs) < want {
		wi := 0
		for i := 1; i < len(runs); i++ {
			if runs[i].e-runs[i].s > runs[wi].e-runs[wi].s {
				wi = i
			}
		}
		r := runs[wi]
		w := r.e - r.s
		if w < 8 {
			return nil
		}
		// search the central 20–80% for the minimum histogram column
		cut, cutVal := -1, 1<<30
		for i := r.s + w/5; i < r.e-w/5; i++ {
			if hist[i] < cutVal {
				cut, cutVal = i, hist[i]
			}
		}
		// only cut a genuine thin BRIDGE — a waist well below the run's bulk.
		// A legitimately wide figure has no such waist and we refuse to cut it.
		if cut < 0 || cutVal > r.peak*35/100 {
			return nil
		}
		left, right := mkRun(r.s, cut), mkRun(cut, r.e)
		runs = append(runs[:wi], append([]run{left, right}, runs[wi+1:]...)...)
	}
	spans := make([][2]int, want)
	for i, r := range runs {
		s0 := 0
		if i > 0 {
			s0 = (runs[i-1].e + r.s) / 2
		}
		e0 := len(hist)
		if i < want-1 {
			e0 = (r.e + runs[i+1].s) / 2
		}
		spans[i] = [2]int{s0, e0}
	}
	return spans
}

// contentGridRects slices a (color-keyed, transparent-background) sheet into
// cols×rows cells by CONTENT instead of fixed grid lines: it finds the empty
// gaps between figures and places every cell boundary at a gap midpoint, so
// a figure drawn slightly across the mathematical grid line is NEVER cut
// (user 2026-06-10: "can't we not cut frames?"). Requires exactly the
// expected number of figures per axis — anything else (touching figures,
// blank cells, ghost specks merging runs) returns nil and the caller falls
// back to proportional slicing, which is never worse than before.
func contentGridRects(img *image.NRGBA, cols, rows int) [][]image.Rectangle {
	b := img.Bounds()
	// 1) split rows by empty pixel rows
	rowHist := make([]int, b.Dy())
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			if img.NRGBAAt(b.Min.X+x, b.Min.Y+y).A > 10 {
				rowHist[y]++
			}
		}
	}
	bands := splitRuns(rowHist, rows)
	if bands == nil {
		return nil
	}
	// 2) split each band's columns by empty pixel columns
	out := make([][]image.Rectangle, rows)
	for ri, band := range bands {
		y0, y1 := b.Min.Y+band[0], b.Min.Y+band[1]
		colHist := make([]int, b.Dx())
		for x := 0; x < b.Dx(); x++ {
			for y := y0; y < y1; y++ {
				if img.NRGBAAt(b.Min.X+x, y).A > 10 {
					colHist[x]++
				}
			}
		}
		spans := splitRuns(colHist, cols)
		if spans == nil {
			return nil
		}
		out[ri] = make([]image.Rectangle, cols)
		for ci, sp := range spans {
			out[ri][ci] = image.Rect(b.Min.X+sp[0], y0, b.Min.X+sp[1], y1)
		}
	}
	return out
}

// gridCellRect returns the rect of cell (c, r) using PROPORTIONAL boundaries:
// boundary i sits at floor(i*W/cols), so a sheet whose dimensions don't divide
// exactly by the grid (e.g. a 1535px-wide 8-column export) distributes the
// remainder across cells instead of truncating a strip off every frame.
// User 2026-06-10: "read the sprite properly" — no padding hacks, no frame
// splitting; the loader adapts to the sheet.
func gridCellRect(bounds image.Rectangle, cols, rows, c, r, inset int) image.Rectangle {
	x0 := bounds.Min.X + c*bounds.Dx()/cols + inset
	y0 := bounds.Min.Y + r*bounds.Dy()/rows + inset
	x1 := bounds.Min.X + (c+1)*bounds.Dx()/cols - inset
	y1 := bounds.Min.Y + (r+1)*bounds.Dy()/rows - inset
	if x1 <= x0 || y1 <= y0 {
		// inset would collapse the cell — fall back to the raw boundaries.
		x0 = bounds.Min.X + c*bounds.Dx()/cols
		y0 = bounds.Min.Y + r*bounds.Dy()/rows
		x1 = bounds.Min.X + (c+1)*bounds.Dx()/cols
		y1 = bounds.Min.Y + (r+1)*bounds.Dy()/rows
	}
	return image.Rect(x0, y0, x1, y1)
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
		fmt.Printf("Warning: could not load PNG grid %s: %v\n", filename, err)
		return emptyGrid(cols, rows)
	}

	bounds := img.Bounds()

	contentRects := contentGridRects(img, cols, rows)
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := gridCellRect(bounds, cols, rows, c, r, 0)
			if contentRects != nil {
				// Gap-detected cell — the whole figure is inside, nothing cut.
				cellRect = contentRects[r][c]
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			ox, oy, ow, oh := opaqueLocal(img, cellRect)
			fcx := footCenterLocal(img, cellRect, ox, oy, ow, oh)
			fry := footRowLocal(img, cellRect, ox, oy, ow, oh)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h, OX: ox, OY: oy, OW: ow, OH: oh, FCX: fcx, FRY: fry}
		}
	}
	return grid
}

func SpriteGridFromPNG(renderer *sdl.Renderer, filename string, cols, rows int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		fmt.Printf("Warning: could not load PNG grid %s: %v\n", filename, err)
		return emptyGrid(cols, rows)
	}
	applyColorKey(img)

	bounds := img.Bounds()

	contentRects := contentGridRects(img, cols, rows)
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := gridCellRect(bounds, cols, rows, c, r, 0)
			if contentRects != nil {
				// Gap-detected cell — the whole figure is inside, nothing cut.
				cellRect = contentRects[r][c]
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			ox, oy, ow, oh := opaqueLocal(img, cellRect)
			fcx := footCenterLocal(img, cellRect, ox, oy, ow, oh)
			fry := footRowLocal(img, cellRect, ox, oy, ow, oh)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h, OX: ox, OY: oy, OW: ow, OH: oh, FCX: fcx, FRY: fry}
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

	isDark := func(c color.NRGBA) bool {
		if c.A < 40 {
			return false
		}
		return c.R < 50 && c.G < 50 && c.B < 50
	}

	scanThickness := 2
	transparent := color.NRGBA{0, 0, 0, 0}

	for c := 1; c < cols; c++ {
		centerX := b.Min.X + c*w/cols
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
		centerY := b.Min.Y + r*h/rows
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
//  1. Removes the white/solid background via color-key sampling.
//  2. Detects and erases horizontal/vertical grid-lines between cells.
//  3. Trims each cell by `inset` pixels to drop any leftover seam.
//  4. Slices the image into [rows][cols] GridFrames with fixed cell sizes so
//     frame-to-frame Y positions stay stable (no apparent floating).
//
// Use inset=2 for typical AI-generated sheets with visible gridlines; inset=0
// for already-clean sheets.
func SpriteGridFromPNGClean(renderer *sdl.Renderer, filename string, cols, rows, inset int) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		fmt.Printf("Warning: could not load PNG grid %s: %v\n", filename, err)
		return emptyGrid(cols, rows)
	}
	applyColorKey(img)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()

	contentRects := contentGridRects(img, cols, rows)
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := gridCellRect(bounds, cols, rows, c, r, inset)
			if contentRects != nil {
				// Gap-detected cell — the whole figure is inside, nothing cut.
				cellRect = contentRects[r][c]
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			ox, oy, ow, oh := opaqueLocal(img, cellRect)
			fcx := footCenterLocal(img, cellRect, ox, oy, ow, oh)
			fry := footRowLocal(img, cellRect, ox, oy, ow, oh)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h, OX: ox, OY: oy, OW: ow, OH: oh, FCX: fcx, FRY: fry}
		}
	}

	return grid
}

// SpriteGridFromPNGCleanConnected is like SpriteGridFromPNGClean, but removes
// only background-colored pixels connected to the image edge. Use it for
// character sheets whose eye whites or teeth match the white background.
func SpriteGridFromPNGCleanConnected(renderer *sdl.Renderer, filename string, cols, rows, inset int) [][]GridFrame {
	return SpriteGridFromPNGCleanConnectedTol(renderer, filename, cols, rows, inset, 8)
}

// SpriteGridFromPNGCleanConnectedTol is SpriteGridFromPNGCleanConnected with a
// caller-chosen match tolerance. Tol 8 leaves a visible fringe on sheets whose
// background was authored with soft anti-aliased edges (the Paris biker,
// 2026-06-12 #12) - a wider tolerance eats the halo while the edge-connected
// flood still protects interior whites.
func SpriteGridFromPNGCleanConnectedTol(renderer *sdl.Renderer, filename string, cols, rows, inset int, tol uint8) [][]GridFrame {
	img, err := loadPNG(filename)
	if err != nil {
		fmt.Printf("Warning: could not load PNG grid %s: %v\n", filename, err)
		return emptyGrid(cols, rows)
	}
	applyColorKeyConnectedTol(img, tol)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()

	contentRects := contentGridRects(img, cols, rows)
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := gridCellRect(bounds, cols, rows, c, r, inset)
			if contentRects != nil {
				// Gap-detected cell — the whole figure is inside, nothing cut.
				cellRect = contentRects[r][c]
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			ox, oy, ow, oh := opaqueLocal(img, cellRect)
			fcx := footCenterLocal(img, cellRect, ox, oy, ow, oh)
			fry := footRowLocal(img, cellRect, ox, oy, ow, oh)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h, OX: ox, OY: oy, OW: ow, OH: oh, FCX: fcx, FRY: fry}
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
		fmt.Printf("Warning: could not load PNG grid %s: %v\n", filename, err)
		return emptyGrid(cols, rows)
	}
	applyColorKeyTol(img, 16)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()

	contentRects := contentGridRects(img, cols, rows)
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := gridCellRect(bounds, cols, rows, c, r, inset)
			if contentRects != nil {
				// Gap-detected cell — the whole figure is inside, nothing cut.
				cellRect = contentRects[r][c]
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			ox, oy, ow, oh := opaqueLocal(img, cellRect)
			fcx := footCenterLocal(img, cellRect, ox, oy, ow, oh)
			fry := footRowLocal(img, cellRect, ox, oy, ow, oh)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h, OX: ox, OY: oy, OW: ow, OH: oh, FCX: fcx, FRY: fry}
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
		fmt.Printf("Warning: could not load PNG grid %s: %v\n", filename, err)
		return emptyGrid(cols, rows)
	}
	applyColorKeyTol(img, 32)
	eraseGridLines(img, cols, rows)

	bounds := img.Bounds()

	contentRects := contentGridRects(img, cols, rows)
	grid := make([][]GridFrame, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]GridFrame, cols)
		for c := 0; c < cols; c++ {
			cellRect := gridCellRect(bounds, cols, rows, c, r, inset)
			if contentRects != nil {
				// Gap-detected cell — the whole figure is inside, nothing cut.
				cellRect = contentRects[r][c]
			}
			tex, w, h := nrgbaToTexture(renderer, img, cellRect)
			ox, oy, ow, oh := opaqueLocal(img, cellRect)
			fcx := footCenterLocal(img, cellRect, ox, oy, ow, oh)
			fry := footRowLocal(img, cellRect, ox, oy, ow, oh)
			grid[r][c] = GridFrame{Tex: tex, W: w, H: h, OX: ox, OY: oy, OW: ow, OH: oh, FCX: fcx, FRY: fry}
		}
	}

	return grid
}

package game

import (
	"image"
	"image/color"
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type bgLayer struct {
	tex      *sdl.Texture
	srcW     int32
	srcH     int32
	parallax float64
}

type background struct {
	tex    *sdl.Texture
	srcW   int32
	srcH   int32
	layers []bgLayer
}

func newPNGBackground(renderer *sdl.Renderer, path string) *background {
	tex, w, h := engine.TextureFromPNGRawClean(renderer, path)
	return &background{tex: tex, srcW: w, srcH: h}
}

// newPNGBackgroundOr loads a PNG if present, otherwise falls back to a solid
// placeholder so scenes can ship before their art does. Missing-art rooms
// still render (as a flat gradient) instead of crashing the whole game.
func newPNGBackgroundOr(renderer *sdl.Renderer, path string, fallback color.NRGBA) *background {
	if _, err := os.Stat(path); err == nil {
		return newPNGBackground(renderer, path)
	}
	return newPlaceholderBackground(renderer, fallback)
}

// newPlaceholderBackground produces a synthetic gradient background at the
// screen's native resolution. Used for city scenes whose PNG art has not
// been authored yet.
func newPlaceholderBackground(renderer *sdl.Renderer, base color.NRGBA) *background {
	w := int(engine.ScreenWidth)
	h := int(engine.ScreenHeight)
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		// Simple top-light / bottom-dark vertical gradient so the scene
		// doesn't look like a debug flat.
		t := float64(y) / float64(h-1)
		r := uint8(float64(base.R)*(1.0-0.4*t) + 20*t)
		g := uint8(float64(base.G)*(1.0-0.4*t) + 20*t)
		b := uint8(float64(base.B)*(1.0-0.4*t) + 20*t)
		base := y*img.Stride
		for x := 0; x < w; x++ {
			img.Pix[base+x*4+0] = r
			img.Pix[base+x*4+1] = g
			img.Pix[base+x*4+2] = b
			img.Pix[base+x*4+3] = 255
		}
	}
	tex, srcW, srcH := engine.TextureFromNRGBA(renderer, img)
	return &background{tex: tex, srcW: srcW, srcH: srcH}
}

func (b *background) draw(renderer *sdl.Renderer, playerX float64) {
	if len(b.layers) == 0 {
		renderer.Copy(
			b.tex,
			&sdl.Rect{X: 0, Y: 0, W: b.srcW, H: b.srcH},
			&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight},
		)
		return
	}
	screenCenter := float64(engine.ScreenWidth) / 2.0
	for _, l := range b.layers {
		offsetX := int32((playerX - screenCenter) * l.parallax)
		renderer.Copy(l.tex, nil,
			&sdl.Rect{X: -offsetX, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})
	}
}

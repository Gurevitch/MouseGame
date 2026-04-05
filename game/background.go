package game

import (
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

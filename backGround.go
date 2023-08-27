package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type background struct {
	tex *sdl.Texture
}

func newBackground(renderer *sdl.Renderer) (newBackground background, err error) {

	newBackground.tex = textureFromBPM(renderer, "images/background.bmp")
	return newBackground, nil
}

func (bkrnd *background) draw(renderer *sdl.Renderer) {

	renderer.Copy(
		bkrnd.tex,
		&sdl.Rect{X: 0, Y: 0, W: 626, H: 626},
		&sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight},
	)
}

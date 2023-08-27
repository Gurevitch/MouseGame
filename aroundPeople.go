package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type aroundPeople struct {
	tex          *sdl.Texture
	xMove, yMove float32
}

func newPaparMan(renderer *sdl.Renderer) (newAroundPeople aroundPeople, err error) {

	newAroundPeople.tex = textureFromBPM(renderer, "images/paparman.bmp")
	return newAroundPeople, nil
}
func (plr *aroundPeople) draw(renderer *sdl.Renderer) {

	renderer.Copy(
		plr.tex,
		&sdl.Rect{X: 0, Y: 0, W: 370, H: 346},
		&sdl.Rect{X: screenWidth / 2.0, Y: screenHeight / 2.0, W: 370, H: 346},
	)
}

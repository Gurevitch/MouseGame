package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type aroundPeople struct {
	tex          *sdl.Texture
	xMove, yMove float32
}

func newPaparMan(renderer *sdl.Renderer) (newPlayer aroundPeople, err error) {

	img, err := sdl.LoadBMP("images/paparman.bmp")
	if err != nil {
		return aroundPeople{}, fmt.Errorf("loading player img: %v", err)
	}
	defer img.Free()
	newPlayer.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		return aroundPeople{}, fmt.Errorf("creating player texture: %v", err)
	}
	return newPlayer, nil
}
func (plr *aroundPeople) draw(renderer *sdl.Renderer) {

	renderer.Copy(
		plr.tex,
		&sdl.Rect{X: 0, Y: 0, W: 370, H: 346},
		&sdl.Rect{X: screenWidth / 2.0, Y: screenHeight / 2.0, W: 370, H: 346},
	)
}

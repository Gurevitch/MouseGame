package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	leftMouseClick   = 1
	middleMouseClick = 2
	rightMouseClick  = 4
	playerSpeed      = 0.3
)

type player struct {
	tex   *sdl.Texture
	xMove float32
	yMove float32
}
type movingPlayerToNewPlace struct {
	xRightOrLeft    int32
	yUpOrDown       int32
	mouseKeyPressed uint32
}

var newPlace movingPlayerToNewPlace

func newPlayer(renderer *sdl.Renderer) (newPlayer player, err error) {

	img, err := sdl.LoadBMP("images/player.bmp")
	if err != nil {
		return player{}, fmt.Errorf("loading player img: %v", err)
	}
	defer img.Free()
	newPlayer.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		return player{}, fmt.Errorf("creating player texture: %v", err)
	}
	return newPlayer, nil
}

func (plr *player) draw(renderer *sdl.Renderer) {

	renderer.Copy(
		plr.tex,
		&sdl.Rect{X: 0, Y: 0, W: 242, H: 580},
		&sdl.Rect{X: int32(plr.xMove), Y: int32(plr.yMove), W: 200, H: 300},
	)
}

var needToMove = false

func (plr *player) update() {
	newPlace.xRightOrLeft, newPlace.yUpOrDown, newPlace.mouseKeyPressed = sdl.GetMouseState()

	if (newPlace.mouseKeyPressed == leftMouseClick) && (int32(plr.xMove) != newPlace.xRightOrLeft || int32(plr.yMove) != newPlace.yUpOrDown) {
		needToMove = true
	}
	if needToMove {

		//moving player right or left
		if int32(plr.xMove) != newPlace.xRightOrLeft {
			if newPlace.xRightOrLeft > int32(plr.xMove) /* && newPlace.xRightOrLeft+242 < 1244 */ {
				plr.xMove += playerSpeed
			}
			if newPlace.xRightOrLeft < int32(plr.xMove) {
				plr.xMove -= playerSpeed
			}
		}
		//moving player up or down
		if int32(plr.yMove) != newPlace.yUpOrDown {
			if int32(plr.yMove) != newPlace.yUpOrDown /*&& newPlace.yUpOrDown+582 < 1070*/ {
				if newPlace.yUpOrDown > int32(plr.yMove) {
					plr.yMove += playerSpeed
				}
				if newPlace.yUpOrDown < int32(plr.yMove) {
					plr.yMove -= playerSpeed
				}
			}
		}
		if int32(plr.xMove) == newPlace.xRightOrLeft && int32(plr.yMove) == newPlace.yUpOrDown {
			needToMove = false
		}
	}

}

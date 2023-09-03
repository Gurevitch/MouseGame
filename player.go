package main

import (
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
	mouseKeyPressed sdl.ButtonMask
}

var (
	newPlace   movingPlayerToNewPlace
	needToMove = false
)

func newPlayer(renderer *sdl.Renderer) (newPlayer player, err error) {

	newPlayer.tex = textureFromBPM(renderer, "images/player.bmp")

	return newPlayer, nil
}

func (plr *player) draw(renderer *sdl.Renderer) {

	renderer.Copy(
		plr.tex,
		&sdl.Rect{X: 0, Y: 0, W: 242, H: 582},
		&sdl.Rect{X: int32(plr.xMove), Y: int32(plr.yMove), W: 105, H: 150},
	)
}

func (plr *player) update() {
	if !needToMove {
		newPlace.xRightOrLeft, newPlace.yUpOrDown, newPlace.mouseKeyPressed = sdl.GetMouseState()
		if newPlace.mouseKeyPressed == leftMouseClick {
			needToMove = true
		}
	}
	keyPressed := sdl.GetKeyboardState()
	if keyPressed[sdl.SCANCODE_SPACE] == truePress {
		plr.xMove = float32(newPlace.xRightOrLeft)
		plr.yMove = float32(newPlace.yUpOrDown)
		needToMove = false

	}
	if needToMove {
		//moving player right or left
		if int32(plr.xMove) != newPlace.xRightOrLeft {
			if newPlace.xRightOrLeft > int32(plr.xMove) {
				plr.xMove += playerSpeed
			}
			if newPlace.xRightOrLeft < int32(plr.xMove) {
				plr.xMove -= playerSpeed
			}
		}
		//moving player up or down
		if int32(plr.yMove) != newPlace.yUpOrDown {
			if int32(plr.yMove) != newPlace.yUpOrDown {
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

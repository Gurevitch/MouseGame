package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 1200
	screenHeight = 800

	truePress  = 1
	falsePress = 0
)

func textureFromBPM(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("failed to load BMP picture: %v", err))
	}
	defer img.Free()
	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("failed to create texture from surface: %v", err))
	}
	return tex
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("Init SDL problem: ", err)
		return
	}
	window, err := sdl.CreateWindow(
		"basic gaming practice in go", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)

	if err != nil {
		fmt.Println("Init Window problem: ", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Init renderer problem: ", err)
		return
	}
	defer renderer.Destroy()

	background, err := newBackground(renderer)
	if err != nil {
		fmt.Println("creating player: ", err)
		return
	}

	plr, err := newPlayer(renderer)
	if err != nil {
		fmt.Println("creating player: ", err)
		return
	}

	pprMan, err := newPaparMan(renderer)
	if err != nil {
		fmt.Println("creating paparMan: ", err)
		return
	}
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()

		background.draw(renderer)
		pprMan.draw(renderer)
		plr.update()
		plr.draw(renderer)

		renderer.Present()
	}
}

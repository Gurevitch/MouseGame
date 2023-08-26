package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

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
		pprMan.draw(renderer)
		plr.update()
		plr.draw(renderer)

		renderer.Present()
	}
}

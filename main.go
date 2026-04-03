package main

import (
	"fmt"

	"bitbucket.org/Local/games/PP/engine"
	"bitbucket.org/Local/games/PP/game"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("Init SDL:", err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Pink Panther Adventure",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		engine.ScreenWidth, engine.ScreenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("Window:", err)
		return
	}
	defer window.Destroy()

	if iconSurf, err := engine.SurfaceFromPNG("assets/images/pp_icon.png"); err == nil {
		window.SetIcon(iconSurf)
		iconSurf.Free()
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Renderer:", err)
		return
	}
	defer renderer.Destroy()
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	sdl.ShowCursor(sdl.DISABLE)

	font := engine.NewBitmapFont(renderer)
	g := game.New(renderer, font)
	defer g.Close()

	var lastTick uint32 = sdl.GetTicks()

	for {
		frameStart := sdl.GetTicks()
		dt := float64(frameStart-lastTick) / 1000.0
		if dt > 0.05 {
			dt = 0.05
		}
		lastTick = frameStart

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.MouseButtonEvent:
				if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == sdl.BUTTON_LEFT {
					g.HandleClick(e.X, e.Y)
				}
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					g.HandleKey(e.Keysym.Scancode)
				}
			}
		}

		mx, my, _ := sdl.GetMouseState()
		g.Update(dt, mx, my)

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		g.Draw(renderer)
		renderer.Present()

		if elapsed := sdl.GetTicks() - frameStart; elapsed < engine.FrameDelay {
			sdl.Delay(engine.FrameDelay - elapsed)
		}
	}
}

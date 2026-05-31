// The Windows taskbar / .exe icon is embedded via go-winres: it reads
// winres/winres.json (which points at assets/icons/pp_app.png) and emits
// rsrc_windows_*.syso at the repo root, which `go build` links into PP.exe
// automatically. Regenerate after changing the icon with:
//
//	go generate ./...        (runs the directive below)
//
//go:generate go-winres make --out rsrc
package main

import (
	"fmt"
	"os"

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

	// Get desktop display size for fullscreen
	dm, _ := sdl.GetDesktopDisplayMode(0)
	window, err := sdl.CreateWindow(
		"Pink Panther Adventure",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		dm.W, dm.H,
		sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		fmt.Println("Window:", err)
		return
	}
	defer window.Destroy()

	// In-app window icon (complements the embedded taskbar icon above). Prefer
	// the new assets/icons/pp_app.png; fall back to the legacy path.
	iconPath := "assets/icons/pp_app.png"
	if _, err := os.Stat(iconPath); err != nil {
		iconPath = "assets/images/pp_icon.png"
	}
	if iconSurf, err := engine.SurfaceFromPNG(iconPath); err == nil {
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
	renderer.SetLogicalSize(engine.ScreenWidth, engine.ScreenHeight)

	// Map raw window/desktop mouse pixels into the game's logical 1400×800
	// space. The window is fullscreen-desktop (e.g. 1920×1080) with a logical
	// render size, letterboxed + uniformly scaled.
	//
	// IMPORTANT SDL2 asymmetry (the source of a long-standing click-offset
	// bug): when SetLogicalSize is active, SDL2 auto-translates coords
	// embedded in SDL_MouseMotion / SDL_MouseButton EVENTS into logical
	// space — but sdl.GetMouseState() still returns WINDOW pixels (it
	// queries the OS directly and bypasses SDL's event translation).
	//
	// So we ONLY convert GetMouseState output, never MouseButtonEvent
	// fields. Previously we converted both, which double-divided click
	// coords and made the player have to click ~273 px right and ~180 px
	// below an NPC to land the hit (user 2026-05-31: "to talk with
	// Higgins, instead of around 855,474 i clicked around 1128,655").
	toLogical := func(x, y int32) (int32, int32) {
		lx, ly := renderer.RenderWindowToLogical(int(x), int(y))
		return int32(lx + 0.5), int32(ly + 0.5)
	}

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
					// e.X/e.Y are already in logical coords — see toLogical
					// comment above. Do NOT pass through toLogical or the
					// click lands far off-target.
					g.HandleClick(e.X, e.Y)
				}
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					g.HandleKey(e.Keysym.Scancode)
				}
			}
		}

		mx, my, _ := sdl.GetMouseState()
		lmx, lmy := toLogical(mx, my)
		g.Update(dt, lmx, lmy)

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		g.Draw(renderer)
		renderer.Present()

		if elapsed := sdl.GetTicks() - frameStart; elapsed < engine.FrameDelay {
			sdl.Delay(engine.FrameDelay - elapsed)
		}
	}
}

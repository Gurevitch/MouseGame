package game

import "github.com/veandco/go-sdl2/sdl"

func drawShadow(renderer *sdl.Renderer, centerX, footY, width int32) {
	for i := int32(0); i < 5; i++ {
		w := width - i*8
		h := int32(10) - i*2
		if w < 4 || h < 1 {
			break
		}
		renderer.SetDrawColor(0, 0, 0, uint8(35-i*7))
		renderer.FillRect(&sdl.Rect{
			X: centerX - w/2, Y: footY - h/2, W: w, H: h,
		})
	}
}

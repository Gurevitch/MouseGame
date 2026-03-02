package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type npc struct {
	tex       *sdl.Texture
	bounds    sdl.Rect
	srcRect   sdl.Rect
	dialog    []dialogEntry
	name      string
	bobTimer  float64
	bobAmount float64
}

func newPaparMan(renderer *sdl.Renderer) *npc {
	tex := engine.TextureFromBMP(renderer, "assets/images/paparman.bmp")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: 370, H: 346},
		bounds:  sdl.Rect{X: 650, Y: 420, W: 180, H: 168},
		name:    "Paparazzi Man",
		dialog: []dialogEntry{
			{speaker: "Paparazzi Man", text: "Hey! You're the Pink Panther! Hold still for a photo!"},
			{speaker: "Pink Panther", text: "..."},
			{speaker: "Paparazzi Man", text: "Come on, just one shot! The tabloids will pay a fortune!"},
			{speaker: "Pink Panther", text: "I'd rather not, thank you very much."},
		},
		bobAmount: 2.0,
	}
}

func (n *npc) update(dt float64) {
	n.bobTimer += dt
}

func (n *npc) draw(renderer *sdl.Renderer) {
	bobOffset := int32(math.Sin(n.bobTimer*1.5) * n.bobAmount)
	dst := sdl.Rect{X: n.bounds.X, Y: n.bounds.Y + bobOffset, W: n.bounds.W, H: n.bounds.H}
	renderer.Copy(n.tex, &n.srcRect, &dst)
}

func (n *npc) containsPoint(x, y int32) bool {
	pt := sdl.Point{X: x, Y: y}
	return pt.InRect(&n.bounds)
}

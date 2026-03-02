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
	flipped   bool
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

func newCryingKid(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/npc_crying_kid.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 120, Y: 380, W: 160, H: 120},
		name:    "Crying Kid",
		dialog: []dialogEntry{
			{speaker: "Crying Kid", text: "*sniff* I... I don't want to be here anymore..."},
			{speaker: "Crying Kid", text: "I miss my mum and dad! I want to go home!"},
			{speaker: "Pink Panther", text: "There there, little one. What happened?"},
			{speaker: "Crying Kid", text: "They sent me to this camp and everyone is so mean!"},
			{speaker: "Crying Kid", text: "Please... can you help me get back home?"},
			{speaker: "Pink Panther", text: "Don't worry. I'll figure something out."},
		},
		bobAmount: 1.2,
	}
}

func newProfessor(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/npc_professor.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 900, Y: 330, W: 130, H: 200},
		name:    "Professor",
		dialog: []dialogEntry{
			{speaker: "Professor", text: "Ah, Pink Panther! Thank goodness you're here!"},
			{speaker: "Professor", text: "This poor child has been crying all day long."},
			{speaker: "Professor", text: "We have to do something! We can't just leave him like this!"},
			{speaker: "Pink Panther", text: "What do you suggest, Professor?"},
			{speaker: "Professor", text: "I've been studying the camp schedules. There might be a way out..."},
			{speaker: "Professor", text: "But we'll need to be clever about it. Very clever indeed!"},
			{speaker: "Pink Panther", text: "Clever is my middle name."},
		},
		bobAmount: 1.5,
		flipped:   true,
	}
}

func (n *npc) update(dt float64) {
	n.bobTimer += dt
}

func (n *npc) draw(renderer *sdl.Renderer) {
	bobOffset := int32(math.Sin(n.bobTimer*1.5) * n.bobAmount)
	dst := sdl.Rect{X: n.bounds.X, Y: n.bounds.Y + bobOffset, W: n.bounds.W, H: n.bounds.H}
	flip := sdl.FLIP_NONE
	if n.flipped {
		flip = sdl.FLIP_HORIZONTAL
	}
	renderer.CopyEx(n.tex, &n.srcRect, &dst, 0, nil, flip)
}

func (n *npc) containsPoint(x, y int32) bool {
	pt := sdl.Point{X: x, Y: y}
	return pt.InRect(&n.bounds)
}

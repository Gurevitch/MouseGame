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
	hovered   bool
}

func newPaparMan(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/npc_paparman.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 650, Y: 380, W: 150, H: 210},
		name:    "Paper Man",
		dialog: []dialogEntry{
			{speaker: "Paper Man", text: "Extra! Extra! Read all about it! Pink Panther spotted in London!"},
			{speaker: "Pink Panther", text: "..."},
			{speaker: "Paper Man", text: "Care to buy a paper, sir? Got all the latest news!"},
			{speaker: "Pink Panther", text: "No thank you, I prefer to make the news, not read it."},
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

	// Breathing scale pulse: ~1% oscillation at 0.8 Hz
	breathScale := 1.0 + 0.01*math.Sin(n.bobTimer*0.8*2*math.Pi)
	dstW := int32(float64(n.bounds.W) * breathScale)
	dstH := int32(float64(n.bounds.H) * breathScale)
	dstX := n.bounds.X - (dstW-n.bounds.W)/2
	dstY := n.bounds.Y + bobOffset - (dstH - n.bounds.H)

	dst := sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH}

	shadowCX := n.bounds.X + n.bounds.W/2
	shadowFY := n.bounds.Y + n.bounds.H
	drawShadow(renderer, shadowCX, shadowFY, n.bounds.W-10)

	if n.hovered {
		renderer.SetDrawColor(255, 220, 100, 35)
		pad := int32(4)
		renderer.FillRect(&sdl.Rect{
			X: dst.X - pad, Y: dst.Y - pad,
			W: dst.W + pad*2, H: dst.H + pad*2,
		})
		renderer.SetDrawColor(255, 220, 100, 90)
		renderer.DrawRect(&sdl.Rect{
			X: dst.X - pad, Y: dst.Y - pad,
			W: dst.W + pad*2, H: dst.H + pad*2,
		})
	}

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

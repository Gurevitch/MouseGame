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
	itemMatch bool
	elevated  bool
	silent    bool
	groupID   string

	dialogDone    bool
	onDialogEnd   func()
	altDialogFunc func() ([]dialogEntry, func())
}

func newPaparMan(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/locations/london/npc/paperman.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 1020, Y: 366, W: 180, H: 145},
		name:    "Paper Man",
		dialog: []dialogEntry{
			{speaker: "Paper Man", text: "Extra! Extra! Read all about it! Pink Panther spotted in London!"},
			{speaker: "Pink Panther", text: "..."},
			{speaker: "Paper Man", text: "Care to buy a paper, sir? Got all the latest news!"},
			{speaker: "Pink Panther", text: "No thank you, I prefer to make the news, not read it."},
			{speaker: "Paper Man", text: "Well then, take this comic at least. Free of charge for a celebrity!"},
			{speaker: "Pink Panther", text: "A comic book? Well... don't mind if I do!"},
		},
		bobAmount: 1.0,
		elevated:  true,
	}
}

var paperManPostComicDialog = []dialogEntry{
	{speaker: "Paper Man", text: "Back again, eh? No more free comics, I'm afraid!"},
	{speaker: "Pink Panther", text: "That's alright. Any interesting headlines today?"},
	{speaker: "Paper Man", text: "Just the usual... thefts, scandals, and the occasional cat burglar."},
	{speaker: "Pink Panther", text: "Cat burglar? I take offense to that."},
}

var streetTalkersDialog = []dialogEntry{
	{speaker: "Woman", text: "Did you hear? There's been a diamond theft at the museum!"},
	{speaker: "Gentleman", text: "Good heavens! Scotland Yard must be in a tizzy."},
	{speaker: "Young Man", text: "I bet it was that inspector... what's his name... Clouseau?"},
	{speaker: "Woman", text: "Oh no dear, he's French. This is a London matter!"},
	{speaker: "Pink Panther", text: "Excuse me, may I pass through?"},
	{speaker: "Gentleman", text: "Sorry old chap, we're rather engrossed in conversation. Try going around!"},
	{speaker: "Pink Panther", text: "Hmm... a diamond theft, you say? Interesting..."},
}

func newStreetTalkers(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/locations/london/npc/3 person.png")
	return &npc{
		tex:       tex,
		srcRect:   sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:    sdl.Rect{X: 40, Y: 500, W: 350, H: 210},
		name:      "Street Talkers",
		dialog:    streetTalkersDialog,
		bobAmount: 1.0,
	}
}

func newGrumpyKid(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/locations/london/npc/grumpy kid.png")
	return &npc{
		tex:       tex,
		srcRect:   sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:    sdl.Rect{X: 5, Y: 520, W: 80, H: 170},
		name:      "Grumpy Kid",
		bobAmount: 0.8,
		silent:    true,
	}
}

var cryingKidDefaultDialog = []dialogEntry{
	{speaker: "Crying Kid", text: "*sniff* I... I don't want to be here anymore..."},
	{speaker: "Crying Kid", text: "I miss my mum and dad! I want to go home!"},
	{speaker: "Pink Panther", text: "There there, little one. What happened?"},
	{speaker: "Crying Kid", text: "They sent me to this camp and everyone is so mean!"},
	{speaker: "Crying Kid", text: "Please... can you help me get back home?"},
	{speaker: "Pink Panther", text: "Don't worry. I'll figure something out."},
}

var cryingKidComicDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Hey there, little one. Look what I got for you!"},
	{speaker: "Crying Kid", text: "*sniff* W-what is it?"},
	{speaker: "Pink Panther", text: "A comic book! It will cheer you up."},
	{speaker: "Crying Kid", text: "A comic book?! Really?! For me?!"},
	{speaker: "Crying Kid", text: "Oh wow, thank you so much! This is my favorite series!"},
	{speaker: "Pink Panther", text: "There you go. A smile suits you much better."},
}

var cryingKidHappyDialog = []dialogEntry{
	{speaker: "Happy Kid", text: "This comic is so cool! Thank you, Pink Panther!"},
	{speaker: "Pink Panther", text: "Glad to see you smiling."},
}

func newCryingKid(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/crying_kid/sprite.png")
	return &npc{
		tex:       tex,
		srcRect:   sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:    sdl.Rect{X: 140, Y: 380, W: 150, H: 130},
		name:      "Crying Kid",
		dialog:    cryingKidDefaultDialog,
		bobAmount: 1.2,
	}
}

func newProfessor(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/professor/sprite.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 930, Y: 500, W: 110, H: 200},
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

	if n.itemMatch {
		pad := int32(6)
		renderer.SetDrawColor(255, 240, 50, 180)
		for i := int32(0); i < 3; i++ {
			renderer.DrawRect(&sdl.Rect{
				X: dst.X - pad + i, Y: dst.Y - pad + i,
				W: dst.W + (pad-i)*2, H: dst.H + (pad-i)*2,
			})
		}
	} else if n.hovered {
		pad := int32(4)
		renderer.SetDrawColor(255, 220, 100, 35)
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

func (n *npc) footY() int32 {
	return n.bounds.Y + n.bounds.H
}

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

	frames     []sdl.Rect
	frameIdx   int
	frameTimer float64
	frameSpeed float64

	talkFrames     []sdl.Rect
	drinkFrames    []sdl.Rect
	talkFrameSpeed float64
	animState      int
	animOnce       bool
}

func newPaperMan(renderer *sdl.Renderer) *npc {
	const (
		cols = 7
		rows = 3
	)
	tex, w, h := engine.TextureFromPNGRaw(renderer, "assets/images/locations/london/npc/kiosk_idle_cropped.png")
	frameW := w / cols
	frameH := h / rows

	rowRects := func(row int32) []sdl.Rect {
		rects := make([]sdl.Rect, cols)
		for i := int32(0); i < cols; i++ {
			rects[i] = sdl.Rect{X: i * frameW, Y: row * frameH, W: frameW, H: frameH}
		}
		return rects
	}

	frameSequence := func(rects []sdl.Rect, indices ...int) []sdl.Rect {
		seq := make([]sdl.Rect, len(indices))
		for i, idx := range indices {
			seq[i] = rects[idx]
		}
		return seq
	}

	idleRow := rowRects(2)
	talkRow := rowRects(1)
	idleFrames := frameSequence(idleRow, 0, 0, 1, 1, 0, 2, 1, 0)
	talkFrames := frameSequence(talkRow, 0, 1, 2, 3, 4, 5, 6, 5, 4, 3, 2, 1)

	return &npc{
		tex:     tex,
		srcRect: idleFrames[0],
		bounds:  sdl.Rect{X: 1036, Y: 408, W: 124, H: 104},
		name:    "Paper Man",
		dialog: []dialogEntry{
			{speaker: "Paper Man", text: "Extra! Extra! Read all about it! Pink Panther spotted in London!"},
			{speaker: "Pink Panther", text: "..."},
			{speaker: "Paper Man", text: "Care to buy a paper, sir? Got all the latest news!"},
			{speaker: "Pink Panther", text: "No thank you, I prefer to make the news, not read it."},
			{speaker: "Paper Man", text: "Well then, take this comic at least. Free of charge for a celebrity!"},
			{speaker: "Pink Panther", text: "A comic book? Well... don't mind if I do!"},
		},
		bobAmount:      0.15,
		elevated:       true,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.24,
		talkFrameSpeed: 0.10,
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
		bounds:    sdl.Rect{X: 40, Y: 520, W: 350, H: 210},
		name:      "Street Talkers",
		dialog:    streetTalkersDialog,
		bobAmount: 1.0,
	}
}

func newGrumpyKid(renderer *sdl.Renderer) *npc {
	const numFrames = 6
	tex, w, h := engine.TextureFromPNGKeyed(renderer, "assets/images/locations/london/npc/kid_idle.png")
	frameW := w / numFrames
	frames := make([]sdl.Rect, numFrames)
	for i := int32(0); i < numFrames; i++ {
		frames[i] = sdl.Rect{X: i * frameW, Y: 0, W: frameW, H: h}
	}
	return &npc{
		tex:        tex,
		srcRect:    frames[0],
		bounds:     sdl.Rect{X: 5, Y: 520, W: 80, H: 170},
		name:       "Grumpy Kid",
		bobAmount:  0.8,
		silent:     true,
		frames:     frames,
		frameSpeed: 1.0,
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

func newBarmaid(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/locations/london/npc/pub_barmaid.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 590, Y: 332, W: 200, H: 235},
		name:    "Barmaid",
		dialog: []dialogEntry{
			{speaker: "Barmaid", text: "Welcome to The Mucky Duck, love! What can I get you?"},
			{speaker: "Pink Panther", text: "Just information, if you don't mind."},
			{speaker: "Barmaid", text: "Information? This ain't the library, darling!"},
			{speaker: "Barmaid", text: "But seein' as you're a celebrity... I might know a thing or two."},
			{speaker: "Pink Panther", text: "I heard there was a diamond theft at the museum."},
			{speaker: "Barmaid", text: "Oh yes, big news that. Scotland Yard's been all over it."},
			{speaker: "Barmaid", text: "Word is, the thief was spotted heading toward the countryside."},
			{speaker: "Pink Panther", text: "The countryside, you say? Most helpful. Thank you."},
			{speaker: "Barmaid", text: "Anytime, love. Now, you sure you don't want a pint?"},
		},
		bobAmount: 0.25,
		elevated:  true,
		flipped:   true,
	}
}

var barmaidPostBeerDialog = []dialogEntry{
	{speaker: "Barmaid", text: "Another round? Sorry love, one per customer!"},
	{speaker: "Pink Panther", text: "That's quite alright. Thank you."},
}

func newButler(renderer *sdl.Renderer) *npc {
	tex, w, h := engine.TextureFromPNG(renderer, "assets/images/locations/london/npc/pub_butler.png")
	return &npc{
		tex:     tex,
		srcRect: sdl.Rect{X: 0, Y: 0, W: w, H: h},
		bounds:  sdl.Rect{X: 250, Y: 342, W: 112, H: 292},
		name:    "Butler",
		dialog: []dialogEntry{
			{speaker: "Butler", text: "Good evening, sir. I trust you are enjoying your stay in London."},
			{speaker: "Pink Panther", text: "Why yes, I am. And you are...?"},
			{speaker: "Butler", text: "I am Jackson, sir. Personal butler to Sir Baldley of Devonshire."},
			{speaker: "Butler", text: "Sir Baldley has taken quite an interest in the recent museum incident."},
			{speaker: "Pink Panther", text: "The diamond theft? What does he know about it?"},
			{speaker: "Butler", text: "I couldn't possibly say, sir. But if you were to visit the countryside..."},
			{speaker: "Butler", text: "...you might find Sir Baldley's estate most... illuminating."},
			{speaker: "Pink Panther", text: "How delightfully cryptic. I'll look into it."},
		},
		bobAmount: 0.15,
	}
}

var bobbyDefaultDialog = []dialogEntry{
	{speaker: "Bobby", text: "Evening, sir. Just keeping an eye on things."},
	{speaker: "Pink Panther", text: "Good evening, officer. Busy night?"},
	{speaker: "Bobby", text: "You have no idea. There's been a theft at the museum."},
	{speaker: "Bobby", text: "We've got every officer in London on it."},
	{speaker: "Pink Panther", text: "A theft? What was stolen?"},
	{speaker: "Bobby", text: "The Pink Diamond, would you believe it?"},
	{speaker: "Bobby", text: "Now if you'll excuse me, I'm parched. I could really use a pint..."},
}

var bobbyBeerDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Here you go, officer. A pint on me."},
	{speaker: "Bobby", text: "Well now! Don't mind if I do!"},
	{speaker: "Bobby", text: "*gulp gulp* ... Ahh, that hits the spot!"},
	{speaker: "Bobby", text: "You're alright, you know that? Let me tell you something..."},
	{speaker: "Bobby", text: "The museum... the back entrance has a dodgy lock. Everyone knows it."},
	{speaker: "Bobby", text: "Not that I'm suggesting anything, mind you!"},
}

var bobbyPostBeerDialog = []dialogEntry{
	{speaker: "Bobby", text: "*burp* ... Lovely pint, that was."},
	{speaker: "Pink Panther", text: "Any more tips about the museum?"},
	{speaker: "Bobby", text: "I've said too much already. Now move along, nothing to see here!"},
}

func newPoliceman(renderer *sdl.Renderer) *npc {
	const (
		cols = 6
		rows = 3
	)
	tex, w, h := engine.TextureFromPNGKeyed(renderer, "assets/images/locations/london/npc/policeman_idle.png")
	frameW := w / cols
	frameH := h / rows

	rowRects := func(row int32) []sdl.Rect {
		rects := make([]sdl.Rect, cols)
		for i := int32(0); i < cols; i++ {
			rects[i] = sdl.Rect{X: i * frameW, Y: row * frameH, W: frameW, H: frameH}
		}
		return rects
	}

	idleFrames := rowRects(0)
	talkFrames := rowRects(1)
	drinkFrames := rowRects(2)

	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 1000, Y: 520, W: 160, H: 225},
		name:           "Bobby",
		dialog:         bobbyDefaultDialog,
		bobAmount:      0.10,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		drinkFrames:    drinkFrames,
		frameSpeed:     1.5,
		talkFrameSpeed: 0.12,
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

const (
	npcAnimIdle  = 0
	npcAnimTalk  = 1
	npcAnimDrink = 2
)

func (n *npc) activeFrames() []sdl.Rect {
	switch n.animState {
	case npcAnimTalk:
		if len(n.talkFrames) > 0 {
			return n.talkFrames
		}
	case npcAnimDrink:
		if len(n.drinkFrames) > 0 {
			return n.drinkFrames
		}
	}
	return n.frames
}

func (n *npc) setAnimState(state int) {
	if n.animState == state {
		return
	}
	n.animState = state
	n.frameIdx = 0
	n.frameTimer = 0
	n.animOnce = false
	af := n.activeFrames()
	if len(af) > 0 {
		n.srcRect = af[0]
	}
}

func (n *npc) currentFrameSpeed() float64 {
	if n.animState != npcAnimIdle && n.talkFrameSpeed > 0 {
		return n.talkFrameSpeed
	}
	return n.frameSpeed
}

func (n *npc) update(dt float64) {
	n.bobTimer += dt

	af := n.activeFrames()
	if len(af) > 1 {
		n.frameTimer += dt
		speed := n.currentFrameSpeed()
		if n.frameTimer >= speed {
			n.frameTimer -= speed
			if n.animOnce && n.frameIdx >= len(af)-1 {
				n.setAnimState(npcAnimIdle)
				return
			}
			n.frameIdx = (n.frameIdx + 1) % len(af)
			n.srcRect = af[n.frameIdx]
		}
	}
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

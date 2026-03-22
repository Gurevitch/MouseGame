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

func frameSequence(rects []sdl.Rect, indices ...int) []sdl.Rect {
	seq := make([]sdl.Rect, len(indices))
	for i, idx := range indices {
		seq[i] = rects[idx%len(rects)]
	}
	return seq
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

// ===== Camp Chilly Wa Wa NPCs =====

func campNPCSheet(renderer *sdl.Renderer, path string, cols, rows int32) (*sdl.Texture, int32, int32, func(row int32) []sdl.Rect) {
	tex, w, h := engine.TextureFromPNGRaw(renderer, path)
	frameW := w / cols
	frameH := h / rows
	rowFn := func(row int32) []sdl.Rect {
		rects := make([]sdl.Rect, cols)
		for i := int32(0); i < cols; i++ {
			rects[i] = sdl.Rect{X: i * frameW, Y: row * frameH, W: frameW, H: frameH}
		}
		return rects
	}
	return tex, frameW, frameH, rowFn
}

// --- Director Higgins ---

var higginsDefaultDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Halt! Who goes there?! This is a PRIVATE camp!"},
	{speaker: "Pink Panther", text: "Good afternoon. I'm the substitute counselor."},
	{speaker: "Director Higgins", text: "Substitute?! I was told nothing about this!"},
	{speaker: "Director Higgins", text: "Where is your appointment letter?! No letter, no entry!"},
	{speaker: "Pink Panther", text: "Hmm... I should find that letter."},
}

var higginsLetterDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Here you go. My official appointment letter."},
	{speaker: "Director Higgins", text: "Let me see that... *adjusts glasses*"},
	{speaker: "Director Higgins", text: "Hmm. This seems to be in order. Very well!"},
	{speaker: "Director Higgins", text: "Welcome to Camp Chilly Wa Wa. Try not to break anything."},
	{speaker: "Director Higgins", text: "The camp grounds are through the gate. The kids are... around."},
}

var higginsPostLetterDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Don't just stand there! The kids need supervision!"},
	{speaker: "Pink Panther", text: "Right. Supervision. My specialty."},
}

func newDirectorHiggins(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_director_higgins.png", 7, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 1, 0, 0, 0, 2, 2, 0, 0, 1, 1, 0, 0)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 700, Y: 360, W: 120, H: 210},
		name:           "Director Higgins",
		dialog:         higginsDefaultDialog,
		bobAmount:      0.2,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.48,
		talkFrameSpeed: 0.10,
	}
}

// --- Tommy (Homesick Kid) ---

var tommyDialog = []dialogEntry{
	{speaker: "Tommy", text: "*sniff* I... I want to go home..."},
	{speaker: "Pink Panther", text: "Hey there, little one. What's wrong?"},
	{speaker: "Tommy", text: "Everything! The food is gross, the bugs are huge, and I miss my dog!"},
	{speaker: "Pink Panther", text: "Camp can be tough. Is there anything that cheers you up?"},
	{speaker: "Tommy", text: "My pen pal lives in Scotland... she says there's a monster in the loch, and someone is trying to catch it!"},
	{speaker: "Tommy", text: "I wish I could visit her instead of being stuck here..."},
	{speaker: "Tommy", text: "If only I had something to read... a comic book or something..."},
}

var tommyComicDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Hey Tommy, look what I found! A comic book."},
	{speaker: "Tommy", text: "R-really?! For me?!"},
	{speaker: "Tommy", text: "Oh wow! This is THE BEST! Thank you so much!"},
	{speaker: "Tommy", text: "You know what... there's a shortcut past Jake."},
	{speaker: "Tommy", text: "Go behind the big cabin and you can get to the lake without him seeing you!"},
}

var tommyHappyDialog = []dialogEntry{
	{speaker: "Tommy", text: "This comic is amazing! Thank you, Pink Panther!"},
	{speaker: "Pink Panther", text: "Anytime, kid. Hang in there."},
}

func newTommy(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_homesick_kid.png", 7, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 1, 0, 0, 0, 0, 2, 2, 0, 0, 1, 0, 0, 0)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 200, Y: 420, W: 120, H: 140},
		name:           "Tommy",
		dialog:         tommyDialog,
		bobAmount:      0.3,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.50,
		talkFrameSpeed: 0.10,
	}
}

// --- Jake (Bully Kid) ---

var jakeDialog = []dialogEntry{
	{speaker: "Jake", text: "Hey! Where do you think YOU'RE going?!"},
	{speaker: "Pink Panther", text: "Just passing through, my young friend."},
	{speaker: "Jake", text: "No way! Nobody gets past ME without paying the toll!"},
	{speaker: "Pink Panther", text: "A toll? What kind of toll?"},
	{speaker: "Jake", text: "Food! I'm STARVING! Get me something good and maybe I'll let you pass."},
	{speaker: "Jake", text: "My dad brought me a weird old coin from Israel. He said it came from some hidden tunnel under an ancient city."},
	{speaker: "Jake", text: "Maybe I'll trade it if you bring me something REALLY good to eat!"},
}

var jakeFedDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "Here, try this."},
	{speaker: "Jake", text: "*munch munch* ... Hey, this is actually good!"},
	{speaker: "Jake", text: "Alright, alright. You can pass. Just don't tell anyone I was nice."},
}

var jakePostFedDialog = []dialogEntry{
	{speaker: "Jake", text: "What? I already let you pass! Scram!"},
	{speaker: "Pink Panther", text: "Charming as always."},
}

func newJake(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_bully_kid.png", 6, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 1, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 700, Y: 400, W: 130, H: 180},
		name:           "Jake",
		dialog:         jakeDialog,
		bobAmount:      0.25,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.45,
		talkFrameSpeed: 0.10,
	}
}

// --- Lily (Shy Girl) ---

var lilyDialog = []dialogEntry{
	{speaker: "Lily", text: "..."},
	{speaker: "Pink Panther", text: "Hello there. Mind if I sit for a moment?"},
	{speaker: "Lily", text: "...o-okay..."},
}

var lilySecondDialog = []dialogEntry{
	{speaker: "Pink Panther", text: "It's a nice campfire, isn't it?"},
	{speaker: "Lily", text: "...I guess so..."},
	{speaker: "Lily", text: "I... I saw a postcard from Japan once... a temple in the mountains with a garden that glows at night..."},
	{speaker: "Lily", text: "I wish I could see it someday..."},
	{speaker: "Pink Panther", text: "A glowing garden? That sounds magical."},
	{speaker: "Lily", text: "*small smile* ...you're... nice..."},
}

func newLily(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_shy_girl.png", 6, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 1, 0)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 550, Y: 440, W: 110, H: 130},
		name:           "Lily",
		dialog:         lilyDialog,
		bobAmount:      0.15,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.55,
		talkFrameSpeed: 0.12,
	}
}

// --- Marcus (Know-It-All) ---

var marcusDialog = []dialogEntry{
	{speaker: "Marcus", text: "Ah, a new counselor! Did you know this camp was founded in 1968?"},
	{speaker: "Pink Panther", text: "I did not. How enlightening."},
	{speaker: "Marcus", text: "Actually, there's a famous art heist case in Paris. The painting was never found."},
	{speaker: "Marcus", text: "I read the whole file! 347 pages! The thief used the catacombs as an escape route!"},
	{speaker: "Pink Panther", text: "A stolen painting in Paris... now THAT is interesting."},
	{speaker: "Marcus", text: "I can tell you more! I've catalogued every unsolved case in my notebook!"},
}

func newMarcus(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_know_it_all.png", 6, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 1, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 2)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 1000, Y: 380, W: 100, H: 200},
		name:           "Marcus",
		dialog:         marcusDialog,
		bobAmount:      0.2,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.46,
		talkFrameSpeed: 0.10,
	}
}

// --- Danny (Prankster) ---

var dannyDialog = []dialogEntry{
	{speaker: "Danny", text: "Psst! Hey! Over here!"},
	{speaker: "Pink Panther", text: "Hmm? What are you doing back there?"},
	{speaker: "Danny", text: "Shh! I'm setting up the ULTIMATE prank! Wanna help?"},
	{speaker: "Pink Panther", text: "I'll pass, thank you."},
	{speaker: "Danny", text: "Your loss! Dude, my cousin snuck into some catacombs in Italy."},
	{speaker: "Danny", text: "Said there's a secret room with gold everywhere! Can you believe it?!"},
	{speaker: "Pink Panther", text: "Catacombs... gold... Italy. Noted."},
}

func newDanny(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_prankster.png", 6, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 1, 0, 0, 2, 0, 0, 0, 1, 1, 0, 0, 2, 0)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 1150, Y: 400, W: 110, H: 180},
		name:           "Danny",
		dialog:         dannyDialog,
		bobAmount:      0.3,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.42,
		talkFrameSpeed: 0.10,
	}
}

// --- Cook Marge ---

var cookMargeDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "Well hello there, sugar! Welcome to my kitchen!"},
	{speaker: "Pink Panther", text: "Good day, madam. Something smells... interesting."},
	{speaker: "Cook Marge", text: "That's my famous mystery stew! Secret recipe!"},
	{speaker: "Cook Marge", text: "Say, you look like a helpful sort. Mind giving me a hand?"},
	{speaker: "Cook Marge", text: "I need someone to stir the pot while I find the salt."},
	{speaker: "Pink Panther", text: "I suppose I could do that..."},
}

var cookMargeHelpedDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "Thank you, dear! Here, take some stew for the road."},
	{speaker: "Cook Marge", text: "It might not look pretty, but it's got kick!"},
	{speaker: "Pink Panther", text: "How... appetizing. Thank you."},
}

var cookMargePostHelpDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "Come back anytime, sugar! The kitchen is always open!"},
	{speaker: "Pink Panther", text: "I'll keep that in mind."},
}

func newCookMarge(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_cook_marge.png", 6, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 1, 0, 0, 2, 0, 0, 0, 1, 1, 0, 0, 0)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 600, Y: 300, W: 220, H: 280},
		name:           "Cook Marge",
		dialog:         cookMargeDialog,
		bobAmount:      0.4,
		elevated:       true,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.48,
		talkFrameSpeed: 0.10,
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

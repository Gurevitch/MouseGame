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

	// Strange state (Day 2) — alternate texture + frames
	strangeTex        *sdl.Texture
	strangeFrames     []sdl.Rect
	strangeTalkFrames []sdl.Rect
	normalTex         *sdl.Texture // saved normal tex for restoring
	normalFrames      []sdl.Rect
	normalTalkFrames  []sdl.Rect
	isStrange         bool
}

func (n *npc) setStrange(strange bool) {
	if strange == n.isStrange {
		return
	}
	n.isStrange = strange
	if strange && n.strangeTex != nil {
		// Save normal state
		n.normalTex = n.tex
		n.normalFrames = n.frames
		n.normalTalkFrames = n.talkFrames
		// Swap to strange
		n.tex = n.strangeTex
		n.frames = n.strangeFrames
		n.talkFrames = n.strangeTalkFrames
	} else if !strange && n.normalTex != nil {
		// Restore normal
		n.tex = n.normalTex
		n.frames = n.normalFrames
		n.talkFrames = n.normalTalkFrames
	}
	n.frameIdx = 0
	n.frameTimer = 0
	n.setAnimState(npcAnimIdle)
}

func frameSequence(rects []sdl.Rect, indices ...int) []sdl.Rect {
	seq := make([]sdl.Rect, len(indices))
	for i, idx := range indices {
		seq[i] = rects[idx%len(rects)]
	}
	return seq
}

// loadStrangeSheet loads a strange-state sprite sheet and sets it on the NPC
func loadStrangeSheet(renderer *sdl.Renderer, n *npc, path string) {
	tex, w, h := engine.SafeTextureFromPNGKeyed(renderer, path)
	if tex == nil {
		return
	}
	cols := int32(8)
	frameW := w / cols
	frameH := h / 2

	rowFn := func(row int32) []sdl.Rect {
		rects := make([]sdl.Rect, cols)
		for i := int32(0); i < cols; i++ {
			rects[i] = sdl.Rect{X: i * frameW, Y: row * frameH, W: frameW, H: frameH}
		}
		return rects
	}

	idleRaw := rowFn(0)
	n.strangeTex = tex
	n.strangeFrames = frameSequence(idleRaw, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7)
	n.strangeTalkFrames = rowFn(1)
}

// ===== Camp Chilly Wa Wa NPCs =====

func campNPCSheet(renderer *sdl.Renderer, path string, cols, rows int32) (*sdl.Texture, int32, int32, func(row int32) []sdl.Rect) {
	tex, w, h := engine.TextureFromPNGKeyed(renderer, path)
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
	{speaker: "Director Higgins", text: "Ah, you must be the new counselor! Finally!"},
	{speaker: "Pink Panther", text: "Good afternoon. Pink Panther, at your service."},
	{speaker: "Director Higgins", text: "Yes, yes. Welcome to Camp Chilly Wa Wa."},
	{speaker: "Director Higgins", text: "The kids are through the gate. Go introduce yourself."},
	{speaker: "Director Higgins", text: "They're a good bunch. A little... eccentric, but good."},
}

var higginsPostDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Go on, the kids are waiting in the camp grounds!"},
	{speaker: "Pink Panther", text: "On my way."},
}

var higginsWorriedDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Something is wrong with the kids."},
	{speaker: "Director Higgins", text: "Marcus has been up all night drawing things he's never seen."},
	{speaker: "Director Higgins", text: "Buildings, paintings, rooftops... from places he's never been!"},
	{speaker: "Pink Panther", text: "That does sound unusual."},
	{speaker: "Director Higgins", text: "I've seen this kind of thing before... well, no I haven't. But it's NOT normal!"},
	{speaker: "Director Higgins", text: "If you think the answer is out there somewhere..."},
	{speaker: "Director Higgins", text: "Camp Chilly Wa Wa Air hasn't flown in years, but I can get the old plane running."},
	{speaker: "Pink Panther", text: "A camp... airline?"},
	{speaker: "Director Higgins", text: "Don't ask questions. Just check on Marcus first."},
}

var higginsPostWorriedDialog = []dialogEntry{
	{speaker: "Director Higgins", text: "Have you talked to Marcus yet? He's in the camp grounds."},
	{speaker: "Director Higgins", text: "And the other kids might have noticed something too."},
}

func newDirectorHiggins(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_director_higgins.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 1, 1, 0, 0, 0, 2, 2, 0, 0, 0, 3, 0)
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
		frameSpeed:     0.7,
		talkFrameSpeed: 0.12,
	}
}

// --- Tommy (Homesick Kid) ---

var tommyDialog = []dialogEntry{
	{speaker: "Tommy", text: "Hi! I'm Tommy! Are you the new counselor?"},
	{speaker: "Pink Panther", text: "That's right. Nice to meet you, Tommy."},
	{speaker: "Tommy", text: "I love telling stories! Did you know there's a legend about a treasure at this camp?"},
	{speaker: "Tommy", text: "Probably not true though... I like making things sound more exciting than they are!"},
	{speaker: "Pink Panther", text: "A natural storyteller. I like that."},
}

var tommyPostDialog = []dialogEntry{
	{speaker: "Tommy", text: "Want to hear another story? I've got HUNDREDS!"},
	{speaker: "Pink Panther", text: "Maybe later, Tommy."},
}

var tommyStrangeDialog = []dialogEntry{
	{speaker: "Tommy", text: "Do you hear that? The music?"},
	{speaker: "Pink Panther", text: "Music? I don't hear anything."},
	{speaker: "Tommy", text: "It's drums and singing... and there's a GIANT STATUE watching over everyone!"},
	{speaker: "Tommy", text: "People are dancing in the streets! It's like the biggest party in the world!"},
	{speaker: "Tommy", text: "And then... tango? Somewhere else, a different city, a wide road with a tall white tower..."},
	{speaker: "Pink Panther", text: "Tommy, are you alright? You've never been to any of these places."},
	{speaker: "Tommy", text: "I KNOW! That's what's so weird! But I can SEE it!"},
}

var tommyPostStrangeDialog = []dialogEntry{
	{speaker: "Tommy", text: "The music won't stop... a giant statue, parades, dancing..."},
	{speaker: "Tommy", text: "It feels like two places at once. I can't explain it."},
}

func newTommy(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_homesick_kid.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0, 3)
	talkFrames := rowFn(1)
	n := &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 200, Y: 420, W: 120, H: 140},
		name:           "Tommy",
		dialog:         tommyDialog,
		bobAmount:      0.3,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.7,
		talkFrameSpeed: 0.12,
	}
	loadStrangeSheet(renderer, n, "assets/images/locations/camp/npc/npc_homesick_kid_strange.png")
	return n
}

// --- Jake (Bully Kid) ---

var jakeDialog = []dialogEntry{
	{speaker: "Jake", text: "Hey! You the new guy?"},
	{speaker: "Pink Panther", text: "That's right. And you are?"},
	{speaker: "Jake", text: "Jake. I'm the toughest kid at camp. Don't forget it."},
	{speaker: "Jake", text: "I collect stuff. Rocks, coins, anything shiny. Check out this coin my dad brought from Israel."},
	{speaker: "Pink Panther", text: "That's a beautiful coin. Where exactly is it from?"},
	{speaker: "Jake", text: "Some old city with tunnels underneath. Jerusalem, I think. Dad said the tunnels are ANCIENT."},
	{speaker: "Pink Panther", text: "Fascinating collection you've got there."},
}

var jakePostDialog = []dialogEntry{
	{speaker: "Jake", text: "Don't touch my collection. I'm watching you."},
	{speaker: "Pink Panther", text: "Wouldn't dream of it."},
}

var jakeStrangeDialog = []dialogEntry{
	{speaker: "Jake", text: "Something's happening to my coins..."},
	{speaker: "Pink Panther", text: "What do you mean?"},
	{speaker: "Jake", text: "I keep hearing echoes. Like tunnels underground. Voices bouncing off stone walls."},
	{speaker: "Jake", text: "And I can't stop rubbing every metal surface for symbols. Look at this bench - I KNOW there's something underneath!"},
	{speaker: "Pink Panther", text: "Jake, that's just a wooden bench."},
	{speaker: "Jake", text: "NO! There are tunnels! Old ones! Under an ancient city! I can FEEL them!"},
}

var jakePostStrangeDialog = []dialogEntry{
	{speaker: "Jake", text: "The echoes won't stop... tunnels under old stone walls..."},
	{speaker: "Jake", text: "It's like I can see a huge wall... and something hidden behind it."},
}

func newJake(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_bully_kid.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 1, 1, 0, 0, 0, 2, 0, 0, 0, 1, 0, 0)
	talkFrames := rowFn(1)
	n := &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 700, Y: 400, W: 130, H: 180},
		name:           "Jake",
		dialog:         jakeDialog,
		bobAmount:      0.25,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.7,
		talkFrameSpeed: 0.12,
	}
	loadStrangeSheet(renderer, n, "assets/images/locations/camp/npc/npc_bully_kid_strange.png")
	return n
}

// --- Lily (Shy Girl) ---

var lilyDialog = []dialogEntry{
	{speaker: "Lily", text: "..."},
	{speaker: "Pink Panther", text: "Hello there. I'm the new counselor."},
	{speaker: "Lily", text: "...o-okay..."},
	{speaker: "Pink Panther", text: "Those are beautiful flowers you're arranging."},
	{speaker: "Lily", text: "...thank you... I like flowers... and quiet places..."},
}

var lilyPostDialog = []dialogEntry{
	{speaker: "Lily", text: "...hi again..."},
	{speaker: "Pink Panther", text: "Hello, Lily. Beautiful day, isn't it?"},
	{speaker: "Lily", text: "*small nod*"},
}

var lilyStrangeDialog = []dialogEntry{
	{speaker: "Lily", text: "...the flowers are glowing..."},
	{speaker: "Pink Panther", text: "Glowing? They look normal to me."},
	{speaker: "Lily", text: "Not these flowers... OTHER flowers. In a garden far away..."},
	{speaker: "Lily", text: "I keep arranging petals into shapes... symbols I don't understand..."},
	{speaker: "Lily", text: "And I hear bells. Temple bells. Very old ones."},
	{speaker: "Lily", text: "There's a red gate... and cherry blossoms falling everywhere..."},
	{speaker: "Pink Panther", text: "That sounds like Japan, Lily. Have you ever been there?"},
	{speaker: "Lily", text: "...never... but I can see it when I close my eyes..."},
}

var lilyPostStrangeDialog = []dialogEntry{
	{speaker: "Lily", text: "...the bells again... and glowing petals..."},
	{speaker: "Lily", text: "...a temple in the mountains... I can almost touch it..."},
}

func newLily(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_shy_girl.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 1, 0)
	talkFrames := rowFn(1)
	n := &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 550, Y: 440, W: 110, H: 130},
		name:           "Lily",
		dialog:         lilyDialog,
		bobAmount:      0.15,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.8,
		talkFrameSpeed: 0.12,
	}
	loadStrangeSheet(renderer, n, "assets/images/locations/camp/npc/npc_shy_girl_strange.png")
	return n
}

// --- Marcus (Know-It-All) ---

var marcusDialog = []dialogEntry{
	{speaker: "Marcus", text: "Ah, a new counselor! Did you know this camp was founded in 1968?"},
	{speaker: "Pink Panther", text: "I did not. And you are?"},
	{speaker: "Marcus", text: "Marcus. I know everything about everything. Ask me anything!"},
	{speaker: "Pink Panther", text: "Alright. What's the most interesting thing about this camp?"},
	{speaker: "Marcus", text: "Statistically, the mess hall food has a 73 percent chance of being inedible."},
	{speaker: "Marcus", text: "But I also love drawing! Want to see my sketches? I drew the whole campfire!"},
	{speaker: "Pink Panther", text: "Very impressive. You've got talent, Marcus."},
}

var marcusPostDialog = []dialogEntry{
	{speaker: "Marcus", text: "Did you know butterflies taste with their feet? It's TRUE!"},
	{speaker: "Pink Panther", text: "You never stop, do you?"},
}

var marcusStrangeDialog = []dialogEntry{
	{speaker: "Marcus", text: "It's WRONG! The picture is WRONG!"},
	{speaker: "Pink Panther", text: "Marcus? What's going on?"},
	{speaker: "Marcus", text: "I keep drawing this woman's face... but I've NEVER seen her before!"},
	{speaker: "Marcus", text: "And these frames... ornate golden frames... and rooftops I've never visited!"},
	{speaker: "Marcus", text: "It's a museum. A HUGE museum. The biggest in the world!"},
	{speaker: "Marcus", text: "There's a glass pyramid in front of it... and inside, a painting that everyone stares at..."},
	{speaker: "Marcus", text: "But something is MISSING from the picture! I can feel it!"},
	{speaker: "Pink Panther", text: "A glass pyramid... the biggest museum... That sounds like the Louvre in Paris."},
	{speaker: "Marcus", text: "I've never been to Paris! But I can't stop drawing it!"},
}

var marcusPostStrangeDialog = []dialogEntry{
	{speaker: "Marcus", text: "The woman's face again... the golden frames... something is missing..."},
	{speaker: "Marcus", text: "I filled twelve pages last night. I can't stop."},
}

func newMarcus(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_know_it_all.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7)
	talkFrames := rowFn(1)
	n := &npc{
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
	loadStrangeSheet(renderer, n, "assets/images/locations/camp/npc/npc_know_it_all_strange.png")
	return n
}

// --- Danny (Prankster) ---

var dannyDialog = []dialogEntry{
	{speaker: "Danny", text: "Psst! Hey! Over here!"},
	{speaker: "Pink Panther", text: "Hmm? What are you doing behind that tree?"},
	{speaker: "Danny", text: "Setting up the ULTIMATE prank! I'm Danny, master of mischief!"},
	{speaker: "Danny", text: "I love treasure stories. My cousin went to Italy once and saw REAL ancient ruins!"},
	{speaker: "Danny", text: "The Colosseum! Gladiators fought there! How cool is that?!"},
	{speaker: "Pink Panther", text: "Very cool, Danny. Try not to prank anyone too badly."},
}

var dannyPostDialog = []dialogEntry{
	{speaker: "Danny", text: "Psst! Want to help me put a frog in Higgins' coffee?"},
	{speaker: "Pink Panther", text: "I'll pretend I didn't hear that."},
}

var dannyStrangeDialog = []dialogEntry{
	{speaker: "Danny", text: "Dude! DUDE! There's treasure EVERYWHERE!"},
	{speaker: "Pink Panther", text: "Danny, calm down. What are you talking about?"},
	{speaker: "Danny", text: "I've been mapping the whole camp! It's just like ancient ruins!"},
	{speaker: "Danny", text: "There are gold paths under the ground... I can FEEL them!"},
	{speaker: "Danny", text: "A huge round arena... with arches... thousands of people cheering..."},
	{speaker: "Danny", text: "And tunnels underneath with hidden rooms full of treasure!"},
	{speaker: "Pink Panther", text: "An arena with arches... that sounds like the Colosseum in Rome."},
	{speaker: "Danny", text: "I've never been to Rome! But I drew a map of it! Look!"},
}

var dannyPostStrangeDialog = []dialogEntry{
	{speaker: "Danny", text: "The gold paths are getting clearer... arches and tunnels everywhere..."},
	{speaker: "Danny", text: "I've dug three holes behind the cabin already. Higgins is NOT happy."},
}

func newDanny(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_prankster.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7)
	talkFrames := rowFn(1)
	n := &npc{
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
	loadStrangeSheet(renderer, n, "assets/images/locations/camp/npc/npc_prankster_strange.png")
	return n
}

// --- Cook Marge ---

var cookMargeDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "Well hello there, sugar! You must be the new counselor!"},
	{speaker: "Pink Panther", text: "Good day, madam. Something smells... interesting."},
	{speaker: "Cook Marge", text: "That's my famous mystery stew! Been cooking it since '68!"},
	{speaker: "Cook Marge", text: "I've seen a LOT of strange things at this camp over the years."},
	{speaker: "Pink Panther", text: "Strange things? Like what?"},
	{speaker: "Cook Marge", text: "Oh, kids acting funny. Drawing things they shouldn't know about."},
	{speaker: "Cook Marge", text: "It happened before, years ago. The old counselor figured it out eventually."},
	{speaker: "Cook Marge", text: "But that's a story for another time. Want some stew?"},
}

var cookMargePostDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "Come back anytime, sugar! The kitchen is always open!"},
	{speaker: "Pink Panther", text: "Thank you, Marge."},
}

var cookMargeWorriedDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "It's happening again, isn't it?"},
	{speaker: "Pink Panther", text: "You mentioned this happened before..."},
	{speaker: "Cook Marge", text: "Years ago, a group of kids started acting the same way."},
	{speaker: "Cook Marge", text: "Seeing places they'd never been, drawing things they shouldn't know."},
	{speaker: "Cook Marge", text: "The old counselor traveled to the places the kids described."},
	{speaker: "Cook Marge", text: "Brought back little souvenirs from each place, and the kids settled right down."},
	{speaker: "Pink Panther", text: "So the answer is out there... in the real places they're seeing."},
	{speaker: "Cook Marge", text: "That's what I believe, sugar. Find what they're seeing, and bring a piece of it back."},
}

var cookMargePostWorriedDialog = []dialogEntry{
	{speaker: "Cook Marge", text: "Find what the kids are seeing, sugar. Bring back a piece of it."},
	{speaker: "Cook Marge", text: "That's the only thing that worked last time."},
}

func newCookMarge(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/camp/npc/npc_cook_marge.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7)
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

// ===== Paris NPCs =====

var frenchGuideDialog = []dialogEntry{
	{speaker: "Madame Colette", text: "Bonjour, monsieur! Welcome to Paris!"},
	{speaker: "Pink Panther", text: "Bonjour, madame. I'm looking for information about the Louvre."},
	{speaker: "Madame Colette", text: "Ah, ze Louvre! Ze largest art museum in ze world!"},
	{speaker: "Madame Colette", text: "It was originally a royal palace, built in ze 12th century."},
	{speaker: "Madame Colette", text: "Today it holds over 380,000 objects and 35,000 works of art!"},
	{speaker: "Pink Panther", text: "Impressive. And what about that glass pyramid?"},
	{speaker: "Madame Colette", text: "Ah, ze Pyramid! Designed by I.M. Pei in 1989. Very controversial at first!"},
	{speaker: "Madame Colette", text: "People said it did not belong. Now it is ze most famous entrance in ze world."},
	{speaker: "Madame Colette", text: "And of course, ze Eiffel Tower behind you — built in 1889 for ze World Fair."},
	{speaker: "Madame Colette", text: "Gustave Eiffel designed it. It was meant to be temporary — just 20 years!"},
	{speaker: "Madame Colette", text: "But zey kept it because it was perfect for radio transmissions."},
	{speaker: "Pink Panther", text: "A temporary tower that became permanent. How fitting."},
	{speaker: "Madame Colette", text: "Ze museum is just down ze street, to ze right. Enjoy, monsieur!"},
}

var frenchGuidePostDialog = []dialogEntry{
	{speaker: "Madame Colette", text: "Ze Louvre is to ze right, monsieur. You cannot miss ze pyramid!"},
	{speaker: "Pink Panther", text: "Merci, madame."},
}

func newFrenchGuide(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/paris/npc/npc_french_guide.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 300, Y: 350, W: 140, H: 240},
		name:           "Madame Colette",
		dialog:         frenchGuideDialog,
		bobAmount:      0.2,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.50,
		talkFrameSpeed: 0.10,
	}
}

var museumCuratorDialog = []dialogEntry{
	{speaker: "Curator Beaumont", text: "Ah, a visitor! Welcome to ze Musee du Louvre."},
	{speaker: "Pink Panther", text: "Thank you. I'm investigating something... unusual."},
	{speaker: "Curator Beaumont", text: "Unusual? In ze Louvre, everything is extraordinary!"},
	{speaker: "Curator Beaumont", text: "Zis hall contains some of ze finest works in history."},
	{speaker: "Curator Beaumont", text: "Ze Mona Lisa, of course — painted by Leonardo da Vinci around 1503."},
	{speaker: "Curator Beaumont", text: "Her smile has puzzled visitors for over 500 years!"},
	{speaker: "Curator Beaumont", text: "And ze Venus de Milo — a Greek sculpture from around 100 BC."},
	{speaker: "Pink Panther", text: "Actually, I'm looking for a specific painting. A boy back at camp keeps drawing it."},
	{speaker: "Curator Beaumont", text: "A boy... drawing paintings he has never seen? How peculiar."},
	{speaker: "Curator Beaumont", text: "Describe what he draws, and perhaps I can identify it."},
	{speaker: "Pink Panther", text: "A woman's face. Ornate golden frames. He says something is 'missing' from it."},
	{speaker: "Curator Beaumont", text: "Mon Dieu... zat sounds like ze portrait in Room 7."},
	{speaker: "Curator Beaumont", text: "A painting zat was recently restored. Ze restorer found a hidden symbol underneath."},
	{speaker: "Curator Beaumont", text: "Perhaps your boy senses what was hidden. Take zis postcard of ze painting."},
	{speaker: "Curator Beaumont", text: "If he sees ze complete image, perhaps his mind will settle."},
}

var museumCuratorPostDialog = []dialogEntry{
	{speaker: "Curator Beaumont", text: "Ze postcard should help your young friend."},
	{speaker: "Curator Beaumont", text: "Ze mysteries of art connect us in ways we do not understand."},
}

func newMuseumCurator(renderer *sdl.Renderer) *npc {
	tex, _, _, rowFn := campNPCSheet(renderer, "assets/images/locations/paris/npc/npc_museum_curator.png", 8, 2)
	raw := rowFn(0)
	idleFrames := frameSequence(raw, 0, 0, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7)
	talkFrames := rowFn(1)
	return &npc{
		tex:            tex,
		srcRect:        idleFrames[0],
		bounds:         sdl.Rect{X: 500, Y: 320, W: 130, H: 250},
		name:           "Curator Beaumont",
		dialog:         museumCuratorDialog,
		bobAmount:      0.15,
		frames:         idleFrames,
		talkFrames:     talkFrames,
		frameSpeed:     0.50,
		talkFrameSpeed: 0.10,
	}
}

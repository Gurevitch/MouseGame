package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type spriteFrame struct {
	tex *sdl.Texture
	w   int32
	h   int32
}

const (
	playerBaseSpeed = 250.0
	playerDstW      = 140
	playerDstH      = 200
	playerMinX      = 10.0
	playerMaxX      = engine.ScreenWidth - playerDstW - 10.0
	playerMinY      = 340.0
	playerMaxY      = engine.ScreenHeight - playerDstH - 60.0
	walkFrameTime   = 0.12
)

type playerState int

const (
	stateIdle playerState = iota
	stateWalking
	stateTalking
)

// Ping-pong walk cycle: frame 0 -> 1 -> 2 -> 1 -> repeat
var walkCycle = [4]int{0, 1, 2, 1}

type player struct {
	idleFrame  spriteFrame
	walkFrames [3]spriteFrame
	talkFrame  spriteFrame
	scale      float64

	x, y       float64
	targetX    float64
	targetY    float64
	moving     bool
	facingLeft bool
	state      playerState

	breathTimer  float64
	walkCycleIdx int
	walkTimer    float64
	talkTimer    float64
	talkOpen     bool

	interactTarget *npc
	dialogSys      *dialogSystem
	onArrival      func()
}

func newPlayer(renderer *sdl.Renderer) *player {
	p := &player{
		x: 200,
		y: float64(engine.ScreenHeight) - playerDstH - 160,
	}

	idleTex, idleW, idleH := engine.TextureFromPNG(renderer, "assets/images/pp_idle.png")
	p.idleFrame = spriteFrame{tex: idleTex, w: idleW, h: idleH}

	walkTexs, walkWs, walkHs := engine.SpriteFramesFromPNG(renderer, "assets/images/pp_walk_sheet.png", 3)
	for i := 0; i < 3; i++ {
		p.walkFrames[i] = spriteFrame{tex: walkTexs[i], w: walkWs[i], h: walkHs[i]}
	}

	talkTex, talkW, talkH := engine.TextureFromPNG(renderer, "assets/images/pp_talk.png")
	p.talkFrame = spriteFrame{tex: talkTex, w: talkW, h: talkH}

	// Uniform scale from the tallest frame so all poses match in size
	maxH := p.idleFrame.h
	if p.talkFrame.h > maxH {
		maxH = p.talkFrame.h
	}
	for _, wf := range p.walkFrames {
		if wf.h > maxH {
			maxH = wf.h
		}
	}
	p.scale = float64(playerDstH) / float64(maxH)

	return p
}

func (p *player) currentSprite() spriteFrame {
	switch p.state {
	case stateWalking:
		return p.walkFrames[walkCycle[p.walkCycleIdx]]
	case stateTalking:
		if p.talkOpen {
			return p.talkFrame
		}
		return p.idleFrame
	default:
		return p.idleFrame
	}
}

func (p *player) setTarget(x, y float64) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(y-playerDstH/2, playerMinY, playerMaxY)
	p.moving = true
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = nil
}

func (p *player) walkToAndInteract(target *npc, ds *dialogSystem) {
	npcCenter := float64(target.bounds.X + target.bounds.W/2)
	var tx float64
	if npcCenter < engine.ScreenWidth/2 {
		tx = float64(target.bounds.X+target.bounds.W) + 20
	} else {
		tx = float64(target.bounds.X) - playerDstW - 20
	}
	ty := float64(target.bounds.Y + target.bounds.H - playerDstH)
	p.targetX = engine.Clamp(tx, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(ty, playerMinY, playerMaxY)
	p.moving = true
	p.state = stateWalking
	p.interactTarget = target
	p.dialogSys = ds
	p.onArrival = nil
}

func (p *player) walkToAndDo(x, y float64, action func()) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(y-playerDstH/2, playerMinY, playerMaxY)
	p.moving = true
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = action
}

func (p *player) update(dt float64) {
	p.breathTimer += dt

	if p.moving {
		p.walkTimer += dt
		if p.walkTimer >= walkFrameTime {
			p.walkTimer -= walkFrameTime
			p.walkCycleIdx = (p.walkCycleIdx + 1) % len(walkCycle)
		}
	} else {
		p.walkCycleIdx = 0
		p.walkTimer = 0
	}

	if p.state == stateTalking {
		p.talkTimer += dt
		if p.talkTimer >= 0.30 {
			p.talkTimer -= 0.30
			p.talkOpen = !p.talkOpen
		}
	} else {
		p.talkTimer = 0
		p.talkOpen = false
	}

	if !p.moving {
		if p.state != stateTalking {
			p.state = stateIdle
		}
		return
	}

	dx := p.targetX - p.x
	dy := p.targetY - p.y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 3.0 {
		p.x = p.targetX
		p.y = p.targetY
		p.moving = false
		p.state = stateIdle
		if p.interactTarget != nil && p.dialogSys != nil {
			p.state = stateTalking
			p.dialogSys.startDialog(p.interactTarget.dialog)
			p.facingLeft = p.x > float64(p.interactTarget.bounds.X)
			p.interactTarget = nil
		}
		if p.onArrival != nil {
			fn := p.onArrival
			p.onArrival = nil
			fn()
		}
		return
	}

	speed := playerBaseSpeed
	if dist < 100 {
		speed = playerBaseSpeed * (0.3 + 0.7*dist/100.0)
	}
	p.x = engine.Clamp(p.x+(dx/dist)*speed*dt, playerMinX, playerMaxX)
	p.y = engine.Clamp(p.y+(dy/dist)*speed*dt, playerMinY, playerMaxY)
	p.facingLeft = dx < 0
}

func (p *player) draw(renderer *sdl.Renderer) {
	frame := p.currentSprite()

	dstW := int32(float64(frame.w) * p.scale)
	dstH := int32(float64(frame.h) * p.scale)

	// Center horizontally within the logical bounding box, bottom-align
	dstX := int32(p.x) + (playerDstW-dstW)/2
	dstY := int32(p.y) + (playerDstH - dstH)

	flip := sdl.FLIP_NONE
	if p.facingLeft {
		flip = sdl.FLIP_HORIZONTAL
	}

	switch p.state {
	case stateIdle:
		breathVal := math.Sin(p.breathTimer * 2.0)
		dstY += int32(breathVal * 2.0)
	case stateTalking:
		bob := math.Sin(p.breathTimer*3.0) * 1.5
		dstY += int32(bob)
	}

	renderer.CopyEx(frame.tex, nil,
		&sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH},
		0, nil, flip)
}

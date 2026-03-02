package game

import (
	"math"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	playerBaseSpeed = 250.0
	playerDstW      = 105
	playerDstH      = 150
	playerMinX      = 10.0
	playerMaxX      = engine.ScreenWidth - playerDstW - 10.0
	playerMinY      = 380.0
	playerMaxY      = engine.ScreenHeight - playerDstH - 80.0
)

type playerState int

const (
	stateIdle playerState = iota
	stateWalking
	stateTalking
)

var walkRotations = [4]float64{-3.0, 0, 3.0, 0}
var walkYOffsets = [4]float64{-3, 0, -3, 0}

type player struct {
	tex            *sdl.Texture
	x, y           float64
	targetX        float64
	targetY        float64
	moving         bool
	facingLeft     bool
	state          playerState
	srcW, srcH     int32
	anim           *engine.Animator
	breathTimer    float64
	walkFrame      int
	walkTimer      float64
	talkTimer      float64
	talkOpen       bool
	interactTarget *npc
	dialogSys      *dialogSystem
	onArrival      func()
}

func newPlayer(renderer *sdl.Renderer) *player {
	p := &player{
		tex:  engine.TextureFromBMP(renderer, "assets/images/player.bmp"),
		x:    200,
		y:    float64(engine.ScreenHeight) - playerDstH - 160,
		srcW: 242,
		srcH: 582,
	}
	p.anim = engine.NewAnimator()
	p.anim.AddAnimation("idle", []sdl.Rect{{X: 0, Y: 0, W: 242, H: 582}}, 1.0)
	p.anim.AddAnimation("walk", []sdl.Rect{{X: 0, Y: 0, W: 242, H: 582}}, 0.15)
	p.anim.Play("idle")
	return p
}

func (p *player) setTarget(x, y float64) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(y-playerDstH/2, playerMinY, playerMaxY)
	p.moving = true
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = nil
	p.anim.Play("walk")
}

func (p *player) walkToAndInteract(target *npc, ds *dialogSystem) {
	tx := float64(target.bounds.X) - playerDstW - 20
	ty := float64(target.bounds.Y + target.bounds.H - playerDstH)
	p.targetX = engine.Clamp(tx, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(ty, playerMinY, playerMaxY)
	p.moving = true
	p.state = stateWalking
	p.facingLeft = false
	p.interactTarget = target
	p.dialogSys = ds
	p.onArrival = nil
	p.anim.Play("walk")
}

func (p *player) walkToAndDo(x, y float64, action func()) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(y-playerDstH/2, playerMinY, playerMaxY)
	p.moving = true
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = action
	p.anim.Play("walk")
}

func (p *player) update(dt float64) {
	p.breathTimer += dt
	p.anim.Update(dt)

	if p.moving {
		p.walkTimer += dt
		if p.walkTimer >= 0.15 {
			p.walkTimer -= 0.15
			p.walkFrame = (p.walkFrame + 1) % len(walkRotations)
		}
	} else {
		p.walkFrame = 0
		p.walkTimer = 0
	}

	if p.state == stateTalking {
		p.talkTimer += dt
		if p.talkTimer >= 0.22 {
			p.talkTimer -= 0.22
			p.talkOpen = !p.talkOpen
		}
	} else {
		p.talkTimer = 0
		p.talkOpen = false
	}

	if !p.moving {
		if p.state != stateTalking {
			p.state = stateIdle
			p.anim.Play("idle")
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
		p.anim.Play("idle")
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
	flip := sdl.FLIP_NONE
	if p.facingLeft {
		flip = sdl.FLIP_HORIZONTAL
	}

	dstW := int32(playerDstW)
	dstH := int32(playerDstH)
	dstY := int32(p.y)
	var rotation float64

	switch p.state {
	case stateWalking:
		rotation = walkRotations[p.walkFrame]
		if p.facingLeft {
			rotation = -rotation
		}
		dstY += int32(walkYOffsets[p.walkFrame])
	case stateTalking:
		bob := math.Sin(p.breathTimer*4.0) * 1.5
		dstY += int32(bob)
		if p.talkOpen {
			w := float64(playerDstW) * 0.97
			h := float64(playerDstH) * 1.01
			dstW = int32(w)
			dstH = int32(h)
		}
	default:
		breathVal := math.Sin(p.breathTimer * 2.0)
		dstY += int32(breathVal * 2.0)
		scale := 1.0 + breathVal*0.008
		dstW = int32(float64(playerDstW) * scale)
		dstH = int32(float64(playerDstH) * scale)
	}

	src := p.anim.CurrentRect()
	if src == nil {
		src = &sdl.Rect{X: 0, Y: 0, W: p.srcW, H: p.srcH}
	}
	renderer.CopyEx(p.tex, src,
		&sdl.Rect{X: int32(p.x), Y: dstY, W: dstW, H: dstH},
		rotation, nil, flip)
}

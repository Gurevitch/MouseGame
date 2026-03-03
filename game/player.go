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
	playerDstW      = 200
	playerDstH      = 280
	playerMinX      = 10.0
	playerMaxX      = engine.ScreenWidth - playerDstW - 10.0
	playerMinY      = 300.0
	playerMaxY      = 430.0
	walkFrameTime   = 0.12
	talkFrameTime   = 0.20
)

type playerState int

const (
	stateIdle playerState = iota
	stateWalking
	stateTalking
)

type direction int

const (
	dirRight direction = iota
	dirLeft
	dirUp
	dirDown
)

type player struct {
	walkSideFrames []spriteFrame
	walkUpFrames   []spriteFrame
	walkDownFrames []spriteFrame
	idleFrames     []spriteFrame
	talkFrames     []spriteFrame

	x, y       float64
	targetX    float64
	targetY    float64
	moving     bool
	facingLeft bool
	dir        direction
	state      playerState

	breathTimer  float64
	walkCycleIdx int
	walkTimer    float64
	talkCycleIdx int
	talkTimer    float64
	idleCycleIdx int
	idleTimer    float64

	interactTarget *npc
	dialogSys      *dialogSystem
	onArrival      func()
}

func gridFrameToSprite(gf engine.GridFrame) spriteFrame {
	return spriteFrame{tex: gf.Tex, w: gf.W, h: gf.H}
}

func newPlayer(renderer *sdl.Renderer) *player {
	p := &player{
		x: 630,
		y: float64(engine.ScreenHeight) - playerDstH - 160,
	}

	posesGrid := engine.SpriteGridFromPNG(renderer,
		"assets/images/player/Gemini_Generated_Image_kkasyqkkasyqkkas.png", 5, 4)

	dirGrid := engine.SpriteGridFromPNG(renderer,
		"assets/images/player/Gemini_Generated_Image_vt2ol9vt2ol9vt2o.png", 4, 3)

	// Side-view walk cycle: poses rows 0-1 (10 frames)
	for r := 0; r < 2; r++ {
		for c := 0; c < 5; c++ {
			p.walkSideFrames = append(p.walkSideFrames, gridFrameToSprite(posesGrid[r][c]))
		}
	}

	// Walk away/back: directions row 0 (4 frames)
	for c := 0; c < 4; c++ {
		p.walkUpFrames = append(p.walkUpFrames, gridFrameToSprite(dirGrid[0][c]))
	}

	// Walk toward camera: directions row 2 (4 frames)
	for c := 0; c < 4; c++ {
		p.walkDownFrames = append(p.walkDownFrames, gridFrameToSprite(dirGrid[2][c]))
	}

	// Idle: single front-facing frame from directions row 1, col 0
	p.idleFrames = append(p.idleFrames, gridFrameToSprite(dirGrid[1][0]))

	// Talk: poses row 2 cols 0-1 (standing + waving)
	p.talkFrames = append(p.talkFrames, gridFrameToSprite(posesGrid[2][0]))
	p.talkFrames = append(p.talkFrames, gridFrameToSprite(posesGrid[2][1]))

	return p
}

func (p *player) currentWalkFrames() []spriteFrame {
	switch p.dir {
	case dirUp:
		return p.walkUpFrames
	case dirDown:
		return p.walkDownFrames
	default:
		return p.walkSideFrames
	}
}

func (p *player) currentSprite() spriteFrame {
	switch p.state {
	case stateWalking:
		frames := p.currentWalkFrames()
		if len(frames) == 0 {
			return p.idleFrames[0]
		}
		return frames[p.walkCycleIdx%len(frames)]
	case stateTalking:
		if len(p.talkFrames) > 0 {
			return p.talkFrames[p.talkCycleIdx%len(p.talkFrames)]
		}
		return p.idleFrames[0]
	default:
		if len(p.idleFrames) > 0 {
			return p.idleFrames[0]
		}
		return p.walkSideFrames[0]
	}
}

func (p *player) setTarget(x, y float64) {
	tx := engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	ty := engine.Clamp(y-playerDstH/2, playerMinY, playerMaxY)
	p.targetX = tx
	p.targetY = ty
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
	var ty float64
	if target.elevated {
		ty = p.y
	} else {
		ty = float64(target.bounds.Y + target.bounds.H - playerDstH)
	}
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

func (p *player) update(dt float64, blockers []sdl.Rect) {
	p.breathTimer += dt

	if p.moving {
		p.walkTimer += dt
		if p.walkTimer >= walkFrameTime {
			p.walkTimer -= walkFrameTime
			frames := p.currentWalkFrames()
			if len(frames) > 0 {
				p.walkCycleIdx = (p.walkCycleIdx + 1) % len(frames)
			}
		}
	} else {
		p.walkCycleIdx = 0
		p.walkTimer = 0
	}

	if p.state == stateTalking {
		p.talkTimer += dt
		if p.talkTimer >= talkFrameTime {
			p.talkTimer -= talkFrameTime
			if len(p.talkFrames) > 0 {
				p.talkCycleIdx = (p.talkCycleIdx + 1) % len(p.talkFrames)
			}
		}
	} else {
		p.talkTimer = 0
		p.talkCycleIdx = 0
	}

	if p.state == stateIdle {
		p.idleCycleIdx = 0
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

	// Determine direction from movement delta
	if math.Abs(dy) > math.Abs(dx)*1.2 {
		if dy < 0 {
			p.dir = dirUp
		} else {
			p.dir = dirDown
		}
	} else {
		if dx < 0 {
			p.dir = dirLeft
		} else {
			p.dir = dirRight
		}
	}
	p.facingLeft = dx < 0

	if dist < 3.0 {
		p.x = p.targetX
		p.y = p.targetY
		p.moving = false
		p.state = stateIdle
		p.dir = dirDown
		p.facingLeft = false
		p.idleCycleIdx = 0
		p.idleTimer = 0
		if p.interactTarget != nil && p.dialogSys != nil {
			p.startNPCDialog()
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

	for _, b := range blockers {
		pr := sdl.Rect{X: int32(p.x), Y: int32(p.y), W: playerDstW, H: playerDstH}
		if pr.HasIntersection(&b) {
			playerCX := p.x + playerDstW/2
			blockerCX := float64(b.X) + float64(b.W)/2
			if playerCX < blockerCX {
				p.x = float64(b.X) - playerDstW
			} else {
				p.x = float64(b.X + b.W)
			}
			if p.targetX < float64(b.X+b.W) && p.targetX+playerDstW > float64(b.X) {
				p.moving = false
				if p.interactTarget != nil && p.dialogSys != nil {
					p.startNPCDialog()
				} else if p.onArrival != nil {
					fn := p.onArrival
					p.onArrival = nil
					p.state = stateIdle
					p.dir = dirDown
					p.facingLeft = false
					fn()
				} else {
					p.state = stateIdle
					p.dir = dirDown
					p.facingLeft = false
				}
				p.idleCycleIdx = 0
				p.idleTimer = 0
			}
		}
	}
}

func (p *player) startNPCDialog() {
	n := p.interactTarget
	ds := p.dialogSys
	if n == nil || ds == nil {
		return
	}
	p.state = stateTalking
	p.facingLeft = p.x > float64(n.bounds.X)

	if n.altDialogFunc != nil {
		entries, cb := n.altDialogFunc()
		if entries != nil {
			ds.startDialogWithCallback(entries, cb)
			p.interactTarget = nil
			return
		}
	}

	cb := n.onDialogEnd
	if !n.dialogDone {
		ds.startDialogWithCallback(n.dialog, func() {
			n.dialogDone = true
			if cb != nil {
				cb()
			}
		})
	} else {
		ds.startDialogWithCallback(n.dialog, nil)
	}
	p.interactTarget = nil
}

func (p *player) containsPoint(x, y int32) bool {
	pt := sdl.Point{X: x, Y: y}
	r := sdl.Rect{X: int32(p.x), Y: int32(p.y), W: playerDstW, H: playerDstH}
	return pt.InRect(&r)
}

func (p *player) footCenter() (int32, int32) {
	cx := int32(p.x) + playerDstW/2
	fy := int32(p.y) + playerDstH
	return cx, fy
}

func (p *player) draw(renderer *sdl.Renderer) {
	frame := p.currentSprite()

	frameScale := float64(playerDstH) / float64(frame.h)
	dstW := int32(float64(frame.w) * frameScale)
	dstH := int32(playerDstH)

	dstX := int32(p.x) + (playerDstW-dstW)/2
	dstY := int32(p.y)

	flip := sdl.FLIP_NONE
	if p.facingLeft && p.state != stateIdle {
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

	cx, fy := p.footCenter()
	drawShadow(renderer, cx, fy, playerDstW-20)

	renderer.CopyEx(frame.tex, nil,
		&sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH},
		0, nil, flip)
}

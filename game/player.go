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
	playerDstH      = 195
	playerMinX      = 10.0
	playerMaxX      = engine.ScreenWidth - playerDstW - 10.0
	playerMinY      = 300.0
	playerMaxY      = 430.0
	walkFrameTime   = 0.12
	talkFrameTime   = 0.07
	actionFrameTime = 0.10
)

type playerState int

const (
	stateIdle playerState = iota
	stateWalking
	stateTalking
	stateGrabbing
	stateUsing
	stateExamining
	stateReacting
	stateShowInventory
)

type direction int

const (
	dirRight direction = iota
	dirLeft
	dirUp
	dirDown
)

type player struct {
	walkSideFrames  []spriteFrame
	walkUpFrames    []spriteFrame
	walkDownFrames  []spriteFrame
	idleFrontFrames []spriteFrame
	idleSideFrames  []spriteFrame
	idleBackFrames  []spriteFrame
	talkFrames      []spriteFrame
	talkSideFrames  []spriteFrame
	grabFrames      []spriteFrame
	useItemFrames   []spriteFrame
	examineFrames   []spriteFrame
	reactFrames     []spriteFrame
	showInvFrames   []spriteFrame

	x, y           float64
	targetX        float64
	targetY        float64
	moving         bool
	allowOffscreen bool
	facingLeft     bool
	dir            direction
	state          playerState

	breathTimer    float64
	walkCycleIdx   int
	walkTimer      float64
	talkCycleIdx   int
	talkTimer      float64
	actionIdx      int
	actionTimer    float64
	actionCallback func()

	interactTarget *npc
	dialogSys      *dialogSystem
	onArrival      func()

	sceneMinY float64
	sceneMaxY float64
}

func stripFrames(renderer *sdl.Renderer, path string, cols int) []spriteFrame {
	return gridFrames(renderer, path, cols, 1)
}

func gridFrames(renderer *sdl.Renderer, path string, cols, rows int) []spriteFrame {
	grid := engine.SpriteGridFromPNG(renderer, path, cols, rows)
	var frames []spriteFrame
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			gf := grid[r][c]
			frames = append(frames, spriteFrame{tex: gf.Tex, w: gf.W, h: gf.H})
		}
	}
	return frames
}

func newPlayer(renderer *sdl.Renderer) *player {
	p := &player{
		x: 630,
		y: float64(engine.ScreenHeight) - playerDstH - 100,
	}

	p.walkSideFrames = gridFrames(renderer, "assets/images/player/PP walk left.png", 8, 2)
	p.walkDownFrames = gridFrames(renderer, "assets/images/player/PP walk front.png", 8, 2)
	p.walkUpFrames = gridFrames(renderer, "assets/images/player/PP walk back.png", 8, 2)

	// Idle images — use all frames for animated idle
	p.idleFrontFrames = gridFrames(renderer, "assets/images/player/PP idle front.png", 8, 2)
	p.idleSideFrames = gridFrames(renderer, "assets/images/player/PP idle side.png", 8, 2)
	p.idleBackFrames = gridFrames(renderer, "assets/images/player/PP idle back.png", 8, 2)

	p.talkFrames = gridFrames(renderer, "assets/images/player/PP talk front.png", 8, 2)
	p.talkSideFrames = gridFrames(renderer, "assets/images/player/PP talk side.png", 8, 2)

	p.grabFrames = gridFrames(renderer, "assets/images/player/PP grab flower.png", 8, 2)

	celebrateFrames := gridFrames(renderer, "assets/images/player/PP celebrate.png", 8, 2)
	p.reactFrames = celebrateFrames
	if len(celebrateFrames) >= 2 {
		p.showInvFrames = celebrateFrames[0:2]
	}

	p.examineFrames = gridFrames(renderer, "assets/images/player/PP sneak examine.png", 8, 2)
	p.useItemFrames = gridFrames(renderer, "assets/images/player/PP sneak use.png", 8, 2)

	p.dir = dirDown

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

func (p *player) currentIdleFrames() []spriteFrame {
	switch p.dir {
	case dirUp:
		return p.idleBackFrames
	case dirLeft, dirRight:
		return p.idleSideFrames
	default:
		return p.idleFrontFrames
	}
}

func (p *player) currentTalkFrames() []spriteFrame {
	switch p.dir {
	case dirLeft, dirRight:
		return p.talkSideFrames
	default:
		return p.talkFrames
	}
}

func firstAvailableFrame(groups ...[]spriteFrame) spriteFrame {
	for _, group := range groups {
		if len(group) > 0 {
			return group[0]
		}
	}
	return spriteFrame{}
}

func (p *player) actionFrames() []spriteFrame {
	switch p.state {
	case stateGrabbing:
		return p.grabFrames
	case stateUsing:
		return p.useItemFrames
	case stateExamining:
		return p.examineFrames
	case stateReacting:
		return p.reactFrames
	case stateShowInventory:
		return p.showInvFrames
	}
	return nil
}

func (p *player) currentSprite() spriteFrame {
	switch p.state {
	case stateWalking:
		frames := p.currentWalkFrames()
		if len(frames) == 0 {
			return firstAvailableFrame(p.currentIdleFrames(), p.walkDownFrames, p.walkSideFrames, p.walkUpFrames)
		}
		return frames[p.walkCycleIdx%len(frames)]
	case stateTalking:
		frames := p.currentTalkFrames()
		if len(frames) > 0 {
			return frames[p.talkCycleIdx%len(frames)]
		}
		return firstAvailableFrame(p.currentIdleFrames(), p.walkDownFrames, p.walkSideFrames, p.walkUpFrames)
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		frames := p.actionFrames()
		if len(frames) > 0 {
			return frames[p.actionIdx%len(frames)]
		}
		return firstAvailableFrame(p.currentIdleFrames(), p.walkDownFrames)
	default:
		frames := p.currentIdleFrames()
		if len(frames) > 0 {
			idx := int(p.breathTimer*4) % len(frames)
			return frames[idx]
		}
		return firstAvailableFrame(p.walkDownFrames, p.walkSideFrames, p.walkUpFrames)
	}
}

func (p *player) minY() float64 {
	if p.sceneMinY > 0 {
		return p.sceneMinY
	}
	return playerMinY
}

func (p *player) maxY() float64 {
	if p.sceneMaxY > 0 {
		return p.sceneMaxY
	}
	return playerMaxY
}

func (p *player) setTarget(x, y float64) {
	tx := engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	ty := engine.Clamp(y-playerDstH/2, p.minY(), p.maxY())
	p.targetX = tx
	p.targetY = ty
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = nil
}

func (p *player) walkToAndInteract(target *npc, ds *dialogSystem) {
	npcCenter := float64(target.bounds.X + target.bounds.W/2)
	npcLeft := float64(target.bounds.X)
	npcRight := float64(target.bounds.X + target.bounds.W)

	pickSide := func(preferRight bool) float64 {
		if preferRight {
			return npcRight + 10
		}
		return npcLeft - playerDstW - 10
	}

	preferred := npcCenter >= engine.ScreenWidth/2
	tx := pickSide(!preferred)
	tx = engine.Clamp(tx, playerMinX, playerMaxX)

	if tx < npcRight && tx+playerDstW > npcLeft {
		tx = pickSide(preferred)
		tx = engine.Clamp(tx, playerMinX, playerMaxX)
	}

	var ty float64
	if target.elevated {
		ty = p.y
	} else {
		npcFootY := float64(target.bounds.Y + target.bounds.H)
		ty = npcFootY - playerDstH + 4
	}
	p.targetX = tx
	p.targetY = engine.Clamp(ty, p.minY(), p.maxY())
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
	p.interactTarget = target
	p.dialogSys = ds
	p.onArrival = nil
}

func (p *player) walkToAndDo(x, y float64, action func()) {
	p.targetX = engine.Clamp(x-playerDstW/2, playerMinX, playerMaxX)
	p.targetY = engine.Clamp(y-playerDstH/2, p.minY(), p.maxY())
	p.moving = true
	p.allowOffscreen = false
	p.state = stateWalking
	p.interactTarget = nil
	p.onArrival = action
}

func (p *player) walkToExit(dir arrowDir, action func()) {
	p.targetY = engine.Clamp(p.y, p.minY(), p.maxY())
	switch dir {
	case arrowLeft:
		p.targetX = -playerDstW
		p.dir = dirLeft
		p.facingLeft = true
	case arrowRight:
		p.targetX = engine.ScreenWidth + playerDstW
		p.dir = dirRight
		p.facingLeft = false
	case arrowDown:
		p.targetX = p.x
		p.targetY = engine.ScreenHeight + playerDstH
		p.dir = dirDown
		p.facingLeft = false
	case arrowUp:
		p.targetX = p.x
		p.targetY = -playerDstH
		p.dir = dirUp
		p.facingLeft = false
	case arrowDownRight:
		p.targetX = engine.ScreenWidth + playerDstW
		p.targetY = engine.ScreenHeight + playerDstH
		p.dir = dirDown
		p.facingLeft = false
	default:
		p.allowOffscreen = false
		p.walkToAndDo(p.x+playerDstW/2, p.y+playerDstH/2, action)
		return
	}
	p.moving = true
	p.allowOffscreen = true
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
		p.allowOffscreen = false
	}

	if p.state == stateTalking {
		p.talkTimer += dt
		if p.talkTimer >= talkFrameTime {
			p.talkTimer -= talkFrameTime
			frames := p.currentTalkFrames()
			if len(frames) > 0 {
				p.talkCycleIdx = (p.talkCycleIdx + 1) % len(frames)
			}
		}
	} else {
		p.talkTimer = 0
		p.talkCycleIdx = 0
	}

	switch p.state {
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		p.actionTimer += dt
		if p.actionTimer >= actionFrameTime {
			p.actionTimer -= actionFrameTime
			frames := p.actionFrames()
			if len(frames) > 0 {
				p.actionIdx++
				if p.actionIdx >= len(frames) {
					p.actionIdx = 0
					p.state = stateIdle
					if p.actionCallback != nil {
						fn := p.actionCallback
						p.actionCallback = nil
						fn()
					}
				}
			}
		}
	default:
		p.actionIdx = 0
		p.actionTimer = 0
	}

	if !p.moving {
		if p.state != stateTalking {
			p.state = stateIdle
			p.dir = dirDown
			p.facingLeft = false
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
	nextX := p.x + (dx/dist)*speed*dt
	nextY := p.y + (dy/dist)*speed*dt
	if p.allowOffscreen {
		p.x = nextX
		p.y = nextY
	} else {
		p.x = engine.Clamp(nextX, playerMinX, playerMaxX)
		p.y = engine.Clamp(nextY, p.minY(), p.maxY())
	}

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
					fn()
				} else {
					p.state = stateIdle
				}
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
	p.talkCycleIdx = 0
	p.talkTimer = 0
	npcCenter := float64(n.bounds.X + n.bounds.W/2)
	playerCenter := p.x + playerDstW/2
	p.facingLeft = playerCenter > npcCenter
	if p.facingLeft {
		p.dir = dirLeft
	} else {
		p.dir = dirRight
	}

	if len(n.talkGrid) > 0 {
		n.setAnimState(npcAnimTalk)
	}

	wrapCb := func(inner func()) func() {
		target := n
		return func() {
			if len(target.talkGrid) > 0 {
				target.setAnimState(npcAnimIdle)
			}
			if inner != nil {
				inner()
			}
		}
	}

	if n.altDialogFunc != nil {
		entries, cb := n.altDialogFunc()
		if entries != nil {
			ds.startDialogWithCallback(entries, wrapCb(cb))
			p.interactTarget = nil
			return
		}
	}

	cb := n.onDialogEnd
	ds.startDialogWithCallback(n.dialog, wrapCb(func() {
		if cb != nil {
			cb()
		}
		n.dialogDone = true
	}))
	p.interactTarget = nil
}

func (p *player) playAction(s playerState, cb func()) {
	p.state = s
	p.actionIdx = 0
	p.actionTimer = 0
	p.actionCallback = cb
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

func (p *player) footY() int32 {
	_, fy := p.footCenter()
	return fy
}

func (p *player) depthScale() float64 {
	progress := engine.Clamp((p.y-playerMinY)/(playerMaxY-playerMinY), 0, 1)
	return 0.88 + progress*0.18
}

func (p *player) draw(renderer *sdl.Renderer) {
	frame := p.currentSprite()
	if frame.tex == nil || frame.h == 0 {
		return
	}

	scaledHeight := int32(float64(playerDstH) * p.depthScale())
	frameScale := float64(scaledHeight) / float64(frame.h)
	dstW := int32(float64(frame.w) * frameScale)
	dstH := scaledHeight

	dstX := int32(p.x) + (playerDstW-dstW)/2
	dstY := p.footY() - dstH

	flip := sdl.FLIP_NONE
	if p.dir == dirLeft {
		flip = sdl.FLIP_HORIZONTAL
	}

	switch p.state {
	case stateIdle:
		breathVal := math.Sin(p.breathTimer * 2.0)
		dstY += int32(breathVal * 2.0)
	case stateTalking:
		bob := math.Sin(p.breathTimer*3.0) * 1.5
		dstY += int32(bob)
	case stateGrabbing, stateUsing, stateExamining, stateReacting, stateShowInventory:
		bob := math.Sin(p.breathTimer*2.5) * 1.0
		dstY += int32(bob)
	}

	cx, fy := p.footCenter()
	drawShadow(renderer, cx, fy, int32(float64(playerDstW-20)*p.depthScale()))

	renderer.CopyEx(frame.tex, nil,
		&sdl.Rect{X: dstX, Y: dstY, W: dstW, H: dstH},
		0, nil, flip)
}

package engine

import "github.com/veandco/go-sdl2/sdl"

type animation struct {
	frames        []sdl.Rect
	frameDuration float64
}

type Animator struct {
	animations   map[string]*animation
	current      string
	currentFrame int
	timer        float64
	Looping      bool
}

func NewAnimator() *Animator {
	return &Animator{
		animations: make(map[string]*animation),
		Looping:    true,
	}
}

func (a *Animator) AddAnimation(name string, frames []sdl.Rect, frameDuration float64) {
	a.animations[name] = &animation{
		frames:        frames,
		frameDuration: frameDuration,
	}
}

func (a *Animator) Play(name string) {
	if a.current == name {
		return
	}
	a.current = name
	a.currentFrame = 0
	a.timer = 0
}

func (a *Animator) Update(dt float64) {
	anim, ok := a.animations[a.current]
	if !ok || len(anim.frames) <= 1 {
		return
	}

	a.timer += dt
	if a.timer >= anim.frameDuration {
		a.timer -= anim.frameDuration
		a.currentFrame++
		if a.currentFrame >= len(anim.frames) {
			if a.Looping {
				a.currentFrame = 0
			} else {
				a.currentFrame = len(anim.frames) - 1
			}
		}
	}
}

func (a *Animator) CurrentRect() *sdl.Rect {
	anim, ok := a.animations[a.current]
	if !ok || len(anim.frames) == 0 {
		return nil
	}
	r := anim.frames[a.currentFrame]
	return &r
}

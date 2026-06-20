package game

import (
	"fmt"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// SeqActionType defines what kind of action a sequence step performs.
type SeqActionType int

const (
	SeqDialog         SeqActionType = iota // Show dialog entries
	SeqWait                                // Wait for duration (seconds)
	SeqTransition                          // Transition to a scene
	SeqCallback                            // Run a callback function
	SeqSetVar                              // Set a variable in VarStore
	SeqNPCAnim                             // Flip a named NPC's animation state (idle/talk)
	SeqNPCStrange                          // Flip a named NPC in / out of the strange state
	SeqSetSceneBG                          // Swap a scene's background for an alternate
	SeqHidePlayer                          // Suppress PP rendering (hide/show)
	SeqPlayerSleep                         // Toggle PP's sleeping sprite overlay
	SeqPlayerWake                          // Kick off PP's waking animation
	SeqStartDay                            // Advance to the next camp day (resets chapter-scope state)
	SeqNPCHidden                           // Toggle an NPC's hidden + silent flags (un-hide on un-silent)
	SeqNPCTeleport                         // Snap an NPC to an absolute (x, y) position
	SeqNPCMove                             // Linearly interpolate an NPC's position over Duration seconds
	SeqGiveItem                            // Add an item to inventory by id (silent - no UI pop)
	SeqPlayerAnim                          // Play a one-shot PP animation (e.g. "receive_map") for Duration seconds
	SeqNPCOneShotAnim                      // Play a one-shot named anim on an NPC ("give_map") for Duration seconds
	SeqTweenItem                           // Lerp a sprite from (FromX,FromY) → (TargetX,TargetY) over Duration (e.g. thrown map)
)

// SeqStep is one step in a sequence. Fields are used by whichever Action
// type the step declares; unused fields stay zero.
type SeqStep struct {
	Action   SeqActionType
	Dialog   []dialogEntry // for SeqDialog
	Duration float64       // for SeqWait (seconds)
	Scene    string        // for SeqTransition, SeqNPCAnim, SeqNPCStrange, SeqSetSceneBG
	Callback func()        // for SeqCallback
	VarScope string        // for SeqSetVar
	VarName  string        // for SeqSetVar
	VarValue int           // for SeqSetVar
	// NPC / scene-action fields
	NPC     string // NPC name (matches npc.name) for SeqNPCAnim, SeqNPCStrange, SeqNPCOneShotAnim
	Anim    string // "idle" | "talk" for SeqNPCAnim; arbitrary anim name for SeqPlayerAnim, SeqNPCOneShotAnim
	Strange bool   // for SeqNPCStrange
	BGKey   string // identifier of the alternate background (e.g. "night")
	Hide    bool   // for SeqHidePlayer, SeqPlayerSleep, SeqNPCHidden
	DayNum  int    // for SeqStartDay (1 or 2)
	ItemID  string // for SeqGiveItem - matches assets/data/items.json id
	// NPC move/teleport targets - absolute pixel coords.
	TargetX int32
	TargetY int32
	// EndScale, when > 0, makes a SeqNPCMove also lerp the NPC's render scale
	// (npc.extraScale) to EndScale over the move - so it shrinks as it walks
	// "into" the scene (Jake into his cabin).
	EndScale float64
	// Runtime-only scratch space for in-progress moves (start pos + elapsed).
	moveStartX     int32
	moveStartY     int32
	moveStartScale float64
	moveElapsed    float64

	// SeqTweenItem - visible thrown/flying sprite that lerps across screen.
	Sprite string // PNG path; loaded lazily on first execute.
	FromX  int32  // start screen x
	FromY  int32  // start screen y
	// Runtime-only: loaded texture + dimensions, set on first execute and
	// cleared when the step finishes so the Draw hook stops rendering.
	tweenTex *sdl.Texture
	tweenW   int32
	tweenH   int32
}

// Sequence is an ordered list of steps that play automatically.
type Sequence struct {
	Name    string
	Steps   []SeqStep
	current int
	timer   float64
	waiting bool
	active  bool
}

// SequencePlayer manages playing sequences.
type SequencePlayer struct {
	current *Sequence
	game    *Game
}

func newSequencePlayer(g *Game) *SequencePlayer {
	return &SequencePlayer{game: g}
}

// Play starts a sequence.
func (sp *SequencePlayer) Play(seq *Sequence) {
	sp.current = seq
	seq.current = 0
	seq.timer = 0
	seq.waiting = false
	seq.active = true
	sp.executeStep()
}

// IsPlaying returns true if a sequence is currently playing.
func (sp *SequencePlayer) IsPlaying() bool {
	return sp.current != nil && sp.current.active
}

// Update advances the sequence each frame.
func (sp *SequencePlayer) Update(dt float64) {
	if sp.current == nil || !sp.current.active {
		return
	}

	seq := sp.current

	if seq.waiting {
		seq.timer += dt
		step := seq.Steps[seq.current]

		switch step.Action {
		case SeqWait:
			if seq.timer >= step.Duration {
				seq.waiting = false
				sp.nextStep()
			}
		case SeqDialog:
			if !sp.game.dialog.active {
				seq.waiting = false
				sp.nextStep()
			}
		case SeqTransition:
			if !sp.game.sceneMgr.transitioning {
				seq.waiting = false
				sp.nextStep()
			}
		case SeqNPCMove:
			// Linearly interpolate from (moveStartX,Y) to (TargetX,Y) over
			// Duration. NPC bounds are updated every tick so the render loop
			// picks up the tween. stepPtr mutates the array in place.
			stepPtr := &seq.Steps[seq.current]
			stepPtr.moveElapsed += dt
			t := stepPtr.moveElapsed / step.Duration
			if t >= 1.0 {
				t = 1.0
			}
			if n := sp.findNPC(step.Scene, step.NPC); n != nil {
				dx := float64(step.TargetX - stepPtr.moveStartX)
				dy := float64(step.TargetY - stepPtr.moveStartY)
				n.bounds.X = stepPtr.moveStartX + int32(dx*t)
				n.bounds.Y = stepPtr.moveStartY + int32(dy*t)
				if step.EndScale > 0 {
					n.extraScale = stepPtr.moveStartScale + (step.EndScale-stepPtr.moveStartScale)*t
				}
			}
			if t >= 1.0 {
				seq.waiting = false
				sp.nextStep()
			}

		case SeqNPCOneShotAnim:
			dur := step.Duration
			if dur <= 0 {
				dur = 1.0
			}
			if seq.timer >= dur {
				if n := sp.findNPC(step.Scene, step.NPC); n != nil {
					n.endOneShotAnim()
				}
				seq.waiting = false
				sp.nextStep()
			}

		case SeqTweenItem:
			// Lerp the item from (FromX,FromY) → (TargetX,TargetY) over
			// Duration seconds. We use moveElapsed on the step pointer to
			// avoid recomputing from the global seq.timer in case other
			// timers stack. Draw hook reads the same elapsed value to
			// position the sprite. Done when elapsed >= Duration.
			stepPtr := &seq.Steps[seq.current]
			stepPtr.moveElapsed += dt
			if stepPtr.moveElapsed >= step.Duration {
				// Clear runtime texture so Draw stops emitting; the texture
				// itself stays cached on the game-wide tweenItemCache so
				// repeat plays don't re-load it.
				stepPtr.tweenTex = nil
				seq.waiting = false
				sp.nextStep()
			}
		}
	}
}

// Draw renders any per-step visuals the sequence player owns. Currently
// just the SeqTweenItem projectile sprite (thrown map, etc.). Called from
// Game.draw AFTER scene actors so the projectile renders on top of the
// world without disturbing NPC rendering. Safe to call when no sequence
// is active.
func (sp *SequencePlayer) Draw(renderer *sdl.Renderer) {
	if sp.current == nil || !sp.current.active || !sp.current.waiting {
		return
	}
	seq := sp.current
	step := &seq.Steps[seq.current]
	if step.Action != SeqTweenItem || step.tweenTex == nil {
		return
	}
	dur := step.Duration
	if dur <= 0 {
		dur = 0.001
	}
	t := step.moveElapsed / dur
	if t > 1.0 {
		t = 1.0
	}
	dx := float64(step.TargetX - step.FromX)
	dy := float64(step.TargetY - step.FromY)
	cx := step.FromX + int32(dx*t)
	// User 2026-05-22: PARABOLIC arc - the projectile flies HIGH at the
	// midpoint and lands at the target. arcHeight is the pixels above the
	// straight-line midpoint the projectile reaches at t=0.5.
	// Formula: arc(t) = 4*h*t*(1-t) → 0 at t=0/1, max=h at t=0.5.
	const arcHeight = 200.0
	arcLift := arcHeight * 4.0 * t * (1.0 - t)
	cy := step.FromY + int32(dy*t) - int32(arcLift)
	dst := sdl.Rect{
		X: cx - step.tweenW/2,
		Y: cy - step.tweenH/2,
		W: step.tweenW,
		H: step.tweenH,
	}
	renderer.Copy(step.tweenTex, nil, &dst)
}

func (sp *SequencePlayer) nextStep() {
	sp.current.current++
	if sp.current.current >= len(sp.current.Steps) {
		sp.current.active = false
		fmt.Printf("Sequence '%s' completed\n", sp.current.Name)
		return
	}
	sp.executeStep()
}

func (sp *SequencePlayer) executeStep() {
	seq := sp.current
	if seq.current >= len(seq.Steps) {
		seq.active = false
		return
	}

	step := seq.Steps[seq.current]

	switch step.Action {
	case SeqDialog:
		sp.game.dialog.startDialog(step.Dialog)
		seq.waiting = true
		seq.timer = 0

	case SeqWait:
		seq.waiting = true
		seq.timer = 0

	case SeqTransition:
		sp.game.sceneMgr.transitionTo(step.Scene, sp.game.player)
		seq.waiting = true
		seq.timer = 0

	case SeqCallback:
		if step.Callback != nil {
			step.Callback()
		}
		sp.nextStep() // callbacks are instant

	case SeqSetVar:
		sp.game.vars.Set(step.VarScope, step.VarName, step.VarValue)
		sp.nextStep()

	case SeqNPCAnim:
		n := sp.findNPC(step.Scene, step.NPC)
		if n == nil {
			fmt.Printf("[SeqNPCAnim] NPC not found: scene=%q npc=%q anim=%q\n",
				step.Scene, step.NPC, step.Anim)
		} else {
			switch step.Anim {
			case "talk":
				n.endOneShotAnim()
				n.restoreSwappedIdle()
				n.setAnimState(npcAnimTalk)
			case "idle":
				n.endOneShotAnim()
				n.restoreSwappedIdle()
				n.setAnimState(npcAnimIdle)
			default:
				// User 2026-05-12: looping named animation (e.g. Higgins's
				// "walk_back" during an npc_move).
				if frames, ok := n.oneShotAnims[step.Anim]; ok && len(frames) > 0 {
					n.swapIdleForOneShot(step.Anim)
					fmt.Printf("[SeqNPCAnim] swapped idle to %q (%d frames) on %q\n",
						step.Anim, len(frames), step.NPC)
				} else {
					fmt.Printf("[SeqNPCAnim] anim %q not registered on %q (frames=%d) - falling back to idle\n",
						step.Anim, step.NPC, len(frames))
					n.setAnimState(npcAnimIdle)
				}
			}
		}
		sp.nextStep()

	case SeqNPCStrange:
		if n := sp.findNPC(step.Scene, step.NPC); n != nil {
			n.setStrange(step.Strange)
		}
		sp.nextStep()

	case SeqSetSceneBG:
		sp.game.setSceneAltBG(step.Scene, step.BGKey)
		sp.nextStep()

	case SeqHidePlayer:
		sp.game.nightHidePlayer = step.Hide
		sp.nextStep()

	case SeqPlayerSleep:
		sp.game.playerSleeping = step.Hide // Hide=true → sleeping on
		if step.Hide {
			sp.game.sleepingFrameIdx = 0
			sp.game.sleepingTimer = 0
			sp.game.wakingPhase = 0
		} else if sp.game.player != nil {
			// When the sleep overlay turns off (waking complete), snap PP's
			// coords so the normal idle picks up at the EXACT spot the wake
			// animation drew at, with no jump.
			// PR#3 (2026-06-12): the wake draws its opaque foot at screen
			// y=565, centre-X 337 (see game.go camp_night draw). The idle
			// path anchors PP's foot at player.y+playerDstH, so set y so the
			// foot lands on 565 (was 650 → PP dropped ~85px on wake).
			sp.game.player.x = float64(337 - playerDstW/2)
			sp.game.player.y = float64(565 - playerDstH)
			sp.game.player.targetX = sp.game.player.x
			sp.game.player.targetY = sp.game.player.y
			sp.game.player.moving = false
			sp.game.player.state = stateIdle
		}
		sp.nextStep()

	case SeqPlayerWake:
		sp.game.wakingPhase = 1
		sp.game.sleepingFrameIdx = 0
		sp.game.sleepingTimer = 0
		sp.nextStep()

	case SeqStartDay:
		if step.DayNum == 2 {
			sp.game.startDay2()
			sp.game.day2Started = true
		}
		sp.nextStep()

	case SeqNPCHidden:
		if n := sp.findNPC(step.Scene, step.NPC); n != nil {
			n.hidden = step.Hide
			// Un-hide also un-silents so the NPC can be clicked; re-hiding
			// keeps the silent flag as-is so callers can opt in to either.
			if !step.Hide {
				n.silent = false
			}
		}
		sp.nextStep()

	case SeqNPCTeleport:
		if n := sp.findNPC(step.Scene, step.NPC); n != nil {
			n.bounds.X = step.TargetX
			n.bounds.Y = step.TargetY
		}
		sp.nextStep()

	case SeqNPCMove:
		// Snapshot the current position onto the step itself so the Update
		// tick can lerp from it. Writing back through &seq.Steps[i] so the
		// array element (not a copy) stores the runtime scratch fields.
		if n := sp.findNPC(step.Scene, step.NPC); n != nil {
			stepPtr := &seq.Steps[seq.current]
			stepPtr.moveStartX = n.bounds.X
			stepPtr.moveStartY = n.bounds.Y
			stepPtr.moveStartScale = n.extraScale
			if stepPtr.moveStartScale <= 0 {
				stepPtr.moveStartScale = 1.0
			}
			stepPtr.moveElapsed = 0
			seq.waiting = true
			seq.timer = 0
			// Zero-duration moves snap immediately.
			if step.Duration <= 0 {
				n.bounds.X = step.TargetX
				n.bounds.Y = step.TargetY
				if step.EndScale > 0 {
					n.extraScale = step.EndScale
				}
				seq.waiting = false
				sp.nextStep()
			}
		} else {
			// NPC not found - skip the step rather than hang.
			sp.nextStep()
		}

	case SeqGiveItem:
		// Silent add - no inventory bar pop. Skips the add if PP already
		// owns the item (idempotent, matches giveMapItem's old guard).
		if step.ItemID != "" && sp.game.items != nil && sp.game.inv != nil {
			def, ok := sp.game.items.getDef(step.ItemID)
			if ok && !sp.game.inv.hasItem(def.Name) {
				if item := sp.game.items.createItem(step.ItemID); item != nil {
					sp.game.inv.addItem(item)
				}
			}
		}
		sp.nextStep()

	case SeqPlayerAnim:
		// One-shot PP animation. Wraps the existing player.playOneShot helper.
		// Currently supports anim names whose frames are loaded into the
		// player's named one-shot frame map (see player.playOneShot).
		dur := step.Duration
		if dur <= 0 {
			dur = 1.0
		}
		seq.waiting = true
		seq.timer = 0
		sp.game.player.playOneShot(step.Anim, dur, func() {
			seq.waiting = false
			sp.nextStep()
		})

	case SeqNPCOneShotAnim:
		// One-shot named anim on a specific NPC (e.g. Higgins's "give_map").
		// Falls back to a fixed wait if the NPC or anim isn't registered, so
		// missing assets don't deadlock the sequence.
		dur := step.Duration
		if dur <= 0 {
			dur = 1.0
		}
		n := sp.findNPC(step.Scene, step.NPC)
		if n != nil {
			n.playOneShotAnim(step.Anim, dur)
		}
		seq.waiting = true
		seq.timer = 0

	case SeqTweenItem:
		// Load the projectile sprite on first entry. Texture lives on the
		// step pointer's runtime scratch so Draw can pick it up. If the
		// PNG is missing, log + skip (don't deadlock the sequence).
		stepPtr := &seq.Steps[seq.current]
		if stepPtr.tweenTex == nil && step.Sprite != "" {
			tex, w, h := engine.SafeTextureFromPNGKeyed(sp.game.renderer, step.Sprite)
			if tex != nil {
				tex.SetBlendMode(sdl.BLENDMODE_BLEND)
				stepPtr.tweenTex = tex
				stepPtr.tweenW = w
				stepPtr.tweenH = h
			} else {
				fmt.Printf("SeqTweenItem: sprite %q missing - skipping projectile draw\n", step.Sprite)
			}
		}
		stepPtr.moveElapsed = 0
		// Zero-duration tweens snap immediately to end.
		if step.Duration <= 0 {
			stepPtr.tweenTex = nil
			sp.nextStep()
			return
		}
		seq.waiting = true
		seq.timer = 0
	}
}

// setSceneAltBG is called by SeqSetSceneBG to swap a scene's background for
// a pre-loaded alternate. Keys are "scene_name/variant" and live in
// g.sceneAltBGs, populated in Game.New. To add a new alt background, just
// load it into the map and emit a scene_bg step from any JSON sequence.
func (g *Game) setSceneAltBG(sceneName, bgKey string) {
	scene, ok := g.sceneMgr.scenes[sceneName]
	if !ok {
		return
	}
	if bg, ok := g.sceneAltBGs[sceneName+"/"+bgKey]; ok {
		scene.bg = bg
	}
	// 2026-06-12: keep Marcus's strange-idle lighting in step with the cabin
	// background (night during the cutscene, day on Day 2).
	if sceneName == "marcus_room" {
		for _, n := range scene.npcs {
			if n.name == "Marcus" {
				n.setStrangeVariant(bgKey == "night")
				break
			}
		}
	}
}

// findNPC locates a named NPC inside a named scene; returns nil if either is
// missing (logs silently - sequences that reference an NPC by typo will
// simply skip the action rather than crash the game).
func (sp *SequencePlayer) findNPC(sceneName, npcName string) *npc {
	s, ok := sp.game.sceneMgr.scenes[sceneName]
	if !ok {
		return nil
	}
	for _, n := range s.npcs {
		if n.name == npcName {
			return n
		}
	}
	return nil
}

// --- Helper constructors for readable sequence building ---

func dialogStep(entries ...dialogEntry) SeqStep {
	return SeqStep{Action: SeqDialog, Dialog: entries}
}

func dialogStepSlice(entries []dialogEntry) SeqStep {
	return SeqStep{Action: SeqDialog, Dialog: entries}
}

func waitStep(seconds float64) SeqStep {
	return SeqStep{Action: SeqWait, Duration: seconds}
}

func transitionStep(scene string) SeqStep {
	return SeqStep{Action: SeqTransition, Scene: scene}
}

func callbackStep(fn func()) SeqStep {
	return SeqStep{Action: SeqCallback, Callback: fn}
}

func setVarStep(scope, name string, value int) SeqStep {
	return SeqStep{Action: SeqSetVar, VarScope: scope, VarName: name, VarValue: value}
}

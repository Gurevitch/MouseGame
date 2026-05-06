package game

import "fmt"

// SeqActionType defines what kind of action a sequence step performs.
type SeqActionType int

const (
	SeqDialog       SeqActionType = iota // Show dialog entries
	SeqWait                              // Wait for duration (seconds)
	SeqTransition                        // Transition to a scene
	SeqCallback                          // Run a callback function
	SeqSetVar                            // Set a variable in VarStore
	SeqNPCAnim                           // Flip a named NPC's animation state (idle/talk)
	SeqNPCStrange                        // Flip a named NPC in / out of the strange state
	SeqSetSceneBG                        // Swap a scene's background for an alternate
	SeqHidePlayer                        // Suppress PP rendering (hide/show)
	SeqPlayerSleep                       // Toggle PP's sleeping sprite overlay
	SeqPlayerWake                        // Kick off PP's waking animation
	SeqStartDay                          // Advance to the next camp day (resets chapter-scope state)
	SeqNPCHidden                         // Toggle an NPC's hidden + silent flags (un-hide on un-silent)
	SeqNPCTeleport                       // Snap an NPC to an absolute (x, y) position
	SeqNPCMove                           // Linearly interpolate an NPC's position over Duration seconds
	SeqGiveItem                          // Add an item to inventory by id (silent — no UI pop)
	SeqPlayerAnim                        // Play a one-shot PP animation (e.g. "receive_map") for Duration seconds
	SeqNPCOneShotAnim                    // Play a one-shot named anim on an NPC ("give_map") for Duration seconds
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
	NPC      string // NPC name (matches npc.name) for SeqNPCAnim, SeqNPCStrange, SeqNPCOneShotAnim
	Anim     string // "idle" | "talk" for SeqNPCAnim; arbitrary anim name for SeqPlayerAnim, SeqNPCOneShotAnim
	Strange  bool   // for SeqNPCStrange
	BGKey    string // identifier of the alternate background (e.g. "night")
	Hide     bool   // for SeqHidePlayer, SeqPlayerSleep, SeqNPCHidden
	DayNum   int    // for SeqStartDay (1 or 2)
	ItemID   string // for SeqGiveItem — matches assets/data/items.json id
	// NPC move/teleport targets — absolute pixel coords.
	TargetX  int32
	TargetY  int32
	// Runtime-only scratch space for in-progress moves (start pos + elapsed).
	moveStartX int32
	moveStartY int32
	moveElapsed float64
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
	current  *Sequence
	game     *Game
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
		}
	}
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
		if n := sp.findNPC(step.Scene, step.NPC); n != nil {
			switch step.Anim {
			case "talk":
				n.setAnimState(npcAnimTalk)
			default:
				n.setAnimState(npcAnimIdle)
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
			stepPtr.moveElapsed = 0
			seq.waiting = true
			seq.timer = 0
			// Zero-duration moves snap immediately.
			if step.Duration <= 0 {
				n.bounds.X = step.TargetX
				n.bounds.Y = step.TargetY
				seq.waiting = false
				sp.nextStep()
			}
		} else {
			// NPC not found — skip the step rather than hang.
			sp.nextStep()
		}

	case SeqGiveItem:
		// Silent add — no inventory bar pop. Skips the add if PP already
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
}

// findNPC locates a named NPC inside a named scene; returns nil if either is
// missing (logs silently — sequences that reference an NPC by typo will
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

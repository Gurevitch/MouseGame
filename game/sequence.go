package game

import "fmt"

// SeqActionType defines what kind of action a sequence step performs.
type SeqActionType int

const (
	SeqDialog     SeqActionType = iota // Show dialog entries
	SeqWait                            // Wait for duration (seconds)
	SeqTransition                      // Transition to a scene
	SeqCallback                        // Run a callback function
	SeqSetVar                          // Set a variable in VarStore
)

// SeqStep is one step in a sequence.
type SeqStep struct {
	Action   SeqActionType
	Dialog   []dialogEntry // for SeqDialog
	Duration float64       // for SeqWait (seconds)
	Scene    string        // for SeqTransition
	Callback func()        // for SeqCallback
	VarScope string        // for SeqSetVar
	VarName  string        // for SeqSetVar
	VarValue int           // for SeqSetVar
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
		sp.nextStep() // var sets are instant
	}
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

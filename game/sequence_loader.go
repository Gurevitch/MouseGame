package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// sequenceJSON is the on-disk shape of a sequence file under
// `assets/data/sequences/*.json`. Each step carries the discriminator `type`
// and whatever subset of fields that step type reads — extra fields are
// ignored so the schema can grow without breaking old files.
type sequenceJSON struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Steps []stepJSON `json:"steps"`
}

type stepJSON struct {
	Type    string          `json:"type"`
	Seconds float64         `json:"seconds,omitempty"`
	Scene   string          `json:"scene,omitempty"`
	Dialog  []dialogEntry   `json:"dialog,omitempty"`
	Scope   string          `json:"scope,omitempty"`
	Name    string          `json:"name,omitempty"`
	Value   int             `json:"value,omitempty"`
	NPC     string          `json:"npc,omitempty"`
	Anim    string          `json:"anim,omitempty"`
	Strange *bool           `json:"strange,omitempty"`
	BG      string          `json:"bg,omitempty"`
	Hide    *bool           `json:"hide,omitempty"`
	Day     int             `json:"day,omitempty"`
	X       int32           `json:"x,omitempty"`
	Y       int32           `json:"y,omitempty"`
	// Raw kept for debugging unknown step types.
	_       json.RawMessage `json:"-"`
}

// sequenceStore caches loaded sequences by ID.
type sequenceStore struct {
	defs map[string]*Sequence
}

func newSequenceStore(dir string, game *Game) *sequenceStore {
	store := &sequenceStore{defs: make(map[string]*Sequence)}
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("sequence_loader: cannot read %s: %v\n", dir, err)
		return store
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		seq, err := loadSequenceFile(path, game)
		if err != nil {
			fmt.Printf("sequence_loader: %s: %v\n", path, err)
			continue
		}
		store.defs[seq.Name] = seq
		fmt.Printf("Loaded sequence: %s (%d steps)\n", seq.Name, len(seq.Steps))
	}
	return store
}

// Get returns a cached sequence by ID. Callers get a shared Sequence value —
// SequencePlayer.Play resets its playback cursor on entry so re-running is safe.
func (s *sequenceStore) Get(id string) *Sequence {
	return s.defs[id]
}

func loadSequenceFile(path string, game *Game) (*Sequence, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sj sequenceJSON
	if err := json.Unmarshal(data, &sj); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	id := sj.ID
	if id == "" {
		id = strings.TrimSuffix(filepath.Base(path), ".json")
	}
	seq := &Sequence{Name: id}
	for i, raw := range sj.Steps {
		step, err := compileStep(raw, game)
		if err != nil {
			return nil, fmt.Errorf("step %d (%s): %w", i, raw.Type, err)
		}
		seq.Steps = append(seq.Steps, step)
	}
	return seq, nil
}

// compileStep converts an on-disk stepJSON into the runtime SeqStep the
// sequence player executes. Unknown step types produce an error so authoring
// mistakes surface at load time, not midway through a cutscene.
func compileStep(j stepJSON, game *Game) (SeqStep, error) {
	switch j.Type {
	case "dialog":
		return SeqStep{Action: SeqDialog, Dialog: j.Dialog}, nil
	case "wait":
		return SeqStep{Action: SeqWait, Duration: j.Seconds}, nil
	case "transition":
		return SeqStep{Action: SeqTransition, Scene: j.Scene}, nil
	case "set_var":
		return SeqStep{Action: SeqSetVar, VarScope: j.Scope, VarName: j.Name, VarValue: j.Value}, nil
	case "npc_anim":
		return SeqStep{Action: SeqNPCAnim, Scene: j.Scene, NPC: j.NPC, Anim: j.Anim}, nil
	case "npc_strange":
		strange := j.Strange != nil && *j.Strange
		return SeqStep{Action: SeqNPCStrange, Scene: j.Scene, NPC: j.NPC, Strange: strange}, nil
	case "scene_bg":
		return SeqStep{Action: SeqSetSceneBG, Scene: j.Scene, BGKey: j.BG}, nil
	case "hide_player":
		hide := j.Hide != nil && *j.Hide
		return SeqStep{Action: SeqHidePlayer, Hide: hide}, nil
	case "player_sleep":
		sleep := j.Hide != nil && *j.Hide
		return SeqStep{Action: SeqPlayerSleep, Hide: sleep}, nil
	case "player_wake":
		return SeqStep{Action: SeqPlayerWake}, nil
	case "start_day":
		return SeqStep{Action: SeqStartDay, DayNum: j.Day}, nil
	case "npc_hidden":
		hide := j.Hide != nil && *j.Hide
		return SeqStep{Action: SeqNPCHidden, Scene: j.Scene, NPC: j.NPC, Hide: hide}, nil
	case "npc_teleport":
		return SeqStep{Action: SeqNPCTeleport, Scene: j.Scene, NPC: j.NPC, TargetX: j.X, TargetY: j.Y}, nil
	case "npc_move":
		return SeqStep{
			Action:   SeqNPCMove,
			Scene:    j.Scene,
			NPC:      j.NPC,
			TargetX:  j.X,
			TargetY:  j.Y,
			Duration: j.Seconds,
		}, nil
	}
	return SeqStep{}, fmt.Errorf("unknown step type %q", j.Type)
}

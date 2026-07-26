package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type hotspotJSON struct {
	Bounds      boundsJSON `json:"bounds"`
	TargetScene string     `json:"targetScene,omitempty"`
	Name        string     `json:"name"`
	Arrow       string     `json:"arrow"`
	// WalkX/WalkY (2026-07-24 #2/#31): optional explicit door-walk anchor
	// (PP CENTER coords) for scripted exit walks — the room exits derive a
	// bad target from bounds alone (right-wall doors have y≈100 bounds).
	// 0 = unset; WalkY 0 with WalkX set keeps PP's current row.
	WalkX float64 `json:"walkX,omitempty"`
	WalkY float64 `json:"walkY,omitempty"`
}

// npcOverrideJSON (2026-07-24, user: "why hardcoded and not in the JSON?"):
// per-scene NPC placement knobs keyed by the factory id from `npcs`. The Go
// constructors stay the defaults; a scene JSON override wins, so positions
// tuned per scene live in ONE place next to the scene's other geometry.
// footY = the sprite's bottom line (bounds.Y is derived); footYDark = the
// bottom line applyCampMood switches to when the camp mood darkens (the
// day-2/3 office art draws the desk higher).
type npcOverrideJSON struct {
	ID        string `json:"id"`
	FootY     int32  `json:"footY,omitempty"`
	FootYDark int32  `json:"footYDark,omitempty"`
}

type walkSegmentJSON struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

type sceneDef struct {
	Name                   string            `json:"name"`
	Background             string            `json:"background"`
	BackgroundFrames       int               `json:"backgroundFrames,omitempty"`
	BackgroundFrameSeconds float64           `json:"backgroundFrameSeconds,omitempty"`
	SpawnX                 float64           `json:"spawnX"`
	SpawnY                 float64           `json:"spawnY"`
	MinY                   float64           `json:"minY"`
	MaxY                   float64           `json:"maxY"`
	MusicPath              string            `json:"musicPath"`
	CharacterScale         float64           `json:"characterScale"`
	NPCs                   []string          `json:"npcs"`
	NPCOverrides           []npcOverrideJSON `json:"npcOverrides,omitempty"`
	Hotspots               []hotspotJSON     `json:"hotspots"`
	Blockers               []boundsJSON      `json:"blockers"`
	FootBlockers           []boundsJSON      `json:"footBlockers"`
	WalkSegments           []walkSegmentJSON `json:"walkSegments"`
}

type sceneConfigStore struct {
	defs map[string]sceneDef
}

func newSceneConfigStore(dir string) *sceneConfigStore {
	store := &sceneConfigStore{
		defs: make(map[string]sceneDef),
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Warning: could not read scene config directory %s: %v\n", dir, err)
		return store
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Warning: could not read scene config %s: %v\n", path, err)
			continue
		}

		var def sceneDef
		if err := json.Unmarshal(data, &def); err != nil {
			fmt.Printf("Warning: could not parse scene config %s: %v\n", path, err)
			continue
		}

		store.defs[def.Name] = def
		fmt.Printf("Loaded scene config: %s\n", def.Name)
	}

	return store
}

// GetDef returns a scene definition by name
func (s *sceneConfigStore) GetDef(name string) (sceneDef, bool) {
	def, ok := s.defs[name]
	return def, ok
}

// parseArrow converts arrow string to arrowDir
func parseArrow(s string) arrowDir {
	switch s {
	case "left":
		return arrowLeft
	case "right":
		return arrowRight
	case "up":
		return arrowUp
	case "down":
		return arrowDown
	case "downRight":
		return arrowDownRight
	default:
		return arrowNone
	}
}

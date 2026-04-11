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
}

type walkSegmentJSON struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

type sceneDef struct {
	Name         string            `json:"name"`
	Background   string            `json:"background"`
	SpawnX       float64           `json:"spawnX"`
	SpawnY       float64           `json:"spawnY"`
	MinY         float64           `json:"minY"`
	MaxY         float64           `json:"maxY"`
	MusicPath    string            `json:"musicPath"`
	NPCs         []string          `json:"npcs"`
	Hotspots     []hotspotJSON     `json:"hotspots"`
	Blockers     []boundsJSON      `json:"blockers"`
	WalkSegments []walkSegmentJSON `json:"walkSegments"`
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

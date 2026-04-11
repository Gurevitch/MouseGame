package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type gridSize struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

type boundsJSON struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
}

type npcDef struct {
	Name           string     `json:"name"`
	IdleSprite     string     `json:"idleSprite"`
	TalkSprite     string     `json:"talkSprite"`
	IdleGrid       gridSize   `json:"idleGrid"`
	TalkGrid       gridSize   `json:"talkGrid"`
	Bounds         boundsJSON `json:"bounds"`
	Dialog         string     `json:"dialog"`
	TalkFrameSpeed float64    `json:"talkFrameSpeed"`
	BobAmount      float64    `json:"bobAmount"`
	Silent         bool       `json:"silent"`

	StrangeIdleSprite string   `json:"strangeIdleSprite,omitempty"`
	StrangeTalkSprite string   `json:"strangeTalkSprite,omitempty"`
	StrangeIdleGrid   gridSize `json:"strangeIdleGrid,omitempty"`
	StrangeTalkGrid   gridSize `json:"strangeTalkGrid,omitempty"`
}

type npcConfigStore struct {
	defs map[string]npcDef
}

func newNPCConfigStore(dir string) *npcConfigStore {
	store := &npcConfigStore{
		defs: make(map[string]npcDef),
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Warning: could not read NPC config directory %s: %v\n", dir, err)
		return store
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Warning: could not read NPC config %s: %v\n", path, err)
			continue
		}

		var file struct {
			NPCs map[string]npcDef `json:"npcs"`
		}
		if err := json.Unmarshal(data, &file); err != nil {
			fmt.Printf("Warning: could not parse NPC config %s: %v\n", path, err)
			continue
		}

		for id, def := range file.NPCs {
			store.defs[id] = def
		}
		fmt.Printf("Loaded NPC config: %s (%d NPCs)\n", e.Name(), len(file.NPCs))
	}

	return store
}

// GetDef returns an NPC definition by ID
func (s *npcConfigStore) GetDef(id string) (npcDef, bool) {
	def, ok := s.defs[id]
	return def, ok
}

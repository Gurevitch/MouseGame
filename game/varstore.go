package game

import (
	"encoding/json"
	"fmt"
	"os"
)

// VarStore manages game state variables at three scopes:
//   - Game: persists across the entire game (e.g., parisUnlocked)
//   - Chapter: persists within a chapter/day (e.g., metKids, day)
//   - Scene: resets on scene change (e.g., dialogPlayed)
type VarStore struct {
	Game    map[string]int `json:"game"`
	Chapter map[string]int `json:"chapter"`
	Scene   map[string]int `json:"scene"`
}

func newVarStore() *VarStore {
	return &VarStore{
		Game:    make(map[string]int),
		Chapter: make(map[string]int),
		Scene:   make(map[string]int),
	}
}

// Get returns a variable value from the specified scope. Returns 0 if not set.
func (vs *VarStore) Get(scope, name string) int {
	switch scope {
	case "game":
		return vs.Game[name]
	case "chapter":
		return vs.Chapter[name]
	case "scene":
		return vs.Scene[name]
	}
	return 0
}

// Set sets a variable value in the specified scope.
func (vs *VarStore) Set(scope, name string, value int) {
	switch scope {
	case "game":
		vs.Game[name] = value
	case "chapter":
		vs.Chapter[name] = value
	case "scene":
		vs.Scene[name] = value
	}
}

// GetBool returns true if the variable is non-zero.
func (vs *VarStore) GetBool(scope, name string) bool {
	return vs.Get(scope, name) != 0
}

// SetBool sets a variable to 1 (true) or 0 (false).
func (vs *VarStore) SetBool(scope, name string, value bool) {
	if value {
		vs.Set(scope, name, 1)
	} else {
		vs.Set(scope, name, 0)
	}
}

// Inc increments a variable by 1 and returns the new value.
func (vs *VarStore) Inc(scope, name string) int {
	v := vs.Get(scope, name) + 1
	vs.Set(scope, name, v)
	return v
}

// ResetScene clears all scene-scoped variables.
func (vs *VarStore) ResetScene() {
	vs.Scene = make(map[string]int)
}

// ResetChapter clears all chapter-scoped variables.
func (vs *VarStore) ResetChapter() {
	vs.Chapter = make(map[string]int)
}

// Save serializes the VarStore to a JSON file.
func (vs *VarStore) Save(path string) error {
	data, err := json.MarshalIndent(vs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal varstore: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Load deserializes a VarStore from a JSON file.
func (vs *VarStore) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read varstore: %w", err)
	}
	return json.Unmarshal(data, vs)
}

// Dump prints all variables for debugging.
func (vs *VarStore) Dump() {
	fmt.Println("=== VarStore ===")
	fmt.Println("Game:", vs.Game)
	fmt.Println("Chapter:", vs.Chapter)
	fmt.Println("Scene:", vs.Scene)
}

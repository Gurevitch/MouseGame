package game

import (
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type itemDef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Texture     string `json:"texture"`
	Description string `json:"description"`
}

type itemRegistry struct {
	defs     map[string]itemDef
	renderer *sdl.Renderer
}

func newItemRegistry(renderer *sdl.Renderer, path string) *itemRegistry {
	reg := &itemRegistry{
		defs:     make(map[string]itemDef),
		renderer: renderer,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Warning: could not load item registry %s: %v\n", path, err)
		return reg
	}

	var file struct {
		Items []itemDef `json:"items"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		fmt.Printf("Warning: could not parse item registry %s: %v\n", path, err)
		return reg
	}

	for _, item := range file.Items {
		reg.defs[item.ID] = item
	}
	fmt.Printf("Loaded %d items from registry\n", len(reg.defs))
	return reg
}

// createItem creates an inventoryItem from the registry by ID
func (reg *itemRegistry) createItem(id string) *inventoryItem {
	def, ok := reg.defs[id]
	if !ok {
		fmt.Printf("Warning: item '%s' not found in registry\n", id)
		return nil
	}

	tex, w, h := engine.SafeTextureFromPNGKeyed(reg.renderer, def.Texture)
	if tex != nil {
		tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	}

	return &inventoryItem{
		name:  def.Name,
		tex:   tex,
		srcW:  w,
		srcH:  h,
		desc:  def.Description,
		owner: "player",
	}
}

// getDef returns the item definition without creating a texture
func (reg *itemRegistry) getDef(id string) (itemDef, bool) {
	def, ok := reg.defs[id]
	return def, ok
}

package game

import (
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type itemDef struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Texture     string  `json:"texture"`
	Description string  `json:"description"`
	IconScale   float64 `json:"iconScale"` // <1 shrinks an oversized bag icon (#17); 0 = full fit
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

// idForName resolves a display name ("Cafe au Lait") back to its registry id
// ("cafe_au_lait") by scanning the defs. 2026-08-08 #8: the save file stores
// DISPLAY names, and the old hardcoded 7-name switch silently dropped every
// other item from the bag on load ("Card", "Confiture", "Coin", ...). Falls
// back to the input so ids saved by future formats still resolve.
func (reg *itemRegistry) idForName(name string) string {
	for id, def := range reg.defs {
		if def.Name == name {
			return id
		}
	}
	return name
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

	// Content box of the artwork (post color-key) so the bag + held cursor
	// center the item by its ART, not its canvas (2026-06-11 #37).
	cbX, cbY, cbW, cbH := engine.ContentBoxKeyed(def.Texture)

	return &inventoryItem{
		name:      def.Name,
		tex:       tex,
		srcW:      w,
		srcH:      h,
		desc:      def.Description,
		owner:     "player",
		iconScale: def.IconScale,
		cbX:       cbX, cbY: cbY, cbW: cbW, cbH: cbH,
	}
}

// getDef returns the item definition without creating a texture
func (reg *itemRegistry) getDef(id string) (itemDef, bool) {
	def, ok := reg.defs[id]
	return def, ok
}

package game

import (
	"fmt"
	"image/color"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// buildSceneFromDef constructs a *scene from a JSON-authored sceneDef.
// It owns the static, declarative parts of the scene (background, spawn,
// bounds, NPCs, hotspots, blockers, walk segments, music). Procedural
// additions (particles, glow effects) are layered on by scene-specific
// decorators called from newSceneManager after this returns.
//
// Unknown NPC ids are skipped. Missing background PNGs use a gradient
// placeholder so JSON-authored scenes can load before their art exists.
func buildSceneFromDef(renderer *sdl.Renderer, def sceneDef) *scene {
	s := &scene{
		name:           def.Name,
		bg:             newPNGBackgroundOr(renderer, def.Background, color.NRGBA{R: 140, G: 110, B: 85, A: 255}),
		npcs:           spawnNPCs(renderer, def.NPCs),
		spawnX:         def.SpawnX,
		spawnY:         def.SpawnY,
		minY:           def.MinY,
		maxY:           def.MaxY,
		musicPath:      def.MusicPath,
		characterScale: def.CharacterScale,
	}

	for _, h := range def.Hotspots {
		s.hotspots = append(s.hotspots, hotspot{
			bounds:      sdl.Rect{X: h.Bounds.X, Y: h.Bounds.Y, W: h.Bounds.W, H: h.Bounds.H},
			name:        h.Name,
			targetScene: h.TargetScene,
			arrow:       parseArrow(h.Arrow),
		})
	}
	for _, b := range def.Blockers {
		s.blockers = append(s.blockers, sdl.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H})
	}
	for _, w := range def.WalkSegments {
		s.walkSegments = append(s.walkSegments, walkSegment{
			x1: w.X1, y1: w.Y1, x2: w.X2, y2: w.Y2,
		})
	}
	return s
}

// loadSceneFromJSON pulls a scene by name from the config store and builds
// it via buildSceneFromDef. Returns nil and logs if the def is missing, so
// newSceneManager can fall through to the hardcoded path during the
// phase-by-phase migration.
func (sm *sceneManager) loadSceneFromJSON(renderer *sdl.Renderer, store *sceneConfigStore, name string) *scene {
	def, ok := store.GetDef(name)
	if !ok {
		fmt.Printf("scene_loader: no JSON def for %q, falling back\n", name)
		return nil
	}
	return buildSceneFromDef(renderer, def)
}

// Compile-time: keep engine package referenced in case buildSceneFromDef
// grows features that need it and linters complain about unused imports.
var _ = engine.ScreenWidth

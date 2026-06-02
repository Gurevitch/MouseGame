package game

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// Packed-atlas system removed 2026-06-02 (user request).
//
// We used to support a "<name>.png + <name>.json" packed atlas under
// assets/sprites/ (built by tools/pack_atlas.py) that bundled every animation
// for a character into a single texture, and the loaders below tried it first
// before falling back to the per-animation PNGs. In practice only Marcus ever
// had an atlas on disk, and its presence overrode freshly-regenerated PNGs
// until someone re-ran the packer — a confusing footgun. Every character now
// loads directly from its per-animation PNGs, so dropping a regenerated sheet
// in takes effect on the next run with no packing step.

// applyKidAtlasOrFallback wires the four canonical kid animations
// (idle / talk / strange_idle / strange_talk) from the per-kid PNGs under
// assets/images/locations/camp/npc/kids/<name>/. Missing sheets are left nil
// so a kid who only has idle/talk still renders.
//
// Color-key tolerance: uses loadNPCGridConnected (connected-edge key, tol=8),
// which clears the soft white sheet background without eating the saturated
// shirt colors or the eye-whites/teeth that match the background.
func applyKidAtlasOrFallback(renderer *sdl.Renderer, n *npc, atlasName string) {
	base := "assets/images/locations/camp/npc/kids/" + atlasName + "/npc_" + atlasName
	loadIfExists := func(path string) []npcFrame {
		if _, err := os.Stat(path); err != nil {
			return nil
		}
		return loadNPCGridConnected(renderer, path, 8, 2)
	}
	n.idleGrid = loadIfExists(base + "_idle.png")
	n.talkGrid = loadIfExists(base + "_talk.png")
	n.strangeIdle = loadIfExists(base + "_strange_idle.png")
	n.strangeTalk = loadIfExists(base + "_strange_talk.png")
}

// applyNPCAtlas is retained as a no-op shim now that the packed-atlas system is
// gone. Adult-NPC constructors call it as `if !applyNPCAtlas(...) { <load
// PNGs> }`; returning false here makes the per-PNG fallback always run.
func applyNPCAtlas(_ *sdl.Renderer, _ *npc, _ string) bool { return false }

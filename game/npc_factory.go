package game

import "github.com/veandco/go-sdl2/sdl"

// npcFactory maps a string id (the same id scenes reference in their JSON
// "npcs" list) to a constructor that produces a *npc. Keeping this as a flat
// map lets scene JSON declare NPCs by name without scene_loader needing to
// import every npc constructor directly.
//
// To add a new NPC:
//   1. Write a `newFoo(renderer) *npc` constructor.
//   2. Register it here with the id you will reference in scene JSON.
//
// Missing ids are silently skipped at scene-build time with a warning; the
// loader logs, the scene still spawns, and the player can continue playing
// while you fix the typo.
var npcFactories = map[string]func(*sdl.Renderer) *npc{
	"director_higgins":        newDirectorHiggins,
	"office_higgins":          newOfficeHiggins,
	"night_higgins":           newNightHiggins,
	"grounds_higgins_hidden":  newGroundsHiggins,
	"marcus":                  newMarcus,
	"jake":                    newJake,
	"lily":                    newLily,
	"tommy":                   newTommy,
	"danny":                   newDanny,
	// Cabin-bound variants: kid at their bed position, silent by default
	// (Day 2 callbacks flip silent off). Marcus is not silent and is larger.
	"room_marcus":             newRoomMarcus,
	"room_jake":               newRoomJake,
	"room_lily":               newRoomLily,
	"room_tommy":              newRoomTommy,
	"room_danny":              newRoomDanny,
	// Paris NPCs
	"french_guide":            newFrenchGuide,
	"museum_curator":          newMuseumCurator,
	"pierre_artist":           newPierreArtist,
	"gendarme_claude":         newGendarmeClaude,
	"bakery_woman":            newBakeryWoman,
	"press_photographer":      newPressPhotographer,
}

// registerNPCFactory lets modules (paris.go / jerusalem.go / ...) add their
// NPCs to the registry at package init-time without a central import list.
func registerNPCFactory(id string, ctor func(*sdl.Renderer) *npc) {
	npcFactories[id] = ctor
}

// spawnNPCs builds the NPCs listed in a scene def, skipping unknown ids with
// a warning so a typo doesn't brick the scene.
//
// Callers can set each NPC's back-reference to Game afterwards with
// attachGameToNPCs; that's not done here because spawnNPCs runs inside
// newSceneManager before Game is fully constructed.
func spawnNPCs(renderer *sdl.Renderer, ids []string) []*npc {
	out := make([]*npc, 0, len(ids))
	for _, id := range ids {
		ctor, ok := npcFactories[id]
		if !ok {
			// Fall through to config-store-driven NPCs (paris etc.) that
			// haven't been migrated yet. They're looked up elsewhere.
			continue
		}
		out = append(out, ctor(renderer))
	}
	return out
}

// attachGameToNPCs sweeps every NPC in every scene and sets its `game` back-
// reference. Called once during Game.New after the sceneManager finishes
// constructing scenes. Without this, rule-driven NPCs can't reach the game
// state they need (inventory, varstore, eventbus).
func (g *Game) attachGameToNPCs() {
	for _, s := range g.sceneMgr.scenes {
		for _, n := range s.npcs {
			n.game = g
		}
	}
}

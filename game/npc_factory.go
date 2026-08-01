package game

import "github.com/veandco/go-sdl2/sdl"

// npcFactory maps a string id (the same id scenes reference in their JSON
// "npcs" list) to a constructor that produces a *npc. Keeping this as a flat
// map lets scene JSON declare NPCs by name without scene_loader needing to
// import every npc constructor directly.
//
// To add a new NPC:
//  1. Write a `newFoo(renderer) *npc` constructor.
//  2. Register it here with the id you will reference in scene JSON.
//
// Missing ids are silently skipped at scene-build time with a warning; the
// loader logs, the scene still spawns, and the player can continue playing
// while you fix the typo.
var npcFactories = map[string]func(*sdl.Renderer) *npc{
	"director_higgins":       newDirectorHiggins,
	"office_higgins":         newOfficeHiggins,
	"night_higgins":          newNightHiggins,
	"grounds_higgins_hidden": newGroundsHiggins,
	"marcus":                 newMarcus,
	"jake":                   newJake,
	"lily":                   newLily,
	"tommy":                  newTommy,
	"danny":                  newDanny,
	// Cabin-bound variants: kid at their bed position, silent by default
	// (Day 2 callbacks flip silent off). Marcus is not silent and is larger.
	"room_marcus": newRoomMarcus,
	"room_jake":   newRoomJake,
	"room_lily":   newRoomLily,
	"room_tommy":  newRoomTommy,
	"room_danny":  newRoomDanny,
	"lily_lake":   newLakeLily,
	// Paris NPCs
	"french_guide":       newFrenchGuide,
	"museum_curator":     newMuseumCurator,
	"pierre_artist":      newPierreArtist,
	"pigeon_lady":        newPigeonLady,
	"gendarme_claude":    newGendarmeClaude,
	"bakery_woman":       newBakeryWoman,
	"press_photographer": newPressPhotographer,
	// Cafe patrons (paris_bakery interior). Henri carries the coffee-jam
	// quest beat; the other 5 are flavor.
	"cafe_patron_yvette":  newCafePatronYvette,
	"cafe_patron_bernard": newCafePatronBernard,
	"cafe_patron_camille": newCafePatronCamille,
	"cafe_patron_henri":   newCafePatronHenri,
	"cafe_patron_lucien":  newCafePatronLucien,
	"cafe_patron_elise":   newCafePatronElise,
	// Jerusalem NPCs — positions baked into the closures so the factory
	// signature stays func(*sdl.Renderer)*npc (no x param).
	"jer_shimon":           func(r *sdl.Renderer) *npc { return newShimon(r, 760) },
	"jer_bagel_seller":     func(r *sdl.Renderer) *npc { return newBagelSeller(r, 250) },
	"jer_praying_man":      func(r *sdl.Renderer) *npc { return newPrayingMan(r, 470) },
	"jer_wall_kid":         func(r *sdl.Renderer) *npc { return newWallKid(r, 880) },
	"jer_spice_seller":     func(r *sdl.Renderer) *npc { return newSpiceSeller(r, 249) },
	"jer_coffee_seller":    func(r *sdl.Renderer) *npc { return newCoffeeSeller(r, 510) },
	"jer_antiques_kid":     func(r *sdl.Renderer) *npc { return newAntiquesKid(r, 900) },
	"jer_antiques_old_man": func(r *sdl.Renderer) *npc { return newAntiquesOldMan(r, 1040) },
	// Japan / Kyoto NPCs
	"jp_gary":       newTouristKyoto,
	"jp_hiro":       newRamenSeller,
	"jp_kenji":      newKenjiStudent,
	"jp_takeshi":    newTakeshi,
	"jp_obachan":    newObachan,
	"jp_kiku":       newDresser,
	"jp_tea_master": newTeaMaster,
}

// spawnNPCs builds the NPCs listed in a scene def, skipping unknown ids with
// a warning so a typo doesn't brick the scene.
//
// Callers can set each NPC's back-reference to Game afterwards with
// attachGameToNPCs; that's not done here because spawnNPCs runs inside
// newSceneManager before Game is fully constructed.
func spawnNPCs(renderer *sdl.Renderer, ids []string, overrides []npcOverrideJSON) []*npc {
	out := make([]*npc, 0, len(ids))
	for _, id := range ids {
		ctor, ok := npcFactories[id]
		if !ok {
			// Fall through to config-store-driven NPCs (paris etc.) that
			// haven't been migrated yet. They're looked up elsewhere.
			continue
		}
		n := ctor(renderer)
		// 2026-07-24 (user): scene-JSON placement overrides beat the
		// constructor defaults — geometry lives with the scene's JSON.
		for _, ov := range overrides {
			if ov.ID != id {
				continue
			}
			if ov.FootY > 0 {
				n.bounds.Y = ov.FootY - n.bounds.H
				n.footYDay = ov.FootY
			}
			if ov.FootYDark > 0 {
				n.footYDark = ov.FootYDark
			}
			break
		}
		out = append(out, n)
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

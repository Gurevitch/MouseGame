package engine

import "testing"

// Exercises gap-based frame detection (contentGridRects) against the real
// player sheets. Informational: logs whether each sheet splits cleanly by its
// gaps or falls back to proportional slicing — fails only on panic/misshape.
func TestContentGridRects_RealSheets(t *testing.T) {
	cases := []struct {
		path string
		cols int
		rows int
	}{
		{"../assets/images/player/PP idle front.png", 8, 2},
		{"../assets/images/player/PP idle side.png", 8, 2},
		{"../assets/images/player/PP idle back.png", 8, 2},
		{"../assets/images/player/PP walk front.png", 8, 2},
		{"../assets/images/player/PP walk back.png", 8, 2},
		{"../assets/images/player/PP walk left.png", 8, 1},
		{"../assets/images/player/PP talk front.png", 8, 2},
		{"../assets/images/player/PP talk side.png", 8, 2},
		{"../assets/images/player/PP grab.png", 8, 2},
		{"../assets/images/player/PP celebrate.png", 8, 2},
		{"../assets/images/player/PP sneak examine.png", 8, 2},
		{"../assets/images/player/PP sneak use.png", 8, 2},
		{"../assets/images/player/PP receive map.png", 4, 2},
		{"../assets/images/player/PP grab flower.png", 6, 1},
		{"../assets/images/player/PP grab rolling pin.png", 6, 1},
		{"../assets/images/player/PP get bagguette.png", 8, 1},
		{"../assets/images/player/PP get jam.png", 8, 1},
		// ---- Camp NPCs ----
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 6, 1},
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png", 6, 2},
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png", 6, 2},
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_give_map.png", 6, 2},
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_shout.png", 8, 2},
		{"../assets/images/locations/camp/npc/higgins/npc_director_higgins_walk_back.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/tommy/npc_tommy_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/tommy/npc_tommy_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/jake/npc_jake_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/jake/npc_jake_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/jake/npc_jake_strange_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/jake/npc_jake_strange_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/lily/npc_lily_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/lily/npc_lily_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_alt.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/danny/npc_danny_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/danny/npc_danny_talk.png", 8, 2},
		// ---- Paris NPCs ----
		{"../assets/images/locations/paris/npc/outside/npc_madame_colette_idle.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_art_vendor.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_security_guard.png", 6, 2},
		{"../assets/images/locations/paris/npc/outside/npc_press_photographer.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_give.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_pigeon_lands.png", 8, 1},
		{"../assets/images/locations/paris/npc/museum/npc_museum_curator_idle.png", 8, 1},
		{"../assets/images/locations/paris/npc/museum/npc_museum_curator_talk.png", 8, 1},
		{"../assets/images/locations/paris/npc/museum/npc_beaumont_give.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/npc_madame_poulain_idle.png", 8, 2},
		{"../assets/images/locations/paris/npc/coffee/npc_madame_poulain_talk.png", 8, 2},
		{"../assets/images/locations/paris/npc/coffee/npc_madame_poulain_work.png", 8, 2},
		{"../assets/images/locations/paris/npc/coffee/npc_madame_poulain_give.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/cafe_patron_yvette.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/cafe_patron_bernard_idle.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/cafe_patron_camille.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/cafe_patron_henri.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/cafe_patron_lucien.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/npc_henri_give_jam.png", 6, 1},
		{"../assets/images/locations/paris/npc/coffee/npc_camille_sketching.png", 8, 1},
	}
	for _, c := range cases {
		img, err := loadPNG(c.path)
		if err != nil {
			t.Logf("skip %s: %v", c.path, err)
			continue
		}
		applyColorKey(img)
		rects := contentGridRects(img, c.cols, c.rows)
		if rects == nil {
			t.Logf("%-70s -> fallback (no clean %dx%d gap split)", c.path, c.cols, c.rows)
			continue
		}
		if len(rects) != c.rows || len(rects[0]) != c.cols {
			t.Errorf("%s: wrong shape %dx%d", c.path, len(rects), len(rects[0]))
			continue
		}
		t.Logf("%-70s -> GAP-DETECTED %dx%d, cell(0,0)=%v cell(0,1)=%v", c.path, c.rows, c.cols, rects[0][0], rects[0][1])
	}
}

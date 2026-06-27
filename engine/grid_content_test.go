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
		{"../assets/images/player/PP write note.png", 6, 1},
		{"../assets/images/player/PP put note in wall.png", 6, 1},
		{"../assets/images/player/PP get bagguette.png", 8, 1},
		{"../assets/images/player/PP get jam.png", 8, 1},
		{"../assets/images/player/PP give flower.png", 8, 1},
		{"../assets/images/player/PP give rolling pin.png", 8, 1},
		{"../assets/images/player/PP give baguette.png", 8, 1},
		{"../assets/images/player/PP give confiture.png", 8, 1},
		{"../assets/images/player/PP give coffee.png", 8, 1},
		{"../assets/images/player/PP give heel.png", 8, 1},
		{"../assets/images/player/PP give pencil.png", 8, 1},
		{"../assets/images/player/PP give sketch.png", 8, 1},
		{"../assets/images/player/PP give postcard.png", 8, 1},
		{"../assets/images/player/PP jump back.png", 8, 1},
		{"../assets/images/player/PP pull map.png", 8, 1},
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
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_idle_day.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_idle_night.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_talk.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_alt.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_going_to_sleep.png", 8, 1},
		{"../assets/images/locations/camp/npc/kids/marcus/npc_marcus_sleeping.png", 8, 1},
		{"../assets/images/locations/camp/npc/kids/danny/npc_danny_idle.png", 8, 2},
		{"../assets/images/locations/camp/npc/kids/danny/npc_danny_talk.png", 8, 2},
		// ---- Paris NPCs ----
		{"../assets/images/locations/paris/npc/outside/npc_madame_colette_idle.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_art_vendor.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_idle.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_talk.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pigeon_lady_idle.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pigeon_lady_give.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_security_guard.png", 6, 2},
		{"../assets/images/locations/paris/npc/outside/npc_press_photographer.png", 8, 2},
		{"../assets/images/locations/paris/npc/outside/npc_press_photographer_idle.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/biker.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/ambient_accordion_player.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/ambient_crumb_lady.png", 8, 1},
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
		// ---- Jerusalem NPCs (grids per game/jerusalem.go) ----
		{"../assets/images/locations/jerusalem/npc/wall/npc_shimon.png", 6, 2},
		{"../assets/images/locations/jerusalem/npc/wall/npc_shimon_give.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_bagel_seller.png", 6, 2},
		{"../assets/images/locations/jerusalem/npc/wall/npc_bagel_seller_give.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_praying_man_idle.png", 8, 2},
		{"../assets/images/locations/jerusalem/npc/wall/npc_praying_man_talk.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_praying_man_give.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_wall_kid_idle.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_wall_kid_talk.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/npc_spice_seller_idle.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/npc_spice_seller_talk.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/npc_spice_seller_give.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/npc_coffee_seller_idle.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/npc_coffee_seller_talk.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/npc_coffee_seller_give.png", 8, 1},
		// ---- 2026-06-24 bug-sweep new sheets (all 8x1) ----
		{"../assets/images/locations/camp/npc/kids/jake/npc_jake_falling_sleep.png", 8, 1},
		{"../assets/images/locations/camp/npc/kids/jake/npc_jake_sleeping.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/kid_antique_idle.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/kid_antique_idle_alter.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/kid_antique_speak.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/market/grandpa_idle.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_praying_man_give_paper.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_shimon_give_coin.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/npc_shimon_give_pen.png", 8, 1},
		{"../assets/images/locations/jerusalem/npc/wall/praying_man.png", 4, 1},
		{"../assets/images/locations/jerusalem/npc/wall/praying_man2.png", 4, 1},
		{"../assets/images/locations/paris/npc/coffee/npc_camille_sketching_portrait.png", 8, 1},
		{"../assets/images/locations/paris/npc/coffee/npc_madame_poulain_give_coffee.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_get_baguette.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_get_jam.png", 8, 1},
		{"../assets/images/locations/paris/npc/outside/npc_pierre_give_pass.png", 8, 1},
		{"../assets/images/player/PP_get_baguette_back.png", 8, 1},
		{"../assets/images/player/PP_get_coffee_back.png", 8, 1},
		{"../assets/images/player/PP_give_rolling_pin_back.png", 8, 1},
		// ---- Japan chapter sheets (NPCs 8x1; leaf 3x1) ----
		{"../assets/images/locations/camp/npc/kids/lily/npc_lily_sad_idle.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_hiro_idle.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_hiro_talk.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_gary_idle.png", 6, 2},
		{"../assets/images/locations/japan/npc/npc_gary_idle_oposite_book.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_gary_talk_oposite_book.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_gary_flip_his_book.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_geisha_idle.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_obachan_idle.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_obachan_talk.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_tea_master_idle.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_tea_master_talk.png", 8, 1},
		{"../assets/images/locations/japan/npc/npc_kenji_idle.png", 8, 1},
		{"../assets/images/locations/japan/props/leaf_fall.png", 3, 1},
		{"../assets/images/player/PP_kimono_spin.png", 8, 1},
		{"../assets/images/player/PP_spin_to_sit.png", 8, 1},
		{"../assets/images/player/PP_tea_ceremony.png", 8, 1},
		{"../assets/images/player/PP_sit_idle.png", 8, 1},
		{"../assets/images/player/PP_sit_talk.png", 8, 1},
		{"../assets/images/locations/camp/npc/higgins/npc_director_front_walk.png", 8, 1},
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

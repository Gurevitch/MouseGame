// jitter_audit scans sprite sheets with their declared grids and measures
// per-cell content drift: foot-line (bottom Y) movement, horizontal center
// drift, content-height pumping, empty cells, and cross-cell bleed. These are
// the root causes of in-game "jitter"/"jumping" sprites.
//
// Run: go run ./tools/jitter_audit
package main

import (
	"fmt"
	"image/png"
	"math"
	"os"
	"sort"
)

type sheet struct {
	label string
	path  string
	cols  int
	rows  int
	row   int // -1 = all rows
}

var sheets = []sheet{
	// ---- Player (game/player.go) ----
	{"PP walk side", "assets/images/player/PP walk left.png", 8, 1, -1},
	{"PP walk front (row0)", "assets/images/player/PP walk front.png", 8, 2, 0},
	{"PP walk back", "assets/images/player/PP walk back.png", 8, 2, -1},
	{"PP idle front", "assets/images/player/PP idle front.png", 8, 2, -1},
	{"PP idle side", "assets/images/player/PP idle side.png", 8, 2, -1},
	{"PP idle back", "assets/images/player/PP idle back.png", 8, 2, -1},
	{"PP talk front", "assets/images/player/PP talk front.png", 8, 2, -1},
	{"PP talk side", "assets/images/player/PP talk side.png", 8, 2, -1},
	{"PP grab", "assets/images/player/PP grab.png", 8, 2, -1},
	{"PP celebrate", "assets/images/player/PP celebrate.png", 8, 2, -1},
	{"PP sneak examine", "assets/images/player/PP sneak examine.png", 8, 2, -1},
	{"PP sneak use", "assets/images/player/PP sneak use.png", 8, 2, -1},
	{"PP receive map", "assets/images/player/PP receive map.png", 4, 2, -1},
	{"PP grab flower", "assets/images/player/PP grab flower.png", 6, 1, -1},
	{"PP grab rolling pin", "assets/images/player/PP grab rolling pin.png", 6, 1, -1},
	{"PP put note in wall", "assets/images/player/PP put note in wall.png", 6, 1, -1},
	{"PP get baguette", "assets/images/player/PP get bagguette.png", 8, 1, -1},
	{"PP get jam", "assets/images/player/PP get jam.png", 8, 1, -1},
	{"PP give flower", "assets/images/player/PP give flower.png", 8, 1, -1},
	{"PP give rolling pin", "assets/images/player/PP give rolling pin.png", 8, 1, -1},
	{"PP give baguette", "assets/images/player/PP give baguette.png", 8, 1, -1},
	{"PP give confiture", "assets/images/player/PP give confiture.png", 8, 1, -1},
	{"PP give coffee", "assets/images/player/PP give coffee.png", 8, 1, -1},
	{"PP give heel", "assets/images/player/PP give heel.png", 8, 1, -1},
	{"PP give pencil", "assets/images/player/PP give pencil.png", 8, 1, -1},
	{"PP give sketch", "assets/images/player/PP give sketch.png", 8, 1, -1},
	{"PP give postcard", "assets/images/player/PP give postcard.png", 8, 1, -1},
	{"PP jump back", "assets/images/player/PP jump back.png", 8, 1, -1},
	{"PP pull map", "assets/images/player/PP pull map.png", 8, 1, -1},

	// ---- Higgins (grids per game/npc.go — JSON grids are NOT used) ----
	{"Higgins idle entrance (6x1)", "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 6, 1, -1},
	{"Higgins idle AS NIGHT loads it (7x1)", "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png", 7, 1, -1},
	{"Higgins talk (row0)", "assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png", 8, 2, 0},
	{"Higgins office idle (row0 of 6x2)", "assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png", 6, 2, 0},
	{"Higgins office talk (6x2)", "assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png", 6, 2, -1},
	{"Higgins give map", "assets/images/locations/camp/npc/higgins/npc_director_higgins_give_map.png", 6, 2, -1},
	{"Higgins shout", "assets/images/locations/camp/npc/higgins/npc_director_higgins_shout.png", 8, 2, -1},
	{"Higgins walk back", "assets/images/locations/camp/npc/higgins/npc_director_higgins_walk_back.png", 8, 2, -1},

	// ---- Camp kids (engine loads ALL kid states as 8x2, atlas.go:34) ----
	{"Tommy idle", "assets/images/locations/camp/npc/kids/tommy/npc_tommy_idle.png", 8, 2, -1},
	{"Tommy talk", "assets/images/locations/camp/npc/kids/tommy/npc_tommy_talk.png", 8, 2, -1},
	{"Tommy walk left (row1 of 8x3)", "assets/images/locations/camp/npc/kids/tommy/npc_tommy_walk_left.png", 8, 3, 1},
	{"Jake idle", "assets/images/locations/camp/npc/kids/jake/npc_jake_idle.png", 8, 2, -1},
	{"Jake talk", "assets/images/locations/camp/npc/kids/jake/npc_jake_talk.png", 8, 2, -1},
	{"Jake strange idle", "assets/images/locations/camp/npc/kids/jake/npc_jake_strange_idle.png", 8, 2, -1},
	{"Jake strange talk", "assets/images/locations/camp/npc/kids/jake/npc_jake_strange_talk.png", 8, 2, -1},
	{"Jake walk back", "assets/images/locations/camp/npc/kids/jake/npc_jake_walk_back.png", 8, 1, -1},
	{"Lily idle", "assets/images/locations/camp/npc/kids/lily/npc_lily_idle.png", 8, 2, -1},
	{"Lily talk", "assets/images/locations/camp/npc/kids/lily/npc_lily_talk.png", 8, 2, -1},
	{"Lily receive flower", "assets/images/locations/camp/npc/kids/lily/npc_lily_receive_flower.png", 8, 1, -1},
	{"Marcus idle", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_idle.png", 8, 2, -1},
	{"Marcus talk", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_talk.png", 8, 2, -1},
	{"Marcus strange idle (day)", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_idle_day.png", 8, 2, -1},
	{"Marcus strange idle (night)", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_idle_night.png", 8, 2, -1},
	{"Marcus strange talk", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_talk.png", 8, 2, -1},
	{"Marcus strange alt", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_alt.png", 8, 2, -1},
	{"Marcus going to sleep", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_going_to_sleep.png", 8, 1, -1},
	{"Marcus sleeping", "assets/images/locations/camp/npc/kids/marcus/npc_marcus_sleeping.png", 8, 1, -1},
	{"Danny idle", "assets/images/locations/camp/npc/kids/danny/npc_danny_idle.png", 8, 2, -1},
	{"Danny talk", "assets/images/locations/camp/npc/kids/danny/npc_danny_talk.png", 8, 2, -1},

	// ---- Paris street (grids per game/npc.go factories) ----
	{"Colette idle (8x2)", "assets/images/locations/paris/npc/outside/npc_madame_colette_idle.png", 8, 2, -1},
	{"Colette talk (8x2)", "assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png", 8, 2, -1},
	{"Pierre idle (art_vendor row0)", "assets/images/locations/paris/npc/outside/npc_art_vendor.png", 8, 2, 0},
	{"Pierre talk (art_vendor row1)", "assets/images/locations/paris/npc/outside/npc_art_vendor.png", 8, 2, 1},
	{"Pierre idle split", "assets/images/locations/paris/npc/outside/npc_pierre_idle.png", 8, 1, -1},
	{"Pierre talk split", "assets/images/locations/paris/npc/outside/npc_pierre_talk.png", 8, 1, -1},
	{"Pigeon lady idle", "assets/images/locations/paris/npc/outside/npc_pigeon_lady_idle.png", 8, 1, -1},
	{"Pigeon lady give", "assets/images/locations/paris/npc/outside/npc_pigeon_lady_give.png", 8, 1, -1},
	{"Claude idle (security_guard row0)", "assets/images/locations/paris/npc/outside/npc_security_guard.png", 6, 2, 0},
	{"Claude talk (security_guard row1)", "assets/images/locations/paris/npc/outside/npc_security_guard.png", 6, 2, 1},
	{"Nicolas idle (row0)", "assets/images/locations/paris/npc/outside/npc_press_photographer.png", 8, 2, 0},
	{"Nicolas talk (row1)", "assets/images/locations/paris/npc/outside/npc_press_photographer.png", 8, 2, 1},
	{"Nicolas idle split", "assets/images/locations/paris/npc/outside/npc_press_photographer_idle.png", 8, 1, -1},
	{"Paris biker", "assets/images/locations/paris/npc/outside/biker.png", 8, 1, -1},
	{"Ambient accordion player", "assets/images/locations/paris/npc/outside/ambient_accordion_player.png", 8, 1, -1},
	{"Ambient crumb lady", "assets/images/locations/paris/npc/outside/ambient_crumb_lady.png", 8, 1, -1},

	// ---- Paris bakery ----
	{"Poulain idle", "assets/images/locations/paris/npc/coffee/npc_madame_poulain_idle.png", 8, 2, -1},
	{"Poulain talk", "assets/images/locations/paris/npc/coffee/npc_madame_poulain_talk.png", 8, 2, -1},
	{"Poulain give", "assets/images/locations/paris/npc/coffee/npc_madame_poulain_give.png", 8, 1, -1},
	{"Poulain work (alt-idle 8x2)", "assets/images/locations/paris/npc/coffee/npc_madame_poulain_work.png", 8, 2, -1},
	{"Poulain bring baguette", "assets/images/locations/paris/npc/coffee/npc_madame_poulain_bring_bagguette.png", 8, 1, -1},
	{"Pierre give", "assets/images/locations/paris/npc/outside/npc_pierre_give.png", 8, 1, -1},
	{"Pierre pigeon lands", "assets/images/locations/paris/npc/outside/npc_pierre_pigeon_lands.png", 8, 1, -1},
	{"Beaumont give", "assets/images/locations/paris/npc/museum/npc_beaumont_give.png", 8, 1, -1},
	{"Henri give jam", "assets/images/locations/paris/npc/coffee/npc_henri_give_jam.png", 6, 1, -1},
	{"Patron Yvette", "assets/images/locations/paris/npc/coffee/cafe_patron_yvette.png", 8, 1, -1},
	{"Patron Yvette talking", "assets/images/locations/paris/npc/coffee/cafe_patron_yvette_talking.png", 8, 1, -1},
	{"Patron Bernard idle", "assets/images/locations/paris/npc/coffee/cafe_patron_bernard_idle.png", 8, 1, -1},
	{"Patron Bernard talking", "assets/images/locations/paris/npc/coffee/cafe_patron_bernard_talking.png", 8, 1, -1},
	{"Patron Camille", "assets/images/locations/paris/npc/coffee/cafe_patron_camille.png", 8, 1, -1},
	{"Patron Camille talking", "assets/images/locations/paris/npc/coffee/cafe_patron_camille_talking.png", 8, 1, -1},
	{"Patron Henri", "assets/images/locations/paris/npc/coffee/cafe_patron_henri.png", 8, 1, -1},
	{"Patron Henri talking", "assets/images/locations/paris/npc/coffee/cafe_patron_henri_talking.png", 8, 1, -1},
	{"Patron Lucien", "assets/images/locations/paris/npc/coffee/cafe_patron_lucien.png", 8, 1, -1},
	{"Patron Lucien talking", "assets/images/locations/paris/npc/coffee/cafe_patron_lucien_talking.png", 8, 1, -1},
	{"Camille sketching", "assets/images/locations/paris/npc/coffee/npc_camille_sketching.png", 8, 1, -1},
	{"Camille lost pencil", "assets/images/locations/paris/npc/coffee/cafe_patron_camille_lostpencil.png", 8, 1, -1},

	// ---- Louvre (npc.go:1549 loads both as 8x1 strips) ----
	{"Curator idle (8x1)", "assets/images/locations/paris/npc/museum/npc_museum_curator_idle.png", 8, 1, -1},
	{"Curator talk (8x1)", "assets/images/locations/paris/npc/museum/npc_museum_curator_talk.png", 8, 1, -1},

	// ---- 2026-06-24 bug-sweep new sheets (all 8x1) ----
	{"Jake falling asleep", "assets/images/locations/camp/npc/kids/jake/npc_jake_falling_sleep.png", 8, 1, -1},
	{"Jake sleeping", "assets/images/locations/camp/npc/kids/jake/npc_jake_sleeping.png", 8, 1, -1},
	{"Antiques kid idle", "assets/images/locations/jerusalem/npc/market/kid_antique_idle.png", 8, 1, -1},
	{"Antiques kid alt idle", "assets/images/locations/jerusalem/npc/market/kid_antique_idle_alter.png", 8, 1, -1},
	{"Antiques kid speak", "assets/images/locations/jerusalem/npc/market/kid_antique_speak.png", 8, 1, -1},
	{"Antiques grandpa idle", "assets/images/locations/jerusalem/npc/market/grandpa_idle.png", 8, 1, -1},
	{"Praying man give paper", "assets/images/locations/jerusalem/npc/wall/npc_praying_man_give_paper.png", 8, 1, -1},
	{"Shimon give coin", "assets/images/locations/jerusalem/npc/wall/npc_shimon_give_coin.png", 8, 1, -1},
	{"Shimon give pen", "assets/images/locations/jerusalem/npc/wall/npc_shimon_give_pen.png", 8, 1, -1},
	{"Wall worshipper 1 (4f sway)", "assets/images/locations/jerusalem/npc/wall/praying_man.png", 4, 1, -1},
	{"Wall worshipper 2 (4f sway)", "assets/images/locations/jerusalem/npc/wall/praying_man2.png", 4, 1, -1},
	{"Camille sketching portrait", "assets/images/locations/paris/npc/coffee/npc_camille_sketching_portrait.png", 8, 1, -1},
	{"Poulain give coffee", "assets/images/locations/paris/npc/coffee/npc_madame_poulain_give_coffee.png", 8, 1, -1},
	{"Pierre get baguette", "assets/images/locations/paris/npc/outside/npc_pierre_get_baguette.png", 8, 1, -1},
	{"Pierre get jam", "assets/images/locations/paris/npc/outside/npc_pierre_get_jam.png", 8, 1, -1},
	{"Pierre give pass", "assets/images/locations/paris/npc/outside/npc_pierre_give_pass.png", 8, 1, -1},
	{"PP get baguette back", "assets/images/player/PP_get_baguette_back.png", 8, 1, -1},
	{"PP get coffee back", "assets/images/player/PP_get_coffee_back.png", 8, 1, -1},
	{"PP give rolling pin back", "assets/images/player/PP_give_rolling_pin_back.png", 8, 1, -1},

	// ---- Japan chapter (NPCs 8x1; leaf 3x1) ----
	{"Lily sad idle", "assets/images/locations/camp/npc/kids/lily/npc_lily_sad_idle.png", 8, 1, -1},
	{"Hiro idle (ramen)", "assets/images/locations/japan/npc/npc_hiro_idle.png", 8, 1, -1},
	{"Hiro talk", "assets/images/locations/japan/npc/npc_hiro_talk.png", 8, 1, -1},
	{"Gary idle (correct book, 6x2)", "assets/images/locations/japan/npc/npc_gary_idle.png", 6, 2, -1},
	{"Gary idle (opposite book)", "assets/images/locations/japan/npc/npc_gary_idle_oposite_book.png", 8, 1, -1},
	{"Gary talk (opposite book)", "assets/images/locations/japan/npc/npc_gary_talk_oposite_book.png", 8, 1, -1},
	{"Gary flip book", "assets/images/locations/japan/npc/npc_gary_flip_his_book.png", 8, 1, -1},
	{"Geisha idle (Kiku)", "assets/images/locations/japan/npc/npc_geisha_idle.png", 8, 1, -1},
	{"Oba-chan idle", "assets/images/locations/japan/npc/npc_obachan_idle.png", 8, 1, -1},
	{"Oba-chan talk", "assets/images/locations/japan/npc/npc_obachan_talk.png", 8, 1, -1},
	{"Tea Master idle", "assets/images/locations/japan/npc/npc_tea_master_idle.png", 8, 1, -1},
	{"Tea Master talk", "assets/images/locations/japan/npc/npc_tea_master_talk.png", 8, 1, -1},
	{"Kenji idle", "assets/images/locations/japan/npc/npc_kenji_idle.png", 8, 1, -1},
	{"Leaf fall (3f)", "assets/images/locations/japan/props/leaf_fall.png", 3, 1, -1},
	{"PP kimono spin", "assets/images/player/PP_kimono_spin.png", 8, 1, -1},
	{"PP spin to sit", "assets/images/player/PP_spin_to_sit.png", 8, 1, -1},
	{"PP tea ceremony", "assets/images/player/PP_tea_ceremony.png", 8, 1, -1},
	{"PP sit idle", "assets/images/player/PP_sit_idle.png", 8, 1, -1},
	{"PP sit talk", "assets/images/player/PP_sit_talk.png", 8, 1, -1},
	{"Higgins front walk", "assets/images/locations/camp/npc/higgins/npc_director_front_walk.png", 8, 1, -1},
}

type cellStat struct {
	empty      bool
	minX, maxX int
	minY, maxY int
	touchesL   bool
	touchesR   bool
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func near(a, b uint32, tol int) bool {
	return absInt(int(a>>8)-int(b>>8)) <= tol
}

// componentAreas labels 4-connected foreground components in a cw×ch mask and
// returns their areas sorted descending.
func componentAreas(fg []bool, cw, ch int) []int {
	seen := make([]bool, len(fg))
	var areas []int
	var stack []int
	for i := range fg {
		if !fg[i] || seen[i] {
			continue
		}
		area := 0
		stack = append(stack[:0], i)
		seen[i] = true
		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			area++
			x, y := p%cw, p/cw
			if x > 0 && fg[p-1] && !seen[p-1] {
				seen[p-1] = true
				stack = append(stack, p-1)
			}
			if x < cw-1 && fg[p+1] && !seen[p+1] {
				seen[p+1] = true
				stack = append(stack, p+1)
			}
			if y > 0 && fg[p-cw] && !seen[p-cw] {
				seen[p-cw] = true
				stack = append(stack, p-cw)
			}
			if y < ch-1 && fg[p+cw] && !seen[p+cw] {
				seen[p+cw] = true
				stack = append(stack, p+cw)
			}
		}
		areas = append(areas, area)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(areas)))
	return areas
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	_ = os.Chdir(root)

	var problems, clean []string

	for _, s := range sheets {
		f, err := os.Open(s.path)
		if err != nil {
			problems = append(problems, fmt.Sprintf("MISSING FILE  %-38s %s", s.label, s.path))
			continue
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			problems = append(problems, fmt.Sprintf("BAD PNG       %-38s %s", s.label, s.path))
			continue
		}
		b := img.Bounds()
		W, H := b.Dx(), b.Dy()
		cw, ch := W/s.cols, H/s.rows

		var issues []string
		// Non-divisible dims are no longer flagged: the engine slices cells
		// PROPORTIONALLY (engine.gridCellRect, 2026-06-10), distributing the
		// remainder instead of truncating frames. This tool mirrors that.

		// background = top-left pixel
		bgR, bgG, bgB, _ := img.At(b.Min.X, b.Min.Y).RGBA()
		tol := 24

		rows := []int{}
		if s.row >= 0 {
			rows = append(rows, s.row)
		} else {
			for r := 0; r < s.rows; r++ {
				rows = append(rows, r)
			}
		}

		var stats []cellStat
		ghostCells := []string{}
		for ri, r := range rows {
			for c := 0; c < s.cols; c++ {
				// Proportional cell boundaries — must match engine.gridCellRect.
				x0 := b.Min.X + c*W/s.cols
				y0 := b.Min.Y + r*H/s.rows
				ccw := b.Min.X + (c+1)*W/s.cols - x0
				cch := b.Min.Y + (r+1)*H/s.rows - y0
				st := cellStat{minX: ccw, minY: cch, maxX: -1, maxY: -1}
				fg := make([]bool, ccw*cch)
				for y := 0; y < cch; y++ {
					for x := 0; x < ccw; x++ {
						pr, pg, pb, pa := img.At(x0+x, y0+y).RGBA()
						if pa < 0x4000 {
							continue
						}
						if near(pr, bgR, tol) && near(pg, bgG, tol) && near(pb, bgB, tol) {
							continue
						}
						fg[y*ccw+x] = true
						if x < st.minX {
							st.minX = x
						}
						if x > st.maxX {
							st.maxX = x
						}
						if y < st.minY {
							st.minY = y
						}
						if y > st.maxY {
							st.maxY = y
						}
					}
				}
				if st.maxX < 0 {
					st.empty = true
				} else {
					st.touchesL = st.minX <= 1
					st.touchesR = st.maxX >= ccw-2
					// Ghost-limb detection (user 2026-06-10: a floating hand
					// from "the previous frame" visible in-game): label the
					// cell's connected components; a second component that is
					// big enough to read as a limb (≥ 300px area) means the
					// generator painted a detached duplicate body part into
					// the cell.
					areas := componentAreas(fg, ccw, cch)
					if len(areas) >= 2 && areas[1] >= 300 {
						ghostCells = append(ghostCells, fmt.Sprintf("r%d c%d (2nd piece %dpx)", ri, c, areas[1]))
					}
				}
				stats = append(stats, st)
			}
		}
		if len(ghostCells) > 0 {
			issues = append(issues, fmt.Sprintf("GHOST PIECES in %d cell(s) -> detached limb/object drawn inside the frame: %v", len(ghostCells), ghostCells))
		}

		// Cross-border continuity (user 2026-06-10: floating hand from the
		// neighbouring frame visible in-game): a shape that SPANS a vertical
		// cell border gets cut by the slicer — its near-border part renders as
		// an orphan limb in the neighbouring frame. Detect: non-bg pixels
		// immediately on BOTH sides of an internal column border at the same y.
		crossBorders := []string{}
		for _, r := range rows {
			y0 := b.Min.Y + r*H/s.rows
			rowH := b.Min.Y + (r+1)*H/s.rows - y0
			for c := 0; c < s.cols-1; c++ {
				bx := b.Min.X + (c+1)*W/s.cols // first column of the right cell
				cont := 0
				for y := 0; y < rowH; y++ {
					lr, lg, lb, la := img.At(bx-1, y0+y).RGBA()
					rr, rg, rb, ra := img.At(bx, y0+y).RGBA()
					lFg := la >= 0x4000 && !(near(lr, bgR, tol) && near(lg, bgG, tol) && near(lb, bgB, tol))
					rFg := ra >= 0x4000 && !(near(rr, bgR, tol) && near(rg, bgG, tol) && near(rb, bgB, tol))
					if lFg && rFg {
						cont++
					}
				}
				if cont > 3 {
					crossBorders = append(crossBorders, fmt.Sprintf("r%d c%d|c%d (%dpx)", r, c, c+1, cont))
				}
			}
		}
		if len(crossBorders) > 0 {
			issues = append(issues, fmt.Sprintf("CONTENT CROSSES %d cell border(s) -> orphan limbs in neighbour frames: %v", len(crossBorders), crossBorders))
		}

		// gather metrics over non-empty cells
		emptyCount := 0
		var bots, cxs, hs, ws []int
		bleed := 0
		for _, st := range stats {
			if st.empty {
				emptyCount++
				continue
			}
			bots = append(bots, st.maxY)
			cxs = append(cxs, (st.minX+st.maxX)/2)
			hs = append(hs, st.maxY-st.minY+1)
			ws = append(ws, st.maxX-st.minX+1)
			if st.touchesL && st.touchesR {
				bleed++
			}
		}
		if emptyCount > 0 {
			issues = append(issues, fmt.Sprintf("%d/%d EMPTY cells (blank frames -> blink in loop)", emptyCount, len(stats)))
		}
		if len(bots) > 1 {
			sort.Ints(bots)
			sort.Ints(cxs)
			sort.Ints(hs)
			sort.Ints(ws)
			footDrift := bots[len(bots)-1] - bots[0]
			cxDrift := cxs[len(cxs)-1] - cxs[0]
			hPump := float64(hs[len(hs)-1]-hs[0]) / math.Max(1, float64(hs[len(hs)-1]))
			if footDrift > ch/24+6 {
				issues = append(issues, fmt.Sprintf("FOOT drift %dpx (cell h %d) -> vertical jumping", footDrift, ch))
			}
			if cxDrift > cw/10+8 {
				issues = append(issues, fmt.Sprintf("CENTER-X drift %dpx (cell w %d) -> horizontal sliding", cxDrift, cw))
			}
			if hPump > 0.12 {
				issues = append(issues, fmt.Sprintf("HEIGHT pump %d..%dpx (%.0f%%) -> size pulsing", hs[0], hs[len(hs)-1], hPump*100))
			}
			if bleed > 0 {
				issues = append(issues, fmt.Sprintf("%d cells touch BOTH side edges -> content bleeds across cells (two-frames-at-once)", bleed))
			}
		}

		if len(issues) == 0 {
			clean = append(clean, fmt.Sprintf("OK   %-38s cell %dx%d, frames %d", s.label, cw, ch, len(stats)))
		} else {
			head := fmt.Sprintf("WARN %-38s %s  [grid %dx%d, cell %dx%d]", s.label, s.path, s.cols, s.rows, cw, ch)
			problems = append(problems, head)
			for _, is := range issues {
				problems = append(problems, "      - "+is)
			}
		}
	}

	fmt.Println("==== PROBLEMS ====")
	for _, p := range problems {
		fmt.Println(p)
	}
	fmt.Println()
	fmt.Println("==== CLEAN ====")
	for _, c := range clean {
		fmt.Println(c)
	}
}

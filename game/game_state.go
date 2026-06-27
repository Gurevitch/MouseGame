package game

import "os"

// Chapter identifiers used as integer values in VarStore scope "game"
// under the key VarChapter. They also map 1:1 to the old `day` concept for
// Camp Sylvania, but allow future cities to have their own numbered arcs.
const (
	ChapterCampDay1  = 1 // Arrival, everything normal
	ChapterCampDay2  = 2 // Weirdness begins at Camp Sylvania
	ChapterParis     = 3 // Marcus's memory arc in Paris
	ChapterJerusalem = 4 // Jake's courage arc
	ChapterTokyo     = 5 // Lily's voice arc
	ChapterRioBA     = 6 // Tommy's family arc (Rio + Buenos Aires)
	ChapterRome      = 7 // Danny's identity arc
	ChapterFinale    = 8 // Mexico City reunion
)

// Canonical VarStore keys. Using string constants prevents typo drift across
// handlers/sequences/save files.
const (
	// "game" scope: persists for the entire playthrough
	VarChapter           = "chapter"            // current chapter number (see above)
	VarDay               = "day"                // legacy camp day counter (1 or 2)
	VarParisUnlocked     = "paris_unlocked"     // travel pin unlocked
	VarJerusalemUnlocked = "jerusalem_unlocked" // travel pin unlocked
	VarTokyoUnlocked     = "tokyo_unlocked"
	VarRioUnlocked       = "rio_unlocked"
	VarRomeUnlocked      = "rome_unlocked"
	VarMexicoUnlocked    = "mexico_unlocked"
	VarMarcusHealed      = "marcus_healed"
	VarJakeHealed        = "jake_healed"
	VarLilyHealed        = "lily_healed"
	VarTommyHealed       = "tommy_healed"
	VarDannyHealed       = "danny_healed"
	VarNightSceneDone    = "night_scene_done"
	VarMonologueIntro    = "monologue_intro_played"
	VarMonologueParis    = "monologue_paris_played"
	VarMonologueLouvre   = "monologue_louvre_played" // museum first-arrival beat (#28)
	VarParisDone         = "paris_done"              // postcard obtained → camp return unlocked (#32)
	VarJerNotePlaced     = "jer_note_placed"         // Jerusalem: note tucked in the Wall → return flight + coin (#26)
	// Japan/Kyoto chapter opening (Lily's arc).
	VarLilyArcStarted   = "lily_arc_started"   // sad Lily revealed at the lake (post-Jake-heal)
	VarLilyLakeMet      = "lily_lake_met"      // PP has talked to Lily at the lake
	VarHigginsRudeDone  = "higgins_rude_done"  // the rude-Higgins + camera aside played → Tokyo unlocks
	VarJpGroveRevealed  = "jp_grove_revealed"  // Oba-chan said "follow me" → the hidden sakura grove exit opens
	VarJpRamenOpen      = "jp_ramen_open"      // PP gave Hiro his fire-striker → stall opens, the waiting line sits
	VarJpTeaLearned     = "jp_tea_learned"     // Kiku the geisha taught PP the tea ceremony → the matcha/bowl shelves unlock
	VarJpTeaDone        = "jp_tea_done"        // PP shared the matcha ceremony with the tea master → grove entry allowed

	// "chapter" scope: resets when a chapter ends (via ResetChapter)
	VarMetKids        = "met_kids" // How many kids PP has talked to on Day 1
	VarTalkedToMarcus = "talked_to_marcus"
	VarDay2Started    = "day2_started"
)

// Scope names (as used by VarStore.Get/Set)
const (
	ScopeGame    = "game"
	ScopeChapter = "chapter"
	ScopeScene   = "scene"
)

// --- Convenience accessors on Game ---
//
// These wrap g.vars so the rest of the code reads like a struct field access
// but the backing storage is the VarStore (which serializes cleanly and can
// be inspected/modified by sequences without knowing Go struct layout).

// Chapter returns the current chapter number.
func (g *Game) Chapter() int {
	if v := g.vars.Get(ScopeGame, VarChapter); v != 0 {
		return v
	}
	return ChapterCampDay1
}

// SetChapter moves the game to a new chapter and clears chapter-scope vars.
func (g *Game) SetChapter(ch int) {
	if g.Chapter() == ch {
		return
	}
	g.vars.Set(ScopeGame, VarChapter, ch)
	g.vars.ResetChapter()
}

// Day mirrors the current chapter when we're in the camp arcs. For non-camp
// chapters it always returns the value of chapter-scoped VarDay (defaults 1).
func (g *Game) Day() int {
	switch g.Chapter() {
	case ChapterCampDay1:
		return 1
	case ChapterCampDay2:
		return 2
	}
	if v := g.vars.Get(ScopeChapter, VarDay); v != 0 {
		return v
	}
	return 1
}

// IsCityChapter returns true once we've left Camp Sylvania.
func (g *Game) IsCityChapter() bool {
	return g.Chapter() >= ChapterParis
}

// syncFlagsToVars pushes the flat runtime flags into the VarStore so save
// files, sequence scripts and debug dumps see a consistent snapshot.
// Runs every frame (cheap map writes) and keeps both views in lockstep
// while we finish migrating call-sites to VarStore directly.
func (g *Game) syncFlagsToVars() {
	if g.vars == nil {
		return
	}
	g.vars.Set(ScopeGame, VarDay, g.day)
	g.vars.SetBool(ScopeGame, VarParisUnlocked, g.parisUnlocked)
	g.vars.SetBool(ScopeGame, VarNightSceneDone, g.nightSceneDone)
	g.vars.SetBool(ScopeGame, VarMarcusHealed, g.marcusHealed)
	g.vars.SetBool(ScopeGame, VarMonologueIntro, g.monologuePlayed)
	g.vars.SetBool(ScopeGame, VarMonologueParis, g.parisMonologuePlayed)

	g.vars.Set(ScopeChapter, VarMetKids, g.metKids)
	g.vars.SetBool(ScopeChapter, VarTalkedToMarcus, g.talkedToMarcus)
	g.vars.SetBool(ScopeChapter, VarDay2Started, g.day2Started)

	// Keep Chapter number aligned with the camp day so the two stay in sync
	// until city chapters take over.
	if g.Chapter() < ChapterParis {
		switch g.day {
		case 1:
			g.vars.Set(ScopeGame, VarChapter, ChapterCampDay1)
		case 2:
			g.vars.Set(ScopeGame, VarChapter, ChapterCampDay2)
		}
	}
}

// campMoodLevel returns the camp's darkness GRADE (2026-06-21 #20/#21):
//   0 = normal (pre-Paris)
//   1 = mid-dark (post-Paris, the Marcus affliction arc)
//   2 = fully dark (deeper - a second kid afflicted; the Jerusalem leg onward)
// Forward-compatible: widen the level-2 condition as later cities land.
func (g *Game) campMoodLevel() int {
	if g == nil || g.vars == nil || !g.vars.GetBool(ScopeGame, VarParisDone) {
		return 0
	}
	if g.vars.GetBool(ScopeGame, VarJerusalemUnlocked) {
		return 2
	}
	return 1
}

// firstExistingPath returns the first path that exists on disk, or "".
func firstExistingPath(paths ...string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// moodBG builds the background for the current darkness grade from the
// day1/day2/day3 folders (2026-06-21 reorg: the same filename lives in each
// folder). Falls back DOWN to day1 when a grade's art isn't on disk yet, so a
// missing day2 (mid-dark) / day3 (full dark) degrades gracefully to day1.
//   level 0 → day1/  (normal)
//   level 1 → day2/  (mid-dark, then day1)
//   level 2 → day3/  (full dark, then day2, then day1)
func (g *Game) moodBG(level int, file string) *background {
	const base = "assets/images/locations/camp/background/"
	var folders []string
	switch level {
	case 2:
		folders = []string{"day3/", "day2/", "day1/"}
	case 1:
		folders = []string{"day2/", "day1/"}
	default:
		folders = []string{"day1/"}
	}
	paths := make([]string, 0, len(folders))
	for _, f := range folders {
		paths = append(paths, base+f+file)
	}
	p := firstExistingPath(paths...)
	if p == "" {
		p = base + "day1/" + file
	}
	return newPNGBackground(g.renderer, p)
}

// applyCampMood swaps the camp backgrounds to the graded "affliction" art once
// PP has returned from France (paris_done), and back to normal otherwise (#34,
// graded + folder-based 2026-06-21). Covers the grounds, the airstrip landing,
// and the cabin interiors. marcus_room is intentionally NOT touched here - its
// bg is driven by the day/night setSceneAltBG system. Safe to call repeatedly.
func (g *Game) applyCampMood() {
	if g == nil || g.sceneMgr == nil {
		return
	}
	level := g.campMoodLevel()
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok && grounds != nil {
		grounds.bg = g.moodBG(level, "camp_grounds.png")
	}
	if landing, ok := g.sceneMgr.scenes["camp_landing"]; ok && landing != nil {
		landing.bg = g.moodBG(level, "camp_landing.png")
	}
	for _, room := range []string{"jake_room", "lily_room", "tommy_room", "danny_room"} {
		if s, ok := g.sceneMgr.scenes[room]; ok && s != nil {
			s.bg = g.moodBG(level, room+".png")
		}
	}
	// Marcus's room darkens with the camp too while he's still afflicted; once
	// he's healed the heal callback brightens it (sceneAltBGs day) and we leave
	// it alone here (so a load after the heal keeps it bright). The night
	// cutscene runs pre-Paris at level 0, so it's unaffected.
	if !g.vars.GetBool(ScopeGame, VarMarcusHealed) {
		if s, ok := g.sceneMgr.scenes["marcus_room"]; ok && s != nil {
			s.bg = g.moodBG(level, "marcus_room.png")
		}
	}
}

// syncVarsToFlags pulls values back from the VarStore into the flat runtime
// flags. Call this after loading a save file or after mutating VarStore via
// a scripted sequence to make sure the engine loop sees the change.
func (g *Game) syncVarsToFlags() {
	if g.vars == nil {
		return
	}
	if v := g.vars.Get(ScopeGame, VarDay); v != 0 {
		g.day = v
	}
	g.parisUnlocked = g.vars.GetBool(ScopeGame, VarParisUnlocked)
	g.nightSceneDone = g.vars.GetBool(ScopeGame, VarNightSceneDone)
	g.marcusHealed = g.vars.GetBool(ScopeGame, VarMarcusHealed)
	g.monologuePlayed = g.vars.GetBool(ScopeGame, VarMonologueIntro)
	g.parisMonologuePlayed = g.vars.GetBool(ScopeGame, VarMonologueParis)

	g.metKids = g.vars.Get(ScopeChapter, VarMetKids)
	g.talkedToMarcus = g.vars.GetBool(ScopeChapter, VarTalkedToMarcus)
	g.day2Started = g.vars.GetBool(ScopeChapter, VarDay2Started)
}

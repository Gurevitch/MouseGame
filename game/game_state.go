package game

// Chapter identifiers used as integer values in VarStore scope "game"
// under the key VarChapter. They also map 1:1 to the old `day` concept for
// Camp Sylvania, but allow future cities to have their own numbered arcs.
const (
	ChapterCampDay1 = 1 // Arrival, everything normal
	ChapterCampDay2 = 2 // Weirdness begins at Camp Sylvania
	ChapterParis    = 3 // Marcus's memory arc in Paris
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
	VarChapter            = "chapter"              // current chapter number (see above)
	VarDay                = "day"                  // legacy camp day counter (1 or 2)
	VarParisUnlocked      = "paris_unlocked"       // travel pin unlocked
	VarJerusalemUnlocked  = "jerusalem_unlocked"   // travel pin unlocked
	VarTokyoUnlocked      = "tokyo_unlocked"
	VarRioUnlocked        = "rio_unlocked"
	VarRomeUnlocked       = "rome_unlocked"
	VarMexicoUnlocked     = "mexico_unlocked"
	VarMarcusHealed       = "marcus_healed"
	VarJakeHealed         = "jake_healed"
	VarLilyHealed         = "lily_healed"
	VarTommyHealed        = "tommy_healed"
	VarDannyHealed        = "danny_healed"
	VarNightSceneDone     = "night_scene_done"
	VarMonologueIntro     = "monologue_intro_played"
	VarMonologueParis     = "monologue_paris_played"

	// "chapter" scope: resets when a chapter ends (via ResetChapter)
	VarMetKids        = "met_kids"        // How many kids PP has talked to on Day 1
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

package game

import (
	"encoding/json"
	"fmt"
	"os"

	"bitbucket.org/Local/games/PP/engine"
)

// SaveState captures the complete game state for serialization.
type SaveState struct {
	// VarStore state
	Vars *VarStore `json:"vars"`

	// Legacy fields (for compatibility during migration)
	Day            int  `json:"day"`
	MetKids        int  `json:"metKids"`
	TalkedToMarcus bool `json:"talkedToMarcus"`
	ParisUnlocked  bool `json:"parisUnlocked"`
	NightSceneDone bool `json:"nightSceneDone"`
	Day2Started    bool `json:"day2Started"`
	MarcusHealed   bool `json:"marcusHealed"`

	// Current position
	CurrentScene string  `json:"currentScene"`
	PlayerX      float64 `json:"playerX"`
	PlayerY      float64 `json:"playerY"`

	// Inventory
	ItemNames []string `json:"items"`

	// Monologue state
	MonologuePlayed      bool `json:"monologuePlayed"`
	ParisMonologuePlayed bool `json:"parisMonologuePlayed"`
}

// SaveGame saves the current game state to a file.
func (g *Game) SaveGame(path string) error {
	g.syncFlagsToVars()
	state := SaveState{
		Vars:                 g.vars,
		Day:                  g.day,
		MetKids:              g.metKids,
		TalkedToMarcus:       g.talkedToMarcus,
		ParisUnlocked:        g.parisUnlocked,
		NightSceneDone:       g.nightSceneDone,
		Day2Started:          g.day2Started,
		MarcusHealed:         g.marcusHealed,
		CurrentScene:         g.sceneMgr.currentName,
		PlayerX:              g.player.x,
		PlayerY:              g.player.y,
		MonologuePlayed:      g.monologuePlayed,
		ParisMonologuePlayed: g.parisMonologuePlayed,
	}

	for _, item := range g.inv.items {
		state.ItemNames = append(state.ItemNames, item.name)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal save state: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write save file: %w", err)
	}

	fmt.Printf("Game saved to %s\n", path)
	return nil
}

// LoadGame restores game state from a file.
func (g *Game) LoadGame(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read save file: %w", err)
	}

	var state SaveState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("parse save file: %w", err)
	}

	if state.Vars != nil {
		g.vars = state.Vars
		// 2026-08-08 #8: nil-map guard — a save missing a scope key yields a
		// nil map, and the next vars.Set would panic on nil-map write.
		if g.vars.Game == nil {
			g.vars.Game = make(map[string]int)
		}
		if g.vars.Chapter == nil {
			g.vars.Chapter = make(map[string]int)
		}
		if g.vars.Scene == nil {
			g.vars.Scene = make(map[string]int)
		}
	}

	// Restore legacy fields
	g.day = state.Day
	g.metKids = state.MetKids
	g.talkedToMarcus = state.TalkedToMarcus
	g.parisUnlocked = state.ParisUnlocked
	g.nightSceneDone = state.NightSceneDone
	g.day2Started = state.Day2Started
	g.marcusHealed = state.MarcusHealed
	g.monologuePlayed = state.MonologuePlayed
	g.parisMonologuePlayed = state.ParisMonologuePlayed

	// If the VarStore is newer than the legacy fields (e.g. save was written
	// by a city chapter build that stopped writing legacy fields) let it win.
	g.syncVarsToFlags()

	// Restore inventory. 2026-08-08 #8: resolve display names through the
	// registry (itemIDFromName's 7-name switch silently dropped every other
	// quest item — a mid-Paris save lost its Card/Confiture/etc on load).
	g.inv.items = nil
	for _, name := range state.ItemNames {
		if item := g.items.createItem(g.items.idForName(name)); item != nil {
			g.inv.addItem(item)
		}
	}

	// Restore scene and player position. 2026-08-08 #8: the transition's
	// fade-in completion used to snap PP to the scene SPAWN, clobbering the
	// direct x/y assignment below a frame later — route the saved coords
	// through the pending-restore hook instead (consumed at fade-in).
	g.sceneMgr.restorePosPending = true
	g.sceneMgr.restoreX = state.PlayerX
	g.sceneMgr.restoreY = state.PlayerY
	g.sceneMgr.transitionTo(state.CurrentScene, g.player)
	g.player.x = state.PlayerX
	g.player.y = state.PlayerY

	if g.parisUnlocked {
		g.travelMap.setUnlocked("paris_street", true)
	}
	if g.marcusHealed {
		g.travelMap.setUnlocked("jerusalem_entrance", true)
		if mRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
			if day, ok := g.sceneAltBGs["marcus_room/day"]; ok {
				mRoom.bg = day
			}
		}
	}
	// #34: restore the camp-return pin and the darkened-grounds mood for saves
	// taken after the France trip.
	if g.vars != nil && g.vars.GetBool(ScopeGame, VarParisDone) {
		g.travelMap.setUnlocked("camp_entrance", true)
	}
	g.applyCampMood()
	// 2026-07-24 (user #4): the see-Marcus-first office gate was armed at
	// setup time, before this load ran — lift it if the save already talked
	// to Marcus.
	if g.talkedToMarcus {
		g.restoreHigginsWorriedDialog()
	}

	g.reconcileLoadedWorld()

	fmt.Printf("Game loaded from %s\n", path)
	return nil
}

// reconcileLoadedWorld re-derives per-scene state that LIVE events set
// during play but a bare load can't reconstruct from the flags alone
// (2026-08-01 user: loading a late save showed the Day-1 camp — the five
// kids lined up on the grounds, the lake daisy back on the ground, and sad
// Lily missing from the dock).
func (g *Game) reconcileLoadedWorld() {
	// Day 2+: the kids are in their cabins, not lined up on the grounds —
	// live play hides them through the day-1 chats and the bedtime beat,
	// all of which are runtime-only state.
	if g.day >= 2 {
		if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
			for _, n := range grounds.npcs {
				switch n.name {
				case "Tommy", "Jake", "Lily", "Marcus", "Danny":
					n.hidden = true
					n.silent = true
				}
			}
		}
	}
	// The lake daisy is a Day-1 beat (it becomes Lily's flower): gone on any
	// later day, and gone whenever PP already carries it.
	if g.day >= 2 || g.inv.hasItem("Flower") {
		if lake, ok := g.sceneMgr.scenes["camp_lake"]; ok {
			for _, fi := range lake.floorItems {
				if fi.name == "Flower" {
					fi.visible = false
					fi.hidden = false
				}
			}
		}
	}
	// Lily's arc: once Jake's heal opened it (and she isn't healed yet), sad
	// Lily sits at the end of the dock and her cabin double is gone — the
	// reveal normally runs inside Jake's heal callback only.
	if g.vars != nil && g.vars.GetBool(ScopeGame, VarLilyArcStarted) &&
		!g.vars.GetBool(ScopeGame, VarLilyHealed) {
		if lk, ok := g.sceneMgr.scenes["camp_lake"]; ok {
			for _, n := range lk.npcs {
				if n.name == "Lily" {
					lakeLily := n
					lakeLily.hidden = false
					lakeLily.silent = false
					lakeLily.setStrange(true)
					if g.vars.GetBool(ScopeGame, VarLilyLakeMet) {
						lakeLily.dialog = lilyLeaveMeAloneDialog
					} else {
						lakeLily.dialog = lilyStrangeDialog
					}
					lakeLily.onDialogEnd = func() {
						g.vars.SetBool(ScopeGame, VarLilyLakeMet, true)
						lakeLily.dialog = lilyLeaveMeAloneDialog
					}
					break
				}
			}
		}
		if lr, ok := g.sceneMgr.scenes["lily_room"]; ok {
			for _, n := range lr.npcs {
				if n.name == "Lily" {
					n.hidden = true
					n.silent = true
					break
				}
			}
		}
	}

	// 2026-08-08 #8/#14: full re-derivation of the Paris chain, the camp
	// cast, and the travel pins from the persisted vars.
	g.applyCampRevealState()
	g.applyParisChainState()
	g.travelMap.restoreUnlockedFromVars(g.vars)
}

// applyCampRevealState re-derives the camp cast's hidden/silent/strange
// state from the persisted day + heal-chain vars (2026-08-08 #14: loading a
// day-2+ save left the WHOLE camp empty — room Marcus is born hidden and
// office Higgins born silent, revealed only by live startDay2 / heal
// callbacks a load never re-runs, so the Postcard could never be delivered).
func (g *Game) applyCampRevealState() {
	if g.day >= 2 {
		if marcusRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok {
			for _, n := range marcusRoom.npcs {
				if n.name == "Marcus" {
					n.hidden = false
					if !g.marcusHealed {
						n.dialog = marcusStrangeDialog
						n.setStrange(true)
					}
					break
				}
			}
		}
		if office, ok := g.sceneMgr.scenes["camp_office"]; ok {
			for _, n := range office.npcs {
				if n.name == "Director Higgins" {
					n.silent = false
					break
				}
			}
		}
	}
	// Heal-chain reveals: each heal callback un-silences the NEXT kid in
	// their room (marcus→jake, lily→tommy, tommy→danny). Strange stays on
	// until that kid's own heal.
	reveal := func(scene, kid string, healed bool) {
		room, ok := g.sceneMgr.scenes[scene]
		if !ok {
			return
		}
		for _, n := range room.npcs {
			if n.name == kid {
				n.silent = false
				n.setStrange(!healed)
				break
			}
		}
	}
	getB := func(k string) bool { return g.vars != nil && g.vars.GetBool(ScopeGame, k) }
	if g.marcusHealed {
		reveal("jake_room", "Jake", getB(VarJakeHealed))
	}
	if getB(VarLilyHealed) {
		reveal("tommy_room", "Tommy", getB(VarTommyHealed))
	}
	if getB(VarTommyHealed) {
		reveal("danny_room", "Danny", getB(VarDannyHealed))
	}
}

// applyParisChainState re-derives every Paris NPC dialog pointer, floor
// item, and gate from the persisted paris_* vars (2026-08-08 #8: the chain
// used to live in closure locals — saving in the Louvre after handing
// Claude the pass hard-stuck the game on relaunch).
func (g *Game) applyParisChainState() {
	if g.vars == nil {
		return
	}
	getB := func(k string) bool { return g.vars.GetBool(ScopeGame, k) }
	if ps, ok := g.sceneMgr.scenes["paris_street"]; ok {
		for _, n := range ps.npcs {
			switch n.name {
			case "Pierre":
				switch g.vars.Get(ScopeGame, VarParisPierreStage) {
				case 1:
					n.hintState = 1
					n.dialog = pierreWaitingSpreadDialog
					n.altDialogRequiresItem = "Confiture"
				case 2:
					n.hintState = 2
					n.dialog = pierreArtistPostDialog
					n.altDialogFunc = nil
					n.altDialogRequiresItem = ""
				}
			case "Claude":
				if getB(VarParisLouvreOpen) {
					n.dialog = gendarmePostDialog
					n.altDialogFunc = nil
					n.altDialogRequiresHeld = false
					n.altDialogRequiresItem = ""
				}
			case "Madame Margaux":
				if getB(VarParisPigeonsClear) {
					n.dialog = pigeonLadyPostDialog
					n.altDialogFunc = nil
					n.altDialogRequiresHeld = false
					n.altDialogRequiresItem = ""
				}
			}
		}
		for _, fi := range ps.floorItems {
			switch fi.name {
			case "Rolling Pin":
				if getB(VarParisPinTaken) {
					fi.visible = false
					fi.hidden = false
				}
			case "Charcoal Pencil":
				if getB(VarParisPigeonsClear) && !getB(VarParisPencilTaken) {
					// The pigeon is gone: show the exposed-pencil pot art
					// (clearPotPigeon's swap, minus the fly-up ambient).
					tex, w, h := engine.SafeTextureFromPNGKeyed(g.renderer,
						"assets/images/locations/paris/props/flower_pot_pencil.png")
					if tex != nil {
						fi.tex = tex
						fi.srcW = w
						fi.srcH = h
					}
				}
				if getB(VarParisPencilTaken) {
					fi.visible = false
					fi.hidden = false
				}
			}
		}
	}
	if bk, ok := g.sceneMgr.scenes["paris_bakery"]; ok {
		for _, n := range bk.npcs {
			switch n.name {
			case "Madame Poulain":
				// Deepest state wins; the click-time selector handles the
				// branch logic, this only restores the resting dialog.
				switch {
				case getB(VarParisSouvenirDone):
					n.dialog = bakeryWomanSouvenirDoneDialog
				case g.marcusHealed && getB(VarParisSouvenirArmed):
					n.dialog = bakeryWomanLouvreSouvenirDialog
				case getB(VarParisPoulainTraded):
					n.dialog = bakeryWomanPostDialog
				}
			case "Monsieur Henri":
				if getB(VarParisHenriTraded) {
					n.dialog = henriPostTradeDialog
					n.altDialogFunc = nil
					n.altDialogRequiresItem = ""
					g.addHenriCupDecor()
				}
			case "Mademoiselle Camille":
				switch {
				case getB(VarParisSketchDone):
					n.dialog = camillePostSketchDialog
				case getB(VarParisCamilleAsked):
					n.dialog = camillePencilReminderDialog
				}
			}
		}
	}
	if lv, ok := g.sceneMgr.scenes["paris_louvre"]; ok {
		for _, n := range lv.npcs {
			if n.name == "Curator Beaumont" {
				switch {
				case getB(VarParisPostcardGiven):
					n.dialog = museumCuratorPostDialog
				case getB(VarParisSketchAsked):
					n.dialog = curatorWaitingDialog
				}
				break
			}
		}
	}
}

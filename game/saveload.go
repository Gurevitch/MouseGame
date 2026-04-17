package game

import (
	"encoding/json"
	"fmt"
	"os"
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

	// Restore inventory
	g.inv.items = nil
	for _, name := range state.ItemNames {
		if item := g.items.createItem(itemIDFromName(name)); item != nil {
			g.inv.addItem(item)
		}
	}

	// Restore scene and player position
	g.sceneMgr.transitionTo(state.CurrentScene, g.player)
	g.player.x = state.PlayerX
	g.player.y = state.PlayerY

	if g.parisUnlocked {
		g.travelMap.setUnlocked("paris_street", true)
	}
	if g.marcusHealed {
		g.travelMap.setUnlocked("jerusalem_street", true)
		if mRoom, ok := g.sceneMgr.scenes["marcus_room"]; ok && g.marcusRoomBg != nil {
			mRoom.bg = g.marcusRoomBg
		}
	}

	fmt.Printf("Game loaded from %s\n", path)
	return nil
}

// itemIDFromName maps item display names to registry IDs
func itemIDFromName(name string) string {
	switch name {
	case "Travel Map":
		return "travel_map"
	case "Flower":
		return "flower"
	case "Postcard":
		return "postcard"
	case "Coin Rubbing":
		return "coin_rubbing"
	case "Pressed Sakura":
		return "pressed_sakura"
	case "Dance Card":
		return "dance_card"
	case "Inscription Rubbing":
		return "inscription_rubbing"
	default:
		return name
	}
}

package game

import (
	"encoding/json"
	"testing"
)

func TestSceneDef_ParsesAnimatedBackgroundFields(t *testing.T) {
	const raw = `{
		"name": "airplane_flight",
		"background": "x.png",
		"backgroundFrames": 6,
		"backgroundFrameSeconds": 0.15,
		"spawnX": 700, "spawnY": 400
	}`
	var def sceneDef
	if err := json.Unmarshal([]byte(raw), &def); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if def.BackgroundFrames != 6 {
		t.Errorf("BackgroundFrames: got %d, want 6", def.BackgroundFrames)
	}
	if def.BackgroundFrameSeconds != 0.15 {
		t.Errorf("BackgroundFrameSeconds: got %v, want 0.15", def.BackgroundFrameSeconds)
	}
}

func TestSceneDef_StaticSceneDefaultsToZero(t *testing.T) {
	// Existing scenes (camp_grounds, paris_street, etc.) don't set the
	// new fields. They must default to zero so buildBackground falls
	// through to the static path.
	const raw = `{
		"name": "camp_grounds",
		"background": "x.png",
		"spawnX": 700, "spawnY": 600
	}`
	var def sceneDef
	if err := json.Unmarshal([]byte(raw), &def); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if def.BackgroundFrames != 0 || def.BackgroundFrameSeconds != 0 {
		t.Errorf("static scene picked up animation fields: frames=%d sec=%v",
			def.BackgroundFrames, def.BackgroundFrameSeconds)
	}
}

func TestParseArrow(t *testing.T) {
	cases := map[string]arrowDir{
		"left":      arrowLeft,
		"right":     arrowRight,
		"up":        arrowUp,
		"down":      arrowDown,
		"downRight": arrowDownRight,
		"":          arrowNone,
		"garbage":   arrowNone,
	}
	for in, want := range cases {
		if got := parseArrow(in); got != want {
			t.Errorf("parseArrow(%q): got %v, want %v", in, got, want)
		}
	}
}

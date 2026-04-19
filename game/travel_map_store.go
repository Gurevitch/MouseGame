package game

import (
	"encoding/json"
	"fmt"
	"os"
)

// travelLocationJSON is the on-disk shape of `assets/data/travel_map.json`.
// Runtime-only fields (landmarkTex, landmarkW, landmarkH) are populated by
// newTravelMap after loading, not stored in JSON.
type travelLocationJSON struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Scene        string `json:"scene"`
	PinX         int32  `json:"pinX"`
	PinY         int32  `json:"pinY"`
	Unlocked     bool   `json:"unlocked"`
	Info         string `json:"info"`
	Landmark     string `json:"landmark"`
	Audio        string `json:"audio"`
	RelevantWhen string `json:"relevantWhen"`
}

type travelMapJSON struct {
	Locations []travelLocationJSON `json:"locations"`
}

// loadTravelLocations reads the travel-map data file and returns the list of
// locations with all JSON-authored fields populated. Returns (nil, err) on
// missing file or parse error so the caller can decide whether to panic or
// fall back to an empty map. Runtime-only fields (textures, dimensions) are
// left zero — newTravelMap fills those after loading textures.
func loadTravelLocations(path string) ([]travelLocation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var raw travelMapJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	out := make([]travelLocation, 0, len(raw.Locations))
	for _, r := range raw.Locations {
		out = append(out, travelLocation{
			id:           r.ID,
			name:         r.Name,
			scene:        r.Scene,
			pinX:         r.PinX,
			pinY:         r.PinY,
			unlocked:     r.Unlocked,
			info:         r.Info,
			landmarkPath: r.Landmark,
			audio:        r.Audio,
			relevantWhen: r.RelevantWhen,
		})
	}
	return out, nil
}

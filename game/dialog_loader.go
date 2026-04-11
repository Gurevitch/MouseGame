package game

import (
	"encoding/json"
	"fmt"
	"os"
)

type dialogFile struct {
	Dialogs map[string][]dialogEntry `json:"dialogs"`
}

// loadDialogFile loads a JSON dialog file and returns a map of dialog name → entries
func loadDialogFile(path string) map[string][]dialogEntry {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Warning: could not load dialog file %s: %v\n", path, err)
		return nil
	}

	var file dialogFile
	if err := json.Unmarshal(data, &file); err != nil {
		fmt.Printf("Warning: could not parse dialog file %s: %v\n", path, err)
		return nil
	}

	return file.Dialogs
}

// getDialog returns a specific dialog from a loaded dialog map, with fallback
func getDialog(dialogs map[string][]dialogEntry, name string) []dialogEntry {
	if dialogs == nil {
		return nil
	}
	if d, ok := dialogs[name]; ok {
		return d
	}
	fmt.Printf("Warning: dialog '%s' not found\n", name)
	return nil
}

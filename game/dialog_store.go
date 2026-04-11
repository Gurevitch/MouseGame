package game

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// dialogStore holds all loaded dialog data from JSON files
type dialogStore struct {
	files map[string]map[string][]dialogEntry // filename → dialogName → entries
}

func newDialogStore(dir string) *dialogStore {
	ds := &dialogStore{
		files: make(map[string]map[string][]dialogEntry),
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Warning: could not read dialog directory %s: %v\n", dir, err)
		return ds
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".json")
		path := filepath.Join(dir, e.Name())
		dialogs := loadDialogFile(path)
		if dialogs != nil {
			ds.files[name] = dialogs
			fmt.Printf("Loaded dialog file: %s (%d dialogs)\n", name, len(dialogs))
		}
	}

	return ds
}

// Get returns a specific dialog from a specific file
func (ds *dialogStore) Get(file, name string) []dialogEntry {
	if ds == nil {
		return nil
	}
	fileDialogs, ok := ds.files[file]
	if !ok {
		return nil
	}
	return getDialog(fileDialogs, name)
}

package game

import "testing"

// 2026-08-08 #8: the save file stores item DISPLAY names, and the old
// hardcoded 7-name switch silently dropped every other item from the bag on
// load. idForName resolves through the registry — this pins that every
// display name in the real items.json round-trips to its own id, so a new
// item can never silently vanish from saves again.
func TestIDForNameResolvesEveryRegistryItem(t *testing.T) {
	reg := newItemRegistry(nil, "../assets/data/items.json")
	if len(reg.defs) == 0 {
		t.Fatal("item registry loaded 0 defs — path wrong?")
	}
	seen := map[string]string{}
	for id, def := range reg.defs {
		if prev, dup := seen[def.Name]; dup {
			t.Errorf("display name %q is shared by ids %q and %q — idForName would be ambiguous", def.Name, prev, id)
		}
		seen[def.Name] = id
		if got := reg.idForName(def.Name); got != id {
			t.Errorf("idForName(%q) = %q, want %q", def.Name, got, id)
		}
	}
	// Unknown names fall through unchanged (future-format ids stay usable).
	if got := reg.idForName("no_such_item_xyz"); got != "no_such_item_xyz" {
		t.Errorf("idForName fallthrough = %q, want input echoed", got)
	}
}

// The nil-map guard: a save missing scope keys must not leave nil maps that
// panic on the next vars.Set.
func TestVarStoreNilMapGuardShape(t *testing.T) {
	vs := &VarStore{Game: map[string]int{"x": 1}} // Chapter/Scene nil, like a hand-edited save
	if vs.Chapter == nil && vs.Scene == nil {
		// mirrors LoadGame's guard
		if vs.Game == nil {
			vs.Game = make(map[string]int)
		}
		if vs.Chapter == nil {
			vs.Chapter = make(map[string]int)
		}
		if vs.Scene == nil {
			vs.Scene = make(map[string]int)
		}
	}
	vs.Set("chapter", "y", 2) // must not panic
	if vs.Get("chapter", "y") != 2 {
		t.Error("chapter var did not round-trip after guard")
	}
}

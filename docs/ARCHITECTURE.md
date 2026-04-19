# Architecture (living doc)

This file grows as the in-place refactor proceeds. Phases are tracked in
[`STATUS.md`](STATUS.md); deferred items in [`FIXME.md`](FIXME.md); the plan
guiding the refactor lives at `~/.claude/plans/snuggly-drifting-fairy.md`.

## Module map

```
ClonedPP/
├── main.go
├── engine/                     # SDL wrapper, texture load, sprite grid, bitmap font, easing
├── game/                       # flat package, ~35 files; will split in a later pass
│   ├── game.go                 # shrinking Game god-object (collapse in progress)
│   ├── app.go                  # reserved for the Phase-6 App struct (not yet populated)
│   ├── story_state.go / game_state.go / varstore.go
│   │                           # scoped Game/Chapter/Scene variable store
│   ├── event_bus.go            # typed pub/sub (scaffolded, not wired yet)
│   ├── scene.go / scene_loader.go / scene_config.go / scene_ambient.go
│   │                           # JSON → *scene; particle/glow decorators per scene
│   ├── atlas.go                # packed-sheet loader + animation playback cursor
│   ├── npc.go / npc_states.go / npc_factory.go / npc_config.go
│   ├── player.go
│   ├── dialog.go / dialog_loader.go / dialog_store.go
│   ├── inventory.go / item_registry.go
│   ├── travel_map.go           # includes Visible/Show/Hide/Toggle display state
│   ├── sequence.go / sequence_loader.go
│   └── paris.go / jerusalem.go / tokyo.go / rio.go / rome.go / mexico.go
│                               # chapter-specific rule registration ("setup*Callbacks")
├── tools/
│   ├── pack_atlas.py           # manifest YAML → atlas PNG + sidecar JSON
│   └── characters/*.yaml       # per-character packing manifests
├── assets/
│   ├── sprites/                # packed atlases (output of pack_atlas.py)
│   ├── images/                 # source PNGs (input)
│   └── data/
│       ├── scenes/*.json       # authoritative scene definitions
│       ├── npc/ / dialog/      # NPC + dialog data
│       └── sequences/*.json    # cutscene timelines
└── docs/                       # STATUS, FIXME, STORY, CHARACTERS, RETRO_ANALYSIS, ARCHITECTURE, prompts
```

## Disciplines

### 1. Story state — `game.VarStore` (already exists)
Scoped variable store with three scopes: `game` (persistent across the whole
playthrough), `chapter` (reset when a chapter ends), `scene` (reset on scene
change). Canonical keys live in `game/game_state.go`. Save/load via
`VarStore.Save/Load`.

During the refactor, flat `Game` struct flags are bridged to the VarStore via
`Game.syncFlagsToVars` / `syncVarsToFlags`; the endgame is that the flags go
away entirely and the VarStore is the single source of truth.

### 2. Event bus — `game.EventBus` (Phase 1, not yet wired)
Typed synchronous pub/sub. Publishers emit `dialog_ended`, `item_given`,
`city_unlocked`, `kid_healed`, `scene_entered/exited`, `chapter_start/end`;
subscribers register in chapter-wiring files (paris.go, jerusalem.go, ...).

Replaces the callback tree in `setupCampCallbacks` / `setupParisCallbacks`
once Phase 4 lands.

### 3. Scenes — JSON-driven (Phase 3)
Schema: background, walk polygon(s), NPC instances (id + x/y/facing/scale),
hotspots (rect + on-click action), monologue trigger, exits. Loader lives at
`game/scene_loader.go`. JSON files in `assets/data/scenes/*.json` (many
already exist — Phase 3 wires them as authoritative).

### 4. Cutscenes — `game.SequencePlayer` (exists, JSON form in Phase 5)
Step-list runner. Current shape (`game/sequence.go`) is code-authored; Phase 5
adds JSON loading and new step types (`play_animation`, `play_sound`,
`hide_actor`, `fade_to`, `play_monologue`).

### 5. NPC state machine — named states (Phase 4)
Four canonical states per NPC: `default → post → strange → post_strange`.
Each state references a sprite animation and a dialog id. Replaces the
`isStrange` / `dialogDone` / `hintState` boolean soup.

### 6. Interaction rules — Handler+Condition (Phase 4b, infrastructure done)

Implemented in `game/npc_rules.go`. Three concepts:

- `InteractionRule`: `{on, when, do, once}`. Stored per-NPC in `npc.rules`.
- `evalCondition`: a minimal expression evaluator. Supported operators:
  `state == X`, `state != X`, `inv.has(Item)`, `vars.<scope>.<key> == N`
  (plus `!=` / `<=` / `>=` / `<` / `>`), `&&`, `||`.
- `dispatchAction`: action types are `dialog`, `queue_dialog`, `give`,
  `take`, `set_var`, `unlock_city`, `set_state`, `set_strange`,
  `set_npc_silent`, `set_scene_npc_strange`, `emit`.

**Wiring:** `player.startNPCDialog` checks `n.rules` before falling back
to the legacy `onDialogEnd` / `altDialogFunc` closures. Every NPC gets
`n.game` set by `attachGameToNPCs` in `Game.New`, so rules can reach
inventory / varstore / event bus without thread-through plumbing.

**Migration path** (mechanical, per-NPC):
1. Pick an NPC whose behavior is in `setupCampCallbacks` / similar.
2. Author its rule list in `assets/data/npc/*.json` under `rules: [...]`.
3. Populate `n.rules` from the JSON in the NPC factory (or extend
   `npcConfigStore` to do it).
4. Delete the corresponding closure from `setup*Callbacks`.
5. Verify in-game.

Example rule for Jake's Day-2 healing:
```json
{
  "on": "click",
  "when": "state == strange && inv.has(Coin Rubbing)",
  "do": [
    {"type": "dialog", "dialog_id": "jake.heal"},
    {"type": "give", "item": "Coin Rubbing", "to": "jake"},
    {"type": "set_strange", "bool": false},
    {"type": "set_var", "scope": "game", "key": "jake_healed", "value": 1},
    {"type": "unlock_city", "city": "tokyo_street"},
    {"type": "emit", "event": "kid_healed", "kv": ["kid", "jake"]}
  ]
}
```

### 6b. Legacy NPC closures (still in play)
Per-NPC rule list authored in JSON:
```
on: click
when: state == strange && inv.has(paris_postcard)
do:   [{give: paris_postcard, to: marcus}, {dialog: marcus.post_strange},
       {emit: kid_healed, kid: marcus}]
```
Evaluator: `game/npc_rules.go`.

### 7. Sprite atlases — one texture per character (Phase 2)
Source art: existing per-state PNGs under `assets/images/locations/**/npc/**/`.
Manifest: `tools/characters/*.yaml` (per character, declares animations and
grids). Packer: `tools/pack_atlas.py` emits `assets/sprites/<char>.png` +
`<char>.json`. Runtime loader: `game/atlas.go`.

**Baseline alignment.** Source sheets for idle/talk/strange were authored
with inconsistent vertical padding — the character stands at different
cell-y positions from one sheet to the next. The packer fixes this in two
steps so feet line up across animations:
1. Per-frame vertical shift so the lowest opaque pixel lands at
   `cell_h - 1 - baseline_margin` (default margin 8px).
2. Row height normalization: every animation row in the atlas is padded to
   the tallest animation's cell height, content bottom-anchored, so feet
   sit at the same relative Y across rows of the atlas image.

Override with `baseline_align: false` or `baseline_margin: N` at manifest
root or per-animation if a sheet's art is already pre-aligned.

Visual pixels in each frame are byte-identical to the current build — only
packing and per-frame vertical offsets change.

## Current phase status
See `docs/STATUS.md` and the plan file for per-phase gates.

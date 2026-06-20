# SKILL.md — How we work in this repo

A working-conventions guide. The other docs describe *what* the game is; this one describes *how to make changes to it without re-deriving the patterns each session.* Read this before touching anything if you're a fresh Claude session or a returning collaborator.

---

## 1. Project shape

- **Engine:** Go + SDL2 (`github.com/veandco/go-sdl2`). Render loop in `engine/`. Game logic in `game/`.
- **Scenes are JSON.** Files live in `assets/data/scenes/<name>.json`, parsed by `scene_config.go` into `sceneDef`, built into a runtime `*scene` by `scene_loader.go:buildSceneFromDef`.
- **Sequences are JSON.** Multi-phase cutscenes in `assets/data/sequences/<name>.json` driven by the sequence player. Used for night-bedtime, Higgins walk-in, map handoff.
- **Sprites are PNG sheets.** Loaded via `engine.SpriteGridFromPNGClean(path, cols, rows, inset)` for grid sheets or `engine.TextureFromPNGRawClean` for single images. `spriteInset` trims the 1-px transparent gutter most generators leave.
- **NPCs go through the atlas pipeline** (`atlas.go`) — one canonical lookup so any code path that needs a character's idle/talk frames gets the same texture.

## 2. The sprite-generation loop (canonical workflow)

We never draw pixel art in code. The loop is:

1. **Identify** the asset (NPC sheet / player pose / static background / animated background sheet / floor item).
2. **Draft a prompt** in `docs/EXTRA_PROMPTS.md` under a new `§NEW` section. Match the established Hokus-Pokus-Pink cartoon style: thick outlines, saturated palette, transparent PNG, no text, no scenery beyond what the asset is.
3. **User generates** the image externally and drops it into the right folder (see §3).
4. **Wire it up** — atlas entry, scene JSON `npcs` array, hotspot, sequence step, or animated-bg fields.
5. **Smoke-test in-game.** Walk to the relevant scene, trigger the relevant action.
6. **Move the prompt** from `§NEW` to a numbered done section in `EXTRA_PROMPTS.md`.

Memory rule: **don't generate pixel art in code; provide prompts for the user.**

## 3. File / directory map

```
assets/
  data/
    scenes/      — one JSON per scene (camp_grounds.json, paris_street.json, …)
    sequences/   — one JSON per cutscene (night_bedtime.json, higgins_give_map.json, …)
    items.json   — inventory items + their icons
  images/
    locations/<city>/background/   — scene backdrops (PNG)
    locations/<city>/npc/          — NPC sprite sheets (PNG)
    locations/<city>/ambient/      — non-interactive crowd / patron sheets
    player/                        — PP poses and the airplane sheet
docs/
  CHARACTERS.md   — visual identity, pixel ceilings, palette
  STORY.md        — narrative
  ARCHITECTURE.md — engine roadmap (Phase 1–3)
  STATUS.md       — implementation tracker
  FIXME.md        — open work, prioritized
  EXTRA_PROMPTS.md — active prompt queue (§NEW = pending, §1–§N = done)
  PROMPTS.md, HIGGINS_PROMPTS.md, RETRO_ANALYSIS.md — historical reference
  CHANGELOG.md    — commit log
  SKILL.md        — this file
```

## 4. Adding a scene

Minimum JSON:

```json
{
  "name": "my_scene",
  "background": "assets/images/locations/foo/background/my_scene.png",
  "spawnX": 700, "spawnY": 600,
  "characterScale": 1.0,
  "musicPath": "assets/audio/music/foo.mp3",
  "npcs": ["npc_id_from_atlas"],
  "hotspots": [{"name":"door","bounds":{"x":0,"y":0,"w":100,"h":100},"targetScene":"other_scene","arrow":"left"}],
  "blockers": [{"x":0,"y":700,"w":1400,"h":100}],
  "walkSegments": [{"x1":100,"y1":600,"x2":1300,"y2":600}]
}
```

Then **register it** in `scene.go` — find the block of `sm.loadSceneFromJSON(renderer, sceneDefs, "...")` calls (~`scene.go:230` onward) and add a line for the new scene name.

Conventions:
- `characterScale` 1.0 for everything post-2026-05-12 rebalance (the old 0.85 cabin / 0.9 office shrink fudges are gone — see CHARACTERS.md).
- `spawnY` should sit so PP's foot lands on the visual ground line. PP H = 270, so `spawnY` = (ground line Y) − 270. The scene transition clamps spawnY against `playerMinY/playerMaxY` (or scene-specific `minY/maxY` if set), so a too-low spawn auto-corrects.
- `walkSegments` define the horizontal corridors PP can stand on. Without them PP can't move.
- Place a `blocker` rect across any visual obstacle to prevent walk-through.
- `musicPath` — set the canonical path even before the MP3 exists (`audioManager.playMusic` no-ops silently on missing files).

### 4a. Adding a NEW DESTINATION (travel-map pin + scene + story flag)

When the new scene is reachable from the travel map (Stonehenge,
Buenos Aires, Mexico City — anything PP flies to), the work is bigger
than a single JSON. Five layers:

1. **Background art** — paint at 1400×800, transparent characters
   layer NOT included (NPCs render on top at runtime). Drop at
   `assets/images/locations/<region>/background/<scene>.png`.

2. **Scene JSON** — minimum template above. Place at
   `assets/data/scenes/<scene>.json` and register the load call in
   `game/scene.go:newSceneManager` next to the existing `paris_*`
   loads. `name` field MUST match the registration string.

3. **Travel-map pin + landmark** — append a location entry to
   `assets/data/travel_map.json` with:
   - `id`, `name`, `scene` (matches the scene JSON `name`)
   - `pinX`, `pinY` (map coords — check existing pins to avoid the
     90×110 hit-rect overlap, see Rome/Jerusalem near-miss fix)
   - `unlocked: false`
   - `facts: [...]` — 2–3 paragraph strings for the info popup
   - `landmark` — path to a transparent-BG PNG (run
     `tools/clean_landmarks.py` to color-key the white if needed)
   - `audio` — voice-clip path (leave empty until recorded)
   - `relevantWhen` — gate expression like
     `vars.game.<region>_unlocked == 1 && vars.game.<kid>_healed == 0`

4. **NPCs for the scene** — add factories in `game/npc.go` (e.g.
   `newStonehengeDruid(renderer)`), register their atlas under
   `assets/sprites/<region>/<name>.{png,json}` via
   `tools/pack_atlas.py`, and list their ids in the scene JSON's
   `npcs` array. Talk + idle dialogs follow the existing
   `parisGuideDialog` shape.

5. **Story flags + setup callback** — add a `setupStonehengeCallbacks`
   in `game/<region>.go` (see `game/rome.go`, `game/jerusalem.go` for
   the pattern). Wire onDialogEnd handlers that set vars and toggle
   pins in `vars` for the `relevantWhen` chain.

Verification: F1 → dev menu chapter-jump to the new region. PP arrives,
NPCs render at the correct foot line, talk works, story flag flips.

**Tracked want:** Stonehenge specifically (per the PtP "part 3" clip
ending — PP flies from London to Stonehenge for a druid puzzle). When
the art lands, follow the 5-step flow above.

## 5. Adding a sequence (multi-phase cutscene)

Sequences are JSON-authored when phases just chain dialog/move/wait with no dynamic parameters. See `night_bedtime.json` for the full pattern (5 phases, dialog → sleep → freakout → wake → next-day).

**Exception — keep hardcoded:** `flightCutscene` (game/flight_cutscene.go) is a typed Go struct *not* a JSON sequence, because it carries a `destination` parameter the sequence schema can't express today. Don't try to "JSON-ify" it without first adding variable substitution to the sequence player. Same exception applies to any future cutscene that takes a runtime parameter.

## 6. Animated backgrounds

Some scenes need a looping animated backdrop (currently: `airplane_flight` plays a 6-frame cloud sky). The pattern is:

1. **Asset:** vertical sprite-sheet PNG, all frames identical width, stacked top-to-bottom, no gutters. The PNG's full height = `frame_height × frame_count`.
2. **Scene JSON:**
   ```json
   "background": "path/to/sheet.png",
   "backgroundFrames": 6,
   "backgroundFrameSeconds": 0.15
   ```
3. **No code changes** — `scene_loader.go:buildBackground` reads the fields, builds an animated `background`, and `Game.Update` already calls `bg.update(dt)` every frame.

Pick `backgroundFrameSeconds` based on motion speed: 0.10–0.12 for fast/punchy loops (campfire, biplane prop), 0.15–0.20 for ambient drifts (clouds, water shimmer). Lower bound of ~0.08 before it looks twitchy.

Foreground sprites (e.g., the biplane in `airplane_flight`) animate on their *own* clocks — both layers run independently and SDL composites them. No merge step.

## 7. Design rules to respect (excerpt from `CHARACTERS.md`)

- **PP canonical size:** 170×235 px in his standing/idle frames.
- **NPCs:** 160–225 px standing; **never above 230 px** in a scene without a character-scale override.
- **Sky clearance:** keep the top 25% of outdoor scene backgrounds free of clutter so dialog boxes and UI don't fight the art.
- **Per-scene scale multipliers:** cabin 0.85, office 0.9, outdoor 1.0. Set via `characterScale` in scene JSON; respected by the player and NPC draw paths.
- **Palette:** bubblegum pink #E88BB5 for PP, khaki for camp uniforms, saturated cartoon outlines. Unsaturated washes are out-of-style — flag in review.
- **Anchor point consistency:** sprites are anchored bottom-center. If you author a sheet where the character's feet sit at different Y across frames, walking will jitter. Fix at the asset, not in code.

## 8. Common pitfalls

- **Forgetting `spriteInset`** when loading a grid sheet — frames bleed into each other. `engine.SpriteGridFromPNGClean` takes it as the last arg; use `spriteInset` constant unless the sheet was authored with zero gutter.
- **Non-transparent backgrounds** on character PNGs — checkerboard or white squares show up in-game. Reject the asset and re-prompt.
- **Adding fields to `sceneDef` without updating loaders** — the JSON parses fine but the field is unused. Always trace from `sceneDef` → `buildSceneFromDef` → runtime `*scene` and add wiring at each layer.
- **"While I'm here" cleanups.** Don't fix unrelated FIXME items inside an unrelated task. Note them, ship the focused change, queue the rest. (User explicitly prefers this.)
- **Don't write trailing summaries** at the end of a response. Diff is visible; describe only what's next or blocked.
- **Don't generate placeholder pixel art in code.** Always prompt the user instead.

## 8a. Quest-item handoff rule (gameplay)

**An NPC must never give the next quest item until the player has *physically
handed over* the prerequisite item.** Having the item sitting in the bag is
NOT enough — the player has to pull it from the inventory (it rides on the
cursor as `heldItem`) and drop it on the NPC. This keeps the fetch-quest loop
legible: pick up → carry → deliver → receive. User-mandated 2026-06-05 ("we
must bring the item to the person to keep the game moving").

How to wire it:

- For NPCs that use the normal click path (`startNPCDialog` → `canTriggerAltDialog`),
  set **`altDialogRequiresHeld = true`** (plus `altDialogRequiresItem = "<Item>"`).
  Clicking the NPC *without* holding the item falls through to their normal
  dialog (which should hint at what's needed); the trade fires only via the
  held-item drop path in `Game.HandleClick`.
- For NPCs with an `onClickOverride` (e.g. Pierre's walk-up choreography), the
  override bypasses `canTriggerAltDialog`, so gate the trade *inside* the
  `altDialogFunc` on `inv.heldItem` (not `inv.hasItem`).
- Canonical examples: Poulain (rolling pin), Henri (café au lait), Pierre
  (baguette → confiture), Claude (press pass) — all in `setupParisCallbacks`.

Same rule applies to camp (Lily's flower already uses `altDialogRequiresHeld`).

**`altDialogFunc` must be PURE** (2026-06-12 #5): the hover probe in
`ui.go updateHover` CALLS it every frame the cursor passes over the NPC to
decide which cursor to show — any side effect (one-shot anims, state flips,
debug prints) fires on hover, before the player ever clicks. Effects belong
in the *returned callback*, which only runs on a real interaction.

## 8b. Item-acquisition animation rule (gameplay + art)

**Every item the player collects must be visibly acquired — never just appear
in the bag.** User-mandated 2026-06-10. Two sides to wire for EVERY new item:

1. **PP side:** PP plays a pickup/receive one-shot when the item lands in his
   inventory — and per the standing design rule, the final frames POCKET the
   item into his invisible hip pocket (PP ends empty-handed).
   - Floor pickups: a dedicated one-shot (`grab_rolling_pin`, `grab_flower`)
     or the generic grab action (`player.playAction(stateGrabbing, cb)`,
     used by the charcoal pencil).
   - NPC hand-overs: a dedicated receive one-shot (`get_baguette`, `get_jam`,
     `receive_map`).
2. **NPC side:** if a PERSON hands the item over, that NPC needs a *give*
   one-shot too (`poulain.playOneShotAnim("give", 1.0)`,
   `henri.playOneShotAnim("give_jam", 1.0)`) played in the same trade
   callback as PP's receive — the two run in parallel, like the
   Poulain-baguette and Henri-confiture trades in `setupParisCallbacks`.

**This means a new quest item is not "done" at wiring time — it has an art
bill:** 1 item icon + 1 PP receive/grab sheet (or a reusable generic) + 1 NPC
give sheet per giving NPC. Add the missing sheets to `EXTRA_PROMPTS.md` in
the SAME pass that wires the quest, and wire the closest existing one-shot as
a placeholder so the beat isn't silent while the art is pending.

Canonical full examples: rolling pin → Poulain trade (PP `get_baguette` +
Poulain `give`), Henri confiture (PP `get_jam` + Henri `give_jam`), Higgins
map throw (tween + PP `receive_map`).

**GIVING is animated too (2026-06-11 #33):** when PP hands an item to an
NPC, play the generic `"give_item"` player one-shot (sheet `PP give.png`,
§PG1; falls back to the grab frames until it lands). For trades that go both
ways, CHAIN them - give first, receive in its callback:
`game.player.playOneShot("give_item", 0.8, func() { game.player.playOneShot("receive_item", 1.0, nil) })`.

## 8c. Generic pickup lines (gameplay)

**Pickup dialogs stay GENERIC** (2026-06-11 #18): use
`genericPickupDialog(flavor)` (game.go) - it appends a rotating PP quip to
the caller's one-line description of what was found. The "who might need
this" hint belongs in NPC dialogs (Poulain's lost-pin line, Camille's
reminder), NOT in the pocket beat - items stay reusable across quests and
the player keeps discovering uses by talking to people.

**Floor-item stand points (2026-06-12 #14/#15):** PP never stands ON a floor
item. `Game.walkToFloorItem` (the HandleClick pickup path) walks him to a
spot BESIDE the item — feet aligned with the item's base, LEFT of it by
default, RIGHT when the item sets `standRight: true` (pick the side that
puts the grab anim's reach hand over the item; the rolling-pin basket is the
canonical standRight case). On arrival PP squares up FRONT (dir down), so
blocked/observation beats ("the pigeons guard the pot") play facing the
camera. New floor items get this behavior for free — only set `standRight`
when the reach hand needs the other side.

## 9. Doc legend (quick reference)

| Doc | Purpose | When to update |
|---|---|---|
| `CHARACTERS.md` | visual identity, pixel ceilings | when art rules change |
| `STORY.md` | narrative beats | when story scope changes |
| `ARCHITECTURE.md` | engine roadmap | when phases progress |
| `STATUS.md` | impl progress | sprint-level updates |
| `FIXME.md` | open work, prioritized | as work is queued / closed |
| `EXTRA_PROMPTS.md` | active prompt queue | each prompt cycle |
| `CHANGELOG.md` | commit log | each commit |
| `SKILL.md` (this) | how-we-work | when conventions evolve |
| `PROMPTS.md`, `HIGGINS_PROMPTS.md`, `RETRO_ANALYSIS.md` | historical | rarely; archival |

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
  "npcs": ["npc_id_from_atlas"],
  "hotspots": [{"name":"door","bounds":{"x":0,"y":0,"w":100,"h":100},"targetScene":"other_scene","arrow":"left"}],
  "blockers": [{"x":0,"y":700,"w":1400,"h":100}],
  "walkSegments": [{"x1":100,"y1":600,"x2":1300,"y2":600}]
}
```

Then **register it** in `scene.go` — find the block of `sm.loadSceneFromJSON(renderer, sceneDefs, "...")` calls (~`scene.go:230` onward) and add a line for the new scene name.

Conventions:
- `characterScale` 0.85 for cabin interiors, 0.9 for offices, 1.0 for outdoor scenes.
- `spawnY` should sit on the visual ground line; ground line is usually y≈600 in 1400×800 scenes.
- `walkSegments` define the horizontal corridors PP can stand on. Without them PP can't move.
- Place a `blocker` rect across any visual obstacle to prevent walk-through.

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

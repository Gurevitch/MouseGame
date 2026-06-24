# Pink Panther: Camp Chilly Wa Wa

A hand-drawn, point-and-click adventure game in the style of the 1990s
Pink Panther titles *Hokus Pokus Pink* (1997) and *Passport to Peril*
(1996). You play the Pink Panther, a substitute counselor at Camp Chilly
Wa Wa, where the kids start having vivid visions of real-world cities
they've never visited. Travel to each city, find the source of the
visions, and bring back a small "anchor object" that settles each child's
mind.

Built in **Go + SDL2**. Original PP IP belongs to MGM; this is a personal,
non-commercial fan project.

---

## What the game does

- **Story-driven adventure.** Classic retro loop: *encounter → blocked →
  hint → collect item → use item → reward*. Each city is a multi-step
  fetch-quest chain that ends in the anchor object.
- **The camp hub.** Day 1 you meet Director Higgins and the five kids
  (Marcus, Tommy, Jake, Lily, Danny). Overnight the visions begin; Day 2+
  you heal the kids one city at a time.
- **Travel map.** A globe with landmark pins. Cities unlock as the story
  progresses; clicking a pin flies the Pink Panther there in his biplane.
- **Currently playable:** the camp chapters, **Paris** (Marcus / Louvre
  postcard) end-to-end, and **Jerusalem** (Jake / Western Wall coin
  rubbing — trivial heal live, full task chain designed). Tokyo, Rome,
  Rio/Buenos Aires and Mexico City are planned.

See [`docs/STORY.md`](docs/STORY.md) for the full narrative flow and
[`docs/STATUS.md`](docs/STATUS.md) for what's implemented.

---

## Running it

Requires Go 1.21+ and the SDL2 development libraries (the project uses
[`go-sdl2`](https://github.com/veandco/go-sdl2), which links against
native SDL2 / SDL2_image / SDL2_ttf / SDL2_mixer).

```sh
go run .
# or build the executable (Windows icon is embedded via go-winres):
go build -o PP.exe .
```

The game opens fullscreen at a logical resolution of **1400×800**
(letterboxed and scaled to the desktop).

### Controls

| Input | Action |
|---|---|
| **Left click** | Move, talk to an NPC, pick up an item, use the cursor |
| **Drag from the bag** | Carry a held item and drop it on an NPC to give it |
| **Space** | Advance dialog |
| **M** | Toggle the travel map |
| **Esc** | Close the map / open–close the pause menu |
| **F1** | Dev menu (jump to any chapter/scene) |
| **F2** | Click-probe diagnostic (alpha hit-test markers) |
| **F3** | Walk-debug overlay (walk segments + foot/snap coords) |

---

## How it's built

```
main.go            — window + render loop; maps mouse pixels to logical 1400×800
engine/            — SDL2 wrapper: texture load, sprite-grid slicing, bitmap font, easing
game/              — all game logic (flat package, ~35 files)
  game.go            — the Game object (scene wiring, click handling, story flags)
  scene*.go          — JSON scene loader + runtime scene
  player.go / npc.go — actors, animation, movement, blockers, walk segments
  paris.go / jerusalem.go / ... — per-chapter rule + dialog wiring
  travel_map.go, inventory.go, dialog*.go, sequence*.go
tools/             — pack_atlas.py (sprite atlases) + sprite audit/repair tools
assets/
  data/scenes/*.json     — authoritative scene definitions (background, walk
                           segments, blockers, NPCs, hotspots, spawn)
  data/sequences/*.json  — cutscene timelines
  images/                — source PNG sprite sheets + backgrounds
  audio/                 — music + voice (silently no-op if a file is absent)
docs/                — design + working docs (start with the table below)
```

Scenes, dialog and cutscenes are **data-driven JSON**, so most tuning
(NPC positions, walkable paths, hotspots, story beats) happens in
`assets/data/` without recompiling. Character art is AI-generated from
the paste-ready prompts in `docs/EXTRA_PROMPTS.md`; the engine slices each
sheet by detecting the empty gaps between figures, color-keys the white
background, and anchors every frame by the feet so imperfect art still
renders cleanly.

### Documentation

| Doc | What it covers |
|---|---|
| [`docs/STORY.md`](docs/STORY.md) | Full narrative, per-city quest flows, item tables |
| [`docs/STATUS.md`](docs/STATUS.md) | Implementation tracker (engine, scenes, NPCs, cities) |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Engine roadmap, scoped var store, scene/sequence systems |
| [`docs/CHARACTERS.md`](docs/CHARACTERS.md) | Visual identity, canonical draw sizes, palette |
| [`docs/SKILL.md`](docs/SKILL.md) | How to make changes (add a scene, sprite, sequence) |
| [`docs/FIXME.md`](docs/FIXME.md) | Open issues, prioritized |
| [`docs/EXTRA_PROMPTS.md`](docs/EXTRA_PROMPTS.md) | Paste-ready art-generation prompts |
| [`docs/RETRO_ANALYSIS.md`](docs/RETRO_ANALYSIS.md) | Architecture notes from the original PP games |

---

## Status

Active development. Camp + Paris are complete; Jerusalem is partway; the
remaining cities are designed but not yet built. The refactor toward a
fully data-driven engine (JSON scenes/sequences, a scoped variable store,
declarative NPC rules) is tracked in `docs/STATUS.md` and
`docs/ARCHITECTURE.md`.

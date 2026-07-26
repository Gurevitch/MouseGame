# ClonedPP — Pink Panther: Camp Chilly Wa Wa

A 2D point-and-click adventure (Go + SDL2). PP is a substitute camp counselor;
the kids have visions of real cities, so PP flies to each city, solves a retro
item daisy-chain, and brings home an "anchor object" that heals the kid.
Built: Camp, Paris (Marcus), Jerusalem (Jake), Tokyo/Kyoto (Lily). Next: Rome
(Danny) / Rio (Tommy).

## Build / run / test (Windows + AV)

- **Always set `GOTMPDIR` into the project** (the antivirus blocks temp-folder
  executables): `GOTMPDIR=.gotmp go build ./...`
- Test verbosely so per-test PASS lines show: `go test -v ./...`
- To run/playtest the app use the **`run`** skill; to verify/repair sprite
  sheets use the **`sprite-check`** skill.

## Working conventions (important)

- **Don't auto stage/commit/branch.** Edit the working tree only; leave staging
  and commits to the user. "New PR" means a new *batch of work*, not a literal
  git PR. Work on `main`.
- **Never generate pixel art in code.** New sprites/BGs are AI-image **prompts**
  written to `docs/EXTRA_PROMPTS.md`; the user generates the PNG, drops it in,
  then we run `sprite-check`. Wire code with **graceful fallbacks** so the build
  runs before art lands (missing sheet = no-op / placeholder).
- **Art rules** (also in `docs/EXTRA_PROMPTS.md` header + the sprite-check
  skill): **NEVER use pure white `#FFFFFF` ANYWHERE on a sprite** — not on
  characters, props, canvases, aprons, paper, teeth, or **eyes** (the engine
  chroma-keys white and punches it into holes). Use cream `#E5DDC8` / vanilla
  `#F2EFE5` / pale-grey `#C4C4C4` instead; only the scene-cell background stays
  pure white. `SPRITE_SCAN=1 go test ./engine -run TestSpriteScan` flags any
  white-on-character. Also: every PP pickup ends with him pocketing the item,
  and every PP GIVE **starts with him pulling it out of that pocket**;
  PP has plain pink paws (no gloves); every speaking NPC gets **separate idle +
  talk** sheets (not one packed sheet); ≥15px gaps between figures so frames
  gap-cut cleanly.
- **Sprite grids:** most one-shot/idle/talk sheets are 1536×1024 loaded **8×1**;
  verify cuts with `go test ./engine -run ContentGrid` + `go run ./tools/jitter_audit .`
  and register new sheets in both manifests.
- **Item/give beats** use the two-stage `handOff` (PP gives → NPC receives → NPC
  gives back → PP receives); per-item `give_<key>` / `receive_<key>` sheets
  auto-upgrade the animation. See `game/player.go playHandOff`.
- Chapter wiring lives in `game/<city>.go` (paris.go, jerusalem.go, tokyo.go…);
  scenes load NPCs/hotspots; save-safe flags live in `game/game_state.go`
  (VarStore). Architecture detail in `docs/ARCHITECTURE.md`.

## Project docs

Core design docs are imported below (auto-loaded each session):

@docs/ARCHITECTURE.md
@docs/STORY.md
@docs/JAPAN_TASKS.md
@docs/JERUSALEM_TASKS.md

Read these on demand (too large / situational to auto-load):

- `docs/EXTRA_PROMPTS.md` — all pending AI-image prompts + the standing art
  rules. Read/append whenever generating or fixing sprites.
- `docs/FIXME.md` — the running bug log (check off items as fixed, with the date).
- `docs/CHARACTERS.md` — character bios/designs; read when writing dialog or NPC art.
- `docs/PROMPTS.md`, `docs/HIGGINS_PROMPTS.md` — older prompt archives.
- `docs/RETRO_ANALYSIS.md` — design patterns from the original PP games.
- `docs/STATUS.md`, `docs/CHANGELOG.md` — phase status / change history.
- `docs/SKILL.md` — sprite/scene authoring workflow notes.

---
name: sprite-check
description: Verify and repair this game's sprite sheets (PP + NPCs) using the project's audit/clean/repair pipeline. Use this skill whenever the user mentions checking or verifying sprites/sheets/animations, jitter, characters jumping or sliding, cut or split frames, floating/ghost limbs or duplicate hands, missing colors or see-through holes in a character, a newly generated sprite sheet landing on disk, art regen results, or wants the art verified before a PR — even if they don't say "sprite-check" explicitly.
---

# Sprite Check — verification & repair pipeline

This project renders AI-generated sprite sheets through a forgiving engine
pipeline (built 2026-06-10). Knowing what the engine already fixes is the key
to interpreting the checkers correctly — most "scary" numbers are harmless,
and the two real problems are easy to spot.

## What the engine already handles (don't fix these by hand)

- **Gap-based slicing** (`engine.contentGridRects`): frames are detected by
  the empty background gaps between figures and cut at gap midpoints, so
  figures are never sliced even when off-center in their cells. Figures
  bridged by a thin touch (tail tip on a neighbour) are separated at the
  thinnest waist. Only when the expected figure count can't be resolved does
  it fall back to proportional grid lines.
- **Per-frame feet anchoring**: every frame's feet (wide bottom row + dense
  bottom-band centre, tail excluded) are detected at load and pinned to the
  ground line at draw time, with a ±6px deadband snapped to the sheet median.
  Art that drifts inside its cells renders perfectly still anyway.
- **Non-divisible dimensions** (e.g. a 1535px-wide 8-column sheet): boundary
  remainders are distributed proportionally. Never pad or resize a sheet.
- **Missing sheets**: loaders warn and degrade to an invisible animation —
  the game never panics on absent art.

## The check (run both, in this order)

```
go test ./engine -v -run ContentGrid    # which sheets gap-detect vs fall back
go run ./tools/jitter_audit .           # ghosts, cross-border, drift, empty cells
SPRITE_SCAN=1 go test ./engine -run TestSpriteScan -v   # pure-WHITE-on-character + leak report
```

`TestSpriteScan` reads every registered sheet once (no duplicate images) and
lists sheets with pure-white-on-character (the chroma-key punches those into
holes) and full-width "leak" frames. NOTE: it flags INTENDED whites too
(a white prop, eye sclera that should already be off-white) — eyeball the top
hits; the fix is almost always an art regen with the no-pure-white rule, not a
code change.

**Background residue** (blue/white spots around a figure or trapped between
arm and body): the art keeps being GENERATED on flat `#B4D7EE` — that stays
the style. Fix residue by BAKING the bg out of that PNG with `tools/bg_bake`
(allowlist + per-sheet mode inside `tools/bg_bake/main.go`) — but ONLY for
sheets with a reported, verified problem; a proactive sweep punched pixels
out of PP's ivory and was rolled back (2026-08-01). Strict `pocketTol: 6` on
near-white bgs, and ALWAYS `Read` the PNG after baking. Baked (transparent)
sheets stay baked in the repo — runtime keys no-op on them; reverting brings
the spots back. Currently baked: PP talk side, pp_sleeping, pp_waking,
npc_pierre_get_baguette, npc_bagel_seller_give.

When a NEW sheet lands on disk, add it to BOTH manifests first:
`engine/grid_content_test.go` (the cases list) and
`tools/jitter_audit/main.go` (the sheets list), with the grid the loader
uses — grids live in `game/player.go` / `game/npc.go` factories; the
`assets/data/npc/*.json` grids are NOT used by the engine.

## Interpreting results

| Finding | Meaning | Action |
|---|---|---|
| `GAP-DETECTED`, no ghosts | Sheet is good | Nothing |
| `fallback (no clean gap split)` | Figures touch with no thin waist — proportional cutting may clip crossing limbs | Queue a re-roll with the ≥15px-gap rule (see prompts below); until then it renders as before |
| `GHOST PIECES` on a **prop-free** sheet | Generator painted a detached duplicate limb inside a frame — visible in-game | Run `tools/sheet_clean` (see safety rules) |
| `GHOST PIECES` on a sheet with a legit prop (thrown map, handed baguette/jam/cup, pigeon, received item) | Usually the prop itself | `Read` the PNG and judge visually — do NOT auto-clean |
| `CONTENT CROSSES` borders | A limb spans the mathematical boundary; harmless if the sheet gap-detects (the cut moves to the gap), real if it falls back | Covered by the fallback/re-roll decision |
| `FOOT drift` / `CENTER-X drift` | Art quality only — the renderer cancels positional drift via feet anchoring | Note for re-roll priority; not urgent |
| `EMPTY cells` | Blank frames blink in the loop | Re-roll |
| Missing colors / see-through holes on a PLAYER sheet in-game | Enclosed pure-white regions get globally color-keyed | Run `tools/sheet_repair` |

## Repair tools — safety rules

- `go run ./tools/sheet_clean .` keeps only content inside/around the main
  body per cell (v2: never touches interior details like belly shading).
  It runs ONLY on its allowlist inside `tools/sheet_clean/main.go` — edit
  the list to add a flagged sheet, and only add **prop-free** sheets (a
  character whose every part connects to one body). Never add sheets with
  legit separate objects.
- `go run ./tools/sheet_repair .` refills enclosed pure-white holes with the
  surrounding color — for PLAYER sheets only (they use a global color key).
  Its list is inside `tools/sheet_repair/main.go`.
- After ANY clean/repair: re-run both checkers AND `Read` the PNG to verify
  visually — a tool pass that "succeeded" can still look wrong (this caught
  straddled Marcus figures that no metric flagged).
- If a regenerated sheet turns out unusable (figures straddling cells), the
  old art is in git: `git restore --source=HEAD --worktree -- <path>`.

## When art needs a regen

Point the user at `docs/EXTRA_PROMPTS.md` — paste-ready prompts live there
(§JIT batch for known re-rolls). Every prompt must carry the standing rules,
which exist because each one fixes a bug we actually shipped:

- **Separate idle + talk sheets per NPC** (user 2026-06-24, firm rule): every
  speaking NPC ships TWO files — `<name>_idle.png` and `<name>_talk.png` — same
  character/size/anchor across both. NEVER pack idle on row 0 + talk on row 1 of
  one sheet; author them as two separate single-row strips. Loaders look for the
  talk sheet and fall back to the idle until it lands, so a missing talk sheet is
  a queued prompt, not a bug.
- **One character per cell, ≥15px empty gaps** between figures and to sheet
  edges — gaps are what the engine cuts at; clear gaps = uncuttable frames.
  No ghost/duplicate limbs, no figures straddling cells.
- **Anchor lock** — feet on the same pixel row, centerline on the same
  column, same size every frame; limbs/tail animate around the anchor.
- **Storyboard the motion** (frame-by-frame beats) and say "every frame must
  be CLEARLY different" — otherwise generators output one frozen pose ×16.
- **NO pure white `#FFFFFF` ANYWHERE on a sprite** (hard rule, user
  2026-06-30) — not on characters, props, painter canvases, aprons/smocks,
  paper, teeth, or **eye sclera**. Pure white gets chroma-keyed into holes /
  renders as white spots. Use cream `#E5DDC8` for fabric/canvas, ivory
  `#F2EFE5` for PP's belly, pale-grey `#C4C4C4` for eyes/steam. Only the
  scene-cell BACKGROUND stays pure white. Verify with `TestSpriteScan` above.
- **PP specifics**: plain pink paws, NO gloves; every pickup ends with PP
  pocketing the item into his invisible hip pocket, and every GIVE starts
  with him pulling the item back OUT of that pocket (user 2026-07-24) —
  frame 1 reaches to the hip, frame 2 the item appears in his paw, then
  the offer. Never open a give with the item already in hand.
- Seated/behind-counter characters (office Higgins, Poulain): upper body
  only, no desk/counter in the sprite, all action at chest height, anchor by
  the waist cutoff row.

After the user generates a sheet: save to the exact path, re-run this skill's
check, and `Read` the PNG. A good sheet shows GAP-DETECTED with no ghost
warnings.

## Report format

End with a short table: sheet → status (clean / cleaned now / needs re-roll /
visual-check) → action taken or queued. Log open re-rolls in `docs/FIXME.md`
and queue prompts in `docs/EXTRA_PROMPTS.md` (move finished ones to its Done
log). Remember SKILL.md §8b when new items are involved: every acquisition
needs a PP receive/grab anim AND an NPC give anim.

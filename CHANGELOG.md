# Changelog

Short per-session log. One section per working pass. Most recent on top.
See `FIXME.md` for open issues and `STATUS.md` for feature-level progress.

## 2026-04-17 — Polish pass 4 (Higgins geometry + positions)

- **Higgins sprite geometry** — aligned every Higgins loader in
  `game/npc.go` to the authoritative grid sizes in PROMPTS.md:
  entrance talk `8x2 row 0` → `6x1`, office idle `6x2` → `7x1`,
  office talk `6x2` → `4x2`, night talk `4x2` → `6x1`. This is the
  real fix for "doubled Higgins" on talk.
- **Office Higgins position + size** — bounds `(942, 357) 240x320`
  so he sits at the user-requested `(1062, 357)` area and reads
  roughly 1.3x his old size.
- **Night Higgins** — moved into the bottom-right corner of the
  campfire frame (`(1120, 430) 200x260`) so the bedtime-speech beat
  has him already on-screen instead of walking in.
- **Marcus room** — bounds `(526, 181) 280x380` in `scene.go` so he
  fills more of the cabin with his foot center near `(666, 561)`.
- **Cursor cleanup** — the proposed `cursorUse` state and
  `cursor_use.png` prompt are removed; only existing cursor PNGs
  remain in use.

## 2026-04-16 — Polish pass 3 (cursor / Higgins / Lily / night)

- **Cursor** — `updateHover` now defaults to `cursorGrab` whenever PP
  is carrying an item so the pointer itself shows the held-item state
  (previously only the ghost icon beside the pointer did). No new
  cursor PNGs were added — the existing set is reused as-is.
- **Higgins idle double-render** — `game/npc.go:247` was loading the
  idle sheet as `7x1` while the PNG is actually multi-row, so each
  "frame" blended half of two neighboring poses. Swapped to
  `loadNPCGridRow(..., 8, 2, 0)`, matching the talk loader. Asset
  regen request logged in PROMPTS.md.
- **Lily flow** — replaced the closure-local `lilyHinted` flag with a
  per-NPC `hintState` field on the `npc` struct. Lily's altDialog is
  now armed at setup but gated internally on `hintState == 1`, so the
  first click always plays `lilyShyDialog` regardless of what's in the
  bag, and the flower handoff fires exactly once after the shy beat.
  Survives scene re-entry cleanly.
- **Night flow** — added `nightHidePlayer` flag on `Game`. Set true
  during phase 3 (Marcus's cabin) and false when transitioning back to
  the campfire. `Draw` now calls `scene.drawActorsNoPlayer` in that
  window so PP is no longer visible in Marcus's room while he's
  freaking out. Also gated `drawWarmTint` off during the campfire
  sleep so the orange overlay stops bleeding the sleep sprite.
- **Docs** — added this changelog. Appended cursor_use and Higgins
  idle regen requests to PROMPTS.md. FIXME.md swept (see below).

# Changelog

Short per-session log. One section per working pass. Most recent on top.
See `FIXME.md` for open issues and `STATUS.md` for feature-level progress.

## 2026-06-05 — Playtest pass (museum + Paris flow)

Engine/JSON side of a 32-item playtest list. Art-bound items (sprite re-cuts /
regens) are queued in `EXTRA_PROMPTS.md` under "Playtest pass — museum".

- **Inventory + cursor** — inventory bag oval enlarged (816×680); redundant
  left/right chevron logos removed (the inv_circle art already has hand grips);
  open-bag hit test tightened to PP's body so it no longer opens from a wide
  radius; the active "pointing" cursor now shows while carrying an item and
  over the relevant travel-map pin (other pins keep the talk icon).
- **Camp** — PP walks into camp centre (755,533) on arrival; Higgins-office
  arrow moved to ~1186,692; room Marcus shrunk (150×205) so he reads shorter
  than PP; office Higgins flipped to face PP.
- **Paris flow** — item trades now require *handing the item over* (held), not
  just carrying it (Poulain/Henri/Pierre/Claude); Pierre talks to the side and
  eases back to full size after dialog instead of popping; Claude talk cadence
  slowed; bakery exit walks PP to the door and back through it (not off-right);
  Poulain/Camille repositioned.
- **Rolling pin** — now hidden inside the bike basket (~539,644): no sprite,
  but the grab cursor reveals it; pickup plays the dedicated grab-rolling-pin
  one-shot. New `floorItem.hidden` mode.
- **Museum** — first arrival walks PP in from the left tunnel (381,481) with a
  one-time monologue; scene `characterScale` 0.7 shrinks PP + Beaumont;
  Beaumont flipped and repositioned (~546,599); bottom-right travel button
  removed; getting the postcard unlocks a "fly back to Camp" travel pin.
- **Verified** — night shout sprite is correctly wired (`camp_night` →
  `night_higgins`); its broken playback is an art issue (queued §SH).

## 2026-04-17 — Polish pass 5 (FIXME batch)

- **Strict NPC click bounds** — `npc.containsPoint` no longer expands
  by 70/50 px. Fixes Danny snap-stealing Marcus clicks and dialogs
  firing from empty ground behind an NPC.
- **Map reveal removed** — Higgins handing over the Travel Map no
  longer triggers a grow-onto-screen zoom. The existing take-map
  animation is the whole beat.
- **Higgins appears for Lily** — new hidden `Director Higgins` NPC on
  `camp_grounds` (`(910, 400)`); unhides when Lily's shy dialog ends
  and delivers the flower clue via `higginsLilyHintDialog`.
- **Bedtime beat restored** — `checkDay1Complete` plays a short
  Higgins bedtime dialog on `camp_grounds` before fading to
  `camp_night`, latched behind `day1BedtimeStarted` so the Lily
  flower callback can't double-trigger it.
- **PP sleeping size + halo** — draw scale 1.8 → 1.1; sleep/wake
  sheets now use `SpriteGridFromPNGCleanAggressive` with inset 4 to
  strip the cream-white rim.
- **Flight hides PP** — `Draw` skips the player actor while in
  `airplane_flight` so PP stops standing on top of the biplane.
- **Office Higgins** — bounds top-left clamped to the user's spec
  `(1062, 357)` with size `220x280`.
- **Kids smaller than PP** — every `camp_grounds` kid bound
  normalized to `150x180` (PP stays at `170x235`).
- **`npc.hidden`** — new flag added and honored by `drawScaled`,
  `scene.checkNPCClick`, and `updateHover` so story-timed NPCs can
  live in the scene list before their cue.

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

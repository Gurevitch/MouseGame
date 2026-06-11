# Jerusalem chapter — task design (plan only, not built)

Status: **design notes.** The chapter currently works the trivial way (walk up to
Miriam, she hands over the Coin Rubbing on first talk). This doc is the plan for
turning it into a proper retro daisy-chain like Paris. We refine it here, then build.

## Theme

Jake's arc is **courage**. His nightmare is the face in the tunnels — a Roman
coin's portrait staring out of the Western Wall. The gameplay should make *PP*
(and Dov, the scared kid who mirrors Jake) face the wall, so the puzzle echoes
the emotional beat.

## Hard rule (user, 2026-06-08)

**The Coin Rubbing — Jake's "key" — is NOT given in the first task.** Like Paris
(rolling pin → baguette → coffee → confiture → press pass → ticket → postcard)
and the retro PP games, Jerusalem is a chain of several *tangential, slightly
absurd* errands (an alley cat, a tourist's sardines, a lost guidebook) that only
pay off in the anchor object at the END. The note-in-the-wall is the FINAL task,
and placing it is what unlocks the flight home.

## Goal / heal item — UNCHANGED

The **Coin Rubbing** is what PP carries home to heal Jake (existing wiring in
`game/game.go`: Jake `altDialogFunc` → `setStrange(false)` + `jake_healed` +
unlock Tokyo). Don't touch the heal mechanism — just gate *earning* the rubbing
behind the chain below.

## Firm decisions (from the user)

1. Coin Rubbing heals Jake; it's the **end reward**, not a first-task handout.
2. The chapter needs **several tangential errands** (retro adventure rhythm), not
   one fetch.
3. **Paper** and a **pen/pencil** come from **two separate NPCs**, each handed
   over (with the rest of the chain) — and PP pockets them **before he knows what
   they're for**, with offhand lines:
   - Paper (PP): *"I don't really know what to do with this currently, but oh well..."*
   - Pen (PP): *"A pen, too. Still no idea what for, but a panther never turns down free stationery."*
4. PP then **writes a note** and **puts it in the wall** (the real Western Wall
   prayer-note custom) → two new PP one-shot sprites: **write** + **put-in-wall**.
5. The note ritual is the **FINAL task**. Placing it **unlocks the travel map's
   flight back to camp** — PP can't leave Jerusalem until the note is in the wall.

## The daisy-chain (draft)

Anchor NPC = **Miriam** (archeologist at the wall), like Beaumont/Poulain in Paris.

```
Arrive at jerusalem_entrance (plaza). Ways on: LEFT arch -> market, RIGHT -> wall.
    |
    v
WALL — Miriam (initial): "A rubbing of the Hadrian coin, for a frightened boy?
  I can do that — but I'm pinned holding this survey rod and my charcoal stick
  is gone. Little Dov ran off with it. Get it back and I'll make your rubbing."
    -> need: Miriam's charcoal stick (Dov has it)
    |
    v
WALL — Dov (scared kid): "Charcoal? I was drawing glow-bugs with it... but I
  dropped my bug jar and the alley cat knocked it under the spice table. I'm
  not going near that cat. Get my jar back and the charcoal's yours."
    -> need: Dov's glow-bug jar (in the market, guarded by the cat)
    |
    v
MARKET — Eli (spice seller): "Your jar? Ha — under my saffron table, and that
  thieving cat guards it like a dragon. Only one thing moves that cat: fish.
  Ask the tourist, he hoards airplane snacks."  (Eli also, in passing, presses
  a sheet of spice-wrapping PAPER on PP as thanks for the company.)
    -> need: something fishy to lure the cat
    -> PICKUP: Paper  (PP: "...dunno what for, but oh well." [pockets it])
    |
    v
MARKET — Gary (tourist): "Sardines? Sure, take the tin — airline food, blech.
  But first, find my guidebook, I put it down somewhere round here." He's been
  sketching the wall, so when he's done he hands PP his PENCIL too.
    -> need: Gary's lost Guidebook (floor item in the market/plaza)
    |
    v
Find the Guidebook (floor item) -> give to Gary -> get Sardine Tin (+ Pencil)
    -> PICKUP: Pencil  (PP: "...still no idea what for." [pockets it])
    |
    v
MARKET — lure the cat with the Sardine Tin -> cat leaves -> grab the Glow-bug Jar
    |
    v
WALL — Dov: return the Glow-bug Jar -> get Miriam's Charcoal Stick
    (Dov, watching PP poke around the dark wall, starts to lose his own fear)
    |
    v
WALL — Miriam: give the Charcoal Stick -> she makes the COIN RUBBING (Jake's key)
    Then she names the custom: "One last thing — you never take from the Wall
    without leaving something. Write a note and tuck it in the stones. Then go home."
    |
    v
PP realizes the paper + pen were for THIS.
WRITE the note (PP one-shot sprite #1)  -> PP: "...Oh. THAT's what the paper's for."
    |
    v
PUT the note in the wall (PP one-shot sprite #2). Dov copies him, brave now.
    |
    v
Note placed -> the travel map UNLOCKS the flight back to camp.
    |
    v
Use Travel Map -> fly to Camp Chilly Wa Wa
    |
    v
CAMP — Jake: give the Coin Rubbing -> heal (existing wiring)
```

## Cast / roles (jerusalem_entrance + jerusalem_market + jerusalem_wall)

| NPC | Scene | Role |
|-----|-------|------|
| Miriam | wall | **QUEST anchor** — makes the Coin Rubbing once she has her charcoal; teaches the note custom |
| Dov | wall | **QUEST** — trades the charcoal stick for his glow-bug jar; courage mirror of Jake |
| Eli (spice seller) | market | **QUEST** — cat-at-the-stall errand; hands over the Paper |
| Gary (tourist) | market | **QUEST** — trades the Sardine Tin (+ Pencil) for his lost Guidebook; comic relief |
| Worshippers | wall | ambient sway overlay (no quest) |
| Alley cat | market | obstacle/prop — guards the jar, lured by sardines |

## Items in this chapter

| Item | Source | Used on | Result |
|------|--------|---------|--------|
| Guidebook | Market floor item (Gary dropped it) | Gary | Trade for Sardine Tin (+ Pencil) |
| Sardine Tin | Gary | the alley cat | Lures the cat off the jar |
| Glow-bug Jar | Under Eli's stall (after cat leaves) | Dov | Trade for the Charcoal Stick |
| Charcoal Stick | Dov (was Miriam's) | Miriam | She makes the Coin Rubbing |
| Paper | Eli (offhand) | — | Written into the Note |
| Pencil | Gary (offhand) | — | Writes the Note |
| Note (written) | PP writes it (paper + pencil) | the wall crack | Placed in the wall → unlocks return flight |
| **Coin Rubbing** | Miriam | **Jake (camp)** | **Heals Jake** (existing) |

## New assets needed

**Items (registry entries + inventory icons → art prompts):** Guidebook, Sardine
Tin, Glow-bug Jar, Charcoal Stick, Paper, Pencil, Note. (Coin Rubbing exists.)

**Alley cat sprite:** small idle + "walk off" frames (a market prop). Could reuse
the ambient-sprite system (ambient_sprite.go) as a `sway`/`travel` mover, or a
proper NPC if it needs click interaction.

**PP one-shot animations (→ EXTRA_PROMPTS.md):**
- `PP write note` — PP writing on paper, ends pocketing it (standing PP rule).
- `PP put note in wall` — PP reaching up and tucking the note into a crack.

## Code touchpoints (when we build)

- `game/jerusalem.go` — NPC give/receive wiring (mirror the Paris bakery trade
  pattern: `altDialogFunc` / held-item gating), Miriam's two-stage dialog, the
  trade-chain hint states.
- A **floor item** in `jerusalem_market` for the Guidebook (pattern: paris_street
  bicycle-basket rolling pin).
- The **alley cat** prop + its lure interaction.
- A wall-crack **hotspot** in `jerusalem_wall` with `onInteract` = place-note,
  gated on having the written Note.
- **One-shot anim** wiring: `player.playOneShot(...)` for write + put (pattern:
  PP get_baguette / grab_rolling_pin).
- New items in the item registry (`game/items*`), icons under `assets/images/items/`.
- VarStore flags for save-safe gating: e.g. `jer_have_charcoal`, `jer_coin_made`,
  `jer_note_placed`.
- **Return-flight gate:** the camp travel pin (currently relevant once
  `jerusalem_unlocked && jake_healed == 0`) must ALSO require `jer_note_placed`
  for the Jerusalem leg, so the map won't fly PP home until the note is in the
  wall. Edit the camp pin `relevantWhen` in `assets/data/travel_map.json` (and/or
  the hitTest/flight path in `game/game.go`). Make sure every step (paper, pen,
  charcoal, jar, sardines, guidebook, note) is reachable from inside Jerusalem so
  PP can't soft-lock.

## Open questions (to refine together)

1. **Cat handling** — pure prop lured by sardines (auto), or a clickable NPC? A
   simple prop that walks off when PP "uses sardines on it" is least work.
2. **Whose note** — does PP write down *Jake's* fear to leave it at the wall, or
   his own? (Affects the writing-beat dialog.) Strong theme: leaving Jake's fear.
3. **Dov's courage payoff** — should placing the note visibly change Dov (he
   places his own, or thanks PP), as the on-screen echo of Jake's healing?
4. **Trim vs keep length** — current chain is ~6 steps (guidebook → sardines →
   cat → jar → charcoal → rubbing → note). Matches Paris density. Trim a link if
   playtest feels long.
5. Exact NPC/floor-item/hotspot positions per scene, fixed once steps are locked.

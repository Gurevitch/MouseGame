# Kyoto / Japan chapter — quest design

Status: **BUILT (2026-06-24)** — the full 8-step chain below is wired in
`game/tokyo.go` (logic + dynamic stall + queue-sit), builds + tests green. Runs
on placeholder art; the new art is queued in EXTRA_PROMPTS §JP. Decisions locked:
**Kiku = optional gag** (not a gate); **the well is a hotspot in the ramen
street**; **the stall stays open with the line seated** after Hiro serves PP.

## Theme

Lily's arc is **voice / being heard**. She's gone quiet inside; she loved
flowers and hears far-off bells. The anchor is a blossom from the grove's
**"Whispering Cherry"**, an old tree that only blooms for someone who brings it
a proper **offering** — so the chapter is about *earning the right to ask*, with
small tangential Kyoto errands along the way (a dead hearth, a dry ink-brush,
well-water, a lost fire-striker) like the Paris press-pass and Jerusalem note
chains.

## User decisions (2026-06-24)

1. **6–8 step chain** (not a light hand-over).
2. **Dynamic ramen store**: starts **CLOSED** with a **static line** waiting
   outside; when PP gets it to **OPEN**, the waiting people **sit down at the
   counter right away**. So we need: a store open/closed prop + queue NPCs with a
   standing→sitting swap.

## Revision (2026-06-24, round 2)

- **Gary opens the stall**, not the fire-striker. His overjoyed "you MUST taste
  the ramen" tip flips `jp_ramen_open` (line sits) the moment you talk to him at
  the torii. The book-upside-down gag plays mid-chat.
- The **fire-striker is "what you give for the BLESSED bowl"** (locked): the
  stall is open for normal ramen, but Hiro needs his striker to light the SACRED
  hearth for the offering bowl. Keeps the Kenji→well→Oba-chan→striker chain.
- **Higgins** comes OUT of his office (bottom-right corner) and strides across
  to PP for the rude intercept.
- **Tea ceremony → its own temple tea-house** (`tokyo_teahouse`, a 5th Kyoto
  scene, reached UP from the flower store). Tea master + ceremony live there;
  matcha/bowl supplies stay in the flower store; the whisk is at the street well.
  (Authentic — the ceremony grew out of Zen temple tea rooms.)
- **Items consumed in-place** (user rule): all Kyoto items are now used up by
  chapter end — incl. the Voice Charm, now placed at the old tree with the
  offering, so nothing lingers in the bag.

## The chain — "The Whispering Cherry's Offering" (8 steps)

```
(arrive Kyoto at the torii, post rude-Higgins)
    |
1.  TORII — Gary: flips his guidebook right-way-up and reads the custom aloud -
    "ze old cherry in ze grove blooms only for an offering blessed at ze ramen
    hearth." Points PP down to the street. (clue; no item.)
    |
2.  STREET — Hiro, stall CLOSED, a static line waiting: "Bless an offering? I
    would - but my hearth is COLD. My fire-striker is gone, a crow took it. No
    fire, no broth, no opening... and look at zis line!"
        -> need: Hiro's fire-striker.
    |
3.  STREET — Kenji: "Ze crow? I saw it. I'll tell you where it dropped - but my
    ink has dried and I cannot work. Bring me water from ze temple well."
        -> need: Well-Water.
    |
4.  TORII (or street) — the temple WELL (hotspot): draw a cup of Well-Water.
        -> PICKUP: Well-Water.
    |
5.  STREET — give Well-Water to Kenji -> he wets his ink, tells PP the striker
    fell in the FLOWER-STORE eaves, and brushes a tiny "voice" charm as thanks
    (ties to Lily).
        -> get: Voice Charm. clue: striker is at Oba-chan's.
    |
6.  FLOWER STORE — Oba-chan: "Ze crow's little treasure? On my shelf. Take it."
        -> get: Fire-Striker.
    (Optional gate — Kiku: "You cannot meet ze Whispering Cherry dressed like
     ZAT!" -> the kimono-spin gag becomes a soft requirement before the grove.)
    |
7.  STREET — give Fire-Striker to Hiro -> THE STORE OPENS: window/noren slides
    up, the waiting LINE SITS at the counter, Hiro lights the hearth and ladles
    a blessed Offering Bowl for PP.
        -> get: Offering Bowl.  [the dynamic open->sit beat]
    |
8.  FLOWER STORE — Oba-chan: now PP carries the Offering Bowl (+ Voice Charm):
    "Now you bring something to give. Follow me." -> opens the grove path.
    |
    GROVE — place the offering at the old tree -> it blooms -> PICK the
    Pressed Sakura. -> Danny's phone call -> fly home -> heal Lily at the lake.
```

Step count: Gary clue (1), Hiro ask (2), Kenji ask (3), well (4), Kenji trade +
charm (5), Oba-chan striker (6), Hiro open + bowl (7), Oba-chan leads + grove
pick (8). Kiku kimono is an optional 9th soft-gate.

## Items

| Item | Source | Used on | Result |
|------|--------|---------|--------|
| Well-Water | temple well (hotspot) | Kenji | unlocks his clue + Voice Charm |
| Voice Charm | Kenji | the offering / Lily flavor | carried to the grove |
| Fire-Striker | Oba-chan's shelf (crow dropped it) | Hiro | opens the store |
| Offering Bowl | Hiro (once the hearth's lit) | the old tree | tree blooms |
| Pressed Sakura | the old tree (pick) | Lily (camp lake) | heals Lily (existing) |

## Dynamic ramen store + sitting queue (tech plan)

- **Store state** — VarStore flag `jp_ramen_open` (default false). A 2-state prop
  drawn over the stall: CLOSED (shutter down / noren furled, dim) ↔ OPEN (noren
  up, lantern lit, steam). Toggled in Hiro's give-fire-striker callback. Same
  pattern as the leaf ambient but state-driven, not looping.
- **Queue NPCs** — a static row of ~4 waiting customers outside the stall
  (ambient `sway`). Each has a STANDING sheet and a SITTING sheet. On
  `jp_ramen_open` flipping true, swap each to its sitting sheet and reposition to
  a counter stool (x near the window, lower y). Implement either as ambient
  sprites with a swap, or as lightweight NPCs; ambient is cheaper.
- Save-safe: on load, if `jp_ramen_open`, build the stall already open + people
  seated.

## New art needed (→ EXTRA_PROMPTS.md, all graceful-fallback)

- Ramen stall **open vs closed** prop (2 states: shutter/noren down dim ↔ up lit + steam).
- **Waiting customer** standing (1-2 variants) + **sitting-at-counter** sheets.
- Item icons: Well-Water, Voice Charm, Fire-Striker, Offering Bowl.
- Temple **well** prop/hotspot art (or reuse a scene element).
- Kenji **talk** sheet (still pending); Hiro talk landed.
- (Carry-over Japan art still open: geisha talk, Lily sad talk, gap re-rolls,
  sakura-grove BG, BG edge-continuity.)

## Code touchpoints (when we build)

- `game/tokyo.go` — Hiro/Kenji/Oba-chan trade wiring (held-item gating like the
  Paris/Jerusalem chains), the store open/close prop + queue-sit swap, the well
  hotspot, the grove "place offering → bloom → pick" rework.
- `game/game_state.go` — flags: `jp_ramen_open`, plus per-step gates as needed.
- Item registry + icons for the 4 new items.
- Reuse: `firstExisting` candidate loading, `newAmbientSway`, the handOff trade
  pattern, the `grab_flower` one-shot for the grove pick.

## Matcha ceremony — required sub-quest (BUILT 2026-06-24)

A short ceremony that GATES the grove (the last BG): the grove exit now needs
both Oba-chan's opened path (`jp_grove_revealed`) AND a still heart
(`jp_tea_done`). Not story-critical, but you must do it to enter. Hosted by a
NEW **Tea Master** NPC in the flower store. No reward item — "just the moment."

```
Flower store: talk to KIKU the geisha - she dresses PP (kimono-spin gag) AND
  TEACHES the way of tea. This sets jp_tea_learned and is what UNLOCKS the
  matcha + bowl shelves (you can't start the ceremony until you've heard her).
    |
Flower store: grab MATCHA (tea-shelf hotspot) + a TEA BOWL (bowl-shelf hotspot,
  a RANDOM chawan via math/rand - cosmetic flavor only). Both gated on jp_tea_learned.
    |
Street WELL: with Matcha + Tea Bowl in hand, whisk them with cool water -> MATCHA BOWL.
    |
Flower store TEA MASTER: bring the Matcha Bowl -> kneel + whisk + sip one-shot
  (PP_tea_ceremony) -> jp_tea_done = true -> the grove will open.
```

Items: Matcha, Tea Bowl, Matcha Bowl (all consumed by the end). Art queued at
EXTRA_PROMPTS §JP-MATCHA (tea master idle/talk, PP sit-and-drink, 3 item icons).

## Resolved (2026-06-24)

1. Kiku the geisha = **teaches the tea ceremony** (and does the kimono gag);
   talking to her sets `jp_tea_learned` and unlocks the matcha + bowl shelves -
   so the tea sub-quest can't begin until PP has heard her. (Superseded the
   earlier "optional gag" call.)
2. The well = **hotspot in the ramen street** (local to Kenji + the matcha whisk).
3. Hiro's stall **stays open, line seated** after serving PP.
4. Matcha ceremony = **required gate** for the grove (`jp_tea_done`); host = a
   **new Tea Master**; reward = **just the moment** (no item).

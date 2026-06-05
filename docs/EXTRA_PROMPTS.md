# Extra Sprite Prompts — everything still needed for the current FIXME sweep

This file is read by ChatGPT inside Cursor. Each paste-ready prompt is wrapped
between `===PROMPT START===` and `===PROMPT END===` markers. **Workflow:**

1. Highlight everything BETWEEN the markers (the blockquote block itself,
   not the marker lines).
2. Paste into ChatGPT. Include the header below (style lock + standing
   rules) as context if ChatGPT doesn't already have it — those rules
   apply to every prompt in the file.
3. Save the resulting PNG at the path shown in that prompt's section.
4. Run the atlas re-pack (or restart the game for legacy loaders):

```
python tools/pack_atlas.py tools/characters/<name>.yaml
```

5. Move the section header into the **Done log** at the bottom of this
   file and delete the prompt body so the working part stays scannable.

---

**Style lock + standing rules below — feed these to ChatGPT once per session
so it doesn't violate them on the next prompt:**

> Hand-drawn 1990s Saturday-morning cartoon, Pink Panther *Hokus Pokus Pink*
> (1997) / *Passport to Peril* (1996). Confident black ink linework ~3 px,
> flat saturated fills, no cross-hatching, no gradients, no airbrush. Two
> cel tones max per color region. Pure #FFFFFF background, zero scenery.
> Every cell is **tall rectangular**, never square.

Canvas dimensions are locked per sheet; do **not** scale down to square.

**Standing PP design rules (apply to EVERY PP prompt):**

1. **No gloves of any color.** Pink Panther in this game has plain
   pink paws/hands — never yellow gloves, never any gloves.
2. **Every pickup sprite ends with PP pocketing the item.** The final
   1-2 frames show the item vanishing into his invisible hip pocket
   (the classic Pink Panther "magic pocket"); PP ends empty-handed in
   a relaxed standing pose with a small secretive smile.
3. **No pure white anywhere on the panther.** Belly uses ivory
   off-white `#F2EFE5`, eye sclera uses pale grey. Pure white pixels
   on PP get chroma-keyed away by the engine.

**Standing rule for ALL characters who need "white" in their design:**
the engine chroma-keys pure `#FFFFFF` plus a tolerance band. Use these in
order of preference for fabric / large white areas:
- **Cream `#E5DDC8`** ← USE THIS for "white shirts" or any large fabric.
- **Bone `#EDE5D3`** — paper, small label panels.
- **Pale grey `#C4C4C4`** — steam wisps, eye sclera.
- **Vanilla `#F2EFE5`** — only safe for tiny accents (a tooth, a button).

The **scene background** in the sprite cell still uses pure
`#FFFFFF` — that IS the chroma key; it's only the character /
foreground objects that must avoid pure white.

---

## Critical separator rule (applies to EVERY multi-frame sheet)

User 2026-05-24: several recent regens (café patrons, Marcus, Higgins office)
came back with **visible thin lines BETWEEN frame cells** — faint grey or
near-white seams that survive the chroma-key and render in-game as dark
verticals between animation frames.

**Fix language to include in every prompt that uses a grid:**

> The sheet is a **flat grid of cells with NO visible separators**: no
> drawn borders, no thin lines, no grey/black strips, no shadow gradients
> between cells. Cell boundaries are conceptual only — neighbouring cells
> meet directly with pure `#FFFFFF` background pixels on both sides. The
> exported PNG must look like ONE continuous canvas where each Nth × Mth
> rectangle happens to hold one frame; if you cropped any cell out you'd
> see only that frame on pure white, with no edge artefacts.

If you see a faint vertical/horizontal line in the preview, the generator
drew a separator — regenerate with the rule above emphasised.

## No extras rule (applies to EVERY sheet)

User 2026-06-02: generators keep adding a large "hero" character **portrait**
beside the frame grid (and sometimes title text / labels). Include in every
prompt:

> Output ONLY the N×M grid of animation frames — nothing else. NO separate
> large character portrait or "hero" reference image beside or above the grid,
> NO title text, NO labels, NO watermark, NO colour swatches. Just the frames
> on pure #FFFFFF.

---

## Open Prompts

All prompts below still need a PNG generated. When one lands, move its row
into the **Done log** at the bottom and delete the body.

No open prompts.

> **Also flagged for in-game audit before regen (don't generate yet):** the
> bakery café-patron sheets under `.../paris/npc/coffee/` (#24 — check each for
> a baked background box / separator lines / wrong frame count; Henri verified
> clean at 8×1), Danny idle/talk halo (#7 — now re-keyed wider in code; regen
> only if the fringe persists), and `pp_sleeping.png`/`pp_waking.png` (#13 —
> regen only if they don't match the current PP design).


## Done log — landed sprites (FYI only, no action needed)

These prompts produced PNGs that are now on disk and wired up. Listed for
record so we don't re-generate them by accident. If you need a variant, the
original prompt is in git history at `docs/EXTRA_PROMPTS.md` pre-2026-05-24.

| § | Sprite | Path | Landed |
|---|--------|------|--------|
| §1 | Higgins entrance idle | `npc_director_higgins_idle.png` | 2026-04 |
| §2 | Higgins walk back | `npc_director_higgins_walk_back.png` | 2026-04 |
| §4 | Marcus strange_alt | `npc_marcus_strange_alt.png` | 2026-04 |
| §6 | Campfire small loop | campfire frames | 2026-04 |
| §8 | Bakery Woman | `npc_bakery_woman.png` | 2026-04 |
| §9 | Press Photographer | `npc_press_photographer.png` | 2026-04 |
| §10 | Higgins entrance talk | `npc_director_higgins_talk.png` | 2026-04 |
| §18 | Higgins office idle + talk (v1) | `npc_director_higgins_office_*.png` | 2026-04 — superseded by §C above |
| §19 | Higgins give_map handoff | `npc_director_higgins_give_map.png` | 2026-04 |
| §Y | Paris Bakery BG v2 (door right + tablecloths + framed counter) | `paris_bakery.png` | 2026-05-23 |
| §E | Tommy walk_left | `npc_tommy_walk_left.png` | 2026-05-21 |
| §F | Jake walk_back | `npc_jake_walk_back.png` | 2026-05-21 |
| §M | Action cursor (cursor_point) | `cursor_point.png` | 2026-05-21 |
| §H | PP airplane (modern Cessna-style + pilot) | `pp_airplane.png` | 2026-05-23 |
| §7 | Café patrons combined sheets (v1, fringe issues) | `cafe_patron_<name>.png` | 2026-05 — superseded by §D above |
| §NEW Paris Clouds | Paris Clouds airplane sky | `paris_clouds.png` | 2026-05-23 |
| §I | Higgins throw-map one-shot | `npc_director_higgins_throw_map.png` | 2026-05-23 |
| §J | PP catch-map one-shot | `pp_catch_map.png` | 2026-05-23 |
| §K | Thrown-map projectile sprite | `inv_travel_map_throw.png` | 2026-05-23 |
| §L | Travel-map inventory icon | `travel_map_icon.png` | 2026-05 |
| §N | Item sprites (8 items) | `assets/images/items/*.png` | 2026-05-23 |
| §R | Café au Lait inventory item | `cafe_au_lait.png` | 2026-05-23 |
| §S | Confiture inventory item | `confiture.png` | 2026-05-23 |
| §T | Camille quick-sketch one-shot | `npc_camille_sketching.png` | 2026-05-23 |
| §V | Henri give-jam one-shot | `npc_henri_give_jam.png` | 2026-05-23 |
| §E | Danny talk clean regen | `npc_danny_talk.png` | 2026-05-31 |
| §F | Madame Poulain work alt-idle | `npc_madame_poulain_work.png` | 2026-05-31 |
| §G | Marcus strange idle + talk clean 8×2 | `npc_marcus_strange_*.png` | 2026-05-31 |
| §E2 | Danny talk idle-match regen | `npc_danny_talk.png` | 2026-06-02 |
| §G2 | Marcus strange talk idle-match regen | `npc_marcus_strange_talk.png` | 2026-06-02 |
| §G3 | Marcus strange talk freakout regen | `npc_marcus_strange_talk.png` | 2026-06-02 |
| §M1 | Marcus strange idle design-match regen | `npc_marcus_strange_idle.png` | 2026-06-02 |
| §M2 | Marcus strange talk design-match regen | `npc_marcus_strange_talk.png` | 2026-06-02 |
| §M3 | Marcus strange alt idle beat regen | `npc_marcus_strange_alt.png` | 2026-06-02 |
| §P1 | PP walk front current-design redraw | `PP walk front.png` | 2026-06-02 |
| §C1 | Madame Colette idle beret/stripes design-lock regen | `npc_madame_colette_idle.png` | 2026-06-02 |
| §C2 | Madame Colette talk beret/stripes design-lock regen | `npc_madame_colette_talk.png` | 2026-06-02 |
| §C3 | Director Higgins green-shirt shout design-lock regen | `npc_director_higgins_shout.png` | 2026-06-02 |
| Madame Colette | **DO NOT REGENERATE** — user 2026-05-23 prefers the current design | `npc_french_guide_*.png` | — |

**Removed in 2026-05-24 cleanup (low-priority / deferred):** previous PP
talk-front + talk-side regen prompts (user reversed direction — see §A),
previous Marcus talk regen (user wants idle recolor first — see §B),
previous Higgins office regen prompt (replaced by §C with the "match
entrance design" instruction), PP grab-flower regen, PP grab rolling pin,
Marcus strange-idle fringe touch-up, Windows .exe icon prompt. The git
history of this file before 2026-05-24 has the bodies if any of these
come back.

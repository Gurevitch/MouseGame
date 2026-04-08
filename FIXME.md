# Pink Panther: Camp Chilly Wa Wa — Known Issues & Fixes

> **Reference:** See [STORY.md](STORY.md) for full story flow and design.
> **Progress:** See [STATUS.md](STATUS.md) for implementation status.

---

## How to Use

Add issues below as you find them. Mark priority: `[P0]` critical, `[P1]` important, `[P2]` minor.
When fixed, move to the **Resolved** section with the date.

---

## Open Issues

### Story / Flow

- [ ] `[P1]` Postcard not added to inventory after Curator dialog (item pickup missing)
- [ ] `[P1]` No Marcus healing flow — giving postcard to Marcus should cure him
- [ ] `[P2]` Higgins bedtime dialog uses existing entrance sprite — needs new camp sprite (see STATUS.md assets)
- [ ] `[P2]` Airplane cutscene uses PP standing idle — needs sitting-in-plane sprite

### Scenes / Navigation

- [ ] `[P2]` Kid rooms (Tommy, Jake, Lily, Danny) have no NPCs inside — should have kids in their rooms on Day 2

### Travel Map

- [ ] `[P2]` Buenos Aires and Rio pins are close together — landmarks may overlap visually

### Assets

- [ ] `[P1]` Need Higgins camp sprite (idle + talk) for bedtime scene
- [ ] `[P1]` Need PP airplane idle sprite for flight cutscene
- [ ] `[P2]` No airplane background — currently using paris_clouds.png as fallback

---

## Resolved

| Issue | Date | Notes |
|-------|------|-------|
| Map used dots instead of landmarks | 2026-04-08 | Replaced with landmark images from ui/landmarks/ |
| Night scene too simple (instant Marcus room) | 2026-04-08 | Reworked to multi-phase: campfire sleep → Marcus freakout → wakeup |
| No airplane transition before cities | 2026-04-08 | Added airplane_flight scene with 4-second cutscene |

---

*When adding issues, check STORY.md to verify expected behavior.*

### Reported (2026-04-08)

- [x] `[P1]` PP sleeping sprite showed in Marcus room + freakout dialog didn't auto-start — FIXED: sleeping sprite hidden during Marcus phase, dialog trigger fixed with flag
- [x] `[P2]` No Higgins visible on camp_grounds during bedtime — FIXED: Higgins NPC now spawns on camp_grounds for bedtime dialog, removed after
- [ ] `[P3]` Camp night background looks too big/zoomed — this is an asset issue (the camp_night.png image itself). Regenerate the image at correct proportions or crop it
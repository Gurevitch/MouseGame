# Pink Panther: Camp Chilly Wa Wa — Status & Leaderboard

> **Reference:** See [STORY.md](STORY.md) for full story flow and design.
> **Issues:** See [FIXME.md](FIXME.md) for known bugs and issues.

---

## Implementation Leaderboard

### Engine & Systems

| System | Status | Notes |
|--------|--------|-------|
| SDL2 rendering | DONE | 1400x750 window |
| Bitmap font | DONE | Custom pixel font with scaling |
| Scene manager | DONE | Fade transitions between scenes |
| Player movement | DONE | Walk, idle, talk, grab, examine, celebrate |
| NPC system | DONE | Idle/talk animations, dialog callbacks, strange variants |
| Dialog system | DONE | Typing animation, speaker sync, callbacks |
| Inventory | DONE | Oval carousel, drag-drop, held items |
| Cursor system | DONE | Normal, talk, grab, arrow cursors |
| Particle system | DONE | Fire, smoke, birds, butterflies, insects, clouds, water, dust, stars |
| Glow effects | DONE | Pulsing ambient lighting |
| Audio manager | DONE | Background music |
| Travel map | DONE | Globe with landmark icons, flight paths |
| Airplane cutscene | DONE | 4-second flight scene with clouds |

### Story Progress

| Story Beat | Status | Scene(s) |
|------------|--------|----------|
| Day 1: Arrival monologue | DONE | camp_entrance |
| Day 1: Meet Higgins | DONE | camp_entrance |
| Day 1: Meet all 5 kids | DONE | camp_grounds |
| Night: Higgins bedtime | DONE | camp_grounds |
| Night: Campfire sleep | DONE | camp_night |
| Night: Marcus freakout | DONE | marcus_room (night) |
| Night: Morning wakeup | DONE | camp_night |
| Day 2: Morning monologue | DONE | camp_grounds |
| Day 2: Talk to Marcus (strange) | DONE | marcus_room (day) |
| Day 2: Higgins gives map | DONE | camp_office |
| Airplane flight to Paris | DONE | airplane_flight |
| Paris: Street monologue | DONE | paris_street |
| Paris: Talk to Madame Colette | DONE | paris_street |
| Paris: Talk to Curator Beaumont | DONE | paris_louvre |
| Paris: Get postcard | DONE | paris_louvre |
| Give postcard to Marcus | TODO | marcus_room |
| Marcus returns to normal | TODO | marcus_room |
| Next kid goes strange | TODO | camp_grounds |

### Scenes

| Scene | Status | NPCs | Hotspots |
|-------|--------|------|----------|
| camp_entrance | DONE | Higgins | Enter Camp, Air strip |
| camp_grounds | DONE | Marcus, Tommy, Jake, Lily, Danny | Lake, Office, 5 cabins |
| camp_office | DONE | Office Higgins | Back to Camp |
| camp_night | DONE | (cutscene) | — |
| camp_lake | DONE | — | Back to Camp |
| marcus_room | DONE | Marcus | Exit Cabin |
| tommy_room | DONE | — | Exit Cabin |
| jake_room | DONE | — | Exit Cabin |
| lily_room | DONE | — | Exit Cabin |
| danny_room | DONE | — | Exit Cabin |
| airplane_flight | DONE | — | (auto-transition) |
| paris_street | DONE | Madame Colette | To Louvre, Travel Map |
| paris_louvre | DONE | Curator Beaumont | Back to Street, Travel Map |

### Cities

| City | Landmark | Map Status | Gameplay |
|------|----------|------------|----------|
| Camp Chilly Wa Wa | Red pin (no landmark) | Always unlocked | DONE |
| Paris | Eiffel Tower | Unlocks after Marcus dialog | DONE |
| Jerusalem | Western Wall | Locked | TODO |
| Tokyo | Torii Gate | Locked | TODO |
| Rome | Colosseum | Locked | TODO |
| Rio de Janeiro | Christ Redeemer | Locked | TODO |
| Buenos Aires | Statue of Liberty | Locked | TODO |
| Mexico City | Pyramids | Locked | TODO |

### NPCs

| NPC | Normal Dialog | Post Dialog | Strange Dialog | Post-Strange | Alt Item Dialog |
|-----|-------------|-------------|----------------|--------------|-----------------|
| Director Higgins (entrance) | DONE | DONE | — | — | — |
| Director Higgins (office) | DONE | DONE | — | — | — |
| Marcus | DONE | DONE | DONE | DONE | TODO (postcard) |
| Tommy | DONE | DONE | DONE | DONE | TODO |
| Jake | DONE | DONE | DONE | DONE | TODO |
| Lily | DONE | DONE | DONE | DONE | TODO |
| Danny | DONE | DONE | DONE | DONE | TODO |
| Madame Colette | DONE | DONE | — | — | — |
| Curator Beaumont | DONE | DONE | — | — | — |

### Assets Needed

| Asset | Type | Status |
|-------|------|--------|
| Higgins camp sprite (idle) | NPC sprite | NEEDED |
| Higgins camp sprite (talk) | NPC sprite | NEEDED |
| PP airplane idle sprite | Player sprite | NEEDED |
| Jerusalem scenes | Backgrounds | TODO |
| Tokyo scenes | Backgrounds | TODO |
| Rome scenes | Backgrounds | TODO |
| Rio scenes | Backgrounds | TODO |
| Buenos Aires scenes | Backgrounds | TODO |
| Mexico City scenes | Backgrounds | TODO |

---

*Last updated: 2026-04-08*

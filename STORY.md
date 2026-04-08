# Pink Panther: Camp Chilly Wa Wa — Story Flow

## Overview
Pink Panther arrives at Camp Chilly Wa Wa as a substitute counselor. The kids start experiencing strange visions of real-world cities they've never visited. PP must travel to each city, find the source of their visions, and bring back an anchor object to help each child.

---

## DAY 1 — Arrival (Normal)

```
Camp Entrance
    |
    v
[Opening monologue plays]
    |
    v
[Auto-walk to Director Higgins]
    |
    v
Higgins: "Welcome! Go meet the kids."
    |
    v (arrow UP — walk up the road)
    |
Camp Grounds -----> Higgins' Office (right arrow)
    |
    v
Meet all 5 kids (normal personalities):
  - Marcus: Know-it-all, loves sketching
  - Tommy: Storyteller, music lover  
  - Jake: Tough kid, coin collector
  - Lily: Shy, loves flowers
  - Danny: Prankster, treasure stories
    |
    v (after talking to all 5)
    |
[Higgins arrives: "It's getting late!"]
    |
    v
```

## NIGHT — The Campfire

```
[Transition to camp_night scene]
    |
    v
[PP falls asleep on log by fire — sleeping animation]
    |
    v
[Marcus appears, frantic]
Marcus: "I can't stop... golden frames... a glass pyramid..."
    |
    v
[PP wakes up — waking animation]
PP: "I need to find Higgins."
    |
    v
```

## DAY 2 — The Weirdness

```
[Transition to Higgins' Office]
    |
    v
[Day 2 monologue: "Something feels different..."]
    |
    v
Higgins: worried, gives Travel Map
    |
    v (go left to camp_grounds)
    |
Camp Grounds (Marcus has strange sprite)
    |
    v
Talk to Marcus (strange dialog):
  "A woman's face... ornate golden frames... a glass pyramid..."
  "The biggest museum in the world!"
    |
    v
[Paris unlocks on travel map]
    |
    v
Use Travel Map (from inventory, click anywhere) --> Globe opens
    |
    v
Click Paris --> Fly to Paris
```

## PARIS — The Louvre

```
Paris Street (Eiffel Tower visible)
    |
    v
Talk to Madame Colette (French cafe owner):
  - Eiffel Tower: Built 1889, was temporary, Gustave Eiffel
  - Louvre: Largest museum, 380,000 objects
  - Glass Pyramid: I.M. Pei, 1989
    |
    v (arrow RIGHT)
    |
Louvre Interior
    |
    v
Talk to Curator Beaumont:
  - Mona Lisa: Leonardo da Vinci, 1503
  - Venus de Milo: Greek, 100 BC
  - Identifies Marcus's painting
  - Gives postcard of the painting (anchor object)
    |
    v
[Use Travel Map --> Return to Camp]
    |
    v
```

- [ ] Give postcard to Marcus --> Marcus calms down
- [ ] Marcus returns to normal sprite
- [ ] Next kid starts showing strange behavior

---

## FUTURE CITIES (Not Yet Implemented)

### Jerusalem (Jake's visions — tunnels, coins, echoes)
- [ ] Jake goes strange: rubbing surfaces, hearing echoes, "tunnels under old city"
- [ ] Travel to Jerusalem
- [ ] Visit Western Wall area
- [ ] Explore ancient tunnels
- [ ] Learn about Old City, Western Wall history
- [ ] Find anchor object: coin rubbing or tunnel map
- [ ] Return to camp, give to Jake

### Tokyo (Lily's visions — glowing garden, temple bells)
- [ ] Lily goes strange: arranging petals in symbols, hearing bells
- [ ] Travel to Tokyo
- [ ] Visit Senso-ji temple
- [ ] Explore cherry blossom garden
- [ ] Learn about temples, cherry blossoms, torii gates
- [ ] Find anchor object: pressed flower or charm card
- [ ] Return to camp, give to Lily

### Rome (Danny's visions — arena, gold paths, ruins)
- [ ] Danny goes strange: mapping camp like ruins, digging holes
- [ ] Travel to Rome
- [ ] Visit Colosseum
- [ ] Explore Roman ruins
- [ ] Learn about gladiators, Colosseum (72 AD), Roman engineering
- [ ] Find anchor object: seal sketch or inscription rubbing
- [ ] Return to camp, give to Danny

### Rio de Janeiro + Buenos Aires (Tommy's visions — music, statue, carnival, tango)
- [ ] Tommy goes strange: hearing music, dancing, mixing landmarks
- [ ] Travel to Rio de Janeiro
  - [ ] Christ the Redeemer statue
  - [ ] Copacabana beach
  - [ ] Carnival
- [ ] Travel to Buenos Aires
  - [ ] La Boca neighborhood
  - [ ] Obelisco
  - [ ] Tango
- [ ] Find anchor object: dance card or paired stamp
- [ ] Return to camp, give to Tommy

### Mexico City (Secret city — discovered during gameplay)
- [ ] Unlocked by clue found in another city (Rome?)
- [ ] Aztec pyramids, Teotihuacan
- [ ] Frida Kahlo museum
- [ ] Learn about Aztec civilization, pyramids

---

## Kid Repair Formula

Each kid follows the same pattern:
1. **Normal** — Kid has a natural personality trait (Day 1)
2. **Strange** — That trait becomes exaggerated/supernatural (Day 2+)
3. **Travel** — PP flies to the real city the kid is seeing
4. **Discover** — PP finds the source of the kid's visions, learns real history
5. **Anchor** — PP brings back a small real-world object from the place
6. **Repair** — Kid sees the object, their mind settles, they return to normal

---

## Travel System

- Camp Chilly Wa Wa Air (funny old propeller plane)
- Globe/world map with city pins
- Unlocked cities glow yellow
- Locked cities show as gray dots with name
- Dotted flight path between camp and destination

---

## Cities on Globe

| City | Kid | Status |
|------|-----|--------|
| Camp Chilly Wa Wa | — | Always unlocked |
| Paris | Marcus | Unlocked after Day 2 Marcus dialog |
| Jerusalem | Jake | Not yet implemented |
| Tokyo | Lily | Not yet implemented |
| Rome | Danny | Not yet implemented |
| Rio de Janeiro | Tommy | Not yet implemented |
| Buenos Aires | Tommy | Not yet implemented |
| Mexico City | Secret | Not yet implemented |

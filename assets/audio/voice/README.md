# Voice clips

Per-dialog-line voice clips. Wired through `dialogEntry.audio`. When a
dialog line starts, the engine plays the file at that path via
`audioManager.playSFX`. Missing files silently no-op so the dialog still
advances on text only.

## Naming convention

`<character>_<scene_or_chapter>_<line_id>.mp3`

Examples:
- `higgins_camp_intro_01.mp3` — Higgins's first welcome line at camp_entrance
- `lily_flower_receive_01.mp3` — Lily's "thank you" when given the flower
- `pp_paris_arrival.mp3` — PP's monologue on arrival in Paris

Keep line IDs lower-case-with-underscores and stable; they appear in
`dialogEntry.audio` fields throughout the Go source.

## Format

MP3 mono is fine for voice (mono saves bandwidth and matches retro mono
voiceover). 22.05 kHz / 44.1 kHz both work. Keep files short — one line
per file.

## Adding a voice line

1. Drop the file at `assets/audio/voice/<character>_<id>.mp3`.
2. In the `dialogEntry` literal (or JSON dialog file), set
   `audio: "assets/audio/voice/<character>_<id>.mp3"`.
3. No code change — `dialogSystem.advance` calls `playSFX` on each
   line start.

# Music files

Looping background tracks, one per region or scene. Engine plays the
file at `scene.musicPath` (set in each scene's JSON) via
`audioManager.playMusic`. Missing files are silently no-op'd, so it's
safe to leave entries pointing at not-yet-authored paths.

## Expected files

| Path | Used by |
|---|---|
| `camp.mp3` | `camp_entrance`, `camp_grounds`, `camp_lake` |
| `camp_night.mp3` | `camp_night` |
| `camp_office.mp3` | `camp_office` |
| `camp_cabin.mp3` | `tommy_room`, `jake_room`, `lily_room`, `marcus_room`, `danny_room` |
| `paris.mp3` | `paris_street`, `paris_bakery`, `paris_louvre` |
| `airplane.mp3` | `airplane_flight` |

## Format

MP3 looped via MCI (`mciSendString play ... repeat`). 44.1 kHz stereo
recommended. Keep tracks short (1–2 min) so the loop point isn't
jarring; the MCI loop is seamless.

## Adding a new region

1. Drop the file at `assets/audio/music/<region>.mp3`.
2. Set `"musicPath": "assets/audio/music/<region>.mp3"` in the
   scene's JSON.
3. No code change needed — `Game.update` plays the path on scene change.

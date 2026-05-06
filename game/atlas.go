package game

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/Local/games/PP/engine"

	"github.com/veandco/go-sdl2/sdl"
)

// atlasCache holds one AtlasSheet per character so multiple NPC instances
// that reference the same character (e.g. Marcus in camp_grounds *and* in
// marcus_room) share a single SDL texture instead of loading it twice.
//
// Lives package-global because the renderer and SDL textures live for the
// game's lifetime — there's no teardown path today. If we ever add one,
// walk this map and call .Destroy() on each AtlasSheet.Texture.
var atlasCache = map[string]*AtlasSheet{}

// GetAtlas returns the cached atlas for `name`, loading it if this is the
// first request. Returns nil (and a one-time log) when the on-disk files
// are missing — callers can fall through to legacy loaders.
func GetAtlas(renderer *sdl.Renderer, name string) *AtlasSheet {
	if s, ok := atlasCache[name]; ok {
		return s
	}
	s := LoadAtlas(renderer, name)
	atlasCache[name] = s // cache nil too so a missing atlas only logs once
	return s
}

// AtlasSheet is a character's packed sprite atlas plus its animation metadata.
// One texture per sheet, N named animations, each with a fixed frame size and
// per-frame source rectangles into the atlas.
//
// Source of truth for the data here is tools/pack_atlas.py, which generates
// one <name>.png + <name>.json pair per character under assets/sprites/. The
// PNG already has color-key transparency applied in the packer, so the Go
// loader does no runtime cleanup.
type AtlasSheet struct {
	Texture    *sdl.Texture
	Animations map[string]*AtlasAnimation
}

// AtlasAnimation is a named animation strip inside an atlas: frame rectangles
// plus playback speed.
type AtlasAnimation struct {
	Name    string
	FPS     float64
	FrameW  int32
	FrameH  int32
	Frames  []sdl.Rect
	sheet   *AtlasSheet
}

// atlasJSON matches the on-disk schema emitted by pack_atlas.py.
type atlasJSON struct {
	Image      string                     `json:"image"`
	Animations map[string]atlasAnimJSON   `json:"animations"`
}

type atlasAnimJSON struct {
	FPS     float64         `json:"fps"`
	FrameW  int32           `json:"frame_w"`
	FrameH  int32           `json:"frame_h"`
	Frames  []atlasFrameJSON `json:"frames"`
}

type atlasFrameJSON struct {
	X, Y, W, H int32
}

// LoadAtlas loads an atlas by base name (no extension), looking up
//   assets/sprites/<name>.png and assets/sprites/<name>.json
// Subfolders inside assets/sprites/ are fine too: LoadAtlas(r, "paris/pierre")
// resolves to assets/sprites/paris/pierre.png + .json.
//
// Returns nil with a logged warning if either file is missing — callers that
// want to fall back to the legacy per-state loaders can nil-check.
func LoadAtlas(renderer *sdl.Renderer, name string) *AtlasSheet {
	base := filepath.Join("assets", "sprites", filepath.FromSlash(name))
	jsonPath := base + ".json"
	pngPath := base + ".png"

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Printf("atlas: missing %s: %v\n", jsonPath, err)
		return nil
	}
	var meta atlasJSON
	if err := json.Unmarshal(data, &meta); err != nil {
		fmt.Printf("atlas: bad JSON %s: %v\n", jsonPath, err)
		return nil
	}

	tex, _, _ := engine.SafeTextureFromPNGRaw(renderer, pngPath)
	if tex == nil {
		return nil
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)

	sheet := &AtlasSheet{
		Texture:    tex,
		Animations: make(map[string]*AtlasAnimation, len(meta.Animations)),
	}
	for animName, a := range meta.Animations {
		frames := make([]sdl.Rect, len(a.Frames))
		for i, f := range a.Frames {
			frames[i] = sdl.Rect{X: f.X, Y: f.Y, W: f.W, H: f.H}
		}
		sheet.Animations[animName] = &AtlasAnimation{
			Name:   animName,
			FPS:    a.FPS,
			FrameW: a.FrameW,
			FrameH: a.FrameH,
			Frames: frames,
			sheet:  sheet,
		}
	}
	return sheet
}

// Animation is a playback cursor into an AtlasAnimation. Multiple cursors can
// share one AtlasAnimation (e.g. two NPCs using the same sheet) without
// interfering — the cursor state lives here, not on the animation.
type Animation struct {
	anim  *AtlasAnimation
	idx   int
	timer float64
	loop  bool
}

// Play returns a cursor initialized to frame 0 of the named animation, looping
// by default. Unknown names fall back to the first animation in the sheet so a
// typo never produces a nil-deref at a bad moment.
func (s *AtlasSheet) Play(name string) *Animation {
	a, ok := s.Animations[name]
	if !ok {
		for _, any := range s.Animations {
			a = any
			break
		}
		if a != nil {
			fmt.Printf("atlas: animation %q not found, falling back to %q\n", name, a.Name)
		}
	}
	if a == nil {
		return nil
	}
	return &Animation{anim: a, loop: true}
}

// Has reports whether the atlas knows the named animation.
func (s *AtlasSheet) Has(name string) bool {
	_, ok := s.Animations[name]
	return ok
}

// Texture returns the underlying SDL texture for blitting.
func (s *AtlasSheet) Tex() *sdl.Texture { return s.Texture }

// Update advances playback by dt seconds. Looping wraps around; non-looping
// clamps to the last frame.
func (a *Animation) Update(dt float64) {
	if a == nil || a.anim == nil || len(a.anim.Frames) == 0 {
		return
	}
	if a.anim.FPS <= 0 {
		return
	}
	a.timer += dt
	frameDur := 1.0 / a.anim.FPS
	for a.timer >= frameDur {
		a.timer -= frameDur
		a.idx++
		if a.idx >= len(a.anim.Frames) {
			if a.loop {
				a.idx = 0
			} else {
				a.idx = len(a.anim.Frames) - 1
				a.timer = 0
				return
			}
		}
	}
}

// Current returns the source rectangle of the current frame in the atlas.
func (a *Animation) Current() sdl.Rect {
	if a == nil || a.anim == nil || len(a.anim.Frames) == 0 {
		return sdl.Rect{}
	}
	return a.anim.Frames[a.idx]
}

// FrameSize returns the (w, h) of this animation's frames.
func (a *Animation) FrameSize() (int32, int32) {
	if a == nil || a.anim == nil {
		return 0, 0
	}
	return a.anim.FrameW, a.anim.FrameH
}

// Reset rewinds to frame 0.
func (a *Animation) Reset() {
	if a == nil {
		return
	}
	a.idx = 0
	a.timer = 0
}

// SetLoop toggles looping behavior.
func (a *Animation) SetLoop(loop bool) {
	if a != nil {
		a.loop = loop
	}
}

// Done reports whether a non-looping animation has reached its last frame.
func (a *Animation) Done() bool {
	if a == nil || a.anim == nil {
		return true
	}
	return !a.loop && a.idx >= len(a.anim.Frames)-1
}

// Tex returns the atlas texture that backs this animation's frames.
func (a *Animation) Tex() *sdl.Texture {
	if a == nil || a.anim == nil || a.anim.sheet == nil {
		return nil
	}
	return a.anim.sheet.Texture
}

// --- helpers for the smoke test / future inspectors ---

// applyKidAtlas wires the four canonical kid animations (idle / talk /
// strange_idle / strange_talk) on an *npc from the named atlas. Returns true
// on success; false if the atlas is missing or an expected animation isn't in
// it, in which case the caller should fall back to legacy per-PNG loading.
//
// This is the one-call replacement for loadNPCGridKids + loadStrangeGridsKids:
//   applyKidAtlas(renderer, n, "tommy")
// covers what used to be 4 separate PNG loads.
func applyKidAtlas(renderer *sdl.Renderer, n *npc, atlasName string) bool {
	sheet := GetAtlas(renderer, atlasName)
	if sheet == nil {
		return false
	}
	idle := sheet.npcFrames("idle")
	talk := sheet.npcFrames("talk")
	si := sheet.npcFrames("strange_idle")
	st := sheet.npcFrames("strange_talk")
	if len(idle) == 0 || len(talk) == 0 {
		return false
	}
	n.idleGrid = idle
	n.talkGrid = talk
	n.strangeIdle = si
	n.strangeTalk = st
	return true
}

// applyNPCAtlas wires just the two canonical adult-NPC animations (idle /
// talk) from the named atlas. Returns true on success; false if the atlas
// is missing or either idle/talk is absent, letting the caller fall through
// to legacy loadNPCGrid*-based PNG loads.
//
// Use this for Paris NPCs and other adult flavor NPCs that don't have the
// strange_* freakout pair. For the kid cast, use applyKidAtlas instead so
// the freakout rows get wired.
//
// Atlas names for Paris live under the "paris/" subfolder:
//   applyNPCAtlas(renderer, n, "paris/bakery_woman")
// resolves to assets/sprites/paris/bakery_woman.png + .json.
func applyNPCAtlas(renderer *sdl.Renderer, n *npc, atlasName string) bool {
	sheet := GetAtlas(renderer, atlasName)
	if sheet == nil {
		return false
	}
	idle := sheet.npcFrames("idle")
	talk := sheet.npcFrames("talk")
	if len(idle) == 0 || len(talk) == 0 {
		return false
	}
	n.idleGrid = idle
	n.talkGrid = talk
	return true
}

// npcFrames returns the animation as a []npcFrame suitable for assignment to
// npc.idleGrid / npc.talkGrid. All frames share the atlas texture; each gets
// its own source rect so the NPC renderer samples the correct cell.
//
// Returns nil if the animation is missing from the atlas — callers that want
// a legacy fallback can nil-check and fall through to loadNPCGrid.
func (s *AtlasSheet) npcFrames(animName string) []npcFrame {
	if s == nil {
		return nil
	}
	a, ok := s.Animations[animName]
	if !ok {
		return nil
	}
	out := make([]npcFrame, len(a.Frames))
	for i, f := range a.Frames {
		r := f
		out[i] = npcFrame{tex: s.Texture, w: f.W, h: f.H, src: &r}
	}
	return out
}

// AnimationNames returns the sorted list of animation names in the sheet.
// Handy for debug overlays and the atlas smoke test.
func (s *AtlasSheet) AnimationNames() []string {
	names := make([]string, 0, len(s.Animations))
	for k := range s.Animations {
		names = append(names, k)
	}
	// simple insertion sort; expected N is small (<10)
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && strings.Compare(names[j-1], names[j]) > 0; j-- {
			names[j-1], names[j] = names[j], names[j-1]
		}
	}
	return names
}

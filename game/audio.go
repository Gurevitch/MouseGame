package game

// Audio support requires SDL2_mixer.
// To enable: pacman -S mingw-w64-x86_64-SDL2_mixer
// Then use github.com/veandco/go-sdl2/mix

type audioManager struct {
	enabled bool
}

func newAudioManager() *audioManager {
	return &audioManager{enabled: false}
}

func (am *audioManager) playMusic(path string) {}
func (am *audioManager) fadeOutMusic(ms int)   {}
func (am *audioManager) stopMusic()           {}

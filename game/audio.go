package game

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	winmm         = syscall.NewLazyDLL("winmm.dll")
	mciSendStringW = winmm.NewProc("mciSendStringW")
)

func mciSend(cmd string) error {
	p, _ := syscall.UTF16PtrFromString(cmd)
	ret, _, _ := mciSendStringW.Call(uintptr(unsafe.Pointer(p)), 0, 0, 0)
	if ret != 0 {
		return fmt.Errorf("MCI error %d for: %s", ret, cmd)
	}
	return nil
}

type audioManager struct {
	currentPath string
	playing     bool
}

func newAudioManager() *audioManager {
	return &audioManager{}
}

func (am *audioManager) playMusic(path string) {
	if path == am.currentPath && am.playing {
		return
	}
	am.stopMusic()
	if path == "" {
		return
	}
	if err := mciSend(fmt.Sprintf(`open "%s" type mpegvideo alias bgmusic`, path)); err != nil {
		fmt.Println("Audio open (non-fatal):", err)
		return
	}
	if err := mciSend("play bgmusic repeat"); err != nil {
		fmt.Println("Audio play (non-fatal):", err)
		mciSend("close bgmusic")
		return
	}
	am.currentPath = path
	am.playing = true
}

func (am *audioManager) stopMusic() {
	if am.playing {
		mciSend("stop bgmusic")
		mciSend("close bgmusic")
		am.playing = false
	}
	am.currentPath = ""
}

func (am *audioManager) close() {
	am.stopMusic()
}

// playSFX plays a one-shot sound effect. Used by travel-map info popups to
// optionally play a short voice clip alongside the text popup. Silently
// no-ops if the file is missing (path = "" or file does not exist) so the
// popup still works when audio files haven't been authored yet.
//
// Uses a distinct MCI alias so it doesn't interrupt background music.
func (am *audioManager) playSFX(path string) {
	if path == "" {
		return
	}
	if _, err := os.Stat(path); err != nil {
		// Silent no-op — authors haven't dropped the file yet.
		return
	}
	// Close any previous SFX instance first so rapid clicks don't stack.
	mciSend("close sfx")
	if err := mciSend(fmt.Sprintf(`open "%s" type mpegvideo alias sfx`, path)); err != nil {
		fmt.Println("Audio SFX open (non-fatal):", err)
		return
	}
	if err := mciSend("play sfx"); err != nil {
		fmt.Println("Audio SFX play (non-fatal):", err)
		mciSend("close sfx")
	}
}

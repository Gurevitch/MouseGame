package game

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	winmm          = syscall.NewLazyDLL("winmm.dll")
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

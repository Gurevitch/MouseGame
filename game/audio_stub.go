//go:build !windows

package game

// mciSend reports success without playing anything: MCI is a Windows-only API,
// and returning nil keeps audioManager's state tracking consistent instead of
// spamming "non-fatal" warnings on every track change.
func mciSend(cmd string) error { return nil }

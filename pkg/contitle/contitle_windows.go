//go:build windows && !linux

package contitle

import (
	"syscall"
	"unsafe"
)

// Set Windows console Title text (Header)
//
// base take here https://github.com/lxi1400/GoTitle
func SetTitle(title string) error {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return err
	}

	defer func() {
		_ = syscall.FreeLibrary(handle)
	}()

	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return err
	}

	pStr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}

	_, _, err = syscall.SyscallN(proc, uintptr(unsafe.Pointer(pStr)), 0, 0)
	if err != syscall.Errno(0) {
		return err
	}

	return nil
}

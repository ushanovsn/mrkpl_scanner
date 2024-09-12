//go:build windows && !linux

package contitle

import (
	"syscall"
	"unsafe"
)

// Set Windows console Title text (Header)
//
// base take here https://github.com/lxi1400/GoTitle
func SetTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = syscall.FreeLibrary(handle)
	}()
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	pStr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return 0, err
	}

	r, _, err := syscall.SyscallN(proc, 1, uintptr(unsafe.Pointer(pStr)), 0, 0)
	//r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)

	return int(r), err
}

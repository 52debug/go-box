package osutil

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle   = kernel32.NewProc("GetStdHandle")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

// DisableConsoleQuickEdit 禁用 Windows 控制台的快速编辑模式
// 防止用户选中文字时程序被挂起
func DisableConsoleQuickEdit() error {
	if runtime.GOOS != "windows" {
		return nil
	}

	// 下面都是 Windows 专用的代码
	const (
		StdInputHandle   = 0xFFFFFFF6 // -10 的 64 位无符号表示
		EnableQuickEdit  = 0x0040
		EnableMouseInput = 0x0010
	)

	handle, _, _ := procGetStdHandle.Call(uintptr(StdInputHandle))
	if handle == uintptr(syscall.InvalidHandle) {
		return errors.New("failed to get stdin handle")
	}

	var mode uint32
	r, _, err := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return errors.New("failed to get console mode: " + err.Error())
	}

	newMode := mode &^ (EnableQuickEdit | EnableMouseInput)

	r, _, err = procSetConsoleMode.Call(handle, uintptr(newMode))
	if r == 0 {
		return errors.New("failed to set console mode: " + err.Error())
	}

	return nil
}

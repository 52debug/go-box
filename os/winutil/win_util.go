//go:build windows

package winutil

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	createMutex = kernel32.NewProc("CreateMutexW")
	closeHandle = kernel32.NewProc("CloseHandle")
)

// ExistsInstance 是否存在指定的实例
func ExistsInstance(mutexName string) bool {
	globalName := "Global\\" + mutexName
	_, alreadyExists, err := CreateAppMutex(globalName)
	if err != nil {
		// 保守处理：当作已经有实例
		return true
	}

	// 注意：这里没有关闭句柄，会造成句柄泄漏
	// 只适合非常短生命周期的程序
	return alreadyExists
}

// CreateAppMutex 创建一个命名互斥体，用于实现程序单实例
// 参数：
//
//	name - 互斥体名称，建议使用全局唯一格式，例如：
//	       "Global\\MyCompany-MyApp-SingleInstance-2025"
//	       使用 Global\\ 前缀可在不同会话（如服务/RDP）之间共享
//
// 返回：
//
//	handle       - 互斥体句柄（成功时有效，需要在程序退出时关闭）
//	alreadyExists - 是否已经有其他实例在运行（true=已经有实例）
//	err          - 错误（创建失败时非nil）
func CreateAppMutex(name string) (syscall.Handle, bool, error) {
	if name == "" {
		return 0, false, errors.New("互斥体名称不能为空")
	}

	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, false, fmt.Errorf("UTF16字符串转换失败: %v", err)
	}

	// 调用 CreateMutexW
	// 参数说明：
	//   lpMutexAttributes = nil (0)    → 默认安全属性
	//   bInitialOwner     = false (0)  → 不初始拥有（更安全）
	//   lpName            = 名称指针
	handle, _, err := createMutex.Call(
		0,
		0,
		uintptr(unsafe.Pointer(namePtr)),
	)

	if handle == 0 { // INVALID_HANDLE_VALUE
		return 0, false, fmt.Errorf("CreateMutexW 失败: %v", err)
	}

	// 判断是否是已经存在的互斥体
	alreadyExists := err != nil && errors.Is(err, syscall.ERROR_ALREADY_EXISTS)

	// 如果是新创建的互斥体，我们拥有它（但这里我们不主动 Release，保持占有直到程序退出）
	return syscall.Handle(handle), alreadyExists, nil
}

// CloseMutex 安全关闭互斥体句柄
func CloseMutex(handle syscall.Handle) error {
	if handle == 0 {
		return nil
	}
	r1, _, err := closeHandle.Call(uintptr(handle))
	if r1 == 0 {
		return fmt.Errorf("CloseHandle 失败: %v", err)
	}
	return nil
}

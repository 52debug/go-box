package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// GetProjectRoot 获取程序根目录
func GetProjectRoot() string {
	exePath, err := os.Executable()
	if err != nil {
		// 回退方案
		wd, _ := os.Getwd()
		return wd
	}

	dir := filepath.Dir(exePath)

	// 如果看起来像是 go run 的临时路径 → 改用当前工作目录
	if strings.Contains(strings.ToLower(dir), "tmp") ||
		strings.Contains(dir, "go-build") {
		wd, _ := os.Getwd()
		return wd
	}

	return dir
}

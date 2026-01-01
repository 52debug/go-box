package fsutil

import (
	"os"
	"path/filepath"
	"strings"
)

// IsJSONFile 判断是否为 JSON 文件
func IsJSONFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".json"
}

// IsFile 判断是否为文件（非目录）
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// IsDir 判断是否为目录（非文件）
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// FileExists 文件是否存在
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// DirExists 目录是否存在
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// PathExists 路径是否存在 文件或目录
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// CreateDir 递归创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

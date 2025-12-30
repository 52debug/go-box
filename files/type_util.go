package files

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

package ioutil

import (
	"bufio"
	"log/slog"
	"os"
)

// ReadFileBytes 读取文件内容，返回 []byte
func ReadFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFileString 读取文件内容，返回 string
func ReadFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadFileStringRemoveBOM 读取文件内容并移除 UTF-8 BOM，返回 string
func ReadFileStringRemoveBOM(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// 检测并移除 UTF-8 BOM (EF BB BF)
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return string(data[3:]), nil
	}

	return string(data), nil
}

// ReadLines 读取文件的所有行
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			slog.Info("Error while closing file %s", path)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

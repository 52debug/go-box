package uuidmgr

import (
	"encoding/binary"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUUIDv4 生成V4 版本UUID
func GenerateUUIDv4(withHyphen bool) string {
	return formatUUID(uuid.New(), withHyphen)
}

// GenerateUUIDv7 生成V7 版本UUID
func GenerateUUIDv7(withHyphen bool) (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return formatUUID(u, withHyphen), nil
}

// ExtractV7Timestamp 解析 V7版本的 时间戳
func ExtractV7Timestamp(uuidStr string) (time.Time, bool) {
	// 解析 UUID 字符串
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return time.Time{}, false
	}

	// 检查是否为 v7 UUID
	if u.Version() != 7 {
		return time.Time{}, false
	}

	// 前 6 字节是大端序的毫秒时间戳
	timestampMs := int64(binary.BigEndian.Uint64(append(u[:6], 0, 0))) >> 16
	return time.UnixMilli(timestampMs), true
}

func formatUUID(u uuid.UUID, withHyphen bool) string {
	if withHyphen {
		return u.String()
	}
	return strings.ReplaceAll(u.String(), "-", "")
}

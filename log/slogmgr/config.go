package slogmgr

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/52debug/go-box/log/logmgr"
	"gopkg.in/natefinch/lumberjack.v2"
)

// getHandlerOption 设置 HandlerOptions
func getHandlerOption(level slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 修改时间格式
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))
			}
			// 修改级别字段为小写
			if a.Key == slog.LevelKey {
				logLevel := a.Value.Any().(slog.Level)
				a.Value = slog.StringValue(strings.ToLower(logLevel.String()))
			}
			return a
		},
	}
}

// newFileWriter 创建文件写入器
func newFileWriter(config logmgr.LogConfig) io.Writer {
	// 确保日志目录存在
	dir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic("创建日志目录失败: " + err.Error())
	}

	// 使用 lumberjack 实现日志轮转
	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

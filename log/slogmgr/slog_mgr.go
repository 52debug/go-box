package slogmgr

import (
	"io"
	"log/slog"
	"os"
)

func Setup(config SLogConfig) {
	var level slog.Level
	switch config.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var writer io.Writer
	switch config.Output {
	case "console":
		writer = os.Stdout
	case "file":
		writer = newFileWriter(config)
	case "both":
		fileWriter := newFileWriter(config)
		writer = io.MultiWriter(os.Stdout, fileWriter)
	default:
		writer = os.Stdout
	}

	handlerOpt := getHandlerOption(level)

	// 创建 JSON 格式处理器
	handler := slog.NewJSONHandler(writer, handlerOpt)
	// 设置默认日志记录器
	slog.SetDefault(slog.New(handler))
}

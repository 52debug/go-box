package slogmgr

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type SLogConfig struct {
	Level      string // 日志级别: debug, info, warn, error
	Output     string // 输出位置: console, file, both
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 最大保留日志文件数
	MaxAge     int    // 最大保留天数
	Compress   bool   // 是否压缩
}

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

	handlerOpt := &slog.HandlerOptions{
		Level:     level,
		AddSource: true, // 添加调用源文件和行号，在高频日志场景会有性能开销
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
			// 只保留文件名和包路径
			if a.Key == slog.SourceKey {
				src := a.Value.Any().(*slog.Source)
				// 只保留文件名
				src.File = filepath.Base(src.File)

				// 只保留函数名部分
				fnName := src.Function
				if idx := strings.LastIndex(fnName, "."); idx != -1 {
					fnName = fnName[idx+1:]
				}
				src.Function = fnName
			}
			return a
		},
	}

	// 创建 JSON 格式处理器
	handler := slog.NewJSONHandler(writer, handlerOpt)
	// 设置默认日志记录器
	slog.SetDefault(slog.New(handler))
}

// newFileWriter 创建文件写入器
func newFileWriter(config SLogConfig) io.Writer {
	// 确保日志目录存在
	dir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
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

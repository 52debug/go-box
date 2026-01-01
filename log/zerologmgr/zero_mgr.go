package zerologmgr

import (
	"io"
	"os"
	"path/filepath"

	"github.com/52debug/go-box/log/logmgr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Setup(config logmgr.LogConfig) {
	// 设置全局日志级别
	var level zerolog.Level
	switch config.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

	// 构造输出 writer
	var writers []io.Writer

	if config.Output == "console" || config.Output == "both" {
		// 彩色控制台输出
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: zerolog.TimeFieldFormat,
		}
		writers = append(writers, consoleWriter)
	}

	if config.Output == "file" || config.Output == "both" {
		// 文件滚动输出
		fileWriter := newFileWriter(config)
		writers = append(writers, fileWriter)
	}

	if len(writers) == 0 {
		// 如果没有输出目标，直接丢弃
		log.Logger = zerolog.Nop()
		return
	}

	// 合并多个输出
	multiWriter := io.MultiWriter(writers...)
	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()
}

// newFileWriter 创建滚动文件写入器
func newFileWriter(config logmgr.LogConfig) io.Writer {
	// 确保日志目录存在
	dir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic("创建日志目录失败: " + err.Error())
	}

	// 使用 lumberjack 实现滚动
	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

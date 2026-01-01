package zaplogmgr

import (
	"github.com/52debug/go-box/log/logmgr"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Setup(config logmgr.LogConfig) {
	level := getLogLevel(config.Level)

	var cores []zapcore.Core

	if config.Output == "console" || config.Output == "both" {
		consoleEncoder := newColoredConsoleEncoder()
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(colorable.NewColorableStdout()), level)
		cores = append(cores, consoleCore)
	}

	if config.Output == "file" || config.Output == "both" {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjackLogger), level)
		cores = append(cores, fileCore)
	}

	if len(cores) == 0 {
		// 创建空 logger
		nopLogger := zap.NewNop()
		zap.ReplaceGlobals(nopLogger) // 替换全局 logger
		return
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)
}

func getLevelColor(level zapcore.Level) string {
	switch level {
	case zapcore.DebugLevel:
		return logmgr.ColorSource
	case zapcore.InfoLevel:
		return logmgr.ColorInfo
	case zapcore.WarnLevel:
		return logmgr.ColorWarn
	case zapcore.ErrorLevel:
		return logmgr.ColorError
	default:
		return logmgr.ColorReset
	}
}

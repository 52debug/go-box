package zaplogmgr

import (
	"github.com/52debug/go-box/log/logmgr"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufPool = buffer.NewPool()

type coloredConsoleEncoder struct {
	zapcore.Encoder
}

func newColoredConsoleEncoder() zapcore.Encoder {
	return &coloredConsoleEncoder{
		Encoder: zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "",
			NameKey:        "",
			CallerKey:      "",
			FunctionKey:    "",
			MessageKey:     "",
			StacktraceKey:  "",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
	}
}

func (e *coloredConsoleEncoder) Clone() zapcore.Encoder {
	return &coloredConsoleEncoder{
		Encoder: e.Encoder.Clone(),
	}
}

func (e *coloredConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := bufPool.Get()

	buf.AppendString(logmgr.ColorDebug)
	buf.AppendString(entry.Time.Format("2006-01-02 15:04:05.000"))
	buf.AppendString(logmgr.ColorReset)
	buf.AppendString(" [")

	buf.AppendString(getLevelColor(entry.Level))
	buf.AppendString(entry.Level.CapitalString())
	buf.AppendString(logmgr.ColorReset)

	buf.AppendString("] ")

	if entry.Caller.Defined {
		buf.AppendString(logmgr.ColorTime)
		buf.AppendString(entry.Caller.TrimmedPath())
		buf.AppendString(logmgr.ColorReset)
		buf.AppendString(" ")
	}

	buf.AppendString(logmgr.ColorMessage)
	buf.AppendString(entry.Message)
	buf.AppendString(logmgr.ColorReset)

	for _, field := range fields {
		buf.AppendString(" ")
		field.AddTo(e.Encoder)
	}

	buf.AppendString("\n")

	return buf, nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

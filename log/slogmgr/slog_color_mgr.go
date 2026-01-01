package slogmgr

import (
	"context"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/52debug/go-box/log/logmgr"
)

func SetupWithColor(config logmgr.LogConfig) {
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

	handlerOpt := getHandlerOption(level)

	var handler slog.Handler

	switch config.Output {
	case "console":
		// 控制台输出使用带颜色的文本格式
		handler = newColorHandler(nil)
	case "file":
		// 文件输出使用 JSON 格式
		writer := newFileWriter(config)
		handler = slog.NewJSONHandler(writer, handlerOpt)
	case "both":
		// 控制台使用带颜色的文本格式，文件使用 JSON 格式
		fileWriter := newFileWriter(config)

		// 控制台处理器（带颜色）
		var consoleHandler slog.Handler
		consoleHandler = newColorHandler(nil)

		// 文件处理器（JSON 格式）
		var fileHandler slog.Handler
		fileHandler = slog.NewJSONHandler(fileWriter, handlerOpt)

		// 合并处理器
		handler = &multiHandler{
			handlers: []slog.Handler{consoleHandler, fileHandler},
		}
	default:
		// 默认使用带颜色的文本格式
		handler = newColorHandler(nil)
	}

	// 设置默认日志记录器
	slog.SetDefault(slog.New(handler))
}

// colorHandler 直接输出带颜色的文本日志

type colorHandler struct{}

func newColorHandler(_ slog.Handler) slog.Handler {
	return &colorHandler{}
}

func (ch *colorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// 始终启用所有级别
	return true
}

func (ch *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	// 对于控制台输出，我们将使用自定义的彩色文本格式
	// 获取颜色
	var levelColor string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = logmgr.ColorDebug
	case slog.LevelInfo:
		levelColor = logmgr.ColorInfo
	case slog.LevelWarn:
		levelColor = logmgr.ColorWarn
	case slog.LevelError:
		levelColor = logmgr.ColorError
	default:
		levelColor = logmgr.ColorReset
	}

	// 格式化时间
	timeStr := r.Time.Format("2006-01-02 15:04:05.000")

	// 格式化源信息
	var sourceStr string
	if r.PC != 0 {
		src := r.Source()
		if src != nil {
			// 只保留文件名和函数名
			filename := filepath.Base(src.File)
			funcName := src.Function
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
			sourceStr = logmgr.ColorSource + " " + filename + ":" + funcName + ":" + strconv.Itoa(src.Line) + logmgr.ColorReset
		}
	}

	// 构建彩色日志行
	levelStr := strings.ToUpper(r.Level.String())
	logLine := logmgr.ColorTime + "[" + timeStr + "]" + logmgr.ColorReset + " " +
		levelColor + "[" + levelStr + "]" + logmgr.ColorReset +
		sourceStr + " " +
		levelColor + r.Message + logmgr.ColorReset

	// 添加属性
	r.Attrs(func(a slog.Attr) bool {
		logLine += " " + logmgr.ColorMessage + a.Key + "=" + logmgr.ColorReset + a.Value.String()
		return true
	})

	// 输出到控制台
	println(logLine)

	return nil
}

func (ch *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 简单返回自身，因为我们不处理属性的持久化
	return ch
}

func (ch *colorHandler) WithGroup(name string) slog.Handler {
	// 简单返回自身，因为我们不处理分组
	return ch
}

// multiHandler 允许多个处理器处理同一个日志记录

type multiHandler struct {
	handlers []slog.Handler
}

func (mh *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range mh.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (mh *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error
	for _, h := range mh.handlers {
		if h.Enabled(ctx, r.Level) {
			if e := h.Handle(ctx, r); e != nil {
				err = e
			}
		}
	}
	return err
}

func (mh *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (mh *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

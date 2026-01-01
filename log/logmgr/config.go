package logmgr

// 颜色常量
const (
	ColorReset   = "\033[0m"
	ColorTime    = "\033[35m" // Magenta
	ColorDebug   = "\033[36m" // Cyan
	ColorInfo    = "\033[32m" // Green
	ColorWarn    = "\033[33m" // Yellow
	ColorError   = "\033[31m" // Red
	ColorSource  = "\033[34m" // Blue
	ColorMessage = "\033[37m" // White
)

// LogConfig 日志配置
type LogConfig struct {
	Level      string // 日志级别: debug, info, warn, error
	Output     string // 输出位置: console, file, both
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 最大保留日志文件数
	MaxAge     int    // 最大保留天数
	Compress   bool   // 是否压缩
}

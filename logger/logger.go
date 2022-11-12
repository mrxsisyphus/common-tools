package logger

import "go.uber.org/zap"

var (
	//defaultLogger 默认的logger 用于打印正常的日志信息
	globalLogger *Logger
)

func init() {
	globalLogger = newDefaultLogger()
}

// Logger
type Logger struct {
	zapLogger *zap.Logger //zap结构化日志对象
	opts      *LogOptions //日志配置
}

// newDefaultLogger 新的默认Logger
func newDefaultLogger() *Logger {
	zapLogger := NewDefaultConsoleZapLogger()
	log := &Logger{
		zapLogger: zapLogger,
	}
	return log
}

// GetLogger 获得全局的logger(默认情况下获得defaultLogger)
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger = newDefaultLogger()
	}
	return globalLogger
}

// SetGlobalLogger 设置全局logger
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

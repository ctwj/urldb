package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String 返回级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 简化的日志器
type Logger struct {
	level  LogLevel
	logger *log.Logger
	file   *os.File
	mu     sync.RWMutex
}

var (
	globalLogger *Logger
	loggerOnce   sync.Once
)

// InitLogger 初始化日志器
func InitLogger() error {
	var err error
	loggerOnce.Do(func() {
		globalLogger = &Logger{
			level:  INFO,
			logger: log.New(os.Stdout, "", log.LstdFlags),
		}

		// 创建日志目录
		logDir := "logs"
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		// 创建日志文件
		logFile := filepath.Join(logDir, "app.log")
		globalLogger.file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}

		// 同时输出到控制台和文件
		globalLogger.logger = log.New(io.MultiWriter(os.Stdout, globalLogger.file), "", log.LstdFlags)
	})

	return err
}

// GetLogger 获取全局日志器
func GetLogger() *Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}

// log 内部日志方法
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] [%s] %s", level.String(), caller, message)

	l.logger.Println(logMessage)

	// Fatal级别终止程序
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// TelegramDebug Telegram调试日志
func (l *Logger) TelegramDebug(format string, args ...interface{}) {
	l.log(DEBUG, "[TELEGRAM] "+format, args...)
}

// TelegramInfo Telegram信息日志
func (l *Logger) TelegramInfo(format string, args ...interface{}) {
	l.log(INFO, "[TELEGRAM] "+format, args...)
}

// TelegramWarn Telegram警告日志
func (l *Logger) TelegramWarn(format string, args ...interface{}) {
	l.log(WARN, "[TELEGRAM] "+format, args...)
}

// TelegramError Telegram错误日志
func (l *Logger) TelegramError(format string, args ...interface{}) {
	l.log(ERROR, "[TELEGRAM] "+format, args...)
}

// DebugWithFields 带字段的调试日志
func (l *Logger) DebugWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if len(fields) > 0 {
		var fieldStrs []string
		for k, v := range fields {
			fieldStrs = append(fieldStrs, fmt.Sprintf("%s=%v", k, v))
		}
		message = fmt.Sprintf("%s [%s]", message, strings.Join(fieldStrs, ", "))
	}
	l.log(DEBUG, message)
}

// InfoWithFields 带字段的信息日志
func (l *Logger) InfoWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if len(fields) > 0 {
		var fieldStrs []string
		for k, v := range fields {
			fieldStrs = append(fieldStrs, fmt.Sprintf("%s=%v", k, v))
		}
		message = fmt.Sprintf("%s [%s]", message, strings.Join(fieldStrs, ", "))
	}
	l.log(INFO, message)
}

// ErrorWithFields 带字段的错误日志
func (l *Logger) ErrorWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if len(fields) > 0 {
		var fieldStrs []string
		for k, v := range fields {
			fieldStrs = append(fieldStrs, fmt.Sprintf("%s=%v", k, v))
		}
		message = fmt.Sprintf("%s [%s]", message, strings.Join(fieldStrs, ", "))
	}
	l.log(ERROR, message)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
	}
}

// 全局便捷函数
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
}

func TelegramDebug(format string, args ...interface{}) {
	GetLogger().TelegramDebug(format, args...)
}

func TelegramInfo(format string, args ...interface{}) {
	GetLogger().TelegramInfo(format, args...)
}

func TelegramWarn(format string, args ...interface{}) {
	GetLogger().TelegramWarn(format, args...)
}

func TelegramError(format string, args ...interface{}) {
	GetLogger().TelegramError(format, args...)
}

func DebugWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().DebugWithFields(fields, format, args...)
}

func InfoWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().InfoWithFields(fields, format, args...)
}

func ErrorWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().ErrorWithFields(fields, format, args...)
}

// Min 返回两个整数中的较小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


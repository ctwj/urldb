package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
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

// String 返回日志级别的字符串表示
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

// Logger 统一日志器
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger

	file   *os.File
	mu     sync.Mutex
	config *LogConfig
}

// LogConfig 日志配置
type LogConfig struct {
	LogDir         string   // 日志目录
	LogLevel       LogLevel // 日志级别
	MaxFileSize    int64    // 单个日志文件最大大小（MB）
	MaxBackups     int      // 最大备份文件数
	MaxAge         int      // 日志文件最大保留天数
	EnableConsole  bool     // 是否启用控制台输出
	EnableFile     bool     // 是否启用文件输出
	EnableRotation bool     // 是否启用日志轮转
}

// DefaultConfig 默认配置
func DefaultConfig() *LogConfig {
	return &LogConfig{
		LogDir:         "logs",
		LogLevel:       INFO,
		MaxFileSize:    100, // 100MB
		MaxBackups:     5,
		MaxAge:         30, // 30天
		EnableConsole:  true,
		EnableFile:     true,
		EnableRotation: true,
	}
}

var (
	globalLogger *Logger
	onceLogger   sync.Once
)

// InitLogger 初始化全局日志器
func InitLogger(config *LogConfig) error {
	var err error
	onceLogger.Do(func() {
		if config == nil {
			config = DefaultConfig()
		}

		globalLogger, err = NewLogger(config)
	})
	return err
}

// GetLogger 获取全局日志器
func GetLogger() *Logger {
	if globalLogger == nil {
		InitLogger(nil)
	}
	return globalLogger
}

// NewLogger 创建新的日志器
func NewLogger(config *LogConfig) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	logger := &Logger{
		config: config,
	}

	// 创建日志目录
	if config.EnableFile {
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			return nil, fmt.Errorf("创建日志目录失败: %v", err)
		}
	}

	// 初始化日志文件
	if err := logger.initLogFile(); err != nil {
		return nil, err
	}

	// 初始化日志器
	logger.initLoggers()

	// 启动日志轮转检查
	if config.EnableRotation {
		go logger.startRotationCheck()
	}

	return logger, nil
}

// initLogFile 初始化日志文件
func (l *Logger) initLogFile() error {
	if !l.config.EnableFile {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 关闭现有文件
	if l.file != nil {
		l.file.Close()
	}

	// 创建新的日志文件
	logFile := filepath.Join(l.config.LogDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}

	l.file = file
	return nil
}

// initLoggers 初始化各个级别的日志器
func (l *Logger) initLoggers() {
	var writers []io.Writer

	// 添加控制台输出
	if l.config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	// 添加文件输出
	if l.config.EnableFile && l.file != nil {
		writers = append(writers, l.file)
	}

	multiWriter := io.MultiWriter(writers...)

	// 创建各个级别的日志器
	l.debugLogger = log.New(multiWriter, "[DEBUG] ", log.LstdFlags)
	l.infoLogger = log.New(multiWriter, "[INFO] ", log.LstdFlags)
	l.warnLogger = log.New(multiWriter, "[WARN] ", log.LstdFlags)
	l.errorLogger = log.New(multiWriter, "[ERROR] ", log.LstdFlags)
	l.fatalLogger = log.New(multiWriter, "[FATAL] ", log.LstdFlags)
}

// log 内部日志方法
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.config.LogLevel {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	// 提取文件名
	fileName := filepath.Base(file)

	// 格式化消息
	message := fmt.Sprintf(format, args...)

	// 添加调用位置信息
	fullMessage := fmt.Sprintf("[%s:%d] %s", fileName, line, message)

	switch level {
	case DEBUG:
		l.debugLogger.Println(fullMessage)
	case INFO:
		l.infoLogger.Println(fullMessage)
	case WARN:
		l.warnLogger.Println(fullMessage)
	case ERROR:
		l.errorLogger.Println(fullMessage)
	case FATAL:
		l.fatalLogger.Println(fullMessage)
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

// startRotationCheck 启动日志轮转检查
func (l *Logger) startRotationCheck() {
	ticker := time.NewTicker(1 * time.Hour) // 每小时检查一次
	defer ticker.Stop()

	for range ticker.C {
		l.checkRotation()
	}
}

// checkRotation 检查是否需要轮转日志
func (l *Logger) checkRotation() {
	if !l.config.EnableFile || l.file == nil {
		return
	}

	// 检查文件大小
	fileInfo, err := l.file.Stat()
	if err != nil {
		return
	}

	// 如果文件超过最大大小，进行轮转
	if fileInfo.Size() > l.config.MaxFileSize*1024*1024 {
		l.rotateLog()
	}

	// 清理旧日志文件
	l.cleanOldLogs()
}

// rotateLog 轮转日志文件
func (l *Logger) rotateLog() {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 关闭当前文件
	if l.file != nil {
		l.file.Close()
	}

	// 重命名当前日志文件
	currentLogFile := filepath.Join(l.config.LogDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	backupLogFile := filepath.Join(l.config.LogDir, fmt.Sprintf("app_%s_%s.log", time.Now().Format("2006-01-02"), time.Now().Format("15-04-05")))

	if _, err := os.Stat(currentLogFile); err == nil {
		os.Rename(currentLogFile, backupLogFile)
	}

	// 创建新的日志文件
	l.initLogFile()
	l.initLoggers()
}

// cleanOldLogs 清理旧日志文件
func (l *Logger) cleanOldLogs() {
	if l.config.MaxAge <= 0 {
		return
	}

	files, err := filepath.Glob(filepath.Join(l.config.LogDir, "app_*.log"))
	if err != nil {
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -l.config.MaxAge)

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			os.Remove(file)
		}
	}
}

// Close 关闭日志器
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
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

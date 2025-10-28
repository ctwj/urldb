package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

// StructuredLogEntry 结构化日志条目
type StructuredLogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Caller    string            `json:"caller"`
	Module    string            `json:"module"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
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
	StructuredLog  bool     // 是否启用结构化日志格式
}

// DefaultConfig 默认配置
func DefaultConfig() *LogConfig {
	// 从环境变量获取日志级别，默认为INFO
	logLevel := getLogLevelFromEnv()

	return &LogConfig{
		LogDir:         "logs",
		LogLevel:       logLevel,
		MaxFileSize:    100, // 100MB
		MaxBackups:     5,
		MaxAge:         30, // 30天
		EnableConsole:  true,
		EnableFile:     true,
		EnableRotation: true,
		StructuredLog:  os.Getenv("STRUCTURED_LOG") == "true", // 从环境变量控制结构化日志
	}
}

// getLogLevelFromEnv 从环境变量获取日志级别
func getLogLevelFromEnv() LogLevel {
	envLogLevel := os.Getenv("LOG_LEVEL")
	envDebug := os.Getenv("DEBUG")

	// 如果设置了DEBUG环境变量为true，则使用DEBUG级别
	if envDebug == "true" || envDebug == "1" {
		return DEBUG
	}

	// 根据LOG_LEVEL环境变量设置日志级别
	switch strings.ToUpper(envLogLevel) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		// 根据运行环境设置默认级别：开发环境DEBUG，生产环境INFO
		if isDevelopment() {
			return DEBUG
		}
		return INFO
	}
}

// isDevelopment 判断是否为开发环境
func isDevelopment() bool {
	env := os.Getenv("GO_ENV")
	return env == "development" || env == "dev" || env == "local" || env == "test"
}

// getEnvironment 获取当前环境类型
func (l *Logger) getEnvironment() string {
	if isDevelopment() {
		return "development"
	}
	return "production"
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

	// 打印日志配置信息
	logger.Info("日志系统初始化完成 - 级别: %s, 环境: %s",
		config.LogLevel.String(),
		logger.getEnvironment())

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
	logFile := filepath.Join(l.config.LogDir, fmt.Sprintf("app_%s.log", GetCurrentTime().Format("2006-01-02")))
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

	// 提取文件名作为模块名
	fileName := filepath.Base(file)
	moduleName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 格式化消息
	message := fmt.Sprintf(format, args...)

	// 添加调用位置信息
	caller := fmt.Sprintf("%s:%d", fileName, line)

	if l.config.StructuredLog {
		// 结构化日志格式
		entry := StructuredLogEntry{
			Timestamp: GetCurrentTime(),
			Level:     level.String(),
			Message:   message,
			Caller:    caller,
			Module:    moduleName,
		}

		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			// 如果JSON序列化失败，回退到普通格式
			fullMessage := fmt.Sprintf("[%s] [%s:%d] %s", level.String(), fileName, line, message)
			l.logToLevel(level, fullMessage)
			return
		}

		l.logToLevel(level, string(jsonBytes))
	} else {
		// 普通文本格式
		fullMessage := fmt.Sprintf("[%s] [%s:%d] %s", level.String(), fileName, line, message)
		l.logToLevel(level, fullMessage)
	}
}

// logToLevel 根据级别输出日志
func (l *Logger) logToLevel(level LogLevel, message string) {
	switch level {
	case DEBUG:
		l.debugLogger.Println(message)
	case INFO:
		l.infoLogger.Println(message)
	case WARN:
		l.warnLogger.Println(message)
	case ERROR:
		l.errorLogger.Println(message)
	case FATAL:
		l.fatalLogger.Println(message)
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
	currentLogFile := filepath.Join(l.config.LogDir, fmt.Sprintf("app_%s.log", GetCurrentTime().Format("2006-01-02")))
	backupLogFile := filepath.Join(l.config.LogDir, fmt.Sprintf("app_%s_%s.log", GetCurrentTime().Format("2006-01-02"), GetCurrentTime().Format("15-04-05")))

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

	cutoffTime := GetCurrentTime().AddDate(0, 0, -l.config.MaxAge)

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

// Min 返回两个整数中的较小值
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 结构化日志方法
func (l *Logger) DebugWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	l.logWithFields(DEBUG, fields, format, args...)
}

func (l *Logger) InfoWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	l.logWithFields(INFO, fields, format, args...)
}

func (l *Logger) WarnWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	l.logWithFields(WARN, fields, format, args...)
}

func (l *Logger) ErrorWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	l.logWithFields(ERROR, fields, format, args...)
}

func (l *Logger) FatalWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	l.logWithFields(FATAL, fields, format, args...)
}

// logWithFields 带字段的结构化日志方法
func (l *Logger) logWithFields(level LogLevel, fields map[string]interface{}, format string, args ...interface{}) {
	if level < l.config.LogLevel {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}

	// 提取文件名作为模块名
	fileName := filepath.Base(file)
	moduleName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// 格式化消息
	message := fmt.Sprintf(format, args...)

	// 添加调用位置信息
	caller := fmt.Sprintf("%s:%d", fileName, line)

	if l.config.StructuredLog {
		// 结构化日志格式
		entry := StructuredLogEntry{
			Timestamp: GetCurrentTime(),
			Level:     level.String(),
			Message:   message,
			Caller:    caller,
			Module:    moduleName,
			Fields:    fields,
		}

		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			// 如果JSON序列化失败，回退到普通格式
			fullMessage := fmt.Sprintf("[%s] [%s:%d] %s - Fields: %v", level.String(), fileName, line, message, fields)
			l.logToLevel(level, fullMessage)
			return
		}

		l.logToLevel(level, string(jsonBytes))
	} else {
		// 普通文本格式
		fullMessage := fmt.Sprintf("[%s] [%s:%d] %s - Fields: %v", level.String(), fileName, line, message, fields)
		l.logToLevel(level, fullMessage)
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

// 全局结构化日志便捷函数
func DebugWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().DebugWithFields(fields, format, args...)
}

func InfoWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().InfoWithFields(fields, format, args...)
}

func WarnWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().WarnWithFields(fields, format, args...)
}

func ErrorWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().ErrorWithFields(fields, format, args...)
}

func FatalWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	GetLogger().FatalWithFields(fields, format, args...)
}

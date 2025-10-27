package monitor

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// ErrorInfo 错误信息结构
type ErrorInfo struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
	StackTrace  string    `json:"stack_trace"`
	RequestInfo *RequestInfo `json:"request_info,omitempty"`
	Level       string    `json:"level"` // error, warn, info
	Count       int       `json:"count"`
}

// RequestInfo 请求信息结构
type RequestInfo struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	RemoteAddr  string            `json:"remote_addr"`
	UserAgent   string            `json:"user_agent"`
	RequestBody string            `json:"request_body"`
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	errors     map[string]*ErrorInfo
	mu         sync.RWMutex
	maxErrors  int
	retention  time.Duration
}

// NewErrorHandler 创建新的错误处理器
func NewErrorHandler(maxErrors int, retention time.Duration) *ErrorHandler {
	eh := &ErrorHandler{
		errors:    make(map[string]*ErrorInfo),
		maxErrors: maxErrors,
		retention: retention,
	}

	// 启动错误清理协程
	go eh.cleanupRoutine()

	return eh
}

// RecoverMiddleware panic恢复中间件
func (eh *ErrorHandler) RecoverMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误信息
				stackTrace := getStackTrace()

				errorInfo := &ErrorInfo{
					ID:         fmt.Sprintf("panic_%d", time.Now().UnixNano()),
					Timestamp:  time.Now(),
					Message:    fmt.Sprintf("%v", err),
					StackTrace: stackTrace,
					RequestInfo: &RequestInfo{
						Method:     c.Request.Method,
						URL:        c.Request.URL.String(),
						RemoteAddr: c.ClientIP(),
						UserAgent:  c.GetHeader("User-Agent"),
					},
					Level: "error",
					Count: 1,
				}

				// 保存错误信息
				eh.saveError(errorInfo)

				utils.Error("Panic recovered: %v\nStack trace: %s", err, stackTrace)

				// 返回错误响应
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"code":  "INTERNAL_ERROR",
				})

				// 不继续处理
				c.Abort()
			}
		}()

		c.Next()
	}
}

// ErrorMiddleware 通用错误处理中间件
func (eh *ErrorHandler) ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				errorInfo := &ErrorInfo{
					ID:        fmt.Sprintf("error_%d_%s", time.Now().UnixNano(), ginErr.Type),
					Timestamp: time.Now(),
					Message:   ginErr.Error(),
					Level:     "error",
					Count:     1,
					RequestInfo: &RequestInfo{
						Method:     c.Request.Method,
						URL:        c.Request.URL.String(),
						RemoteAddr: c.ClientIP(),
						UserAgent:  c.GetHeader("User-Agent"),
					},
				}

				eh.saveError(errorInfo)
			}
		}
	}
}

// saveError 保存错误信息
func (eh *ErrorHandler) saveError(errorInfo *ErrorInfo) {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	key := errorInfo.Message
	if existing, exists := eh.errors[key]; exists {
		// 如果错误已存在，增加计数
		existing.Count++
		existing.Timestamp = time.Now()
	} else {
		// 如果是新错误，添加到映射中
		eh.errors[key] = errorInfo
	}

	// 如果错误数量超过限制，清理旧错误
	if len(eh.errors) > eh.maxErrors {
		eh.cleanupOldErrors()
	}
}

// GetErrors 获取错误列表
func (eh *ErrorHandler) GetErrors() []*ErrorInfo {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	errors := make([]*ErrorInfo, 0, len(eh.errors))
	for _, errorInfo := range eh.errors {
		errors = append(errors, errorInfo)
	}

	return errors
}

// GetErrorByID 根据ID获取错误
func (eh *ErrorHandler) GetErrorByID(id string) (*ErrorInfo, bool) {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	for _, errorInfo := range eh.errors {
		if errorInfo.ID == id {
			return errorInfo, true
		}
	}

	return nil, false
}

// ClearErrors 清空所有错误
func (eh *ErrorHandler) ClearErrors() {
	eh.mu.Lock()
	defer eh.mu.Unlock()

	eh.errors = make(map[string]*ErrorInfo)
}

// cleanupOldErrors 清理旧错误
func (eh *ErrorHandler) cleanupOldErrors() {
	// 简单策略：保留最近的错误，删除旧的
	errors := make([]*ErrorInfo, 0, len(eh.errors))
	for _, errorInfo := range eh.errors {
		errors = append(errors, errorInfo)
	}

	// 按时间戳排序
	for i := 0; i < len(errors)-1; i++ {
		for j := i + 1; j < len(errors); j++ {
			if errors[i].Timestamp.Before(errors[j].Timestamp) {
				errors[i], errors[j] = errors[j], errors[i]
			}
		}
	}

	// 保留最新的maxErrors/2个错误
	keep := eh.maxErrors / 2
	if keep < 1 {
		keep = 1
	}

	if len(errors) > keep {
		// 重建错误映射
		newErrors := make(map[string]*ErrorInfo)
		for i := 0; i < keep; i++ {
			newErrors[errors[i].Message] = errors[i]
		}
		eh.errors = newErrors
	}
}

// cleanupRoutine 定期清理过期错误的协程
func (eh *ErrorHandler) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟清理一次
	defer ticker.Stop()

	for range ticker.C {
		eh.mu.Lock()
		for key, errorInfo := range eh.errors {
			if time.Since(errorInfo.Timestamp) > eh.retention {
				delete(eh.errors, key)
			}
		}
		eh.mu.Unlock()
	}
}

// getStackTrace 获取堆栈跟踪信息
func getStackTrace() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

// GetErrorStatistics 获取错误统计信息
func (eh *ErrorHandler) GetErrorStatistics() map[string]interface{} {
	eh.mu.RLock()
	defer eh.mu.RUnlock()

	totalErrors := len(eh.errors)
	totalCount := 0
	errorTypes := make(map[string]int)

	for _, errorInfo := range eh.errors {
		totalCount += errorInfo.Count
		// 提取错误类型（基于错误消息的前几个单词）
		parts := strings.Split(errorInfo.Message, " ")
		if len(parts) > 0 {
			errorType := strings.Join(parts[:min(3, len(parts))], " ")
			errorTypes[errorType]++
		}
	}

	return map[string]interface{}{
		"total_errors":  totalErrors,
		"total_count":   totalCount,
		"error_types":   errorTypes,
		"max_errors":    eh.maxErrors,
		"retention":     eh.retention,
		"active_errors": len(eh.errors),
	}
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GlobalErrorHandler 全局错误处理器
var globalErrorHandler *ErrorHandler

// InitGlobalErrorHandler 初始化全局错误处理器
func InitGlobalErrorHandler(maxErrors int, retention time.Duration) {
	globalErrorHandler = NewErrorHandler(maxErrors, retention)
}

// GetGlobalErrorHandler 获取全局错误处理器
func GetGlobalErrorHandler() *ErrorHandler {
	if globalErrorHandler == nil {
		InitGlobalErrorHandler(100, 24*time.Hour)
	}
	return globalErrorHandler
}

// Recover 全局panic恢复函数
func Recover() gin.HandlerFunc {
	if globalErrorHandler == nil {
		InitGlobalErrorHandler(100, 24*time.Hour)
	}
	return globalErrorHandler.RecoverMiddleware()
}

// Error 全局错误处理函数
func Error() gin.HandlerFunc {
	if globalErrorHandler == nil {
		InitGlobalErrorHandler(100, 24*time.Hour)
	}
	return globalErrorHandler.ErrorMiddleware()
}

// RecordError 记录错误（全局函数）
func RecordError(message string, level string) {
	if globalErrorHandler == nil {
		InitGlobalErrorHandler(100, 24*time.Hour)
		return
	}

	errorInfo := &ErrorInfo{
		ID:        fmt.Sprintf("%s_%d", level, time.Now().UnixNano()),
		Timestamp: time.Now(),
		Message:   message,
		Level:     level,
		Count:     1,
	}

	globalErrorHandler.saveError(errorInfo)
}
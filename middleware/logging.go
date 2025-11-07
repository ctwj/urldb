package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ctwj/urldb/utils"
)

// responseWriter 包装http.ResponseWriter以捕获响应状态码和内容
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware HTTP请求日志中间件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包装ResponseWriter
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     200, // 默认状态码
		}

		// 读取请求体
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		next.ServeHTTP(rw, r)

		// 计算处理时间
		duration := time.Since(start)

		// 记录请求日志
		logRequest(r, rw, duration, requestBody)
	})
}

// logRequest 记录请求日志 - 优化后仅记录异常和关键请求
func logRequest(r *http.Request, rw *responseWriter, duration time.Duration, requestBody []byte) {
	// 获取客户端IP
	clientIP := getClientIP(r)

	// 判断是否需要记录日志的条件
	shouldLog := rw.statusCode >= 400 || // 错误状态码
		duration > 5*time.Second || // 耗时过长
		shouldLogPath(r.URL.Path) || // 关键路径
		isAdminPath(r.URL.Path) // 管理员路径

	if !shouldLog {
		return // 正常请求不记录日志，减少日志噪音
	}

	// 简化的日志格式，移除User-Agent以减少噪音
	if rw.statusCode >= 400 {
		// 错误请求记录详细信息
		utils.Error("HTTP异常 - %s %s - IP: %s - 状态码: %d - 耗时: %v",
			r.Method, r.URL.Path, clientIP, rw.statusCode, duration)

		// 仅在错误状态下记录简要的请求信息
		if len(requestBody) > 0 && len(requestBody) <= 500 {
			utils.Error("请求详情: %s", string(requestBody))
		}
	} else if duration > 5*time.Second {
		// 慢请求警告
		utils.Warn("HTTP慢请求 - %s %s - IP: %s - 耗时: %v",
			r.Method, r.URL.Path, clientIP, duration)
	} else {
		// 关键路径的正常请求
		utils.Info("HTTP关键请求 - %s %s - IP: %s - 状态码: %d - 耗时: %v",
			r.Method, r.URL.Path, clientIP, rw.statusCode, duration)
	}
}

// shouldLogPath 判断路径是否需要记录日志
func shouldLogPath(path string) bool {
	// 定义需要记录日志的关键路径
	keyPaths := []string{
		"/api/public/resources",
		"/api/admin/config",
		"/api/admin/users",
		"/telegram/webhook",
	}

	for _, keyPath := range keyPaths {
		if strings.HasPrefix(path, keyPath) {
			return true
		}
	}
	return false
}

// isAdminPath 判断是否为管理员路径
func isAdminPath(path string) bool {
	return strings.HasPrefix(path, "/api/admin/") ||
		   strings.HasPrefix(path, "/admin/")
}

// getClientIP 获取客户端真实IP地址
func getClientIP(r *http.Request) string {
	// 检查X-Forwarded-For头
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	// 检查X-Real-IP头
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// 检查X-Client-IP头
	if ip := r.Header.Get("X-Client-IP"); ip != "" {
		return ip
	}

	// 返回远程地址
	return r.RemoteAddr
}

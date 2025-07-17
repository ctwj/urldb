package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/ctwj/panResManage/utils"
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

// logRequest 记录请求日志
func logRequest(r *http.Request, rw *responseWriter, duration time.Duration, requestBody []byte) {
	// 获取客户端IP
	clientIP := getClientIP(r)

	// 获取用户代理
	userAgent := r.UserAgent()
	if userAgent == "" {
		userAgent = "Unknown"
	}

	// 记录请求信息
	utils.Info("HTTP请求 - %s %s - IP: %s - User-Agent: %s - 状态码: %d - 耗时: %v",
		r.Method, r.URL.Path, clientIP, userAgent, rw.statusCode, duration)

	// 如果是错误状态码，记录详细信息
	if rw.statusCode >= 400 {
		utils.Error("HTTP错误 - %s %s - 状态码: %d - 响应体: %s",
			r.Method, r.URL.Path, rw.statusCode, rw.body.String())
	}

	// 记录请求参数（仅对POST/PUT请求）
	if (r.Method == "POST" || r.Method == "PUT") && len(requestBody) > 0 {
		// 限制日志长度，避免日志文件过大
		if len(requestBody) > 1000 {
			utils.Debug("请求体(截断): %s...", string(requestBody[:1000]))
		} else {
			utils.Debug("请求体: %s", string(requestBody))
		}
	}

	// 记录查询参数
	if len(r.URL.RawQuery) > 0 {
		utils.Debug("查询参数: %s", r.URL.RawQuery)
	}
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

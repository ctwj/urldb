package monitor

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics 监控指标
type Metrics struct {
	// HTTP请求指标
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestSize      *prometheus.SummaryVec
	ResponseSize     *prometheus.SummaryVec

	// 数据库指标
	DatabaseQueries  *prometheus.CounterVec
	DatabaseErrors   *prometheus.CounterVec
	DatabaseDuration *prometheus.HistogramVec

	// 系统指标
	MemoryUsage      prometheus.Gauge
	Goroutines       prometheus.Gauge
	GCStats          *prometheus.CounterVec

	// 业务指标
	ResourcesCreated *prometheus.CounterVec
	ResourcesViewed  *prometheus.CounterVec
	Searches         *prometheus.CounterVec
	Transfers        *prometheus.CounterVec

	// 错误指标
	ErrorsTotal      *prometheus.CounterVec

	// 自定义指标
	CustomCounters   map[string]prometheus.Counter
	CustomGauges     map[string]prometheus.Gauge
	mu               sync.RWMutex
}

// MetricsConfig 监控配置
type MetricsConfig struct {
	Enabled        bool
	ListenAddress  string
	MetricsPath    string
	Namespace      string
	Subsystem      string
}

// DefaultMetricsConfig 默认监控配置
func DefaultMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		Enabled:       true,
		ListenAddress: ":9090",
		MetricsPath:   "/metrics",
		Namespace:     "urldb",
		Subsystem:     "api",
	}
}

// GlobalMetrics 全局监控实例
var (
	globalMetrics *Metrics
	once          sync.Once
)

// NewMetrics 创建新的监控指标
func NewMetrics(config *MetricsConfig) *Metrics {
	if config == nil {
		config = DefaultMetricsConfig()
	}

	namespace := config.Namespace
	subsystem := config.Subsystem

	m := &Metrics{
		// HTTP请求指标
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Subsystem:  subsystem,
				Name:       "http_request_size_bytes",
				Help:       "HTTP request size in bytes",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"method", "endpoint"},
		),
		ResponseSize: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Subsystem:  subsystem,
				Name:       "http_response_size_bytes",
				Help:       "HTTP response size in bytes",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"method", "endpoint"},
		),

		// 数据库指标
		DatabaseQueries: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"table", "operation"},
		),
		DatabaseErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "errors_total",
				Help:      "Total number of database errors",
			},
			[]string{"table", "operation", "error"},
		),
		DatabaseDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"table", "operation"},
		),

		// 系统指标
		MemoryUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "system",
				Name:      "memory_usage_bytes",
				Help:      "Current memory usage in bytes",
			},
		),
		Goroutines: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "system",
				Name:      "goroutines",
				Help:      "Number of goroutines",
			},
		),
		GCStats: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "system",
				Name:      "gc_stats_total",
				Help:      "Garbage collection statistics",
			},
			[]string{"type"},
		),

		// 业务指标
		ResourcesCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "resources_created_total",
				Help:      "Total number of resources created",
			},
			[]string{"category", "platform"},
		),
		ResourcesViewed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "resources_viewed_total",
				Help:      "Total number of resources viewed",
			},
			[]string{"category"},
		),
		Searches: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "searches_total",
				Help:      "Total number of searches",
			},
			[]string{"platform"},
		),
		Transfers: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "transfers_total",
				Help:      "Total number of transfers",
			},
			[]string{"platform", "status"},
		),

		// 错误指标
		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "errors",
				Name:      "total",
				Help:      "Total number of errors",
			},
			[]string{"type", "endpoint"},
		),

		// 自定义指标
		CustomCounters: make(map[string]prometheus.Counter),
		CustomGauges:   make(map[string]prometheus.Gauge),
	}

	// 启动系统指标收集
	go m.collectSystemMetrics()

	return m
}

// GetGlobalMetrics 获取全局监控实例
func GetGlobalMetrics() *Metrics {
	once.Do(func() {
		globalMetrics = NewMetrics(DefaultMetricsConfig())
	})
	return globalMetrics
}

// SetGlobalMetrics 设置全局监控实例
func SetGlobalMetrics(metrics *Metrics) {
	globalMetrics = metrics
}

// collectSystemMetrics 收集系统指标
func (m *Metrics) collectSystemMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 收集内存使用情况
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		m.MemoryUsage.Set(float64(ms.Alloc))

		// 收集goroutine数量
		m.Goroutines.Set(float64(runtime.NumGoroutine()))

		// 收集GC统计
		m.GCStats.WithLabelValues("alloc").Add(float64(ms.TotalAlloc))
		m.GCStats.WithLabelValues("sys").Add(float64(ms.Sys))
		m.GCStats.WithLabelValues("lookups").Add(float64(ms.Lookups))
		m.GCStats.WithLabelValues("mallocs").Add(float64(ms.Mallocs))
		m.GCStats.WithLabelValues("frees").Add(float64(ms.Frees))
	}
}

// MetricsMiddleware 监控中间件
func (m *Metrics) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		// 如果没有匹配的路由，使用请求路径
		if path == "" {
			path = c.Request.URL.Path
		}

		// 记录请求大小
		requestSize := float64(c.Request.ContentLength)
		m.RequestSize.WithLabelValues(c.Request.Method, path).Observe(requestSize)

		c.Next()

		// 记录响应信息
		status := c.Writer.Status()
		latency := time.Since(start).Seconds()
		responseSize := float64(c.Writer.Size())

		// 更新指标
		m.RequestsTotal.WithLabelValues(c.Request.Method, path, fmt.Sprintf("%d", status)).Inc()
		m.RequestDuration.WithLabelValues(c.Request.Method, path, fmt.Sprintf("%d", status)).Observe(latency)
		m.ResponseSize.WithLabelValues(c.Request.Method, path).Observe(responseSize)

		// 如果是错误状态码，记录错误
		if status >= 400 {
			m.ErrorsTotal.WithLabelValues("http", path).Inc()
		}
	}
}

// StartMetricsServer 启动监控服务器
func (m *Metrics) StartMetricsServer(config *MetricsConfig) {
	if config == nil {
		config = DefaultMetricsConfig()
	}

	if !config.Enabled {
		utils.Info("监控服务器未启用")
		return
	}

	// 创建新的Gin路由器
	router := gin.New()
	router.Use(gin.Recovery())

	// 注册Prometheus指标端点
	router.GET(config.MetricsPath, gin.WrapH(promhttp.Handler()))

	// 启动HTTP服务器
	go func() {
		utils.Info("监控服务器启动在 %s", config.ListenAddress)
		if err := router.Run(config.ListenAddress); err != nil {
			utils.Error("监控服务器启动失败: %v", err)
		}
	}()

	utils.Info("监控服务器已启动，指标路径: %s%s", config.ListenAddress, config.MetricsPath)
}

// IncrementDatabaseQuery 增加数据库查询计数
func (m *Metrics) IncrementDatabaseQuery(table, operation string) {
	m.DatabaseQueries.WithLabelValues(table, operation).Inc()
}

// IncrementDatabaseError 增加数据库错误计数
func (m *Metrics) IncrementDatabaseError(table, operation, error string) {
	m.DatabaseErrors.WithLabelValues(table, operation, error).Inc()
}

// ObserveDatabaseDuration 记录数据库查询耗时
func (m *Metrics) ObserveDatabaseDuration(table, operation string, duration float64) {
	m.DatabaseDuration.WithLabelValues(table, operation).Observe(duration)
}

// IncrementResourceCreated 增加资源创建计数
func (m *Metrics) IncrementResourceCreated(category, platform string) {
	m.ResourcesCreated.WithLabelValues(category, platform).Inc()
}

// IncrementResourceViewed 增加资源查看计数
func (m *Metrics) IncrementResourceViewed(category string) {
	m.ResourcesViewed.WithLabelValues(category).Inc()
}

// IncrementSearch 增加搜索计数
func (m *Metrics) IncrementSearch(platform string) {
	m.Searches.WithLabelValues(platform).Inc()
}

// IncrementTransfer 增加转存计数
func (m *Metrics) IncrementTransfer(platform, status string) {
	m.Transfers.WithLabelValues(platform, status).Inc()
}

// IncrementError 增加错误计数
func (m *Metrics) IncrementError(errorType, endpoint string) {
	m.ErrorsTotal.WithLabelValues(errorType, endpoint).Inc()
}

// AddCustomCounter 添加自定义计数器
func (m *Metrics) AddCustomCounter(name, help string, labels []string) prometheus.Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s_%v", name, labels)
	if counter, exists := m.CustomCounters[key]; exists {
		return counter
	}

	counter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "urldb",
			Name:      name,
			Help:      help,
		},
		labels,
	).WithLabelValues() // 如果没有标签，返回默认实例

	m.CustomCounters[key] = counter
	return counter
}

// AddCustomGauge 添加自定义仪表盘
func (m *Metrics) AddCustomGauge(name, help string, labels []string) prometheus.Gauge {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s_%v", name, labels)
	if gauge, exists := m.CustomGauges[key]; exists {
		return gauge
	}

	gauge := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "urldb",
			Name:      name,
			Help:      help,
		},
		labels,
	).WithLabelValues() // 如果没有标签，返回默认实例

	m.CustomGauges[key] = gauge
	return gauge
}

// GetMetricsSummary 获取指标摘要
func (m *Metrics) GetMetricsSummary() map[string]interface{} {
	// 这里可以实现获取当前指标摘要的逻辑
	// 由于Prometheus指标不能直接读取，我们只能返回一些基本的统计信息
	return map[string]interface{}{
		"timestamp": time.Now(),
		"status":    "running",
		"info":      "使用 /metrics 端点获取详细指标",
	}
}

// HealthCheck 健康检查
func (m *Metrics) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

// SetupHealthCheck 设置健康检查端点
func (m *Metrics) SetupHealthCheck(router *gin.Engine) {
	router.GET("/health", m.HealthCheck)
	router.GET("/healthz", m.HealthCheck)
}

// MetricsHandler 指标处理器
func (m *Metrics) MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}
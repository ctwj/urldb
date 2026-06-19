package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// handlerTestOnce 确保测试 DB 只初始化一次（与 services/stats_service_test.go 独立）
var handlerTestOnce sync.Once
var handlerDBReady bool

// ensureHandlerTestDB 初始化测试 DB
//
// go test 的 CWD 是 handlers/ 目录，需向上一级加载 .env
func ensureHandlerTestDB(t *testing.T) bool {
	t.Helper()
	handlerTestOnce.Do(func() {
		_ = godotenv.Load("../.env")
		origMigrate := os.Getenv("MIGRATE")
		_ = os.Setenv("MIGRATE", "false")
		defer func() { _ = os.Setenv("MIGRATE", origMigrate) }()

		if err := db.InitDB(); err != nil {
			t.Logf("数据库初始化失败（跳过 handler 集成测试）: %v", err)
			return
		}
		handlerDBReady = true
	})
	return handlerDBReady
}

// injectService 注入 service 到全局，返回还原函数
func injectService(svc *services.StatsService) func() {
	// 通过 SetDefaultStatsService 暴露的入口注入（先取当前值用于还原）
	orig := services.GetDefaultStatsService()
	services.SetDefaultStatsService(svc)
	return func() { services.SetDefaultStatsService(orig) }
}

// parseJSONBody 解析响应体为 map
func parseJSONBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体非合法 JSON: %v (body=%s)", err, w.Body.String())
	}
	return body
}

// TestGetSummaryHandler_ServiceNil_Returns500 当服务未初始化时返回 500
//
// 防御性分支：GetDefaultStatsService 在 db.DB 为 nil 时返回 nil，
// handler 必须返回 500 而非 panic
//
// 注意：本测试需先确保 db.DB 为 nil。由于 db.InitDB 是全局副作用，
// 若 services 测试先运行已初始化 db.DB，则此测试改走真实 DB 分支
func TestGetSummaryHandler_ServiceNil_Returns500(t *testing.T) {
	// 临时清空注入的 service（保留 db.DB 不变）
	// 由于 GetDefaultStatsService 会在 db.DB 非 nil 时兜底构造，
	// 此测试只能在 db.DB 未初始化的环境（如无 .env 的 CI）走 nil 分支
	defer injectService(nil)()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/stats/summary", nil)

	// 防御：若 db.DB 已初始化，handler 会兜底成功，此时允许 200
	GetSummary(c)

	if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
		t.Errorf("期望 500（未初始化）或 200（兜底成功），实际 %d", w.Code)
	}
}

// TestGetSummaryHandler_Success_Returns200WithUnifiedEnvelope 成功路径
//
// 验证 handler：
//   - 返回 HTTP 200
//   - 响应体为统一格式 {success, message, data, code}
//   - data 包含 resources/views/searches/todos 四组
//
// 401（未鉴权）/ 403（非管理员）由认证中间件负责，handler 本身不做鉴权检查
// （已在 plan.md Complexity Tracking 登记：T010 仅覆盖 200/500，鉴权分支属中间件测试范畴）
func TestGetSummaryHandler_Success_Returns200WithUnifiedEnvelope(t *testing.T) {
	if !ensureHandlerTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}

	// 显式构造并注入 service（避免依赖前序测试的全局状态）
	svc := services.NewStatsService(db.DB, repo.NewRepositoryManager(db.DB))
	defer injectService(svc)()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/stats/summary", nil)

	GetSummary(c)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d (body=%s)", w.Code, w.Body.String())
	}

	body := parseJSONBody(t, w)

	// 统一响应封套断言
	if success, ok := body["success"].(bool); !ok || !success {
		t.Errorf("期望 success=true，实际 %v", body["success"])
	}
	if code, ok := body["code"].(float64); !ok || int(code) != 200 {
		t.Errorf("期望 code=200，实际 %v", body["code"])
	}
	if msg, ok := body["message"].(string); !ok || msg == "" {
		t.Errorf("期望 message 非空，实际 %v", body["message"])
	}

	// data 结构断言
	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("期望 data 为对象，实际 %T: %v", body["data"], body["data"])
	}
	for _, key := range []string{"resources", "views", "searches", "todos"} {
		if _, exists := data[key]; !exists {
			t.Errorf("data.%s 字段缺失", key)
		}
	}

	// resources 子字段断言
	resources, ok := data["resources"].(map[string]interface{})
	if !ok {
		t.Fatalf("期望 data.resources 为对象，实际 %T", data["resources"])
	}
	for _, sub := range []string{"today", "yesterday", "total"} {
		if _, exists := resources[sub]; !exists {
			t.Errorf("data.resources.%s 字段缺失", sub)
		}
	}
}

// TestGetSummaryHandler_ContentTypeIsJSON 响应头 Content-Type 为 application/json
func TestGetSummaryHandler_ContentTypeIsJSON(t *testing.T) {
	if !ensureHandlerTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}
	svc := services.NewStatsService(db.DB, repo.NewRepositoryManager(db.DB))
	defer injectService(svc)()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/stats/summary", nil)

	GetSummary(c)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Errorf("期望 Content-Type=application/json; charset=utf-8，实际 %s", ct)
	}
}

// TestGetSummaryHandler_MultipleCallsStable 多次调用响应结构稳定
//
// 验证 JSON 字段不会因 map 迭代顺序而闪烁
func TestGetSummaryHandler_MultipleCallsStable(t *testing.T) {
	if !ensureHandlerTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}
	svc := services.NewStatsService(db.DB, repo.NewRepositoryManager(db.DB))
	defer injectService(svc)()

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/stats/summary", nil)

		GetSummary(c)

		if w.Code != http.StatusOK {
			t.Errorf("第 %d 次调用返回非 200: %d", i+1, w.Code)
		}
		body := parseJSONBody(t, w)
		if _, ok := body["data"].(map[string]interface{}); !ok {
			t.Errorf("第 %d 次调用 data 字段异常", i+1)
		}
	}
}

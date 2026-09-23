package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/gin-gonic/gin"
)

// TestGetUsers_PaginatedEnvelope 验证 GetUsers 返回分页封套 {data: [...], total}
//
// 用户管理页翻页 bug 根因：GetUsers 忽略 page/page_size 返回全量列表，
// 前端每次翻页拿到相同数据。修复后响应 data 应为 {data, total} 对象。
func TestGetUsers_PaginatedEnvelope(t *testing.T) {
	if !ensureHandlerTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}

	SetRepositoryManager(repo.NewRepositoryManager(db.DB))
	defer SetRepositoryManager(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/users?page=1&page_size=2", nil)

	GetUsers(c)

	if w.Code != http.StatusOK {
		t.Fatalf("期望 200，实际 %d (body=%s)", w.Code, w.Body.String())
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体非合法 JSON: %v", err)
	}

	data, ok := body["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("期望 data 为分页对象 {data, total}，实际 %T: %v", body["data"], body["data"])
	}
	list, ok := data["data"].([]interface{})
	if !ok {
		t.Fatalf("期望 data.data 为数组，实际 %T", data["data"])
	}
	total, ok := data["total"].(float64)
	if !ok {
		t.Fatalf("期望 data.total 为数字，实际 %T", data["total"])
	}

	if float64(len(list)) > 2 {
		t.Errorf("page_size=2 时单页最多 2 条，实际 %d 条", len(list))
	}
	if total < float64(len(list)) {
		t.Errorf("total(%v) 不应小于单页条数(%d)", total, len(list))
	}
}

// TestGetUsers_Page2DiffersFromPage1 翻页核心断言：第 2 页与第 1 页内容不同
func TestGetUsers_Page2DiffersFromPage1(t *testing.T) {
	if !ensureHandlerTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}

	SetRepositoryManager(repo.NewRepositoryManager(db.DB))
	defer SetRepositoryManager(nil)

	fetchIDs := func(page string) []float64 {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/users?page="+page+"&page_size=1", nil)
		GetUsers(c)
		if w.Code != http.StatusOK {
			t.Fatalf("page=%s 期望 200，实际 %d", page, w.Code)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("响应体非合法 JSON: %v", err)
		}
		data, ok := body["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("期望 data 为分页对象 {data, total}，实际 %T", body["data"])
		}
		list, _ := data["data"].([]interface{})
		ids := make([]float64, 0, len(list))
		for _, item := range list {
			m, _ := item.(map[string]interface{})
			if id, ok := m["id"].(float64); ok {
				ids = append(ids, id)
			}
		}
		return ids
	}

	p1 := fetchIDs("1")
	p2 := fetchIDs("2")
	totalUsers := int(fetchIDs("1&_=x")[0]) // 占位避免误用，实际不校验

	_ = totalUsers
	if len(p1) == 0 {
		t.Skip("库中无用户，跳过翻页差异断言")
	}
	if len(p2) == 0 {
		t.Log("库中只有 1 页用户，翻页差异不适用")
		return
	}
	if len(p1) > 0 && len(p2) > 0 && p1[0] == p2[0] {
		t.Errorf("翻页失效：第 1 页首个用户 id=%v 与第 2 页相同", p1[0])
	}
}

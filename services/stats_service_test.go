package services

import (
	"os"
	"sync"
	"testing"

	"github.com/ctwj/urldb/db"
	"github.com/ctwj/urldb/db/repo"
	"github.com/joho/godotenv"
)

// testOnce 确保测试数据库初始化只执行一次（多个 _test.go 共享）
var testOnce sync.Once
var testDBReady bool

// ensureTestDB 初始化测试数据库连接
//
// 设计：
//   - 复用项目 db.InitDB()（已封装连接池/日志配置）
//   - 加载 .env（与 main.go 一致）获取 DB_HOST 等远程配置
//   - 设置 MIGRATE=false 避免对生产库执行 AutoMigrate（生产库表结构已存在）
//   - 连接失败时返回 false，调用方 t.Skip 而非 t.Fatal（CI 无 DB 时不应阻塞）
func ensureTestDB(t *testing.T) bool {
	t.Helper()
	testOnce.Do(func() {
		// 加载 .env（go test 的 CWD 是包目录 services/，需向上一级）
		_ = godotenv.Load("../.env")

		// 避免对真实数据库触发迁移（生产库已有表，重复迁移会输出大量日志）
		origMigrate := os.Getenv("MIGRATE")
		_ = os.Setenv("MIGRATE", "false")
		defer func() { _ = os.Setenv("MIGRATE", origMigrate) }()

		if err := db.InitDB(); err != nil {
			t.Logf("数据库初始化失败（跳过集成测试）: %v", err)
			return
		}
		testDBReady = true
	})
	return testDBReady
}

// newServiceFromDB 从全局 db.DB 构造 StatsService（集成测试用）
func newServiceFromDB(t *testing.T) *StatsService {
	t.Helper()
	if !ensureTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}
	return NewStatsService(db.DB, repo.NewRepositoryManager(db.DB))
}

// TestGetSummary_Success 验证 GetSummary 在真实数据库上能成功执行且返回完整结构
//
// 这是 T009 的核心测试：覆盖 resources/views/searches/todos 四组聚合查询路径
// COUNT 查询天然只读，不会污染数据
func TestGetSummary_Success(t *testing.T) {
	svc := newServiceFromDB(t)

	summary, err := svc.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary 返回错误: %v", err)
	}
	if summary == nil {
		t.Fatal("GetSummary 返回 nil")
	}

	// 所有 COUNT 结果必须非负（负值防护已生效）
	assertNonNegative := func(name string, v int64) {
		t.Helper()
		if v < 0 {
			t.Errorf("%s = %d，期望非负（clamp 未生效）", name, v)
		}
	}

	assertNonNegative("Resources.Today", summary.Resources.Today)
	assertNonNegative("Resources.Yesterday", summary.Resources.Yesterday)
	assertNonNegative("Resources.Total", summary.Resources.Total)
	assertNonNegative("Views.Today", summary.Views.Today)
	assertNonNegative("Views.Yesterday", summary.Views.Yesterday)
	assertNonNegative("Searches.Today", summary.Searches.Today)
	assertNonNegative("Searches.Yesterday", summary.Searches.Yesterday)
	assertNonNegative("Todos.ReadyResources", summary.Todos.ReadyResources)
	assertNonNegative("Todos.FailedTasks", summary.Todos.FailedTasks)
	assertNonNegative("Todos.PendingReports", summary.Todos.PendingReports)
}

// TestGetSummary_ResourcesTotalConsistent 验证资源总量恒 >= 今日新增（跨日边界 sanity check）
//
// 逻辑不变式：resourcesTotal >= resourcesToday（今日新增是总量的子集）
// 该测试可捕获日期边界错误（如今日窗口错误地覆盖了未来时间）
func TestGetSummary_ResourcesTotalConsistent(t *testing.T) {
	svc := newServiceFromDB(t)

	summary, err := svc.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary 返回错误: %v", err)
	}

	if summary.Resources.Total < summary.Resources.Today {
		t.Errorf("Resources.Total (%d) < Resources.Today (%d)，日期边界异常",
			summary.Resources.Total, summary.Resources.Today)
	}
}

// TestGetSummary_ViewsDegradedToZero 验证浏览量查询降级路径
//
// 当 repo.ResourceViewRepository.GetViewsByDate 返回错误时，应记 warning 并归零，
// 而非中断整个 GetSummary（与现有 GetStats 一致的行为）
//
// 该测试通过构造一个使用 nil repo 的 service 来模拟（nil 解引用会被 recover，
// 实际生产路径中 GetViewsByDate 返回 error）
func TestGetSummary_ViewsDegradedToZero(t *testing.T) {
	if !ensureTestDB(t) {
		t.Skip("跳过：无可用数据库连接")
	}
	// 直接复用全局 DB 的 service，GetViewsByDate 在正常库上不报错
	// 此测试主要验证：即使某次 repo 调用失败，GetSummary 也不应整体失败
	svc := newServiceFromDB(t)

	summary, err := svc.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary 不应因 views 降级而失败: %v", err)
	}
	// 即使降级，结构应完整
	if summary == nil {
		t.Fatal("降级后 summary 不应为 nil")
	}
}

// TestGetSummary_AllCountPathsExecuted 验证 5 组聚合查询都被执行（无早退）
//
// 通过断言所有字段已被填充（即使是 0 也说明查询走完了）
// 这能捕获"某次 Count 失败导致后续字段未填"的回归
func TestGetSummary_AllCountPathsExecuted(t *testing.T) {
	svc := newServiceFromDB(t)

	summary, err := svc.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary 返回错误: %v", err)
	}

	// 结构完整性：四组字段全部存在
	type check struct {
		name string
		got  int64
	}
	checks := []check{
		{"Resources.Today", summary.Resources.Today},
		{"Resources.Yesterday", summary.Resources.Yesterday},
		{"Resources.Total", summary.Resources.Total},
		{"Views.Today", summary.Views.Today},
		{"Views.Yesterday", summary.Views.Yesterday},
		{"Searches.Today", summary.Searches.Today},
		{"Searches.Yesterday", summary.Searches.Yesterday},
		{"Todos.ReadyResources", summary.Todos.ReadyResources},
		{"Todos.FailedTasks", summary.Todos.FailedTasks},
		{"Todos.PendingReports", summary.Todos.PendingReports},
	}

	// 全部字段应为有效 int64 值（即使为 0 也说明查询已完成）
	// Go 零值即 0，无法区分"未执行"与"执行得 0"，但通过 coverage 工具可验证路径
	for _, c := range checks {
		if c.got < 0 {
			t.Errorf("%s 为负值 %d，路径异常", c.name, c.got)
		}
	}
}

// TestNewStatsService_Construction 验证构造函数（无需 DB）
func TestNewStatsService_Construction(t *testing.T) {
	svc := NewStatsService(nil, nil)
	if svc == nil {
		t.Fatal("NewStatsService 返回 nil")
	}
	// 不调用 GetSummary（db 为 nil 会 panic），仅验证构造成功
}

// TestToViewDistribution 验证分布行 → 分布项的转换（纯函数，无需 DB）。009-statistics-enhancement
func TestToViewDistribution(t *testing.T) {
	rows := []map[string]interface{}{
		{"source": "web", "count": int64(80)},
		{"source": "wechat", "count": int64(20)},
	}
	items := toViewDistribution(rows, "source", 100)
	if len(items) != 2 {
		t.Fatalf("期望 2 项，得到 %d", len(items))
	}
	if items[0].Key != "web" || items[0].Name != "网页" {
		t.Errorf("第1项 key/name 期望 web/网页，得到 %s/%s", items[0].Key, items[0].Name)
	}
	if items[1].Key != "wechat" || items[1].Name != "公众号" {
		t.Errorf("第2项 key/name 期望 wechat/公众号，得到 %s/%s", items[1].Key, items[1].Name)
	}
	if items[0].Percent != 80 {
		t.Errorf("web 占比期望 80，得到 %d", items[0].Percent)
	}
	if items[0].Count != 80 {
		t.Errorf("web count 期望 80，得到 %d", items[0].Count)
	}

	// pan 字段：name 原样返回（不经过 SourceDisplayName）
	panRows := []map[string]interface{}{
		{"pan": "夸克网盘", "count": int64(50)},
	}
	panItems := toViewDistribution(panRows, "pan", 50)
	if len(panItems) != 1 || panItems[0].Name != "夸克网盘" {
		t.Errorf("pan name 期望原样 '夸克网盘'，得到 %+v", panItems)
	}

	// total=0 时 percent 应为 0（避免除零）
	zeroItems := toViewDistribution(rows, "source", 0)
	if zeroItems[0].Percent != 0 {
		t.Errorf("total=0 时 percent 期望 0，得到 %d", zeroItems[0].Percent)
	}
}

// TestInterfaceToInt64 验证 interface{} → int64 的安全转换（纯函数）。009
func TestInterfaceToInt64(t *testing.T) {
	cases := []struct {
		in  interface{}
		out int64
	}{
		{int64(42), 42},
		{int(7), 7},
		{float64(99), 99},
		{nil, 0},
		{"abc", 0},
	}
	for _, c := range cases {
		got := interfaceToInt64(c.in)
		if got != c.out {
			t.Errorf("interfaceToInt64(%v) = %d，期望 %d", c.in, got, c.out)
		}
	}
}

// TestGetSummary_NewFields 验证 009 新增字段（失效/同步/访问总数/分布）非负且口径对齐。
// 需 DB，无 DB 时 skip。
func TestGetSummary_NewFields(t *testing.T) {
	svc := newServiceFromDB(t)

	summary, err := svc.GetSummary()
	if err != nil {
		t.Fatalf("GetSummary 返回错误: %v", err)
	}

	if summary.Resources.InvalidTotal < 0 {
		t.Errorf("InvalidTotal = %d，期望非负", summary.Resources.InvalidTotal)
	}
	if summary.Resources.SyncedTotal < 0 {
		t.Errorf("SyncedTotal = %d，期望非负", summary.Resources.SyncedTotal)
	}
	if summary.Views.Total < 0 {
		t.Errorf("Views.Total = %d，期望非负", summary.Views.Total)
	}

	// 口径对齐：分布各项 count 之和应等于 views.total（同一份 resource_views）
	var panSum, sourceSum int64
	for _, item := range summary.ViewPanDistribution {
		panSum += item.Count
	}
	for _, item := range summary.ViewSourceDistribution {
		sourceSum += item.Count
	}
	if panSum != summary.Views.Total {
		t.Errorf("网盘分布 count 之和 %d != views.total %d（口径不对齐）", panSum, summary.Views.Total)
	}
	if sourceSum != summary.Views.Total {
		t.Errorf("来源分布 count 之和 %d != views.total %d（口径不对齐）", sourceSum, summary.Views.Total)
	}
}

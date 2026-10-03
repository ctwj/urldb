package db

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/ctwj/urldb/db/entity"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ============================================================================
// 种子逻辑测试基建（specs/005-fix-admin-login T002）
//
// 安全红线：本测试严禁触碰 ../.env 默认库（当前指向真实远程库）。
// 连接参数可取自 ../.env（host/port/user/password），但库名强制覆盖为
// 独立测试库（默认 urldb_seed_test，可用 SEED_TEST_DB_NAME 覆盖），
// 连接失败按仓库惯例 t.Skip。
// ============================================================================

const seedTestDBDefaultName = "urldb_seed_test"

var (
	seedTestOnce  sync.Once
	seedTestDB    *gorm.DB
	seedTestReady bool
)

// seedTestDBName 测试库名（SEED_TEST_DB_NAME 可覆盖，便于 CI 隔离）
func seedTestDBName() string {
	if name := os.Getenv("SEED_TEST_DB_NAME"); name != "" {
		return name
	}
	return seedTestDBDefaultName
}

// initSeedTestDB 惰性建立测试库连接：必要时建库，AutoMigrate 种子涉及的四张表
func initSeedTestDB() error {
	// go test 的 CWD 是 db/，向上一级加载 .env（仅取连接参数，库名随后覆盖）
	_ = godotenv.Load("../.env")

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	name := seedTestDBName()

	// 连接维护库，确保测试库存在（CREATE DATABASE 无 IF NOT EXISTS，容忍已存在报错）
	admin, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)))
	if err != nil {
		return fmt.Errorf("连接维护库失败: %w", err)
	}
	if err := admin.Exec(fmt.Sprintf(`CREATE DATABASE %s`, name)).Error; err != nil {
		// 42P07 = duplicate_database，视为成功
		if !isDuplicateDatabaseErr(err) {
			return fmt.Errorf("创建测试库 %s 失败: %w", name, err)
		}
	}

	gdb, err := gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)))
	if err != nil {
		return fmt.Errorf("连接测试库 %s 失败: %w", name, err)
	}

	if err := gdb.AutoMigrate(&entity.User{}, &entity.Pan{}, &entity.Category{}, &entity.SystemConfig{}); err != nil {
		return fmt.Errorf("测试库 AutoMigrate 失败: %w", err)
	}

	seedTestDB = gdb
	seedTestReady = true
	return nil
}

func isDuplicateDatabaseErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "already exists") || strings.Contains(msg, "42P07")
}

// setupSeedTest 每个用例的入口：清空四张种子表后，将全局 DB 指向测试库并返回
func setupSeedTest(t *testing.T) *gorm.DB {
	t.Helper()

	seedTestOnce.Do(func() {
		if err := initSeedTestDB(); err != nil {
			t.Logf("种子测试库不可用（跳过种子测试）: %v", err)
			return
		}
	})
	if !seedTestReady {
		t.Skip("跳过：无可用种子测试库")
	}

	if err := seedTestDB.Exec(`TRUNCATE TABLE users, pan, categories, system_configs RESTART IDENTITY CASCADE`).Error; err != nil {
		t.Fatalf("清空种子表失败: %v", err)
	}

	// 种子函数使用包级 DB，测试期间指向测试库，用例结束还原
	orig := DB
	DB = seedTestDB
	t.Cleanup(func() { DB = orig })

	return seedTestDB
}

// assertAdminPasswordValid 断言给定哈希与明文密码匹配
func assertAdminPasswordValid(t *testing.T, hash, plaintext string) {
	t.Helper()
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext)); err != nil {
		t.Errorf("密码哈希与明文 %q 不匹配: %v", plaintext, err)
	}
}

// countRows 统计表行数
func countRows(t *testing.T, gdb *gorm.DB, model interface{}) int64 {
	t.Helper()
	var n int64
	if err := gdb.Model(model).Count(&n).Error; err != nil {
		t.Fatalf("统计行数失败: %v", err)
	}
	return n
}

// ============================================================================
// US1：空库首启全量补种（quickstart 场景 1）
//
// 回归背景：ensureDefaultPans 先行填充 pan 表，insertDefaultDataIfEmpty 以
// pan 行数判定"是否空库"而误判跳过，全新安装 admin/password 无法登录。
// ============================================================================

func TestEnsureDefault_FreshInstall(t *testing.T) {
	gdb := setupSeedTest(t)

	seedDefaultData()

	// 管理员：唯一、管理员角色、激活、密码与明文 password 匹配
	var admins []entity.User
	if err := gdb.Where("username = ?", "admin").Find(&admins).Error; err != nil {
		t.Fatalf("查询 admin 失败: %v", err)
	}
	if len(admins) != 1 {
		t.Fatalf("admin 用户数 = %d, 期望 1", len(admins))
	}
	a := admins[0]
	if a.Role != "admin" || !a.IsActive {
		t.Errorf("admin 角色或激活状态错误: role=%q is_active=%v", a.Role, a.IsActive)
	}
	assertAdminPasswordValid(t, a.Password, "password")

	// 默认分类 13 项、平台 10 项、系统配置非空
	if n := countRows(t, gdb, &entity.Category{}); n != 13 {
		t.Errorf("默认分类行数 = %d, 期望 13", n)
	}
	if n := countRows(t, gdb, &entity.Pan{}); n != 10 {
		t.Errorf("默认平台行数 = %d, 期望 10", n)
	}
	if n := countRows(t, gdb, &entity.SystemConfig{}); n == 0 {
		t.Error("默认系统配置为空")
	}
}

// ============================================================================
// US2：受影响存量库自愈 + 已有数据零覆盖（quickstart 场景 4 等效）
// 复现线上库状态：pan 已有数据、users 为空 —— 升级修复版后重启应补建 admin，
// 且既有数据（自定义配置值/自定义分类/平台）原样保留。
// ============================================================================

func TestEnsureDefault_SelfHeal(t *testing.T) {
	gdb := setupSeedTest(t)

	// 复现"半初始化"状态：平台已补种、配置/分类已有自定义数据、管理员缺失
	if err := ensureDefaultPans(); err != nil {
		t.Fatalf("预置平台数据失败: %v", err)
	}
	customConfig := entity.SystemConfig{Key: entity.ConfigKeySiteTitle, Value: "我的自定义站名", Type: entity.ConfigTypeString}
	if err := gdb.Create(&customConfig).Error; err != nil {
		t.Fatalf("预置自定义配置失败: %v", err)
	}
	customCategory := entity.Category{Name: "自定义分类", Description: "用户自建"}
	if err := gdb.Create(&customCategory).Error; err != nil {
		t.Fatalf("预置自定义分类失败: %v", err)
	}

	seedDefaultData()

	// 管理员被补建
	var admin entity.User
	if err := gdb.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatalf("存量库自愈后应补建 admin: %v", err)
	}
	assertAdminPasswordValid(t, admin.Password, "password")

	// 既有数据零覆盖：自定义配置值保留、自定义分类保留、平台不重复
	var cfg entity.SystemConfig
	if err := gdb.Where("key = ?", entity.ConfigKeySiteTitle).First(&cfg).Error; err != nil {
		t.Fatalf("查询站点标题配置失败: %v", err)
	}
	if cfg.Value != "我的自定义站名" {
		t.Errorf("自定义配置值被覆盖: %q", cfg.Value)
	}
	var catCount int64
	gdb.Model(&entity.Category{}).Where("name = ?", "自定义分类").Count(&catCount)
	if catCount != 1 {
		t.Errorf("自定义分类应保留且不重复，实际 %d 行", catCount)
	}
	if n := countRows(t, gdb, &entity.Pan{}); n != 10 {
		t.Errorf("平台行数 = %d, 期望保持 10", n)
	}
	// 默认键共 33 个（见 ensureDefaultSystemConfigs 清单），其中 site_title 已预置
	// 不重复插入 → 预置 1 + 补齐 32 = 33 行。若默认清单增删，此处同步调整。
	if n := countRows(t, gdb, &entity.SystemConfig{}); n != 33 {
		t.Errorf("配置行数 = %d, 期望 33（预置 1 + 补齐缺失的 32 个默认键）", n)
	}
}

// TestEnsureDefault_NoReset 已有管理员（含改密）绝不重置（data-model PRESENT_KEPT）
func TestEnsureDefault_NoReset(t *testing.T) {
	gdb := setupSeedTest(t)

	// 预置改过密码的管理员与一个普通用户
	changedHash, err := bcrypt.GenerateFromPassword([]byte("changedpass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("生成改密哈希失败: %v", err)
	}
	existingAdmin := entity.User{
		Username: "boss", Password: string(changedHash),
		Email: "boss@example.com", Role: "admin", IsActive: true,
	}
	if err := gdb.Create(&existingAdmin).Error; err != nil {
		t.Fatalf("预置管理员失败: %v", err)
	}
	normalUser := entity.User{
		Username: "user1", Password: string(changedHash),
		Email: "u1@example.com", Role: "user", IsActive: true,
	}
	if err := gdb.Create(&normalUser).Error; err != nil {
		t.Fatalf("预置普通用户失败: %v", err)
	}

	seedDefaultData()

	// 用户行数不变：不植入默认 admin、不产生重复
	if n := countRows(t, gdb, &entity.User{}); n != 2 {
		t.Errorf("用户行数 = %d, 期望 2（不应植入默认 admin）", n)
	}
	// 改过的密码保持有效（未被重置为 password）
	var boss entity.User
	if err := gdb.Where("username = ?", "boss").First(&boss).Error; err != nil {
		t.Fatalf("查询既有管理员失败: %v", err)
	}
	assertAdminPasswordValid(t, boss.Password, "changedpass")
	if err := bcrypt.CompareHashAndPassword([]byte(boss.Password), []byte("password")); err == nil {
		t.Error("既有管理员密码被重置为默认值")
	}
}

// ============================================================================
// US3：幂等与软删除冲突（quickstart 场景 2/3，SC-003）
// ============================================================================

// TestEnsureDefault_Idempotent 连续多轮种子结果一致，改密后不再被重置
func TestEnsureDefault_Idempotent(t *testing.T) {
	gdb := setupSeedTest(t)

	seedDefaultData()
	// 修改管理员密码后再次种子（模拟"重启"）
	var admin entity.User
	if err := gdb.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatalf("首轮种子后应存在 admin: %v", err)
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte("restarted-pass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("生成新密码哈希失败: %v", err)
	}
	if err := gdb.Model(&admin).Update("password", string(newHash)).Error; err != nil {
		t.Fatalf("修改管理员密码失败: %v", err)
	}

	seedDefaultData()
	seedDefaultData()

	if n := countRows(t, gdb, &entity.User{}); n != 1 {
		t.Errorf("三轮种子后用户行数 = %d, 期望 1（无重复）", n)
	}
	if err := gdb.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatalf("查询 admin 失败: %v", err)
	}
	assertAdminPasswordValid(t, admin.Password, "restarted-pass")
	if n := countRows(t, gdb, &entity.Category{}); n != 13 {
		t.Errorf("分类行数 = %d, 期望 13（无重复）", n)
	}
	if n := countRows(t, gdb, &entity.Pan{}); n != 10 {
		t.Errorf("平台行数 = %d, 期望 10（无重复）", n)
	}
}

// TestEnsureDefault_SoftDeletedAdmin 软删除 admin 占用用户名：
// 显著告警、跳过补建、不中断（返回 nil）、不复活、不新增（BLOCKED_SOFT_DELETED）
func TestEnsureDefault_SoftDeletedAdmin(t *testing.T) {
	gdb := setupSeedTest(t)

	// 预置一个被软删除的 admin
	seedAdmin := entity.User{
		Username: "admin", Password: "x", Email: "admin@example.com",
		Role: "admin", IsActive: true,
	}
	if err := gdb.Create(&seedAdmin).Error; err != nil {
		t.Fatalf("预置 admin 失败: %v", err)
	}
	if err := gdb.Delete(&seedAdmin).Error; err != nil {
		t.Fatalf("软删除 admin 失败: %v", err)
	}

	// 当前不存在任何未软删除的管理员 → 会尝试补建，但用户名被软删除行占用
	if err := ensureDefaultAdmin(false); err != nil {
		t.Fatalf("软删除冲突时应跳过补建并正常返回（不中断启动），实际返回错误: %v", err)
	}

	// 无新增行、未复活（Unscoped 计数含软删除行，恒为 1）
	var unscopedCount int64
	if err := gdb.Unscoped().Model(&entity.User{}).Where("username = ?", "admin").Count(&unscopedCount).Error; err != nil {
		t.Fatalf("Unscoped 统计失败: %v", err)
	}
	if unscopedCount != 1 {
		t.Errorf("username=admin 的 Unscoped 行数 = %d, 期望 1（不新增、不复活）", unscopedCount)
	}
	// 活跃管理员仍不存在（软删除行不可登录）
	var activeCount int64
	gdb.Model(&entity.User{}).Where("role = ?", "admin").Count(&activeCount)
	if activeCount != 0 {
		t.Errorf("活跃管理员行数 = %d, 期望 0（软删除行不应被复活为可登录）", activeCount)
	}
}

// TestEnsureDefault_SingleFailureIsolated 单项补种失败只影响该实体，
// 不阻断其余默认数据与管理员补建（FR-006，data-model 不变量 3）
func TestEnsureDefault_SingleFailureIsolated(t *testing.T) {
	gdb := setupSeedTest(t)

	// 模拟分类表缺失（单项失败源：13 条 FirstOrCreate 全部报错）
	if err := gdb.Migrator().DropTable(&entity.Category{}); err != nil {
		t.Fatalf("删除分类表失败: %v", err)
	}

	seedDefaultData()

	// 管理员、配置、平台均应照常补齐
	var admin entity.User
	if err := gdb.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatalf("分类失败不应阻断管理员补建: %v", err)
	}
	if n := countRows(t, gdb, &entity.SystemConfig{}); n == 0 {
		t.Error("分类失败不应阻断系统配置补种")
	}
	if n := countRows(t, gdb, &entity.Pan{}); n != 10 {
		t.Errorf("平台行数 = %d, 期望 10", n)
	}
}

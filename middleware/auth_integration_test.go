package middleware

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ============================================================================
// JWT密钥初始化集成测试
//
// 安全红线：与 db 包种子测试一致，严禁触碰 ../.env 默认库（可能指向真实库）。
// 连接参数取自 ../.env，库名强制为独立测试库 urldb_jwt_test，连接失败 t.Skip。
// ============================================================================

const jwtTestDBName = "urldb_jwt_test"

var (
	jwtTestOnce  sync.Once
	jwtTestDB    *gorm.DB
	jwtTestReady bool
)

func initJWTTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	jwtTestOnce.Do(func() {
		// go test 的 CWD 是 middleware/，向上一级加载 .env（仅取连接参数，库名随后覆盖）
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

		// 连接维护库，确保测试库存在（CREATE DATABASE 无 IF NOT EXISTS，容忍已存在报错）
		admin, err := gorm.Open(postgres.Open(fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)))
		if err != nil {
			return
		}
		if err := admin.Exec(fmt.Sprintf("CREATE DATABASE %s", jwtTestDBName)).Error; err != nil {
			// 库已存在等情况，继续尝试连接
		}

		dbConn, err := gorm.Open(postgres.Open(fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, jwtTestDBName)))
		if err != nil {
			return
		}
		if err := dbConn.AutoMigrate(&entity.SystemConfig{}); err != nil {
			return
		}

		jwtTestDB = dbConn
		jwtTestReady = true
	})

	if !jwtTestReady {
		t.Skip("测试数据库不可用，跳过JWT密钥集成测试")
	}
	return jwtTestDB
}

// TestInitJWTSecretIntegration 验证完整链路：随机生成→持久化→签发/校验闭环
func TestInitJWTSecretIntegration(t *testing.T) {
	dbConn := initJWTTestDB(t)
	t.Setenv("JWT_SECRET", "") // 确保走"无环境变量"路径

	testRepo := repo.NewRepositoryManager(dbConn)
	SetRepositoryManager(testRepo)

	if err := InitJWTSecret(); err != nil {
		t.Fatalf("InitJWTSecret 失败: %v", err)
	}

	// 密钥已持久化到系统配置
	stored, err := testRepo.SystemConfigRepository.GetConfigValue(entity.ConfigKeyJWTSecret)
	if err != nil {
		t.Fatalf("持久化的密钥应可读回: %v", err)
	}
	if len(stored) != 64 {
		t.Errorf("持久化的密钥应为64字符hex，得到 %d 字符", len(stored))
	}

	// 内存中的密钥与持久化值一致
	if string(jwtSecret) != stored {
		t.Error("内存密钥应与持久化值一致")
	}

	// 签发/校验闭环：动态密钥签出的令牌可被校验
	user := &entity.User{ID: 42, Username: "tester", Role: "admin"}
	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}
	claims, err := parseToken(token)
	if err != nil {
		t.Fatalf("parseToken 校验失败: %v", err)
	}
	if claims.UserID != 42 || claims.Username != "tester" || claims.Role != "admin" {
		t.Errorf("claims 不匹配: %+v", claims)
	}
}

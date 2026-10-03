package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
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

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "url_db"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	// 配置慢查询日志
	slowThreshold := getEnvInt("DB_SLOW_THRESHOLD_MS", 200)
	logLevel := logger.Info
	if os.Getenv("ENV") == "production" {
		logLevel = logger.Warn
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Duration(slowThreshold) * time.Millisecond,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		return err
	}

	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 优化数据库连接池参数
	maxOpenConns := getEnvInt("DB_MAX_OPEN_CONNS", 50)
	maxIdleConns := getEnvInt("DB_MAX_IDLE_CONNS", 20)
	connMaxLifetime := getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)

	sqlDB.SetMaxOpenConns(maxOpenConns)                                    // 最大打开连接数
	sqlDB.SetMaxIdleConns(maxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute) // 连接最大生命周期

	utils.Info("数据库连接池配置 - 最大连接: %d, 空闲连接: %d, 生命周期: %d分钟",
		maxOpenConns, maxIdleConns, connMaxLifetime)

	// 检查是否需要迁移（只在开发环境或首次启动时）
	if shouldRunMigration() {
		utils.Info("开始数据库迁移...")
		err = DB.AutoMigrate(
			&entity.User{},
			&entity.Category{},
			&entity.Pan{},
			&entity.Cks{},
			&entity.Tag{},
			&entity.Resource{},
			&entity.ResourceTag{},
			&entity.ReadyResource{},
			&entity.SearchStat{},
			&entity.SystemConfig{},
			&entity.HotDrama{},
			&entity.ResourceView{},
			&entity.Task{},
			&entity.TaskItem{},
			&entity.File{},
			&entity.TelegramChannel{},
			&entity.APIAccessLog{},
			&entity.APIAccessLogStats{},
			&entity.APIAccessLogSummary{},
			&entity.Report{},
			&entity.CopyrightClaim{},
			&entity.DownloadHistory{},
			&entity.UserResource{},
			&entity.ApiApplication{},
			&entity.ApiCredential{},
			// 插件系统相关表
			&entity.PluginConfig{},
			&entity.PluginLog{},
			&entity.CustomEvent{},
			&entity.CronJob{},
			&entity.UserPreference{},
			&entity.UserInterest{},
			&entity.URLStats{},
			&entity.AccessStats{},
			&entity.PopularResources{},
			&entity.DomainPattern{},
			&entity.PathPattern{},
			&entity.RealTimeMetrics{},
			&entity.ClassificationStats{},
			&entity.DailyReport{},
			&entity.SystemHealth{},
		)
		if err != nil {
			utils.Fatal("数据库迁移失败: %v", err)
		}
		utils.Info("数据库迁移完成")
	} else {
		utils.Info("跳过数据库迁移（表结构已是最新）")
	}

	// 创建索引以提高查询性能（只在需要迁移时）
	if shouldRunMigration() {
		createIndexes(DB)
	}

	// 009-statistics-enhancement: 历史来源数据回填（幂等，仅在迁移时执行）
	if shouldRunMigration() {
		backfillSourceColumn(DB)
	}

	// 幂等补种全部默认数据（specs/005-fix-admin-login）
	// 各 ensure 互相独立，顺序仅用于日志可读性；单项失败只影响该实体，不阻断启动
	seedDefaultData()

	utils.Info("数据库连接成功")
	return nil
}

// shouldRunMigration 检查是否需要运行数据库迁移
func shouldRunMigration() bool {
	// 优先检查新的 MIGRATE 配置
	migrate := os.Getenv("MIGRATE")
	if migrate == "false" {
		utils.Info("MIGRATE=false，跳过数据库迁移")
		return false
	}

	// 如果明确设置 MIGRATE=true，则执行迁移
	if migrate == "true" {
		utils.Info("MIGRATE=true，执行数据库迁移")
		return true
	}

	// 兼容旧的 SKIP_MIGRATION 配置
	skipMigration := os.Getenv("SKIP_MIGRATION")
	if skipMigration == "true" {
		utils.Info("SKIP_MIGRATION=true，跳过数据库迁移")
		return false
	}

	// 检查环境变量
	env := os.Getenv("ENV")
	if env == "production" {
		// 生产环境：检查是否有迁移标记
		var count int64
		DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'schema_migrations'").Count(&count)
		if count == 0 {
			// 没有迁移表，说明是首次部署
			utils.Info("生产环境首次部署，执行数据库迁移")
			return true
		}
		// 有迁移表，检查是否需要迁移（这里可以添加更复杂的逻辑）
		utils.Info("生产环境已有迁移表，跳过数据库迁移")
		return false
	}

	// 开发环境：总是运行迁移
	utils.Info("开发环境，执行数据库迁移")
	return true
}

// autoMigrate 自动迁移表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&entity.SystemConfig{}, // 系统配置表（独立表，先创建）
		&entity.Pan{},
		&entity.Cks{},
		&entity.Category{},
		&entity.Tag{},
		&entity.Resource{},
		&entity.ResourceTag{},
		&entity.ReadyResource{},
		&entity.User{},
		&entity.ApiApplication{},
		&entity.ApiCredential{},
		&entity.SearchStat{},
		&entity.HotDrama{},
		&entity.File{},
		&entity.TelegramChannel{},
		// 插件系统相关表
		&entity.PluginConfig{},
		&entity.PluginLog{},
		&entity.CustomEvent{},
		&entity.CronJob{},
		&entity.UserPreference{},
		&entity.UserInterest{},
		&entity.URLStats{},
		&entity.AccessStats{},
		&entity.PopularResources{},
		&entity.DomainPattern{},
		&entity.PathPattern{},
		&entity.RealTimeMetrics{},
		&entity.ClassificationStats{},
		&entity.DailyReport{},
		&entity.SystemHealth{},
	)
}

// createIndexes 创建数据库索引以提高查询性能
func createIndexes(db *gorm.DB) {
	// 资源表索引（移除全文搜索索引，使用Meilisearch替代）
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_category_id ON resources(category_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_pan_id ON resources(pan_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_created_at ON resources(created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_updated_at ON resources(updated_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_is_valid ON resources(is_valid)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_is_public ON resources(is_public)")

	// 为Meilisearch准备的基础文本索引（用于精确匹配）
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_title ON resources(title)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resources_description ON resources(description)")

	// 待处理资源表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ready_resource_key ON ready_resource(key)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ready_resource_url ON ready_resource(url)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_ready_resource_create_time ON ready_resource(create_time DESC)")

	// 搜索统计表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_search_stats_keyword ON search_stats(keyword)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_search_stats_created_at ON search_stats(created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_search_stats_source ON search_stats(source)") // 009 来源分布

	// 资源访问记录表索引（009-statistics-enhancement：获取资源来源/网盘分布）
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_views_resource_id ON resource_views(resource_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_views_source ON resource_views(source)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_views_created_at ON resource_views(created_at DESC)")

	// 热播剧表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_hot_dramas_title ON hot_dramas(title)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_hot_dramas_category ON hot_dramas(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_hot_dramas_created_at ON hot_dramas(created_at DESC)")

	// 资源标签关联表索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_tags_resource_id ON resource_tags(resource_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_resource_tags_tag_id ON resource_tags(tag_id)")

	// API访问日志表索引 - 高性能查询优化
	db.Exec("CREATE INDEX IF NOT EXISTS idx_api_access_logs_created_at ON api_access_logs(created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_api_access_logs_endpoint_status ON api_access_logs(endpoint, response_status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_api_access_logs_ip_created ON api_access_logs(ip, created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_api_access_logs_method_endpoint ON api_access_logs(method, endpoint)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_api_access_logs_response_time ON api_access_logs(processing_time)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_api_access_logs_error_logs ON api_access_logs(response_status, created_at DESC) WHERE response_status >= 400")

	// 任务和任务项表索引 - Google索引功能优化
	db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_type_status ON tasks(type, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC)")

	// task_items表的关键索引 - 支持高效去重和状态查询
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_items_url ON task_items(url)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_items_url_status ON task_items(url, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_items_task_id ON task_items(task_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_items_status ON task_items(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_items_created_at ON task_items(created_at DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_task_items_status_created ON task_items(status, created_at)")

	// 创建插件系统相关索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_plugin_configs_plugin_name ON plugin_configs(plugin_name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_plugin_configs_enabled ON plugin_configs(enabled)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_plugin_logs_plugin_name ON plugin_logs(plugin_name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_plugin_logs_hook_name ON plugin_logs(hook_name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_plugin_logs_success ON plugin_logs(success)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_custom_events_event_name ON custom_events(event_name)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_custom_events_processed ON custom_events(processed)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_cron_jobs_enabled ON cron_jobs(enabled)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_cron_jobs_next_run ON cron_jobs(next_run)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_preferences_category ON user_preferences(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_interests_user_id ON user_interests(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_user_interests_category ON user_interests(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_url_stats_domain ON url_stats(domain)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_url_stats_category ON url_stats(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_access_stats_url_id ON access_stats(url_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_access_stats_user_id ON access_stats(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_access_stats_access_time ON access_stats(access_time)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_popular_resources_access_count ON popular_resources(access_count DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_domain_pattern_count ON domain_patterns(count DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_path_pattern_count ON path_patterns(count DESC)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_real_time_metrics_metric ON real_time_metrics(metric)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_real_time_metrics_timestamp ON real_time_metrics(timestamp)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_classification_stats_category ON classification_stats(category)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_daily_reports_type ON daily_reports(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_system_health_status ON system_health(status)")

	utils.Info("数据库索引创建完成（已移除全文搜索索引，准备使用Meilisearch，新增API访问日志性能索引，任务项表索引优化，插件系统索引）")
}

// backfillSourceColumn 009-statistics-enhancement: 将历史搜索/获取记录的来源回填为 web。
// search_stats / resource_views 新增 source 列带 default:'web'，但 AutoMigrate 不会回填历史已存在行，
// 故显式将 NULL/空值记录置为 web（代码可确定历史记录均来自网页前端）。幂等，重复执行无害。
func backfillSourceColumn(db *gorm.DB) {
	if res := db.Exec("UPDATE search_stats SET source = ? WHERE source IS NULL OR source = ''", entity.SourceWeb); res.Error != nil {
		utils.Error("回填 search_stats.source 失败: %v", res.Error)
	} else {
		utils.Info("回填 search_stats.source=web，影响行数: %d", res.RowsAffected)
	}
	if res := db.Exec("UPDATE resource_views SET source = ? WHERE source IS NULL OR source = ''", entity.SourceWeb); res.Error != nil {
		utils.Error("回填 resource_views.source 失败: %v", res.Error)
	} else {
		utils.Info("回填 resource_views.source=web，影响行数: %d", res.RowsAffected)
	}
}

// seedDefaultData 幂等补种全部默认数据（specs/005-fix-admin-login）
//
// 历史缺陷：曾以 pan 表行数判定"是否空库"，而 ensureDefaultPans 先于该判定执行
// 并填充了 pan 表，导致全新安装被误判为老库、默认管理员永不创建（admin/password
// 无法登录）。现改为各实体独立幂等 ensure：顺序无关、缺失即补、已有不覆盖。
// 各 ensure 单项失败只影响该实体并记录错误，不阻断启动。
func seedDefaultData() {
	// D3：必须在任何补种发生前快照"库中是否已有默认数据"，用于区分全新首装与
	// 存量库自愈——若在补种后检测，pan/分类已被本轮填充，两者将无法区分。
	selfHeal := hasAnyDefaultData()
	if err := ensureDefaultPans(); err != nil {
		utils.Error("补种默认平台失败: %v", err)
	}
	if err := ensureDefaultCategories(); err != nil {
		utils.Error("补种默认分类失败: %v", err)
	}
	if err := ensureDefaultSystemConfigs(); err != nil {
		utils.Error("补种默认系统配置失败: %v", err)
	}
	if err := ensureDefaultAdmin(selfHeal); err != nil {
		utils.Error("补种默认管理员失败: %v", err)
	}
}

// hasAnyDefaultData 库中是否已存在任一默认数据（平台/分类/系统配置）
func hasAnyDefaultData() bool {
	var n int64
	for _, model := range []interface{}{&entity.Pan{}, &entity.Category{}, &entity.SystemConfig{}} {
		if err := DB.Model(model).Count(&n).Error; err == nil && n > 0 {
			return true
		}
	}
	return false
}

// ensureDefaultAdmin 幂等保障默认管理员存在（specs/005-fix-admin-login，research.md D1）
//
// 仅当不存在任何未软删除的 admin 角色用户时，补建默认管理员
// （username=admin / 初始密码 password）；已存在管理员（含改名/改密）则静默
// 跳过，绝不重置密码。判定语义是"系统是否仍可管理"，而非"用户名 admin 是否
// 存在"——避免向已另建管理员的部署静默植入公开默认凭据账号。
// selfHeal 由调用方在补种开始前快照传入（见 seedDefaultData），决定补建日志级别。
func ensureDefaultAdmin(selfHeal bool) error {
	var adminCount int64
	if err := DB.Model(&entity.User{}).Where("role = ?", "admin").Count(&adminCount).Error; err != nil {
		return fmt.Errorf("检查管理员是否存在失败: %w", err)
	}
	if adminCount > 0 {
		return nil
	}

	// D2（specs/005-fix-admin-login research.md）：username 唯一索引在数据库级生效且
	// 包含软删除行。补建前先 Unscoped 预检用户名 admin——命中软删除行时，是否复活
	// 属于人工处置决定：输出显著告警、跳过补建、不中断启动（不硬删旧行、不静默复活）。
	// 命中活跃行则必为非管理员角色占用用户名（管理员角色已在上方静默返回），同样告警。
	var existing entity.User
	if err := DB.Unscoped().Where("username = ?", "admin").First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			utils.Error("默认管理员补建受阻：用户名 admin 已被软删除记录占用（数据库唯一索引包含软删除行）。" +
				"处置建议：恢复该账号（清除 deleted_at）或将其改名后重启；本次启动跳过补建并继续运行")
		} else {
			utils.Error("默认管理员补建受阻：用户名 admin 已被非管理员角色的用户占用。" +
				"处置建议：将该账号升级为管理员或改名后重启；本次启动跳过补建并继续运行")
		}
		return nil
	}

	// D3（specs/005-fix-admin-login research.md）：存量库补建公开默认凭据账号属安全
	// 相关事件，selfHeal=true 时以 WARN 级提示运维知晓并尽快改密（区别于首装 INFO）。
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("生成默认管理员密码哈希失败: %w", err)
	}

	defaultAdmin := entity.User{
		Username: "admin",
		Password: string(hash),
		Email:    "admin@example.com",
		Role:     "admin",
		IsActive: true,
	}
	if err := DB.Create(&defaultAdmin).Error; err != nil {
		return fmt.Errorf("补建默认管理员失败: %w", err)
	}

	if selfHeal {
		utils.Warn("检测到存量数据但管理员缺失（历史缺陷的半初始化状态），已补建默认管理员 admin（初始密码 password）——默认凭据为公开信息，请立即登录并修改密码")
	} else {
		utils.Info("已创建默认管理员 admin（初始密码 password，请尽快修改）")
	}
	return nil
}

// ensureDefaultPans 幂等补种默认网盘平台
// 老库升级场景：新版本新增的默认平台（如 guangya）不会写入已有数据库，导致前端
// "添加账号"下拉缺少该平台。这里每次启动按 name 逐个 FirstOrCreate：
// 已有的（含管理员修改过的）不动，缺失的补齐。
func ensureDefaultPans() error {
	defaultPans := []entity.Pan{
		{Name: "baidu", Key: 1, Icon: "<i class=\"fas fa-cloud text-blue-500\"></i>", Remark: "百度网盘"},
		{Name: "aliyun", Key: 2, Icon: "<i class=\"fas fa-cloud text-orange-500\"></i>", Remark: "阿里云盘"},
		{Name: "quark", Key: 3, Icon: "<i class=\"fas fa-atom text-purple-500\"></i>", Remark: "夸克网盘"},
		{Name: "tianyi", Key: 4, Icon: "<i class=\"fas fa-cloud text-cyan-500\"></i>", Remark: "天翼云盘"},
		{Name: "xunlei", Key: 5, Icon: "<i class=\"fas fa-bolt text-yellow-500\"></i>", Remark: "迅雷云盘"},
		{Name: "123pan", Key: 8, Icon: "<i class=\"fas fa-folder text-red-500\"></i>", Remark: "123云盘"},
		{Name: "115", Key: 12, Icon: "<i class=\"fas fa-cloud-upload-alt text-green-600\"></i>", Remark: "115网盘"},
		{Name: "uc", Key: 14, Icon: "<i class=\"fas fa-cloud-download-alt text-purple-600\"></i>", Remark: "UC网盘"},
		{Name: "guangya", Key: 16, Icon: "<i class=\"fas fa-dove text-pink-500\"></i>", Remark: "光鸭云盘"},
		{Name: "other", Key: 15, Icon: "<i class=\"fas fa-cloud text-gray-500\"></i>", Remark: "其他"},
	}

	for _, pan := range defaultPans {
		if err := DB.Where("name = ?", pan.Name).FirstOrCreate(&pan).Error; err != nil {
			utils.Error("插入平台 %s 失败: %v", pan.Name, err)
			// 继续执行，不因为单个平台失败而停止
		}
	}
	return nil
}

// ensureDefaultCategories 幂等补种默认分类（specs/005-fix-admin-login）
// 已存在的分类（含管理员修改过的）不动，缺失的补齐；单项失败仅记录错误不阻断。
func ensureDefaultCategories() error {
	defaultCategories := []entity.Category{
		{Name: "电影", Description: "电影"},
		{Name: "电视剧", Description: "电视剧"},
		{Name: "短剧", Description: "短剧"},
		{Name: "综艺", Description: "综艺"},
		{Name: "动漫", Description: "动漫"},
		{Name: "纪录片", Description: "纪录片"},
		{Name: "视频教程", Description: "视频教程"},
		{Name: "学习资料", Description: "学习资料"},
		{Name: "游戏", Description: "其他游戏资源"},
		{Name: "软件", Description: "软件"},
		{Name: "APP", Description: "APP"},
		{Name: "AI", Description: "AI"},
		{Name: "其他", Description: "其他资源"},
	}

	for _, category := range defaultCategories {
		if err := DB.Where("name = ?", category.Name).FirstOrCreate(&category).Error; err != nil {
			utils.Error("插入分类 %s 失败: %v", category.Name, err)
			// 继续执行，不因为单个分类失败而停止
		}
	}
	return nil
}

// ensureDefaultSystemConfigs 幂等补种默认系统配置（specs/005-fix-admin-login）
// 已存在的配置（含管理员修改过的值）不动，缺失的键补默认值；单项失败仅记录错误不阻断。
// 与运行期 SystemConfigRepository.GetOrCreateDefault 的按需补种同源同口径。
func ensureDefaultSystemConfigs() error {
	defaultSystemConfigs := []entity.SystemConfig{
		{Key: entity.ConfigKeySiteTitle, Value: entity.ConfigDefaultSiteTitle, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeySiteDescription, Value: entity.ConfigDefaultSiteDescription, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyKeywords, Value: entity.ConfigDefaultKeywords, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyAuthor, Value: entity.ConfigDefaultAuthor, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyCopyright, Value: entity.ConfigDefaultCopyright, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyAutoProcessReadyResources, Value: entity.ConfigDefaultAutoProcessReadyResources, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyAutoProcessInterval, Value: entity.ConfigDefaultAutoProcessInterval, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyAutoTransferEnabled, Value: entity.ConfigDefaultAutoTransferEnabled, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyAutoTransferLimitDays, Value: entity.ConfigDefaultAutoTransferLimitDays, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyAutoTransferMinSpace, Value: entity.ConfigDefaultAutoTransferMinSpace, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyAutoFetchHotDramaEnabled, Value: entity.ConfigDefaultAutoFetchHotDramaEnabled, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyApiToken, Value: entity.ConfigDefaultApiToken, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyForbiddenWords, Value: entity.ConfigDefaultForbiddenWords, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyAdKeywords, Value: entity.ConfigDefaultAdKeywords, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyAutoInsertAd, Value: entity.ConfigDefaultAutoInsertAd, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyPageSize, Value: entity.ConfigDefaultPageSize, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyMaintenanceMode, Value: entity.ConfigDefaultMaintenanceMode, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyEnableRegister, Value: entity.ConfigDefaultEnableRegister, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyThirdPartyStatsCode, Value: entity.ConfigDefaultThirdPartyStatsCode, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyMeilisearchEnabled, Value: entity.ConfigDefaultMeilisearchEnabled, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyMeilisearchHost, Value: entity.ConfigDefaultMeilisearchHost, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyMeilisearchPort, Value: entity.ConfigDefaultMeilisearchPort, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyMeilisearchMasterKey, Value: entity.ConfigDefaultMeilisearchMasterKey, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyMeilisearchIndexName, Value: entity.ConfigDefaultMeilisearchIndexName, Type: entity.ConfigTypeString},
		// PanCheck 链接检测服务默认配置
		{Key: entity.ConfigKeyPanCheckEnabled, Value: entity.ConfigDefaultPanCheckEnabled, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyPanCheckHost, Value: entity.ConfigDefaultPanCheckHost, Type: entity.ConfigTypeString},
		{Key: entity.ConfigKeyPanCheckTimeoutSeconds, Value: entity.ConfigDefaultPanCheckTimeoutSeconds, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyPanCheckBatchSize, Value: entity.ConfigDefaultPanCheckBatchSize, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyPanCheckConcurrency, Value: entity.ConfigDefaultPanCheckConcurrency, Type: entity.ConfigTypeInt},
		// 自动清理转存文件默认配置（002-auto-cleanup-transfer）
		{Key: entity.ConfigKeyAutoCleanupEnabled, Value: entity.ConfigDefaultAutoCleanupEnabled, Type: entity.ConfigTypeBool},
		{Key: entity.ConfigKeyAutoCleanupRetentionDays, Value: entity.ConfigDefaultAutoCleanupRetentionDays, Type: entity.ConfigTypeInt},
		{Key: entity.ConfigKeyAutoCleanupIntervalMinutes, Value: entity.ConfigDefaultAutoCleanupIntervalMinutes, Type: entity.ConfigTypeInt},
		// 用户上传资源默认配置（015-user-resource-upload）：每用户每日提交上限，0 表示不限制
		{Key: entity.ConfigKeyUserUploadDailyLimit, Value: entity.ConfigDefaultUserUploadDailyLimit, Type: entity.ConfigTypeInt},
	}

	for _, config := range defaultSystemConfigs {
		if err := DB.Where("key = ?", config.Key).FirstOrCreate(&config).Error; err != nil {
			utils.Error("插入系统配置 %s 失败: %v", config.Key, err)
			// 继续执行，不因为单个配置失败而停止
		}
	}
	return nil
}

// getEnvInt 获取环境变量中的整数值，如果不存在则返回默认值
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		utils.Warn("环境变量 %s 的值 '%s' 不是有效的整数，使用默认值 %d", key, value, defaultValue)
		return defaultValue
	}

	return intValue
}

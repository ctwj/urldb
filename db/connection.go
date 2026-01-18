package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/utils"
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

	sqlDB.SetMaxOpenConns(maxOpenConns)                               // 最大打开连接数
	sqlDB.SetMaxIdleConns(maxIdleConns)                               // 最大空闲连接数
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
			// AI 提示词管理表
			&entity.AIPrompt{},
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

	// 插入默认数据（只在数据库为空时）
	if err := insertDefaultDataIfEmpty(); err != nil {
		utils.Error("插入默认数据失败: %v", err)
	}

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

// insertDefaultDataIfEmpty 只在数据库为空时插入默认数据
func insertDefaultDataIfEmpty() error {
	// 检查是否已有数据
	var panCount int64
	if err := DB.Model(&entity.Pan{}).Count(&panCount).Error; err != nil {
		return err
	}

	// 如果pan表已有数据，跳过插入
	if panCount > 0 {
		utils.Info("pan表已有数据，跳过默认数据插入")
		return nil
	}

	utils.Info("pan表为空，开始插入默认数据...")

	// 插入默认分类（使用FirstOrCreate避免重复）
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

	// 插入默认网盘平台（使用FirstOrCreate避免重复）
	defaultPans := []entity.Pan{
		{Name: "baidu", Key: 1, Icon: "<i class=\"fas fa-cloud text-blue-500\"></i>", Remark: "百度网盘"},
		{Name: "aliyun", Key: 2, Icon: "<i class=\"fas fa-cloud text-orange-500\"></i>", Remark: "阿里云盘"},
		{Name: "quark", Key: 3, Icon: "<i class=\"fas fa-atom text-purple-500\"></i>", Remark: "夸克网盘"},
		{Name: "tianyi", Key: 4, Icon: "<i class=\"fas fa-cloud text-cyan-500\"></i>", Remark: "天翼云盘"},
		{Name: "xunlei", Key: 5, Icon: "<i class=\"fas fa-bolt text-yellow-500\"></i>", Remark: "迅雷云盘"},
		{Name: "123pan", Key: 8, Icon: "<i class=\"fas fa-folder text-red-500\"></i>", Remark: "123云盘"},
		{Name: "115", Key: 12, Icon: "<i class=\"fas fa-cloud-upload-alt text-green-600\"></i>", Remark: "115网盘"},
		{Name: "uc", Key: 14, Icon: "<i class=\"fas fa-cloud-download-alt text-purple-600\"></i>", Remark: "UC网盘"},
		{Name: "other", Key: 15, Icon: "<i class=\"fas fa-cloud text-gray-500\"></i>", Remark: "其他"},
	}

	for _, pan := range defaultPans {
		if err := DB.Where("name = ?", pan.Name).FirstOrCreate(&pan).Error; err != nil {
			utils.Error("插入平台 %s 失败: %v", pan.Name, err)
			// 继续执行，不因为单个平台失败而停止
		}
	}

	// 插入默认系统配置
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
	}

	for _, config := range defaultSystemConfigs {
		if err := DB.Where("key = ?", config.Key).FirstOrCreate(&config).Error; err != nil {
			utils.Error("插入系统配置 %s 失败: %v", config.Key, err)
			// 继续执行，不因为单个配置失败而停止
		}
	}

	// 插入默认管理员用户
	defaultAdmin := entity.User{
		Username: "admin",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Email:    "admin@example.com",
		Role:     "admin",
		IsActive: true,
	}

	if err := DB.Create(&defaultAdmin).Error; err != nil {
		return err
	}

	utils.Info("默认数据插入完成")
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

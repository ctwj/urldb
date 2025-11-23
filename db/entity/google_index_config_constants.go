package entity

// GoogleIndexConfigKeys Google索引配置键常量
const (
	// 通用配置
	GoogleIndexConfigKeyEnabled            = "google_index_enabled"             // 是否启用Google索引功能
	GoogleIndexConfigKeySiteURL            = "google_index_site_url"            // 网站URL
	GoogleIndexConfigKeySiteName           = "google_index_site_name"           // 网站名称

	// 认证配置
	GoogleIndexConfigKeyCredentialsFile    = "google_index_credentials_file"    // 凭证文件路径
	GoogleIndexConfigKeyClientEmail        = "google_index_client_email"        // 客户端邮箱
	GoogleIndexConfigKeyClientID           = "google_index_client_id"           // 客户端ID
	GoogleIndexConfigKeyPrivateKey         = "google_index_private_key"         // 私钥（加密存储）
	GoogleIndexConfigKeyToken              = "google_index_token"               // 访问令牌

	// 调度配置
	GoogleIndexConfigKeyCheckInterval     = "google_index_check_interval"      // 检查间隔（分钟）
	GoogleIndexConfigKeyBatchSize         = "google_index_batch_size"          // 批处理大小
	GoogleIndexConfigKeyConcurrency       = "google_index_concurrency"         // 并发数
	GoogleIndexConfigKeyRetryAttempts     = "google_index_retry_attempts"      // 重试次数
	GoogleIndexConfigKeyRetryDelay        = "google_index_retry_delay"         // 重试延迟（秒）

	// 网站地图配置
	GoogleIndexConfigKeyAutoSitemap       = "google_index_auto_sitemap"        // 自动提交网站地图
	GoogleIndexConfigKeySitemapPath       = "google_index_sitemap_path"        // 网站地图路径
	GoogleIndexConfigKeySitemapSchedule   = "google_index_sitemap_schedule"    // 网站地图调度
)
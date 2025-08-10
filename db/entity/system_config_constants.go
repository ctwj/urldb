package entity

// ConfigKey 配置键常量
const (
	// SEO 配置
	ConfigKeySiteTitle       = "site_title"
	ConfigKeySiteDescription = "site_description"
	ConfigKeyKeywords        = "keywords"
	ConfigKeyAuthor          = "author"
	ConfigKeyCopyright       = "copyright"

	// 自动处理配置组
	ConfigKeyAutoProcessReadyResources = "auto_process_ready_resources"
	ConfigKeyAutoProcessInterval       = "auto_process_interval"
	ConfigKeyAutoTransferEnabled       = "auto_transfer_enabled"
	ConfigKeyAutoTransferLimitDays     = "auto_transfer_limit_days"
	ConfigKeyAutoTransferMinSpace      = "auto_transfer_min_space"
	ConfigKeyAutoFetchHotDramaEnabled  = "auto_fetch_hot_drama_enabled"

	// API配置
	ConfigKeyApiToken = "api_token"

	// 违禁词配置
	ConfigKeyForbiddenWords = "forbidden_words"

	// 其他配置
	ConfigKeyPageSize        = "page_size"
	ConfigKeyMaintenanceMode = "maintenance_mode"
	ConfigKeyEnableRegister  = "enable_register"

	// 三方统计配置
	ConfigKeyThirdPartyStatsCode = "third_party_stats_code"
)

// ConfigType 配置类型常量
const (
	ConfigTypeString = "string"
	ConfigTypeInt    = "int"
	ConfigTypeBool   = "bool"
	ConfigTypeJSON   = "json"
)

// ConfigResponseField API响应字段名常量
const (
	// 基础字段
	ConfigResponseFieldID        = "id"
	ConfigResponseFieldCreatedAt = "created_at"
	ConfigResponseFieldUpdatedAt = "updated_at"

	// SEO 配置字段
	ConfigResponseFieldSiteTitle       = "site_title"
	ConfigResponseFieldSiteDescription = "site_description"
	ConfigResponseFieldKeywords        = "keywords"
	ConfigResponseFieldAuthor          = "author"
	ConfigResponseFieldCopyright       = "copyright"

	// 自动处理配置字段
	ConfigResponseFieldAutoProcessReadyResources = "auto_process_ready_resources"
	ConfigResponseFieldAutoProcessInterval       = "auto_process_interval"
	ConfigResponseFieldAutoTransferEnabled       = "auto_transfer_enabled"
	ConfigResponseFieldAutoTransferLimitDays     = "auto_transfer_limit_days"
	ConfigResponseFieldAutoTransferMinSpace      = "auto_transfer_min_space"
	ConfigResponseFieldAutoFetchHotDramaEnabled  = "auto_fetch_hot_drama_enabled"

	// API配置字段
	ConfigResponseFieldApiToken = "api_token"

	// 违禁词配置字段
	ConfigResponseFieldForbiddenWords = "forbidden_words"

	// 其他配置字段
	ConfigResponseFieldPageSize        = "page_size"
	ConfigResponseFieldMaintenanceMode = "maintenance_mode"
	ConfigResponseFieldEnableRegister  = "enable_register"

	// 三方统计配置字段
	ConfigResponseFieldThirdPartyStatsCode = "third_party_stats_code"
)

// ConfigDefaultValue 配置默认值常量
const (
	// SEO 配置默认值
	ConfigDefaultSiteTitle       = "老九网盘资源数据库"
	ConfigDefaultSiteDescription = "专业的老九网盘资源数据库"
	ConfigDefaultKeywords        = "网盘,资源管理,文件分享"
	ConfigDefaultAuthor          = "系统管理员"
	ConfigDefaultCopyright       = "© 2024 老九网盘资源数据库"

	// 自动处理配置默认值
	ConfigDefaultAutoProcessReadyResources = "false"
	ConfigDefaultAutoProcessInterval       = "30"
	ConfigDefaultAutoTransferEnabled       = "false"
	ConfigDefaultAutoTransferLimitDays     = "0"
	ConfigDefaultAutoTransferMinSpace      = "100"
	ConfigDefaultAutoFetchHotDramaEnabled  = "false"

	// API配置默认值
	ConfigDefaultApiToken = ""

	// 违禁词配置默认值
	ConfigDefaultForbiddenWords = ""

	// 其他配置默认值
	ConfigDefaultPageSize        = "100"
	ConfigDefaultMaintenanceMode = "false"
	ConfigDefaultEnableRegister  = "true"

	// 三方统计配置默认值
	ConfigDefaultThirdPartyStatsCode = ""
)

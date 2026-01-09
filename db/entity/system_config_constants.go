package entity

// ConfigKey 配置键常量
const (
	// SEO 配置
	ConfigKeySiteTitle       = "site_title"
	ConfigKeySiteDescription = "site_description"
	ConfigKeyKeywords        = "keywords"
	ConfigKeyAuthor          = "author"
	ConfigKeyCopyright       = "copyright"
	ConfigKeySiteLogo        = "site_logo"

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

	// 广告配置
	ConfigKeyAdKeywords   = "ad_keywords"    // 广告关键词
	ConfigKeyAutoInsertAd = "auto_insert_ad" // 自动插入广告

	// 其他配置
	ConfigKeyPageSize        = "page_size"
	ConfigKeyMaintenanceMode = "maintenance_mode"
	ConfigKeyEnableRegister  = "enable_register"

	// 三方统计配置
	ConfigKeyThirdPartyStatsCode = "third_party_stats_code"

	// Meilisearch配置
	ConfigKeyMeilisearchEnabled   = "meilisearch_enabled"
	ConfigKeyMeilisearchHost      = "meilisearch_host"
	ConfigKeyMeilisearchPort      = "meilisearch_port"
	ConfigKeyMeilisearchMasterKey = "meilisearch_master_key"
	ConfigKeyMeilisearchIndexName = "meilisearch_index_name"

	// Telegram配置
	ConfigKeyTelegramBotEnabled         = "telegram_bot_enabled"
	ConfigKeyTelegramBotApiKey          = "telegram_bot_api_key"
	ConfigKeyTelegramAutoReplyEnabled   = "telegram_auto_reply_enabled"
	ConfigKeyTelegramAutoReplyTemplate  = "telegram_auto_reply_template"
	ConfigKeyTelegramAutoDeleteEnabled  = "telegram_auto_delete_enabled"
	ConfigKeyTelegramAutoDeleteInterval = "telegram_auto_delete_interval"
	ConfigKeyTelegramProxyEnabled       = "telegram_proxy_enabled"
	ConfigKeyTelegramProxyType          = "telegram_proxy_type"
	ConfigKeyTelegramProxyHost          = "telegram_proxy_host"
	ConfigKeyTelegramProxyPort          = "telegram_proxy_port"
	ConfigKeyTelegramProxyUsername      = "telegram_proxy_username"
	ConfigKeyTelegramProxyPassword      = "telegram_proxy_password"

	// 微信公众号配置
	ConfigKeyWechatBotEnabled         = "wechat_bot_enabled"
	ConfigKeyWechatAppId              = "wechat_app_id"
	ConfigKeyWechatAppSecret          = "wechat_app_secret"
	ConfigKeyWechatToken              = "wechat_token"
	ConfigKeyWechatEncodingAesKey     = "wechat_encoding_aes_key"
	ConfigKeyWechatWelcomeMessage     = "wechat_welcome_message"
	ConfigKeyWechatAutoReplyEnabled   = "wechat_auto_reply_enabled"
	ConfigKeyWechatSearchLimit        = "wechat_search_limit"

	// 界面配置
	ConfigKeyEnableAnnouncements = "enable_announcements"
	ConfigKeyAnnouncements       = "announcements"
	ConfigKeyEnableFloatButtons  = "enable_float_buttons"
	ConfigKeyWechatSearchImage   = "wechat_search_image"
	ConfigKeyTelegramQrImage     = "telegram_qr_image"
	ConfigKeyQrCodeStyle         = "qr_code_style"

	// Sitemap配置
	ConfigKeySitemapConfig              = "sitemap_config"
	ConfigKeySitemapLastGenerateTime    = "sitemap_last_generate_time"
	ConfigKeySitemapAutoGenerateEnabled = "sitemap_auto_generate_enabled"

	// 网站URL配置
	ConfigKeyWebsiteURL = "website_url"

	// AI 配置
	ConfigKeyAIAPIKey       = "ai_api_key"
	ConfigKeyAIAPIURL       = "ai_api_url"
	ConfigKeyAIModel        = "ai_model"
	ConfigKeyAIMaxTokens    = "ai_max_tokens"
	ConfigKeyAITemperature  = "ai_temperature"
	ConfigKeyAIOrganization = "ai_organization"
	ConfigKeyAIProxy        = "ai_proxy"
	ConfigKeyAITimeout      = "ai_timeout"
	ConfigKeyAIRetryCount   = "ai_retry_count"

	// MCP 配置
	ConfigKeyMCPEnabled     = "mcp_enabled"
	ConfigKeyMCPConfigPath  = "mcp_config_path"
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

	// 广告配置字段
	ConfigResponseFieldAdKeywords   = "ad_keywords"
	ConfigResponseFieldAutoInsertAd = "auto_insert_ad"

	// 其他配置字段
	ConfigResponseFieldPageSize        = "page_size"
	ConfigResponseFieldMaintenanceMode = "maintenance_mode"
	ConfigResponseFieldEnableRegister  = "enable_register"

	// 三方统计配置字段
	ConfigResponseFieldThirdPartyStatsCode = "third_party_stats_code"

	// Meilisearch配置字段
	ConfigResponseFieldMeilisearchEnabled   = "meilisearch_enabled"
	ConfigResponseFieldMeilisearchHost      = "meilisearch_host"
	ConfigResponseFieldMeilisearchPort      = "meilisearch_port"
	ConfigResponseFieldMeilisearchMasterKey = "meilisearch_master_key"
	ConfigResponseFieldMeilisearchIndexName = "meilisearch_index_name"

	// Telegram配置字段
	ConfigResponseFieldTelegramBotEnabled         = "telegram_bot_enabled"
	ConfigResponseFieldTelegramBotApiKey          = "telegram_bot_api_key"
	ConfigResponseFieldTelegramAutoReplyEnabled   = "telegram_auto_reply_enabled"
	ConfigResponseFieldTelegramAutoReplyTemplate  = "telegram_auto_reply_template"
	ConfigResponseFieldTelegramAutoDeleteEnabled  = "telegram_auto_delete_enabled"
	ConfigResponseFieldTelegramAutoDeleteInterval = "telegram_auto_delete_interval"
	ConfigResponseFieldTelegramProxyEnabled       = "telegram_proxy_enabled"
	ConfigResponseFieldTelegramProxyType          = "telegram_proxy_type"
	ConfigResponseFieldTelegramProxyHost          = "telegram_proxy_host"
	ConfigResponseFieldTelegramProxyPort          = "telegram_proxy_port"
	ConfigResponseFieldTelegramProxyUsername      = "telegram_proxy_username"
	ConfigResponseFieldTelegramProxyPassword      = "telegram_proxy_password"

	// 微信公众号配置字段
	ConfigResponseFieldWechatBotEnabled         = "wechat_bot_enabled"
	ConfigResponseFieldWechatAppId              = "wechat_app_id"
	ConfigResponseFieldWechatAppSecret          = "wechat_app_secret"
	ConfigResponseFieldWechatToken              = "wechat_token"
	ConfigResponseFieldWechatEncodingAesKey     = "wechat_encoding_aes_key"
	ConfigResponseFieldWechatWelcomeMessage     = "wechat_welcome_message"
	ConfigResponseFieldWechatAutoReplyEnabled   = "wechat_auto_reply_enabled"
	ConfigResponseFieldWechatSearchLimit        = "wechat_search_limit"

	// 界面配置字段
	ConfigResponseFieldEnableAnnouncements = "enable_announcements"
	ConfigResponseFieldAnnouncements       = "announcements"
	ConfigResponseFieldEnableFloatButtons  = "enable_float_buttons"
	ConfigResponseFieldWechatSearchImage   = "wechat_search_image"
	ConfigResponseFieldTelegramQrImage     = "telegram_qr_image"
	ConfigResponseFieldQrCodeStyle         = "qr_code_style"

	// 网站URL配置字段
	ConfigResponseFieldWebsiteURL = "site_url"

	// AI 配置字段
	ConfigResponseFieldAIAPIKey       = "ai_api_key"
	ConfigResponseFieldAIAPIURL       = "ai_api_url"
	ConfigResponseFieldAIModel        = "ai_model"
	ConfigResponseFieldAIMaxTokens    = "ai_max_tokens"
	ConfigResponseFieldAITemperature  = "ai_temperature"
	ConfigResponseFieldAIOrganization = "ai_organization"
	ConfigResponseFieldAIProxy        = "ai_proxy"
	ConfigResponseFieldAITimeout      = "ai_timeout"
	ConfigResponseFieldAIRetryCount   = "ai_retry_count"

	// MCP 配置字段
	ConfigResponseFieldMCPEnabled    = "mcp_enabled"
	ConfigResponseFieldMCPConfigPath = "mcp_config_path"
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

	// 广告配置默认值
	ConfigDefaultAdKeywords   = ""
	ConfigDefaultAutoInsertAd = ""

	// 其他配置默认值
	ConfigDefaultPageSize        = "100"
	ConfigDefaultMaintenanceMode = "false"
	ConfigDefaultEnableRegister  = "true"

	// 三方统计配置默认值
	ConfigDefaultThirdPartyStatsCode = ""

	// Meilisearch配置默认值
	ConfigDefaultMeilisearchEnabled   = "false"
	ConfigDefaultMeilisearchHost      = "localhost"
	ConfigDefaultMeilisearchPort      = "7700"
	ConfigDefaultMeilisearchMasterKey = ""
	ConfigDefaultMeilisearchIndexName = "resources"

	// Telegram配置默认值
	ConfigDefaultTelegramBotEnabled         = "false"
	ConfigDefaultTelegramBotApiKey          = ""
	ConfigDefaultTelegramAutoReplyEnabled   = "true"
	ConfigDefaultTelegramAutoReplyTemplate  = "您好！我可以帮您搜索网盘资源，请输入您要搜索的内容。"
	ConfigDefaultTelegramAutoDeleteEnabled  = "false"
	ConfigDefaultTelegramAutoDeleteInterval = "60"
	ConfigDefaultTelegramProxyEnabled       = "false"
	ConfigDefaultTelegramProxyType          = "http"
	ConfigDefaultTelegramProxyHost          = ""
	ConfigDefaultTelegramProxyPort          = "8080"
	ConfigDefaultTelegramProxyUsername      = ""
	ConfigDefaultTelegramProxyPassword      = ""

	// 微信公众号配置默认值
	ConfigDefaultWechatBotEnabled         = "false"
	ConfigDefaultWechatAppId              = ""
	ConfigDefaultWechatAppSecret          = ""
	ConfigDefaultWechatToken              = ""
	ConfigDefaultWechatEncodingAesKey     = ""
	ConfigDefaultWechatWelcomeMessage     = "欢迎关注老九网盘资源库！发送关键词即可搜索资源。"
	ConfigDefaultWechatAutoReplyEnabled   = "true"
	ConfigDefaultWechatSearchLimit        = "5"

	// 界面配置默认值
	ConfigDefaultEnableAnnouncements = "false"
	ConfigDefaultAnnouncements       = ""
	ConfigDefaultEnableFloatButtons  = "false"
	ConfigDefaultWechatSearchImage   = ""
	ConfigDefaultTelegramQrImage     = ""
	ConfigDefaultQrCodeStyle         = "Plain"

	// 网站URL配置默认值
	ConfigDefaultWebsiteURL = "https://example.com"

	// AI 配置默认值
	ConfigDefaultAIAPIKey       = ""
	ConfigDefaultAIAPIURL       = "https://api.openai.com/v1"
	ConfigDefaultAIModel        = "gpt-3.5-turbo"
	ConfigDefaultAIMaxTokens    = "1000"
	ConfigDefaultAITemperature  = "0.7"
	ConfigDefaultAIOrganization = ""
	ConfigDefaultAIProxy        = ""
	ConfigDefaultAITimeout      = "30"
	ConfigDefaultAIRetryCount   = "3"

	// MCP 配置默认值
	ConfigDefaultMCPEnabled    = "false"
	ConfigDefaultMCPConfigPath = "./data/mcp_config.json"
)

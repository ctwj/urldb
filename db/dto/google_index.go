package dto

import (
	"time"
)

// GoogleIndexConfigInput Google索引配置输入
type GoogleIndexConfigInput struct {
	Group   string `json:"group" binding:"required"`
	Key     string `json:"key" binding:"required"`
	Value   string `json:"value" binding:"required"`
	Type    string `json:"type" default:"string"`
}

// GoogleIndexConfigOutput Google索引配置输出
type GoogleIndexConfigOutput struct {
	ID        uint      `json:"id"`
	Group     string    `json:"group"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GoogleIndexConfigGeneral 通用配置
type GoogleIndexConfigGeneral struct {
	Enabled  bool   `json:"enabled" binding:"required"`
	SiteURL  string `json:"site_url" binding:"required"`
	SiteName string `json:"site_name"`
}

// GoogleIndexConfigAuth 认证配置
type GoogleIndexConfigAuth struct {
	CredentialsFile string `json:"credentials_file"`
	ClientEmail     string `json:"client_email"`
	ClientID        string `json:"client_id"`
	PrivateKey      string `json:"private_key"`
	Token           string `json:"token"`
}

// GoogleIndexConfigSchedule 调度配置
type GoogleIndexConfigSchedule struct {
	CheckInterval int `json:"check_interval" binding:"required,min=1,max=1440"` // 检查间隔（分钟），1-24小时
	BatchSize     int `json:"batch_size" binding:"required,min=1,max=1000"`     // 批处理大小
	Concurrency   int `json:"concurrency" binding:"required,min=1,max=10"`      // 并发数
	RetryAttempts int `json:"retry_attempts" binding:"required,min=0,max=10"`    // 重试次数
	RetryDelay    int `json:"retry_delay" binding:"required,min=1,max=60"`       // 重试延迟（秒）
}

// GoogleIndexConfigSitemap 网站地图配置
type GoogleIndexConfigSitemap struct {
	AutoSitemap     bool   `json:"auto_sitemap"`
	SitemapPath     string `json:"sitemap_path" default:"/sitemap.xml"`
	SitemapSchedule string `json:"sitemap_schedule" default:"@daily"` // cron表达式
}

// GoogleIndexTaskInput Google索引任务输入
type GoogleIndexTaskInput struct {
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description"`
	URLs        []string `json:"urls,omitempty"` // 用于URL索引检查任务
	SitemapURL  string `json:"sitemap_url,omitempty"` // 用于网站地图提交任务
	ConfigID    *uint  `json:"config_id,omitempty"`
}

// GoogleIndexTaskOutput Google索引任务输出
type GoogleIndexTaskOutput struct {
	ID                  uint                   `json:"id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	Type                string                 `json:"type"`
	Status              string                 `json:"status"`
	Progress            float64                `json:"progress"`
	TotalItems          int                    `json:"total_items"`
	ProcessedItems      int                    `json:"processed_items"`
	SuccessfulItems     int                    `json:"successful_items"`
	FailedItems         int                    `json:"failed_items"`
	PendingItems        int                    `json:"pending_items"`
	ProcessingItems     int                    `json:"processing_items"`
	IndexedURLs         int                    `json:"indexed_urls"`
	FailedURLs          int                    `json:"failed_urls"`
	ConfigID            *uint                  `json:"config_id"`
	ProgressData        map[string]interface{} `json:"progress_data"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	StartedAt           *time.Time             `json:"started_at"`
	CompletedAt         *time.Time             `json:"completed_at"`
}

// GoogleIndexTaskItemInput Google索引任务项输入
type GoogleIndexTaskItemInput struct {
	TaskID uint `json:"task_id" binding:"required"`
	URL    string `json:"url" binding:"required,url"`
}

// GoogleIndexTaskItemOutput Google索引任务项输出
type GoogleIndexTaskItemOutput struct {
	ID           uint       `json:"id"`
	TaskID       uint       `json:"task_id"`
	URL          string     `json:"url"`
	Status       string     `json:"status"`
	IndexStatus  string     `json:"index_status"`
	ErrorMessage string     `json:"error_message"`
	InspectResult string    `json:"inspect_result"`
	MobileFriendly bool     `json:"mobile_friendly"`
	LastCrawled  *time.Time `json:"last_crawled"`
	StatusCode   int        `json:"status_code"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// GoogleIndexURLStatusInput Google索引URL状态输入
type GoogleIndexURLStatusInput struct {
	URL         string `json:"url" binding:"required,url"`
	IndexStatus string `json:"index_status" binding:"required"`
}

// GoogleIndexURLStatusOutput Google索引URL状态输出
type GoogleIndexURLStatusOutput struct {
	ID             uint      `json:"id"`
	URL            string    `json:"url"`
	IndexStatus    string    `json:"index_status"`
	LastChecked    time.Time `json:"last_checked"`
	CanonicalURL   *string   `json:"canonical_url"`
	LastCrawled    *time.Time `json:"last_crawled"`
	ChangeFreq     *string   `json:"change_freq"`
	Priority       *float64  `json:"priority"`
	MobileFriendly bool      `json:"mobile_friendly"`
	RobotsBlocked  bool      `json:"robots_blocked"`
	LastError      *string   `json:"last_error"`
	StatusCode     int       `json:"status_code"`
	StatusCodeText string    `json:"status_code_text"`
	CheckCount     int       `json:"check_count"`
	SuccessCount   int       `json:"success_count"`
	FailureCount   int       `json:"failure_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GoogleIndexBatchRequest 批量处理请求
type GoogleIndexBatchRequest struct {
	URLs      []string `json:"urls" binding:"required,min=1,max=1000"`
	Operation string   `json:"operation" binding:"required,oneof=check_index submit_sitemap ping"` // 操作类型
}

// GoogleIndexBatchResponse 批量处理响应
type GoogleIndexBatchResponse struct {
	Success bool     `json:"success"`
	Results []string `json:"results,omitempty"`
	Message string   `json:"message,omitempty"`
	Total   int      `json:"total"`
	Processed int    `json:"processed"`
	Failed  int      `json:"failed"`
}

// GoogleIndexStatusResponse 索引状态响应
type GoogleIndexStatusResponse struct {
	Enabled           bool      `json:"enabled"`
	SiteURL           string    `json:"site_url"`
	LastCheckTime     time.Time `json:"last_check_time"`
	TotalURLs         int       `json:"total_urls"`
	IndexedURLs       int       `json:"indexed_urls"`
	NotIndexedURLs    int       `json:"not_indexed_urls"`
	ErrorURLs         int       `json:"error_urls"`
	LastSitemapSubmit time.Time `json:"last_sitemap_submit"`
	AuthValid         bool      `json:"auth_valid"`
}

// GoogleIndexTaskListResponse 任务列表响应
type GoogleIndexTaskListResponse struct {
	Tasks       []GoogleIndexTaskOutput `json:"tasks"`
	Total       int64                   `json:"total"`
	Page        int                     `json:"page"`
	PageSize    int                     `json:"page_size"`
	TotalPages  int                     `json:"total_pages"`
}

// GoogleIndexTaskItemPageResponse 任务项分页响应
type GoogleIndexTaskItemPageResponse struct {
	Items []GoogleIndexTaskItemOutput `json:"items"`
	Total int64                       `json:"total"`
	Page  int                         `json:"page"`
	Size  int                         `json:"size"`
}
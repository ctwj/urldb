package dto

// UpdateAIConfigRequest 更新AI配置请求
type UpdateAIConfigRequest struct {
	APIKey       *string  `json:"api_key,omitempty"`
	APIURL       *string  `json:"ai_api_url,omitempty"`
	Model        *string  `json:"ai_model,omitempty"`
	MaxTokens    *int     `json:"ai_max_tokens,omitempty"`
	Temperature  *float32 `json:"ai_temperature,omitempty"`
	Organization *string  `json:"ai_organization,omitempty"`
	Proxy        *string  `json:"ai_proxy,omitempty"`
	Timeout      *int     `json:"ai_timeout,omitempty"`
	RetryCount   *int     `json:"ai_retry_count,omitempty"`
}

// TestAIConnectionRequest 测试AI连接请求
type TestAIConnectionRequest struct {
	APIKey      *string  `json:"api_key,omitempty"`
	APIURL      *string  `json:"ai_api_url,omitempty"`
	Model       *string  `json:"ai_model,omitempty"`
	MaxTokens   *int     `json:"ai_max_tokens,omitempty"`
	Temperature *float32 `json:"ai_temperature,omitempty"`
	Timeout     *int     `json:"ai_timeout,omitempty"`
	RetryCount  *int     `json:"ai_retry_count,omitempty"`
}

// GenerateTextRequest 通用文本生成请求
type GenerateTextRequest struct {
	Prompt  string        `json:"prompt" binding:"required"`
	Options []ChatOption  `json:"options,omitempty"`
}

// ChatOption 对话选项
type ChatOption struct {
	Type  string      `json:"type" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

// AskQuestionRequest 问答请求
type AskQuestionRequest struct {
	Question string `json:"question" binding:"required"`
	Context  string `json:"context" binding:"required"`
}

// AnalyzeTextRequest 文本分析请求
type AnalyzeTextRequest struct {
	Text        string `json:"text" binding:"required"`
	AnalysisType string `json:"analysis_type" binding:"required"`
}

// GenerateContentRequest 生成内容请求
type GenerateContentRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
}

// ApplyGeneratedContentRequest 应用生成内容请求
type ApplyGeneratedContentRequest struct {
	SessionID             string   `json:"session_id" binding:"required"`
	ResourceID            uint     `json:"resource_id" binding:"required"`
	GeneratedTitle        string   `json:"generated_title"`
	GeneratedDescription  string   `json:"generated_description"`
	GeneratedSEOTitle     string   `json:"generated_seo_title"`
	GeneratedSEODescription string `json:"generated_seo_description"`
	GeneratedSEOKeywords  []string `json:"generated_seo_keywords"`
	AIModelUsed           string   `json:"ai_model_used"`
}

// ClassifyRequest 分类请求
type ClassifyRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
}

// ApplyClassificationRequest 应用分类请求
type ApplyClassificationRequest struct {
	SessionID             string  `json:"session_id" binding:"required"`
	ResourceID            uint    `json:"resource_id" binding:"required"`
	SuggestedCategoryID   uint    `json:"suggested_category_id" binding:"required"`
	SuggestedCategoryName string  `json:"suggested_category_name" binding:"required"`
	Confidence            float64 `json:"confidence" binding:"required"`
	AIModelUsed           string  `json:"ai_model_used"`
}

// ToolCallRequest 工具调用请求
type ToolCallRequest struct {
	Params map[string]interface{} `json:"params" binding:"required"`
}
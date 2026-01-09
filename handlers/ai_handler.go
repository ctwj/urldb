package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/ai/mcp"
	"github.com/ctwj/urldb/pkg/ai/service"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// AIHandler AI处理器
type AIHandler struct {
	aiService  *service.AIService
	mcpManager *mcp.MCPManager
}

// NewAIHandler 创建AI处理器
func NewAIHandler(aiService *service.AIService, mcpManager *mcp.MCPManager) *AIHandler {
	return &AIHandler{
		aiService:  aiService,
		mcpManager: mcpManager,
	}
}

// NewAIHandlerWithConfig 创建AI处理器（使用配置管理器）
func NewAIHandlerWithConfig(configManager service.ConfigManager, repoManager *repo.RepositoryManager, mcpManager *mcp.MCPManager) (*AIHandler, error) {
	aiService, err := service.NewAIServiceWithConfig(configManager, repoManager)
	if err != nil {
		return nil, fmt.Errorf("创建AI服务失败: %v", err)
	}

	return &AIHandler{
		aiService:  aiService,
		mcpManager: mcpManager,
	}, nil
}

// GetAIConfig 获取AI配置
func (h *AIHandler) GetAIConfig(c *gin.Context) {
	// 获取配置值
	apiURL, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAIAPIURL)
	model, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAIModel)
	maxTokens, _ := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyAIMaxTokens)
	temperature, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAITemperature)
	organization, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAIOrganization)
	proxy, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAIProxy)
	timeout, _ := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyAITimeout)
	retryCount, _ := repoManager.SystemConfigRepository.GetConfigInt(entity.ConfigKeyAIRetryCount)

	// 检查API Key是否已配置（不返回实际的API Key值）
	apiKey, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAIAPIKey)
	apiKeyConfigured := apiKey != ""

	config := map[string]interface{}{
		"ai_api_url":        apiURL,
		"ai_model":          model,
		"ai_max_tokens":     maxTokens,
		"ai_temperature":    temperature,
		"ai_organization":   organization,
		"ai_proxy":          proxy,
		"ai_timeout":        timeout,
		"ai_retry_count":    retryCount,
		"ai_api_key_configured": apiKeyConfigured,
	}

	SuccessResponse(c, config)
}

// UpdateAIConfig 更新AI配置
func (h *AIHandler) UpdateAIConfig(c *gin.Context) {
	var req dto.UpdateAIConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error("JSON绑定失败: %v", err)
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	adminUsername, _ := c.Get("username")
	clientIP, _ := c.Get("client_ip")
	utils.Info("UpdateAIConfig - 管理员更新AI配置 - 管理员: %s, IP: %s", adminUsername, clientIP)

	// 更新配置
	var configs []entity.SystemConfig
	if req.APIURL != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIAPIURL,
			Value: *req.APIURL,
			Type:  entity.ConfigTypeString,
		})
	}
	if req.Model != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIModel,
			Value: *req.Model,
			Type:  entity.ConfigTypeString,
		})
	}
	if req.MaxTokens != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIMaxTokens,
			Value: strconv.Itoa(*req.MaxTokens),
			Type:  entity.ConfigTypeInt,
		})
	}
	if req.Temperature != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAITemperature,
			Value: strconv.FormatFloat(float64(*req.Temperature), 'f', 2, 64),
			Type:  entity.ConfigTypeString,
		})
	}
	if req.Organization != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIOrganization,
			Value: *req.Organization,
			Type:  entity.ConfigTypeString,
		})
	}
	if req.Proxy != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIProxy,
			Value: *req.Proxy,
			Type:  entity.ConfigTypeString,
		})
	}
	if req.Timeout != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAITimeout,
			Value: strconv.Itoa(*req.Timeout),
			Type:  entity.ConfigTypeInt,
		})
	}
	if req.RetryCount != nil {
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIRetryCount,
			Value: strconv.Itoa(*req.RetryCount),
			Type:  entity.ConfigTypeInt,
		})
	}
	if req.APIKey != nil {
		// TODO: 实现加密存储 API Key
		// 目前暂时使用明文存储，生产环境需要实现加密
		configs = append(configs, entity.SystemConfig{
			Key:   entity.ConfigKeyAIAPIKey,
			Value: *req.APIKey,
			Type:  entity.ConfigTypeString,
		})
	}

	if len(configs) > 0 {
		if err := repoManager.SystemConfigRepository.UpsertConfigs(configs); err != nil {
			utils.Error("更新 AI 配置失败: %v", err)
			ErrorResponse(c, "更新配置失败", http.StatusInternalServerError)
			return
		}
	}

	// 重新加载 AI 客户端配置
	if err := h.aiService.ReloadClient(); err != nil {
		utils.Error("重新加载 AI 客户端配置失败: %v", err)
		ErrorResponse(c, "重新加载配置失败", http.StatusInternalServerError)
		return
	}

	utils.Info("AI 配置更新成功 - 管理员: %s", adminUsername)
	SuccessResponse(c, "配置更新成功")
}

// TestAIConnection 测试AI连接
func (h *AIHandler) TestAIConnection(c *gin.Context) {
	err := h.aiService.TestConnection()
	if err != nil {
		utils.Error("AI 连接测试失败: %v", err)
		ErrorResponse(c, "连接测试失败，请检查配置", http.StatusBadRequest)
		return
	}

	utils.Info("AI 连接测试成功")
	SuccessResponse(c, "连接测试成功")
}

// GenerateText 通用文本生成功能
func (h *AIHandler) GenerateText(c *gin.Context) {
	var req dto.GenerateTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	// 转换 DTO ChatOption 到 service ChatOption
	options := make([]service.ChatOption, 0, len(req.Options))
	for _, opt := range req.Options {
		switch opt.Type {
		case "max_tokens":
			if tokens, ok := opt.Value.(float64); ok {
				options = append(options, service.WithMaxTokens(int(tokens)))
			}
		case "temperature":
			if temp, ok := opt.Value.(float64); ok {
				options = append(options, service.WithTemperature(float32(temp)))
			}
		case "system_prompt":
			if prompt, ok := opt.Value.(string); ok {
				options = append(options, service.WithSystemPrompt(prompt))
			}
		}
	}

	result, err := h.aiService.GenerateText(req.Prompt, options...)
	if err != nil {
		utils.Error("文本生成失败: %v", err)
		ErrorResponse(c, "文本生成失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"result": result,
	})
}

// AskQuestion 通用问答功能
func (h *AIHandler) AskQuestion(c *gin.Context) {
	var req dto.AskQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	result, err := h.aiService.AskQuestion(req.Question, req.Context)
	if err != nil {
		utils.Error("问答处理失败: %v", err)
		ErrorResponse(c, "问答处理失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"answer": result,
	})
}

// AnalyzeText 通用文本分析功能
func (h *AIHandler) AnalyzeText(c *gin.Context) {
	var req dto.AnalyzeTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	result, err := h.aiService.AnalyzeText(req.Text, req.AnalysisType)
	if err != nil {
		utils.Error("文本分析失败: %v", err)
		ErrorResponse(c, "文本分析失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"result": result,
	})
}

// GenerateContentPreview 生成内容预览 - 用户发起后 AI 生成内容，但不立即写入数据库
func (h *AIHandler) GenerateContentPreview(c *gin.Context) {
	var req dto.GenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	preview, err := h.aiService.GenerateContentPreview(req.ResourceID)
	if err != nil {
		utils.Error("内容预览生成失败: %v", err)
		ErrorResponse(c, "内容预览生成失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, preview)
}

// ApplyGeneratedContent 应用生成的内容 - 用户确认后才保存到数据库
func (h *AIHandler) ApplyGeneratedContent(c *gin.Context) {
	var req dto.ApplyGeneratedContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	preview := &service.GeneratedContentPreview{
		SessionID:             req.SessionID,
		ResourceID:            req.ResourceID,
		GeneratedTitle:        req.GeneratedTitle,
		GeneratedDescription:  req.GeneratedDescription,
		GeneratedSEOTitle:     req.GeneratedSEOTitle,
		GeneratedSEODescription: req.GeneratedSEODescription,
		GeneratedSEOKeywords:  req.GeneratedSEOKeywords,
		AIModelUsed:           req.AIModelUsed,
	}

	err := h.aiService.ApplyGeneratedContent(preview)
	if err != nil {
		utils.Error("应用生成内容失败: %v", err)
		ErrorResponse(c, "应用生成内容失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "内容已成功应用",
	})
}

// ClassifyResourcePreview 分类资源预览 - AI 生成分类建议，但不立即应用到数据库
func (h *AIHandler) ClassifyResourcePreview(c *gin.Context) {
	var req dto.ClassifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	preview, err := h.aiService.ClassifyResourcePreview(req.ResourceID)
	if err != nil {
		utils.Error("分类预览生成失败: %v", err)
		ErrorResponse(c, "分类预览生成失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, preview)
}

// ApplyClassification 应用分类建议 - 用户确认后才写入数据库
func (h *AIHandler) ApplyClassification(c *gin.Context) {
	var req dto.ApplyClassificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	preview := &service.ClassificationPreview{
		SessionID:             req.SessionID,
		ResourceID:            req.ResourceID,
		SuggestedCategoryID:   req.SuggestedCategoryID,
		SuggestedCategoryName: req.SuggestedCategoryName,
		Confidence:            req.Confidence,
		AIModelUsed:           req.AIModelUsed,
	}

	err := h.aiService.ApplyClassification(preview)
	if err != nil {
		utils.Error("应用分类失败: %v", err)
		ErrorResponse(c, "应用分类失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "分类已成功应用",
	})
}
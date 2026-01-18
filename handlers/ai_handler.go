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
	aiService    *service.AIService
	promptService *service.PromptService
	mcpManager   *mcp.MCPManager
}

// NewAIHandler 创建AI处理器
func NewAIHandler(aiService *service.AIService, promptService *service.PromptService, mcpManager *mcp.MCPManager) *AIHandler {
	return &AIHandler{
		aiService:    aiService,
		promptService: promptService,
		mcpManager:   mcpManager,
	}
}

// NewAIHandlerWithConfig 创建AI处理器（使用配置管理器）
func NewAIHandlerWithConfig(configManager service.ConfigManager, repoManager *repo.RepositoryManager, mcpManager *mcp.MCPManager) (*AIHandler, error) {
	// 先创建OpenAI客户端
	client, err := service.NewOpenAIClientWithConfig(configManager)
	if err != nil {
		return nil, fmt.Errorf("创建OpenAI客户端失败: %v", err)
	}

	// 创建支持MCP的AI服务
	aiService, err := service.NewAIServiceWithMCP(client, repoManager, mcpManager)
	if err != nil {
		return nil, fmt.Errorf("创建AI服务失败: %v", err)
	}

	// 创建提示词服务
	promptService := service.NewPromptService(repoManager.GetDB())

	return &AIHandler{
		aiService:    aiService,
		promptService: promptService,
		mcpManager:   mcpManager,
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

	// 获取API Key
	apiKey, _ := repoManager.SystemConfigRepository.GetConfigValue(entity.ConfigKeyAIAPIKey)
	apiKeyConfigured := apiKey != ""

	config := map[string]interface{}{
		"api_key":           apiKey,
		"ai_api_key_configured": apiKeyConfigured,
		"ai_api_url":        apiURL,
		"ai_model":          model,
		"ai_max_tokens":     maxTokens,
		"ai_temperature":    temperature,
		"ai_organization":   organization,
		"ai_proxy":          proxy,
		"ai_timeout":        timeout,
		"ai_retry_count":    retryCount,
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
	var req dto.TestAIConnectionRequest
	var response string
	var err error

	// 尝试解析请求体，如果没有提供则使用保存的配置
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		// 如果没有提供配置参数，使用保存的配置进行测试
		utils.Info("使用保存的配置进行AI连接测试")
		response, err = h.aiService.TestConnectionWithResponse()
		if err != nil {
			utils.Error("AI 连接测试失败: %v", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "连接测试失败",
				"error":   err.Error(),
			})
			return
		}
	} else {
		// 使用提供的临时配置进行测试
		utils.Info("使用临时配置进行AI连接测试")
		response, err = h.aiService.TestConnectionWithConfigAndResponse(&service.AIConfig{
			APIKey:      req.APIKey,
			APIURL:      req.APIURL,
			Model:       req.Model,
			MaxTokens:   req.MaxTokens,
			Temperature: req.Temperature,
			Timeout:     req.Timeout,
			RetryCount:  req.RetryCount,
		})
		if err != nil {
			utils.Error("AI 连接测试失败: %v", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"message": "连接测试失败",
				"error":   err.Error(),
			})
			return
		}
	}

	utils.Info("AI 连接测试成功，返回响应: %s", response)
	SuccessResponse(c, map[string]interface{}{
		"success": true,
		"message": "连接测试成功",
		"response": response,
	})
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

// GetAvailableTools 获取所有可用的MCP工具
func (h *AIHandler) GetAvailableTools(c *gin.Context) {
	tools, err := h.aiService.GetAvailableTools()
	if err != nil {
		utils.Error("获取可用工具失败: %v", err)
		ErrorResponse(c, "获取工具失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"tools": tools,
		"count": len(tools),
	})
}

// CallTool 调用指定的MCP工具
func (h *AIHandler) CallTool(c *gin.Context) {
	var req dto.CallToolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	result, err := h.aiService.CallTool(req.ToolName, req.Params)
	if err != nil {
		utils.Error("工具调用失败: %v", err)
		ErrorResponse(c, "工具调用失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, result)
}

// GenerateTextWithTools 使用工具的文本生成
func (h *AIHandler) GenerateTextWithTools(c *gin.Context) {
	var req dto.GenerateTextWithToolsRequest
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

	result, err := h.aiService.GenerateTextWithTools(req.Prompt, options...)
	if err != nil {
		utils.Error("工具增强文本生成失败: %v", err)
		ErrorResponse(c, "文本生成失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"result": result,
		"with_tools": true,
	})
}

// GetPrompts 获取所有提示词
func (h *AIHandler) GetPrompts(c *gin.Context) {
	prompts, err := h.promptService.GetAllPrompts()
	if err != nil {
		utils.Error("获取提示词列表失败: %v", err)
		ErrorResponse(c, "获取提示词列表失败", http.StatusInternalServerError)
		return
	}

	// 为每个提示词添加类型信息
	result := make([]map[string]interface{}, 0, len(prompts))
	for _, prompt := range prompts {
		promptData := map[string]interface{}{
			"id":            prompt.ID,
			"name":          prompt.Name,
			"type":          prompt.Type,
			"system_content": prompt.SystemContent,
			"user_content":   prompt.UserContent,
			"description":   prompt.Description,
			"is_active":     prompt.IsActive,
			"created_at":    prompt.CreatedAt,
			"updated_at":    prompt.UpdatedAt,
			"type_name":        entity.GetPromptTypeName(prompt.Type),
			"type_description": entity.GetPromptTypeDescription(prompt.Type),
		}

		// 获取变量列表
		variables, _ := h.promptService.GetPromptVariables(&prompt)
		promptData["variables"] = variables

		result = append(result, promptData)
	}

	SuccessResponse(c, result)
}

// UpdatePrompt 更新提示词
func (h *AIHandler) UpdatePrompt(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的提示词ID", http.StatusBadRequest)
		return
	}

	var req struct {
		SystemContent *string `json:"system_content"`
		UserContent   *string `json:"user_content"`
		Description   *string `json:"description"`
	}

	utils.Info("开始解析提示词更新请求")
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error("参数绑定失败: %v", err)
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}
	utils.Info("参数绑定成功: SystemContent=%v, UserContent=%v, Description=%v",
		req.SystemContent != nil, req.UserContent != nil, req.Description != nil)

	// 检查是否提供了任何更新内容
	if req.SystemContent == nil && req.UserContent == nil && req.Description == nil {
		ErrorResponse(c, "请提供要更新的内容", http.StatusBadRequest)
		return
	}

	// 根据提供的字段进行智能更新
	var systemContent, userContent, description string
	if req.SystemContent != nil {
		systemContent = *req.SystemContent
	}
	if req.UserContent != nil {
		userContent = *req.UserContent
	}
	if req.Description != nil {
		description = *req.Description
	}

	// 使用新的完整更新方法
	err = h.promptService.UpdateFullPrompt(uint(id), systemContent, userContent, description)

	if err != nil {
		utils.Error("更新提示词失败: %v", err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "提示词更新成功",
		"id":      id,
	})
}

// TogglePromptStatus 切换提示词状态
func (h *AIHandler) TogglePromptStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ErrorResponse(c, "无效的提示词ID", http.StatusBadRequest)
		return
	}

	if err := h.promptService.TogglePromptStatus(uint(id)); err != nil {
		utils.Error("切换提示词状态失败: %v", err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "提示词状态切换成功",
		"id":      id,
	})
}

// TestPrompt 测试提示词
func (h *AIHandler) TestPrompt(c *gin.Context) {
	var req struct {
		Type string                 `json:"type" binding:"required"`
		Data map[string]interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	result, err := h.promptService.TestPrompt(req.Type, req.Data)
	if err != nil {
		utils.Error("测试提示词失败: %v", err)
		ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"type":   req.Type,
		"result": result,
		"data":   req.Data,
	})
}

// InitDefaultPrompts 初始化默认提示词
func (h *AIHandler) InitDefaultPrompts(c *gin.Context) {
	if err := h.promptService.CreateDefaultPrompts(); err != nil {
		utils.Error("初始化默认提示词失败: %v", err)
		ErrorResponse(c, "初始化默认提示词失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "默认提示词初始化成功",
	})
}
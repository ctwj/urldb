package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/ai/mcp"
	"github.com/ctwj/urldb/utils"
	"github.com/sashabaranov/go-openai"
)

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// AIConfig AI配置结构
type AIConfig struct {
	APIKey      *string
	APIURL      *string
	Model       *string
	MaxTokens   *int
	Temperature *float32
	Timeout     *int
	RetryCount  *int
}

// ToolDefinition OpenAI工具定义结构
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCallResult 工具调用结果
type ToolCallResult struct {
	ToolName string      `json:"tool_name"`
	Result   interface{} `json:"result"`
	Error    string      `json:"error,omitempty"`
}

// AIService 主AI服务，提供通用AI能力供其他模块调用
type AIService struct {
	client        *OpenAIClient
	contentGen    *ContentGenerator
	classifier    *Classifier
	promptService *PromptService
	repoManager   *repo.RepositoryManager
	mcpManager    *mcp.MCPManager
}

// NewAIServiceWithConfig 创建AI服务
func NewAIServiceWithConfig(configManager ConfigManager, repoManager *repo.RepositoryManager) (*AIService, error) {
	client, err := NewOpenAIClientWithConfig(configManager)
	if err != nil {
		return nil, fmt.Errorf("创建OpenAI客户端失败: %v", err)
	}

	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)
	promptService := NewPromptService(repoManager.GetDB())

	return &AIService{
		client:        client,
		contentGen:    contentGen,
		classifier:    classifier,
		promptService: promptService,
		repoManager:   repoManager,
	}, nil
}

// NewAIService 创建AI服务
func NewAIService(client *OpenAIClient, repoManager *repo.RepositoryManager) (*AIService, error) {
	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)
	promptService := NewPromptService(repoManager.GetDB())

	return &AIService{
		client:        client,
		contentGen:    contentGen,
		classifier:    classifier,
		promptService: promptService,
		repoManager:   repoManager,
	}, nil
}

// NewAIServiceWithMCP 创建支持MCP的AI服务
func NewAIServiceWithMCP(client *OpenAIClient, repoManager *repo.RepositoryManager, mcpManager *mcp.MCPManager) (*AIService, error) {
	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)
	promptService := NewPromptService(repoManager.GetDB())

	return &AIService{
		client:        client,
		contentGen:    contentGen,
		classifier:    classifier,
		promptService: promptService,
		repoManager:   repoManager,
		mcpManager:    mcpManager,
	}, nil
}

// GenerateText 通用文本生成 - 供其他模块调用
func (as *AIService) GenerateText(prompt string, options ...ChatOption) (string, error) {
	traceID := generateAITraceID()
	utils.Info("[AI][%s] 开始处理文本生成请求", traceID)

	// 如果有 MCP 管理器，尝试使用工具增强的生成
	if as.mcpManager != nil {
		utils.Debug("[AI][%s] MCP 管理器已初始化，尝试使用工具增强生成", traceID)
		result, err := as.GenerateTextWithTools(prompt, options...)
		if err != nil {
			utils.Warn("[AI][%s] 工具增强生成失败: %v", traceID, err)
			return "", err
		}
		utils.Info("[AI][%s] 工具增强生成成功", traceID)
		return result, nil
	} else {
		utils.Debug("[AI][%s] MCP 管理器未初始化，使用普通生成", traceID)
	}

	return as.generatePlainTextWithFallback(traceID, prompt, options...)
}

// GenerateTextWithoutTools 显式禁用工具增强流程，直接走普通生成
func (as *AIService) GenerateTextWithoutTools(prompt string, options ...ChatOption) (string, error) {
	traceID := generateAITraceID()
	utils.Info("[AI][%s] 开始处理文本生成请求（禁用工具）", traceID)
	return as.generatePlainTextWithFallback(traceID, prompt, options...)
}

// getToolListSummary 生成工具列表摘要
func getToolListSummary(tools []ToolDefinition) string {
	var summary string
	for _, tool := range tools {
		summary += fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description)
	}
	return summary
}

// getToolListWithParams 生成包含参数信息的工具列表
func getToolListWithParams(tools []ToolDefinition) string {
	var summary string
	for _, tool := range tools {
		summary += fmt.Sprintf("🔹 %s: %s\n", tool.Name, tool.Description)

		// 解析参数信息
		if tool.Parameters != nil {
			if properties, ok := tool.Parameters["properties"].(map[string]interface{}); ok {
				var required []interface{}
				if req, ok := tool.Parameters["required"].([]interface{}); ok {
					required = req
				}

				if len(required) > 0 {
					summary += "   【必需参数】："
					for i, req := range required {
						if i > 0 {
							summary += "、"
						}
						summary += fmt.Sprintf("%v", req)
					}
					summary += "\n"
				}

				// 显示每个参数的详细信息
				for paramName, paramInfo := range properties {
					if paramMap, ok := paramInfo.(map[string]interface{}); ok {
						var isRequired bool
						for _, req := range required {
							if req == paramName {
								isRequired = true
								break
							}
						}

						reqMark := "可选"
						if isRequired {
							reqMark = "【必需】"
						}

						desc := ""
						if description, ok := paramMap["description"].(string); ok {
							desc = fmt.Sprintf(" - %s", description)
						}

						summary += fmt.Sprintf("   - %s (%s)%s\n", paramName, reqMark, desc)
					}
				}
			}
		}
		summary += "\n"
	}
	return summary
}

// needsTools 判断用户问题是否需要使用工具
func needsTools(prompt string) bool {
	normalized := strings.ToLower(strings.TrimSpace(prompt))
	if normalized == "" {
		return false
	}

	// 结构化内容生成提示（如资源标题优化）不应触发工具调用
	structuredMarkers := []string{
		"根据以下资源信息", "资源标题:", "资源描述:", "资源url:", "现有描述:",
		"可用分类列表", "请生成", "请直接返回最适合的分类",
	}
	markerHits := 0
	for _, marker := range structuredMarkers {
		if strings.Contains(normalized, marker) {
			markerHits++
		}
	}
	if markerHits >= 2 {
		return false
	}

	// 只有明确涉及外部实时信息/检索时才启用工具
	toolKeywords := []string{
		"现在几点", "当前时间", "今天几号", "日期", "天气", "预报", "新闻", "实时",
		"汇率", "股票", "价格", "搜索", "查询", "查一下", "网页", "网站", "抓取",
		"http://", "https://", "fetch ",
	}
	for _, keyword := range toolKeywords {
		if strings.Contains(normalized, keyword) {
			return true
		}
	}

	return false
}

func parseRequiredParams(requiredRaw interface{}) []string {
	required := make([]string, 0)

	switch req := requiredRaw.(type) {
	case []string:
		required = append(required, req...)
	case []interface{}:
		for _, item := range req {
			if reqStr, ok := item.(string); ok {
				required = append(required, reqStr)
			}
		}
	}

	return required
}

// getToolsAsNaturalLanguage 将工具定义转换为自然语言描述
func getToolsAsNaturalLanguage(tools []ToolDefinition) string {
	var description string
	description += "你可以使用以下工具来回答用户的问题：\n\n"

	for i, tool := range tools {
		description += fmt.Sprintf("工具%d：%s\n", i+1, tool.Name)
		description += fmt.Sprintf("- 描述：%s\n", tool.Description)

		// 解析参数信息
		if tool.Parameters != nil {
			if properties, ok := tool.Parameters["properties"].(map[string]interface{}); ok {
				requiredSet := make(map[string]struct{})
				for _, req := range parseRequiredParams(tool.Parameters["required"]) {
					requiredSet[req] = struct{}{}
				}

				// 显示每个参数的详细信息
				for paramName, paramInfo := range properties {
					if paramMap, ok := paramInfo.(map[string]interface{}); ok {
						_, isRequired := requiredSet[paramName]

						reqMark := "可选"
						if isRequired {
							reqMark = "必需"
						}

						desc := ""
						if description, ok := paramMap["description"].(string); ok {
							desc = fmt.Sprintf(" - %s", description)
						}

						// 添加枚举值信息（如果有）
						enumInfo := ""
						if enumValues, ok := paramMap["enum"].([]interface{}); ok {
							enumInfo = " (可选值: "
							for j, enum := range enumValues {
								if j > 0 {
									enumInfo += ", "
								}
								enumInfo += fmt.Sprintf("%v", enum)
							}
							enumInfo += ")"
						}

						description += fmt.Sprintf("- 参数：%s (%s)%s%s\n", paramName, reqMark, desc, enumInfo)
					}
				}
			}
		}
		description += "\n"
	}

	description += "工具调用格式：<工具名称: {\"参数名\": \"参数值\"}>\n"
	description += "通用示例：<工具名称: {\"参数1\": \"值1\", \"参数2\": \"值2\"}>\n\n"

	return description
}

// ToolCallFromContent 从内容解析出的工具调用
type ToolCallFromContent struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

// parseToolCallsFromContent 从响应内容中解析工具调用标记
// 支持 GLM 格式：<tool_name/> 或 <tool_name param1="value1" param2="value2"/>
// 也支持：<tool_name: {}> 或 <tool_name: {param1: value1, param2: value2}>
// 也支持跨行格式：<tool_name\n: {}>
// 也支持特殊字符格式：<tool_name\n⟶
// toolNameSet: 已注册的工具名称集合，用于过滤无效的工具调用
func parseToolCallsFromContent(content string, toolNameSet map[string]bool) []ToolCallFromContent {
	var toolCalls []ToolCallFromContent

	utils.Debug("[AI] 工具解析 原始内容: %q", content)

	toolNamePattern := `[A-Za-z0-9_-]+`

	// 先尝试匹配 JSON 格式的工具调用：<tool_name: {}>
	jsonRe := regexp.MustCompile(`(?s)<(` + toolNamePattern + `):\s*({[^}]*})>`)
	jsonMatches := jsonRe.FindAllStringSubmatch(content, -1)
	utils.Debug("[AI] 工具解析 JSON 格式匹配到 %d 个结果", len(jsonMatches))

	for i, match := range jsonMatches {
		utils.Debug("[AI] 工具解析 JSON 匹配 %d: %v", i, match)
		if len(match) < 3 {
			continue
		}

		toolName := match[1]

		// 检查工具名称是否在已注册的工具列表中
		if !toolNameSet[toolName] {
			utils.Debug("[AI] 工具解析 工具 %s 未注册，跳过", toolName)
			continue
		}

		jsonStr := match[2]
		params := make(map[string]interface{})

		if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
			utils.Debug("[AI] 工具解析 解析 JSON 参数失败: %v", err)
			params = map[string]interface{}{"args": jsonStr}
		}

		toolCalls = append(toolCalls, ToolCallFromContent{
			Name:   toolName,
			Params: params,
		})
		utils.Debug("[AI] 工具解析 解析工具: %s, 参数: %v", toolName, params)
	}

	// 如果没有匹配到 JSON 格式，尝试匹配简单标签格式：<tool_name> 或 <tool_name/> 或 <tool_name\n
	if len(toolCalls) == 0 {
		simpleRe := regexp.MustCompile(`<(` + toolNamePattern + `)[\s\n>]`)
		simpleMatches := simpleRe.FindAllStringSubmatch(content, -1)
		utils.Debug("[AI] 工具解析 简单标签格式匹配到 %d 个结果", len(simpleMatches))

		for i, match := range simpleMatches {
			utils.Debug("[AI] 工具解析 简单标签匹配 %d: %v", i, match)
			if len(match) < 2 {
				continue
			}

			toolName := match[1]

			// 检查工具名称是否在已注册的工具列表中
			if !toolNameSet[toolName] {
				utils.Debug("[AI] 工具解析 工具 %s 未注册，跳过", toolName)
				continue
			}

			toolCalls = append(toolCalls, ToolCallFromContent{
				Name:   toolName,
				Params: map[string]interface{}{},
			})
			utils.Debug("[AI] 工具解析 解析工具: %s, 参数: map[]", toolName)
		}
	}

	// 如果还没有匹配到，尝试匹配 HTML 属性格式：<tool_name param1="value1"/>
	if len(toolCalls) == 0 {
		htmlRe := regexp.MustCompile(`<(` + toolNamePattern + `)(\s+[^>]*)>`)
		htmlMatches := htmlRe.FindAllStringSubmatch(content, -1)
		utils.Debug("[AI] 工具解析 HTML 格式匹配到 %d 个结果", len(htmlMatches))

		for i, match := range htmlMatches {
			utils.Debug("[AI] 工具解析 HTML 匹配 %d: %v", i, match)
			if len(match) < 3 {
				continue
			}

			toolName := match[1]

			// 检查工具名称是否在已注册的工具列表中
			if !toolNameSet[toolName] {
				utils.Debug("[AI] 工具解析 工具 %s 未注册，跳过", toolName)
				continue
			}

			paramsStr := match[2]
			params := make(map[string]interface{})

			paramRe := regexp.MustCompile(`(\w+)="([^"]*)"`)
			paramMatches := paramRe.FindAllStringSubmatch(paramsStr, -1)
			for _, paramMatch := range paramMatches {
				if len(paramMatch) >= 3 {
					params[paramMatch[1]] = paramMatch[2]
				}
			}

			toolCalls = append(toolCalls, ToolCallFromContent{
				Name:   toolName,
				Params: params,
			})
			utils.Debug("[AI] 工具解析 解析工具: %s, 参数: %v", toolName, params)
		}
	}

	// 检查是否已经包含了工具结果
	// 如果响应内容中已经包含了详细的工具结果（如日期时间信息），说明 AI 已经自己处理了工具调用
	// 这种情况下，我们不应该再调用工具
	if len(toolCalls) > 0 {
		// 检查响应中是否包含具体的时间数据（不仅仅是关键词）
		// 例如：2025年6月17日、10:32:15、timestamp: 1718601135 等
		hasResult := false

		// 检查是否包含具体的时间格式
		timePatterns := []string{
			`\d{4}年\d{1,2}月\d{1,2}日`, // 中文日期格式
			`\d{4}-\d{1,2}-\d{1,2}`,  // 英文日期格式
			`\d{1,2}:\d{2}:\d{2}`,    // 时间格式
			`timestamp:\s*\d+`,       // 时间戳格式
		}

		for _, pattern := range timePatterns {
			if matched, _ := regexp.MatchString(pattern, content); matched {
				hasResult = true
				utils.Debug("[AI] 工具解析 检测到时间数据: %s", pattern)
				break
			}
		}

		if hasResult {
			utils.Debug("[AI] 工具解析 检测到响应中已包含工具结果，忽略工具调用")
			return []ToolCallFromContent{}
		}

		// 检查响应长度，如果很短（比如只有工具调用标记），说明没有工具结果
		// 去除工具调用标记后的内容长度
		cleanContent := regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")
		cleanContent = strings.TrimSpace(cleanContent)
		if len(cleanContent) < 10 {
			utils.Debug("[AI] 工具解析 响应内容过短，没有工具结果")
		} else {
			utils.Debug("[AI] 工具解析 响应内容长度: %d", len(cleanContent))
		}
	}

	return toolCalls
}

// cleanToolCallMarkers 清理响应内容中的工具调用标记
func cleanToolCallMarkers(content string) string {
	// 移除工具调用标记：<tool_name>...</tool_name> 或 <tool_name/> 或 <tool_name: {}> 等
	// 也支持没有闭合标签的格式：<tool_name\n⟶
	re := regexp.MustCompile(`<[A-Za-z0-9_-]+(?::\s*{[^}]*})?\s*/?>\s*</[A-Za-z0-9_-]+>|<[A-Za-z0-9_-]+(?::\s*{[^}]*})?\s*/?>|<[A-Za-z0-9_-]+>|<[A-Za-z0-9_-]+[\s\n]`)
	cleanContent := re.ReplaceAllString(content, "")

	// 清理多余的空行
	cleanContent = regexp.MustCompile(`\n\s*\n\s*\n`).ReplaceAllString(cleanContent, "\n\n")

	// 去除首尾空白
	cleanContent = strings.TrimSpace(cleanContent)

	return cleanContent
}

// AskQuestion 通用问答 - 供其他模块调用
func (as *AIService) AskQuestion(question string, context string) (string, error) {
	// 获取系统提示词
	systemPrompt, err := as.promptService.RenderSystemPromptByType(entity.PromptTypeQATemplate, nil)
	if err != nil {
		// 如果获取失败，使用默认系统提示词
		systemPrompt = "你是一个专业的问答助手，擅长基于提供的上下文信息给出准确的回答。你需要严格根据上下文信息回答问题，不要编造或推测信息。如果上下文中没有相关信息，请明确说明。"
	}

	// 获取用户提示词
	userPrompt, err := as.promptService.RenderUserPromptByType(entity.PromptTypeQATemplate, map[string]interface{}{
		"Context":  context,
		"Question": question,
	})
	if err != nil {
		// 如果获取失败，使用默认用户提示词
		userPrompt = fmt.Sprintf("根据以下上下文回答问题：\n\n上下文：%s\n\n问题：%s\n\n请基于提供的上下文信息给出准确的回答。", context, question)
	}

	// 直接构建消息，不通过GenerateText避免重复添加系统提示词
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	resp, err := as.client.Chat(messages, WithMaxTokens(500), WithTemperature(0.7))
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	return resp.Choices[0].Message.Content, nil
}

// AnalyzeText 通用文本分析 - 供其他模块调用
func (as *AIService) AnalyzeText(text string, analysisType string) (string, error) {
	// 获取系统提示词
	systemPrompt, err := as.promptService.RenderSystemPromptByType(entity.PromptTypeAnalysisTemplate, nil)
	if err != nil {
		// 如果获取失败，使用默认系统提示词
		systemPrompt = "你是一个专业的文本分析专家，擅长对各类文本进行深入分析。你需要根据用户指定的分析类型，对提供的文本进行全面、准确的分析，并提供有价值的见解。"
	}

	// 获取用户提示词
	userPrompt, err := as.promptService.RenderUserPromptByType(entity.PromptTypeAnalysisTemplate, map[string]interface{}{
		"Text":         text,
		"AnalysisType": analysisType,
	})
	if err != nil {
		// 如果获取失败，使用默认用户提示词
		userPrompt = fmt.Sprintf("请对以下文本进行%s分析：\n\n%s", analysisType, text)
	}

	// 直接构建消息，不通过GenerateText避免重复添加系统提示词
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	resp, err := as.client.Chat(messages, WithMaxTokens(300), WithTemperature(0.5))
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateContentPreview 生成内容预览
func (as *AIService) GenerateContentPreview(resourceID uint) (*GeneratedContentPreview, error) {
	return as.contentGen.GenerateContentPreview(resourceID)
}

// ApplyGeneratedContent 应用生成的内容
func (as *AIService) ApplyGeneratedContent(preview *GeneratedContentPreview) error {
	return as.contentGen.ApplyGeneratedContent(preview)
}

// ClassifyResourcePreview 分类资源预览
func (as *AIService) ClassifyResourcePreview(resourceID uint) (*ClassificationPreview, error) {
	return as.classifier.ClassifyResourcePreview(resourceID)
}

// ApplyClassification 应用分类建议
func (as *AIService) ApplyClassification(preview *ClassificationPreview) error {
	return as.classifier.ApplyClassification(preview)
}

// TestConnection 测试AI连接
func (as *AIService) TestConnection() error {
	_, err := as.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return err
}

// TestConnectionWithResponse 测试AI连接并返回响应
func (as *AIService) TestConnectionWithResponse() (string, error) {
	response, err := as.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return response, err
}

// TestConnectionWithConfig 使用临时配置测试AI连接
func (as *AIService) TestConnectionWithConfig(config *AIConfig) error {
	// 创建临时客户端
	tempClient, err := as.createTempClient(config)
	if err != nil {
		return fmt.Errorf("创建临时客户端失败: %v", err)
	}

	// 创建临时 AIService
	tempAIService := &AIService{
		client:      tempClient,
		contentGen:  NewContentGenerator(tempClient, as.repoManager),
		classifier:  NewClassifier(tempClient, as.repoManager),
		repoManager: as.repoManager,
	}

	// 使用临时 AIService 询问模型信息
	_, err = tempAIService.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return err
}

// TestConnectionWithConfigAndResponse 使用临时配置测试AI连接并返回响应
func (as *AIService) TestConnectionWithConfigAndResponse(config *AIConfig) (string, error) {
	// 创建临时客户端
	tempClient, err := as.createTempClient(config)
	if err != nil {
		return "", fmt.Errorf("创建临时客户端失败: %v", err)
	}

	// 创建临时 AIService
	tempAIService := &AIService{
		client:      tempClient,
		contentGen:  NewContentGenerator(tempClient, as.repoManager),
		classifier:  NewClassifier(tempClient, as.repoManager),
		repoManager: as.repoManager,
	}

	// 使用临时 AIService 询问模型信息
	response, err := tempAIService.GenerateText("你是什么AI模型？请详细介绍你的名称、版本和能力。")
	return response, err
}

// createTempClient 创建临时客户端
func (as *AIService) createTempClient(config *AIConfig) (*OpenAIClient, error) {
	if config.APIKey == nil || *config.APIKey == "" {
		return nil, fmt.Errorf("API Key 不能为空")
	}

	// 设置默认值
	baseURL := "https://api.openai.com/v1"
	if config.APIURL != nil && *config.APIURL != "" {
		baseURL = *config.APIURL
	}

	model := "gpt-3.5-turbo"
	if config.Model != nil && *config.Model != "" {
		model = *config.Model
	}

	timeout := 30 * time.Second
	if config.Timeout != nil {
		timeout = time.Duration(*config.Timeout) * time.Second
	}

	retryCount := 3
	if config.RetryCount != nil {
		retryCount = *config.RetryCount
	}

	// 创建 OpenAI 客户端配置
	clientConfig := openai.DefaultConfig(*config.APIKey)
	clientConfig.BaseURL = baseURL

	// 设置超时
	clientConfig.HTTPClient = &http.Client{
		Timeout: timeout,
	}

	// 创建 OpenAI 客户端
	openaiClient := openai.NewClientWithConfig(clientConfig)

	// 创建临时 OpenAI 客户端包装器
	tempOpenAIClient := &OpenAIClient{
		apiKey:       *config.APIKey,
		baseURL:      baseURL,
		model:        model,
		organization: "",
		proxy:        "",
		timeout:      timeout,
		retryCount:   retryCount,
		client:       openaiClient,
		config:       nil,
	}

	return tempOpenAIClient, nil
}

// ReloadClient 重新加载客户端配置
func (as *AIService) ReloadClient() error {
	return as.client.ReloadConfig()
}

// GetModel 获取当前使用的模型
func (as *AIService) GetModel() string {
	return as.client.model
}

// GetAvailableTools 获取所有可用的MCP工具
func (as *AIService) GetAvailableTools() ([]ToolDefinition, error) {
	if as.mcpManager == nil {
		return nil, fmt.Errorf("MCP管理器未初始化")
	}

	var tools []ToolDefinition
	services := as.mcpManager.ListServices()

	utils.Info("[AI] 检查 %d 个服务", len(services))

	for _, serviceName := range services {
		// 检查服务健康状态
		if !as.mcpManager.CheckServiceHealth(serviceName) {
			utils.Info("[AI] 服务 %s 不健康，跳过", serviceName)
			continue
		}

		mcpTools := as.mcpManager.GetToolRegistry().GetTools(serviceName)
		utils.Info("[AI] 服务 %s 有 %d 个工具", serviceName, len(mcpTools))

		for _, tool := range mcpTools {
			// 转换为OpenAI工具定义格式
			toolDef := ToolDefinition{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			}
			tools = append(tools, toolDef)
		}
	}

	utils.Info("[AI] 获取到 %d 个可用工具", len(tools))
	return tools, nil
}

// validateToolCallParams 验证工具调用参数
func (as *AIService) validateToolCallParams(toolName string, params map[string]interface{}) error {
	if as.mcpManager == nil {
		return fmt.Errorf("MCP管理器未初始化")
	}

	// 查找工具定义
	services := as.mcpManager.ListServices()
	for _, serviceName := range services {
		tools := as.mcpManager.GetToolRegistry().GetTools(serviceName)
		for _, tool := range tools {
			if tool.Name == toolName {
				// 将Tool转换为ToolDefinition
				toolDef := ToolDefinition{
					Name:        tool.Name,
					Description: tool.Description,
					Parameters:  tool.InputSchema,
				}
				return as.validateParams(toolDef, params)
			}
		}
	}

	return fmt.Errorf("未找到工具定义: %s", toolName)
}

// validateParams 验证单个工具的参数
func (as *AIService) validateParams(tool ToolDefinition, params map[string]interface{}) error {
	if tool.Parameters == nil {
		return nil // 没有参数定义，跳过验证
	}

	utils.Debug("[AI] 验证工具 %s 的参数: %+v", tool.Name, params)

	// 检查必需参数
	required := parseRequiredParams(tool.Parameters["required"])

	utils.Debug("[AI] 工具 %s 的必需参数: %v", tool.Name, required)

	// 验证所有必需参数是否都提供了
	for _, reqParam := range required {
		if _, exists := params[reqParam]; !exists {
			return fmt.Errorf("缺少必需参数: %s (工具: %s)", reqParam, tool.Name)
		}
		if params[reqParam] == nil || params[reqParam] == "" {
			return fmt.Errorf("必需参数 %s 不能为空 (工具: %s)", reqParam, tool.Name)
		}
	}

	// 验证参数类型（如果有定义）
	if properties, ok := tool.Parameters["properties"].(map[string]interface{}); ok {
		for paramName, paramValue := range params {
			if propDef, exists := properties[paramName]; exists {
				if err := as.validateParamType(paramName, paramValue, propDef); err != nil {
					return err
				}
			}
		}
	}

	utils.Debug("[AI] 工具 %s 参数验证通过", tool.Name)
	return nil
}

// validateParamType 验证参数类型
func (as *AIService) validateParamType(paramName string, value interface{}, propDef interface{}) error {
	// 这里可以添加更复杂的类型验证逻辑
	// 目前只做基本的非空验证
	if value == nil {
		return fmt.Errorf("参数 %s 不能为 null", paramName)
	}

	if str, ok := value.(string); ok && str == "" {
		return fmt.Errorf("参数 %s 不能为空字符串", paramName)
	}

	return nil
}

// CallTool 调用指定的MCP工具
func (as *AIService) CallTool(toolName string, params map[string]interface{}) (*ToolCallResult, error) {
	if as.mcpManager == nil {
		return nil, fmt.Errorf("MCP管理器未初始化")
	}

	utils.Info("[AI] 调用工具: %s, 参数数量: %d", toolName, len(params))

	// 验证工具调用参数
	if err := as.validateToolCallParams(toolName, params); err != nil {
		utils.Warn("[AI] 工具参数验证失败: %v", err)
		return &ToolCallResult{
			ToolName: toolName,
			Error:    err.Error(),
		}, err
	}

	// 查找包含该工具的服务
	services := as.mcpManager.ListServices()
	for _, serviceName := range services {
		tools := as.mcpManager.GetToolRegistry().GetTools(serviceName)
		for _, tool := range tools {
			if tool.Name == toolName {
				// 调用工具
				result, err := as.mcpManager.CallTool(serviceName, toolName, params)
				if err != nil {
					utils.Error("[AI] 工具调用失败: %v", err)
					return &ToolCallResult{
						ToolName: toolName,
						Error:    err.Error(),
					}, err
				}

				utils.Info("[AI] 工具调用成功: %s", toolName)
				return &ToolCallResult{
					ToolName: toolName,
					Result:   result,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("未找到工具: %s", toolName)
}

// GenerateTextWithTools 使用工具的文本生成
func (as *AIService) GenerateTextWithTools(prompt string, options ...ChatOption) (string, error) {
	traceID := generateAITraceID()
	utils.Info("[AI][%s] GenerateTextWithTools 开始: prompt_runes=%d, options_count=%d, prompt_preview=%q",
		traceID, len([]rune(prompt)), len(options), previewRunesForLog(strings.TrimSpace(prompt), 120))

	// 获取可用工具
	tools, err := as.GetAvailableTools()
	if err != nil {
		utils.Warn("[AI][%s] 获取工具失败，使用普通生成: %v", traceID, err)
		return as.generatePlainTextWithFallback(traceID, prompt, options...)
	}

	if len(tools) == 0 {
		utils.Info("[AI][%s] 没有可用工具，使用普通生成", traceID)
		return as.generatePlainTextWithFallback(traceID, prompt, options...)
	}

	// 由模型进行工具路由判定，避免关键词误判
	shouldUseTools, routeErr := as.shouldUseToolsByModel(traceID, prompt, tools)
	if routeErr != nil {
		utils.Warn("[AI][%s] 工具路由判定失败，回退普通生成: %v", traceID, routeErr)
		return as.generatePlainTextWithFallback(traceID, prompt, options...)
	}
	if !shouldUseTools {
		utils.Info("[AI][%s] 模型判定无需工具，使用普通生成", traceID)
		return as.generatePlainTextWithFallback(traceID, prompt, options...)
	}

	utils.Info("[AI] === 新方案：将工具定义移到用户提示词中 ===")

	// 从数据库获取工具系统提示词
	utils.Info("[AI] 开始获取系统提示词，类型: %s", entity.PromptTypeToolSystem)
	systemPrompt, err := as.promptService.RenderSystemPromptByType(entity.PromptTypeToolSystem, nil)
	if err != nil {
		utils.Info("[AI] 获取系统提示词失败，使用默认提示词: %v", err)
		// 如果获取失败，使用默认提示词
		systemPrompt = `你叫 老九助手，你是一个充满智慧的辅助专家，可以回答用户的各种问题问题，并且可以调用各种mcp工具为用户获取更加专业的回答。

重要规则：
1. 当用户的问题需要使用工具才能获得准确信息时，你必须调用相应的工具
2. 不要猜测或编造信息，对于需要实时数据或外部验证的问题，必须使用工具
3. 调用工具后，根据工具返回的结果给用户准确的回答
4. 调用工具时，必须提供所有必需的参数，不要省略任何 required 参数
5. 根据工具的参数定义和用户的问题，智能选择合适的参数值
6. 如果工具返回错误或无效结果，可以尝试调整参数或尝试其他相关工具

工具调用格式要求：

【主要格式 - JSON格式】
- 推荐格式：<工具名称: {"参数名": "参数值"}>
- 支持跨行格式：<工具名称
: {"参数名": "参数值"}>

【重要约束】
- 必须提供所有必需的参数
- 确保工具名称与可用工具列表中的名称完全一致
- JSON格式的参数值必须用双引号包裹
- 根据用户问题的具体需求，选择最合适的参数值
- 时间格式建议：用户问"今天几号"用"YYYY-MM-DD"，问"现在几点"用"HH:mm:ss"

工具选择原则：
1. 仔细分析用户问题，选择最相关的工具
2. 如果多个工具相关，选择最具体的工具
3. 如果不知道使用哪个工具，可以向用户询问更多细节
4. 对于复杂任务，可以按顺序调用多个工具

响应格式：
1. 直接调用工具，使用上述格式
2. 工具返回结果后，总结或直接展示结果
3. 如果结果需要进一步分析或处理，可以进行解释
4. 保持回答简洁但完整`
	} else {
		utils.Info("[AI] 成功获取系统提示词，长度: %d", len(systemPrompt))
	}

	// 生成工具信息的自然语言描述
	toolsDescription := getToolsAsNaturalLanguage(tools)
	utils.Info("[AI][%s] 模型判定需要工具，生成工具描述，长度: %d", traceID, len(toolsDescription))
	// 组合用户提示词：工具描述 + 用户问题
	fullUserPrompt := toolsDescription + fmt.Sprintf("\n用户问题：%s\n\n请根据用户的问题使用相应的工具来获取准确信息并回答。", prompt)

	// 创建消息（不包含functions参数）
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fullUserPrompt,
		},
	}

	// 关键提示词信息调试（保留用于验证提示词使用情况）
	utils.Debug("[AI] 请求调试摘要")
	utils.Debug("[AI] 模型: %s", as.client.GetModel())
	utils.Debug("[AI] 消息数量: %d", len(messages))
	utils.Debug("[AI] 系统提示词长度: %d 字符", len(systemPrompt))
	utils.Debug("[AI] 完整用户提示词长度: %d 字符", len(fullUserPrompt))
	utils.Debug("[AI] 可用工具数量: %d", len(tools))
	for i, tool := range tools {
		utils.Debug("[AI] 工具 %d: %s", i+1, tool.Name)
	}

	utils.Info("[AI] 发送请求到 AI（新方案：不使用functions参数）")

	// 调用AI（不传递functions参数）
	resp, err := as.client.Chat(messages, options...)
	if err != nil {
		utils.Info("[AI] AI 调用失败: %v", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		utils.Info("[AI] AI 未返回任何内容")
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	choice := resp.Choices[0]
	utils.Info("[AI][%s] AI 返回结果，FinishReason: %s", traceID, resp.Choices[0].FinishReason)

	content := strings.TrimSpace(choice.Message.Content)
	utils.Debug("[AI][%s] 响应调试摘要: choices=%d, finish_reason=%s, content_length=%d, usage(prompt=%d completion=%d total=%d)",
		traceID, len(resp.Choices), choice.FinishReason, len(content), resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)

	// 某些模型在 token 耗尽时可能返回 finish_reason=length 但 content 为空，触发一次简化重试
	if content == "" {
		return as.retryPlainResponseForEmptyContent(traceID, prompt, choice.FinishReason)
	}

	// 检查响应内容中是否包含工具调用标记
	if content != "" {
		utils.Info("[AI] 检查响应内容中的工具调用标记")

		// 创建工具名称集合用于快速查找
		toolNameSet := make(map[string]bool)
		for _, tool := range tools {
			toolNameSet[tool.Name] = true
		}

		toolCalls := parseToolCallsFromContent(content, toolNameSet)
		if len(toolCalls) > 0 {
			utils.Info("[AI] 从响应内容中解析到 %d 个工具调用", len(toolCalls))

			// 处理所有工具调用
			for _, toolCall := range toolCalls {
				utils.Info("[AI] 调用工具: %s, 参数数量: %d", toolCall.Name, len(toolCall.Params))

				// 调用工具
				toolResult, err := as.CallTool(toolCall.Name, toolCall.Params)
				if err != nil {
					utils.Info("[AI] 工具调用失败: %v", err)
					return "", fmt.Errorf("工具调用失败: %v", err)
				}

				utils.Info("[AI] 工具调用成功: %s", toolCall.Name)

				// 将工具结果添加到对话中
				messages = append(messages,
					openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: fmt.Sprintf("<%s/>", toolCall.Name),
					},
					openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Content: fmt.Sprintf("工具 %s 的返回结果：%v", toolCall.Name, toolResult.Result),
					},
				)
			}

			// 再次调用AI处理工具结果
			utils.Info("[AI] 再次调用 AI 处理工具结果")
			resp, err = as.client.Chat(messages, options...)
			if err != nil {
				utils.Info("[AI] AI 处理工具结果失败: %v", err)
				return "", err
			}

			if len(resp.Choices) == 0 {
				utils.Info("[AI] AI 处理工具结果后未返回任何内容")
				return "", fmt.Errorf("AI 处理工具结果后未返回任何内容")
			}

			utils.Info("[AI] AI 处理工具结果成功")
			utils.Debug("[AI][%s] 二次响应usage: prompt=%d completion=%d total=%d",
				traceID, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
			finalContent := strings.TrimSpace(resp.Choices[0].Message.Content)
			if finalContent == "" {
				return as.retryPlainResponseForEmptyContent(traceID, prompt, resp.Choices[0].FinishReason)
			}
			return finalContent, nil
		}
	}

	// 清理响应内容中的工具调用标记
	cleanContent := strings.TrimSpace(cleanToolCallMarkers(content))
	if cleanContent != content {
		utils.Info("[AI] 清理了工具调用标记")
	}
	if cleanContent == "" {
		return as.retryPlainResponseForEmptyContent(traceID, prompt, choice.FinishReason)
	}

	utils.Info("[AI] AI 没有调用工具，直接返回内容")
	return cleanContent, nil
}

func (as *AIService) generatePlainTextWithFallback(traceID string, prompt string, options ...ChatOption) (string, error) {
	systemPrompt := "你是一个有用的 AI 助手，擅长理解和回答各种问题。请提供准确、有帮助的回答。"
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	resp, err := as.client.Chat(messages, options...)
	if err != nil {
		if isLikelyResourceTitleTask(prompt) && isTimeoutLikeError(err) {
			if fallbackTitle, ok := buildLocalTitleFallback(prompt); ok {
				utils.Warn("[AI][%s] 普通生成超时，使用本地标题回退: %q", traceID, fallbackTitle)
				return fallbackTitle, nil
			}
		}
		return "", err
	}

	if len(resp.Choices) == 0 {
		utils.Warn("[AI][%s] 普通生成返回 choices=0，触发空内容重试", traceID)
		return as.retryPlainResponseForEmptyContent(traceID, prompt, openai.FinishReason("no_choice"))
	}

	first := resp.Choices[0]
	content := strings.TrimSpace(first.Message.Content)
	if content == "" {
		return as.retryPlainResponseForEmptyContent(traceID, prompt, first.FinishReason)
	}

	utils.Debug("[AI][%s] 普通生成成功: finish_reason=%s, content_length=%d, usage(prompt=%d completion=%d total=%d)",
		traceID, first.FinishReason, len(content), resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	return content, nil
}

func (as *AIService) shouldUseToolsByModel(traceID string, prompt string, tools []ToolDefinition) (bool, error) {
	if len(tools) == 0 {
		return false, nil
	}

	toolNames := make([]string, 0, len(tools))
	for _, tool := range tools {
		toolNames = append(toolNames, tool.Name)
	}

	routerSystemPrompt := `你是工具路由器。你只做一件事：判断是否需要调用外部工具。
仅返回 NEED_TOOL 或 NO_TOOL，禁止输出其他任何内容。
判断标准：
- 需要外部实时信息、网页抓取、搜索、时间/天气/新闻等 -> NEED_TOOL
- 文本改写、内容生成、标题优化、摘要、分类等纯语言任务 -> NO_TOOL`

	routerUserPrompt := fmt.Sprintf("用户输入：%s\n可用工具：%s",
		prompt, strings.Join(toolNames, ", "))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: routerSystemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: routerUserPrompt,
		},
	}

	resp, err := as.client.Chat(messages, WithMaxTokens(32), WithTemperature(0))
	if err != nil {
		return false, err
	}
	if len(resp.Choices) == 0 {
		return false, fmt.Errorf("工具路由无返回结果")
	}

	decisionRaw := strings.ToUpper(strings.TrimSpace(resp.Choices[0].Message.Content))
	utils.Info("[AI][%s] 模型工具路由结果: %q (finish_reason=%s, usage_total=%d)",
		traceID, decisionRaw, resp.Choices[0].FinishReason, resp.Usage.TotalTokens)
	if decisionRaw == "" {
		utils.Warn("[AI][%s] 模型工具路由返回空结果，默认 NO_TOOL", traceID)
		return false, nil
	}

	switch {
	case strings.HasPrefix(decisionRaw, "NEED_TOOL"):
		return true, nil
	case strings.HasPrefix(decisionRaw, "NO_TOOL"):
		return false, nil
	default:
		utils.Warn("[AI][%s] 模型工具路由结果无法解析，默认 NO_TOOL: %q", traceID, decisionRaw)
		return false, nil
	}
}

// retryPlainResponseForEmptyContent 当模型返回空内容时，用简化提示词重试一次，避免返回 200 + 空字符串
func (as *AIService) retryPlainResponseForEmptyContent(traceID string, prompt string, finishReason openai.FinishReason) (string, error) {
	utils.Warn("[AI][%s] AI 返回空内容，触发简化重试，finish_reason=%s", traceID, finishReason)

	sanitizedPrompt := strings.TrimSpace(prompt)
	if sanitizedPrompt == "" {
		return "", fmt.Errorf("AI 返回空内容，且用户提示词为空")
	}

	// 防止超长提示词在部分兼容模型中触发 finish_reason=length 且 content 为空
	sanitizedPrompt = truncateRunes(sanitizedPrompt, 2000)

	attempts := []struct {
		name    string
		options []ChatOption
	}{
		{
			name:    "中等max_tokens重试",
			options: []ChatOption{WithMaxTokens(256), WithTemperature(0.2)},
		},
		{
			name:    "较高max_tokens重试",
			options: []ChatOption{WithMaxTokens(512), WithTemperature(0.1)},
		},
	}

	lastFinishReason := finishReason
	var lastErr error

	for _, attempt := range attempts {
		utils.Info("[AI][%s] 空内容重试开始[%s]: prompt_runes=%d", traceID, attempt.name, len([]rune(sanitizedPrompt)))

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个有用的 AI 助手。请直接给出最终答案，不要输出工具调用标记或思考过程。",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("请直接回答以下问题，控制在200字以内：\n%s", sanitizedPrompt),
			},
		}

		resp, err := as.client.Chat(messages, attempt.options...)
		if err != nil {
			lastErr = err
			utils.Warn("[AI][%s] 空内容重试失败[%s]: %v", traceID, attempt.name, err)
			if isLikelyResourceTitleTask(sanitizedPrompt) && isTimeoutLikeError(err) {
				if fallbackTitle, ok := buildLocalTitleFallback(sanitizedPrompt); ok {
					utils.Warn("[AI][%s] 标题任务重试超时，立即使用本地标题回退: %q", traceID, fallbackTitle)
					return fallbackTitle, nil
				}
			}
			continue
		}
		if len(resp.Choices) == 0 {
			lastErr = fmt.Errorf("无返回结果")
			utils.Warn("[AI][%s] 空内容重试失败[%s]: choices为空", traceID, attempt.name)
			continue
		}

		lastFinishReason = resp.Choices[0].FinishReason
		retryContent := strings.TrimSpace(resp.Choices[0].Message.Content)
		utils.Info("[AI][%s] 空内容重试结果[%s]: finish_reason=%s, content_length=%d, usage(prompt=%d completion=%d total=%d)",
			traceID, attempt.name, lastFinishReason, len(retryContent), resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)

		if retryContent != "" {
			utils.Info("[AI][%s] 空内容重试成功[%s]", traceID, attempt.name)
			return retryContent, nil
		}
	}

	if lastErr != nil {
		if fallbackTitle, ok := buildLocalTitleFallback(sanitizedPrompt); ok {
			utils.Warn("[AI][%s] AI重试失败，使用本地标题回退: %q", traceID, fallbackTitle)
			return fallbackTitle, nil
		}
		return "", fmt.Errorf("AI 返回空内容，重试后仍失败（finish_reason=%s）: %w", lastFinishReason, lastErr)
	}

	if fallbackTitle, ok := buildLocalTitleFallback(sanitizedPrompt); ok {
		utils.Warn("[AI][%s] AI重试仍为空，使用本地标题回退: %q", traceID, fallbackTitle)
		return fallbackTitle, nil
	}

	return "", fmt.Errorf("AI 返回空内容，重试后仍为空（finish_reason=%s）", lastFinishReason)
}

func generateAITraceID() string {
	return fmt.Sprintf("t_%d", time.Now().UnixNano())
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes])
}

func buildLocalTitleFallback(prompt string) (string, bool) {
	if !isLikelyResourceTitleTask(prompt) {
		return "", false
	}

	description := parseResourceFieldValueAny(prompt, []string{"资源描述", "现有描述"})
	if titleFromDesc := extractQuotedTitleCandidate(description); titleFromDesc != "" {
		return titleFromDesc, true
	}

	tags := parseResourceFieldValueAny(prompt, []string{"资源标签", "标签"})
	if titleFromTags := extractQuotedTitleCandidate(tags); titleFromTags != "" {
		return titleFromTags, true
	}

	rawTitle := parseResourceFieldValueAny(prompt, []string{"资源标题", "资源标题(原始，可能含噪声)", "资源标题(清洗候选)"})
	cleanedTitle := cleanNoisyResourceTitle(rawTitle)
	if cleanedTitle == "" {
		return "", false
	}
	return cleanedTitle, true
}

func isLikelyResourceTitleTask(prompt string) bool {
	normalized := strings.TrimSpace(prompt)
	if normalized == "" {
		return false
	}
	hasTitleField := strings.Contains(normalized, "资源标题:") || strings.Contains(normalized, "资源标题(原始")
	return hasTitleField && strings.Contains(normalized, "生成一个") && strings.Contains(normalized, "标题")
}

func parseResourceFieldValue(prompt string, fieldName string) string {
	fieldPrefix := fieldName + ":"
	lines := strings.Split(prompt, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, fieldPrefix) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, fieldPrefix))
		}
	}
	return ""
}

func parseResourceFieldValueAny(prompt string, fieldNames []string) string {
	for _, fieldName := range fieldNames {
		if value := parseResourceFieldValue(prompt, fieldName); value != "" {
			return value
		}
	}
	return ""
}

func extractQuotedTitleCandidate(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`《([^》]{2,80})》`),
		regexp.MustCompile(`“([^”]{2,80})”`),
		regexp.MustCompile(`"([^"]{2,80})"`),
	}
	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(text)
		if len(match) >= 2 {
			candidate := strings.TrimSpace(match[1])
			if candidate != "" {
				return "《" + candidate + "》"
			}
		}
	}
	return ""
}

func cleanNoisyResourceTitle(title string) string {
	normalized := strings.TrimSpace(title)
	if normalized == "" {
		return ""
	}

	// 常见来源会在标题前附带数字编号，如 044651-xxx
	normalized = regexp.MustCompile(`^\d{3,}-`).ReplaceAllString(normalized, "")
	// 清理夹杂在中文标题中的连续大写噪声，如 NNN/MMM/ZZZ
	normalized = regexp.MustCompile(`[A-Z]{2,}`).ReplaceAllString(normalized, "")
	normalized = regexp.MustCompile(`\s+`).ReplaceAllString(normalized, " ")
	normalized = strings.TrimSpace(normalized)

	// 去除清理后可能出现的前导符号
	normalized = strings.TrimLeft(normalized, "-_.,，。:：;； ")
	return strings.TrimSpace(normalized)
}

func isTimeoutLikeError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "deadline exceeded") || strings.Contains(errMsg, "client.timeout")
}

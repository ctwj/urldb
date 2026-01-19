package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/ai/mcp"
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
	ToolName string                 `json:"tool_name"`
	Result   interface{}            `json:"result"`
	Error    string                 `json:"error,omitempty"`
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
	log.Printf("[GenerateText] 开始处理请求，prompt: %s", prompt)

	// 如果有 MCP 管理器，尝试使用工具增强的生成
	if as.mcpManager != nil {
		log.Printf("[GenerateText] MCP 管理器已初始化，尝试使用工具增强生成")
		result, err := as.GenerateTextWithTools(prompt, options...)
		if err != nil {
			log.Printf("[GenerateText] 工具增强生成失败，回退到普通生成: %v", err)
		} else {
			log.Printf("[GenerateText] 工具增强生成成功")
			return result, nil
		}
	} else {
		log.Printf("[GenerateText] MCP 管理器未初始化，使用普通生成")
	}

	// 使用通用的系统提示词
	systemPrompt := "你是一个有用的 AI 助手，擅长理解和回答各种问题。请提供准确、有帮助的回答。"

	// 创建标准system+user消息结构
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
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	return resp.Choices[0].Message.Content, nil
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
	// 将提示词转换为小写进行匹配
	lowerPrompt := strings.ToLower(prompt)

	// 工具需求关键词
	toolKeywords := []string{
		"时间", "几点", "现在", "今天", "日期", "当前",
		"搜索", "查询", "找", "搜索信息", "google", "百度",
		"网页", "网站", "内容", "获取", "抓取",
		"天气", "温度", "气候", "预报",
		"新闻", "资讯", "动态", "最新",
		"翻译", "英文", "中文", "语言",
		"计算", "换算", "转换", "公式",
		"汇率", "价格", "股票", "金融",
	}

	// 检查是否包含工具相关关键词
	for _, keyword := range toolKeywords {
		if strings.Contains(lowerPrompt, keyword) {
			return true
		}
	}

	// 检查是否是问句（通常需要查询信息）
	questionPatterns := []string{
		"什么", "怎么", "如何", "为什么", "哪里", "哪个", "谁",
		"吗", "呢", "？", "?",
	}

	for _, pattern := range questionPatterns {
		if strings.Contains(lowerPrompt, pattern) {
			return true
		}
	}

	// 检查是否包含数字相关的查询（如时间、日期等）
	if regexp.MustCompile(`\d+`).MatchString(prompt) {
		return true
	}

	return false
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
				var required []interface{}
				if req, ok := tool.Parameters["required"].([]interface{}); ok {
					required = req
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

	log.Printf("[parseToolCallsFromContent] 原始内容: %q", content)

	// 先尝试匹配 JSON 格式的工具调用：<tool_name: {}>
	jsonRe := regexp.MustCompile(`(?s)<(\w+):\s*({[^}]*})>`)
	jsonMatches := jsonRe.FindAllStringSubmatch(content, -1)
	log.Printf("[parseToolCallsFromContent] JSON 格式匹配到 %d 个结果", len(jsonMatches))

	for i, match := range jsonMatches {
		log.Printf("[parseToolCallsFromContent] JSON 匹配 %d: %v", i, match)
		if len(match) < 3 {
			continue
		}

		toolName := match[1]

		// 检查工具名称是否在已注册的工具列表中
		if !toolNameSet[toolName] {
			log.Printf("[parseToolCallsFromContent] 工具 %s 未注册，跳过", toolName)
			continue
		}

		jsonStr := match[2]
		params := make(map[string]interface{})

		if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
			log.Printf("[parseToolCallsFromContent] 解析 JSON 参数失败: %v", err)
			params = map[string]interface{}{"args": jsonStr}
		}

		toolCalls = append(toolCalls, ToolCallFromContent{
			Name:   toolName,
			Params: params,
		})
		log.Printf("[parseToolCallsFromContent] 解析工具: %s, 参数: %v", toolName, params)
	}

	// 如果没有匹配到 JSON 格式，尝试匹配简单标签格式：<tool_name> 或 <tool_name/> 或 <tool_name\n
	if len(toolCalls) == 0 {
		simpleRe := regexp.MustCompile(`<(\w+)[\s\n>]`)
		simpleMatches := simpleRe.FindAllStringSubmatch(content, -1)
		log.Printf("[parseToolCallsFromContent] 简单标签格式匹配到 %d 个结果", len(simpleMatches))

		for i, match := range simpleMatches {
			log.Printf("[parseToolCallsFromContent] 简单标签匹配 %d: %v", i, match)
			if len(match) < 2 {
				continue
			}

			toolName := match[1]

			// 检查工具名称是否在已注册的工具列表中
			if !toolNameSet[toolName] {
				log.Printf("[parseToolCallsFromContent] 工具 %s 未注册，跳过", toolName)
				continue
			}

			toolCalls = append(toolCalls, ToolCallFromContent{
				Name:   toolName,
				Params: map[string]interface{}{},
			})
			log.Printf("[parseToolCallsFromContent] 解析工具: %s, 参数: map[]", toolName)
		}
	}

	// 如果还没有匹配到，尝试匹配 HTML 属性格式：<tool_name param1="value1"/>
	if len(toolCalls) == 0 {
		htmlRe := regexp.MustCompile(`<(\w+)(\s+[^>]*)>`)
		htmlMatches := htmlRe.FindAllStringSubmatch(content, -1)
		log.Printf("[parseToolCallsFromContent] HTML 格式匹配到 %d 个结果", len(htmlMatches))

		for i, match := range htmlMatches {
			log.Printf("[parseToolCallsFromContent] HTML 匹配 %d: %v", i, match)
			if len(match) < 3 {
				continue
			}

			toolName := match[1]

			// 检查工具名称是否在已注册的工具列表中
			if !toolNameSet[toolName] {
				log.Printf("[parseToolCallsFromContent] 工具 %s 未注册，跳过", toolName)
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
			log.Printf("[parseToolCallsFromContent] 解析工具: %s, 参数: %v", toolName, params)
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
			`\d{4}年\d{1,2}月\d{1,2}日`,  // 中文日期格式
			`\d{4}-\d{1,2}-\d{1,2}`,        // 英文日期格式
			`\d{1,2}:\d{2}:\d{2}`,          // 时间格式
			`timestamp:\s*\d+`,             // 时间戳格式
		}

		for _, pattern := range timePatterns {
			if matched, _ := regexp.MatchString(pattern, content); matched {
				hasResult = true
				log.Printf("[parseToolCallsFromContent] 检测到时间数据: %s", pattern)
				break
			}
		}

		if hasResult {
			log.Printf("[parseToolCallsFromContent] 检测到响应中已包含工具结果，忽略工具调用")
			return []ToolCallFromContent{}
		}

		// 检查响应长度，如果很短（比如只有工具调用标记），说明没有工具结果
		// 去除工具调用标记后的内容长度
		cleanContent := regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")
		cleanContent = strings.TrimSpace(cleanContent)
		if len(cleanContent) < 10 {
			log.Printf("[parseToolCallsFromContent] 响应内容过短，没有工具结果")
		} else {
			log.Printf("[parseToolCallsFromContent] 响应内容长度: %d", len(cleanContent))
		}
	}

	return toolCalls
}

// cleanToolCallMarkers 清理响应内容中的工具调用标记
func cleanToolCallMarkers(content string) string {
	// 移除工具调用标记：<tool_name>...</tool_name> 或 <tool_name/> 或 <tool_name: {}> 等
	// 也支持没有闭合标签的格式：<tool_name\n⟶
	re := regexp.MustCompile(`<\w+(?::\s*{[^}]*})?\s*/?>\s*</\w+>|<\w+(?::\s*{[^}]*})?\s*/?>|<\w+>|<\w+[\s\n]`)
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

	log.Printf("[GetAvailableTools] 检查 %d 个服务", len(services))

	for _, serviceName := range services {
		// 检查服务健康状态
		if !as.mcpManager.CheckServiceHealth(serviceName) {
			log.Printf("[GetAvailableTools] 服务 %s 不健康，跳过", serviceName)
			continue
		}

		mcpTools := as.mcpManager.GetToolRegistry().GetTools(serviceName)
		log.Printf("[GetAvailableTools] 服务 %s 有 %d 个工具", serviceName, len(mcpTools))

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

	log.Printf("[GetAvailableTools] 获取到 %d 个可用工具", len(tools))
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

	log.Printf("[validateParams] 验证工具 %s 的参数: %+v", tool.Name, params)

	// 检查必需参数
	required := []string{}
	if reqArray, ok := tool.Parameters["required"].([]interface{}); ok {
		for _, req := range reqArray {
			if reqStr, ok := req.(string); ok {
				required = append(required, reqStr)
			}
		}
	}

	log.Printf("[validateParams] 工具 %s 的必需参数: %v", tool.Name, required)

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

	log.Printf("[validateParams] 工具 %s 参数验证通过", tool.Name)
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

	log.Printf("调用工具: %s, 参数: %+v", toolName, params)

	// 验证工具调用参数
	if err := as.validateToolCallParams(toolName, params); err != nil {
		log.Printf("工具参数验证失败: %v", err)
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
					log.Printf("工具调用失败: %v", err)
					return &ToolCallResult{
						ToolName: toolName,
						Error:    err.Error(),
					}, err
				}

				log.Printf("工具调用成功: %s", toolName)
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
	// 获取可用工具
	tools, err := as.GetAvailableTools()
	if err != nil {
		log.Printf("获取工具失败，使用普通生成: %v", err)
		// 直接使用 OpenAI 客户端生成，避免循环调用
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
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("AI 未返回任何内容")
		}
		return resp.Choices[0].Message.Content, nil
	}

	if len(tools) == 0 {
		log.Printf("没有可用工具，使用普通生成")
		// 直接使用 OpenAI 客户端生成，避免循环调用
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
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("AI 未返回任何内容")
		}
		return resp.Choices[0].Message.Content, nil
	}

	log.Printf("[GenerateTextWithTools] === 新方案：将工具定义移到用户提示词中 ===")

	// 从数据库获取工具系统提示词
	log.Printf("[GenerateTextWithTools] 开始获取系统提示词，类型: %s", entity.PromptTypeToolSystem)
	systemPrompt, err := as.promptService.RenderSystemPromptByType(entity.PromptTypeToolSystem, nil)
	if err != nil {
		log.Printf("[GenerateTextWithTools] 获取系统提示词失败，使用默认提示词: %v", err)
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
		log.Printf("[GenerateTextWithTools] 成功获取系统提示词，长度: %d", len(systemPrompt))
	}

	// 智能判断是否需要工具描述
	var fullUserPrompt string
	if needsTools(prompt) {
		// 生成工具信息的自然语言描述
		toolsDescription := getToolsAsNaturalLanguage(tools)
		log.Printf("[GenerateTextWithTools] 检测到工具需求，生成工具描述，长度: %d", len(toolsDescription))
		// 组合用户提示词：工具描述 + 用户问题
		fullUserPrompt = toolsDescription + fmt.Sprintf("\n用户问题：%s\n\n请根据用户的问题使用相应的工具来获取准确信息并回答。", prompt)
	} else {
		log.Printf("[GenerateTextWithTools] 未检测到工具需求，使用简洁提示词")
		// 简洁的用户提示词，不包含工具描述
		fullUserPrompt = fmt.Sprintf("用户问题：%s\n\n请直接回答用户的问题。", prompt)
	}

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

	// ===== 完整的AI接口请求数据调试日志 =====
	log.Printf("=== [GenerateTextWithTools] 完整AI接口请求数据（新方案） ===")

	// 1. 打印完整的请求结构（不包含functions）
	requestData := map[string]interface{}{
		"model":    as.client.GetModel(),
		"messages": messages,
	}

	if requestJSON, err := json.MarshalIndent(requestData, "", "  "); err == nil {
		log.Printf("完整OpenAI请求JSON（新方案）:\n%s", string(requestJSON))
	} else {
		log.Printf("序列化请求JSON失败: %v", err)
	}

	// 2. 分别打印各个部分以便调试
	log.Printf("--- 系统提示词完整内容 ---")
	log.Printf("%s", systemPrompt)
	log.Printf("--- 用户提示词完整内容 ---")
	log.Printf("%s", fullUserPrompt)
	// 只在需要时显示工具描述
	if needsTools(prompt) {
		log.Printf("--- 工具自然语言描述 ---")
		log.Printf("%s", getToolsAsNaturalLanguage(tools))
	}
	log.Printf("========================================")

	// 关键提示词信息调试（保留用于验证提示词使用情况）
	log.Printf("=== [GenerateTextWithTools] 提示词调试信息（新方案） ===")
	log.Printf("用户原始输入: %q", prompt)
	log.Printf("系统提示词长度: %d 字符", len(systemPrompt))
	log.Printf("完整用户提示词长度: %d 字符", len(fullUserPrompt))
	log.Printf("可用工具数量: %d", len(tools))
	for i, tool := range tools {
		log.Printf("工具 %d: %s", i+1, tool.Name)
	}
	log.Printf("===========================================")

	log.Printf("[GenerateTextWithTools] 发送请求到 AI（新方案：不使用functions参数）")

	// 调用AI（不传递functions参数）
	resp, err := as.client.Chat(messages, options...)
	if err != nil {
		log.Printf("[GenerateTextWithTools] AI 调用失败: %v", err)
		return "", err
	}

	if len(resp.Choices) == 0 {
		log.Printf("[GenerateTextWithTools] AI 未返回任何内容")
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	choice := resp.Choices[0]
	log.Printf("[GenerateTextWithTools] AI 返回结果，FinishReason: %s", resp.Choices[0].FinishReason)

	// ===== 调试信息打印 - 完整的AI响应数据 =====
	log.Printf("=== [GenerateTextWithTools] 完整AI响应数据（新方案） ===")
	if responseJSON, err := json.MarshalIndent(resp, "", "  "); err == nil {
		log.Printf("完整AI响应JSON:\n%s", string(responseJSON))
	} else {
		log.Printf("序列化响应JSON失败: %v", err)
	}
	log.Printf("===========================================")

	// 检查响应内容中是否包含工具调用标记
	content := choice.Message.Content
	if content != "" {
		log.Printf("[GenerateTextWithTools] 检查响应内容中的工具调用标记")

		// 创建工具名称集合用于快速查找
		toolNameSet := make(map[string]bool)
		for _, tool := range tools {
			toolNameSet[tool.Name] = true
		}

		toolCalls := parseToolCallsFromContent(content, toolNameSet)
		if len(toolCalls) > 0 {
			log.Printf("[GenerateTextWithTools] 从响应内容中解析到 %d 个工具调用", len(toolCalls))

			// 处理所有工具调用
			for _, toolCall := range toolCalls {
				log.Printf("[GenerateTextWithTools] 调用工具: %s, 参数: %v", toolCall.Name, toolCall.Params)

				// 调用工具
				toolResult, err := as.CallTool(toolCall.Name, toolCall.Params)
				if err != nil {
					log.Printf("[GenerateTextWithTools] 工具调用失败: %v", err)
					return "", fmt.Errorf("工具调用失败: %v", err)
				}

				log.Printf("[GenerateTextWithTools] 工具调用成功，结果: %v", toolResult.Result)

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
			log.Printf("[GenerateTextWithTools] 再次调用 AI 处理工具结果")
			resp, err = as.client.Chat(messages, options...)
			if err != nil {
				log.Printf("[GenerateTextWithTools] AI 处理工具结果失败: %v", err)
				return "", err
			}

			if len(resp.Choices) == 0 {
				log.Printf("[GenerateTextWithTools] AI 处理工具结果后未返回任何内容")
				return "", fmt.Errorf("AI 处理工具结果后未返回任何内容")
			}

			log.Printf("[GenerateTextWithTools] AI 处理工具结果成功")
			return resp.Choices[0].Message.Content, nil
		}
	}

	// 清理响应内容中的工具调用标记
	cleanContent := cleanToolCallMarkers(content)
	if cleanContent != content {
		log.Printf("[GenerateTextWithTools] 清理了工具调用标记")
	}

	log.Printf("[GenerateTextWithTools] AI 没有调用工具，直接返回内容")
	return cleanContent, nil
}
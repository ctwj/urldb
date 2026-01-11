package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ctwj/urldb/db/repo"
	"github.com/ctwj/urldb/pkg/ai/mcp"
	"github.com/sashabaranov/go-openai"
)

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

	return &AIService{
		client:      client,
		contentGen:  contentGen,
		classifier:  classifier,
		repoManager: repoManager,
	}, nil
}

// NewAIService 创建AI服务
func NewAIService(client *OpenAIClient, repoManager *repo.RepositoryManager) (*AIService, error) {
	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)

	return &AIService{
		client:      client,
		contentGen:  contentGen,
		classifier:  classifier,
		repoManager: repoManager,
	}, nil
}

// NewAIServiceWithMCP 创建支持MCP的AI服务
func NewAIServiceWithMCP(client *OpenAIClient, repoManager *repo.RepositoryManager, mcpManager *mcp.MCPManager) (*AIService, error) {
	contentGen := NewContentGenerator(client, repoManager)
	classifier := NewClassifier(client, repoManager)

	return &AIService{
		client:      client,
		contentGen:  contentGen,
		classifier:  classifier,
		repoManager: repoManager,
		mcpManager:  mcpManager,
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

	// 普通文本生成（无工具）
	messages := []openai.ChatCompletionMessage{
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
		simpleRe := regexp.MustCompile(`<(\w+)(?=[\s\n>]|$)`)
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
	re := regexp.MustCompile(`(?s)<\w+(?::\s*{[^}]*})?\s*/?>\s*</\w+>|<\w+(?::\s*{[^}]*})?\s*/?>|<\w+>|<\w+(?=[\s\n]|$)`)
	cleanContent := re.ReplaceAllString(content, "")

	// 清理多余的空行
	cleanContent = regexp.MustCompile(`\n\s*\n\s*\n`).ReplaceAllString(cleanContent, "\n\n")

	// 去除首尾空白
	cleanContent = strings.TrimSpace(cleanContent)

	return cleanContent
}

// AskQuestion 通用问答 - 供其他模块调用
func (as *AIService) AskQuestion(question string, context string) (string, error) {
	prompt := fmt.Sprintf("根据以下上下文回答问题：\n\n上下文：%s\n\n问题：%s\n\n请基于提供的上下文信息给出准确的回答。", context, question)
	return as.GenerateText(prompt, WithMaxTokens(500), WithTemperature(0.7))
}

// AnalyzeText 通用文本分析 - 供其他模块调用
func (as *AIService) AnalyzeText(text string, analysisType string) (string, error) {
	prompt := fmt.Sprintf("请对以下文本进行%s分析：\n\n%s", analysisType, text)
	return as.GenerateText(prompt, WithMaxTokens(300), WithTemperature(0.5))
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

// CallTool 调用指定的MCP工具
func (as *AIService) CallTool(toolName string, params map[string]interface{}) (*ToolCallResult, error) {
	if as.mcpManager == nil {
		return nil, fmt.Errorf("MCP管理器未初始化")
	}

	log.Printf("调用工具: %s, 参数: %+v", toolName, params)

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
		return as.GenerateText(prompt, options...)
	}

	if len(tools) == 0 {
		log.Printf("没有可用工具，使用普通生成")
		return as.GenerateText(prompt, options...)
	}

	// 创建OpenAI函数调用格式
	var functions []openai.FunctionDefinition
	for _, tool := range tools {
		functionDef := openai.FunctionDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.Parameters,
		}
		functions = append(functions, functionDef)
	}

	// 创建消息，添加系统提示告诉 AI 必须使用工具
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `你是一个有用的 AI 助手。你可以使用提供的工具来获取信息并回答问题。

重要规则：
1. 如果用户询问时间、日期、搜索信息或其他需要实时数据的问题，你必须使用相应的工具
2. 不要猜测或编造信息，必须使用工具获取准确的数据
3. 调用工具后，根据工具返回的结果给用户准确的回答

工具调用格式要求：
- 当需要调用工具时，请使用以下格式：<工具名称: {}>
- 示例：<current_time: {}> 或 <search: {"query": "北京现在几点"}>
- 不要使用其他格式，如 <工具名称> 或 <工具名称/> 等
- 确保工具名称与可用工具列表中的名称完全一致

可用工具列表：` + getToolListSummary(tools),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// 添加工具调用选项
	toolOptions := append(options, WithFunctions(functions))

	log.Printf("[GenerateTextWithTools] 发送请求到 AI，包含 %d 个工具", len(functions))

	// 调用AI
	resp, err := as.client.Chat(messages, toolOptions...)
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

	// 检查是否有函数调用
	if choice.Message.FunctionCall != nil {
		toolCall := choice.Message.FunctionCall
		log.Printf("[GenerateTextWithTools] AI 决定调用函数: %s, 参数: %s", toolCall.Name, toolCall.Arguments)

		// 解析参数
		var params map[string]interface{}
		if toolCall.Arguments != "" {
			if err := json.Unmarshal([]byte(toolCall.Arguments), &params); err != nil {
				return "", fmt.Errorf("解析工具参数失败: %v", err)
			}
		}

		// 调用工具
		log.Printf("[GenerateTextWithTools] 开始调用工具: %s", toolCall.Name)
		toolResult, err := as.CallTool(toolCall.Name, params)
		if err != nil {
			log.Printf("[GenerateTextWithTools] 工具调用失败: %v", err)
			return "", fmt.Errorf("工具调用失败: %v", err)
		}

		log.Printf("[GenerateTextWithTools] 工具调用成功，结果: %v", toolResult.Result)

		// 将工具结果添加到对话中
		messages = append(messages,
			openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleAssistant,
				Content: "",
				FunctionCall: &openai.FunctionCall{
					Name:      toolCall.Name,
					Arguments: toolCall.Arguments,
				},
			},
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleFunction,
				Name:    toolCall.Name,
				Content: fmt.Sprintf("%v", toolResult.Result),
			},
		)

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

	// 检查响应内容中是否包含工具调用标记（GLM 特有的格式）
	content := choice.Message.Content
	if content != "" {
		log.Printf("[GenerateTextWithTools] 检查响应内容中的工具调用标记")

		// 获取所有已注册的工具名称列表
		availableTools, err := as.GetAvailableTools()
		if err != nil {
			log.Printf("[GenerateTextWithTools] 获取可用工具失败: %v", err)
		}

		// 创建工具名称集合用于快速查找
		toolNameSet := make(map[string]bool)
		for _, tool := range availableTools {
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
						Role:    openai.ChatMessageRoleFunction,
						Name:    toolCall.Name,
						Content: fmt.Sprintf("%v", toolResult.Result),
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
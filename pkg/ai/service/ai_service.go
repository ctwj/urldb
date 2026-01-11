package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	for _, serviceName := range services {
		mcpTools := as.mcpManager.GetToolRegistry().GetTools(serviceName)
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

	log.Printf("获取到 %d 个可用工具", len(tools))
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

	// 创建消息
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	// 添加工具调用选项
	toolOptions := append(options, WithFunctions(functions))

	// 调用AI
	resp, err := as.client.Chat(messages, toolOptions...)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("AI 未返回任何内容")
	}

	choice := resp.Choices[0]

	// 检查是否有工具调用
	if choice.Message.FunctionCall != nil {
		toolCall := choice.Message.FunctionCall
		log.Printf("AI决定调用工具: %s", toolCall.Name)

		// 解析参数
		var params map[string]interface{}
		if toolCall.Arguments != "" {
			if err := json.Unmarshal([]byte(toolCall.Arguments), &params); err != nil {
				return "", fmt.Errorf("解析工具参数失败: %v", err)
			}
		}

		// 调用工具
		toolResult, err := as.CallTool(toolCall.Name, params)
		if err != nil {
			return "", fmt.Errorf("工具调用失败: %v", err)
		}

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
		resp, err = as.client.Chat(messages, options...)
		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("AI 处理工具结果后未返回任何内容")
		}

		return resp.Choices[0].Message.Content, nil
	}

	return choice.Message.Content, nil
}
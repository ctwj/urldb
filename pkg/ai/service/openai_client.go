package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ctwj/urldb/utils"
	"github.com/sashabaranov/go-openai"
)

// OpenAIClient 包装 OpenAI 客户端
type OpenAIClient struct {
	apiKey       string
	baseURL      string
	model        string
	organization string
	proxy        string
	timeout      time.Duration
	retryCount   int
	client       *openai.Client
	config       ConfigManager // 添加配置管理器
}

// NewOpenAIClientWithConfig 创建新的 OpenAI 客户端
func NewOpenAIClientWithConfig(configManager ConfigManager) (*OpenAIClient, error) {
	// 从数据库配置系统读取参数
	apiKey, err := configManager.GetConfigString("ai_api_key")
	if err != nil || apiKey == "" {
		// 尝试从环境变量读取
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("AI API Key 未配置")
		}
	}

	baseURL, _ := configManager.GetConfigString("ai_api_url")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	model, _ := configManager.GetConfigString("ai_model")
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	organization, _ := configManager.GetConfigString("ai_organization")
	proxy, _ := configManager.GetConfigString("ai_proxy")

	timeoutSec, _ := configManager.GetConfigInt("ai_timeout")
	if timeoutSec == 0 {
		timeoutSec = 30
	}

	retryCount, _ := configManager.GetConfigInt("ai_retry_count")
	if retryCount == 0 {
		retryCount = 3
	}

	// 构建 OpenAI 客户端配置
	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = baseURL
	if organization != "" {
		clientConfig.OrgID = organization
	}

	// 设置代理
	if proxy != "" {
		clientConfig.HTTPClient = &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return url.Parse(proxy)
				},
			},
		}
	} else {
		clientConfig.HTTPClient = &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
		}
	}

	return &OpenAIClient{
		apiKey:       apiKey,
		baseURL:      baseURL,
		model:        model,
		organization: organization,
		proxy:        proxy,
		timeout:      time.Duration(timeoutSec) * time.Second,
		retryCount:   retryCount,
		client:       openai.NewClientWithConfig(clientConfig),
		config:       configManager, // 设置配置管理器
	}, nil
}

// NewOpenAIClient 创建新的 OpenAI 客户端 (兼容旧方法，使用默认配置)
func NewOpenAIClient() (*OpenAIClient, error) {
	// 使用空的配置管理器，实际使用时应使用 NewOpenAIClientWithConfig
	apiKey := os.Getenv("OPENAI_API_KEY")
	clientConfig := openai.DefaultConfig(apiKey)
	clientConfig.BaseURL = "https://api.openai.com/v1"

	return &OpenAIClient{
		apiKey:     apiKey,
		baseURL:    "https://api.openai.com/v1",
		model:      "gpt-3.5-turbo",
		timeout:    30 * time.Second,
		retryCount: 3,
		client:     openai.NewClientWithConfig(clientConfig),
	}, nil
}

// 从数据库重新加载配置
func (c *OpenAIClient) ReloadConfig() error {
	if c.config == nil {
		return fmt.Errorf("配置管理器未初始化")
	}

	newClient, err := NewOpenAIClientWithConfig(c.config)
	if err != nil {
		return err
	}

	c.client = newClient.client
	c.apiKey = newClient.apiKey
	c.baseURL = newClient.baseURL
	c.model = newClient.model
	c.organization = newClient.organization
	c.proxy = newClient.proxy
	c.timeout = newClient.timeout
	c.retryCount = newClient.retryCount

	return nil
}

// Chat 调用对话接口
func (c *OpenAIClient) Chat(messages []openai.ChatCompletionMessage, options ...ChatOption) (*openai.ChatCompletionResponse, error) {
	request := openai.ChatCompletionRequest{
		Model:    c.model,
		Messages: messages,
	}

	// 应用选项
	for _, option := range options {
		option(&request)
	}

	// 添加调试日志
	utils.Debug("[AI] OpenAI请求 - Model: %s, Messages: %d, Functions: %d",
		request.Model, len(request.Messages), len(request.Functions))

	if len(request.Functions) > 0 {
		for i, fn := range request.Functions {
			utils.Debug("[AI] Function %d: %s - %s", i, fn.Name, fn.Description)
		}
	}

	// 记录 FunctionCall 设置
	if request.FunctionCall != nil {
		utils.Debug("[OpenAI] FunctionCall 设置: %v", request.FunctionCall)
	}

	// 详细的请求参数调试信息
	utils.Debug("[AI] 请求参数摘要: max_tokens=%d temperature=%.2f top_p=%.2f function_call=%v stream=%v n=%d",
		request.MaxTokens, request.Temperature, request.TopP, request.FunctionCall, request.Stream, request.N)
	utils.Debug("[AI] 请求消息摘要: %s", summarizeMessagesForLog(request.Messages))

	resp, err := c.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		utils.Debug("[OpenAI] 请求失败: %v", err)
		return nil, err
	}

	utils.Debug("[OpenAI] 响应元数据: id=%s model=%s usage(prompt=%d completion=%d total=%d)",
		resp.ID, resp.Model, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)

	// 记录响应
	if len(resp.Choices) > 0 {
		for i, choice := range resp.Choices {
			contentLen := len([]rune(strings.TrimSpace(choice.Message.Content)))
			utils.Debug("[OpenAI] choice[%d] - finish_reason=%s, content_length=%d, has_function_call=%v",
				i, choice.FinishReason, contentLen, choice.Message.FunctionCall != nil)
			if choice.Message.FunctionCall != nil {
				utils.Debug("[OpenAI] choice[%d] 函数调用: name=%s, arguments_length=%d",
					i, choice.Message.FunctionCall.Name, len(choice.Message.FunctionCall.Arguments))
			}
		}
		firstPreview := previewRunesForLog(strings.TrimSpace(resp.Choices[0].Message.Content), 120)
		if firstPreview != "" {
			utils.Debug("[OpenAI] 首条响应内容预览: %q", firstPreview)
		} else {
			utils.Debug("[OpenAI] 首条响应内容预览为空")
		}
	}

	return &resp, nil
}

func summarizeMessagesForLog(messages []openai.ChatCompletionMessage) string {
	if len(messages) == 0 {
		return "empty"
	}

	parts := make([]string, 0, len(messages))
	for i, message := range messages {
		content := strings.TrimSpace(message.Content)
		preview := strings.ReplaceAll(previewRunesForLog(content, 80), "\n", " ")
		preview = strings.ReplaceAll(preview, "\r", " ")
		parts = append(parts, fmt.Sprintf("%d:%s(len=%d,preview=%q)",
			i, message.Role, len([]rune(content)), preview))
	}

	return strings.Join(parts, "; ")
}

func previewRunesForLog(s string, limit int) string {
	if limit <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit]) + "...(truncated)"
}

// ChatOption 对话选项类型
type ChatOption func(*openai.ChatCompletionRequest)

// WithMaxTokens 设置最大令牌数
func WithMaxTokens(tokens int) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		req.MaxTokens = tokens
	}
}

// WithTemperature 设置温度
func WithTemperature(temp float32) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		req.Temperature = temp
	}
}

// WithSystemPrompt 设置系统提示
func WithSystemPrompt(prompt string) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		// 确保消息数组以系统消息开始
		if len(req.Messages) == 0 {
			req.Messages = append([]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			}, req.Messages...)
		} else {
			// 如果已有消息，替换或添加系统消息
			req.Messages[0] = openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			}
		}
	}
}

// WithFunctions 设置函数调用
func WithFunctions(functions []openai.FunctionDefinition) ChatOption {
	return func(req *openai.ChatCompletionRequest) {
		if len(functions) > 0 {
			req.Functions = functions
			// 设置 FunctionCall 为 auto，让 AI 自动决定是否使用函数
			req.FunctionCall = "auto"
			utils.Debug("[OpenAI] 设置了 %d 个函数，FunctionCall: auto", len(functions))
		}
	}
}

// GetModel 获取当前模型名称
func (c *OpenAIClient) GetModel() string {
	return c.model
}

// CreateChatCompletion 创建对话完成
func (c *OpenAIClient) CreateChatCompletion(prompt string) (*openai.ChatCompletionResponse, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	return c.Chat(messages)
}

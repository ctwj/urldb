package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

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
		apiKey:       apiKey,
		baseURL:      "https://api.openai.com/v1",
		model:        "gpt-3.5-turbo",
		timeout:      30 * time.Second,
		retryCount:   3,
		client:       openai.NewClientWithConfig(clientConfig),
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
	log.Printf("[OpenAI] 发送请求 - Model: %s, Messages: %d, Functions: %d",
		request.Model, len(request.Messages), len(request.Functions))

	if len(request.Functions) > 0 {
		for i, fn := range request.Functions {
			log.Printf("[OpenAI] Function %d: %s - %s", i, fn.Name, fn.Description)
		}
	}

	// 记录 FunctionCall 设置
	if request.FunctionCall != "" {
		log.Printf("[OpenAI] FunctionCall 设置: %s", request.FunctionCall)
	}

	// 详细的请求参数调试信息
	log.Printf("[OpenAI] 请求参数详情:")
	log.Printf("  Model: %s", request.Model)
	log.Printf("  MaxTokens: %d", request.MaxTokens)
	log.Printf("  Temperature: %.2f", request.Temperature)
	log.Printf("  TopP: %.2f", request.TopP)
	log.Printf("  FunctionCall: %v", request.FunctionCall)
	log.Printf("  Stream: %v", request.Stream)
	log.Printf("  N: %d", request.N)

	// 完整的请求调试信息
	log.Printf("=== [OpenAI] 完整API请求数据 ===")
	if requestJSON, err := json.MarshalIndent(request, "", "  "); err == nil {
		log.Printf("完整请求JSON:\n%s", string(requestJSON))
	} else {
		log.Printf("序列化请求JSON失败: %v", err)
	}
	log.Printf("==============================")

	resp, err := c.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		log.Printf("[OpenAI] 请求失败: %v", err)
		return nil, err
	}

	// 记录响应
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		log.Printf("[OpenAI] 收到响应 - FinishReason: %s", choice.FinishReason)
		log.Printf("[OpenAI] 响应内容: %s", choice.Message.Content)
		if choice.Message.FunctionCall != nil {
			log.Printf("[OpenAI] AI 调用了函数: %s, 参数: %s", choice.Message.FunctionCall.Name, choice.Message.FunctionCall.Arguments)
		} else {
			log.Printf("[OpenAI] AI 没有调用函数")
		}
	}

	return &resp, nil
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
			log.Printf("[OpenAI] 设置了 %d 个函数，FunctionCall: auto", len(functions))
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
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/ctwj/urldb/utils"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool 表示MCP工具
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ConvertMCPToolToTool 将MCP工具转换为内部Tool结构
func ConvertMCPToolToTool(mcpTool mcp.Tool) Tool {
	var inputSchema map[string]interface{}

	// 检查InputSchema是否为空
	var interfaceToCheck interface{} = mcpTool.InputSchema
	if interfaceToCheck != nil {
		inputSchema = map[string]interface{}{
			"type":       mcpTool.InputSchema.Type,
			"properties": mcpTool.InputSchema.Properties,
			"required":   mcpTool.InputSchema.Required,
		}
		if mcpTool.InputSchema.Defs != nil {
			inputSchema["$defs"] = mcpTool.InputSchema.Defs
		}
	}

	return Tool{
		Name:        mcpTool.Name,
		Description: mcpTool.Description,
		InputSchema: inputSchema,
	}
}

// Client 表示MCP客户端
type Client struct {
	Name     string
	MCPCli   *client.Client
	Context  context.Context
	Cancel   context.CancelFunc
	Tools    []mcp.Tool
	Status   string // "running", "stopped", "error"
}

// ServiceStatus 服务状态信息
type ServiceStatus struct {
	Name      string   `json:"name"`
	Status    string   `json:"status"`    // "running", "stopped", "error"
	Tools     []Tool   `json:"tools"`
	Config    *MCPServerConfig `json:"config,omitempty"`
	LastError string   `json:"last_error,omitempty"`
}

// MCPServerConfig MCP服务器配置结构
type MCPServerConfig struct {
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Transport   string            `json:"transport"` // stdio, http, sse
	URL         string            `json:"url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Enabled     bool              `json:"enabled"`
	AutoStart   bool              `json:"auto_start"`
}

// MCPConfig MCP配置结构
type MCPConfig struct {
	MCPServers map[string]*MCPServerConfig `json:"mcpServers"`
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	registry map[string][]Tool
	mutex    sync.RWMutex
}

// Register 注册工具
func (tr *ToolRegistry) Register(serviceName string, tool Tool) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if _, exists := tr.registry[serviceName]; exists {
		tr.registry[serviceName] = append(tr.registry[serviceName], tool)
	} else {
		tr.registry[serviceName] = []Tool{tool}
	}
}

// GetTools 获取服务的所有工具
func (tr *ToolRegistry) GetTools(serviceName string) []Tool {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	tools, exists := tr.registry[serviceName]
	if !exists {
		return nil
	}
	return tools
}

// UnregisterService 注销服务的所有工具
func (tr *ToolRegistry) UnregisterService(serviceName string) {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	delete(tr.registry, serviceName)
}

// MCPManager MCP管理器
type MCPManager struct {
	clients      map[string]*Client
	configs      map[string]*MCPServerConfig
	toolReg      *ToolRegistry
	services     map[string]*ServiceStatus // 新增服务状态跟踪
	mutex        sync.RWMutex
	configPath   string  // 新增配置文件路径
}

// NewMCPManager 创建MCP管理器
func NewMCPManager() *MCPManager {
	return &MCPManager{
		clients: make(map[string]*Client),
		configs: make(map[string]*MCPServerConfig),
		services: make(map[string]*ServiceStatus),
		toolReg: &ToolRegistry{
			registry: make(map[string][]Tool),
		},
		configPath: "./data/mcp_config.json", // 默认配置路径
	}
}

// NewMCPManagerWithConfigPath 创建MCP管理器并指定配置文件路径
func NewMCPManagerWithConfigPath(configPath string) *MCPManager {
	manager := NewMCPManager()
	manager.configPath = configPath
	return manager
}

// LoadConfig 从配置文件加载MCP配置
func (m *MCPManager) LoadConfig(configFile string) error {
	// 读取配置文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("读取 MCP 配置文件失败: %v", err)
	}

	// 环境变量替换
	expandedData := m.expandEnvVars(string(data))

	var config MCPConfig
	if err := json.Unmarshal([]byte(expandedData), &config); err != nil {
		return fmt.Errorf("解析 MCP 配置失败: %v", err)
	}

	// 加载启用的服务器配置
	for name, serverConfig := range config.MCPServers {
		if serverConfig.Enabled {
			m.configs[name] = serverConfig

			// 初始化服务状态
			m.services[name] = &ServiceStatus{
				Name:   name,
				Status: "stopped",
				Config: serverConfig,
				Tools:  []Tool{},
			}

			// 自动启动配置
			if serverConfig.AutoStart {
				if err := m.StartClient(name); err != nil {
					fmt.Printf("自动启动 MCP 服务器 %s 失败: %v\n", name, err)
					m.services[name].Status = "error"
					m.services[name].LastError = err.Error()
				}
			}
		}
	}

	return nil
}

// expandEnvVars 环境变量替换，支持默认值
func (m *MCPManager) expandEnvVars(data string) string {
	// 匹配 ${VAR} 或 ${VAR:-default} 格式
	re := regexp.MustCompile(`\$\{([^}:]+)(?::([^}]*))?\}`)

	return re.ReplaceAllStringFunc(data, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		varName := parts[1]
		defaultValue := ""
		if len(parts) > 2 && parts[2] != "" {
			defaultValue = parts[2]
		}

		if value := os.Getenv(varName); value != "" {
			return value
		}
		return defaultValue
	})
}

// getDefaultTools 获取默认工具列表（当无法获取真实工具时使用）
func getDefaultTools(serviceName string) []mcp.Tool {
	switch serviceName {
	case "duckduckgo":
		return []mcp.Tool{
			{
				Name:        "duckduckgo_search",
				Description: "使用DuckDuckGo搜索引擎进行网络搜索",
				InputSchema: mcp.ToolInputSchema(mcp.ToolArgumentsSchema{
					Type: "object",
					Properties: map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "搜索查询词",
						},
						"max_results": map[string]any{
							"type":        "integer",
							"description": "返回结果的最大数量，默认为10",
							"default":     10,
						},
						"safe_search": map[string]any{
							"type":        "string",
							"description": "安全搜索级别：'off', 'moderate', 'strict'",
							"enum":        []string{"off", "moderate", "strict"},
							"default":     "moderate",
						},
						"time_range": map[string]any{
							"type":        "string",
							"description": "时间范围限制：'d' (天), 'w' (周), 'm' (月), 'y' (年)",
							"enum":        []string{"d", "w", "m", "y"},
						},
					},
					Required: []string{"query"},
				}),
			},
		}
	default:
		return []mcp.Tool{
			{
				Name:        "generic_tool",
				Description: "通用工具",
				InputSchema: mcp.ToolInputSchema(mcp.ToolArgumentsSchema{
					Type: "object",
					Properties: map[string]any{
						"action": map[string]any{
							"type":        "string",
							"description": "要执行的操作",
						},
					},
					Required: []string{"action"},
				}),
			},
		}
	}
}

// StartClient 启动MCP客户端
func (m *MCPManager) StartClient(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	config, exists := m.configs[name]
	if !exists {
		log.Printf("MCP 启动失败: 服务器配置 %s 未找到", name)
		return fmt.Errorf("MCP 服务器 %s 未找到", name)
	}

	log.Printf("开始启动 MCP 服务器: %s", name)
	log.Printf("配置 - 命令: %s, 参数: %v, 传输: %s", config.Command, config.Args, config.Transport)

	// 创建上下文 - 用于初始化，工具调用时会使用单独的context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// 根据传输类型创建客户端
	var mcpClient *client.Client
	var err error

	switch config.Transport {
	case "stdio":
		// 解析命令和参数
		args := config.Args
		if len(args) == 0 {
			cancel()
			log.Printf("MCP 启动失败: stdio 传输需要指定命令")
			return fmt.Errorf("stdio 传输需要指定命令")
		}

		command := args[0]
		cmdArgs := args[1:]

		log.Printf("创建 stdio 传输 - 命令: %s, 参数: %v", command, cmdArgs)

		// 设置环境变量
		env := os.Environ()
		for k, v := range config.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
			log.Printf("设置环境变量: %s=%s", k, v)
		}

		// 创建 stdio 传输
		stdioTransport := transport.NewStdio(command, env, cmdArgs...)

		// 创建客户端
		mcpClient = client.NewClient(stdioTransport)

	default:
		cancel()
		log.Printf("MCP 启动失败: 不支持的传输类型: %s", config.Transport)
		return fmt.Errorf("不支持的传输类型: %s", config.Transport)
	}

	// 启动客户端
	log.Printf("正在启动 MCP 客户端...")
	if err := mcpClient.Start(ctx); err != nil {
		cancel()
		log.Printf("MCP 启动失败: 客户端启动失败 - %v", err)
		return fmt.Errorf("启动MCP客户端失败: %v", err)
	}
	log.Printf("MCP 客户端启动成功")

	// 初始化客户端
	log.Printf("正在初始化 MCP 客户端...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "URLDB MCP Client",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := mcpClient.Initialize(ctx, initRequest)
	if err != nil {
		mcpClient.Close()
		cancel()
		log.Printf("MCP 启动失败: 客户端初始化失败 - %v", err)
		return fmt.Errorf("初始化MCP客户端失败: %v", err)
	}

	log.Printf("MCP 服务器 %s 连接成功: %s (版本 %s)",
		name, serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	// 获取工具列表
	var tools []mcp.Tool
	if serverInfo.Capabilities.Tools != nil {
		log.Printf("服务器支持工具，正在获取工具列表...")
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := mcpClient.ListTools(ctx, toolsRequest)
		if err != nil {
			log.Printf("获取工具列表失败: %v，使用默认工具", err)
			// 如果获取真实工具失败，使用默认工具
			tools = getDefaultTools(name)
		} else {
			tools = toolsResult.Tools
			log.Printf("成功获取到 %d 个真实工具", len(tools))
			for i, tool := range tools {
				log.Printf("工具 %d: %s - %s", i+1, tool.Name, tool.Description)
			}
		}
	} else {
		log.Printf("服务器不支持工具，使用默认工具")
		// 如果服务器不支持工具，使用默认工具
		tools = getDefaultTools(name)
	}

	// 创建客户端实例
	client := &Client{
		Name:    name,
		MCPCli:  mcpClient,
		Context: ctx,
		Cancel:  cancel,
		Tools:   tools,
	}

	// 注册工具到工具注册表
	for _, tool := range client.Tools {
		// 将 mcp.Tool 转换为内部 Tool 结构
		// 需要将 ToolArgumentsSchema 转换为 map[string]interface{}
		inputSchema := make(map[string]interface{})
		if tool.InputSchema.Type != "" {
			inputSchema["type"] = tool.InputSchema.Type
		}
		if tool.InputSchema.Properties != nil {
			inputSchema["properties"] = tool.InputSchema.Properties
		}
		if tool.InputSchema.Required != nil {
			inputSchema["required"] = tool.InputSchema.Required
		}

		internalTool := Tool{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: inputSchema,
		}
		m.toolReg.Register(name, internalTool)
	}

	m.clients[name] = client

	// 更新服务状态
	if serviceStatus, exists := m.services[name]; exists {
		serviceStatus.Status = "running"
		serviceStatus.LastError = ""

		// 转换工具
		var tools []Tool
		for _, tool := range client.Tools {
			tools = append(tools, ConvertMCPToolToTool(tool))
		}
		serviceStatus.Tools = tools
	}

	fmt.Printf("MCP 服务器 %s 启动成功\n", name)
	return nil
}

// StopClient 停止MCP客户端
func (m *MCPManager) StopClient(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.clients[name]
	if !exists {
		// 检查服务是否在配置中但未启动
		if serviceStatus, serviceExists := m.services[name]; serviceExists {
			if serviceStatus.Status == "stopped" {
				return fmt.Errorf("MCP 客户端 %s 已经停止", name)
			}
		}
		return fmt.Errorf("MCP 客户端 %s 未找到", name)
	}

	// 关闭MCP客户端连接
	if client.MCPCli != nil {
		client.MCPCli.Close()
	}

	// 取消上下文
	if client.Cancel != nil {
		client.Cancel()
	}

	delete(m.clients, name)

	// 更新服务状态为停止，但保留配置
	if serviceStatus, exists := m.services[name]; exists {
		serviceStatus.Status = "stopped"
		serviceStatus.LastError = ""
		serviceStatus.Tools = []Tool{} // 清空工具列表，因为服务已停止
	}

	fmt.Printf("MCP 服务器 %s 已停止\n", name)
	return nil
}

// CallTool 调用工具
func (m *MCPManager) CallTool(serviceName, toolName string, params map[string]interface{}) (interface{}, error) {
	m.mutex.RLock()
	client, exists := m.clients[serviceName]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("服务 %s 未启动", serviceName)
	}

	// 为工具调用创建单独的context，避免使用初始化时的短超时context
	toolCtx, toolCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer toolCancel()

	// 调用真实的MCP工具
	result, err := client.MCPCli.CallTool(toolCtx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: params,
		},
	})

	if err != nil {
		utils.Error("MCP工具调用失败 - 服务: %s, 工具: %s, 错误: %v", serviceName, toolName, err)
		return nil, fmt.Errorf("工具调用失败: %v", err)
	}

	// 处理返回结果
	if result.IsError {
		return nil, fmt.Errorf("工具返回错误")
	}

	// 提取文本内容
	var contentResults []interface{}
	for _, content := range result.Content {
		if textContent, ok := mcp.AsTextContent(content); ok {
			contentResults = append(contentResults, map[string]interface{}{
				"type": "text",
				"text": textContent.Text,
			})
		} else if imgContent, ok := mcp.AsImageContent(content); ok {
			contentResults = append(contentResults, map[string]interface{}{
				"type":      "image",
				"mime_type": imgContent.MIMEType,
				"data":      imgContent.Data,
			})
		} else {
			// 其他类型的内容
			contentResults = append(contentResults, map[string]interface{}{
				"type":    "unknown",
				"content": content,
			})
		}
	}

	// 返回格式化的结果
	return map[string]interface{}{
		"service":  serviceName,
		"tool":     toolName,
		"params":   params,
		"result":   contentResults,
		"status":   "success",
		"raw":      result,
	}, nil
}

// ListServices 列出所有服务
func (m *MCPManager) ListServices() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var services []string
	// 返回所有配置的服务，而不仅仅是运行中的
	for name := range m.services {
		services = append(services, name)
	}
	return services
}

// GetServiceStatus 获取服务状态
func (m *MCPManager) GetServiceStatus(serviceName string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	serviceStatus, exists := m.services[serviceName]
	if !exists {
		return false, fmt.Errorf("服务 %s 未找到", serviceName)
	}
	return serviceStatus.Status == "running", nil
}

// GetServiceInfo 获取服务信息
func (m *MCPManager) GetServiceInfo(serviceName string) (*Client, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[serviceName]
	if !exists {
		return nil, fmt.Errorf("服务 %s 未找到", serviceName)
	}
	return client, nil
}

// ListServiceStatuses 列出所有服务及其状态
func (m *MCPManager) ListServiceStatuses() map[string]*ServiceStatus {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]*ServiceStatus)
	for name, status := range m.services {
		// 复制状态信息
		statusCopy := &ServiceStatus{
			Name:      status.Name,
			Status:    status.Status,
			Config:    status.Config,
			LastError: status.LastError,
		}
		// 只有运行中的服务才有工具
		if status.Status == "running" {
			statusCopy.Tools = status.Tools
		}
		result[name] = statusCopy
	}
	return result
}

// GetToolRegistry 获取工具注册表
func (m *MCPManager) GetToolRegistry() *ToolRegistry {
	return m.toolReg
}

// RestartClient 重启MCP客户端
func (m *MCPManager) RestartClient(name string) error {
	if err := m.StopClient(name); err != nil {
		fmt.Printf("停止 MCP 客户端 %s 时出错: %v\n", name, err)
	}
	return m.StartClient(name)
}

// GetConfigFileContent 获取配置文件内容
func (m *MCPManager) GetConfigFileContent() (string, error) {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return "", fmt.Errorf("读取 MCP 配置文件失败: %v", err)
	}
	return string(data), nil
}

// UpdateConfigFileContent 更新配置文件内容
func (m *MCPManager) UpdateConfigFileContent(content string) error {
	// 验证JSON格式
	var configData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &configData); err != nil {
		return fmt.Errorf("配置JSON格式错误: %v", err)
	}

	// 备份原文件
	backupPath := m.configPath + ".backup"
	if _, err := os.Stat(m.configPath); err == nil {
		// 原文件存在，创建备份
		originalData, err := os.ReadFile(m.configPath)
		if err != nil {
			return fmt.Errorf("读取原配置文件失败: %v", err)
		}
		if err := os.WriteFile(backupPath, originalData, 0644); err != nil {
			return fmt.Errorf("创建备份文件失败: %v", err)
		}
	}

	// 写入新配置
	if err := os.WriteFile(m.configPath, []byte(content), 0644); err != nil {
		// 如果写入失败，尝试恢复备份
		if _, backupErr := os.Stat(backupPath); backupErr == nil {
			originalData, _ := os.ReadFile(backupPath)
			os.WriteFile(m.configPath, originalData, 0644) // 尝试恢复
		}
		return fmt.Errorf("更新配置文件失败: %v", err)
	}

	// 验证新配置是否可以加载
	tempManager := NewMCPManagerWithConfigPath(m.configPath)
	if err := tempManager.LoadConfig(m.configPath); err != nil {
		// 配置加载失败，恢复备份
		if _, backupErr := os.Stat(backupPath); backupErr == nil {
			originalData, _ := os.ReadFile(backupPath)
			os.WriteFile(m.configPath, originalData, 0644)
			os.Remove(backupPath) // 清理备份文件
		}
		return fmt.Errorf("新配置文件无法加载: %v", err)
	}

	// 清理备份文件
	os.Remove(backupPath)

	// 如果当前配置已加载，重新加载配置
	if len(m.configs) > 0 {
		// 停止所有客户端
		for name := range m.clients {
			m.StopClient(name)
		}
		// 清空当前配置
		m.configs = make(map[string]*MCPServerConfig)
		// 重新加载配置
		if err := m.LoadConfig(m.configPath); err != nil {
			return fmt.Errorf("重新加载配置失败: %v", err)
		}
	}

	return nil
}

// ReloadConfig 动态重新加载配置并更新服务
func (m *MCPManager) ReloadConfig(configContent string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 解析新配置
	var newConfig MCPConfig
	if err := json.Unmarshal([]byte(configContent), &newConfig); err != nil {
		return fmt.Errorf("解析新配置失败: %v", err)
	}

	// 获取当前运行的服务
	currentServices := make(map[string]bool)
	for name := range m.clients {
		currentServices[name] = true
	}

	// 获取新配置中的服务
	newServices := make(map[string]bool)
	for name := range newConfig.MCPServers {
		newServices[name] = true
	}

	// 停止不再需要服务的工具
	for name := range currentServices {
		if !newServices[name] {
			m.toolReg.UnregisterService(name)
		}
	}

	// 停止不再需要的服务
	for name := range currentServices {
		if !newServices[name] {
			if err := m.stopClientUnsafe(name); err != nil {
				fmt.Printf("停止服务 %s 时出错: %v\n", name, err)
			}
		}
	}

	// 更新配置
	m.configs = newConfig.MCPServers

	// 启动新服务或重启现有服务
	for name, serverConfig := range newConfig.MCPServers {
		if serverConfig.Enabled {
			if currentServices[name] {
				// 服务已存在，重启它
				if err := m.stopClientUnsafe(name); err != nil {
					fmt.Printf("停止服务 %s 时出错: %v\n", name, err)
				}
			}
			if err := m.startClientUnsafe(name); err != nil {
				fmt.Printf("启动服务 %s 时出错: %v\n", name, err)
			}
		}
	}

	fmt.Printf("MCP配置已动态重新加载，当前运行的服务: %v\n", m.ListServices())
	return nil
}

// startClientUnsafe 不加锁启动客户端（内部方法）
func (m *MCPManager) startClientUnsafe(name string) error {
	// 检查服务配置是否存在
	serverConfig, exists := m.configs[name]
	if !exists {
		return fmt.Errorf("服务配置 %s 不存在", name)
	}

	if !serverConfig.Enabled {
		return fmt.Errorf("服务 %s 已禁用", name)
	}

	// 检查是否已经启动
	if _, exists := m.clients[name]; exists {
		return fmt.Errorf("MCP 客户端 %s 已经启动", name)
	}

	// 创建模拟客户端
	m.clients[name] = &Client{Name: name}

	// 注册模拟工具
	tool := Tool{
		Name:        "mock-tool",
		Description: "模拟工具",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type": "string",
					"description": "查询参数",
				},
			},
			"required": []string{"query"},
		},
	}
	m.toolReg.Register(name, tool)

	fmt.Printf("MCP 服务器 %s 已启动\n", name)
	return nil
}

// stopClientUnsafe 不加锁停止客户端（内部方法）
func (m *MCPManager) stopClientUnsafe(name string) error {
	_, exists := m.clients[name]
	if !exists {
		return fmt.Errorf("MCP 客户端 %s 未找到", name)
	}

	delete(m.clients, name)
	fmt.Printf("MCP 服务器 %s 已停止\n", name)
	return nil
}
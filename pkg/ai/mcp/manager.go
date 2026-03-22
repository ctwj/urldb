package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
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
	Name    string
	MCPCli  *client.Client
	Context context.Context
	Cancel  context.CancelFunc
	Tools   []mcp.Tool
	Status  string // "running", "stopped", "error"
}

// ServiceStatus 服务状态信息
type ServiceStatus struct {
	Name      string           `json:"name"`
	Status    string           `json:"status"` // "running", "stopped", "error"
	Tools     []Tool           `json:"tools"`
	Config    *MCPServerConfig `json:"config,omitempty"`
	LastError string           `json:"last_error,omitempty"`
}

// MCPServerConfig MCP服务器配置结构
type MCPServerConfig struct {
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	Transport string            `json:"transport"` // stdio, http, sse
	Endpoint  string            `json:"endpoint,omitempty"`
	URL       string            `json:"url,omitempty"` // 保留以向后兼容
	Headers   map[string]string `json:"headers,omitempty"`
	Enabled   bool              `json:"enabled"`
	AutoStart bool              `json:"auto_start"`
}

// MCPConfig MCP配置结构
type MCPConfig struct {
	MCPServers map[string]*MCPServerConfig `json:"mcpServers"`
}

const (
	mcpStartupTotalTimeout = 60 * time.Second
	mcpStartTimeout        = 30 * time.Second
	mcpInitializeTimeout   = 30 * time.Second
	mcpListToolsTimeout    = 20 * time.Second
	mcpCloseTimeout        = 3 * time.Second
)

func validateMCPConfig(config *MCPConfig) error {
	if config == nil {
		return fmt.Errorf("配置不能为空")
	}

	for name, serverConfig := range config.MCPServers {
		if serverConfig == nil {
			return fmt.Errorf("服务 %s 配置为空", name)
		}

		// 与加载逻辑保持一致，仅校验启用的服务。
		if !serverConfig.Enabled {
			continue
		}

		switch serverConfig.Transport {
		case "stdio":
			if serverConfig.Command == "" {
				return fmt.Errorf("服务 %s 的 stdio 配置缺少 command", name)
			}
		case "http", "https", "sse":
			endpoint := serverConfig.Endpoint
			if endpoint == "" {
				endpoint = serverConfig.URL
			}
			if endpoint == "" {
				return fmt.Errorf("服务 %s 的 %s 配置缺少 endpoint", name, serverConfig.Transport)
			}
		default:
			return fmt.Errorf("服务 %s 使用了不支持的传输类型: %s", name, serverConfig.Transport)
		}
	}

	return nil
}

func initializeClientWithTimeout(ctx context.Context, mcpClient *client.Client, initRequest mcp.InitializeRequest) (*mcp.InitializeResult, error) {
	type initResult struct {
		result *mcp.InitializeResult
		err    error
	}

	resultCh := make(chan initResult, 1)
	go func() {
		result, err := mcpClient.Initialize(ctx, initRequest)
		resultCh <- initResult{result: result, err: err}
	}()

	select {
	case result := <-resultCh:
		return result.result, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func listToolsWithTimeout(ctx context.Context, mcpClient *client.Client, request mcp.ListToolsRequest) (*mcp.ListToolsResult, error) {
	type toolsResult struct {
		result *mcp.ListToolsResult
		err    error
	}

	resultCh := make(chan toolsResult, 1)
	go func() {
		result, err := mcpClient.ListTools(ctx, request)
		resultCh <- toolsResult{result: result, err: err}
	}()

	select {
	case result := <-resultCh:
		return result.result, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func closeClientWithTimeout(mcpClient *client.Client, timeout time.Duration) {
	if mcpClient == nil {
		return
	}

	done := make(chan struct{}, 1)
	go func() {
		mcpClient.Close()
		done <- struct{}{}
	}()

	select {
	case <-done:
		utils.Debug("[MCP] 客户端关闭完成")
	case <-time.After(timeout):
		utils.Warn("[MCP] 客户端关闭超时 (>%s)，将异步继续关闭", timeout)
	}
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
	clients    map[string]*Client
	configs    map[string]*MCPServerConfig
	toolReg    *ToolRegistry
	services   map[string]*ServiceStatus // 新增服务状态跟踪
	mutex      sync.RWMutex
	configPath string // 新增配置文件路径
}

// NewMCPManager 创建MCP管理器
func NewMCPManager() *MCPManager {
	return &MCPManager{
		clients:  make(map[string]*Client),
		configs:  make(map[string]*MCPServerConfig),
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

	// 读取配置成功后更新当前配置路径。
	m.configPath = configFile

	return m.ReloadConfig(string(data))
}

func (m *MCPManager) loadConfigFromStruct(config MCPConfig) {
	m.mutex.Lock()
	m.configs = make(map[string]*MCPServerConfig)
	m.services = make(map[string]*ServiceStatus)
	m.toolReg = &ToolRegistry{
		registry: make(map[string][]Tool),
	}
	for name, serverConfig := range config.MCPServers {
		if !serverConfig.Enabled {
			continue
		}
		m.configs[name] = serverConfig
		m.services[name] = &ServiceStatus{
			Name:   name,
			Status: "stopped",
			Config: serverConfig,
			Tools:  []Tool{},
		}
	}
	m.mutex.Unlock()
}

func (m *MCPManager) autoStartEnabledServices() {
	m.mutex.RLock()
	autoStartServices := make(map[string]*MCPServerConfig, len(m.configs))
	for name, serverConfig := range m.configs {
		autoStartServices[name] = serverConfig
	}
	m.mutex.RUnlock()

	for name, serverConfig := range autoStartServices {
		if !serverConfig.Enabled || !serverConfig.AutoStart {
			continue
		}

		// 启动流程不等待 MCP 初始化，后台异步完成连接与工具发现。
		go func(serviceName string) {
			utils.Info("[MCP] 异步自动启动服务: %s", serviceName)
			if err := m.StartClient(serviceName); err != nil {
				utils.Error("自动启动 MCP 服务器 %s 失败: %v", serviceName, err)
				m.mutex.Lock()
				if status, exists := m.services[serviceName]; exists {
					status.Status = "error"
					status.LastError = err.Error()
				}
				m.mutex.Unlock()
			}
		}(name)
	}
}

func (m *MCPManager) stopAllClients() {
	m.mutex.RLock()
	clientNames := make([]string, 0, len(m.clients))
	for name := range m.clients {
		clientNames = append(clientNames, name)
	}
	m.mutex.RUnlock()

	for _, name := range clientNames {
		if err := m.StopClient(name); err != nil {
			utils.Warn("停止 MCP 客户端 %s 时出错: %v", name, err)
		}
	}
}

// parseConfigContent 解析并验证配置内容，不产生启动副作用。
func parseConfigContent(content string) (MCPConfig, error) {
	var config MCPConfig
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		return MCPConfig{}, fmt.Errorf("解析 MCP 配置失败: %v", err)
	}
	if err := validateMCPConfig(&config); err != nil {
		return MCPConfig{}, err
	}
	return config, nil
}

// expandEnvVars 环境变量替换，支持默认值

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

	if existingClient, running := m.clients[name]; running && existingClient != nil && existingClient.MCPCli != nil {
		utils.Info("[MCP] 服务器 %s 已在运行，跳过重复启动", name)
		if serviceStatus, exists := m.services[name]; exists {
			serviceStatus.Status = "running"
		}
		return nil
	}

	config, exists := m.configs[name]
	if !exists {
		utils.Debug("MCP 启动失败: 服务器配置 %s 未找到", name)
		return fmt.Errorf("MCP 服务器 %s 未找到", name)
	}

	utils.Info("[MCP] 启动服务器: %s", name)
	utils.Debug("[MCP] 配置 - 命令: %s, 参数: %v, 传输: %s", config.Command, config.Args, config.Transport)

	// 启动阶段总超时，避免分阶段超时叠加导致整体启动时间过长。
	startupCtx, startupCancel := context.WithTimeout(context.Background(), mcpStartupTotalTimeout)
	defer startupCancel()

	// 运行期上下文用于统一管理客户端生命周期。
	runtimeCtx, runtimeCancel := context.WithCancel(context.Background())

	// 根据传输类型创建客户端
	var mcpClient *client.Client
	var err error

	switch config.Transport {
	case "stdio":
		// 解析命令和参数
		// 使用 config.Command 作为命令，config.Args 作为参数
		if config.Command == "" {
			runtimeCancel()
			utils.Debug("MCP 启动失败: stdio 传输需要指定 command")
			return fmt.Errorf("stdio 传输需要指定 command")
		}

		command := config.Command
		cmdArgs := config.Args

		// 在Windows上处理PowerShell执行策略问题
		if runtime.GOOS == "windows" {
			// 检查命令是否为npx，如果是，使用npm.cmd exec绕过执行策略
			if command == "npx" {
				utils.Debug("检测到Windows系统和npx命令，尝试使用npm.cmd exec")
				utils.Debug("原始命令: %s, 参数: %v", command, cmdArgs)

				// 使用npm.cmd exec作为包装器
				newArgs := []string{"exec", "npx", "--", cmdArgs[0]}
				newArgs = append(newArgs, cmdArgs[1:]...)
				utils.Debug("修改后的命令: npm.cmd, 参数: %v", newArgs)

				// 设置环境变量
				env := os.Environ()
				for k, v := range config.Env {
					env = append(env, fmt.Sprintf("%s=%s", k, v))
					utils.Debug("设置环境变量: %s=***", k)
				}

				// 使用 NewStdioMCPClientWithOptions 创建客户端（自动启动子进程）
				var err error
				mcpClient, err = client.NewStdioMCPClientWithOptions("npm.cmd", env, newArgs)
				if err != nil {
					runtimeCancel()
					utils.Debug("MCP 启动失败: 创建客户端失败 - %v", err)
					return fmt.Errorf("创建客户端失败: %v", err)
				}
			} else {
				utils.Debug("创建 stdio 传输 - 命令: %s, 参数: %v", command, cmdArgs)

				// 设置环境变量
				env := os.Environ()
				for k, v := range config.Env {
					env = append(env, fmt.Sprintf("%s=%s", k, v))
					utils.Debug("设置环境变量: %s=***", k)
				}

				// 使用 NewStdioMCPClientWithOptions 创建客户端（自动启动子进程）
				var err error
				mcpClient, err = client.NewStdioMCPClientWithOptions(command, env, cmdArgs)
				if err != nil {
					runtimeCancel()
					utils.Debug("MCP 启动失败: 创建客户端失败 - %v", err)
					return fmt.Errorf("创建客户端失败: %v", err)
				}
			}
		} else {
			utils.Debug("创建 stdio 传输 - 命令: %s, 参数: %v", command, cmdArgs)

			// 设置环境变量
			env := os.Environ()
			for k, v := range config.Env {
				env = append(env, fmt.Sprintf("%s=%s", k, v))
				utils.Debug("设置环境变量: %s=***", k)
			}

			// 使用 NewStdioMCPClientWithOptions 创建客户端（自动启动子进程）
			var err error
			mcpClient, err = client.NewStdioMCPClientWithOptions(command, env, cmdArgs)
			if err != nil {
				runtimeCancel()
				utils.Debug("MCP 启动失败: 创建客户端失败 - %v", err)
				return fmt.Errorf("创建客户端失败: %v", err)
			}
		}
	case "http", "https":
		// HTTP/HTTPS 传输：连接到远程MCP服务
		endpoint := config.Endpoint
		if endpoint == "" {
			endpoint = config.URL // 向后兼容旧的URL字段
		}
		if endpoint == "" {
			runtimeCancel()
			utils.Debug("MCP 启动失败: HTTP/HTTPS 传输需要指定 endpoint")
			return fmt.Errorf("HTTP/HTTPS 传输需要指定 endpoint")
		}

		utils.Debug("创建 HTTP/HTTPS 传输 - 端点: %s", endpoint)

		// 使用 NewStreamableHttpClient 连接到远程MCP服务
		var err error
		var options []transport.StreamableHTTPCOption
		if len(config.Headers) > 0 {
			options = append(options, transport.WithHTTPHeaders(config.Headers))
			utils.Debug("使用自定义 headers: %v", config.Headers)
		}
		mcpClient, err = client.NewStreamableHttpClient(endpoint, options...)
		if err != nil {
			runtimeCancel()
			utils.Debug("MCP 启动失败: 创建HTTP客户端失败 - %v", err)
			return fmt.Errorf("创建HTTP客户端失败: %v", err)
		}
	case "sse":
		// SSE (Server-Sent Events) 传输：连接到远程MCP服务
		endpoint := config.Endpoint
		if endpoint == "" {
			endpoint = config.URL // 向后兼容旧的URL字段
		}
		if endpoint == "" {
			runtimeCancel()
			utils.Debug("MCP 启动失败: SSE 传输需要指定 endpoint")
			return fmt.Errorf("SSE 传输需要指定 endpoint")
		}

		utils.Debug("创建 SSE 传输 - 端点: %s", endpoint)

		// 使用 NewSSEMCPClient 连接到远程MCP服务
		var err error
		mcpClient, err = client.NewSSEMCPClient(endpoint)
		if err != nil {
			runtimeCancel()
			utils.Debug("MCP 启动失败: 创建SSE客户端失败 - %v", err)
			return fmt.Errorf("创建SSE客户端失败: %v", err)
		}
	default:
		runtimeCancel()
		utils.Debug("MCP 启动失败: 不支持的传输类型: %s", config.Transport)
		return fmt.Errorf("不支持的传输类型: %s", config.Transport)
	}

	// 启动客户端
	// 注意：NewStdioMCPClientWithOptions 已经启动了子进程，但 client.Start() 仍然需要调用以完成客户端初始化
	startCtx, startCancel := context.WithTimeout(startupCtx, mcpStartTimeout)
	defer startCancel()

	utils.Debug("正在启动 MCP 客户端... (超时: %s)", mcpStartTimeout)
	startStageStart := time.Now()
	if err := mcpClient.Start(startCtx); err != nil {
		runtimeCancel()
		if errors.Is(err, context.DeadlineExceeded) {
			utils.Debug("MCP 启动失败: 客户端启动超时 (>%s，总超时: %s)", mcpStartTimeout, mcpStartupTotalTimeout)
			closeClientWithTimeout(mcpClient, mcpCloseTimeout)
			return fmt.Errorf("启动MCP客户端超时 (>%s)", mcpStartTimeout)
		}
		utils.Debug("MCP 启动失败: 客户端启动失败 - %v", err)
		closeClientWithTimeout(mcpClient, mcpCloseTimeout)
		return fmt.Errorf("启动MCP客户端失败: %v", err)
	}
	utils.Debug("MCP 客户端启动成功，耗时: %s", time.Since(startStageStart))

	// 初始化客户端
	initCtx, initCancel := context.WithTimeout(startupCtx, mcpInitializeTimeout)
	defer initCancel()

	utils.Debug("正在初始化 MCP 客户端... (超时: %s)", mcpInitializeTimeout)
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "URLDB MCP Client",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	initStart := time.Now()
	serverInfo, err := initializeClientWithTimeout(initCtx, mcpClient, initRequest)
	if err != nil {
		runtimeCancel()
		if errors.Is(err, context.DeadlineExceeded) {
			utils.Debug("MCP 启动失败: 客户端初始化超时 (>%s，总超时: %s)", mcpInitializeTimeout, mcpStartupTotalTimeout)
			closeClientWithTimeout(mcpClient, mcpCloseTimeout)
			return fmt.Errorf("初始化MCP客户端超时 (>%s)", mcpInitializeTimeout)
		}
		utils.Debug("MCP 启动失败: 客户端初始化失败 - %v", err)
		closeClientWithTimeout(mcpClient, mcpCloseTimeout)
		return fmt.Errorf("初始化MCP客户端失败: %v", err)
	}
	if serverInfo == nil {
		runtimeCancel()
		utils.Debug("MCP 启动失败: 客户端初始化返回空结果")
		closeClientWithTimeout(mcpClient, mcpCloseTimeout)
		return fmt.Errorf("初始化MCP客户端失败: 返回空结果")
	}
	utils.Debug("MCP 客户端初始化完成，耗时: %s", time.Since(initStart))

	utils.Debug("MCP 服务器 %s 连接成功: %s (版本 %s)",
		name, serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	// 获取工具列表
	var tools []mcp.Tool
	if serverInfo.Capabilities.Tools != nil {
		listToolsCtx, listToolsCancel := context.WithTimeout(startupCtx, mcpListToolsTimeout)
		defer listToolsCancel()

		utils.Debug("服务器支持工具，正在获取工具列表... (超时: %s)", mcpListToolsTimeout)
		toolsRequest := mcp.ListToolsRequest{}
		listToolsStart := time.Now()
		toolsResult, err := listToolsWithTimeout(listToolsCtx, mcpClient, toolsRequest)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				utils.Debug("获取工具列表超时 (>%s，总超时: %s)，使用默认工具", mcpListToolsTimeout, mcpStartupTotalTimeout)
			} else {
				utils.Debug("获取工具列表失败: %v，使用默认工具", err)
			}
			// 如果获取真实工具失败，使用默认工具
			tools = getDefaultTools(name)
		} else if toolsResult != nil {
			tools = toolsResult.Tools
			utils.Debug("成功获取到 %d 个真实工具，耗时: %s", len(tools), time.Since(listToolsStart))
			for i, tool := range tools {
				utils.Debug("工具 %d: %s - %s", i+1, tool.Name, tool.Description)
			}
		} else {
			utils.Debug("获取工具列表返回空结果，使用默认工具")
			tools = getDefaultTools(name)
		}
	} else {
		utils.Debug("服务器不支持工具，使用默认工具")
		// 如果服务器不支持工具，使用默认工具
		tools = getDefaultTools(name)
	}

	// 创建客户端实例
	client := &Client{
		Name:    name,
		MCPCli:  mcpClient,
		Context: runtimeCtx,
		Cancel:  runtimeCancel,
		Tools:   tools,
	}

	// 注册工具到工具注册表
	m.toolReg.UnregisterService(name)
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

	utils.Info("MCP 服务器 %s 启动成功\n", name)
	return nil
}

// StopClient 停止MCP客户端
func (m *MCPManager) StopClient(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.clients[name]
	if !exists {
		m.toolReg.UnregisterService(name)
		if serviceStatus, serviceExists := m.services[name]; serviceExists {
			serviceStatus.Status = "stopped"
			serviceStatus.LastError = ""
			serviceStatus.Tools = []Tool{}
		}
		utils.Info("[MCP] 服务器 %s 已是停止状态", name)
		return nil
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
	m.toolReg.UnregisterService(name)

	// 更新服务状态为停止，但保留配置
	if serviceStatus, exists := m.services[name]; exists {
		serviceStatus.Status = "stopped"
		serviceStatus.LastError = ""
		serviceStatus.Tools = []Tool{} // 清空工具列表，因为服务已停止
	}

	utils.Info("MCP 服务器 %s 已停止\n", name)
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
	if client == nil || client.MCPCli == nil {
		return nil, fmt.Errorf("服务 %s 客户端未初始化", serviceName)
	}

	// 为工具调用创建单独的context，避免使用初始化时的短超时context
	toolCtx, toolCancel := context.WithTimeout(context.Background(), 30*time.Second)
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
		"service": serviceName,
		"tool":    toolName,
		"params":  params,
		"result":  contentResults,
		"status":  "success",
		"raw":     result,
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

// CheckServiceHealth 检查服务健康状态
func (m *MCPManager) CheckServiceHealth(serviceName string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[serviceName]
	if !exists {
		utils.Debug("[MCP] 服务 %s 未找到", serviceName)
		return false
	}

	// 检查客户端是否已初始化
	if client.MCPCli == nil {
		utils.Debug("[MCP] 服务 %s 的客户端未初始化", serviceName)
		return false
	}

	return true
}

// RestartClient 重启MCP客户端
func (m *MCPManager) RestartClient(name string) error {
	if err := m.StopClient(name); err != nil {
		utils.Info("停止 MCP 客户端 %s 时出错: %v\n", name, err)
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
	// 无副作用地解析并验证配置。
	if _, err := parseConfigContent(content); err != nil {
		utils.Error("MCP配置验证失败: %v", err)
		return err
	}

	// 备份原文件
	backupPath := m.configPath + ".backup"
	var backupData []byte
	if _, err := os.Stat(m.configPath); err == nil {
		// 原文件存在，创建备份
		originalData, err := os.ReadFile(m.configPath)
		if err != nil {
			utils.Error("读取原配置文件失败: %v", err)
			return fmt.Errorf("读取原配置文件失败: %v", err)
		}
		backupData = originalData
		if err := os.WriteFile(backupPath, originalData, 0644); err != nil {
			utils.Error("创建备份文件失败: %v", err)
			return fmt.Errorf("创建备份文件失败: %v", err)
		}
	}

	// 写入新配置
	if err := os.WriteFile(m.configPath, []byte(content), 0644); err != nil {
		utils.Error("更新配置文件失败: %v", err)
		// 如果写入失败，尝试恢复备份
		if _, backupErr := os.Stat(backupPath); backupErr == nil {
			originalData, _ := os.ReadFile(backupPath)
			os.WriteFile(m.configPath, originalData, 0644) // 尝试恢复
		}
		return fmt.Errorf("更新配置文件失败: %v", err)
	}

	// 使用安全重载语义加载新配置。
	if err := m.ReloadConfig(content); err != nil {
		utils.Error("重新加载MCP配置失败: %v", err)
		if len(backupData) > 0 {
			if restoreErr := os.WriteFile(m.configPath, backupData, 0644); restoreErr != nil {
				utils.Error("恢复原配置文件失败: %v", restoreErr)
			} else if runtimeErr := m.ReloadConfig(string(backupData)); runtimeErr != nil {
				utils.Error("恢复运行时配置失败: %v", runtimeErr)
			}
		}
		return fmt.Errorf("重新加载配置失败: %v", err)
	}

	// 清理备份文件
	_ = os.Remove(backupPath)

	return nil
}

// RemoveServiceFromConfig 从配置文件中删除服务
func (m *MCPManager) RemoveServiceFromConfig(serviceName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 1. 读取当前配置
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config MCPConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 2. 检查服务是否存在
	if _, exists := config.MCPServers[serviceName]; !exists {
		return fmt.Errorf("服务 %s 不存在", serviceName)
	}

	// 3. 从配置中删除服务
	delete(config.MCPServers, serviceName)

	// 4. 从内存中删除服务
	delete(m.configs, serviceName)
	delete(m.services, serviceName)
	m.toolReg.UnregisterService(serviceName)

	// 5. 写回配置文件
	newConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(m.configPath, newConfig, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	utils.Info("MCP服务 %s 已从配置中删除", serviceName)
	return nil
}

// ReloadConfig 动态重新加载配置并更新服务
func (m *MCPManager) ReloadConfig(configContent string) error {
	newConfig, err := parseConfigContent(configContent)
	if err != nil {
		return err
	}

	m.stopAllClients()
	m.loadConfigFromStruct(newConfig)
	m.autoStartEnabledServices()

	utils.Info("MCP配置已重新加载，当前服务数量: %d", len(m.ListServices()))
	return nil
}

// startClientUnsafe 不加锁启动客户端（内部方法）
func (m *MCPManager) startClientUnsafe(name string) error {
	return fmt.Errorf("startClientUnsafe 已废弃，请使用 StartClient")
}

// stopClientUnsafe 不加锁停止客户端（内部方法）
func (m *MCPManager) stopClientUnsafe(name string) error {
	return fmt.Errorf("stopClientUnsafe 已废弃，请使用 StopClient")
}

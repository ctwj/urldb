package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sync"
)

// Tool 表示MCP工具
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Client 表示MCP客户端
type Client struct {
	Name string
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

// MCPManager MCP管理器
type MCPManager struct {
	clients      map[string]*Client
	configs      map[string]*MCPServerConfig
	toolReg      *ToolRegistry
	mutex        sync.RWMutex
	configPath   string  // 新增配置文件路径
}

// NewMCPManager 创建MCP管理器
func NewMCPManager() *MCPManager {
	return &MCPManager{
		clients: make(map[string]*Client),
		configs: make(map[string]*MCPServerConfig),
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

			// 自动启动配置
			if serverConfig.AutoStart {
				if err := m.StartClient(name); err != nil {
					fmt.Printf("自动启动 MCP 服务器 %s 失败: %v\n", name, err)
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

// StartClient 启动MCP客户端
func (m *MCPManager) StartClient(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.configs[name]
	if !exists {
		return fmt.Errorf("MCP 服务器 %s 未找到", name)
	}

	// 创建一个模拟客户端
	client := &Client{Name: name}

	// 注册一些模拟工具（用于演示）
	tools := []Tool{
		{
			Name:        "mock-tool",
			Description: "模拟工具",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"input": map[string]interface{}{
						"type":        "string",
						"description": "输入参数",
					},
				},
				"required": []string{"input"},
			},
		},
	}

	for _, tool := range tools {
		m.toolReg.Register(name, tool)
	}

	m.clients[name] = client
	fmt.Printf("MCP 服务器 %s 启动成功\n", name)
	return nil
}

// StopClient 停止MCP客户端
func (m *MCPManager) StopClient(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, exists := m.clients[name]
	if !exists {
		return fmt.Errorf("MCP 客户端 %s 未找到", name)
	}

	delete(m.clients, name)
	fmt.Printf("MCP 服务器 %s 已停止\n", name)
	return nil
}

// CallTool 调用工具
func (m *MCPManager) CallTool(serviceName, toolName string, params map[string]interface{}) (interface{}, error) {
	m.mutex.RLock()
	_, exists := m.clients[serviceName]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("服务 %s 未启动", serviceName)
	}

	// 模拟工具调用
	result := map[string]interface{}{
		"service": serviceName,
		"tool":    toolName,
		"params":  params,
		"result":  "模拟工具调用结果",
		"status":  "success",
	}

	return result, nil
}

// ListServices 列出所有服务
func (m *MCPManager) ListServices() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var services []string
	for name := range m.clients {
		services = append(services, name)
	}
	return services
}

// GetServiceStatus 获取服务状态
func (m *MCPManager) GetServiceStatus(serviceName string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, exists := m.clients[serviceName]
	if !exists {
		return false, fmt.Errorf("服务 %s 未找到", serviceName)
	}
	return true, nil
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
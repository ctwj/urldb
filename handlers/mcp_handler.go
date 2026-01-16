package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctwj/urldb/db/dto"
	"github.com/ctwj/urldb/pkg/ai/mcp"
	"github.com/ctwj/urldb/utils"

	"github.com/gin-gonic/gin"
)

// MCPHandler MCP处理器
type MCPHandler struct {
	mcpManager *mcp.MCPManager
}

// NewMCPHandler 创建MCP处理器
func NewMCPHandler(mcpManager *mcp.MCPManager) *MCPHandler {
	return &MCPHandler{
		mcpManager: mcpManager,
	}
}

// ListServices 列出所有MCP服务
func (h *MCPHandler) ListServices(c *gin.Context) {
	serviceStatuses := h.mcpManager.ListServiceStatuses()
	SuccessResponse(c, serviceStatuses)
}

// GetServiceInfo 获取服务信息
func (h *MCPHandler) GetServiceInfo(c *gin.Context) {
	serviceName := c.Param("name")

	serviceStatuses := h.mcpManager.ListServiceStatuses()
	serviceStatus, exists := serviceStatuses[serviceName]
	if !exists {
		ErrorResponse(c, "服务不存在", http.StatusNotFound)
		return
	}

	serviceInfo := map[string]interface{}{
		"name":       serviceName,
		"status":     serviceStatus.Status,
		"ready":      serviceStatus.Status == "running",
		"tools":      serviceStatus.Tools,
		"config":     serviceStatus.Config,
		"last_error": serviceStatus.LastError,
	}

	SuccessResponse(c, serviceInfo)
}

// StartService 启动服务
func (h *MCPHandler) StartService(c *gin.Context) {
	serviceName := c.Param("name")

	err := h.mcpManager.StartClient(serviceName)
	if err != nil {
		utils.Error("启动MCP服务失败 - 服务: %s, 错误: %v", serviceName, err)
		ErrorResponse(c, "启动服务失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "服务启动成功",
		"service": serviceName,
	})
}

// StopService 停止服务
func (h *MCPHandler) StopService(c *gin.Context) {
	serviceName := c.Param("name")

	err := h.mcpManager.StopClient(serviceName)
	if err != nil {
		utils.Error("停止MCP服务失败 - 服务: %s, 错误: %v", serviceName, err)
		ErrorResponse(c, "停止服务失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "服务停止成功",
		"service": serviceName,
	})
}

// RestartService 重启服务
func (h *MCPHandler) RestartService(c *gin.Context) {
	serviceName := c.Param("name")

	err := h.mcpManager.RestartClient(serviceName)
	if err != nil {
		utils.Error("重启MCP服务失败 - 服务: %s, 错误: %v", serviceName, err)
		ErrorResponse(c, "重启服务失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"message": "服务重启成功",
		"service": serviceName,
	})
}

// DeleteService 删除服务
func (h *MCPHandler) DeleteService(c *gin.Context) {
	serviceName := c.Param("name")

	// 1. 停止服务（如果正在运行）
	if err := h.mcpManager.StopClient(serviceName); err != nil {
		// 服务可能已经停止，忽略错误
		utils.Warn("停止MCP服务失败: %v", err)
	}

	// 2. 从配置中删除服务
	err := h.mcpManager.RemoveServiceFromConfig(serviceName)
	if err != nil {
		utils.Error("删除MCP服务失败 - 服务: %s, 错误: %v", serviceName, err)
		ErrorResponse(c, "删除服务失败", http.StatusInternalServerError)
		return
	}

	utils.Info("MCP服务删除成功 - 服务: %s", serviceName)
	SuccessResponse(c, map[string]interface{}{
		"message": "服务删除成功",
		"service": serviceName,
	})
}

// CallTool 调用工具
func (h *MCPHandler) CallTool(c *gin.Context) {
	var req dto.ToolCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	serviceName := c.Param("service")
	toolName := c.Param("tool")

	result, err := h.mcpManager.CallTool(serviceName, toolName, req.Params)
	if err != nil {
		utils.Error("调用MCP工具失败 - 服务: %s, 工具: %s, 错误: %v", serviceName, toolName, err)
		ErrorResponse(c, "工具调用失败", http.StatusInternalServerError)
		return
	}

	SuccessResponse(c, result)
}

// ListAllTools 列出所有工具
func (h *MCPHandler) ListAllTools(c *gin.Context) {
	services := h.mcpManager.ListServices()
	allTools := make(map[string][]interface{})

	for _, serviceName := range services {
		tools := h.mcpManager.GetToolRegistry().GetTools(serviceName)
		// 将 []mcp.Tool 转换为 []interface{}
		toolInterfaces := make([]interface{}, len(tools))
		for i, tool := range tools {
			toolInterfaces[i] = tool
		}
		allTools[serviceName] = toolInterfaces
	}

	SuccessResponse(c, allTools)
}

// ListServiceTools 列出服务的工具
func (h *MCPHandler) ListServiceTools(c *gin.Context) {
	serviceName := c.Param("service")

	tools := h.mcpManager.GetToolRegistry().GetTools(serviceName)
	// 将 []mcp.Tool 转换为 []interface{}
	toolInterfaces := make([]interface{}, len(tools))
	for i, tool := range tools {
		toolInterfaces[i] = tool
	}

	SuccessResponse(c, toolInterfaces)
}

// GetConfig 获取MCP配置文件内容
func (h *MCPHandler) GetConfig(c *gin.Context) {
	configContent, err := h.mcpManager.GetConfigFileContent()
	if err != nil {
		// 如果文件不存在或读取失败，返回一个默认配置
		utils.Warn("读取MCP配置文件失败: %v, 使用默认配置", err)
		configContent = `{
  "mcpServers": {
    "web-search": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-web-search"],
      "env": {
        "BING_API_KEY": "${BING_API_KEY}"
      },
      "transport": "stdio",
      "enabled": true,
      "auto_start": true
    },
    "filesystem": {
      "command": "npx",
      "args": ["@modelcontextprotocol/server-filesystem", "/tmp"],
      "transport": "stdio",
      "enabled": true,
      "auto_start": false
    }
  }
}`
		SuccessResponse(c, map[string]interface{}{
			"config": configContent,
		})
		return
	}

	SuccessResponse(c, map[string]interface{}{
		"config": configContent,
	})
}

// UpdateConfig 更新MCP配置文件内容
func (h *MCPHandler) UpdateConfig(c *gin.Context) {
	var req struct {
		Config string `json:"config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, "参数错误", http.StatusBadRequest)
		return
	}

	// 验证JSON格式
	var configData map[string]interface{}
	if err := json.Unmarshal([]byte(req.Config), &configData); err != nil {
		ErrorResponse(c, "配置JSON格式错误", http.StatusBadRequest)
		return
	}

	// 实际保存配置到文件
	err := h.mcpManager.UpdateConfigFileContent(req.Config)
	if err != nil {
		utils.Error("保存MCP配置文件失败: %v", err)
		ErrorResponse(c, "保存配置失败", http.StatusInternalServerError)
		return
	}

	utils.Info("MCP配置文件已成功更新")

	SuccessResponse(c, map[string]interface{}{
		"message": "配置保存成功",
		"services_updated": true,  // 标记服务已更新，前端可以刷新列表
	})
}
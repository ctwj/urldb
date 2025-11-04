package market

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// MarketManager 插件市场管理器
type MarketManager struct {
	client     *MarketClient
	manager    *manager.Manager
	pluginDir  string
	database   *gorm.DB
}

// NewMarketManager 创建新的市场管理器
func NewMarketManager(client *MarketClient, mgr *manager.Manager, pluginDir string, db *gorm.DB) *MarketManager {
	return &MarketManager{
		client:    client,
		manager:   mgr,
		pluginDir: pluginDir,
		database:  db,
	}
}

// SearchPlugins 搜索插件
func (mm *MarketManager) SearchPlugins(req SearchRequest) (*SearchResponse, error) {
	return mm.client.Search(req)
}

// GetPluginInfo 获取插件详细信息
func (mm *MarketManager) GetPluginInfo(pluginID string) (*PluginInfo, error) {
	return mm.client.GetPlugin(pluginID)
}

// InstallPlugin 安装插件
func (mm *MarketManager) InstallPlugin(req InstallRequest) (*InstallResponse, error) {
	resp := &InstallResponse{
		Success: false,
	}

	// 获取插件信息
	pluginInfo, err := mm.client.GetPlugin(req.PluginID)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to get plugin info: %v", err)
		return resp, nil
	}

	// 检查插件是否已安装
	_, err = mm.manager.GetPlugin(pluginInfo.Name)
	if err == nil && !req.Force {
		resp.Error = fmt.Sprintf("Plugin %s is already installed", pluginInfo.Name)
		return resp, nil
	}

	// 检查兼容性
	if !mm.checkCompatibility(pluginInfo.Compatibility) {
		resp.Error = fmt.Sprintf("Plugin %s is not compatible with current system", pluginInfo.Name)
		return resp, nil
	}

	// 下载插件
	pluginData, err := mm.client.DownloadPlugin(req.PluginID, req.Version)
	if err != nil {
		resp.Error = fmt.Sprintf("Failed to download plugin: %v", err)
		return resp, nil
	}

	// 验证校验和
	if pluginInfo.Checksum != "" {
		if !mm.verifyChecksum(pluginData, pluginInfo.Checksum) {
			resp.Error = "Plugin checksum verification failed"
			return resp, nil
		}
	}

	// 保存插件文件
	pluginPath := filepath.Join(mm.pluginDir, fmt.Sprintf("%s.so", pluginInfo.Name))
	if err := mm.savePluginFile(pluginPath, pluginData); err != nil {
		resp.Error = fmt.Sprintf("Failed to save plugin file: %v", err)
		return resp, nil
	}

	// 如果插件已存在且强制安装，先卸载
	if err == nil && req.Force {
		if err := mm.manager.UninstallPlugin(pluginInfo.Name, true); err != nil {
			utils.Warn("Failed to uninstall existing plugin: %v", err)
		}
	}

	// 加载并注册插件
	// 注意：这里需要根据具体的插件加载机制来实现
	// plugin, err := mm.loadPlugin(pluginPath)
	// if err != nil {
	// 	resp.Error = fmt.Sprintf("Failed to load plugin: %v", err)
	// 	return resp, nil
	// }

	// 注册插件
	// if err := mm.manager.RegisterPlugin(plugin); err != nil {
	// 	resp.Error = fmt.Sprintf("Failed to register plugin: %v", err)
	// 	return resp, nil
	// }

	// 初始化插件
	config := make(map[string]interface{})
	if err := mm.manager.InitializePlugin(pluginInfo.Name, config); err != nil {
		resp.Error = fmt.Sprintf("Failed to initialize plugin: %v", err)
		return resp, nil
	}

	// 记录安装信息到数据库
	if err := mm.recordInstallation(pluginInfo); err != nil {
		utils.Warn("Failed to record installation: %v", err)
	}

	resp.Success = true
	resp.Message = fmt.Sprintf("Plugin %s installed successfully", pluginInfo.Name)
	resp.PluginName = pluginInfo.Name
	resp.Version = pluginInfo.Version

	return resp, nil
}

// UninstallPlugin 卸载插件
func (mm *MarketManager) UninstallPlugin(pluginName string, force bool) (*InstallResponse, error) {
	resp := &InstallResponse{
		Success: false,
	}

	// 检查插件是否存在
	_, err := mm.manager.GetPlugin(pluginName)
	if err != nil {
		resp.Error = fmt.Sprintf("Plugin %s not found", pluginName)
		return resp, nil
	}

	// 卸载插件
	if err := mm.manager.UninstallPlugin(pluginName, force); err != nil {
		resp.Error = fmt.Sprintf("Failed to uninstall plugin: %v", err)
		return resp, nil
	}

	// 删除插件文件
	pluginPath := filepath.Join(mm.pluginDir, fmt.Sprintf("%s.so", pluginName))
	if _, err := os.Stat(pluginPath); err == nil {
		if err := os.Remove(pluginPath); err != nil {
			utils.Warn("Failed to remove plugin file: %v", err)
		}
	}

	// 从数据库中删除安装记录
	if err := mm.removeInstallationRecord(pluginName); err != nil {
		utils.Warn("Failed to remove installation record: %v", err)
	}

	resp.Success = true
	resp.Message = fmt.Sprintf("Plugin %s uninstalled successfully", pluginName)
	resp.PluginName = pluginName

	return resp, nil
}

// UpdatePlugin 更新插件
func (mm *MarketManager) UpdatePlugin(pluginName string) (*InstallResponse, error) {
	// 获取当前安装的插件信息
	currentPlugin, err := mm.manager.GetPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("plugin %s not found: %v", pluginName, err)
	}

	// 在市场中搜索同名插件
	searchReq := SearchRequest{
		Query: pluginName,
	}
	searchResp, err := mm.client.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search for plugin updates: %v", err)
	}

	// 查找匹配的插件
	var marketPlugin *PluginInfo
	for _, plugin := range searchResp.Plugins {
		if plugin.Name == pluginName {
			marketPlugin = &plugin
			break
		}
	}

	if marketPlugin == nil {
		return nil, fmt.Errorf("plugin %s not found in market", pluginName)
	}

	// 比较版本
	if marketPlugin.Version == currentPlugin.Version() {
		return &InstallResponse{
			Success:    true,
			Message:    "Plugin is already up to date",
			PluginName: pluginName,
			Version:    marketPlugin.Version,
		}, nil
	}

	// 执行更新（相当于重新安装）
	installReq := InstallRequest{
		PluginID: marketPlugin.ID,
		Version:  marketPlugin.Version,
		Force:    true,
	}

	return mm.InstallPlugin(installReq)
}

// ListInstalledPlugins 列出已安装的插件
func (mm *MarketManager) ListInstalledPlugins() []types.PluginInfo {
	return mm.manager.ListPlugins()
}

// GetPluginReviews 获取插件评价
func (mm *MarketManager) GetPluginReviews(pluginID string, page, pageSize int) ([]PluginReview, error) {
	return mm.client.GetPluginReviews(pluginID, page, pageSize)
}

// SubmitReview 提交评价
func (mm *MarketManager) SubmitReview(review PluginReview) error {
	return mm.client.SubmitReview(review)
}

// GetCategories 获取插件分类
func (mm *MarketManager) GetCategories() ([]PluginCategory, error) {
	return mm.client.GetCategories()
}

// checkCompatibility 检查系统兼容性
func (mm *MarketManager) checkCompatibility(compatibility []string) bool {
	// 简化实现，实际需要检查操作系统、架构等
	if len(compatibility) == 0 {
		return true // 没有指定兼容性要求，默认兼容
	}

	// 获取当前系统信息
	// 这里需要根据实际情况实现
	currentOS := "linux" // 简化示例
	currentArch := "amd64" // 简化示例

	currentSystem := fmt.Sprintf("%s/%s", currentOS, currentArch)

	for _, compat := range compatibility {
		if compat == "all" || compat == currentSystem {
			return true
		}
	}

	return false
}

// verifyChecksum 验证文件校验和
func (mm *MarketManager) verifyChecksum(data []byte, expectedChecksum string) bool {
	hash := md5.Sum(data)
	actualChecksum := hex.EncodeToString(hash[:])
	return strings.EqualFold(actualChecksum, expectedChecksum)
}

// savePluginFile 保存插件文件
func (mm *MarketManager) savePluginFile(path string, data []byte) error {
	// 确保插件目录存在
	pluginDir := filepath.Dir(path)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// 写入文件
	return os.WriteFile(path, data, 0644)
}

// loadPlugin 加载插件（具体实现依赖于插件加载机制）
// func (mm *MarketManager) loadPlugin(pluginPath string) (types.Plugin, error) {
// 	// 这里需要根据具体的插件加载机制来实现
// 	// 可能使用 Go 的 plugin 包或其他自定义机制
// 	return nil, fmt.Errorf("not implemented: loadPlugin")
// }

// recordInstallation 记录安装信息
func (mm *MarketManager) recordInstallation(plugin *PluginInfo) error {
	// 这里需要实现将安装信息保存到数据库的逻辑
	// 可以创建一个新的表来存储插件市场相关的安装信息
	utils.Info("Recording installation of plugin: %s", plugin.Name)
	return nil
}

// removeInstallationRecord 删除安装记录
func (mm *MarketManager) removeInstallationRecord(pluginName string) error {
	// 删除数据库中的安装记录
	utils.Info("Removing installation record for plugin: %s", pluginName)
	return nil
}

// GetInstalledPluginInfo 获取已安装插件的市场信息
func (mm *MarketManager) GetInstalledPluginInfo(pluginName string) (*PluginInfo, error) {
	// 这里可以实现从数据库或其他存储中获取已安装插件的市场信息
	// 用于显示插件的详细信息、检查更新等
	return nil, fmt.Errorf("not implemented: GetInstalledPluginInfo")
}
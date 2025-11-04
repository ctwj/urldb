package version

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ctwj/urldb/plugin/manager"
	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
	"gorm.io/gorm"
)

// VersionManager 插件版本管理器
type VersionManager struct {
	manager  *manager.Manager
	database *gorm.DB
	history  map[string][]PluginVersion // plugin_name -> versions
}

// NewVersionManager 创建新的版本管理器
func NewVersionManager(mgr *manager.Manager, db *gorm.DB) *VersionManager {
	vm := &VersionManager{
		manager:  mgr,
		database: db,
		history:  make(map[string][]PluginVersion),
	}

	// 从数据库加载版本历史
	vm.loadVersionHistory()

	return vm
}

// RegisterVersion 注册插件新版本
func (vm *VersionManager) RegisterVersion(pluginName, version, description, author string) error {
	// 检查插件是否存在
	plugin, err := vm.manager.GetPlugin(pluginName)
	if err != nil {
		return fmt.Errorf("plugin %s not found: %v", pluginName, err)
	}

	// 检查版本是否已存在
	if vm.hasVersion(pluginName, version) {
		return fmt.Errorf("version %s already exists for plugin %s", version, pluginName)
	}

	// 创建版本信息
	versionInfo := PluginVersion{
		ID:          fmt.Sprintf("%s_%s_%d", pluginName, version, time.Now().Unix()),
		PluginName:  pluginName,
		Version:     version,
		Description: description,
		Author:      author,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ReleaseDate: time.Now(),
		Status:      VersionStatusActive,
		Changelog:   fmt.Sprintf("Initial version %s for plugin %s", version, pluginName),
		Dependencies: plugin.Dependencies(),
	}

	// 添加到版本历史
	vm.addVersionToHistory(pluginName, versionInfo)

	// 保存到数据库
	if err := vm.saveVersionToDB(versionInfo); err != nil {
		return fmt.Errorf("failed to save version to database: %v", err)
	}

	utils.Info("Registered new version %s for plugin %s", version, pluginName)
	return nil
}

// GetVersion 获取特定版本的插件信息
func (vm *VersionManager) GetVersion(pluginName, version string) (*PluginVersion, error) {
	vm.loadVersionHistory() // 确保历史记录是最新的

	versions, exists := vm.history[pluginName]
	if !exists {
		return nil, fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	for _, v := range versions {
		if v.Version == version {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("version %s not found for plugin %s", version, pluginName)
}

// GetLatestVersion 获取最新版本
func (vm *VersionManager) GetLatestVersion(pluginName string) (*PluginVersion, error) {
	vm.loadVersionHistory()

	versions, exists := vm.history[pluginName]
	if !exists || len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	// 按版本号排序，获取最新的
	latest := versions[0]
	for _, v := range versions[1:] {
		if vm.compareVersions(v.Version, latest.Version) > 0 {
			latest = v
		}
	}

	return &latest, nil
}

// GetStableVersion 获取最新稳定版本
func (vm *VersionManager) GetStableVersion(pluginName string) (*PluginVersion, error) {
	vm.loadVersionHistory()

	versions, exists := vm.history[pluginName]
	if !exists || len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	var stableVersions []PluginVersion
	for _, v := range versions {
		if v.Status == VersionStatusStable {
			stableVersions = append(stableVersions, v)
		}
	}

	if len(stableVersions) == 0 {
		// 如果没有稳定版本，返回最新版本
		return vm.GetLatestVersion(pluginName)
	}

	// 按版本号排序，获取最新的稳定版本
	latest := stableVersions[0]
	for _, v := range stableVersions[1:] {
		if vm.compareVersions(v.Version, latest.Version) > 0 {
			latest = v
		}
	}

	return &latest, nil
}

// ListVersions 列出插件所有版本
func (vm *VersionManager) ListVersions(pluginName string) ([]PluginVersion, error) {
	vm.loadVersionHistory()

	versions, exists := vm.history[pluginName]
	if !exists {
		return nil, fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	// 按版本号排序
	sortedVersions := make([]PluginVersion, len(versions))
	copy(sortedVersions, versions)
	sort.Slice(sortedVersions, func(i, j int) bool {
		return vm.compareVersions(sortedVersions[i].Version, sortedVersions[j].Version) > 0
	})

	return sortedVersions, nil
}

// CheckVersionCompatibility 检查版本兼容性
func (vm *VersionManager) CheckVersionCompatibility(pluginName, version string) (*CompatibilityInfo, error) {
	pluginVersion, err := vm.GetVersion(pluginName, version)
	if err != nil {
		return nil, err
	}

	// 获取当前系统信息
	systemInfo := vm.getCurrentSystemInfo()

	// 检查兼容性要求
	compatible := true
	reason := ""

	// 检查依赖兼容性
	if len(pluginVersion.Dependencies) > 0 {
		// 这里可以实现更复杂的依赖检查逻辑
		for _, dep := range pluginVersion.Dependencies {
			// 检查依赖是否存在且版本兼容
			if _, err := vm.manager.GetPlugin(dep); err != nil {
				compatible = false
				reason = fmt.Sprintf("dependency %s not found", dep)
				break
			}
		}
	}

	// 检查系统兼容性
	if len(pluginVersion.Compatibility) > 0 {
		currentSystem := fmt.Sprintf("%s/%s", systemInfo.OS, systemInfo.Arch)
		isCompatible := false
		for _, compat := range pluginVersion.Compatibility {
			if compat == "all" || compat == currentSystem {
				isCompatible = true
				break
			}
		}
		if !isCompatible {
			compatible = false
			reason = fmt.Sprintf("not compatible with system %s", currentSystem)
		}
	}

	return &CompatibilityInfo{
		PluginName: pluginName,
		Version:    version,
		Compatible: compatible,
		Reason:     reason,
		SystemInfo: systemInfo,
	}, nil
}

// CompareVersions 比较两个版本
func (vm *VersionManager) CompareVersions(version1, version2 string) (*VersionComparison, error) {
	result := &VersionComparison{
		Current: version1,
		Available: version2,
	}

	comp := vm.compareVersions(version1, version2)
	if comp > 0 {
		result.IsNewer = true
	} else if comp < 0 {
		result.IsOlder = true
	} else {
		result.IsSame = true
	}

	return result, nil
}

// GetVersionRequirements 获取版本要求
func (vm *VersionManager) GetVersionRequirements(pluginName string) (*VersionRequirement, error) {
	// 这里可以实现获取插件版本依赖关系的逻辑
	// 检查插件及其依赖的版本要求
	plugin, err := vm.manager.GetPlugin(pluginName)
	if err != nil {
		return nil, fmt.Errorf("plugin %s not found: %v", pluginName, err)
	}

	// 获取当前版本
	current, err := vm.GetLatestVersion(pluginName)
	if err != nil {
		return nil, err
	}

	return &VersionRequirement{
		PluginName: pluginName,
		MinVersion: current.Version,
		MaxVersion: current.Version,
		ExactVersion: current.Version,
	}, nil
}

// UpdateVersionStatus 更新版本状态
func (vm *VersionManager) UpdateVersionStatus(pluginName, version string, status VersionStatus) error {
	vm.loadVersionHistory()

	versions, exists := vm.history[pluginName]
	if !exists {
		return fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	for i, v := range versions {
		if v.Version == version {
			// 更新状态
			vm.history[pluginName][i].Status = status
			vm.history[pluginName][i].UpdatedAt = time.Now()

			// 更新数据库
			if err := vm.updateVersionInDB(pluginName, version, status); err != nil {
				return fmt.Errorf("failed to update version in database: %v", err)
			}

			utils.Info("Updated status of version %s for plugin %s to %s", version, pluginName, status)
			return nil
		}
	}

	return fmt.Errorf("version %s not found for plugin %s", version, pluginName)
}

// GetVersionHistory 获取版本历史
func (vm *VersionManager) GetVersionHistory(pluginName string) (*VersionHistory, error) {
	vm.loadVersionHistory()

	versions, exists := vm.history[pluginName]
	if !exists {
		return nil, fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	// 获取当前运行的版本
	var currentVersion string
	plugin, err := vm.manager.GetPlugin(pluginName)
	if err == nil {
		currentVersion = plugin.Version()
	}

	return &VersionHistory{
		PluginName: pluginName,
		Versions:   versions,
		Current:    currentVersion,
	}, nil
}

// DeleteVersion 删除版本（通常不推荐，可以改为标记为废弃）
func (vm *VersionManager) DeleteVersion(pluginName, version string) error {
	vm.loadVersionHistory()

	versions, exists := vm.history[pluginName]
	if !exists {
		return fmt.Errorf("no versions found for plugin %s", pluginName)
	}

	// 找到要删除的版本
	found := false
	var newVersions []PluginVersion
	for _, v := range versions {
		if v.Version == version {
			found = true
		} else {
			newVersions = append(newVersions, v)
		}
	}

	if !found {
		return fmt.Errorf("version %s not found for plugin %s", version, pluginName)
	}

	// 更新历史记录
	vm.history[pluginName] = newVersions

	// 从数据库中删除版本记录
	if err := vm.deleteVersionFromDB(pluginName, version); err != nil {
		return fmt.Errorf("failed to delete version from database: %v", err)
	}

	utils.Info("Deleted version %s for plugin %s", version, pluginName)
	return nil
}

// LoadPluginByVersion 从特定版本加载插件
func (vm *VersionManager) LoadPluginByVersion(pluginName, version string) (types.Plugin, error) {
	// 获取版本信息
	versionInfo, err := vm.GetVersion(pluginName, version)
	if err != nil {
		return nil, err
	}

	// 检查兼容性
	compatibility, err := vm.CheckVersionCompatibility(pluginName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to check compatibility: %v", err)
	}
	if !compatibility.Compatible {
		return nil, fmt.Errorf("version %s is not compatible: %s", version, compatibility.Reason)
	}

	// 根据版本信息加载插件文件
	// 这里需要根据具体插件加载机制实现
	// plugin, err := vm.loadPluginByVersionInfo(versionInfo)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to load plugin version: %v", err)
	// }

	// 返回插件实例
	return nil, fmt.Errorf("not implemented: LoadPluginByVersion")
}

// loadVersionHistory 从数据库加载版本历史
func (vm *VersionManager) loadVersionHistory() error {
	// 从数据库加载版本历史信息
	// 由于我们没有定义数据库表，这里简化实现
	return nil
}

// addVersionToHistory 添加版本到历史记录
func (vm *VersionManager) addVersionToHistory(pluginName string, version PluginVersion) {
	if _, exists := vm.history[pluginName]; !exists {
		vm.history[pluginName] = make([]PluginVersion, 0)
	}
	vm.history[pluginName] = append(vm.history[pluginName], version)
}

// hasVersion 检查是否存在指定版本
func (vm *VersionManager) hasVersion(pluginName, version string) bool {
	versions, exists := vm.history[pluginName]
	if !exists {
		return false
	}

	for _, v := range versions {
		if v.Version == version {
			return true
		}
	}
	return false
}

// saveVersionToDB 保存版本到数据库
func (vm *VersionManager) saveVersionToDB(versionInfo PluginVersion) error {
	// 这里需要实现保存到数据库的逻辑
	// 需要先创建对应的数据库表
	utils.Debug("Saving version to database: %s - %s", versionInfo.PluginName, versionInfo.Version)
	return nil
}

// updateVersionInDB 更新数据库中的版本信息
func (vm *VersionManager) updateVersionInDB(pluginName, version string, status VersionStatus) error {
	// 需要实现更新数据库中的版本状态
	utils.Debug("Updating version status in database: %s - %s -> %s", pluginName, version, status)
	return nil
}

// deleteVersionFromDB 从数据库删除版本记录
func (vm *VersionManager) deleteVersionFromDB(pluginName, version string) error {
	// 需要实现从数据库删除版本记录
	utils.Debug("Deleting version from database: %s - %s", pluginName, version)
	return nil
}

// compareVersions 比较两个版本号 (返回: 1 if v1 > v2, -1 if v1 < v2, 0 if equal)
func (vm *VersionManager) compareVersions(v1, v2 string) int {
	// 简化版本比较逻辑，实际实现可能需要更复杂的语义化版本比较
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var num1, num2 int

		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &num1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &num2)
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	return 0
}

// getCurrentSystemInfo 获取当前系统信息
func (vm *VersionManager) getCurrentSystemInfo() SystemInfo {
	// 在实际实现中，这里应该返回真实的系统信息
	return SystemInfo{
		OS:      "linux",   // 简化示例
		Arch:    "amd64",   // 简化示例
		Version: "1.0.0",   // 简化示例
	}
}

// loadPluginByVersionInfo 根据版本信息加载插件（具体实现依赖于插件加载机制）
// func (vm *VersionManager) loadPluginByVersionInfo(versionInfo PluginVersion) (types.Plugin, error) {
// 	// 这里需要根据具体的插件加载机制来实现
// 	// 可能需要使用 Go 的 plugin 包或其他自定义机制
// 	return nil, fmt.Errorf("not implemented: loadPluginByVersionInfo")
// }
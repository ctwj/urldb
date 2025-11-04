package demo

import (
	"fmt"
	"time"

	"github.com/ctwj/urldb/plugin/types"
)

// DatabasePlugin 模拟一个数据库插件
type DatabasePlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
}

// NewDatabasePlugin 创建数据库插件实例
func NewDatabasePlugin() *DatabasePlugin {
	return &DatabasePlugin{
		name:        "database-plugin",
		version:     "1.0.0",
		description: "A database plugin that provides database access services",
		author:      "urlDB Team",
	}
}

// Name 返回插件名称
func (p *DatabasePlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *DatabasePlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *DatabasePlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *DatabasePlugin) Author() string {
	return p.author
}

// Initialize 初始化插件
func (p *DatabasePlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Database plugin initialized")
	return nil
}

// Start 启动插件
func (p *DatabasePlugin) Start() error {
	p.context.LogInfo("Database plugin started")
	return nil
}

// Stop 停止插件
func (p *DatabasePlugin) Stop() error {
	p.context.LogInfo("Database plugin stopped")
	return nil
}

// Cleanup 清理插件
func (p *DatabasePlugin) Cleanup() error {
	p.context.LogInfo("Database plugin cleaned up")
	return nil
}

// Dependencies 返回插件依赖
func (p *DatabasePlugin) Dependencies() []string {
	return []string{} // 数据库插件没有依赖
}

// CheckDependencies 检查插件依赖状态
func (p *DatabasePlugin) CheckDependencies() map[string]bool {
	return make(map[string]bool) // 没有依赖，返回空map
}

// AuthPlugin 模拟一个认证插件，依赖数据库插件
type AuthPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
	dependencies []string
}

// NewAuthPlugin 创建认证插件实例
func NewAuthPlugin() *AuthPlugin {
	return &AuthPlugin{
		name:        "auth-plugin",
		version:     "1.0.0",
		description: "An authentication plugin that depends on database plugin",
		author:      "urlDB Team",
		dependencies: []string{"database-plugin"}, // 依赖数据库插件
	}
}

// Name 返回插件名称
func (p *AuthPlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *AuthPlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *AuthPlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *AuthPlugin) Author() string {
	return p.author
}

// Initialize 初始化插件
func (p *AuthPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Auth plugin initialized")

	// 检查依赖
	satisfied, unresolved, err := checkDependenciesWithManager(p, ctx)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %v", err)
	}

	if !satisfied {
		return fmt.Errorf("unsatisfied dependencies: %v", unresolved)
	}

	return nil
}

// Start 启动插件
func (p *AuthPlugin) Start() error {
	p.context.LogInfo("Auth plugin started")

	// 检查依赖状态
	depStatus := p.CheckDependencies()
	for dep, satisfied := range depStatus {
		if satisfied {
			p.context.LogInfo("Dependency %s is satisfied", dep)
		} else {
			p.context.LogWarn("Dependency %s is NOT satisfied", dep)
		}
	}

	return nil
}

// Stop 停止插件
func (p *AuthPlugin) Stop() error {
	p.context.LogInfo("Auth plugin stopped")
	return nil
}

// Cleanup 清理插件
func (p *AuthPlugin) Cleanup() error {
	p.context.LogInfo("Auth plugin cleaned up")
	return nil
}

// Dependencies 返回插件依赖
func (p *AuthPlugin) Dependencies() []string {
	return p.dependencies
}

// CheckDependencies 检查插件依赖状态
func (p *AuthPlugin) CheckDependencies() map[string]bool {
	dependencies := p.Dependencies()
	result := make(map[string]bool)

	// 在实际实现中，这会与插件管理器通信检查依赖状态
	// 这里我们模拟检查
	for _, dep := range dependencies {
		// 模拟检查依赖是否满足
		result[dep] = isDependencySatisfied(dep)
	}

	return result
}

// BusinessPlugin 模拟一个业务插件，依赖认证和数据库插件
type BusinessPlugin struct {
	name        string
	version     string
	description string
	author      string
	context     types.PluginContext
	dependencies []string
}

// NewBusinessPlugin 创建业务插件实例
func NewBusinessPlugin() *BusinessPlugin {
	return &BusinessPlugin{
		name:        "business-plugin",
		version:     "1.0.0",
		description: "A business logic plugin that depends on auth and database plugins",
		author:      "urlDB Team",
		dependencies: []string{"auth-plugin", "database-plugin"}, // 依赖认证和数据库插件
	}
}

// Name 返回插件名称
func (p *BusinessPlugin) Name() string {
	return p.name
}

// Version 返回插件版本
func (p *BusinessPlugin) Version() string {
	return p.version
}

// Description 返回插件描述
func (p *BusinessPlugin) Description() string {
	return p.description
}

// Author 返回插件作者
func (p *BusinessPlugin) Author() string {
	return p.author
}

// Initialize 初始化插件
func (p *BusinessPlugin) Initialize(ctx types.PluginContext) error {
	p.context = ctx
	p.context.LogInfo("Business plugin initialized")

	// 检查依赖
	satisfied, unresolved, err := checkDependenciesWithManager(p, ctx)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %v", err)
	}

	if !satisfied {
		return fmt.Errorf("unsatisfied dependencies: %v", unresolved)
	}

	return nil
}

// Start 启动插件
func (p *BusinessPlugin) Start() error {
	p.context.LogInfo("Business plugin started at %s", time.Now().Format(time.RFC3339))

	// 检查依赖状态
	depStatus := p.CheckDependencies()
	for dep, satisfied := range depStatus {
		if satisfied {
			p.context.LogInfo("Dependency %s is satisfied", dep)
		} else {
			p.context.LogWarn("Dependency %s is NOT satisfied", dep)
		}
	}

	// 模拟业务逻辑
	p.context.LogInfo("Business plugin is processing data...")

	return nil
}

// Stop 停止插件
func (p *BusinessPlugin) Stop() error {
	p.context.LogInfo("Business plugin stopped")
	return nil
}

// Cleanup 清理插件
func (p *BusinessPlugin) Cleanup() error {
	p.context.LogInfo("Business plugin cleaned up")
	return nil
}

// Dependencies 返回插件依赖
func (p *BusinessPlugin) Dependencies() []string {
	return p.dependencies
}

// CheckDependencies 检查插件依赖状态
func (p *BusinessPlugin) CheckDependencies() map[string]bool {
	dependencies := p.Dependencies()
	result := make(map[string]bool)

	// 在实际实现中，这会与插件管理器通信检查依赖状态
	// 这里我们模拟检查
	for _, dep := range dependencies {
		// 模拟检查依赖是否满足
		result[dep] = isDependencySatisfied(dep)
	}

	return result
}

// 辅助函数：检查依赖是否满足（模拟实现）
func isDependencySatisfied(dependencyName string) bool {
	// 在实际实现中，这会检查依赖插件是否已注册并正在运行
	// 这里我们模拟总是满足依赖
	switch dependencyName {
	case "database-plugin", "auth-plugin":
		return true
	default:
		return false
	}
}

// 辅助函数：与插件管理器检查依赖（模拟实现）
func checkDependenciesWithManager(plugin types.Plugin, ctx types.PluginContext) (bool, []string, error) {
	// 在实际实现中，这会与插件管理器通信检查依赖
	// 这里我们模拟检查结果
	dependencies := plugin.Dependencies()
	unresolved := []string{}

	for _, dep := range dependencies {
		if !isDependencySatisfied(dep) {
			unresolved = append(unresolved, dep)
		}
	}

	return len(unresolved) == 0, unresolved, nil
}
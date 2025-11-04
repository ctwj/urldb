// Package market 实现插件市场功能
//
// 该包提供了与插件市场交互的功能，包括搜索、安装、更新和卸载插件。
//
// 主要组件:
//   - MarketClient: 与插件市场API交互的客户端
//   - MarketManager: 管理插件安装和更新的管理器
//   - PluginInfo: 插件信息结构
//
// 使用方法:
//   1. 创建 MarketClient 连接到插件市场
//   2. 创建 MarketManager 管理本地插件
//   3. 使用 Search, Install, Update 等方法管理插件
//
// 注意事项:
//   - 插件市场功能需要网络连接
//   - 安装插件前应验证校验和
//   - 需要处理插件依赖关系
package market
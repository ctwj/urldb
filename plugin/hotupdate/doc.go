// Package hotupdate 实现插件热更新功能
//
// 该包提供了插件文件监视和自动更新功能，允许在不重启主应用的情况下
// 更新插件实现。
//
// 主要组件:
//   - PluginWatcher: 监视插件文件变化
//   - PluginUpdater: 执行插件更新操作
//
// 使用方法:
//   1. 创建 PluginUpdater 实例
//   2. 调用 StartUpdaterWithWatcher 开始监视
//   3. 插件文件发生变化时会自动触发更新
//
// 注意事项:
//   - 插件热更新依赖于操作系统的文件监视机制
//   - 更新过程中插件会短暂停止服务
//   - 建议在低峰期进行插件更新
package hotupdate
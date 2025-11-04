// Package version 实现插件版本管理功能
//
// 该包提供了插件版本管理的核心功能，包括版本注册、查询、比较和兼容性检查。
//
// 主要组件:
//   - VersionManager: 版本管理器，负责版本的增删改查
//   - PluginVersion: 插件版本信息结构
//   - VersionComparison: 版本比较结果
//   - CompatibilityInfo: 兼容性信息
//
// 特性:
//   - 支持语义化版本控制 (Semantic Versioning)
//   - 版本历史记录和查询
//   - 版本兼容性检查
//   - 版本状态管理（激活、废弃等）
//
// 使用方法:
//   1. 创建 VersionManager 实例
//   2. 使用 RegisterVersion 注册新版本
//   3. 使用 GetVersion/GetLatestVersion 查询版本
//   4. 使用 CompareVersions 比较版本
//   5. 使用 CheckVersionCompatibility 检查兼容性
//
// 注意事项:
//   - 版本信息会持久化到数据库
//   - 支持预发布版本和构建元数据
//   - 提供版本号递增工具函数
package version
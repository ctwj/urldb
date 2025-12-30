/// <reference path="../pb_data/types.d.ts" />

/**
 * test_migration_plugin 钩子
 * 创建时间: 2025-12-30 22:15:00
 *
 * @name test_migration_plugin
 * @display_name 测试迁移插件
 * @author URLDB开发团队
 * @description 测试压缩包插件的迁移功能
 * @version 1.0.0
 * @category test
 * @license MIT
 */

console.log("测试迁移插件已加载");

// 添加测试路由
routerAdd("GET", "/api/test-migration-plugin", (e) => {
    return e.json(200, {
        message: "测试迁移插件正常工作",
        plugin: "test_migration_plugin",
        timestamp: new Date().toISOString()
    });
});
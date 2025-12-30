/// <reference path="../pb_data/types.d.ts" />

// 测试迁移文件
// 这个迁移将在 migrations 目录中被处理

migrate(
    // up 迁移函数
    (app) => {
        console.log("执行测试迁移 - up: 创建测试表");

        try {
            // 使用 db.raw 执行 SQL 创建测试表
            const result = app.db().raw(`
                CREATE TABLE IF NOT EXISTS migration_test_table (
                    id INTEGER PRIMARY KEY,
                    message TEXT,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            `);

            console.log("迁移测试表创建成功");
            return null; // 成功返回 null
        } catch (error) {
            console.log("迁移测试表创建失败: " + error.message);
            return error;
        }
    },
    // down 迁移函数
    (app) => {
        console.log("执行测试迁移 - down: 删除测试表");

        try {
            const result = app.db().raw("DROP TABLE IF EXISTS migration_test_table");
            console.log("迁移测试表删除成功");
            return null; // 成功返回 null
        } catch (error) {
            console.log("迁移测试表删除失败: " + error.message);
            return error;
        }
    }
);

console.log("测试迁移文件加载完成");
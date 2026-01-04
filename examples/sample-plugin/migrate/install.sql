-- 测试迁移插件的安装脚本
-- 创建测试表用于验证迁移功能

CREATE TABLE IF NOT EXISTS plugin_migration_test (
    id INTEGER PRIMARY KEY,
    plugin_name VARCHAR(100) NOT NULL,
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入测试数据（明确指定ID值）
INSERT INTO plugin_migration_test (id, plugin_name, message)
VALUES (1, 'test_migration_plugin', '插件安装时创建的测试数据');
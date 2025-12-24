package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ctwj/urldb/utils"
)

// PluginUtils 插件系统工具函数
type PluginUtils struct{}

// NewPluginUtils 创建插件工具实例
func NewPluginUtils() *PluginUtils {
	return &PluginUtils{}
}

// CreatePluginTemplate 创建插件模板文件
func (p *PluginUtils) CreatePluginTemplate(pluginName, pluginType string) error {
	var template string
	var filename string

	switch pluginType {
	case "hook":
		template = p.generateHookTemplate(pluginName)
		filename = filepath.Join("./hooks", pluginName+".pb.js")
	case "migration":
		template = p.generateMigrationTemplate(pluginName)
		filename = filepath.Join("./migrations", pluginName+".js")
	default:
		return fmt.Errorf("Unsupported plugin type: %s", pluginType)
	}

	// 确保目录存在
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("Failed to create directory %s: %v", dir, err)
	}

	// 写入模板文件
	if err := os.WriteFile(filename, []byte(template), 0644); err != nil {
		return fmt.Errorf("Failed to write template file %s: %v", filename, err)
	}

	utils.Info("Plugin template created: %s", filename)
	return nil
}

// generateHookTemplate 生成钩子模板
func (p *PluginUtils) generateHookTemplate(pluginName string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var template strings.Builder
	template.WriteString(`/// <reference path="../pb_data/types.d.ts" />

/**
 * ` + pluginName + ` 钩子
 * 创建时间: ` + timestamp + `
 */

// 示例：监听 URL 添加事件
onURLAdd((e) => {
    console.log("URL 添加触发:", e.url.url);

    // 在这里添加你的自定义逻辑
    // 例如：自动分类、标签提取、通知等

    return e.next();
});

// 示例：监听用户登录事件
onUserLogin((e) => {
    console.log("用户登录:", e.user.username);

    // 在这里添加登录后处理逻辑
    // 例如：日志记录、欢迎消息、权限检查等

    return e.next();
});

// 示例：添加自定义路由
routerAdd("GET", "/api/custom", (e) => {
    return e.json(200, {
        message: "来自 ` + pluginName + ` 插件的自定义 API",
        timestamp: new Date().toISOString()
    });
});

// 示例：添加定时任务
cronAdd("` + pluginName + `_task", "0 */6 * * *", () => {
    console.log("执行定时任务: ` + pluginName + `");
    // 在这里添加定时任务逻辑
});
`)
	return template.String()
}

// generateMigrationTemplate 生成迁移模板
func (p *PluginUtils) generateMigrationTemplate(pluginName string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var template strings.Builder
	template.WriteString(`/**
 * ` + pluginName + ` 数据库迁移
 * 创建时间: ` + timestamp + `
 */

migrate(
    // up 迁移
    (app) => {
        console.log("执行向上迁移: ` + pluginName + `");

        // 在这里添加向上迁移逻辑
        // 例如：创建表、添加字段、插入数据等

        // 示例：创建新表
        /*
        const createTable = "CREATE TABLE IF NOT EXISTS example_table (" +
            "id INTEGER PRIMARY KEY AUTOINCREMENT," +
            "name TEXT NOT NULL," +
            "created_at DATETIME DEFAULT CURRENT_TIMESTAMP" +
            ")";

        app.db.run(createTable);
        */
    },

    // down 回滚
    (app) => {
        console.log("执行回滚迁移: ` + pluginName + `");

        // 在这里添加回滚逻辑
        // 例如：删除表、移除字段、清理数据等

        // 示例：删除表
        /*
        app.db.run("DROP TABLE IF EXISTS example_table");
        */
    }
);
`)
	return template.String()
}

// ListPlugins 列出所有插件文件
func (p *PluginUtils) ListPlugins() error {
	utils.Info("=== 钩子文件 ===")
	if err := p.listFilesInDir("./hooks", "*.pb.js"); err != nil {
		return err
	}

	utils.Info("=== 迁移文件 ===")
	if err := p.listFilesInDir("./migrations", "*.js"); err != nil {
		return err
	}

	return nil
}

// listFilesInDir 列出目录中的文件
func (p *PluginUtils) listFilesInDir(dir, pattern string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Info("目录不存在: %s", dir)
			return nil
		}
		return err
	}

	if len(files) == 0 {
		utils.Info("目录为空: %s", dir)
		return nil
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				utils.Error("无法获取文件信息: %v", err)
				continue
			}
			utils.Info("  - %s (大小: %d, 修改时间: %s)",
				file.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
		}
	}

	return nil
}

// ValidatePlugin 验证插件文件
func (p *PluginUtils) ValidatePlugin(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("Plugin file does not exist: %s", filePath)
	}

	// 检查文件扩展名
	ext := filepath.Ext(filePath)
	if ext != ".js" && ext != ".ts" {
		return fmt.Errorf("Plugin file must have .js or .ts extension: %s", filePath)
	}

	// 检查文件大小
	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("Failed to get plugin file info: %v", err)
	}

	if info.Size() == 0 {
		return fmt.Errorf("Plugin file is empty: %s", filePath)
	}

	if info.Size() > 10*1024*1024 { // 10MB
		return fmt.Errorf("Plugin file is too large (>10MB): %s", filePath)
	}

	utils.Info("Plugin file validation passed: %s", filePath)
	return nil
}

// GetPluginStats 获取插件统计信息
func (p *PluginUtils) GetPluginStats() map[string]interface{} {
	stats := make(map[string]interface{})

	// 统计钩子文件
	hooksDir := "./hooks"
	hooksCount := p.countFiles(hooksDir, "*.pb.js")
	stats["hooks_count"] = hooksCount

	// 统计迁移文件
	migrationsDir := "./migrations"
	migrationsCount := p.countFiles(migrationsDir, "*.js")
	stats["migrations_count"] = migrationsCount

	// 统计类型文件
	typesDir := "./pb_data"
	typesFile := filepath.Join(typesDir, "types.d.ts")
	if _, err := os.Stat(typesFile); err == nil {
		stats["types_file_exists"] = true
	} else {
		stats["types_file_exists"] = false
	}

	stats["last_updated"] = time.Now().Format("2006-01-02 15:04:05")

	return stats
}

// countFiles 统计目录中匹配模式的文件数量
func (p *PluginUtils) countFiles(dir, pattern string) int {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	count := 0
	for _, file := range files {
		if !file.IsDir() {
			// 简单的模式匹配
			if pattern == "*.pb.js" && len(file.Name()) > 6 && file.Name()[len(file.Name())-6:] == ".pb.js" {
				count++
			} else if pattern == "*.js" && len(file.Name()) > 3 && file.Name()[len(file.Name())-3:] == ".js" {
				count++
			}
		}
	}

	return count
}
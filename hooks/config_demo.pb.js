/// <reference path="../pb_data/types.d.ts" />

/**
 * 插件配置演示
 * 展示如何在插件中使用配置系统
 */

// 1. 从系统配置中读取插件配置
function getPluginConfig() {
    // 从数据库中读取插件配置
    const configRecord = app.db.queryOne(`
        SELECT config_json FROM plugin_configs
        WHERE plugin_name = 'config_demo' AND enabled = true
    `);

    if (!configRecord) {
        // 返回默认配置
        return {
            enabled: true,
            log_level: "info",
            max_retries: 3,
            timeout: 5000,
            custom_settings: {
                greeting_message: "Hello from config demo!",
                feature_flags: {
                    enable_notifications: true,
                    enable_analytics: false
                }
            }
        };
    }

    return JSON.parse(configRecord.config_json);
}

// 2. 从环境配置中读取系统级配置
function getSystemConfig() {
    return {
        debug_mode: app.config().GetBool("PLUGIN_DEBUG"),
        hooks_dir: app.config().GetString("PLUGIN_HOOKS_DIR"),
        vm_pool_size: app.config().GetInt("PLUGIN_VM_POOL_SIZE")
    };
}

// 监听URL添加事件，演示配置使用
onURLAdd((e) => {
    const pluginConfig = getPluginConfig();
    const systemConfig = getSystemConfig();

    // 根据配置决定是否处理
    if (!pluginConfig.enabled) {
        console.log("[配置演示] 插件已禁用，跳过处理");
        return e.next();
    }

    // 根据日志级别记录
    if (pluginConfig.log_level === "debug" || systemConfig.debug_mode) {
        console.log("[配置演示] Debug模式 - URL:", e.url.url);
        console.log("[配置演示] 插件配置:", JSON.stringify(pluginConfig, null, 2));
        console.log("[配置演示] 系统配置:", JSON.stringify(systemConfig, null, 2));
    }

    // 使用配置中的超时设置
    const startTime = Date.now();

    try {
        // 模拟处理逻辑
        processURLWithConfig(e.url, pluginConfig);

        const processingTime = Date.now() - startTime;
        if (processingTime > pluginConfig.timeout) {
            console.warn("[配置演示] 处理超时:", processingTime, "ms >", pluginConfig.timeout, "ms");
        }

    } catch (error) {
        // 使用配置中的重试次数
        if (shouldRetry(error, pluginConfig.max_retries)) {
            console.log("[配置演示] 将重试，剩余次数:", pluginConfig.max_retries);
        }
    }

    return e.next();
});

// 监听用户登录事件，演示个性化配置
onUserLogin((e) => {
    const pluginConfig = getPluginConfig();

    // 根据用户偏好显示个性化消息
    if (pluginConfig.custom_settings.feature_flags.enable_notifications) {
        const greeting = pluginConfig.custom_settings.greeting_message;
        console.log("[配置演示] 发送欢迎消息:", greeting, "给用户:", e.user.username);

        // 可以从用户偏好表中读取个性化设置
        const userPrefs = app.db.queryOne(`
            SELECT * FROM user_preferences
            WHERE user_id = ? AND category = 'notifications'
        `, [e.user.id]);

        if (userPrefs) {
            const prefs = JSON.parse(userPrefs.preferences_json);
            console.log("[配置演示] 用户通知偏好:", prefs);
        }
    }

    return e.next();
});

// 添加配置管理API
routerAdd("GET", "/api/plugins/config_demo/config", (e) => {
    try {
        const config = getPluginConfig();
        const systemConfig = getSystemConfig();

        return e.json(200, {
            success: true,
            data: {
                plugin_config: config,
                system_config: systemConfig,
                config_source: "database"
            }
        });
    } catch (error) {
        console.error("[配置演示] 获取配置失败:", error);
        return e.json(500, {
            success: false,
            error: "获取配置失败"
        });
    }
});

// 添加配置更新API
routerAdd("POST", "/api/plugins/config_demo/config", (e) => {
    try {
        const { config } = e.body;

        if (!config) {
            return e.json(400, {
                success: false,
                error: "配置不能为空"
            });
        }

        // 验证配置
        if (!validatePluginConfig(config)) {
            return e.json(400, {
                success: false,
                error: "配置格式无效"
            });
        }

        // 更新数据库中的配置
        app.db.execute(`
            INSERT OR REPLACE INTO plugin_configs (plugin_name, config_json, enabled, updated_at)
            VALUES (?, ?, ?, datetime('now'))
        `, ['config_demo', JSON.stringify(config), config.enabled !== false]);

        console.log("[配置演示] 配置已更新:", JSON.stringify(config, null, 2));

        return e.json(200, {
            success: true,
            message: "配置更新成功"
        });
    } catch (error) {
        console.error("[配置演示] 更新配置失败:", error);
        return e.json(500, {
            success: false,
            error: "更新配置失败"
        });
    }
});

// 添加用户偏好管理API
routerAdd("POST", "/api/plugins/config_demo/preferences", (e) => {
    try {
        const { user_id, category, preferences } = e.body;

        if (!user_id || !category || !preferences) {
            return e.json(400, {
                success: false,
                error: "用户ID、类别和偏好设置不能为空"
            });
        }

        // 保存用户偏好
        app.db.execute(`
            INSERT OR REPLACE INTO user_preferences
            (user_id, category, preferences_json, updated_at)
            VALUES (?, ?, ?, datetime('now'))
        `, [user_id, category, JSON.stringify(preferences)]);

        console.log("[配置演示] 用户偏好已保存:", user_id, category);

        return e.json(200, {
            success: true,
            message: "用户偏好保存成功"
        });
    } catch (error) {
        console.error("[配置演示] 保存用户偏好失败:", error);
        return e.json(500, {
            success: false,
            error: "保存用户偏好失败"
        });
    }
});

// 工具函数：使用配置处理URL
function processURLWithConfig(url, config) {
    console.log("[配置演示] 使用配置处理URL:", url.url);

    // 根据配置决定处理方式
    if (config.custom_settings.feature_flags.enable_analytics) {
        // 启用分析功能
        recordAnalytics(url);
    }

    // 可以根据配置执行不同的业务逻辑
    return true;
}

// 工具函数：验证插件配置
function validatePluginConfig(config) {
    const requiredFields = ['enabled', 'log_level', 'max_retries', 'timeout'];

    for (const field of requiredFields) {
        if (config[field] === undefined || config[field] === null) {
            console.error("[配置演示] 缺少必需字段:", field);
            return false;
        }
    }

    // 验证数据类型
    if (typeof config.enabled !== 'boolean') {
        console.error("[配置演示] enabled 必须是布尔值");
        return false;
    }

    if (typeof config.max_retries !== 'number' || config.max_retries < 0) {
        console.error("[配置演示] max_retries 必须是非负整数");
        return false;
    }

    return true;
}

// 工具函数：判断是否应该重试
function shouldRetry(error, maxRetries) {
    // 简单的重试逻辑
    return maxRetries > 0 && error.retriable !== false;
}

// 工具函数：记录分析数据
function recordAnalytics(url) {
    app.db.execute(`
        INSERT INTO url_stats (url_id, domain, category, created_at)
        VALUES (?, ?, ?, datetime('now'))
    `, [url.id, extractDomain(url.url), url.category || 'uncategorized']);
}

// 工具函数：提取域名
function extractDomain(url) {
    try {
        const urlObj = new URL(url);
        return urlObj.hostname.replace(/^www\./, '');
    } catch {
        return 'unknown';
    }
}

// 添加定时任务：配置清理
cronAdd("config_demo_cleanup", "0 2 * * *", () => {
    console.log("[配置演示] 执行配置清理任务");

    try {
        // 清理过期的日志
        const daysToKeep = 30;
        app.db.execute(`
            DELETE FROM plugin_logs
            WHERE created_at < datetime('now', '-${daysToKeep} days')
        `);

        console.log("[配置演示] 配置清理完成");
    } catch (error) {
        console.error("[配置演示] 配置清理失败:", error);
    }
});

console.log("[配置演示] 插件初始化完成 - 配置系统演示");
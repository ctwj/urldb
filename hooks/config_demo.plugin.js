/// <reference path="../pb_data/types.d.ts" />

//  * // 非必填项，使用 @optional
//  * // 默认值，使用 @default
//  * // @field {类型} 名字 表单名字 placeholder [@optional] [@default]


/**
 * config_demo 钩子
 * 创建时间: 2025-12-25 23:23:19
 *
 * @name config_demo
 * @display_name 配置演示插件
 * @author URLDB开发团队
 * @description 演示插件配置功能的示例插件，包含Webhook通知、日志级别设置等配置选项
 * @version 1.0.1
 * @category demo
 * @license MIT
 *
 * @config
 * @field {string} webhook_url Webhook URL "通知发送的Webhook地址" @default "https://hooks.slack.com/services/YOUR/DEFAULT/WEBHOOK"
 * @field {boolean} enable_notification 启用通知 "是否启用通知功能" @default true
 * @field {number} retry_count 重试次数 "通知失败时的重试次数" @default 3
 * @field {select} log_level 日志级别 "日志输出级别" ["debug", "info", "warn", "error"] @default "info"
 * @field {text} custom_message 自定义消息 "自定义通知消息内容" @optional @default "这是来自 config_demo 插件的默认消息"
 * @config
 */

// 提取的配置处理函数
function processConfigDemo() {
    try {
        // 获取插件配置
        const config = getPluginConfig("config_demo");

        // 最简化处理，避免所有console.log
        if (config) {
            return {
                success: true,
                config: config,
                timestamp: new Date().toISOString()
            };
        } else {
            return {
                success: false,
                error: "未找到插件配置",
                timestamp: new Date().toISOString()
            };
        }
    } catch (error) {
        return {
            success: false,
            error: error.message,
            timestamp: new Date().toISOString()
        };
    }
}

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

// 示例：添加自定义路由 - 获取配置信息
routerAdd("GET", "/api/config-demo", (e) => {
    const result = processConfigDemo();

    return e.json(200, {
        message: "来自 config_demo 插件的自定义 API",
        data: result,
        timestamp: new Date().toISOString()
    });
});

// 添加新的路由 - 手动触发配置处理
routerAdd("POST", "/api/config-demo/refresh", (e) => {
    console.log("手动触发配置处理");

    const result = processConfigDemo();

    return e.json(200, {
        message: "配置处理完成",
        data: result,
        timestamp: new Date().toISOString()
    });
});

// 添加新的路由 - 获取配置摘要
routerAdd("GET", "/api/config-demo/summary", (e) => {
    const config = getPluginConfig("config_demo");

    if (!config) {
        return e.json(404, {
            error: "未找到插件配置",
            timestamp: new Date().toISOString()
        });
    }

    return e.json(200, {
        plugin_name: "config_demo",
        display_name: "配置演示插件",
        version: "1.0.1",
        webhook_configured: config.webhook_url && config.webhook_url !== "https://hooks.slack.com/services/YOUR/DEFAULT/WEBHOOK",
        notification_enabled: config.enable_notification || false,
        log_level: config.log_level || "info",
        retry_count: config.retry_count || 0,
        timestamp: new Date().toISOString()
    });
});

// 示例：添加定时任务 - 重新启用，每1分钟执行一次
cronAdd("config_demo_task", "*/1 * * * *", () => {
    console.log("执行定时任务: config_demo - 每1分钟执行一次");

    try {
        // 调用提取的函数
        const result = processConfigDemo();

        // 记录执行结果
        if (result && result.success) {
            console.log("定时任务执行成功，配置已处理");
        } else {
            console.log("定时任务执行失败:", (result && result.error) || "未知错误");
        }
    } catch (error) {
        console.log("定时任务执行异常:", error.message || error);
    }
});

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
 * @version 1.0.0
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
routerAdd("GET", "/api/config-demo", (e) => {
    return e.json(200, {
        message: "来自 config_demo 插件的自定义 API",
        timestamp: new Date().toISOString()
    });
});

// 示例：添加定时任务
cronAdd("config_demo_task", "0 */6 * * *", () => {
    console.log("执行定时任务: config_demo");
    // 在这里添加定时任务逻辑
});

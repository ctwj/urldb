/// <reference path="../pb_data/types.d.ts" />

/**
 * notification_center 钩子
 * 创建时间: 2025-12-25 23:23:19
 *
 * @name notification_center
 * @display_name 通知中心
 * @author URLDB开发团队
 * @description 统一的通知管理中心，支持邮件、短信等多种通知方式，可配置通知时间和频率
 * @version 1.0.0
 * @category notification
 * @license MIT
 *
 * @config
 * @field {string} email_address 邮箱地址 "接收通知的邮箱地址" @optional
 * @field {boolean} enable_email 启用邮件通知 "是否启用邮件通知功能" @default false
 * @field {boolean} enable_sms 启用短信通知 "是否启用短信通知功能" @default false
 * @field {select} notification_time 通知时间 "偏好接收通知的时间" ["immediate", "daily", "weekly"] @default "immediate"
 * @field {number} max_notifications_per_day 每日最大通知数 "限制每日发送的通知数量" @default 10
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
routerAdd("GET", "/api/notification-center", (e) => {
    return e.json(200, {
        message: "来自 notification_center 插件的自定义 API",
        timestamp: new Date().toISOString()
    });
});

// 示例：添加定时任务
cronAdd("notification_center_task", "0 */6 * * *", () => {
    console.log("执行定时任务: notification_center");
    // 在这里添加定时任务逻辑
});

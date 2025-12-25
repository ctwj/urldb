/// <reference path="../pb_data/types.d.ts" />

/**
 * test_demo 钩子
 * 创建时间: 2025-12-25 08:16:09
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
        message: "来自 test_demo 插件的自定义 API",
        timestamp: new Date().toISOString()
    });
});

// 示例：添加定时任务
cronAdd("test_demo_task", "0 */6 * * *", () => {
    console.log("执行定时任务: test_demo");
    // 在这里添加定时任务逻辑
});

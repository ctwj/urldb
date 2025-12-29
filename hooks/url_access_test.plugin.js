/*
 * URL访问测试插件
 * 用于测试 onURLAccess 事件钩子功能
 */

// 插件元信息
plugin({
    name: "url_access_test",
    version: "1.0.0",
    description: "URL访问事件测试插件",
    author: "URLDB Team"
});

// 记录插件加载日志
log("info", "URL访问测试插件已加载", "url_access_test");

// 测试 onURLAccess 事件钩子
onURLAccess(function(event) {
    log("info", "=== onURLAccess 事件触发 ===", "url_access_test");
    log("info", "URL ID: " + event.url.id, "url_access_test");
    log("info", "URL Title: " + event.url.title, "url_access_test");
    log("info", "URL: " + event.url.url, "url_access_test");
    log("info", "=== onURLAccess 事件处理完成 ===", "url_access_test");
});

// 注册一个测试路由
router.get("/api/url-access-test", function() {
    return {
        success: true,
        message: "URL访问测试插件运行正常",
        timestamp: new Date().toISOString(),
        plugin: "url_access_test"
    };
});

log("info", "URL访问测试插件初始化完成", "url_access_test");
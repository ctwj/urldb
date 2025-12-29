/**
 * 用户登录事件测试插件
 * 用于测试 onUserLogin 钩子函数的功能
 *
 * @hooks ["onUserLogin"]
 */

// 用户登录事件处理
onUserLogin((e) => {
    log("info", "=== onUserLogin 事件触发 ===", "login_test");

    try {
        // 记录用户基本信息
        try {
            log("info", `用户登录: ${e.user.username} (ID: ${e.user.id})`, "login_test");
            log("info", `用户邮箱: ${e.user.email || '未设置'}`, "login_test");
            log("info", `用户角色: ${e.user.role || '普通用户'}`, "login_test");
            log("info", `用户状态: ${e.user.is_active ? '活跃' : '禁用'}`, "login_test");
        } catch (err) {
            log("error", `用户信息处理错误: ${err.message}`, "login_test");
        }

        // 记录登录事件数据
        if (e.data) {
            log("info", `登录IP: ${e.data.ip || '未知'}`, "login_test");
            log("info", `用户代理: ${e.data.user_agent || '未知'}`, "login_test");
            log("info", `登录时间: ${e.data.login_time || new Date().toISOString()}`, "login_test");
        }

        // 记录应用程序信息
        if (e.app) {
            log("info", `应用名称: ${e.app.name || 'URLDB'}`, "login_test");
            log("info", `应用版本: ${e.app.version || '1.0.0'}`, "login_test");
        }

        // 模拟一些业务逻辑
        log("info", "正在检查用户权限...", "login_test");

        // 模拟权限检查
        if (e.user && e.user.role === 'admin') {
            log("info", "管理员用户登录，拥有所有权限", "login_test");
        } else {
            log("info", "普通用户登录，拥有基本权限", "login_test");
        }

        // 模拟登录统计
        log("info", "更新用户登录统计信息", "login_test");

        // 模拟安全检查
        if (e.data && e.data.ip) {
            const ip = e.data.ip;
            if (ip.startsWith('192.168.') || ip.startsWith('127.') || ip.startsWith('localhost')) {
                log("info", "内网IP登录，安全级别: 低", "login_test");
            } else {
                log("warn", `外网IP登录: ${ip}，安全级别: 中`, "login_test");
            }
        }

        log("info", "=== onUserLogin 事件处理完成 ===", "login_test");

    } catch (error) {
        log("error", `onUserLogin 处理错误: ${error.message}`, "login_test");
    }
});

// 添加一个测试路由，用于手动触发登录事件测试
routerAdd('get', '/api/login-test', (ctx) => {
    log("info", "访问登录测试接口", "login_test");

    // 返回测试信息
    ctx.body = {
        success: true,
        message: "登录测试插件运行中",
        data: {
            plugin: "login_test",
            hooks: ["onUserLogin"],
            description: "用户登录事件测试插件",
            test_method: "使用 admin/password1 登录系统以触发 onUserLogin 事件"
        }
    };
});

log("info", "登录测试插件已加载", "login_test");
log("info", "请使用账号 admin/password1 登录系统以测试 onUserLogin 功能", "login_test");
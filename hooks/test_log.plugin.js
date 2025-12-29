/// <reference path="../pb_data/types.d.ts" />

/**
 * test_log 钩子
 * 创建时间: 2025-12-29 00:17:00
 *
 * @name test_log
 * @display_name 测试日志插件
 * @author URLDB开发团队
 * @description 测试新的log函数是否正常工作
 * @version 1.0.0
 * @category demo
 * @license MIT
 */

// 测试基本log函数
routerAdd("GET", "/api/test-log", (e) => {
    log("info", "这是一条info级别的测试日志", "test_log");
    log("warn", "这是一条warn级别的测试日志", "test_log");
    log("error", "这是一条error级别的测试日志", "test_log");
    log("debug", "这是一条debug级别的测试日志", "test_log");

    return e.json(200, {
        message: "测试日志已记录",
        timestamp: new Date().toISOString()
    });
});

// 测试定时任务中的log函数
cronAdd("test_log_task", "*/2 * * * *", () => {
    log("info", "定时任务测试日志：每2分钟执行一次", "test_log");
});
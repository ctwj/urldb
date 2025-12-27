/// <reference path="../types.d.ts" />

// 简单的JavaScript执行测试插件
console.log("JavaScript插件正在加载...");

// 测试基础函数调用
log("info", "JavaScript插件日志系统正常工作");

// 测试JSON处理
const testData = { message: "Hello from JavaScript!", timestamp: timestamp() };
const jsonString = jsonStringify(testData);
log("info", "JSON序列化测试: " + jsonString);

// 测试钩子注册
onURLAdd((e) => {
    log("info", "URL添加钩子被触发: " + e.url);
    console.log("URL添加事件:", e);
});

// 测试定时任务
cron.add("test-job", "*/1 * * * *", () => {
    log("info", "定时任务执行: " + new Date().toISOString());
});

console.log("JavaScript插件加载完成!");
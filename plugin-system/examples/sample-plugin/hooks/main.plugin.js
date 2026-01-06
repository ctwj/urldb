/**
 * @name sample-plugin
 * @display_name 压缩包演示插件
 * @description A sample plugin for demonstration
 * @version 1.0.0
 * @author URLDB Team
 * 
 * @config
 * @field {string} webhook_url Webhook URL "通知发送的Webhook地址" @default "https://hooks.slack.com/services/YOUR/DEFAULT/WEBHOOK"
 * @field {boolean} enable_notification 启用通知 "是否启用通知功能" @default true
 * @field {number} retry_count 重试次数 "通知失败时的重试次数" @default 3
 * @field {select} log_level 日志级别 "日志输出级别" ["debug", "info", "warn", "error"] @default "info"
 * @field {text} custom_message 自定义消息 "自定义通知消息内容" @optional @default "这是来自 config_demo 插件的默认消息"
 * @config
 */

// Plugin initialization
console.log('Sample plugin hook loaded!');

// Register URL event handler
onURLAdd((e) => {
    console.log('Sample plugin: URL added', e.url);
    e.next();
});

// Register API request handler
onAPIRequest((e) => {
    if (e.path.startsWith('/api/sample')) {
        console.log('Sample plugin: API request intercepted', e.method, e.path);
    }
    e.next();
});

// Add a custom route
routerAdd('GET', '/api/sample/hello', (ctx) => {
    ctx.json({
        success: true,
        message: getConfig('custom_message') || 'Hello from sample plugin!',
        timestamp: new Date().toISOString()
    });
});

// Add a cron job
cronAdd('sample-cleanup', '0 */5 * * * *', () => {
    console.log('Sample plugin: Running cleanup task');
});

console.log('Sample plugin initialization completed');
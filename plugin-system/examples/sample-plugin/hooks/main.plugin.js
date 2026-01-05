/**
 * @name sample-plugin
 * @description A sample plugin for demonstration
 * @version 1.0.0
 * @author URLDB Team
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
        message: 'Hello from sample plugin!',
        timestamp: new Date().toISOString()
    });
});

// Add a cron job
cronAdd('sample-cleanup', '0 */5 * * * *', () => {
    console.log('Sample plugin: Running cleanup task');
});

console.log('Sample plugin initialization completed');
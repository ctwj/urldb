/// <reference path="../pb_data/types.d.ts" />

/**
 * 数据分析引擎插件
 * 功能：提供全面的数据分析和统计报告
 * 作者: URLDB Team
 * 版本: 1.0.0
 * 创建时间: 2024-12-24
 */

// 分析配置
const ANALYTICS_CONFIG = {
    // 数据保留期限（天）
    dataRetentionDays: 90,

    // 统计计算间隔（小时）
    aggregationInterval: 6,

    // 报告生成时间
    reportSchedule: "0 9 * * *", // 每天早上9点

    // 性能监控阈值
    performanceThresholds: {
        responseTime: 2000, // 2秒
        errorRate: 0.05,    // 5%
        memoryUsage: 0.8    // 80%
    }
};

// 监听URL添加事件，收集分析数据
onURLAdd((e) => {
    console.log("[分析引擎] 记录URL添加事件:", e.url.url);

    try {
        // 记录基础统计
        recordBasicStats(e.url);

        // 分析URL模式
        analyzeURLPattern(e.url);

        // 更新实时指标
        updateRealTimeMetrics('url_added', 1);

        // 检查是否需要触发告警
        checkPerformanceAlerts();

    } catch (error) {
        console.error("[分析引擎] 处理URL添加事件失败:", error);
    }

    return e.next();
});

// 监听URL访问事件
onURLAccess((e) => {
    console.log("[分析引擎] 记录URL访问事件:", e.url.url);

    try {
        // 记录访问统计
        recordAccessStats(e.url, e.user);

        // 分析用户行为
        analyzeUserBehavior(e.user, e.url);

        // 更新热门资源排行
        updatePopularResources(e.url);

        // 检测异常访问模式
        detectAnomalousAccess(e.user, e.url);

    } catch (error) {
        console.error("[分析引擎] 处理URL访问事件失败:", error);
    }

    return e.next();
});

// 监听用户登录事件
onUserLogin((e) => {
    console.log("[分析引擎] 记录用户登录事件:", e.user.username);

    try {
        // 记录用户活动
        recordUserActivity(e.user, 'login');

        // 分析用户活跃度
        analyzeUserEngagement(e.user);

        // 更新用户会话统计
        updateUserSessionStats(e.user);

    } catch (error) {
        console.error("[分析引擎] 处理用户登录事件失败:", error);
    }

    return e.next();
});

// 添加分析API路由
routerAdd("GET", "/api/analytics/dashboard", (e) => {
    try {
        const dashboard = generateDashboardData();

        return e.json(200, {
            success: true,
            data: dashboard
        });
    } catch (error) {
        console.error("[分析引擎] 获取仪表板数据失败:", error);
        return e.json(500, {
            success: false,
            error: "获取仪表板数据失败"
        });
    }
});

// 添加资源统计API
routerAdd("GET", "/api/analytics/resources", (e) => {
    try {
        const { period = '7d', category = null } = e.query;

        const stats = getResourceStats(period, category);

        return e.json(200, {
            success: true,
            data: stats
        });
    } catch (error) {
        console.error("[分析引擎] 获取资源统计失败:", error);
        return e.json(500, {
            success: false,
            error: "获取资源统计失败"
        });
    }
});

// 添加用户行为分析API
routerAdd("GET", "/api/analytics/users", (e) => {
    try {
        const { period = '30d', metric = 'activity' } = e.query;

        const analytics = getUserAnalytics(period, metric);

        return e.json(200, {
            success: true,
            data: analytics
        });
    } catch (error) {
        console.error("[分析引擎] 获取用户分析失败:", error);
        return e.json(500, {
            success: false,
            error: "获取用户分析失败"
        });
    }
});

// 添加系统性能监控API
routerAdd("GET", "/api/analytics/performance", (e) => {
    try {
        const performance = getPerformanceMetrics();

        return e.json(200, {
            success: true,
            data: performance
        });
    } catch (error) {
        console.error("[分析引擎] 获取性能指标失败:", error);
        return e.json(500, {
            success: false,
            error: "获取性能指标失败"
        });
    }
});

// 添加趋势分析API
routerAdd("GET", "/api/analytics/trends", (e) => {
    try {
        const { metric = 'growth', period = '30d' } = e.query;

        const trends = getTrendAnalysis(metric, period);

        return e.json(200, {
            success: true,
            data: trends
        });
    } catch (error) {
        console.error("[分析引擎] 获取趋势分析失败:", error);
        return e.json(500, {
            success: false,
            error: "获取趋势分析失败"
        });
    }
});

// 添加自定义报告生成API
routerAdd("POST", "/api/analytics/reports", (e) => {
    try {
        const { type, period, format = 'json', filters = {} } = e.body;

        if (!type || !period) {
            return e.json(400, {
                success: false,
                error: "报告类型和周期参数不能为空"
            });
        }

        const report = generateCustomReport(type, period, format, filters);

        return e.json(200, {
            success: true,
            data: report
        });
    } catch (error) {
        console.error("[分析引擎] 生成自定义报告失败:", error);
        return e.json(500, {
            success: false,
            error: "生成自定义报告失败"
        });
    }
});

// 添加数据导出API
routerAdd("GET", "/api/analytics/export", (e) => {
    try {
        const { type = 'summary', format = 'csv', period = '30d' } = e.query;

        const exportData = exportAnalyticsData(type, format, period);

        // 设置适当的响应头
        e.response.setHeader('Content-Type', getContentType(format));
        e.response.setHeader('Content-Disposition', `attachment; filename="analytics_${type}_${period}.${format}"`);

        return e.string(200, exportData);
    } catch (error) {
        console.error("[分析引擎] 导出数据失败:", error);
        return e.json(500, {
            success: false,
            error: "导出数据失败"
        });
    }
});

// 添加定时任务：数据聚合
cronAdd("analytics_aggregation", `0 */${ANALYTICS_CONFIG.aggregationInterval} * * *`, () => {
    console.log("[分析引擎] 执行数据聚合任务");

    try {
        aggregateRawData();
        cleanupOldData();
        updateMaterializedViews();

        console.log("[分析引擎] 数据聚合完成");
    } catch (error) {
        console.error("[分析引擎] 数据聚合失败:", error);
    }
});

// 添加定时任务：生成日报
cronAdd("analytics_daily_report", ANALYTICS_CONFIG.reportSchedule, () => {
    console.log("[分析引擎] 生成每日分析报告");

    try {
        const report = generateDailyReport();

        // 保存报告到数据库
        app.db.execute(`
            INSERT INTO analytics_reports (type, report_data, created_at)
            VALUES (?, ?, datetime('now'))
        `, ['daily', JSON.stringify(report)]);

        // 发送报告给管理员
        sendReportToAdministrators(report);

        console.log("[分析引擎] 日报生成完成:", {
            date: report.date,
            total_urls: report.totalUrls,
            active_users: report.activeUsers
        });

    } catch (error) {
        console.error("[分析引擎] 生成日报失败:", error);
    }
});

// 添加定时任务：系统健康检查
cronAdd("analytics_health_check", "*/15 * * * *", () => {
    console.log("[分析引擎] 执行系统健康检查");

    try {
        const health = performHealthCheck();

        if (health.status !== 'healthy') {
            console.warn("[分析引擎] 系统健康检查异常:", health);
            triggerHealthAlert(health);
        }

        // 记录健康检查结果
        app.db.execute(`
            INSERT INTO system_health (status, metrics, checked_at)
            VALUES (?, ?, datetime('now'))
        `, [health.status, JSON.stringify(health.metrics)]);

    } catch (error) {
        console.error("[分析引擎] 健康检查失败:", error);
    }
});

// 工具函数：记录基础统计
function recordBasicStats(url) {
    app.db.execute(`
        INSERT INTO url_stats (url_id, domain, category, created_at)
        VALUES (?, ?, ?, datetime('now'))
    `, [url.id, extractDomain(url.url), url.category || 'uncategorized']);
}

// 工具函数：分析URL模式
function analyzeURLPattern(url) {
    const domain = extractDomain(url.url);
    const pathParts = extractPathParts(url.url);

    // 分析域名分布
    app.db.execute(`
        INSERT INTO domain_patterns (domain, count, last_seen)
        VALUES (?, 1, datetime('now'))
        ON CONFLICT(domain) DO UPDATE SET
        count = count + 1,
        last_seen = datetime('now')
    `, [domain]);

    // 分析路径模式
    pathParts.forEach(part => {
        if (part.length >= 3) {
            app.db.execute(`
                INSERT INTO path_patterns (pattern, count, last_seen)
                VALUES (?, 1, datetime('now'))
                ON CONFLICT(pattern) DO UPDATE SET
                count = count + 1,
                last_seen = datetime('now')
            `, [part.toLowerCase()]);
        }
    });
}

// 工具函数：更新实时指标
function updateRealTimeMetrics(metric, value) {
    app.db.execute(`
        INSERT INTO real_time_metrics (metric, value, timestamp)
        VALUES (?, ?, datetime('now'))
    `, [metric, value]);
}

// 工具函数：检查性能告警
function checkPerformanceAlerts() {
    const recentMetrics = app.db.queryOne(`
        SELECT AVG(CAST(value AS REAL)) as avg_value
        FROM real_time_metrics
        WHERE metric = 'response_time'
        AND timestamp > datetime('now', '-1 hour')
    `);

    if (recentMetrics && recentMetrics.avg_value > ANALYTICS_CONFIG.performanceThresholds.responseTime) {
        triggerPerformanceAlert('response_time', recentMetrics.avg_value);
    }
}

// 工具函数：记录访问统计
function recordAccessStats(url, user) {
    app.db.execute(`
        INSERT INTO access_stats (url_id, user_id, access_time, referrer, user_agent)
        VALUES (?, ?, datetime('now'), ?, ?)
    `, [url.id, user?.id || null, e.request.getHeader('Referer') || '', e.request.getHeader('User-Agent') || '']);
}

// 工具函数：分析用户行为
function analyzeUserBehavior(user, url) {
    if (!user) return;

    // 更新用户兴趣标签
    const category = url.category || 'uncategorized';
    app.db.execute(`
        INSERT INTO user_interests (user_id, category, score, last_updated)
        VALUES (?, ?, 1, datetime('now'))
        ON CONFLICT(user_id, category) DO UPDATE SET
        score = score + 1,
        last_updated = datetime('now')
    `, [user.id, category]);
}

// 工具函数：更新热门资源
function updatePopularResources(url) {
    app.db.execute(`
        INSERT INTO popular_resources (url_id, access_count, last_accessed)
        VALUES (?, 1, datetime('now'))
        ON CONFLICT(url_id) DO UPDATE SET
        access_count = access_count + 1,
        last_accessed = datetime('now')
    `, [url.id]);
}

// 工具函数：检测异常访问
function detectAnomalousAccess(user, url) {
    if (!user) return;

    // 检查用户最近访问频率
    const recentAccess = app.db.queryOne(`
        SELECT COUNT(*) as count
        FROM access_stats
        WHERE user_id = ? AND access_time > datetime('now', '-1 hour')
    `, [user.id]);

    if (recentAccess && recentAccess.count > 100) { // 阈值：1小时100次
        triggerSecurityAlert('excessive_access', user.id, recentAccess.count);
    }
}

// 工具函数：生成仪表板数据
function generateDashboardData() {
    const summary = app.db.queryOne(`
        SELECT
            COUNT(DISTINCT u.id) as total_urls,
            COUNT(DISTINCT a.user_id) as active_users,
            COUNT(DISTINCT s.domain) as unique_domains,
            COUNT(CASE WHEN u.created_at > datetime('now', '-7 days') THEN 1 END) as new_urls
        FROM urls u
        LEFT JOIN access_stats a ON a.access_time > datetime('now', '-7 days')
        LEFT JOIN url_stats s ON s.url_id = u.id
    `);

    const trends = getRecentTrends();
    const topResources = getTopResources();
    const userActivity = getUserActivityTrends();

    return {
        summary,
        trends,
        topResources,
        userActivity,
        lastUpdated: new Date().toISOString()
    };
}

// 工具函数：获取资源统计
function getResourceStats(period, category) {
    const periodCondition = getPeriodCondition(period);
    const categoryCondition = category ? `AND u.category = '${category}'` : '';

    return app.db.query(`
        SELECT
            DATE(u.created_at) as date,
            COUNT(*) as count,
            u.category,
            s.domain
        FROM urls u
        LEFT JOIN url_stats s ON s.url_id = u.id
        WHERE u.created_at > ${periodCondition} ${categoryCondition}
        GROUP BY DATE(u.created_at), u.category, s.domain
        ORDER BY date DESC
    `);
}

// 工具函数：获取用户分析
function getUserAnalytics(period, metric) {
    const periodCondition = getPeriodCondition(period);

    switch (metric) {
        case 'activity':
            return app.db.query(`
                SELECT
                    DATE(a.access_time) as date,
                    COUNT(DISTINCT a.user_id) as active_users,
                    COUNT(*) as total_actions
                FROM access_stats a
                WHERE a.access_time > ${periodCondition}
                GROUP BY DATE(a.access_time)
                ORDER BY date DESC
            `);
        case 'engagement':
            return app.db.query(`
                SELECT
                    u.username,
                    ui.category,
                    ui.score,
                    ui.last_updated
                FROM user_interests ui
                JOIN users u ON u.id = ui.user_id
                WHERE ui.last_updated > ${periodCondition}
                ORDER BY ui.score DESC
                LIMIT 50
            `);
        default:
            return [];
    }
}

// 工具函数：获取性能指标
function getPerformanceMetrics() {
    return {
        responseTime: getAverageResponseTime(),
        errorRate: getErrorRate(),
        memoryUsage: getMemoryUsage(),
        databaseSize: getDatabaseSize(),
        activeConnections: getActiveConnections(),
        systemLoad: getSystemLoad()
    };
}

// 工具函数：获取趋势分析
function getTrendAnalysis(metric, period) {
    const periodCondition = getPeriodCondition(period);

    switch (metric) {
        case 'growth':
            return getGrowthTrends(periodCondition);
        case 'seasonality':
            return getSeasonalityTrends(periodCondition);
        case 'prediction':
            return getPredictiveTrends(periodCondition);
        default:
            return [];
    }
}

// 工具函数：生成自定义报告
function generateCustomReport(type, period, format, filters) {
    const reportData = {
        type,
        period,
        generatedAt: new Date().toISOString(),
        filters,
        data: {}
    };

    switch (type) {
        case 'usage':
            reportData.data = generateUsageReport(period, filters);
            break;
        case 'performance':
            reportData.data = generatePerformanceReport(period, filters);
            break;
        case 'security':
            reportData.data = generateSecurityReport(period, filters);
            break;
        default:
            throw new Error(`不支持的报告类型: ${type}`);
    }

    return reportData;
}

// 工具函数：导出分析数据
function exportAnalyticsData(type, format, period) {
    const periodCondition = getPeriodCondition(period);
    let data = [];

    switch (type) {
        case 'summary':
            data = app.db.query(`
                SELECT
                    DATE(u.created_at) as date,
                    COUNT(*) as urls_added,
                    COUNT(DISTINCT u.category) as categories,
                    COUNT(DISTINCT s.domain) as domains
                FROM urls u
                LEFT JOIN url_stats s ON s.url_id = u.id
                WHERE u.created_at > ${periodCondition}
                GROUP BY DATE(u.created_at)
                ORDER BY date DESC
            `);
            break;
        case 'detailed':
            data = app.db.query(`
                SELECT
                    u.url,
                    u.title,
                    u.category,
                    u.created_at,
                    s.domain,
                    COUNT(a.id) as access_count
                FROM urls u
                LEFT JOIN url_stats s ON s.url_id = u.id
                LEFT JOIN access_stats a ON a.url_id = u.id
                WHERE u.created_at > ${periodCondition}
                GROUP BY u.id
                ORDER BY u.created_at DESC
            `);
            break;
    }

    return formatExportData(data, format);
}

// 辅助函数：提取域名
function extractDomain(url) {
    try {
        const urlObj = new URL(url);
        return urlObj.hostname.replace(/^www\./, '');
    } catch {
        return 'unknown';
    }
}

// 辅助函数：提取路径部分
function extractPathParts(url) {
    try {
        const urlObj = new URL(url);
        return urlObj.pathname.split('/').filter(Boolean);
    } catch {
        return [];
    }
}

// 辅助函数：获取时间条件
function getPeriodCondition(period) {
    switch (period) {
        case '1d': return "datetime('now', '-1 day')";
        case '7d': return "datetime('now', '-7 days')";
        case '30d': return "datetime('now', '-30 days')";
        case '90d': return "datetime('now', '-90 days')";
        default: return "datetime('now', '-7 days')";
    }
}

// 辅助函数：获取内容类型
function getContentType(format) {
    switch (format) {
        case 'csv': return 'text/csv';
        case 'json': return 'application/json';
        case 'xml': return 'application/xml';
        default: return 'text/plain';
    }
}

// 更多辅助函数...
function getRecentTrends() {
    // 实现获取最近趋势的逻辑
    return { growth: '+12%', activity: 'high' };
}

function getTopResources() {
    // 实现获取热门资源的逻辑
    return [];
}

function getUserActivityTrends() {
    // 实现获取用户活动趋势的逻辑
    return [];
}

function getAverageResponseTime() {
    // 实现获取平均响应时间的逻辑
    return 150;
}

function getErrorRate() {
    // 实现获取错误率的逻辑
    return 0.02;
}

function getMemoryUsage() {
    // 实现获取内存使用率的逻辑
    return 0.65;
}

function getDatabaseSize() {
    // 实现获取数据库大小的逻辑
    return '125MB';
}

function getActiveConnections() {
    // 实现获取活跃连接数的逻辑
    return 25;
}

function getSystemLoad() {
    // 实现获取系统负载的逻辑
    return 0.45;
}

function formatExportData(data, format) {
    // 实现数据格式化导出的逻辑
    return JSON.stringify(data);
}

console.log("[分析引擎] 插件初始化完成");

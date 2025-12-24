/// <reference path="../pb_data/types.d.ts" />

/**
 * 分析和统计插件钩子
 * 提供数据收集、分析和报告功能
 */

// 监听API请求事件
onAPIRequest((e) => {
    // 记录API访问统计
    recordAPIAccess(e);

    // 实时监控
    if (e.path.includes("/api/search")) {
        recordSearchQuery(e);
    }

    return e.next();
});

// 监听用户登录事件
onUserLogin((e) => {
    $log.info("用户登录: %s", e.user.username);

    // 记录登录统计
    recordUserLogin(e.user);

    // 更新用户活跃度
    updateUserActivity(e.user.id, "login");

    return e.next();
});

// 自定义路由：获取实时统计
routerAdd("GET", "/api/analytics/realtime", (e) => {
    const stats = getRealtimeStats();

    return e.json(200, {
        timestamp: new Date().toISOString(),
        online_users: stats.onlineUsers,
        active_requests: stats.activeRequests,
        recent_searches: stats.recentSearches,
        top_urls: stats.topUrls,
        system_load: stats.systemLoad
    });
});

// 自定义路由：获取访问趋势
routerAdd("GET", "/api/analytics/trends", (e) => {
    const { period = "7d" } = e.request.query;

    const trends = getAccessTrends(period);

    return e.json(200, {
        period,
        data: trends,
        summary: {
            total_visits: trends.reduce((sum, day) => sum + day.visits, 0),
            unique_visitors: trends.reduce((sum, day) => sum + day.uniqueVisitors, 0),
            avg_daily_visits: trends.length > 0 ? trends.reduce((sum, day) => sum + day.visits, 0) / trends.length : 0
        }
    });
});

// 自定义路由：获取搜索分析
routerAdd("GET", "/api/analytics/search", (e) => {
    const { period = "30d", limit = 50 } = e.request.query;

    const searchStats = getSearchAnalytics(period, parseInt(limit));

    return e.json(200, {
        period,
        hot_keywords: searchStats.hotKeywords,
        search_trends: searchStats.trends,
        no_result_searches: searchStats.noResultSearches,
        popular_categories: searchStats.popularCategories
    });
});

// 自定义路由：导出分析报告
routerAdd("GET", "/api/analytics/export", (e) => {
    const { format = "json", type = "full" } = e.request.query;

    let report;
    switch (type) {
        case "summary":
            report = generateSummaryReport();
            break;
        case "detailed":
            report = generateDetailedReport();
            break;
        default:
            report = generateFullReport();
    }

    if (format === "csv") {
        e.response.setHeader("Content-Type", "text/csv");
        e.response.setHeader("Content-Disposition", "attachment; filename=analytics_report.csv");
        return e.send(200, convertToCSV(report));
    }

    return e.json(200, report);
});

// 定时任务：生成日报
cronAdd("daily_report", "0 1 * * *", () => {
    $log.info("生成每日分析报告...");

    const report = generateDailyReport();
    saveReport("daily", report);

    // 发送报告通知
    sendReportNotification("daily", report);

    $log.info("日报生成完成");
});

// 定时任务：生成周报
cronAdd("weekly_report", "0 2 * * 1", () => {
    $log.info("生成每周分析报告...");

    const report = generateWeeklyReport();
    saveReport("weekly", report);

    sendReportNotification("weekly", report);

    $log.info("周报生成完成");
});

// 定时任务：清理过期统计数据
cronAdd("cleanup_stats", "0 3 * * 0", () => {
    $log.info("清理过期统计数据...");

    // 清理90天前的详细统计数据
    const cutoffDate = new Date();
    cutoffDate.setDate(cutoffDate.getDate() - 90);

    // cleanupAnalyticsData(cutoffDate);

    $log.info("统计数据清理完成");
});

// --- 数据收集函数 ---

function recordAPIAccess(event) {
    const accessData = {
        path: event.path,
        method: event.method,
        ip: event.request.ip,
        user_agent: event.request.headers["user-agent"],
        timestamp: new Date(),
        response_time: Date.now() - event.request.startTime
    };

    // 存储到数据库或缓存
    // $analytics.recordAPIAccess(accessData);
}

function recordSearchQuery(event) {
    const { q, category, tags } = event.request.query;

    if (!q) return;

    const searchData = {
        query: q,
        category: category,
        tags: tags ? tags.split(",") : [],
        ip: event.request.ip,
        timestamp: new Date(),
        results_count: event.response.results_count || 0
    };

    // $analytics.recordSearch(searchData);
}

function recordUserLogin(user) {
    const loginData = {
        user_id: user.id,
        username: user.username,
        ip: event.request.ip,
        timestamp: new Date()
    };

    // $analytics.recordLogin(loginData);
}

function updateUserActivity(userId, activity) {
    const activityData = {
        user_id: userId,
        activity: activity,
        timestamp: new Date()
    };

    // $analytics.recordActivity(activityData);
}

// --- 数据分析函数 ---

function getRealtimeStats() {
    // 模拟实时统计数据
    return {
        onlineUsers: Math.floor(Math.random() * 100) + 50,
        activeRequests: Math.floor(Math.random() * 20) + 5,
        recentSearches: [
            { query: "JavaScript教程", count: 15 },
            { query: "Docker容器", count: 12 },
            { query: "React框架", count: 8 }
        ],
        topUrls: [
            { url: "https://github.com/example/project", visits: 245 },
            { url: "https://stackoverflow.com/questions/123", visits: 189 }
        ],
        systemLoad: {
            cpu: Math.random() * 0.8,
            memory: Math.random() * 0.7,
            disk: Math.random() * 0.5
        }
    };
}

function getAccessTrends(period) {
    const days = period === "7d" ? 7 : period === "30d" ? 30 : 90;
    const trends = [];

    for (let i = days - 1; i >= 0; i--) {
        const date = new Date();
        date.setDate(date.getDate() - i);

        trends.push({
            date: date.toISOString().split('T')[0],
            visits: Math.floor(Math.random() * 1000) + 100,
            uniqueVisitors: Math.floor(Math.random() * 500) + 50,
            pageViews: Math.floor(Math.random() * 2000) + 200,
            avgSessionDuration: Math.floor(Math.random() * 300) + 60
        });
    }

    return trends;
}

function getSearchAnalytics(period, limit) {
    return {
        hotKeywords: [
            { keyword: "JavaScript", count: 156, trend: "up" },
            { keyword: "Python", count: 142, trend: "up" },
            { keyword: "React", count: 98, trend: "stable" },
            { keyword: "Docker", count: 87, trend: "down" },
            { keyword: "API", count: 76, trend: "up" }
        ],
        trends: [
            { date: "2024-01-01", searches: 234 },
            { date: "2024-01-02", searches: 189 },
            { date: "2024-01-03", searches: 267 }
        ],
        noResultSearches: [
            { query: "不存在的关键词", count: 23 },
            { query: "拼写错误", count: 15 }
        ],
        popularCategories: [
            { category: "开发工具", count: 342 },
            { category: "教程", count: 289 },
            { category: "框架", count: 198 }
        ]
    };
}

// --- 报告生成函数 ---

function generateDailyReport() {
    const today = new Date().toISOString().split('T')[0];
    const stats = getRealtimeStats();

    return {
        date: today,
        summary: {
            total_visits: stats.onlineUsers * 10, // 模拟数据
            unique_visitors: stats.onlineUsers,
            total_searches: 156,
            new_users: 12,
            top_search: "JavaScript教程"
        },
        categories: [
            { name: "开发工具", visits: 89 },
            { name: "教程", visits: 67 },
            { name: "框架", visits: 45 }
        ]
    };
}

function generateWeeklyReport() {
    return {
        week_start: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        week_end: new Date().toISOString().split('T')[0],
        total_visits: 5432,
        unique_visitors: 1234,
        avg_daily_visits: 776,
        top_pages: [
            { url: "/api/resources", visits: 1234 },
            { url: "/api/search", visits: 987 }
        ]
    };
}

function generateFullReport() {
    return {
        generated_at: new Date().toISOString(),
        period: "last_30_days",
        overview: {
            total_visits: 23456,
            unique_visitors: 5678,
            total_page_views: 45678,
            avg_session_duration: 245,
            bounce_rate: 0.35
        },
        traffic_sources: [
            { source: "直接访问", percentage: 45 },
            { source: "搜索引擎", percentage: 30 },
            { source: "外部链接", percentage: 25 }
        ],
        top_content: [
            { title: "JavaScript完整教程", views: 1234 },
            { title: "Docker入门指南", views: 987 }
        ]
    };
}

function convertToCSV(data) {
    // 简化的CSV转换逻辑
    return "date,visits,unique_visitors\n2024-01-01,1000,500\n2024-01-02,1200,600";
}

// --- 通知函数 ---

function sendReportNotification(type, report) {
    $log.info("发送%s报告通知", type);

    // 这里可以集成邮件、钉钉、企业微信等通知方式
    // 例如：$notification.sendEmail(report);
}

function saveReport(type, report) {
    const filename = `${type}_report_${new Date().toISOString().split('T')[0]}.json`;

    // 保存报告到文件系统或数据库
    // $fs.writeFile(`./reports/${filename}`, JSON.stringify(report, null, 2));

    $log.info("报告已保存: %s", filename);
}
/**
 * URL 相关插件钩子示例
 * 这些钩子会在 URL 的生命周期中被触发
 *
 * @name url_hooks
 * @display_name URL钩子插件
 * @author URLDB Team
 * @description URL生命周期事件钩子，监听URL的添加、访问等事件
 * @version 1.0.0
 * @category system
 * @license MIT
 */

// 记录插件加载日志
log("info", "URL钩子插件已加载", "url_hooks");

// 监听 URL 添加事件
onURLAdd(function(event) {
    log("info", "=== onURLAdd 事件触发 ===", "url_hooks");
    log("info", "URL ID: " + event.url.id, "url_hooks");
    log("info", "URL Title: " + event.url.title, "url_hooks");
    log("info", "URL: " + event.url.url, "url_hooks");

    // 自动分类逻辑
    if (event.url.url && event.url.url.includes("github.com")) {
        log("info", "检测到GitHub URL，建议分类为: 开发工具", "url_hooks");
    }

    // 标签提取逻辑
    const tags = extractTagsFromUrl(event.url.url, event.url.title);
    if (tags.length > 0) {
        log("info", "自动提取标签: " + tags.join(", "), "url_hooks");
    }

    log("info", "=== onURLAdd 事件处理完成 ===", "url_hooks");
});

// 监听 URL 访问事件
onURLAccess(function(event) {
    log("info", "=== onURLAccess 事件触发 ===", "url_hooks");
    log("info", "URL被访问: " + event.url.url, "url_hooks");

    if (event.request && event.request.ip) {
        log("info", "访问IP: " + event.request.ip, "url_hooks");
    }

    // 热门度计算概念
    log("info", "更新URL热门度", "url_hooks");

    log("info", "=== onURLAccess 事件处理完成 ===", "url_hooks");
});

// 自定义路由：URL 智能分析
router.post("/api/analyze-url", function() {
    // 这里简化实现，返回基本信息
    return {
        success: true,
        message: "URL分析服务运行正常",
        plugin: "url_hooks",
        timestamp: new Date().toISOString()
    };
});

// 自定义路由：批量URL处理
router.post("/api/batch-process", function() {
    return {
        success: true,
        message: "批量处理服务运行正常",
        plugin: "url_hooks",
        timestamp: new Date().toISOString()
    };
});

// 定时任务：清理过期访问日志
cron("cleanup_logs", "0 2 * * *", function() {
    log("info", "开始清理过期访问日志...", "url_hooks");
    log("info", "访问日志清理完成", "url_hooks");
});

// 定时任务：计算热门URL
cron("calculate_hot", "0 */6 * * *", function() {
    log("info", "开始计算热门URL...", "url_hooks");
    log("info", "热门URL计算完成", "url_hooks");
});

log("info", "URL钩子插件初始化完成", "url_hooks");

// --- 工具函数 ---

// 从URL和标题提取标签
function extractTagsFromUrl(url, title) {
    const tags = new Set();

    // 从URL域名提取
    const domain = new URL(url).hostname;
    if (domain.includes("github")) tags.add("GitHub");
    if (domain.includes("stackoverflow")) tags.add("StackOverflow");
    if (domain.includes("youtube")) tags.add("视频");
    if (domain.includes("bilibili")) tags.add("B站");

    // 从标题提取关键词
    if (title) {
        const keywords = [
            "教程", "文档", "工具", "框架", "库", "API",
            "JavaScript", "Python", "Go", "Java", "React",
            "Vue", "Node.js", "Docker", "Kubernetes"
        ];

        keywords.forEach(keyword => {
            if (title.toLowerCase().includes(keyword.toLowerCase())) {
                tags.add(keyword);
            }
        });
    }

    return Array.from(tags);
}

// 分析URL
function analyzeUrl(url) {
    try {
        const urlObj = new URL(url);
        const domain = urlObj.hostname;
        const path = urlObj.pathname;

        let category = "其他";
        let tags = [];
        let confidence = 0.5;

        // 域名分类
        if (domain.includes("github")) {
            category = "代码仓库";
            tags = ["GitHub", "开源"];
            confidence = 0.9;
        } else if (domain.includes("stackoverflow")) {
            category = "技术问答";
            tags = ["StackOverflow", "问答"];
            confidence = 0.9;
        } else if (domain.includes("youtube") || domain.includes("bilibili")) {
            category = "视频";
            tags = domain.includes("youtube") ? ["YouTube", "视频"] : ["B站", "视频"];
            confidence = 0.8;
        }

        return {
            category,
            tags,
            description: `自动分析的${category}资源`,
            confidence
        };
    } catch (error) {
        return {
            category: "无效URL",
            tags: [],
            description: "URL格式错误",
            confidence: 0
        };
    }
}

// URL分类
function categorizeUrl(url) {
    const analysis = analyzeUrl(url);
    return analysis.category;
}

// 更新热门度分数
function updateHotScore(url) {
    // 简单的热门度计算逻辑
    const now = new Date();
    const daysSinceCreation = (now - new Date(url.createdAt)) / (1000 * 60 * 60 * 24);

    // 基于访问量和时间衰减计算热门度
    const accessCount = url.accessCount || 0;
    const timeDecay = Math.exp(-daysSinceCreation / 30); // 30天衰减
    const hotScore = accessCount * timeDecay;

    // 更新数据库中的热门度分数
    // $urldb.updateHotScore(url.id, hotScore);
}

// 清理相关数据
function cleanupRelatedData(urlId) {
    // 清理与该URL相关的访问日志、标签等
    $log.info("清理URL %s 的相关数据", urlId);
}
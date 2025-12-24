/// <reference path="../pb_data/types.d.ts" />

/**
 * URL 相关插件钩子示例
 * 这些钩子会在 URL 的生命周期中被触发
 */

// 监听 URL 添加事件
onURLAdd((e) => {
    $log.info("新URL添加: %s", e.url.url);

    // 自动分类逻辑
    if (e.url.url.includes("github.com")) {
        e.url.category = "开发工具";
        $log.info("自动分类为: 开发工具");
    }

    // 标签提取逻辑
    const tags = extractTagsFromUrl(e.url.url, e.url.title);
    if (tags.length > 0) {
        e.url.tags = tags;
        $log.info("自动提取标签: %s", tags.join(", "));
    }

    return e.next();
});

// 监听 URL 访问事件
onURLAccess((e) => {
    $log.info("URL被访问: %s - IP: %s", e.url.url, e.request.ip);

    // 访问统计
    $urldb.incrementAccess(e.url.id);

    // 热门度计算
    updateHotScore(e.url);

    return e.next();
});

// 监听 URL 更新事件
onURLUpdate((e) => {
    $log.info("URL更新: %s", e.url.url);

    // 重新计算标签
    const tags = extractTagsFromUrl(e.url.url, e.url.title);
    e.url.tags = tags;

    return e.next();
});

// 监听 URL 删除事件
onURLDelete((e) => {
    $log.info("URL删除: %s", e.url.url);

    // 清理相关数据
    cleanupRelatedData(e.url.id);

    return e.next();
});

// 自定义路由：URL 智能分析
routerAdd("POST", "/api/analyze-url", (e) => {
    const { url } = e.request.body;

    if (!url) {
        return e.json(400, { error: "URL参数缺失" });
    }

    const analysis = analyzeUrl(url);

    return e.json(200, {
        url: url,
        category: analysis.category,
        tags: analysis.tags,
        description: analysis.description,
        confidence: analysis.confidence
    });
});

// 自定义路由：批量URL处理
routerAdd("POST", "/api/batch-process", (e) => {
    const { urls, operation } = e.request.body;

    if (!urls || !Array.isArray(urls)) {
        return e.json(400, { error: "URLs参数必须是数组" });
    }

    const results = [];

    for (const url of urls) {
        try {
            let result;
            switch (operation) {
                case "analyze":
                    result = analyzeUrl(url);
                    break;
                case "categorize":
                    result = { category: categorizeUrl(url) };
                    break;
                case "extract_tags":
                    result = { tags: extractTagsFromUrl(url) };
                    break;
                default:
                    result = { error: "不支持的操作" };
            }
            results.push({ url, ...result });
        } catch (error) {
            results.push({ url, error: error.message });
        }
    }

    return e.json(200, { results });
});

// 定时任务：清理过期访问日志
cronAdd("cleanup_logs", "0 2 * * *", () => {
    $log.info("开始清理过期访问日志...");

    // 清理30天前的访问日志
    const cutoffDate = new Date();
    cutoffDate.setDate(cutoffDate.getDate() - 30);

    // 这里应该调用数据库清理方法
    // $urldb.cleanupAccessLogs(cutoffDate);

    $log.info("访问日志清理完成");
});

// 定时任务：计算热门URL
cronAdd("calculate_hot", "0 */6 * * *", () => {
    $log.info("开始计算热门URL...");

    // 计算最近24小时的热门URL
    // const hotUrls = $urldb.calculateHotUrls(24);

    $log.info("热门URL计算完成");
});

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
/// <reference path="../pb_data/types.d.ts" />

/**
 * URL智能分类和标签提取插件
 * 功能：自动为URL分类并提取相关标签
 * 作者: URLDB Team
 * 版本: 1.0.0
 * 创建时间: 2024-12-24
 */

// 分类规则配置
const CLASSIFICATION_RULES = {
    // 开发相关
    development: {
        domains: ['github.com', 'gitlab.com', 'bitbucket.org', 'stackoverflow.com', 'dev.to', 'medium.com'],
        paths: ['docs', 'api', 'tutorial', 'guide'],
        keywords: ['code', 'programming', 'development', 'tutorial', 'documentation']
    },

    // 视频相关
    video: {
        domains: ['youtube.com', 'bilibili.com', 'vimeo.com', 'dailymotion.com', 'twitch.tv'],
        paths: ['watch', 'video', 'play'],
        keywords: ['video', 'watch', 'stream', 'play']
    },

    // 社交媒体
    social: {
        domains: ['twitter.com', 'facebook.com', 'instagram.com', 'linkedin.com', 'weibo.com'],
        paths: ['profile', 'post', 'timeline'],
        keywords: ['social', 'network', 'profile', 'post']
    },

    // 新闻资讯
    news: {
        domains: ['cnn.com', 'bbc.com', 'reuters.com', 'xinhuanet.com', 'people.com.cn'],
        paths: ['news', 'article', 'report'],
        keywords: ['news', 'article', 'report', 'breaking']
    },

    // 购物电商
    shopping: {
        domains: ['amazon.com', 'taobao.com', 'jd.com', 'tmall.com', 'pinduoduo.com'],
        paths: ['product', 'shop', 'buy', 'item'],
        keywords: ['shop', 'buy', 'product', 'price', 'deal']
    },

    // 教育学习
    education: {
        domains: ['coursera.org', 'edx.org', 'udemy.com', 'khanacademy.org', 'study.com'],
        paths: ['course', 'learn', 'lesson', 'tutorial'],
        keywords: ['course', 'learn', 'education', 'tutorial', 'lesson']
    }
};

// 监听URL添加事件
onURLAdd((e) => {
    console.log("[URL分类器] 开始处理URL:", e.url.url);

    try {
        // 提取URL信息
        const urlInfo = extractURLInfo(e.url.url);

        // 智能分类
        const category = classifyURL(urlInfo);
        e.url.category = category;

        // 提取标签
        const tags = extractTags(urlInfo, category);
        e.url.tags = [...new Set([...(e.url.tags || []), ...tags])]; // 合并去重

        // 计算置信度
        const confidence = calculateConfidence(urlInfo, category, tags);
        e.url.confidence = confidence;

        // 添加元数据
        e.url.metadata = {
            domain: urlInfo.domain,
            path_parts: urlInfo.pathParts,
            query_params: urlInfo.queryParams,
            protocol: urlInfo.protocol,
            classified_at: new Date().toISOString(),
            classifier_version: '1.0.0'
        };

        // 记录分类结果
        console.log("[URL分类器] 分类完成:", {
            url: e.url.url,
            category: e.url.category,
            tags: e.url.tags,
            confidence: e.url.confidence
        });

        // 保存分类统计
        saveClassificationStats(category, tags.length);

    } catch (error) {
        console.error("[URL分类器] 处理失败:", error);
        // 设置默认分类
        e.url.category = 'uncategorized';
        e.url.tags = e.url.tags || ['auto-classified'];
    }

    return e.next();
});

// 监听用户登录事件，展示用户偏好
onUserLogin((e) => {
    console.log("[URL分类器] 用户登录:", e.user.username);

    // 获取用户分类偏好
    const userPreferences = getUserCategoryPreferences(e.user.id);

    if (userPreferences && Object.keys(userPreferences).length > 0) {
        console.log("[URL分类器] 用户分类偏好:", userPreferences);

        // 将用户偏好存储到事件数据中，供其他插件使用
        e.data.user_category_preferences = userPreferences;
    }

    return e.next();
});

// 添加自定义API路由
routerAdd("GET", "/api/plugins/classifier/stats", (e) => {
    try {
        // 获取分类统计
        const stats = getClassificationStats();

        return e.json(200, {
            success: true,
            data: {
                total_classified: stats.total,
                category_distribution: stats.categories,
                most_common_tags: stats.tags,
                avg_confidence: stats.avgConfidence,
                last_updated: stats.lastUpdated
            }
        });
    } catch (error) {
        console.error("[URL分类器] 获取统计失败:", error);
        return e.json(500, {
            success: false,
            error: "获取统计信息失败"
        });
    }
});

// 添加手动分类API
routerAdd("POST", "/api/plugins/classifier/classify", (e) => {
    try {
        const { url } = e.body;

        if (!url) {
            return e.json(400, {
                success: false,
                error: "URL参数不能为空"
            });
        }

        // 验证URL格式
        let urlObj;
        try {
            urlObj = new URL(url);
        } catch (error) {
            return e.json(400, {
                success: false,
                error: "URL格式无效"
            });
        }

        // 执行分类
        const urlInfo = extractURLInfo(url);
        const category = classifyURL(urlInfo);
        const tags = extractTags(urlInfo, category);
        const confidence = calculateConfidence(urlInfo, category, tags);

        return e.json(200, {
            success: true,
            data: {
                url: url,
                category: category,
                tags: tags,
                confidence: confidence,
                metadata: {
                    domain: urlInfo.domain,
                    path_parts: urlInfo.pathParts,
                    query_params: urlInfo.queryParams
                }
            }
        });

    } catch (error) {
        console.error("[URL分类器] 手动分类失败:", error);
        return e.json(500, {
            success: false,
            error: "分类处理失败"
        });
    }
});

// 添加定时任务，每日生成分类报告
cronAdd("classifier_daily_report", "0 8 * * *", () => {
    console.log("[URL分类器] 生成每日分类报告");

    try {
        const report = generateDailyReport();

        // 将报告保存到数据库
        app.db.execute(`
            INSERT INTO daily_reports (type, report_data, created_at)
            VALUES (?, ?, datetime('now'))
        `, ['classifier', JSON.stringify(report)]);

        console.log("[URL分类器] 每日报告生成完成:", {
            date: report.date,
            total_classified: report.totalClassified,
            top_category: report.topCategory
        });

    } catch (error) {
        console.error("[URL分类器] 生成报告失败:", error);
    }
});

// 工具函数：提取URL信息
function extractURLInfo(url) {
    try {
        const urlObj = new URL(url);

        return {
            originalUrl: url,
            protocol: urlObj.protocol.replace(':', ''),
            domain: urlObj.hostname.replace(/^www\./, ''),
            path: urlObj.pathname,
            pathParts: urlObj.pathname.split('/').filter(Boolean),
            queryParams: Object.fromEntries(urlObj.searchParams),
            hash: urlObj.hash,
            port: urlObj.port || (urlObj.protocol === 'https:' ? 443 : 80)
        };
    } catch (error) {
        console.error("[URL分类器] URL解析失败:", error);
        return {
            originalUrl: url,
            protocol: 'unknown',
            domain: 'unknown',
            path: '',
            pathParts: [],
            queryParams: {},
            hash: '',
            port: 0
        };
    }
}

// 工具函数：URL分类
function classifyURL(urlInfo) {
    let bestMatch = { category: 'uncategorized', score: 0 };

    // 遍历所有分类规则
    for (const [category, rules] of Object.entries(CLASSIFICATION_RULES)) {
        let score = 0;

        // 域名匹配（权重最高）
        if (rules.domains.some(domain => urlInfo.domain.includes(domain))) {
            score += 10;
        }

        // 路径匹配
        if (rules.paths.some(path => urlInfo.pathParts.some(part => part.includes(path)))) {
            score += 5;
        }

        // 查询参数匹配
        if (rules.keywords.some(keyword =>
            Object.values(urlInfo.queryParams).some(value =>
                String(value).toLowerCase().includes(keyword)
            )
        )) {
            score += 3;
        }

        // 路径部分关键词匹配
        if (rules.keywords.some(keyword =>
            urlInfo.pathParts.some(part => part.toLowerCase().includes(keyword))
        )) {
            score += 2;
        }

        // 更新最佳匹配
        if (score > bestMatch.score) {
            bestMatch = { category, score };
        }
    }

    return bestMatch.category;
}

// 工具函数：提取标签
function extractTags(urlInfo, category) {
    const tags = new Set();

    // 添加域名标签
    tags.add(urlInfo.domain);

    // 添加路径部分作为标签（过滤掉太短或纯数字的）
    urlInfo.pathParts.forEach(part => {
        if (part.length >= 3 && !/^\d+$/.test(part)) {
            tags.add(part.toLowerCase());
        }
    });

    // 添加查询参数作为标签
    Object.keys(urlInfo.queryParams).forEach(key => {
        if (key.length >= 3) {
            tags.add(key.toLowerCase());
        }
    });

    // 添加分类特定标签
    const categoryRules = CLASSIFICATION_RULES[category];
    if (categoryRules) {
        categoryRules.keywords.forEach(keyword => {
            if (urlInfo.path.includes(keyword) ||
                Object.values(urlInfo.queryParams).some(value =>
                    String(value).toLowerCase().includes(keyword)
                )) {
                tags.add(keyword);
            }
        });
    }

    // 添加协议标签
    tags.add(urlInfo.protocol);

    // 添加特殊标签
    if (urlInfo.queryParams.length > 0) tags.add('has-params');
    if (urlInfo.hash) tags.add('has-hash');
    if (urlInfo.port !== 80 && urlInfo.port !== 443) tags.add('custom-port');

    return Array.from(tags);
}

// 工具函数：计算置信度
function calculateConfidence(urlInfo, category, tags) {
    let confidence = 0.3; // 基础置信度

    // 分类匹配度
    const categoryRules = CLASSIFICATION_RULES[category];
    if (categoryRules) {
        if (categoryRules.domains.includes(urlInfo.domain)) {
            confidence += 0.4;
        }

        if (categoryRules.paths.some(path =>
            urlInfo.pathParts.some(part => part.includes(path))
        )) {
            confidence += 0.2;
        }
    }

    // 标签数量（但不要太高）
    const tagScore = Math.min(tags.length * 0.05, 0.1);
    confidence += tagScore;

    return Math.min(confidence, 1.0);
}

// 工具函数：保存分类统计
function saveClassificationStats(category, tagCount) {
    try {
        app.db.execute(`
            INSERT INTO classification_stats (category, tag_count, created_at)
            VALUES (?, ?, datetime('now'))
            ON CONFLICT(category) DO UPDATE SET
            tag_count = tag_count + excluded.tag_count,
            updated_at = datetime('now')
        `, [category, tagCount]);
    } catch (error) {
        console.error("[URL分类器] 保存统计失败:", error);
    }
}

// 工具函数：获取分类统计
function getClassificationStats() {
    try {
        const totalResult = app.db.queryOne(`
            SELECT COUNT(*) as total FROM urls WHERE category IS NOT NULL AND category != 'uncategorized'
        `);

        const categoryResults = app.db.query(`
            SELECT category, COUNT(*) as count
            FROM urls
            WHERE category IS NOT NULL AND category != 'uncategorized'
            GROUP BY category
            ORDER BY count DESC
        `);

        const tagResults = app.db.query(`
            SELECT value as tag, COUNT(*) as count
            FROM (
                SELECT json_each(value) as value FROM urls WHERE tags IS NOT NULL
            )
            GROUP BY value
            ORDER BY count DESC
            LIMIT 10
        `);

        const confidenceResult = app.db.queryOne(`
            SELECT AVG(confidence) as avg_confidence
            FROM urls
            WHERE confidence IS NOT NULL
        `);

        return {
            total: totalResult.total || 0,
            categories: categoryResults,
            tags: tagResults,
            avgConfidence: confidenceResult.avg_confidence || 0,
            lastUpdated: new Date().toISOString()
        };
    } catch (error) {
        console.error("[URL分类器] 获取统计失败:", error);
        return {
            total: 0,
            categories: [],
            tags: [],
            avgConfidence: 0,
            lastUpdated: new Date().toISOString()
        };
    }
}

// 工具函数：获取用户分类偏好
function getUserCategoryPreferences(userId) {
    try {
        const preferences = app.db.queryOne(`
            SELECT category_preferences FROM user_preferences WHERE user_id = ?
        `, [userId]);

        return preferences ? JSON.parse(preferences.category_preferences || '{}') : null;
    } catch (error) {
        console.error("[URL分类器] 获取用户偏好失败:", error);
        return null;
    }
}

// 工具函数：生成每日报告
function generateDailyReport() {
    const today = new Date().toISOString().split('T')[0];

    const todayStats = app.db.queryOne(`
        SELECT
            COUNT(*) as total_classified,
            category,
            COUNT(*) as category_count
        FROM urls
        WHERE DATE(created_at) = ? AND category IS NOT NULL AND category != 'uncategorized'
        GROUP BY category
        ORDER BY category_count DESC
        LIMIT 1
    `, [today]);

    const allCategories = app.db.query(`
        SELECT category, COUNT(*) as count
        FROM urls
        WHERE DATE(created_at) = ? AND category IS NOT NULL AND category != 'uncategorized'
        GROUP BY category
        ORDER BY count DESC
    `, [today]);

    return {
        date: today,
        totalClassified: todayStats.total_classified || 0,
        topCategory: todayStats.category || 'none',
        categoryDistribution: allCategories,
        generatedAt: new Date().toISOString()
    };
}

console.log("[URL分类器] 插件初始化完成");
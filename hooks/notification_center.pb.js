/// <reference path="../pb_data/types.d.ts" />

/**
 * 通知中心插件
 * 功能：统一管理系统通知，支持邮件、站内消息、推送等多种通知方式
 * 作者: URLDB Team
 * 版本: 1.0.0
 * 创建时间: 2024-12-24
 */

// 通知类型定义
const NOTIFICATION_TYPES = {
    SYSTEM: 'system',           // 系统通知
    USER: 'user',             // 用户相关通知
    URL: 'url',               // URL相关通知
    SECURITY: 'security',     // 安全通知
    MARKETING: 'marketing'     // 营销通知
};

// 通知优先级
const PRIORITY_LEVELS = {
    LOW: 1,
    NORMAL: 2,
    HIGH: 3,
    URGENT: 4
};

// 通知渠道
const NOTIFICATION_CHANNELS = {
    EMAIL: 'email',
    IN_APP: 'in_app',
    PUSH: 'push',
    WEBHOOK: 'webhook'
};

// 监听URL添加事件，发送通知
onURLAdd((e) => {
    console.log("[通知中心] URL添加事件触发:", e.url.url);

    try {
        // 检查用户通知偏好
        const userPrefs = getUserNotificationPreferences(e.user?.id);

        if (!userPrefs?.enable_url_notifications) {
            return e.next();
        }

        // 创建URL添加通知
        const notification = createNotification({
            type: NOTIFICATION_TYPES.URL,
            priority: PRIORITY_LEVELS.NORMAL,
            title: '新URL已添加',
            message: `成功添加URL: ${e.url.title || e.url.url}`,
            data: {
                url_id: e.url.id,
                url: e.url.url,
                category: e.url.category,
                user_id: e.user?.id
            },
            channels: userPrefs.url_channels || [NOTIFICATION_CHANNELS.IN_APP]
        });

        // 发送通知
        sendNotification(notification);

        console.log("[通知中心] URL添加通知已发送");

    } catch (error) {
        console.error("[通知中心] URL通知发送失败:", error);
    }

    return e.next();
});

// 监听用户登录事件，发送欢迎通知
onUserLogin((e) => {
    console.log("[通知中心] 用户登录事件触发:", e.user.username);

    try {
        // 检查是否是首次登录
        const isFirstLogin = checkFirstLogin(e.user.id);

        if (isFirstLogin) {
            // 发送欢迎通知
            const welcomeNotification = createNotification({
                type: NOTIFICATION_TYPES.USER,
                priority: PRIORITY_LEVELS.NORMAL,
                title: '欢迎来到URLDB！',
                message: `欢迎您，${e.user.username}！感谢您注册使用我们的服务。`,
                data: {
                    user_id: e.user.id,
                    username: e.user.username,
                    email: e.user.email
                },
                channels: [NOTIFICATION_CHANNELS.IN_APP, NOTIFICATION_CHANNELS.EMAIL]
            });

            sendNotification(welcomeNotification);

            // 记录首次登录
            recordFirstLogin(e.user.id);
        } else {
            // 发送登录通知（可选）
            const loginNotification = createNotification({
                type: NOTIFICATION_TYPES.USER,
                priority: PRIORITY_LEVELS.LOW,
                title: '登录成功',
                message: `欢迎回来，${e.user.username}！`,
                data: {
                    user_id: e.user.id,
                    login_time: new Date().toISOString()
                },
                channels: [NOTIFICATION_CHANNELS.IN_APP]
            });

            sendNotification(loginNotification);
        }

    } catch (error) {
        console.error("[通知中心] 登录通知发送失败:", error);
    }

    return e.next();
});

// 监听API请求事件，记录可疑活动
onAPIRequest((e) => {
    // 只监控特定路径的API请求
    const monitoredPaths = ['/api/admin', '/api/user/delete', '/api/batch'];
    const isMonitored = monitoredPaths.some(path => e.path.startsWith(path));

    if (!isMonitored) {
        return e.next();
    }

    console.log("[通知中心] API请求监控:", e.method, e.path);

    try {
        // 检查是否有可疑活动
        const suspicious = detectSuspiciousActivity(e);

        if (suspicious.isSuspicious) {
            // 创建安全通知
            const securityNotification = createNotification({
                type: NOTIFICATION_TYPES.SECURITY,
                priority: PRIORITY_LEVELS.HIGH,
                title: '检测到可疑活动',
                message: `检测到来自 ${e.ip} 的可疑${e.method}请求: ${e.path}`,
                data: {
                    ip: e.ip,
                    method: e.method,
                    path: e.path,
                    user_agent: e.headers['user-agent'],
                    timestamp: new Date().toISOString(),
                    risk_score: suspicious.riskScore
                },
                channels: [NOTIFICATION_CHANNELS.IN_APP, NOTIFICATION_CHANNELS.EMAIL]
            });

            sendNotification(securityNotification);

            // 记录安全事件
            recordSecurityEvent(suspicious);
        }

    } catch (error) {
        console.error("[通知中心] 安全监控失败:", error);
    }

    return e.next();
});

// 添加通知管理API路由
routerAdd("GET", "/api/plugins/notifications", (e) => {
    try {
        const userId = e.user?.id;
        const { page = 1, limit = 20, unread_only = false, type } = e.query;

        if (!userId) {
            return e.json(401, {
                success: false,
                error: "需要登录"
            });
        }

        const notifications = getUserNotifications(userId, {
            page: parseInt(page),
            limit: parseInt(limit),
            unreadOnly: unread_only === 'true',
            type
        });

        return e.json(200, {
            success: true,
            data: notifications
        });

    } catch (error) {
        console.error("[通知中心] 获取通知失败:", error);
        return e.json(500, {
            success: false,
            error: "获取通知失败"
        });
    }
});

// 标记通知为已读
routerAdd("PUT", "/api/plugins/notifications/:id/read", (e) => {
    try {
        const userId = e.user?.id;
        const notificationId = e.pathParams.id;

        if (!userId) {
            return e.json(401, {
                success: false,
                error: "需要登录"
            });
        }

        markNotificationAsRead(notificationId, userId);

        return e.json(200, {
            success: true,
            message: "通知已标记为已读"
        });

    } catch (error) {
        console.error("[通知中心] 标记已读失败:", error);
        return e.json(500, {
            success: false,
            error: "标记已读失败"
        });
    }
});

// 批量标记通知为已读
routerAdd("PUT", "/api/plugins/notifications/read-all", (e) => {
    try {
        const userId = e.user?.id;

        if (!userId) {
            return e.json(401, {
                success: false,
                error: "需要登录"
            });
        }

        const result = markAllNotificationsAsRead(userId);

        return e.json(200, {
            success: true,
            data: {
                marked_count: result.markedCount
            }
        });

    } catch (error) {
        console.error("[通知中心] 批量标记已读失败:", error);
        return e.json(500, {
            success: false,
            error: "批量标记已读失败"
        });
    }
});

// 获取通知统计
routerAdd("GET", "/api/plugins/notifications/stats", (e) => {
    try {
        const userId = e.user?.id;

        if (!userId) {
            return e.json(401, {
                success: false,
                error: "需要登录"
            });
        }

        const stats = getNotificationStats(userId);

        return e.json(200, {
            success: true,
            data: stats
        });

    } catch (error) {
        console.error("[通知中心] 获取统计失败:", error);
        return e.json(500, {
            success: false,
            error: "获取统计失败"
        });
    }
});

// 发送自定义通知
routerAdd("POST", "/api/plugins/notifications/send", (e) => {
    try {
        const userId = e.user?.id;
        const { title, message, type = NOTIFICATION_TYPES.USER, priority = PRIORITY_LEVELS.NORMAL } = e.body;

        if (!userId) {
            return e.json(401, {
                success: false,
                error: "需要登录"
            });
        }

        if (!title || !message) {
            return e.json(400, {
                success: false,
                error: "标题和消息不能为空"
            });
        }

        const notification = createNotification({
            type,
            priority,
            title,
            message,
            data: {
                user_id: userId,
                source: 'user_manual'
            },
            channels: [NOTIFICATION_CHANNELS.IN_APP]
        });

        sendNotification(notification);

        return e.json(200, {
            success: true,
            message: "通知发送成功",
            data: {
                notification_id: notification.id
            }
        });

    } catch (error) {
        console.error("[通知中心] 发送通知失败:", error);
        return e.json(500, {
            success: false,
            error: "发送通知失败"
        });
    }
});

// 添加定时任务，清理过期通知
cronAdd("cleanup_notifications", "0 2 * * *", () => {
    console.log("[通知中心] 开始清理过期通知");

    try {
        const result = cleanupExpiredNotifications();

        console.log("[通知中心] 清理完成:", {
            deleted_count: result.deletedCount
        });

    } catch (error) {
        console.error("[通知中心] 清理通知失败:", error);
    }
});

// 添加定时任务，发送每日摘要
cronAdd("daily_summary", "0 9 * * *", () => {
    console.log("[通知中心] 开始生成每日摘要");

    try {
        const summary = generateDailySummary();

        // 为活跃用户发送每日摘要
        sendDailySummaryToUsers(summary);

        console.log("[通知中心] 每日摘要发送完成");

    } catch (error) {
        console.error("[通知中心] 每日摘要失败:", error);
    }
});

// 工具函数：创建通知对象
function createNotification(options) {
    return {
        id: generateUUID(),
        type: options.type,
        priority: options.priority,
        title: options.title,
        message: options.message,
        data: options.data || {},
        channels: options.channels || [NOTIFICATION_CHANNELS.IN_APP],
        created_at: new Date().toISOString(),
        read: false,
        read_at: null
    };
}

// 工具函数：发送通知
function sendNotification(notification) {
    try {
        // 保存通知到数据库
        saveNotificationToDB(notification);

        // 根据渠道发送通知
        notification.channels.forEach(channel => {
            switch (channel) {
                case NOTIFICATION_CHANNELS.EMAIL:
                    sendEmailNotification(notification);
                    break;
                case NOTIFICATION_CHANNELS.IN_APP:
                    // 站内通知已通过数据库保存
                    break;
                case NOTIFICATION_CHANNELS.PUSH:
                    sendPushNotification(notification);
                    break;
                case NOTIFICATION_CHANNELS.WEBHOOK:
                    sendWebhookNotification(notification);
                    break;
            }
        });

        console.log("[通知中心] 通知发送成功:", notification.id);

    } catch (error) {
        console.error("[通知中心] 发送通知失败:", error);
    }
}

// 工具函数：保存通知到数据库
function saveNotificationToDB(notification) {
    try {
        app.db.execute(`
            INSERT INTO notifications (
                id, type, priority, title, message, data, channels,
                user_id, created_at, read, read_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `, [
            notification.id,
            notification.type,
            notification.priority,
            notification.title,
            notification.message,
            JSON.stringify(notification.data),
            JSON.stringify(notification.channels),
            notification.data.user_id,
            notification.created_at,
            notification.read ? 1 : 0,
            notification.read_at
        ]);
    } catch (error) {
        console.error("[通知中心] 保存通知失败:", error);
    }
}

// 工具函数：发送邮件通知
function sendEmailNotification(notification) {
    try {
        const userEmail = getUserEmail(notification.data.user_id);

        if (!userEmail) {
            console.warn("[通知中心] 用户邮箱不存在，跳过邮件通知");
            return;
        }

        // 模拟邮件发送（实际项目中应使用真实的邮件服务）
        console.log("[通知中心] 发送邮件通知:", {
            to: userEmail,
            subject: notification.title,
            body: notification.message,
            timestamp: new Date().toISOString()
        });

        // 实际邮件发送代码示例：
        // app.mail.send({
        //     to: userEmail,
        //     subject: notification.title,
        //     body: notification.message,
        //     html: `<p>${notification.message}</p>`
        // });

    } catch (error) {
        console.error("[通知中心] 邮件发送失败:", error);
    }
}

// 工具函数：发送推送通知
function sendPushNotification(notification) {
    try {
        console.log("[通知中心] 发送推送通知:", {
            title: notification.title,
            message: notification.message,
            user_id: notification.data.user_id
        });

        // 实际推送服务集成代码
        // pushService.send(notification.data.user_id, {
        //     title: notification.title,
        //     body: notification.message,
        //     data: notification.data
        // });

    } catch (error) {
        console.error("[通知中心] 推送发送失败:", error);
    }
}

// 工具函数：发送Webhook通知
function sendWebhookNotification(notification) {
    try {
        console.log("[通知中心] 发送Webhook通知:", {
            url: "https://webhook.example.com/notifications",
            data: notification
        });

        // 实际Webhook发送代码
        // fetch("https://webhook.example.com/notifications", {
        //     method: "POST",
        //     headers: {
        //         "Content-Type": "application/json"
        //     },
        //     body: JSON.stringify(notification)
        // });

    } catch (error) {
        console.error("[通知中心] Webhook发送失败:", error);
    }
}

// 工具函数：获取用户通知偏好
function getUserNotificationPreferences(userId) {
    try {
        const prefs = app.db.queryOne(`
            SELECT notification_preferences FROM user_preferences WHERE user_id = ?
        `, [userId]);

        return prefs ? JSON.parse(prefs.notification_preferences || '{}') : {
            enable_url_notifications: true,
            enable_security_notifications: true,
            url_channels: [NOTIFICATION_CHANNELS.IN_APP],
            security_channels: [NOTIFICATION_CHANNELS.IN_APP, NOTIFICATION_CHANNELS.EMAIL]
        };
    } catch (error) {
        console.error("[通知中心] 获取用户偏好失败:", error);
        return null;
    }
}

// 工具函数：检查是否首次登录
function checkFirstLogin(userId) {
    try {
        const result = app.db.queryOne(`
            SELECT COUNT(*) as login_count FROM user_login_logs WHERE user_id = ?
        `, [userId]);

        return result.login_count === 0;
    } catch (error) {
        console.error("[通知中心] 检查首次登录失败:", error);
        return false;
    }
}

// 工具函数：记录首次登录
function recordFirstLogin(userId) {
    try {
        app.db.execute(`
            INSERT INTO user_login_logs (user_id, login_time, is_first_login)
            VALUES (?, datetime('now'), 1)
        `, [userId]);
    } catch (error) {
        console.error("[通知中心] 记录首次登录失败:", error);
    }
}

// 工具函数：检测可疑活动
function detectSuspiciousActivity(e) {
    const suspicious = {
        isSuspicious: false,
        riskScore: 0,
        reasons: []
    };

    // 检查频繁请求
    if (isHighFrequencyRequest(e.ip)) {
        suspicious.isSuspicious = true;
        suspicious.riskScore += 3;
        suspicious.reasons.push('高频请求');
    }

    // 检查异常时间
    if (isUnusualTime()) {
        suspicious.isSuspicious = true;
        suspicious.riskScore += 1;
        suspicious.reasons.push('异常时间访问');
    }

    // 检查敏感操作
    if (isSensitiveOperation(e.path)) {
        suspicious.isSuspicious = true;
        suspicious.riskScore += 2;
        suspicious.reasons.push('敏感操作');
    }

    return suspicious;
}

// 工具函数：记录安全事件
function recordSecurityEvent(suspicious) {
    try {
        app.db.execute(`
            INSERT INTO security_events (
                ip, path, method, risk_score, reasons, created_at
            ) VALUES (?, ?, ?, ?, ?, datetime('now'))
        `, [
            suspicious.ip,
            suspicious.path,
            suspicious.method,
            suspicious.riskScore,
            JSON.stringify(suspicious.reasons)
        ]);
    } catch (error) {
        console.error("[通知中心] 记录安全事件失败:", error);
    }
}

// 工具函数：获取用户通知
function getUserNotifications(userId, options = {}) {
    try {
        let whereClause = 'WHERE user_id = ?';
        const params = [userId];

        if (options.unreadOnly) {
            whereClause += ' AND read = 0';
        }

        if (options.type) {
            whereClause += ' AND type = ?';
            params.push(options.type);
        }

        const offset = (options.page - 1) * options.limit;

        const notifications = app.db.query(`
            SELECT * FROM notifications
            ${whereClause}
            ORDER BY priority DESC, created_at DESC
            LIMIT ? OFFSET ?
        `, [...params, options.limit, offset]);

        const totalCount = app.db.queryOne(`
            SELECT COUNT(*) as total FROM notifications ${whereClause}
        `, params);

        return {
            notifications: notifications.map(n => ({
                ...n,
                data: JSON.parse(n.data || '{}'),
                channels: JSON.parse(n.channels || '[]')
            })),
            pagination: {
                page: options.page,
                limit: options.limit,
                total: totalCount.total,
                totalPages: Math.ceil(totalCount.total / options.limit)
            }
        };
    } catch (error) {
        console.error("[通知中心] 获取通知失败:", error);
        return { notifications: [], pagination: { total: 0, totalPages: 0 } };
    }
}

// 工具函数：标记通知为已读
function markNotificationAsRead(notificationId, userId) {
    try {
        app.db.execute(`
            UPDATE notifications
            SET read = 1, read_at = datetime('now')
            WHERE id = ? AND user_id = ?
        `, [notificationId, userId]);
    } catch (error) {
        console.error("[通知中心] 标记已读失败:", error);
    }
}

// 工具函数：标记所有通知为已读
function markAllNotificationsAsRead(userId) {
    try {
        const result = app.db.execute(`
            UPDATE notifications
            SET read = 1, read_at = datetime('now')
            WHERE user_id = ? AND read = 0
        `, [userId]);

        return { markedCount: result.changes };
    } catch (error) {
        console.error("[通知中心] 批量标记已读失败:", error);
        return { markedCount: 0 };
    }
}

// 工具函数：获取通知统计
function getNotificationStats(userId) {
    try {
        const total = app.db.queryOne(`
            SELECT COUNT(*) as count FROM notifications WHERE user_id = ?
        `, [userId]);

        const unread = app.db.queryOne(`
            SELECT COUNT(*) as count FROM notifications WHERE user_id = ? AND read = 0
        `, [userId]);

        const byType = app.db.query(`
            SELECT type, COUNT(*) as count FROM notifications
            WHERE user_id = ?
            GROUP BY type
        `, [userId]);

        return {
            total: total.count,
            unread: unread.count,
            by_type: byType
        };
    } catch (error) {
        console.error("[通知中心] 获取统计失败:", error);
        return { total: 0, unread: 0, by_type: [] };
    }
}

// 工具函数：清理过期通知
function cleanupExpiredNotifications() {
    try {
        const thirtyDaysAgo = new Date();
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);

        const result = app.db.execute(`
            DELETE FROM notifications
            WHERE created_at < ? AND read = 1
        `, [thirtyDaysAgo.toISOString()]);

        return { deletedCount: result.changes };
    } catch (error) {
        console.error("[通知中心] 清理通知失败:", error);
        return { deletedCount: 0 };
    }
}

// 工具函数：生成每日摘要
function generateDailySummary() {
    try {
        const yesterday = new Date();
        yesterday.setDate(yesterday.getDate() - 1);
        const yesterdayStr = yesterday.toISOString().split('T')[0];

        const stats = app.db.queryOne(`
            SELECT
                COUNT(*) as total_notifications,
                COUNT(CASE WHEN read = 0 THEN 1 END) as unread_notifications,
                COUNT(CASE WHEN type = 'security' THEN 1 END) as security_notifications
            FROM notifications
            WHERE DATE(created_at) = ?
        `, [yesterdayStr]);

        return {
            date: yesterdayStr,
            total: stats.total_notifications || 0,
            unread: stats.unread_notifications || 0,
            security: stats.security_notifications || 0
        };
    } catch (error) {
        console.error("[通知中心] 生成每日摘要失败:", error);
        return null;
    }
}

// 工具函数：发送每日摘要给用户
function sendDailySummaryToUsers(summary) {
    try {
        if (!summary || summary.total === 0) {
            return;
        }

        // 获取活跃用户列表
        const activeUsers = app.db.query(`
            SELECT DISTINCT user_id FROM notifications
            WHERE DATE(created_at) = ?
            LIMIT 100
        `, [summary.date]);

        activeUsers.forEach(user => {
            const dailyNotification = createNotification({
                type: NOTIFICATION_TYPES.SYSTEM,
                priority: PRIORITY_LEVELS.LOW,
                title: '每日通知摘要',
                message: `昨日您收到 ${summary.total} 条通知，其中 ${summary.unread} 条未读`,
                data: {
                    user_id: user.user_id,
                    summary: summary
                },
                channels: [NOTIFICATION_CHANNELS.IN_APP]
            });

            sendNotification(dailyNotification);
        });

    } catch (error) {
        console.error("[通知中心] 发送每日摘要失败:", error);
    }
}

// 辅助函数
function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

function getUserEmail(userId) {
    try {
        const user = app.db.queryOne(`
            SELECT email FROM users WHERE id = ?
        `, [userId]);
        return user ? user.email : null;
    } catch (error) {
        console.error("[通知中心] 获取用户邮箱失败:", error);
        return null;
    }
}

function isHighFrequencyRequest(ip) {
    // 简化的高频请求检测
    // 实际项目中应该使用更复杂的算法
    return false;
}

function isUnusualTime() {
    const hour = new Date().getHours();
    return hour < 6 || hour > 23;
}

function isSensitiveOperation(path) {
    const sensitivePaths = ['/admin', '/delete', '/batch', '/security'];
    return sensitivePaths.some(sensitivePath => path.includes(sensitivePath));
}

console.log("[通知中心] 插件初始化完成");
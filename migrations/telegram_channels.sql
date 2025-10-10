-- 创建 Telegram 频道/群组表
CREATE TABLE telegram_channels (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- Telegram 频道/群组信息
    chat_id BIGINT NOT NULL COMMENT 'Telegram 聊天ID',
    chat_name VARCHAR(255) NOT NULL COMMENT '聊天名称',
    chat_type VARCHAR(50) NOT NULL COMMENT '类型：channel/group',

    -- 推送配置
    push_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用推送',
    push_frequency INT DEFAULT 24 COMMENT '推送频率（小时）',
    content_categories TEXT COMMENT '推送的内容分类，用逗号分隔',
    content_tags TEXT COMMENT '推送的标签，用逗号分隔',

    -- 频道状态
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否活跃',
    last_push_at TIMESTAMP NULL COMMENT '最后推送时间',

    -- 注册信息
    registered_by VARCHAR(100) COMMENT '注册者用户名',
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
   
    -- API配置
    api VARCHAR(255) COMMENT 'API地址',
    token VARCHAR(255) COMMENT '访问令牌',
    api_type VARCHAR(50) COMMENT 'API类型',
    is_push_saved_info BOOLEAN DEFAULT FALSE COMMENT '是否只推送已转存资源',

    -- 资源策略和时间限制配置
    resource_strategy VARCHAR(20) DEFAULT 'random' COMMENT '资源策略：latest-最新优先,transferred-已转存优先,random-纯随机',
    time_limit VARCHAR(20) DEFAULT 'none' COMMENT '时间限制：none-无限制,week-一周内,month-一月内',
    push_start_time VARCHAR(10) COMMENT '推送开始时间，格式HH:mm',
    push_end_time VARCHAR(10) COMMENT '推送结束时间，格式HH:mm',
   
    -- 索引
    INDEX idx_chat_id (chat_id),
    INDEX idx_chat_type (chat_type),
    INDEX idx_is_active (is_active),
    INDEX idx_push_enabled (push_enabled),
    INDEX idx_registered_at (registered_at),
    INDEX idx_last_push_at (last_push_at),

    UNIQUE KEY uk_chat_id (chat_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Telegram 频道/群组表';
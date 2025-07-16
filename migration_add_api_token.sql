-- 为system_configs表添加api_token字段
-- 执行时间: 2024-12-19

-- 添加api_token字段
ALTER TABLE system_configs 
ADD COLUMN api_token VARCHAR(100) UNIQUE;

-- 为现有记录生成默认的api_token
UPDATE system_configs 
SET api_token = CONCAT('api_', MD5(RANDOM()::text), '_', EXTRACT(EPOCH FROM NOW())::bigint)
WHERE api_token IS NULL;

-- 添加索引以提高查询性能
CREATE INDEX idx_system_configs_api_token ON system_configs(api_token);

-- 添加注释
COMMENT ON COLUMN system_configs.api_token IS '公开API访问令牌，用于API认证'; 
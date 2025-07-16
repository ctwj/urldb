-- 添加used_space字段
ALTER TABLE cks ADD COLUMN IF NOT EXISTS used_space BIGINT DEFAULT 0;

-- 更新现有数据，将字符串类型的容量字段转换为bigint
-- 注意：这里需要根据实际情况调整转换逻辑
UPDATE cks SET 
    used_space = CASE 
        WHEN used_space_text IS NOT NULL AND used_space_text != '' THEN 
            CASE 
                WHEN used_space_text LIKE '%GB%' THEN 
                    CAST(REPLACE(REPLACE(used_space_text, 'GB', ''), ' ', '') AS BIGINT) * 1024 * 1024 * 1024
                WHEN used_space_text LIKE '%MB%' THEN 
                    CAST(REPLACE(REPLACE(used_space_text, 'MB', ''), ' ', '') AS BIGINT) * 1024 * 1024
                WHEN used_space_text LIKE '%KB%' THEN 
                    CAST(REPLACE(REPLACE(used_space_text, 'KB', ''), ' ', '') AS BIGINT) * 1024
                ELSE 0
            END
        ELSE 0
    END
WHERE used_space_text IS NOT NULL;

-- 删除旧的字符串类型字段（可选，建议先备份数据）
-- ALTER TABLE cks DROP COLUMN IF EXISTS capacity;
-- ALTER TABLE cks DROP COLUMN IF EXISTS used_space_text;
-- ALTER TABLE cks DROP COLUMN IF EXISTS total_space_text;

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_cks_used_space ON cks(used_space);
CREATE INDEX IF NOT EXISTS idx_cks_space_left_space ON cks(space, left_space); 
-- 添加文件哈希字段
ALTER TABLE files ADD COLUMN file_hash VARCHAR(64) COMMENT '文件哈希值';
CREATE UNIQUE INDEX idx_files_hash ON files(file_hash);

-- 添加同步状态字段
ALTER TABLE resources ADD COLUMN synced_to_meilisearch BOOLEAN DEFAULT FALSE;
ALTER TABLE resources ADD COLUMN synced_at TIMESTAMP NULL;

-- 创建索引以提高查询性能
CREATE INDEX idx_resources_synced ON resources(synced_to_meilisearch, synced_at);

-- 添加注释
COMMENT ON COLUMN resources.synced_to_meilisearch IS '是否已同步到Meilisearch';
COMMENT ON COLUMN resources.synced_at IS '同步时间'; 
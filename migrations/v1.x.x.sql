-- 添加文件哈希字段
ALTER TABLE files ADD COLUMN file_hash VARCHAR(64) COMMENT '文件哈希值';
CREATE UNIQUE INDEX idx_files_hash ON files(file_hash);
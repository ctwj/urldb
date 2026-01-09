-- 添加AI和SEO相关的字段到resources表
ALTER TABLE resources ADD COLUMN IF NOT EXISTS seo_keywords TEXT[];
ALTER TABLE resources ADD COLUMN IF NOT EXISTS seo_title VARCHAR(500);
ALTER TABLE resources ADD COLUMN IF NOT EXISTS seo_description TEXT;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS ai_model_used VARCHAR(100);
ALTER TABLE resources ADD COLUMN IF NOT EXISTS ai_generation_status VARCHAR(20) DEFAULT 'none';
ALTER TABLE resources ADD COLUMN IF NOT EXISTS ai_generation_timestamp TIMESTAMP;
ALTER TABLE resources ADD COLUMN IF NOT EXISTS ai_last_regeneration TIMESTAMP;
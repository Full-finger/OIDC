-- 创建番剧表和用户收藏表
-- 番剧表
CREATE TABLE animes (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    episode_count INTEGER,
    director VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 为animes表创建触发器函数（如果不存在）
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为animes表创建触发器，自动更新updated_at字段
CREATE TRIGGER trigger_animes_updated_at
    BEFORE UPDATE ON animes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 用户收藏表
CREATE TABLE user_collections (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    anime_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    rating INTEGER CHECK (rating >= 1 AND rating <= 10),
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (anime_id) REFERENCES animes(id) ON DELETE CASCADE
);

-- 为user_collections表创建触发器，自动更新updated_at字段
CREATE TRIGGER trigger_user_collections_updated_at
    BEFORE UPDATE ON user_collections
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 为user_collections表创建唯一索引，防止用户重复收藏同一番剧
CREATE UNIQUE INDEX idx_user_collections_user_anime ON user_collections(user_id, anime_id);
-- docs/schema.sql

-- 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    nickname VARCHAR(100),
    avatar_url VARCHAR(255),
    bio TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 自动更新 updated_at 的函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为 users 表创建触发器
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- OAuth客户端表
CREATE TABLE oauth_clients (
    id BIGSERIAL PRIMARY KEY,
    client_id VARCHAR(255) UNIQUE NOT NULL,
    client_secret_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    redirect_uris TEXT[] NOT NULL,
    scopes TEXT[] NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 为oauth_clients表的client_id字段创建索引以提高查询性能
CREATE INDEX idx_oauth_clients_client_id ON oauth_clients(client_id);

-- OAuth授权码表
CREATE TABLE oauth_authorization_codes (
    code VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    redirect_uri VARCHAR(255) NOT NULL,
    scopes TEXT[] NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    code_challenge VARCHAR(255),  -- PKCE, 可选但推荐
    code_challenge_method VARCHAR(10)
);

-- 为oauth_authorization_codes表的相关字段创建索引以提高查询性能
CREATE INDEX idx_oauth_authorization_codes_client_id ON oauth_authorization_codes(client_id);
CREATE INDEX idx_oauth_authorization_codes_user_id ON oauth_authorization_codes(user_id);
CREATE INDEX idx_oauth_authorization_codes_expires_at ON oauth_authorization_codes(expires_at);

-- OAuth刷新令牌表
CREATE TABLE oauth_refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    client_id VARCHAR(255) NOT NULL,
    scopes TEXT[] NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ
);

-- 为oauth_refresh_tokens表的相关字段创建索引以提高查询性能
CREATE INDEX idx_oauth_refresh_tokens_token_hash ON oauth_refresh_tokens(token_hash);
CREATE INDEX idx_oauth_refresh_tokens_client_id ON oauth_refresh_tokens(client_id);
CREATE INDEX idx_oauth_refresh_tokens_user_id ON oauth_refresh_tokens(user_id);
CREATE INDEX idx_oauth_refresh_tokens_expires_at ON oauth_refresh_tokens(expires_at);

-- 表结构说明注释
COMMENT ON TABLE users IS '用户基本信息表';
COMMENT ON TABLE oauth_clients IS 'OAuth客户端信息表';
COMMENT ON TABLE oauth_authorization_codes IS 'OAuth授权码表';

COMMENT ON COLUMN oauth_clients.client_id IS '客户端公开标识';
COMMENT ON COLUMN oauth_clients.client_secret_hash IS '加密后的客户端密钥';
COMMENT ON COLUMN oauth_clients.name IS '应用名称';
COMMENT ON COLUMN oauth_clients.redirect_uris IS '允许的重定向URI列表';
COMMENT ON COLUMN oauth_clients.scopes IS '允许的scope列表';

COMMENT ON COLUMN oauth_authorization_codes.code IS '授权码';
COMMENT ON COLUMN oauth_authorization_codes.client_id IS '关联的客户端ID';
COMMENT ON COLUMN oauth_authorization_codes.user_id IS '关联的用户ID';
COMMENT ON COLUMN oauth_authorization_codes.redirect_uri IS '重定向URI';
COMMENT ON COLUMN oauth_authorization_codes.scopes IS '授权的scope列表';
COMMENT ON COLUMN oauth_authorization_codes.expires_at IS '过期时间';
COMMENT ON COLUMN oauth_authorization_codes.code_challenge IS 'PKCE挑战码';
COMMENT ON COLUMN oauth_authorization_codes.code_challenge_method IS 'PKCE挑战码方法';

-- 表结构说明注释
COMMENT ON TABLE oauth_refresh_tokens IS 'OAuth刷新令牌表';
COMMENT ON COLUMN oauth_refresh_tokens.token_hash IS '加密后的刷新令牌';
COMMENT ON COLUMN oauth_refresh_tokens.user_id IS '关联的用户ID';
COMMENT ON COLUMN oauth_refresh_tokens.client_id IS '关联的客户端ID';
COMMENT ON COLUMN oauth_refresh_tokens.scopes IS '授权的scope列表';
COMMENT ON COLUMN oauth_refresh_tokens.expires_at IS '过期时间';
COMMENT ON COLUMN oauth_refresh_tokens.revoked_at IS '撤销时间';

-- 番剧表
CREATE TABLE animes (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    episode_count INTEGER,
    director VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 为animes表创建触发器，自动更新updated_at字段
CREATE TRIGGER trigger_animes_updated_at
    BEFORE UPDATE ON animes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 表结构说明注释
COMMENT ON TABLE animes IS '番剧信息表';
COMMENT ON COLUMN animes.id IS '番剧ID';
COMMENT ON COLUMN animes.title IS '番剧名';
COMMENT ON COLUMN animes.episode_count IS '话数';
COMMENT ON COLUMN animes.director IS '导演';
COMMENT ON COLUMN animes.created_at IS '创建时间';
COMMENT ON COLUMN animes.updated_at IS '更新时间';

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

-- 表结构说明注释
COMMENT ON TABLE user_collections IS '用户番剧收藏表';
COMMENT ON COLUMN user_collections.id IS '收藏ID';
COMMENT ON COLUMN user_collections.user_id IS '用户ID';
COMMENT ON COLUMN user_collections.anime_id IS '番剧ID';
COMMENT ON COLUMN user_collections.type IS '收藏类型（想看、在看、已看）';
COMMENT ON COLUMN user_collections.rating IS '评分（1-10）';
COMMENT ON COLUMN user_collections.comment IS '评论';
COMMENT ON COLUMN user_collections.created_at IS '创建时间';
COMMENT ON COLUMN user_collections.updated_at IS '更新时间';

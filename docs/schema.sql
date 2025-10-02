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
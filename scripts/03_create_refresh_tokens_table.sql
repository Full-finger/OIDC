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
COMMENT ON TABLE oauth_refresh_tokens IS 'OAuth刷新令牌表';
COMMENT ON COLUMN oauth_refresh_tokens.token_hash IS '加密后的刷新令牌';
COMMENT ON COLUMN oauth_refresh_tokens.user_id IS '关联的用户ID';
COMMENT ON COLUMN oauth_refresh_tokens.client_id IS '关联的客户端ID';
COMMENT ON COLUMN oauth_refresh_tokens.scopes IS '授权的scope列表';
COMMENT ON COLUMN oauth_refresh_tokens.expires_at IS '过期时间';
COMMENT ON COLUMN oauth_refresh_tokens.revoked_at IS '撤销时间';
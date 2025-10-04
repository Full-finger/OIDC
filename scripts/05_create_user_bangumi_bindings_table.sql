-- 创建用户与Bangumi账号绑定表
-- 用于存储用户与Bangumi账号的绑定关系和认证信息

CREATE TABLE user_bangumi_bindings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    bangumi_user_id BIGINT NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    token_expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 为user_bangumi_bindings表创建索引
-- 提高根据用户ID查询绑定信息的性能
CREATE INDEX idx_user_bangumi_bindings_user_id ON user_bangumi_bindings(user_id);

-- 为user_bangumi_bindings表创建索引
-- 提高根据Bangumi用户ID查询绑定信息的性能
CREATE INDEX idx_user_bangumi_bindings_bangumi_user_id ON user_bangumi_bindings(bangumi_user_id);

-- 为user_bangumi_bindings表创建唯一索引
-- 确保每个用户只能绑定一个Bangumi账号
CREATE UNIQUE INDEX idx_user_bangumi_bindings_user_unique ON user_bangumi_bindings(user_id);

-- 为user_bangumi_bindings表创建唯一索引
-- 确保每个Bangumi账号只能被一个用户绑定
CREATE UNIQUE INDEX idx_user_bangumi_bindings_bangumi_user_unique ON user_bangumi_bindings(bangumi_user_id);

-- 为user_bangumi_bindings表创建触发器函数（如果不存在）
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为user_bangumi_bindings表创建触发器，自动更新updated_at字段
CREATE TRIGGER trigger_user_bangumi_bindings_updated_at
    BEFORE UPDATE ON user_bangumi_bindings
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();
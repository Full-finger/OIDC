# OIDC服务实现

这是一个基于Go语言和Gin框架实现的OpenID Connect(OIDC)服务，同时包含番剧收藏和Bangumi账号绑定功能。

## 功能特性

1. 用户注册、登录和邮箱验证
2. OAuth 2.0 授权码流程实现
3. OpenID Connect核心功能（提供ID Token）
4. 番剧收藏管理
5. Bangumi账号绑定和数据同步

## 技术栈

- Go 1.24.5
- Gin Web框架
- PostgreSQL数据库
- Redis缓存
- JWT token
- Docker容器化部署

## 项目结构

```
OIDC/
├── cmd/                 # 应用入口
├── config/              # 配置文件
├── docs/                # 文档
├── frontend/            # 前端代码
├── internal/            # 核心代码
│   ├── handler/         # 控制器层
│   ├── service/         # 服务层
│   ├── repository/      # 仓库层
│   ├── mapper/          # 数据映射层
│   ├── model/           # 数据模型
│   ├── util/            # 工具类
│   ├── middleware/      # 中间件
│   └── router/          # 路由配置
├── scripts/             # 脚本文件
└── docs/                # 项目文档
```

## 快速开始

### 环境准备

1. 安装Docker和Docker Compose
2. 安装Go 1.24.5或更高版本

### 启动服务

1. 启动数据库和Redis：
```bash
docker-compose up -d
```

2. 生成JWT密钥对：
```bash
go run scripts/generate_jwt_keys.go
```

3. 启动应用：
```bash
go run cmd/main.go
```

## API端点

### 用户相关
- `POST /api/v1/register` - 用户注册
- `POST /api/v1/login` - 用户登录
- `GET /api/v1/verify` - 邮箱验证

### OAuth 2.0 / OIDC相关
- `GET /.well-known/openid-configuration` - OIDC服务发现
- `GET /oauth/authorize` - 授权端点
- `POST /oauth/token` - 令牌端点
- `GET /oauth/userinfo` - 用户信息端点

### 番剧相关
- `GET /api/v1/anime/:id` - 获取番剧详情
- `GET /api/v1/anime/search` - 搜索番剧
- `GET /api/v1/anime/list` - 列出所有番剧
- `GET /api/v1/anime/status` - 根据状态列出番剧
- `POST /api/v1/anime/` - 创建番剧
- `PUT /api/v1/anime/:id` - 更新番剧
- `DELETE /api/v1/anime/:id` - 删除番剧

### 收藏相关
- `POST /api/v1/collection/` - 添加番剧到收藏
- `GET /api/v1/collection/:anime_id` - 获取用户对某个番剧的收藏
- `PUT /api/v1/collection/:anime_id` - 更新收藏信息
- `DELETE /api/v1/collection/:anime_id` - 从收藏中移除番剧
- `GET /api/v1/collection/` - 列出用户的所有收藏
- `GET /api/v1/collection/status` - 根据状态列出用户的收藏
- `GET /api/v1/collection/favorites` - 列出用户的收藏夹

### Bangumi绑定相关
- `GET /api/v1/bangumi/authorize` - 发起Bangumi授权
- `GET /api/v1/bangumi/callback` - Bangumi授权回调
- `DELETE /api/v1/bangumi/unbind` - 解绑Bangumi账号
- `GET /api/v1/bangumi/account` - 获取已绑定的Bangumi账号信息
- `POST /api/v1/bangumi/sync` - 同步Bangumi收藏数据

## 测试

项目包含多种测试脚本用于验证各功能模块：

1. `go run cmd/test_oidc_client.go` - 测试OIDC客户端完整流程
2. `go run cmd/test_anime_collection.go` - 测试番剧和收藏功能
3. `go run cmd/test_anime_crud.go` - 测试番剧增删改查功能

注意：运行测试前需要：
1. 确保服务已启动
2. 生成JWT密钥对（`go run scripts/generate_jwt_keys.go`）
3. 数据库服务正常运行

## 环境变量配置

请参考`.env.example`文件配置必要的环境变量。

## 数据库设计

数据库表结构定义在`scripts/init.sql`文件中，包含用户、验证令牌、OAuth客户端、授权码、刷新令牌、番剧、收藏和Bangumi账号绑定等表。

## 开发说明

项目遵循严格的分层架构设计：
1. Model层：定义数据结构
2. Mapper层：数据访问接口
3. Repository层：业务数据访问
4. Service层：业务逻辑实现
5. Handler层：HTTP请求处理
6. Router层：路由配置

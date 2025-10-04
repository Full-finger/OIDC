# OIDC ID Token 结构设计文档

## 1. OIDC 概述

OpenID Connect (OIDC) 是建立在 OAuth 2.0 协议之上的简单身份层。它允许客户端根据授权服务器执行的身份验证来验证最终用户的身份，并以可互操作和类似 REST 的方式获取有关最终用户的基本个人资料信息。

OIDC 扩展了 OAuth 2.0，添加了以下核心组件：
1. ID Token - 一个 JSON Web Token (JWT)，包含用户身份信息
2. UserInfo Endpoint - 一个标准化的 HTTP API，用于获取经过身份验证的用户详细信息
3. Standardized Scopes - 标准化的作用域，用于请求特定类型的用户信息

## 2. ID Token 结构

ID Token 是一个 JWT 格式的令牌，包含以下组成部分：

### 2.1 标准 Claims

以下 Claims 是 OIDC ID Token 中的标准 Claims：

| Claim | 必需 | 描述 |
|-------|------|------|
| iss | 是 | Issuer Identifier，标识 ID Token 的签发者 |
| sub | 是 | Subject Identifier，用户的唯一标识符 |
| aud | 是 | Audience(s)，标识 ID Token 的目标受众，通常是 client_id |
| exp | 是 | Expiration time，ID Token 的过期时间戳 |
| iat | 是 | Issued at time，ID Token 的签发时间戳 |
| auth_time | 推荐 | Authentication time，用户身份验证发生的时间 |
| nonce | 条件性 | 如果请求中包含 nonce 参数，则必需包含在 ID Token 中 |
| acr | 否 | Authentication Context Class Reference，认证上下文类引用 |
| amr | 否 | Authentication Methods References，认证方法引用 |
| azp | 条件性 | Authorized party，当 Audience 包含多个值时，标识授权的一方 |

### 2.2 根据 Scope 扩展的 Claims

根据请求的 Scope，ID Token 中可以包含额外的 Claims：

#### 2.2.1 openid Scope

openid 是 OIDC 的核心 Scope，必须在所有 OIDC 请求中包含。它表示客户端是一个 OIDC 客户端并请求一个 ID Token。

#### 2.2.2 profile Scope

当请求 profile Scope 时，ID Token 中可以包含以下 Claims：

| Claim | 描述 |
|-------|------|
| name | 用户的全名 |
| family_name | 用户的姓氏 |
| given_name | 用户的名字 |
| middle_name | 用户的中间名 |
| nickname | 用户的昵称 |
| preferred_username | 用户的首选用户名 |
| profile | 用户个人资料页面的 URL |
| picture | 用户头像的 URL |
| website | 用户网站的 URL |
| gender | 用户的性别 |
| birthdate | 用户的生日 |
| zoneinfo | 用户的时区 |
| locale | 用户的语言环境 |
| updated_at | 用户信息最后一次更新的时间戳 |

#### 2.2.3 email Scope

当请求 email Scope 时，ID Token 中可以包含以下 Claims：

| Claim | 描述 |
|-------|------|
| email | 用户的电子邮件地址 |
| email_verified | 指示用户的电子邮件地址是否已验证 |

#### 2.2.4 address Scope

当请求 address Scope 时，ID Token 中可以包含以下 Claims：

| Claim | 描述 |
|-------|------|
| address | 用户的邮寄地址，是一个结构化的 JSON 对象 |

#### 2.2.5 phone Scope

当请求 phone Scope 时，ID Token 中可以包含以下 Claims：

| Claim | 描述 |
|-------|------|
| phone_number | 用户的电话号码 |
| phone_number_verified | 指示用户的电话号码是否已验证 |

## 3. ID Token 示例

一个典型的 OIDC ID Token 可能如下所示：

```json
{
  "iss": "http://localhost:8080",
  "sub": "1234567890",
  "aud": "test_client",
  "exp": 1688888888,
  "iat": 1688880000,
  "auth_time": 1688880000,
  "nonce": "abc123",
  "name": "张三",
  "nickname": "小张",
  "preferred_username": "zhangsan",
  "picture": "http://example.com/avatar.jpg",
  "email": "zhangsan@example.com",
  "email_verified": true
}
```

## 4. 项目中的 ID Token 设计

基于 OIDC 标准和项目需求，我们的 ID Token 将包含以下 Claims：

### 4.1 标准 Claims（必需）

| Claim | 值 | 说明 |
|-------|----|------|
| iss | http://localhost:8080 | 本地开发环境的 Issuer |
| sub | 用户的唯一标识符 | 数据库中的用户 ID |
| aud | client_id | OAuth 客户端的 ID |
| exp | 当前时间 + 有效期 | 通常设置为 1 小时 |
| iat | 当前时间戳 | ID Token 的签发时间 |
| auth_time | 用户认证时间 | 用户登录的时间戳 |
| nonce | 请求中的 nonce 值 | 如果提供，则必须包含 |

### 4.2 根据 Scope 的扩展 Claims

#### 4.2.1 profile Scope Claims

当请求包含 profile Scope 时，ID Token 将包含以下 Claims：

| Claim | 数据来源 | 说明 |
|-------|----------|------|
| name | 用户的全名 | 从用户表中获取 |
| nickname | 用户昵称 | 从用户表中获取 |
| picture | 用户头像URL | 从用户表中获取 |
| preferred_username | 用户首选用户名 | 从用户表中获取 |

#### 4.2.2 email Scope Claims

当请求包含 email Scope 时，ID Token 将包含以下 Claims：

| Claim | 数据来源 | 说明 |
|-------|----------|------|
| email | 用户邮箱地址 | 从用户表中获取 |
| email_verified | 邮箱验证状态 | 从用户表中获取，默认为 false |

## 5. UserInfo Endpoint

UserInfo Endpoint 是一个标准化的 HTTP API，客户端可以使用 Access Token 来获取用户的详细信息。

### 5.1 端点地址

在我们的实现中，UserInfo Endpoint 的地址为：`http://localhost:8080/oauth/userinfo`

### 5.2 请求方式

使用 HTTP GET 或 POST 方法，通过 Authorization Header 传递 Access Token：

```
Authorization: Bearer <access_token>
```

### 5.3 响应格式

响应为 JSON 格式，包含与 ID Token 中相同的 Claims，但可能更加详细。

```json
{
  "sub": "1234567890",
  "name": "张三",
  "nickname": "小张",
  "preferred_username": "zhangsan",
  "picture": "http://example.com/avatar.jpg",
  "email": "zhangsan@example.com",
  "email_verified": true,
  "updated_at": 1688880000
}
```

## 6. 安全考虑

1. ID Token 必须使用 JWT 签名，防止篡改
2. ID Token 应该有合理的过期时间
3. 使用 HTTPS 传输 ID Token 和 Access Token
4. 验证 ID Token 的签名和 Claims
5. 防止重放攻击，通过检查 nonce 和 auth_time
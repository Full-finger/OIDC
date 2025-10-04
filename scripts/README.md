# API 测试脚本说明

## 批处理脚本 (Windows CMD)

### 单独运行测试

1. `test_register.bat` - 测试用户注册
2. `test_login.bat` - 测试用户登录
3. `test_get_profile.bat` - 测试获取用户信息（需要提供JWT token）
4. `test_update_profile.bat` - 测试更新用户信息（需要提供JWT token）
5. `test_oauth_token.bat` - 测试OAuth令牌交换（需要提供授权码）

### 运行所有测试

运行 `test_all.bat` 执行所有测试。注意：由于获取和更新用户信息需要JWT token，该脚本只执行注册和登录测试。

## PowerShell 脚本

1. `test_all.ps1` - 完整的测试流程，包括注册、登录、获取用户信息和更新用户信息，自动处理JWT token传递
2. `test_oauth_token.ps1` - 测试OAuth令牌交换功能，需要提供授权码
3. `oauth_flow_test_simple.ps1` - 完整的OAuth流程测试（简化版）
4. `setup_oauth_client.go` - 设置OAuth测试环境的Go脚本

## Go 测试应用

1. `oauth_client_test.go` - OAuth2.0客户端测试应用，模拟完整的OAuth授权码流程

PowerShell脚本会自动处理JWT token的传递。

## 使用方法

### 方法1：使用CMD
```cmd
# 启动服务器（在项目根目录）
go run cmd/main.go

# 在另一个终端中运行测试脚本（在scripts目录）
test_register.bat
test_login.bat
# 复制登录返回的token，用于以下测试
test_get_profile.bat
test_update_profile.bat
# 复制授权码，用于OAuth测试
test_oauth_token.bat
```

### 方法2：使用PowerShell
```powershell
# 启动服务器（在项目根目录）
go run cmd/main.go

# 在另一个终端中运行测试脚本（在scripts目录）
.\test_all.ps1
# 获取授权码后测试OAuth功能
.\test_oauth_token.ps1 -Code "your_authorization_code"

# 或者运行完整的OAuth流程测试
.\oauth_flow_test_simple.ps1
```

### 方法3：使用OAuth客户端测试应用
```bash
# 启动OIDC服务器（在项目根目录）
go run cmd/main.go

# 在另一个终端中运行OAuth客户端测试应用（在scripts目录）
go run oauth_client_test.go

# 然后在浏览器中访问 http://localhost:9999 开始OAuth流程测试
```

## 完整OAuth 2.0流程测试说明

要测试完整的OAuth 2.0授权码流程，请按照以下步骤操作：

1. 确保服务器正在运行：
   ```cmd
   go run cmd/main.go
   ```

2. 运行简化版的OAuth流程测试脚本：
   ```powershell
   .\oauth_flow_test_simple.ps1
   ```

3. 按照脚本输出的说明，在浏览器中完成OAuth流程

4. 获取授权码后，使用以下命令测试令牌交换：
   ```powershell
   .\test_oauth_token.ps1 -Code "your_authorization_code"
   ```

## OAuth 2.0 客户端测试应用

我们提供了一个完整的OAuth 2.0客户端测试应用，可以模拟整个授权码流程：

1. 启动OIDC服务：
   ```bash
   go run cmd/main.go
   ```

2. 在另一个终端中启动测试客户端：
   ```bash
   cd scripts
   go run oauth_client_test.go
   ```

3. 在浏览器中访问 `http://localhost:9999` 开始测试

4. 点击 "Login with OIDC Service" 链接开始OAuth流程

5. 系统会引导您完成以下步骤：
   - 跳转到OIDC服务登录页面
   - 登录（使用测试用户凭证）
   - 授权应用访问您的数据
   - 重定向回测试客户端
   - 自动使用授权码换取访问令牌
   - 使用访问令牌获取用户信息

## 注意事项

1. 确保服务器正在运行（默认端口8080）
2. 确保数据库连接正常
3. 第一次运行测试时，注册应该会成功
4. 如果使用相同的用户名再次注册，会返回用户已存在的错误
5. 登录时需要提供正确的用户名和密码
6. 获取和更新用户信息需要有效的JWT token
7. OAuth令牌交换需要有效的授权码
8. 在测试OAuth功能之前，需要先通过授权流程获取授权码
9. 测试OAuth功能时，确保先运行`setup_oauth_client.go`来创建测试客户端
10. OAuth客户端测试应用运行在端口9999，确保该端口未被占用
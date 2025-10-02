# API 测试脚本说明

## 批处理脚本 (Windows CMD)

### 单独运行测试

1. `test_register.bat` - 测试用户注册
2. `test_login.bat` - 测试用户登录
3. `test_get_profile.bat` - 测试获取用户信息（需要提供JWT token）
4. `test_update_profile.bat` - 测试更新用户信息（需要提供JWT token）

### 运行所有测试

运行 `test_all.bat` 执行所有测试。注意：由于获取和更新用户信息需要JWT token，该脚本只执行注册和登录测试。

## PowerShell 脚本

运行 `test_all.ps1` 可以执行完整的测试流程，包括：
1. 用户注册
2. 用户登录
3. 获取用户信息
4. 更新用户信息

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
```

### 方法2：使用PowerShell
```powershell
# 启动服务器（在项目根目录）
go run cmd/main.go

# 在另一个终端中运行测试脚本（在scripts目录）
.\test_all.ps1
```

## 注意事项

1. 确保服务器正在运行（默认端口8080）
2. 确保数据库连接正常
3. 第一次运行测试时，注册应该会成功
4. 如果使用相同的用户名再次注册，会返回用户已存在的错误
5. 登录时需要提供正确的用户名和密码
6. 获取和更新用户信息需要有效的JWT token
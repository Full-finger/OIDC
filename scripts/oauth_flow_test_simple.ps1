# OAuth 2.0 Authorization Code Flow Test Script (Simplified)

Write-Host "OAuth 2.0 Authorization Code Flow Test" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""

# Step 1: Setup test environment
Write-Host "Step 1: Setting up test environment..." -ForegroundColor Yellow
go run scripts/setup_oauth_test.go
Write-Host ""

# Step 2: User registration
Write-Host "Step 2: Registering a test user..." -ForegroundColor Yellow
try {
    $registerResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/register" -Method POST -Body @{
        username = "oauth_test_user"
        email = "oauth_test@example.com"
        password = "oauth_test_password"
    } -ContentType "application/json"

    Write-Host "Registration response:" -ForegroundColor Cyan
    $registerResponse | ConvertTo-Json -Depth 5
    Write-Host ""
} catch {
    Write-Host "Registration failed (user may already exist):" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
}

# Step 3: User login to get JWT token
Write-Host "Step 3: Logging in to get JWT token..." -ForegroundColor Yellow
try {
    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" -Method POST -Body @{
        username = "oauth_test_user"
        password = "oauth_test_password"
    } -ContentType "application/json"

    Write-Host "Login response:" -ForegroundColor Cyan
    $loginResponse | ConvertTo-Json -Depth 5
    Write-Host ""

    $jwtToken = $loginResponse.token
    Write-Host "JWT Token: $jwtToken" -ForegroundColor Magenta
    Write-Host ""
} catch {
    Write-Host "Login failed:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host ""
    exit 1
}

# Step 4: Instructions for OAuth flow
Write-Host "Step 4: Manual OAuth Flow Instructions" -ForegroundColor Yellow
Write-Host ""
Write-Host "1. Open your browser and navigate to:" -ForegroundColor Gray
Write-Host "   http://localhost:8080/oauth/authorize?response_type=code&client_id=test_client&redirect_uri=http://localhost:3000/callback&scope=read&state=xyz" -ForegroundColor Gray
Write-Host ""
Write-Host "2. You will be redirected to the login page. Use these credentials:" -ForegroundColor Gray
Write-Host "   Username: oauth_test_user" -ForegroundColor Gray
Write-Host "   Password: oauth_test_password" -ForegroundColor Gray
Write-Host ""
Write-Host "3. After logging in, you'll see the authorization page. Click '同意' (Agree)" -ForegroundColor Gray
Write-Host ""
Write-Host "4. You will be redirected to:" -ForegroundColor Gray
Write-Host "   http://localhost:3000/callback?code=AUTHORIZATION_CODE&state=xyz" -ForegroundColor Gray
Write-Host ""
Write-Host "5. Copy the AUTHORIZATION_CODE from the URL" -ForegroundColor Gray
Write-Host ""
Write-Host "6. Run the token exchange test script with the code:" -ForegroundColor Gray
Write-Host "   .\test_oauth_token.ps1 -Code `"YOUR_AUTHORIZATION_CODE`"" -ForegroundColor Gray
Write-Host ""

Write-Host "Follow the instructions above to complete the OAuth flow test!" -ForegroundColor Green
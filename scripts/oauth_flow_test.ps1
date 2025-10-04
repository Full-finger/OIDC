# OAuth 2.0 Authorization Code Flow Test Script

Write-Host "OAuth 2.0 Authorization Code Flow Test" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""

# Step 1: Setup test environment
Write-Host "Step 1: Setting up test environment..." -ForegroundColor Yellow
go run scripts/setup_oauth_test.go
Write-Host ""

# Step 2: User registration
Write-Host "Step 2: Registering a test user..." -ForegroundColor Yellow
$registerResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/register" -Method POST -Body @{
    username = "oauth_test_user"
    email = "oauth_test@example.com"
    password = "oauth_test_password"
} -ContentType "application/json"

Write-Host "Registration response:" -ForegroundColor Cyan
$registerResponse | ConvertTo-Json -Depth 5
Write-Host ""

# Step 3: User login to get JWT token
Write-Host "Step 3: Logging in to get JWT token..." -ForegroundColor Yellow
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

# Step 4: Simulate OAuth authorization request
Write-Host "Step 4: Simulating OAuth authorization request..." -ForegroundColor Yellow
Write-Host "In a real scenario, the user would be redirected to:" -ForegroundColor Gray
Write-Host "http://localhost:8080/oauth/authorize?response_type=code&client_id=test_client&redirect_uri=http://localhost:3000/callback&scope=read&state=xyz" -ForegroundColor Gray
Write-Host ""

# Step 5: Manually create an authorization code for testing
Write-Host "Step 5: Creating a test authorization code in the database..." -ForegroundColor Yellow

# Generate a random code
$code = -join ((65..90) + (97..122) | Get-Random -Count 32 | % {[char]$_})

# Insert test authorization code directly into the database
# In a real scenario, this would be done by the OAuth service
$env:PGHOST = "localhost"
$env:PGUSER = "oidc_user"
$env:PGPASSWORD = "oidc_password"
$env:PGDATABASE = "oidc_db"

psql -c "INSERT INTO oauth_authorization_codes (code, client_id, user_id, redirect_uri, scopes, expires_at) VALUES ('$code', 'test_client', 1, 'http://localhost:3000/callback', '{read}', NOW() + INTERVAL '10 minutes');"

Write-Host "Created authorization code: $code" -ForegroundColor Magenta
Write-Host ""

# Step 6: Exchange authorization code for access token
Write-Host "Step 6: Exchanging authorization code for access token..." -ForegroundColor Yellow

# Test with form data
Write-Host "Testing with form data:" -ForegroundColor Gray
$tokenResponse = Invoke-RestMethod -Uri "http://localhost:8080/oauth/token" -Method POST -Body @{
    grant_type = "authorization_code"
    code = $code
    redirect_uri = "http://localhost:3000/callback"
    client_id = "test_client"
    client_secret = "test_secret"
} -ContentType "application/x-www-form-urlencoded"

Write-Host "Token response:" -ForegroundColor Cyan
$tokenResponse | ConvertTo-Json -Depth 5
Write-Host ""

# Test with Basic Auth
Write-Host "Testing with Basic Auth:" -ForegroundColor Gray
$base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes(("test_client:test_secret")))
$headers = @{
    Authorization = "Basic $base64AuthInfo"
}

$tokenResponse2 = Invoke-RestMethod -Uri "http://localhost:8080/oauth/token" -Method POST -Headers $headers -Body @{
    grant_type = "authorization_code"
    code = $code
    redirect_uri = "http://localhost:3000/callback"
} -ContentType "application/x-www-form-urlencoded"

Write-Host "Token response:" -ForegroundColor Cyan
$tokenResponse2 | ConvertTo-Json -Depth 5
Write-Host ""

Write-Host "OAuth 2.0 flow test completed successfully!" -ForegroundColor Green
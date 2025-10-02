# OIDC API Testing with PowerShell

Write-Host "OIDC API Testing" -ForegroundColor Green
Write-Host "================" -ForegroundColor Green
Write-Host ""

# Test 1: User Registration
Write-Host "1. Testing User Registration" -ForegroundColor Yellow
Write-Host "----------------------------" -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    email = "test@example.com"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/register" -Method POST -Body $registerBody -ContentType "application/json"
Write-Host "Registration Response:" -ForegroundColor Cyan
$response | ConvertTo-Json -Depth 5
Write-Host ""

# Test 2: User Login
Write-Host "2. Testing User Login" -ForegroundColor Yellow
Write-Host "---------------------" -ForegroundColor Yellow
$loginBody = @{
    username = "testuser"
    password = "password123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/login" -Method POST -Body $loginBody -ContentType "application/json"
Write-Host "Login Response:" -ForegroundColor Cyan
$loginResponse | ConvertTo-Json -Depth 5
Write-Host ""

# Extract token for subsequent requests
$token = $loginResponse.token
Write-Host "JWT Token: $token" -ForegroundColor Magenta
Write-Host ""

# Test 3: Get User Profile
Write-Host "3. Testing Get User Profile" -ForegroundColor Yellow
Write-Host "---------------------------" -ForegroundColor Yellow
$headers = @{
    Authorization = "Bearer $token"
}

$profileResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/profile" -Method GET -Headers $headers
Write-Host "Profile Response:" -ForegroundColor Cyan
$profileResponse | ConvertTo-Json -Depth 5
Write-Host ""

# Test 4: Update User Profile
Write-Host "4. Testing Update User Profile" -ForegroundColor Yellow
Write-Host "-----------------------------" -ForegroundColor Yellow
$updateBody = @{
    nickname = "New Nickname"
    bio = "This is my updated bio"
} | ConvertTo-Json

$updateResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/profile" -Method PUT -Body $updateBody -ContentType "application/json" -Headers $headers
Write-Host "Update Profile Response:" -ForegroundColor Cyan
$updateResponse | ConvertTo-Json -Depth 5
Write-Host ""

Write-Host "All tests completed!" -ForegroundColor Green
# OAuth Token Exchange Test Script

param(
    [Parameter(Mandatory=$true)]
    [string]$Code
)

Write-Host "Testing OAuth Token Exchange" -ForegroundColor Green
Write-Host "===========================" -ForegroundColor Green
Write-Host ""

# Test with form data
Write-Host "1. Testing with form data:" -ForegroundColor Yellow
$response = Invoke-RestMethod -Uri "http://localhost:8080/oauth/token" -Method POST -Body @{
    grant_type = "authorization_code"
    code = $Code
    redirect_uri = "http://localhost:3000/callback"
    client_id = "test_client"
    client_secret = "test_secret"
}

Write-Host "Response:" -ForegroundColor Cyan
$response | ConvertTo-Json -Depth 5
Write-Host ""

# Test with Basic Auth
Write-Host "2. Testing with Basic Auth:" -ForegroundColor Yellow
$base64AuthInfo = [Convert]::ToBase64String([Text.Encoding]::ASCII.GetBytes(("test_client:test_secret")))
$headers = @{
    Authorization = "Basic $base64AuthInfo"
}

$response2 = Invoke-RestMethod -Uri "http://localhost:8080/oauth/token" -Method POST -Headers $headers -Body @{
    grant_type = "authorization_code"
    code = $Code
    redirect_uri = "http://localhost:3000/callback"
}

Write-Host "Response:" -ForegroundColor Cyan
$response2 | ConvertTo-Json -Depth 5
Write-Host ""

Write-Host "Test completed!" -ForegroundColor Green
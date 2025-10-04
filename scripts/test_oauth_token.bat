@echo off
set /p code=Enter authorization code: 
echo Testing OAuth Token Exchange
echo ===========================
curl -X POST http://localhost:8080/oauth/token ^
  -H "Content-Type: application/x-www-form-urlencoded" ^
  -d "grant_type=authorization_code" ^
  -d "code=%code%" ^
  -d "redirect_uri=http://localhost:3000/callback" ^
  -d "client_id=test_client" ^
  -d "client_secret=test_secret"
echo.
echo.
pause
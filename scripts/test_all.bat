@echo off
echo OIDC API Testing
echo ================
echo.

echo 1. Testing User Registration
echo ----------------------------
curl -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"email\":\"test@example.com\",\"password\":\"password123\"}"
echo.
echo.

echo 2. Testing User Login
echo ---------------------
curl -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"password\":\"password123\"}"
echo.
echo.

echo 3. Please copy the token from the login response above and use it in the following tests
echo    You can run test_get_profile.bat and test_update_profile.bat separately with the token
echo.

pause
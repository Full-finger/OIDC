@echo off
set /p token=Enter JWT Token: 
echo Testing Get User Profile
echo =======================
curl -X GET http://localhost:8080/api/v1/profile ^
  -H "Authorization: Bearer %token%"
echo.
echo.
pause
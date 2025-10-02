@echo off
set /p token=Enter JWT Token: 
echo Testing Update User Profile
echo ==========================
curl -X PUT http://localhost:8080/api/v1/profile ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %token%" ^
  -d "{\"nickname\":\"New Nickname\",\"bio\":\"This is my bio\"}"
echo.
echo.
pause
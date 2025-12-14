@echo off
echo ====================================
echo   INICIANDO BACKEND - ANDARIA
echo ====================================
echo.
echo Directorio: %CD%
echo.
echo Compilando y ejecutando...
echo.

cd /d "%~dp0"
go run cmd/api/main.go

pause

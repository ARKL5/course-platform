@echo off
chcp 65001 >nul
title 课程平台 - 微服务启动管理器

echo.
echo ================================================
echo           🎯 课程平台微服务启动管理器
echo ================================================
echo.

:: 检查端口函数
call :check_port 3306 "MySQL"
if %errorlevel% neq 0 (
    echo ❌ MySQL 未运行，请先启动 MySQL 服务
    echo 💡 建议：启动 XAMPP 或其他 MySQL 服务
    pause
    exit /b 1
)

call :check_port 6379 "Redis"
if %errorlevel% neq 0 (
    echo ⚠️  Redis 未运行，服务将在无缓存模式下运行
    timeout /t 3 >nul
)

echo.
echo 🚀 开始启动微服务...
echo.

:: 创建日志目录
if not exist "logs" mkdir logs

:: 启动用户服务
echo 🔐 启动用户服务...
start "用户服务" /min powershell -NoExit -Command "cd '%~dp0'; Write-Host '🔐 用户服务启动中...' -ForegroundColor Green; go run cmd/user-service/main.go"
timeout /t 3 >nul

:: 启动课程服务
echo 📖 启动课程服务...
start "课程服务" /min powershell -NoExit -Command "cd '%~dp0'; Write-Host '📖 课程服务启动中...' -ForegroundColor Green; go run cmd/course-service/main.go"
timeout /t 3 >nul

:: 启动内容服务
echo 📄 启动内容服务...
start "内容服务" /min powershell -NoExit -Command "cd '%~dp0'; Write-Host '📄 内容服务启动中...' -ForegroundColor Green; go run cmd/content-service/main.go"
timeout /t 3 >nul

:: 启动API网关
echo 🌐 启动API网关...
start "API网关" /min powershell -NoExit -Command "cd '%~dp0'; Write-Host '🌐 API网关启动中...' -ForegroundColor Green; go run cmd/server/main.go"

echo.
echo ⏳ 等待所有服务启动完成...
timeout /t 10 >nul

echo.
echo 🔍 检查服务状态...

:: 检查各服务端口
call :check_service_status 50051 "用户服务"
call :check_service_status 50052 "课程服务"
call :check_service_status 50053 "内容服务"
call :check_service_status 8083 "API网关"

echo.
echo ========================================
echo 🎉 服务启动完成！
echo.
echo 📱 访问地址：
echo    🏠 首页: http://localhost:8083
echo    📚 课程: http://localhost:8083/dashboard
echo    ✍️  创作者: http://localhost:8083/creator-dashboard
echo    🔐 登录: http://localhost:8083/login
echo    📝 注册: http://localhost:8083/register
echo    📖 API文档: http://localhost:8083/swagger/index.html
echo.
echo 📋 管理命令：
echo    💾 初始化数据: scripts\run-init-data.ps1
echo    🧪 API测试: api-test.ps1
echo    🐳 Docker部署: docker-start.ps1
echo.
echo ⚠️  注意: 关闭此窗口会停止所有服务
echo 🛑 停止服务: 关闭各服务窗口或按 Ctrl+C
echo ========================================
echo.

pause
goto :eof

:: 检查端口是否开放的函数
:check_port
set port=%1
set service_name=%2
netstat -an | find ":%port% " | find "LISTENING" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ %service_name% 端口 %port% 可用
    exit /b 0
) else (
    exit /b 1
)

:: 检查服务状态的函数
:check_service_status
set port=%1
set service_name=%2
netstat -an | find ":%port% " | find "LISTENING" >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ %service_name% 运行正常 ^(端口: %port%^)
) else (
    echo ❌ %service_name% 启动失败 ^(端口: %port%^)
)
goto :eof 
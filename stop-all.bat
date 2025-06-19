@echo off
chcp 65001 >nul
title 课程平台 - 停止所有服务

echo.
echo ================================================
echo           🛑 课程平台微服务停止管理器
echo ================================================
echo.

echo 🔍 正在查找运行中的服务...

:: 停止特定端口的进程
call :kill_process_by_port 50051 "用户服务"
call :kill_process_by_port 50052 "课程服务"
call :kill_process_by_port 50053 "内容服务"
call :kill_process_by_port 8083 "API网关"

echo.
echo 🧹 清理服务窗口...

:: 关闭特定标题的命令窗口
taskkill /f /fi "WINDOWTITLE eq 用户服务*" >nul 2>&1
taskkill /f /fi "WINDOWTITLE eq 课程服务*" >nul 2>&1
taskkill /f /fi "WINDOWTITLE eq 内容服务*" >nul 2>&1
taskkill /f /fi "WINDOWTITLE eq API网关*" >nul 2>&1

echo.
echo ✅ 所有服务已停止
echo.

pause
goto :eof

:: 根据端口杀死进程的函数
:kill_process_by_port
set port=%1
set service_name=%2

for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%port% "') do (
    if not "%%a"=="0" (
        echo 🛑 停止 %service_name% ^(PID: %%a^)
        taskkill /f /pid %%a >nul 2>&1
    )
)
goto :eof 
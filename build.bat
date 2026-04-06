@echo off
chcp 65001 >nul
echo ========================================
echo   抖音工具箱 - Wails 构建脚本
echo ========================================
echo.

:: 检查 Go 版本
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到 Go，请先安装 Go
    pause
    exit /b 1
)

:: 检查 Node.js 版本
node --version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到 Node.js，请先安装 Node.js
    pause
    exit /b 1
)

echo [1/4] 安装前端依赖...
cd frontend
call npm install
if errorlevel 1 (
    echo [错误] 前端依赖安装失败
    pause
    exit /b 1
)
cd ..

echo.
echo [2/4] 安装后端依赖...
call go mod tidy
if errorlevel 1 (
    echo [错误] 后端依赖安装失败
    pause
    exit /b 1
)

echo.
echo [3/4] 生成 Wails 绑定...
call wails generate module
if errorlevel 1 (
    echo [警告] 绑定生成失败，请确保已安装 Wails CLI
    echo   安装命令: go install github.com/wailsapp/wails/v2/cmd/wails@latest
)

echo.
echo [4/4] 开始构建桌面程序...
call wails build
if errorlevel 1 (
    echo [错误] 构建失败
    pause
    exit /b 1
)

echo.
echo ========================================
echo   构建完成！
echo   输出目录: build/bin/
echo ========================================
pause

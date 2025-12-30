@echo off
setlocal enabledelayedexpansion

REM fzjjyz 发布脚本 (Windows)
REM 用于自动化发布流程

set VERSION=%1
set SKIP_TEST=false
set SKIP_BUILD=false
set DRY_RUN=false

REM 颜色代码
set "RED=0C"
set "GREEN=0A"
set "YELLOW=0E"
set "BLUE=09"
set "WHITE=07"

REM 检查参数
if "%VERSION%"=="" (
    echo 用法: release.bat ^<版本号^> [选项]
    echo.
    echo 选项:
    echo   --skip-test    跳过测试
    echo   --skip-build   跳过构建
    echo   --dry-run      试运行
    echo.
    echo 示例:
    echo   release.bat v0.1.1
    echo   release.bat v0.1.1 --dry-run
    exit /b 1
)

REM 解析选项
shift
:parse_args
if "%~1"=="" goto :args_parsed
if "%~1"=="--skip-test" (
    set SKIP_TEST=true
    shift
    goto :parse_args
)
if "%~1"=="--skip-build" (
    set SKIP_BUILD=true
    shift
    goto :parse_args
)
if "%~1"=="--dry-run" (
    set DRY_RUN=true
    shift
    goto :parse_args
)
echo 未知选项: %~1
exit /b 1

:args_parsed

REM 验证版本号格式
echo %VERSION% | findstr /R "^v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*$" >nul
if errorlevel 1 (
    echo [ERROR] 版本号格式错误: %VERSION% ^(应为 v0.1.0 格式^)
    exit /b 1
)

call :color %BLUE% "[INFO] 开始发布流程: %VERSION%"
echo.
echo 选项: 跳过测试=%SKIP_TEST%, 跳过构建=%SKIP_BUILD%, 试运行=%DRY_RUN%
echo.

if "%DRY_RUN%"=="true" (
    call :color %YELLOW% "[WARNING] 这是试运行模式，不会执行实际操作"
    echo.
)

REM 步骤 1: 检查当前分支
call :color %BLUE% "[INFO] 步骤 1: 检查当前分支"
for /f "tokens=*" %%i in ('git branch --show-current') do set CURRENT_BRANCH=%%i
echo 当前分支: %CURRENT_BRANCH%
echo.

REM 步骤 2: 检查未提交的更改
call :color %BLUE% "[INFO] 步骤 2: 检查工作区状态"
if "%DRY_RUN%"=="false" (
    git diff --quiet >nul 2>&1
    if errorlevel 1 (
        call :color %YELLOW% "[WARNING] 检测到未提交的更改"
        git status --short
        set /p "CONTINUE=是否提交这些更改? (y/N): "
        if /i "!CONTINUE!"=="y" (
            git add .
            git commit -m "chore: 发布前准备 %VERSION%"
        ) else (
            call :color %RED% "[ERROR] 请先清理工作区再发布"
            exit /b 1
        )
    )
)

REM 步骤 3: 运行测试
if "%SKIP_TEST%"=="false" (
    call :color %BLUE% "[INFO] 步骤 3: 运行测试"
    if "%DRY_RUN%"=="false" (
        go test ./... -cover
        if errorlevel 1 (
            call :color %RED% "[ERROR] 测试失败"
            exit /b 1
        )
        call :color %GREEN% "[SUCCESS] 测试通过"
    ) else (
        echo [DRY RUN] 跳过测试执行
    )
) else (
    call :color %YELLOW% "[WARNING] 跳过测试"
)
echo.

REM 步骤 4: 更新版本号
call :color %BLUE% "[INFO] 步骤 4: 更新版本号"
set VERSION_NUM=%VERSION:v=%
echo 版本号: %VERSION_NUM%

if "%DRY_RUN%"=="false" (
    if exist "cmd\fzjjyz\main.go" (
        REM 备份原文件
        copy cmd\fzjjyz\main.go cmd\fzjjyz\main.go.bak >nul

        REM 更新版本号 (使用 PowerShell 进行替换)
        powershell -Command "(Get-Content cmd\fzjjyz\main.go) -replace 'const Version = \".*\"', 'const Version = \"%VERSION_NUM%\"' | Set-Content cmd\fzjjyz\main.go"

        REM 验证更新
        findstr /C:"const Version = \"%VERSION_NUM%\"" cmd\fzjjyz\main.go >nul
        if errorlevel 1 (
            call :color %RED% "[ERROR] 版本号更新失败"
            move /Y cmd\fzjjyz\main.go.bak cmd\fzjjyz\main.go >nul
            exit /b 1
        ) else (
            call :color %GREEN% "[SUCCESS] 版本号已更新为 %VERSION_NUM%"
        )

        REM 清理备份
        del cmd\fzjjyz\main.go.bak
    ) else (
        call :color %YELLOW% "[WARNING] 未找到版本定义，跳过更新"
    )
) else (
    echo [DRY RUN] 将更新版本号为 %VERSION_NUM%
)
echo.

REM 步骤 5: 构建二进制文件
if "%SKIP_BUILD%"=="false" (
    call :color %BLUE% "[INFO] 步骤 5: 构建跨平台二进制"
    if "%DRY_RUN%"=="false" (
        if not exist "release\%VERSION%" mkdir release\%VERSION%

        echo 构建 Linux amd64...
        set GOOS=linux
        set GOARCH=amd64
        go build -o release\%VERSION%\fzjjyz_linux_amd64 cmd\fzjjyz

        echo 构建 Windows amd64...
        set GOOS=windows
        set GOARCH=amd64
        go build -o release\%VERSION%\fzjjyz_windows_amd64.exe cmd\fzjjyz

        echo 构建 macOS Intel...
        set GOOS=darwin
        set GOARCH=amd64
        go build -o release\%VERSION%\fzjjyz_darwin_amd64 cmd\fzjjyz

        echo 构建 macOS Apple Silicon...
        set GOOS=darwin
        set GOARCH=arm64
        go build -o release\%VERSION%\fzjjyz_darwin_arm64 cmd\fzjjyz

        echo 生成校验和...
        cd release\%VERSION%
        certutil -hashfile fzjjyz_linux_amd64 SHA256 > checksums.txt
        certutil -hashfile fzjjyz_windows_amd64.exe SHA256 >> checksums.txt
        certutil -hashfile fzjjyz_darwin_amd64 SHA256 >> checksums.txt
        certutil -hashfile fzjjyz_darwin_arm64 SHA256 >> checksums.txt
        cd ..\..

        call :color %GREEN% "[SUCCESS] 构建完成"
        dir release\%VERSION%
    ) else (
        echo [DRY RUN] 将构建各平台二进制
    )
) else (
    call :color %YELLOW% "[WARNING] 跳过构建"
)
echo.

REM 步骤 6: 提交版本更新
call :color %BLUE% "[INFO] 步骤 6: 提交版本更新"
if "%DRY_RUN%"=="false" (
    git add cmd\fzjjyz\main.go
    git commit -m "chore: 发布 %VERSION%" >nul 2>&1
    if errorlevel 0 (
        call :color %GREEN% "[SUCCESS] 版本更新已提交"
    ) else (
        call :color %YELLOW% "[WARNING] 没有需要提交的更改"
    )
) else (
    echo [DRY RUN] 将提交版本更新
)
echo.

REM 步骤 7: 创建 Git 标签
call :color %BLUE% "[INFO] 步骤 7: 创建 Git 标签"
if "%DRY_RUN%"=="false" (
    git tag -a "%VERSION%" -m "Release %VERSION%" >nul 2>&1
    if errorlevel 0 (
        call :color %GREEN% "[SUCCESS] Git 标签 %VERSION% 已创建"
        echo 推送标签到远程...
        git push origin %VERSION% >nul 2>&1
        if errorlevel 0 (
            call :color %GREEN% "[SUCCESS] 标签已推送到远程"
        ) else (
            call :color %YELLOW% "[WARNING] 推送标签失败"
        )
        echo 推送提交到远程...
        git push >nul 2>&1
        if errorlevel 0 (
            call :color %GREEN% "[SUCCESS] 提交已推送到远程"
        ) else (
            call :color %YELLOW% "[WARNING] 推送提交失败，请手动推送"
        )
    ) else (
        call :color %RED% "[ERROR] 创建标签失败"
        exit /b 1
    )
) else (
    echo [DRY RUN] 将创建并推送标签 %VERSION%
)
echo.

REM 步骤 8: 生成发布说明
call :color %BLUE% "[INFO] 步骤 8: 生成发布说明"
if "%DRY_RUN%"=="false" (
    if exist "RELEASE_NOTES.md" (
        echo 使用 RELEASE_NOTES.md 作为发布说明
    ) else (
        echo 生成基础发布说明...
        (
            echo # Release %VERSION%
            echo.
            echo ## 🎉 发布概述
            echo.
            echo **版本**: %VERSION%
            echo **日期**: %date%
            echo **状态**: ✅ 生产就绪
            echo.
            echo ## 📦 下载
            echo.
            echo 从附件下载对应平台的二进制文件：
            echo.
            echo - `fzjjyz_linux_amd64` - Linux 64位
            echo - `fzjjyz_windows_amd64.exe` - Windows 64位
            echo - `fzjjyz_darwin_amd64` - macOS Intel
            echo - `fzjjyz_darwin_arm64` - macOS Apple Silicon
            echo - `checksums.txt` - SHA256 校验和
            echo.
            echo ## 🔍 验证完整性
            echo.
            echo 下载后请验证校验和：
            echo.
            echo `sha256sum -c checksums.txt`
            echo.
            echo ## 🚀 快速开始
            echo.
            echo ### 1. 生成密钥
            echo `./fzjjyz keygen -d ./keys -n mykey`
            echo.
            echo ### 2. 加密文件
            echo `./fzjjyz encrypt -i secret.txt -o secret.fzj -p keys/mykey_public.pem -s keys/mykey_dilithium_private.pem`
            echo.
            echo ### 3. 解密文件
            echo `./fzjjyz decrypt -i secret.fzj -o recovered.txt -p keys/mykey_private.pem -s keys/mykey_dilithium_public.pem`
            echo.
            echo ## 📊 变更详情
            echo.
            echo 请查看 [CHANGELOG.md](../CHANGELOG.md) 获取详细的变更记录。
            echo.
            echo ## 🔐 安全说明
            echo.
            echo 本版本包含后量子加密实现。
            echo 详细安全信息请参考 [SECURITY.md](../SECURITY.md)
            echo.
            echo ## 🤝 贡献
            echo.
            echo 欢迎贡献！请阅读 [CONTRIBUTING.md](../CONTRIBUTING.md) 了解如何参与。
            echo.
            echo ## 📄 许可证
            echo.
            echo MIT License - 详见 [LICENSE](../LICENSE)
            echo.
            echo ---
            echo.
            echo **发布者**: %USERNAME%
            echo **构建时间**: %date% %time%
        ) > release\%VERSION%\release_notes.md
        call :color %GREEN% "[SUCCESS] 发布说明已生成"
    )
) else (
    echo [DRY RUN] 将生成发布说明
)
echo.

REM 步骤 9: 总结
call :color %BLUE% "[INFO] 步骤 9: 发布总结"
echo.
echo ========================================
echo 发布流程完成: %VERSION%
echo ========================================
echo.
if "%DRY_RUN%"=="false" (
    echo ✅ 版本号已更新: %VERSION_NUM%
    echo ✅ 测试已运行
    echo ✅ 二进制已构建
    echo ✅ 校验和已生成
    echo ✅ Git 标签已创建
    echo.
    echo 📁 发布文件位置: release\%VERSION%\
    echo.
    echo 📦 需要手动上传到 GitHub Release:
    dir /b release\%VERSION%\
    echo.
    echo 🔗 下一步:
    echo    1. 访问 GitHub Releases 页面
    echo    2. 创建新 Release: %VERSION%
    echo    3. 上传 release\%VERSION%\ 目录中的所有文件
    echo    4. 使用生成的发布说明或 RELEASE_NOTES.md
) else (
    echo ⚠️  试运行模式 - 未执行实际操作
    echo 💡 使用以下命令执行实际发布:
    echo    release.bat %VERSION%
)
echo.
echo ========================================
call :color %GREEN% "[SUCCESS] 发布流程结束"
exit /b 0

REM 颜色输出函数
:color
set "color=%~1"
set "text=%~2"
powershell -Command "Write-Host '%text%' -ForegroundColor (Get-Host).UI.RawUI.ForegroundColor -NoNewline; Write-Host ''" 2>nul || echo %text%
goto :eof
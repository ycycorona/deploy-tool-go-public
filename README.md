# Deploy-tool-go

一个golang编写的简单服务器部署工具

A simple server deploy tool wrtite by golang.

## 描述

把部署业务流程分为：本地构建->打压缩包->scp发送压缩包到服务器->服务器解压缩包到部署路径->执行服务器部署脚本

针对上述业务流程，分别开发了4个工具方法:

- zip_file #压缩
- un_zip_file #解压缩
- scp #文件传输
- sh #ssh 远程命令执行

### 使用

#### 使用例子 windows下

```cmd
@echo off

set toolWin=.\deploy\zip-tool-win.exe
set toolLinux=.\deploy\zip-tool-linux
@REM 服务器端工具路径
set remoteToolName=~/zip-tool-linux 
set user=username
set host=x.x.x.x
set passwd=passwd
set localDist=.\dist
set localZipName=.\dist.zip

@REM 修改以下设置，适配自己的项目
set remoteCWD=/path-to-your-program
set remoteDist=/path-to-your-program/dist
set remoteZipName=/path-to-your-program/dist.zip

@REM 前端生成最新部署文件，使用call，防止生成后关闭命令行窗口
call yarn run build || exit /B

@REM 压缩
%toolWin% -act zip ^
  -src %localDist% ^
  -dst %localZipName% || exit /B

@REM 发送工具包
%toolWin% -act scp ^
  -src %toolLinux% ^
  -dst %remoteToolName% ^
  -user %user% ^
  -host %host% ^
  -passwd %passwd% || exit /B

@REM 增加权限+清理目录
%toolWin% -act sh ^
  -cmd "chmod +x %remoteToolName% && cd %remoteCWD% && ls -alh && rm -rf ./dist/*" ^
  -user %user% ^
  -host %host% ^
  -passwd %passwd% || exit /B

@REM 发送部署文件
%toolWin% -act scp ^
  -src %localZipName% ^
  -dst %remoteZipName% ^
  -user %user% ^
  -host %host% ^
  -passwd %passwd% || exit /B

@REM 远端解压
%toolWin% -act sh ^
  -cmd "%remoteToolName% -act unzip -src %remoteZipName% -dst %remoteDist%" ^
  -user %user% ^
  -host %host% ^
  -passwd %passwd% || exit /B

@REM 远端部署脚本
%toolWin% -act sh ^
  -cmd "cd %remoteCWD% && ./easyfix-app-3-hr.sh easyfix-app-3-hr-auto" ^
  -user %user% ^
  -host %host% ^
  -passwd %passwd% || exit /B

PAUSE

```

## 构建
win x64
```bash
#!/bin/bash -e
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -x -o ./dist/zip-tool-win.exe ./src/main.go
```
linux x64
```bash
#!/bin/bash -e
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -x -o ./dist/zip-tool-linux ./src/main.go
```
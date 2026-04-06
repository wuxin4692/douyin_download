---
name: douyin-downloader-go-react
overview: 使用 Go 语言重构抖音下载工具，新增 React 前端界面，支持本地运行。保留原有核心功能：视频解析下载、无水印提取、音频转文字、语义分段。
design:
  architecture:
    framework: react
    component: shadcn
  styleKeywords:
    - 深色主题
    - 简约现代
    - 卡片式布局
  fontSystem:
    fontFamily: Inter
    heading:
      size: 24px
      weight: 600
    subheading:
      size: 18px
      weight: 500
    body:
      size: 14px
      weight: 400
  colorSystem:
    primary:
      - "#6366f1"
      - "#8b5cf6"
    background:
      - "#0f172a"
      - "#1e293b"
    text:
      - "#f8fafc"
      - "#94a3b8"
    functional:
      - "#22c55e"
      - "#ef4444"
      - "#f59e0b"
todos:
  - id: create-go-backend
    content: 创建 Go 后端项目结构，配置 go.mod 和依赖
    status: completed
  - id: implement-http-utility
    content: 实现 HTTP 请求工具函数和配置文件
    status: completed
    dependencies:
      - create-go-backend
  - id: implement-video-service
    content: 实现抖音链接解析服务（parseShareUrl, getVideoInfoByModalId）
    status: completed
    dependencies:
      - implement-http-utility
  - id: implement-download-service
    content: 实现视频下载和音频提取服务
    status: completed
    dependencies:
      - implement-video-service
  - id: implement-transcription-service
    content: 实现语音转写和语义分段服务
    status: completed
    dependencies:
      - implement-download-service
  - id: implement-api-handlers
    content: 实现 HTTP API 处理器和路由
    status: completed
    dependencies:
      - implement-transcription-service
  - id: create-react-frontend
    content: 创建 React 前端项目（Vite + TypeScript + Tailwind）
    status: completed
    dependencies:
      - create-go-backend
  - id: implement-ui-components
    content: 实现前端 UI 组件（输入框、进度条、结果展示）
    status: completed
    dependencies:
      - create-react-frontend
  - id: integrate-frontend-backend
    content: 实现前端 API 服务集成和状态管理
    status: completed
    dependencies:
      - implement-ui-components
      - implement-api-handlers
  - id: add-settings-history
    content: 添加设置管理和历史记录功能
    status: completed
    dependencies:
      - integrate-frontend-backend
---

## 产品概述

将现有的 Node.js 抖音视频下载工具用 Go 语言重写，并添加 React 前端界面，打造一个本地化的抖音视频处理工具。

## 核心功能

1. **视频解析**: 从抖音分享链接获取无水印视频信息（标题、下载链接、视频ID）
2. **视频下载**: 下载无水印视频到本地，支持进度显示
3. **音频提取**: 使用 ffmpeg 从视频提取 MP3 音频
4. **语音转文字**: 调用 Silicon Flow API (FunAudioLLM/SenseVoiceSmall) 将音频转为文字
5. **语义分段**: 调用 MiniMax API 对文本进行智能分段，添加小标题
6. **文案管理**: 保存文案到 Markdown 文件，支持复制到剪贴板
7. **ffmpeg 集成**: 程序启动时自动检测 ffmpeg，不存在则自动下载对应平台的最新版本

## 前端功能

- 链接输入框：支持分享链接、modal_id 多种格式
- 视频信息展示：显示标题、ID、封面
- 下载进度：实时显示下载、转换进度
- 文案展示：显示识别结果和分段结果
- 一键复制：复制文案到剪贴板
- 下载历史：显示最近处理记录（本地存储）
- 设置面板：配置 API 密钥

## 技术栈

- **后端**: Go 1.21+ / Gin (Web 框架) / 标准库 HTTP 客户端
- **前端**: React 18 + TypeScript + Vite + Tailwind CSS
- **音视频处理**: ffmpeg（自动检测，无则自动下载）
- **跨平台**: Windows/macOS/Linux

## 架构设计

采用前后分离架构，Go 后端提供 HTTP API，前端通过 API 调用后端服务：

```
┌─────────────────────────────────────────────────────┐
│                    React 前端                        │
│  (localhost:5173)  ←─── HTTP/REST ───→  (localhost:8080)  │
│                                                    │
│  • 链接输入 & 参数配置                    Go 后端    │
│  • 进度展示 & 结果展示                    • 视频解析   │
│  • 本地存储 (历史记录/设置)               • 视频下载   │
│                                           • 音频提取  │
└───────────────────────────────────────────• 语音转写  │
                                            • 语义分段  │
                                            • 文件管理  │
                                            └──────────┘
```

## 项目结构

```
douyin-tool/
├── backend/                    # Go 后端
│   ├── main.go                 # 程序入口
│   ├── go.mod                  # Go 模块
│   ├── config/                 # 配置
│   │   └── config.go
│   ├── handler/                # HTTP 处理器
│   │   ├── video.go            # 视频相关 API
│   │   ├── task.go             # 任务处理 API
│   │   └── settings.go         # 设置 API
│   ├── service/                # 业务逻辑
│   │   ├── douyin.go           # 抖音解析服务
│   │   ├── downloader.go       # 下载服务
│   │   ├── audio.go            # 音频处理
│   │   ├── transcription.go    # 语音转写
│   │   └── segment.go          # 语义分段
│   └── utils/                  # 工具函数
│       ├── http.go             # HTTP 请求
│       ├── file.go             # 文件操作
│       └── ffmpeg.go           # ffmpeg 自动下载与管理
│
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── components/         # 组件
│   │   │   ├── LinkInput.tsx   # 链接输入
│   │   │   ├── VideoInfo.tsx    # 视频信息
│   │   │   ├── TaskPanel.tsx   # 任务面板
│   │   │   ├── ResultView.tsx  # 结果展示
│   │   │   └── Settings.tsx    # 设置面板
│   │   ├── hooks/              # 自定义 Hooks
│   │   ├── services/           # API 服务
│   │   ├── types/              # 类型定义
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   └── vite.config.ts
│
├── README.md
└── Makefile                    # 构建脚本
```

## 实现要点

### 后端 API 设计

| 方法 | 路径 | 描述 |
| --- | --- | --- |
| POST | /api/video/parse | 解析抖音链接 |
| POST | /api/video/download | 下载视频 |
| POST | /api/audio/extract | 提取音频 |
| POST | /api/task/transcribe | 语音转写 |
| POST | /api/task/segment | 语义分段 |
| GET | /api/settings | 获取设置 |
| POST | /api/settings | 保存设置 |


### 前端状态管理

- 使用 React Hooks 管理状态
- 本地存储保存 API 密钥和历史记录
- SSE (Server-Sent Events) 实现进度推送

### ffmpeg 自动集成

- 程序启动时自动检测系统是否有 ffmpeg
- 如未检测到，自动从 GitHub 下载对应平台的 ffmpeg static build
- Windows: ffmpeg.ico + ffmpeg.exe + ffprobe.exe
- macOS/Linux: 下载对应的 static build 版本
- 下载目录: ~/.douyin-tool/ffmpeg/
- 提供 API: GET /api/ffmpeg/status - 获取 ffmpeg 安装状态

### 错误处理

- 统一的错误响应格式
- 前端友好错误提示
- 超时和重试机制

## 设计风格

采用现代简约风格，以深色主题为主，营造专业工具感。界面分为左侧输入区和右侧结果展示区，清晰直观。

## 布局结构

- **顶部导航栏**: Logo、标题、设置按钮
- **主内容区**: 
- 左侧卡片：链接输入 + 操作按钮
- 右侧卡片：结果展示 + 进度显示
- **底部状态栏**: 下载目录、版本信息

## 设计特点

- 深色主题 (#0f172a 背景)
- 渐变强调色 (#6366f1 → #8b5cf6)
- 圆角卡片布局
- 微动画交互反馈
- 清晰的视觉层次

## Agent Extensions

### douyin-downloader

- **用途**: 参考抖音下载工具的最佳实践和现有实现逻辑
- **预期效果**: 确保 Go 重写版本保持与原版相同的功能完整性和稳定性
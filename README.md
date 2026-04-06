# 抖音工具箱

> ⚠️ **免责声明**：本项目仅供学习交流使用，请勿用于商业目的或任何可能侵犯版权的活动。使用者需自行承担一切责任。

抖音视频下载与文案提取工具，基于 Wails 打包成桌面程序。

## 功能特性

### 核心功能
- 解析抖音视频信息（标题、作者、封面）
- 下载无水印视频
- 提取音频（MP3）
- 语音转文字（Silicon Flow API）

### AI 增强功能
- 语义分段 - 智能识别文案主题和段落结构
- AI 文案分析 - 分析内容结构、情感基调、受众定位、爆款元素
- AI 一键改写 - 支持多种风格（激励人心/幽默风趣/专业严谨/日常随意/情感丰富）
- 自定义改写指令 - 加入自己的想法定制改写结果

### AI 模型支持
| 模式 | 说明 |
|------|------|
| OpenAI | GPT-3.5 Turbo / GPT-4 / GPT-4 Turbo |
| MiniMax | MiniMax-M2 / MiniMax-M2.5 |
| OpenAI 兼容接口 | 支持硅基流动、阿里云通义等第三方 API |

## 技术栈

- **桌面框架**: Wails v2
- **前端**: React + TypeScript + Vite + TailwindCSS
- **后端**: Go

## 环境要求

- Go 1.21+
- Node.js 18+
- FFmpeg（需手动配置）
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

## 安装

### 方式一：下载预编译版本

前往 [Releases](https://github.com/wuxin4692/douyin_download/releases) 下载最新版本。

### 方式二：从源码构建

```bash
# 克隆项目
git clone https://github.com/wuxin4692/douyin_download.git
cd douyin-tool

# 使用构建脚本（Windows）
.\build.bat

# 或手动构建
wails build
```

构建完成后，可执行文件位于 `build/bin/` 目录。

## FFmpeg 配置

程序不包含 FFmpeg，需要用户手动配置：

1. 下载 FFmpeg：https://ffmpeg.org/download.html
2. 打开应用设置，选择 FFmpeg 所在目录
3. 程序会自动检测目录中的 ffmpeg.exe 和 ffprobe.exe

## API Key 配置

在应用设置界面配置以下项：

| 配置项 | 用途 | 必填 |
|--------|------|------|
| Silicon Flow API Key | 语音转文字 | 是 |
| AI 模型 | AI 分析/改写 | 是 |
| API Key | AI 服务密钥 | 是 |
| API 接口地址 | 仅兼容模式需要 | 视情况 |

### 获取 API Key

- **Silicon Flow**: https://cloud.siliconflow.cn/
- **OpenAI**: https://platform.openai.com/
- **MiniMax**: https://www.minimaxi.com/
- **硅基流动**: https://cloud.siliconflow.cn/

## 使用流程

1. 粘贴抖音分享链接或视频 ID
2. 点击解析获取视频信息
3. 下载视频
4. 提取音频
5. 语音转文字
6. 语义分段 / AI 分析 / 一键改写
7. 保存文案到 Markdown

## 项目结构

```
douyin-tool/
├── app.go              # Wails 应用主文件
├── main.go             # 应用入口
├── wails.json          # Wails 配置
├── frontend/            # 前端代码
│   ├── src/
│   │   ├── App.tsx     # 主应用组件
│   │   ├── bindings.ts # Wails 绑定
│   │   ├── services/   # API 服务
│   │   └── wailsjs/    # Wails 生成代码
│   └── index.html
├── backend/             # 后端代码（备用/独立运行）
├── service/             # 核心服务
│   ├── douyin.go       # 抖音 API
│   ├── downloader.go   # 视频下载
│   ├── audio.go        # 音频提取
│   ├── transcription.go # 语音转写
│   ├── segment.go      # 语义分段
│   ├── llm.go          # AI 服务
│   └── file_service.go # 文件操作
├── handler/             # HTTP 处理器
├── config/             # 配置管理
├── utils/              # 工具函数
└── build/              # 构建输出
```

## 开发

### 安装依赖

```bash
# 安装前端依赖
cd frontend
npm install

# 安装后端依赖
go mod tidy
```

### 开发模式

```bash
# 启动前端开发服务器
cd frontend
npm run dev

# 在另一个终端启动 Wails 开发模式
wails dev
```

### 构建发布版本

```bash
wails build
```

## 配置存储

配置文件保存在用户目录：
- Windows: `%USERPROFILE%\.douyin-tool\settings.json`
- macOS: `~/.douyin-tool/settings.json`
- Linux: `~/.douyin-tool/settings.json`

## 常见问题

### Q: 提示 "FFmpeg 未配置"
A: 请在设置中选择 FFmpeg 所在目录，或将 ffmpeg 添加到系统 PATH。

### Q: 语音转写失败
A: 检查 Silicon Flow API Key 是否正确，或账户余额是否充足。

### Q: AI 功能无法使用
A: 检查 AI 模型和 API Key 配置，确保网络可以访问对应 API 服务。

## 许可证

GNU GENERAL PUBLIC LICENSE v3

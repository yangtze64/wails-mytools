# MyTools

MyTools 是一个基于 Wails 3 + Go + Vue 3 构建的桌面工具箱应用，集成了多种日常开发和文本处理工具，所有数据本地存储，无需联网。

[English](README.en.md)

## 功能列表

### 文本工具

| 工具 | 说明 |
|------|------|
| 文本去重 | 按行删除重复内容 |
| 随机字符串 | 生成密码和 Token |
| 文本差异对比 | 查看新增和删除行 |
| JSON 工具 | 格式化、压缩和校验 JSON |
| Markdown 编辑器 | 实时编辑和预览 |
| Base64 编解码 | 文本和文件编解码 |

### 开发工具

| 工具 | 说明 |
|------|------|
| 时间戳转换 | 时间和时间戳互转 |
| Cron 表达式生成器 | 可视化配置定时任务 |
| JWT 解析器 | 解析 JSON Web Token |
| URL 编解码 | URL 编码解码和解析 |
| 浏览器书签 | 保存和打开常用网站 |
| 账号密钥管理 | 本地加密保存密码和 Key |

### 图像工具

| 工具 | 说明 |
|------|------|
| 生成二维码 | 文本和链接转 PNG |
| 解析二维码 | 图片提取二维码内容 |
| 图片压缩 | 调整质量压缩图片 |

### 图形工具

| 工具 | 说明 |
|------|------|
| Mermaid 编辑器 | 图表源码实时预览 |
| 脑图 | 缩进文本生成脑图 |
| Markmap 脑图 | Markdown 生成精美脑图 |

## 快速开始

### 前提条件

- Go 1.25+
- Node.js 18+
- pnpm（推荐通过 Corepack 启用：`corepack enable`）
- Wails 3 CLI（`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`）

### 安装步骤

```sh
git clone https://github.com/your-username/wails-mytools.git
cd wails-mytools
cd frontend && pnpm install && cd ..
task dev
```

也可以直接使用 Wails 命令：

```sh
wails3 dev -config ./build/config.yml
```

## 常用命令

```sh
task dev          # 启动开发模式
task build        # 构建应用
task package      # 打包应用（生成安装包）
```

前端命令在 `frontend/` 下执行：

```sh
pnpm run dev      # 启动前端开发服务器
pnpm run build    # 构建前端生产版本
```

## 目录结构

```
wails-mytools/
├── main.go                    # Wails 应用入口
├── internal/
│   └── services/              # Go 后端服务（导出给前端调用）
├── frontend/
│   ├── src/
│   │   ├── api/               # 前端调用 Wails bindings 的封装
│   │   ├── components/        # 可复用 Vue 组件
│   │   ├── composables/       # Vue 组合式函数
│   │   ├── config/            # 工具注册表等配置
│   │   ├── styles/            # 全局样式
│   │   ├── types/             # 前端共享类型
│   │   └── views/             # 页面级 Vue 组件
│   ├── bindings/              # Wails 自动生成的 TypeScript bindings
│   └── package.json
├── build/
│   ├── config.yml             # Wails 3 构建配置
│   ├── Taskfile.yml           # 通用构建任务
│   ├── darwin/                # macOS 构建资源
│   ├── windows/               # Windows 构建资源
│   └── linux/                 # Linux 构建资源
├── Taskfile.yml               # 项目根 Taskfile
├── .github/
│   └── workflows/             # GitHub Actions 工作流
│       ├── ci.yml             # CI 检查（类型检查、lint）
│       └── release.yml        # 自动构建发布
└── README.md
```

## 开发约定

- 新增 Go 服务后运行 `wails3 generate bindings -clean=true -ts` 更新前端 bindings。
- 前端依赖安装、开发和构建统一使用 pnpm，不使用 npm。
- 页面不要直接 import `frontend/bindings`，优先在 `frontend/src/api/` 做一次封装。
- 工具数量较少时保持简单结构；需要多页面导航、历史记录或深链时再引入 Vue Router。
- 共享状态变复杂时再引入 Pinia，避免过早增加全局状态层。

## 技术栈

| 层级 | 技术 |
|------|------|
| 框架 | Wails 3 |
| 后端 | Go 1.25+ |
| 前端 | Vue 3 + TypeScript |
| UI 组件库 | Element Plus |
| 包管理 | pnpm |
| 构建工具 | Vite 5 |
| 图表 | Mermaid |
| 脑图 | mind-elixir / markmap |
| 二维码 | qrcode / jsqr |

## License

MIT

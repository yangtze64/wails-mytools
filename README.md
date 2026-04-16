# MyTools

MyTools 是一个基于 Wails 3、Go 和 Vue 3 的桌面工具箱项目。前端包管理统一使用 pnpm。

## 快速开始

```sh
cd frontend
pnpm install
cd ..
task dev
```

也可以直接使用 Wails 命令：

```sh
wails3 dev -config ./build/config.yml
```

## 常用命令

```sh
task dev
task build
task package
```

前端命令在 `frontend/` 下执行：

```sh
pnpm run dev
pnpm run build
```

## 目录规划

- `main.go`: Wails 应用入口，负责窗口、资源和服务注册。
- `internal/services/`: Go 后端服务，导出给 Vue 调用。新增工具能力优先按业务服务拆分到这里。
- `frontend/src/api/`: 前端调用 Wails bindings 的轻封装，避免页面直接依赖生成代码。
- `frontend/src/components/`: 可复用 Vue 组件。
- `frontend/src/views/`: 页面级 Vue 组件。
- `frontend/src/styles/`: 全局样式。
- `frontend/src/types/`: 前端共享类型。
- `frontend/bindings/`: Wails 自动生成的 TypeScript bindings，不手动修改。
- `build/`: Wails 3 构建、平台打包和应用元信息配置。

## 开发约定

- 新增 Go 服务后运行 `wails3 generate bindings -clean=true -ts` 更新前端 bindings。
- 前端依赖安装、开发和构建都使用 pnpm，不使用 npm。
- 页面不要直接 import `frontend/bindings`，优先在 `frontend/src/api/` 做一次封装。
- 工具数量较少时保持简单结构；需要多页面导航、历史记录或深链时再引入 Vue Router。
- 共享状态变复杂时再引入 Pinia，避免过早增加全局状态层。

## 当前基础能力

- Go 服务 `AppService` 提供应用信息和健康检查。
- Vue 首页提供工具箱布局和后续工具入口占位。
- `build/Taskfile.yml` 已切换到 pnpm。

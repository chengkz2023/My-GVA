# BoyKing Admin Web

Vue 3 + Vite + Element Plus + Pinia 前端，与 `server/` 配套使用（动态路由 + 后端菜单 + JWT 鉴权）。

## 本地运行

```bash
npm install
npm run dev       # 默认 :8080，/api 代理到 :8888
```

需要后端已启动（登录依赖真实接口，无 Mock 登录）。

## 开发指南

前端开发规范（加页面、request 约定、错误处理、i18n、pinia）见 **`../docs/frontend-guide.md`**。
架构与模块开发见 `../docs/backend-handoff.md` 与 `../server/docs/how-to-add-module.md`。

## 构建

```bash
npm run build
npm run preview
```

镜像构建用仓库内 `web/Dockerfile`（npm 多阶段构建 + Nginx，同源代理见 `../docs/deployment.md`）。

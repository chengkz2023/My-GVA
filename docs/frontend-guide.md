# 前端开发指南

Vue 3 + Vite + Element Plus + Pinia 前端约定。后端菜单驱动动态路由 + JWT 鉴权（无 Mock 登录，`src/core/config.js` 中 `useStaticMenu/useMockLogin` 均为 `false`）。

## 目录结构

```
web/src/
├── api/               # 每个后端模块一个 js 文件，只做请求封装
├── assets/            # 静态资源
├── components/        # 公共组件（warningBar 等）
├── core/              # 启动配置（config.js、error-handel、gin-vue-admin.js）
├── directive/         # 自定义指令（auth 按钮权限、clickOutSide）
├── i18n/              # vue-i18n 语言包（默认 zh-CN）
├── pinia/modules/     # user / router / app 状态
├── router/            # 静态路由（登录、404 等）+ 动态路由装配
├── style/             # main.scss、reset、transition 等
├── utils/             # request.js、bus.js、格式化工具等
└── view/              # layout/ 布局；superAdmin/ 管理页；error/ 错误页
```

## 加一个页面的四步

1. **api 封装**：`web/src/api/<module>.js`，用 `service`（`utils/request.js`）写请求函数，路径不带 `/api` 前缀（`request.js` 已带 baseURL）。
2. **页面组件**：`web/src/view/superAdmin/<module>/<page>.vue`。组件 `name` 与菜单 `name` 对应；文案走 `t()`（见下）。
3. **菜单注册**：后端菜单管理新增菜单（或种子 `seed_data.go` 的 `seedMenus`），`name` + `component: "view/superAdmin/<module>/<page>.vue"`（该字符串由 `import.meta.glob` 解析）；再给角色分配菜单（角色管理 → 设置权限）。
4. **API 权限**：API 管理 → 同步 API 自动发现新路由；给角色勾选新 API 策略（888 超级管理员直通）。

前端无需手写路由表——登录后 `permission.js` 依据后端返回的菜单树生成动态路由。

## 请求封装约定（utils/request.js）

- 统一响应结构 `{code, message, msg, data}`；**逻辑判断用 `code === 0` 与 HTTP 状态码，禁止解析 message 文案**（契约见 `docs/i18n.md`）。
- 401 由拦截器统一处理：清空登录态并跳登录页（页面不要自己实现登出跳转）。
- 页面代码模式：

```js
const res = await getXxxList(params)
if (res.code === 0) {
  // res.data ...
} else {
  ElMessage.error(res.msg || '请求失败')   // 文案仅用于展示
}
```

## 错误处理

- 后端错误已按 Kind 映射 HTTP 状态：400 校验 / 401 未登录 / 403 无权限 / 404 不存在 / 409 冲突 / 429 限流 / 500 内部错误。
- 表单校验：与后端策略保持一致（如密码规则必须 ≥8 位含字母数字，见 `superAdmin/user/user.vue` 的 `passwordPolicyValidator`）。

## i18n 规则（vue-i18n v11）

- 语言包在 `web/src/i18n/index.js`，默认 `zh-CN`。
- **新增 UI 文案一律走 `t()`**（示范页：`view/superAdmin/dictionary/dictionary.vue`）。
- 存量硬编码中文不抽取；有海外交付需求的项目追加语言包并抽取存量文案。
- 动态/带参文案用 `t('key', {name: row.name})`。

## Pinia

- `pinia/modules/user.js`：登录态（token、userInfo、权限）、logout、ClearStorage。
- `pinia/modules/router.js`：动态路由装配状态（SetAsyncRouter / resetRouter）。
- `pinia/modules/app.js`：布局/主题等 UI 状态（drawerSize 等）。
- 新状态按模块建文件，注意登出时清理与路由重置的一致性。

## 常用约定

- 表格页面骨架：`gva-table-box` + `gva-btn-list` + `el-table` + `gva-pagination`（见字典页）。
- 删除等危险操作必须 `ElMessageBox.confirm`。
- 外链跳转：`window.open` 必须 http/https 白名单 + `noopener,noreferrer`（先例：`components/warningBar/warningBar.vue`、`layout/iframe.vue`）。
- 事件用 `e.target`（不要 `e.srcElement`）；组件卸载时清理监听（`onUnmounted`）。
- 大列表 `KeepAlive` 与 tabs 联动遵循 `layout/tabs` 现有实现。

## 相关文档

- 前端与服务端联调问题 → `docs/v1-v2-diff.md`（v1 老习惯对照）
- 评审对照 → `docs/code-review-checklist.md`（前端小节）
- 构建/部署 → `docs/deployment.md`（同源代理）

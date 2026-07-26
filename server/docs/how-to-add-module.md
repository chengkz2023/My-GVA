# 如何新增业务模块

本文档说明在 BoyKing Admin 后端 V2 架构中，如何新增一个标准的业务模块。

---

## 目录

1. [架构概览](#架构概览)
2. [目录结构](#目录结构)
3. [两种模块风格](#两种模块风格)
4. [逐步指南（完整 DDD 风格）](#逐步指南完整-ddd-风格)
5. [依赖边界规则](#依赖边界规则)
6. [模块注册](#模块注册)
7. [代码模板速查](#代码模板速查)
8. [检查清单](#检查清单)

---

## 架构概览

```
请求 → transport/http (Handler, 参数绑定/校验)
      → application (Service, 业务编排)
        → domain (Repository 接口, 实体定义)
        ← infrastructure/mysql (Repository 实现, GORM 映射)
      ← application (DTO 转换)
      ← transport/http (响应序列化)
```

**核心原则：domain 层不依赖任何外部库，infrastructure 层实现 domain 定义的接口。**

```
┌────────────────────────────────────────────┐
│  transport/http/   Gin Handler, 路由注册     │
│  依赖: application, platform/response       │
├────────────────────────────────────────────┤
│  application/      Service, DTO, 业务编排   │
│  依赖: domain, platform/errors              │
├────────────────────────────────────────────┤
│  domain/           实体, Repository 接口    │
│  依赖: 无（或仅 platform/pagination）       │
├────────────────────────────────────────────┤
│  infrastructure/   GORM 仓储实现             │
│  mysql/            依赖: domain, model/     │
├────────────────────────────────────────────┤
│  module.go         依赖注入 & 组装           │
│  依赖: container, 以上四层                  │
└────────────────────────────────────────────┘
```

---

## 目录结构

### 完整 DDD 风格（推荐用于 CRUD 模块）

```
internal/modules/{domain}/{module}/
├── module.go                    # 模块组装 & 注册
├── domain/
│   ├── {entity}.go              # 领域实体 struct
│   └── repository.go            # Repository 接口 + 查询/输入类型
├── application/
│   ├── dto.go                   # 响应 DTO, Command 对象
│   └── service.go               # 业务逻辑编排
├── infrastructure/
│   └── mysql/
│       └── repository.go        # Repository 的 GORM 实现
└── transport/
    └── http/
        └── handler.go           # Gin Handler, 路由注册, 请求体
```

### 轻量风格（用于简单查询或公开接口）

```
internal/modules/{domain}/{module}/
├── module.go                    # 模块组装 & 注册
├── model.go                     # 领域模型
├── dto.go                       # 响应 DTO
├── repository.go                # Repository 接口 + 内存实现
├── service.go                   # 业务逻辑
└── handler.go                   # Gin Handler
```

**适用场景：** 逻辑简单、不依赖数据库、或仅做数据聚合的模块（如 version、config、status）。

---

## 两种模块风格

| 特性 | 完整 DDD | 轻量 |
|------|----------|------|
| 目录层级 | 4 层子目录 | 单层平铺 |
| Repository | 接口在 domain，实现在 infrastructure | 接口 + 实现在同一文件 |
| 数据库访问 | GORM，通过 infra/mysql | 通常无 DB，或直接拿 `*gorm.DB` |
| 适用场景 | CRUD 业务模块 | 公开接口、状态查询 |
| 示例 | user, menu, role, api, file | auth, config, status, version |

**选择建议：** 新业务模块默认选择完整 DDD 风格。只在确实不需要持久化、逻辑少于 3 个方法时考虑轻量风格。

---

## 逐步指南（完整 DDD 风格）

以下以创建一个 **"文章管理 (article)"** 模块为例，一步一步说明。

### 第 1 步：创建目录结构

```
internal/modules/business/article/
├── module.go
├── domain/
│   ├── article.go
│   └── repository.go
├── application/
│   ├── dto.go
│   └── service.go
├── infrastructure/
│   └── mysql/
│       └── repository.go
└── transport/
    └── http/
        └── handler.go
```

### 第 2 步：domain 层 — 定义实体

```go
// internal/modules/business/article/domain/article.go
package domain

import "time"

type Article struct {
    ID        uint
    Title     string
    Content   string
    AuthorID  uint
    Status    string   // draft, published, archived
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**规则：**
- 只放纯数据结构，不放任何 ORM 标签
- 不依赖 `gorm`、`gin`、`global` 等任何外部包
- 字段使用大写导出，与数据库表结构对应

### 第 3 步：domain 层 — 定义 Repository 接口

```go
// internal/modules/business/article/domain/repository.go
package domain

import (
    "context"
    "errors"

    "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

var (
    ErrArticleNotFound      = errors.New("article not found")
    ErrRepositoryUnavailable = errors.New("article repository unavailable")
)

type Repository interface {
    FindByID(ctx context.Context, id uint) (Article, error)
    List(ctx context.Context, query ListQuery) (pagination.Result[Article], error)
    Create(ctx context.Context, input CreateArticleInput) (Article, error)
    Update(ctx context.Context, id uint, input UpdateArticleInput) (Article, error)
    Delete(ctx context.Context, id uint) error
}

type ListQuery struct {
    Page     pagination.Page
    Title    string
    AuthorID uint
    Status   string
}

type CreateArticleInput struct {
    Title    string
    Content  string
    AuthorID uint
    Status   string
}

type UpdateArticleInput struct {
    Title   *string
    Content *string
    Status  *string
}
```

**规则：**
- Repository 接口只使用 domain 自身类型 + platform 基础类型
- 每个仓储方法接收 `context.Context` 作为第一个参数
- 查询输入用独立的 struct（`ListQuery`、`CreateXxxInput`），不直接传 ORM 对象
- 定义领域错误变量（`ErrXxxNotFound`），供上层判断

### 第 4 步：application 层 — 定义 DTO

```go
// internal/modules/business/article/application/dto.go
package application

import "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"

// ---- 查询入参 ----
type ListArticlesQuery struct {
    Page     pagination.Page
    Title    string
    AuthorID uint
    Status   string
}

// ---- 响应类型 ----
type ArticleResponse struct {
    ID        uint   `json:"ID"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    AuthorID  uint   `json:"authorId"`
    Status    string `json:"status"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
}

type ListArticlesResponse = pagination.Result[ArticleResponse]

// ---- 命令对象 ----
type CreateArticleCommand struct {
    Title    string
    Content  string
    Status   string
}

type CreateArticleResponse struct {
    Article ArticleResponse `json:"article"`
}

type UpdateArticleCommand struct {
    Title   *string
    Content *string
    Status  *string
}

type UpdateArticleResponse struct {
    Article ArticleResponse `json:"article"`
}

type DeleteArticleResponse struct {
    DeletedID uint `json:"deletedId"`
}
```

**规则：**
- DTO 带 `json` 标签，控制序列化格式
- `ListXxxResponse` 直接使用 `pagination.Result[T]` 类型别名
- Command 对象用值类型，不依赖 `gin` 的 `context`/`binding` 标签

### 第 5 步：application 层 — 编写 Service

```go
// internal/modules/business/article/application/service.go
package application

import (
    "context"

    "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/domain"
    platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
    apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
    "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

type Service struct {
    repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, query ListArticlesQuery) (ListArticlesResponse, error) {
    if _, ok := platformauth.ActorFromContext(ctx); !ok {
        return ListArticlesResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
    }

    page := pagination.Normalize(query.Page)
    if s.repo == nil {
        return pagination.Result[ArticleResponse]{
            List: []ArticleResponse{}, Total: 0,
            Page: page.Page, PageSize: page.PageSize,
        }, nil
    }

    result, err := s.repo.List(ctx, domain.ListQuery{
        Page: page, Title: query.Title,
        AuthorID: query.AuthorID, Status: query.Status,
    })
    if err == domain.ErrRepositoryUnavailable {
        return pagination.Result[ArticleResponse]{
            List: []ArticleResponse{}, Total: 0,
            Page: page.Page, PageSize: page.PageSize,
        }, nil
    }
    if err != nil {
        return ListArticlesResponse{}, apperrors.New(apperrors.Internal, 0, "list articles failed", err)
    }

    items := make([]ArticleResponse, 0, len(result.List))
    for _, a := range result.List {
        items = append(items, articleToResponse(a))
    }
    return ListArticlesResponse{
        List: items, Total: result.Total,
        Page: result.Page, PageSize: result.PageSize,
    }, nil
}

func (s *Service) Create(ctx context.Context, cmd CreateArticleCommand) (CreateArticleResponse, error) {
    actor, ok := platformauth.ActorFromContext(ctx)
    if !ok {
        return CreateArticleResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
    }
    if cmd.Title == "" {
        return CreateArticleResponse{}, apperrors.WithMessage(apperrors.Validation, "title is required")
    }
    if s.repo == nil {
        return CreateArticleResponse{}, apperrors.WithMessage(apperrors.Internal, "article repository unavailable")
    }

    status := cmd.Status
    if status == "" {
        status = "draft"
    }

    article, err := s.repo.Create(ctx, domain.CreateArticleInput{
        Title: cmd.Title, Content: cmd.Content,
        AuthorID: actor.UserID, Status: status,
    })
    if err == domain.ErrRepositoryUnavailable {
        return CreateArticleResponse{}, apperrors.WithMessage(apperrors.Internal, "article repository unavailable")
    }
    if err != nil {
        return CreateArticleResponse{}, apperrors.New(apperrors.Internal, 0, "create article failed", err)
    }
    return CreateArticleResponse{Article: articleToResponse(article)}, nil
}

// 其余方法（Update, Delete, GetByID）遵循相同模式
// - 从 context 提取 actor
// - 校验入参
// - 检查 repo 是否可用
// - 调用 repo 方法
// - 映射 domain → DTO 返回

func articleToResponse(a domain.Article) ArticleResponse {
    return ArticleResponse{
        ID: a.ID, Title: a.Title, Content: a.Content,
        AuthorID: a.AuthorID, Status: a.Status,
        CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
        UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
    }
}
```

**规则：**
- Service 依赖 `domain.Repository` 接口，不依赖具体实现
- 始终从 `context` 提取当前用户：`platformauth.ActorFromContext(ctx)`
- repo 为 nil 时优雅降级（返回空列表），不 panic
- 用 `apperrors.WithMessage` / `apperrors.New` 创建结构化错误
- 写一个私有 `xxxToResponse(domain.X) XResponse` 映射函数

### 第 6 步：infrastructure 层 — GORM 实现

```go
// internal/modules/business/article/infrastructure/mysql/repository.go
package mysql

import (
    "context"
    "errors"

    "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/domain"
    "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
    "gorm.io/gorm"
)

// ArticleModel 是数据库表的 GORM 映射。
// 如果 model/ 中已有对应的表结构，可直接复用。
type ArticleModel struct {
    ID        uint   `gorm:"primarykey"`
    Title     string `gorm:"column:title;type:varchar(200)"`
    Content   string `gorm:"column:content;type:text"`
    AuthorID  uint   `gorm:"column:author_id"`
    Status    string `gorm:"column:status;type:varchar(20);default:draft"`
    CreatedAt string `gorm:"column:created_at"`
    UpdatedAt string `gorm:"column:updated_at"`
}

func (ArticleModel) TableName() string {
    return "articles"
}

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.Article, error) {
    if r == nil || r.db == nil {
        return domain.Article{}, domain.ErrRepositoryUnavailable
    }
    var m ArticleModel
    err := r.db.WithContext(ctx).First(&m, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return domain.Article{}, domain.ErrArticleNotFound
    }
    if err != nil {
        return domain.Article{}, err
    }
    return domain.Article{
        ID: m.ID, Title: m.Title, Content: m.Content,
        AuthorID: m.AuthorID, Status: m.Status,
    }, nil
}

func (r *Repository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Article], error) {
    page := pagination.Normalize(query.Page)
    if r == nil || r.db == nil {
        return pagination.Result[domain.Article]{
            List: []domain.Article{}, Page: page.Page, PageSize: page.PageSize,
        }, domain.ErrRepositoryUnavailable
    }

    db := r.db.WithContext(ctx).Model(&ArticleModel{})
    if query.Title != "" {
        db = db.Where("title LIKE ?", "%"+query.Title+"%")
    }
    if query.AuthorID != 0 {
        db = db.Where("author_id = ?", query.AuthorID)
    }
    if query.Status != "" {
        db = db.Where("status = ?", query.Status)
    }

    var total int64
    if err := db.Count(&total).Error; err != nil {
        return pagination.Result[domain.Article]{}, err
    }

    var models []ArticleModel
    err := db.Limit(page.Limit()).Offset(page.Offset()).Order("id desc").Find(&models).Error
    if err != nil {
        return pagination.Result[domain.Article]{}, err
    }

    items := make([]domain.Article, 0, len(models))
    for _, m := range models {
        items = append(items, domain.Article{
            ID: m.ID, Title: m.Title, Content: m.Content,
            AuthorID: m.AuthorID, Status: m.Status,
        })
    }
    return pagination.Result[domain.Article]{
        List: items, Total: total, Page: page.Page, PageSize: page.PageSize,
    }, nil
}

func (r *Repository) Create(ctx context.Context, input domain.CreateArticleInput) (domain.Article, error) {
    if r == nil || r.db == nil {
        return domain.Article{}, domain.ErrRepositoryUnavailable
    }
    m := ArticleModel{
        Title: input.Title, Content: input.Content,
        AuthorID: input.AuthorID, Status: input.Status,
    }
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return domain.Article{}, err
    }
    return r.FindByID(ctx, m.ID)
}

// Update, Delete 实现类似，省略…
```

**规则：**
- infrastructure 层是**唯一**可以 import `gorm` 和 `model/` 的模块层
- nil DB 时返回 `domain.ErrRepositoryUnavailable`，不 panic
- 负责 `domain.Xxx` ↔ ORM Model 之间的双向映射
- 分页逻辑：先 Count 再 Limit/Offset/Order

### 第 7 步：transport 层 — Gin Handler

```go
// internal/modules/business/article/transport/http/handler.go
package http

import (
    "strconv"

    "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/application"
    apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
    "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
    "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    service *application.Service
}

func NewHandler(service *application.Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Register(group *gin.RouterGroup) {
    articles := group.Group("/article")
    articles.GET("/list", h.List)
    articles.GET("/:id", h.GetByID)
    articles.POST("", h.Create)
    articles.PUT("/:id", h.Update)
    articles.DELETE("/:id", h.Delete)
}

// ---- 请求体 ----
type listArticlesRequest struct {
    pagination.Page
    Title    string `form:"title"`
    AuthorID uint   `form:"authorId"`
    Status   string `form:"status"`
}

type createArticleRequest struct {
    Title   string `json:"title" binding:"required"`
    Content string `json:"content"`
    Status  string `json:"status"`
}

type updateArticleRequest struct {
    Title   *string `json:"title"`
    Content *string `json:"content"`
    Status  *string `json:"status"`
}

// ---- 处理函数 ----
func (h *Handler) List(c *gin.Context) {
    var req listArticlesRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, err)
        return
    }
    result, err := h.service.List(c.Request.Context(), application.ListArticlesQuery{
        Page: req.Page, Title: req.Title,
        AuthorID: req.AuthorID, Status: req.Status,
    })
    if err != nil {
        response.Error(c, err)
        return
    }
    response.OK(c, result)
}

func (h *Handler) GetByID(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, err)
        return
    }
    // 假设有 GetByID service 方法
    response.OK(c, gin.H{"id": id})
}

func (h *Handler) Create(c *gin.Context) {
    var req createArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
        return
    }
    result, err := h.service.Create(c.Request.Context(), application.CreateArticleCommand{
        Title: req.Title, Content: req.Content, Status: req.Status,
    })
    if err != nil {
        response.Error(c, err)
        return
    }
    response.OK(c, result)
}

func (h *Handler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, err)
        return
    }
    var req updateArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
        return
    }
    // 调用 service.Update(...)
    _ = id
    response.OK(c, nil)
}

func (h *Handler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        response.Error(c, err)
        return
    }
    result, err := h.service.Delete(c.Request.Context(), uint(id))
    if err != nil {
        response.Error(c, err)
        return
    }
    response.OK(c, result)
}
```

**规则：**
- Handler 只做：参数绑定 → 调用 service → 响应序列化
- 请求体 struct 放在 handler 文件内（不导出）
- Query 参数用 `form` 标签 + `ShouldBindQuery`
- JSON body 用 `json` 标签 + `ShouldBindJSON`
- 路径参数用 `c.Param("id")` + `strconv.ParseUint`
- 统一用 `response.OK(c, data)` / `response.Error(c, err)` 返回

### 第 8 步：module.go — 组装 & 注入

```go
// internal/modules/business/article/module.go
package article

import (
    "github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
    v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
    "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/domain"
    "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/application"
    articlemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/infrastructure/mysql"
    articlehttp "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/article/transport/http"
)

type Module struct {
    handler *articlehttp.Handler
}

func NewModule(c *container.Container) *Module {
    var repo domain.Repository
    if c != nil {
        repo = articlemysql.NewRepository(c.DB)
    }
    service := application.NewService(repo)
    return &Module{
        handler: articlehttp.NewHandler(service),
    }
}

func (m *Module) RegisterHTTP(routes v2http.Routes) {
    m.handler.Register(routes.Authenticated)  // 需要登录
}
```

**如果需要公开接口（无需 JWT）：**

```go
func (m *Module) RegisterHTTP(routes v2http.Routes) {
    m.handler.Register(routes.Public)  // 无需认证
}
```

**如果需要同时注册公开和认证路由：** 在 Handler 中提供两个 Register 方法，或直接传两个 group。

### 第 9 步：注册到模块列表

```go
// internal/modules/modules.go
func HTTPModules(c *container.Container) []v2http.Module {
    return []v2http.Module{
        // ... 已有模块 ...
        businessarticle.NewModule(c),  // ← 新增这一行
    }
}
```

---

## 依赖边界规则

这是最重要的部分。每一层**只能**依赖下面列出的包：

```
┌──────────────────────────────────────────────────────┐
│  层               允许 import                         │
├──────────────────────────────────────────────────────┤
│  domain/          无（纯 Go 类型）                     │
│                   platform/pagination（仅分页类型）    │
├──────────────────────────────────────────────────────┤
│  application/     domain（自身模块）                   │
│                   platform/auth（提取当前用户）         │
│                   platform/errors（结构化错误）         │
│                   platform/pagination                │
├──────────────────────────────────────────────────────┤
│  infrastructure/  domain（自身模块，实现接口）          │
│  mysql/           model/*（legacy 数据库模型）         │
│                   platform/pagination                │
│                   gorm.io/gorm                       │
├──────────────────────────────────────────────────────┤
│  transport/http/  application（自身模块）              │
│                   platform/errors                    │
│                   platform/pagination（分页绑定）      │
│                   platform/response（序列化）          │
│                   github.com/gin-gonic/gin           │
├──────────────────────────────────────────────────────┤
│  module.go        container（DI 容器）                 │
│                   v2http（路由注册接口）                │
│                   自身模块的四层子包                    │
└──────────────────────────────────────────────────────┘
```

### 硬性规则

1. **domain/ 不能 import 任何 application/、infrastructure/、transport/ 的包**
2. **application/ 不能 import infrastructure/ 或 transport/ 的包**
3. **infrastructure/ 不能 import application/ 或 transport/ 的包**
4. **transport/ 不能 import infrastructure/ 的包**
5. **任何模块层不能 import `global` 包** — 需要的依赖通过构造函数注入
6. **不能跨模块 import domain 类型** — 每个模块的 domain 是私有的
7. **不能在 application 层直接使用 `*gorm.DB`** — 必须通过 Repository 接口

### 为什么这样设计

- **可测试性：** application 层的 Service 可以通过 mock Repository 接口来单元测试
- **可替换性：** 可以换用 PostgreSQL/Redis 实现 Repository 接口，不改 application 代码
- **可读性：** 打开一个文件就知道它在架构中的角色
- **防腐败：** ORM 标签、Gin binding 标签不会污染领域模型

---

## 模块注册

### 路由前缀

所有 V2 模块的路由自动带 `/v2` 前缀：

```
模块注册:    group.Group("/article").GET("/list", ...)
实际路由:    GET /v2/article/list
```

### 两种路由组

| 路由组 | 中间件 | 用途 |
|--------|--------|------|
| `routes.Public` | 无 | 登录、验证码、健康检查等公开接口 |
| `routes.Authenticated` | JWT + Casbin | 需要登录 + 权限校验的业务接口 |

### 认证模块的特殊处理

如果模块同时需要公开和认证路由（如 auth 模块）：

```go
// handler.go
func (h *Handler) Register(authenticated, public *gin.RouterGroup) {
    public.POST("/login", h.Login)
    public.POST("/base/captcha", h.Captcha)
    authenticated.GET("/system/auth/me", h.Me)
}

// module.go
func (m *Module) RegisterHTTP(routes v2http.Routes) {
    m.handler.Register(routes.Authenticated, routes.Public)
}
```

---

## 代码模板速查

### 错误处理模式

```go
// Service 中统一的错误处理：
if s.repo == nil {
    return XxxResponse{}, apperrors.WithMessage(apperrors.Internal, "xxx repository unavailable")
}
if err == domain.ErrXxxNotFound {
    return XxxResponse{}, apperrors.WithMessage(apperrors.NotFound, "xxx not found")
}
if err == domain.ErrRepositoryUnavailable {
    return XxxResponse{}, apperrors.WithMessage(apperrors.Internal, "xxx repository unavailable")
}
if err != nil {
    return XxxResponse{}, apperrors.New(apperrors.Internal, 0, "do something failed", err)
}
```

### 分页查询模式

```go
// Handler:
var req listRequest
c.ShouldBindQuery(&req)

// Service:
page := pagination.Normalize(query.Page)
result, err := s.repo.List(ctx, domain.ListQuery{Page: page, ...})

// Repository:
page := pagination.Normalize(query.Page)
db.Count(&total)
db.Limit(page.Limit()).Offset(page.Offset()).Order("id desc").Find(&models)
```

### nil-repo 降级模式

```go
// 当 DB 不可用时，优雅返回空列表而不是 panic：
if s.repo == nil {
    return pagination.Result[T]{
        List: []T{}, Total: 0,
        Page: page.Page, PageSize: page.PageSize,
    }, nil
}
```

---

## 检查清单

新增业务模块时，逐项确认：

- [ ] 目录结构符合规范（domain / application / infrastructure / transport）
- [ ] domain 层无外部依赖（无 gorm、gin、global 等 import）
- [ ] Repository 接口定义在 domain 层
- [ ] 领域错误变量定义在 domain 层（`ErrXxxNotFound`、`ErrRepositoryUnavailable`）
- [ ] application Service 依赖 Repository 接口，不依赖具体实现
- [ ] DTO 带 `json` 标签，放在 application 层
- [ ] infrastructure 层实现 domain.Repository 接口
- [ ] infrastructure 层 nil DB 时返回 `ErrRepositoryUnavailable`
- [ ] Handler 只做参数绑定 + 调用 service + 响应
- [ ] 请求体 struct 在 handler 文件内，不导出
- [ ] module.go 通过 `c *container.Container` 完成依赖注入
- [ ] 在 `internal/modules/modules.go` 中注册模块
- [ ] 路由注册到正确的 group（Public vs Authenticated）
- [ ] 模块包可编译：`go build ./internal/modules/business/xxx/...`

package http

import (
	"strconv"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/api/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/api/domain"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *application.Service
	routes  func() gin.RoutesInfo
}

func NewHandler(service *application.Service, routes func() gin.RoutesInfo) *Handler {
	return &Handler{service: service, routes: routes}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	apiGroup := group.Group("/system/api")
	apiGroup.GET("/list", h.List)
	apiGroup.GET("/all", h.GetAll)
	apiGroup.GET("/groups", h.Groups)
	apiGroup.GET("/policies/:authorityId", h.Policies)
	apiGroup.PUT("/policies", h.UpdatePolicies)
	apiGroup.GET("/:id", h.GetByID)
	apiGroup.POST("", h.Create)
	apiGroup.PUT("/:id", h.Update)
	apiGroup.DELETE("/:id", h.Delete)
	apiGroup.POST("/batch-delete", h.BatchDelete)
	apiGroup.POST("/fresh-casbin", h.FreshCasbin)
	apiGroup.GET("/sync", h.Sync)
	apiGroup.POST("/ignore", h.IgnoreApi)
	apiGroup.POST("/batch-sync", h.BatchSync)
}

type listQuery struct {
	Page        int    `form:"page"`
	PageSize    int    `form:"pageSize"`
	Path        string `form:"path"`
	Description string `form:"description"`
	ApiGroup    string `form:"apiGroup"`
	Method      string `form:"method"`
	OrderKey    string `form:"orderKey"`
	Desc        bool   `form:"desc"`
}

func (h *Handler) List(c *gin.Context) {
	var q listQuery
	_ = c.ShouldBindQuery(&q)
	result, err := h.service.List(c.Request.Context(), domain.ListQuery{
		Page:        pagination.Page{Page: q.Page, PageSize: q.PageSize},
		Path:        q.Path,
		Description: q.Description,
		ApiGroup:    q.ApiGroup,
		Method:      q.Method,
		OrderKey:    q.OrderKey,
		Desc:        q.Desc,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) GetAll(c *gin.Context) {
	result, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Groups(c *gin.Context) {
	result, err := h.service.Groups(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Policies(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("authorityId"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Policies(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

type updatePoliciesRequest struct {
	AuthorityID uint                        `json:"authorityId"`
	Policies    []application.PolicyResponse `json:"policies"`
}

func (h *Handler) UpdatePolicies(c *gin.Context) {
	var req updatePoliciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	if err := h.service.UpdatePolicies(c.Request.Context(), req.AuthorityID, req.Policies); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req domain.SaveApiInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Create(c.Request.Context(), req)
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
	var req domain.SaveApiInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	req.ID = uint(id)
	result, err := h.service.Update(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

type batchDeleteRequest struct {
	IDs []int `json:"ids"`
}

func (h *Handler) BatchDelete(c *gin.Context) {
	var req batchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.BatchDelete(c.Request.Context(), req.IDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) FreshCasbin(c *gin.Context) {
	if err := h.service.FreshCasbin(); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Sync(c *gin.Context) {
	routes := make([]domain.RouteInfo, 0)
	for _, r := range h.routes() {
		routes = append(routes, domain.RouteInfo{Path: r.Path, Method: r.Method})
	}
	result, err := h.service.Sync(routes)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) IgnoreApi(c *gin.Context) {
	var req domain.IgnoreApiInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.IgnoreApi(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) BatchSync(c *gin.Context) {
	var req domain.SyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.BatchSync(c.Request.Context(), req); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

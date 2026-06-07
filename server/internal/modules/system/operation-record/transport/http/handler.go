package http

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/domain"
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
	r := group.Group("/system/operation-record")
	r.GET("/list", h.List)
	r.GET("/:id", h.FindByID)
	r.DELETE("/:id", h.Delete)
	r.POST("/batch-delete", h.DeleteByIds)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	result, err := h.service.List(c.Request.Context(), domain.ListQuery{
		Page:   pagination.Page{Page: page, PageSize: pageSize},
		Method: c.Query("method"),
		Path:   c.Query("path"),
		Status: status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.FindByID(c.Request.Context(), uint(id))
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

func (h *Handler) DeleteByIds(c *gin.Context) {
	var req batchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.DeleteByIds(c.Request.Context(), req.IDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

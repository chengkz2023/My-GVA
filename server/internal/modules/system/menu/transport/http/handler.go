package http

import (
	"strconv"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/application"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	menuGroup := group.Group("/system/menu")
	menuGroup.GET("/tree", h.Tree)
	menuGroup.POST("/authority", h.AssignAuthority)
	menuGroup.POST("", h.Save)
	menuGroup.GET("/:id", h.GetByID)
	menuGroup.DELETE("/:id", h.Delete)
}

func (h *Handler) Tree(c *gin.Context) {
	if c.Query("all") == "true" {
		tree, err := h.service.All(c.Request.Context())
		if err != nil {
			response.Error(c, err)
			return
		}
		response.OK(c, tree)
		return
	}
	if authorityIDStr := c.Query("authorityId"); authorityIDStr != "" {
		id, err := strconv.ParseUint(authorityIDStr, 10, 64)
		if err != nil {
			response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid authorityId"))
			return
		}
		tree, err := h.service.TreeForAuthority(c.Request.Context(), uint(id))
		if err != nil {
			response.Error(c, err)
			return
		}
		response.OK(c, tree)
		return
	}
	tree, err := h.service.Tree(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, tree)
}

type assignAuthorityRequest struct {
	AuthorityID uint   `json:"authorityId"`
	MenuIDs     []uint `json:"menuIds"`
}

func (h *Handler) AssignAuthority(c *gin.Context) {
	var req assignAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.OK(c, nil)
		return
	}
	if err := h.service.AssignAuthority(c.Request.Context(), req.AuthorityID, req.MenuIDs); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Save(c *gin.Context) {
	var req application.SaveMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Save(c.Request.Context(), req)
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

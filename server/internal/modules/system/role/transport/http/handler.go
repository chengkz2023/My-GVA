package http

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/domain"
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
	roleGroup := group.Group("/system/role")
	roleGroup.GET("/tree", h.Tree)
	roleGroup.GET("/:id", h.GetByID)
	roleGroup.GET("/:id/data-authorities", h.GetDataAuthorities)
	roleGroup.POST("", h.Create)
	roleGroup.PUT("/:id", h.Update)
	roleGroup.DELETE("/:id", h.Delete)
	roleGroup.POST("/copy", h.Copy)
	roleGroup.POST("/data-authority", h.SetDataAuthority)
}

func (h *Handler) Tree(c *gin.Context) {
	tree, err := h.service.Tree(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, tree)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	role, err := h.service.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	dataIDs, _ := h.service.GetDataAuthorities(c.Request.Context(), uint(id))
	response.OK(c, application.RoleDetailResponse{
		AuthorityID:      role.AuthorityID,
		AuthorityName:    role.AuthorityName,
		ParentID:         role.ParentID,
		DefaultRouter:    role.DefaultRouter,
		DataAuthorityIDs: dataIDs,
	})
}

func (h *Handler) GetDataAuthorities(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	ids, err := h.service.GetDataAuthorities(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	if ids == nil {
		ids = []uint{}
	}
	response.OK(c, gin.H{"dataAuthorityIds": ids})
}

func (h *Handler) Create(c *gin.Context) {
	var req application.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Create(c.Request.Context(), domain.SaveRoleInput{
		AuthorityID:   req.AuthorityID,
		AuthorityName: req.AuthorityName,
		ParentID:      req.ParentID,
		DefaultRouter: req.DefaultRouter,
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
	var req application.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	req.AuthorityID = uint(id)
	result, err := h.service.Update(c.Request.Context(), domain.SaveRoleInput{
		AuthorityID:   uint(id),
		AuthorityName: req.AuthorityName,
		ParentID:      req.ParentID,
		DefaultRouter: req.DefaultRouter,
	})
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

func (h *Handler) Copy(c *gin.Context) {
	var req application.CopyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Copy(c.Request.Context(), req.OldAuthorityID, domain.SaveRoleInput{
		AuthorityID:   req.Authority.AuthorityID,
		AuthorityName: req.Authority.AuthorityName,
		ParentID:      req.Authority.ParentID,
		DefaultRouter: req.Authority.DefaultRouter,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) SetDataAuthority(c *gin.Context) {
	var req application.DataAuthorityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.SetDataAuthority(c.Request.Context(), domain.DataAuthorityInput{
		AuthorityID:       req.AuthorityID,
		DataAuthorityIDs: req.DataAuthorityIDs,
	}); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

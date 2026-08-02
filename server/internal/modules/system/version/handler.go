package version

import (
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	versionGroup := group.Group("/system/version")
	versionGroup.GET("/info", h.Info)
}

func (h *Handler) Info(c *gin.Context) {
	info, err := h.service.Info(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, info)
}

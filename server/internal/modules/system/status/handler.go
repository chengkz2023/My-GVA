package status

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	statusGroup := group.Group("/system/status")
	statusGroup.GET("/info", h.Info)
}

func (h *Handler) Info(c *gin.Context) {
	info, err := h.service.Info(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, info)
}

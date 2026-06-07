package example

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
	exampleGroup := group.Group("/example")
	exampleGroup.GET("/info", h.Info)
	exampleGroup.GET("/missing", h.Missing)
}

func (h *Handler) Info(c *gin.Context) {
	info, err := h.service.Info(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, info)
}

func (h *Handler) Missing(c *gin.Context) {
	_, err := h.service.Missing(c.Request.Context())
	response.Error(c, err)
}

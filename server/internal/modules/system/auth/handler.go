package auth

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

func (h *Handler) Register(authenticated *gin.RouterGroup, public *gin.RouterGroup) {
	authenticated.Group("/system/auth").GET("/me", h.Me)
	public.POST("/login", h.Login)
}

func (h *Handler) Me(c *gin.Context) {
	me, err := h.service.Me(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, me)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

package http

import (
	"strconv"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/application"
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

// Register 注册示例路由。传入的 group 已带 /api 前缀。
func (h *Handler) Register(group *gin.RouterGroup) {
	g := group.Group("/example/greetings")
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
}

type createGreetingRequest struct {
	Message string `json:"message" binding:"required"`
	Author  string `json:"author"`
}

func (h *Handler) List(c *gin.Context) {
	result, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid id"))
		return
	}
	result, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req createGreetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.Create(c.Request.Context(), application.CreateGreetingCommand{
		Message: req.Message,
		Author:  req.Author,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

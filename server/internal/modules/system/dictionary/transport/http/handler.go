package http

import (
	"strconv"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/application"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

// Register 注册字典路由。传入的 group 已带 /api 前缀，路由挂在 Authenticated（需登录）。
func (h *Handler) Register(group *gin.RouterGroup) {
	g := group.Group("/dictionary")
	g.GET("/list", h.List)
	g.GET("/types", h.Types)
	g.POST("", h.Create)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.GET("/details", h.ListDetails)
	g.POST("/details", h.CreateDetail)
	g.PUT("/details/:id", h.UpdateDetail)
	g.DELETE("/details/:id", h.DeleteDetail)
}

type saveDictionaryRequest struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type saveDetailRequest struct {
	DictionaryID uint   `json:"dictionaryId"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Sort         int    `json:"sort"`
	Status       int    `json:"status"`
}

func (h *Handler) List(c *gin.Context) {
	page := pagination.Page{Page: atoiDefault(c.Query("page"), 1), PageSize: atoiDefault(c.Query("pageSize"), 10)}
	result, err := h.service.List(c.Request.Context(), application.ListDictionaryQuery{
		Page: page,
		Type: c.Query("type"),
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Types(c *gin.Context) {
	result, err := h.service.Types(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req saveDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.Create(c.Request.Context(), application.SaveDictionaryCommand{
		Type: req.Type, Name: req.Name, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req saveDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.Update(c.Request.Context(), application.SaveDictionaryCommand{
		ID: id, Type: req.Type, Name: req.Name, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"deletedId": id})
}

func (h *Handler) ListDetails(c *gin.Context) {
	dictionaryID, err := parseUint(c.Query("dictionaryId"))
	if err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid dictionaryId"))
		return
	}
	result, err := h.service.ListDetails(c.Request.Context(), dictionaryID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"list": result})
}

func (h *Handler) CreateDetail(c *gin.Context) {
	var req saveDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.CreateDetail(c.Request.Context(), application.SaveDetailCommand{
		DictionaryID: req.DictionaryID, Label: req.Label, Value: req.Value, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) UpdateDetail(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req saveDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.UpdateDetail(c.Request.Context(), application.SaveDetailCommand{
		ID: id, DictionaryID: req.DictionaryID, Label: req.Label, Value: req.Value, Sort: req.Sort, Status: req.Status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) DeleteDetail(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.service.DeleteDetail(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"deletedId": id})
}

func parseID(c *gin.Context) (uint, error) {
	return parseUint(c.Param("id"))
}

func parseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil || v == 0 {
		return 0, apperrors.WithMessage(apperrors.Validation, "invalid id")
	}
	return uint(v), nil
}

func atoiDefault(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

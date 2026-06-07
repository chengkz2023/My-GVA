package http

import (
	"io"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/domain"
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
	fileGroup := group.Group("/file")
	fileGroup.POST("/upload", h.Upload)
	fileGroup.GET("/list", h.List)
	fileGroup.DELETE("/:id", h.Delete)
	fileGroup.PUT("/:id", h.Update)
}

func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, err)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, err)
		return
	}

	classID := 0
	if v := c.PostForm("classId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			classID = n
		}
	}

	result, err := h.service.Upload(c.Request.Context(), application.FileHeader{
		Filename: header.Filename,
		Data:     data,
	}, classID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	classID, _ := strconv.Atoi(c.DefaultQuery("classId", "0"))

	result, err := h.service.List(c.Request.Context(), domain.ListQuery{
		Page:    pagination.Page{Page: page, PageSize: pageSize},
		ClassID: classID,
		Name:    c.Query("name"),
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

type updateFileRequest struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req updateFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	result, err := h.service.Update(c.Request.Context(), uint(id), req.Name, req.Tag)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

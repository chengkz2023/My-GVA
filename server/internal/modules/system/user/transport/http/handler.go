package http

import (
	"strconv"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/application"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

type ListUsersRequest struct {
	pagination.Page
	Username string `form:"username"`
	NickName string `form:"nickName"`
	Phone    string `form:"phone"`
	Email    string `form:"email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateProfileRequest struct {
	NickName  string `json:"nickName"`
	HeaderImg string `json:"headerImg"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type CreateUserRequest struct {
	Username     string `json:"userName"`
	Password     string `json:"passWord"`
	NickName     string `json:"nickName"`
	HeaderImg    string `json:"headerImg"`
	AuthorityID  uint   `json:"authorityId"`
	AuthorityIDs []uint `json:"authorityIds"`
	Enable       int    `json:"enable"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

type SetAuthoritiesRequest struct {
	AuthorityIDs []uint `json:"authorityIds"`
}

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(group *gin.RouterGroup) {
	userGroup := group.Group("/system/user")
	userGroup.GET("/me", h.Me)
	userGroup.GET("/list", h.List)
	userGroup.POST("/password", h.ChangePassword)
	userGroup.PUT("/profile", h.UpdateProfile)
	userGroup.POST("", h.Create)
	userGroup.DELETE("/:id", h.Delete)
	userGroup.POST("/:id/reset-password", h.ResetPassword)
	userGroup.PUT("/:id/authorities", h.SetAuthorities)
}

func (h *Handler) Me(c *gin.Context) {
	current, err := h.service.Current(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, current)
}

func (h *Handler) List(c *gin.Context) {
	var req ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, err)
		return
	}
	users, err := h.service.List(c.Request.Context(), application.ListUsersQuery{
		Page:     req.Page,
		Username: req.Username,
		NickName: req.NickName,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, users)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.ChangePassword(c.Request.Context(), application.ChangePasswordCommand{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.UpdateProfile(c.Request.Context(), application.UpdateProfileCommand{
		NickName:  req.NickName,
		HeaderImg: req.HeaderImg,
		Phone:     req.Phone,
		Email:     req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.Create(c.Request.Context(), application.CreateUserCommand{
		Username:     req.Username,
		Password:     req.Password,
		NickName:     req.NickName,
		HeaderImg:    req.HeaderImg,
		AuthorityID:  req.AuthorityID,
		AuthorityIDs: req.AuthorityIDs,
		Enable:       req.Enable,
		Phone:        req.Phone,
		Email:        req.Email,
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
	result, err := h.service.Delete(c.Request.Context(), uint(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.ResetPassword(c.Request.Context(), application.ResetPasswordCommand{
		UserID:   uint(id),
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

func (h *Handler) SetAuthorities(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req SetAuthoritiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.WithMessage(apperrors.Validation, "invalid request body"))
		return
	}
	result, err := h.service.SetAuthorities(c.Request.Context(), application.SetAuthoritiesCommand{
		UserID:       uint(id),
		AuthorityIDs: req.AuthorityIDs,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, result)
}

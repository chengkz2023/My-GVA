package auth

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service    *Service
	captchaCfg config.Captcha
	log        *zap.Logger
}

func NewHandler(service *Service, captchaCfg config.Captcha, log *zap.Logger) *Handler {
	return &Handler{service: service, captchaCfg: captchaCfg, log: log}
}

func (h *Handler) Register(authenticated *gin.RouterGroup, public *gin.RouterGroup) {
	authenticated.Group("/system/auth").GET("/me", h.Me)
	public.POST("/login", h.Login)
	public.POST("/base/captcha", h.Captcha)
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
	if !VerifyCaptcha(req.CaptchaId, req.Captcha, h.captchaCfg.OpenCaptcha) {
		response.Fail(c, 400, 7, "验证码错误")
		return
	}
	result, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	platformauth.SetToken(c, result.Token, 7*24*60*60)
	response.OK(c, result)
}

func (h *Handler) Captcha(c *gin.Context) {
	result := GenerateCaptcha(h.captchaCfg, h.log)
	response.OK(c, result)
}

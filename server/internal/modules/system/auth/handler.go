package auth

import (
	"time"

	"github.com/chengkz2023/My-GVA/server/config"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/ratelimit"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service      *Service
	captchaCfg   config.Captcha
	log          *zap.Logger
	loginLimiter *ratelimit.Limiter
	tokenMaxAge  int // 登录 cookie 有效期（秒），读自 jwt.expires-time
}

func NewHandler(service *Service, captchaCfg config.Captcha, log *zap.Logger, limiter *ratelimit.Limiter, tokenMaxAge int) *Handler {
	return &Handler{service: service, captchaCfg: captchaCfg, log: log, loginLimiter: limiter, tokenMaxAge: tokenMaxAge}
}

func (h *Handler) Register(authenticated *gin.RouterGroup, public *gin.RouterGroup) {
	authenticated.Group("/system/auth").GET("/me", h.Me)
	authenticated.Group("/system/auth").POST("/logout", h.Logout)
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
	ipKey := "ip:" + c.ClientIP()
	userKey := "user:" + req.Username
	if h.loginLimiter != nil && (!h.loginLimiter.Allow(ipKey) || !h.loginLimiter.Allow(userKey)) {
		response.Fail(c, 429, 7, "请求过于频繁，请稍后再试")
		return
	}
	if !VerifyCaptcha(req.CaptchaId, req.Captcha, h.captchaCfg.OpenCaptcha) {
		if h.loginLimiter != nil {
			h.loginLimiter.Take(ipKey)
			h.loginLimiter.Take(userKey)
		}
		response.Fail(c, 400, 7, "验证码错误")
		return
	}
	result, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if h.loginLimiter != nil {
			h.loginLimiter.Take(ipKey)
			h.loginLimiter.Take(userKey)
		}
		response.Error(c, err)
		return
	}
	if h.loginLimiter != nil {
		h.loginLimiter.Reset(userKey)
	}
	platformauth.SetToken(c, result.Token, h.tokenMaxAge)
	response.OK(c, result)
}

func (h *Handler) Captcha(c *gin.Context) {
	result := GenerateCaptcha(h.captchaCfg, h.log)
	response.OK(c, result)
}

func (h *Handler) Logout(c *gin.Context) {
	token := platformauth.GetToken(c)
	claims, _ := platformauth.GetClaimsFromContext(c)
	var expiresAt time.Time
	if claims != nil && claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}
	if err := h.service.Logout(c.Request.Context(), token, expiresAt); err != nil {
		h.log.Error("logout failed", zap.Error(err))
	}
	platformauth.ClearToken(c)
	response.OK(c, gin.H{"logout": true})
}

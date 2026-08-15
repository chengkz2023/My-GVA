package middleware

import (
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

const superAdminAuthorityID = 888

// CasbinHandlerWithPrefix 基于注入的 enforcer 做 RBAC 鉴权。
// routerPrefix 用于从请求路径中剥离前缀（如 "/api"），使其与策略中的 path 口径一致。
func CasbinHandlerWithPrefix(enforcer *casbin.SyncedCachedEnforcer, routerPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := platformauth.GetClaimsFromContext(c)
		if !ok {
			response.Fail(c, 403, 7, "permission denied")
			c.Abort()
			return
		}
		if claims.AuthorityId == superAdminAuthorityID {
			c.Next()
			return
		}
		if enforcer == nil {
			response.Fail(c, 403, 7, "permission denied")
			c.Abort()
			return
		}
		obj := strings.TrimPrefix(c.Request.URL.Path, routerPrefix)
		sub := strconv.Itoa(int(claims.AuthorityId))
		success, err := enforcer.Enforce(sub, obj, c.Request.Method)
		if err != nil || !success {
			response.Fail(c, 403, 7, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

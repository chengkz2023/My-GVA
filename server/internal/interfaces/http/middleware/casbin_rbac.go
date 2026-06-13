package middleware

import (
	"strconv"
	"strings"

	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	casbinauthz "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz/casbin"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func CasbinHandler() gin.HandlerFunc {
	return CasbinHandlerWithPrefix("")
}

const superAdminAuthorityID = 888

func CasbinHandlerWithPrefix(routerPrefix string) gin.HandlerFunc {
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
		path := c.Request.URL.Path
		obj := strings.TrimPrefix(path, routerPrefix)
		act := c.Request.Method
		sub := strconv.Itoa(int(claims.AuthorityId))
		e := casbinauthz.GetEnforcer(nil)
		if e == nil {
			response.Fail(c, 403, 7, "permission denied")
			c.Abort()
			return
		}
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			response.Fail(c, 403, 7, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

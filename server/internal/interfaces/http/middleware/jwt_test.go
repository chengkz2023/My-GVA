package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	legacyrequest "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

func TestJWTAuthInjectsActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	global.GVA_CONFIG.JWT.SigningKey = "abcdefgh"
	global.GVA_CONFIG.JWT.ExpiresTime = "1h"
	global.GVA_CONFIG.JWT.BufferTime = "30m"
	global.GVA_CONFIG.JWT.Issuer = "boyking-admin"
	blacklist := local_cache.NewCache(local_cache.SetDefaultExpire(600))

	j := utils.NewJWT()
	token, err := j.CreateToken(j.CreateClaims(legacyrequest.BaseClaims{
		ID:          1,
		Username:    "admin",
		NickName:    "Admin",
		AuthorityId: 888,
	}))
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	engine := gin.New()
	engine.Use(JWTAuthWithConfig(JWTConfig{
		ExpiresTime: "1h",
		BufferTime:  "30m",
		SigningKey:  "abcdefgh",
		BlacklistCheck: func(token string) bool {
			_, ok := blacklist.Get(token)
			return ok
		},
	}))
	engine.GET("/me", func(c *gin.Context) {
		actor, ok := platformauth.ActorFromContext(c.Request.Context())
		if !ok {
			t.Fatal("missing actor")
		}
		if actor.UserID != 1 || actor.AuthorityID != 888 || actor.Username != "admin" || actor.NickName != "Admin" {
			t.Fatalf("actor = %+v", actor)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("x-token", token)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

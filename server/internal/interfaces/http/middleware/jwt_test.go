package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/gin-gonic/gin"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

func TestJWTAuthInjectsActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	blacklist := local_cache.NewCache(local_cache.SetDefaultExpire(600))

	pjwt := platformauth.NewJWT(platformauth.JWTConfig{
		SigningKey:  "abcdefgh",
		ExpiresTime: "1h",
		BufferTime:  "30m",
		Issuer:      "boyking-admin",
	})
	token, err := pjwt.CreateToken(pjwt.CreateClaims(platformauth.BaseClaims{
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

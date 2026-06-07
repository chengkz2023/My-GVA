package auth

import platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"

type Me struct {
	Actor platformauth.Actor
}

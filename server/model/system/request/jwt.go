package request

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
)

// CustomClaims re-exports platform/auth.CustomClaims for backward compatibility.
type CustomClaims = platformauth.CustomClaims

// BaseClaims re-exports platform/auth.BaseClaims for backward compatibility.
type BaseClaims = platformauth.BaseClaims

var _ = jwt.ErrTokenExpired
var _ = uuid.Nil

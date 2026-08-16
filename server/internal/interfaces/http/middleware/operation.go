package middleware

import (
	"net/http"
	"time"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/audit"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationRecord 将写操作（POST/PUT/DELETE/PATCH）交给审计 Sink 记录。
// sink 为审计输出接缝（默认 MySQLSink，等保项目可替换为 syslog/ELK）。
func OperationRecord(sink audit.Sink, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		if sink == nil {
			return
		}
		switch method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return
		}

		userID := 0
		if claims, ok := platformauth.GetClaimsFromContext(c); ok {
			userID = int(claims.BaseClaims.ID)
		}

		entry := audit.Entry{
			IP:      c.ClientIP(),
			Method:  method,
			Path:    path,
			Status:  c.Writer.Status(),
			Latency: time.Since(start),
			Agent:   c.Request.UserAgent(),
			UserID:  userID,
		}
		if len(c.Errors) > 0 {
			entry.ErrorMessage = c.Errors.String()
		}

		if err := sink.Record(c.Request.Context(), entry); err != nil {
			log.Error("operation record create failed", zap.Error(err))
		}
	}
}

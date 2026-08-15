package middleware

import (
	"net/http"
	"time"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperationRecord 记录写操作（POST/PUT/DELETE/PATCH）到 sys_operation_records 审计表。
func OperationRecord(db *gorm.DB, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		if db == nil {
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

		record := platformdb.SysOperationRecord{
			Ip:      c.ClientIP(),
			Method:  method,
			Path:    path,
			Status:  c.Writer.Status(),
			Latency: time.Since(start),
			Agent:   c.Request.UserAgent(),
			UserID:  userID,
		}
		if len(c.Errors) > 0 {
			record.ErrorMessage = c.Errors.String()
		}

		if err := db.Create(&record).Error; err != nil {
			log.Error("operation record create failed", zap.Error(err))
		}
	}
}

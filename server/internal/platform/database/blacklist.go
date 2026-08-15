package database

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

// JwtBlacklist 记录已吊销的 JWT（只存 sha256 哈希），用于登出/改密后作废旧 token。
type JwtBlacklist struct {
	GVA_MODEL
	JwtHash   string    `json:"jwtHash" gorm:"column:jwt_hash;type:char(64);uniqueIndex;comment:jwt sha256"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"column:expires_at;index;comment:过期时间"`
}

func (JwtBlacklist) TableName() string {
	return "jwt_blacklists"
}

func jwtHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// AddToJwtBlacklist 将 token 的哈希写入黑名单。db 或 token 为空时静默跳过。
func AddToJwtBlacklist(db *gorm.DB, token string, expiresAt time.Time) error {
	if db == nil || token == "" {
		return nil
	}
	return db.Create(&JwtBlacklist{JwtHash: jwtHash(token), ExpiresAt: expiresAt}).Error
}

// IsJwtBlacklisted 判断 token 是否已进入黑名单。
func IsJwtBlacklisted(db *gorm.DB, token string) bool {
	if db == nil || token == "" {
		return false
	}
	var count int64
	if err := db.Model(&JwtBlacklist{}).Where("jwt_hash = ?", jwtHash(token)).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

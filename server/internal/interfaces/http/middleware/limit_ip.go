package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type LimitConfig struct {
	GenerationKey func(c *gin.Context) string
	CheckOrMark   func(key string, expire int, limit int) error
	Expire        int
	Limit         int
}

func (l LimitConfig) LimitWithTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := l.CheckOrMark(l.GenerationKey(c), l.Expire, l.Limit); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": response.Failure, "msg": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func DefaultGenerationKey(c *gin.Context) string {
	return "GVA_Limit" + c.ClientIP()
}

func NewLimiter(rdb redis.UniversalClient, log *zap.Logger, expire, limit int) gin.HandlerFunc {
	return LimitConfig{
		GenerationKey: DefaultGenerationKey,
		CheckOrMark:   checkOrMark(rdb, log),
		Expire:        expire,
		Limit:         limit,
	}.LimitWithTime()
}

func checkOrMark(rdb redis.UniversalClient, log *zap.Logger) func(key string, expire int, limit int) error {
	return func(key string, expire int, limit int) error {
		if rdb == nil {
			return nil
		}
		if err := setLimitWithTime(rdb, key, limit, time.Duration(expire)*time.Second); err != nil {
			log.Error("limit", zap.Error(err))
			return err
		}
		return nil
	}
}

func setLimitWithTime(rdb redis.UniversalClient, key string, limit int, expiration time.Duration) error {
	count, err := rdb.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}
	if count == 0 {
		pipe := rdb.TxPipeline()
		pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, expiration)
		_, err = pipe.Exec(context.Background())
		return err
	}
	times, err := rdb.Get(context.Background(), key).Int()
	if err != nil {
		return err
	}
	if times >= limit {
		if t, err := rdb.PTTL(context.Background(), key).Result(); err != nil {
			return errors.New("请求太过频繁，请稍后再试")
		} else {
			return errors.New("请求太过频繁, 请 " + t.String() + " 秒后尝试")
		}
	}
	return rdb.Incr(context.Background(), key).Err()
}

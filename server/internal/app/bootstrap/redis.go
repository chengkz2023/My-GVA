package bootstrap

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func initRedisClient(redisCfg config.Redis, log *zap.Logger) (redis.UniversalClient, error) {
	var client redis.UniversalClient
	if redisCfg.UseCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    redisCfg.ClusterAddrs,
			Password: redisCfg.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     redisCfg.Addr,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		})
	}
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Error("redis connect ping failed, err:", zap.String("name", redisCfg.Name), zap.Error(err))
		return nil, err
	}

	log.Info("redis connect ping response:", zap.String("name", redisCfg.Name), zap.String("pong", pong))
	return client, nil
}

func InitRedis(cfg config.Redis, log *zap.Logger) redis.UniversalClient {
	redisClient, err := initRedisClient(cfg, log)
	if err != nil {
		panic(err)
	}
	return redisClient
}

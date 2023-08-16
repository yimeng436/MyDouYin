package rediscli

import (
	"fmt"
	redis "github.com/redis/go-redis/v9"
	"sync"
	"time"
	"usersvr/config"
	"usersvr/log"
)

var (
	redisConn   *redis.Client
	redisOnce   sync.Once
	ValueExpirt = time.Hour * 24 * 7
)

func initRedis() {
	redisConfig := config.GetGlobalConfig().RedisConfig
	log.Infof("redisConfig=======%+v", redisConfig)
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	opt := &redis.Options{
		Addr:     addr,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	}
	redisConn = redis.NewClient(opt)

	if redisConn == nil {
		panic("failed to call rediscli.NewClient")
	}
}

func CloseRedis() {
	if redisConn != nil {
		redisConn.Close()
	}
}

// GetRedisCli 获取数据库连接
func GetRedisCli() *redis.Client {
	redisOnce.Do(initRedis)

	return redisConn
}

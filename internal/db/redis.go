package db

import (
	"context"
	"custom_insurance/configs"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

func InitRedis(redisConfig configs.RedisConfig) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
		Protocol: 3, // specify 2 for RESP 2 or 3 for RESP 3
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}
}

func GetRdb() *redis.Client {
	return rdb
}

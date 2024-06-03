package redis

import (
	"github.com/go-redis/redis/v8"
)

func ConnectRedis(redisAddr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	return rdb
}

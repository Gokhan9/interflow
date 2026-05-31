package cache

import (
	"context"

	"github.com/redis/go-redis/v8"
)

var RedisClient *redis.Client

func InitRedis(url string) error {
	options, err := redis.ParseURL("redis://" + url)
	if err != nil {
		options = &redis.Options{
			Addr: url,
		}
	}

	RedisClient = redis.NewClient(options)

	return RedisClient.Ping(context.Background()).Err()
}

package cache

import (
	"context"

	"github.com/redis/go-redis/v8"
)

var RedisClient *redis.Client

func InitRedis(url string) error {
	options, err := redis.ParseURL("redis://" + url)
	if err != nil {
		// Eğer URL parse edilemezse direkt adresi deneyelim
		options = &redis.Options{
			Addr: url,
		}
	}

	RedisClient = redis.NewClient(options)

	// Bağlantı Testi için ping atıyoruz.
	return RedisClient.Ping(context.Background()).Err()
}

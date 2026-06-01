package cache

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(url string) error {

	if !strings.HasPrefix(url, "redis://") {
		url = "redis://" + url
	}

	options, err := redis.ParseURL(url)
	if err != nil {
		return err
	}

	RDB = redis.NewClient(options)

	// Bağlantı Testi için ping atıyoruz.
	return RDB.Ping(context.Background()).Err()
}

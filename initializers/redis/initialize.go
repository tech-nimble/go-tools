package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/envy"
)

func Initialize() *redis.Client {
	addr := envy.Get("REDIS_ADDR", "localhost:6379")
	pass := envy.Get("REDIS_PASS", "")

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})
}

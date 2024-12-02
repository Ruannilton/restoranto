package infrastructure

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(ctx context.Context) RedisCache {

	addr := os.Getenv("REDIS_URL")
	pass := os.Getenv("REDIS_PASSWORD")

	options := redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	}

	instance := RedisCache{
		client: redis.NewClient(&options),
		ctx:    ctx,
	}

	if err := instance.client.Ping(instance.ctx).Err(); err != nil {
		fmt.Println("Failed to connect to redis:", err)
	}

	return instance
}

func (cache RedisCache) Set(key string, value string, duration time.Duration) error {
	status := cache.client.Set(cache.ctx, key, value, duration)

	redisErr := status.Err()

	if redisErr != nil {
		// TODO: handle error
		return redisErr
	}

	return nil
}

func (cache RedisCache) Get(key string) (string, error) {
	data, redisErr := cache.client.Get(cache.ctx, key).Result()

	if redisErr == redis.Nil {
		// TODO: handle error
		return "", redisErr
	}

	if redisErr != nil {
		// TODO: handle error
		return "", redisErr
	}

	return data, nil
}

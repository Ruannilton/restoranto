package repositories

import (
	"context"
	"costumers-api/domain/domain_errors"
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

	addr := os.Getenv("REDIS_ADDR")
	port := os.Getenv("REDIS_PORT")
	pass := os.Getenv("REDIS_PASSWORD")

	options := redis.Options{
		Addr:     addr + ":" + port,
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
		err := domain_errors.NewInfrastructureError(DistributedCacheServiceName, domain_errors.CodeInvalidOperation, "failed to add entry to cache", redisErr)
		return err
	}

	return nil
}

func (cache RedisCache) Get(key string) (string, error) {
	data, redisErr := cache.client.Get(cache.ctx, key).Result()

	if redisErr == redis.Nil {
		err := domain_errors.NewInfrastructureError(DistributedCacheServiceName, domain_errors.CodeEntryNotFound, "failed to get entry from cache", redisErr)
		return "", err
	}

	if redisErr != nil {
		err := domain_errors.NewInfrastructureError(DistributedCacheServiceName, domain_errors.CodeInvalidOperation, "failed to get entry from cache", redisErr)
		return "", err
	}

	return data, nil
}

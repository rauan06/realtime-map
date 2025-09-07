package cache

import (
	"context"
	"time"

	"github.com/rauan06/realtime-map/producer/config"
	"github.com/rauan06/realtime-map/producer/internal/repo"
	"github.com/redis/go-redis/v9"
)

var _ (repo.ICache) = &Cache{}

type Cache struct {
	client *redis.Client
}

func New(cfg config.Config) (repo.ICache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.RedisURI,
		Password: cfg.Redis.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Cache{client}, nil
}

func (r *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Result()
	bytes := []byte(res)
	return bytes, err
}

func (r *Cache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Cache) Close() error {
	return r.client.Close()
}

package cache

import (
	"context"
	"time"

	"github.com/rauan06/realtime-map/analytics/config"
	"github.com/rauan06/realtime-map/analytics/internal/repo"
	"github.com/redis/go-redis/v9"
)

var _ (repo.ICache) = &Cache{}

type Cache struct {
	Client *redis.Client
}

func New(client *redis.Client, cfg config.Config) (repo.ICache, error) {
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &Cache{client}, nil
}

func (r *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.Client.Get(ctx, key).Result()
	bytes := []byte(res)
	return bytes, err
}

func (r *Cache) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Cache) Close() error {
	return r.Client.Close()
}

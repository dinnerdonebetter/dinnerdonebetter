package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prixfixeco/backend/internal/cache"

	"github.com/go-redis/redis/v8"
)

type redisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type redisCacheImpl[T cache.Cacheable] struct {
	client     redisClient
	expiration time.Duration
}

// NewRedisCache builds a new redis-backed cache.
func NewRedisCache[T cache.Cacheable](cfg *Config, expiration time.Duration) cache.Cache[T] {
	return &redisCacheImpl[T]{
		client:     buildRedisClient(cfg),
		expiration: expiration,
	}
}

func (i *redisCacheImpl[T]) Get(ctx context.Context, key string) (*T, error) {
	res, err := i.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var x *T
	if unmarshalErr := json.Unmarshal([]byte(res), &x); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	if x == nil {
		return nil, cache.ErrNotFound
	}

	return x, nil
}

func (i *redisCacheImpl[T]) Set(ctx context.Context, key string, value *T) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if setErr := i.client.Set(ctx, key, content, i.expiration).Err(); setErr != nil {
		return setErr
	}

	return nil
}

func (i *redisCacheImpl[T]) Delete(ctx context.Context, key string) error {
	return i.client.Del(ctx, key).Err()
}

// buildRedisClient returns a PublisherProvider for a given address.
func buildRedisClient(cfg *Config) redisClient {
	var c redisClient
	if len(cfg.QueueAddresses) > 1 {
		c = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.QueueAddresses,
			Username:     cfg.Username,
			Password:     cfg.Password,
			DialTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		})
	} else if len(cfg.QueueAddresses) == 1 {
		c = redis.NewClient(&redis.Options{
			Addr:         cfg.QueueAddresses[0],
			Username:     cfg.Username,
			Password:     cfg.Password,
			DialTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		})
	}

	return c
}

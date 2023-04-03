package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prixfixeco/backend/internal/cache"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"

	"github.com/go-redis/redis/v8"
)

type redisCacheImpl[T cache.Cacheable] struct {
	client     *redis.Client
	expiration time.Duration
}

func newRedisCache[T cache.Cacheable](client *redis.Client, expiration time.Duration) cache.Cache[T] {
	return &redisCacheImpl[T]{
		client:     client,
		expiration: expiration,
	}
}

func (i *redisCacheImpl[T]) Get(ctx context.Context, key string) (*T, error) {
	res, err := i.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var x *T
	if err = json.Unmarshal([]byte(res), &x); err != nil {
		return nil, err
	}

	return x, cache.ErrNotFound
}

func (i *redisCacheImpl[T]) Set(ctx context.Context, key string, value *T) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = i.client.Set(ctx, key, content, i.expiration); err != nil {
		return err
	}

	return nil
}

func (i *redisCacheImpl[T]) Delete(ctx context.Context, key string) error {
	return i.client.Del(ctx, key).Err()
}

// Config configures a Redis-backed consumer.
type Config struct {
	Username       string   `json:"username" mapstructure:"username" toml:"username,omitempty"`
	Password       string   `json:"password,omitempty" mapstructure:"password" toml:"password,omitempty"`
	QueueAddresses []string `json:"queueAddress" mapstructure:"queue_address" toml:"queue_address,omitempty"`
	DB             int      `json:"database,omitempty" mapstructure:"database" toml:"database,omitempty"`
}

// buildRedisClient returns a PublisherProvider for a given address.
func buildRedisClient(logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config) *redis.Client {
	logger.WithValue("queue_addresses", cfg.QueueAddresses).
		WithValue("username", cfg.Username).
		WithValue("password", cfg.Password).Info("setting up redis client")

	var redisClient *redis.Client
	if len(cfg.QueueAddresses) > 1 {
		//redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		//	Addrs:        cfg.QueueAddresses,
		//	Username:     cfg.Username,
		//	Password:     cfg.Password,
		//	DialTimeout:  1 * time.Second,
		//	WriteTimeout: 1 * time.Second,
		//})
	} else if len(cfg.QueueAddresses) == 1 {
		redisClient = redis.NewClient(&redis.Options{
			Addr:         cfg.QueueAddresses[0],
			Username:     cfg.Username,
			Password:     cfg.Password,
			DialTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		})
	}

	return redisClient
}

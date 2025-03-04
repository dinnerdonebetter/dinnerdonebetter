package redis

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rediscontainers "github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	exampleKey = "example"
)

type example struct {
	Name string `json:"name"`
}

func buildContainerBackedRedisConfig(t *testing.T, ctx context.Context) (config *Config, shutdownFunction func(context.Context) error) {
	t.Helper()

	redisContainer, err := rediscontainers.Run(ctx,
		"docker.io/redis:7-bullseye",
		rediscontainers.WithLogLevel(rediscontainers.LogLevelNotice),
	)
	if err != nil {
		panic(err)
	}

	redisAddress, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	cfg := &Config{
		QueueAddresses: []string{
			strings.TrimPrefix(redisAddress, "redis://"),
		},
	}

	return cfg, redisContainer.Terminate
}

func Test_redisCacheImpl_Get(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		cfg, containerShutdown := buildContainerBackedRedisConfig(t, ctx)
		defer func() {
			assert.NoError(t, containerShutdown(ctx))
		}()
		c := NewRedisCache[example](cfg, 0)

		exampleContent := &example{Name: t.Name()}
		assert.NoError(t, c.Set(ctx, exampleKey, exampleContent))

		actual, err := c.Get(ctx, exampleKey)
		assert.Equal(t, exampleContent, actual)
		assert.NoError(t, err)
	})
}

func Test_redisCacheImpl_Set(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		cfg, containerShutdown := buildContainerBackedRedisConfig(t, ctx)
		defer func() {
			assert.NoError(t, containerShutdown(ctx))
		}()
		c := NewRedisCache[example](cfg, 0)

		exampleContent := &example{Name: t.Name()}
		assert.NoError(t, c.Set(ctx, exampleKey, exampleContent))
	})
}

func Test_redisCacheImpl_Delete(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()

		cfg, containerShutdown := buildContainerBackedRedisConfig(t, ctx)
		defer func() {
			assert.NoError(t, containerShutdown(ctx))
		}()
		c := NewRedisCache[example](cfg, 0)

		exampleContent := &example{Name: t.Name()}
		assert.NoError(t, c.Set(ctx, exampleKey, exampleContent))

		assert.NoError(t, c.Delete(ctx, exampleKey))
	})
}

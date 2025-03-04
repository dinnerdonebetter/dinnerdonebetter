package redis

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	rediscontainers "github.com/testcontainers/testcontainers-go/modules/redis"
)

func buildContainerBackedRedisConfig(t *testing.T, ctx context.Context) (config *Config, shutdownFunction func(context.Context) error) {
	t.Helper()

	redisContainer, err := rediscontainers.Run(
		ctx,
		"redis:7-bullseye",
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

func TestConfig_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		cfg := &Config{
			Username:       t.Name(),
			Password:       t.Name(),
			QueueAddresses: []string{t.Name()},
		}

		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}

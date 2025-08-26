package redis

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	rediscontainers "github.com/testcontainers/testcontainers-go/modules/redis"
	"strings"
	"testing"
)

const (
	redisProtocolPrefix      = "redis://"
	redisContainerImageToUse = "redis:7-bullseye"
)

func BuildContainerBackedRedisConfigForTest(t *testing.T) (config *Config, shutdownFunc func(context.Context) error) {
	t.Helper()

	cfg, sdf, err := BuildContainerBackedRedisConfig(t.Context())
	require.NoError(t, err)

	return cfg, sdf
}

func BuildContainerBackedRedisConfig(ctx context.Context) (config *Config, shutdownFunc func(context.Context) error, err error) {
	redisContainer, err := rediscontainers.Run(
		ctx,
		redisContainerImageToUse,
		rediscontainers.WithLogLevel(rediscontainers.LogLevelNotice),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build redis container: %w", err)
	}

	redisAddress, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build redis connection string: %w", err)
	}

	cfg := &Config{
		QueueAddresses: []string{strings.TrimPrefix(redisAddress, redisProtocolPrefix)},
	}

	return cfg, redisContainer.Terminate, nil
}

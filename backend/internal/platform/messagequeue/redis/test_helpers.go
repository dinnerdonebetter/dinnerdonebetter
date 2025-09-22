package redis

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	rediscontainers "github.com/testcontainers/testcontainers-go/modules/redis"
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

	// Wait a small amount to ensure container is fully ready
	time.Sleep(100 * time.Millisecond)

	redisAddress, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		redisContainer.Terminate(ctx) // cleanup on error
		return nil, nil, fmt.Errorf("failed to build redis connection string: %w", err)
	}

	cfg := &Config{
		QueueAddresses: []string{strings.TrimPrefix(redisAddress, redisProtocolPrefix)},
	}

	shutdownFunction := func(shutdownCtx context.Context) error {
		// Use a reasonable timeout for shutdown
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer timeoutCancel()
		return redisContainer.Terminate(timeoutCtx)
	}

	return cfg, shutdownFunction, nil
}

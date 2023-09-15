package redis

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	exampleKey = "example"
)

type mockRedisClient struct {
	mock.Mock
}

func (m *mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.Called(ctx, key).Get(0).(*redis.StringCmd)
}

func (m *mockRedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return m.Called(ctx, key, value, expiration).Get(0).(*redis.StatusCmd)
}

func (m *mockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return m.Called(ctx, keys).Get(0).(*redis.IntCmd)
}

func Test_newRedisCache(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		c := NewRedisCache[types.SessionContextData](cfg, 0)

		assert.NotNil(t, c)
	})
}

// TODO: use testcontainers to properly test this: https://golang.testcontainers.org/modules/redis/

func Test_redisCacheImpl_Get(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}
		c := NewRedisCache[types.SessionContextData](cfg, 0)

		mockClient := &mockRedisClient{}
		mockClient.On("Get", testutils.ContextMatcher, exampleKey).Return(redis.NewStringResult("{}", nil))

		c.(*redisCacheImpl[types.SessionContextData]).client = mockClient

		actual, err := c.Get(ctx, exampleKey)
		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}

func Test_redisCacheImpl_Set(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}
		exampleDuration := time.Hour

		c := NewRedisCache[types.SessionContextData](cfg, 0)

		rawInput := &types.SessionContextData{}
		input, err := json.Marshal(rawInput)
		require.NoError(t, err)
		require.NotNil(t, input)

		mockClient := &mockRedisClient{}
		mockClient.On("Set", testutils.ContextMatcher, exampleKey, input, exampleDuration).Return(redis.NewStatusResult("{}", nil))

		c.(*redisCacheImpl[types.SessionContextData]).client = mockClient
		c.(*redisCacheImpl[types.SessionContextData]).expiration = exampleDuration

		assert.NoError(t, c.Set(ctx, exampleKey, rawInput))
	})
}

func Test_redisCacheImpl_Delete(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}
		c := NewRedisCache[types.SessionContextData](cfg, 0)

		mockClient := &mockRedisClient{}
		mockClient.On("Del", testutils.ContextMatcher, []string{exampleKey}).Return(redis.NewIntResult(1, nil))

		c.(*redisCacheImpl[types.SessionContextData]).client = mockClient

		err := c.Delete(ctx, exampleKey)
		assert.NoError(t, err)
	})
}

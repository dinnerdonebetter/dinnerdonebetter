package elasticsearch

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	elasticsearchcontainers "github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

var (
	runningContainerTests = strings.ToLower(os.Getenv("RUN_CONTAINER_TESTS")) == "true"
)

func TestConfig_provideElasticsearchClient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}

		esClient, err := provideElasticsearchClient(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, esClient)
	})
}

func buildContainerBackedElasticsearchConfig(t *testing.T, ctx context.Context) (config *Config, shutdownFunction func(context.Context) error) {
	t.Helper()

	elasticsearchContainer, err := elasticsearchcontainers.Run(
		ctx,
		"elasticsearch:8.10.2",
		elasticsearchcontainers.WithPassword("arbitraryPassword"),
	)
	require.NoError(t, err)
	require.NotNil(t, elasticsearchContainer)

	cfg := &Config{
		Address:               elasticsearchContainer.Settings.Address,
		IndexOperationTimeout: 0,
		Username:              "elastic",
		Password:              elasticsearchContainer.Settings.Password,
		CACert:                elasticsearchContainer.Settings.CACert,
	}

	return cfg, elasticsearchContainer.Terminate
}

func Test_ProvideIndexManager(T *testing.T) {
	T.Parallel()

	if !runningContainerTests {
		T.SkipNow()
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg, shutdownFunc := buildContainerBackedElasticsearchConfig(t, ctx)
		defer func() {
			require.NoError(t, shutdownFunc(ctx))
		}()

		im, err := ProvideIndexManager[example](ctx, nil, nil, cfg, t.Name(), circuitbreaking.NewNoopCircuitBreaker())
		assert.NoError(t, err)
		assert.NotNil(t, im)
	})

	T.Run("without available instance", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}

		im, err := ProvideIndexManager[example](ctx, nil, nil, cfg, t.Name(), circuitbreaking.NewNoopCircuitBreaker())
		assert.Error(t, err)
		assert.Nil(t, im)
	})
}

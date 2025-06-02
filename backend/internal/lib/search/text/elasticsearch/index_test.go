package elasticsearch

import (
	"context"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type example struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Test_indexManager_CompleteLifecycle(T *testing.T) {
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

		im, err := ProvideIndexManager[example](ctx, nil, nil, cfg, "index_test", circuitbreaking.NewNoopCircuitBreaker())
		assert.NoError(t, err)
		assert.NotNil(t, im)

		searchable := &example{
			ID:   identifiers.New(),
			Name: t.Name(),
		}

		assert.NoError(t, im.Index(ctx, searchable.ID, searchable))

		time.Sleep(5 * time.Second)

		results, err := im.Search(ctx, searchable.Name[0:2])
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, searchable, results[0])

		assert.NoError(t, im.Delete(ctx, searchable.ID))
	})
}

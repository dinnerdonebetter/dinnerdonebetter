package algolia

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/circuitbreaking"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestProvideIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()

		im, err := ProvideIndexManager[types.UserSearchSubset](ctx, logger, tracerProvider, &Config{}, "test", circuitbreaking.NewNoopCircuitBreaker())
		assert.NoError(t, err)
		assert.NotNil(t, im)
	})
}

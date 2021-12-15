package jaeger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

func Test_tracingErrorHandler_Handle(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		errorHandler{logger: logging.NewNoopLogger()}.Handle(errors.New("blah"))
	})
}

func TestConfig_SetupJaeger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			CollectorEndpoint:         "blah blah blah",
			ServiceName:               t.Name(),
			SpanCollectionProbability: 0,
		}

		actual, err := SetupJaeger(ctx, cfg)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

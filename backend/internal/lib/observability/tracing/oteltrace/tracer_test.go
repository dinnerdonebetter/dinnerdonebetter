package oteltrace

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"

	"github.com/stretchr/testify/assert"
)

func Test_tracingErrorHandler_Handle(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		errorHandler{logger: logging.NewNoopLogger()}.Handle(errors.New("blah"))
	})
}

func TestConfig_SetupOtelHTTP(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{
			CollectorEndpoint: "blah blah blah",
		}

		actual, err := SetupOtelGRPC(ctx, t.Name(), 0, cfg)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

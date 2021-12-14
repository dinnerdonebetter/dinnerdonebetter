package jaeger

import (
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

		cfg := &Config{
			CollectorEndpoint:         "blah blah blah",
			ServiceName:               t.Name(),
			SpanCollectionProbability: 0,
		}

		actual, flush, err := SetupJaeger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.NotNil(t, flush)
	})

	T.Run("with empty collector endpoint", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			CollectorEndpoint:         "",
			ServiceName:               t.Name(),
			SpanCollectionProbability: 0,
		}

		actual, flush, err := SetupJaeger(cfg)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Nil(t, flush)
	})
}

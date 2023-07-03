package wasm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

func buildTestService() *Service {
	return &Service{
		logger:         logging.NewNoopLogger(),
		encoderDecoder: mockencoding.NewMockEncoderDecoder(),
		tracer:         tracing.NewTracerForTest("test"),
		cfg:            &Config{},
	}
}

func TestProvideValidIngredientsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := &Config{}

		s, err := ProvideService(
			ctx,
			logger,
			cfg,
			mockencoding.NewMockEncoderDecoder(),
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})
}

package oauth2

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/stretchr/testify/assert"
)

func buildTestService() *Service {
	return &Service{
		logger:         logging.NewNoopLogger(),
		encoderDecoder: mockencoding.NewMockEncoderDecoder(),
		tracer:         tracing.NewTracerForTest("test"),
	}
}

func TestProvideWebhooksService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		cfg := &Config{}

		actual, err := ProvideOAuth2Service(
			ctx,
			logging.NewNoopLogger(),
			cfg,
			database.NewMockDatabase(),
			mockencoding.NewMockEncoderDecoder(),
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}

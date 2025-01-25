package capitalism

import (
	"testing"

	capitalismmock "github.com/dinnerdonebetter/backend/internal/lib/capitalism/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	tracing "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/lib/routing/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:         logging.NewNoopLogger(),
		tracer:         tracing.NewTracerForTest("test"),
		paymentManager: capitalismmock.NewMockPaymentManager(),
	}
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		mpm := capitalismmock.NewMockPaymentManager()

		s := ProvideService(
			logger,
			tracing.NewNoopTracerProvider(),
			mpm,
		)
		assert.NotNil(t, s)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}

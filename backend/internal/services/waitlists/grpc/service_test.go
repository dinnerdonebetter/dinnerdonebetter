package grpc

import (
	"testing"

	waitlistmock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/waitlists/mock"
	waitlistssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/waitlists"

	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v5/observability/tracing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		waitlistsManager := &waitlistmock.Repository{}

		service := NewService(logger, tracerProvider, waitlistsManager)

		assert.NotNil(t, service)
		assert.Implements(t, (*waitlistssvc.WaitlistsServiceServer)(nil), service)
	})
}

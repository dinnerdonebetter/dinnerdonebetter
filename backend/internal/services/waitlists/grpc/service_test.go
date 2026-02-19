package grpc

import (
	"testing"

	waitlistmock "github.com/dinnerdonebetter/backend/internal/domain/waitlists/mock"
	waitlistssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/waitlists"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

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

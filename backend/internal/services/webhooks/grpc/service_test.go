package grpc

import (
	"testing"

	webhookmgrmock "github.com/dinnerdonebetter/backend/internal/domain/webhooks/manager/mock"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
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
		webhookManager := &webhookmgrmock.WebhookDataManager{}

		service := NewService(logger, tracerProvider, webhookManager)

		assert.NotNil(t, service)
		assert.Implements(t, (*webhookssvc.WebhooksServiceServer)(nil), service)
	})
}

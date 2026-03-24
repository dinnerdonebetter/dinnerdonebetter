package grpc

import (
	"testing"

	webhookmgrmock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/webhooks/manager/mock"
	webhookssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"

	"github.com/stretchr/testify/assert"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
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

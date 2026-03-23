package grpc

import (
	"testing"

	settingssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"

	"github.com/stretchr/testify/assert"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v2/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		msgConfig := &msgconfig.Config{}
		service := NewService(logger, tracerProvider, msgConfig, nil)

		assert.NotNil(t, service)
		assert.Implements(t, (*settingssvc.InternalOperationsServer)(nil), service)

		impl, ok := service.(*serviceImpl)
		assert.True(t, ok)
		assert.NotNil(t, impl.logger)
		assert.NotNil(t, impl.tracer)
		assert.Equal(t, msgConfig, impl.msgConfig)
	})
}

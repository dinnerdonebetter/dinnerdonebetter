package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func buildTestService(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	msgCfg := &msgconfig.Config{}

	service := &serviceImpl{
		tracer:    tracer,
		logger:    logger,
		msgConfig: msgCfg,
		sessionContextDataFetcher: func(ctx context.Context) (sessions.ContextData, error) {
			return sessions.ContextData{
				Requester: sessions.RequesterInfo{
					ServicePermissions: authorization.NewServiceRolePermissionChecker(authorization.ServiceAdminRole.String()),
				},
			}, nil
		},
	}

	return service
}

func buildTestServiceWithSessionError(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	msgCfg := &msgconfig.Config{}

	service := &serviceImpl{
		tracer:    tracer,
		logger:    logger,
		msgConfig: msgCfg,
		sessionContextDataFetcher: func(ctx context.Context) (sessions.ContextData, error) {
			return sessions.ContextData{}, errors.New("session error")
		},
	}

	return service
}

func buildTestServiceWithoutPermission(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	msgCfg := &msgconfig.Config{}

	service := &serviceImpl{
		tracer:    tracer,
		logger:    logger,
		msgConfig: msgCfg,
		sessionContextDataFetcher: func(ctx context.Context) (sessions.ContextData, error) {
			return sessions.ContextData{
				Requester: sessions.RequesterInfo{
					ServicePermissions: authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
				},
			}, nil
		},
	}

	return service
}

func TestNewService(t *testing.T) {
	t.Parallel()

	t.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider()
		msgConfig := &msgconfig.Config{}

		service := NewService(
			logger,
			tracerProvider,
			msgConfig,
		)

		assert.NotNil(t, service)
		assert.Implements(t, (*settingssvc.InternalOperationsServer)(nil), service)

		// Type assertion to ensure we get the correct implementation
		impl, ok := service.(*serviceImpl)
		assert.True(t, ok)
		assert.NotNil(t, impl.logger)
		assert.NotNil(t, impl.tracer)
		assert.Equal(t, msgConfig, impl.msgConfig)
	})
}

func TestServiceImpl_PublishArbitraryQueueMessage(t *testing.T) {
	t.Parallel()

	// Note: This test cannot easily test the full flow since it involves
	// creating actual message queue publishers, which would require complex
	// infrastructure setup. We test the validation and error paths instead.

	t.Run("with session error", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		service := buildTestServiceWithSessionError(t)

		request := &settingssvc.PublishArbitraryQueueMessageRequest{
			QueueName: "test-queue",
			Body:      "test message body",
		}

		response, err := service.PublishArbitraryQueueMessage(ctx, request)

		assert.Nil(t, response)
		assert.Error(t, err)

		grpcStatus, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, grpcStatus.Code())
	})
}

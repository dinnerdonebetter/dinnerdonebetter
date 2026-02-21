package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/assert"
)

func buildTestService(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	msgConfig := &msgconfig.Config{}

	service := &serviceImpl{
		tracer:          tracer,
		logger:          logger,
		msgConfig:       msgConfig,
		internalOpsRepo: nil,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					ServicePermissions: nil,
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
	msgConfig := &msgconfig.Config{}

	service := &serviceImpl{
		tracer:          tracer,
		logger:          logger,
		msgConfig:       msgConfig,
		internalOpsRepo: nil,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return nil, errors.New("session error")
		},
	}

	return service
}

func buildTestServiceWithInsufficientPermissions(t *testing.T) *serviceImpl {
	t.Helper()

	logger := logging.NewNoopLogger()
	tracer := tracing.NewTracerForTest(t.Name())
	msgConfig := &msgconfig.Config{}

	service := &serviceImpl{
		tracer:          tracer,
		logger:          logger,
		msgConfig:       msgConfig,
		internalOpsRepo: nil,
		sessionContextDataFetcher: func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{
					ServicePermissions: nil,
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

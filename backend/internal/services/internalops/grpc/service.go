package grpc

import (
	"context"

	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "configuration_service"
)

var _ settingssvc.InternalOperationsServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		settingssvc.UnimplementedInternalOperationsServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (sessions.ContextData, error)
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) settingssvc.InternalOperationsServer {
	return &serviceImpl{
		logger: logging.EnsureLogger(logger).WithName(o11yName),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
	}
}

func (s serviceImpl) PublishArbitraryQueueMessage(ctx context.Context, request *settingssvc.PublishArbitraryQueueMessageRequest) (*settingssvc.PublishArbitraryQueueMessageResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	x := &settingssvc.PublishArbitraryQueueMessageResponse{
		ResponseDetails: &types.ResponseDetails{TraceID: span.SpanContext().TraceID().String()},
		Success:         false,
	}

	return x, nil
}

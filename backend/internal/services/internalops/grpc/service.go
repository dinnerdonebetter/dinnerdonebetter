package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"google.golang.org/grpc/codes"
)

const (
	o11yName = "internalops_service"
)

var _ settingssvc.InternalOperationsServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		settingssvc.UnimplementedInternalOperationsServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		msgConfig                 *msgconfig.Config
		sessionContextDataFetcher func(context.Context) (sessions.ContextData, error)
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	msgConfig *msgconfig.Config,
) settingssvc.InternalOperationsServer {
	return &serviceImpl{
		msgConfig: msgConfig,
		logger:    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
	}
}

func (s serviceImpl) PublishArbitraryQueueMessage(ctx context.Context, request *settingssvc.PublishArbitraryQueueMessageRequest) (*settingssvc.PublishArbitraryQueueMessageResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	if !sessionContextData.ServiceRolePermissionChecker().HasPermission(authorization.PublishArbitraryQueueMessagePermission) {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "user is not permitted")
	}

	pp, err := msgconfig.ProvidePublisherProvider(ctx, s.logger, tracing.NewNoopTracerProvider(), s.msgConfig)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "establishing publisher provider")
	}

	publisher, err := pp.ProvidePublisher(request.QueueName)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "initializing publisher")
	}

	if err = publisher.Publish(ctx, []byte(request.Body)); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "publishing message")
	}

	x := &settingssvc.PublishArbitraryQueueMessageResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Success: true,
	}

	return x, nil
}

package grpc

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/internalops"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"google.golang.org/grpc/codes"
)

const (
	o11yName = "internalops_service"

	testQueueMessagePollInterval = 500 * time.Millisecond
	testQueueMessageTimeout      = 30 * time.Second
)

var _ settingssvc.InternalOperationsServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		settingssvc.UnimplementedInternalOperationsServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		encoder                   encoding.ServerEncoderDecoder
		msgConfig                 *msgconfig.Config
		internalOpsRepo           internalops.InternalOpsDataManager
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	msgConfig *msgconfig.Config,
	repo internalops.InternalOpsDataManager,
) settingssvc.InternalOperationsServer {
	return &serviceImpl{
		msgConfig:                 msgConfig,
		internalOpsRepo:           repo,
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
	}
}

func (s *serviceImpl) TestQueueMessage(ctx context.Context, request *settingssvc.TestQueueMessageRequest) (*settingssvc.TestQueueMessageResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue("queue_name", request.QueueName)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "fetching session context data")
	}

	if !sessionContextData.ServiceRolePermissionChecker().HasPermission(authorization.PublishArbitraryQueueMessagePermission) {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.PermissionDenied, "user is not permitted")
	}

	if request.QueueName == "" {
		return nil, observability.PrepareAndLogGRPCStatus(nil, logger, span, codes.InvalidArgument, "queue_name is required")
	}

	testID := identifiers.New()
	start := time.Now()

	if err = s.internalOpsRepo.CreateQueueTestMessage(ctx, testID, request.QueueName); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "creating queue test message record")
	}

	msg := &audit.DataChangeMessage{
		EventType: internalops.QueueTestMessageEventType,
		Context: map[string]any{
			"test_id":    testID,
			"queue_name": request.QueueName,
		},
		UserID: sessionContextData.Requester.UserID,
	}

	msg := &audit.DataChangeMessage{
		EventType: internalops.QueueTestMessageEventType,
		Context: map[string]any{
			"test_id":    testID,
			"queue_name": request.QueueName,
		},
		UserID: sessionContextData.Requester.UserID,
	}

	msgBytes := s.encoder.MustEncodeJSON(ctx, msg)

	pp, err := msgconfig.ProvidePublisherProvider(ctx, s.logger, tracing.NewNoopTracerProvider(), s.msgConfig)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "establishing publisher provider")
	}

	publisher, err := pp.ProvidePublisher(ctx, request.QueueName)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "initializing publisher")
	}

	if err = publisher.Publish(ctx, msg); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "publishing test message")
	}

	ticker := time.NewTicker(testQueueMessagePollInterval)
	defer ticker.Stop()

	deadline := time.After(testQueueMessageTimeout)

	for {
		select {
		case <-deadline:
			return nil, observability.PrepareAndLogGRPCStatus(nil, logger, span, codes.DeadlineExceeded, "timed out waiting for queue test message acknowledgment")
		case <-ctx.Done():
			return nil, observability.PrepareAndLogGRPCStatus(ctx.Err(), logger, span, codes.Canceled, "context canceled while waiting for acknowledgment")
		case <-ticker.C:
			record, pollErr := s.internalOpsRepo.GetQueueTestMessage(ctx, testID)
			if pollErr != nil {
				logger.Error("polling for queue test message acknowledgment", pollErr)
				continue
			}

			if record.AcknowledgedAt != nil {
				roundTripMs := time.Since(start).Milliseconds()
				return &settingssvc.TestQueueMessageResponse{
					ResponseDetails: &types.ResponseDetails{
						TraceId: span.SpanContext().TraceID().String(),
					},
					Success:     true,
					TestId:      testID,
					RoundTripMs: roundTripMs,
				}, nil
			}
		}
	}
}

package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"

	"google.golang.org/grpc/codes"
)

func (s *ServiceImpl) ArchiveWebhook(ctx context.Context, request *configurationsvc.ArchiveWebhookRequest) (*configurationsvc.ArchiveWebhookResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	if err = s.webhookRepository.ArchiveWebhook(ctx, request.WebhookID, sessionContextData.ActiveAccountID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive webhook")
	}

	x := &configurationsvc.ArchiveWebhookResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) ArchiveWebhookTriggerEvent(ctx context.Context, request *configurationsvc.ArchiveWebhookTriggerEventRequest) (*configurationsvc.ArchiveWebhookTriggerEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.WebhookIDKey, request.WebhookID).WithValue(keys.WebhookTriggerEventIDKey, request.WebhookTriggerEventID)

	if err := s.webhookRepository.ArchiveWebhookTriggerEvent(ctx, request.WebhookID, request.WebhookTriggerEventID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive webhook trigger event")
	}

	x := &configurationsvc.ArchiveWebhookTriggerEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *ServiceImpl) CreateWebhook(ctx context.Context, request *configurationsvc.CreateWebhookRequest) (*configurationsvc.CreateWebhookResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	created, err := s.webhookRepository.CreateWebhook(ctx, ConvertGRPCWebhookCreationRequestInputToWebhookDatabaseCreationInput(request.Input, sessionContextData.ActiveAccountID))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create webhook")
	}

	x := &configurationsvc.CreateWebhookResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: ConvertWebhookToGRPCWebhook(created),
	}

	return x, nil
}

func (s *ServiceImpl) AddWebhookTriggerEvent(ctx context.Context, request *configurationsvc.AddWebhookTriggerEventRequest) (*configurationsvc.AddWebhookTriggerEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	created, err := s.webhookRepository.AddWebhookTriggerEvent(ctx, sessionContextData.ActiveAccountID, ConvertGRPCWebhookTriggerEventDatabaseCreationInputToWebhookTriggerEventDatabaseCreationInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to add webhook trigger event")
	}

	x := &configurationsvc.AddWebhookTriggerEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(created),
	}

	return x, nil
}

func (s *ServiceImpl) GetWebhook(ctx context.Context, request *configurationsvc.GetWebhookRequest) (*configurationsvc.GetWebhookResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span).WithValue(keys.WebhookIDKey, request.WebhookID)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	webhook, err := s.webhookRepository.GetWebhook(ctx, request.WebhookID, sessionContextData.ActiveAccountID)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch webhook")
	}

	x := &configurationsvc.GetWebhookResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: ConvertWebhookToGRPCWebhook(webhook),
	}

	return x, nil
}

func (s *ServiceImpl) GetWebhooks(ctx context.Context, request *configurationsvc.GetWebhooksRequest) (*configurationsvc.GetWebhooksResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	retrieved, err := s.webhookRepository.GetWebhooks(ctx, sessionContextData.ActiveAccountID, grpcconverters.ConvertGRPCQueryFilterToQueryFilter(request.Filter))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to fetch webhooks")
	}

	x := &configurationsvc.GetWebhooksResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, webhook := range retrieved.Data {
		x.Results = append(x.Results, ConvertWebhookToGRPCWebhook(webhook))
	}

	return x, nil
}

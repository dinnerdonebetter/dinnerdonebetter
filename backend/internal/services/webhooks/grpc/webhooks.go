package grpc

import (
	"context"

	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/grpc/generated/types"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/services/webhooks/grpc/converters"

	"google.golang.org/grpc/codes"
)

func (s *serviceImpl) CreateWebhook(ctx context.Context, request *webhookssvc.CreateWebhookRequest) (*webhookssvc.CreateWebhookResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	created, err := s.webhookRepository.CreateWebhook(ctx, converters.ConvertGRPCWebhookCreationRequestInputToWebhookDatabaseCreationInput(request.Input, sessionContextData.ActiveAccountID))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to create webhook")
	}

	x := &webhookssvc.CreateWebhookResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertWebhookToGRPCWebhook(created),
	}

	return x, nil
}

func (s *serviceImpl) AddWebhookTriggerEvent(ctx context.Context, request *webhookssvc.AddWebhookTriggerEventRequest) (*webhookssvc.AddWebhookTriggerEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	sessionContextData, err := s.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Unauthenticated, "failed to fetch session context data")
	}
	logger = logger.WithValue(keys.AccountIDKey, sessionContextData.ActiveAccountID)

	created, err := s.webhookRepository.AddWebhookTriggerEvent(ctx, sessionContextData.ActiveAccountID, converters.ConvertGRPCWebhookTriggerEventDatabaseCreationInputToWebhookTriggerEventDatabaseCreationInput(request.Input))
	if err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to add webhook trigger event")
	}

	x := &webhookssvc.AddWebhookTriggerEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Created: converters.ConvertWebhookTriggerEventToGRPCWebhookTriggerEvent(created),
	}

	return x, nil
}

func (s *serviceImpl) GetWebhook(ctx context.Context, request *webhookssvc.GetWebhookRequest) (*webhookssvc.GetWebhookResponse, error) {
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

	x := &webhookssvc.GetWebhookResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
		Result: converters.ConvertWebhookToGRPCWebhook(webhook),
	}

	return x, nil
}

func (s *serviceImpl) GetWebhooks(ctx context.Context, request *webhookssvc.GetWebhooksRequest) (*webhookssvc.GetWebhooksResponse, error) {
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

	x := &webhookssvc.GetWebhooksResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	for _, webhook := range retrieved.Data {
		x.Results = append(x.Results, converters.ConvertWebhookToGRPCWebhook(webhook))
	}

	return x, nil
}

func (s *serviceImpl) ArchiveWebhook(ctx context.Context, request *webhookssvc.ArchiveWebhookRequest) (*webhookssvc.ArchiveWebhookResponse, error) {
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

	x := &webhookssvc.ArchiveWebhookResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

func (s *serviceImpl) ArchiveWebhookTriggerEvent(ctx context.Context, request *webhookssvc.ArchiveWebhookTriggerEventRequest) (*webhookssvc.ArchiveWebhookTriggerEventResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.WebhookIDKey, request.WebhookID).WithValue(keys.WebhookTriggerEventIDKey, request.WebhookTriggerEventID)

	if err := s.webhookRepository.ArchiveWebhookTriggerEvent(ctx, request.WebhookID, request.WebhookTriggerEventID); err != nil {
		return nil, observability.PrepareAndLogGRPCStatus(err, logger, span, codes.Internal, "failed to archive webhook trigger event")
	}

	x := &webhookssvc.ArchiveWebhookTriggerEventResponse{
		ResponseDetails: &types.ResponseDetails{
			TraceID: span.SpanContext().TraceID().String(),
		},
	}

	return x, nil
}

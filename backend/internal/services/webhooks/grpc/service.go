package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks/manager"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "configuration_service"
)

var _ webhookssvc.WebhooksServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		webhookssvc.UnimplementedWebhooksServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		webhookManager            manager.WebhookDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	webhookManager manager.WebhookDataManager,
) webhookssvc.WebhooksServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
		webhookManager:            webhookManager,
	}
}

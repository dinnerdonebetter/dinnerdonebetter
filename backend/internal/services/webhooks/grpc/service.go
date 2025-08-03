package grpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "configuration_service"
)

var _ settingssvc.UserConfigurationServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		settingssvc.UnimplementedUserConfigurationServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (sessions.ContextData, error)
		webhookRepository         webhooks.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	webhookRepository webhooks.Repository,
) settingssvc.UserConfigurationServiceServer {
	return &serviceImpl{
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		webhookRepository: webhookRepository,
	}
}

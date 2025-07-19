package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "configuration_service"
)

var _ configurationsvc.UserConfigurationServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		configurationsvc.UnimplementedUserConfigurationServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		accountIDFetcher          func(x any) (string, error)
		webhookRepository         webhooks.Repository
		serviceSettingsRepository settings.Repository
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	webhookRepository webhooks.Repository,
	settingsRepository settings.Repository,
) configurationsvc.UserConfigurationServiceServer {
	return &ServiceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		webhookRepository:         webhookRepository,
		serviceSettingsRepository: settingsRepository,
	}
}

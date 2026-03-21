package grpc

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/sessions"
	settingsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/settings/manager"
	settingssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/settings"

	"github.com/verygoodsoftwarenotvirus/platform/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/observability/tracing"
)

const (
	o11yName = "configuration_service"
)

var _ settingssvc.SettingsServiceServer = (*serviceImpl)(nil)

type (
	serviceImpl struct {
		settingssvc.UnimplementedSettingsServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		sessionContextDataFetcher func(context.Context) (*sessions.ContextData, error)
		settingsManager           settingsmanager.SettingsDataManager
	}
)

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	settingsManager settingsmanager.SettingsDataManager,
) settingssvc.SettingsServiceServer {
	return &serviceImpl{
		logger:                    logging.EnsureLogger(logger).WithName(o11yName),
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		settingsManager:           settingsManager,
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
	}
}

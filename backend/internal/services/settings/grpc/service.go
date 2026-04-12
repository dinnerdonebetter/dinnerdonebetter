package grpc

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/sessions"
	settingsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/settings/manager"
	settingssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/settings"

	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"
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
		logger:                    logging.NewNamedLogger(logger, o11yName),
		tracer:                    tracing.NewNamedTracer(tracerProvider, o11yName),
		settingsManager:           settingsManager,
		sessionContextDataFetcher: sessions.FetchContextDataFromContext,
	}
}

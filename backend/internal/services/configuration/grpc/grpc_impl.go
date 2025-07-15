package grpc

import (
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	configurationsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/configuration"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var _ configurationsvc.UserConfigurationServiceServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		configurationsvc.UnimplementedUserConfigurationServiceServer
		tracer                    tracing.Tracer
		logger                    logging.Logger
		webhookRepository         webhooks.Repository
		serviceSettingsRepository settings.Repository
	}
)

func NewService(
	webhookRepository webhooks.Repository,
	settingsRepository settings.Repository,
) configurationsvc.UserConfigurationServiceServer {
	return &ServiceImpl{
		webhookRepository:         webhookRepository,
		serviceSettingsRepository: settingsRepository,
	}
}

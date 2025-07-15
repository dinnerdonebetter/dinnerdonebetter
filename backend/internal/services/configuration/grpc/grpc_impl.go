package grpc

import (
	"context"

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

func (s *ServiceImpl) GetServiceSettingConfigurationsForAccount(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationsForAccountRequest) (*configurationsvc.GetServiceSettingConfigurationsForAccountResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetServiceSettingConfigurationsForUser(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationsForUserRequest) (*configurationsvc.GetServiceSettingConfigurationsForUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) SearchForServiceSettings(ctx context.Context, request *configurationsvc.SearchForServiceSettingsRequest) (*configurationsvc.SearchForServiceSettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveServiceSetting(ctx context.Context, request *configurationsvc.ArchiveServiceSettingRequest) (*configurationsvc.ArchiveServiceSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveServiceSettingConfiguration(ctx context.Context, request *configurationsvc.ArchiveServiceSettingConfigurationRequest) (*configurationsvc.ArchiveServiceSettingConfigurationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateServiceSetting(ctx context.Context, request *configurationsvc.CreateServiceSettingRequest) (*configurationsvc.CreateServiceSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateServiceSettingConfiguration(ctx context.Context, request *configurationsvc.CreateServiceSettingConfigurationRequest) (*configurationsvc.CreateServiceSettingConfigurationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetServiceSetting(ctx context.Context, request *configurationsvc.GetServiceSettingRequest) (*configurationsvc.GetServiceSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetServiceSettingConfigurationByName(ctx context.Context, request *configurationsvc.GetServiceSettingConfigurationByNameRequest) (*configurationsvc.GetServiceSettingConfigurationByNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetServiceSettings(ctx context.Context, request *configurationsvc.GetServiceSettingsRequest) (*configurationsvc.GetServiceSettingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) UpdateServiceSettingConfiguration(ctx context.Context, request *configurationsvc.UpdateServiceSettingConfigurationRequest) (*configurationsvc.UpdateServiceSettingConfigurationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveWebhook(ctx context.Context, request *configurationsvc.ArchiveWebhookRequest) (*configurationsvc.ArchiveWebhookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) ArchiveWebhookTriggerEvent(ctx context.Context, request *configurationsvc.ArchiveWebhookTriggerEventRequest) (*configurationsvc.ArchiveWebhookTriggerEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateWebhook(ctx context.Context, request *configurationsvc.CreateWebhookRequest) (*configurationsvc.CreateWebhookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) CreateWebhookTriggerEvent(ctx context.Context, request *configurationsvc.CreateWebhookTriggerEventRequest) (*configurationsvc.CreateWebhookTriggerEventResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetWebhook(ctx context.Context, request *configurationsvc.GetWebhookRequest) (*configurationsvc.GetWebhookResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServiceImpl) GetWebhooks(ctx context.Context, request *configurationsvc.GetWebhooksRequest) (*configurationsvc.GetWebhooksResponse, error) {
	//TODO implement me
	panic("implement me")
}

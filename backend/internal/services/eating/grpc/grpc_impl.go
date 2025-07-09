package eatinggrpc

import (
	"context"
	"errors"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	eatingsvc "github.com/dinnerdonebetter/backend/internal/grpc/service/eating"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

type (
	ServiceImpl struct {
		eatingsvc.UnimplementedEatingServer
		tracer                   tracing.Tracer
		logger                   logging.Logger
		recipeManager            managers.RecipeManager
		validEnumerationsManager managers.ValidEnumerationsManager
		mealPlanningManager      managers.MealPlanningManager
	}
)

func NewService(
	recipeManager managers.RecipeManager,
	validEnumerationsManager managers.ValidEnumerationsManager,
	mealPlanningManager managers.MealPlanningManager,
) *ServiceImpl {
	return &ServiceImpl{
		recipeManager:            recipeManager,
		validEnumerationsManager: validEnumerationsManager,
		mealPlanningManager:      mealPlanningManager,
	}
}

var (
	errUnimplemented = errors.New("unimplemented")
)

func (s *ServiceImpl) Ping(ctx context.Context, request *messages.PingRequest) (*messages.PingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return &messages.PingResponse{}, nil
}

func (s *ServiceImpl) PublishArbitraryQueueMessage(ctx context.Context, request *messages.PublishArbitraryQueueMessageRequest) (*messages.PublishArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) DestroyAllUserData(ctx context.Context, request *messages.DestroyAllUserDataRequest) (*messages.DestroyAllUserDataResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) FetchUserDataReport(ctx context.Context, request *messages.FetchUserDataReportRequest) (*messages.FetchUserDataReportResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) AdminUpdateUserStatus(ctx context.Context, request *messages.AdminUpdateUserStatusRequest) (*messages.AdminUpdateUserStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) AggregateUserDataReport(ctx context.Context, request *messages.AggregateUserDataReportRequest) (*messages.AggregateUserDataReportResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*messages.ArchiveOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*messages.ArchiveServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*messages.ArchiveServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*messages.ArchiveWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*messages.ArchiveWebhookTriggerEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateOAuth2Client(ctx context.Context, request *messages.CreateOAuth2ClientRequest) (*messages.CreateOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateServiceSetting(ctx context.Context, request *messages.CreateServiceSettingRequest) (*messages.CreateServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateServiceSettingConfiguration(ctx context.Context, request *messages.CreateServiceSettingConfigurationRequest) (*messages.CreateServiceSettingConfigurationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateWebhook(ctx context.Context, request *messages.CreateWebhookRequest) (*messages.CreateWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.CreateWebhookTriggerEventResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.GetAuditLogEntryByIDResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetAuthStatus(ctx context.Context, request *messages.GetAuthStatusRequest) (*messages.GetAuthStatusResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.GetOAuth2ClientResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.GetOAuth2ClientsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.GetServiceSettingResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.GetServiceSettingConfigurationByNameResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.GetServiceSettingsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.GetWebhookResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.GetWebhooksResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.SearchForServiceSettingsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*messages.ArchiveUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) CreateUserIngredientPreference(ctx context.Context, request *messages.CreateUserIngredientPreferenceRequest) (*messages.CreateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.GetUserIngredientPreferencesResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

func (s *ServiceImpl) UpdateUserIngredientPreference(ctx context.Context, request *messages.UpdateUserIngredientPreferenceRequest) (*messages.UpdateUserIngredientPreferenceResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithSpan(span)

	return nil, observability.PrepareAndLogError(errUnimplemented, logger, span, "unimplemented")
}

package serverimpl

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) AcceptHouseholdInvitation(ctx context.Context, request *messages.AcceptHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveHousehold(ctx context.Context, request *messages.ArchiveHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveHouseholdInstrumentOwnership(ctx context.Context, request *messages.ArchiveHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveUserIngredientPreference(ctx context.Context, request *messages.ArchiveUserIngredientPreferenceRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*messages.WebhookTriggerEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CancelHouseholdInvitation(ctx context.Context, request *messages.CancelHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CheckForReadiness(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateHousehold(ctx context.Context, input *messages.HouseholdCreationRequestInput) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateHouseholdInstrumentOwnership(ctx context.Context, input *messages.HouseholdInstrumentOwnershipCreationRequestInput) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateHouseholdInvitation(ctx context.Context, request *messages.CreateHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateOAuth2Client(ctx context.Context, input *messages.OAuth2ClientCreationRequestInput) (*messages.OAuth2ClientCreationResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateServiceSetting(ctx context.Context, input *messages.ServiceSettingCreationRequestInput) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateServiceSettingConfiguration(ctx context.Context, input *messages.ServiceSettingConfigurationCreationRequestInput) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateUserIngredientPreference(ctx context.Context, input *messages.UserIngredientPreferenceCreationRequestInput) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateWebhook(ctx context.Context, input *messages.WebhookCreationRequestInput) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.WebhookTriggerEvent, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetActiveHousehold(ctx context.Context, _ *emptypb.Empty) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntriesForHousehold(ctx context.Context, request *messages.GetAuditLogEntriesForHouseholdRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.AuditLogEntry, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHousehold(ctx context.Context, request *messages.GetHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInstrumentOwnership(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInstrumentOwnerships(ctx context.Context, request *messages.GetHouseholdInstrumentOwnershipsRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInvitation(ctx context.Context, request *messages.GetHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholdInvitationByID(ctx context.Context, request *messages.GetHouseholdInvitationByIDRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetHouseholds(ctx context.Context, request *messages.GetHouseholdsRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.OAuth2Client, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetRandomValidInstrument(ctx context.Context, _ *emptypb.Empty) (*messages.ValidInstrument, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetReceivedHouseholdInvitations(ctx context.Context, request *messages.GetReceivedHouseholdInvitationsRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetSelf(ctx context.Context, _ *emptypb.Empty) (*messages.User, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetSentHouseholdInvitations(ctx context.Context, request *messages.GetSentHouseholdInvitationsRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettingConfigurationsForHousehold(ctx context.Context, request *messages.GetServiceSettingConfigurationsForHouseholdRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserIngredientPreferences(ctx context.Context, request *messages.GetUserIngredientPreferencesRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.Webhook, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) PublishArbitraryQueueMessage(ctx context.Context, input *messages.ArbitraryQueueMessageRequestInput) (*messages.ArbitraryQueueMessageResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) RejectHouseholdInvitation(ctx context.Context, request *messages.RejectHouseholdInvitationRequest) (*messages.HouseholdInvitation, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.ServiceSetting, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) SetDefaultHousehold(ctx context.Context, request *messages.SetDefaultHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) TransferHouseholdOwnership(ctx context.Context, request *messages.TransferHouseholdOwnershipRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHousehold(ctx context.Context, request *messages.UpdateHouseholdRequest) (*messages.Household, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHouseholdInstrumentOwnership(ctx context.Context, request *messages.UpdateHouseholdInstrumentOwnershipRequest) (*messages.HouseholdInstrumentOwnership, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateHouseholdMemberPermissions(ctx context.Context, request *messages.UpdateHouseholdMemberPermissionsRequest) (*messages.UserPermissionsResponse, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateServiceSettingConfiguration(ctx context.Context, request *messages.UpdateServiceSettingConfigurationRequest) (*messages.ServiceSettingConfiguration, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) UpdateUserIngredientPreference(ctx context.Context, request *messages.UpdateUserIngredientPreferenceRequest) (*messages.UserIngredientPreference, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserNotification(ctx context.Context, request *messages.GetUserNotificationRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

func (s *Server) GetUserNotifications(ctx context.Context, request *messages.GetUserNotificationsRequest) (*messages.UserNotification, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return nil, Unimplemented()
}

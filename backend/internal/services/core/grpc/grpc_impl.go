package coregrpc

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	coresvc "github.com/dinnerdonebetter/backend/internal/grpc/service/core"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	serviceName = "core_service"
)

var _ coresvc.CoreServer = (*ServiceImpl)(nil)

type (
	ServiceImpl struct {
		coresvc.UnimplementedCoreServer
		tracer tracing.Tracer
		logger logging.Logger
		authentication.Authenticator
	}
)

func NewCoreService(tracerProvider tracing.TracerProvider, logger logging.Logger) *ServiceImpl {
	return &ServiceImpl{
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		logger: logger,
	}
}

func (s *ServiceImpl) Ping(ctx context.Context, request *messages.PingRequest) (*messages.PingResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ExchangeToken(ctx context.Context, request *messages.ExchangeTokenRequest) (*messages.ExchangeTokenResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) AcceptHouseholdInvitation(ctx context.Context, request *messages.AcceptHouseholdInvitationRequest) (*messages.AcceptHouseholdInvitationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) AdminLoginForToken(ctx context.Context, request *messages.AdminLoginForTokenRequest) (*messages.AdminLoginForTokenResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) AdminUpdateUserStatus(ctx context.Context, request *messages.AdminUpdateUserStatusRequest) (*messages.AdminUpdateUserStatusResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) AggregateUserDataReport(ctx context.Context, request *messages.AggregateUserDataReportRequest) (*messages.AggregateUserDataReportResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveHousehold(ctx context.Context, request *messages.ArchiveHouseholdRequest) (*messages.ArchiveHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveOAuth2Client(ctx context.Context, request *messages.ArchiveOAuth2ClientRequest) (*messages.ArchiveOAuth2ClientResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveServiceSetting(ctx context.Context, request *messages.ArchiveServiceSettingRequest) (*messages.ArchiveServiceSettingResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveServiceSettingConfiguration(ctx context.Context, request *messages.ArchiveServiceSettingConfigurationRequest) (*messages.ArchiveServiceSettingConfigurationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveUser(ctx context.Context, request *messages.ArchiveUserRequest) (*messages.ArchiveUserResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveUserMembership(ctx context.Context, request *messages.ArchiveUserMembershipRequest) (*messages.ArchiveUserMembershipResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveWebhook(ctx context.Context, request *messages.ArchiveWebhookRequest) (*messages.ArchiveWebhookResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) ArchiveWebhookTriggerEvent(ctx context.Context, request *messages.ArchiveWebhookTriggerEventRequest) (*messages.ArchiveWebhookTriggerEventResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CancelHouseholdInvitation(ctx context.Context, request *messages.CancelHouseholdInvitationRequest) (*messages.CancelHouseholdInvitationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CheckPermissions(ctx context.Context, request *messages.CheckPermissionsRequest) (*messages.CheckPermissionsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateHousehold(ctx context.Context, request *messages.CreateHouseholdRequest) (*messages.CreateHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateHouseholdInvitation(ctx context.Context, request *messages.CreateHouseholdInvitationRequest) (*messages.CreateHouseholdInvitationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateOAuth2Client(ctx context.Context, request *messages.CreateOAuth2ClientRequest) (*messages.CreateOAuth2ClientResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateServiceSetting(ctx context.Context, request *messages.CreateServiceSettingRequest) (*messages.CreateServiceSettingResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateServiceSettingConfiguration(ctx context.Context, request *messages.CreateServiceSettingConfigurationRequest) (*messages.CreateServiceSettingConfigurationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateUser(ctx context.Context, request *messages.CreateUserRequest) (*messages.CreateUserResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateUserNotification(ctx context.Context, request *messages.CreateUserNotificationRequest) (*messages.CreateUserNotificationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateWebhook(ctx context.Context, request *messages.CreateWebhookRequest) (*messages.CreateWebhookResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) CreateWebhookTriggerEvent(ctx context.Context, request *messages.CreateWebhookTriggerEventRequest) (*messages.CreateWebhookTriggerEventResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) DestroyAllUserData(ctx context.Context, request *messages.DestroyAllUserDataRequest) (*messages.DestroyAllUserDataResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) FetchUserDataReport(ctx context.Context, request *messages.FetchUserDataReportRequest) (*messages.FetchUserDataReportResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetActiveHousehold(ctx context.Context, request *messages.GetActiveHouseholdRequest) (*messages.GetActiveHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetAuditLogEntriesForHousehold(ctx context.Context, request *messages.GetAuditLogEntriesForHouseholdRequest) (*messages.GetAuditLogEntriesForHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetAuditLogEntriesForUser(ctx context.Context, request *messages.GetAuditLogEntriesForUserRequest) (*messages.GetAuditLogEntriesForUserResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetAuditLogEntryByID(ctx context.Context, request *messages.GetAuditLogEntryByIDRequest) (*messages.GetAuditLogEntryByIDResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetAuthStatus(ctx context.Context, request *messages.GetAuthStatusRequest) (*messages.GetAuthStatusResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetHousehold(ctx context.Context, request *messages.GetHouseholdRequest) (*messages.GetHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetHouseholdInvitation(ctx context.Context, request *messages.GetHouseholdInvitationRequest) (*messages.GetHouseholdInvitationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetHouseholdInvitationByID(ctx context.Context, request *messages.GetHouseholdInvitationByIDRequest) (*messages.GetHouseholdInvitationByIDResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetHouseholds(ctx context.Context, request *messages.GetHouseholdsRequest) (*messages.GetHouseholdsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetOAuth2Client(ctx context.Context, request *messages.GetOAuth2ClientRequest) (*messages.GetOAuth2ClientResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetOAuth2Clients(ctx context.Context, request *messages.GetOAuth2ClientsRequest) (*messages.GetOAuth2ClientsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetReceivedHouseholdInvitations(ctx context.Context, request *messages.GetReceivedHouseholdInvitationsRequest) (*messages.GetReceivedHouseholdInvitationsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetSelf(ctx context.Context, request *messages.GetSelfRequest) (*messages.GetSelfResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetSentHouseholdInvitations(ctx context.Context, request *messages.GetSentHouseholdInvitationsRequest) (*messages.GetSentHouseholdInvitationsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetServiceSetting(ctx context.Context, request *messages.GetServiceSettingRequest) (*messages.GetServiceSettingResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetServiceSettingConfigurationByName(ctx context.Context, request *messages.GetServiceSettingConfigurationByNameRequest) (*messages.GetServiceSettingConfigurationByNameResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetServiceSettingConfigurationsForHousehold(ctx context.Context, request *messages.GetServiceSettingConfigurationsForHouseholdRequest) (*messages.GetServiceSettingConfigurationsForHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetServiceSettingConfigurationsForUser(ctx context.Context, request *messages.GetServiceSettingConfigurationsForUserRequest) (*messages.GetServiceSettingConfigurationsForUserResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetServiceSettings(ctx context.Context, request *messages.GetServiceSettingsRequest) (*messages.GetServiceSettingsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetUser(ctx context.Context, request *messages.GetUserRequest) (*messages.GetUserResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetUserNotification(ctx context.Context, request *messages.GetUserNotificationRequest) (*messages.GetUserNotificationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetUserNotifications(ctx context.Context, request *messages.GetUserNotificationsRequest) (*messages.GetUserNotificationsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetUsers(ctx context.Context, request *messages.GetUsersRequest) (*messages.GetUsersResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetWebhook(ctx context.Context, request *messages.GetWebhookRequest) (*messages.GetWebhookResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) GetWebhooks(ctx context.Context, request *messages.GetWebhooksRequest) (*messages.GetWebhooksResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) LoginForToken(ctx context.Context, request *messages.LoginForTokenRequest) (*messages.LoginForTokenResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) PublishArbitraryQueueMessage(ctx context.Context, request *messages.PublishArbitraryQueueMessageRequest) (*messages.PublishArbitraryQueueMessageResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) RedeemPasswordResetToken(ctx context.Context, request *messages.RedeemPasswordResetTokenRequest) (*messages.RedeemPasswordResetTokenResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) RefreshTOTPSecret(ctx context.Context, request *messages.RefreshTOTPSecretRequest) (*messages.RefreshTOTPSecretResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) RejectHouseholdInvitation(ctx context.Context, request *messages.RejectHouseholdInvitationRequest) (*messages.RejectHouseholdInvitationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) RequestEmailVerificationEmail(ctx context.Context, request *messages.RequestEmailVerificationEmailRequest) (*messages.RequestEmailVerificationEmailResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) RequestPasswordResetToken(ctx context.Context, request *messages.RequestPasswordResetTokenRequest) (*messages.RequestPasswordResetTokenResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) RequestUsernameReminder(ctx context.Context, request *messages.RequestUsernameReminderRequest) (*messages.RequestUsernameReminderResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) SearchForServiceSettings(ctx context.Context, request *messages.SearchForServiceSettingsRequest) (*messages.SearchForServiceSettingsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) SearchForUsers(ctx context.Context, request *messages.SearchForUsersRequest) (*messages.SearchForUsersResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) SetDefaultHousehold(ctx context.Context, request *messages.SetDefaultHouseholdRequest) (*messages.SetDefaultHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) TransferHouseholdOwnership(ctx context.Context, request *messages.TransferHouseholdOwnershipRequest) (*messages.TransferHouseholdOwnershipResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateHousehold(ctx context.Context, request *messages.UpdateHouseholdRequest) (*messages.UpdateHouseholdResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateHouseholdMemberPermissions(ctx context.Context, request *messages.UpdateHouseholdMemberPermissionsRequest) (*messages.UpdateHouseholdMemberPermissionsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdatePassword(ctx context.Context, request *messages.UpdatePasswordRequest) (*messages.UpdatePasswordResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateServiceSettingConfiguration(ctx context.Context, request *messages.UpdateServiceSettingConfigurationRequest) (*messages.UpdateServiceSettingConfigurationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateUserDetails(ctx context.Context, request *messages.UpdateUserDetailsRequest) (*messages.UpdateUserDetailsResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateUserEmailAddress(ctx context.Context, request *messages.UpdateUserEmailAddressRequest) (*messages.UpdateUserEmailAddressResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateUserNotification(ctx context.Context, request *messages.UpdateUserNotificationRequest) (*messages.UpdateUserNotificationResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UpdateUserUsername(ctx context.Context, request *messages.UpdateUserUsernameRequest) (*messages.UpdateUserUsernameResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) UploadUserAvatar(ctx context.Context, request *messages.UploadUserAvatarRequest) (*messages.UploadUserAvatarResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) VerifyEmailAddress(ctx context.Context, request *messages.VerifyEmailAddressRequest) (*messages.VerifyEmailAddressResponse, error) {
	return nil, nil
}

func (s *ServiceImpl) VerifyTOTPSecret(ctx context.Context, request *messages.VerifyTOTPSecretRequest) (*messages.VerifyTOTPSecretResponse, error) {
	return nil, nil
}

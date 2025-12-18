package grpcapi

import (
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	issuereportssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/issue_reports"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"
	waitlistssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/waitlists"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
)

type GRPCService struct {
	auditsvc.AuditServiceServer
	authsvc.AuthServiceServer
	dataprivacysvc.DataPrivacyServiceServer
	identitysvc.IdentityServiceServer
	internalopssvc.InternalOperationsServer
	issuereportssvc.IssueReportsServiceServer
	mealplanningsvc.MealPlanningServiceServer
	notificationssvc.UserNotificationsServiceServer
	oauthsvc.OAuthServiceServer
	settingssvc.SettingsServiceServer
	uploadedmediasvc.UploadedMediaServiceServer
	waitlistssvc.WaitlistsServiceServer
	webhookssvc.WebhooksServiceServer
	*grpc.Server
}

func NewGRPCService(
	auditServiceServer auditsvc.AuditServiceServer,
	authServiceServer authsvc.AuthServiceServer,
	dataPrivacyServiceServer dataprivacysvc.DataPrivacyServiceServer,
	identityServiceServer identitysvc.IdentityServiceServer,
	internalOperationsServer internalopssvc.InternalOperationsServer,
	issueReportsServiceServer issuereportssvc.IssueReportsServiceServer,
	mealPlanningServiceServer mealplanningsvc.MealPlanningServiceServer,
	userNotificationsServiceServer notificationssvc.UserNotificationsServiceServer,
	oauthServiceServer oauthsvc.OAuthServiceServer,
	settingsServiceServer settingssvc.SettingsServiceServer,
	uploadedMediaServiceServer uploadedmediasvc.UploadedMediaServiceServer,
	webhooksServiceServer webhookssvc.WebhooksServiceServer,
	waitlistsServiceServer waitlistssvc.WaitlistsServiceServer,
	server *grpc.Server,
) *GRPCService {
	return &GRPCService{
		Server:                         server,
		AuditServiceServer:             auditServiceServer,
		AuthServiceServer:              authServiceServer,
		DataPrivacyServiceServer:       dataPrivacyServiceServer,
		IdentityServiceServer:          identityServiceServer,
		InternalOperationsServer:       internalOperationsServer,
		IssueReportsServiceServer:      issueReportsServiceServer,
		MealPlanningServiceServer:      mealPlanningServiceServer,
		UserNotificationsServiceServer: userNotificationsServiceServer,
		OAuthServiceServer:             oauthServiceServer,
		SettingsServiceServer:          settingsServiceServer,
		UploadedMediaServiceServer:     uploadedMediaServiceServer,
		WebhooksServiceServer:          webhooksServiceServer,
		WaitlistsServiceServer:         waitlistsServiceServer,
	}
}

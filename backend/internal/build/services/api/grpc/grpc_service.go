package grpcapi

import (
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
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
	mealplanningsvc.MealPlanningServiceServer
	notificationssvc.UserNotificationsServiceServer
	oauthsvc.OAuthServiceServer
	settingssvc.SettingsServiceServer
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
	mealPlanningServiceServer mealplanningsvc.MealPlanningServiceServer,
	userNotificationsServiceServer notificationssvc.UserNotificationsServiceServer,
	oauthServiceServer oauthsvc.OAuthServiceServer,
	settingsServiceServer settingssvc.SettingsServiceServer,
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
		MealPlanningServiceServer:      mealPlanningServiceServer,
		UserNotificationsServiceServer: userNotificationsServiceServer,
		OAuthServiceServer:             oauthServiceServer,
		SettingsServiceServer:          settingsServiceServer,
		WebhooksServiceServer:          webhooksServiceServer,
		WaitlistsServiceServer:         waitlistsServiceServer,
	}
}

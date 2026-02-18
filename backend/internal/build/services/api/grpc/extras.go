package grpcapi

import (
	"context"
	"maps"

	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/dataprivacy"
	identitysvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/internalops"
	issuereportssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/issue_reports"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/notifications"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/oauth"
	paymentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/payments"
	settingssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/settings"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/uploaded_media"
	waitlistssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/waitlists"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	platformgrpc "github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	authgrpc "github.com/dinnerdonebetter/backend/internal/services/auth/grpc"
	"github.com/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	identitygrpc "github.com/dinnerdonebetter/backend/internal/services/identity/grpc"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	issuereportsgrpc "github.com/dinnerdonebetter/backend/internal/services/issuereports/grpc"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc"
	notificationsgrpc "github.com/dinnerdonebetter/backend/internal/services/notifications/grpc"
	oauthgrpc "github.com/dinnerdonebetter/backend/internal/services/oauth/grpc"
	paymentsgrpc "github.com/dinnerdonebetter/backend/internal/services/payments/grpc"
	settingsgrpc "github.com/dinnerdonebetter/backend/internal/services/settings/grpc"
	uploadedmediagrpc "github.com/dinnerdonebetter/backend/internal/services/uploadedmedia/grpc"
	waitlistsgrpc "github.com/dinnerdonebetter/backend/internal/services/waitlists/grpc"
	webhooksgrpc "github.com/dinnerdonebetter/backend/internal/services/webhooks/grpc"

	grpc "google.golang.org/grpc"
)

func BuildRegistrationFuncs(
	auditLogService auditsvc.AuditServiceServer,
	authService authsvc.AuthServiceServer,
	dataPrivacyServer dataprivacysvc.DataPrivacyServiceServer,
	identityServiceServer identitysvc.IdentityServiceServer,
	internalOpsService internalopssvc.InternalOperationsServer,
	issueReportsService issuereportssvc.IssueReportsServiceServer,
	mealPlanningService mealplanningsvc.MealPlanningServiceServer,
	notificationsService notificationssvc.UserNotificationsServiceServer,
	oauthService oauthsvc.OAuthServiceServer,
	paymentsService paymentssvc.PaymentsServiceServer,
	settingsService settingssvc.SettingsServiceServer,
	uploadedMediaService uploadedmediasvc.UploadedMediaServiceServer,
	waitlistsService waitlistssvc.WaitlistsServiceServer,
	webhooksService webhookssvc.WebhooksServiceServer,
) []platformgrpc.RegistrationFunc {
	return []platformgrpc.RegistrationFunc{
		func(server *grpc.Server) {
			auditsvc.RegisterAuditServiceServer(server, auditLogService)
			authsvc.RegisterAuthServiceServer(server, authService)
			dataprivacysvc.RegisterDataPrivacyServiceServer(server, dataPrivacyServer)
			identitysvc.RegisterIdentityServiceServer(server, identityServiceServer)
			internalopssvc.RegisterInternalOperationsServer(server, internalOpsService)
			issuereportssvc.RegisterIssueReportsServiceServer(server, issueReportsService)
			mealplanningsvc.RegisterMealPlanningServiceServer(server, mealPlanningService)
			notificationssvc.RegisterUserNotificationsServiceServer(server, notificationsService)
			oauthsvc.RegisterOAuthServiceServer(server, oauthService)
			paymentssvc.RegisterPaymentsServiceServer(server, paymentsService)
			settingssvc.RegisterSettingsServiceServer(server, settingsService)
			uploadedmediasvc.RegisterUploadedMediaServiceServer(server, uploadedMediaService)
			waitlistssvc.RegisterWaitlistsServiceServer(server, waitlistsService)
			webhookssvc.RegisterWebhooksServiceServer(server, webhooksService)
		},
	}
}

func BuildUnaryServerInterceptors(authInterceptor *interceptors.AuthInterceptor) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		authInterceptor.UnaryServerInterceptor(),
	}
}

func BuildStreamServerInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		//
	}
}

func ProvideUserTextSearcher(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	metricsProvider metrics.Provider,
	cfg *textsearchcfg.Config,
) (identityindexing.UserTextSearcher, error) {
	return textsearchcfg.ProvideIndex[identityindexing.UserSearchSubset](
		ctx,
		logger,
		tracerProvider, metricsProvider,
		cfg,
		identityindexing.IndexTypeUsers,
	)
}

// AggregateMethodPermissions combines method permissions from all services into a single map.
// Each service provides its permissions via a typed map (e.g., SettingsMethodPermissions),
// which are then aggregated here for the auth interceptor.
func AggregateMethodPermissions(
	authPermissions authgrpc.AuthMethodPermissions,
	identityPermissions identitygrpc.IdentityMethodPermissions,
	issuereportsPermissions issuereportsgrpc.IssueReportsMethodPermissions,
	mealplanningPermissions mealplanninggrpc.MealPlanningMethodPermissions,
	notificationsPermissions notificationsgrpc.NotificationsMethodPermissions,
	oauthPermissions oauthgrpc.OAuthMethodPermissions,
	paymentsPermissions paymentsgrpc.PaymentsMethodPermissions,
	settingsPermissions settingsgrpc.SettingsMethodPermissions,
	uploadedmediaPermissions uploadedmediagrpc.UploadedMediaMethodPermissions,
	waitlistsPermissions waitlistsgrpc.WaitlistsMethodPermissions,
	webhooksPermissions webhooksgrpc.WebhooksMethodPermissions,
) interceptors.MethodPermissionsMap {
	result := make(interceptors.MethodPermissionsMap)

	// Copy all service permissions into the aggregated map
	maps.Copy(result, authPermissions)
	maps.Copy(result, identityPermissions)
	maps.Copy(result, issuereportsPermissions)
	maps.Copy(result, mealplanningPermissions)
	maps.Copy(result, notificationsPermissions)
	maps.Copy(result, oauthPermissions)
	maps.Copy(result, paymentsPermissions)
	maps.Copy(result, settingsPermissions)
	maps.Copy(result, uploadedmediaPermissions)
	maps.Copy(result, waitlistsPermissions)
	maps.Copy(result, webhooksPermissions)

	return result
}

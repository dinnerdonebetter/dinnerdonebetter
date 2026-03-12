package grpcapi

import (
	"context"
	"maps"

	"github.com/dinnerdonebetter/backend/internal/config"
	analyticspb "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/analytics"
	auditsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/audit"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	commentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/comments"
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
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	errorsgrpc "github.com/dinnerdonebetter/backend/internal/platform/errors/grpc"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	platformgrpc "github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	analyticsgrpc "github.com/dinnerdonebetter/backend/internal/services/analytics/grpc"
	auditgrpc "github.com/dinnerdonebetter/backend/internal/services/audit/grpc"
	authgrpc "github.com/dinnerdonebetter/backend/internal/services/auth/grpc"
	"github.com/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	commentsgrpc "github.com/dinnerdonebetter/backend/internal/services/comments/grpc"
	identitygrpc "github.com/dinnerdonebetter/backend/internal/services/identity/grpc"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	internalopsgrpc "github.com/dinnerdonebetter/backend/internal/services/internalops/grpc"
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
	analyticsService analyticspb.AnalyticsServiceServer,
	auditLogService auditsvc.AuditServiceServer,
	authService authsvc.AuthServiceServer,
	commentsService commentssvc.CommentsServiceServer,
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
			analyticspb.RegisterAnalyticsServiceServer(server, analyticsService)
			auditsvc.RegisterAuditServiceServer(server, auditLogService)
			authsvc.RegisterAuthServiceServer(server, authService)
			commentssvc.RegisterCommentsServiceServer(server, commentsService)
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
		errorsgrpc.UnaryErrorEncodingInterceptor(),
	}
}

func BuildStreamServerInterceptors(authInterceptor *interceptors.AuthInterceptor) []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		authInterceptor.StreamServerInterceptor(),
		errorsgrpc.StreamErrorEncodingInterceptor(),
	}
}

// ProvideAnalyticsProxySources extracts proxy sources config for the multisource reporter.
func ProvideAnalyticsProxySources(cfg *config.APIServiceConfig) map[string]*analyticscfg.SourceConfig {
	return cfg.Analytics.ProxySources.ToMap()
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
	analyticsPermissions analyticsgrpc.AnalyticsMethodPermissions,
	auditPermissions auditgrpc.AuditMethodPermissions,
	authPermissions authgrpc.AuthMethodPermissions,
	commentsPermissions commentsgrpc.CommentsMethodPermissions,
	identityPermissions identitygrpc.IdentityMethodPermissions,
	internalopsPermissions internalopsgrpc.InternalOpsMethodPermissions,
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
	maps.Copy(result, analyticsPermissions)
	maps.Copy(result, auditPermissions)
	maps.Copy(result, authPermissions)
	maps.Copy(result, commentsPermissions)
	maps.Copy(result, identityPermissions)
	maps.Copy(result, internalopsPermissions)
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

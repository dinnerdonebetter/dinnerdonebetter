package grpcapi

import (
	"context"

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
	waitlistssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/waitlists"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	platformgrpc "github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	"github.com/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"

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
	settingsService settingssvc.SettingsServiceServer,
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
			settingssvc.RegisterSettingsServiceServer(server, settingsService)
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

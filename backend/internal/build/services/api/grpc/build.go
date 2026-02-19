//go:build wireinject

package grpcapi

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	auditmanager "github.com/dinnerdonebetter/backend/internal/domain/audit/manager"
	authmgr "github.com/dinnerdonebetter/backend/internal/domain/auth/managers"
	commentsmanager "github.com/dinnerdonebetter/backend/internal/domain/comments/manager"
	dataprivacymanager "github.com/dinnerdonebetter/backend/internal/domain/dataprivacy/manager"
	identitymgr "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	issuereportsmanager "github.com/dinnerdonebetter/backend/internal/domain/issuereports/manager"
	grocerylistpreparation "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmgr "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	recipeanalysis "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	notificationsmanager "github.com/dinnerdonebetter/backend/internal/domain/notifications/manager"
	oauthmgr "github.com/dinnerdonebetter/backend/internal/domain/oauth/manager"
	paymentsmanager "github.com/dinnerdonebetter/backend/internal/domain/payments/manager"
	settingsmanager "github.com/dinnerdonebetter/backend/internal/domain/settings/manager"
	uploadedmediamanager "github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia/manager"
	waitlistsmanager "github.com/dinnerdonebetter/backend/internal/domain/waitlists/manager"
	webhooksmanager "github.com/dinnerdonebetter/backend/internal/domain/webhooks/manager"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/qrcodes"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/platform/server/grpc"
	uploadscfg "github.com/dinnerdonebetter/backend/internal/platform/uploads/config"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads/objectstorage"
	auditrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	authrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auth"
	commentsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/comments"
	dataprivacyrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	issuereportsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	paymentsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments"
	uploadedmediarepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	webhooksrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	auditsvc "github.com/dinnerdonebetter/backend/internal/services/audit/grpc"
	authsvc "github.com/dinnerdonebetter/backend/internal/services/auth/grpc"
	"github.com/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	authhttpsvc "github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
	commentssvc "github.com/dinnerdonebetter/backend/internal/services/comments/grpc"
	dataprivacysvc "github.com/dinnerdonebetter/backend/internal/services/dataprivacy/grpc"
	identitysvc "github.com/dinnerdonebetter/backend/internal/services/identity/grpc"
	internalopssvc "github.com/dinnerdonebetter/backend/internal/services/internalops/grpc"
	issuereportssvc "github.com/dinnerdonebetter/backend/internal/services/issuereports/grpc"
	mealplanningsvc "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc"
	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"
	notificationssvc "github.com/dinnerdonebetter/backend/internal/services/notifications/grpc"
	oauthsvc "github.com/dinnerdonebetter/backend/internal/services/oauth/grpc"
	paymentsadapters "github.com/dinnerdonebetter/backend/internal/services/payments/adapters"
	paymentssvc "github.com/dinnerdonebetter/backend/internal/services/payments/grpc"
	settingssvc "github.com/dinnerdonebetter/backend/internal/services/settings/grpc"
	uploadedmediacfg "github.com/dinnerdonebetter/backend/internal/services/uploadedmedia/config"
	uploadedmediasvc "github.com/dinnerdonebetter/backend/internal/services/uploadedmedia/grpc"
	waitlistssvc "github.com/dinnerdonebetter/backend/internal/services/waitlists/grpc"
	webhookssvc "github.com/dinnerdonebetter/backend/internal/services/webhooks/grpc"

	"github.com/google/wire"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (*GRPCService, error) {
	wire.Build(ConfigProviders,
		// core
		metricscfg.MetricsConfigProviders,
		loggingcfg.LogConfigProviders,
		tracingcfg.TracingConfigProviders,
		msgconfig.MessageQueueProviders,
		authentication.AuthProviders,
		sessions.SessionProviders,
		observability.O11yProviders,
		random.RandProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		grpc.ProvidersGRPC,
		qrcodes.QRCodeProviders,
		tokenscfg.TokenIssuerProviders,
		interceptors.InterceptorProviders,
		uploadscfg.Providers,
		objectstorage.Providers,
		// repos
		auditrepo.AuditRepoProviders,
		auditmanager.AuditManagerProviders,
		authrepo.AuthRepoProviders,
		commentsrepo.CommentsRepoProviders,
		identityrepo.IDRepoProviders,
		issuereportsrepo.IssueReportsRepoProviders,
		issuereportsmanager.IssueReportsManagerProviders,
		uploadedmediarepo.UploadedMediaRepoProviders,
		uploadedmediamanager.UploadedMediaManagerProviders,
		webhooksrepo.WebhookProviders,
		oauthrepo.OAuthRepoProviders,
		paymentsrepo.PaymentsRepoProviders,
		mealplanningrepo.MPRepoProviders,
		dataprivacyrepo.DataPrivProviders,
		dataprivacymanager.DataPrivacyManagerProviders,
		// services
		authhttpsvc.AuthHTTPServiceProviders,
		auditsvc.AuditSvcProviders,
		commentssvc.CommentsSvcProviders,
		authsvc.AuthSvcProviders,
		dataprivacysvc.DataPrivSvcProviders,
		identitysvc.IDSvcProviders,
		internalopssvc.InternalOpsSvcProviders,
		issuereportssvc.IssueReportSvcProviders,
		notificationssvc.NotifsSvcProviders,
		settingssvc.SettingSvcProviders,
		uploadedmediasvc.UploadedMediaSvcProviders,
		webhookssvc.WebhookSvcProviders,
		oauthsvc.OAuthSvcProviders,
		paymentssvc.PaymentsSvcProviders,
		paymentsadapters.PaymentsAdapterProviders,
		mealplanningsvc.MPSvcProviders,
		waitlistssvc.WaitlistsSvcProviders,
		uploadedmediacfg.UploadedMediaConfigProviders,
		// manager
		commentsmanager.CommentsManagerProviders,
		identitymgr.IDManagerProviders,
		notificationsmanager.NotificationsManagerProviders,
		settingsmanager.SettingsManagerProviders,
		paymentsmanager.PaymentsManagerProviders,
		oauthmgr.OAuthManagerProviders,
		mealplanningmgr.MPManagerProviders,
		authmgr.AuthManagerProviders,
		webhooksmanager.WebhookManagerProviders,
		waitlistsmanager.WaitlistManagerProviders,
		// workers
		mealplanfinalizer.ProvidersMealPlanFinalizer,
		mealplangrocerylistinitializer.ProvidersMealPlanGroceryListInitializer,
		mealplantaskcreator.ProvidersMealPlanTaskCreator,
		// misc
		recipeanalysis.ProvidersRecipeAnalysis,
		grocerylistpreparation.ProvidersGroceryListPreparation,
		ProvideUserTextSearcher,
		AggregateMethodPermissions,
		BuildUnaryServerInterceptors,
		BuildStreamServerInterceptors,
		BuildRegistrationFuncs,
		NewGRPCService,
	)

	return nil, nil
}

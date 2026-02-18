//go:build wireinject

package grpcapi

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	authmgr "github.com/dinnerdonebetter/backend/internal/domain/auth/managers"
	identitymgr "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	grocerylistpreparation "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmgr "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	recipeanalysis "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	oauthmgr "github.com/dinnerdonebetter/backend/internal/domain/oauth/manager"
	paymentsmanager "github.com/dinnerdonebetter/backend/internal/domain/payments/manager"
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
	dataprivacyrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	issuereportsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	notificationsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	paymentsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments"
	settingsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/settings"
	uploadedmediarepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	waitlistsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/waitlists"
	webhooksrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	auditsvc "github.com/dinnerdonebetter/backend/internal/services/audit/grpc"
	authsvc "github.com/dinnerdonebetter/backend/internal/services/auth/grpc"
	"github.com/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	authhttpsvc "github.com/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
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
		authrepo.AuthRepoProviders,
		identityrepo.IDRepoProviders,
		issuereportsrepo.IssueReportsRepoProviders,
		notificationsrepo.NotifRepoProviders,
		settingsrepo.SettingsRepoProviders,
		uploadedmediarepo.UploadedMediaRepoProviders,
		webhooksrepo.WebhookProviders,
		oauthrepo.OAuthRepoProviders,
		paymentsrepo.PaymentsRepoProviders,
		mealplanningrepo.MPRepoProviders,
		waitlistsrepo.WaitlistsRepoProviders,
		dataprivacyrepo.DataPrivProviders,
		// services
		authhttpsvc.AuthHTTPServiceProviders,
		auditsvc.AuditSvcProviders,
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
		identitymgr.IDManagerProviders,
		paymentsmanager.PaymentsManagerProviders,
		oauthmgr.OAuthManagerProviders,
		mealplanningmgr.MPManagerProviders,
		authmgr.AuthManagerProviders,
		webhooksmanager.WebhookManagerProviders,
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

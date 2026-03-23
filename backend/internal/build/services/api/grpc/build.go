package grpcapi

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/sessions"
	tokenscfg "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/tokens/config"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	auditmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/audit/manager"
	authmgr "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/auth/managers"
	commentsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/comments/manager"
	dataprivacymanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/dataprivacy/manager"
	identitymgr "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity/manager"
	issuereportsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/issuereports/manager"
	grocerylistpreparation "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/grocerylistpreparation"
	mealplanningmgr "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	recipeanalysis "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	notificationsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/notifications/manager"
	oauthmgr "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/oauth/manager"
	paymentsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/payments/manager"
	settingsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/settings/manager"
	uploadedmediamanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/uploadedmedia/manager"
	waitlistsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/waitlists/manager"
	webhooksmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/webhooks/manager"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories"
	auditrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	authrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auth"
	commentsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/comments"
	dataprivacyrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	identityrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	internalopsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	issuereportsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	mealplanningrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	oauthrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	paymentsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/payments"
	uploadedmediarepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	webhooksrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	analyticssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/analytics/grpc"
	auditsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/audit/grpc"
	authsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/auth/grpc"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/auth/grpc/interceptors"
	authhttpsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/auth/handlers/authentication"
	commentssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/comments/grpc"
	dataprivacysvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/dataprivacy/grpc"
	identitysvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/identity/grpc"
	internalopssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/internalops/grpc"
	issuereportssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/issuereports/grpc"
	mealplanningsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/grpc"
	mealplanfinalizer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_finalizer"
	mealplangrocerylistinitializer "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_grocery_list_initializer"
	mealplantaskcreator "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/workers/meal_plan_task_creator"
	notificationssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/notifications/grpc"
	oauthsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/oauth/grpc"
	paymentsadapters "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/payments/adapters"
	paymentssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/payments/grpc"
	settingssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/settings/grpc"
	uploadedmediacfg "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/uploadedmedia/config"
	uploadedmediasvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/uploadedmedia/grpc"
	waitlistssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/waitlists/grpc"
	webhookssvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/webhooks/grpc"

	"github.com/samber/do/v2"
	"github.com/verygoodsoftwarenotvirus/platform/v2/analytics/multisource"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v2/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/database/postgres"
	featureflagscfg "github.com/verygoodsoftwarenotvirus/platform/v2/featureflags/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/httpclient"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v2/messagequeue/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v2/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v2/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v2/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/qrcodes"
	"github.com/verygoodsoftwarenotvirus/platform/v2/random"
	"github.com/verygoodsoftwarenotvirus/platform/v2/server/grpc"
	uploadscfg "github.com/verygoodsoftwarenotvirus/platform/v2/uploads/config"
	"github.com/verygoodsoftwarenotvirus/platform/v2/uploads/objectstorage"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	// config field extraction
	RegisterConfigs(i)

	// platform providers
	observability.RegisterO11yConfigs(i)
	metricscfg.RegisterMetricsProvider(i)
	loggingcfg.RegisterLogger(i)
	tracingcfg.RegisterTracerProvider(i)
	httpclient.RegisterHTTPClient(i)
	msgconfig.RegisterMessageQueue(i)
	random.RegisterGenerator(i)
	databasecfg.RegisterClientConfig(i)
	postgres.RegisterDatabaseClient(i)
	grpc.RegisterGRPCServer(i)
	qrcodes.RegisterBuilder(i)
	uploadscfg.RegisterStorageConfig(i)
	objectstorage.RegisterUploadManager(i)
	featureflagscfg.RegisterFeatureFlagManager(i)
	multisource.RegisterMultiSourceEventReporter(i)

	// authentication
	authentication.RegisterAuth(i)
	sessions.RegisterSessionProviders(i)
	tokenscfg.RegisterTokenIssuer(i)
	interceptors.RegisterAuthInterceptor(i)

	// repositories
	repositories.RegisterMigrator(i)
	auditrepo.RegisterAuditLogRepository(i)
	authrepo.RegisterAuthRepository(i)
	commentsrepo.RegisterCommentsRepository(i)
	identityrepo.RegisterIdentityRepository(i)
	issuereportsrepo.RegisterIssueReportsRepository(i)
	uploadedmediarepo.RegisterUploadedMediaRepository(i)
	webhooksrepo.RegisterWebhooksRepository(i)
	oauthrepo.RegisterOAuthRepository(i)
	paymentsrepo.RegisterPaymentsRepository(i)
	mealplanningrepo.RegisterMealPlanningRepository(i)
	dataprivacyrepo.RegisterDataPrivacyRepository(i)
	internalopsrepo.RegisterInternalOpsRepository(i)

	// managers
	auditmanager.RegisterAuditDataManager(i)
	authmgr.RegisterAuthManager(i)
	commentsmanager.RegisterCommentsDataManager(i)
	identitymgr.RegisterIdentityDataManager(i)
	notificationsmanager.RegisterNotificationsDataManager(i)
	settingsmanager.RegisterSettingsDataManager(i)
	paymentsmanager.RegisterPaymentsDataManager(i)
	oauthmgr.RegisterOAuth2Manager(i)
	mealplanningmgr.RegisterManagers(i)
	webhooksmanager.RegisterWebhookDataManager(i)
	waitlistsmanager.RegisterWaitlistDataManager(i)
	issuereportsmanager.RegisterIssueReportsDataManager(i)
	uploadedmediamanager.RegisterUploadedMediaManager(i)
	dataprivacymanager.RegisterDataPrivacyManager(i)
	paymentsadapters.RegisterPaymentProcessorRegistry(i)

	// services
	authsvc.RegisterAuthService(i)
	authhttpsvc.RegisterAuthHTTPService(i)
	analyticssvc.RegisterAnalyticsService(i)
	auditsvc.RegisterAuditService(i)
	commentssvc.RegisterCommentsService(i)
	dataprivacysvc.RegisterDataPrivacyService(i)
	identitysvc.RegisterIdentityService(i)
	internalopssvc.RegisterInternalOpsService(i)
	issuereportssvc.RegisterIssueReportsService(i)
	notificationssvc.RegisterNotificationsService(i)
	settingssvc.RegisterSettingsService(i)
	uploadedmediasvc.RegisterUploadedMediaService(i)
	webhookssvc.RegisterWebhooksService(i)
	oauthsvc.RegisterOAuthService(i)
	paymentssvc.RegisterPaymentsService(i)
	mealplanningsvc.RegisterMealPlanningService(i)
	waitlistssvc.RegisterWaitlistsService(i)
	uploadedmediacfg.RegisterUploadedMediaConfig(i)

	// workers
	mealplanfinalizer.RegisterMealPlanFinalizer(i)
	mealplangrocerylistinitializer.RegisterMealPlanGroceryListInitializer(i)
	mealplantaskcreator.RegisterMealPlanTaskCreator(i)

	// misc
	recipeanalysis.RegisterRecipeAnalyzer(i)
	grocerylistpreparation.RegisterGroceryListCreator(i)

	// extras (functions from extras.go)
	RegisterExtras(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.APIServiceConfig,
) (*GRPCService, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*GRPCService](i), nil
}

package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	notificationsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/notifications/manager"
	settingsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/settings/manager"
	waitlistsmanager "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/waitlists/manager"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/functions/datachangemessagehandler"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/auth"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	internalopsrepo "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	issue_reports "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	identityindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/samber/do/v2"
	analyticscfg "github.com/verygoodsoftwarenotvirus/platform/v4/analytics/config"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/v4/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/database/postgres"
	emailcfg "github.com/verygoodsoftwarenotvirus/platform/v4/email/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/v4/httpclient"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/v4/messagequeue/config"
	notificationscfg "github.com/verygoodsoftwarenotvirus/platform/v4/notifications/mobile/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/v4/uploads/objectstorage"
)

// BuildInjector creates and configures the dependency injection container.
func BuildInjector(
	ctx context.Context,
	cfg *config.AsyncMessageHandlerConfig,
) *do.RootScope {
	i := do.New()

	do.ProvideValue(i, ctx)
	do.ProvideValue(i, cfg)

	// config field extraction
	RegisterConfigs(i)

	// platform providers
	observability.RegisterO11yConfigs(i)
	tracingcfg.RegisterTracerProvider(i)
	loggingcfg.RegisterLogger(i)
	metricscfg.RegisterMetricsProvider(i)
	msgconfig.RegisterMessageQueue(i)
	httpclient.RegisterHTTPClient(i)
	encoding.RegisterServerEncoderDecoder(i)
	analyticscfg.RegisterEventReporter(i)
	emailcfg.RegisterEmailer(i)
	databasecfg.RegisterClientConfig(i)
	postgres.RegisterDatabaseClient(i)
	objectstorage.RegisterUploadManager(i)
	notificationscfg.RegisterPushSender(i)

	// repos
	auditlogentries.RegisterAuditLogRepository(i)
	auth.RegisterAuthRepository(i)
	dataprivacy.RegisterDataPrivacyRepository(i)
	identity.RegisterIdentityRepository(i)
	issue_reports.RegisterIssueReportsRepository(i)
	mealplanning.RegisterMealPlanningRepository(i)
	uploadedmedia.RegisterUploadedMediaRepository(i)
	webhooks.RegisterWebhooksRepository(i)
	internalopsrepo.RegisterInternalOpsRepository(i)

	// managers
	notificationsmanager.RegisterNotificationsDataManager(i)
	settingsmanager.RegisterSettingsDataManager(i)
	waitlistsmanager.RegisterWaitlistDataManager(i)

	// indexing
	identityindexing.RegisterCoreDataIndexer(i)
	eatingindexing.RegisterMealPlanningDataIndexer(i)

	// searchers
	RegisterSearchers(i)

	// main handler
	datachangemessagehandler.RegisterAsyncDataChangeMessageHandler(i)

	return i
}

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.AsyncMessageHandlerConfig,
) (*datachangemessagehandler.AsyncDataChangeMessageHandler, error) {
	i := BuildInjector(ctx, cfg)
	return do.MustInvoke[*datachangemessagehandler.AsyncDataChangeMessageHandler](i), nil
}

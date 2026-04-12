package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	mealplanningregistration "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/mealplanning/registration"
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
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	identityindexing "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/services/identity/indexing"

	analyticscfg "github.com/primandproper/platform/analytics/config"
	databasecfg "github.com/primandproper/platform/database/config"
	"github.com/primandproper/platform/database/postgres"
	emailcfg "github.com/primandproper/platform/email/config"
	"github.com/primandproper/platform/encoding"
	"github.com/primandproper/platform/httpclient"
	msgconfig "github.com/primandproper/platform/messagequeue/config"
	notificationscfg "github.com/primandproper/platform/notifications/mobile/config"
	"github.com/primandproper/platform/observability"
	loggingcfg "github.com/primandproper/platform/observability/logging/config"
	metricscfg "github.com/primandproper/platform/observability/metrics/config"
	tracingcfg "github.com/primandproper/platform/observability/tracing/config"
	"github.com/primandproper/platform/uploads/objectstorage"

	"github.com/samber/do/v2"
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

	// Domain: mealplanning
	mealplanningregistration.RegisterForDataChangeHandler(i)

	// repos
	auditlogentries.RegisterAuditLogRepository(i)
	auth.RegisterAuthRepository(i)
	dataprivacy.RegisterDataPrivacyRepository(i)
	identity.RegisterIdentityRepository(i)
	issue_reports.RegisterIssueReportsRepository(i)
	uploadedmedia.RegisterUploadedMediaRepository(i)
	webhooks.RegisterWebhooksRepository(i)
	internalopsrepo.RegisterInternalOpsRepository(i)

	// managers
	notificationsmanager.RegisterNotificationsDataManager(i)
	settingsmanager.RegisterSettingsDataManager(i)
	waitlistsmanager.RegisterWaitlistDataManager(i)

	// indexing
	identityindexing.RegisterCoreDataIndexer(i)

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

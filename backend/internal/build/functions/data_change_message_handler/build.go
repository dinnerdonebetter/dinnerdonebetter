//go:build wireinject

package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	notificationsmanager "github.com/dinnerdonebetter/backend/internal/domain/notifications/manager"
	settingsmanager "github.com/dinnerdonebetter/backend/internal/domain/settings/manager"
	waitlistsmanager "github.com/dinnerdonebetter/backend/internal/domain/waitlists/manager"
	"github.com/dinnerdonebetter/backend/internal/functions/datachangemessagehandler"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auth"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	internalopsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/internalops"
	issue_reports "github.com/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/google/wire"
	analyticscfg "github.com/verygoodsoftwarenotvirus/platform/analytics/config"
	databasecfg "github.com/verygoodsoftwarenotvirus/platform/database/config"
	"github.com/verygoodsoftwarenotvirus/platform/database/postgres"
	emailcfg "github.com/verygoodsoftwarenotvirus/platform/email/config"
	"github.com/verygoodsoftwarenotvirus/platform/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/httpclient"
	msgconfig "github.com/verygoodsoftwarenotvirus/platform/messagequeue/config"
	notificationscfg "github.com/verygoodsoftwarenotvirus/platform/notifications/config"
	"github.com/verygoodsoftwarenotvirus/platform/observability"
	loggingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/logging/config"
	metricscfg "github.com/verygoodsoftwarenotvirus/platform/observability/metrics/config"
	tracingcfg "github.com/verygoodsoftwarenotvirus/platform/observability/tracing/config"
	"github.com/verygoodsoftwarenotvirus/platform/uploads/objectstorage"
)

// Build builds a server.
func Build(
	ctx context.Context,
	cfg *config.AsyncMessageHandlerConfig,
) (*datachangemessagehandler.AsyncDataChangeMessageHandler, error) {
	wire.Build(
		datachangemessagehandler.Providers,
		msgconfig.MessageQueueProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		auditlogentries.AuditRepoProviders,
		auth.AuthRepoProviders,
		dataprivacy.DataPrivProviders,
		identity.IDRepoProviders,
		issue_reports.IssueReportsRepoProviders,
		mealplanning.MPRepoProviders,
		notificationsmanager.NotificationsManagerProviders,
		settingsmanager.SettingsManagerProviders,
		uploadedmedia.UploadedMediaRepoProviders,
		waitlistsmanager.WaitlistManagerProviders,
		webhooks.WebhookProviders,
		internalopsrepo.Providers,
		analyticscfg.Providers,
		emailcfg.Providers,
		metricscfg.MetricsConfigProviders,
		encoding.Providers,
		loggingcfg.LogConfigProviders,
		httpclient.Providers,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		objectstorage.Providers,
		identityindexing.Providers,
		eatingindexing.Providers,
		ConfigProviders,
		SearcherProviders,
		notificationscfg.Providers,
	)

	return nil, nil
}

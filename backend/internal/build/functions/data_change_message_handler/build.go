//go:build wireinject
// +build wireinject

package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/functions/datachangemessagehandler"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	emailcfg "github.com/dinnerdonebetter/backend/internal/platform/email/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	issue_reports "github.com/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/settings"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/waitlists"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"
	identityindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"github.com/google/wire"
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
		dataprivacy.DataPrivProviders,
		identity.IDRepoProviders,
		issue_reports.IssueReportsRepoProviders,
		mealplanning.MPRepoProviders,
		notifications.NotifRepoProviders,
		settings.SettingsRepoProviders,
		uploadedmedia.UploadedMediaRepoProviders,
		waitlists.WaitlistsRepoProviders,
		webhooks.WebhookProviders,
		analyticscfg.Providers,
		emailcfg.Providers,
		metricscfg.MetricsConfigProviders,
		encoding.Providers,
		loggingcfg.LogConfigProviders,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		tracing.ProvidersTracing,
		objectstorage.Providers,
		identityindexing.Providers,
		eatingindexing.Providers,
		ConfigProviders,
		SearcherProviders,
	)

	return nil, nil
}

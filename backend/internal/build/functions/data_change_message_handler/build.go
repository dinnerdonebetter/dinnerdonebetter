//go:build wireinject

package datachangemessagehandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	notificationsmanager "github.com/dinnerdonebetter/backend/internal/domain/notifications/manager"
	settingsmanager "github.com/dinnerdonebetter/backend/internal/domain/settings/manager"
	waitlistsmanager "github.com/dinnerdonebetter/backend/internal/domain/waitlists/manager"
	"github.com/dinnerdonebetter/backend/internal/functions/datachangemessagehandler"
	analyticscfg "github.com/dinnerdonebetter/backend/internal/platform/analytics/config"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	emailcfg "github.com/dinnerdonebetter/backend/internal/platform/email/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/httpclient"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	notificationscfg "github.com/dinnerdonebetter/backend/internal/platform/notifications/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads/objectstorage"
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

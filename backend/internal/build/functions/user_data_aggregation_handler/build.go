//go:build wireinject

package userdataaggregationhandler

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/config"
	notificationsmanager "github.com/dinnerdonebetter/backend/internal/domain/notifications/manager"
	settingsmanager "github.com/dinnerdonebetter/backend/internal/domain/settings/manager"
	waitlistsmanager "github.com/dinnerdonebetter/backend/internal/domain/waitlists/manager"
	"github.com/dinnerdonebetter/backend/internal/functions/userdataaggregationhandler"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/logging/config"
	metricscfg "github.com/dinnerdonebetter/backend/internal/platform/observability/metrics/config"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/platform/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/platform/uploads/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	issue_reports "github.com/dinnerdonebetter/backend/internal/repositories/postgres/issuereports"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/uploadedmedia"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/webhooks"

	"github.com/google/wire"
)

// Build builds the user data aggregation handler.
func Build(
	ctx context.Context,
	cfg *config.UserDataAggregationHandlerConfig,
) (*userdataaggregationhandler.UserDataAggregationHandler, error) {
	wire.Build(
		userdataaggregationhandler.Providers,
		msgconfig.MessageQueueProviders,
		databasecfg.ClientConfigProviders,
		postgres.PGProviders,
		auditlogentries.AuditRepoProviders,
		dataprivacy.DataPrivProviders,
		identity.IDRepoProviders,
		issue_reports.IssueReportsRepoProviders,
		mealplanning.MPRepoProviders,
		notificationsmanager.NotificationsManagerProviders,
		settingsmanager.SettingsManagerProviders,
		uploadedmedia.UploadedMediaRepoProviders,
		waitlistsmanager.WaitlistManagerProviders,
		webhooks.WebhookProviders,
		encoding.Providers,
		metricscfg.MetricsConfigProviders,
		loggingcfg.LogConfigProviders,
		tracingcfg.TracingConfigProviders,
		observability.O11yProviders,
		objectstorage.Providers,
		ConfigProviders,
	)

	return nil, nil
}

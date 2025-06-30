package webhooks

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/webhooks/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "webhook_db_client"
)

// Querier is the webhook repository client.
type Querier struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	auditLogEntryRepo auditlogentries.Repository
	db                *sql.DB
}

// ProvideWebhooksRepository provides a new client.
func ProvideWebhooksRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo auditlogentries.Repository,
	client database.Client,
) webhooks.Repository {
	c := &Querier{
		Client:            client,
		db:                client.DB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}

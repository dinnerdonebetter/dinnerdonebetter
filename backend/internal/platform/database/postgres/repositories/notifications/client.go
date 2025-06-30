package notifications

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/notifications/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
)

const (
	o11yName = "notifications_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	database.Client
	tracer            tracing.Tracer
	logger            logging.Logger
	generatedQuerier  generated.Querier
	auditLogEntryRepo auditlogentries.Repository
	secretGenerator   random.Generator
	db                *sql.DB
}

// ProvideSettingsRepository provides a new client.
func ProvideSettingsRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	auditLogEntryRepo auditlogentries.Repository,
	client database.Client,
) notifications.Repository {
	c := &Querier{
		Client:            client,
		db:                client.DB(),
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier:  generated.New(),
		auditLogEntryRepo: auditLogEntryRepo,
		secretGenerator:   random.NewGenerator(logger, tracerProvider),
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}

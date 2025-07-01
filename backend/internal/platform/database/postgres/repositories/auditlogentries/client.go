package auditlogentries

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/repositories/auditlogentries/generated"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
)

const (
	o11yName = "audit_log_entries_db_client"
)

// Querier is the audit log entry client.
type Querier struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated.Querier
	secretGenerator  random.Generator
	db               *sql.DB
}

// ProvideAuditLogRepository provides a new client.
func ProvideAuditLogRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	client database.Client,
) audit.Repository {
	c := &Querier{
		Client:           client,
		db:               client.DB(),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier: generated.New(),
		secretGenerator:  random.NewGenerator(logger, tracerProvider),
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}

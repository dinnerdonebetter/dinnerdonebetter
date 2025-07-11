package auditlogentries

import (
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	generated2 "github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries/generated"
)

const (
	o11yName = "audit_log_entries_db_client"
)

// repository is the audit log entry repository implementation.
type repository struct {
	database.Client
	tracer           tracing.Tracer
	logger           logging.Logger
	generatedQuerier generated2.Querier
	secretGenerator  random.Generator
	db               *sql.DB
}

// ProvideAuditLogRepository provides a new repository.
func ProvideAuditLogRepository(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	client database.Client,
) audit.Repository {
	c := &repository{
		Client:           client,
		db:               client.DB(),
		tracer:           tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		generatedQuerier: generated2.New(),
		secretGenerator:  random.NewGenerator(logger, tracerProvider),
		logger:           logging.EnsureLogger(logger).WithName(o11yName),
	}

	return c
}
